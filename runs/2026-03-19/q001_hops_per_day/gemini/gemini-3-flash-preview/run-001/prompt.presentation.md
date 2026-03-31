Connect to clickhouse server though MCP connection 
Do not use direct HTTP by any tools like curl.
Generate only the artifacts requested in this prompt.
Stay within the configured dataset scope.

Dataset semantic layer:

Use `ontime.ontime` as the primary fact table for flight operations.

Use `ontime.airports_latest` as the semantic airport dimension for current airport reference data, including:

- `code`
- `name`
- `latitude`
- `longitude`
- `utc_local_time_variation`

`ontime.airports_latest` can be used directly to get coordinates and other columns for enrichment in application code (such as JavaScript or python)
or used for SQL JOINs.

Preferred sql joins:

- `ontime.ontime.OriginAirportID = ontime.airports_latest.airport_id`
- `ontime.ontime.DestAirportID = ontime.airports_latest.airport_id`

For `ontime.ontime`, use `OriginAirportID` and `DestAirportID` when enriching airport names, coordinates, or other airport attributes inside SQL.
Do not join `ontime.ontime` legs directly to `ontime.airports_latest` by airport code when the airport id columns are already available.

Fallback sql joins:

- use `ontime.ontime.Origin = ontime.airports_latest.code`
- use `ontime.ontime.Dest = ontime.airports_latest.code`

Use the `airport_id` joins when those columns are available.
Use code-based joins only when the analytical result exposes route strings or airport codes but not airport IDs.

Generate only the visual artifact.

The analytical run already produced:

- `query.sql`
- `report.template.md`
- `report.md`
- `result.json`

Use `query.sql`, `report.template.md`, `report.md`, and `result.json` as authoritative inputs.
Do not regenerate SQL or report artifacts.
Do not respond with a prose summary of what you created.

Return exactly this fenced section:

```html
<!doctype html>
<html>...</html>
```

Visual input context:

- Question title: `Highest daily hops for one aircraft on one flight number`
- Result columns: `AircraftID, FlightNumber, Carrier, Date, Hops, Route, DepTimes`

Saved report template to respect:

```report
# Highest daily hops for one aircraft on one flight number

{{data_overview_md}}

The highest number of hops recorded for a single aircraft using a consistent flight number in one day is **{{metric.primary_value}}**. This pattern is observed across multiple dates, notably with Southwest Airlines (WN) flight 3149, which frequently operates an 8-hop route such as `CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN`.

## Longest Recent Itinerary
The most recent instance of an 8-hop itinerary was on **{{metric.most_recent_date}}** by carrier **{{metric.most_recent_carrier}}** (Flight **{{metric.most_recent_flight}}**), using aircraft **{{metric.most_recent_tail}}**. The full route was:
`{{metric.most_recent_route}}`

Departure times from each origin were: `{{metric.most_recent_dep_times}}`.

## Operational Analysis
The top 10 results show that 8 hops is the operational limit for a single flight number/aircraft combination in this dataset. Most of these occurrences belong to Southwest Airlines, indicating a high-utilization point-to-point routing model. The repetition of specific routes (e.g., flight 3149 and flight 2787) suggests these are scheduled multi-stop service patterns rather than one-off occurrences.
```

Create `visual.html` using the `ontime-analyst-dashboard` skill.

The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.

General visual rules:

- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.
- Keep report logic out of `visual.html`; qforge already renders `report.md` separately.

Saved SQL to preserve in the final page contract:

```sql
SELECT Tail_Number AS AircraftID, Flight_Number_Reporting_Airline AS FlightNumber, Reporting_Airline AS Carrier, FlightDate AS Date, count() AS Hops, arrayStringConcat(arrayPushBack(arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((DepTime, Origin, Dest)))), (arraySort(x -> x.1, groupArray((DepTime, Origin, Dest)))[-1]).3), ' -> ') AS Route, arrayStringConcat(arrayMap(x -> lpad(toString(x.1), 4, '0'), arraySort(x -> x.1, groupArray((DepTime, Origin, Dest)))), ', ') AS DepTimes FROM ontime.ontime WHERE Cancelled = 0 AND Tail_Number != '' GROUP BY AircraftID, FlightNumber, Carrier, Date ORDER BY Hops DESC, Date DESC LIMIT 10
```

Visual context:

- Question title: `Highest daily hops for one aircraft on one flight number`
- Visual mode: `dynamic`
- Visual type: `html_map`

Question-specific visual guidance:

The page must:

- show a lead-itinerary map that remains present even before airport-coordinate enrichment succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- label the map as airport-coordinate enrichment in the query ledger
- reuse the enrichment results for any itinerary selected from the primary result set without issuing a new per-click enrichment query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip anchored to the top-ranked result even when the selected itinerary changes
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis

Dynamic-mode additions:

- Build the page in dynamic mode using the `ontime-analyst-dashboard` skill contract.
- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE_TOKEN}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Include the footer control block for JWE and SQL controls.
- Execute the embedded saved SQL in the browser as the primary query.
- Keep the embedded saved SQL authoritative for the artifact.
- Surface every browser query in a visible query ledger.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Keep additional browser queries limited to explicit enrichment or drill-down that materially improves the visualization.
- Never call ClickHouse directly through localhost or any hardcoded server URL.