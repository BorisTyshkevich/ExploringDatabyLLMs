package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadVisualArtifactRejectsStaleFallbackFile(t *testing.T) {
	tmpDir := t.TempDir()
	htmlPath := filepath.Join(tmpDir, "visual.html")
	if err := os.WriteFile(htmlPath, []byte("<html>old</html>"), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	notBefore := time.Now().Add(2 * time.Second)
	if _, err := loadVisualArtifact("", tmpDir, notBefore); err == nil {
		t.Fatalf("expected stale fallback files to be rejected")
	}
}

func TestLoadVisualArtifactAcceptsFreshFallbackFile(t *testing.T) {
	tmpDir := t.TempDir()
	htmlPath := filepath.Join(tmpDir, "visual.html")
	notBefore := time.Now()
	time.Sleep(20 * time.Millisecond)

	html := "<!doctype html>\n<html><body>fresh</body></html>\n"
	if err := os.WriteFile(htmlPath, []byte(html), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	gotHTML, err := loadVisualArtifact("", tmpDir, notBefore)
	if err != nil {
		t.Fatalf("expected fresh fallback files to be accepted: %v", err)
	}
	if gotHTML != "<!doctype html>\n<html><body>fresh</body></html>" {
		t.Fatalf("unexpected html content: %q", gotHTML)
	}
}
