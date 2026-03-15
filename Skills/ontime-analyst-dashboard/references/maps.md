# Map Support

Maps belong inside this skill. Do not create a separate map skill unless the repo later adopts a separate geospatial stack.

## When to use a map

Use maps when geography is part of the reasoning:

- itineraries and hop chains
- route rankings with geographic context
- hub spillovers
- airport hotspots
- regional concentration
- network or corridor views

Do not use a map as the primary visual for purely temporal or purely ranked questions unless geography adds real signal.

## Validator split

- If `visual_type` is `html_map`:
  - Leaflet is allowed
  - remote basemap tiles are allowed
  - remote Leaflet CSS/JS is allowed under current repo validation
- If `visual_type` is not `html_map`:
  - keep maps out, or use inline SVG geography only
  - no remote assets

## Required data blocks

Typical map dashboards need:

- airport coordinates:
  `airport_code`, `lat`, `lon`, optional `city`, `label`
- route or leg rows:
  `origin`, `dest`, optional `rank`, `distance`, `metric`, `highlight_flag`
- optional annotations:
  `note`, `group`, `severity`

## Standard map behavior

- Fit bounds to displayed points
- Provide legend
- Highlight the single most important route, airport, or event
- Use popups or tooltips for details
- Keep routes semi-transparent unless highlighted
- Make markers clickable or hoverable

## Styling rules

- Keep the same hero, KPI, and table system as non-map dashboards
- Use `--sky` for contextual routes
- Use `--red` for highlighted route/event
- Use `--navy` for airport markers
- Popups and labels should inherit the dashboard typography

## Fallback behavior

If coordinates are missing:

- render KPIs, table, and supporting chart anyway
- replace the map with a visible warning card explaining which coordinates are missing

## Good fit for current OnTime reports

- `q001_hops_per_day`
- `q007_highest_delay_route_season_top50` as an optional secondary view
- `q009_hub_disruption_spillover_days`
