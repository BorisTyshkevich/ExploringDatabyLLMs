package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"qforge/internal/model"
)

const (
	commonPromptFile              = "common.md"
	commonSQLPromptFile           = "common_sql.md"
	commonPresentationPromptFile  = "common_presentation.md"
	commonVisualPromptFile        = "common_visual.md"
	commonVisualStaticPromptFile  = "common_visual_static.md"
	commonVisualDynamicPromptFile = "common_visual_dynamic.md"
)

func BuildSQLPrompt(question model.Question, dataset model.DatasetConfig) (string, error) {
	common, err := loadCommonPrompt(question, commonPromptFile)
	if err != nil {
		return "", err
	}
	commonSQL, err := loadCommonPrompt(question, commonSQLPromptFile)
	if err != nil {
		return "", err
	}
	values := map[string]string{
		"dataset_semantic_layer_md": datasetSemanticLayerMarkdown(dataset),
		"question_title":            question.Meta.Title,
		"report_prompt_md":          question.ReportPrompt,
		"report_placeholders":       "{{row_count}}, {{generated_at}}, {{columns_csv}}, {{question_title}}, {{data_overview_md}}, {{result_table_md}}",
	}
	sections := []string{
		RenderTemplate(common, values),
		RenderTemplate(commonSQL, values),
		question.Prompt,
	}
	return joinSections(sections), nil
}

func BuildPresentationPrompt(question model.Question, dataset model.DatasetConfig, result model.CanonicalResult, savedSQL string) (string, error) {
	return BuildVisualPrompt(question, dataset, result, savedSQL, "", "")
}

func BuildVisualPrompt(question model.Question, dataset model.DatasetConfig, result model.CanonicalResult, savedSQL, reportTemplate, savedAnalysisJSON string) (string, error) {
	common, err := loadCommonPrompt(question, commonPromptFile)
	if err != nil {
		return "", err
	}
	commonPresentation, err := loadCommonPrompt(question, commonPresentationPromptFile)
	if err != nil {
		return "", err
	}
	commonVisual, err := loadCommonPrompt(question, commonVisualPromptFile)
	if err != nil {
		return "", err
	}
	modePromptFile := commonVisualDynamicPromptFile
	if strings.EqualFold(strings.TrimSpace(question.Meta.VisualMode), "static") {
		modePromptFile = commonVisualStaticPromptFile
	}
	modeVisual, err := loadCommonPrompt(question, modePromptFile)
	if err != nil {
		return "", err
	}
	values := map[string]string{
		"dataset_semantic_layer_md": datasetSemanticLayerMarkdown(dataset),
		"question_title":            question.Meta.Title,
		"visual_mode":               strings.TrimSpace(question.Meta.VisualMode),
		"visual_type":               question.Meta.VisualType,
		"result_columns_csv":        strings.Join(result.Columns, ", "),
		"saved_sql":                 strings.TrimSpace(savedSQL),
		"saved_report_template":     strings.TrimSpace(reportTemplate),
		"saved_analysis_json":       strings.TrimSpace(savedAnalysisJSON),
		"visual_prompt_md":          question.VisualPrompt,
	}
	sections := []string{
		RenderTemplate(common, values),
		RenderTemplate(commonPresentation, values),
		RenderTemplate(commonVisual, values),
		RenderTemplate(modeVisual, values),
	}
	return joinSections(sections), nil
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

func datasetSemanticLayerMarkdown(dataset model.DatasetConfig) string {
	if strings.TrimSpace(dataset.SemanticLayer) == "" {
		return ""
	}
	return "Dataset semantic layer:\n\n" + strings.TrimSpace(dataset.SemanticLayer)
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
