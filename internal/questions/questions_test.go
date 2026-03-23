package questions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsAnalysisModeToJSONArtifact(t *testing.T) {
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
	if question.Meta.AnalysisMode != "json_artifact" {
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
