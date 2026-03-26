package model

import "time"

type Phase string

const (
	PhaseSQL          Phase = "sql"
	PhasePresentation Phase = "presentation"
)

type AnalysisMode string

const (
	AnalysisModeMultiQuery    AnalysisMode = "multi_query"
	AnalysisModeTemplateFiles AnalysisMode = "template_files"
)

type RunStatus string

const (
	RunStatusOK         RunStatus = "ok"
	RunStatusPartial    RunStatus = "partial"
	RunStatusFailed     RunStatus = "failed"
	RunStatusAuthFailed RunStatus = "auth_failed"
)

type PhaseStatus string

const (
	PhaseStatusNotRun  PhaseStatus = "not_run"
	PhaseStatusOK      PhaseStatus = "ok"
	PhaseStatusFailed  PhaseStatus = "failed"
	PhaseStatusSkipped PhaseStatus = "skipped"
)

type DatasetConfig struct {
	Name                 string `yaml:"dataset" json:"name"`
	DefaultMCPServerName string `yaml:"default_mcp_server_name" json:"default_mcp_server_name"`
	MCPURL               string `yaml:"mcp_url" json:"mcp_url"`
	MCPBaseURL           string `yaml:"mcp_base_url" json:"mcp_base_url"`
	MCPJWETokenEnv       string `yaml:"mcp_jwe_token_env" json:"mcp_jwe_token_env"`
	AuthMode             string `yaml:"auth_mode" json:"auth_mode"`
	DefaultDatabase      string `yaml:"default_database" json:"default_database"`
	Notes                string `yaml:"notes" json:"notes"`
	SemanticLayer        string `yaml:"-" json:"semantic_layer"`
}

type QuestionMeta struct {
	ID                 string `yaml:"id" json:"id"`
	Slug               string `yaml:"slug" json:"slug"`
	Title              string `yaml:"title" json:"title"`
	Dataset            string `yaml:"dataset" json:"dataset"`
	AnalysisMode       string `yaml:"analysis_mode" json:"analysis_mode"`
	ArtifactsRequired  string `yaml:"artifacts_required" json:"artifacts_required"`
	VisualMode         string `yaml:"visual_mode" json:"visual_mode"`
	PresentationTarget string `yaml:"presentation_target" json:"presentation_target"`
	VisualType         string `yaml:"visual_type" json:"visual_type"`
	Tags               string `yaml:"tags" json:"tags"`
	ReferencePolicy    string `yaml:"reference_policy" json:"reference_policy"`
	CommandTimeoutSec  int    `yaml:"command_timeout_sec" json:"command_timeout_sec"`
}

type QuestionSubquestion struct {
	ID   string `yaml:"id" json:"id"`
	Text string `yaml:"text" json:"text"`
}

type Question struct {
	Dir                 string                `json:"dir"`
	Meta                QuestionMeta          `json:"meta"`
	Prompt              string                `json:"prompt"`
	VisualPrompt        string                `json:"visual_prompt"`
	Subquestions        []QuestionSubquestion `json:"subquestions,omitempty"`
	PresentationEnabled bool                  `json:"presentation_enabled"`
	ReportEnabled       bool                  `json:"report_enabled"`
	VisualEnabled       bool                  `json:"visual_enabled"`
}

type ArtifactPaths struct {
	PromptReportRaw       string `json:"prompt_report_raw"`
	AnswerReportRaw       string `json:"answer_report_raw"`
	AnswerRawJSON         string `json:"answer_raw_json,omitempty"`
	AnalysisJSON          string `json:"analysis_json,omitempty"`
	PromptReviewRaw       string `json:"prompt_review_raw,omitempty"`
	AnswerReviewRaw       string `json:"answer_review_raw,omitempty"`
	ReviewMD              string `json:"review_md,omitempty"`
	QuerySQL              string `json:"query_sql"`
	ResultTSV             string `json:"result_tsv,omitempty"`
	ResultJSON            string `json:"result_json"`
	VisualInputJSON       string `json:"visual_input_json,omitempty"`
	ManifestJSON          string `json:"manifest_json"`
	StdoutLog             string `json:"stdout_log"`
	StderrLog             string `json:"stderr_log"`
	PromptPresentationRaw string `json:"prompt_presentation_raw,omitempty"`
	AnswerPresentationRaw string `json:"answer_presentation_raw,omitempty"`
	ReportTemplateMD      string `json:"report_template_md,omitempty"`
	ReportMD              string `json:"report_md,omitempty"`
	VisualHTML            string `json:"visual_html,omitempty"`
	VisualSourceDir       string `json:"visual_source_dir,omitempty"`
	VisualBuildDir        string `json:"visual_build_dir,omitempty"`
	VisualAssetsDir       string `json:"visual_assets_dir,omitempty"`
	VisualPackageJSON     string `json:"visual_package_json,omitempty"`
}

