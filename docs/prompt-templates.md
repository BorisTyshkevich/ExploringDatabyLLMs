# Prompt Templates And Variables

This document describes how qforge assembles shared prompt assets and which template variables are available in those assets.

## Prompt Assets

Shared prompt assets live under [`/Users/bvt/work/ExploringDatabyLLMs/prompts`](/Users/bvt/work/ExploringDatabyLLMs/prompts):

- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md)
  - shared qforge and dataset-scope guidance used by both SQL and presentation phases
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_sql.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_sql.md)
  - SQL-only rules such as schema inspection, self-verification, and the fenced `sql` output contract
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_presentation.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_presentation.md)
  - report/template rules for the presentation phase
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual.md)
  - shared visual rules used by all presentation prompts
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md)
  - dynamic `visual.html` requirements
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_static.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_static.md)
  - static `visual.html` requirements

Question-specific files such as `prompts/qXXX.../prompt.md`, `report_prompt.md`, and `visual_prompt.md` should contain task logic, not repeated dataset boilerplate.

## Prompt Composition

Prompt builders are implemented in [`/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go).

Current composition order:

- SQL phase: `common.md` + `common_sql.md` + question `prompt.md`
- Presentation phase: `common.md` + `common_presentation.md` + `common_visual.md` + mode-specific visual asset + question `visual_prompt.md`

The shared `common.md` file is rendered in both phases, so any template variables used there must be available to both builders.

## Variable Sources

Template variables are rendered in [`/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go).

Dataset-related variables:

- `{{dataset_primary_table}}`
  - The dataset's main fact table.
  - Source: `DatasetConfig.PrimaryTable`.
- `{{dataset_constraints_md}}`
  - Markdown bullet list describing allowed and forbidden table usage.
  - Built from `DatasetConfig.PrimaryTable` and `DatasetConfig.ForbiddenTables`.

Presentation/report variables:

- `{{question_title}}`
  - Source: `question.Meta.Title`.
- `{{visual_type}}`
  - Source: `question.Meta.VisualType`.
- `{{visual_mode}}`
  - Source: `question.Meta.VisualMode`.
- `{{result_columns_csv}}`
  - Source: `strings.Join(result.Columns, ", ")`.
- `{{saved_sql}}`
  - Source: the saved `query.sql` text passed into `BuildPresentationPrompt`.
- `{{report_prompt_md}}`
  - Source: `question.ReportPrompt`.
- `{{visual_prompt_md}}`
  - Source: `question.VisualPrompt`.
- `{{report_placeholders}}`
  - Literal CSV string listing which report placeholders the model may emit.

## Report Placeholders

These placeholders are not expanded when the prompt is built. They are passed through so the model can emit a template-style `report`, and qforge fills them later when rendering `report.md`.

- `{{row_count}}`
  - Total number of rows in `result.json`.
- `{{generated_at}}`
  - Timestamp from `result.json`.
- `{{columns_csv}}`
  - Comma-separated output column list.
- `{{question_title}}`
  - Question title at report-render time.
- `{{data_overview_md}}`
  - Markdown summary block derived from `result.json`.
- `{{result_table_md}}`
  - Markdown table derived from `result.json`.

## Example Dataset Mapping

For `ontime_v2`, the source config is [`/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime_v2/mcp.yaml`](/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime_v2/mcp.yaml).

Current values:

- `primary_table: default.ontime_v2`
- `forbidden_tables: default.ontime`

`{{dataset_constraints_md}}` is rendered from that config into Markdown bullets such as:

- `Use default.ontime_v2 as the primary fact table.`
- `Do not reference default.ontime.`
