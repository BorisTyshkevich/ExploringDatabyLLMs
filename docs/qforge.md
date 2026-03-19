# qforge 

`qforge` is a Go-based harness for model-generated analytics code.

It is designed around a strict split of responsibilities:

- the model writes code
- the harness executes SQL itself
- the canonical data artifact is JSON
- performance metrics come from `system.query_log`
- comparison happens after runs, not during them

This repository still contains historical Bash and Python benchmark code, but the active path documented here is `qforge`.

## Core Model

`qforge` uses a two-phase design:

1. SQL generation
   - the model is prompted to inspect schema and self-verify its SQL before writing artifacts
   - the model writes `answer.raw.json` in the run directory
   - `answer.raw.json` contains the analysis artifact with `sql`, `report_markdown`, and `metrics`
2. Optional presentation generation
   - the model writes final `html`
   - the final `report.md` is rendered by the harness from the saved analysis artifact plus JSON-derived sections
   - `visual.html` is model-authored final output and is not patched by `qforge`
   - `visual.html` is validated first against the visual contract and then in a headless browser via `chromedp`
   - this can be done later with `process-visual` or immediately with `run --with-visual`

The model should not emit result rows directly.

The harness always executes the final SQL itself and writes `result.json`.

## Prompt Assets

Prompt assembly is split into shared and phase-specific assets under [`/Users/bvt/work/ExploringDatabyLLMs/prompts`](/Users/bvt/work/ExploringDatabyLLMs/prompts):

- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md)
  - shared qforge and dataset-scope guidance used by both SQL and presentation phases
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_sql.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_sql.md)
  - SQL-only rules such as schema inspection, self-verification, and the `answer.raw.json` analysis-artifact contract
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_presentation.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_presentation.md)
  - report/template rules for the presentation phase
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual.md)
  - shared `visual.html` rules
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md)
  - dynamic `visual.html` requirements
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_static.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_static.md)
  - static `visual.html` requirements

Question-specific files such as `prompts/qXXX.../prompt.md`, `report_prompt.md`, and `visual_prompt.md` should contain task logic, not repeated dataset boilerplate.

Question metadata may also declare `visual_mode`:

- `dynamic`
  - live browser dashboard that uses saved SQL plus browser-side MCP fetch
- `static`
  - self-contained benchmark artifact with embedded analytical data

If `visual_mode` is absent, qforge treats the question as `dynamic` for backward compatibility.

Template variables, prompt composition, and dataset mapping are documented in [`/Users/bvt/work/ExploringDatabyLLMs/docs/prompt-templates.md`](/Users/bvt/work/ExploringDatabyLLMs/docs/prompt-templates.md).

End-to-end phase behavior and artifact ownership are documented in [`/Users/bvt/work/ExploringDatabyLLMs/docs/processing-logic.md`](/Users/bvt/work/ExploringDatabyLLMs/docs/processing-logic.md).

## Setup

Set the demo JWE token:

```bash
export MCP_JWE_TOKEN="YOUR_PUBLIC_DEMO_JWE_TOKEN"
```

Current public demo example from this repo:

```bash
export MCP_JWE_TOKEN="eyJhbGciOiJBMjU2S1ciLCJjdHkiOiJKU09OIiwiZW5jIjoiQTI1NkdDTSIsInR5cCI6IkpXRSJ9.1Zhu2eydw0lNdOwL81KM0Z3_Q9hgpKCgqlyAtDkyMMzf39tuz0tnYQ.2ZWKRwXcebF2f-Zy.SQLoNUAExT0uf7GhTdOKfK9i4yHZRN77Bxa4yQT1lAKUHvEY9vZgaUCD3FXYEOz5y_Njt5S9ERVDW1qdFI8EdT1bQfO7tJaI_VeU51xFDETygTWMs9NTACNxVQFJHsvfo9ZY4vrT7HamA1UD-bH1erFfKug6YsLf2j-Pa6DvjI4-ODZpX1HBNKm2uU__8qkwC-a09IaU1QYSgXb2kKFMAqLkWgrMQ041CkFNUA.NGTvmP8D7i6CWg9V67nkNw"
```