type RunPhases struct {
	SQLGeneration          PhaseStatus `json:"sql_generation"`
	SQLExecution           PhaseStatus `json:"sql_execution"`
	Review                 PhaseStatus `json:"review"`
	PresentationGeneration PhaseStatus `json:"presentation_generation"`
	PresentationRender     PhaseStatus `json:"presentation_render"`
}

type RunManifest struct {
	SchemaVersion                   string            `json:"schema_version"`
	Status                          RunStatus         `json:"status"`
	QuestionID                      string            `json:"question_id"`
	QuestionSlug                    string            `json:"question_slug"`
	QuestionTitle                   string            `json:"question_title"`
	Dataset                         string            `json:"dataset"`
	Runner                          string            `json:"runner"`
	Model                           string            `json:"model"`
	AnalysisMode                    string            `json:"analysis_mode,omitempty"`
	PresentationTarget              string            `json:"presentation_target,omitempty"`
	ReviewRunner                    string            `json:"review_runner,omitempty"`
	ReviewModel                     string            `json:"review_model,omitempty"`
	ReviewVerdict                   string            `json:"review_verdict,omitempty"`
	CLIBin                          string            `json:"cli_bin"`
	MCPServerName                   string            `json:"mcp_server_name"`
	MCPConfigSource                 string            `json:"mcp_config_source"`
	StartedAt                       time.Time         `json:"started_at"`
	FinishedAt                      time.Time         `json:"finished_at"`
	DurationSec                     int64             `json:"duration_sec"`
	SQLGenerationProviderDurationMs int64             `json:"sql_generation_provider_duration_ms,omitempty"`
	PresentationProviderDurationMs  int64             `json:"presentation_provider_duration_ms,omitempty"`
	PresentationBuildDurationMs     int64             `json:"presentation_build_duration_ms,omitempty"`
	LogComment                      string            `json:"log_comment"`
	QuerySHA256                     string            `json:"query_sha256"`
	ResultRowCount                  int               `json:"result_row_count"`
	Phases                          RunPhases         `json:"phases"`
	Artifacts                       ArtifactPaths     `json:"artifacts"`
	Metadata                        map[string]string `json:"metadata,omitempty"`
}

type CanonicalResult struct {
	Columns           []string         `json:"columns"`
	Rows              []map[string]any `json:"rows"`
	RowCount          int              `json:"row_count"`
	GeneratedAt       time.Time        `json:"generated_at"`
	SourceQuerySHA256 string           `json:"source_query_sha256"`
	LogComment        string           `json:"log_comment"`
}

type ProviderRequest struct {
	Question      Question
	Dataset       DatasetConfig
	Prompt        string
	OutDir        string
	Model         string
	AnalysisMode  string
	MCPURL        string
	MCPServerName string
	MCPToken      string
	CLIBin        string
	Verbose       bool
}

type ProviderResponse struct {
	RawOutput string
	Stdout    string
	Stderr    string
	CLIBin    string
}

type AnalysisArtifact struct {
	SQL            string                `json:"sql"`
	ReportMarkdown string                `json:"report_markdown"`
	Subquestions   []AnalysisSubquestion `json:"subquestions,omitempty"`
}

type AnalysisSubquestion struct {
	ID             string `json:"id,omitempty"`
	Subquestion    string `json:"subquestion,omitempty"`
	AnswerMarkdown string `json:"answer_markdown"`
	SQL            string `json:"sql"`
}

type QueryResultSummary struct {
	ID             string         `json:"id"`
	Subquestion    string         `json:"subquestion"`
	AnswerMarkdown string         `json:"answer_markdown"`
	SQL            string         `json:"sql"`
	RowCount       int            `json:"row_count"`
	ResultColumns  []string       `json:"result_columns"`
	FirstRow       map[string]any `json:"first_row,omitempty"`
}

type VisualInputSummary struct {
	QuestionTitle   string               `json:"question_title"`
	ResultColumns   []string             `json:"result_columns"`
	RowCount        int                  `json:"row_count"`
	SampleRows      []map[string]any     `json:"sample_rows,omitempty"`
	FieldShapeNotes map[string]string    `json:"field_shape_notes,omitempty"`
	ModeHint        string               `json:"mode_hint,omitempty"`
	QuerySummaries  []QueryResultSummary `json:"query_summaries,omitempty"`
}

type NumericColumnSpec struct {
	Name         string  `yaml:"name" json:"name"`
	ToleranceAbs float64 `yaml:"tolerance_abs" json:"tolerance_abs"`
}

type QueryLogMetrics struct {
	LogComment      string `json:"log_comment"`
	QueryID         string `json:"query_id"`
	QueryDurationMS int64  `json:"query_duration_ms"`
	ReadRows        int64  `json:"read_rows"`
	ReadBytes       int64  `json:"read_bytes"`
	ResultRows      int64  `json:"result_rows"`
	ResultBytes     int64  `json:"result_bytes"`
	MemoryUsage     int64  `json:"memory_usage"`
	PeakThreads     int64  `json:"peak_threads"`
	Query           string `json:"query"`
	EventTime       string `json:"event_time"`
	Type            string `json:"type"`
}
