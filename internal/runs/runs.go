package runs

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"qforge/internal/model"
)

func NextRunDir(repoRoot string, question model.Question, runner, modelName string, now time.Time) (string, error) {
	baseDir := filepath.Join(repoRoot, now.Format("2006-01-02"), question.Meta.Slug, runner, modelName)
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", err
	}
	var runNames []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "run-") {
			runNames = append(runNames, entry.Name())
		}
	}
	sort.Strings(runNames)
	next := 1
	if len(runNames) > 0 {
		last := runNames[len(runNames)-1]
		fmt.Sscanf(last, "run-%03d", &next)
		next++
	}
	dir := filepath.Join(baseDir, fmt.Sprintf("run-%03d", next))
	return dir, os.MkdirAll(dir, 0o755)
}

func DefaultArtifacts(outDir string, presentation bool) model.ArtifactPaths {
	artifacts := model.ArtifactPaths{
		PromptSQLRaw:    filepath.Join(outDir, "prompt.sql.md"),
		AnswerSQLRaw:    filepath.Join(outDir, "answer.sql.raw.md"),
		AnswerRawJSON:   filepath.Join(outDir, "answer.raw.json"),
		AnalysisJSON:    filepath.Join(outDir, "analysis.json"),
		QuerySQL:        filepath.Join(outDir, "query.sql"),
		ResultTSV:       filepath.Join(outDir, "result.tsv"),
		ResultJSON:      filepath.Join(outDir, "result.json"),
		VisualInputJSON: filepath.Join(outDir, "visual_input.json"),
		ManifestJSON:    filepath.Join(outDir, "manifest.json"),
		StdoutLog:       filepath.Join(outDir, "stdout.log"),
		StderrLog:       filepath.Join(outDir, "stderr.log"),
	}
	if presentation {
		artifacts.PromptPresentationRaw = filepath.Join(outDir, "prompt.presentation.md")
		artifacts.AnswerPresentationRaw = filepath.Join(outDir, "answer.presentation.raw.md")
		artifacts.ReportTemplateMD = filepath.Join(outDir, "report.template.md")
		artifacts.ReportMD = filepath.Join(outDir, "report.md")
		artifacts.VisualHTML = filepath.Join(outDir, "visual.html")
	}
	return artifacts
}

func QuerySHA256(sql string) string {
	sum := sha256.Sum256([]byte(sql))
	return hex.EncodeToString(sum[:])
}

func WriteManifest(path string, manifest model.RunManifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func ReadManifest(path string) (model.RunManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.RunManifest{}, err
	}
	var manifest model.RunManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return model.RunManifest{}, err
	}
	return manifest, nil
}
