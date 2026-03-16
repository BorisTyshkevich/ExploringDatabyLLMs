package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"

	"qforge/internal/compare"
	"qforge/internal/datasets"
	"qforge/internal/execute"
	"qforge/internal/extract"
	"qforge/internal/model"
	"qforge/internal/prompts"
	"qforge/internal/providers"
	"qforge/internal/questions"
	"qforge/internal/render"
	"qforge/internal/runs"
)

const defaultCommandTimeoutSec = 900

func Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		printRootUsage(os.Stdout)
		return nil
	}
	switch args[0] {
	case "help", "-h", "--help":
		printRootUsage(os.Stdout)
		return nil
	}
	switch args[0] {
	case "list-questions":
		return runListQuestions(args[1:])
	case "run":
		return runRun(ctx, args[1:])
	case "process-visual":
		return runProcessVisual(ctx, args[1:])
	case "compare":
		return runCompare(ctx, args[1:])
	case "inspect-run":
		return runInspectRun(args[1:])
	default:
		return usageError()
	}
}

func usageError() error {
	return errors.New("usage: qforge <run|process-visual|compare|list-questions|inspect-run> ...")
}

func printRootUsage(out *os.File) {
	fmt.Fprintln(out, "qforge orchestrates model-generated SQL, harness-owned execution, and compare reporting.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out, "  qforge <command> [flags]")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Commands:")
	fmt.Fprintln(out, "  list-questions   List available benchmark questions")
	fmt.Fprintln(out, "  run              Run one question for one or more providers")
	fmt.Fprintln(out, "  process-visual   Generate report/html for an existing run directory")
	fmt.Fprintln(out, "  compare          Compare runs and fetch query_log metrics")
	fmt.Fprintln(out, "  inspect-run      Print one run manifest")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Common workflow:")
	fmt.Fprintln(out, "  1. list questions")
	fmt.Fprintln(out, "  2. run one question for one or more providers")
	fmt.Fprintln(out, "  3. compare runs and fetch deferred performance metrics")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Examples:")
	fmt.Fprintln(out, "  qforge list-questions")
	fmt.Fprintln(out, "  qforge run --question q001 --runner claude --verbose")
	fmt.Fprintln(out, "  qforge run --question q001 --runner claude --with-visual --verbose")
	fmt.Fprintln(out, "  qforge run --question q001 --runner codex --runner claude --verbose")
	fmt.Fprintln(out, "  qforge run --question q001 --verbose")
	fmt.Fprintln(out, "  qforge process-visual --run-dir runs/2026-03-15/q001_hops_per_day/claude/opus/run-004 --verbose")
	fmt.Fprintln(out, "  qforge compare --day 2026-03-15 --question q003 --runner codex --verbose")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Use `qforge <command> --help` for detailed subcommand help.")
}

func repoRoot() (string, error) {
	return os.Getwd()
}

func runListQuestions(args []string) error {
	fs := flag.NewFlagSet("list-questions", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Fprintln(os.Stdout, "Usage: qforge list-questions")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "List question metadata discovered under questions/.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Prints tab-separated columns:")
		fmt.Fprintln(os.Stdout, "  question_id, question_slug, title, dataset, presentation_enabled")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Example:")
		fmt.Fprintln(os.Stdout, "  qforge list-questions")
	}
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	root, err := repoRoot()
	if err != nil {
		return err
	}
	items, err := questions.LoadAll(root)
	if err != nil {
		return err
	}
	for _, item := range items {
		fmt.Printf("%s\t%s\t%s\t%s\t%t\n", item.Meta.ID, item.Meta.Slug, item.Meta.Title, item.Meta.Dataset, item.PresentationEnabled)
	}
	return nil
}

