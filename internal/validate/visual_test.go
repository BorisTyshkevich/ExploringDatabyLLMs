package validate

import (
	"strings"
	"testing"
)

func TestValidateVisualHTML_ValidHTML(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>
        :root {
            --navy: #0e3a52;
            --sky: #3c88b5;
            --teal: #1f8a70;
            --amber: #d48a1f;
        }
        .ledger-entry { border-bottom: 1px solid var(--border); }
        .ledger-sql { display: none; }
        .toggle-icon { font-family: monospace; }
    </style>
</head>
<body>
    <div id="query-ledger">
        <div class="ledger-entry" data-expanded="false">
            <div class="ledger-row" onclick="toggleLedgerEntry(this)">
                <span class="toggle-icon">▶</span>
                <span>Primary Query</span>
            </div>
            <div class="ledger-sql hidden"><pre>SELECT * FROM flights</pre></div>
        </div>
    </div>
    <footer>
        <input type="password" id="jwe-token" placeholder="JWE Token">
        <textarea id="sql-query">SELECT * FROM flights</textarea>
    </footer>
    <script>
        const key = 'OnTimeAnalystDashboard::auth::jwe';
        localStorage.getItem(key);
        function toggleLedgerEntry(el) {}
    </script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if !result.Valid {
		t.Errorf("expected Valid=true, got false with errors: %v", result.Errors)
	}
	if len(result.Errors) > 0 {
		t.Errorf("expected no errors, got: %v", result.Errors)
	}
	if len(result.Warnings) > 0 {
		t.Errorf("expected no warnings, got: %v", result.Warnings)
	}
}

func TestValidateVisualHTML_MissingLeafletWithMap(t *testing.T) {
	// Dashboard with map content but no Leaflet should fail
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="map"></div>
    <div id="query-ledger"></div>
    <input type="password"><textarea id="sql"></textarea>
    <script>
        localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');
        const map = L.map('map');
    </script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if result.Valid {
		t.Error("expected Valid=false for missing Leaflet with map content")
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(result.Errors), result.Errors)
	}
	if !strings.Contains(result.Errors[0], "Leaflet") {
		t.Errorf("expected Leaflet error, got: %s", result.Errors[0])
	}
}

func TestValidateVisualHTML_NoLeafletNoMap(t *testing.T) {
	// Dashboard without map content doesn't need Leaflet
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="heatmap"></div>
    <div id="query-ledger"></div>
    <input type="password"><textarea id="sql"></textarea>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if !result.Valid {
		t.Errorf("non-map dashboard should be valid without Leaflet: %v", result.Errors)
	}
}

func TestValidateVisualHTML_EmbeddedJWT(t *testing.T) {
	// Simulated JWT token in HTML
	html := `<!DOCTYPE html>
<html>
<head>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="query-ledger"></div>
    <input type="password"><textarea id="sql"></textarea>
    <script>
        const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";
        localStorage.setItem('OnTimeAnalystDashboard::auth::jwe', token);
    </script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if result.Valid {
		t.Error("expected Valid=false for embedded token")
	}
	foundTokenError := false
	for _, err := range result.Errors {
		if strings.Contains(err, "embedded") || strings.Contains(err, "JWE") {
			foundTokenError = true
			break
		}
	}
	if !foundTokenError {
		t.Errorf("expected embedded token error, got: %v", result.Errors)
	}
}

func TestValidateVisualHTML_MissingThemeTokens(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>:root { --navy: #0e3a52; }</style>
</head>
<body>
    <div id="query-ledger"></div>
    <input type="password"><textarea id="sql"></textarea>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if !result.Valid {
		t.Errorf("missing theme tokens should be warning, not error: %v", result.Errors)
	}
	if len(result.Warnings) == 0 {
		t.Error("expected warnings for missing theme tokens")
	}

	foundThemeWarning := false
	for _, warn := range result.Warnings {
		if strings.Contains(warn, "theme tokens") {
			foundThemeWarning = true
			if !strings.Contains(warn, "--sky") || !strings.Contains(warn, "--teal") || !strings.Contains(warn, "--amber") {
				t.Errorf("expected missing tokens in warning, got: %s", warn)
			}
			break
		}
	}
	if !foundThemeWarning {
		t.Errorf("expected theme token warning, got: %v", result.Warnings)
	}
}

func TestValidateVisualHTML_MissingLedger(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="main-content"></div>
    <input type="password"><textarea id="sql"></textarea>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if !result.Valid {
		t.Errorf("missing ledger should be warning, not error: %v", result.Errors)
	}

	foundLedgerWarning := false
	for _, warn := range result.Warnings {
		if strings.Contains(warn, "no query ledger element") {
			foundLedgerWarning = true
			break
		}
	}
	if !foundLedgerWarning {
		t.Errorf("expected ledger warning, got: %v", result.Warnings)
	}
}

func TestValidateVisualHTML_LedgerWithoutExpandableSQL(t *testing.T) {
	// Ledger exists but lacks expandable SQL structure
	html := `<!DOCTYPE html>
<html>
<head>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="query-ledger">
        <div>Primary Query | OK | 10</div>
    </div>
    <input type="password"><textarea id="sql"></textarea>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if !result.Valid {
		t.Errorf("missing expandable SQL should be warning, not error: %v", result.Errors)
	}

	foundExpandableWarning := false
	for _, warn := range result.Warnings {
		if strings.Contains(warn, "expandable SQL") {
			foundExpandableWarning = true
			break
		}
	}
	if !foundExpandableWarning {
		t.Errorf("expected expandable SQL warning, got: %v", result.Warnings)
	}
}

