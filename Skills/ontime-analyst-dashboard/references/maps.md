# Map Support

Use maps when geography adds analytical value.

## When to use a map

Use maps when geography is part of the reasoning:

- itineraries and hop chains
- route rankings with geographic context
- hub spillovers
- airport hotspots
- regional concentration
- network or corridor views

Do not use a map as the primary visual for purely temporal or purely ranked questions unless geography adds real signal.

## Allowed assets

- For `html_map`, Leaflet JS/CSS from CDN is allowed
- Remote basemap tiles are allowed for `html_map`
- Non-map dashboards should not pull in Leaflet or map tiles just because the library is available

## Required data blocks

Typical map dashboards need:

- airport coordinates:
  `airport_code`, `lat`, `lon`, optional `city`, `label`
- route or leg rows:
  `origin`, `dest`, optional `rank`, `distance`, `metric`, `highlight_flag`
- optional annotations:
  `note`, `group`, `severity`

## Standard map behavior

- For `html_map`, the map card and map container must be present in the initial HTML structure. Do not create or reveal the map section only after enrichment succeeds.
- Fit bounds to displayed points
- Provide legend
- Highlight the single most important route, airport, or event
- Use popups or tooltips for details
- Keep routes semi-transparent unless highlighted
- Make markers clickable or hoverable
- Do not place narrative cards, KPI summaries, or lead-itinerary copy as overlays on top of the map surface; keep that content in normal panels adjacent to or below the map
- Check dataset-native coordinate and lookup tables first when the primary query lacks lat/lon
- Coordinate enrichment queries are acceptable when they are explicit, user-visible, and tied to the map
- Cache coordinate enrichment results and reuse them for selection changes instead of issuing new per-click queries
- If the map card is cloned from a template, initialize Leaflet against a live scoped element reference, not a duplicated global `id` that may still exist inside the template
- Prefer initializing Leaflet only after the map container is visible in layout; do not create the map while its parent region is hidden if delayed init is practical
- If layout or visibility changes after map creation, call `invalidateSize()` after reveal and after major layout moves as a safety fallback
- Compute bounds from raw lat/lon arrays or concrete marker/polyline layers; do not compute map bounds from parent layer-group containers
- Row-selection redraws must update the selected route or its degraded state without throwing uncaught JS errors
- If a selected route cannot be plotted after the map already exists, replace the current route with a visible degraded message instead of leaving stale content or a blank map surface

## Styling rules

- Keep the same hero, KPI, and table system as non-map dashboards
- Use `--sky` for contextual routes
- Use `--red` for highlighted route/event
- Use `--navy` for airport markers
- Popups and labels should inherit the dashboard typography

## Fallback behavior

If coordinates are missing:

- render KPIs, table, and supporting chart anyway
- attempt an explicit enrichment query against dataset-native lookup tables from the allowed dataset scope before degrading
- keep the map card visible and degrade the map state with a visible warning or placeholder message if enrichment is unavailable or fails
- do not keep a slippy map only to place summary text on top of it
