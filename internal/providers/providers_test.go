package providers

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"qforge/internal/model"
)

func TestRunCodexRecoversFromStableVisualOutputFile(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "fake-codex.sh")
	script := "#!/usr/bin/env bash\n" +
		"set -euo pipefail\n" +
		"out=\"\"\n" +
		"while [[ $# -gt 0 ]]; do\n" +
		"  if [[ \"$1\" == \"--output-last-message\" ]]; then\n" +
		"    out=\"$2\"\n" +
		"    shift 2\n" +
		"    continue\n" +
		"  fi\n" +
		"  shift\n" +
		"done\n" +
		"cat >/dev/null\n" +
		"cat >\"$out\" <<'EOF'\n" +
		"```html\n" +
		"<!doctype html>\n" +
		"<html><body>ok</body></html>\n" +
		"```\n" +
		"EOF\n" +
		"sleep 30\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake codex: %v", err)
	}

	req := model.ProviderRequest{
		OutDir:        tmpDir,
		Model:         "gpt-5.4",
		MCPURL:        "https://example.invalid/http",
		MCPServerName: "altinity_ontime_demo",
		CLIBin:        scriptPath,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	start := time.Now()
	resp, err := cliProvider{name: "codex", defaultBin: scriptPath}.GeneratePresentation(ctx, req)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("GeneratePresentation returned error: %v", err)
	}
	if elapsed > 8*time.Second {
		t.Fatalf("expected recovery before context timeout, elapsed=%s", elapsed)
	}
	if !codexVisualComplete(resp.RawOutput) {
		t.Fatalf("expected complete presentation output, got: %s", resp.RawOutput)
	}
	if !strings.Contains(resp.RawOutput, "<!doctype html>") {
		t.Fatalf("expected html in raw output, got: %s", resp.RawOutput)
	}
}

func TestCodexCompletionChecks(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "answer.raw.json"), []byte("{\"sql\":\"SELECT 1\",\"report_markdown\":\"# Title\\n\\n{{data_overview_md}}\"}"), 0o644); err != nil {
		t.Fatalf("write answer.raw.json: %v", err)
	}
	if !codexAnalysisComplete(tmpDir)("") {
		t.Fatalf("expected analysis completion checker to accept answer.raw.json")
	}

	presentationRaw := "```html\n<!doctype html>\n<html></html>\n```"
	if !codexVisualComplete(presentationRaw) {
		t.Fatalf("expected presentation completion checker to accept fenced html")
	}

	if codexAnalysisComplete(t.TempDir())("") {
		t.Fatalf("did not expect analysis checker to accept incomplete json")
	}
	if codexVisualComplete("```report\nonly report\n```") {
		t.Fatalf("did not expect visual checker to accept non-html output")
	}
}

func TestRunCodexRecoversFromStableAnalysisFile(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "fake-codex.sh")
	script := "#!/usr/bin/env bash\n" +
		"set -euo pipefail\n" +
		"while [[ $# -gt 0 ]]; do shift; done\n" +
		"cat >/dev/null\n" +
		"cat > answer.raw.json <<'EOF'\n" +
		"{\"sql\":\"SELECT 1\",\"report_markdown\":\"# Title\\n\\n{{data_overview_md}}\",\"metrics\":{\"named_values\":{\"max_hops\":\"8\"}}}\n" +
		"EOF\n" +
		"echo 'analysis artifact written'\n" +
		"sleep 30\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake codex: %v", err)
	}

	req := model.ProviderRequest{
		OutDir:        tmpDir,
		Model:         "gpt-5.4",
		MCPURL:        "https://example.invalid/http",
		MCPServerName: "altinity_ontime_demo",
		CLIBin:        scriptPath,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	start := time.Now()
	resp, err := cliProvider{name: "codex", defaultBin: scriptPath}.GenerateSQL(ctx, req)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("GenerateSQL returned error: %v", err)
	}
	if elapsed > 8*time.Second {
		t.Fatalf("expected recovery before context timeout, elapsed=%s", elapsed)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "answer.raw.json")); err != nil {
		t.Fatalf("expected answer.raw.json in outDir: %v", err)
	}
	if !strings.Contains(resp.Stdout, "analysis artifact written") {
		t.Fatalf("unexpected stdout: %q", resp.Stdout)
	}
}

func TestRunClaudeUsesOutDirAsWorkingDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "fake-claude.sh")
	script := "#!/usr/bin/env bash\n" +
		"set -euo pipefail\n" +
		"printf '# Report\\n{{data_overview_md}}\\n' > report.md\n" +
		"printf '<!doctype html>\\n<html><body>ok</body></html>\\n' > visual.html\n" +
		"printf 'generated files in cwd\\n'\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake claude: %v", err)
	}

	outDir := filepath.Join(tmpDir, "run-001")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir outDir: %v", err)
	}

	req := model.ProviderRequest{
		OutDir:        outDir,
		Model:         "sonnet",
		MCPURL:        "https://example.invalid/http",
		MCPServerName: "altinity_ontime_demo",
		CLIBin:        scriptPath,
	}

	resp, err := cliProvider{name: "claude", defaultBin: scriptPath}.GeneratePresentation(context.Background(), req)
	if err != nil {
		t.Fatalf("GeneratePresentation returned error: %v", err)
	}
	if !strings.Contains(resp.RawOutput, "generated files in cwd") {
		t.Fatalf("unexpected raw output: %q", resp.RawOutput)
	}
	if _, err := os.Stat(filepath.Join(outDir, "report.md")); err != nil {
		t.Fatalf("expected report.md in outDir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "visual.html")); err != nil {
		t.Fatalf("expected visual.html in outDir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "report.md")); !os.IsNotExist(err) {
		t.Fatalf("did not expect report.md outside outDir, err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "visual.html")); !os.IsNotExist(err) {
		t.Fatalf("did not expect visual.html outside outDir, err=%v", err)
	}
}
