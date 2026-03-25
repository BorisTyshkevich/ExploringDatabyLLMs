package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"qforge/internal/model"
	verbosepkg "qforge/internal/verbose"
)

func TestVerbosePrefixFormat(t *testing.T) {
	got := verbosepkg.PrefixAt(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.FixedZone("UTC+1", 3600)), "opus")
	want := "2026-01-01 00:00:00 opus"
	if got != want {
		t.Fatalf("unexpected prefix: got %q want %q", got, want)
	}
}

func TestModelLabelForRunnersUsesResolvedModels(t *testing.T) {
	got, err := modelLabelForRunners([]string{"codex", "claude", "gemini"}, []string{"", "sonnet"})
	if err != nil {
		t.Fatalf("modelLabelForRunners returned error: %v", err)
	}
	want := "gpt-5.4,sonnet,gemini-3.1-pro-preview"
	if got != want {
		t.Fatalf("unexpected label: got %q want %q", got, want)
	}
}

func TestDefaultModelForRunnerUsesSonnetForClaude(t *testing.T) {
	got, err := defaultModelForRunner("claude")
	if err != nil {
		t.Fatalf("defaultModelForRunner returned error: %v", err)
	}
	if got != "sonnet" {
		t.Fatalf("unexpected claude default model: got %q want %q", got, "sonnet")
	}
}

func TestModelLabelForRunnersHandlesDuplicateRunnersWithDistinctModels(t *testing.T) {
	got, err := modelLabelForRunners([]string{"claude", "claude", "codex"}, []string{"opus", "sonnet", "gpt-5.4"})
	if err != nil {
		t.Fatalf("modelLabelForRunners returned error: %v", err)
	}
	want := "opus,sonnet,gpt-5.4"
	if got != want {
		t.Fatalf("unexpected label: got %q want %q", got, want)
	}
}

func TestPrintRootUsageMentionsProcessPresentation(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	printRootUsage(w)
	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read usage: %v", err)
	}
	if !strings.Contains(buf.String(), "process-presentation") {
		t.Fatalf("expected root usage to mention process-presentation, got: %s", buf.String())
	}
}

func TestMarkPresentationDeferredSkipsUnsetPhases(t *testing.T) {
	got := markPresentationDeferred(model.RunPhases{
		SQLGeneration: model.PhaseStatusOK,
		SQLExecution:  model.PhaseStatusOK,
	})
	if got.PresentationGeneration != model.PhaseStatusSkipped {
		t.Fatalf("expected presentation_generation skipped, got %q", got.PresentationGeneration)
	}
	if got.PresentationRender != model.PhaseStatusSkipped {
		t.Fatalf("expected presentation_render skipped, got %q", got.PresentationRender)
	}
}

func TestResolveRequestedAnalysisModeUsesManualAlias(t *testing.T) {
	got := resolveRequestedAnalysisMode("", true)
	if got != "manual_templates" {
		t.Fatalf("expected manual alias to force manual_templates, got %q", got)
	}
}

func TestResolveRequestedAnalysisModeManualAliasOverridesExplicitValue(t *testing.T) {
	got := resolveRequestedAnalysisMode("template_files", true)
	if got != "manual_templates" {
		t.Fatalf("expected manual alias to override explicit mode, got %q", got)
	}
}

func TestResolveRequestedAnalysisModeKeepsExplicitValueWithoutManualAlias(t *testing.T) {
	got := resolveRequestedAnalysisMode("template_files", false)
	if got != "template_files" {
		t.Fatalf("expected explicit mode to be preserved, got %q", got)
	}
}

func TestApplyPresentationTargetOverride(t *testing.T) {
	question := model.Question{Meta: model.QuestionMeta{PresentationTarget: "html"}}
	if err := applyPresentationTargetOverride(&question, "react"); err != nil {
		t.Fatalf("applyPresentationTargetOverride returned error: %v", err)
	}
	if question.Meta.PresentationTarget != "react" {
		t.Fatalf("expected presentation target override to apply, got %q", question.Meta.PresentationTarget)
	}
	if err := applyPresentationTargetOverride(&question, "bad"); err == nil {
		t.Fatalf("expected unsupported presentation target override to fail")
	}
}

