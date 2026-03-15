package prompts

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"qforge/internal/model"
)

func TestBuildSQLPromptLoadsMarkdownAssets(t *testing.T) {
	question := model.Question{
		Dir:    filepath.Join("..", "..", "questions", "q003_delta_atl_departure_delay_hotspots"),
		Prompt: "Question-specific SQL guidance.",
	}
	dataset := model.DatasetConfig{
		PrimaryTable:    "default.ontime_v2",
		ForbiddenTables: "default.ontime",
	}
	got, err := BuildSQLPrompt(question, dataset)
	if err != nil {
		t.Fatalf("BuildSQLPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "You are running inside qforge.") {
		t.Fatalf("expected common SQL scaffold, got: %s", got)
	}
	if !strings.Contains(got, "Question-specific SQL guidance.") {
		t.Fatalf("expected question prompt section, got: %s", got)
	}
	if !strings.Contains(got, "default.ontime_v2") || !strings.Contains(got, "default.ontime") {
		t.Fatalf("expected dataset substitutions, got: %s", got)
	}
}

func TestBuildPresentationPromptLoadsMarkdownAssets(t *testing.T) {
	question := model.Question{
		Dir:          filepath.Join("..", "..", "questions", "q003_delta_atl_departure_delay_hotspots"),
		Meta:         model.QuestionMeta{ID: "q003", Title: "Delta ATL", VisualType: "html_heatmap"},
		ReportPrompt: "Report guidance.",
		VisualPrompt: "Visual guidance.",
	}
	result := model.CanonicalResult{
		Columns:     []string{"RowType", "Dest"},
		GeneratedAt: time.Now(),
	}
	got, err := BuildPresentationPrompt(question, result, "SELECT *\nFROM default.ontime_v2")
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "ontime-analyst-dashboard") {
		t.Fatalf("expected skill reference, got: %s", got)
	}
	if !strings.Contains(got, "SELECT *") || !strings.Contains(got, "FROM default.ontime_v2") {
		t.Fatalf("expected saved sql to be embedded in prompt, got: %s", got)
	}
	if strings.Contains(got, "qforge-result-data") || strings.Contains(got, "__QFORGE_DEFAULT_SQL__") {
		t.Fatalf("did not expect legacy injected JSON contract, got: %s", got)
	}
	if !strings.Contains(got, "Report guidance.") || !strings.Contains(got, "Visual guidance.") {
		t.Fatalf("expected question prompt sections, got: %s", got)
	}
}
