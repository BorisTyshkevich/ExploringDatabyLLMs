package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"qforge/internal/model"
	"qforge/internal/validate"
	browservalidate "qforge/internal/validate/browser"
)

type presentationValidationOptions struct {
	RunDir               string
	HTMLPath             string
	HTML                 string
	Model                string
	VisualMode           string
	VisualType           string
	Token                string
	SkipVisualValidation bool
	SkipBrowserLiveFetch bool
	Verbose              bool
}

type presentationValidationResult struct {
	Valid    bool
	Metadata map[string]string
}

type presentationArtifactResult struct {
	HTML            string
	BuildDurationMS int64
	Metadata        map[string]string
}

type reactPackageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

type presentationArtifactError struct {
	Stage string
	Err   error
}

func (e presentationArtifactError) Error() string {
	if e.Err == nil {
		return e.Stage
	}
	return e.Err.Error()
}

func validatePresentationHTML(ctx context.Context, opts presentationValidationOptions) presentationValidationResult {
	result := presentationValidationResult{
		Valid:    true,
		Metadata: map[string]string{},
	}
	if opts.SkipVisualValidation {
		result.Metadata["visual_validation"] = "skipped"
		result.Metadata["browser_validation"] = "skipped"
		result.Metadata["browser_validation_live_fetch"] = "skipped"
		logf(opts.Verbose, opts.Model, "phase=visual_validation status=skipped")
		logf(opts.Verbose, opts.Model, "phase=browser_validation status=skipped")
		return result
	}

	contractResult := validate.ValidateVisualHTML(opts.HTML, opts.VisualMode, opts.VisualType)
	if !contractResult.Valid {
		result.Valid = false
		result.Metadata["visual_validation"] = "failed"
		result.Metadata["visual_validation_errors"] = strings.Join(contractResult.Errors, "; ")
		result.Metadata["browser_validation"] = "skipped_contract_failed"
		logf(opts.Verbose, opts.Model, "phase=visual_validation status=failed errors=%q", result.Metadata["visual_validation_errors"])
	} else {
		result.Metadata["visual_validation"] = "ok"
		logf(opts.Verbose, opts.Model, "phase=visual_validation status=ok")
	}
	if len(contractResult.Warnings) > 0 {
		result.Metadata["visual_validation_warnings"] = strings.Join(contractResult.Warnings, "; ")
		logf(opts.Verbose, opts.Model, "visual_validation_warnings=%q", result.Metadata["visual_validation_warnings"])
	}
	if !contractResult.Valid {
		return result
	}

	logf(opts.Verbose, opts.Model, "phase=browser_validation status=started")
	browserResult := browservalidate.Validate(ctx, browservalidate.Options{
		HTMLPath:      opts.HTMLPath,
		VisualMode:    opts.VisualMode,
		Token:         opts.Token,
		SkipLiveFetch: opts.SkipBrowserLiveFetch,
	})
	if !browserResult.Valid {
		result.Valid = false
		result.Metadata["browser_validation"] = "failed"
		result.Metadata["browser_validation_errors"] = strings.Join(browserResult.Errors, "; ")
		logf(opts.Verbose, opts.Model, "phase=browser_validation status=failed errors=%q", result.Metadata["browser_validation_errors"])
	} else {
		result.Metadata["browser_validation"] = "ok"
		logf(opts.Verbose, opts.Model, "phase=browser_validation status=ok")
	}
	if len(browserResult.Warnings) > 0 {
		result.Metadata["browser_validation_warnings"] = strings.Join(browserResult.Warnings, "; ")
		logf(opts.Verbose, opts.Model, "browser_validation_warnings=%q", result.Metadata["browser_validation_warnings"])
	}
	if len(browserResult.ConsoleErrors) > 0 {
		result.Metadata["browser_validation_console_errors"] = strings.Join(browserResult.ConsoleErrors, "; ")
	}
	if browserResult.MatchedRequestURL != "" {
		result.Metadata["browser_validation_request_url"] = browserResult.MatchedRequestURL
	}
	if browserResult.MatchedResponseCode > 0 {
		result.Metadata["browser_validation_response_code"] = strconv.FormatInt(browserResult.MatchedResponseCode, 10)
	}
	if browserResult.StatusText != "" {
		result.Metadata["browser_validation_status"] = browserResult.StatusText
	}
	switch {
	case browserResult.LiveFetchSucceeded:
		result.Metadata["browser_validation_live_fetch"] = "ok"
		logf(opts.Verbose, opts.Model, "phase=browser_validation_live_fetch status=ok")
	case browserResult.LiveFetchAttempted:
		result.Metadata["browser_validation_live_fetch"] = "failed"
		logf(opts.Verbose, opts.Model, "phase=browser_validation_live_fetch status=failed")
	case browserResult.LiveFetchSkipped:
		result.Metadata["browser_validation_live_fetch"] = browserResult.SkipReason
		logf(opts.Verbose, opts.Model, "phase=browser_validation_live_fetch status=%s", browserResult.SkipReason)
	default:
		result.Metadata["browser_validation_live_fetch"] = "not_attempted"
	}
	return result
}

