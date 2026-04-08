package compare

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"qforge/internal/model"
	"qforge/internal/prompts"
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
		"compare_json_path":     runsRelativePath(repoRoot, compareJSONPath),
		"question_prompt_path":  githubPromptURL(repoRoot, filepath.Join(question.Dir, "prompt.md")),
		"report_prompt_path":    optionalGithubPromptURL(repoRoot, filepath.Join(question.Dir, "report_prompt.md")),
		"visual_prompt_path":    optionalGithubPromptURL(repoRoot, filepath.Join(question.Dir, "visual_prompt.md")),
		"compare_contract_path": optionalGithubPromptURL(repoRoot, filepath.Join(question.Dir, "compare.yaml")),
		"run_dirs_md":           bulletList(runsRelativePaths(repoRoot, runDirs(report.Runs))),
		"query_sql_paths_md":    bulletList(runsRelativePaths(repoRoot, querySQLPaths(report.Runs))),
		"report_md_paths_md":    bulletList(runsRelativePaths(repoRoot, reportMDPaths(report.Runs))),
		"visual_html_paths_md":  bulletList(runsRelativePaths(repoRoot, visualHTMLPaths(report.Runs))),
		"result_json_paths_md":  bulletList(runsRelativePaths(repoRoot, resultJSONPaths(report.Runs))),
		"compare_summary_md":    renderMarkdown(report),
	}
	return prompts.RenderTemplate(string(data), values), nil
}

func publishPath(repoRoot, path string) string {
	if path == "" {
		return path
	}
	rel, err := filepath.Rel(repoRoot, path)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(rel)
	}
	return path
}

// runsRelativePath converts absolute path to path relative to question dir in runs repo.
// e.g., /path/to/runs/2026-03-17/q001_slug/claude/opus/run-001 -> claude/opus/run-001
func runsRelativePath(repoRoot, path string) string {
	rel := publishPath(repoRoot, path)
	// Strip "runs/YYYY-MM-DD/qXXX_slug/" prefix to get path relative to question dir
	parts := strings.SplitN(rel, "/", 4) // runs / date / qslug / rest
	if len(parts) >= 4 && parts[0] == "runs" {
		return parts[3]
	}
	return rel
}

// runsRelativePaths applies runsRelativePath to a slice of paths.
func runsRelativePaths(repoRoot string, values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, runsRelativePath(repoRoot, value))
	}
	return out
}

// githubPromptURL converts prompt path to absolute GitHub URL.
func githubPromptURL(repoRoot, path string) string {
	rel := publishPath(repoRoot, path)
	if strings.HasPrefix(rel, "prompts/") {
		return "https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/" + rel
	}
	return rel
}

// optionalGithubPromptURL returns a GitHub URL for the prompt if it exists, otherwise "(not present)".
func optionalGithubPromptURL(repoRoot, path string) string {
	if _, err := os.Stat(path); err == nil {
		return githubPromptURL(repoRoot, path)
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
