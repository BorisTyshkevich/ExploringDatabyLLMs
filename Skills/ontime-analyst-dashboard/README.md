# OnTime Analyst Dashboard

Human-oriented overview for the `ontime-analyst-dashboard` skill.

## What this skill is for

This skill standardizes how `visual.html` artifacts are built for the OnTime reporting questions in this repository.

It is intended to give all dashboards:

- one coherent aviation-oriented corporate design
- one consistent data-handling model
- one clear split between static benchmark artifacts and live browser dashboards
- optional map support for the report types where geography is part of the analysis

The skill is designed for the OnTime question set already in this repo, including:

- carrier leadership and market-share shifts
- airport and carrier delay rankings
- hotspot heatmaps
- seasonality analysis
- peak-month contribution analysis
- chronic schedule stress reports
- hub spillover and network effects
- itinerary and route maps

## Design goals

The main design direction is **aviation corporate**, not alert-triage.

That means:

- navy, slate, sky, teal, amber, and signal-red instead of the old purple-heavy alert palette
- clear KPI cards and readable tables
- strong but restrained visual hierarchy
- dashboards that feel analytical and operational, not incident-console oriented

The same theme should be used across:

- static dashboards
- dynamic dashboards
- map dashboards
- non-map dashboards

## Two modes

### 1. Static mode

Use static mode for benchmark artifacts and any self-contained `visual.html`.

Key properties:

- inline CSS
- inline JavaScript
- embedded CSV data
- no external JS/CSS for non-map dashboards
- validator-safe for the current repo rules

This is the default mode for questions that generate artifact outputs under `runs/...`.

### 2. Dynamic mode

Use dynamic mode for live dashboards opened in a browser and connected to the MCP OpenAPI endpoint.

Key properties:

- browser fetches data directly from MCP
- JWE token is stored in browser `localStorage`
- same client-side normalization pipeline as static mode
- suitable for interactive investigation after the benchmark/reporting flow
- do not embed analytical CSV or result JSON in the page

## Map policy

Maps are part of this skill. They are **not** split into a separate map skill.

Why:

- many OnTime reports mix charts, tables, and geography
- route, itinerary, spillover, and hub analyses benefit from a shared visual system
- splitting maps out would duplicate theme, auth, loader, and layout guidance

Important constraint:

- for `visual_type: html_map`, Leaflet and remote basemap assets are acceptable under the current repo validator
- for non-map visual types, dashboards should stay self-contained and use inline SVG only

So the skill supports both:

- non-map validator-safe dashboards
- map dashboards when the question explicitly uses the `html_map` path

## What every dashboard should include

Regardless of mode, the intended structure is:

- a branded header with a useful analytical subtitle
- KPI strip derived from the loaded data
- primary visual
- supporting visual or detail table
- filters when the data naturally supports slicing
- provenance or source context
- export for filtered rows

The skill also expects a lightweight edit mode so card layout can be rearranged without changing source code.

## Data model philosophy

The core idea is that **all data flows through one normalized client-side row model**.

That means:

- static CSV data is parsed and normalized
- dynamic MCP results are transformed into row objects and normalized the same way
- KPIs, filters, charts, and tables are always derived from normalized data

This reduces drift between benchmark visuals and live investigative dashboards.

## Browser storage keys

The skill reserves these key patterns:

- auth token:
  `OnTimeAnalystDashboard::auth::jwe`
- layout:
  `OnTimeAnalystDashboard::<dashboardId>::layout`
- optional saved query:
  `OnTimeAnalystDashboard::<dashboardId>::query`

This keeps layouts and queries isolated per dashboard while sharing a single auth token across dashboards.

## Good fit for current report types

### Best fit for standard dashboards

- yearly carrier leadership
- worst origin airport OTP
- winter carrier-origin rankings
- peak delay month contribution analysis
- chronic schedule stress

### Best fit for heatmap-style dashboards

- Delta ATL destination/time-block hotspots
- route-season delay analysis

### Best fit for map-enabled dashboards

- q001 itinerary/hops map
- hub disruption spillovers
- route-heavy analyses where geography materially improves understanding

## Why this skill exists

The older `alert-analyst-dashboard` skill proved the value of:

- a reusable dashboard structure
- a consistent visual language
- a loader mode
- editable layout behavior

But it is still centered on alerts and severity concepts, and it assumes library patterns that do not cleanly fit the current non-map validator rules.

This new skill keeps the useful ideas, but adapts them for:

- OnTime aviation analytics
- current `visual.html` validation constraints
- live MCP OpenAPI access
- mixed map and non-map reporting

## Files in this skill

- `SKILL.md`
  Agent-facing entrypoint and usage rules
- `references/theme.md`
  design tokens and visual system
- `references/static-mode.md`
  self-contained dashboard guidance
- `references/dynamic-mode.md`
  live MCP loading guidance
- `references/maps.md`
  map support and constraints
- `references/edit-mode.md`
  native edit-mode layout rules
- `references/patterns.md`
  recommended dashboard shapes for the current OnTime reports

## Practical summary

Use this skill whenever an OnTime question needs `visual.html`.

Think of it as:

- one dashboard system
- two data modes
- optional map support
- one unified aviation design language

If the artifact is for benchmark validation, prefer static mode.  
If the artifact is for exploratory browser use, dynamic mode is appropriate.  
If geography is a core part of the reasoning and the question is `html_map`, bring in the map guidance.
