package cli

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"qforge/internal/model"
	"qforge/internal/runs"
)

func TestResolveRunAnalysisModeAllowsTemplateManualOverride(t *testing.T) {
	got, err := resolveRunAnalysisMode("template_files", "manual_templates")
	if err != nil {
		t.Fatalf("resolveRunAnalysisMode returned error: %v", err)
	}
	if got != model.AnalysisModeManualTemplate {
		t.Fatalf("unexpected override result: %q", got)
	}

	got, err = resolveRunAnalysisMode("manual_templates", "template_files")
	if err != nil {
		t.Fatalf("resolveRunAnalysisMode returned error: %v", err)
	}
	if got != model.AnalysisModeTemplateFiles {
		t.Fatalf("unexpected override result: %q", got)
	}
}

func TestResolveRunAnalysisModeRejectsJSONCrossover(t *testing.T) {
	if _, err := resolveRunAnalysisMode("json_artifact", "template_files"); err == nil {
		t.Fatalf("expected json_artifact crossover to fail")
	}
	if _, err := resolveRunAnalysisMode("manual_templates", "json_artifact"); err == nil {
		t.Fatalf("expected manual_templates -> json_artifact to fail")
	}
}

func TestLoadTemplateAnalysisArtifact(t *testing.T) {
	dir := t.TempDir()
	artifacts := runs.DefaultArtifacts(dir, true)
	if err := os.WriteFile(artifacts.QuerySQL, []byte("SELECT 1"), 0o644); err != nil {
		t.Fatalf("write query.sql: %v", err)
	}
	if err := os.WriteFile(artifacts.ReportTemplateMD, []byte("# Title\n\n{{data_overview_md}}"), 0o644); err != nil {
		t.Fatalf("write report.template.md: %v", err)
	}
	got, err := loadTemplateAnalysisArtifact(artifacts)
	if err != nil {
		t.Fatalf("loadTemplateAnalysisArtifact returned error: %v", err)
	}
	if got.SQL != "SELECT 1" || got.ReportMarkdown == "" {
		t.Fatalf("unexpected artifact: %+v", got)
	}
}

