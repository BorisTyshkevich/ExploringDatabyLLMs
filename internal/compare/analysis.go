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
		"question_id":                question.Meta.ID,
		"question_slug":              question.Meta.Slug,
		"question_title":             question.Meta.Title,
		"compare_day":                report.Day,
		"compare_json_path":          repoRelativePath(repoRoot, compareJSONPath),
		"compare_json_url":           publishedBlobURL(publishedRelativePath(repoRoot, compareJSONPath)),
		"question_prompt_path":       repoRelativePath(repoRoot, filepath.Join(question.Dir, "prompt.md")),
		"question_prompt_url":        qforgeRepoURL(repoRoot, filepath.Join(question.Dir, "prompt.md")),
		"report_prompt_path":         optionalPath(repoRoot, filepath.Join(question.Dir, "report_prompt.md")),
		"report_prompt_url":          optionalURL(repoRoot, filepath.Join(question.Dir, "report_prompt.md")),
		"visual_prompt_path":         optionalPath(repoRoot, filepath.Join(question.Dir, "visual_prompt.md")),
		"visual_prompt_url":          optionalURL(repoRoot, filepath.Join(question.Dir, "visual_prompt.md")),
		"compare_contract_path":      optionalPath(repoRoot, filepath.Join(question.Dir, "compare.yaml")),
		"compare_contract_url":       optionalURL(repoRoot, filepath.Join(question.Dir, "compare.yaml")),
		"run_dirs_md":                bulletList(repoRelativePaths(repoRoot, runDirs(report.Runs))),
		"query_sql_paths_md":         bulletList(repoRelativePaths(repoRoot, querySQLPaths(report.Runs))),
		"report_md_paths_md":         bulletList(repoRelativePaths(repoRoot, reportMDPaths(report.Runs))),
		"visual_html_paths_md":       bulletList(repoRelativePaths(repoRoot, visualHTMLPaths(report.Runs))),
		"result_json_paths_md":       bulletList(repoRelativePaths(repoRoot, resultJSONPaths(report.Runs))),
		"published_run_artifacts_md": renderPublishedRunArtifacts(report.Runs),
		"compare_summary_md":         renderMarkdown(report),
	}
	return prompts.RenderTemplate(string(data), values), nil
}

func optionalPath(repoRoot, path string) string {
	if _, err := os.Stat(path); err == nil {
		return repoRelativePath(repoRoot, path)
	}
	return "(not present)"
}

func optionalURL(repoRoot, path string) string {
	if _, err := os.Stat(path); err == nil {
		if url := qforgeRepoURL(repoRoot, path); url != "" {
			return url
		}
	}
	return "(not present)"
}

func repoRelativePaths(repoRoot string, values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, repoRelativePath(repoRoot, value))
	}
	return out
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

func renderPublishedRunArtifacts(items []RunSummary) string {
	if len(items) == 0 {
		return "- none"
	}
	var b strings.Builder
	lastGroup := ""
	for _, item := range items {
		group := item.Runner + " / " + item.Model
		if group != lastGroup {
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString("### ")
			b.WriteString(group)
			b.WriteString("\n")
			lastGroup = group
		}
		runLabel := item.RunID
		if runLabel == "" {
			runLabel = "run"
		}
		b.WriteString("- `")
		b.WriteString(runLabel)
		b.WriteString("`\n")
		b.WriteString("  - Published links: ")
		b.WriteString(renderPublishedLinks(item.Artifacts))
		b.WriteString("\n")
		b.WriteString("  - Local verification: `")
		b.WriteString(strings.Join(existingLocalPaths(item.Artifacts), "`, `"))
		b.WriteString("`\n")
	}
	return strings.TrimSpace(b.String())
}

func renderPublishedLinks(links ArtifactLinks) string {
	var parts []string
	if links.QuerySQL.URL != "" {
		parts = append(parts, fmt.Sprintf("query.sql: %s", links.QuerySQL.URL))
	}
	if links.ReportMD.URL != "" {
		parts = append(parts, fmt.Sprintf("report.md: %s", links.ReportMD.URL))
	}
	if links.ResultJSON.URL != "" {
		parts = append(parts, fmt.Sprintf("result.json: %s", links.ResultJSON.URL))
	}
	if links.VisualHTML.URL != "" {
		parts = append(parts, fmt.Sprintf("visual.html: %s", links.VisualHTML.URL))
	}
	if len(parts) == 0 {
		return "(no published artifacts found)"
	}
	return strings.Join(parts, " | ")
}

func existingLocalPaths(links ArtifactLinks) []string {
	var out []string
	for _, ref := range []ArtifactRef{links.QuerySQL, links.ReportMD, links.ResultJSON, links.VisualHTML} {
		if ref.LocalPath != "" {
			out = append(out, ref.LocalPath)
		}
	}
	if len(out) == 0 {
		return []string{"(no local artifact files found)"}
	}
	return out
}
