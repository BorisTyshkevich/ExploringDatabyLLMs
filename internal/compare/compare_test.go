package compare

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
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
	got := ArtifactPathsForQuestion("/repo-runs", "2026-03-16", "q003_delta_atl_departure_delay_hotspots")
	if got.JSON != "/repo-runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/compare/compare.json" {
		t.Fatalf("unexpected compare json path: %s", got.JSON)
	}
	if got.ReportMD != "/repo-runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/compare_report.md" {
		t.Fatalf("unexpected compare report path: %s", got.ReportMD)
	}
}

func TestBuildAnalysisPromptIncludesPresentationArtifacts(t *testing.T) {
	codeRoot := t.TempDir()
	runsRoot := t.TempDir()
	promptDir := filepath.Join(codeRoot, "prompts")
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
		"PUBLISHED:",
		"{{published_run_artifacts_md}}",
	}, "\n")
	if err := os.WriteFile(filepath.Join(promptDir, analysisPromptFile), []byte(template), 0o644); err != nil {
		t.Fatalf("write analysis prompt: %v", err)
	}

	questionDir := filepath.Join(codeRoot, "prompts", "q003_delta_atl_departure_delay_hotspots")
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
			{
				RunDir:    filepath.Join(runsRoot, "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "claude", "opus", "run-001"),
				RunID:     "run-001",
				Runner:    "claude",
				Model:     "opus",
				Artifacts: buildRunArtifactLinks(runsRoot, filepath.Join(runsRoot, "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "claude", "opus", "run-001")),
			},
			{
				RunDir:    filepath.Join(runsRoot, "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "gemini", "gemini-3.1-pro-preview", "run-001"),
				RunID:     "run-001",
				Runner:    "gemini",
				Model:     "gemini-3.1-pro-preview",
				Artifacts: buildRunArtifactLinks(runsRoot, filepath.Join(runsRoot, "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "gemini", "gemini-3.1-pro-preview", "run-001")),
			},
		},
	}
	for _, runDir := range []string{
		report.Runs[0].RunDir,
		report.Runs[1].RunDir,
	} {
		if err := os.MkdirAll(runDir, 0o755); err != nil {
			t.Fatalf("mkdir run dir: %v", err)
		}
		for _, name := range []string{"query.sql", "report.md", "visual.html", "result.json"} {
			if err := os.WriteFile(filepath.Join(runDir, name), []byte("x"), 0o644); err != nil {
				t.Fatalf("write artifact %s: %v", name, err)
			}
		}
	}
	report.Runs[0].Artifacts = buildRunArtifactLinks(runsRoot, report.Runs[0].RunDir)
	report.Runs[1].Artifacts = buildRunArtifactLinks(runsRoot, report.Runs[1].RunDir)

	got, err := BuildAnalysisPrompt(codeRoot, runsRoot, question, report, filepath.Join(runsRoot, "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "compare", "compare.json"))
	if err != nil {
		t.Fatalf("BuildAnalysisPrompt returned error: %v", err)
	}

	for _, want := range []string{
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/report.md",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual.html",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/report.md",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/visual.html",
		"https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-16%2Fq003_delta_atl_departure_delay_hotspots%2Fclaude%2Fopus%2Frun-001%2Freport.md",
		"https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql",
		"https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual.html",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected prompt to contain %q, got:\n%s", want, got)
		}
	}
}

func TestPublishedRelativePathAndURLs(t *testing.T) {
	runsRoot := t.TempDir()
	local := filepath.Join(runsRoot, "2026-03-17", "q001_hops_per_day", "codex", "gpt-5.4", "run-003", "report.md")
	if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
		t.Fatalf("mkdir report dir: %v", err)
	}
	if err := os.WriteFile(local, []byte("x"), 0o644); err != nil {
		t.Fatalf("write report: %v", err)
	}
	got := buildArtifactRef(runsRoot, local, "md")
	if got.LocalPath != "2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/report.md" {
		t.Fatalf("unexpected local path: %s", got.LocalPath)
	}
	if got.PublishedPath != "2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/report.md" {
		t.Fatalf("unexpected published path: %s", got.PublishedPath)
	}
	if got.URL != "https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-003%2Freport.md" {
		t.Fatalf("unexpected report URL: %s", got.URL)
	}
}

func TestRunSortingUsesRunnerModelAndRunNumber(t *testing.T) {
	report := Report{
		Runs: []RunSummary{
			{Runner: "codex", Model: "gpt-5.4", RunID: "run-003", RunNumber: 3},
			{Runner: "claude", Model: "opus", RunID: "run-002", RunNumber: 2},
			{Runner: "codex", Model: "gpt-5.4", RunID: "run-001", RunNumber: 1},
			{Runner: "claude", Model: "opus", RunID: "run-001", RunNumber: 1},
		},
	}
	sort.Slice(report.Runs, func(i, j int) bool {
		if report.Runs[i].Runner != report.Runs[j].Runner {
			return report.Runs[i].Runner < report.Runs[j].Runner
		}
		if report.Runs[i].Model != report.Runs[j].Model {
			return report.Runs[i].Model < report.Runs[j].Model
		}
		if report.Runs[i].RunNumber != report.Runs[j].RunNumber {
			return report.Runs[i].RunNumber < report.Runs[j].RunNumber
		}
		return report.Runs[i].RunID < report.Runs[j].RunID
	})
	got := []string{
		report.Runs[0].Runner + "/" + report.Runs[0].Model + "/" + report.Runs[0].RunID,
		report.Runs[1].Runner + "/" + report.Runs[1].Model + "/" + report.Runs[1].RunID,
		report.Runs[2].Runner + "/" + report.Runs[2].Model + "/" + report.Runs[2].RunID,
		report.Runs[3].Runner + "/" + report.Runs[3].Model + "/" + report.Runs[3].RunID,
	}
	want := []string{
		"claude/opus/run-001",
		"claude/opus/run-002",
		"codex/gpt-5.4/run-001",
		"codex/gpt-5.4/run-003",
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("unexpected sort order: got %v want %v", got, want)
	}
}
