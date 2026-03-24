package prompts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"qforge/internal/model"
)

const (
	commonPromptFile                 = "common.md"
	commonReportPromptFile           = "common_report_json.md"
	commonReportMultiQueryPromptFile = "common_report_multi_query_json.md"
	commonReportTemplatePromptFile   = "common_report_templates.md"
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
	contractPromptFile := commonReportPromptFile
	if mode == model.AnalysisModeMultiQueryJSON {
		contractPromptFile = commonReportMultiQueryPromptFile
	} else if mode == model.AnalysisModeTemplateFiles || mode == model.AnalysisModeManualTemplate {
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
	if analysisMode == model.AnalysisModeMultiQueryJSON {
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
	if analysisMode == model.AnalysisModeMultiQueryJSON {
		sections = append(sections, RenderTemplate(multiQueryVisual, values))
	}
	sections = append(sections, RenderTemplate(modeVisual, values))
	sections = append(sections, RenderTemplate(targetVisual, values))
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
	prompt := strings.TrimSpace(question.Prompt)
	if !strings.EqualFold(strings.TrimSpace(question.Meta.AnalysisMode), string(model.AnalysisModeMultiQueryJSON)) {
		return prompt
	}
	if strings.Contains(strings.ToLower(prompt), "## dashboard questions") {
		return prompt
	}
	if len(question.Subquestions) == 0 {
		return prompt
	}
	lines := make([]string, 0, len(question.Subquestions)+2)
	lines = append(lines, "## Dashboard Questions", "")
	for _, item := range question.Subquestions {
		text := strings.TrimSpace(item.Text)
		if text == "" {
			continue
		}
		lines = append(lines, "- "+text)
	}
	section := strings.TrimSpace(strings.Join(lines, "\n"))
	if section == "## Dashboard Questions" {
		return prompt
	}
	if prompt == "" {
		return section
	}
	return prompt + "\n\n" + section
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
