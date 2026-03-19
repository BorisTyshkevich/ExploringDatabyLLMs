package browser

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestShouldAttemptLiveFetch(t *testing.T) {
	if ShouldAttemptLiveFetch("dynamic", "token", true) {
		t.Fatal("expected explicit skip to disable live fetch")
	}
	if ShouldAttemptLiveFetch("static", "token", false) {
		t.Fatal("expected non-dynamic mode to skip live fetch")
	}
	if ShouldAttemptLiveFetch("dynamic", "", false) {
		t.Fatal("expected missing token to skip live fetch")
	}
	if !ShouldAttemptLiveFetch("dynamic", "token", false) {
		t.Fatal("expected dynamic mode with token to attempt live fetch")
	}
}

func TestClassifyError(t *testing.T) {
	got := classifyError("live_fetch", context.DeadlineExceeded)
	if !strings.Contains(got, "timeout") {
		t.Fatalf("expected timeout classification, got %q", got)
	}
}

func TestValidateLocalFixtureSuccess(t *testing.T) {
	srv := newFixtureServer(http.StatusOK, `{"columns":["value"],"rows":[[1]],"count":1}`, false, false)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	result := Validate(ctx, Options{
		BaseURL:                  srv.URL + "/visual.html",
		VisualMode:               "dynamic",
		Token:                    "public-token",
		ExpectedRequestSubstring: "/api/test",
	})
	if !result.Valid {
		t.Fatalf("expected valid result, got errors: %v warnings: %v", result.Errors, result.Warnings)
	}
	if !result.LiveFetchSucceeded {
		t.Fatalf("expected live fetch success, got %+v", result)
	}
}

func TestValidateFailsOnPageException(t *testing.T) {
	srv := newFixtureServer(http.StatusOK, `{"columns":["value"],"rows":[[1]],"count":1}`, true, false)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result := Validate(ctx, Options{
		BaseURL:                  srv.URL + "/visual.html",
		VisualMode:               "dynamic",
		Token:                    "public-token",
		ExpectedRequestSubstring: "/api/test",
	})
	if result.Valid {
		t.Fatalf("expected invalid result for runtime exception, got %+v", result)
	}
}

func TestValidateFailsWhenControlsMissing(t *testing.T) {
	srv := newFixtureServer(http.StatusOK, `{"columns":["value"],"rows":[[1]],"count":1}`, false, true)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result := Validate(ctx, Options{
		BaseURL:                  srv.URL + "/visual.html",
		VisualMode:               "dynamic",
		Token:                    "public-token",
		ExpectedRequestSubstring: "/api/test",
	})
	if result.Valid {
		t.Fatalf("expected invalid result for missing controls, got %+v", result)
	}
	if len(result.MissingControls) == 0 {
		t.Fatal("expected missing controls to be reported")
	}
}

func TestValidateStaticModeDoesNotRequireDynamicControls(t *testing.T) {
	srv := newFixtureServer(http.StatusOK, `{"columns":["value"],"rows":[[1]],"count":1}`, false, true)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result := Validate(ctx, Options{
		BaseURL:                  srv.URL + "/visual.html",
		VisualMode:               "static",
		Token:                    "public-token",
		ExpectedRequestSubstring: "/api/test",
	})
	if !result.Valid {
		t.Fatalf("expected static mode to skip dynamic controls, got %+v", result)
	}
	if len(result.MissingControls) != 0 {
		t.Fatalf("did not expect missing dynamic controls for static mode, got %+v", result.MissingControls)
	}
	if !result.LiveFetchSkipped || result.SkipReason != "non_dynamic_mode" {
		t.Fatalf("expected static mode live fetch skip, got %+v", result)
	}
}

func TestValidatePublicLiveFetch(t *testing.T) {
	token := strings.TrimSpace(os.Getenv("MCP_JWE_TOKEN"))
	if token == "" {
		t.Skip("MCP_JWE_TOKEN not set")
	}
	html := `<!doctype html>
<html><body>
<section id="query-ledger"><span class="ledger-status">pending</span></section>
<footer>
  <input type="password" id="tokenInput">
  <textarea id="sqlTextarea">SELECT 1 AS value</textarea>
  <button id="fetchBtn" type="button" onclick="runQuery()">Fetch</button>
  <div id="statusText">ready</div>
</footer>
<script>
function runQuery() {
  const token = document.getElementById('tokenInput').value.trim();
  const sql = document.getElementById('sqlTextarea').value.trim();
  localStorage.setItem('OnTimeAnalystDashboard::auth::jwe', token);
  fetch('https://mcp.demo.altinity.cloud/' + encodeURIComponent(token) + '/openapi/execute_query?query=' + encodeURIComponent(sql))
    .then(async (resp) => {
      if (!resp.ok) {
        throw new Error(await resp.text());
      }
      const payload = await resp.json();
      document.querySelector('.ledger-status').textContent = 'loaded';
      document.getElementById('statusText').textContent = 'loaded ' + (payload.count ?? 0) + ' rows';
      const card = document.createElement('div');
      card.className = 'card';
      card.textContent = JSON.stringify(payload.rows || []);
      document.body.insertBefore(card, document.querySelector('footer'));
    })
    .catch((err) => {
      document.getElementById('statusText').textContent = 'error ' + err.message;
    });
}
</script>
</body></html>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	result := Validate(ctx, Options{
		BaseURL:    srv.URL,
		VisualMode: "dynamic",
		Token:      token,
	})
	if !result.Valid {
		t.Fatalf("expected public live fetch to succeed, got errors: %v warnings: %v", result.Errors, result.Warnings)
	}
}

func newFixtureServer(statusCode int, body string, throwException bool, missingControls bool) *httptest.Server {
	html := fixtureHTML(throwException, missingControls)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/visual.html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(html))
		case "/api/test":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte(body))
		default:
			http.NotFound(w, r)
		}
	}))
}

func fixtureHTML(throwException bool, missingControls bool) string {
	if missingControls {
		return `<!doctype html><html><body><section id="query-ledger"></section><script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe')</script></body></html>`
	}
	extraScript := ""
	if throwException {
		extraScript = `setTimeout(() => { throw new Error("fixture boom"); }, 0);`
	}
	return fmt.Sprintf(`<!doctype html>
<html>
<body>
  <section class="card" id="query-ledger">
    <span class="ledger-status">pending</span>
  </section>
  <footer>
    <input type="password" id="tokenInput">
    <textarea id="sqlTextarea">SELECT 1</textarea>
    <button id="fetchBtn" type="button" onclick="runQuery()">Fetch</button>
    <div id="statusText">idle</div>
  </footer>
  <script>
    %s
    async function runQuery() {
      const token = document.getElementById('tokenInput').value.trim();
      localStorage.setItem('OnTimeAnalystDashboard::auth::jwe', token);
      const resp = await fetch('/api/test?query=' + encodeURIComponent(document.getElementById('sqlTextarea').value.trim()));
      if (!resp.ok) {
        document.getElementById('statusText').textContent = 'error ' + resp.status;
        return;
      }
      const payload = await resp.json();
      document.querySelector('.ledger-status').textContent = 'loaded';
      document.getElementById('statusText').textContent = 'loaded';
      const card = document.createElement('article');
      card.className = 'card';
      card.textContent = JSON.stringify(payload.rows || []);
      document.body.insertBefore(card, document.querySelector('footer'));
    }
  </script>
</body>
</html>`, extraScript)
}
