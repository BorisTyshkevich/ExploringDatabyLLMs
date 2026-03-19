package render

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"qforge/internal/model"
)

var reportPlaceholderPattern = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.]+)\s*\}\}`)

var allowedReportPlaceholders = map[string]struct{}{
	"row_count":        {},
	"generated_at":     {},
	"columns_csv":      {},
	"question_title":   {},
	"data_overview_md": {},
	"result_table_md":  {},
}

func ValidateReportTemplate(template string, metrics model.AnalysisMetrics) error {
	matches := reportPlaceholderPattern.FindAllStringSubmatch(template, -1)
	var unknown []string
	for _, match := range matches {
		name := match[1]
		if strings.HasPrefix(name, "metric.") {
			metricName := strings.TrimPrefix(name, "metric.")
			if metricName == "" {
				unknown = appendUnique(unknown, name)
				continue
			}
			if _, ok := metrics.NamedValues[metricName]; !ok {
				unknown = appendUnique(unknown, name)
			}
			continue
		}
		if _, ok := allowedReportPlaceholders[name]; !ok {
			unknown = appendUnique(unknown, name)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return fmt.Errorf("report template uses unsupported placeholders: %s", strings.Join(unknown, ", "))
	}
	if strings.Count(template, "{{result_table_md}}") > 1 {
		return fmt.Errorf("report template may include {{result_table_md}} at most once")
	}
	return nil
}

func RenderReport(template string, question model.Question, result model.CanonicalResult, metrics model.AnalysisMetrics) string {
	dataOverviewMD := renderDataOverviewMarkdown(result)
	resultTableMD := renderResultTableMarkdown(result, 20)
	replacements := []string{
		"{{row_count}}", fmt.Sprintf("%d", result.RowCount),
		"{{generated_at}}", result.GeneratedAt.Format("2006-01-02T15:04:05Z"),
		"{{columns_csv}}", strings.Join(result.Columns, ", "),
		"{{question_title}}", question.Meta.Title,
		"{{data_overview_md}}", dataOverviewMD,
		"{{result_table_md}}", resultTableMD,
	}
	for key, value := range metrics.NamedValues {
		replacements = append(replacements, "{{metric."+key+"}}", value)
	}
	replacer := strings.NewReplacer(replacements...)
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

func appendUnique(items []string, value string) []string {
	for _, item := range items {
		if item == value {
			return items
		}
	}
	return append(items, value)
}
