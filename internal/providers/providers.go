package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"qforge/internal/extract"
	"qforge/internal/model"
	verbosepkg "qforge/internal/verbose"
)

type Provider interface {
	GenerateSQL(context.Context, model.ProviderRequest) (model.ProviderResponse, error)
	GeneratePresentation(context.Context, model.ProviderRequest) (model.ProviderResponse, error)
}

func New(name string) (Provider, error) {
	switch name {
	case "codex":
		return cliProvider{name: name, defaultBin: "codex"}, nil
	case "claude":
		return cliProvider{name: name, defaultBin: "claude"}, nil
	case "gemini":
		return cliProvider{name: name, defaultBin: "gemini"}, nil
	default:
		return nil, fmt.Errorf("unsupported runner: %s", name)
	}
}

type cliProvider struct {
	name       string
	defaultBin string
}

func (p cliProvider) GenerateSQL(ctx context.Context, req model.ProviderRequest) (model.ProviderResponse, error) {
	return p.run(ctx, req, req.Prompt, codexAnalysisComplete(req.OutDir))
}

func (p cliProvider) GeneratePresentation(ctx context.Context, req model.ProviderRequest) (model.ProviderResponse, error) {
	return p.run(ctx, req, req.Prompt, codexVisualComplete)
}

func (p cliProvider) run(ctx context.Context, req model.ProviderRequest, prompt string, codexComplete func(string) bool) (model.ProviderResponse, error) {
	switch p.name {
	case "codex":
		return p.runCodex(ctx, req, prompt, codexComplete)
	case "claude":
		return p.runClaude(ctx, req, prompt)
	case "gemini":
		return p.runGemini(ctx, req, prompt)
	default:
		return model.ProviderResponse{}, fmt.Errorf("unsupported runner: %s", p.name)
	}
}

func (p cliProvider) bin(req model.ProviderRequest) string {
	if req.CLIBin != "" {
		return req.CLIBin
	}
	return p.defaultBin
}

