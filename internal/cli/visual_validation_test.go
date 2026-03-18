package cli

import (
	"context"
	"testing"
)

func TestValidatePresentationHTMLSkipAll(t *testing.T) {
	result := validatePresentationHTML(context.Background(), presentationValidationOptions{
		SkipVisualValidation: true,
	})
	if !result.Valid {
		t.Fatalf("expected skip result to remain valid, got %+v", result)
	}
	if result.Metadata["visual_validation"] != "skipped" {
		t.Fatalf("unexpected metadata: %+v", result.Metadata)
	}
}

func TestValidatePresentationHTMLContractFailureSkipsBrowser(t *testing.T) {
	result := validatePresentationHTML(context.Background(), presentationValidationOptions{
		HTML:       "<html><body>broken</body></html>",
		VisualMode: "dynamic",
		VisualType: "html_map",
	})
	if result.Valid {
		t.Fatalf("expected invalid result, got %+v", result)
	}
	if result.Metadata["browser_validation"] != "skipped_contract_failed" {
		t.Fatalf("expected browser validation skip after contract failure, got %+v", result.Metadata)
	}
}
