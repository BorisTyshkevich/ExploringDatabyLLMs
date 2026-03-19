# qforge Processing Logic

This document describes the current end-to-end processing flow for `qforge`: what the model produces, what the harness owns, which artifacts are authoritative, and how the optional visual phase is assembled.

## Overview

`qforge` runs in two phases:

1. analysis phase
   - the provider inspects schema, reasons about the question, and writes `answer.raw.json`
   - the harness loads `answer.raw.json`, extracts SQL and report inputs, executes SQL itself, and renders `report.md`
2. visual phase
   - optional, controlled by `run --with-visual` or by `process-visual`
   - the provider receives saved analytical artifacts and generates only `visual.html`

The model never executes the final SQL. The harness always executes SQL and writes the canonical `result.json`.

## Phase 1: Analysis

Phase 1 prompt assembly is implemented in [`/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go).

Current composition order:

- `prompts/common.md`
- `prompts/common_sql.md`
- question `prompt.md`

The shared dataset semantic layer is loaded from `datasets/<dataset>/semantic_layer.md` when present and inlined into the prompt.

### Provider contract

The provider must write `answer.raw.json` in the run directory.

`answer.raw.json` must contain raw JSON bytes only with this shape:

```json
{
  "sql": "-- one SQL statement",
  "report_markdown": "# Report template using built-in placeholders and {{metric.<name>}}",
  "metrics": {
    "summary_facts": [],
    "named_values": {},
    "named_lists": {}
  }
}
```

Important rules:

- `sql` is the only executable query
- `report_markdown` is a template, not a filled report
- `metrics` carries report-only derived facts
- stdout is diagnostic only and is not used for phase-1 artifact loading

### Harness behavior

The analysis flow is orchestrated in [`/Users/bvt/work/ExploringDatabyLLMs/internal/cli/cli.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/cli/cli.go).

After the provider returns, qforge:

1. saves provider stdout/stderr logs
2. reads `answer.raw.json`
3. validates the JSON artifact
4. writes normalized `analysis.json`
5. writes `query.sql`
6. writes `report.template.md`
7. executes the SQL itself
8. writes `result.json`
9. renders final `report.md`

Phase 1 fails if:

- `answer.raw.json` is missing
- `answer.raw.json` is not valid raw JSON
- required fields like `sql` or `report_markdown` are empty
- the report template uses unsupported placeholders
- SQL execution fails

## Report Rendering

Report rendering is implemented in [`/Users/bvt/work/ExploringDatabyLLMs/internal/render/render.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/render/render.go).

The renderer owns the final Markdown assembly. The provider supplies only:

- `report_markdown`
- `metrics`

Built-in placeholders currently include:

- `{{row_count}}`
- `{{generated_at}}`
- `{{columns_csv}}`
- `{{question_title}}`
- `{{data_overview_md}}`
- `{{result_table_md}}`

Metric placeholders use:

- `{{metric.<name>}}`

`metrics.named_values` is the source for metric placeholder substitution.

## Phase 2: Visual

Phase 2 prompt assembly is also implemented in [`/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go).

Current composition order:

- `prompts/common.md`
- `prompts/common_presentation.md`
- `prompts/common_visual.md`
- mode-specific visual asset
- question `visual_prompt.md`

### Inputs

The visual provider receives these saved artifacts as prompt context:

- `analysis.json`
- `query.sql`
- `report.template.md`
- `report.md`
- `result.json`

The visual phase should treat:

- `result.json` as the authoritative data result
- `query.sql` as the authoritative executed query
- `analysis.json` as authoritative structured analysis context

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

Typical run artifacts under `runs/YYYY-MM-DD/<question>/<runner>/<model>/run-XXX/`:

- `prompt.sql.md`
- `answer.sql.raw.md`
- `answer.raw.json`
- `analysis.json`
- `query.sql`
- `report.template.md`
- `result.json`
- `report.md`
- `prompt.presentation.md`
- `answer.presentation.raw.md`
- `visual.html`
- `manifest.json`
- `stdout.log`
- `stderr.log`

Source-of-truth artifacts:

- analysis artifact: `answer.raw.json`
- normalized analysis snapshot: `analysis.json`
- executed SQL: `query.sql`
- canonical result: `result.json`
- final rendered report: `report.md`
- final visual artifact: `visual.html`
