package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"qforge/internal/model"
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

func TestLoadAnalysisArtifactAcceptsValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	payload, err := json.Marshal(model.AnalysisArtifact{
		SQL:            "SELECT 1",
		ReportMarkdown: "# Title\n\n{{data_overview_md}}",
		Metrics: model.AnalysisMetrics{
			NamedValues: map[string]string{"max_hops": "8"},
		},
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	path := filepath.Join(tmpDir, "answer.raw.json")
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write answer.raw.json: %v", err)
	}
	got, err := loadAnalysisArtifact(path)
	if err != nil {
		t.Fatalf("expected valid analysis json: %v", err)
	}
	if got.SQL != "SELECT 1" || got.ReportMarkdown == "" {
		t.Fatalf("unexpected artifact: %+v", got)
	}
	if got.Metrics.NamedValues["max_hops"] != "8" {
		t.Fatalf("expected metrics to be preserved: %+v", got)
	}
}

func TestLoadAnalysisArtifactNormalizesEscapedMultilineStrings(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "answer.raw.json")
	payload := `{
  "sql": "SELECT\\n  1",
  "report_markdown": "# Title\\n\\n{{data_overview_md}}"
}`
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatalf("write answer.raw.json: %v", err)
	}
	got, err := loadAnalysisArtifact(path)
	if err != nil {
		t.Fatalf("expected escaped multiline strings to be accepted: %v", err)
	}
	if got.SQL != "SELECT\n  1" {
		t.Fatalf("expected sql newlines to be normalized, got: %q", got.SQL)
	}
	if got.ReportMarkdown != "# Title\n\n{{data_overview_md}}" {
		t.Fatalf("expected report markdown newlines to be normalized, got: %q", got.ReportMarkdown)
	}
}

func TestLoadAnalysisArtifactRejectsMissingFields(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "answer.raw.json")
	if err := os.WriteFile(path, []byte("{\"sql\":\"SELECT 1\"}"), 0o644); err != nil {
		t.Fatalf("write answer.raw.json: %v", err)
	}
	_, err := loadAnalysisArtifact(path)
	if err == nil {
		t.Fatalf("expected missing report_markdown to fail")
	}
}

func TestLoadAnalysisArtifactRejectsFencedJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "answer.raw.json")
	if err := os.WriteFile(path, []byte("```json\n{\"sql\":\"SELECT 1\",\"report_markdown\":\"# Title\"}\n```"), 0o644); err != nil {
		t.Fatalf("write answer.raw.json: %v", err)
	}
	_, err := loadAnalysisArtifact(path)
	if err == nil {
		t.Fatalf("expected fenced json file to fail")
	}
}
