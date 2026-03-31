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

Generate only the visual artifact.

The analytical run already produced:

- `query.sql`
- `visual_input.json`

Use `query.sql` as the authoritative input for the primary data query.
Use `visual_input.json` to understand the result shape before building visuals.
You may construct additional queries when needed for enrichment or drill-down, but do not regenerate the saved SQL.

Return exactly this fenced section:

```html
<!doctype html>
<html>...</html>
```

Visual input summary:

```json
{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": [
    "Tail_Number",
    "Flight_Number_Reporting_Airline",
    "IATA_CODE_Reporting_Airline",
    "FlightDate",
    "Hops",
    "Route",
    "DepTimes"
  ],
  "row_count": 10,
  "sample_rows": [
    {
      "DepTimes": "5:43, 8:10, 10:20, 11:42, 14:01, 16:43, 18:28, 20:41",
      "FlightDate": "2024-12-01T00:00:00Z",
      "Flight_Number_Reporting_Airline": "366",
      "Hops": 8,
      "IATA_CODE_Reporting_Airline": "WN",
      "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
      "Tail_Number": "N957WN"
    },
    {
      "DepTimes": "6:21, 8:01, 10:07, 12:34, 15:14, 17:47, 19:02, 21:17",
      "FlightDate": "2024-02-18T00:00:00Z",
      "Flight_Number_Reporting_Airline": "3149",
      "Hops": 8,
      "IATA_CODE_Reporting_Airline": "WN",
      "Route": "CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN",
      "Tail_Number": "N7835A"
    }
  ],
  "field_shape_notes": {
    "FlightDate": "ISO-like timestamp string"
  },
  "mode_hint": "Dynamic mode still fetches live data in the browser via query.sql and the configured endpoint."
}
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
        DepTime
    FROM ontime.ontime
    WHERE Tail_Number != '' AND Cancelled = 0 AND DepTime IS NOT NULL
),
counts AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops
    FROM legs
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    ORDER BY Hops DESC, FlightDate DESC
    LIMIT 10
),
top_legs AS (
    SELECT l.*
    FROM legs l
    INNER JOIN counts c
        ON l.Tail_Number = c.Tail_Number
        AND l.Flight_Number_Reporting_Airline = c.Flight_Number_Reporting_Airline
        AND l.IATA_CODE_Reporting_Airline = c.IATA_CODE_Reporting_Airline
        AND l.FlightDate = c.FlightDate
    ORDER BY l.Tail_Number, l.Flight_Number_Reporting_Airline, l.FlightDate, l.DepTime
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    count() AS Hops,
    arrayStringConcat(
        arrayConcat(
            groupArray(Origin),
            [arrayElement(groupArray(Dest), toUInt32(count()))]
        ),
        '-'
    ) AS Route,
    arrayStringConcat(
        arrayMap(x -> concat(toString(intDiv(x, 100)), ':', leftPad(toString(x % 100), 2, '0')),
            groupArray(DepTime)
        ),
        ', '
    ) AS DepTimes
FROM top_legs
GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
ORDER BY Hops DESC, FlightDate DESC
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

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.