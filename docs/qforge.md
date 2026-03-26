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

`qforge` uses a three-step design:

1. SQL generation
   - the model is prompted to inspect schema and self-verify its SQL before writing artifacts
   - the exact artifact contract is selected by the question's `analysis_mode`
   - `multi_query`: the model writes `answer.raw.json` with ordered section answers and proof queries keyed by `main`, `qN`, or implicit `q0`
   - `template_files`: the model writes `query.sql` and `report.template.md`
   - `--manual`: runtime staging flag that skips provider invocation for the selected analysis contract; it works with both `multi_query` and `template_files`
2. Report and review materialization
   - `qforge run` executes SQL, renders `report.md`, and makes a separate review-model call
   - `qforge review` is the manual-run completion path for an existing run directory; it materializes `report.md` and writes `review.md`
   - `process-presentation` is report-only materialization from saved analysis artifacts and does not write `review.md`
   - the reviewer writes `review.md` with `Verdict: PASS|WARN|FAIL`
   - `WARN` means the run is materially usable but has limited issues and is marked non-clean
   - `qforge run` stops before visual generation if the review verdict is `FAIL`
3. Optional presentation generation
   - the model writes final `html`
   - the final `report.md` is rendered by the harness from the saved analysis artifact plus JSON-derived sections
   - `visual.html` is model-authored final output and is not patched by `qforge`
   - `visual.html` is validated first against the visual contract and then in a headless browser via `chromedp`
   - this can be done later with `visual` or immediately with `run --with-visual`

The model should not emit result rows directly.

The harness always executes the final SQL itself and writes `result.json`.

## Prompt Assets

Prompt assembly is split into shared and phase-specific assets under [`/Users/bvt/work/ExploringDatabyLLMs/prompts`](/Users/bvt/work/ExploringDatabyLLMs/prompts):

- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md)
  - shared qforge and dataset-scope guidance used by both SQL and presentation phases
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_report_multi_query.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_report_multi_query.md)
  - SQL-only rules such as schema inspection, self-verification, and the `answer.raw.json` multi-query contract
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_report_templates.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_report_templates.md)
  - SQL-only rules for direct `query.sql` plus `report.template.md` output
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_review.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_review.md)
  - shared analysis-review contract for the mandatory post-analysis reviewer
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual.md)
  - shared presentation and `visual.html` rules
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_html.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_html.md)
  - HTML-target presentation requirements
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_react.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_react.md)
  - React-target presentation requirements
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md)
  - dynamic `visual.html` requirements
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_static.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_static.md)
  - static `visual.html` requirements

Question-specific files such as `prompts/qXXX.../report_prompt.md` and `visual_prompt.md` should contain task logic, not repeated dataset boilerplate.

Question metadata may also declare `visual_mode`:

- `dynamic`
  - live browser dashboard that uses saved SQL plus browser-side MCP fetch
- `static`
  - self-contained benchmark artifact with embedded analytical data

If `visual_mode` is absent, qforge treats the question as `dynamic` for backward compatibility.

Question metadata may also declare `presentation_target`:

- `html`
  - model writes final `visual.html`
- `react`
  - model writes source under `visual_src/`, qforge builds it, and the built browser artifact is still published as `visual.html`

If `presentation_target` is absent, qforge treats the question as `html` for backward compatibility.

Template variables, prompt composition, and dataset mapping are documented in [`/Users/bvt/work/ExploringDatabyLLMs/docs/prompt-templates.md`](/Users/bvt/work/ExploringDatabyLLMs/docs/prompt-templates.md).

Practical guidance on writing question prompts is documented in [`/Users/bvt/work/ExploringDatabyLLMs/docs/writing-prompts.md`](/Users/bvt/work/ExploringDatabyLLMs/docs/writing-prompts.md).

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

Run one question with an explicit separate reviewer:

```bash
./scripts/qforge run -q q001 -r claude --model sonnet --review-runner codex --review-model gpt-5.4 -v
```

Stage a manual run without invoking the provider:

```bash
./scripts/qforge run -q q001 -r claude --manual -v
```

The same `--manual` flag stages a run for either analysis contract without changing that contract.
After the provider artifacts are saved into the run directory, finish the run with `qforge review`.

Finish an existing manual run:

```bash
./scripts/qforge review --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
```

Optionally generate visual output later:

