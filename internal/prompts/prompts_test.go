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
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildSQLPrompt(question, dataset, model.AnalysisModeJSONArtifact)
	if err != nil {
		t.Fatalf("BuildSQLPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Use the `ontime` database to answer analytical questions") {
		t.Fatalf("expected shared core scaffold, got: %s", got)
	}
	if !strings.Contains(got, "Use `ontime-semantic-layer` skill for schema inspection") {
		t.Fatalf("expected semantic skill guidance, got: %s", got)
	}
	if !strings.Contains(got, "Question-specific SQL guidance.") {
		t.Fatalf("expected question prompt section, got: %s", got)
	}
	if strings.Contains(got, "Dataset semantic layer:") {
		t.Fatalf("did not expect inlined semantic layer guidance, got: %s", got)
	}
	if !strings.Contains(got, "answer.raw.json") || !strings.Contains(got, "\"sql\"") || !strings.Contains(got, "\"report_markdown\"") || !strings.Contains(got, "\"metrics\"") {
		t.Fatalf("expected file-based single json analysis contract, got: %s", got)
	}
	if !strings.Contains(got, "raw JSON, not fenced Markdown") || !strings.Contains(got, "Do not emit result rows") {
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

func TestBuildSQLPromptTemplateModeUsesDirectFileContract(t *testing.T) {
	question := model.Question{
		Dir:    filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Prompt: "Question-specific SQL guidance.",
	}
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildSQLPrompt(question, dataset, model.AnalysisModeTemplateFiles)
	if err != nil {
		t.Fatalf("BuildSQLPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Write the final verified SQL to `query.sql`.") {
		t.Fatalf("expected query.sql contract, got: %s", got)
	}
	if !strings.Contains(got, "Write the Markdown report template to `report.template.md`.") {
		t.Fatalf("expected report.template.md contract, got: %s", got)
	}
	if strings.Contains(got, "Write one JSON object containing the final verified SQL") {
		t.Fatalf("did not expect answer.raw.json contract in template mode, got: %s", got)
	}
	if !strings.Contains(got, "Allowed built-in placeholders:") || !strings.Contains(got, "Do not invent any placeholder outside the built-in list.") {
		t.Fatalf("expected placeholder constraints in template mode, got: %s", got)
	}
}

func TestBuildSQLPromptMultiQueryModeUsesStructuredJSONContract(t *testing.T) {
	question := model.Question{
		Dir:    filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Prompt: "Question-specific SQL guidance.\n\n## Dashboard Questions\n\n- Which hotspot is worst?\n- Is it persistent?",
		Subquestions: []model.QuestionSubquestion{
			{ID: "worst_hotspot", Text: "Which hotspot is worst?"},
			{ID: "persistence", Text: "Is it persistent?"},
		},
	}
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildSQLPrompt(question, dataset, model.AnalysisModeMultiQueryJSON)
	if err != nil {
		t.Fatalf("BuildSQLPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "\"subquestions\"") || !strings.Contains(got, "\"answer_markdown\"") {
		t.Fatalf("expected multi-query json artifact contract, got: %s", got)
	}
	if !strings.Contains(got, "Read every bullet under `## Dashboard Questions`") {
		t.Fatalf("expected dashboard-question instruction in prompt, got: %s", got)
	}
	if !strings.Contains(got, "## Dashboard Questions") || !strings.Contains(got, "- Which hotspot is worst?") || !strings.Contains(got, "- Is it persistent?") {
		t.Fatalf("expected dashboard questions to be present in prompt, got: %s", got)
	}
}

func TestBuildSQLPromptMultiQueryModeAppendsDashboardQuestionsFromSidecarContract(t *testing.T) {
	question := model.Question{
		Dir:    filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta:   model.QuestionMeta{AnalysisMode: string(model.AnalysisModeMultiQueryJSON)},
		Prompt: "Question-specific SQL guidance.",
		Subquestions: []model.QuestionSubquestion{
			{ID: "worst_hotspot", Text: "Which hotspot is worst?"},
			{ID: "persistence", Text: "Is it persistent?"},
		},
	}
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildSQLPrompt(question, dataset, model.AnalysisModeMultiQueryJSON)
	if err != nil {
		t.Fatalf("BuildSQLPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "## Dashboard Questions") || !strings.Contains(got, "- Which hotspot is worst?") || !strings.Contains(got, "- Is it persistent?") {
		t.Fatalf("expected injected dashboard questions in prompt, got: %s", got)
	}
}

func TestBuildVisualPromptMultiQueryModeUsesDynamicDashboardContract(t *testing.T) {
	question := model.Question{
		Dir:          filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta:         model.QuestionMeta{ID: "q003", Title: "Delta ATL", VisualMode: "dynamic", VisualType: "html_heatmap", AnalysisMode: string(model.AnalysisModeMultiQueryJSON)},
		VisualPrompt: "Visual guidance.",
	}
	visualInput := model.VisualInputSummary{
		QuestionTitle: "Delta ATL",
		RowCount:      3,
		QuerySummaries: []model.QueryResultSummary{
			{ID: "worst_hotspot", SQL: "SELECT * FROM hotspots"},
			{ID: "persistence", SQL: "SELECT * FROM persistence"},
		},
		ModeHint: "First pass produced named proof queries.",
	}
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildVisualPrompt(question, dataset, model.CanonicalResult{}, "SELECT * FROM hotspots", "https://mcp.example.invalid/{JWE}/openapi/execute_query?query=...", visualInput)
	if err != nil {
		t.Fatalf("BuildVisualPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Use the `ontime` database to answer analytical questions") {
		t.Fatalf("expected shared core scaffold in multi-query visual prompt, got: %s", got)
	}
	if !strings.Contains(got, "The saved SQL shown below is the primary dashboard query for this page.") {
		t.Fatalf("expected multi-query visual supplement, got: %s", got)
	}
	if !strings.Contains(got, "Use this endpoint template for every browser query") || !strings.Contains(got, "OnTimeAnalystDashboard::auth::jwe") {
		t.Fatalf("expected dynamic-mode contract in multi-query visual prompt, got: %s", got)
	}
	if !strings.Contains(got, "SQL query for primary data source:") || !strings.Contains(got, "SELECT * FROM hotspots") {
		t.Fatalf("expected primary saved SQL in multi-query visual prompt, got: %s", got)
	}
	if strings.Contains(got, "Do not query ClickHouse") || strings.Contains(got, "Embed the verified analysis package") {
		t.Fatalf("did not expect legacy static-only multi-query guidance, got: %s", got)
	}
}

func TestBuildPresentationPromptLoadsMarkdownAssets(t *testing.T) {
	question := model.Question{
		Dir:          filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta:         model.QuestionMeta{ID: "q003", Title: "Delta ATL", VisualMode: "dynamic", VisualType: "html_heatmap"},
		VisualPrompt: "Visual guidance.",
	}
	result := model.CanonicalResult{
		Columns:     []string{"RowType", "DestCode"},
		GeneratedAt: time.Now(),
	}
	visualInput := model.VisualInputSummary{
		QuestionTitle: "Delta ATL",
		ResultColumns: []string{"RowType", "DestCode"},
		RowCount:      2,
		SampleRows: []map[string]any{
			{"RowType": "summary", "DestCode": "LAX"},
		},
		FieldShapeNotes: map[string]string{"FlightDate": "ISO-like timestamp string"},
		ModeHint:        "Dynamic mode still fetches live data in the browser via query.sql and the configured endpoint.",
	}
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT *\nFROM ontime.fact_ontime", "https://mcp.example.invalid/{JWE}/openapi/execute_query?query=...", visualInput)
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Use the `ontime` database to answer analytical questions") {
		t.Fatalf("expected shared core scaffold, got: %s", got)
	}
	if !strings.Contains(got, "Create browser-ready HTML `visual.html`") {
		t.Fatalf("expected merged visual scaffold, got: %s", got)
	}
	if !strings.Contains(got, "Write the file or provide a download link.") {
		t.Fatalf("expected file-writing visual contract, got: %s", got)
	}
	if !strings.Contains(got, "*-analyst-dashboard") || !strings.Contains(got, "Use `ontime-semantic-layer` skill for schema inspection") {
		t.Fatalf("expected dashboard and semantic skill guidance, got: %s", got)
	}
	if !strings.Contains(got, "Visual mode: `dynamic`") {
		t.Fatalf("expected visual mode in prompt, got: %s", got)
	}
	if !strings.Contains(got, "https://mcp.example.invalid/{JWE}/openapi/execute_query?query=...") || !strings.Contains(got, "OnTimeAnalystDashboard::auth::jwe") {
		t.Fatalf("expected dynamic endpoint/auth contract, got: %s", got)
	}
	if !strings.Contains(got, "Do not embed the primary analytical dataset") {
		t.Fatalf("expected no-embedded-dataset contract, got: %s", got)
	}
	if !strings.Contains(got, "Data example/snippet:") || !strings.Contains(got, "\"sample_rows\"") {
		t.Fatalf("expected visual input summary context, got: %s", got)
	}
	if !strings.Contains(got, "SELECT *") || !strings.Contains(got, "FROM ontime.fact_ontime") {
		t.Fatalf("expected saved sql to be embedded in prompt, got: %s", got)
	}
	if !strings.Contains(got, "Do not include the HTML source in the response.") || strings.Contains(got, "```html") {
		t.Fatalf("expected write-file output contract in presentation prompt, got: %s", got)
	}
	if strings.Contains(got, "qforge-result-data") || strings.Contains(got, "__QFORGE_DEFAULT_SQL__") {
		t.Fatalf("did not expect legacy injected JSON contract, got: %s", got)
	}
	if strings.Contains(got, "saved report template") || strings.Contains(strings.ToLower(got), "saved report template") {
		t.Fatalf("did not expect saved report template context in visual prompt, got: %s", got)
	}
	if strings.Contains(got, "Saved analysis artifact:") || strings.Contains(got, "\"report_markdown\"") {
		t.Fatalf("did not expect saved analysis json context in visual prompt, got: %s", got)
	}
	if strings.Contains(got, "Dataset semantic layer:") {
		t.Fatalf("did not expect inlined semantic-layer guidance, got: %s", got)
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
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT 1", "https://mcp.example.invalid/{JWE}/openapi/execute_query?query=...", model.VisualInputSummary{})
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "*-analyst-dashboard") || !strings.Contains(got, "Use `ontime-semantic-layer` skill for schema inspection") {
		t.Fatalf("expected q001 prompt to reference dashboard and semantic skill guidance, got: %s", got)
	}
	if !strings.Contains(got, "airport-coordinate enrichment") {
		t.Fatalf("expected q001 prompt to label map enrichment clearly, got: %s", got)
	}
	if !strings.Contains(got, "query ledger") {
		t.Fatalf("expected q001 prompt to inherit ledger contract, got: %s", got)
	}
	if !strings.Contains(got, "keep the map card visible with degraded-state messaging") {
		t.Fatalf("expected q001 prompt to require degraded map state, got: %s", got)
	}
	if strings.Contains(got, "Dataset semantic layer:") {
		t.Fatalf("did not expect semantic layer heading when using skill references, got: %s", got)
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
		VisualPrompt: "Visual guidance.",
	}
	result := model.CanonicalResult{
		Columns:     []string{"Carrier", "Flights"},
		GeneratedAt: time.Now(),
	}
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT Carrier, Flights FROM ontime.fact_ontime", "", model.VisualInputSummary{
		QuestionTitle: "Static Fixture",
		ResultColumns: []string{"Carrier", "Flights"},
		RowCount:      1,
		SampleRows:    []map[string]any{{"Carrier": "DL", "Flights": 42}},
		ModeHint:      "Static mode embeds analytical data from result.json directly in the page.",
	})
	if err != nil {
		t.Fatalf("BuildPresentationPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Visual mode: `static`") {
		t.Fatalf("expected static visual mode in prompt, got: %s", got)
	}
	if !strings.Contains(got, "Build a self-contained benchmark artifact") {
		t.Fatalf("expected static artifact contract, got: %s", got)
	}
	if !strings.Contains(got, "Embed or bake the analytical data needed by the page into the final browser artifact") {
		t.Fatalf("expected embedded-data contract, got: %s", got)
	}
	if !strings.Contains(got, "Data example/snippet:") {
		t.Fatalf("expected visual input summary in static prompt, got: %s", got)
	}
	if !strings.Contains(got, "Use `ontime-semantic-layer` skill for schema inspection") {
		t.Fatalf("expected semantic skill guidance in static prompt, got: %s", got)
	}
	if strings.Contains(got, "OnTimeAnalystDashboard::auth::jwe") {
		t.Fatalf("did not expect dynamic JWE contract in static prompt, got: %s", got)
	}
}

func TestBuildPresentationPromptReactDynamicUsesSourceContract(t *testing.T) {
	question := model.Question{
		Dir: filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta: model.QuestionMeta{
			ID:                 "q901",
			Title:              "React Dynamic Fixture",
			VisualMode:         "dynamic",
			PresentationTarget: "react",
			VisualType:         "html_heatmap",
		},
		VisualPrompt: "Visual guidance.",
	}
	result := model.CanonicalResult{
		Columns:     []string{"Carrier", "Flights"},
		GeneratedAt: time.Now(),
	}
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT 1", "https://mcp.example.invalid/{JWE}/openapi/execute_query?query=...", model.VisualInputSummary{})
	if err != nil {
		t.Fatalf("BuildVisualPrompt returned error: %v", err)
	}
	for _, want := range []string{
		"Presentation target: `react`",
		"Create a React source artifact under `visual_src/`",
		"`visual_src/package.json`",
		"`visual_src/src/main.jsx`",
		"Do not emit the source code inline in the response",
		"Use this endpoint template for every browser query",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected react dynamic prompt to contain %q, got: %s", want, got)
		}
	}
	if strings.Contains(got, "Create browser-ready HTML `visual.html`") {
		t.Fatalf("did not expect html-only contract in react prompt, got: %s", got)
	}
}

func TestBuildPresentationPromptReactStaticAvoidsDynamicTokenFlow(t *testing.T) {
	question := model.Question{
		Dir: filepath.Join("..", "..", "prompts", "q003_delta_atl_departure_delay_hotspots"),
		Meta: model.QuestionMeta{
			ID:                 "q902",
			Title:              "React Static Fixture",
			VisualMode:         "static",
			PresentationTarget: "react",
			VisualType:         "html_ranked_dashboard",
		},
		VisualPrompt: "Visual guidance.",
	}
	result := model.CanonicalResult{
		Columns:     []string{"Carrier", "Flights"},
		GeneratedAt: time.Now(),
	}
	dataset := model.DatasetConfig{DefaultDatabase: "ontime"}
	got, err := BuildVisualPrompt(question, dataset, result, "SELECT 1", "", model.VisualInputSummary{})
	if err != nil {
		t.Fatalf("BuildVisualPrompt returned error: %v", err)
	}
	if !strings.Contains(got, "Presentation target: `react`") || !strings.Contains(got, "Create a React source artifact under `visual_src/`") {
		t.Fatalf("expected react static source contract, got: %s", got)
	}
	if !strings.Contains(got, "Build a self-contained benchmark artifact") {
		t.Fatalf("expected static runtime guidance, got: %s", got)
	}
	if strings.Contains(got, "OnTimeAnalystDashboard::auth::jwe") {
		t.Fatalf("did not expect dynamic token flow in static react prompt, got: %s", got)
	}
}