func runRun(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Fprintln(os.Stdout, "Usage: qforge run --question <id|slug> [--runner <codex|claude|gemini> ...] [flags]")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Run one question through one or more providers.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Behavior:")
		fmt.Fprintln(os.Stdout, "  - loads question and dataset metadata")
		fmt.Fprintln(os.Stdout, "  - prompts each selected provider for final SQL only")
		fmt.Fprintln(os.Stdout, "  - enforces forbidden-table policy")
		fmt.Fprintln(os.Stdout, "  - executes SQL itself and writes result.json")
		fmt.Fprintln(os.Stdout, "  - runs providers concurrently when more than one is selected")
		fmt.Fprintln(os.Stdout, "  - optionally performs a separate follow-up provider call for report/html")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Important:")
		fmt.Fprintln(os.Stdout, "  Presentation is handled separately by `qforge process-visual`, or by `--with-visual`.")
		fmt.Fprintln(os.Stdout, "  `--with-visual` makes a second independent provider call after SQL execution succeeds.")
		fmt.Fprintln(os.Stdout, "  If --runner is omitted, qforge runs codex, claude, and gemini.")
		fmt.Fprintln(os.Stdout, "  Repeated --model flags are matched positionally to repeated --runner flags.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Flags:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Examples:")
		fmt.Fprintln(os.Stdout, "  qforge run --question q001 --runner claude")
		fmt.Fprintln(os.Stdout, "  qforge run --question q001 --runner claude --with-visual")
		fmt.Fprintln(os.Stdout, "  qforge run --question q001 --runner codex --runner claude")
		fmt.Fprintln(os.Stdout, "  qforge run --question q001")
		fmt.Fprintln(os.Stdout, "  qforge run --question q002 --runner claude --verbose")
		fmt.Fprintln(os.Stdout, "  qforge run --question q003 --runner codex --model gpt-5.4 --mcp-url https://.../http")
		fmt.Fprintln(os.Stdout, "  qforge run --question q003 --runner codex --model gpt-5.4 --runner claude --model opus --verbose")
	}
	questionRef := fs.String("question", "", "Question id, slug, or folder name")
	datasetName := fs.String("dataset", "", "Override the dataset from question metadata")
	mcpURL := fs.String("mcp-url", "", "Explicit MCP base URL ending in /http")
	mcpServer := fs.String("mcp-server-name", "", "Explicit MCP server name for provider config")
	mcpToken := fs.String("mcp-token", "", "Explicit MCP bearer token")
	mcpTokenFile := fs.String("mcp-token-file", "", "Read MCP token from a file")
	cliBin := fs.String("cli-bin", "", "Override the provider CLI executable")
	withVisual := fs.Bool("with-visual", false, "After SQL succeeds, make a separate presentation call for report.md and visual.html")
	verbose := fs.Bool("verbose", false, "Print phase-level progress logs")
	var runners multiFlag
	var models multiFlag
	fs.Var(&runners, "runner", "Runner to include; repeat for multiple providers; default: codex, claude, gemini")
	fs.Var(&models, "model", "Model override aligned positionally with repeated --runner values")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *questionRef == "" {
		return errors.New("run requires --question")
	}
	if len(runners) == 0 {
		runners = multiFlag{"codex", "claude", "gemini"}
	}
	if *verbose {
		fmt.Printf("[qforge] run question=%s runners=%s\n", *questionRef, strings.Join(runners, ","))
	}
	modelByRunner := map[string]string{}
	for i, runner := range runners {
		if i < len(models) {
			modelByRunner[runner] = models[i]
		}
	}
	var wg sync.WaitGroup
	errCh := make(chan error, len(runners))
	for _, runner := range runners {
		opts := runOptions{
			QuestionRef:  *questionRef,
			Runner:       runner,
			Model:        modelByRunner[runner],
			Dataset:      *datasetName,
			MCPURL:       *mcpURL,
			MCPServer:    *mcpServer,
			MCPToken:     *mcpToken,
			MCPTokenFile: *mcpTokenFile,
			CLIBin:       *cliBin,
			WithVisual:   *withVisual,
			Verbose:      *verbose,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- executeRun(ctx, opts)
		}()
	}
	wg.Wait()
	close(errCh)
	var errs []string
	for err := range errCh {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func runCompare(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Fprintln(os.Stdout, "Usage: qforge compare [flags]")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Compare existing qforge runs one question at a time and generate a rich compare report.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Behavior:")
		fmt.Fprintln(os.Stdout, "  - when --question is set, compares only that question for the selected day")
		fmt.Fprintln(os.Stdout, "  - when --question is omitted, iterates over each question found for the selected day")
		fmt.Fprintln(os.Stdout, "  - writes compare/compare.json and compare_report.md under each question directory")
		fmt.Fprintln(os.Stdout, "  - performs one provider call per question to write the rich compare report")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Flags:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Examples:")
		fmt.Fprintln(os.Stdout, "  qforge compare --day 2026-03-15")
		fmt.Fprintln(os.Stdout, "  qforge compare --day 2026-03-15 --question q001 --runner codex --verbose")
	}
	questionRef := fs.String("question", "", "Restrict compare output to one question id or slug")
	day := fs.String("day", time.Now().Format("2006-01-02"), "Run day in YYYY-MM-DD")
	mcpURL := fs.String("mcp-url", "", "Explicit MCP base URL ending in /http for query_log fetches")
	mcpServer := fs.String("mcp-server-name", "", "Explicit MCP server name for provider config")
	mcpToken := fs.String("mcp-token", "", "Explicit MCP bearer token")
	mcpTokenFile := fs.String("mcp-token-file", "", "Read MCP token from a file")
	runner := fs.String("runner", "codex", "Provider runner for compare_report.md: codex, claude, or gemini")
	modelName := fs.String("model", "", "Model override for the compare report provider")
	cliBin := fs.String("cli-bin", "", "Override the provider CLI executable")
	verbose := fs.Bool("verbose", false, "Print compare progress logs")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	root, err := repoRoot()
	if err != nil {
		return err
	}
	token := *mcpToken
	if *mcpTokenFile != "" && token == "" {
		bytes, err := os.ReadFile(*mcpTokenFile)
		if err != nil {
			return err
		}
		token = strings.TrimSpace(string(bytes))
	}
	if *modelName == "" {
		*modelName, err = defaultModelForRunner(*runner)
		if err != nil {
			return err
		}
	}

	var questionRefs []string
	if *questionRef != "" {
		questionRefs = []string{*questionRef}
	} else {
		questionRefs, err = compare.DiscoverQuestionRefs(root, *day)
		if err != nil {
			return err
		}
	}
	if len(questionRefs) == 0 {
		return fmt.Errorf("no question runs found for %s", *day)
	}

	var errs []string
	for _, ref := range questionRefs {
		if *verbose {
			fmt.Printf("[qforge] compare day=%s question=%s runner=%s model=%s\n", *day, ref, *runner, *modelName)
		}
		if err := executeCompare(ctx, compareOptions{
			QuestionRef: ref,
			Day:         *day,
			Runner:      *runner,
			Model:       *modelName,
			MCPURL:      *mcpURL,
			MCPServer:   *mcpServer,
			MCPToken:    token,
			CLIBin:      *cliBin,
			Verbose:     *verbose,
		}); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

type compareOptions struct {
	QuestionRef string
	Day         string
	Runner      string
	Model       string
	MCPURL      string
	MCPServer   string
	MCPToken    string
	CLIBin      string
	Verbose     bool
}

func executeCompare(ctx context.Context, opts compareOptions) error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	question, err := questions.Resolve(root, opts.QuestionRef)
	if err != nil {
		return err
	}
	cfg, err := datasets.Load(root, question.Meta.Dataset)
	if err != nil {
		return err
	}
	paths := compare.ArtifactPathsForQuestion(root, opts.Day, question.Meta.Slug)
	report, err := compare.Generate(ctx, root, paths.Dir, opts.Day, question.Meta.ID, opts.MCPURL, opts.MCPToken)
	if err != nil {
		return err
	}
	prompt, err := compare.BuildAnalysisPrompt(root, question, report, paths.JSON)
	if err != nil {
		return err
	}
	if err := os.WriteFile(paths.PromptMD, []byte(prompt), 0o644); err != nil {
		return err
	}

	mcpURL, token, err := compareResolveMCPURL(cfg, opts.MCPURL, opts.MCPToken)
	if err != nil {
		return err
	}
	req := model.ProviderRequest{
		Question:      question,
		Dataset:       cfg,
		Prompt:        prompt,
		OutDir:        paths.Dir,
		Model:         opts.Model,
		MCPURL:        mcpURL,
		MCPServerName: datasets.ResolveMCPServerName(cfg, opts.MCPServer),
		MCPToken:      token,
		CLIBin:        opts.CLIBin,
		Verbose:       opts.Verbose,
	}
	provider, err := providers.New(opts.Runner)
	if err != nil {
		return err
	}
	resp, providerErr := provider.GeneratePresentation(ctx, req)
	_ = os.WriteFile(paths.RawAnalysis, []byte(resp.RawOutput), 0o644)
	if providerErr != nil {
		return fmt.Errorf("compare report generation for %s: %w", question.Meta.ID, providerErr)
	}
	reportMD, err := extractCompareMarkdown(resp.RawOutput)
	if err != nil {
		return fmt.Errorf("extract compare report for %s: %w", question.Meta.ID, err)
	}
	return os.WriteFile(paths.ReportMD, []byte(reportMD+"\n"), 0o644)
}

func compareResolveMCPURL(cfg model.DatasetConfig, explicitURL, explicitToken string) (string, string, error) {
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

func extractCompareMarkdown(raw string) (string, error) {
	if block, err := extract.Block(raw, "markdown"); err == nil {
		return block, nil
	}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("empty compare report output")
	}
	return trimmed, nil
}

func runProcessVisual(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("process-visual", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Fprintln(os.Stdout, "Usage: qforge process-visual --run-dir <path> [flags]")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Generate report/html artifacts for an existing run that already has result.json.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Behavior:")
		fmt.Fprintln(os.Stdout, "  - loads manifest.json and result.json from the selected run")
		fmt.Fprintln(os.Stdout, "  - rebuilds the presentation prompt from question metadata and result schema")
		fmt.Fprintln(os.Stdout, "  - invokes the original provider again for report/html template output")
		fmt.Fprintln(os.Stdout, "  - renders final report.md and visual.html in the same run directory")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Flags:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Example:")
		fmt.Fprintln(os.Stdout, "  qforge process-visual --run-dir runs/2026-03-15/q001_hops_per_day/claude/opus/run-004 --verbose")
	}
	runDir := fs.String("run-dir", "", "Path to an existing qforge run directory")
	mcpURL := fs.String("mcp-url", "", "Explicit MCP base URL ending in /http")
	mcpServer := fs.String("mcp-server-name", "", "Explicit MCP server name for provider config")
	mcpToken := fs.String("mcp-token", "", "Explicit MCP bearer token")
	mcpTokenFile := fs.String("mcp-token-file", "", "Read MCP token from a file")
	cliBin := fs.String("cli-bin", "", "Override the provider CLI executable")
	verbose := fs.Bool("verbose", false, "Print phase-level progress logs")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *runDir == "" {
		return errors.New("process-visual requires --run-dir")
	}
	return processVisual(ctx, processVisualOptions{
		RunDir:       *runDir,
		MCPURL:       *mcpURL,
		MCPServer:    *mcpServer,
		MCPToken:     *mcpToken,
		MCPTokenFile: *mcpTokenFile,
		CLIBin:       *cliBin,
		Verbose:      *verbose,
	})
}

func runInspectRun(args []string) error {
	fs := flag.NewFlagSet("inspect-run", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Fprintln(os.Stdout, "Usage: qforge inspect-run --run-dir <path>")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Print one qforge run manifest.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Flags:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Example:")
		fmt.Fprintln(os.Stdout, "  qforge inspect-run --run-dir runs/2026-03-15/q001_hops_per_day/claude/opus/run-004")
	}
	path := fs.String("run-dir", "", "Absolute or relative path to a qforge run directory")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *path == "" {
		return errors.New("inspect-run requires --run-dir")
	}
	manifest, err := runs.ReadManifest(filepath.Join(*path, "manifest.json"))
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

type runOptions struct {
	QuestionRef  string
	Runner       string
	Model        string
	Dataset      string
	MCPURL       string
	MCPServer    string
	MCPToken     string
	MCPTokenFile string
	CLIBin       string
	WithVisual   bool
	Verbose      bool
}

type processVisualOptions struct {
	RunDir       string
	MCPURL       string
	MCPServer    string
	MCPToken     string
	MCPTokenFile string
	CLIBin       string
	Verbose      bool
}

func executeRun(ctx context.Context, opts runOptions) error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	question, err := questions.Resolve(root, opts.QuestionRef)
	if err != nil {
		return err
	}
	datasetName := question.Meta.Dataset
	if opts.Dataset != "" {
		datasetName = opts.Dataset
	}
	cfg, err := datasets.Load(root, datasetName)
	if err != nil {
		return err
	}
	mcpURL, token, err := datasets.ResolveMCPURL(cfg, opts.MCPURL)
	if err != nil {
		return err
	}
	if opts.MCPTokenFile != "" && opts.MCPToken == "" {
		bytes, err := os.ReadFile(opts.MCPTokenFile)
		if err != nil {
			return err
		}
		opts.MCPToken = strings.TrimSpace(string(bytes))
	}
	if opts.MCPToken != "" {
		token = opts.MCPToken
	}
	if opts.Model == "" {
		opts.Model, err = defaultModelForRunner(opts.Runner)
		if err != nil {
			return err
		}
	}
	commandTimeoutSec := question.Meta.CommandTimeoutSec
	if commandTimeoutSec <= 0 {
		commandTimeoutSec = defaultCommandTimeoutSec
	}
	logf(opts.Verbose, "run question=%s runner=%s model=%s dataset=%s", question.Meta.ID, opts.Runner, opts.Model, datasetName)
	outDir, err := runs.NextRunDir(root, question, opts.Runner, opts.Model, time.Now())
	if err != nil {
		return err
	}
	presentationEnabled := question.PresentationEnabled
	logf(opts.Verbose, "out_dir=%s presentation=%t timeout_sec=%d", outDir, presentationEnabled, commandTimeoutSec)
	artifacts := runs.DefaultArtifacts(outDir, presentationEnabled)
	startedAt := time.Now().UTC()
	manifest := model.RunManifest{
		SchemaVersion:   "2",
		Status:          model.RunStatusFailed,
		QuestionID:      question.Meta.ID,
		QuestionSlug:    question.Meta.Slug,
		QuestionTitle:   question.Meta.Title,
		Dataset:         datasetName,
		Runner:          opts.Runner,
		Model:           opts.Model,
		MCPServerName:   datasets.ResolveMCPServerName(cfg, opts.MCPServer),
		MCPConfigSource: filepath.Join("datasets", datasetName, "mcp.yaml"),
		StartedAt:       startedAt,
		Artifacts:       artifacts,
		Phases: model.RunPhases{
			SQLGeneration:          model.PhaseStatusNotRun,
			SQLExecution:           model.PhaseStatusNotRun,
			PresentationGeneration: model.PhaseStatusNotRun,
			PresentationRender:     model.PhaseStatusNotRun,
		},
	}
	defer func() {
		manifest.FinishedAt = time.Now().UTC()
		manifest.DurationSec = int64(manifest.FinishedAt.Sub(startedAt).Seconds())
		_ = runs.WriteManifest(artifacts.ManifestJSON, manifest)
	}()

	sqlPrompt, err := prompts.BuildSQLPrompt(question, cfg)
	if err != nil {
		return err
	}
	logf(opts.Verbose, "phase=sql_generation status=started")
	if err := os.WriteFile(artifacts.PromptSQLRaw, []byte(sqlPrompt), 0o644); err != nil {
		return err
	}
	provider, err := providers.New(opts.Runner)
	if err != nil {
		return err
	}
	req := model.ProviderRequest{
		Question:      question,
		Dataset:       cfg,
		Prompt:        sqlPrompt,
		OutDir:        outDir,
		Model:         opts.Model,
		MCPURL:        mcpURL,
		MCPServerName: manifest.MCPServerName,
		MCPToken:      token,
		CLIBin:        opts.CLIBin,
		Verbose:       opts.Verbose,
	}
	sqlCtx, cancelSQL := context.WithTimeout(ctx, time.Duration(commandTimeoutSec)*time.Second)
	defer cancelSQL()
	sqlResponse, providerErr := provider.GenerateSQL(sqlCtx, req)
	manifest.CLIBin = sqlResponse.CLIBin
	_ = os.WriteFile(artifacts.AnswerSQLRaw, []byte(sqlResponse.RawOutput), 0o644)
	_ = os.WriteFile(artifacts.StdoutLog, []byte(sqlResponse.Stdout), 0o644)
	_ = os.WriteFile(artifacts.StderrLog, []byte(sqlResponse.Stderr), 0o644)
	sqlBlock, err := extract.Block(sqlResponse.RawOutput, "sql")
	if err != nil {
		manifest.Status = model.RunStatusFailed
		manifest.Phases.SQLGeneration = model.PhaseStatusFailed
		if sqlResponse.RawOutput == "" {
			if providerErr != nil {
				return fmt.Errorf("provider %s sql generation: %w", opts.Runner, providerErr)
			}
		}
		return err
	}
	if providerErr != nil {
		manifest.Metadata = addMetadata(manifest.Metadata, "sql_generation_warning", providerErr.Error())
	}
	if err := enforceSQLPolicy(sqlBlock, cfg); err != nil {
		manifest.Status = model.RunStatusPartial
		manifest.Phases.SQLGeneration = model.PhaseStatusFailed
		return err
	}
	if err := os.WriteFile(artifacts.QuerySQL, []byte(sqlBlock+"\n"), 0o644); err != nil {
		return err
	}
	manifest.QuerySHA256 = runs.QuerySHA256(sqlBlock)
	manifest.Phases.SQLGeneration = model.PhaseStatusOK
	logf(opts.Verbose, "phase=sql_generation status=ok query_sha=%s", manifest.QuerySHA256[:12])
	manifest.LogComment = fmt.Sprintf("qforge|question=%s|run=%s|runner=%s|model=%s|phase=full", question.Meta.ID, filepath.Base(outDir), opts.Runner, opts.Model)
	logf(opts.Verbose, "phase=sql_execution status=started log_comment=%s", manifest.LogComment)
	rawDB, result, err := execute.ExecuteSQL(ctx, mcpURL, token, sqlBlock, manifest.LogComment)
	if err != nil {
		manifest.Status = model.RunStatusPartial
		manifest.Phases.SQLExecution = model.PhaseStatusFailed
		return err
	}
	manifest.Phases.SQLExecution = model.PhaseStatusOK
	manifest.ResultRowCount = result.RowCount
	logf(opts.Verbose, "phase=sql_execution status=ok row_count=%d", result.RowCount)
	if err := execute.WriteJSON(artifacts.ResultJSON, result); err != nil {
		return err
	}
	manifest.Metadata = addMetadata(manifest.Metadata, "execution_response_bytes", fmt.Sprintf("%d", len(rawDB)))

	if !presentationEnabled || !opts.WithVisual {
		manifest.Status = model.RunStatusOK
		manifest.Phases.PresentationGeneration = model.PhaseStatusSkipped
		manifest.Phases.PresentationRender = model.PhaseStatusSkipped
		if presentationEnabled && !opts.WithVisual {
			logf(opts.Verbose, "run status=ok presentation=deferred")
		} else {
			logf(opts.Verbose, "run status=ok presentation=skipped")
		}
		return nil
	}

	logf(opts.Verbose, "phase=presentation_generation status=started")
	querySQL, err := os.ReadFile(artifacts.QuerySQL)
	if err != nil {
		return fmt.Errorf("read query.sql for presentation prompt: %w", err)
	}
	prompt, err := prompts.BuildPresentationPrompt(question, cfg, result, string(querySQL))
	if err != nil {
		return err
	}
	if err := os.WriteFile(artifacts.PromptPresentationRaw, []byte(prompt), 0o644); err != nil {
		return err
	}
	req.Prompt = prompt
	presentationCtx, cancelPresentation := context.WithTimeout(ctx, time.Duration(commandTimeoutSec)*time.Second)
	defer cancelPresentation()
	presentationResponse, presentationErr := provider.GeneratePresentation(presentationCtx, req)
	_ = os.WriteFile(artifacts.AnswerPresentationRaw, []byte(presentationResponse.RawOutput), 0o644)
	reportTemplate, htmlTemplate, err := loadPresentationArtifacts(presentationResponse.RawOutput, outDir)
	if err != nil {
		manifest.Status = model.RunStatusPartial
		manifest.Phases.PresentationGeneration = model.PhaseStatusFailed
		return err
	}
	if presentationErr != nil {
		manifest.Metadata = addMetadata(manifest.Metadata, "presentation_generation_warning", presentationErr.Error())
	}
	manifest.Phases.PresentationGeneration = model.PhaseStatusOK
	logf(opts.Verbose, "phase=presentation_generation status=ok")
	if err := os.WriteFile(artifacts.ReportTemplateMD, []byte(reportTemplate), 0o644); err != nil {
		return err
	}
	renderedReport := render.RenderReport(reportTemplate, question, result)
	if err := os.WriteFile(artifacts.ReportMD, []byte(renderedReport), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(artifacts.VisualHTML, []byte(htmlTemplate), 0o644); err != nil {
		return err
	}
	manifest.Phases.PresentationRender = model.PhaseStatusOK
	manifest.Status = model.RunStatusOK
	logf(opts.Verbose, "run status=ok presentation=rendered mode=with-visual")
	return nil
}

func enforceSQLPolicy(sql string, cfg model.DatasetConfig) error {
	lowered := strings.ToLower(sql)
	tokens := tokenizeSQL(lowered)
	tokenSet := make(map[string]struct{}, len(tokens))
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
	}
	for _, forbidden := range datasets.ForbiddenTables(cfg) {
		if _, found := tokenSet[forbidden]; found {
			return fmt.Errorf("sql policy violation: forbidden table %s", forbidden)
		}
	}
	return nil
}

func tokenizeSQL(input string) []string {
	var tokens []string
	var current strings.Builder
	flush := func() {
		if current.Len() == 0 {
			return
		}
		tokens = append(tokens, current.String())
		current.Reset()
	}
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' {
			current.WriteRune(r)
			continue
		}
		flush()
	}
	flush()
	return tokens
}

