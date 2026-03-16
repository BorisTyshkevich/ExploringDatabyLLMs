---
name: ontime-analyst-dashboard
description: Builds `visual.html` dashboards for analytics reports in static or dynamic mode, with browser-side normalization and visualization derived from the provided result set. Use for dashboard-style visual artifacts in this repository.
---

Use this skill for `visual.html` artifacts and dashboard generation in this repository.
Prefer simple contracts that keep the dashboard behavior aligned with the saved SQL and the question prompt.

## When to use

- Any prompt under `prompts/` that requires `visual.html`
- Static benchmark artifacts with embedded CSV data
- Live browser dashboards that fetch from the MCP OpenAPI endpoint
- Geographic or network-style reports where maps materially improve the analysis

## First decision: mode

- If the artifact must be self-contained for benchmark output, use **static mode**
  Read: [references/static-mode.md](references/static-mode.md)
- If the dashboard should fetch live data in the browser, use **dynamic mode**
  Read: [references/dynamic-mode.md](references/dynamic-mode.md)

## Second decision: map or non-map

- If `visual_type` is `html_map`, or geography is the primary analytical view, add **map support**
  Read: [references/maps.md](references/maps.md)
- Otherwise keep the dashboard validator-safe with inline CSS, inline JS, and inline SVG only

## Core contract

- Never hardcode KPIs or chart values; derive them from parsed data
- In dynamic mode, treat the saved SQL prefilled into the page as the primary analytical query
- Additional browser queries are allowed when they serve a concrete enrichment or drill-down purpose and stay within dataset constraints
- Keep every extra query explicit, intentional, and explainable to the user; do not hide follow-up queries behind unexplained UI behavior
- Prefer dataset-native dimensions and lookup tables from the allowed dataset scope when geography or enrichment fields are needed
- Keep JS responsibilities to normalization, filtering, sorting, aggregation, formatting, and visualization
- Use optional chaining (`?.`) and nullish coalescing (`??`) in client-side JS
- Normalize rows before rendering charts or tables
- Normalize temporal fields explicitly before rendering or comparison logic
- If a ClickHouse `Date` column may arrive as an ISO timestamp, derive a stable `YYYY-MM-DD` field in JS and reuse it everywhere
- When cloning dashboard cards from `<template>`, do not rely on duplicate global `id` attributes inside template content; bind JS to scoped live-card selectors or unique IDs instead

## Presentation guidance

- Use the aviation design system and tokens from [references/theme.md](references/theme.md)
- Use report-specific layout guidance from [references/patterns.md](references/patterns.md)
- Use native edit mode from [references/edit-mode.md](references/edit-mode.md)
- Keep one coherent look across map and non-map dashboards

## Hard constraints

- Non-map dashboards:
  - no remote `<script src>`
  - no remote `<link href>`
  - must contain inline `<svg>`
- Map dashboards:
  - allowed only when the target artifact is `html_map`
  - may use Leaflet and remote basemap assets, matching current repo validation rules
- Dynamic mode:
  - store JWE in browser `localStorage`
  - load any previously stored JWE token on page init so subsequent dashboards can reuse it
  - use `https://mcp.demo.altinity.cloud/{JWE_TOKEN}/openapi/execute_query?query=...`
  - never embed real tokens in examples or committed artifacts
  - prefill the SQL textarea with the saved SQL and treat it as the dashboard's primary analytical data source
  - do not embed analytical CSV data or result JSON for the fetched dataset
  - surface every browser-executed query in a visible query ledger with its purpose and status
  - handle empty result payloads where `count = 0` and `rows = null` without treating them as malformed
  - place the JWE token input, SQL textarea, fetch button, forget button, and status text together in a separate control block at the very end of the page
  - keep that control block at the bottom after data loads; it must not jump into the hero or empty-state area

## Minimal output structure

Every dashboard should include:

- Branded header with title and analytical subtitle
- KPI strip derived from current data
- Primary visual
- Secondary visual or supporting table
- Filters when the data supports slicing
- Provenance block showing mode, data source, and query ledger
- Export control for current filtered rows

## Reference order

1. [references/theme.md](references/theme.md)
2. One of:
   - [references/static-mode.md](references/static-mode.md)
   - [references/dynamic-mode.md](references/dynamic-mode.md)
3. Optional:
   - [references/maps.md](references/maps.md)
4. [references/edit-mode.md](references/edit-mode.md)
5. [references/patterns.md](references/patterns.md)
