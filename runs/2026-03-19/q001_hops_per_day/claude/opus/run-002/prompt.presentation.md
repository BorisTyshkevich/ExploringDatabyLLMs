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
- Result columns: `Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate, Hops, Route, DepSchedule`

Saved analysis artifact:

```json
{
  "sql": "WITH sorted AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline,\n        IATA_CODE_Reporting_Airline,\n        FlightDate,\n        count() AS Hops,\n        arraySort(x -\u003e x.1, groupArray(tuple(DepTime, Origin, Dest))) AS legs\n    FROM ontime.ontime\n    WHERE Tail_Number != '' AND DepTime IS NOT NULL\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate\n    HAVING Hops \u003e= 6\n)\nSELECT\n    Tail_Number,\n    Flight_Number_Reporting_Airline,\n    IATA_CODE_Reporting_Airline,\n    FlightDate,\n    Hops,\n    arrayStringConcat(arrayMap(x -\u003e x.2, legs), ' -\u003e ')\n        || ' -\u003e ' || legs[length(legs)].3 AS Route,\n    arrayStringConcat(\n        arrayMap(\n            x -\u003e concat(x.2, ' ', lpad(toString(x.1 DIV 100), 2, '0'), ':', lpad(toString(x.1 MOD 100), 2, '0')),\n            legs\n        ), ', '\n    ) AS DepSchedule\nFROM sorted\nORDER BY Hops DESC, FlightDate DESC\nLIMIT 10",
  "report_markdown": "# {{question_title}}\n\n{{data_overview_md}}\n\n## Key finding\n\nThe maximum observed hop count is **{{metric.max_hops}}** legs in a single day for one aircraft operating under one flight number. All {{metric.max_hops}}-hop itineraries belong to **{{metric.carrier}}** (Southwest Airlines), reflecting its characteristic point-to-point, multi-stop scheduling model.\n\n## Most recent maximum-hop itinerary\n\nThe most recent {{metric.max_hops}}-hop itinerary was **{{metric.most_recent_carrier}} {{metric.most_recent_flight}}** on **{{metric.most_recent_date}}**, aircraft **{{metric.most_recent_tail}}**, flying the route:\n\n\u003e {{metric.most_recent_route}}\n\nDeparture times from each origin: {{metric.most_recent_schedule}}\n\nThe aircraft began its day at ISP (Long Island MacArthur) at 05:43 and reached SEA (Seattle-Tacoma) after 20:41, covering eight legs coast-to-coast over roughly 15 hours.\n\n## Route repetition and clustering\n\nAmong the top 10 longest itineraries, strong route repetition is evident:\n\n- **WN 3149** (CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN) appears **4 times** across January-February 2024, indicating a recurring weekly schedule pattern.\n- **WN 2787** (MSY -\u003e ATL -\u003e CMH -\u003e BWI -\u003e RDU -\u003e BNA -\u003e DTW -\u003e MDW -\u003e LAX) appears **3 times** across September-October 2022.\n- **WN 154** (ELP -\u003e DAL -\u003e LIT -\u003e ATL -\u003e RIC -\u003e MDW -\u003e MCI -\u003e PHX -\u003e SAN) appears **twice** in April 2023.\n- **WN 366** appears once with a distinct route (ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA).\n\nThese are clearly **repeating weekly operating patterns**, not one-off itineraries. Southwest regularly assigns a single flight number to an aircraft that hops across the country through 8 or more cities in a single day.\n\n## Top 10 results\n\n{{result_table_md}}",
  "metrics": {
    "summary_facts": [
      "The maximum hop count is 8 legs in a single day for one aircraft on one flight number.",
      "All top-10 itineraries belong to Southwest Airlines (WN).",
      "WN 3149 repeats the same 8-hop CLE-to-DEN route 4 times in Jan-Feb 2024.",
      "WN 2787 repeats the same 8-hop MSY-to-LAX route 3 times in Sep-Oct 2022."
    ],
    "named_values": {
      "carrier": "WN (Southwest Airlines)",
      "max_hops": "8",
      "most_recent_carrier": "WN",
      "most_recent_date": "2024-12-01",
      "most_recent_flight": "366",
      "most_recent_route": "ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA",
      "most_recent_schedule": "ISP 05:43, BWI 08:10, MYR 10:20, BNA 11:42, VPS 14:01, DAL 16:43, LAS 18:28, OAK 20:41",
      "most_recent_tail": "N957WN"
    },
    "named_lists": {
      "repeated_routes": [
        "WN 3149: CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN (4 occurrences)",
        "WN 2787: MSY -\u003e ATL -\u003e CMH -\u003e BWI -\u003e RDU -\u003e BNA -\u003e DTW -\u003e MDW -\u003e LAX (3 occurrences)",
        "WN 154: ELP -\u003e DAL -\u003e LIT -\u003e ATL -\u003e RIC -\u003e MDW -\u003e MCI -\u003e PHX -\u003e SAN (2 occurrences)",
        "WN 366: ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA (1 occurrence)"
      ]
    }
  }
}
```

Saved report template to respect:

```report
# {{question_title}}

{{data_overview_md}}

## Key finding

The maximum observed hop count is **{{metric.max_hops}}** legs in a single day for one aircraft operating under one flight number. All {{metric.max_hops}}-hop itineraries belong to **{{metric.carrier}}** (Southwest Airlines), reflecting its characteristic point-to-point, multi-stop scheduling model.

## Most recent maximum-hop itinerary

The most recent {{metric.max_hops}}-hop itinerary was **{{metric.most_recent_carrier}} {{metric.most_recent_flight}}** on **{{metric.most_recent_date}}**, aircraft **{{metric.most_recent_tail}}**, flying the route:

> {{metric.most_recent_route}}

Departure times from each origin: {{metric.most_recent_schedule}}

The aircraft began its day at ISP (Long Island MacArthur) at 05:43 and reached SEA (Seattle-Tacoma) after 20:41, covering eight legs coast-to-coast over roughly 15 hours.

## Route repetition and clustering

Among the top 10 longest itineraries, strong route repetition is evident:

- **WN 3149** (CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN) appears **4 times** across January-February 2024, indicating a recurring weekly schedule pattern.
- **WN 2787** (MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX) appears **3 times** across September-October 2022.
- **WN 154** (ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN) appears **twice** in April 2023.
- **WN 366** appears once with a distinct route (ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA).

These are clearly **repeating weekly operating patterns**, not one-off itineraries. Southwest regularly assigns a single flight number to an aircraft that hops across the country through 8 or more cities in a single day.

## Top 10 results

{{result_table_md}}
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
WITH sorted AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arraySort(x -> x.1, groupArray(tuple(DepTime, Origin, Dest))) AS legs
    FROM ontime.ontime
    WHERE Tail_Number != '' AND DepTime IS NOT NULL
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    HAVING Hops >= 6
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    Hops,
    arrayStringConcat(arrayMap(x -> x.2, legs), ' -> ')
        || ' -> ' || legs[length(legs)].3 AS Route,
    arrayStringConcat(
        arrayMap(
            x -> concat(x.2, ' ', lpad(toString(x.1 DIV 100), 2, '0'), ':', lpad(toString(x.1 MOD 100), 2, '0')),
            legs
        ), ', '
    ) AS DepSchedule
FROM sorted
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