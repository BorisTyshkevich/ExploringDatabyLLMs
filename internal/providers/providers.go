package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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
	return p.run(ctx, req, req.Prompt)
}

func (p cliProvider) GeneratePresentation(ctx context.Context, req model.ProviderRequest) (model.ProviderResponse, error) {
	return p.run(ctx, req, req.Prompt)
}

func (p cliProvider) run(ctx context.Context, req model.ProviderRequest, prompt string) (model.ProviderResponse, error) {
	switch p.name {
	case "codex":
		return p.runCodex(ctx, req, prompt)
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

func (p cliProvider) runCodex(ctx context.Context, req model.ProviderRequest, prompt string) (model.ProviderResponse, error) {
	bin := p.bin(req)
	answerFile := filepath.Join(req.OutDir, "provider.raw.md")
	args := []string{"-c", fmt.Sprintf("mcp_servers.%s.url=%q", req.MCPServerName, req.MCPURL)}
	if req.MCPToken != "" {
		args = append(args, "-c", fmt.Sprintf("mcp_servers.%s.headers.Authorization=%q", req.MCPServerName, "Bearer "+req.MCPToken))
	}
	args = append(args, "exec", "--color", "never", "--output-last-message", answerFile, "--model", req.Model, "-")
	logf(req.Verbose, "provider=codex phase=start bin=%s model=%s output=%s", bin, req.Model, answerFile)
	startedAt := time.Now()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdin = bytes.NewBufferString(prompt)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	raw := stdout.String()
	if data, readErr := os.ReadFile(answerFile); readErr == nil && len(data) > 0 {
		raw = string(data)
	}
	if err != nil {
		logf(req.Verbose, "provider=codex phase=done status=warning elapsed=%s err=%v", time.Since(startedAt).Round(time.Millisecond), err)
		return model.ProviderResponse{RawOutput: raw, Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, err
	}
	logf(req.Verbose, "provider=codex phase=done status=ok elapsed=%s", time.Since(startedAt).Round(time.Millisecond))
	return model.ProviderResponse{RawOutput: raw, Stdout: stdout.String(), Stderr: stderr.String(), CLIBin: bin}, nil
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
