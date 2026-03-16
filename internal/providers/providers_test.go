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

func TestRunCodexRecoversFromStablePresentationOutputFile(t *testing.T) {
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
		"```report\n" +
		"# Report\n" +
		"{{data_overview_md}}\n" +
		"```\n\n" +
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
	if !codexPresentationComplete(resp.RawOutput) {
		t.Fatalf("expected complete presentation output, got: %s", resp.RawOutput)
	}
	if !strings.Contains(resp.RawOutput, "<!doctype html>") {
		t.Fatalf("expected html in raw output, got: %s", resp.RawOutput)
	}
}

func TestCodexCompletionChecks(t *testing.T) {
	sqlRaw := "```sql\nSELECT 1\n```"
	if !codexSQLComplete(sqlRaw) {
		t.Fatalf("expected sql completion checker to accept fenced sql")
	}

	presentationRaw := "```report\n{{data_overview_md}}\n```\n\n```html\n<!doctype html>\n<html></html>\n```"
	if !codexPresentationComplete(presentationRaw) {
		t.Fatalf("expected presentation completion checker to accept fenced report/html")
	}

	if codexPresentationComplete("```report\nonly report\n```") {
		t.Fatalf("did not expect presentation checker to accept incomplete output")
	}
}
