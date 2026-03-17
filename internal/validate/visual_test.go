package validate

import (
	"strings"
	"testing"
)

func TestValidateVisualHTML_DynamicValidHTML(t *testing.T) {
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

	result := ValidateVisualHTML(html, "dynamic", "html_timeseries")
	if !result.Valid {
		t.Fatalf("expected Valid=true, got errors: %v", result.Errors)
	}
	if len(result.Warnings) > 0 {
		t.Fatalf("expected no warnings, got: %v", result.Warnings)
	}
}

func TestValidateVisualHTML_DynamicMissingLedgerFails(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <main>dashboard</main>
    <footer>
        <input type="password">
        <textarea id="sql-query">SELECT 1</textarea>
    </footer>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script>
</body>
</html>`

	result := ValidateVisualHTML(html, "dynamic", "html_timeseries")
	if result.Valid {
		t.Fatal("expected missing ledger to fail validation")
	}
	assertContains(t, result.Errors, "missing query ledger")
}

func TestValidateVisualHTML_DynamicMissingLocalStorageKeyFails(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="query-ledger"></div>
    <footer>
        <input type="password">
        <textarea id="sql-query">SELECT 1</textarea>
    </footer>
</body>
</html>`

	result := ValidateVisualHTML(html, "dynamic", "html_timeseries")
	if result.Valid {
		t.Fatal("expected missing localStorage key to fail validation")
	}
	assertContains(t, result.Errors, "OnTimeAnalystDashboard::auth::jwe")
}

func TestValidateVisualHTML_DynamicMissingFooterControlsFail(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="query-ledger"></div>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script>
</body>
</html>`

	result := ValidateVisualHTML(html, "dynamic", "html_timeseries")
	if result.Valid {
		t.Fatal("expected missing footer controls to fail validation")
	}
	assertContains(t, result.Errors, "missing footer control block")
}

func TestValidateVisualHTML_DynamicMissingLeafletWithMapFails(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="query-ledger"></div>
    <div id="map"></div>
    <footer>
        <input type="password">
        <textarea id="sql-query">SELECT 1</textarea>
    </footer>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe'); L.map('map');</script>
</body>
</html>`

	result := ValidateVisualHTML(html, "dynamic", "html_map")
	if result.Valid {
		t.Fatal("expected map without Leaflet to fail validation")
	}
	assertContains(t, result.Errors, "Leaflet")
}

func TestValidateVisualHTML_DynamicEmbeddedTokenFails(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="query-ledger"></div>
    <footer>
        <input type="password">
        <textarea id="sql-query">SELECT 1</textarea>
    </footer>
    <script>
        const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U";
        localStorage.setItem('OnTimeAnalystDashboard::auth::jwe', token);
    </script>
</body>
</html>`

	result := ValidateVisualHTML(html, "dynamic", "html_timeseries")
	if result.Valid {
		t.Fatal("expected embedded token to fail validation")
	}
	assertContains(t, result.Errors, "embedded JWE/JWT token")
}

func TestValidateVisualHTML_DynamicWarningsRemainAdvisory(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; }</style>
</head>
<body>
    <div id="query-ledger">
        <div>Primary Query</div>
    </div>
    <footer>
        <input type="password">
        <textarea id="sql-query">SELECT 1</textarea>
    </footer>
    <script>localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');</script>
</body>
</html>`

	result := ValidateVisualHTML(html, "dynamic", "html_timeseries")
	if !result.Valid {
		t.Fatalf("expected warnings-only case to remain valid, got errors: %v", result.Errors)
	}
	assertContains(t, result.Warnings, "missing theme tokens")
	assertContains(t, result.Warnings, "expandable SQL")
}

func TestValidateVisualHTML_StaticNonMapRejectsRemoteAssets(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <script src="https://cdn.example.com/chart.js"></script>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <script type="application/json" id="result-data">{"rows":[{"carrier":"AA"}]}</script>
</body>
</html>`

	result := ValidateVisualHTML(html, "static", "html_ranked_dashboard")
	if result.Valid {
		t.Fatal("expected static non-map remote assets to fail validation")
	}
	assertContains(t, result.Errors, "remote scripts")
}

func TestValidateVisualHTML_StaticRequiresEmbeddedData(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <main>No embedded data</main>
</body>
</html>`

	result := ValidateVisualHTML(html, "static", "html_ranked_dashboard")
	if result.Valid {
		t.Fatal("expected static dashboard without embedded data to fail validation")
	}
	assertContains(t, result.Errors, "embedded analytical data block")
}

func TestValidateVisualHTML_StaticRejectsLiveRuntimeDependencies(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <script type="application/json" id="result-data">{"rows":[]}</script>
    <script>
        localStorage.getItem('OnTimeAnalystDashboard::auth::jwe');
        fetch('https://mcp.demo.altinity.cloud/token/openapi/execute_query?query=SELECT+1');
    </script>
</body>
</html>`

	result := ValidateVisualHTML(html, "static", "html_ranked_dashboard")
	if result.Valid {
		t.Fatal("expected static dashboard with live runtime dependencies to fail validation")
	}
	assertContains(t, result.Errors, "must not depend on live MCP fetch")
}

func TestValidateVisualHTML_StaticMapAllowsLeafletWithEmbeddedData(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>:root { --navy: #0e3a52; --sky: #3c88b5; --teal: #1f8a70; --amber: #d48a1f; }</style>
</head>
<body>
    <div id="map"></div>
    <script type="application/json" id="airports-data">{"rows":[{"code":"ATL","lat":33.64,"lon":-84.42}]}</script>
</body>
</html>`

	result := ValidateVisualHTML(html, "static", "html_map")
	if !result.Valid {
		t.Fatalf("expected static map with embedded data to be valid, got errors: %v", result.Errors)
	}
}

func assertContains(t *testing.T, values []string, fragment string) {
	t.Helper()
	for _, value := range values {
		if strings.Contains(value, fragment) {
			return
		}
	}
	t.Fatalf("expected %q in %v", fragment, values)
}