func materializePresentationArtifact(outDir string, question model.Question, artifacts model.ArtifactPaths, rawOutput string, notBefore time.Time, modelName string, verbose bool) (presentationArtifactResult, error) {
	result := presentationArtifactResult{
		Metadata: map[string]string{
			"presentation_target": normalizePresentationTarget(question.Meta.PresentationTarget),
		},
	}
	if normalizePresentationTarget(question.Meta.PresentationTarget) != "react" {
		html, err := loadVisualArtifact(rawOutput, outDir, notBefore)
		if err != nil {
			return result, presentationArtifactError{Stage: "generation", Err: err}
		}
		result.HTML = html
		return result, nil
	}

	if err := loadReactSourceArtifact(artifacts.VisualSourceDir, notBefore); err != nil {
		return result, presentationArtifactError{Stage: "generation", Err: err}
	}
	if err := validateReactSourceArtifacts(artifacts.VisualSourceDir); err != nil {
		result.Metadata["react_source_validation"] = "failed"
		result.Metadata["react_source_validation_errors"] = err.Error()
		return result, presentationArtifactError{Stage: "render", Err: err}
	}
	result.Metadata["react_source_validation"] = "ok"

	buildStartedAt := time.Now()
	if err := buildReactPresentation(artifacts.VisualSourceDir, artifacts.VisualBuildDir, modelName, verbose); err != nil {
		result.Metadata["react_build"] = "failed"
		result.Metadata["react_build_errors"] = err.Error()
		return result, presentationArtifactError{Stage: "render", Err: err}
	}
	result.BuildDurationMS = time.Since(buildStartedAt).Milliseconds()
	result.Metadata["react_build"] = "ok"
	result.Metadata["react_build_duration_ms"] = strconv.FormatInt(result.BuildDurationMS, 10)

	html, err := finalizeBuiltReactArtifact(artifacts)
	if err != nil {
		result.Metadata["react_build"] = "failed"
		result.Metadata["react_build_errors"] = err.Error()
		return result, presentationArtifactError{Stage: "render", Err: err}
	}
	result.HTML = html
	if size, err := dirSizeBytes(artifacts.VisualSourceDir); err == nil {
		result.Metadata["react_source_bytes"] = strconv.FormatInt(size, 10)
	}
	if size, err := dirSizeBytes(artifacts.VisualBuildDir); err == nil {
		result.Metadata["react_build_bytes"] = strconv.FormatInt(size, 10)
	}
	return result, nil
}

func normalizePresentationTarget(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "react") {
		return "react"
	}
	return "html"
}

func loadReactSourceArtifact(sourceDir string, notBefore time.Time) error {
	required := []string{
		filepath.Join(sourceDir, "package.json"),
		filepath.Join(sourceDir, "index.html"),
		filepath.Join(sourceDir, "src", "main.jsx"),
		filepath.Join(sourceDir, "src", "App.jsx"),
	}
	for _, path := range required {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("missing react source artifact %s: %w", path, err)
		}
		if info.ModTime().Before(notBefore) {
			return fmt.Errorf("stale react source artifact %s", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read react source artifact %s: %w", path, err)
		}
		if strings.TrimSpace(string(data)) == "" {
			return fmt.Errorf("empty react source artifact %s", path)
		}
	}
	return nil
}

func validateReactSourceArtifacts(sourceDir string) error {
	required := []string{
		filepath.Join(sourceDir, "package.json"),
		filepath.Join(sourceDir, "index.html"),
		filepath.Join(sourceDir, "src", "main.jsx"),
		filepath.Join(sourceDir, "src", "App.jsx"),
	}
	for _, path := range required {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read required react source file %s: %w", path, err)
		}
		if strings.TrimSpace(string(data)) == "" {
			return fmt.Errorf("required react source file %s is empty", path)
		}
	}
	packageJSONPath := filepath.Join(sourceDir, "package.json")
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return err
	}
	var pkg reactPackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("parse %s: %w", packageJSONPath, err)
	}
	if strings.TrimSpace(pkg.Scripts["build"]) == "" {
		return fmt.Errorf("%s must declare a build script", packageJSONPath)
	}
	lower := strings.ToLower(string(data))
	if strings.Contains(lower, "\"next\"") || strings.Contains(lower, "\"express\"") {
		return fmt.Errorf("%s must not depend on server-side runtimes", packageJSONPath)
	}
	return nil
}

func buildReactPresentation(sourceDir, buildDir, modelName string, verbose bool) error {
	_ = os.RemoveAll(buildDir)
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		return err
	}
	install := exec.Command("npm", "install", "--no-fund", "--no-audit")
	install.Dir = sourceDir
	if output, err := install.CombinedOutput(); err != nil {
		logf(verbose, modelName, "phase=react_build status=failed step=install output=%q", string(output))
		return fmt.Errorf("npm install failed: %w", err)
	}
	build := exec.Command("npm", "run", "build")
	build.Dir = sourceDir
	if output, err := build.CombinedOutput(); err != nil {
		logf(verbose, modelName, "phase=react_build status=failed step=build output=%q", string(output))
		return fmt.Errorf("npm run build failed: %w", err)
	}
	indexPath := filepath.Join(buildDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		return fmt.Errorf("react build missing %s: %w", indexPath, err)
	}
	return nil
}

func finalizeBuiltReactArtifact(artifacts model.ArtifactPaths) (string, error) {
	buildHTMLPath := filepath.Join(artifacts.VisualBuildDir, "index.html")
	buildHTML, err := os.ReadFile(buildHTMLPath)
	if err != nil {
		return "", fmt.Errorf("read built react html: %w", err)
	}
	if err := os.RemoveAll(artifacts.VisualAssetsDir); err != nil && !os.IsNotExist(err) {
		return "", err
	}
	buildAssetsDir := filepath.Join(artifacts.VisualBuildDir, "visual_assets")
	if _, err := os.Stat(buildAssetsDir); err == nil {
		if err := copyDir(buildAssetsDir, artifacts.VisualAssetsDir); err != nil {
			return "", err
		}
	}
	if err := os.WriteFile(artifacts.VisualHTML, buildHTML, 0o644); err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buildHTML)), nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info != nil && info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

func dirSizeBytes(root string) (int64, error) {
	var total int64
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total, err
}
