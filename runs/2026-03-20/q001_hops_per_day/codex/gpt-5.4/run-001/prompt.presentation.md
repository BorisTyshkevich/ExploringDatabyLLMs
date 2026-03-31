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
    "Aircraft ID",
    "Flight Number",
    "Carrier",
    "Date",
    "Hops",
    "Route",
    "Departure Times From Origin",
    "Max Hops Overall",
    "Max Hop Itinerary Count",
    "Same-Hops Route Count",
    "Route Frequency In Top 10"
  ],
  "row_count": 10,
  "sample_rows": [
    {
      "Aircraft ID": "N957WN",
      "Carrier": "WN",
      "Date": "2024-12-01T00:00:00Z",
      "Departure Times From Origin": "05:43 ISP | 08:10 BWI | 10:20 MYR | 11:42 BNA | 14:01 VPS | 16:43 DAL | 18:28 LAS | 20:41 OAK",
      "Flight Number": "366",
      "Hops": 8,
      "Max Hop Itinerary Count": 9859,
      "Max Hops Overall": 8,
      "Route": "ISP -\u003e BWI -\u003e MYR -\u003e BNA -\u003e VPS -\u003e DAL -\u003e LAS -\u003e OAK -\u003e SEA",
      "Route Frequency In Top 10": 1,
      "Same-Hops Route Count": 1
    },
    {
      "Aircraft ID": "N7835A",
      "Carrier": "WN",
      "Date": "2024-02-18T00:00:00Z",
      "Departure Times From Origin": "06:21 CLE | 08:01 BNA | 10:07 PNS | 12:34 HOU | 15:14 MCI | 17:47 PHX | 19:02 BUR | 21:17 OAK",
      "Flight Number": "3149",
      "Hops": 8,
      "Max Hop Itinerary Count": 9859,
      "Max Hops Overall": 8,
      "Route": "CLE -\u003e BNA -\u003e PNS -\u003e HOU -\u003e MCI -\u003e PHX -\u003e BUR -\u003e OAK -\u003e DEN",
      "Route Frequency In Top 10": 4,
      "Same-Hops Route Count": 4
    }
  ],
  "field_shape_notes": {
    "Date": "ISO-like timestamp string"
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
        o.Tail_Number,
        o.Flight_Number_Reporting_Airline,
        o.IATA_CODE_Reporting_Airline,
        o.FlightDate,
        toDateTime(
            concat(
                toString(o.FlightDate),
                ' ',
                leftPad(toString(intDiv(o.DepTime, 100)), 2, '0'),
                ':',
                leftPad(toString(o.DepTime % 100), 2, '0'),
                ':00'
            )
        ) AS dep_ts,
        concat(
            leftPad(toString(intDiv(o.DepTime, 100)), 2, '0'),
            ':',
            leftPad(toString(o.DepTime % 100), 2, '0')
        ) AS dep_hhmm,
        o.Origin,
        o.Dest
    FROM ontime.ontime AS o
    WHERE o.Cancelled = 0
      AND o.Diverted = 0
      AND o.DepTime IS NOT NULL
      AND o.Tail_Number != ''
      AND o.Flight_Number_Reporting_Airline != ''
),
itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arraySort(
            x -> (x.1, x.2, x.3),
            groupArray((dep_ts, toString(Origin), toString(Dest), dep_hhmm))
        ) AS ordered_legs,
        arrayStringConcat(
            arrayConcat([ordered_legs[1].2], arrayMap(x -> x.3, ordered_legs)),
            ' -> '
        ) AS Route,
        arrayStringConcat(
            arrayMap(x -> concat(x.4, ' ', x.2), ordered_legs),
            ' | '
        ) AS `Departure Times From Origin`
    FROM legs
    GROUP BY
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate
),
max_hops AS (
    SELECT max(Hops) AS value
    FROM itineraries
),
scored AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Hops,
        Route,
        `Departure Times From Origin`,
        (SELECT value FROM max_hops) AS `Max Hops Overall`,
        (SELECT count() FROM itineraries WHERE Hops = (SELECT value FROM max_hops)) AS `Max Hop Itinerary Count`,
        count() OVER (PARTITION BY Hops, Route) AS `Same-Hops Route Count`
    FROM itineraries
),
top10 AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Hops,
        Route,
        `Departure Times From Origin`,
        `Max Hops Overall`,
        `Max Hop Itinerary Count`,
        `Same-Hops Route Count`
    FROM scored
    ORDER BY Hops DESC, FlightDate DESC, Tail_Number DESC
    LIMIT 10
)
SELECT
    Tail_Number AS `Aircraft ID`,
    Flight_Number_Reporting_Airline AS `Flight Number`,
    IATA_CODE_Reporting_Airline AS Carrier,
    FlightDate AS Date,
    Hops,
    Route,
    `Departure Times From Origin`,
    `Max Hops Overall`,
    `Max Hop Itinerary Count`,
    `Same-Hops Route Count`,
    count() OVER (PARTITION BY Route) AS `Route Frequency In Top 10`
FROM top10
ORDER BY Hops DESC, Date DESC, `Aircraft ID` DESC
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