```bash
./scripts/qforge visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
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

Process report-only artifacts for an existing run:

```bash
./scripts/qforge process-presentation --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
```

Process visual for an existing run:

```bash
./scripts/qforge visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
```

Process report and visual but skip all visual validation:

```bash
./scripts/qforge visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 --skip-visual-validation -v
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
  - default when omitted: `claude/opus`, `claude/sonnet`, `codex/gpt-5.4`
- `--model`
  - optional, repeatable
  - override the default model for the selected runner
  - current defaults: `codex -> gpt-5.4`, `claude -> sonnet`, `gemini -> gemini-3.1-pro-preview`
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
- `--analysis-mode`
  - optional
  - may override only for non-`multi_query` questions
  - cannot switch to or from `multi_query`
  - does not control whether the run is staged manually
- `--presentation-target`
  - optional
  - override the question presentation target with `html` or `react`
  - useful for benchmarking the same question in both presentation formats without editing prompt metadata
- `--manual`
  - optional
  - runtime staging flag that skips provider invocation and leaves the selected analysis contract unchanged
  - works with both `template_files` and `multi_query`
- `-m`
  - shorthand for `--manual`
- `--verbose`
  - optional
  - print phase-level progress logs and provider subprocess timing
- `-v`
  - shorthand for `--verbose`
- `--with-visual`
  - optional
  - after SQL succeeds and `report.md` is rendered, make a second independent provider call for `visual.html`
  - this is equivalent in behavior to running `visual` after a successful run
- `--skip-visual-validation`
  - optional
  - skip both contract validation and browser validation for `visual.html`
- `--skip-browser-live-fetch`
  - optional
  - keep contract validation and browser smoke checks, but skip the token-entry and live fetch interaction

What `run` does:

1. resolves question metadata
2. selects one or more providers
3. resolves the effective analysis mode from question metadata plus any allowed CLI override
4. builds the SQL prompt for each selected provider
5. either invokes the provider or stages a manual run, depending on `--manual`
6. loads the saved analysis artifact for that mode
7. executes SQL directly against the OpenAPI endpoint when analysis artifacts are available
8. writes canonical `result.json`
9. writes `visual_input.json`
10. renders `report.md` from the saved report template
11. writes `prompt.visual.md` for visual-capable questions so visual generation can be run manually later
12. writes `manifest.json`
13. optionally makes a second independent provider call for `visual.html` when `--with-visual` is set

What `run` does not do:

- it does not produce `visual.html` unless `--with-visual` is set
- it does not complete a staged manual run; use `qforge review` after saving the provider artifacts into the run directory
- it may still prebuild `prompt.visual.md` for visual-capable questions
- use `qforge visual` for HTML generation later

Exception:

- when `--with-visual` is set, `run` performs the SQL phase first and then a separate presentation call for successful runs

### `qforge review`

Usage:

```bash
./scripts/qforge review --run-dir <path> [flags]
```

What `review` does:

- loads `manifest.json` and the saved analysis artifacts from an existing run directory
- regenerates `report.md` from the saved analysis output
- writes `review.md` using the same review contract as `qforge run`
- is the recommended completion step after `qforge run --manual`
- does not generate `visual.html`

### `qforge process-presentation`

Usage:

```bash
./scripts/qforge process-presentation --run-dir <path> [flags]
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
- `--verbose`
  - optional
  - print phase-level progress logs

What `process-presentation` does:

- loads `manifest.json` and the saved analysis artifacts for the question's declared `analysis_mode`, whether they were provider-generated or staged with `--manual`
- extracts SQL and report inputs from the saved analysis artifact
- executes SQL directly against the OpenAPI endpoint
- rewrites:
  - `query.sql`
  - `analysis.json`
  - `result.json`
  - `visual_input.json`
  - `prompt.visual.md` when the question has visual artifacts
  - `report.template.md`
  - `report.md`
- does not invoke a provider again
- does not write `review.md`
- does not generate `visual.html`

### `qforge visual`

Usage:

```bash
./scripts/qforge visual --run-dir <path> [flags]
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

What `visual` does:

- loads `manifest.json`, the saved primary SQL artifact (`query.sql` or the selected `queries/<id>.sql`), and `visual_input.json` from an existing run
- for static mode, also loads `result.json`
- rebuilds the presentation prompt from question metadata and the saved artifacts
- invokes the original provider again for `html`
- allows the provider to author visual-only lookup queries directly into `visual.html`; those queries are self-verified in the visual pass and are not part of reviewed analysis artifacts
- validates `visual.html` in two stages unless `--skip-visual-validation` is set:
  - contract validation against the shared visual rules
  - browser validation using `chromedp`
