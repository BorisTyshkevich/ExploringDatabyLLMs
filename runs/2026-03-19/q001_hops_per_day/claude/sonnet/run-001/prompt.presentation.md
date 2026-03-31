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
- Result columns: `Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate, hops, Route, DepartureTimes`

Saved analysis artifact:

```json
{
  "sql": "WITH leg_data AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline,\n        FlightDate,\n        count() AS hops,\n        groupArray((toUInt32(ifNull(DepTime, ifNull(CRSDepTime, 9999))), Origin, Dest)) AS legs_raw\n    FROM ontime.ontime\n    WHERE Cancelled = 0\n      AND Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline,\n        FlightDate\n),\nitineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline,\n        FlightDate,\n        hops,\n        arraySort(x -\u003e x.1, legs_raw) AS legs_sorted\n    FROM leg_data\n)\nSELECT\n    Tail_Number,\n    Flight_Number_Reporting_Airline,\n    IATA_CODE_Reporting_Airline,\n    FlightDate,\n    hops,\n    arrayStringConcat(\n        arrayConcat(\n            [legs_sorted[1].2],\n            arrayMap(x -\u003e x.3, legs_sorted)\n        ),\n        ' -\u003e '\n    ) AS Route,\n    arrayStringConcat(\n        arrayMap(\n            x -\u003e concat(x.2, ' ', lpad(toString(intDiv(x.1, 100)), 2, '0'), ':', lpad(toString(x.1 % 100), 2, '0')),\n            legs_sorted\n        ),\n        ' | '\n    ) AS DepartureTimes\nFROM itineraries\nORDER BY hops DESC, FlightDate DESC\nLIMIT 10",
  "report_markdown": "# Highest Daily Hops for One Aircraft on One Flight Number\n\n## Overview\n\nThis report identifies the longest single-day itineraries flown by one aircraft tail under one flight number. A **hop** is one leg (one Origin→Dest segment); an itinerary of N hops contains N+1 distinct airports.\n\nThe maximum observed hop count across all carriers and dates is **{{metric.max_hops}}**. All top-10 itineraries reach this ceiling, meaning no aircraft exceeded 8 legs in a single day under a single flight number.\n\n## Most Recent Maximum-Hop Itinerary\n\n| Field | Value |\n|---|---|\n| Carrier | {{metric.most_recent_carrier}} |\n| Flight Number | {{metric.most_recent_flight}} |\n| Tail Number | {{metric.most_recent_tail}} |\n| Date | {{metric.most_recent_date}} |\n| Route | {{metric.most_recent_route}} |\n| Departure Times | {{metric.most_recent_dep_times}} |\n\n## Repeated Operating Pattern vs. One-Off\n\nThe 8-hop ceiling is **not a one-off event** — it appears across multiple carriers, years, and flight numbers. The single most repeated pattern in the top 10 is **{{metric.top_recurring_flight}}**, which flew the identical {{metric.max_hops}}-hop route {{metric.top_recurring_count}} times:\n\n\u003e {{metric.top_recurring_route}}\n\nThis signals a scheduled, recurring turnaround pattern rather than an irregular re-routing.\n\n## Top 10 Longest Itineraries\n\n{{result_table_md}}\n\n## Key Observations\n\n- All top-10 itineraries are operated by **{{metric.dominant_carrier}}**, consistent with Southwest's point-to-point network model where a single aircraft and flight number chain many short segments in one day.\n- Route clustering around certain city pairs (e.g., the recurring CLE–DEN spine) indicates structured schedule blocks, not ad-hoc assignments.\n- Departure times span roughly 05:00–23:00 local, confirming full-day aircraft utilization on these itineraries.\n\n---\n_Generated: {{generated_at}}_",
  "metrics": {
    "summary_facts": [
      "The maximum hop count is 8, reached by multiple Southwest Airlines (WN) flights.",
      "The most recent 8-hop itinerary is WN flight 366 on 2024-12-01 (tail N957WN): ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA.",
      "WN flight 3149 flew the identical CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN route 4 times in January–February 2024, confirming a recurring scheduled pattern.",
      "WN flight 2787 and WN flight 154 each appear 3 times and 2 times respectively in the top 10, also with identical repeated routes.",
      "All 10 itineraries are Southwest Airlines (WN), reflecting their point-to-point model of chaining many short legs under one flight number."
    ],
    "named_values": {
      "dominant_carrier": "WN (Southwest Airlines)",
      "max_hops": "8",
      "most_recent_carrier": "WN (Southwest Airlines)",
      "most_recent_date": "2024-12-01",
      "most_recent_dep_times": "ISP 05:43 | BWI 08:10 | MYR 10:20 | BNA 11:42 | VPS 14:01 | DAL 16:43 | LAS 18:28 | OAK 20:41",
      "most_recent_flight": "366",
      "most_recent_route": "ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA",
      "most_recent_tail": "N957WN",
      "top_recurring_count": "4",
      "top_recurring_flight": "WN 3149",
      "top_recurring_route": "CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN"
    },
    "named_lists": {
      "repeated_routes": [
        "WN 3149 (×4): CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN",
        "WN 2787 (×3): MSY -\u003e ATL -\u003e CMH -\u003e BWI -\u003e RDU -\u003e BNA -\u003e DTW -\u003e MDW -\u003e LAX",
        "WN 154 (×2): ELP -\u003e DAL -\u003e LIT -\u003e ATL -\u003e RIC -\u003e MDW -\u003e MCI -\u003e PHX -\u003e SAN"
      ]
    }
  }
}
```

