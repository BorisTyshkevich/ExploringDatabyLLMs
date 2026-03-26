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
	"qforge/internal/runs"
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
				Artifacts: ArtifactLinks{
					ReviewMD: ArtifactRef{URL: "https://example.invalid/run-001/review.md"},
				},
				PresentationTarget: "react",
				ReviewVerdict:      "PASS",
				Status:             model.RunStatusOK,
				ResultRowCount:     832,
				SQLGenMS:           2400,
				VisualGenMS:        5100,
				VisualBuildMS:      1300,
				Metrics: &RunMetrics{
					QueryDurationMS: 900,
					ReadRows:        1146680615,
					ReadBytes:       3221225472,
					MemoryUsage:     283105876,
				},
			},
			{
				QuestionID:    "q003",
				QuestionTitle: "Delta ATL departure delay hotspots by destination and time block",
				Runner:        "gemini",
				Model:         "gemini-2.5-pro",
				Artifacts: ArtifactLinks{
					ReviewMD: ArtifactRef{URL: "https://example.invalid/run-002/review.md"},
				},
				ReviewVerdict:  "FAIL",
				Status:         model.RunStatusPartial,
				ResultRowCount: 0,
				Warnings:       []string{"gemini/gemini-2.5-pro: query_log metrics not found"},
			},
		},
	}

	got := renderMarkdown(report)
	for _, want := range []string{
		"## q003: Delta ATL departure delay hotspots by destination and time block",
		"- Status: 1 run(s) did not finish cleanly: gemini/gemini-2.5-pro.",
		"- Result rows (manifest): mismatch (0, 832).",
		"- Fastest successful run: codex/gpt-5.4 at 900 ms.",
		"| runner | model | run | target | review | review md | status | result rows (manifest) | sql gen | visual gen | build | query time | read rows | bytes read | peak memory | warnings |",
		"| codex | gpt-5.4 | n/a | react | PASS | [review.md](https://example.invalid/run-001/review.md) | ok | 832 | 2.40 s | 5.10 s | 1.30 s | 900 ms | 1,146,680,615 | 3.0 GiB | 270.0 MiB | 0 |",
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
				RunDir:                "/tmp/run-001",
				QuestionID:            "q004",
				QuestionTitle:         "Worst origin airports by departure on-time performance",
				Runner:                "claude",
				Model:                 "opus",
				Status:                model.RunStatusOK,
				StartedAt:             time.Unix(0, 0).UTC(),
				FinishedAt:            time.Unix(1, 0).UTC(),
				ResultRowCount:        25,
				VisualArtifactPresent: true,
				Columns:               []string{"OriginCode", "DepartureOtpPct"},
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
	if !strings.Contains(string(data), "\"columns\"") || !strings.Contains(string(data), "\"result_row_count\"") {
		t.Fatalf("expected summary fields in compare json, got: %s", string(data))
	}
	if !strings.Contains(string(data), "\"visual_artifact_present\"") {
		t.Fatalf("expected visual artifact presence in compare json, got: %s", string(data))
	}
}

