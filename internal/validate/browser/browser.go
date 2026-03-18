package browser

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/log"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

const (
	defaultLoadTimeout      = 20 * time.Second
	defaultLiveFetchTimeout = 30 * time.Second
	defaultRequestMatch     = "/openapi/execute_query"
)

type Options struct {
	HTMLPath                 string
	BaseURL                  string
	VisualMode               string
	Token                    string
	SkipLiveFetch            bool
	LoadTimeout              time.Duration
	LiveFetchTimeout         time.Duration
	ExpectedRequestSubstring string
}

type Result struct {
	Valid               bool
	Errors              []string
	Warnings            []string
	ConsoleErrors       []string
	MissingControls     []string
	LiveFetchAttempted  bool
	LiveFetchSucceeded  bool
	LiveFetchSkipped    bool
	SkipReason          string
	MatchedRequestURL   string
	MatchedResponseCode int64
	StatusText          string
}

type controlsState struct {
	HasFooter      bool     `json:"hasFooter"`
	HasLedger      bool     `json:"hasLedger"`
	HasTokenInput  bool     `json:"hasTokenInput"`
	HasTextarea    bool     `json:"hasTextarea"`
	HasFetchButton bool     `json:"hasFetchButton"`
	HasStatus      bool     `json:"hasStatus"`
	Missing        []string `json:"missing"`
}

type pageState struct {
	StatusText        string   `json:"statusText"`
	FatalMessages     []string `json:"fatalMessages"`
	LoadingVisible    int      `json:"loadingVisible"`
	LedgerStatuses    []string `json:"ledgerStatuses"`
	AnalyticalSignals int      `json:"analyticalSignals"`
	TokenStored       bool     `json:"tokenStored"`
}

type tracker struct {
	mu                sync.Mutex
	match             string
	requests          map[network.RequestID]string
	matchedRequestURL string
	matchedStatusCode int64
	loadingFailures   []string
	consoleErrors     []string
	exceptionErrors   []string
}

func Validate(ctx context.Context, opts Options) Result {
	result := Result{Valid: true}
	loadTimeout := opts.LoadTimeout
	if loadTimeout <= 0 {
		loadTimeout = defaultLoadTimeout
	}
	liveFetchTimeout := opts.LiveFetchTimeout
	if liveFetchTimeout <= 0 {
		liveFetchTimeout = defaultLiveFetchTimeout
	}
	match := opts.ExpectedRequestSubstring
	if match == "" {
		match = defaultRequestMatch
	}

	targetURL, cleanup, err := resolveTargetURL(opts)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, err.Error())
		return result
	}
	if cleanup != nil {
		defer cleanup()
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
	)...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	tr := &tracker{
		match:    match,
		requests: map[network.RequestID]string{},
	}
	tr.attach(browserCtx)

	loadCtx, cancelLoad := context.WithTimeout(browserCtx, loadTimeout)
	defer cancelLoad()
	if err := chromedp.Run(loadCtx,
		network.Enable(),
		runtime.Enable(),
		log.Enable(),
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, classifyError("page_load", err))
		return result
	}

	controls, err := discoverControls(browserCtx)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, classifyError("control_discovery", err))
		return result
	}
	result.MissingControls = append(result.MissingControls, controls.Missing...)
	if len(controls.Missing) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "browser validation missing required controls: "+strings.Join(controls.Missing, ", "))
	}

	initialState, err := collectPageState(browserCtx)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, classifyError("initial_state", err))
		return result
	}
	result.StatusText = initialState.StatusText

	consoleErrors := tr.consoleAndExceptionErrors()
	result.ConsoleErrors = append(result.ConsoleErrors, consoleErrors...)
	if len(consoleErrors) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "browser console/runtime errors: "+strings.Join(consoleErrors, " | "))
	}
	if len(initialState.FatalMessages) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "page reported fatal UI errors: "+strings.Join(initialState.FatalMessages, " | "))
	}
	if !result.Valid {
		return result
	}

	if !ShouldAttemptLiveFetch(opts.VisualMode, opts.Token, opts.SkipLiveFetch) {
		result.LiveFetchSkipped = true
		result.SkipReason = liveFetchSkipReason(opts.VisualMode, opts.Token, opts.SkipLiveFetch)
		return result
	}

	result.LiveFetchAttempted = true
	if err := setTokenAndClick(browserCtx, opts.Token); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, classifyError("click_fetch", err))
		return result
	}

	if err := waitForMatchedResponse(browserCtx, tr, liveFetchTimeout); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, classifyError("live_fetch", err))
		result.ConsoleErrors = append(result.ConsoleErrors, tr.consoleAndExceptionErrors()...)
		return result
	}

	result.MatchedRequestURL, result.MatchedResponseCode = tr.matchedResponse()
	if result.MatchedRequestURL == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "live fetch did not issue a matching browser request")
		return result
	}
	if result.MatchedResponseCode < 200 || result.MatchedResponseCode >= 300 {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("live fetch returned HTTP %d for %s", result.MatchedResponseCode, result.MatchedRequestURL))
		return result
	}

	settleState, err := waitForSettledState(browserCtx, liveFetchTimeout)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, classifyError("post_fetch_render", err))
		return result
	}
	result.StatusText = settleState.StatusText
	if len(settleState.FatalMessages) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "page reported fatal UI errors after live fetch: "+strings.Join(settleState.FatalMessages, " | "))
	}
	if !hasSuccessSignal(settleState) {
		result.Valid = false
		result.Errors = append(result.Errors, "browser validation did not observe a successful rendered state after live fetch")
	}
	postConsoleErrors := tr.consoleAndExceptionErrors()
	result.ConsoleErrors = append(result.ConsoleErrors[:0], postConsoleErrors...)
	if len(postConsoleErrors) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "browser console/runtime errors after live fetch: "+strings.Join(postConsoleErrors, " | "))
	}
	if !settleState.TokenStored {
		result.Warnings = append(result.Warnings, "page did not persist the shared JWE token to localStorage")
	}
	if result.Valid {
		result.LiveFetchSucceeded = true
	}
	return result
}

