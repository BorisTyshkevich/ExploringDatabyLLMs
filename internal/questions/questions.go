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
	root := filepath.Join(repoRoot, "questions")
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
	promptPath := filepath.Join(dir, "prompt.md")
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
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		return model.Question{}, err
	}
	reportPromptBytes, _ := os.ReadFile(reportPromptPath)
	visualPromptBytes, _ := os.ReadFile(visualPromptPath)
	return model.Question{
		Dir:                 dir,
		Meta:                meta,
		Prompt:              strings.TrimSpace(string(promptBytes)),
		ReportPrompt:        strings.TrimSpace(string(reportPromptBytes)),
		VisualPrompt:        strings.TrimSpace(string(visualPromptBytes)),
		PresentationEnabled: requiresArtifact(meta.ArtifactsRequired, "report.md") || requiresArtifact(meta.ArtifactsRequired, "visual.html"),
	}, nil
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