func TestSummarizeRunWarnsWhenVisualExistsButManifestSaysSkipped(t *testing.T) {
	runsRoot := t.TempDir()
	codeRoot := t.TempDir()
	runDir := filepath.Join(runsRoot, "2026-03-26", "q001_hops_per_day", "codex", "gpt-5.4", "run-001")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir run dir: %v", err)
	}
	manifest := model.RunManifest{
		QuestionID:    "q001",
		QuestionSlug:  "q001_hops_per_day",
		QuestionTitle: "Hops per day",
		Dataset:       "ontime",
		Runner:        "codex",
		Model:         "gpt-5.4",
		Status:        model.RunStatusOK,
		StartedAt:     time.Unix(0, 0).UTC(),
		FinishedAt:    time.Unix(1, 0).UTC(),
		DurationSec:   1,
		Phases: model.RunPhases{
			PresentationGeneration: model.PhaseStatusSkipped,
			PresentationRender:     model.PhaseStatusSkipped,
		},
		Artifacts: model.ArtifactPaths{
			ManifestJSON: filepath.Join(runDir, "manifest.json"),
		},
	}
	if err := runs.WriteManifest(filepath.Join(runDir, "manifest.json"), manifest); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "visual.html"), []byte("<html></html>"), 0o644); err != nil {
		t.Fatalf("write visual.html: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "result.json"), []byte(`{"columns":["c"],"row_count":1}`), 0o644); err != nil {
		t.Fatalf("write result.json: %v", err)
	}

	got, warnings, err := summarizeRun(t.Context(), codeRoot, runsRoot, runDir, "", "")
	if err != nil {
		t.Fatalf("summarizeRun returned error: %v", err)
	}
	if !got.VisualArtifactPresent {
		t.Fatalf("expected visual artifact to be present")
	}
	want := "codex/gpt-5.4/run-001: visual.html exists even though manifest presentation phases are marked skipped"
	if strings.Join(warnings, "\n") == "" || !strings.Contains(strings.Join(warnings, "\n"), want) {
		t.Fatalf("expected warning %q, got %v", want, warnings)
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
		"PROMPT_REPORT:",
		"{{prompt_report_paths_md}}",
		"PROMPT_VISUAL:",
		"{{prompt_visual_paths_md}}",
		"SQL:",
		"{{query_sql_paths_md}}",
		"REPORT:",
		"{{report_md_paths_md}}",
		"REVIEW:",
		"{{review_md_paths_md}}",
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
		for _, name := range []string{"prompt.report.md", "query.sql", "report.md", "review.md", "visual.html", "result.json"} {
			if err := os.WriteFile(filepath.Join(runDir, name), []byte("x"), 0o644); err != nil {
				t.Fatalf("write artifact %s: %v", name, err)
			}
		}
		if runDir == report.Runs[0].RunDir {
			if err := os.WriteFile(filepath.Join(runDir, "prompt.visual.md"), []byte("x"), 0o644); err != nil {
				t.Fatalf("write prompt.visual.md: %v", err)
			}
		}
		if err := os.MkdirAll(filepath.Join(runDir, "visual_src"), 0o755); err != nil {
			t.Fatalf("mkdir visual_src: %v", err)
		}
		if err := os.MkdirAll(filepath.Join(runDir, "visual_build"), 0o755); err != nil {
			t.Fatalf("mkdir visual_build: %v", err)
		}
	}
	report.Runs[0].Artifacts = buildRunArtifactLinks(runsRoot, report.Runs[0].RunDir)
	report.Runs[1].Artifacts = buildRunArtifactLinks(runsRoot, report.Runs[1].RunDir)

	got, err := BuildAnalysisPrompt(codeRoot, runsRoot, question, report, filepath.Join(runsRoot, "2026-03-16", "q003_delta_atl_departure_delay_hotspots", "compare", "compare.json"))
	if err != nil {
		t.Fatalf("BuildAnalysisPrompt returned error: %v", err)
	}

	for _, want := range []string{
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/prompt.report.md",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/prompt.report.md",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/prompt.visual.md",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/report.md",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/review.md",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual.html",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/report.md",
		"2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/visual.html",
		"https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-16%2Fq003_delta_atl_departure_delay_hotspots%2Fclaude%2Fopus%2Frun-001%2Freport.md",
		"https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-16%2Fq003_delta_atl_departure_delay_hotspots%2Fclaude%2Fopus%2Frun-001%2Freview.md",
		"https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-16%2Fq003_delta_atl_departure_delay_hotspots%2Fclaude%2Fopus%2Frun-001%2Fprompt.report.md",
		"https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-16%2Fq003_delta_atl_departure_delay_hotspots%2Fclaude%2Fopus%2Frun-001%2Fprompt.visual.md",
		"https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql",
		"https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual.html",
		"https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/tree/main/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual_src",
		"https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/tree/main/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual_build",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected prompt to contain %q, got:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{
		"prompts/q003_delta_atl_departure_delay_hotspots/report_prompt.md",
		"prompts/q003_delta_atl_departure_delay_hotspots/visual_prompt.md",
		"## q003: Delta ATL departure delay hotspots",
		"Fastest successful run:",
		"gemini/gemini-3.1-pro-preview/run-001/prompt.visual.md",
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("expected prompt not to contain %q, got:\n%s", unwanted, got)
		}
	}
}

