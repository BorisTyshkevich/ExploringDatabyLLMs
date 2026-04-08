package render

import (
	"fmt"
	"sort"
	"strings"

	"qforge/internal/model"
)

func RenderMonitoringReport(question model.Question, summaries []model.QueryResultSummary) string {
	lines := []string{"# " + question.Meta.Title, ""}
	for i, item := range summaries {
		if i > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, "> "+strings.TrimSpace(item.Subquestion), "")
		lines = append(lines, strings.TrimSpace(item.AnswerMarkdown), "")
		lines = append(lines, fmt.Sprintf("- Rows returned: %d", item.RowCount))
		if len(item.ResultColumns) > 0 {
			lines = append(lines, fmt.Sprintf("- Columns: %s", strings.Join(item.ResultColumns, ", ")))
		}
		if len(item.FirstRow) > 0 {
			result := model.CanonicalResult{
				Columns: item.ResultColumns,
				Rows:    []map[string]any{item.FirstRow},
			}
			lines = append(lines, "", renderResultTableMarkdown(result, 1))
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n")) + "\n"
}

func renderResultTableMarkdown(result model.CanonicalResult, limit int) string {
	if len(result.Columns) == 0 {
		return "No columns returned."
	}
	header := "| " + strings.Join(result.Columns, " | ") + " |"
	separatorParts := make([]string, len(result.Columns))
	for i := range separatorParts {
		separatorParts[i] = "---"
	}
	separator := "| " + strings.Join(separatorParts, " | ") + " |"
	lines := []string{header, separator}
	if len(result.Rows) == 0 {
		emptyCells := make([]string, len(result.Columns))
		for i := range emptyCells {
			emptyCells[i] = ""
		}
		if len(emptyCells) > 0 {
			emptyCells[0] = "_no rows_"
		}
		lines = append(lines, "| "+strings.Join(emptyCells, " | ")+" |")
		return strings.Join(lines, "\n")
	}
	rowLimit := minInt(len(result.Rows), limit)
	for i := 0; i < rowLimit; i++ {
		lines = append(lines, "| "+strings.Join(markdownCells(result.Columns, result.Rows[i]), " | ")+" |")
	}
	if len(result.Rows) > rowLimit {
		lines = append(lines, "", fmt.Sprintf("_Showing %d of %d rows._", rowLimit, len(result.Rows)))
	}
	return strings.Join(lines, "\n")
}

func markdownCells(columns []string, row map[string]any) []string {
	cells := make([]string, len(columns))
	for i, column := range columns {
		cells[i] = markdownEscapeCell(formatValue(row[column]))
	}
	return cells
}

func formatValue(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case []any:
		parts := make([]string, len(v))
		for i := range v {
			parts[i] = formatValue(v[i])
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]any:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			parts = append(parts, fmt.Sprintf("%s=%s", key, formatValue(v[key])))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		return fmt.Sprint(v)
	}
}

func markdownEscapeCell(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	value = strings.ReplaceAll(value, "\n", "<br>")
	value = strings.ReplaceAll(value, "|", "\\|")
	return value
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

