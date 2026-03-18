package datasets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"qforge/internal/model"
)

func Load(repoRoot, name string) (model.DatasetConfig, error) {
	path := filepath.Join(repoRoot, "datasets", name, "mcp.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return model.DatasetConfig{}, err
	}
	var cfg model.DatasetConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return model.DatasetConfig{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if cfg.Name == "" {
		cfg.Name = name
	}
	if cfg.MCPJWETokenEnv == "" {
		cfg.MCPJWETokenEnv = "MCP_JWE_TOKEN"
	}
	semanticLayerPath := filepath.Join(repoRoot, "datasets", name, "semantic_layer.md")
	semanticLayerBytes, err := os.ReadFile(semanticLayerPath)
	if err == nil {
		cfg.SemanticLayer = strings.TrimSpace(string(semanticLayerBytes))
	} else if !os.IsNotExist(err) {
		return model.DatasetConfig{}, err
	}
	return cfg, nil
}

func ResolveMCPURL(cfg model.DatasetConfig, explicitURL string) (string, string, error) {
	if explicitURL != "" {
		return explicitURL, resolveTokenFromURLOrEnv(explicitURL, cfg), nil
	}
	if cfg.MCPURL != "" {
		return cfg.MCPURL, resolveTokenFromURLOrEnv(cfg.MCPURL, cfg), nil
	}
	token := os.Getenv(cfg.MCPJWETokenEnv)
	if token == "" {
		return "", "", fmt.Errorf("missing %s", cfg.MCPJWETokenEnv)
	}
	baseURL := cfg.MCPBaseURL
	if baseURL == "" {
		baseURL = "https://mcp.demo.altinity.cloud"
	}
	return fmt.Sprintf("%s/%s/http", strings.TrimRight(baseURL, "/"), token), token, nil
}

func ResolveMCPServerName(cfg model.DatasetConfig, explicit string) string {
	if explicit != "" {
		return explicit
	}
	if cfg.DefaultMCPServerName != "" {
		return cfg.DefaultMCPServerName
	}
	return "mcp"
}

func resolveTokenFromURLOrEnv(mcpURL string, cfg model.DatasetConfig) string {
	if strings.HasSuffix(mcpURL, "/http") {
		trimmed := strings.TrimSuffix(mcpURL, "/http")
		parts := strings.Split(trimmed, "/")
		if len(parts) > 0 {
			last := parts[len(parts)-1]
			if last != "" && !strings.Contains(last, ".") {
				return last
			}
		}
	}
	return os.Getenv(cfg.MCPJWETokenEnv)
}
