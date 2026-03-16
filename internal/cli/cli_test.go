package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadPresentationArtifactsRejectsStaleFallbackFiles(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "report.md")
	htmlPath := filepath.Join(tmpDir, "visual.html")
	if err := os.WriteFile(reportPath, []byte("old report"), 0o644); err != nil {
		t.Fatalf("write report: %v", err)
	}
	if err := os.WriteFile(htmlPath, []byte("<html>old</html>"), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	notBefore := time.Now().Add(2 * time.Second)
	if _, _, err := loadPresentationArtifacts("", tmpDir, notBefore); err == nil {
		t.Fatalf("expected stale fallback files to be rejected")
	}
}

func TestLoadPresentationArtifactsAcceptsFreshFallbackFiles(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "report.md")
	htmlPath := filepath.Join(tmpDir, "visual.html")
	notBefore := time.Now()
	time.Sleep(20 * time.Millisecond)

	report := "```report\n# Fresh\n```\n"
	html := "<!doctype html>\n<html><body>fresh</body></html>\n"
	if err := os.WriteFile(reportPath, []byte(report), 0o644); err != nil {
		t.Fatalf("write report: %v", err)
	}
	if err := os.WriteFile(htmlPath, []byte(html), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	gotReport, gotHTML, err := loadPresentationArtifacts("", tmpDir, notBefore)
	if err != nil {
		t.Fatalf("expected fresh fallback files to be accepted: %v", err)
	}
	if gotReport != "# Fresh" {
		t.Fatalf("unexpected report content: %q", gotReport)
	}
	if gotHTML != "<!doctype html>\n<html><body>fresh</body></html>" {
		t.Fatalf("unexpected html content: %q", gotHTML)
	}
}
