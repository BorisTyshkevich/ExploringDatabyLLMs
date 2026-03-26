# qforge Processing Logic

This document describes the current end-to-end processing flow for `qforge`: what the model produces, what the harness owns, which artifacts are authoritative, and how the optional visual phase is assembled.

## Overview

`qforge` runs in three steps:

1. analysis phase
   - the exact contract is selected by the question's `analysis_mode`
   - `multi_query`: the provider writes `answer.raw.json`
   - `template_files`: the provider writes `query.sql` and `report.template.md`
   - `--manual`: runtime staging flag that skips provider invocation for the selected analysis contract; it works with both `multi_query` and `template_files`
   - during a normal `qforge run`, the harness loads the saved analysis artifacts, executes SQL itself, and renders `report.md`
2. report and review phase
   - `qforge review` is the manual-run completion path for an existing run directory
   - it loads the saved analysis artifacts, regenerates `report.md`, and writes `review.md`
   - `qforge process-presentation` is report-only materialization from the same saved artifacts and does not write `review.md`
   - during a normal `qforge run`, the harness then performs the review step automatically and writes `review.md`
3. visual phase
   - optional, controlled by `run --with-visual` or by `visual`
   - the provider receives the saved primary SQL artifact (`query.sql` or the selected `queries/<id>.sql`) plus a harness-generated visual input summary and generates only `visual.html`

The model never executes the final SQL. The harness always executes SQL and writes the canonical `result.json`.

## Phase 1: Analysis

Phase 1 prompt assembly is implemented in [`/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go).

Current composition order:

- `prompts/common.md`
- `prompts/common_report_multi_query.md` for `multi_query`
- `prompts/common_report_templates.md` for `template_files`
- question `report_prompt.md`

`--manual` does not change prompt assembly. It only skips the provider call and leaves the selected analysis contract in place for later human completion.

Shared prompt assets now reference dataset-specific skills directly. For OnTime, schema inspection and join guidance come from the `ontime-semantic-layer` skill rather than an inlined `semantic_layer.md` block.

### Provider contract

Question metadata selects one of these analysis contracts via `analysis_mode`.

#### `multi_query`

The provider must write `answer.raw.json` in the run directory.

`answer.raw.json` must contain raw JSON bytes only with this shape:

```json
{
  "subquestions": [
    {
      "id": "main, q1, q2, or q0",
      "answer_markdown": "Direct prose answer",
      "sql": "-- one proof query for this section"
    }
  ]
}
```

Important rules:

- top-level `### main` / `### qN` sections are the source of truth for ids and ordering
- each section carries one direct prose answer and one proof query
- if a prompt has no `### main` and no `### qN`, qforge treats the full prompt as one implicit `q0` section
- the visual phase uses `queries/main.sql` when `main` exists, otherwise the sole section query, otherwise `queries/q1.sql`
- stdout is diagnostic only and is not used for phase-1 artifact loading

#### `template_files`

The provider must write these files in the run directory:

- `query.sql`
- `report.template.md`

Important rules:

- `query.sql` is the only executable query
- `report.template.md` is a template, not a filled report
- qforge does not parse a phase-1 JSON artifact in this mode
- metric placeholders are not supported in this mode because no metrics artifact is collected

#### `--manual`

`--manual` is a runtime staging flag, not a separate analysis mode.

When it is set, `qforge run` saves the analysis prompt but does not call a provider. The human later uses that prompt in ChatGPT, Claude, or another external UI and saves:

- for `template_files`: `query.sql` and `report.template.md`
- for `multi_query`: `answer.raw.json`

`qforge review` continues from those saved files when you want to finish a staged manual run.
`qforge process-presentation` continues from those saved files when you only need report-only materialization.

### Harness behavior

The analysis flow is orchestrated in [`/Users/bvt/work/ExploringDatabyLLMs/internal/cli/cli.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/cli/cli.go).

After the provider returns, qforge:

1. saves provider stdout/stderr logs
2. reads the saved analysis artifact for the selected mode
3. validates the JSON artifact
4. writes normalized `analysis.json`
5. writes normalized proof-query artifacts
6. executes the SQL itself
7. writes `result.json` or named per-query result summaries
8. writes `visual_input.json`
9. renders final `report.md`
10. during `qforge run`, calls the review model and writes `review.md`
11. during `qforge review`, performs the same review step from an existing run directory and writes `review.md`

Phase 1 fails if:

- the required saved analysis artifact is missing
- `answer.raw.json` is not valid raw JSON when `analysis_mode: multi_query`
- required section fields like `id`, `answer_markdown`, or `sql` are empty
- the report template uses unsupported placeholders
- SQL execution fails

## Report Rendering

Report rendering is implemented in [`/Users/bvt/work/ExploringDatabyLLMs/internal/render/render.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/render/render.go).

The renderer owns the final Markdown assembly. The provider supplies only:

- `report.template.md` in template modes, or ordered subquestion answers in multi-query mode

Built-in placeholders currently include:

- `{{row_count}}`
- `{{generated_at}}`
- `{{columns_csv}}`
- `{{question_title}}`
- `{{data_overview_md}}`
- `{{result_table_md}}`

## Phase 2: Visual

Phase 2 prompt assembly is also implemented in [`/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go).

Current composition order:

- `prompts/common.md`
- `prompts/common_visual.md`
- mode-specific visual asset
- question `visual_prompt.md`

### Inputs

The visual provider receives these saved artifacts as prompt context:

- the saved primary SQL artifact (`query.sql` for single-query modes, or the selected `queries/<id>.sql` for `multi_query`)
- `visual_input.json`
- `result.json` in static mode only

The visual phase should treat:

- the saved primary SQL artifact as the authoritative query for the page
- `visual_input.json` as the compact data-shape summary
- `result.json` as the authoritative embedded data source in static mode

### Output

The provider generates only:

- `visual.html`

It must not regenerate SQL or report artifacts.

### Validation

After `visual.html` is written, qforge validates it in two layers:

1. contract validation
2. browser validation with `chromedp`

Dynamic visuals may also perform live browser-side fetches using the same MCP token already used by qforge.

## Artifact Summary

Typical run artifacts under `YYYY-MM-DD/<question>/<runner>/<model>/run-XXX/`:

- `prompt.report.md`
- `answer.report.raw.md`
- `answer.raw.json`
- `analysis.json`
- `prompt.review.md`
- `answer.review.raw.md`
- `review.md`
- `query.sql`
- `queries/main.sql`
- `queries/q1.sql`
- `report.template.md`
- `result.json`
- `visual_input.json`
- `report.md`
- `prompt.visual.md`
- `answer.presentation.raw.md`
- `visual.html`
- `manifest.json`
- `stdout.log`
- `stderr.log`

Source-of-truth artifacts:

- analysis artifact:
  - `answer.raw.json` for `multi_query`
  - `query.sql` + `report.template.md` for `template_files`
- normalized analysis snapshot: `analysis.json`
- review artifact: `review.md`
- executed SQL:
  - `queries/*.sql` for `multi_query`
  - `query.sql` for `template_files`
- canonical result: `result.json`
- visual grounding summary: `visual_input.json`
- final rendered report: `report.md`
- final visual artifact: `visual.html`
