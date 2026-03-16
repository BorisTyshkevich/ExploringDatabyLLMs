package compare

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"qforge/internal/datasets"
	"qforge/internal/model"
	"qforge/internal/querylog"
	"qforge/internal/runs"
)

type RunMetrics struct {
	QueryDurationMS int64  `json:"query_duration_ms,omitempty"`
	ReadRows        int64  `json:"read_rows,omitempty"`
	ReadBytes       int64  `json:"read_bytes,omitempty"`
	ResultRows      int64  `json:"result_rows,omitempty"`
	ResultBytes     int64  `json:"result_bytes,omitempty"`
	MemoryUsage     int64  `json:"memory_usage,omitempty"`
	PeakThreads     int64  `json:"peak_threads,omitempty"`
	EventTime       string `json:"event_time,omitempty"`
	Type            string `json:"type,omitempty"`
}

type RunSummary struct {
	RunDir        string          `json:"run_dir"`
	QuestionID    string          `json:"question_id"`
	QuestionSlug  string          `json:"question_slug"`
	QuestionTitle string          `json:"question_title"`
	Dataset       string          `json:"dataset"`
	Runner        string          `json:"runner"`
	Model         string          `json:"model"`
	Status        model.RunStatus `json:"status"`
	Phases        model.RunPhases `json:"phases"`
	StartedAt     time.Time       `json:"started_at"`
	FinishedAt    time.Time       `json:"finished_at"`
	DurationSec   int64           `json:"duration_sec"`
	QuerySHA256   string          `json:"query_sha256,omitempty"`
	RowCount      int             `json:"row_count"`
	Columns       []string        `json:"columns,omitempty"`
	Metrics       *RunMetrics     `json:"metrics,omitempty"`
	Warnings      []string        `json:"warnings,omitempty"`
}

type Report struct {
	GeneratedAt string       `json:"generated_at"`
	Day         string       `json:"day"`
	Question    string       `json:"question,omitempty"`
	Warnings    []string     `json:"warnings,omitempty"`
	Runs        []RunSummary `json:"runs"`
}

type ArtifactPaths struct {
	Dir         string
	JSON        string
	PromptMD    string
	RawAnalysis string
	ReportMD    string
}

type resultSummary struct {
	Columns  []string `json:"columns"`
	RowCount int      `json:"row_count"`
}

func ArtifactPathsForQuestion(repoRoot, day, questionSlug string) ArtifactPaths {
	baseDir := filepath.Join(repoRoot, "runs", day, questionSlug)
	compareDir := filepath.Join(baseDir, "compare")
	return ArtifactPaths{
		Dir:         compareDir,
		JSON:        filepath.Join(compareDir, "compare.json"),
		PromptMD:    filepath.Join(compareDir, "analysis.prompt.md"),
		RawAnalysis: filepath.Join(compareDir, "analysis.raw.md"),
		ReportMD:    filepath.Join(baseDir, "compare_report.md"),
	}
}

