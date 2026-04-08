# Writing Good qforge Prompts

This manual explains how to write question prompts that produce reliable qforge runs.

The short version:

- write the question like a BI analyst or manager would ask it
- keep execution mechanics in shared prompt assets, not in the question text
- be explicit about business decisions the model must answer
- be equally explicit about what the model must not invent
- for visuals, describe the business story and intended dashboard behavior, not HTML mechanics

## Prompt Structure

Each question lives under `prompts/qXXX_.../` and usually includes:

- `meta.yaml`
- `report_prompt.md`
- `visual_prompt.md`
- optional `subquestions.yaml` for `analysis_mode: structured`

Use them for different jobs:

- `meta.yaml`
  - selects the technical contract
  - examples: `analysis_mode`, `visual_mode`, `visual_type`
- `report_prompt.md`
  - defines the business question and analytical boundaries
- `visual_prompt.md`
  - defines what the human-facing dashboard should help the reader see and do
- `subquestions.yaml`
  - lists required business questions when one prompt must answer several distinct points

Do not move technical contract details from shared prompt assets into question prompts unless the question genuinely needs an exception.

## Write `report_prompt.md` Like A BI Brief

Good `report_prompt.md` files sound like an analyst request, not a warehouse implementation note.

Prefer:

- the business problem
- the scope
- the grain
- the success criteria
- the interpretation the report must deliver

Avoid:

- function names
- exact CTE instructions
- row-contract language
- implementation tie-break rules
- internal harness mechanics

A good report prompt should answer these questions clearly:

- What business issue is being investigated?
- What slice of the data is in scope?
- What is the analytical grain?
- What counts as credible signal versus noise?
- What business questions must be answered directly in prose?

### Good pattern

Use wording like:

- “Find the worst recurring hotspots”
- “Focus on combinations with enough volume to be credible”
- “State whether the issue is persistent or concentrated in a narrower period”
- “Summarize what the leading hotspots imply operationally”

### Bad pattern

Avoid wording like:

- “Use `quantile(0.9)`”
- “Return `RowType = summary | trend`”
- “Use exactly these columns in exactly this order”
- “Implement deterministic tie-breaking by X then Y”

Those instructions belong in shared scaffolding only if the harness truly requires them.

## Be Explicit About Anti-Drift Rules

Models drift when the prompt leaves key business definitions open.

For report prompts, explicitly constrain:

- analysis window
  - use the most recent 5 years unless the question asks for a different period
- business definition
  - do not redefine “hotspot”, “leadership”, “persistent”, or similar terms unless asked
- ranking logic
  - do not invent composite scores unless the question asks for one
- noise filtering
  - allow reasonable volume filters, but frame them as business guardrails

This is usually enough:

- “Use the most recent 5 years unless the question explicitly asks otherwise.”
- “Do not invent a custom score or narrower time window.”
- “Exclude low-volume noise before identifying leaders.”

## Use `subquestions.yaml` When The Prompt Really Has Several Questions

If one business prompt actually contains several required answers, do not rely on one free-form paragraph to cover them all.

Use `subquestions.yaml` when:

- the report must answer several distinct questions
- each question may need different evidence
- the visual narrative needs several analytical panels

Good examples:

- “Which hotspot is worst?”
- “Is it persistent or concentrated?”
- “What broader pattern do the leaders suggest?”

This gives qforge a validation target without forcing question text to become technical.

## Write `visual_prompt.md` For The Reader, Not The Harness

`visual_prompt.md` should describe the human-facing artifact:

- what the page should emphasize
- what visuals should exist
- what interactions matter
- what the reader should understand after viewing it

It should not restate shared runtime mechanics already covered by:

- `common_visual.md`
- `common_visual_dynamic.md`
- `common_visual_static.md`
- `common_visual_structured.md`
- the `ontime-analyst-dashboard` skill

### For dynamic visuals

Assume the shared contract already handles:

- JWE token storage
- footer controls
- date selector
- primary and supporting SQL editors
- live fetch
- query ledger
- validation expectations

