package model

import "time"

type Phase string

const (
	PhaseSQL          Phase = "sql"
	PhasePresentation Phase = "presentation"
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
	ID                string `yaml:"id" json:"id"`
	Slug              string `yaml:"slug" json:"slug"`
	Title             string `yaml:"title" json:"title"`
	Dataset           string `yaml:"dataset" json:"dataset"`
	ArtifactsRequired string `yaml:"artifacts_required" json:"artifacts_required"`
	VisualMode        string `yaml:"visual_mode" json:"visual_mode"`
	VisualType        string `yaml:"visual_type" json:"visual_type"`
	Tags              string `yaml:"tags" json:"tags"`
	ReferencePolicy   string `yaml:"reference_policy" json:"reference_policy"`
	CommandTimeoutSec int    `yaml:"command_timeout_sec" json:"command_timeout_sec"`
}

type Question struct {
	Dir                 string       `json:"dir"`
	Meta                QuestionMeta `json:"meta"`
	Prompt              string       `json:"prompt"`
	ReportPrompt        string       `json:"report_prompt"`
	VisualPrompt        string       `json:"visual_prompt"`
	PresentationEnabled bool         `json:"presentation_enabled"`
	ReportEnabled       bool         `json:"report_enabled"`
	VisualEnabled       bool         `json:"visual_enabled"`
}

type ArtifactPaths struct {
	PromptSQLRaw          string `json:"prompt_sql_raw"`
	AnswerSQLRaw          string `json:"answer_sql_raw"`
	QuerySQL              string `json:"query_sql"`
	ResultTSV             string `json:"result_tsv,omitempty"`
	ResultJSON            string `json:"result_json"`
	ManifestJSON          string `json:"manifest_json"`
	StdoutLog             string `json:"stdout_log"`
	StderrLog             string `json:"stderr_log"`
	PromptPresentationRaw string `json:"prompt_presentation_raw,omitempty"`
	AnswerPresentationRaw string `json:"answer_presentation_raw,omitempty"`
	ReportTemplateMD      string `json:"report_template_md,omitempty"`
	ReportMD              string `json:"report_md,omitempty"`
	VisualHTML            string `json:"visual_html,omitempty"`
}

type RunPhases struct {
	SQLGeneration          PhaseStatus `json:"sql_generation"`
	SQLExecution           PhaseStatus `json:"sql_execution"`
	PresentationGeneration PhaseStatus `json:"presentation_generation"`
	PresentationRender     PhaseStatus `json:"presentation_render"`
}

type RunManifest struct {
	SchemaVersion   string            `json:"schema_version"`
	Status          RunStatus         `json:"status"`
	QuestionID      string            `json:"question_id"`
	QuestionSlug    string            `json:"question_slug"`
	QuestionTitle   string            `json:"question_title"`
	Dataset         string            `json:"dataset"`
	Runner          string            `json:"runner"`
	Model           string            `json:"model"`
	CLIBin          string            `json:"cli_bin"`
	MCPServerName   string            `json:"mcp_server_name"`
	MCPConfigSource string            `json:"mcp_config_source"`
	StartedAt       time.Time         `json:"started_at"`
	FinishedAt      time.Time         `json:"finished_at"`
	DurationSec     int64             `json:"duration_sec"`
	LogComment      string            `json:"log_comment"`
	QuerySHA256     string            `json:"query_sha256"`
	ResultRowCount  int               `json:"result_row_count"`
	Phases          RunPhases         `json:"phases"`
	Artifacts       ArtifactPaths     `json:"artifacts"`
	Metadata        map[string]string `json:"metadata,omitempty"`
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

type NumericColumnSpec struct {
	Name         string  `yaml:"name" json:"name"`
	ToleranceAbs float64 `yaml:"tolerance_abs" json:"tolerance_abs"`
}

type RowFilter struct {
	IncludeRowTypes []string       `yaml:"include_row_types" json:"include_row_types"`
	Where           map[string]any `yaml:"where" json:"where"`
}

type NormalizationSpec struct {
	TrimStrings     bool     `yaml:"trim_strings" json:"trim_strings"`
	NullEquivalents []string `yaml:"null_equivalents" json:"null_equivalents"`
}

type ComplianceSpec struct {
	RequireNonemptyRows   bool `yaml:"require_nonempty_rows" json:"require_nonempty_rows"`
	RequireUniqueKeys     bool `yaml:"require_unique_keys" json:"require_unique_keys"`
	RequireReferenceMatch bool `yaml:"require_reference_match" json:"require_reference_match"`
}

type CompareColumns struct {
	Exact   []string            `yaml:"exact" json:"exact"`
	Numeric []NumericColumnSpec `yaml:"numeric" json:"numeric"`
}

type CompareContract struct {
	Version         int                 `yaml:"version" json:"version"`
	RowFilter       RowFilter           `yaml:"row_filter" json:"row_filter"`
	KeyColumns      []string            `yaml:"key_columns" json:"key_columns"`
	CompareColumns  CompareColumns      `yaml:"compare_columns" json:"compare_columns"`
	OptionalColumns []string            `yaml:"optional_columns" json:"optional_columns"`
	HeaderAliases   map[string][]string `yaml:"header_aliases" json:"header_aliases"`
	Normalization   NormalizationSpec   `yaml:"normalization" json:"normalization"`
	Compliance      ComplianceSpec      `yaml:"compliance" json:"compliance"`
}

type CompareContractFile struct {
	CompareContract CompareContract `yaml:"compare_contract" json:"compare_contract"`
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
