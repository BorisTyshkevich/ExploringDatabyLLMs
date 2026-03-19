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
		SemanticLayer: "Use `ontime.ontime` and `ontime.airports_latest`.",
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
	if !strings.Contains(got, "Dataset semantic layer:") || !strings.Contains(got, "ontime.airports_latest") {
		t.Fatalf("expected dataset semantic layer guidance, got: %s", got)
	}
	if !strings.Contains(got, "answer.raw.json") || !strings.Contains(got, "\"sql\"") || !strings.Contains(got, "\"report_markdown\"") || !strings.Contains(got, "\"metrics\"") {
		t.Fatalf("expected file-based single json analysis contract, got: %s", got)
	}
	if !strings.Contains(got, "plain JSON bytes only") || !strings.Contains(got, "ignore stdout") {
		t.Fatalf("expected strict answer.raw.json file rules, got: %s", got)
	}
	if !strings.Contains(got, "{{metric.<name>}}") || !strings.Contains(got, "Do not invent any placeholder") {
		t.Fatalf("expected explicit metric placeholder guidance, got: %s", got)
	}
	if strings.Contains(got, "CLE -> BNA -> PNS") || strings.Contains(got, "\"max_hops\": \"8\"") {
		t.Fatalf("did not expect question-specific example values in shared analysis prompt, got: %s", got)
	}
	if strings.Contains(got, "dataset constraints") || strings.Contains(got, "ontime_semantic") {
		t.Fatalf("did not expect legacy dataset constraints or semantic db references, got: %s", got)
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
		SemanticLayer: "Use `ontime.ontime` and `ontime.airports_latest`.",
	}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT *\nFROM ontime.ontime", "# Report\n{{data_overview_md}}", "{\"sql\":\"SELECT * FROM ontime.ontime\",\"report_markdown\":\"# Report\\n\\n{{data_overview_md}}\",\"metrics\":{\"named_values\":{\"max_hops\":\"8\"}}}")
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Stay within the configured dataset scope.") {
		t.Fatalf("expected shared core scaffold, got: %s", got)
	}
	if !strings.Contains(got, "Generate only the visual artifact.") {
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
	if !strings.Contains(got, "SELECT *") || !strings.Contains(got, "FROM ontime.ontime") {
		t.Fatalf("expected saved sql to be embedded in prompt, got: %s", got)
	}
	if !strings.Contains(got, "```html") || strings.Contains(got, "```report\nUse placeholders only") {
		t.Fatalf("expected html-only output contract in presentation prompt, got: %s", got)
	}
	if strings.Contains(got, "qforge-result-data") || strings.Contains(got, "__QFORGE_DEFAULT_SQL__") {
		t.Fatalf("did not expect legacy injected JSON contract, got: %s", got)
	}
	if !strings.Contains(got, "saved report template") && !strings.Contains(strings.ToLower(got), "saved report template") {
		t.Fatalf("expected saved report template context, got: %s", got)
	}
	if !strings.Contains(got, "Saved analysis artifact:") || !strings.Contains(got, "\"report_markdown\"") {
		t.Fatalf("expected saved analysis json context, got: %s", got)
	}
	if !strings.Contains(got, "Dataset semantic layer:") || !strings.Contains(got, "ontime.airports_latest") {
		t.Fatalf("expected inlined semantic-layer guidance, got: %s", got)
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
		SemanticLayer: "Use `ontime.ontime` and `ontime.airports_latest`.",
	}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT 1", "# Report\n{{data_overview_md}}", "{\"sql\":\"SELECT 1\",\"report_markdown\":\"# Report\\n\\n{{data_overview_md}}\",\"metrics\":{\"named_values\":{\"max_hops\":\"8\"}}}")
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "ontime-analyst-dashboard") {
		t.Fatalf("expected q001 prompt to reference skill, got: %s", got)
	}
	if !strings.Contains(got, "ontime.airports_latest") {
		t.Fatalf("expected q001 prompt to mention airport enrichment, got: %s", got)
	}
	if !strings.Contains(got, "airport-coordinate enrichment") {
		t.Fatalf("expected q001 prompt to label map enrichment clearly, got: %s", got)
	}
	if !strings.Contains(got, "query ledger") {
		t.Fatalf("expected q001 prompt to inherit ledger contract, got: %s", got)
	}
	if !strings.Contains(got, "saved SQL in the browser") {
		t.Fatalf("expected q001 prompt to inherit primary-query contract, got: %s", got)
	}
	if !strings.Contains(got, "keep the map card visible with degraded-state messaging") {
		t.Fatalf("expected q001 prompt to require degraded map state, got: %s", got)
	}
	if !strings.Contains(got, "Dataset semantic layer:") {
		t.Fatalf("expected semantic layer heading when semantic guidance is present, got: %s", got)
	}
}

func TestBuildPresentationPromptQ007NoLongerUsesSemanticDiscovery(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	question, err := questions.Resolve(repoRoot, "q007")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if strings.Contains(question.VisualPrompt, "ontime.airports_latest") {
		t.Fatalf("did not expect q007 local visual prompt to hardcode airport enrichment source: %s", question.VisualPrompt)
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
		SemanticLayer: "Use `ontime.ontime` and `ontime.airports_latest`.",
	}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT 1", "# Report\n{{data_overview_md}}", "{\"sql\":\"SELECT 1\",\"report_markdown\":\"# Report\\n\\n{{data_overview_md}}\",\"metrics\":{\"named_values\":{\"max_hops\":\"8\"}}}")
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if strings.Contains(got, "ontime_semantic") || strings.Contains(got, "Dataset discovery:") {
		t.Fatalf("did not expect semantic discovery guidance, got: %s", got)
	}
	if strings.Contains(got, "run an explicit airport-coordinate enrichment query against `ontime.airports_latest` using airport codes parsed from the route strings") {
		t.Fatalf("did not expect q007 built prompt to include the old explicit airport instruction, got: %s", got)
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
		SemanticLayer: "Use `ontime.ontime` and `ontime.airports_latest`.",
	}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT Carrier, Flights FROM ontime.ontime", "# Report\n{{data_overview_md}}", "{\"sql\":\"SELECT Carrier, Flights FROM ontime.ontime\",\"report_markdown\":\"# Report\\n\\n{{data_overview_md}}\",\"metrics\":{\"named_values\":{\"max_hops\":\"8\"}}}")
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