func TestParseReviewVerdict(t *testing.T) {
	got, err := parseReviewVerdict("# Analysis Review\nVerdict: PASS\n")
	if err != nil {
		t.Fatalf("parseReviewVerdict returned error: %v", err)
	}
	if got != "PASS" {
		t.Fatalf("unexpected verdict: %q", got)
	}
	got, err = parseReviewVerdict("# Analysis Review\nVerdict: WARN\n")
	if err != nil {
		t.Fatalf("parseReviewVerdict returned error for WARN: %v", err)
	}
	if got != "WARN" {
		t.Fatalf("unexpected WARN verdict: %q", got)
	}
	if _, err := parseReviewVerdict("# Analysis Review\nVerdict: MAYBE\n"); err == nil {
		t.Fatalf("expected unsupported verdict to fail")
	}
	if _, err := parseReviewVerdict("# Analysis Review\n## Summary\n"); err == nil {
		t.Fatalf("expected missing verdict to fail")
	}
}

func TestReviewVerdictHelpers(t *testing.T) {
	if !reviewVerdictBlocksRun("FAIL") {
		t.Fatalf("expected FAIL to block run")
	}
	if reviewVerdictBlocksRun("WARN") {
		t.Fatalf("did not expect WARN to block run")
	}
	if got := successfulRunStatus("PASS"); got != model.RunStatusOK {
		t.Fatalf("expected PASS to keep ok status, got %q", got)
	}
	if got := successfulRunStatus("WARN"); got != model.RunStatusPartial {
		t.Fatalf("expected WARN to downgrade to partial, got %q", got)
	}
}

func TestPresentationPhasesOKRejectsFailures(t *testing.T) {
	if presentationPhasesOK(model.RunPhases{
		PresentationGeneration: model.PhaseStatusSkipped,
		PresentationRender:     model.PhaseStatusFailed,
	}) {
		t.Fatalf("expected failed presentation phase to make status not ok")
	}
}

func TestReadOrInferRunManifestWithoutManifest(t *testing.T) {
	runDir := "/tmp/2026-03-20/q001_hops_per_day/claude/opus/run-002"
	manifest, question, err := readOrInferRunManifest("/Users/bvt/work/ExploringDatabyLLMs", runDir)
	if err != nil {
		t.Fatalf("expected manifest inference to succeed: %v", err)
	}
	if question.Meta.ID != "q001" {
		t.Fatalf("expected q001 question, got %q", question.Meta.ID)
	}
	if manifest.Runner != "claude" || manifest.Model != "opus" {
		t.Fatalf("unexpected manifest runner/model: %+v", manifest)
	}
	if manifest.QuestionSlug != "q001_hops_per_day" {
		t.Fatalf("unexpected manifest slug: %+v", manifest)
	}
	if manifest.Metadata["manifest_inferred"] != "true" {
		t.Fatalf("expected inferred manifest metadata, got %+v", manifest.Metadata)
	}
}

