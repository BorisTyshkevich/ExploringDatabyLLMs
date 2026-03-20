package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	verbosepkg "qforge/internal/verbose"
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
	fmt.Fprintln(out, "  process-visual   Generate visual.html for an existing run directory")
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
	fmt.Fprintln(out, "  qforge run -q q001 -r claude -v")
	fmt.Fprintln(out, "  qforge run -q q001 -r claude --with-visual -v")
	fmt.Fprintln(out, "  qforge run -q q001 -r codex -r claude -v")
	fmt.Fprintln(out, "  qforge run -q q001 -v")
	fmt.Fprintln(out, "  qforge process-visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v")
	fmt.Fprintln(out, "  qforge compare --question q003 -r codex -v")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Use `qforge <command> --help` for detailed subcommand help.")
}

func repoRoot() (string, error) {
	if override := strings.TrimSpace(os.Getenv("QFORGE_CODE_ROOT")); override != "" {
		return override, nil
	}
	return os.Getwd()
}

func runsRoot(codeRoot string) string {
	if override := strings.TrimSpace(os.Getenv("QFORGE_RUN_ROOT")); override != "" {
		return override
	}
	return codeRoot
}

func runListQuestions(args []string) error {
	fs := flag.NewFlagSet("list-questions", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Fprintln(os.Stdout, "Usage: qforge list-questions")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "List question metadata discovered under prompts/.")
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
		fmt.Fprintln(os.Stdout, "Usage: qforge run [--question|-q <id|slug>] [--runner|-r <codex|claude|gemini> ...] [flags]")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Run one question through one or more providers.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Behavior:")
		fmt.Fprintln(os.Stdout, "  - loads question and dataset metadata")
		fmt.Fprintln(os.Stdout, "  - prompts each selected provider for final SQL and a report template")
		fmt.Fprintln(os.Stdout, "  - executes SQL itself and writes result.json")
		fmt.Fprintln(os.Stdout, "  - renders final report.md from the saved report template")
		fmt.Fprintln(os.Stdout, "  - runs providers concurrently when more than one is selected")
		fmt.Fprintln(os.Stdout, "  - optionally performs a separate follow-up provider call for visual.html only")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Important:")
		fmt.Fprintln(os.Stdout, "  Visual generation is handled separately by `qforge process-visual`, or by `--with-visual`.")
		fmt.Fprintln(os.Stdout, "  `--with-visual` makes a second independent provider call after SQL execution and report rendering succeed.")
		fmt.Fprintln(os.Stdout, "  If --runner is omitted, qforge runs codex, claude, and gemini.")
		fmt.Fprintln(os.Stdout, "  Repeated --model flags are matched positionally to repeated --runner flags.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Flags:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Examples:")
		fmt.Fprintln(os.Stdout, "  qforge run -q q001 -r claude")
		fmt.Fprintln(os.Stdout, "  qforge run -q q001 -r claude --with-visual")
		fmt.Fprintln(os.Stdout, "  qforge run -q q001 -r codex -r claude")
		fmt.Fprintln(os.Stdout, "  qforge run -q q001")
		fmt.Fprintln(os.Stdout, "  qforge run -q q002 -r claude -v")
		fmt.Fprintln(os.Stdout, "  qforge run -q q003 -r codex --model gpt-5.4 --mcp-url https://.../http")
		fmt.Fprintln(os.Stdout, "  qforge run -q q003 -r codex --model gpt-5.4 -r claude --model opus -v")
	}
	questionRef := fs.String("question", "", "Question id, slug, or folder name")
	fs.StringVar(questionRef, "q", "", "Question id, slug, or folder name (shorthand)")
	datasetName := fs.String("dataset", "", "Override the dataset from question metadata")
	mcpURL := fs.String("mcp-url", "", "Explicit MCP base URL ending in /http")
	mcpServer := fs.String("mcp-server-name", "", "Explicit MCP server name for provider config")
	mcpToken := fs.String("mcp-token", "", "Explicit MCP bearer token")
	mcpTokenFile := fs.String("mcp-token-file", "", "Read MCP token from a file")
	cliBin := fs.String("cli-bin", "", "Override the provider CLI executable")
	withVisual := fs.Bool("with-visual", false, "After SQL and report rendering succeed, make a separate presentation call for visual.html")
	skipVisualValidation := fs.Bool("skip-visual-validation", false, "Skip contract and browser validation for visual.html")
	skipBrowserLiveFetch := fs.Bool("skip-browser-live-fetch", false, "Skip only the browser live-fetch step during visual validation")
	verbose := fs.Bool("verbose", false, "Print phase-level progress logs")
	fs.BoolVar(verbose, "v", false, "Print phase-level progress logs (shorthand)")
	var runners multiFlag
	var models multiFlag
	fs.Var(&runners, "runner", "Runner to include; repeat for multiple providers; default: codex, claude, gemini")
	fs.Var(&runners, "r", "Runner to include; repeat for multiple providers; default: codex, claude, gemini (shorthand)")
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
	modelLabel, err := modelLabelForRunners(runners, models)
	if err != nil {
		return err
	}
	if *verbose {
		logf(true, modelLabel, "run question=%s runners=%s", *questionRef, strings.Join(runners, ","))
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
			QuestionRef:          *questionRef,
			Runner:               runner,
			Model:                modelByRunner[runner],
			Dataset:              *datasetName,
			MCPURL:               *mcpURL,
			MCPServer:            *mcpServer,
			MCPToken:             *mcpToken,
			MCPTokenFile:         *mcpTokenFile,
			CLIBin:               *cliBin,
			WithVisual:           *withVisual,
			SkipVisualValidation: *skipVisualValidation,
			SkipBrowserLiveFetch: *skipBrowserLiveFetch,
			Verbose:              *verbose,
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
		fmt.Fprintln(os.Stdout, "  qforge compare")
		fmt.Fprintln(os.Stdout, "  qforge compare -q q001 -r codex -v")
	}
	questionRef := fs.String("question", "", "Restrict compare output to one question id or slug")
	fs.StringVar(questionRef, "q", "", "Restrict compare output to one question id or slug (shorthand)")
	day := fs.String("day", time.Now().Format("2006-01-02"), "Run day in YYYY-MM-DD")
	mcpURL := fs.String("mcp-url", "", "Explicit MCP base URL ending in /http for compare-report provider config")
	mcpServer := fs.String("mcp-server-name", "", "Explicit MCP server name for provider config")
	mcpToken := fs.String("mcp-token", "", "Explicit MCP bearer token")
	mcpTokenFile := fs.String("mcp-token-file", "", "Read MCP token from a file")
	runner := fs.String("runner", "codex", "Provider runner for compare_report.md: codex, claude, or gemini")
	fs.StringVar(runner, "r", "codex", "Provider runner for compare_report.md: codex, claude, or gemini (shorthand)")
	modelName := fs.String("model", "", "Model override for the compare report provider")
	cliBin := fs.String("cli-bin", "", "Override the provider CLI executable")
	verbose := fs.Bool("verbose", false, "Print compare progress logs")
	fs.BoolVar(verbose, "v", false, "Print compare progress logs (shorthand)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	codeRoot, err := repoRoot()
	if err != nil {
		return err
	}
	runRoot := runsRoot(codeRoot)
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
		questionRefs, err = compare.DiscoverQuestionRefs(runRoot, *day)
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
			logf(true, *modelName, "compare day=%s question=%s runner=%s model=%s", *day, ref, *runner, *modelName)
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
	codeRoot, err := repoRoot()
	if err != nil {
		return err
	}
	runRoot := runsRoot(codeRoot)
	question, err := questions.Resolve(codeRoot, opts.QuestionRef)
	if err != nil {
		return err
	}
	cfg, err := datasets.Load(codeRoot, question.Meta.Dataset)
	if err != nil {
		return err
	}
	paths := compare.ArtifactPathsForQuestion(runRoot, opts.Day, question.Meta.Slug)
	report, err := compare.Generate(ctx, codeRoot, runRoot, paths.Dir, opts.Day, question.Meta.ID, opts.MCPURL, opts.MCPToken)
	if err != nil {
		return err
	}
	prompt, err := compare.BuildAnalysisPrompt(codeRoot, runRoot, question, report, paths.JSON)
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
		fmt.Fprintln(os.Stdout, "Generate visual.html for an existing run that already has query.sql and any mode-specific visual inputs.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Behavior:")
		fmt.Fprintln(os.Stdout, "  - loads manifest.json and query.sql from the selected run")
		fmt.Fprintln(os.Stdout, "  - for static mode, also loads result.json")
		fmt.Fprintln(os.Stdout, "  - rebuilds the visual prompt from the original question and saved artifacts")
		fmt.Fprintln(os.Stdout, "  - invokes the original provider again for visual.html only")
		fmt.Fprintln(os.Stdout, "  - writes final visual.html in the same run directory")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Flags:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Example:")
		fmt.Fprintln(os.Stdout, "  qforge process-visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v")
	}
	runDir := fs.String("run-dir", "", "Path to an existing qforge run directory")
	mcpURL := fs.String("mcp-url", "", "Explicit MCP base URL ending in /http")
	mcpServer := fs.String("mcp-server-name", "", "Explicit MCP server name for provider config")
	mcpToken := fs.String("mcp-token", "", "Explicit MCP bearer token")
	mcpTokenFile := fs.String("mcp-token-file", "", "Read MCP token from a file")
	cliBin := fs.String("cli-bin", "", "Override the provider CLI executable")
	skipVisualValidation := fs.Bool("skip-visual-validation", false, "Skip contract and browser validation for visual.html")
	skipBrowserLiveFetch := fs.Bool("skip-browser-live-fetch", false, "Skip only the browser live-fetch step during visual validation")
	verbose := fs.Bool("verbose", false, "Print phase-level progress logs")
	fs.BoolVar(verbose, "v", false, "Print phase-level progress logs (shorthand)")
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
		RunDir:               *runDir,
		MCPURL:               *mcpURL,
		MCPServer:            *mcpServer,
		MCPToken:             *mcpToken,
		MCPTokenFile:         *mcpTokenFile,
		CLIBin:               *cliBin,
		SkipVisualValidation: *skipVisualValidation,
		SkipBrowserLiveFetch: *skipBrowserLiveFetch,
		Verbose:              *verbose,
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
		fmt.Fprintln(os.Stdout, "  qforge inspect-run --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004")
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
	QuestionRef          string
	Runner               string
	Model                string
	Dataset              string
	MCPURL               string
	MCPServer            string
	MCPToken             string
	MCPTokenFile         string
	CLIBin               string
	WithVisual           bool
	SkipVisualValidation bool
	SkipBrowserLiveFetch bool
	Verbose              bool
}

