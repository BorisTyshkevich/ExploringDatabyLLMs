package cli

import (
	"context"
	"strconv"
	"strings"

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