- for dynamic dashboards, browser validation automatically attempts the live MCP fetch path when a token is available unless `--skip-browser-live-fetch` is set
- writes:
  - `prompt.visual.md`
  - `answer.presentation.raw.md`
  - `visual.html`

Manual visual generation in ChatGPT.com:

1. First make sure the run already has up-to-date analysis artifacts.
   - For staged manual runs, use `qforge review --run-dir <path>` first.
   - If you only want report-only materialization, use `qforge process-presentation --run-dir <path>` first.
2. Open `<run-dir>/prompt.visual.md`.
   - This is the complete visual-generation prompt assembled by qforge from question metadata plus the saved SQL and visual grounding artifacts.
3. Paste the full contents of `prompt.visual.md` into ChatGPT.com.
4. Ask ChatGPT.com to return only the final HTML artifact.
   - For normal `html` presentation targets, the response should be the full `visual.html` document.
5. Save the returned HTML into `<run-dir>/visual.html`.
6. If you want qforge's automated validation, run `qforge visual --run-dir <path>` instead of this manual path.

Files required before a manual visual pass:

- always:
  - `prompt.visual.md`
- produced by the earlier qforge materialization step:
  - `visual_input.json`
  - the primary saved SQL artifact:
    - `query.sql` for `template_files`
    - `queries/<id>.sql` for `multi_query`
- for static visual mode:
  - `result.json`

Practical notes:

- In the normal manual ChatGPT.com workflow, `prompt.visual.md` is the only file you need to paste into the external UI.
- You do not need to separately paste `visual_input.json` or `query.sql` when `prompt.visual.md` is already current; qforge has already embedded the needed context into that prompt.
- This manual path is most straightforward for `presentation_target: html`.
- If a question uses `presentation_target: react`, prefer `qforge visual` because the automated path can handle source-artifact extraction and build steps.

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
  - explicit MCP URL ending in `/http` for compare-report provider config
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
- fetches deferred performance metrics from `system.query_log` using `clickhouse-client --connection demo` and `log_comment`
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
2026-01-01 00:00:00 opus run question=q002 runners=claude
2026-01-01 00:00:00 opus run question=q002 runner=claude model=opus dataset=ontime
2026-01-01 00:00:00 opus out_dir=... visual=false timeout_sec=900
2026-01-01 00:00:00 opus phase=sql_generation status=started
2026-01-01 00:00:00 opus provider=claude phase=start ...
2026-01-01 00:01:29 opus provider=claude phase=done status=ok elapsed=1m29.067s
2026-01-01 00:01:29 opus phase=sql_execution status=ok row_count=195
```

## Artifacts

Each run is stored under:

```text
YYYY-MM-DD/<question-slug>/<runner>/<model>/run-XXX/
```

Typical SQL-only run artifacts:

- `prompt.report.md`
- `answer.report.raw.md`
- `answer.raw.json`
- `report.template.md`
- `query.sql`
- `result.json`
- `visual_input.json`
- `manifest.json`
- `stdout.log`
- `stderr.log`

When presentation is processed later with `qforge visual`:

- `prompt.visual.md`
- `answer.presentation.raw.md`
- `visual.html`

Typical analysis-phase artifacts now include:

- `analysis.json`
- `query.sql`
- `visual_input.json`

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
- Presentation is a separate explicit step by default. Use `qforge review` to finish a staged manual run, `qforge process-presentation` if you only need report-only materialization, and `run --with-visual` if you want `run` to also make the follow-up presentation call.

## Recommended Workflows

SQL/JSON verification only:

```bash
./scripts/qforge run -q q001 -r claude -v
```

Manual run completion:

```bash
./scripts/qforge run -q q001 -r claude --manual -v
./scripts/qforge review --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
```

Optionally follow with visual generation:

```bash
./scripts/qforge visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
```

Manual visual generation in ChatGPT.com after materialization:

```bash
./scripts/qforge run -q q001 -r claude --manual -v
./scripts/qforge review --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
# open 2026-03-15/q001_hops_per_day/claude/opus/run-004/prompt.visual.md
# paste it into ChatGPT.com
# save the returned HTML to 2026-03-15/q001_hops_per_day/claude/opus/run-004/visual.html
```

Report-only materialization for an existing run:

```bash
./scripts/qforge process-presentation --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
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

Process visual for one completed run:

```bash
./scripts/qforge visual --run-dir 2026-03-15/q001_hops_per_day/claude/opus/run-004 -v
```
