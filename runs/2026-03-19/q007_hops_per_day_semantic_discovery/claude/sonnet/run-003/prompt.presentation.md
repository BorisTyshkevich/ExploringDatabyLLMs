Generate only the artifacts requested in this prompt.
Use the configured MCP server for all data access.
Do not construct raw OpenAPI URLs manually.
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

- `analysis.json`
- `query.sql`
- `report.template.md`
- `report.md`
- `result.json`

Use `analysis.json`, `query.sql`, `report.template.md`, and `result.json` as authoritative inputs.
Do not regenerate SQL or report artifacts.
Do not respond with a prose summary of what you created.

Return exactly this fenced section:

```html
<!doctype html>
<html>...</html>
```

Visual input context:

- Question title: `Highest daily hops for one aircraft on one flight number`
- Result columns: `Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate, Hops, Route`

Saved analysis artifact:

```json
{
  "sql": "WITH legs AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline,\n        FlightDate,\n        Origin,\n        Dest,\n        assumeNotNull(DepTime) AS DepTime\n    FROM ontime.ontime\n    WHERE DepTime IS NOT NULL\n      AND Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n),\ndaily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline,\n        FlightDate,\n        count() AS Hops,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(\n                    x -\u003e concat(x.2, '(', toString(x.1), ')'),\n                    arraySort(x -\u003e x.1, groupArray((DepTime, Origin)))\n                ),\n                [argMax(Dest, DepTime)]\n            ),\n            ' -\u003e '\n        ) AS Route\n    FROM legs\n    GROUP BY\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline,\n        FlightDate\n)\nSELECT\n    Tail_Number,\n    Flight_Number_Reporting_Airline,\n    IATA_CODE_Reporting_Airline,\n    FlightDate,\n    Hops,\n    Route\nFROM daily_itineraries\nORDER BY Hops DESC, FlightDate DESC\nLIMIT 10",
  "report_markdown": "# Highest Daily Hops for One Aircraft on One Flight Number\n\n{{data_overview_md}}\n\n## Maximum Hop Count\n\nThe highest number of hops recorded for a single aircraft operating under a single flight number on one calendar day is **{{metric.max_hops}}**. This means the aircraft completed {{metric.max_hops}} individual flight legs under the same flight number within a single day.\n\n## Most Recent Maximum-Hop Itinerary\n\n| Field | Value |\n|---|---|\n| Carrier | {{metric.top_carrier}} |\n| Flight Number | {{metric.top_flight_number}} |\n| Aircraft (Tail) | {{metric.top_tail_number}} |\n| Date | {{metric.top_date}} |\n| Route | {{metric.top_route}} |\n\nThe route above shows each origin airport followed by its actual departure time (HHMM) in parentheses, ending at the final destination.\n\n## Operating Pattern Analysis\n\n{{metric.pattern_summary}}\n\n## Top 10 Longest Itineraries\n\n{{result_table_md}}\n\n### Notable Observations\n\n- All top-10 itineraries belong to **{{metric.dominant_carrier}}**, reflecting that carrier's hub-and-spoke multi-stop routing strategy.\n- The most repeated flight number in the top 10 is **{{metric.most_repeated_flight}}**, suggesting a recurrent scheduled itinerary rather than an ad-hoc assignment.\n- Route clustering across dates for the same flight number indicates a **{{metric.recurrence_pattern}}** operating schedule.\n\n_Generated: {{generated_at}}_",
  "metrics": {
    "summary_facts": [
      "The maximum daily hop count for a single aircraft on one flight number is 8, achieved exclusively by Southwest Airlines (WN) across multiple flight numbers and dates.",
      "Flight WN3149 accounts for 4 of the top 10 itineraries, all in January–February 2024, confirming a recurring weekly 8-hop schedule.",
      "The most recent maximum-hop itinerary is N957WN operating WN366 on 2024-12-01: ISP(543) -\u003e BWI(810) -\u003e MYR(1020) -\u003e BNA(1142) -\u003e VPS(1401) -\u003e DAL(1643) -\u003e LAS(1828) -\u003e OAK(2041) -\u003e SEA."
    ],
    "named_values": {
      "dominant_carrier": "WN (Southwest Airlines)",
      "max_hops": "8",
      "most_repeated_flight": "WN3149 (4 occurrences in top 10)",
      "pattern_summary": "The 8-hop maximum appears to be a repeated scheduled pattern, not a one-off. Flight WN3149 ran the same 8-leg CLE→BNA→PNS→HOU→MCI→PHX→BUR→OAK→DEN itinerary on four consecutive Sundays in January–February 2024, operated by different aircraft each week. Flight WN154 similarly repeated an 8-hop ELP-based itinerary across two Sundays in April 2023. This strongly indicates these are scheduled multi-stop routes, not exceptional occurrences.",
      "recurrence_pattern": "weekly (approximately every 7 days for WN3149 in Jan–Feb 2024)",
      "top_carrier": "WN (Southwest Airlines)",
      "top_date": "2024-12-01",
      "top_flight_number": "366",
      "top_route": "ISP(543) -\u003e BWI(810) -\u003e MYR(1020) -\u003e BNA(1142) -\u003e VPS(1401) -\u003e DAL(1643) -\u003e LAS(1828) -\u003e OAK(2041) -\u003e SEA",
      "top_tail_number": "N957WN"
    },
    "named_lists": {
      "top_flight_numbers": [
        "WN366 (2024-12-01, N957WN)",
        "WN3149 (2024-02-18, N7835A)",
        "WN3149 (2024-02-04, N263WN)",
        "WN3149 (2024-01-28, N429WN)",
        "WN3149 (2024-01-21, N228WN)",
        "WN3149 (2024-01-14, N569WN)",
        "WN154 (2023-04-30, N7742B)",
        "WN154 (2023-04-16, N929WN)",
        "WN2787 (2022-10-23, N8631A)",
        "WN2787 (2022-10-02, N8809L)"
      ]
    }
  }
}
```