func DiscoverQuestionRefs(repoRoot, day string) ([]string, error) {
	runDirs, err := discoverRunDirs(repoRoot, day, "")
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	var refs []string
	for _, runDir := range runDirs {
		manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
		if err != nil {
			return nil, err
		}
		ref := manifest.QuestionSlug
		if ref == "" {
			ref = manifest.QuestionID
		}
		if ref == "" {
			continue
		}
		if _, ok := seen[ref]; ok {
			continue
		}
		seen[ref] = struct{}{}
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	return refs, nil
}

func Generate(ctx context.Context, repoRoot, outDir, day, questionFilter, explicitMCPURL, explicitMCPToken string) (Report, error) {
	runDirs, err := discoverRunDirs(repoRoot, day, questionFilter)
	if err != nil {
		return Report{}, err
	}

	items := make([]RunSummary, 0, len(runDirs))
	var reportWarnings []string
	for _, runDir := range runDirs {
		item, warnings, err := summarizeRun(ctx, repoRoot, runDir, explicitMCPURL, explicitMCPToken)
		if err != nil {
			return Report{}, err
		}
		if len(warnings) > 0 {
			reportWarnings = append(reportWarnings, warnings...)
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Status != items[j].Status {
			return statusRank(items[i].Status) < statusRank(items[j].Status)
		}
		if items[i].Runner == items[j].Runner {
			return items[i].Model < items[j].Model
		}
		return items[i].Runner < items[j].Runner
	})

	report := Report{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Day:         day,
		Question:    questionFilter,
		Warnings:    dedupeStrings(reportWarnings),
		Runs:        items,
	}
	if len(items) == 0 {
		report.Warnings = append(report.Warnings, fmt.Sprintf("no runs found under runs/%s for %s", day, questionFilter))
	}
	if err := writeOutputs(outDir, report); err != nil {
		return Report{}, err
	}
	return report, nil
}

func summarizeRun(ctx context.Context, repoRoot, runDir, explicitMCPURL, explicitMCPToken string) (RunSummary, []string, error) {
	manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
	if err != nil {
		return RunSummary{}, nil, err
	}

	item := RunSummary{
		RunDir:        runDir,
		QuestionID:    manifest.QuestionID,
		QuestionSlug:  manifest.QuestionSlug,
		QuestionTitle: manifest.QuestionTitle,
		Dataset:       manifest.Dataset,
		Runner:        manifest.Runner,
		Model:         manifest.Model,
		Status:        manifest.Status,
		Phases:        manifest.Phases,
		StartedAt:     manifest.StartedAt,
		FinishedAt:    manifest.FinishedAt,
		DurationSec:   manifest.DurationSec,
		QuerySHA256:   manifest.QuerySHA256,
		RowCount:      manifest.ResultRowCount,
	}
	var warnings []string

	resultPath := filepath.Join(runDir, "result.json")
	resultBytes, err := os.ReadFile(resultPath)
	if err != nil {
		if os.IsNotExist(err) {
			warnings = append(warnings, fmt.Sprintf("%s: missing result.json", runID(item)))
		} else {
			warnings = append(warnings, fmt.Sprintf("%s: failed to read result.json: %v", runID(item), err))
		}
	} else {
		var summary resultSummary
		if err := json.Unmarshal(resultBytes, &summary); err != nil {
			warnings = append(warnings, fmt.Sprintf("%s: failed to parse result.json: %v", runID(item), err))
		} else {
			item.Columns = summary.Columns
			if item.RowCount == 0 {
				item.RowCount = summary.RowCount
			} else if summary.RowCount != 0 && summary.RowCount != item.RowCount {
				warnings = append(warnings, fmt.Sprintf("%s: manifest row count %d differs from result.json row count %d", runID(item), item.RowCount, summary.RowCount))
			}
		}
	}

	if manifest.LogComment != "" {
		cfg, err := datasets.Load(repoRoot, manifest.Dataset)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s: failed to load dataset config for metrics: %v", runID(item), err))
		} else {
			mcpURL, token, err := resolveMCPAccess(cfg, explicitMCPURL, explicitMCPToken)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("%s: failed to resolve MCP URL for metrics: %v", runID(item), err))
			} else {
				metrics, err := querylog.FetchLatest(ctx, mcpURL, token, manifest.LogComment)
				if err != nil {
					warnings = append(warnings, fmt.Sprintf("%s: failed to fetch query_log metrics: %v", runID(item), err))
				} else if metrics == nil {
					warnings = append(warnings, fmt.Sprintf("%s: query_log metrics not found", runID(item)))
				} else {
					item.Metrics = &RunMetrics{
						QueryDurationMS: metrics.QueryDurationMS,
						ReadRows:        metrics.ReadRows,
						ReadBytes:       metrics.ReadBytes,
						ResultRows:      metrics.ResultRows,
						ResultBytes:     metrics.ResultBytes,
						MemoryUsage:     metrics.MemoryUsage,
						PeakThreads:     metrics.PeakThreads,
						EventTime:       metrics.EventTime,
						Type:            metrics.Type,
					}
				}
			}
		}
	}

	item.Warnings = dedupeStrings(warnings)
	return item, item.Warnings, nil
}

func discoverRunDirs(repoRoot, day, questionFilter string) ([]string, error) {
	root := filepath.Join(repoRoot, "runs", day)
	var runDirs []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == "." {
			return nil
		}
		manifestPath := filepath.Join(path, "manifest.json")
		if _, err := os.Stat(manifestPath); err == nil {
			manifest, err := runs.ReadManifest(manifestPath)
			if err != nil {
				return err
			}
			if questionFilter == "" || manifest.QuestionID == questionFilter || manifest.QuestionSlug == questionFilter {
				runDirs = append(runDirs, path)
			}
			return filepath.SkipDir
		}
		return nil
	})
	if os.IsNotExist(err) {
		return nil, nil
	}
	return runDirs, err
}