func ShouldAttemptLiveFetch(visualMode, token string, skip bool) bool {
	if skip {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(visualMode), "dynamic") {
		return false
	}
	return strings.TrimSpace(token) != ""
}

func liveFetchSkipReason(visualMode, token string, skip bool) string {
	switch {
	case skip:
		return "skipped_by_flag"
	case !strings.EqualFold(strings.TrimSpace(visualMode), "dynamic"):
		return "non_dynamic_mode"
	case strings.TrimSpace(token) == "":
		return "missing_token"
	default:
		return "not_requested"
	}
}

func classifyError(stage string, err error) string {
	if err == nil {
		return ""
	}
	msg := strings.TrimSpace(err.Error())
	switch {
	case strings.Contains(strings.ToLower(msg), "context deadline exceeded"):
		return fmt.Sprintf("%s timeout: %s", stage, msg)
	default:
		return fmt.Sprintf("%s failed: %s", stage, msg)
	}
}

func resolveTargetURL(opts Options) (string, func(), error) {
	if opts.BaseURL != "" {
		return opts.BaseURL, nil, nil
	}
	if opts.HTMLPath == "" {
		return "", nil, fmt.Errorf("browser validation requires HTMLPath or BaseURL")
	}
	absPath, err := filepath.Abs(opts.HTMLPath)
	if err != nil {
		return "", nil, fmt.Errorf("resolve html path: %w", err)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return "", nil, fmt.Errorf("stat html path: %w", err)
	}
	if info.IsDir() {
		return "", nil, fmt.Errorf("html path must be a file: %s", absPath)
	}
	dir := filepath.Dir(absPath)
	base := filepath.Base(absPath)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, fmt.Errorf("start local file server: %w", err)
	}
	server := &http.Server{Handler: http.FileServer(http.Dir(dir))}
	go func() {
		_ = server.Serve(ln)
	}()
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}
	return fmt.Sprintf("http://%s/%s", ln.Addr().String(), url.PathEscape(base)), cleanup, nil
}

func (t *tracker) attach(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev any) {
		t.mu.Lock()
		defer t.mu.Unlock()
		switch e := ev.(type) {
		case *runtime.EventExceptionThrown:
			if e.ExceptionDetails.Text != "" {
				t.exceptionErrors = append(t.exceptionErrors, strings.TrimSpace(e.ExceptionDetails.Text))
			} else if e.ExceptionDetails.Exception != nil {
				t.exceptionErrors = append(t.exceptionErrors, strings.TrimSpace(e.ExceptionDetails.Exception.Description))
			} else {
				t.exceptionErrors = append(t.exceptionErrors, "uncaught runtime exception")
			}
		case *runtime.EventConsoleAPICalled:
			if e.Type != runtime.APITypeError {
				return
			}
			var parts []string
			for _, arg := range e.Args {
				switch {
				case arg.Value != nil:
					parts = append(parts, fmt.Sprint(arg.Value))
				case arg.Description != "":
					parts = append(parts, arg.Description)
				}
			}
			if len(parts) == 0 {
				parts = append(parts, "console.error")
			}
			t.consoleErrors = append(t.consoleErrors, strings.TrimSpace(strings.Join(parts, " ")))
		case *network.EventRequestWillBeSent:
			if strings.Contains(e.Request.URL, t.match) {
				t.requests[e.RequestID] = e.Request.URL
				t.matchedRequestURL = e.Request.URL
			}
		case *network.EventResponseReceived:
			if url, ok := t.requests[e.RequestID]; ok {
				t.matchedRequestURL = url
				t.matchedStatusCode = int64(e.Response.Status)
			}
		case *network.EventLoadingFailed:
			if url, ok := t.requests[e.RequestID]; ok {
				t.loadingFailures = append(t.loadingFailures, fmt.Sprintf("%s: %s", url, e.ErrorText))
			}
		}
	})
}

