# Static Mode

Use static mode for self-contained benchmark artifacts and any `visual.html` that must render without browser-side MCP access.

## Required constraints

- Inline CSS and inline JavaScript
- Analytical data embedded directly in the page
- No live MCP fetches, no token flow, and no dependency on browser `localStorage`
- For non-map dashboards, no remote scripts or stylesheets
- For `html_map`, Leaflet and basemap assets are allowed, but the analytical dataset itself must still be embedded

## Data embedding

Embed analytical rows in `<script type="application/json">` or `<script type="text/csv">` blocks and parse them in browser code.

## Normalization rules

- Normalize parsed rows before deriving KPIs, filters, charts, or tables
- Guard optional fields with `?.` / `??`
- If a date may be serialized as an ISO datetime, derive a stable `YYYY-MM-DD` key once and reuse it
- Show a clear warning panel when required fields are missing or no rows remain after parsing

## Dashboard behavior

- Derive visuals from the embedded data only
- Keep export behavior local to the current embedded/filtered row set
- If maps are present, render them from embedded coordinates or route geometry rather than browser-side enrichment fetches