func TestLoadVisualArtifactRejectsStaleFallbackFile(t *testing.T) {
	tmpDir := t.TempDir()
	htmlPath := filepath.Join(tmpDir, "visual.html")
	if err := os.WriteFile(htmlPath, []byte("<html>old</html>"), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	notBefore := time.Now().Add(2 * time.Second)
	if _, err := loadVisualArtifact("", tmpDir, notBefore); err == nil {
		t.Fatalf("expected stale fallback files to be rejected")
	}
}

func TestLoadVisualArtifactAcceptsFreshFallbackFile(t *testing.T) {
	tmpDir := t.TempDir()
	htmlPath := filepath.Join(tmpDir, "visual.html")
	notBefore := time.Now()
	time.Sleep(20 * time.Millisecond)

	html := "<!doctype html>\n<html><body>fresh</body></html>\n"
	if err := os.WriteFile(htmlPath, []byte(html), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}

	gotHTML, err := loadVisualArtifact("", tmpDir, notBefore)
	if err != nil {
		t.Fatalf("expected fresh fallback files to be accepted: %v", err)
	}
	if gotHTML != "<!doctype html>\n<html><body>fresh</body></html>" {
		t.Fatalf("unexpected html content: %q", gotHTML)
	}
}

func TestValidateReactSourceArtifactsRejectsMissingBuildScript(t *testing.T) {
	sourceDir := filepath.Join(t.TempDir(), "visual_src")
	if err := os.MkdirAll(filepath.Join(sourceDir, "src"), 0o755); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}
	files := map[string]string{
		"package.json": `{"name":"fixture","scripts":{}}`,
		"index.html":   "<!doctype html><html><body><div id=\"root\"></div></body></html>",
		"src/main.jsx": "console.log('main')",
		"src/App.jsx":  "export default function App(){ return null }",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(sourceDir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	if err := validateReactSourceArtifacts(sourceDir); err == nil {
		t.Fatalf("expected missing build script to fail source validation")
	}
}

func TestMaterializePresentationArtifactBuildsReactOutput(t *testing.T) {
	runDir := t.TempDir()
	artifacts := model.ArtifactPaths{
		VisualSourceDir:   filepath.Join(runDir, "visual_src"),
		VisualBuildDir:    filepath.Join(runDir, "visual_build"),
		VisualAssetsDir:   filepath.Join(runDir, "visual_assets"),
		VisualPackageJSON: filepath.Join(runDir, "visual_src", "package.json"),
		VisualHTML:        filepath.Join(runDir, "visual.html"),
	}
	if err := os.MkdirAll(filepath.Join(artifacts.VisualSourceDir, "src"), 0o755); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}
	files := map[string]string{
		"package.json": `{"name":"fixture","private":true,"scripts":{"build":"node build.mjs"}}`,
		"index.html":   "<!doctype html><html><body><div id=\"root\"></div></body></html>",
		"build.mjs": `import { mkdirSync, writeFileSync } from 'node:fs';
mkdirSync('../visual_build/visual_assets', { recursive: true });
writeFileSync('../visual_build/index.html', '<!doctype html><html><head><script type="module" src="./visual_assets/app.js"></script></head><body><footer><input type="password"/><textarea>SELECT 1</textarea><button>Fetch</button><div id="query-ledger"></div><div data-status>ok</div></footer></body></html>');
writeFileSync('../visual_build/visual_assets/app.js', 'console.log("ok")');
`,
		"src/main.jsx": "console.log('main')",
		"src/App.jsx":  "export default function App(){ return null }",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(artifacts.VisualSourceDir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	notBefore := time.Now().Add(-1 * time.Second)
	question := model.Question{Meta: model.QuestionMeta{PresentationTarget: "react", VisualMode: "dynamic"}}
	got, err := materializePresentationArtifact(runDir, question, artifacts, "", notBefore, "test-model", false)
	if err != nil {
		t.Fatalf("materializePresentationArtifact returned error: %v", err)
	}
	if got.BuildDurationMS <= 0 {
		t.Fatalf("expected build duration to be recorded, got %+v", got)
	}
	if !strings.Contains(got.HTML, "./visual_assets/app.js") {
		t.Fatalf("expected built html to reference visual_assets, got: %s", got.HTML)
	}
	if _, err := os.Stat(artifacts.VisualHTML); err != nil {
		t.Fatalf("expected visual.html to be written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(artifacts.VisualAssetsDir, "app.js")); err != nil {
		t.Fatalf("expected built visual asset to be copied: %v", err)
	}
	if got.Metadata["react_build"] != "ok" || got.Metadata["react_source_validation"] != "ok" {
		t.Fatalf("expected react metadata to record successful build, got %+v", got.Metadata)
	}
}

func TestLoadCompareReportArtifactPrefersFreshFile(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "compare_report.md")
	notBefore := time.Now()
	time.Sleep(20 * time.Millisecond)

	report := "# Long Report\n\nFull compare body.\n"
	if err := os.WriteFile(reportPath, []byte(report), 0o644); err != nil {
		t.Fatalf("write compare_report.md: %v", err)
	}

	got, err := loadCompareReportArtifact("Short stdout note", tmpDir, notBefore)
	if err != nil {
		t.Fatalf("expected fresh compare_report.md to be accepted: %v", err)
	}
	if got != "# Long Report\n\nFull compare body." {
		t.Fatalf("unexpected compare report content: %q", got)
	}
}

func TestLoadCompareReportArtifactFallsBackToStdout(t *testing.T) {
	tmpDir := t.TempDir()
	got, err := loadCompareReportArtifact("Short stdout note", tmpDir, time.Now())
	if err != nil {
		t.Fatalf("expected stdout fallback to work: %v", err)
	}
	if got != "Short stdout note" {
		t.Fatalf("unexpected stdout fallback content: %q", got)
	}
}

func TestBuildVisualInputSummaryCapturesShapeHints(t *testing.T) {
	question := model.Question{
		Meta: model.QuestionMeta{Title: "Q001", VisualMode: "dynamic"},
	}
	result := model.CanonicalResult{
		Columns:  []string{"Date", "DepTimes", "Route"},
		RowCount: 2,
		Rows: []map[string]any{
			{"Date": "2024-12-01T00:00:00Z", "DepTimes": []any{543.0, 810.0}, "Route": "ISP-BWI-SEA"},
			{"Date": "2024-02-18T00:00:00Z", "DepTimes": []any{621.0, 801.0}, "Route": "CLE-BNA-DEN"},
		},
	}

	got := buildVisualInputSummary(question, result)
	if got.QuestionTitle != "Q001" || got.RowCount != 2 {
		t.Fatalf("unexpected summary header: %+v", got)
	}
	if len(got.SampleRows) != 2 {
		t.Fatalf("expected two sample rows, got %+v", got.SampleRows)
	}
	if got.FieldShapeNotes["Date"] != "ISO-like timestamp string" {
		t.Fatalf("expected timestamp shape note, got %+v", got.FieldShapeNotes)
	}
	if got.FieldShapeNotes["DepTimes"] != "array field" {
		t.Fatalf("expected array shape note, got %+v", got.FieldShapeNotes)
	}
	if got.ModeHint == "" {
		t.Fatalf("expected mode hint in summary: %+v", got)
	}
}

func TestEnsureVisualInputSummaryBackfillsMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "visual_input.json")
	question := model.Question{
		Meta: model.QuestionMeta{Title: "Q001", VisualMode: "dynamic"},
	}
	result := model.CanonicalResult{
		Columns:  []string{"Date", "DepTimes", "Route"},
		RowCount: 1,
		Rows: []map[string]any{
			{"Date": "2024-12-01T00:00:00Z", "DepTimes": []any{543.0, 810.0}, "Route": "ISP-BWI-SEA"},
		},
	}

	got, err := ensureVisualInputSummary(path, question, result)
	if err != nil {
		t.Fatalf("expected missing visual_input.json to be backfilled: %v", err)
	}
	if got.QuestionTitle != "Q001" || got.RowCount != 1 {
		t.Fatalf("unexpected summary: %+v", got)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected visual_input.json to be written: %v", err)
	}
	var persisted model.VisualInputSummary
	if err := json.Unmarshal(bytes, &persisted); err != nil {
		t.Fatalf("expected visual_input.json to be valid json: %v", err)
	}
	if persisted.FieldShapeNotes["DepTimes"] != "array field" {
		t.Fatalf("expected persisted shape notes, got %+v", persisted.FieldShapeNotes)
	}
}

func TestWritePresentationPromptFromSummaryUsesProvidedPrimarySQLForMultiQuery(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "prompt.visual.md")
	question := model.Question{
		Dir: filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta: model.QuestionMeta{
			ID:           "q003",
			Title:        "Delta ATL",
			VisualMode:   "dynamic",
			VisualType:   "html_heatmap",
			AnalysisMode: string(model.AnalysisModeMultiQuery),
		},
		VisualPrompt: "Visual guidance.",
	}
	cfg := model.DatasetConfig{DefaultDatabase: "ontime"}
	visualInput := model.VisualInputSummary{
		QuestionTitle: "Delta ATL",
		QuerySummaries: []model.QueryResultSummary{
			{ID: "worst_hotspot", SQL: "SELECT * FROM supporting"},
		},
	}

	if err := writePresentationPromptFromSummary(path, question, cfg, model.CanonicalResult{}, "SELECT * FROM main_query", "https://mcp.example.invalid/http", "token", visualInput); err != nil {
		t.Fatalf("writePresentationPromptFromSummary returned error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read prompt: %v", err)
	}
	got := string(data)
	if !strings.Contains(got, "SELECT * FROM main_query") {
		t.Fatalf("expected prompt to embed provided primary SQL content, got: %s", got)
	}
	if strings.Contains(got, "SELECT * FROM supporting") && !strings.Contains(got, "\"query_summaries\"") {
		t.Fatalf("did not expect primary SQL to be inferred from supporting query summaries, got: %s", got)
	}
}

