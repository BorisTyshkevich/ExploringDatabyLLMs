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
- Result columns: `Date, AircraftID, FlightNumber, Carrier, Hops, Route`

Saved report template to respect:

```report
# Highest daily hops for one aircraft on one flight number

{{data_overview_md}}

The maximum number of daily hops observed for a single aircraft under the same flight number is {{metric.max_hops}}. This behavior does not appear to be a one-off itinerary, but rather a repeated operating pattern, particularly noticeable with Southwest Airlines (WN).

The most recent itinerary achieving this maximum is flight {{metric.recent_flight}} on {{metric.recent_date}} operated by {{metric.recent_carrier}} (Aircraft {{metric.recent_aircraft}}). The route was:
**{{metric.recent_route}}**

Across the top 10 longest itineraries, there is notable route repetition. For instance, Southwest's flight 3149 flew the exact same 8-hop route (`CLE - BNA - PNS - HOU - MCI - PHX - BUR - OAK - DEN`) on multiple dates in early 2024. Similar identical 8-hop patterns are seen for flight 154 in April 2023 and flight 2787 in Fall 2022. This indicates intentional, cyclic aircraft scheduling across numerous stations within a single day.
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
SELECT FlightDate                                                AS Date, Tail_Number AS AircraftID,
       Flight_Number_Reporting_Airline                           AS FlightNumber, Reporting_Airline AS Carrier,
       count()                                                   AS Hops,
       arrayStringConcat(arrayMap(x -> x.2 || ' (' || leftPad(toString(x.1), 4, '0') || ')',
                                  arraySort(groupArray(tuple(DepTime, Origin, Dest)))), ' - ') || ' - ' ||
       arraySort(groupArray(tuple(DepTime, Origin, Dest)))[-1].3 AS Route
FROM ontime.ontime
WHERE Cancelled = 0 AND Tail_Number != '' AND Tail_Number != 'UNKNOW' AND DepTime IS NOT NULL
GROUP BY Date, AircraftID, FlightNumber, Carrier
ORDER BY Hops DESC, Date DESC
LIMIT 10
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
- Execute the embedded saved SQL in the browser as the primary query.
- Keep the embedded saved SQL authoritative for the artifact.
- Surface every browser query in a visible query ledger.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Keep additional browser queries limited to explicit enrichment or drill-down that materially improves the visualization.