Saved report template to respect:

```report
# Highest Daily Hops for One Aircraft on One Flight Number

{{data_overview_md}}

## Maximum Hop Count

The highest number of hops recorded for a single aircraft operating under a single flight number on one calendar day is **{{metric.max_hops}}**. This means the aircraft completed {{metric.max_hops}} individual flight legs under the same flight number within a single day.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | {{metric.top_carrier}} |
| Flight Number | {{metric.top_flight_number}} |
| Aircraft (Tail) | {{metric.top_tail_number}} |
| Date | {{metric.top_date}} |
| Route | {{metric.top_route}} |

The route above shows each origin airport followed by its actual departure time (HHMM) in parentheses, ending at the final destination.

## Operating Pattern Analysis

{{metric.pattern_summary}}

## Top 10 Longest Itineraries

{{result_table_md}}

### Notable Observations

- All top-10 itineraries belong to **{{metric.dominant_carrier}}**, reflecting that carrier's hub-and-spoke multi-stop routing strategy.
- The most repeated flight number in the top 10 is **{{metric.most_repeated_flight}}**, suggesting a recurrent scheduled itinerary rather than an ad-hoc assignment.
- Route clustering across dates for the same flight number indicates a **{{metric.recurrence_pattern}}** operating schedule.

_Generated: {{generated_at}}_
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
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Origin,
        Dest,
        assumeNotNull(DepTime) AS DepTime
    FROM ontime.ontime
    WHERE DepTime IS NOT NULL
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
),
daily_itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arrayStringConcat(
            arrayConcat(
                arrayMap(
                    x -> concat(x.2, '(', toString(x.1), ')'),
                    arraySort(x -> x.1, groupArray((DepTime, Origin)))
                ),
                [argMax(Dest, DepTime)]
            ),
            ' -> '
        ) AS Route
    FROM legs
    GROUP BY
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    Hops,
    Route
FROM daily_itineraries
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
```

Visual context:

- Question title: `Highest daily hops for one aircraft on one flight number`
- Visual mode: `static`
- Visual type: `html_map`

Question-specific visual guidance:

The page must:

- show a lead-itinerary map using embedded airport coordinates from the static artifact
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- reuse embedded airport metadata for any itinerary selected from the primary result set without issuing browser-side data fetches
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip anchored to the top-ranked result even when the selected itinerary changes
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- if the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary and continue rendering the non-map analysis

Static-mode requirements:

- Build a self-contained benchmark artifact. Do not require browser-side MCP access, tokens, localStorage, or live fetches to render the analytical content.
- Embed the analytical data needed by the page directly in the HTML using inline data blocks such as `<script type="application/json">` or `<script type="text/csv">`.
- Derive KPIs, filters, charts, and tables from the embedded data after parsing and normalization in browser JavaScript.
- Keep CSS and JavaScript inline.
- For non-map visuals, do not use remote `<script>` or `<link>` assets.
- For `html_map`, Leaflet and remote basemap assets are allowed, but the analytical dataset itself must still be embedded rather than fetched.
- Show clear empty-state or missing-field warnings instead of rendering broken visuals.
- Normalize temporal fields before grouping, filtering, or comparison logic. If a ClickHouse `Date` may appear as ISO datetime text, derive a stable `YYYY-MM-DD` key and reuse it consistently.
- If the page uses `<template>` cloning, bind behavior to scoped selectors or per-instance references rather than duplicated global `id` values.