func TestRunComparePositionalQuestionRefScopesOutputs(t *testing.T) {
	codeRoot := t.TempDir()
	runRoot := t.TempDir()
	t.Setenv("QFORGE_CODE_ROOT", codeRoot)
	t.Setenv("QFORGE_RUN_ROOT", runRoot)

	mustWriteFile(t, filepath.Join(codeRoot, "prompts", "analysis_prompt.md"), "Return exactly this fenced section:\n\n```markdown\n# Compare Report\n```\n")
	mustWriteFile(t, filepath.Join(codeRoot, "datasets", "ontime", "mcp.yaml"), "dataset: ontime\ndefault_mcp_server_name: demo\nmcp_url: https://example.invalid/http\n")
	for _, item := range []struct {
		dir   string
		id    string
		slug  string
		title string
	}{
		{dir: "q001_hops_per_day", id: "q001", slug: "q001_hops_per_day", title: "Hops Per Day"},
		{dir: "q002_top_carrier", id: "q002", slug: "q002_top_carrier", title: "Top Carrier"},
	} {
		mustWriteFile(t, filepath.Join(codeRoot, "prompts", item.dir, "meta.yaml"), fmt.Sprintf("id: %s\nslug: %s\ntitle: %q\ndataset: ontime\nanalysis_mode: template_files\nartifacts_required: report.md\npresentation_target: html\n", item.id, item.slug, item.title))
		mustWriteFile(t, filepath.Join(codeRoot, "prompts", item.dir, "report_prompt.md"), "# Prompt\n")
	}

	for _, item := range []struct {
		slug string
		id   string
	}{
		{slug: "q001_hops_per_day", id: "q001"},
		{slug: "q002_top_carrier", id: "q002"},
	} {
		runDir := filepath.Join(runRoot, "2026-03-24", item.slug, "claude", "sonnet", "run-001")
		if err := os.MkdirAll(runDir, 0o755); err != nil {
			t.Fatalf("mkdir run dir: %v", err)
		}
		manifest := model.RunManifest{
			SchemaVersion: "4",
			Status:        model.RunStatusOK,
			QuestionID:    item.id,
			QuestionSlug:  item.slug,
			QuestionTitle: item.slug,
			Dataset:       "ontime",
			Runner:        "claude",
			Model:         "sonnet",
			StartedAt:     time.Unix(0, 0).UTC(),
			FinishedAt:    time.Unix(1, 0).UTC(),
			DurationSec:   1,
			Artifacts:     model.ArtifactPaths{},
		}
		data, err := json.Marshal(manifest)
		if err != nil {
			t.Fatalf("marshal manifest: %v", err)
		}
		if err := os.WriteFile(filepath.Join(runDir, "manifest.json"), data, 0o644); err != nil {
			t.Fatalf("write manifest: %v", err)
		}
		if err := os.WriteFile(filepath.Join(runDir, "result.json"), []byte(`{"columns":["x"],"row_count":1}`), 0o644); err != nil {
			t.Fatalf("write result: %v", err)
		}
		if err := os.WriteFile(filepath.Join(runDir, "query.sql"), []byte("SELECT 1\n"), 0o644); err != nil {
			t.Fatalf("write query: %v", err)
		}
		if err := os.WriteFile(filepath.Join(runDir, "report.md"), []byte("# report\n"), 0o644); err != nil {
			t.Fatalf("write report: %v", err)
		}
	}

	fakeClaude := filepath.Join(t.TempDir(), "claude")
	script := "#!/bin/sh\ncat <<'EOF'\n```markdown\n# Compare Report\n\nScoped output.\n```\nEOF\n"
	if err := os.WriteFile(fakeClaude, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake claude: %v", err)
	}

	if err := runCompare(context.Background(), []string{"q001", "--day", "2026-03-24", "--runner", "claude", "--model", "sonnet", "--cli-bin", fakeClaude, "--mcp-url", "https://example.invalid/http"}); err != nil {
		t.Fatalf("runCompare returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(runRoot, "2026-03-24", "q001_hops_per_day", "compare", "compare.json")); err != nil {
		t.Fatalf("expected q001 compare.json to be written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runRoot, "2026-03-24", "q001_hops_per_day", "compare", "analysis.prompt.md")); err != nil {
		t.Fatalf("expected q001 analysis.prompt.md to be written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runRoot, "2026-03-24", "q001_hops_per_day", "compare_report.md")); err != nil {
		t.Fatalf("expected q001 compare_report.md to be written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runRoot, "2026-03-24", "q002_top_carrier", "compare", "compare.json")); !os.IsNotExist(err) {
		t.Fatalf("expected q002 compare artifacts to be untouched, got err=%v", err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
