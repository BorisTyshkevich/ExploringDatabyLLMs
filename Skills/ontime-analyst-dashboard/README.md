# OnTime Analyst Dashboard

Human-oriented overview for the `ontime-analyst-dashboard` skill.

## What this skill is for

This skill standardizes how `visual.html` artifacts are built for analytics questions in this repository.

It is intended to give all dashboards:

- one coherent aviation-oriented corporate design
- one consistent data-handling model
- clear mode selection between self-contained benchmark artifacts and live browser dashboards
- Leaflet-based visualizations with optional map support

## Design goals

The main design direction is **aviation corporate**, not alert-triage.

That means:

- navy, slate, sky, teal, amber, and signal-red instead of the old purple-heavy alert palette
- clear KPI cards and readable tables
- strong but restrained visual hierarchy
- dashboards that feel analytical and operational, not incident-console oriented

## How it works

All dashboards:

- normalize and render data client-side
- use one shared visual design system
- choose a runtime mode first, then a map/non-map presentation

Mode split:

- static mode embeds the analytical data directly in `visual.html` so benchmark artifacts stay self-contained
- dynamic mode fetches analytical data through a runtime-provided tokenized HTTP SQL endpoint and shows a visible query ledger

Maps are used when geography adds analytical value (routes, airports, regions). For non-geographic questions, use charts, heatmaps, or tables as the primary visual.

## What every dashboard should include

- a branded header with a useful analytical subtitle
- KPI strip derived from the loaded data
- primary visual
- supporting visual or detail table
- filters when the data naturally supports slicing
- query ledger showing data provenance
- export for filtered rows

## Data model philosophy

All data flows through one normalized client-side row model:

- embedded static data or MCP results are transformed into row objects
- KPIs, filters, charts, and tables are derived from normalized data
- Temporal fields are normalized to `YYYY-MM-DD` before use

## Browser storage keys

- `OnTimeAnalystDashboard::<dashboardId>::layout` — layout state

## Files in this skill

- `SKILL.md` — agent-facing entrypoint
- `references/theme.md` — design tokens
- `references/dynamic-mode.md` — fetch and render logic
- `references/static-mode.md` — self-contained artifact rules
- `references/maps.md` — Leaflet and geography
- `references/patterns.md` — layout shapes
- `references/edit-mode.md` — optional layout editing
