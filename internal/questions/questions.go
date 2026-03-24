package questions

import (
	"fmt"
	"os"
	"path/filepath"
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
	case model.AnalysisModeMultiQueryJSON, model.AnalysisModeTemplateFiles, model.AnalysisModeManualTemplate:
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
			if mode == model.AnalysisModeMultiQueryJSON {
				items := extractDashboardQuestions(reportPrompt)
				if len(items) == 0 {
					return nil, fmt.Errorf("load %s: missing ## Dashboard Questions section for analysis_mode %q", dir, mode)
				}
				return items, nil
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
	if mode == model.AnalysisModeMultiQueryJSON && len(file.Subquestions) == 0 {
		return nil, fmt.Errorf("parse %s: subquestions list is required for analysis_mode %q", path, mode)
	}
	return file.Subquestions, nil
}

func extractDashboardQuestions(reportPrompt string) []model.QuestionSubquestion {
	lines := strings.Split(reportPrompt, "\n")
	inSection := false
	var items []model.QuestionSubquestion
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "## ") {
			if strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(line, "## ")), "Dashboard Questions") {
				inSection = true
				continue
			}
			if inSection {
				break
			}
		}
		if !inSection {
			continue
		}
		if !strings.HasPrefix(line, "- ") {
			continue
		}
		text := strings.TrimSpace(strings.TrimPrefix(line, "- "))
		if text == "" {
			continue
		}
		items = append(items, model.QuestionSubquestion{Text: text})
	}
	return items
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

func LoadCompareContract(question model.Question) (*model.CompareContract, error) {
	path := filepath.Join(question.Dir, "compare.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var file model.CompareContractFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if len(file.CompareContract.Normalization.NullEquivalents) == 0 {
		file.CompareContract.Normalization.NullEquivalents = []string{"", "NULL", "null"}
	}
	return &file.CompareContract, nil
}