func TestValidateVisualHTML_MissingLocalStorageKey(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="query-ledger"></div>
    <input type="password"><textarea id="sql"></textarea>
    <script>localStorage.getItem('some-other-key');</script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if !result.Valid {
		t.Errorf("missing localStorage key should be warning, not error: %v", result.Errors)
	}

	foundKeyWarning := false
	for _, warn := range result.Warnings {
		if strings.Contains(warn, "OnTimeAnalystDashboard::auth::jwe") {
			foundKeyWarning = true
			break
		}
	}
	if !foundKeyWarning {
		t.Errorf("expected localStorage key warning, got: %v", result.Warnings)
	}
}

func TestValidateVisualHTML_MissingFooterControls(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="query-ledger"></div>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if !result.Valid {
		t.Errorf("missing footer controls should be warning, not error: %v", result.Errors)
	}

	foundTokenInputWarning := false
	foundSQLTextareaWarning := false
	for _, warn := range result.Warnings {
		if strings.Contains(warn, "token input") {
			foundTokenInputWarning = true
		}
		if strings.Contains(warn, "SQL textarea") {
			foundSQLTextareaWarning = true
		}
	}
	if !foundTokenInputWarning {
		t.Errorf("expected token input warning, got: %v", result.Warnings)
	}
	if !foundSQLTextareaWarning {
		t.Errorf("expected SQL textarea warning, got: %v", result.Warnings)
	}
}

func TestValidateVisualHTML_LeafletVariants(t *testing.T) {
	variants := []string{
		`<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>`,
		`<script src="https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/leaflet.min.js"></script>`,
		`<script src="/vendor/leaflet.js"></script>`,
		`<link href="leaflet.min.js">`,
	}

	baseHTML := `<!DOCTYPE html><html><head>%s
<style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head><body><div id="query-ledger"></div>
<input type="password"><textarea id="sql"></textarea>
<script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script></body></html>`

	for _, variant := range variants {
		html := strings.Replace(baseHTML, "%s", variant, 1)
		result := ValidateVisualHTML(html)
		if !result.Valid {
			t.Errorf("Leaflet variant should be valid: %s, errors: %v", variant, result.Errors)
		}
	}
}

func TestValidateVisualHTML_MultipleErrors(t *testing.T) {
	// Map content without Leaflet and has embedded token
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; }</style>
</head>
<body>
    <div id="map"></div>
    <script>
        const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U";
        L.map('map');
    </script>
</body>
</html>`

	result := ValidateVisualHTML(html)

	if result.Valid {
		t.Error("expected Valid=false with multiple errors")
	}
	if len(result.Errors) < 2 {
		t.Errorf("expected at least 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}
