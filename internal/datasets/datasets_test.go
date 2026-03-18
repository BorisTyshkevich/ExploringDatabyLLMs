package datasets

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadReadsSemanticLayerFile(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	cfg, err := Load(repoRoot, "ontime")
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if !strings.Contains(cfg.SemanticLayer, "ontime.ontime") {
		t.Fatalf("expected semantic layer to include fact table guidance, got: %s", cfg.SemanticLayer)
	}
	if strings.Contains(cfg.SemanticLayer, "ontime_semantic") {
		t.Fatalf("did not expect semantic-layer indirection, got: %s", cfg.SemanticLayer)
	}
}
