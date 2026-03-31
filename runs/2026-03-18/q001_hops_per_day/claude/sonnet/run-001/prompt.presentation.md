Generate only the artifacts requested in this prompt.
Use the configured MCP server for all data access.
Do not construct raw OpenAPI URLs manually.
Stay within the configured dataset scope.

Dataset semantic layer:

Use `ontime.ontime` as the primary fact table for flight operations.

Use `ontime.airports_latest` for current airport reference data such as:

- `code`
- `name`
- `latitude`
- `longitude`
- `utc_local_time_variation`

Preferred joins:

- `ontime.ontime.OriginAirportID = ontime.airports_latest.airport_id`
- `ontime.ontime.DestAirportID = ontime.airports_latest.airport_id`

Fallback joins:

- use `ontime.ontime.Origin = ontime.airports_latest.code`
- use `ontime.ontime.Dest = ontime.airports_latest.code`

Use the `airport_id` joins when those columns are available.
Use code-based joins only when the analytical result exposes route strings or airport codes but not airport IDs.

For dynamic map visuals, explicit airport-coordinate enrichment queries against `ontime.airports_latest` are allowed when the primary analytical result does not already include usable coordinates.

Generate only the visual artifact.

The analytical run already produced:

- `query.sql`
- `report.template.md`
- `report.md`
- `result.json`

Use `query.sql`, `report.template.md`, and `result.json` as authoritative inputs.
Do not regenerate SQL or report artifacts.
Do not respond with a prose summary of what you created.

Return exactly this fenced section:

```html
<!doctype html>
<html>...</html>
```

Visual input context:

- Question title: `Highest daily hops for one aircraft on one flight number`
- Result columns: `Tail_Number, Carrier, FlightNum, FlightDate, Hops, Route, DepartureTimes`

Saved report template to respect:

```report
# {{question_title}}

_Generated: {{generated_at}}_

## Overview

{{data_overview_md}}

## Top 10 Longest Daily Itineraries

{{result_table_md}}

## Analysis

### Maximum Hop Count

The maximum observed is **{{max_hops}} hops** in a single calendar day under a single flight number, producing an itinerary of {{max_hops_plus_one}} airports. This is not a one-off anomaly: the top 10 table contains multiple distinct dates with the same hop count, confirming it reflects a repeating scheduled operation rather than an exceptional day.

### Most Recent Maximum-Hop Itinerary

The single most recent entry in the top 10 is:

- **Carrier / Flight:** {{most_recent_carrier}} {{most_recent_flight_num}}
- **Date:** {{most_recent_date}}
- **Tail Number:** {{most_recent_tail}}
- **Route:** {{most_recent_route}}
- **Leg departure times:** {{most_recent_dep_times}}

### Route Repetition and Clustering

Across the top 10 itineraries, distinct named flight numbers recur on multiple dates with identical or nearly identical routing:

{{result_table_md}}

All carriers represented are Southwest Airlines (WN), consistent with WN's point-to-point network model in which a single flight number can chain many short segments across one day. Route patterns appear to repeat weekly or on a regular schedule, suggesting these are published schedule rotations rather than ad-hoc operations.
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
WITH legs AS (
    SELECT
        Tail_Number,
        IATA_CODE_Reporting_Airline AS Carrier,
        Flight_Number_Reporting_Airline AS FlightNum,
        FlightDate,
        Origin,
        Dest,
        DepTime
    FROM ontime.ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND DepTime IS NOT NULL
),
grouped AS (
    SELECT
        Tail_Number,
        Carrier,
        FlightNum,
        FlightDate,
        count() AS Hops,
        arraySort(x -> tupleElement(x, 1), groupArray((DepTime, Origin, Dest))) AS sorted_legs
    FROM legs
    GROUP BY Tail_Number, Carrier, FlightNum, FlightDate
)
SELECT
    Tail_Number,
    Carrier,
    FlightNum,
    FlightDate,
    Hops,
    arrayStringConcat(
        arrayConcat(
            arrayMap(x -> tupleElement(x, 2), sorted_legs),
            [tupleElement(sorted_legs[length(sorted_legs)], 3)]
        ),
        ' -> '
    ) AS Route,
    arrayStringConcat(
        arrayMap(x -> concat(
            leftPad(toString(intDiv(tupleElement(x, 1), 100)), 2, '0'),
            ':',
            leftPad(toString(tupleElement(x, 1) % 100), 2, '0')
        ), sorted_legs),
        ', '
    ) AS DepartureTimes
FROM grouped
ORDER BY Hops DESC, FlightDate DESC
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
- run an explicit airport-coordinate enrichment query against `ontime.airports_latest` using airport codes parsed from the route strings
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