func writeOutputs(outDir string, report Report) error {
	if outDir == "" {
		outDir = "runs/compare"
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "compare.json"), data, 0o644)
}

func renderMarkdown(report Report) string {
	var md strings.Builder
	md.WriteString("# qforge Compare Report\n\n")
	md.WriteString(fmt.Sprintf("- Generated: `%s`\n", report.GeneratedAt))
	md.WriteString(fmt.Sprintf("- Day: `%s`\n", report.Day))
	if report.Question != "" {
		md.WriteString(fmt.Sprintf("- Question: `%s`\n", report.Question))
	}
	if len(report.Warnings) > 0 {
		md.WriteString(fmt.Sprintf("- Warnings: `%d`\n", len(report.Warnings)))
	}
	md.WriteString("\n")

	if len(report.Runs) == 0 {
		md.WriteString("No runs found.\n")
		return md.String()
	}

	title := report.Question
	if title == "" && len(report.Runs) > 0 {
		title = report.Runs[0].QuestionID
	}
	if len(report.Runs) > 0 && report.Runs[0].QuestionTitle != "" {
		title = fmt.Sprintf("%s: %s", report.Runs[0].QuestionID, report.Runs[0].QuestionTitle)
	}
	md.WriteString("## ")
	md.WriteString(title)
	md.WriteString("\n\n")
	md.WriteString(renderQuestionSummary(report.Runs))
	md.WriteString("\n")
	md.WriteString("| runner | model | status | rows | duration | read rows | memory | warnings |\n")
	md.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | ---: |\n")
	for _, item := range report.Runs {
		duration := "n/a"
		readRows := "n/a"
		memory := "n/a"
		if item.Metrics != nil {
			duration = formatDurationMS(item.Metrics.QueryDurationMS)
			readRows = formatInt(item.Metrics.ReadRows)
			memory = formatBytes(item.Metrics.MemoryUsage)
		}
		md.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %s | %s | %s | %d |\n",
			item.Runner,
			item.Model,
			item.Status,
			item.RowCount,
			duration,
			readRows,
			memory,
			len(item.Warnings),
		))
	}
	md.WriteString("\n")
	if warnings := collectWarnings(report.Runs); len(warnings) > 0 {
		md.WriteString("### Warnings\n\n")
		for _, warning := range warnings {
			md.WriteString("- ")
			md.WriteString(warning)
			md.WriteString("\n")
		}
		md.WriteString("\n")
	}
	return md.String()
}

func renderQuestionSummary(items []RunSummary) string {
	var lines []string
	if failed := failedRuns(items); len(failed) == 0 {
		lines = append(lines, "- Status: all runs succeeded.")
	} else {
		lines = append(lines, fmt.Sprintf("- Status: %d run(s) did not finish cleanly: %s.", len(failed), strings.Join(failed, ", ")))
	}
	lines = append(lines, "- Row counts: "+rowCountSummary(items)+".")
	if fastest := bestByDuration(items); fastest != "" {
		lines = append(lines, "- Fastest successful run: "+fastest+".")
	}
	if leanest := bestByReadRows(items); leanest != "" {
		lines = append(lines, "- Lowest read volume: "+leanest+".")
	}
	if lowestMem := bestByMemory(items); lowestMem != "" {
		lines = append(lines, "- Lowest memory usage: "+lowestMem+".")
	}
	if warnings := countWarnings(items); warnings > 0 {
		lines = append(lines, fmt.Sprintf("- Warnings: %d.", warnings))
	}
	return strings.Join(lines, "\n")
}

func failedRuns(items []RunSummary) []string {
	var out []string
	for _, item := range items {
		if item.Status != model.RunStatusOK {
			out = append(out, runID(item))
		}
	}
	return out
}

