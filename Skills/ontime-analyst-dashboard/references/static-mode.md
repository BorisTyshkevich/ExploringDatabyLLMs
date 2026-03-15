# Static Mode

Use static mode for benchmark artifacts and any self-contained `visual.html`.

## Required constraints

- Inline CSS only
- Inline JS only
- For non-map dashboards, inline SVG only
- No remote scripts or stylesheets
- Data must be embedded in the page

## Data embedding

Embed CSV in `<script type="text/csv">` tags.

Example:

```html
<script type="text/csv" id="summaryCsv">
label,value
Flights,1200
DelayPct,18.4
</script>
```

## Parsing pattern

Use a small built-in CSV parser in plain JS. Keep it simple and deterministic.

Minimum behavior:

- split header row
- support quoted values
- trim cells
- skip empty trailing rows
- return array of objects

After parsing:

- drop fully empty rows
- coerce obvious numerics
- preserve strings for labels and IDs
- normalize null-like values to `null` or `''` consistently

## Normalization rules

- `row[key] = row[key] ?? ''` for string fields before string operations
- numeric metrics should be coerced with `Number(...)` and guarded with `Number.isFinite`
- never assume a column exists
- if required fields are missing, show a visible warning panel in the dashboard

## Dashboard behavior

- derive KPIs from parsed data, never from hardcoded literals
- derive filter options from current data
- export currently filtered rows as CSV
- if multiple CSV blocks exist, name them by role:
  - `summaryCsv`
  - `trendCsv`
  - `detailCsv`
  - `lookupCsv`
  - `airportsCsv`
  - `routesCsv`

## Recommended visuals

- ranking: horizontal SVG bar chart
- time series: SVG line/area chart
- heatmap: SVG matrix with text labels and legend
- scatter: SVG circles with scales derived in JS
- tables: semantic HTML table, not SVG

## Empty and error states

- If no rows remain after parsing, keep header visible and render a clear “No data available” panel
- If required columns are missing, explain which fields are absent
- Do not let a broken chart suppress the table and KPI sections