func (t *tracker) consoleAndExceptionErrors() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	combined := append([]string{}, t.consoleErrors...)
	combined = append(combined, t.exceptionErrors...)
	return combined
}

func (t *tracker) matchedResponse() (string, int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.matchedRequestURL, t.matchedStatusCode
}

func (t *tracker) snapshot() (string, int64, []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.matchedRequestURL, t.matchedStatusCode, append([]string{}, t.loadingFailures...)
}

func waitForMatchedResponse(ctx context.Context, t *tracker, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			return err
		}
		url, status, loadingFailures := t.snapshot()
		if len(loadingFailures) > 0 {
			return fmt.Errorf("matching network request failed: %s", strings.Join(loadingFailures, " | "))
		}
		if url != "" && status > 0 {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return context.DeadlineExceeded
}

func waitForSettledState(ctx context.Context, timeout time.Duration) (pageState, error) {
	deadline := time.Now().Add(timeout)
	var last pageState
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			return pageState{}, err
		}
		state, err := collectPageState(ctx)
		if err != nil {
			return pageState{}, err
		}
		last = state
		if len(state.FatalMessages) == 0 && state.LoadingVisible == 0 && hasSuccessSignal(state) {
			return state, nil
		}
		time.Sleep(150 * time.Millisecond)
	}
	return last, context.DeadlineExceeded
}

func hasSuccessSignal(state pageState) bool {
	for _, status := range state.LedgerStatuses {
		lower := strings.ToLower(status)
		if strings.Contains(lower, "ok") || strings.Contains(lower, "success") || strings.Contains(lower, "loaded") || strings.Contains(lower, "complete") || strings.Contains(lower, "ready") || strings.Contains(lower, "done") {
			return true
		}
	}
	return state.AnalyticalSignals > 0
}

func discoverControls(ctx context.Context) (controlsState, error) {
	var state controlsState
	if err := chromedp.Run(ctx, chromedp.Evaluate(discoverControlsJS, &state)); err != nil {
		return controlsState{}, err
	}
	return state, nil
}

func collectPageState(ctx context.Context) (pageState, error) {
	var state pageState
	if err := chromedp.Run(ctx, chromedp.Evaluate(collectPageStateJS, &state)); err != nil {
		return pageState{}, err
	}
	return state, nil
}

func setTokenAndClick(ctx context.Context, token string) error {
	return chromedp.Run(ctx, chromedp.Evaluate(setTokenAndClickJS(token), nil))
}

func discoverControlsJSFunc() string {
	return `(() => {
		const footer = document.querySelector('footer') || document.querySelector('.footer-controls') || document.querySelector('[data-role="controls"]');
		const scope = footer || document;
		const input = scope.querySelector('input[type="password"]') || document.querySelector('input[type="password"]');
		const textarea = scope.querySelector('textarea') || document.querySelector('textarea');
		const buttons = Array.from(scope.querySelectorAll('button,input[type="button"],input[type="submit"]'));
		const fetchButton = buttons.find((btn) => /fetch|run|execute|load|query|refresh|save/i.test([btn.innerText, btn.value, btn.id, btn.className].join(' ')));
		const status = scope.querySelector('#statusText,[id*="status"],[class*="status"],[data-role*="status"]');
		const ledger = document.querySelector('#query-ledger,.ledger,[id*="ledger"],[class*="ledger"]');
		const missing = [];
		if (!footer) missing.push('footer control block');
		if (!ledger) missing.push('query ledger');
		if (!input) missing.push('password token input');
		if (!textarea) missing.push('SQL textarea');
		if (!fetchButton) missing.push('fetch action button');
		if (!status) missing.push('status text');
		return {
			hasFooter: Boolean(footer),
			hasLedger: Boolean(ledger),
			hasTokenInput: Boolean(input),
			hasTextarea: Boolean(textarea),
			hasFetchButton: Boolean(fetchButton),
			hasStatus: Boolean(status),
			missing
		};
	})()`
}

