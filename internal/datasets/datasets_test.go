package datasets

import (
	"path/filepath"
	"testing"
)

func TestLoadReadsDatasetYAMLOnly(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	cfg, err := Load(repoRoot, "ontime")
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.Name != "ontime" {
		t.Fatalf("expected dataset name to come from yaml or fallback, got: %s", cfg.Name)
	}
	if cfg.DefaultDatabase != "ontime" {
		t.Fatalf("expected default database from yaml, got: %s", cfg.DefaultDatabase)
	}
}