func TestExecuteRunManualTemplatesStagesOnly(t *testing.T) {
	repoRoot := t.TempDir()
	writeTestQuestionRepo(t, repoRoot, "manual_templates")
	t.Setenv("QFORGE_CODE_ROOT", repoRoot)
	t.Setenv("QFORGE_RUN_ROOT", repoRoot)

	err := executeRun(context.Background(), runOptions{
		QuestionRef: "q901",
		Runner:      "claude",
		Model:       "opus",
		CLIBin:      "/path/that/should/not/run",
	})
	if err != nil {
		t.Fatalf("executeRun returned error: %v", err)
	}

	runDir := latestRunDir(t, repoRoot)
	if _, err := os.Stat(filepath.Join(runDir, "prompt.report.md")); err != nil {
		t.Fatalf("expected staged prompt.report.md: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runDir, "query.sql")); !os.IsNotExist(err) {
		t.Fatalf("did not expect query.sql in manual staging run, err=%v", err)
	}
	manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if manifest.AnalysisMode != "manual_templates" {
		t.Fatalf("unexpected manifest analysis mode: %q", manifest.AnalysisMode)
	}
}

func TestExecuteRunTemplateFilesOverrideInvokesProvider(t *testing.T) {
	repoRoot := t.TempDir()
	server := newExecuteQueryServer()
	defer server.Close()
	writeTestQuestionRepo(t, repoRoot, "manual_templates")
	writeFakeTemplateProvider(t, repoRoot)
	t.Setenv("QFORGE_CODE_ROOT", repoRoot)
	t.Setenv("QFORGE_RUN_ROOT", repoRoot)

	err := executeRun(context.Background(), runOptions{
		QuestionRef:          "q901",
		Runner:               "claude",
		Model:                "opus",
		AnalysisModeOverride: "template_files",
		CLIBin:               filepath.Join(repoRoot, "fake-provider.sh"),
	})
	if err != nil {
		t.Fatalf("executeRun returned error: %v", err)
	}

	runDir := latestRunDir(t, repoRoot)
	if _, err := os.Stat(filepath.Join(runDir, "query.sql")); err != nil {
		t.Fatalf("expected query.sql written by provider: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runDir, "report.template.md")); err != nil {
		t.Fatalf("expected report.template.md written by provider: %v", err)
	}
	manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if manifest.AnalysisMode != "template_files" {
		t.Fatalf("unexpected manifest analysis mode: %q", manifest.AnalysisMode)
	}
	if manifest.Status != model.RunStatusOK {
		t.Fatalf("expected successful automated template run, got %q", manifest.Status)
	}
	_ = server
}

func TestExecuteRunTemplateFilesOverrideToManualStagesOnly(t *testing.T) {
	repoRoot := t.TempDir()
	writeTestQuestionRepo(t, repoRoot, "template_files")
	t.Setenv("QFORGE_CODE_ROOT", repoRoot)
	t.Setenv("QFORGE_RUN_ROOT", repoRoot)

	err := executeRun(context.Background(), runOptions{
		QuestionRef:          "q901",
		Runner:               "claude",
		Model:                "opus",
		AnalysisModeOverride: "manual_templates",
		CLIBin:               "/path/that/should/not/run",
	})
	if err != nil {
		t.Fatalf("executeRun returned error: %v", err)
	}

	runDir := latestRunDir(t, repoRoot)
	if _, err := os.Stat(filepath.Join(runDir, "query.sql")); !os.IsNotExist(err) {
		t.Fatalf("did not expect query.sql in manual override run, err=%v", err)
	}
	manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if manifest.AnalysisMode != "manual_templates" {
		t.Fatalf("unexpected manifest analysis mode: %q", manifest.AnalysisMode)
	}
}

func TestExecuteRunStagesPresentationPromptWithoutWithVisual(t *testing.T) {
	repoRoot := t.TempDir()
	server := newExecuteQueryServer()
	defer server.Close()
	writeTestQuestionRepoWithArtifacts(t, repoRoot, "template_files", "report.md,visual.html")
	writeFakeTemplateProvider(t, repoRoot)
	t.Setenv("QFORGE_CODE_ROOT", repoRoot)
	t.Setenv("QFORGE_RUN_ROOT", repoRoot)

	err := executeRun(context.Background(), runOptions{
		QuestionRef: "q901",
		Runner:      "claude",
		Model:       "opus",
		CLIBin:      filepath.Join(repoRoot, "fake-provider.sh"),
	})
	if err != nil {
		t.Fatalf("executeRun returned error: %v", err)
	}

	runDir := latestRunDir(t, repoRoot)
	if _, err := os.Stat(filepath.Join(runDir, "prompt.visual.md")); err != nil {
		t.Fatalf("expected staged prompt.visual.md: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runDir, "visual.html")); !os.IsNotExist(err) {
		t.Fatalf("did not expect visual.html without --with-visual, err=%v", err)
	}
	manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if manifest.Phases.PresentationGeneration != model.PhaseStatusSkipped || manifest.Phases.PresentationRender != model.PhaseStatusSkipped {
		t.Fatalf("expected deferred presentation phases, got %+v", manifest.Phases)
	}
	_ = server
}

func TestExecuteRunManualTemplatesStagesVisualPromptForVisualQuestions(t *testing.T) {
	repoRoot := t.TempDir()
	writeTestQuestionRepoWithArtifacts(t, repoRoot, "manual_templates", "report.md,visual.html")
	t.Setenv("QFORGE_CODE_ROOT", repoRoot)
	t.Setenv("QFORGE_RUN_ROOT", repoRoot)

	err := executeRun(context.Background(), runOptions{
		QuestionRef: "q901",
		Runner:      "claude",
		Model:       "opus",
		CLIBin:      "/path/that/should/not/run",
	})
	if err != nil {
		t.Fatalf("executeRun returned error: %v", err)
	}

	runDir := latestRunDir(t, repoRoot)
	if _, err := os.Stat(filepath.Join(runDir, "prompt.visual.md")); err != nil {
		t.Fatalf("expected staged prompt.visual.md: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runDir, "visual.html")); !os.IsNotExist(err) {
		t.Fatalf("did not expect visual.html in manual staging run, err=%v", err)
	}
}

func TestProcessPresentationTemplateFiles(t *testing.T) {
	repoRoot := t.TempDir()
	server := newExecuteQueryServer()
	defer server.Close()
	writeTestQuestionRepo(t, repoRoot, "template_files")
	t.Setenv("QFORGE_CODE_ROOT", repoRoot)
	t.Setenv("QFORGE_RUN_ROOT", repoRoot)

	runDir := filepath.Join(repoRoot, "2026-03-21", "q901_test_question", "claude", "opus", "run-001")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("mkdir runDir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "query.sql"), []byte("SELECT 1"), 0o644); err != nil {
		t.Fatalf("write query.sql: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "report.template.md"), []byte("# Report\n\n{{data_overview_md}}"), 0o644); err != nil {
		t.Fatalf("write report.template.md: %v", err)
	}

	if err := processPresentation(context.Background(), processPresentationOptions{RunDir: runDir}); err != nil {
		t.Fatalf("processPresentation returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runDir, "report.md")); err != nil {
		t.Fatalf("expected report.md: %v", err)
	}
	if _, err := os.Stat(filepath.Join(runDir, "result.json")); err != nil {
		t.Fatalf("expected result.json: %v", err)
	}
	manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if manifest.AnalysisMode != "template_files" {
		t.Fatalf("unexpected manifest analysis mode: %q", manifest.AnalysisMode)
	}
	_ = server
}

func writeTestQuestionRepo(t *testing.T, repoRoot, analysisMode string) {
	t.Helper()
	writeTestQuestionRepoWithArtifacts(t, repoRoot, analysisMode, "report.md")
}

func writeTestQuestionRepoWithArtifacts(t *testing.T, repoRoot, analysisMode, artifactsRequired string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(repoRoot, "prompts", "q901_test_question"), 0o755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repoRoot, "datasets", "ontime"), 0o755); err != nil {
		t.Fatalf("mkdir datasets: %v", err)
	}
	meta := "id: q901\nslug: q901_test_question\ntitle: Test Question\ndataset: ontime\nanalysis_mode: " + analysisMode + "\nartifacts_required: " + artifactsRequired + "\n"
	if err := os.WriteFile(filepath.Join(repoRoot, "prompts", "q901_test_question", "meta.yaml"), []byte(meta), 0o644); err != nil {
		t.Fatalf("write meta.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "prompts", "q901_test_question", "report_prompt.md"), []byte("Return a simple query."), 0o644); err != nil {
		t.Fatalf("write report_prompt.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "datasets", "ontime", "mcp.yaml"), []byte("dataset: ontime\nmcp_url: "+newExecuteQueryServerURL(t)+"\ndefault_mcp_server_name: demo\n"), 0o644); err != nil {
		t.Fatalf("write mcp.yaml: %v", err)
	}
	for _, name := range []string{"common.md", "common_report.md", "common_report_templates.md", "common_visual.md", "common_visual_dynamic.md", "common_visual_static.md"} {
		src := filepath.Join("/Users/bvt/work/ExploringDatabyLLMs", "prompts", name)
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("read shared prompt %s: %v", name, err)
		}
		if err := os.WriteFile(filepath.Join(repoRoot, "prompts", name), data, 0o644); err != nil {
			t.Fatalf("write shared prompt %s: %v", name, err)
		}
	}
}

func writeFakeTemplateProvider(t *testing.T, repoRoot string) {
	t.Helper()
	script := "#!/usr/bin/env bash\nset -euo pipefail\ncat >/dev/null\nprintf 'SELECT 1\\n' > query.sql\nprintf '# Report\\n\\n{{data_overview_md}}\\n' > report.template.md\nprintf 'provider wrote template files\\n'\n"
	if err := os.WriteFile(filepath.Join(repoRoot, "fake-provider.sh"), []byte(script), 0o755); err != nil {
		t.Fatalf("write fake provider: %v", err)
	}
}

func latestRunDir(t *testing.T, repoRoot string) string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(repoRoot, "*", "q901_test_question", "claude", "opus", "run-*"))
	if err != nil {
		t.Fatalf("glob run dir: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one run dir, got %v", matches)
	}
	return matches[0]
}

func newExecuteQueryServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/openapi/execute_query") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"columns":["value"],"rows":[[1]]}`))
	}))
}

func newExecuteQueryServerURL(t *testing.T) string {
	t.Helper()
	server := newExecuteQueryServer()
	t.Cleanup(server.Close)
	return server.URL + "/token/http"
}
