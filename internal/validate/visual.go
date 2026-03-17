package validate

import (
	"regexp"
	"strings"
)

// VisualValidationResult contains the outcome of visual.html validation.
type VisualValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

var (
	// themeTokens are the CSS custom properties required by the design system.
	themeTokens = []string{"--navy", "--sky", "--teal", "--amber"}

	// jwePattern matches hardcoded JWE/JWT tokens (dot-separated base64url segments).
	// JWE compact format: header.encryptedKey.iv.ciphertext.tag (5 parts)
	// JWT format: header.payload.signature (3 parts)
	// We look for long base64url strings in quotes starting with eyJ (base64 for {"...).
	jwePattern = regexp.MustCompile(`["']eyJ[A-Za-z0-9_-]{20,}\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)
)

// ValidateVisualHTML checks a visual.html file against the skill contract.
// Errors indicate contract violations that should fail the run.
// Warnings indicate deviations that are recorded but don't fail the run.
func ValidateVisualHTML(html string) VisualValidationResult {
	result := VisualValidationResult{Valid: true}

	// Error checks
	checkLeafletCDN(html, &result)
	checkNoEmbeddedTokens(html, &result)

	// Warning checks
	checkThemeTokens(html, &result)
	checkQueryLedger(html, &result)
	checkLocalStorageKey(html, &result)
	checkFooterControls(html, &result)

	return result
}

// checkLeafletCDN verifies Leaflet library is loaded from CDN when the dashboard uses maps.
// If no map-related content is detected, missing Leaflet is a warning instead of an error.
func checkLeafletCDN(html string, result *VisualValidationResult) {
	lower := strings.ToLower(html)
	hasLeaflet := strings.Contains(lower, "leaflet@") ||
		strings.Contains(lower, "leaflet.js") ||
		strings.Contains(lower, "leaflet.min.js") ||
		strings.Contains(lower, "unpkg.com/leaflet") ||
		strings.Contains(lower, "cdnjs.cloudflare.com/ajax/libs/leaflet")

	if hasLeaflet {
		return
	}

	// Check if the dashboard appears to use maps
	usesMap := strings.Contains(lower, `id="map"`) ||
		strings.Contains(lower, `id="map-`) ||
		strings.Contains(lower, `class="map"`) ||
		strings.Contains(lower, `class="map-`) ||
		strings.Contains(html, "L.map") ||
		strings.Contains(html, "L.marker") ||
		strings.Contains(html, "L.tileLayer") ||
		strings.Contains(lower, "latitude") ||
		strings.Contains(lower, "longitude") ||
		strings.Contains(html, "LatLng")

	if usesMap {
		// Map content without Leaflet is an error
		result.Valid = false
		result.Errors = append(result.Errors, "map content detected but missing Leaflet CDN reference")
	}
	// Non-map dashboards don't require Leaflet - no warning needed
}

// checkNoEmbeddedTokens ensures no hardcoded JWE/JWT tokens appear in the HTML.
func checkNoEmbeddedTokens(html string, result *VisualValidationResult) {
	if jwePattern.MatchString(html) {
		result.Valid = false
		result.Errors = append(result.Errors, "embedded JWE/JWT token detected; tokens must not be hardcoded")
	}
}

// checkThemeTokens verifies the dashboard uses the required design system tokens.
func checkThemeTokens(html string, result *VisualValidationResult) {
	var missing []string
	for _, token := range themeTokens {
		if !strings.Contains(html, token) {
			missing = append(missing, token)
		}
	}
	if len(missing) > 0 {
		result.Warnings = append(result.Warnings, "missing theme tokens: "+strings.Join(missing, ", "))
	}
}

// checkQueryLedger verifies the presence of a query ledger element with expandable SQL support.
func checkQueryLedger(html string, result *VisualValidationResult) {
	// Look for common ledger patterns
	hasLedger := strings.Contains(html, `id="ledger`) ||
		strings.Contains(html, `id="query-ledger`) ||
		strings.Contains(html, `id="queryLedger`) ||
		strings.Contains(html, `class="ledger`) ||
		strings.Contains(strings.ToLower(html), "query ledger")

	if !hasLedger {
		result.Warnings = append(result.Warnings, "no query ledger element found (expected id containing 'ledger')")
		return
	}

	// Check for expandable SQL structure (new contract)
	hasExpandableSQL := strings.Contains(html, "ledger-entry") ||
		strings.Contains(html, "toggleLedgerEntry") ||
		strings.Contains(html, "ledger-sql") ||
		strings.Contains(html, "toggle-icon")

	if !hasExpandableSQL {
		result.Warnings = append(result.Warnings, "query ledger may not have expandable SQL (expected ledger-entry or toggleLedgerEntry)")
	}
}

// checkLocalStorageKey verifies the correct localStorage key is used for JWE.
func checkLocalStorageKey(html string, result *VisualValidationResult) {
	if !strings.Contains(html, "OnTimeAnalystDashboard::auth::jwe") {
		result.Warnings = append(result.Warnings, "missing localStorage key 'OnTimeAnalystDashboard::auth::jwe'")
	}
}

// checkFooterControls verifies the presence of token input and SQL textarea controls.
func checkFooterControls(html string, result *VisualValidationResult) {
	lower := strings.ToLower(html)

	// Check for token input
	hasTokenInput := strings.Contains(lower, `type="password"`) ||
		strings.Contains(lower, "token") && strings.Contains(lower, "<input")

	// Check for SQL textarea
	hasSQLTextarea := strings.Contains(lower, "<textarea") &&
		(strings.Contains(lower, "sql") || strings.Contains(lower, "query"))

	if !hasTokenInput {
		result.Warnings = append(result.Warnings, "no token input control found in footer")
	}
	if !hasSQLTextarea {
		result.Warnings = append(result.Warnings, "no SQL textarea control found in footer")
	}
}
