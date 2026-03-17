package compare

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"qforge/internal/model"
)

func TestRenderMarkdownHighlightsSummaryAndWarnings(t *testing.T) {
	report := Report{
		GeneratedAt: "2026-03-16T12:00:00Z",
		Day:         "2026-03-16",
		Runs: []RunSummary{
			{
				QuestionID:    "q003",
				QuestionTitle: "Delta ATL departure delay hotspots by destination and time block",
				Runner:        "codex",
				Model:         "gpt-5.4",
				Status:        model.RunStatusOK,
				RowCount:      832,
				Metrics: &RunMetrics{
					QueryDurationMS: 900,
					ReadRows:        1146680615,
					MemoryUsage:     283105876,
				},
			},
			{
				QuestionID:    "q003",
				QuestionTitle: "Delta ATL departure delay hotspots by destination and time block",
				Runner:        "gemini",
				Model:         "gemini-2.5-pro",
				Status:        model.RunStatusPartial,
				RowCount:      0,
				Warnings:      []string{"gemini/gemini-2.5-pro: query_log metrics not found"},
			},
		},
	}

	got := renderMarkdown(report)
	for _, want := range []string{
		"## q003: Delta ATL departure delay hotspots by destination and time block",
		"- Status: 1 run(s) did not finish cleanly: gemini/gemini-2.5-pro.",
		"- Row counts: mismatch (0, 832).",
		"- Fastest successful run: codex/gpt-5.4 at 900 ms.",
		"### Warnings",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected markdown to contain %q, got:\n%s", want, got)
		}
	}
}

func TestWriteOutputsWritesCompactJSON(t *testing.T) {
	dir := t.TempDir()
	report := Report{
		GeneratedAt: "2026-03-16T12:00:00Z",
		Day:         "2026-03-16",
		Runs: []RunSummary{
			{
				RunDir:        "/tmp/run-001",
				QuestionID:    "q004",
				QuestionTitle: "Worst origin airports by departure on-time performance",
				Runner:        "claude",
				Model:         "opus",
				Status:        model.RunStatusOK,
				StartedAt:     time.Unix(0, 0).UTC(),
				FinishedAt:    time.Unix(1, 0).UTC(),
				RowCount:      25,
				Columns:       []string{"Origin", "DepartureOtpPct"},
			},
		},
	}

	outDir := filepath.Join(dir, "compare")
	if err := writeOutputs(outDir, report); err != nil {
		t.Fatalf("writeOutputs returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "compare.json"))
	if err != nil {
		t.Fatalf("read compare json: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal compare json: %v", err)
	}
	if strings.Contains(string(data), "\"manifest\"") || strings.Contains(string(data), "\"result\"") {
		t.Fatalf("expected compact compare json, got: %s", string(data))
	}
	if !strings.Contains(string(data), "\"columns\"") || !strings.Contains(string(data), "\"row_count\"") {
		t.Fatalf("expected summary fields in compare json, got: %s", string(data))
	}
}

func TestArtifactPathsForQuestion(t *testing.T) {
	got := ArtifactPathsForQuestion("/repo", "2026-03-16", "q003_delta_atl_departure_delay_hotspots")
	if got.JSON != "/repo/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/compare/compare.json" {
		t.Fatalf("unexpected compare json path: %s", got.JSON)
	}
	if got.ReportMD != "/repo/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/compare_report.md" {
		t.Fatalf("unexpected compare report path: %s", got.ReportMD)
	}
}

func TestBuildAnalysisPromptIncludesPresentationArtifacts(t *testing.T) {
	repoRoot := t.TempDir()
	promptDir := filepath.Join(repoRoot, "prompts")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}
	template := strings.Join([]string{
		"SQL:",
		"{{query_sql_paths_md}}",
		"REPORT:",
		"{{report_md_paths_md}}",
		"VISUAL:",
		"{{visual_html_paths_md}}",
	}, "\n")
	if err := os.WriteFile(filepath.Join(promptDir, analysisPromptFile), []byte(template), 0o644); err != nil {
		t.Fatalf("write analysis prompt: %v", err)
	}

	questionDir := filepath.Join(repoRoot, "prompts", "q003_delta_atl_departure_delay_hotspots")
	if err := os.MkdirAll(questionDir, 0o755); err != nil {
		t.Fatalf("mkdir question dir: %v", err)
	}
	question := model.Question{
		Dir: questionDir,
		Meta: model.QuestionMeta{
			ID:    "q003",
			Slug:  "q003_delta_atl_departure_delay_hotspots",
			Title: "Delta ATL departure delay hotspots",
		},
	}
	report := Report{
		Day: "2026-03-16",
		Runs: []RunSummary{
			{RunDir: filepath.Join(repoRoot, "runs", "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "claude", "opus", "run-001")},
			{RunDir: filepath.Join(repoRoot, "runs", "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "gemini", "gemini-3.1-pro-preview", "run-001")},
		},
	}

	got, err := BuildAnalysisPrompt(repoRoot, question, report, filepath.Join(repoRoot, "runs", "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "compare", "compare.json"))
	if err != nil {
		t.Fatalf("BuildAnalysisPrompt returned error: %v", err)
	}

	for _, want := range []string{
		"runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql",
		"runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/report.md",
		"runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual.html",
		"runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/report.md",
		"runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/visual.html",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected prompt to contain %q, got:\n%s", want, got)
		}
	}
}
