package questions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsAnalysisModeToTemplateFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "meta.yaml"), []byte("id: qx\nslug: qx\ntitle: Test\ndataset: ontime\nartifacts_required: report.md\n"), 0o644); err != nil {
		t.Fatalf("write meta.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "report_prompt.md"), []byte("report prompt"), 0o644); err != nil {
		t.Fatalf("write report_prompt.md: %v", err)
	}
	question, err := Load(dir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if question.Meta.AnalysisMode != "template_files" {
		t.Fatalf("expected default analysis mode, got %q", question.Meta.AnalysisMode)
	}
}

func TestLoadRejectsUnsupportedAnalysisMode(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "meta.yaml"), []byte("id: qx\nslug: qx\ntitle: Test\ndataset: ontime\nanalysis_mode: bad_mode\nartifacts_required: report.md\n"), 0o644); err != nil {
		t.Fatalf("write meta.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "report_prompt.md"), []byte("report prompt"), 0o644); err != nil {
		t.Fatalf("write report_prompt.md: %v", err)
	}
	if _, err := Load(dir); err == nil {
		t.Fatalf("expected unsupported analysis mode to fail")
	}
}

func TestLoadDefaultsPresentationTargetToHTML(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "meta.yaml"), []byte("id: qx\nslug: qx\ntitle: Test\ndataset: ontime\nartifacts_required: report.md,visual.html\n"), 0o644); err != nil {
		t.Fatalf("write meta.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "report_prompt.md"), []byte("report prompt"), 0o644); err != nil {
		t.Fatalf("write report_prompt.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "visual_prompt.md"), []byte("visual prompt"), 0o644); err != nil {
		t.Fatalf("write visual_prompt.md: %v", err)
	}
	question, err := Load(dir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if question.Meta.PresentationTarget != "html" {
		t.Fatalf("expected default presentation target html, got %q", question.Meta.PresentationTarget)
	}
}

func TestLoadRejectsUnsupportedPresentationTarget(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "meta.yaml"), []byte("id: qx\nslug: qx\ntitle: Test\ndataset: ontime\npresentation_target: svelte\nartifacts_required: report.md,visual.html\n"), 0o644); err != nil {
		t.Fatalf("write meta.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "report_prompt.md"), []byte("report prompt"), 0o644); err != nil {
		t.Fatalf("write report_prompt.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "visual_prompt.md"), []byte("visual prompt"), 0o644); err != nil {
		t.Fatalf("write visual_prompt.md: %v", err)
	}
	if _, err := Load(dir); err == nil {
		t.Fatalf("expected unsupported presentation target to fail")
	}
}

func TestLoadMultiQueryModeFallsBackToImplicitQ0WhenNoSections(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "meta.yaml"), []byte("id: qx\nslug: qx\ntitle: Test\ndataset: ontime\nanalysis_mode: multi_query\nartifacts_required: report.md\n"), 0o644); err != nil {
		t.Fatalf("write meta.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "report_prompt.md"), []byte("report prompt"), 0o644); err != nil {
		t.Fatalf("write report_prompt.md: %v", err)
	}
	question, err := Load(dir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(question.Subquestions) != 1 || question.Subquestions[0].ID != "q0" || question.Subquestions[0].Text != "report prompt" {
		t.Fatalf("unexpected implicit question set: %+v", question.Subquestions)
	}
}

func TestLoadMultiQueryModeParsesPromptSectionsFromReportPrompt(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "meta.yaml"), []byte("id: qx\nslug: qx\ntitle: Test\ndataset: ontime\nanalysis_mode: multi_query\nartifacts_required: report.md\n"), 0o644); err != nil {
		t.Fatalf("write meta.yaml: %v", err)
	}
	report := "### main\nMain question.\n\n### q1\nFirst question?\n\n### q2\nSecond question?"
	if err := os.WriteFile(filepath.Join(dir, "report_prompt.md"), []byte(report), 0o644); err != nil {
		t.Fatalf("write report_prompt.md: %v", err)
	}
	question, err := Load(dir)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(question.Subquestions) != 3 {
		t.Fatalf("expected 3 parsed sections, got %d", len(question.Subquestions))
	}
	if question.Subquestions[0].ID != "main" || question.Subquestions[0].Text != "Main question." {
		t.Fatalf("unexpected first section: %+v", question.Subquestions[0])
	}
	if question.Subquestions[1].ID != "q1" || question.Subquestions[1].Text != "First question?" {
		t.Fatalf("unexpected second section: %+v", question.Subquestions[1])
	}
	if question.Subquestions[2].ID != "q2" || question.Subquestions[2].Text != "Second question?" {
		t.Fatalf("unexpected third section: %+v", question.Subquestions[2])
	}
}
