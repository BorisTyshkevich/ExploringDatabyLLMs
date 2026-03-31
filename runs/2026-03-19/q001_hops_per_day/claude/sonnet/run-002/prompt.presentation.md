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
- Result columns: `FlightDate, Tail_Number, Carrier, FlightNum, Hops, Route, DepartureTimes`

Saved analysis artifact:

```json
{
  "sql": "WITH itineraries AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        count() AS Hops,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(\n                    x -\u003e x.2,\n                    arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), Origin)))\n                ),\n                [arrayElement(\n                    arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), Dest)))),\n                    -1\n                )]\n            ),\n            ' -\u003e '\n        ) AS Route,\n        arrayStringConcat(\n            arrayMap(\n                x -\u003e x.2 || '@' || toString(x.1),\n                arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), Origin)))\n            ),\n            ', '\n        ) AS DepartureTimes\n    FROM ontime.ontime\n    WHERE Cancelled = 0\n      AND Tail_Number != ''\n      AND DepTime IS NOT NULL\n    GROUP BY FlightDate, Tail_Number, Carrier, FlightNum\n)\nSELECT\n    FlightDate,\n    Tail_Number,\n    Carrier,\n    FlightNum,\n    Hops,\n    Route,\n    DepartureTimes\nFROM itineraries\nORDER BY Hops DESC, FlightDate DESC\nLIMIT 10",
  "report_markdown": "# Highest Daily Hops for One Aircraft on One Flight Number\n\n{{data_overview_md}}\n\n## Maximum Hop Count\n\nThe highest recorded hop count is **{{metric.max_hops}}** legs flown by a single aircraft under a single flight number in one day. All top-10 entries belong to {{metric.dominant_carrier}}. This is not a one-off anomaly: the same 8-hop routing patterns recur across multiple weeks on the same flight numbers, indicating deliberate scheduled network itineraries.\n\n## Most Recent Maximum-Hop Itinerary\n\n| Field | Value |\n|---|---|\n| Date | {{metric.most_recent_date}} |\n| Carrier | {{metric.most_recent_carrier}} |\n| Flight Number | {{metric.most_recent_flight_num}} |\n| Aircraft (Tail) | {{metric.most_recent_aircraft}} |\n| Route | {{metric.most_recent_route}} |\n\n## Top 10 Longest and Most Recent Itineraries\n\n{{result_table_md}}\n\n*DepartureTimes column shows `AIRPORT@HHMM` for each origin leg in chronological order.*\n\n## Route Repetition and Clustering\n\nThree distinct 8-hop route patterns recur across different aircraft on the same flight number, confirming them as recurring scheduled itineraries. WN 3149 (CLE → DEN corridor) is the most frequently repeated, appearing on consecutive Sundays in early 2024 with consistent city pairs from the Midwest through the South and out to the West Coast. WN 2787 (MSY → LAX corridor) and WN 154 (ELP → SAN corridor) each recur across separate weeks as well. This tight clustering by flight number and consistent routing is characteristic of deliberate point-to-point network scheduling.",
  "metrics": {
    "summary_facts": [
      "The maximum hop count is 8, observed across multiple Southwest Airlines itineraries.",
      "WN flight 3149 operated the same 8-hop CLE-\u003eDEN route on 4 separate Sundays in Jan-Feb 2024.",
      "All top-10 maximum-hop itineraries are operated by Southwest Airlines (WN).",
      "Three distinct 8-hop routes recur: WN 3149, WN 2787, and WN 154."
    ],
    "named_values": {
      "dominant_carrier": "Southwest Airlines (WN)",
      "max_hops": "8",
      "most_recent_aircraft": "N957WN",
      "most_recent_carrier": "WN",
      "most_recent_date": "2024-12-01",
      "most_recent_flight_num": "366",
      "most_recent_route": "ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA"
    },
    "named_lists": {
      "repeated_routes": [
        "WN 3149: CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN (4 occurrences, Jan-Feb 2024)",
        "WN 2787: MSY -\u003e ATL -\u003e CMH -\u003e BWI -\u003e RDU -\u003e BNA -\u003e DTW -\u003e MDW -\u003e LAX (3 occurrences, Sep-Oct 2022)",
        "WN 154: ELP -\u003e DAL -\u003e LIT -\u003e ATL -\u003e RIC -\u003e MDW -\u003e MCI -\u003e PHX -\u003e SAN (2 occurrences, Apr 2023)"
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

The highest recorded hop count is **{{metric.max_hops}}** legs flown by a single aircraft under a single flight number in one day. All top-10 entries belong to {{metric.dominant_carrier}}. This is not a one-off anomaly: the same 8-hop routing patterns recur across multiple weeks on the same flight numbers, indicating deliberate scheduled network itineraries.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Date | {{metric.most_recent_date}} |
| Carrier | {{metric.most_recent_carrier}} |
| Flight Number | {{metric.most_recent_flight_num}} |
| Aircraft (Tail) | {{metric.most_recent_aircraft}} |
| Route | {{metric.most_recent_route}} |

## Top 10 Longest and Most Recent Itineraries

{{result_table_md}}

*DepartureTimes column shows `AIRPORT@HHMM` for each origin leg in chronological order.*

## Route Repetition and Clustering

Three distinct 8-hop route patterns recur across different aircraft on the same flight number, confirming them as recurring scheduled itineraries. WN 3149 (CLE → DEN corridor) is the most frequently repeated, appearing on consecutive Sundays in early 2024 with consistent city pairs from the Midwest through the South and out to the West Coast. WN 2787 (MSY → LAX corridor) and WN 154 (ELP → SAN corridor) each recur across separate weeks as well. This tight clustering by flight number and consistent routing is characteristic of deliberate point-to-point network scheduling.
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
        FlightDate,
        Tail_Number,
        IATA_CODE_Reporting_Airline AS Carrier,
        Flight_Number_Reporting_Airline AS FlightNum,
        count() AS Hops,
        arrayStringConcat(
            arrayConcat(
                arrayMap(
                    x -> x.2,
                    arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Origin)))
                ),
                [arrayElement(
                    arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Dest)))),
                    -1
                )]
            ),
            ' -> '
        ) AS Route,
        arrayStringConcat(
            arrayMap(
                x -> x.2 || '@' || toString(x.1),
                arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Origin)))
            ),
            ', '
        ) AS DepartureTimes
    FROM ontime.ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND DepTime IS NOT NULL
    GROUP BY FlightDate, Tail_Number, Carrier, FlightNum
)
SELECT
    FlightDate,
    Tail_Number,
    Carrier,
    FlightNum,
    Hops,
    Route,
    DepartureTimes
FROM itineraries
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