Optional:

```bash
export MCP_BASE_URL="https://mcp.demo.altinity.cloud"
```

The effective MCP URL is:

```text
https://mcp.demo.altinity.cloud/$MCP_JWE_TOKEN/http
```

Provider CLIs expected on `PATH`:

- `codex`
- `claude`
- `gemini`

## Quick Start

List questions:

```bash
./scripts/qforge list-questions
```

Run one question:

```bash
./scripts/qforge run -q q001 -r claude -v
```

Run one question and immediately follow with a separate presentation call:

```bash
./scripts/qforge run -q q001 -r claude --with-visual -v
```

Run one question with presentation generation but skip only the live browser fetch step:

```bash
./scripts/qforge run -q q001 -r claude --with-visual --skip-browser-live-fetch -v
```

Run one question across all default providers:

```bash
./scripts/qforge run -q q001 -v
```

Run one question across selected providers:

```bash
./scripts/qforge run -q q001 -r codex -r claude -v
```

Process report and visual for an existing run:

```bash
./scripts/qforge process-visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
```

Process report and visual but skip all visual validation:

```bash
./scripts/qforge process-visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 --skip-visual-validation -v
```

Compare runs for a day:

```bash
./scripts/qforge compare -v
```

Inspect one run directory:

```bash
./scripts/qforge inspect-run --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004
```

## Command Reference

All commands support `--help`.

Top-level help:

```bash
./scripts/qforge --help
```

### `qforge list-questions`

Usage:

```bash
./scripts/qforge list-questions
```

Behavior:

- prints one tab-separated line per question
- output columns are:
  - `question_id`
  - `question_slug`
  - `title`
  - `dataset`
  - `presentation_enabled`

### `qforge run`

Usage:

```bash
./scripts/qforge run [--question|-q <id|slug>] [--runner|-r <codex|claude|gemini> ...] [flags]
```

Flags:

- `--question`
  - required
  - question id, slug, or folder name
- `-q`
  - shorthand for `--question`
- `--runner`
  - optional, repeatable
  - provider runner: `codex`, `claude`, or `gemini`
- `-r`
  - shorthand for `--runner`
  - default when omitted: `codex`, `claude`, `gemini`
- `--model`
  - optional, repeatable
  - override the default model for the selected runner
  - current defaults: `codex -> gpt-5.4`, `claude -> opus`, `gemini -> gemini-3.1-pro-preview`
  - matched positionally with repeated `--runner` flags
  - example:
    - `--runner codex --model gpt-5.4 --runner claude --model opus`
- `--dataset`
  - optional
  - override the dataset from question metadata
- `--mcp-url`
  - optional
  - explicit MCP URL ending in `/http`
- `--mcp-server-name`
  - optional
  - explicit MCP server name for provider config
- `--mcp-token`
  - optional
  - explicit MCP bearer token
- `--mcp-token-file`
  - optional
  - read MCP token from a file
- `--cli-bin`
  - optional
  - override the provider CLI executable
- `--verbose`
  - optional
  - print phase-level progress logs and provider subprocess timing
- `-v`
  - shorthand for `--verbose`
- `--with-visual`
  - optional
  - after SQL succeeds and `report.md` is rendered, make a second independent provider call for `visual.html`
  - this is equivalent in behavior to running `process-visual` after a successful run
- `--skip-visual-validation`
  - optional
  - skip both contract validation and browser validation for `visual.html`
- `--skip-browser-live-fetch`
  - optional
  - keep contract validation and browser smoke checks, but skip the token-entry and live fetch interaction

What `run` does:

1. resolves question metadata
2. selects one or more providers
3. builds the SQL prompt for each selected provider
4. invokes those providers, concurrently when more than one is selected
5. extracts fenced SQL and fenced report template
6. executes SQL directly against the OpenAPI endpoint
7. writes canonical `result.json`
8. renders `report.md` from the saved report template
9. writes `manifest.json`
10. optionally makes a second independent provider call for `visual.html` when `--with-visual` is set