Saved report template to respect:

```report
# Highest Daily Hops for One Aircraft on One Flight Number

## Overview

This report identifies the longest single-day itineraries flown by one aircraft tail under one flight number. A **hop** is one leg (one Origin→Dest segment); an itinerary of N hops contains N+1 distinct airports.

The maximum observed hop count across all carriers and dates is **{{metric.max_hops}}**. All top-10 itineraries reach this ceiling, meaning no aircraft exceeded 8 legs in a single day under a single flight number.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | {{metric.most_recent_carrier}} |
| Flight Number | {{metric.most_recent_flight}} |
| Tail Number | {{metric.most_recent_tail}} |
| Date | {{metric.most_recent_date}} |
| Route | {{metric.most_recent_route}} |
| Departure Times | {{metric.most_recent_dep_times}} |

## Repeated Operating Pattern vs. One-Off

The 8-hop ceiling is **not a one-off event** — it appears across multiple carriers, years, and flight numbers. The single most repeated pattern in the top 10 is **{{metric.top_recurring_flight}}**, which flew the identical {{metric.max_hops}}-hop route {{metric.top_recurring_count}} times:

> {{metric.top_recurring_route}}

This signals a scheduled, recurring turnaround pattern rather than an irregular re-routing.

## Top 10 Longest Itineraries

{{result_table_md}}

## Key Observations

- All top-10 itineraries are operated by **{{metric.dominant_carrier}}**, consistent with Southwest's point-to-point network model where a single aircraft and flight number chain many short segments in one day.
- Route clustering around certain city pairs (e.g., the recurring CLE–DEN spine) indicates structured schedule blocks, not ad-hoc assignments.
- Departure times span roughly 05:00–23:00 local, confirming full-day aircraft utilization on these itineraries.

---
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
WITH leg_data AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS hops,
        groupArray((toUInt32(ifNull(DepTime, ifNull(CRSDepTime, 9999))), Origin, Dest)) AS legs_raw
    FROM ontime.ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
    GROUP BY
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate
),
itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        hops,
        arraySort(x -> x.1, legs_raw) AS legs_sorted
    FROM leg_data
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    hops,
    arrayStringConcat(
        arrayConcat(
            [legs_sorted[1].2],
            arrayMap(x -> x.3, legs_sorted)
        ),
        ' -> '
    ) AS Route,
    arrayStringConcat(
        arrayMap(
            x -> concat(x.2, ' ', lpad(toString(intDiv(x.1, 100)), 2, '0'), ':', lpad(toString(x.1 % 100), 2, '0')),
            legs_sorted
        ),
        ' | '
    ) AS DepartureTimes
FROM itineraries
ORDER BY hops DESC, FlightDate DESC
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