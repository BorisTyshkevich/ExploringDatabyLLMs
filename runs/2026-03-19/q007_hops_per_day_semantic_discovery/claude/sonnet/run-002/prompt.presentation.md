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

Airport codes are not guaranteed to be unique in `ontime.airports_latest`.
If a code-based lookup is unavoidable, first reduce `ontime.airports_latest` to one deterministic row per `code` in a subquery before joining.

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
- Result columns: `Tail_Number, Flight_Number_Reporting_Airline, Carrier, FlightDate, Hops, Route, DepartureSchedule`

Saved analysis artifact:

```json
{
  "sql": "WITH itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS Hops,\n        groupArray((toUInt32(coalesce(DepTime, CRSDepTime, 0)), Origin, Dest)) AS legs_raw\n    FROM ontime.ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate\n),\nitineraries_sorted AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        Carrier,\n        FlightDate,\n        Hops,\n        arraySort(x -\u003e x.1, legs_raw) AS legs_sorted\n    FROM itineraries\n)\nSELECT\n    Tail_Number,\n    Flight_Number_Reporting_Airline,\n    Carrier,\n    FlightDate,\n    Hops,\n    arrayStringConcat(\n        arrayConcat(\n            arrayMap(x -\u003e x.2, legs_sorted),\n            [arrayElement(arrayMap(x -\u003e x.3, legs_sorted), -1)]\n        ),\n        ' -\u003e '\n    ) AS Route,\n    arrayStringConcat(\n        arrayMap(\n            x -\u003e concat(x.2, ' ', printf('%02d:%02d', intDiv(x.1, 100), x.1 % 100)),\n            legs_sorted\n        ),\n        ', '\n    ) AS DepartureSchedule\nFROM itineraries_sorted\nORDER BY Hops DESC, FlightDate DESC\nLIMIT 10",
  "report_markdown": "# Highest Daily Hops for One Aircraft on One Flight Number\n\n## Overview\n\n{{data_overview_md}}\n\nThe maximum number of hops flown by a single aircraft under the same flight number in a single day is **{{metric.max_hops}}**, achieved exclusively by carrier **{{metric.max_hops_carrier}}**. All 10 entries in the top-10 list share this hop count, indicating that {{metric.max_hops}}-hop itineraries are a structured, repeating feature of Southwest Airlines scheduling rather than a one-off occurrence.\n\n## Most Recent Maximum-Hop Itinerary\n\n| Field | Value |\n|---|---|\n| Carrier | {{metric.most_recent_carrier}} |\n| Flight Number | {{metric.most_recent_flight_number}} |\n| Aircraft (Tail) | {{metric.most_recent_tail}} |\n| Date | {{metric.most_recent_date}} |\n| Hops | {{metric.max_hops}} |\n| Route | {{metric.most_recent_route}} |\n| Departure Schedule | {{metric.most_recent_departure_schedule}} |\n\nThe route spans {{metric.most_recent_airport_count}} airports across the country, with departures beginning before 06:00 local time and the last segment departing after 20:00, covering roughly 16 hours of continuous flying operations.\n\n## Top 10 Longest and Most Recent Itineraries\n\n{{result_table_md}}\n\n## Route Patterns and Clustering\n\nThe top-10 itineraries reveal strong route clustering by flight number:\n\n- **Flight {{metric.flight_3149_label}}** operates the identical {{metric.max_hops}}-stop sequence {{metric.flight_3149_occurrences}} times across different dates and tail numbers, confirming it as a regularly scheduled transcontinental multi-hop service.\n- **Flight {{metric.flight_154_label}}** and **Flight {{metric.flight_2787_label}}** each repeat their respective fixed routes on multiple dates, further demonstrating that {{metric.max_hops}}-hop itineraries are codified schedule patterns.\n- All top-10 entries are operated by {{metric.max_hops_carrier}}, which historically employs a point-to-point multi-stop model suited to high daily utilization of individual aircraft.",
  "metrics": {
    "summary_facts": [
      "The maximum hops flown by a single aircraft on one flight number in one day is 8, achieved exclusively by Southwest Airlines (WN). Three distinct flight numbers (366, 3149, 154, 2787) each repeat the same 8-stop route on multiple dates, confirming these are structured schedule patterns rather than anomalies."
    ],
    "named_values": {
      "flight_154_label": "154 (ELP -\u003e DAL -\u003e LIT -\u003e ATL -\u003e RIC -\u003e MDW -\u003e MCI -\u003e PHX -\u003e SAN)",
      "flight_2787_label": "2787 (MSY -\u003e ATL -\u003e CMH -\u003e BWI -\u003e RDU -\u003e BNA -\u003e DTW -\u003e MDW -\u003e LAX)",
      "flight_3149_label": "3149 (CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN)",
      "flight_3149_occurrences": "3",
      "max_hops": "8",
      "max_hops_carrier": "WN (Southwest Airlines)",
      "most_recent_airport_count": "9",
      "most_recent_carrier": "WN",
      "most_recent_date": "2024-12-01",
      "most_recent_departure_schedule": "ISP 05:43, BWI 08:10, MYR 10:20, BNA 11:42, VPS 14:01, DAL 16:43, LAS 18:28, OAK 20:41",
      "most_recent_flight_number": "366",
      "most_recent_route": "ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA",
      "most_recent_tail": "N957WN"
    },
    "named_lists": {
      "top_flight_numbers_8_hops": [
        "WN 366 — ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA (2024-12-01)",
        "WN 3149 — CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN (3 occurrences, latest 2024-02-18)",
        "WN 154 — ELP -\u003e DAL -\u003e LIT -\u003e ATL -\u003e RIC -\u003e MDW -\u003e MCI -\u003e PHX -\u003e SAN (2 occurrences, latest 2023-04-30)",
        "WN 2787 — MSY -\u003e ATL -\u003e CMH -\u003e BWI -\u003e RDU -\u003e BNA -\u003e DTW -\u003e MDW -\u003e LAX (3 occurrences, latest 2022-10-23)"
      ]
    }
  }
}
```

