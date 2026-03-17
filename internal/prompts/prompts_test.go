package prompts

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"qforge/internal/model"
	"qforge/internal/questions"
)

func TestBuildSQLPromptLoadsMarkdownAssets(t *testing.T) {
	question := model.Question{
		Dir:    filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Prompt: "Question-specific SQL guidance.",
	}
	dataset := model.DatasetConfig{
		PrimaryTable:    "default.ontime_v2",
		ForbiddenTables: "default.ontime",
	}
	got, err := BuildSQLPrompt(question, dataset)
	if err != nil {
		t.Fatalf("BuildSQLPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Stay within the configured dataset scope.") {
		t.Fatalf("expected shared core scaffold, got: %s", got)
	}
	if !strings.Contains(got, "Inspect the live schema first") {
		t.Fatalf("expected SQL-specific scaffold, got: %s", got)
	}
	if !strings.Contains(got, "Question-specific SQL guidance.") {
		t.Fatalf("expected question prompt section, got: %s", got)
	}
	if !strings.Contains(got, "default.ontime_v2") || !strings.Contains(got, "default.ontime") {
		t.Fatalf("expected dataset substitutions, got: %s", got)
	}
}

func TestBuildPresentationPromptLoadsMarkdownAssets(t *testing.T) {
	question := model.Question{
		Dir:          filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta:         model.QuestionMeta{ID: "q003", Title: "Delta ATL", VisualMode: "dynamic", VisualType: "html_heatmap"},
		ReportPrompt: "Report guidance.",
		VisualPrompt: "Visual guidance.",
	}
	result := model.CanonicalResult{
		Columns:     []string{"RowType", "Dest"},
		GeneratedAt: time.Now(),
	}
	dataset := model.DatasetConfig{
		PrimaryTable:    "default.ontime_v2",
		ForbiddenTables: "default.ontime",
	}
	got, err := BuildPresentationPrompt(question, dataset, result, "SELECT *\nFROM default.ontime_v2")
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Stay within the configured dataset scope.") {
		t.Fatalf("expected shared core scaffold, got: %s", got)
	}
	if !strings.Contains(got, "default.ontime_v2") || !strings.Contains(got, "default.ontime") {
		t.Fatalf("expected dataset substitutions, got: %s", got)
	}
	if !strings.Contains(got, "The report is a Markdown template") {
		t.Fatalf("expected presentation-specific scaffold, got: %s", got)
	}
	if !strings.Contains(got, "ontime-analyst-dashboard") {
		t.Fatalf("expected skill reference, got: %s", got)
	}
	if !strings.Contains(got, "Visual mode: `dynamic`") {
		t.Fatalf("expected visual mode in prompt, got: %s", got)
	}
	if !strings.Contains(got, "Execute the embedded saved SQL in the browser as the primary query.") {
		t.Fatalf("expected dynamic primary-query contract, got: %s", got)
	}
	if !strings.Contains(got, "visible query ledger") {
		t.Fatalf("expected dynamic ledger contract, got: %s", got)
	}
	if !strings.Contains(got, "Do not embed the primary analytical dataset") {
		t.Fatalf("expected no-embedded-dataset contract, got: %s", got)
	}
	if !strings.Contains(got, "SELECT *") || !strings.Contains(got, "FROM default.ontime_v2") {
		t.Fatalf("expected saved sql to be embedded in prompt, got: %s", got)
	}
	if strings.Contains(got, "Return exactly this fenced section:") {
		t.Fatalf("did not expect SQL-only fenced section instructions in presentation prompt, got: %s", got)
	}
	if strings.Contains(got, "qforge-result-data") || strings.Contains(got, "__QFORGE_DEFAULT_SQL__") {
		t.Fatalf("did not expect legacy injected JSON contract, got: %s", got)
	}
	if !strings.Contains(got, "Report guidance.") || !strings.Contains(got, "Visual guidance.") {
		t.Fatalf("expected question prompt sections, got: %s", got)
	}
}

func TestBuildPresentationPromptQ001UsesEnrichmentContract(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	question, err := questions.Resolve(repoRoot, "q001")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	result := model.CanonicalResult{
		Columns: []string{
			"Tail_Number",
			"Flight_Number_Reporting_Airline",
			"IATA_CODE_Reporting_Airline",
			"FlightDate",
			"Route",
		},
		GeneratedAt: time.Now(),
	}
	dataset := model.DatasetConfig{
		PrimaryTable:    "default.ontime_v2",
		ForbiddenTables: "default.ontime",
	}
	got, err := BuildPresentationPrompt(question, dataset, result, "SELECT 1")
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "ontime-analyst-dashboard") {
		t.Fatalf("expected q001 prompt to reference skill, got: %s", got)
	}
	if !strings.Contains(got, "default.airports_bts") {
		t.Fatalf("expected q001 prompt to mention airport enrichment, got: %s", got)
	}
	if !strings.Contains(got, "airport-coordinate enrichment") {
		t.Fatalf("expected q001 prompt to label map enrichment clearly, got: %s", got)
	}
	if !strings.Contains(got, "Leaflet") {
		t.Fatalf("expected q001 prompt to mention Leaflet, got: %s", got)
	}
	if !strings.Contains(got, "query ledger") {
		t.Fatalf("expected q001 prompt to inherit ledger contract, got: %s", got)
	}
	if !strings.Contains(got, "saved SQL in the browser") {
		t.Fatalf("expected q001 prompt to inherit primary-query contract, got: %s", got)
	}
	if !strings.Contains(got, "Lead Itinerary Map card and visible map container") {
		t.Fatalf("expected q001 prompt to require persistent map card, got: %s", got)
	}
	if !strings.Contains(got, "enrichment only as an upgrade to the map") {
		t.Fatalf("expected q001 prompt to treat enrichment as an upgrade, got: %s", got)
	}
	if !strings.Contains(got, "keep the map card visible with degraded-state messaging") {
		t.Fatalf("expected q001 prompt to require degraded map state, got: %s", got)
	}
}

func TestBuildPresentationPromptStaticModeUsesEmbeddedDataContract(t *testing.T) {
	question := model.Question{
		Dir: filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta: model.QuestionMeta{
			ID:         "q900",
			Title:      "Static Fixture",
			VisualMode: "static",
			VisualType: "html_ranked_dashboard",
		},
		ReportPrompt: "Report guidance.",
		VisualPrompt: "Visual guidance.",
	}
	result := model.CanonicalResult{
		Columns:     []string{"Carrier", "Flights"},
		GeneratedAt: time.Now(),
	}
	dataset := model.DatasetConfig{
		PrimaryTable:    "default.ontime_v2",
		ForbiddenTables: "default.ontime",
	}
	got, err := BuildPresentationPrompt(question, dataset, result, "SELECT Carrier, Flights FROM default.ontime_v2")
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Visual mode: `static`") {
		t.Fatalf("expected static visual mode in prompt, got: %s", got)
	}
	if !strings.Contains(got, "Build a self-contained benchmark artifact") {
		t.Fatalf("expected static artifact contract, got: %s", got)
	}
	if !strings.Contains(got, "Embed the analytical data needed by the page directly in the HTML") {
		t.Fatalf("expected embedded-data contract, got: %s", got)
	}
	if strings.Contains(got, "OnTimeAnalystDashboard::auth::jwe") {
		t.Fatalf("did not expect dynamic JWE contract in static prompt, got: %s", got)
	}
}
