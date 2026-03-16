package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"qforge/internal/extract"
	"qforge/internal/model"
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
	return p.run(ctx, req, req.Prompt, codexSQLComplete)
}

func (p cliProvider) GeneratePresentation(ctx context.Context, req model.ProviderRequest) (model.ProviderResponse, error) {
	return p.run(ctx, req, req.Prompt, codexPresentationComplete)
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
	logf(req.Verbose, "provider=codex phase=start bin=%s model=%s output=%s", bin, req.Model, answerFile)
	startedAt := time.Now()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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
					logf(req.Verbose, "provider=codex phase=done status=recovered elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				} else {
					logf(req.Verbose, "provider=codex phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
				}
				return resp, nil
			}
			if err != nil {
				logf(req.Verbose, "provider=codex phase=done status=warning elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				return resp, err
			}
			logf(req.Verbose, "provider=codex phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
			return resp, nil
		case <-ticker.C:
			raw := codexRawOutput("", answerFile)
			if !isComplete(raw) {
				completedRaw = ""
				completedAt = time.Time{}
				continue
			}
			if raw != completedRaw {
				completedRaw = raw
				completedAt = time.Now()
				continue
			}
			if completedAt.IsZero() || time.Since(completedAt) < 2*time.Second {
				continue
			}
			logf(req.Verbose, "provider=codex phase=completion-file status=stable elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
			terminateProcess(cmd.Process)
			select {
			case err := <-waitCh:
				resp := model.ProviderResponse{RawOutput: completedRaw, Stdout: readFileText(stdoutFile.Name()), Stderr: readFileText(stderrFile.Name()), CLIBin: bin}
				if err != nil {
					logf(req.Verbose, "provider=codex phase=done status=recovered elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				} else {
					logf(req.Verbose, "provider=codex phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
				}
				return resp, nil
			case <-time.After(2 * time.Second):
				terminateProcess(cmd.Process)
				err := <-waitCh
				resp := model.ProviderResponse{RawOutput: completedRaw, Stdout: readFileText(stdoutFile.Name()), Stderr: readFileText(stderrFile.Name()), CLIBin: bin}
				logf(req.Verbose, "provider=codex phase=done status=recovered elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
				return resp, nil
			}
		case <-ctx.Done():
			err := <-waitCh
			stdoutText := readFileText(stdoutFile.Name())
			stderrText := readFileText(stderrFile.Name())
			raw := codexRawOutput(stdoutText, answerFile)
			resp := model.ProviderResponse{RawOutput: raw, Stdout: stdoutText, Stderr: stderrText, CLIBin: bin}
			if isComplete(raw) && errors.Is(ctx.Err(), context.DeadlineExceeded) {
				logf(req.Verbose, "provider=codex phase=done status=recovered elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
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

func codexSQLComplete(raw string) bool {
	_, err := extract.Block(raw, "sql")
	return err == nil
}

func codexPresentationComplete(raw string) bool {
	_, reportErr := extract.Block(raw, "report")
	_, htmlErr := extract.Block(raw, "html")
	return reportErr == nil && htmlErr == nil
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
	logf(req.Verbose, "provider=claude phase=start bin=%s model=%s config=%s", bin, req.Model, configPath)
	startedAt := time.Now()
	cmd := exec.CommandContext(ctx, bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		logf(req.Verbose, "provider=claude phase=done status=warning elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
		return model.ProviderResponse{RawOutput: stdout.String(), Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
	}
	logf(req.Verbose, "provider=claude phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
	return model.ProviderResponse{RawOutput: stdout.String(), Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
}

func (p cliProvider) runGemini(ctx context.Context, req model.ProviderRequest, prompt string) (model.ProviderResponse, error) {
	bin := p.bin(req)
	args := []string{
		"--model", req.Model,
		"--prompt", prompt,
		"--allowed-mcp-server-names", req.MCPServerName,
		"--approval-mode", "yolo",
		"--output-format", "text",
	}
	logf(req.Verbose, "provider=gemini phase=start bin=%s model=%s workdir=%s", bin, req.Model, req.OutDir)
	startedAt := time.Now()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = req.OutDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		logf(req.Verbose, "provider=gemini phase=done status=warning elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
		return model.ProviderResponse{RawOutput: stdout.String(), Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
	}
	logf(req.Verbose, "provider=gemini phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
	return model.ProviderResponse{RawOutput: stdout.String(), Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
}

func logf(verbose bool, format string, args ...any) {
	if !verbose {
		return
	}
	fmt.Printf("[qforge] "+format+"\n", args...)
}
