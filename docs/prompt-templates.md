# Prompt Templates And Variables

This document describes how qforge assembles shared prompt assets and which template variables are available in those assets.

## Prompt Assets

Shared prompt assets live under [`/Users/bvt/work/ExploringDatabyLLMs/prompts`](/Users/bvt/work/ExploringDatabyLLMs/prompts):

- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common.md)
  - shared qforge and dataset-scope guidance used by both SQL and presentation phases
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_report.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_report.md)
  - SQL-only rules for `analysis_mode: json_artifact`
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_report_templates.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_report_templates.md)
  - SQL-only rules for `analysis_mode: template_files` and `analysis_mode: manual_templates`
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual.md)
  - shared presentation and visual rules used by all presentation prompts
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_dynamic.md)
  - dynamic `visual.html` requirements
- [`/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_static.md`](/Users/bvt/work/ExploringDatabyLLMs/prompts/common_visual_static.md)
  - static `visual.html` requirements

Question-specific files such as `prompts/qXXX.../report_prompt.md` and `visual_prompt.md` should contain task logic, not repeated dataset boilerplate.

## Prompt Composition

Prompt builders are implemented in [`/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go).

Current composition order:

- SQL phase:
  - `json_artifact`: `common.md` + `common_report.md` + question `report_prompt.md`
  - `template_files` / `manual_templates`: `common.md` + `common_report_templates.md` + question `report_prompt.md`
- Presentation phase: `common.md` + `common_visual.md` + mode-specific visual asset + question `visual_prompt.md`

The shared `common.md` file is rendered in both phases, so any template variables used there must be available to both builders.

## Variable Sources

Template variables are rendered in [`/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go`](/Users/bvt/work/ExploringDatabyLLMs/internal/prompts/prompts.go).

Dataset-related variables:

- `{{dataset_name}}`
  - Dataset/database name to mention in shared prompts.
  - Source: `DatasetConfig.DefaultDatabase`, falling back to `DatasetConfig.Name`.

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
- `{{question_prompt_md}}`
  - Source: `question.Prompt`.
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

`{{metric.<name>}}` placeholders are available only in `analysis_mode: json_artifact`, because that mode persists `metrics` from the provider artifact.

## Example Dataset Mapping

For `ontime`, the source config is [`/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime/mcp.yaml`](/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime/mcp.yaml).

Dataset semantic-layer documentation remains in [`/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime/semantic_layer.md`](/Users/bvt/work/ExploringDatabyLLMs/datasets/ontime/semantic_layer.md), but prompt generation no longer injects that file directly.

OnTime prompt generation now references the [`ontime-semantic-layer`](../Skills/ontime-semantic-layer/SKILL.md) skill for schema inspection, joins, and airport-dimension semantics.