var discoverControlsJS = discoverControlsJSFunc()

func collectPageStateJSFunc() string {
	return `(() => {
		const footer = document.querySelector('footer') || document.querySelector('.footer-controls') || document.querySelector('[data-role="controls"]');
		const status = (footer || document).querySelector('#statusText,[id*="status"],[class*="status"],[data-role*="status"]');
		const isVisible = (el) => {
			if (!el) return false;
			const style = window.getComputedStyle(el);
			return style.display !== 'none' && style.visibility !== 'hidden' && el.getClientRects().length > 0;
		};
		const fatalMessages = [];
		const fatalNodes = document.querySelectorAll('[role="alert"], .error, .fatal, .alert-error, .status.error, [data-state="error"]');
		fatalNodes.forEach((el) => {
			if (isVisible(el)) {
				const text = (el.innerText || el.textContent || '').trim();
				if (text) fatalMessages.push(text);
			}
		});
		if (status) {
			const text = (status.innerText || status.textContent || '').trim();
			if (/(error|failed|exception|unauthorized|forbidden|denied|invalid)/i.test(text)) {
				fatalMessages.push(text);
			}
		}
		const loadingVisible = Array.from(document.querySelectorAll('[aria-busy="true"], .loading, .spinner, [data-state="loading"]'))
			.filter(isVisible).length;
		const ledgerStatuses = Array.from(document.querySelectorAll('.ledger-status,[class*="ledger-status"],[data-role="ledger-status"]'))
			.filter(isVisible)
			.map((el) => (el.innerText || el.textContent || '').trim())
			.filter(Boolean);
		const analyticalSignals = (() => {
			let count = 0;
			count += Array.from(document.querySelectorAll('tbody tr')).filter(isVisible).length;
			count += Array.from(document.querySelectorAll('svg')).filter((el) => isVisible(el) && el.querySelectorAll('*').length > 5).length;
			count += Array.from(document.querySelectorAll('canvas')).filter(isVisible).length;
			count += Array.from(document.querySelectorAll('[id*="map"],[class*="map"]')).filter((el) => {
				if (!isVisible(el)) return false;
				return el.querySelector('.leaflet-pane,[class*="leaflet"]') || el.childElementCount > 0;
			}).length;
			count += Array.from(document.querySelectorAll('.kpi,[class*="kpi"],.card,.panel')).filter((el) => {
				if (!isVisible(el)) return false;
				if (footer && footer.contains(el)) return false;
				const ident = [el.id, el.className].join(' ');
				if (/(ledger|hero|header|footer|hint|status)/i.test(ident)) return false;
				const text = (el.innerText || el.textContent || '').trim();
				return text.length > 60;
			}).length;
			return count;
		})();
		const stored = (() => {
			try {
				return Boolean(localStorage.getItem('OnTimeAnalystDashboard::auth::jwe'));
			} catch (err) {
				return false;
			}
		})();
		return {
			statusText: status ? (status.innerText || status.textContent || '').trim() : '',
			fatalMessages,
			loadingVisible,
			ledgerStatuses,
			analyticalSignals,
			tokenStored: stored
		};
	})()`
}

var collectPageStateJS = collectPageStateJSFunc()

func setTokenAndClickJS(token string) string {
	return fmt.Sprintf(`(() => {
		const footer = document.querySelector('footer') || document.querySelector('.footer-controls') || document.querySelector('[data-role="controls"]');
		const scope = footer || document;
		const input = scope.querySelector('input[type="password"]') || document.querySelector('input[type="password"]');
		const textarea = scope.querySelector('textarea') || document.querySelector('textarea');
		const buttons = Array.from(scope.querySelectorAll('button,input[type="button"],input[type="submit"]'));
		const fetchButton = buttons.find((btn) => /fetch|run|execute|load|query|refresh|save/i.test([btn.innerText, btn.value, btn.id, btn.className].join(' ')));
		if (!input) throw new Error('missing password token input');
		if (!textarea) throw new Error('missing SQL textarea');
		if (!fetchButton) throw new Error('missing fetch action button');
		input.focus();
		input.value = %q;
		input.dispatchEvent(new Event('input', {bubbles: true}));
		input.dispatchEvent(new Event('change', {bubbles: true}));
		fetchButton.click();
		return true;
	})()`, token)
}