What `run` does not do:

- it does not produce `visual.html` unless `--with-visual` is set
- use `qforge process-visual` for HTML generation later

Exception:

- when `--with-visual` is set, `run` performs the SQL phase first and then a separate presentation call for successful runs

### `qforge process-visual`

Usage:

```bash
./scripts/qforge process-visual --run-dir <path> [flags]
```

Flags:

- `--run-dir`
  - required
  - path to an existing qforge run directory
- `--mcp-url`
  - optional
  - explicit MCP URL ending in `/http`
- `--mcp-server-name`
  - optional
  - explicit MCP server name for provider config
- `--mcp-token`
  - optional
  - explicit MCP bearer token
- `--mcp-token-file`
  - optional
  - read MCP token from a file
- `--cli-bin`
  - optional
  - override the provider CLI executable
- `--skip-visual-validation`
  - optional
  - skip both contract validation and browser validation for `visual.html`
- `--skip-browser-live-fetch`
  - optional
  - keep browser smoke checks but skip the live token-entry and fetch step
- `--verbose`
  - optional
  - print phase-level progress logs

What `process-visual` does:

- loads `manifest.json`, `analysis.json`, `query.sql`, `report.template.md`, `report.md`, and `result.json` from an existing run
- rebuilds the presentation prompt from question metadata and the saved artifacts
- invokes the original provider again for `html`
- validates `visual.html` in two stages unless `--skip-visual-validation` is set:
  - contract validation against the shared visual rules
  - browser validation using `chromedp`
- for dynamic dashboards, browser validation automatically attempts the live MCP fetch path when a token is available unless `--skip-browser-live-fetch` is set
- writes:
  - `prompt.presentation.md`
  - `answer.presentation.raw.md`
  - `report.template.md`
  - `report.md`
  - `visual.html`

Report placeholder contract:

- scalar placeholders:
  - `{{row_count}}`
  - `{{generated_at}}`
  - `{{columns_csv}}`
  - `{{question_title}}`
- Markdown placeholders:
  - `{{data_overview_md}}`
  - `{{result_table_md}}`

If the report template omits the Markdown placeholders, `qforge` appends:

- `## Data Overview`
- `## Result Rows`

### `qforge compare`

Usage:

```bash
./scripts/qforge compare [flags]
```

Flags:

- `--day`
  - optional
  - run day in `YYYY-MM-DD`
  - default: current local day
- omitting `--day` compares runs for today
- `--question`
  - optional
  - restrict compare to one question id or slug
  - if omitted, `compare` iterates all questions found for that day and runs one compare pass per question
- `-q`
  - shorthand for `--question`
- `--runner`
  - optional
  - provider runner used for `compare_report.md`
  - default: `codex`
- `-r`
  - shorthand for `--runner`
- `--model`
  - optional
  - override the default model for the compare report provider
- `--cli-bin`
  - optional
  - override the provider CLI executable for the compare report provider
- `--mcp-url`
  - optional
  - explicit MCP URL ending in `/http` for `system.query_log` fetches and any direct validation queries used during report generation
- `--mcp-server-name`
  - optional
  - explicit MCP server name for provider config
- `--mcp-token`
  - optional
  - explicit MCP bearer token
- `--mcp-token-file`
  - optional
  - read MCP token from a file
- `--verbose`
  - optional
  - print compare progress logs
- `-v`
  - shorthand for `--verbose`

What `compare` writes:

- `<day>/<question-slug>/compare/compare.json`
- `<day>/<question-slug>/compare/analysis.prompt.md`
- `<day>/<question-slug>/compare/analysis.raw.md`
- `<day>/<question-slug>/compare_report.md`

What `compare` does:

- resolves one question per compare pass
- loads `manifest.json` and `result.json` from matching runs
- fetches deferred performance metrics from `system.query_log` using `log_comment`
- writes compact structured compare data to `compare/compare.json`
- runs one provider call using the shared analysis prompt to generate `compare_report.md`
- includes partial runs in the structured compare output when possible

### `qforge inspect-run`

Usage:

