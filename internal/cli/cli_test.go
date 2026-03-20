package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"qforge/internal/model"
	verbosepkg "qforge/internal/verbose"
)

func TestVerbosePrefixFormat(t *testing.T) {
	got := verbosepkg.PrefixAt(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.FixedZone("UTC+1", 3600)), "opus")
	want := "2026-01-01 00:00:00 opus"
	if got != want {
		t.Fatalf("unexpected prefix: got %q want %q", got, want)
	}
}

func TestModelLabelForRunnersUsesResolvedModels(t *testing.T) {
	got, err := modelLabelForRunners([]string{"codex", "claude", "gemini"}, []string{"", "sonnet"})
	if err != nil {
		t.Fatalf("modelLabelForRunners returned error: %v", err)
	}
	want := "gpt-5.4,sonnet,gemini-3.1-pro-preview"
	if got != want {
		t.Fatalf("unexpected label: got %q want %q", got, want)
	}
}

func TestModelLabelForRunnersHandlesDuplicateRunnersWithDistinctModels(t *testing.T) {
	got, err := modelLabelForRunners([]string{"claude", "claude", "codex"}, []string{"opus", "sonnet", "gpt-5.4"})
	if err != nil {
		t.Fatalf("modelLabelForRunners returned error: %v", err)
	}
	want := "opus,sonnet,gpt-5.4"
	if got != want {
		t.Fatalf("unexpected label: got %q want %q", got, want)
	}
}

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

func TestLoadCompareReportArtifactPrefersFreshFile(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "compare_report.md")
	notBefore := time.Now()
	time.Sleep(20 * time.Millisecond)

	report := "# Long Report\n\nFull compare body.\n"
	if err := os.WriteFile(reportPath, []byte(report), 0o644); err != nil {
		t.Fatalf("write compare_report.md: %v", err)
	}

	got, err := loadCompareReportArtifact("Short stdout note", tmpDir, notBefore)
	if err != nil {
		t.Fatalf("expected fresh compare_report.md to be accepted: %v", err)
	}
	if got != "# Long Report\n\nFull compare body." {
		t.Fatalf("unexpected compare report content: %q", got)
	}
}

func TestLoadCompareReportArtifactFallsBackToStdout(t *testing.T) {
	tmpDir := t.TempDir()
	got, err := loadCompareReportArtifact("Short stdout note", tmpDir, time.Now())
	if err != nil {
		t.Fatalf("expected stdout fallback to work: %v", err)
	}
	if got != "Short stdout note" {
		t.Fatalf("unexpected stdout fallback content: %q", got)
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

func TestLoadAnalysisArtifactAcceptsNumericMetricValues(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "answer.raw.json")
	payload := `{
  "sql": "SELECT 1",
  "report_markdown": "# Title\n\n{{metric.max_hops}}",
  "metrics": {
    "named_values": {
      "max_hops": 8
    }
  }
}`
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatalf("write answer.raw.json: %v", err)
	}
	got, err := loadAnalysisArtifact(path)
	if err != nil {
		t.Fatalf("expected numeric metrics to be accepted: %v", err)
	}
	if got.Metrics.NamedValues["max_hops"] != "8" {
		t.Fatalf("expected numeric metric to be coerced to string, got: %#v", got.Metrics.NamedValues)
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

func TestBuildVisualInputSummaryCapturesShapeHints(t *testing.T) {
	question := model.Question{
		Meta: model.QuestionMeta{Title: "Q001", VisualMode: "dynamic"},
	}
	result := model.CanonicalResult{
		Columns:  []string{"Date", "DepTimes", "Route"},
		RowCount: 2,
		Rows: []map[string]any{
			{"Date": "2024-12-01T00:00:00Z", "DepTimes": []any{543.0, 810.0}, "Route": "ISP-BWI-SEA"},
			{"Date": "2024-02-18T00:00:00Z", "DepTimes": []any{621.0, 801.0}, "Route": "CLE-BNA-DEN"},
		},
	}

	got := buildVisualInputSummary(question, result)
	if got.QuestionTitle != "Q001" || got.RowCount != 2 {
		t.Fatalf("unexpected summary header: %+v", got)
	}
	if len(got.SampleRows) != 2 {
		t.Fatalf("expected two sample rows, got %+v", got.SampleRows)
	}
	if got.FieldShapeNotes["Date"] != "ISO-like timestamp string" {
		t.Fatalf("expected timestamp shape note, got %+v", got.FieldShapeNotes)
	}
	if got.FieldShapeNotes["DepTimes"] != "array field" {
		t.Fatalf("expected array shape note, got %+v", got.FieldShapeNotes)
	}
	if got.ModeHint == "" {
		t.Fatalf("expected mode hint in summary: %+v", got)
	}
}

func TestEnsureVisualInputSummaryBackfillsMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "visual_input.json")
	question := model.Question{
		Meta: model.QuestionMeta{Title: "Q001", VisualMode: "dynamic"},
	}
	result := model.CanonicalResult{
		Columns:  []string{"Date", "DepTimes", "Route"},
		RowCount: 1,
		Rows: []map[string]any{
			{"Date": "2024-12-01T00:00:00Z", "DepTimes": []any{543.0, 810.0}, "Route": "ISP-BWI-SEA"},
		},
	}

	got, err := ensureVisualInputSummary(path, question, result)
	if err != nil {
		t.Fatalf("expected missing visual_input.json to be backfilled: %v", err)
	}
	if got.QuestionTitle != "Q001" || got.RowCount != 1 {
		t.Fatalf("unexpected summary: %+v", got)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected visual_input.json to be written: %v", err)
	}
	var persisted model.VisualInputSummary
	if err := json.Unmarshal(bytes, &persisted); err != nil {
		t.Fatalf("expected visual_input.json to be valid json: %v", err)
	}
	if persisted.FieldShapeNotes["DepTimes"] != "array field" {
		t.Fatalf("expected persisted shape notes, got %+v", persisted.FieldShapeNotes)
	}
}
