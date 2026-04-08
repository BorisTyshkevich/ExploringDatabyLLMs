package render

import (
	"strings"
	"testing"

	"qforge/internal/model"
)

func TestRenderMonitoringReport(t *testing.T) {
	question := model.Question{Meta: model.QuestionMeta{Title: "Test Question"}}
	summaries := []model.QueryResultSummary{
		{
			ID:             "main",
			Subquestion:    "Which airport is worst?",
			AnswerMarkdown: "MDW is the worst.",
			SQL:            "SELECT 1",
			RowCount:       3,
			ResultColumns:  []string{"Origin", "OTP"},
			FirstRow:       map[string]any{"Origin": "MDW", "OTP": "76.5"},
		},
	}
	got := RenderMonitoringReport(question, summaries)
	if !strings.Contains(got, "# Test Question") {
		t.Fatalf("expected title, got: %s", got)
	}
	if !strings.Contains(got, "MDW is the worst.") {
		t.Fatalf("expected answer markdown, got: %s", got)
	}
	if !strings.Contains(got, "| Origin | OTP |") {
		t.Fatalf("expected result table, got: %s", got)
	}
}
