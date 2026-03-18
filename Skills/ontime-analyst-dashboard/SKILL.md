---
name: ontime-analyst-dashboard
description: Builds `visual.html` dashboards for qforge questions in either static or dynamic mode, with optional Leaflet maps when geography materially improves the analysis.
---

Use this skill for `visual.html` artifacts and dashboard generation in this repository.

## When to use

- Any prompt under `prompts/` that requires `visual.html`
- Live browser dashboards that fetch from the MCP OpenAPI endpoint
- Self-contained benchmark artifacts that must render without live fetches

## First decision: mode

- If question metadata says `visual_mode: dynamic`, use [references/dynamic-mode.md](references/dynamic-mode.md)
- If question metadata says `visual_mode: static`, use [references/static-mode.md](references/static-mode.md)
- If no mode is declared, treat the question as dynamic for backward compatibility

## Core contract

- Never hardcode KPIs or chart values; derive them from parsed data
- Treat the saved SQL and the selected visual mode as part of the page contract
- Additional browser queries are allowed only in dynamic mode, and only for explicit enrichment or drill-down
- Keep every extra query explicit and visible in a query ledger when using dynamic mode
- Prefer dataset-native dimensions and lookup tables when enrichment fields are needed
- Use optional chaining (`?.`) and nullish coalescing (`??`) in client-side JS
- Normalize temporal fields explicitly; if a ClickHouse `Date` arrives as ISO timestamp, derive `YYYY-MM-DD` in JS
- Keep implementation detail in the skill and references, not in question prompts
- Treat question prompts as business intent and analytical behavior, then supply HTML structure and rendering details from this skill
- Keep runtime failures scoped to the component that failed; a map-render failure must not be reported as a primary-query failure
- If a dashboard supports row/table drill-down, keep selection state explicit and redraw only the dependent detail/map panels while leaving top-level summary KPIs anchored unless the prompt says otherwise

## Presentation guidance

- Use the aviation design system from [references/theme.md](references/theme.md)
- Use layout guidance from [references/patterns.md](references/patterns.md)
- If the question requests layout editing, use [references/edit-mode.md](references/edit-mode.md)

## Technical rules

- Leaflet JS/CSS from CDN is allowed for `html_map` dashboards
- Dynamic mode stores JWE in browser `localStorage` key `OnTimeAnalystDashboard::auth::jwe`
- Dynamic mode loads stored JWE on page init and persists it after successful entry
- Dynamic-mode endpoint: `https://mcp.demo.altinity.cloud/{JWE_TOKEN}/openapi/execute_query?query=...`
- Never embed real tokens in committed artifacts
- Dynamic mode prefills the SQL textarea with the saved SQL and places JWE/SQL controls in a footer block at the very end of the page
- Static mode embeds its analytical data directly in the page and must not depend on live MCP fetch or token flow
- Surface every browser query in a visible ledger with purpose, status, and row count in dynamic mode
- Handle empty results (`count = 0`, `rows = null`) gracefully in dynamic mode

## Minimal output structure

- Branded header with title and analytical subtitle
- KPI strip derived from fetched data
- Primary visual (map, chart, or heatmap)
- Secondary visual or supporting table
- Filters when data supports slicing
- Query ledger showing data provenance
- Export control for filtered rows

## References

1. [references/theme.md](references/theme.md) — design tokens
2. [references/dynamic-mode.md](references/dynamic-mode.md) — live fetch and render logic
3. [references/static-mode.md](references/static-mode.md) — self-contained artifact rules
4. [references/maps.md](references/maps.md) — Leaflet and geography
5. [references/patterns.md](references/patterns.md) — layout shapes
