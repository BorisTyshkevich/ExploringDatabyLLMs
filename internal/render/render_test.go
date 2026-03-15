package render

import (
	"strings"
	"testing"
	"time"

	"qforge/internal/model"
)

func TestRenderReport(t *testing.T) {
	question := model.Question{Meta: model.QuestionMeta{Title: "Test"}}
	result := model.CanonicalResult{
		RowCount:    7,
		Columns:     []string{"a", "b"},
		GeneratedAt: time.Date(2026, 3, 14, 12, 0, 0, 0, time.UTC),
		Rows: []map[string]any{
			{"a": "x", "b": 42},
		},
	}
	got := RenderReport("Rows={{row_count}} Columns={{columns_csv}} Title={{question_title}}", question, result)
	if !strings.Contains(got, "Rows=7") || !strings.Contains(got, "Columns=a, b") || !strings.Contains(got, "Title=Test") {
		t.Fatalf("unexpected output: %s", got)
	}
	if !strings.Contains(got, "## Data Overview") || !strings.Contains(got, "## Result Rows") {
		t.Fatalf("expected default markdown data sections, got: %s", got)
	}
	if !strings.Contains(got, "| a | b |") || !strings.Contains(got, "| x | 42 |") {
		t.Fatalf("expected markdown table rows, got: %s", got)
	}
}

func TestRenderReportWithExplicitMarkdownPlaceholders(t *testing.T) {
	question := model.Question{Meta: model.QuestionMeta{Title: "Test"}}
	result := model.CanonicalResult{
		RowCount:    2,
		Columns:     []string{"airport", "hops"},
		GeneratedAt: time.Date(2026, 3, 14, 12, 0, 0, 0, time.UTC),
		Rows: []map[string]any{
			{"airport": "ATL", "hops": 5},
			{"airport": "DFW", "hops": 4},
		},
	}
	template := "# {{question_title}}\n\n{{data_overview_md}}\n\n{{result_table_md}}\n"
	got := RenderReport(template, question, result)
	if strings.Count(got, "## Data Overview") != 0 {
		t.Fatalf("did not expect default sections to be appended when placeholders exist: %s", got)
	}
	if !strings.Contains(got, "- Rows returned: 2") {
		t.Fatalf("expected markdown overview bullet list, got: %s", got)
	}
	if !strings.Contains(got, "| airport | hops |") || !strings.Contains(got, "| ATL | 5 |") {
		t.Fatalf("expected rendered markdown table, got: %s", got)
	}
}
