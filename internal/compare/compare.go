package compare

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"qforge/internal/datasets"
	"qforge/internal/model"
	"qforge/internal/querylog"
	"qforge/internal/runs"
)

type RunSummary struct {
	Manifest model.RunManifest      `json:"manifest"`
	Result   model.CanonicalResult  `json:"result"`
	Metrics  *model.QueryLogMetrics `json:"metrics,omitempty"`
}

type Report struct {
	GeneratedAt string       `json:"generated_at"`
	Day         string       `json:"day"`
	Question    string       `json:"question,omitempty"`
	Runs        []RunSummary `json:"runs"`
}

func Generate(ctx context.Context, repoRoot, outPrefix, day, questionFilter, explicitMCPURL string) (Report, error) {
	runDirs, err := discoverRunDirs(repoRoot, day, questionFilter)
	if err != nil {
		return Report{}, err
	}
	items := make([]RunSummary, 0, len(runDirs))
	for _, runDir := range runDirs {
		manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
		if err != nil {
			return Report{}, err
		}
		var result model.CanonicalResult
		resultBytes, err := os.ReadFile(filepath.Join(runDir, "result.json"))
		if err == nil {
			if err := json.Unmarshal(resultBytes, &result); err != nil {
				return Report{}, err
			}
		}
		var metrics *model.QueryLogMetrics
		if manifest.LogComment != "" {
			cfg, err := datasets.Load(repoRoot, manifest.Dataset)
			if err != nil {
				return Report{}, err
			}
			mcpURL, token, err := datasets.ResolveMCPURL(cfg, explicitMCPURL)
			if err != nil {
				return Report{}, err
			}
			metrics, err = querylog.FetchLatest(ctx, mcpURL, token, manifest.LogComment)
			if err != nil {
				return Report{}, err
			}
		}
		items = append(items, RunSummary{
			Manifest: manifest,
			Result:   result,
			Metrics:  metrics,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Manifest.QuestionID == items[j].Manifest.QuestionID {
			if items[i].Manifest.Runner == items[j].Manifest.Runner {
				return items[i].Manifest.Model < items[j].Manifest.Model
			}
			return items[i].Manifest.Runner < items[j].Manifest.Runner
		}
		return items[i].Manifest.QuestionID < items[j].Manifest.QuestionID
	})
	report := Report{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Day:         day,
		Question:    questionFilter,
		Runs:        items,
	}
	if err := writeOutputs(outPrefix, report); err != nil {
		return Report{}, err
	}
	return report, nil
}

func discoverRunDirs(repoRoot, day, questionFilter string) ([]string, error) {
	root := filepath.Join(repoRoot, "runs", day)
	var runDirs []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == "." {
			return nil
		}
		manifestPath := filepath.Join(path, "manifest.json")
		if _, err := os.Stat(manifestPath); err == nil {
			manifest, err := runs.ReadManifest(manifestPath)
			if err != nil {
				return err
			}
			if questionFilter == "" || manifest.QuestionID == questionFilter || manifest.QuestionSlug == questionFilter {
				runDirs = append(runDirs, path)
			}
			return filepath.SkipDir
		}
		return nil
	})
	if os.IsNotExist(err) {
		return nil, nil
	}
	return runDirs, err
}

func writeOutputs(outPrefix string, report Report) error {
	if outPrefix == "" {
		outPrefix = "runs/compare"
	}
	jsonPath := outPrefix + ".json"
	mdPath := outPrefix + ".md"
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		return err
	}
	var md strings.Builder
	md.WriteString("# qforge Compare Report\n\n")
	md.WriteString(fmt.Sprintf("- Generated: `%s`\n", report.GeneratedAt))
	md.WriteString(fmt.Sprintf("- Day: `%s`\n", report.Day))
	if report.Question != "" {
		md.WriteString(fmt.Sprintf("- Question: `%s`\n", report.Question))
	}
	md.WriteString("\n| question | runner | model | status | rows | duration_ms | read_rows | memory_usage |\n")
	md.WriteString("| --- | --- | --- | --- | ---: | ---: | ---: | ---: |\n")
	for _, item := range report.Runs {
		duration := int64(0)
		readRows := int64(0)
		memory := int64(0)
		if item.Metrics != nil {
			duration = item.Metrics.QueryDurationMS
			readRows = item.Metrics.ReadRows
			memory = item.Metrics.MemoryUsage
		}
		md.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %d | %d | %d | %d |\n",
			item.Manifest.QuestionID,
			item.Manifest.Runner,
			item.Manifest.Model,
			item.Manifest.Status,
			item.Result.RowCount,
			duration,
			readRows,
			memory,
		))
	}
	return os.WriteFile(mdPath, []byte(md.String()), 0o644)
}
