package questions

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"qforge/internal/model"
)

func LoadAll(repoRoot string) ([]model.Question, error) {
	root := filepath.Join(repoRoot, "prompts")
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var items []model.Question
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		item, err := Load(filepath.Join(root, entry.Name()))
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Meta.ID < items[j].Meta.ID
	})
	return items, nil
}

func Resolve(repoRoot, ref string) (model.Question, error) {
	items, err := LoadAll(repoRoot)
	if err != nil {
		return model.Question{}, err
	}
	for _, item := range items {
		if item.Meta.ID == ref || item.Meta.Slug == ref || filepath.Base(item.Dir) == ref {
			return item, nil
		}
	}
	return model.Question{}, fmt.Errorf("unknown question: %s", ref)
}

func Load(dir string) (model.Question, error) {
	metaPath := filepath.Join(dir, "meta.yaml")
	reportPromptPath := filepath.Join(dir, "report_prompt.md")
	visualPromptPath := filepath.Join(dir, "visual_prompt.md")

	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		return model.Question{}, err
	}
	var meta model.QuestionMeta
	if err := yaml.Unmarshal(metaBytes, &meta); err != nil {
		return model.Question{}, fmt.Errorf("parse %s: %w", metaPath, err)
	}
	if strings.TrimSpace(meta.AnalysisMode) == "" {
		meta.AnalysisMode = string(model.AnalysisModeTemplateFiles)
	}
	switch model.AnalysisMode(strings.TrimSpace(meta.AnalysisMode)) {
	case model.AnalysisModeMultiQuery, model.AnalysisModeTemplateFiles:
	default:
		return model.Question{}, fmt.Errorf("parse %s: unsupported analysis_mode %q", metaPath, meta.AnalysisMode)
	}
	if strings.TrimSpace(meta.VisualMode) == "" {
		meta.VisualMode = "dynamic"
	}
	if strings.TrimSpace(meta.PresentationTarget) == "" {
		meta.PresentationTarget = "html"
	}
	switch strings.TrimSpace(meta.PresentationTarget) {
	case "html", "react":
	default:
		return model.Question{}, fmt.Errorf("parse %s: unsupported presentation_target %q", metaPath, meta.PresentationTarget)
	}
	reportPromptBytes, err := os.ReadFile(reportPromptPath)
	if err != nil {
		return model.Question{}, err
	}
	visualPromptBytes, _ := os.ReadFile(visualPromptPath)
	reportPrompt := strings.TrimSpace(string(reportPromptBytes))
	subquestions, err := loadSubquestions(dir, model.AnalysisMode(strings.TrimSpace(meta.AnalysisMode)), reportPrompt)
	if err != nil {
		return model.Question{}, err
	}
	reportEnabled := requiresArtifact(meta.ArtifactsRequired, "report.md")
	visualEnabled := requiresArtifact(meta.ArtifactsRequired, "visual.html")
	return model.Question{
		Dir:                 dir,
		Meta:                meta,
		Prompt:              reportPrompt,
		VisualPrompt:        strings.TrimSpace(string(visualPromptBytes)),
		Subquestions:        subquestions,
		PresentationEnabled: reportEnabled || visualEnabled,
		ReportEnabled:       reportEnabled,
		VisualEnabled:       visualEnabled,
	}, nil
}

func loadSubquestions(dir string, mode model.AnalysisMode, reportPrompt string) ([]model.QuestionSubquestion, error) {
	path := filepath.Join(dir, "subquestions.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if mode == model.AnalysisModeMultiQuery {
				return extractPromptSections(reportPrompt)
			}
			return nil, nil
		}
		return nil, err
	}
	var file struct {
		Subquestions []model.QuestionSubquestion `yaml:"subquestions"`
	}
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if mode == model.AnalysisModeMultiQuery && len(file.Subquestions) == 0 {
		return nil, fmt.Errorf("parse %s: subquestions list is required for analysis_mode %q", path, mode)
	}
	return file.Subquestions, nil
}

var sectionIDPattern = regexp.MustCompile(`^q[1-9][0-9]*$`)

func extractPromptSections(reportPrompt string) ([]model.QuestionSubquestion, error) {
	lines := strings.Split(reportPrompt, "\n")
	type section struct {
		id    string
		lines []string
	}
	var current *section
	var sections []section
	flush := func() {
		if current == nil {
			return
		}
		sections = append(sections, *current)
		current = nil
	}
	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "### ") {
			flush()
			id := strings.TrimSpace(strings.TrimPrefix(trimmed, "### "))
			current = &section{id: id}
			continue
		}
		if current != nil {
			current.lines = append(current.lines, line)
		}
	}
	flush()

	if len(sections) == 0 {
		text := strings.TrimSpace(reportPrompt)
		if text == "" {
			return nil, fmt.Errorf("missing question sections for multi_query prompt")
		}
		return []model.QuestionSubquestion{{ID: "q0", Text: text}}, nil
	}

	var items []model.QuestionSubquestion
	seen := map[string]struct{}{}
	mainCount := 0
	for _, sec := range sections {
		id := strings.TrimSpace(sec.id)
		if id == "" {
			return nil, fmt.Errorf("empty section id in multi_query prompt")
		}
		switch {
		case strings.EqualFold(id, "main"):
			id = "main"
			mainCount++
			if mainCount > 1 {
				return nil, fmt.Errorf("multi_query prompt may contain at most one ### main section")
			}
		case sectionIDPattern.MatchString(id):
		default:
			return nil, fmt.Errorf("unsupported multi_query section id %q", id)
		}
		if _, ok := seen[id]; ok {
			return nil, fmt.Errorf("duplicate multi_query section id %q", id)
		}
		seen[id] = struct{}{}
		text := strings.TrimSpace(strings.Join(sec.lines, "\n"))
		if text == "" {
			return nil, fmt.Errorf("multi_query section %q is empty", id)
		}
		items = append(items, model.QuestionSubquestion{ID: id, Text: text})
	}
	return items, nil
}

func requiresArtifact(required, name string) bool {
	parts := strings.Split(required, ",")
	for _, part := range parts {
		if strings.TrimSpace(part) == name {
			return true
		}
	}
	return false
}