type processVisualOptions struct {
	RunDir               string
	MCPURL               string
	MCPServer            string
	MCPToken             string
	MCPTokenFile         string
	CLIBin               string
	SkipVisualValidation bool
	SkipBrowserLiveFetch bool
	Verbose              bool
}

func executeRun(ctx context.Context, opts runOptions) error {
	codeRoot, err := repoRoot()
	if err != nil {
		return err
	}
	runRoot := runsRoot(codeRoot)
	question, err := questions.Resolve(codeRoot, opts.QuestionRef)
	if err != nil {
		return err
	}
	datasetName := question.Meta.Dataset
	if opts.Dataset != "" {
		datasetName = opts.Dataset
	}
	cfg, err := datasets.Load(codeRoot, datasetName)
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
	logf(opts.Verbose, opts.Model, "run question=%s runner=%s model=%s dataset=%s", question.Meta.ID, opts.Runner, opts.Model, datasetName)
	outDir, err := runs.NextRunDir(runRoot, question, opts.Runner, opts.Model, time.Now())
	if err != nil {
		return err
	}
	logf(opts.Verbose, opts.Model, "out_dir=%s visual=%t timeout_sec=%d", outDir, question.VisualEnabled, commandTimeoutSec)
	artifacts := runs.DefaultArtifacts(outDir, question.PresentationEnabled)
	startedAt := time.Now().UTC()
	manifest := model.RunManifest{
		SchemaVersion:   "3",
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
	logf(opts.Verbose, opts.Model, "phase=sql_generation status=started")
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
	sqlProviderStartedAt := time.Now()
	sqlResponse, providerErr := provider.GenerateSQL(sqlCtx, req)
	manifest.SQLGenerationProviderDurationMs = time.Since(sqlProviderStartedAt).Milliseconds()
	manifest.CLIBin = sqlResponse.CLIBin
	_ = os.WriteFile(artifacts.AnswerSQLRaw, []byte(sqlResponse.RawOutput), 0o644)
	_ = os.WriteFile(artifacts.StdoutLog, []byte(sqlResponse.Stdout), 0o644)
	_ = os.WriteFile(artifacts.StderrLog, []byte(sqlResponse.Stderr), 0o644)
	analysisArtifact, err := loadAnalysisArtifact(artifacts.AnswerRawJSON)
	if err != nil {
		manifest.Status = model.RunStatusFailed
		manifest.Phases.SQLGeneration = model.PhaseStatusFailed
		if _, statErr := os.Stat(artifacts.AnswerRawJSON); statErr != nil {
			if providerErr != nil {
				return fmt.Errorf("provider %s sql generation: %w", opts.Runner, providerErr)
			}
		}
		return err
	}
	sqlBlock := analysisArtifact.SQL
	reportTemplate := analysisArtifact.ReportMarkdown
	if err := render.ValidateReportTemplate(reportTemplate, analysisArtifact.Metrics); err != nil {
		manifest.Status = model.RunStatusPartial
		manifest.Phases.SQLGeneration = model.PhaseStatusFailed
		return err
	}
	if providerErr != nil {
		manifest.Metadata = addMetadata(manifest.Metadata, "sql_generation_warning", providerErr.Error())
	}
	if err := os.WriteFile(artifacts.QuerySQL, []byte(sqlBlock+"\n"), 0o644); err != nil {
		return err
	}
	analysisJSON, err := json.MarshalIndent(analysisArtifact, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(artifacts.AnalysisJSON, analysisJSON, 0o644); err != nil {
		return err
	}
	if question.ReportEnabled {
		if err := os.WriteFile(artifacts.ReportTemplateMD, []byte(reportTemplate), 0o644); err != nil {
			return err
		}
	}
	manifest.QuerySHA256 = runs.QuerySHA256(sqlBlock)
	manifest.Phases.SQLGeneration = model.PhaseStatusOK
	logf(opts.Verbose, opts.Model, "phase=sql_generation status=ok query_sha=%s", manifest.QuerySHA256[:12])
	manifest.LogComment = fmt.Sprintf("qforge|question=%s|run=%s|runner=%s|model=%s|phase=full", question.Meta.ID, filepath.Base(outDir), opts.Runner, opts.Model)
	logf(opts.Verbose, opts.Model, "phase=sql_execution status=started log_comment=%s", manifest.LogComment)
	rawDB, result, err := execute.ExecuteSQL(ctx, mcpURL, token, sqlBlock, manifest.LogComment)
	if err != nil {
		manifest.Status = model.RunStatusPartial
		manifest.Phases.SQLExecution = model.PhaseStatusFailed
		return err
	}
	manifest.Phases.SQLExecution = model.PhaseStatusOK
	manifest.ResultRowCount = result.RowCount
	logf(opts.Verbose, opts.Model, "phase=sql_execution status=ok row_count=%d", result.RowCount)
	if err := execute.WriteJSON(artifacts.ResultJSON, result); err != nil {
		return err
	}
	visualSummary := buildVisualInputSummary(question, result)
	if err := execute.WriteJSON(artifacts.VisualInputJSON, visualSummary); err != nil {
		return err
	}
	if question.ReportEnabled {
		renderedReport := render.RenderReport(reportTemplate, question, result, analysisArtifact.Metrics)
		if err := os.WriteFile(artifacts.ReportMD, []byte(renderedReport), 0o644); err != nil {
			return err
		}
	}
	manifest.Metadata = addMetadata(manifest.Metadata, "execution_response_bytes", fmt.Sprintf("%d", len(rawDB)))

	if !question.VisualEnabled || !opts.WithVisual {
		manifest.Status = model.RunStatusOK
		manifest.Phases.PresentationGeneration = model.PhaseStatusSkipped
		manifest.Phases.PresentationRender = model.PhaseStatusSkipped
		if question.VisualEnabled && !opts.WithVisual {
			logf(opts.Verbose, opts.Model, "run status=ok visual=deferred")
		} else {
			logf(opts.Verbose, opts.Model, "run status=ok visual=skipped")
		}
		return nil
	}

	logf(opts.Verbose, opts.Model, "phase=presentation_generation status=started")
	querySQL, err := os.ReadFile(artifacts.QuerySQL)
	if err != nil {
		return fmt.Errorf("read query.sql for presentation prompt: %w", err)
	}
	visualInputBytes, err := os.ReadFile(artifacts.VisualInputJSON)
	if err != nil {
		return fmt.Errorf("read visual_input.json for visual prompt: %w", err)
	}
	var visualInput model.VisualInputSummary
	if err := json.Unmarshal(visualInputBytes, &visualInput); err != nil {
		return fmt.Errorf("parse visual_input.json: %w", err)
	}
	prompt, err := prompts.BuildVisualPrompt(question, cfg, result, string(querySQL), dynamicQueryEndpointTemplate(mcpURL, token, cfg), visualInput)
	if err != nil {
		return err
	}
	if err := os.WriteFile(artifacts.PromptPresentationRaw, []byte(prompt), 0o644); err != nil {
		return err
	}
	req.Prompt = prompt
	presentationCtx, cancelPresentation := context.WithTimeout(ctx, time.Duration(commandTimeoutSec)*time.Second)
	defer cancelPresentation()
	presentationStartedAt := time.Now()
	presentationProviderStartedAt := time.Now()
	presentationResponse, presentationErr := provider.GeneratePresentation(presentationCtx, req)
	manifest.PresentationProviderDurationMs = time.Since(presentationProviderStartedAt).Milliseconds()
	_ = os.WriteFile(artifacts.AnswerPresentationRaw, []byte(presentationResponse.RawOutput), 0o644)
	htmlTemplate, err := loadVisualArtifact(presentationResponse.RawOutput, outDir, presentationStartedAt)
	if err != nil {
		logPresentationFailure(opts.Verbose, opts.Model, err, presentationResponse)
		manifest.Status = model.RunStatusPartial
		manifest.Phases.PresentationGeneration = model.PhaseStatusFailed
		return err
	}
	if presentationErr != nil {
		manifest.Metadata = addMetadata(manifest.Metadata, "presentation_generation_warning", presentationErr.Error())
	}
	manifest.Phases.PresentationGeneration = model.PhaseStatusOK
	logf(opts.Verbose, opts.Model, "phase=presentation_generation status=ok")
	if err := os.WriteFile(artifacts.VisualHTML, []byte(htmlTemplate), 0o644); err != nil {
		return err
	}

	validationResult := validatePresentationHTML(ctx, presentationValidationOptions{
		RunDir:               outDir,
		HTMLPath:             artifacts.VisualHTML,
		HTML:                 htmlTemplate,
		Model:                opts.Model,
		VisualMode:           question.Meta.VisualMode,
		VisualType:           question.Meta.VisualType,
		Token:                token,
		SkipVisualValidation: opts.SkipVisualValidation,
		SkipBrowserLiveFetch: opts.SkipBrowserLiveFetch,
		Verbose:              opts.Verbose,
	})
	for key, value := range validationResult.Metadata {
		manifest.Metadata = addMetadata(manifest.Metadata, key, value)
	}
	if !validationResult.Valid {
		manifest.Status = model.RunStatusPartial
		manifest.Phases.PresentationRender = model.PhaseStatusFailed
		logf(opts.Verbose, opts.Model, "run status=partial visual=validation_failed mode=with-visual")
		return nil
	}

	manifest.Phases.PresentationRender = model.PhaseStatusOK
	manifest.Status = model.RunStatusOK
	logf(opts.Verbose, opts.Model, "run status=ok visual=rendered mode=with-visual")
	return nil
}

func enforceSQLPolicy(sql string, cfg model.DatasetConfig) error {
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

func logf(enabled bool, model, format string, args ...any) {
	if !enabled {
		return
	}
	verbosepkg.Printf(os.Stdout, time.Now, model, format, args...)
}

func logPresentationFailure(enabled bool, modelName string, err error, resp model.ProviderResponse) {
	if !enabled {
		return
	}
	logf(true, modelName, "phase=presentation_generation status=failed reason=%q", err.Error())
	if summary := summarizeProviderFailure(resp.Stderr, resp.Stdout, resp.RawOutput); summary != "" {
		logf(true, modelName, "presentation_provider_detail=%q", summary)
	}
}

func summarizeProviderFailure(parts ...string) string {
	for _, part := range parts {
		s := strings.TrimSpace(part)
		if s == "" {
			continue
		}
		lower := strings.ToLower(s)
		switch {
		case strings.Contains(lower, "terminalquotaerror"):
			return firstMatchingLine(s, "TerminalQuotaError")
		case strings.Contains(lower, "quota will reset"):
			return firstMatchingLine(s, "quota will reset")
		case strings.Contains(lower, "quota"):
			return firstMatchingLine(s, "quota")
		case strings.Contains(lower, "rate limit"):
			return firstMatchingLine(s, "rate limit")
		case strings.Contains(lower, "authentication"):
			return firstMatchingLine(s, "authentication")
		case strings.Contains(lower, "unauthorized"):
			return firstMatchingLine(s, "unauthorized")
		case strings.Contains(lower, "forbidden"):
			return firstMatchingLine(s, "forbidden")
		case strings.Contains(lower, "error when talking to"):
			return firstMatchingLine(s, "Error when talking to")
		}

		for _, line := range strings.Split(s, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			l := strings.ToLower(line)
			if strings.Contains(l, "yolo mode is enabled") || strings.Contains(l, "loaded cached credentials") {
				continue
			}
			return truncate(line, 240)
		}
	}
	return ""
}

func firstMatchingLine(s, needle string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(strings.ToLower(line), strings.ToLower(needle)) {
			return truncate(line, 240)
		}
	}
	return truncate(strings.TrimSpace(s), 240)
}

func truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "...(truncated)"
}

func processVisual(ctx context.Context, opts processVisualOptions) error {
	codeRoot, err := repoRoot()
	if err != nil {
		return err
	}
	runRoot := runsRoot(codeRoot)
	runDir := opts.RunDir
	if !filepath.IsAbs(runDir) {
		runDir = filepath.Join(runRoot, runDir)
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
	question, err := questions.Resolve(codeRoot, manifest.QuestionID)
	if err != nil {
		return err
	}
	if !question.VisualEnabled {
		return fmt.Errorf("question %s does not declare visual artifacts", manifest.QuestionID)
	}
	cfg, err := datasets.Load(codeRoot, manifest.Dataset)
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
	manifest.SchemaVersion = "3"
	logf(opts.Verbose, manifest.Model, "process-visual run_dir=%s question=%s runner=%s model=%s", runDir, manifest.QuestionID, manifest.Runner, manifest.Model)
	querySQL, err := os.ReadFile(filepath.Join(runDir, "query.sql"))
	if err != nil {
		return fmt.Errorf("process-visual requires query.sql: %w", err)
	}
	if strings.EqualFold(strings.TrimSpace(question.Meta.VisualMode), "static") {
		if _, err := os.Stat(filepath.Join(runDir, "result.json")); err != nil {
			return fmt.Errorf("process-visual requires result.json for static mode: %w", err)
		}
	}
	visualInput, err := ensureVisualInputSummary(filepath.Join(runDir, "visual_input.json"), question, result)
	if err != nil {
		return err
	}
	prompt, err := prompts.BuildVisualPrompt(question, cfg, result, string(querySQL), dynamicQueryEndpointTemplate(mcpURL, token, cfg), visualInput)
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
	logf(opts.Verbose, manifest.Model, "phase=presentation_generation status=started")
	presentationCtx, cancel := context.WithTimeout(ctx, time.Duration(commandTimeoutSec)*time.Second)
	defer cancel()
	presentationStartedAt := time.Now()
	presentationProviderStartedAt := time.Now()
	resp, providerErr := provider.GeneratePresentation(presentationCtx, req)
	manifest.PresentationProviderDurationMs = time.Since(presentationProviderStartedAt).Milliseconds()
	_ = os.WriteFile(manifest.Artifacts.AnswerPresentationRaw, []byte(resp.RawOutput), 0o644)
	if providerErr != nil {
		manifest.Metadata = addMetadata(manifest.Metadata, "presentation_generation_warning", providerErr.Error())
	}
	htmlTemplate, err := loadVisualArtifact(resp.RawOutput, runDir, presentationStartedAt)
	if err != nil {
		logPresentationFailure(opts.Verbose, manifest.Model, err, resp)
		manifest.Status = model.RunStatusPartial
		manifest.Phases.PresentationGeneration = model.PhaseStatusFailed
		_ = runs.WriteManifest(manifest.Artifacts.ManifestJSON, manifest)
		return err
	}
	manifest.Phases.PresentationGeneration = model.PhaseStatusOK
	logf(opts.Verbose, manifest.Model, "phase=presentation_generation status=ok")
	if err := os.WriteFile(manifest.Artifacts.VisualHTML, []byte(htmlTemplate), 0o644); err != nil {
		return err
	}

	validationResult := validatePresentationHTML(ctx, presentationValidationOptions{
		RunDir:               runDir,
		HTMLPath:             manifest.Artifacts.VisualHTML,
		HTML:                 htmlTemplate,
		Model:                manifest.Model,
		VisualMode:           question.Meta.VisualMode,
		VisualType:           question.Meta.VisualType,
		Token:                token,
		SkipVisualValidation: opts.SkipVisualValidation,
		SkipBrowserLiveFetch: opts.SkipBrowserLiveFetch,
		Verbose:              opts.Verbose,
	})
	for key, value := range validationResult.Metadata {
		manifest.Metadata = addMetadata(manifest.Metadata, key, value)
	}
	if !validationResult.Valid {
		manifest.Status = model.RunStatusPartial
		manifest.Phases.PresentationRender = model.PhaseStatusFailed
		logf(opts.Verbose, manifest.Model, "phase=presentation_render status=failed")
		return runs.WriteManifest(manifest.Artifacts.ManifestJSON, manifest)
	}

	manifest.Phases.PresentationRender = model.PhaseStatusOK
	if manifest.Phases.SQLGeneration == model.PhaseStatusOK && manifest.Phases.SQLExecution == model.PhaseStatusOK {
		manifest.Status = model.RunStatusOK
	}
	logf(opts.Verbose, manifest.Model, "phase=presentation_render status=ok")
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

func modelLabelForRunners(runners, explicitModels []string) (string, error) {
	labels := make([]string, 0, len(runners))
	seen := map[string]struct{}{}
	for i, runner := range runners {
		modelName := ""
		if i < len(explicitModels) {
			modelName = explicitModels[i]
		}
		if modelName == "" {
			var err error
			modelName, err = defaultModelForRunner(runner)
			if err != nil {
				return "", err
			}
		}
		if _, ok := seen[modelName]; ok {
			continue
		}
		seen[modelName] = struct{}{}
		labels = append(labels, modelName)
	}
	return strings.Join(labels, ","), nil
}

func dynamicQueryEndpointTemplate(mcpURL, token string, cfg model.DatasetConfig) string {
	baseURL := strings.TrimRight(cfg.MCPBaseURL, "/")
	if strings.HasSuffix(mcpURL, "/http") {
		trimmed := strings.TrimSuffix(mcpURL, "/http")
		if token != "" && strings.HasSuffix(trimmed, "/"+token) {
			trimmed = strings.TrimSuffix(trimmed, "/"+token)
		}
		if trimmed != "" {
			baseURL = strings.TrimRight(trimmed, "/")
		}
	}
	if baseURL == "" {
		baseURL = "https://mcp.demo.altinity.cloud"
	}
	return baseURL + "/{JWE}/openapi/execute_query?query=..."
}

func buildVisualInputSummary(question model.Question, result model.CanonicalResult) model.VisualInputSummary {
	summary := model.VisualInputSummary{
		QuestionTitle: question.Meta.Title,
		ResultColumns: append([]string(nil), result.Columns...),
		RowCount:      result.RowCount,
		ModeHint:      visualModeHint(question.Meta.VisualMode),
	}
	if len(result.Rows) > 0 {
		limit := 2
		if len(result.Rows) < limit {
			limit = len(result.Rows)
		}
		summary.SampleRows = make([]map[string]any, 0, limit)
		for i := 0; i < limit; i++ {
			summary.SampleRows = append(summary.SampleRows, result.Rows[i])
		}
	}
	notes := map[string]string{}
	for _, col := range result.Columns {
		note := detectFieldShapeNote(col, result.Rows)
		if note != "" {
			notes[col] = note
		}
	}
	if len(notes) > 0 {
		summary.FieldShapeNotes = notes
	}
	return summary
}

func ensureVisualInputSummary(path string, question model.Question, result model.CanonicalResult) (model.VisualInputSummary, error) {
	var visualInput model.VisualInputSummary
	visualInputBytes, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(visualInputBytes, &visualInput); err != nil {
			return model.VisualInputSummary{}, fmt.Errorf("parse visual_input.json: %w", err)
		}
		return visualInput, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return model.VisualInputSummary{}, fmt.Errorf("read visual_input.json: %w", err)
	}

	visualInput = buildVisualInputSummary(question, result)
	if err := execute.WriteJSON(path, visualInput); err != nil {
		return model.VisualInputSummary{}, fmt.Errorf("write visual_input.json: %w", err)
	}
	return visualInput, nil
}

func visualModeHint(mode string) string {
	if strings.EqualFold(strings.TrimSpace(mode), "static") {
		return "Static mode embeds analytical data from result.json directly in the page."
	}
	return "Dynamic mode still fetches live data in the browser via query.sql and the configured endpoint."
}

func detectFieldShapeNote(column string, rows []map[string]any) string {
	for _, row := range rows {
		value, ok := row[column]
		if !ok || value == nil {
			continue
		}
		switch v := value.(type) {
		case []any:
			return "array field"
		case []string:
			return "array field"
		case string:
			if len(v) >= len("2006-01-02T15:04:05Z") && strings.Contains(v, "T") && strings.HasSuffix(v, "Z") {
				return "ISO-like timestamp string"
			}
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

func loadVisualArtifact(rawOutput, outDir string, notBefore time.Time) (string, error) {
	htmlTemplate, htmlErr := extract.Block(rawOutput, "html")
	if htmlErr == nil {
		return htmlTemplate, nil
	}

	htmlPath := filepath.Join(outDir, "visual.html")
	htmlInfo, htmlStatErr := os.Stat(htmlPath)
	htmlBytes, readHTMLErr := os.ReadFile(htmlPath)
	if htmlStatErr == nil && readHTMLErr == nil && !htmlInfo.ModTime().Before(notBefore) {
		return strings.TrimSpace(string(htmlBytes)), nil
	}

	return "", htmlErr
}

func loadAnalysisArtifact(path string) (model.AnalysisArtifact, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return model.AnalysisArtifact{}, fmt.Errorf("read answer.raw.json: %w", err)
	}
	var artifact model.AnalysisArtifact
	if err := json.Unmarshal(payload, &artifact); err != nil {
		return model.AnalysisArtifact{}, fmt.Errorf("invalid analysis json: %w", err)
	}
	artifact.SQL = normalizeEscapedMultiline(strings.TrimSpace(artifact.SQL))
	artifact.ReportMarkdown = normalizeEscapedMultiline(strings.TrimSpace(artifact.ReportMarkdown))
	if artifact.SQL == "" {
		return model.AnalysisArtifact{}, fmt.Errorf("analysis json missing non-empty sql")
	}
	if artifact.ReportMarkdown == "" {
		return model.AnalysisArtifact{}, fmt.Errorf("analysis json missing non-empty report_markdown")
	}
	return artifact, nil
}

func normalizeEscapedMultiline(value string) string {
	if value == "" {
		return value
	}
	if !strings.Contains(value, `\`) {
		return value
	}
	unquoted, err := strconv.Unquote(`"` + value + `"`)
	if err != nil {
		return value
	}
	return strings.TrimSpace(unquoted)
}
