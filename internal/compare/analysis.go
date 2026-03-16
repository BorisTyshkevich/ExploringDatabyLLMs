package compare

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"qforge/internal/model"
)

const analysisPromptFile = "analysis_prompt.md"

func BuildAnalysisPrompt(repoRoot string, question model.Question, report Report, compareJSONPath string) (string, error) {
	templatePath := filepath.Join(repoRoot, "prompts", analysisPromptFile)
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("load analysis prompt %s: %w", templatePath, err)
	}
	values := map[string]string{
		"question_id":           question.Meta.ID,
		"question_slug":         question.Meta.Slug,
		"question_title":        question.Meta.Title,
		"compare_day":           report.Day,
		"compare_json_path":     compareJSONPath,
		"question_prompt_path":  filepath.Join(question.Dir, "prompt.md"),
		"report_prompt_path":    optionalPath(filepath.Join(question.Dir, "report_prompt.md")),
		"visual_prompt_path":    optionalPath(filepath.Join(question.Dir, "visual_prompt.md")),
		"compare_contract_path": optionalPath(filepath.Join(question.Dir, "compare.yaml")),
		"run_dirs_md":           bulletList(runDirs(report.Runs)),
		"query_sql_paths_md":    bulletList(querySQLPaths(report.Runs)),
		"report_md_paths_md":    bulletList(reportMDPaths(report.Runs)),
		"visual_html_paths_md":  bulletList(visualHTMLPaths(report.Runs)),
		"result_json_paths_md":  bulletList(resultJSONPaths(report.Runs)),
		"compare_summary_md":    renderMarkdown(report),
	}
	return renderTemplate(string(data), values), nil
}

func renderTemplate(template string, values map[string]string) string {
	replacements := make([]string, 0, len(values)*2)
	for key, value := range values {
		replacements = append(replacements, "{{"+key+"}}", strings.TrimSpace(value))
	}
	return strings.TrimSpace(strings.NewReplacer(replacements...).Replace(template))
}

func optionalPath(path string) string {
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return "(not present)"
}

func bulletList(values []string) string {
	if len(values) == 0 {
		return "- none"
	}
	var lines []string
	for _, value := range values {
		lines = append(lines, "- "+value)
	}
	return strings.Join(lines, "\n")
}

func runDirs(items []RunSummary) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.RunDir)
	}
	return out
}

func querySQLPaths(items []RunSummary) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, filepath.Join(item.RunDir, "query.sql"))
	}
	return out
}

func resultJSONPaths(items []RunSummary) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, filepath.Join(item.RunDir, "result.json"))
	}
	return out
}

func reportMDPaths(items []RunSummary) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, filepath.Join(item.RunDir, "report.md"))
	}
	return out
}

func visualHTMLPaths(items []RunSummary) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, filepath.Join(item.RunDir, "visual.html"))
	}
	return out
}