```bash
./scripts/qforge inspect-run --run-dir <path>
```

Flags:

- `--run-dir`
  - required
  - absolute or relative path to a `qforge` run directory

Behavior:

- prints the run `manifest.json`

## Verbose Mode

`--verbose` enables operational logging.

Current verbose output includes:

- selected question, runners, model, dataset
- allocated run directory
- whether presentation is enabled
- SQL generation start/finish
- provider subprocess start/finish
- provider elapsed time
- SQL execution start/finish
- final row count
- compare startup summary

Example:

```text
[qforge] run question=q002 runners=claude
[qforge] run question=q002 runner=claude model=opus dataset=ontime
[qforge] out_dir=... presentation=false timeout_sec=900
[qforge] phase=sql_generation status=started
[qforge] provider=claude phase=start ...
[qforge] provider=claude phase=done status=ok elapsed=1m29.067s
[qforge] phase=sql_execution status=ok row_count=195
```

## Artifacts

Each run is stored under:

```text
YYYY-MM-DD/<question-slug>/<runner>/<model>/run-XXX/
```

Typical SQL-only run artifacts:

- `prompt.sql.md`
- `answer.sql.raw.md`
- `answer.raw.json`
- `query.sql`
- `result.json`
- `manifest.json`
- `stdout.log`
- `stderr.log`

When presentation is processed later with `qforge process-visual`:

- `prompt.presentation.md`
- `answer.presentation.raw.md`
- `visual.html`

Typical analysis-phase artifacts now include:

- `analysis.json`
- `query.sql`
- `report.template.md`
- `report.md`

## Browser Validation

Generated `visual.html` artifacts now go through a headless browser check with `chromedp`.

What this browser phase verifies:

- the page opens and reaches a ready DOM state
- required dashboard controls exist
- runtime exceptions and fatal console errors are absent
- dynamic dashboards can accept the JWE token through the page UI, click the footer action button, issue the expected MCP request, and settle without obvious failure UI

Token behavior:

- the browser phase reuses the same resolved JWE token that `qforge` already uses for SQL execution
- `--skip-browser-live-fetch` keeps the browser load/runtime checks but skips the real fetch interaction
- `--skip-visual-validation` skips both contract validation and browser validation

Failure behavior:

- contract validation failures mark presentation as partial/failed
- browser validation failures also mark presentation as partial/failed
- `manifest.json` records browser validation details under `browser_validation*` metadata keys

## Canonical Output

The canonical data artifact is `result.json`.

It is written by the harness, not the model.

Current shape:

- `columns`
- `rows`
- `row_count`
- `generated_at`
- `source_query_sha256`
- `log_comment`

## `log_comment` And Performance Metrics

Every full SQL execution is run with a deterministic `log_comment`.

That comment is later used by `qforge compare` to fetch execution metrics from `system.query_log`, including:

- query duration
- read rows
- read bytes
- result rows
- result bytes
- memory usage
- peak threads

This is why performance metrics are not fetched during `run` itself.

## Current Scope And Caveats

Important caveats:

- The repository still contains older Bash and Python harness code and old docs; those are not the primary path documented here.
- Provider behavior varies. Some providers can spend a long time in self-verification loops before returning SQL.
- Presentation is a separate explicit step by default. Use `run --with-visual` if you want `run` to also make the follow-up presentation call.

## Recommended Workflows

SQL/JSON verification only:

```bash
./scripts/qforge run -q q001 -r claude -v
```

Single command with SQL plus follow-up report/visual generation:

```bash
./scripts/qforge run -q q003 -r claude --with-visual -v
```

Three-question verification pass:

```bash
./scripts/qforge run -q q001 -r claude -v
./scripts/qforge run -q q002 -r claude -v
./scripts/qforge run -q q003 -r claude -v
./scripts/qforge compare -v
```

Multi-provider comparison:

```bash
./scripts/qforge run -q q001 -r codex -r claude -v
./scripts/qforge compare -q q001 -r codex -v
```

Process presentation for one completed run:

```bash
./scripts/qforge process-visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
```