func (p cliProvider) runCodex(ctx context.Context, req model.ProviderRequest, prompt string, isComplete func(string) bool) (model.ProviderResponse, error) {
	bin := p.bin(req)
	answerFile := filepath.Join(req.OutDir, "provider.raw.md")
	_ = os.Remove(answerFile)
	args := []string{"-c", fmt.Sprintf("mcp_servers.%s.url=%q", req.MCPServerName, req.MCPURL)}
	if req.MCPToken != "" {
		args = append(args, "-c", fmt.Sprintf("mcp_servers.%s.headers.Authorization=%q", req.MCPServerName, "Bearer "+req.MCPToken))
	}
	args = append(args, "exec", "--color", "never", "--output-last-message", answerFile, "--model", req.Model, "-")
	logf(req.Verbose, req.Model, "provider=codex phase=start bin=%s model=%s output=%s", bin, req.Model, answerFile)
	startedAt := time.Now()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Dir = req.OutDir
	cmd.Stdin = strings.NewReader(prompt)
	stdoutFile, err := os.CreateTemp("", "qforge-codex-stdout-*")
	if err != nil {
		return model.ProviderResponse{}, err
	}
	defer func() {
		_ = stdoutFile.Close()
		_ = os.Remove(stdoutFile.Name())
	}()
	stderrFile, err := os.CreateTemp("", "qforge-codex-stderr-*")
	if err != nil {
		return model.ProviderResponse{}, err
	}
	defer func() {
		_ = stderrFile.Close()
		_ = os.Remove(stderrFile.Name())
	}()
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	if err := cmd.Start(); err != nil {
		return model.ProviderResponse{}, err
	}
	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var completedRaw string
	var completedAt time.Time

	for {
		select {
		case err := <-waitCh:
			stdoutText := readFileText(stdoutFile.Name())
			stderrText := readFileText(stderrFile.Name())
			raw := codexRawOutput(stdoutText, answerFile)
			resp := model.ProviderResponse{RawOutput: raw, Stdout: stdoutText, Stderr: stderrText, CLIBin: bin}
			if isComplete(raw) {
				if err != nil {
					logf(req.Verbose, req.Model, "provider=codex phase=done status=recovered elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				} else {
					logf(req.Verbose, req.Model, "provider=codex phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
				}
				return resp, nil
			}
			if err != nil {
				logf(req.Verbose, req.Model, "provider=codex phase=done status=warning elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				return resp, err
			}
			logf(req.Verbose, req.Model, "provider=codex phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
			return resp, nil
		case <-ticker.C:
			raw := codexRawOutput("", answerFile)
			if !isComplete(raw) {
				completedRaw = ""
				completedAt = time.Time{}
				continue
			}
			if completedAt.IsZero() {
				completedRaw = raw
				completedAt = time.Now()
				continue
			}
			if raw != completedRaw {
				completedRaw = raw
				completedAt = time.Now()
				continue
			}
			if time.Since(completedAt) < 2*time.Second {
				continue
			}
			logf(req.Verbose, req.Model, "provider=codex phase=completion-file status=stable elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
			terminateProcess(cmd.Process)
			select {
			case err := <-waitCh:
				resp := model.ProviderResponse{RawOutput: completedRaw, Stdout: readFileText(stdoutFile.Name()), Stderr: readFileText(stderrFile.Name()), CLIBin: bin}
				if err != nil {
					logf(req.Verbose, req.Model, "provider=codex phase=done status=recovered elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				} else {
					logf(req.Verbose, req.Model, "provider=codex phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
				}
				return resp, nil
			case <-time.After(2 * time.Second):
				terminateProcess(cmd.Process)
				err := <-waitCh
				resp := model.ProviderResponse{RawOutput: completedRaw, Stdout: readFileText(stdoutFile.Name()), Stderr: readFileText(stderrFile.Name()), CLIBin: bin}
				logf(req.Verbose, req.Model, "provider=codex phase=done status=recovered elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				return resp, nil
			}
		case <-ctx.Done():
			err := <-waitCh
			stdoutText := readFileText(stdoutFile.Name())
			stderrText := readFileText(stderrFile.Name())
			raw := codexRawOutput(stdoutText, answerFile)
			resp := model.ProviderResponse{RawOutput: raw, Stdout: stdoutText, Stderr: stderrText, CLIBin: bin}
			if isComplete(raw) && errors.Is(ctx.Err(), context.DeadlineExceeded) {
				logf(req.Verbose, req.Model, "provider=codex phase=done status=recovered elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				return resp, nil
			}
			return resp, err
		}
	}
}

func codexRawOutput(stdoutText, answerFile string) string {
	raw := stdoutText
	if data, readErr := os.ReadFile(answerFile); readErr == nil && len(data) > 0 {
		raw = string(data)
	}
	return raw
}

func readFileText(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func codexAnalysisComplete(outDir string) func(string) bool {
	answerPath := filepath.Join(outDir, "answer.raw.json")
	return func(string) bool {
		data, err := os.ReadFile(answerPath)
		if err != nil {
			return false
		}
		var artifact model.AnalysisArtifact
		if err := json.Unmarshal(data, &artifact); err != nil {
			return false
		}
		return strings.TrimSpace(artifact.SQL) != "" && strings.TrimSpace(artifact.ReportMarkdown) != ""
	}
}

func codexVisualComplete(raw string) bool {
	_, htmlErr := extract.Block(raw, "html")
	return htmlErr == nil
}

func terminateProcess(proc *os.Process) {
	if proc == nil {
		return
	}
	_ = syscall.Kill(-proc.Pid, syscall.SIGINT)
	time.Sleep(200 * time.Millisecond)
	_ = syscall.Kill(-proc.Pid, syscall.SIGKILL)
}

func (p cliProvider) runClaude(ctx context.Context, req model.ProviderRequest, prompt string) (model.ProviderResponse, error) {
	bin := p.bin(req)
	tmpDir, err := os.MkdirTemp("", "qforge-claude-*")
	if err != nil {
		return model.ProviderResponse{}, err
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "claude-mcp.json")
	payload := map[string]any{
		"mcpServers": map[string]any{
			req.MCPServerName: map[string]any{
				"type": "http",
				"url":  req.MCPURL,
			},
		},
	}
	if req.MCPToken != "" {
		payload["mcpServers"].(map[string]any)[req.MCPServerName].(map[string]any)["headers"] = map[string]string{
			"Authorization": "Bearer " + req.MCPToken,
		}
	}
	configBytes, _ := json.Marshal(payload)
	if err := os.WriteFile(configPath, configBytes, 0o644); err != nil {
		return model.ProviderResponse{}, err
	}
	args := []string{
		"--print",
		"--model", req.Model,
		"--permission-mode", "bypassPermissions",
		"--output-format", "text",
		"--setting-sources", "user,project,local",
		"--mcp-config", configPath,
		"--strict-mcp-config",
		"--no-session-persistence",
		prompt,
	}
	logf(req.Verbose, req.Model, "provider=claude phase=start bin=%s model=%s config=%s", bin, req.Model, configPath)
	startedAt := time.Now()
	cmd := exec.CommandContext(ctx, bin, args...)
	var stdout, stderr bytes.Buffer
	liveStdout := newLiveLogWriter(req.Verbose, req.Model, "claude", "stdout", os.Stdout)
	liveStderr := newLiveLogWriter(req.Verbose, req.Model, "claude", "stderr", os.Stderr)
	cmd.Dir = req.OutDir
	cmd.Stdout = io.MultiWriter(&stdout, liveStdout)
	cmd.Stderr = io.MultiWriter(&stderr, liveStderr)
	err = cmd.Run()
	liveStdout.Flush()
	liveStderr.Flush()
	if err != nil {
		logf(req.Verbose, req.Model, "provider=claude phase=done status=warning elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
		logProviderDetails(req.Verbose, req.Model, "claude", stdout.String(), stderr.String())
		return model.ProviderResponse{RawOutput: stdout.String(), Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
	}
	logf(req.Verbose, req.Model, "provider=claude phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
	return model.ProviderResponse{RawOutput: stdout.String(), Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
}

func (p cliProvider) runGemini(ctx context.Context, req model.ProviderRequest, prompt string) (model.ProviderResponse, error) {
	bin := p.bin(req)
	if err := p.runGeminiMCP(ctx, req, "add"); err != nil {
		return model.ProviderResponse{}, err
	}
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := p.runGeminiMCP(cleanupCtx, req, "remove"); err != nil {
			logf(req.Verbose, req.Model, "provider=gemini phase=mcp_cleanup status=warning err=%v", err)
		}
	}()
	args := []string{
		"--model", req.Model,
		"--prompt", prompt,
		"--allowed-mcp-server-names", req.MCPServerName,
		"--approval-mode", "yolo",
		"--output-format", "text",
	}
	logf(req.Verbose, req.Model, "provider=gemini phase=start bin=%s model=%s workdir=%s", bin, req.Model, req.OutDir)
	startedAt := time.Now()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = req.OutDir
	var stdout, stderr bytes.Buffer
	liveStdout := newLiveLogWriter(req.Verbose, req.Model, "gemini", "stdout", os.Stdout)
	liveStderr := newLiveLogWriter(req.Verbose, req.Model, "gemini", "stderr", os.Stderr)
	cmd.Stdout = io.MultiWriter(&stdout, liveStdout)
	cmd.Stderr = io.MultiWriter(&stderr, liveStderr)
	err := cmd.Run()
	liveStdout.Flush()
	liveStderr.Flush()
	if err != nil {
		logf(req.Verbose, req.Model, "provider=gemini phase=done status=warning elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
		logProviderDetails(req.Verbose, req.Model, "gemini", stdout.String(), stderr.String())
		return model.ProviderResponse{RawOutput: stdout.String(), Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
	}
	logf(req.Verbose, req.Model, "provider=gemini phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
	return model.ProviderResponse{RawOutput: stdout.String(), Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
}

func (p cliProvider) runGeminiMCP(ctx context.Context, req model.ProviderRequest, action string) error {
	bin := p.bin(req)
	var args []string
	switch action {
	case "add":
		args = []string{"mcp", "--transport", "http", "add", req.MCPServerName, req.MCPURL}
	case "remove":
		args = []string{"mcp", "remove", req.MCPServerName}
	default:
		return fmt.Errorf("unsupported gemini mcp action: %s", action)
	}
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = req.OutDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("gemini mcp %s %s: %w: %s%s", action, req.MCPServerName, err, stdout.String(), stderr.String())
	}
	return nil
}

func logf(enabled bool, model, format string, args ...any) {
	if !enabled {
		return
	}
	verbosepkg.Printf(os.Stdout, time.Now, model, format, args...)
}

func logProviderDetails(enabled bool, model, providerName, stdoutText, stderrText string) {
	if !enabled {
		return
	}
	if summary := summarizeProviderFailure(stderrText, stdoutText); summary != "" {
		verbosepkg.Printf(os.Stdout, time.Now, model, "provider=%s detail=%q", providerName, summary)
	}
}

func summarizeProviderFailure(parts ...string) string {
	for _, part := range parts {
		s := strings.TrimSpace(part)
		if s == "" {
			continue
		}
		lower := strings.ToLower(s)
		switch {
		case strings.Contains(lower, "terminalquotaerror"):
			return firstMatchingLine(s, "TerminalQuotaError")
		case strings.Contains(lower, "quota will reset"):
			return firstMatchingLine(s, "quota will reset")
		case strings.Contains(lower, "quota"):
			return firstMatchingLine(s, "quota")
		case strings.Contains(lower, "rate limit"):
			return firstMatchingLine(s, "rate limit")
		case strings.Contains(lower, "authentication"):
			return firstMatchingLine(s, "authentication")
		case strings.Contains(lower, "unauthorized"):
			return firstMatchingLine(s, "unauthorized")
		case strings.Contains(lower, "forbidden"):
			return firstMatchingLine(s, "forbidden")
		case strings.Contains(lower, "error when talking to"):
			return firstMatchingLine(s, "Error when talking to")
		}

		for _, line := range strings.Split(s, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			l := strings.ToLower(line)
			if strings.Contains(l, "yolo mode is enabled") || strings.Contains(l, "loaded cached credentials") {
				continue
			}
			return truncate(line, 240)
		}
	}
	return ""
}

func firstMatchingLine(s, needle string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(strings.ToLower(line), strings.ToLower(needle)) {
			return truncate(line, 240)
		}
	}
	return truncate(strings.TrimSpace(s), 240)
}

func truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "...(truncated)"
}

type liveLogWriter struct {
	enabled  bool
	model    string
	provider string
	stream   string
	target   io.Writer
	mu       sync.Mutex
	buf      bytes.Buffer
}

func newLiveLogWriter(enabled bool, model, provider, stream string, target io.Writer) *liveLogWriter {
	return &liveLogWriter{
		enabled:  enabled,
		model:    model,
		provider: provider,
		stream:   stream,
		target:   target,
	}
}

func (w *liveLogWriter) Write(p []byte) (int, error) {
	if !w.enabled {
		return len(p), nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, err := w.buf.Write(p); err != nil {
		return 0, err
	}
	for {
		line, err := w.buf.ReadString('\n')
		if err != nil {
			w.buf.WriteString(line)
			return len(p), nil
		}
		verbosepkg.Printf(w.target, time.Now, w.model, "provider=%s stream=%s %s", w.provider, w.stream, strings.TrimRight(line, "\n"))
	}
}

func (w *liveLogWriter) Flush() {
	if !w.enabled {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buf.Len() == 0 {
		return
	}
	verbosepkg.Printf(w.target, time.Now, w.model, "provider=%s stream=%s %s", w.provider, w.stream, strings.TrimRight(w.buf.String(), "\n"))
	w.buf.Reset()
}