func addMetadata(metadata map[string]string, key, value string) map[string]string {
	if metadata == nil {
		metadata = map[string]string{}
	}
	metadata[key] = value
	return metadata
}

func logf(verbose bool, format string, args ...any) {
	if !verbose {
		return
	}
	fmt.Printf("[qforge] "+format+"\n", args...)
}

func processVisual(ctx context.Context, opts processVisualOptions) error {
	root, err := repoRoot()
	if err != nil {
		return err
	}
	runDir := opts.RunDir
	if !filepath.IsAbs(runDir) {
		runDir = filepath.Join(root, runDir)
	}
	manifest, err := runs.ReadManifest(filepath.Join(runDir, "manifest.json"))
	if err != nil {
		return err
	}
	resultBytes, err := os.ReadFile(filepath.Join(runDir, "result.json"))
	if err != nil {
		return fmt.Errorf("process-visual requires result.json: %w", err)
	}
	var result model.CanonicalResult
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		return err
	}
	question, err := questions.Resolve(root, manifest.QuestionID)
	if err != nil {
		return err
	}
	if !question.PresentationEnabled {
		return fmt.Errorf("question %s does not declare presentation artifacts", manifest.QuestionID)
	}
	cfg, err := datasets.Load(root, manifest.Dataset)
	if err != nil {
		return err
	}
	mcpURL, token, err := datasets.ResolveMCPURL(cfg, opts.MCPURL)
	if err != nil {
		return err
	}
	if opts.MCPTokenFile != "" && opts.MCPToken == "" {
		bytes, err := os.ReadFile(opts.MCPTokenFile)
		if err != nil {
			return err
		}
		opts.MCPToken = strings.TrimSpace(string(bytes))
	}
	if opts.MCPToken != "" {
		token = opts.MCPToken
	}
	manifest.Artifacts = runs.DefaultArtifacts(runDir, true)
	manifest.MCPServerName = datasets.ResolveMCPServerName(cfg, opts.MCPServer)
	logf(opts.Verbose, "process-visual run_dir=%s question=%s runner=%s model=%s", runDir, manifest.QuestionID, manifest.Runner, manifest.Model)
	querySQL, err := os.ReadFile(filepath.Join(runDir, "query.sql"))
	if err != nil {
		return fmt.Errorf("process-visual requires query.sql: %w", err)
	}
	prompt, err := prompts.BuildPresentationPrompt(question, cfg, result, string(querySQL))
	if err != nil {
		return err
	}
	if err := os.WriteFile(manifest.Artifacts.PromptPresentationRaw, []byte(prompt), 0o644); err != nil {
		return err
	}
	provider, err := providers.New(manifest.Runner)
	if err != nil {
		return err
	}
	commandTimeoutSec := question.Meta.CommandTimeoutSec
	if commandTimeoutSec <= 0 {
		commandTimeoutSec = defaultCommandTimeoutSec
	}
	req := model.ProviderRequest{
		Question:      question,
		Dataset:       cfg,
		Prompt:        prompt,
		OutDir:        runDir,
		Model:         manifest.Model,
		MCPURL:        mcpURL,
		MCPServerName: manifest.MCPServerName,
		MCPToken:      token,
		CLIBin:        firstNonEmpty(opts.CLIBin, manifest.CLIBin),
		Verbose:       opts.Verbose,
	}
	logf(opts.Verbose, "phase=presentation_generation status=started")
	presentationCtx, cancel := context.WithTimeout(ctx, time.Duration(commandTimeoutSec)*time.Second)
	defer cancel()
	resp, providerErr := provider.GeneratePresentation(presentationCtx, req)
	_ = os.WriteFile(manifest.Artifacts.AnswerPresentationRaw, []byte(resp.RawOutput), 0o644)
	if providerErr != nil {
		manifest.Metadata = addMetadata(manifest.Metadata, "presentation_generation_warning", providerErr.Error())
	}
	reportTemplate, htmlTemplate, err := loadPresentationArtifacts(resp.RawOutput, runDir)
	if err != nil {
		manifest.Status = model.RunStatusPartial
		manifest.Phases.PresentationGeneration = model.PhaseStatusFailed
		_ = runs.WriteManifest(manifest.Artifacts.ManifestJSON, manifest)
		return err
	}
	manifest.Phases.PresentationGeneration = model.PhaseStatusOK
	logf(opts.Verbose, "phase=presentation_generation status=ok")
	if err := os.WriteFile(manifest.Artifacts.ReportTemplateMD, []byte(reportTemplate), 0o644); err != nil {
		return err
	}
	renderedReport := render.RenderReport(reportTemplate, question, result)
	if err := os.WriteFile(manifest.Artifacts.ReportMD, []byte(renderedReport), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(manifest.Artifacts.VisualHTML, []byte(htmlTemplate), 0o644); err != nil {
		return err
	}
	manifest.Phases.PresentationRender = model.PhaseStatusOK
	if manifest.Phases.SQLGeneration == model.PhaseStatusOK && manifest.Phases.SQLExecution == model.PhaseStatusOK {
		manifest.Status = model.RunStatusOK
	}
	logf(opts.Verbose, "phase=presentation_render status=ok")
	return runs.WriteManifest(manifest.Artifacts.ManifestJSON, manifest)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func defaultModelForRunner(runner string) (string, error) {
	switch runner {
	case "codex":
		return "gpt-5.4", nil
	case "claude":
		return "opus", nil
	case "gemini":
		return "gemini-3.1-pro-preview", nil
	default:
		return "", fmt.Errorf("no default model for %s", runner)
	}
}

func loadPresentationArtifacts(rawOutput, outDir string) (string, string, error) {
	reportTemplate, reportErr := extract.Block(rawOutput, "report")
	htmlTemplate, htmlErr := extract.Block(rawOutput, "html")
	if reportErr == nil && htmlErr == nil {
		return reportTemplate, htmlTemplate, nil
	}

	reportPath := filepath.Join(outDir, "report.md")
	htmlPath := filepath.Join(outDir, "visual.html")
	reportBytes, readReportErr := os.ReadFile(reportPath)
	htmlBytes, readHTMLErr := os.ReadFile(htmlPath)
	if readReportErr == nil && readHTMLErr == nil {
		reportTemplate = strings.TrimSpace(string(reportBytes))
		if fencedReport, err := extract.Block(reportTemplate, "report"); err == nil {
			reportTemplate = fencedReport
		}
		return reportTemplate, strings.TrimSpace(string(htmlBytes)), nil
	}

	if reportErr != nil {
		return "", "", reportErr
	}
	return "", "", htmlErr
}