func rowCountSummary(items []RunSummary) string {
	seen := map[int]struct{}{}
	var counts []int
	for _, item := range items {
		if _, ok := seen[item.RowCount]; ok {
			continue
		}
		seen[item.RowCount] = struct{}{}
		counts = append(counts, item.RowCount)
	}
	sort.Ints(counts)
	if len(counts) <= 1 {
		if len(counts) == 0 {
			return "no row counts available"
		}
		return fmt.Sprintf("all runs returned %d rows", counts[0])
	}
	parts := make([]string, len(counts))
	for i, count := range counts {
		parts[i] = strconv.Itoa(count)
	}
	return "mismatch (" + strings.Join(parts, ", ") + ")"
}

func bestByDuration(items []RunSummary) string {
	best, ok := bestMetricRun(items, func(item RunSummary) int64 {
		return item.Metrics.QueryDurationMS
	})
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s at %s", runID(best), formatDurationMS(best.Metrics.QueryDurationMS))
}

func bestByReadRows(items []RunSummary) string {
	best, ok := bestMetricRun(items, func(item RunSummary) int64 {
		return item.Metrics.ReadRows
	})
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s at %s rows", runID(best), formatInt(best.Metrics.ReadRows))
}

func bestByMemory(items []RunSummary) string {
	best, ok := bestMetricRun(items, func(item RunSummary) int64 {
		return item.Metrics.MemoryUsage
	})
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s at %s", runID(best), formatBytes(best.Metrics.MemoryUsage))
}

func bestMetricRun(items []RunSummary, metric func(RunSummary) int64) (RunSummary, bool) {
	var best RunSummary
	found := false
	for _, item := range items {
		if item.Status != model.RunStatusOK || item.Metrics == nil {
			continue
		}
		value := metric(item)
		if value <= 0 {
			continue
		}
		if !found || value < metric(best) {
			best = item
			found = true
		}
	}
	return best, found
}

func countWarnings(items []RunSummary) int {
	total := 0
	for _, item := range items {
		total += len(item.Warnings)
	}
	return total
}

func collectWarnings(items []RunSummary) []string {
	var all []string
	for _, item := range items {
		all = append(all, item.Warnings...)
	}
	return dedupeStrings(all)
}

func dedupeStrings(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func statusRank(status model.RunStatus) int {
	switch status {
	case model.RunStatusOK:
		return 0
	case model.RunStatusPartial:
		return 1
	case model.RunStatusAuthFailed:
		return 2
	case model.RunStatusFailed:
		return 3
	default:
		return 4
	}
}

func runID(item RunSummary) string {
	return item.Runner + "/" + item.Model
}

func formatDurationMS(ms int64) string {
	switch {
	case ms <= 0:
		return "n/a"
	case ms < 1000:
		return fmt.Sprintf("%d ms", ms)
	default:
		return fmt.Sprintf("%.2f s", float64(ms)/1000.0)
	}
}

func formatInt(v int64) string {
	if v == 0 {
		return "0"
	}
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}
	s := strconv.FormatInt(v, 10)
	if len(s) <= 3 {
		return sign + s
	}
	var out []byte
	prefix := len(s) % 3
	if prefix == 0 {
		prefix = 3
	}
	out = append(out, s[:prefix]...)
	for i := prefix; i < len(s); i += 3 {
		out = append(out, ',')
		out = append(out, s[i:i+3]...)
	}
	return sign + string(out)
}

func formatBytes(v int64) string {
	if v <= 0 {
		return "n/a"
	}
	const unit = 1024
	if v < unit {
		return fmt.Sprintf("%d B", v)
	}
	exp := int(math.Log(float64(v)) / math.Log(unit))
	if exp > 4 {
		exp = 4
	}
	value := float64(v) / math.Pow(unit, float64(exp))
	suffixes := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	return fmt.Sprintf("%.1f %s", value, suffixes[exp])
}

func resolveMCPAccess(cfg model.DatasetConfig, explicitURL, explicitToken string) (string, string, error) {
	if explicitToken != "" && explicitURL == "" {
		baseURL := cfg.MCPBaseURL
		if baseURL == "" {
			baseURL = "https://mcp.demo.altinity.cloud"
		}
		return fmt.Sprintf("%s/%s/http", strings.TrimRight(baseURL, "/"), explicitToken), explicitToken, nil
	}
	url, token, err := datasets.ResolveMCPURL(cfg, explicitURL)
	if err != nil {
		return "", "", err
	}
	if explicitToken != "" {
		token = explicitToken
	}
	return url, token, nil
}
