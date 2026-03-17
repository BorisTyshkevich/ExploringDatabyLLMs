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
	themeTokens = []string{"--navy", "--sky", "--teal", "--amber"}

	// JWE compact format has 5 parts; JWT has 3 parts.
	jwePattern = regexp.MustCompile(`["']eyJ[A-Za-z0-9_-]{20,}\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)
)

// ValidateVisualHTML checks a visual.html file against the selected skill contract.
func ValidateVisualHTML(html, visualMode, visualType string) VisualValidationResult {
	result := VisualValidationResult{Valid: true}

	checkNoEmbeddedTokens(html, &result)
	checkThemeTokens(html, &result)

	mode := strings.ToLower(strings.TrimSpace(visualMode))
	if mode == "" {
		mode = "dynamic"
	}
	isMap := strings.EqualFold(strings.TrimSpace(visualType), "html_map")

	switch mode {
	case "static":
		checkStaticHTML(html, isMap, &result)
	default:
		checkDynamicHTML(html, isMap, &result)
	}

	return result
}

func checkDynamicHTML(html string, isMap bool, result *VisualValidationResult) {
	checkLeaflet(html, isMap, result)
	checkQueryLedger(html, true, result)
	checkLocalStorageKey(html, true, result)
	checkFooterControls(html, true, result)
}

func checkStaticHTML(html string, isMap bool, result *VisualValidationResult) {
	checkLeaflet(html, isMap, result)
	checkStaticRemoteAssets(html, isMap, result)
	checkStaticEmbeddedData(html, result)
	checkStaticRuntimeDependencies(html, result)
}

func checkLeaflet(html string, isMap bool, result *VisualValidationResult) {
	if !isMap {
		return
	}
	lower := strings.ToLower(html)
	hasLeaflet := strings.Contains(lower, "leaflet@") ||
		strings.Contains(lower, "leaflet.js") ||
		strings.Contains(lower, "leaflet.min.js") ||
		strings.Contains(lower, "unpkg.com/leaflet") ||
		strings.Contains(lower, "cdnjs.cloudflare.com/ajax/libs/leaflet")
	if !hasLeaflet {
		result.Valid = false
		result.Errors = append(result.Errors, "map dashboard missing Leaflet asset")
	}
}

func checkNoEmbeddedTokens(html string, result *VisualValidationResult) {
	if jwePattern.MatchString(html) {
		result.Valid = false
		result.Errors = append(result.Errors, "embedded JWE/JWT token detected; tokens must not be hardcoded")
	}
}

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

func checkQueryLedger(html string, required bool, result *VisualValidationResult) {
	lower := strings.ToLower(html)
	hasLedger := strings.Contains(html, `id="ledger`) ||
		strings.Contains(html, `id="query-ledger`) ||
		strings.Contains(html, `id="queryLedger`) ||
		strings.Contains(html, `class="ledger`) ||
		strings.Contains(lower, "query ledger")
	if !hasLedger {
		if required {
			result.Valid = false
			result.Errors = append(result.Errors, "missing query ledger element")
		} else {
			result.Warnings = append(result.Warnings, "no query ledger element found")
		}
		return
	}

	hasExpandableSQL := strings.Contains(html, "ledger-entry") ||
		strings.Contains(html, "toggleLedgerEntry") ||
		strings.Contains(html, "ledger-sql") ||
		strings.Contains(html, "toggle-icon")
	if !hasExpandableSQL {
		result.Warnings = append(result.Warnings, "query ledger may not have expandable SQL")
	}
}

func checkLocalStorageKey(html string, required bool, result *VisualValidationResult) {
	if strings.Contains(html, "OnTimeAnalystDashboard::auth::jwe") {
		return
	}
	if required {
		result.Valid = false
		result.Errors = append(result.Errors, "missing localStorage key 'OnTimeAnalystDashboard::auth::jwe'")
		return
	}
	result.Warnings = append(result.Warnings, "missing localStorage key 'OnTimeAnalystDashboard::auth::jwe'")
}

func checkFooterControls(html string, required bool, result *VisualValidationResult) {
	lower := strings.ToLower(html)
	footerStart := strings.LastIndex(lower, "<footer")
	if footerStart < 0 {
		footerStart = strings.LastIndex(lower, `data-role="controls"`)
	}
	if footerStart < 0 {
		if required {
			result.Valid = false
			result.Errors = append(result.Errors, "missing footer control block")
		} else {
			result.Warnings = append(result.Warnings, "missing footer control block")
		}
		return
	}

	footerHTML := lower[footerStart:]
	hasTokenInput := strings.Contains(footerHTML, `type="password"`) ||
		(strings.Contains(footerHTML, "token") && strings.Contains(footerHTML, "<input"))
	hasSQLTextarea := strings.Contains(footerHTML, "<textarea") &&
		(strings.Contains(footerHTML, "sql") || strings.Contains(footerHTML, "query"))

	if !hasTokenInput {
		result.Valid = false
		result.Errors = append(result.Errors, "missing token input control in footer")
	}
	if !hasSQLTextarea {
		result.Valid = false
		result.Errors = append(result.Errors, "missing SQL textarea control in footer")
	}
}

func checkStaticRemoteAssets(html string, isMap bool, result *VisualValidationResult) {
	lower := strings.ToLower(html)
	if isMap {
		return
	}
	if strings.Contains(lower, "<script src=") {
		result.Valid = false
		result.Errors = append(result.Errors, "static non-map dashboard must not load remote scripts")
	}
	if strings.Contains(lower, "<link") {
		result.Valid = false
		result.Errors = append(result.Errors, "static non-map dashboard must not load remote stylesheets")
	}
}

func checkStaticEmbeddedData(html string, result *VisualValidationResult) {
	lower := strings.ToLower(html)
	hasEmbeddedJSON := strings.Contains(lower, `type="application/json"`)
	hasEmbeddedCSV := strings.Contains(lower, `type="text/csv"`)
	if hasEmbeddedJSON || hasEmbeddedCSV {
		return
	}
	result.Valid = false
	result.Errors = append(result.Errors, "static dashboard missing embedded analytical data block")
}

func checkStaticRuntimeDependencies(html string, result *VisualValidationResult) {
	lower := strings.ToLower(html)
	if strings.Contains(lower, "openapi/execute_query") || strings.Contains(lower, "mcp.demo.altinity.cloud") || strings.Contains(lower, "localstorage") {
		result.Valid = false
		result.Errors = append(result.Errors, "static dashboard must not depend on live MCP fetch or localStorage token flow")
	}
}
