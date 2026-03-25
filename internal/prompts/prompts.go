package prompts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"qforge/internal/model"
)

const (
	commonPromptFile                 = "common.md"
	commonReportMultiQueryPromptFile = "common_report_multi_query.md"
	commonReportTemplatePromptFile   = "common_report_templates.md"
	commonReviewPromptFile           = "common_review.md"
	commonVisualPromptFile           = "common_visual.md"
	commonVisualHTMLPromptFile       = "common_visual_html.md"
	commonVisualReactPromptFile      = "common_visual_react.md"
	commonVisualMultiQueryPromptFile = "common_visual_multi_query.md"
	commonVisualStaticPromptFile     = "common_visual_static.md"
	commonVisualDynamicPromptFile    = "common_visual_dynamic.md"
)

func BuildSQLPrompt(question model.Question, dataset model.DatasetConfig, mode model.AnalysisMode) (string, error) {
	common, err := loadCommonPrompt(question, commonPromptFile)
	if err != nil {
		return "", err
	}
	contractPromptFile := commonReportMultiQueryPromptFile
	if mode == model.AnalysisModeTemplateFiles || mode == model.AnalysisModeManualTemplate {
		contractPromptFile = commonReportTemplatePromptFile
	}
	commonReport, err := loadCommonPrompt(question, contractPromptFile)
	if err != nil {
		return "", err
	}
	values := map[string]string{
		"dataset_name":        datasetPromptName(dataset),
		"question_title":      question.Meta.Title,
		"question_prompt_md":  questionPromptForAnalysis(question),
		"report_placeholders": "{{row_count}}, {{generated_at}}, {{columns_csv}}, {{question_title}}, {{data_overview_md}}, {{result_table_md}}",
	}
	sections := []string{
		RenderTemplate(common, values),
		RenderTemplate(commonReport, values),
	}
	return joinSections(sections), nil
}

func BuildPresentationPrompt(question model.Question, dataset model.DatasetConfig, result model.CanonicalResult, savedSQL, dynamicQueryEndpointTemplate string) (string, error) {
	return BuildVisualPrompt(question, dataset, result, savedSQL, dynamicQueryEndpointTemplate, model.VisualInputSummary{})
}

func BuildVisualPrompt(question model.Question, dataset model.DatasetConfig, result model.CanonicalResult, savedSQL, dynamicQueryEndpointTemplate string, visualInput model.VisualInputSummary) (string, error) {
	commonVisualFile := commonVisualPromptFile
	analysisMode := model.AnalysisMode(strings.TrimSpace(question.Meta.AnalysisMode))
	common, err := loadCommonPrompt(question, commonPromptFile)
	if err != nil {
		return "", err
	}
	commonVisual, err := loadCommonPrompt(question, commonVisualFile)
	if err != nil {
		return "", err
	}
	var multiQueryVisual string
	if analysisMode == model.AnalysisModeMultiQuery {
		multiQueryVisual, err = loadCommonPrompt(question, commonVisualMultiQueryPromptFile)
		if err != nil {
			return "", err
		}
	}
	modePromptFile := commonVisualDynamicPromptFile
	if strings.EqualFold(strings.TrimSpace(question.Meta.VisualMode), "static") {
		modePromptFile = commonVisualStaticPromptFile
	}
	modeVisual, err := loadCommonPrompt(question, modePromptFile)
	if err != nil {
		return "", err
	}
	targetPromptFile := commonVisualHTMLPromptFile
	if strings.EqualFold(strings.TrimSpace(question.Meta.PresentationTarget), "react") {
		targetPromptFile = commonVisualReactPromptFile
	}
	targetVisual, err := loadCommonPrompt(question, targetPromptFile)
	if err != nil {
		return "", err
	}
	values := map[string]string{
		"dataset_name":                    datasetPromptName(dataset),
		"question_title":                  question.Meta.Title,
		"visual_mode":                     strings.TrimSpace(question.Meta.VisualMode),
		"presentation_target":             presentationTarget(question.Meta.PresentationTarget),
		"visual_type":                     question.Meta.VisualType,
		"result_columns_csv":              strings.Join(result.Columns, ", "),
		"saved_sql":                       strings.TrimSpace(savedSQL),
		"dynamic_query_endpoint_template": strings.TrimSpace(dynamicQueryEndpointTemplate),
		"visual_input_summary_json":       visualInputSummaryJSON(visualInput),
		"visual_prompt_md":                question.VisualPrompt,
	}
	sections := []string{RenderTemplate(commonVisual, values)}
	sections = append([]string{RenderTemplate(common, values)}, sections...)
	if analysisMode == model.AnalysisModeMultiQuery {
		sections = append(sections, RenderTemplate(multiQueryVisual, values))
	}
	sections = append(sections, RenderTemplate(modeVisual, values))
	sections = append(sections, RenderTemplate(targetVisual, values))
	return joinSections(sections), nil
}

