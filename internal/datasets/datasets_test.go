package datasets

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadReadsSemanticPromptFile(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	cfg, err := Load(repoRoot, "ontime")
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.PrimaryTable != "ontime.ontime" {
		t.Fatalf("expected primary table to be ontime.ontime, got: %s", cfg.PrimaryTable)
	}
	if !strings.Contains(cfg.DiscoveryPrompt, "ontime_semantic.active_joins") {
		t.Fatalf("expected discovery prompt to include semantic joins view, got: %s", cfg.DiscoveryPrompt)
	}
}