func TestBuildAnalysisPromptIncludesMultiQuerySQLArtifacts(t *testing.T) {
	codeRoot := t.TempDir()
	runsRoot := t.TempDir()
	promptDir := filepath.Join(codeRoot, "prompts")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}
	template := "SQL:\n{{query_sql_paths_md}}\nPUBLISHED:\n{{published_run_artifacts_md}}\n"
	if err := os.WriteFile(filepath.Join(promptDir, analysisPromptFile), []byte(template), 0o644); err != nil {
		t.Fatalf("write analysis prompt: %v", err)
	}

	questionDir := filepath.Join(codeRoot, "prompts", "q006_peak_aa_delay_month_network")
	if err := os.MkdirAll(questionDir, 0o755); err != nil {
		t.Fatalf("mkdir question dir: %v", err)
	}
	question := model.Question{
		Dir: questionDir,
		Meta: model.QuestionMeta{
			ID:    "q006",
			Slug:  "q006_peak_aa_delay_month_network",
			Title: "American Airlines peak network delay month and contributors",
		},
	}
	runDir := filepath.Join(runsRoot, "2026-03-24", "q006_peak_aa_delay_month_network", "codex", "gpt-5.4", "run-002")
	if err := os.MkdirAll(filepath.Join(runDir, "queries"), 0o755); err != nil {
		t.Fatalf("mkdir queries dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "queries", "main.sql"), []byte("SELECT peak_month"), 0o644); err != nil {
		t.Fatalf("write queries/main.sql: %v", err)
	}
	for _, name := range []string{"q1.sql", "q2.sql"} {
		if err := os.WriteFile(filepath.Join(runDir, "queries", name), []byte("SELECT 1"), 0o644); err != nil {
			t.Fatalf("write query artifact %s: %v", name, err)
		}
	}
	report := Report{
		Day: "2026-03-24",
		Runs: []RunSummary{
			{
				RunDir:    runDir,
				RunID:     "run-002",
				Runner:    "codex",
				Model:     "gpt-5.4",
				Artifacts: buildRunArtifactLinks(runsRoot, runDir),
			},
		},
	}

	got, err := BuildAnalysisPrompt(codeRoot, runsRoot, question, report, filepath.Join(runsRoot, "2026-03-24", "q006_peak_aa_delay_month_network", "compare", "compare.json"))
	if err != nil {
		t.Fatalf("BuildAnalysisPrompt returned error: %v", err)
	}
	for _, want := range []string{
		"2026-03-24/q006_peak_aa_delay_month_network/codex/gpt-5.4/run-002/queries/main.sql",
		"2026-03-24/q006_peak_aa_delay_month_network/codex/gpt-5.4/run-002/queries/q1.sql",
		"2026-03-24/q006_peak_aa_delay_month_network/codex/gpt-5.4/run-002/queries/q2.sql",
		"main.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-24/q006_peak_aa_delay_month_network/codex/gpt-5.4/run-002/queries/main.sql",
		"q1.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-24/q006_peak_aa_delay_month_network/codex/gpt-5.4/run-002/queries/q1.sql",
		"q2.sql: https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-24/q006_peak_aa_delay_month_network/codex/gpt-5.4/run-002/queries/q2.sql",
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