So the question-specific visual prompt should focus on:

- which query is the primary dashboard view
- which supporting queries matter
- whether the page needs extra visual-only questions that may lead the model to author second-pass lookup or enrichment queries
- which charts or panels should exist
- what narrative should anchor the page
- what degraded behavior is acceptable if a supporting or lookup query fails

Additional question headings in `visual_prompt.md` are prose-only author guidance for the presentation phase.

- qforge passes those sections through unchanged in `BuildVisualPrompt`
- qforge does not parse them into structured visual subquestions or saved analysis artifacts
- any SQL they trigger is visual-only second-pass SQL owned by the authored dashboard, not `queries/*.sql`, `answer.raw.json`, or `review.md`
- tests should cover shared prompt-contract guarantees and stable wording boundaries, not exact per-panel behavior implied by those prose sections

### Good pattern

Use wording like:

- “Use `worst_hotspot` as the primary saved SQL”
- “Use `persistence` as a supporting query for the monthly trend panel”
- “Ask an additional visual-only question for the map coordinate panel”
- “Show supporting queries in the ledger when used”
- “The dashboard does not need to mirror `report.md`”
- “Use the prose answers as framing, and fetched query results for charts and tables”

### Bad pattern

Avoid wording like:

- “Use this exact HTML structure”
- “Embed JSON in a script tag” when the visual is dynamic
- “Do not query” when the page is supposed to live-fetch
- “Show only the first row” unless the visual is intentionally a monitoring page

## Dynamic Visuals Vs Static Visuals

Choose `visual_mode` based on what the final artifact is supposed to do.

Use `dynamic` when:

- the dashboard should fetch live data in the browser
- the page should expose token controls and saved SQL
- the reader benefits from ledger-backed provenance
- supporting or drill-down queries are useful

Use `static` when:

- the artifact must render fully offline
- all analytical data should be embedded in the HTML
- the page is benchmark-style and non-interactive beyond local UI state

Do not mix the two in the question prompt. If the page should live-fetch, write a dynamic visual prompt and let the shared dynamic contract do its job.

## HTML Target Vs React Target

Choose `presentation_target` based on what the model should author.

Use `presentation_target: html` when:

- the model should write final `visual.html`
- you want the simplest presentation path
- you do not need React source as a benchmark artifact

Use `presentation_target: react` when:

- the model should write source under `visual_src/`
- qforge should build that source into the final browser artifact
- you want to compare model authoring of React against direct HTML generation

Keep `visual_mode` and `presentation_target` separate:

- `visual_mode` controls runtime semantics such as live fetch vs embedded data
- `presentation_target` controls whether the model writes HTML directly or React source

## What To Keep Out Of Question Prompts

Do not repeat:

- MCP endpoint details
- JWE storage key details
- browser validation rules
- report placeholder rules
- JSON artifact syntax
- provider file-writing rules

Those belong in shared prompt assets and docs.

Question prompts should carry business intent, not plumbing.

## A Simple Checklist

Before finalizing a new question prompt, verify:

- The business question is understandable without reading qforge code.
- The analytical grain is clear.
- Required business answers are explicit.
- Anti-drift constraints are explicit.
- The visual prompt describes the page the reader needs, not the implementation internals.
- Dynamic vs static mode matches the intended artifact behavior.
- Nothing in the question prompt contradicts the selected shared contract.

## `q003` Lesson

The main `q003` failure mode was not SQL syntax. It was prompt-shape mismatch:

- the report prompt asked a BI question but left too much freedom in answer behavior
- the first visual version reduced the page to preview rows, which was fine for monitoring but wrong for a real dashboard
- the correct fix was to keep first-pass monitoring artifacts simple and let the second pass build a real dynamic dashboard from the saved primary query plus named supporting queries

That pattern is the right default for similar multi-question, dashboard-oriented prompts:

- first pass: answer the business questions and save the proof queries
- report: compact monitoring artifact
- second pass: dynamic dashboard with narrative plus live query-backed visuals
