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
		Dir:    filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
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
	if !strings.Contains(got, "Stay within the configured dataset scope.") {
		t.Fatalf("expected shared core scaffold, got: %s", got)
	}
	if !strings.Contains(got, "Inspect the live schema first") {
		t.Fatalf("expected SQL-specific scaffold, got: %s", got)
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
		Dir:          filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta:         model.QuestionMeta{ID: "q003", Title: "Delta ATL", VisualType: "html_heatmap"},
		ReportPrompt: "Report guidance.",
		VisualPrompt: "Visual guidance.",
	}
	result := model.CanonicalResult{
		Columns:     []string{"RowType", "Dest"},
		GeneratedAt: time.Now(),
	}
	dataset := model.DatasetConfig{
		PrimaryTable:    "default.ontime_v2",
		ForbiddenTables: "default.ontime",
	}
	got, err := BuildPresentationPrompt(question, dataset, result, "SELECT *\nFROM default.ontime_v2")
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Stay within the configured dataset scope.") {
		t.Fatalf("expected shared core scaffold, got: %s", got)
	}
	if !strings.Contains(got, "default.ontime_v2") || !strings.Contains(got, "default.ontime") {
		t.Fatalf("expected dataset substitutions, got: %s", got)
	}
	if !strings.Contains(got, "The report is a Markdown template") {
		t.Fatalf("expected presentation-specific scaffold, got: %s", got)
	}
	if !strings.Contains(got, "ontime-analyst-dashboard") {
		t.Fatalf("expected skill reference, got: %s", got)
	}
	if !strings.Contains(got, "SELECT *") || !strings.Contains(got, "FROM default.ontime_v2") {
		t.Fatalf("expected saved sql to be embedded in prompt, got: %s", got)
	}
	if strings.Contains(got, "Return exactly this fenced section:") {
		t.Fatalf("did not expect SQL-only fenced section instructions in presentation prompt, got: %s", got)
	}
	if strings.Contains(got, "qforge-result-data") || strings.Contains(got, "__QFORGE_DEFAULT_SQL__") {
		t.Fatalf("did not expect legacy injected JSON contract, got: %s", got)
	}
	if !strings.Contains(got, "Report guidance.") || !strings.Contains(got, "Visual guidance.") {
		t.Fatalf("expected question prompt sections, got: %s", got)
	}
}
