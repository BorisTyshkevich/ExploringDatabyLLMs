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
   - the model is prompted to inspect schema and self-verify its SQL before returning it
   - the model emits only a fenced `sql` block
2. Optional presentation generation
   - the model writes final `html` and template-style `report`
   - the final `report.md` is rendered as Markdown from the template plus JSON-derived sections
   - `visual.html` is model-authored final output and is not patched by `qforge`
   - this can be done later with `process-visual` or immediately with `run --with-visual`

The model should not emit result rows directly.

The harness always executes the final SQL itself and writes `result.json`.

## Prompt Assets

Prompt assembly is split into shared and phase-specific assets under [`/Users/bvt/work/ExploringDatabyLLMs/prompts`](/Users/bvt/work/ExploringDatabyLLMs/prompts):

- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md)
  - shared qforge and dataset-scope guidance used by both SQL and presentation phases
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_sql.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_sql.md)
  - SQL-only rules such as schema inspection, self-verification, and the fenced `sql` output contract
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_presentation.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_presentation.md)
  - report/template rules for the presentation phase
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md)
  - dynamic `visual.html` requirements

Question-specific files such as `prompts/qXXX.../prompt.md`, `report_prompt.md`, and `visual_prompt.md` should contain task logic, not repeated dataset boilerplate.

Template variables, prompt composition, and dataset mapping are documented in [`/Users/bvt/work/ExploringDatabyLLMs/docs/prompt-templates.md`](/Users/bvt/work/ExploringDatabyLLMs/docs/prompt-templates.md).

## Setup

Set the demo JWE token:

```bash
export MCP_JWE_TOKEN="YOUR_PUBLIC_DEMO_JWE_TOKEN"
```

Current public demo example from this repo:

```bash
export MCP_JWE_TOKEN="eyJhbGciOiJBMjU2S1ciLCJjdHkiOiJKU09OIiwiZW5jIjoiQTI1NkdDTSIsInR5cCI6IkpXRSJ9.dRhJ0FFe_7MpwMisrxo3z1pTd1xCHp5KUUK7kgNt2GTPtTVxGvU1mw.ZQzRGK7Wre9_REqj.pybQI8aY4pvKel7wR1THrOmkdpIWHGbjJ-cAflJKRodC1AUJsiFY1L13Gxh8L_dnG4j0oSuRdX4n6_W7S6mOTGKW5cAbQ2DS-T4Z7YULrBfebew2Px5QXgrLKRv7EOcO_2fBqpZSbpD5wRfRtYHZvLxuAVOjRAt4E2_WVKhL6yURWbKacK6iqBwndLWhdDKNvhtmwxEowtfiOhFTMP2870iiGeNJ25k1P70nJKYQ9qnROx9P.DIjBhAvis-jqhS02WrOXHg"
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
./scripts/qforge run --question q001 --runner claude --verbose
```

Run one question and immediately follow with a separate presentation call:

```bash
./scripts/qforge run --question q001 --runner claude --with-visual --verbose
```

Run one question across all default providers:

```bash
./scripts/qforge run --question q001 --verbose
```

Run one question across selected providers:

```bash
./scripts/qforge run --question q001 --runner codex --runner claude --verbose
```

Process report and visual for an existing run:

```bash
./scripts/qforge process-visual --run-dir runs/2026-03-15/q001_hops_per_day/claude/opus/run-004 --verbose
```

Compare runs for a day:

```bash
./scripts/qforge compare --day "$(date +%F)" --verbose
```

Inspect one run directory:

```bash
./scripts/qforge inspect-run --run-dir runs/2026-03-15/q001_hops_per_day/claude/opus/run-004
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
./scripts/qforge run --question <id|slug> [--runner <codex|claude|gemini> ...] [flags]
```

Flags:

- `--question`
  - required
  - question id, slug, or folder name
- `--runner`
  - optional, repeatable
  - provider runner: `codex`, `claude`, or `gemini`
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
- `--with-visual`
  - optional
  - after SQL succeeds, make a second independent provider call for `report.md` and `visual.html`
  - this is equivalent in behavior to running `process-visual` after a successful run

What `run` does:

1. resolves question metadata
2. selects one or more providers
3. builds the SQL prompt for each selected provider
4. invokes those providers, concurrently when more than one is selected
5. extracts fenced SQL
6. enforces the dataset table policy
7. executes SQL directly against the OpenAPI endpoint
8. writes canonical `result.json`
9. writes `manifest.json`
10. optionally makes a second independent provider call for `report.md` and `visual.html` when `--with-visual` is set

What `run` does not do:

- it does not produce `report.md` or `visual.html`
- use `qforge process-visual` for Markdown report and HTML generation later

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
- `--verbose`
  - optional
  - print phase-level progress logs

What `process-visual` does:

- loads `manifest.json` and `result.json` from an existing run
- rebuilds the presentation prompt from question metadata and the JSON schema
- invokes the original provider again for `report` and `html`
- fills report placeholders from `result.json`
- injects Markdown data sections from `result.json`
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
- `--question`
  - optional
  - restrict compare to one question id or slug
  - if omitted, `compare` iterates all questions found for that day and runs one compare pass per question
- `--runner`
  - optional
  - provider runner used for `compare_report.md`
  - default: `codex`
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

What `compare` writes:

- `runs/<day>/<question-slug>/compare/compare.json`
- `runs/<day>/<question-slug>/compare/analysis.prompt.md`
- `runs/<day>/<question-slug>/compare/analysis.raw.md`
- `runs/<day>/<question-slug>/compare_report.md`

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
[qforge] run question=q002 runner=claude model=opus dataset=ontime_v2
[qforge] out_dir=... presentation=false timeout_sec=900
[qforge] phase=sql_generation status=started
[qforge] provider=claude phase=start ...
[qforge] provider=claude phase=done status=ok elapsed=1m29.067s
[qforge] phase=sql_execution status=ok row_count=195
```

## Artifacts

Each run is stored under:

```text
runs/YYYY-MM-DD/<question-slug>/<runner>/<model>/run-XXX/
```

Typical SQL-only run artifacts:

- `prompt.sql.md`
- `answer.sql.raw.md`
- `query.sql`
- `result.json`
- `manifest.json`
- `stdout.log`
- `stderr.log`

When presentation is processed later with `qforge process-visual`:

- `prompt.presentation.md`
- `answer.presentation.raw.md`
- `report.template.md`
- `report.md`
- `visual.html`

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
./scripts/qforge run --question q001 --runner claude --verbose
```

Single command with SQL plus follow-up report/visual generation:

```bash
./scripts/qforge run --question q003 --runner claude --with-visual --verbose
```

Three-question verification pass:

```bash
./scripts/qforge run --question q001 --runner claude --verbose
./scripts/qforge run --question q002 --runner claude --verbose
./scripts/qforge run --question q003 --runner claude --verbose
./scripts/qforge compare --day "$(date +%F)" --verbose
```

Multi-provider comparison:

```bash
./scripts/qforge run --question q001 --runner codex --runner claude --verbose
./scripts/qforge compare --day "$(date +%F)" --question q001 --runner codex --verbose
```

Process presentation for one completed run:

```bash
./scripts/qforge process-visual --run-dir runs/2026-03-15/q001_hops_per_day/claude/opus/run-004 --verbose
```
