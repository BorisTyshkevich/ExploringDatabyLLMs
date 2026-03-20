---
name: ontime-analyst-dashboard
description: Builds visual dashboards in either static or dynamic mode, with optional Leaflet maps.
---

## First decision: mode

- Use dynamic mode for browser dashboards that execute the saved SQL through the configured tokenized HTTP SQL endpoint and may run explicit enrichment or drill-down queries. See  `references/dynamic-mode.md`
- Use static mode when embedding data into HTML dashboard. see `references/static-mode.md`

## Core contract

- Never hardcode KPIs or chart values; derive them from parsed data
- Treat the saved SQL and the selected visual mode as part of the page contract
- Prefer dataset-native dimensions and lookup tables when enrichment fields are needed
- Use optional chaining (`?.`) and nullish coalescing (`??`) in client-side JS
- Normalize temporal fields explicitly; if a ClickHouse `Date` arrives as ISO timestamp, derive `YYYY-MM-DD` in JS
- Keep implementation detail in the skill and references, not in question prompts
- Treat question prompts as business intent and analytical behavior, then supply HTML structure and rendering details from this skill
- Keep runtime failures scoped to the component that failed; a map-render failure must not be reported as a primary-query failure
- If a dashboard supports row/table drill-down, keep selection state explicit and redraw only the dependent detail/map panels while leaving top-level summary KPIs anchored unless the prompt says otherwise

## Presentation guidance

- Use the shared aviation visual language from `references/theme.md`, including `--navy`, `--sky`, `--teal`, and `--amber`; do not replace it with an ad hoc palette unless the question explicitly asks for it
- Use layout guidance from `references/patterns.md`
- If requested layout editing, use `references/edit-mode.md`
- Use `references/maps.md` when geography materially improves the analysis

## Minimal output structure

- Branded header with title and analytical subtitle
- KPI strip derived from fetched data
- Primary visual (map, chart, or heatmap)
- Secondary visual or supporting table
- Filters when data supports slicing
- SQL query ledger showing data provenance
- Export control for filtered rows
- additional elements required by static or dynamic mode