Saved report template to respect:

```report
# Highest Daily Hops for One Aircraft on One Flight Number

## Overview

{{data_overview_md}}

The maximum number of hops flown by a single aircraft under the same flight number in a single day is **{{metric.max_hops}}**, achieved exclusively by carrier **{{metric.max_hops_carrier}}**. All 10 entries in the top-10 list share this hop count, indicating that {{metric.max_hops}}-hop itineraries are a structured, repeating feature of Southwest Airlines scheduling rather than a one-off occurrence.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | {{metric.most_recent_carrier}} |
| Flight Number | {{metric.most_recent_flight_number}} |
| Aircraft (Tail) | {{metric.most_recent_tail}} |
| Date | {{metric.most_recent_date}} |
| Hops | {{metric.max_hops}} |
| Route | {{metric.most_recent_route}} |
| Departure Schedule | {{metric.most_recent_departure_schedule}} |

The route spans {{metric.most_recent_airport_count}} airports across the country, with departures beginning before 06:00 local time and the last segment departing after 20:00, covering roughly 16 hours of continuous flying operations.

## Top 10 Longest and Most Recent Itineraries

{{result_table_md}}

## Route Patterns and Clustering

The top-10 itineraries reveal strong route clustering by flight number:

- **Flight {{metric.flight_3149_label}}** operates the identical {{metric.max_hops}}-stop sequence {{metric.flight_3149_occurrences}} times across different dates and tail numbers, confirming it as a regularly scheduled transcontinental multi-hop service.
- **Flight {{metric.flight_154_label}}** and **Flight {{metric.flight_2787_label}}** each repeat their respective fixed routes on multiple dates, further demonstrating that {{metric.max_hops}}-hop itineraries are codified schedule patterns.
- All top-10 entries are operated by {{metric.max_hops_carrier}}, which historically employs a point-to-point multi-stop model suited to high daily utilization of individual aircraft.
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
WITH itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline AS Carrier,
        FlightDate,
        count() AS Hops,
        groupArray((toUInt32(coalesce(DepTime, CRSDepTime, 0)), Origin, Dest)) AS legs_raw
    FROM ontime.ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
),
itineraries_sorted AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        Carrier,
        FlightDate,
        Hops,
        arraySort(x -> x.1, legs_raw) AS legs_sorted
    FROM itineraries
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    Carrier,
    FlightDate,
    Hops,
    arrayStringConcat(
        arrayConcat(
            arrayMap(x -> x.2, legs_sorted),
            [arrayElement(arrayMap(x -> x.3, legs_sorted), -1)]
        ),
        ' -> '
    ) AS Route,
    arrayStringConcat(
        arrayMap(
            x -> concat(x.2, ' ', printf('%02d:%02d', intDiv(x.1, 100), x.1 % 100)),
            legs_sorted
        ),
        ', '
    ) AS DepartureSchedule
FROM itineraries_sorted
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