package render

import (
	"fmt"
	"sort"
	"strings"

	"qforge/internal/model"
)

func RenderReport(template string, question model.Question, result model.CanonicalResult) string {
	dataOverviewMD := renderDataOverviewMarkdown(result)
	resultTableMD := renderResultTableMarkdown(result, 20)
	replacer := strings.NewReplacer(
		"{{row_count}}", fmt.Sprintf("%d", result.RowCount),
		"{{generated_at}}", result.GeneratedAt.Format("2006-01-02T15:04:05Z"),
		"{{columns_csv}}", strings.Join(result.Columns, ", "),
		"{{question_title}}", question.Meta.Title,
		"{{data_overview_md}}", dataOverviewMD,
		"{{result_table_md}}", resultTableMD,
	)
	rendered := replacer.Replace(template)
	if !strings.Contains(template, "{{data_overview_md}}") && !strings.Contains(template, "{{result_table_md}}") {
		rendered = strings.TrimRight(rendered, "\n") + "\n\n## Data Overview\n\n" + dataOverviewMD + "\n\n## Result Rows\n\n" + resultTableMD + "\n"
	}
	return rendered
}

func renderDataOverviewMarkdown(result model.CanonicalResult) string {
	lines := []string{
		fmt.Sprintf("- Rows returned: %d", result.RowCount),
		fmt.Sprintf("- Generated at: %s", result.GeneratedAt.Format("2006-01-02T15:04:05Z")),
		fmt.Sprintf("- Columns: %s", strings.Join(result.Columns, ", ")),
	}
	if len(result.Rows) > 0 {
		firstRow := summarizeRow(result)
		if firstRow != "" {
			lines = append(lines, fmt.Sprintf("- First row snapshot: %s", firstRow))
		}
	}
	return strings.Join(lines, "\n")
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

func summarizeRow(result model.CanonicalResult) string {
	row := result.Rows[0]
	parts := make([]string, 0, len(result.Columns))
	for _, column := range result.Columns {
		value := formatValue(row[column])
		if value == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", column, value))
		if len(parts) == 3 {
			break
		}
	}
	return strings.Join(parts, ", ")
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