type ReviewPromptInputs struct {
	Question        model.Question
	AnalysisMode    model.AnalysisMode
	ReportMarkdown  string
	AnswerRawJSON   string
	AnalysisJSON    string
	QuerySQL        string
	ResultJSON      string
	VisualInputJSON string
	QueryFiles      []string
	ResultFiles     []string
}

func BuildReviewPrompt(inputs ReviewPromptInputs) (string, error) {
	common, err := loadCommonPrompt(inputs.Question, commonPromptFile)
	if err != nil {
		return "", err
	}
	review, err := loadCommonPrompt(inputs.Question, commonReviewPromptFile)
	if err != nil {
		return "", err
	}
	values := map[string]string{
		"question_title":     inputs.Question.Meta.Title,
		"dataset_name":       strings.TrimSpace(inputs.Question.Meta.Dataset),
		"question_prompt_md": questionPromptForAnalysis(inputs.Question),
	}
	sections := []string{RenderTemplate(common, values), RenderTemplate(review, values)}
	sections = append(sections,
		"Question-specific guidance:\n\n"+strings.TrimSpace(inputs.Question.Prompt),
		"Generated report.md:\n\n```md\n"+strings.TrimSpace(inputs.ReportMarkdown)+"\n```",
	)
	if inputs.AnalysisMode == model.AnalysisModeMultiQuery {
		if strings.TrimSpace(inputs.AnswerRawJSON) != "" {
			sections = append(sections, "Saved answer.raw.json:\n\n```json\n"+strings.TrimSpace(inputs.AnswerRawJSON)+"\n```")
		}
		if strings.TrimSpace(inputs.AnalysisJSON) != "" {
			sections = append(sections, "Saved analysis.json:\n\n```json\n"+strings.TrimSpace(inputs.AnalysisJSON)+"\n```")
		}
		if strings.TrimSpace(inputs.VisualInputJSON) != "" {
			sections = append(sections, "Saved visual_input.json:\n\n```json\n"+strings.TrimSpace(inputs.VisualInputJSON)+"\n```")
		}
		if len(inputs.QueryFiles) > 0 {
			files := append([]string(nil), inputs.QueryFiles...)
			sort.Strings(files)
			parts := make([]string, 0, len(files))
			for _, path := range files {
				parts = append(parts, "- `"+strings.TrimSpace(path)+"`")
			}
			sections = append(sections, "Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:\n\n"+strings.Join(parts, "\n"))
		}
		if len(inputs.ResultFiles) > 0 {
			files := append([]string(nil), inputs.ResultFiles...)
			sort.Strings(files)
			parts := make([]string, 0, len(files))
			for _, path := range files {
				parts = append(parts, "- `"+strings.TrimSpace(path)+"`")
			}
			sections = append(sections, "Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:\n\n"+strings.Join(parts, "\n"))
		}
	} else {
		if strings.TrimSpace(inputs.QuerySQL) != "" {
			sections = append(sections, "Saved query.sql:\n\n```sql\n"+strings.TrimSpace(inputs.QuerySQL)+"\n```")
		}
		if strings.TrimSpace(inputs.ResultJSON) != "" {
			sections = append(sections, "Saved result.json:\n\n```json\n"+strings.TrimSpace(inputs.ResultJSON)+"\n```")
		}
	}
	return joinSections(sections), nil
}

func visualInputSummaryJSON(summary model.VisualInputSummary) string {
	if summary.QuestionTitle == "" && len(summary.ResultColumns) == 0 && summary.RowCount == 0 && len(summary.SampleRows) == 0 && len(summary.FieldShapeNotes) == 0 && summary.ModeHint == "" {
		return "{}"
	}
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

func questionPromptForAnalysis(question model.Question) string {
	return strings.TrimSpace(question.Prompt)
}

func loadCommonPrompt(question model.Question, name string) (string, error) {
	questionsDir := filepath.Dir(question.Dir)
	path := filepath.Join(questionsDir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("load prompt asset %s: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// RenderTemplate substitutes {{key}} placeholders with values from the map.
func RenderTemplate(template string, values map[string]string) string {
	replacements := make([]string, 0, len(values)*2)
	for key, value := range values {
		replacements = append(replacements, "{{"+key+"}}", strings.TrimSpace(value))
	}
	return strings.TrimSpace(strings.NewReplacer(replacements...).Replace(template))
}

func joinSections(sections []string) string {
	var cleaned []string
	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}
		cleaned = append(cleaned, section)
	}
	return strings.Join(cleaned, "\n\n")
}

func datasetPromptName(dataset model.DatasetConfig) string {
	if strings.TrimSpace(dataset.DefaultDatabase) != "" {
		return strings.TrimSpace(dataset.DefaultDatabase)
	}
	if strings.TrimSpace(dataset.Name) != "" {
		return strings.TrimSpace(dataset.Name)
	}
	return "configured"
}

func presentationTarget(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "react") {
		return "react"
	}
	return "html"
}
