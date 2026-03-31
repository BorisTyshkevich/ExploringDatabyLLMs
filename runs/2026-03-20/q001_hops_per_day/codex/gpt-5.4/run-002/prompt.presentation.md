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
    "Maximum Hops Observed",
    "Maximum-Hop Itinerary Count"
  ],
  "row_count": 10,
  "sample_rows": [
    {
      "Aircraft ID": "N957WN",
      "Carrier": "WN",
      "Date": "2024-12-01T00:00:00Z",
      "Flight Number": "366",
      "Hops": 8,
      "Maximum Hops Observed": 8,
      "Maximum-Hop Itinerary Count": 9859,
      "Route": "05:43 ISP-\u003eBWI | 08:10 BWI-\u003eMYR | 10:20 MYR-\u003eBNA | 11:42 BNA-\u003eVPS | 14:01 VPS-\u003eDAL | 16:43 DAL-\u003eLAS | 18:28 LAS-\u003eOAK | 20:41 OAK-\u003eSEA"
    },
    {
      "Aircraft ID": "N7835A",
      "Carrier": "WN",
      "Date": "2024-02-18T00:00:00Z",
      "Flight Number": "3149",
      "Hops": 8,
      "Maximum Hops Observed": 8,
      "Maximum-Hop Itinerary Count": 9859,
      "Route": "06:21 CLE-\u003eBNA | 08:01 BNA-\u003ePNS | 10:07 PNS-\u003eHOU | 12:34 HOU-\u003eMCI | 15:14 MCI-\u003ePHX | 17:47 PHX-\u003eBUR | 19:02 BUR-\u003eOAK | 21:17 OAK-\u003eDEN"
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
WITH airport_offsets AS
(
    SELECT
        airport_id,
        any(utc_local_time_variation) AS utc_local_time_variation
    FROM ontime.airports_latest
    GROUP BY airport_id
),
legs AS
(
    SELECT
        o.FlightDate,
        o.TailNum,
        o.FlightNum,
        o.Carrier,
        o.OriginAirportID,
        o.DestAirportID,
        o.Origin,
        o.Dest,
        o.DepTime,
        toDateTime(o.FlightDate)
            + toIntervalDay(if(o.DepTime = 2400, 1, 0))
            + toIntervalMinute(intDiv(if(o.DepTime = 2400, 0, o.DepTime), 100) * 60 + modulo(if(o.DepTime = 2400, 0, o.DepTime), 100)) AS dep_local_ts,
        (
            if(length(a1o.utc_local_time_variation) = 5,
                if(substring(a1o.utc_local_time_variation, 1, 1) = '-', -1, 1)
                * (toInt32(substring(a1o.utc_local_time_variation, 2, 2)) * 60 + toInt32(substring(a1o.utc_local_time_variation, 4, 2))),
                0
            )
        ) AS origin_offset_minutes
    FROM ontime.ontime AS o
    LEFT JOIN airport_offsets AS a1o ON o.OriginAirportID = a1o.airport_id
    WHERE o.Cancelled = 0
      AND o.Diverted = 0
      AND o.DepTime IS NOT NULL
      AND o.TailNum != ''
      AND o.FlightNum != ''
),
itineraries AS
(
    SELECT
        FlightDate,
        TailNum,
        FlightNum,
        Carrier,
        length(ordered_legs) AS Hops,
        arrayStringConcat(
            arrayMap(x -> concat(formatDateTime(x.2, '%H:%i'), ' ', x.3, '->', x.4), ordered_legs),
            ' | '
        ) AS Route
    FROM
    (
        SELECT
            FlightDate,
            TailNum,
            FlightNum,
            Carrier,
            arraySort(x -> x.1, groupArray((
                dep_local_ts - toIntervalMinute(origin_offset_minutes),
                dep_local_ts,
                Origin,
                Dest
            ))) AS ordered_legs
        FROM legs
        GROUP BY
            FlightDate,
            TailNum,
            FlightNum,
            Carrier
    )
),
max_hops AS
(
    SELECT max(Hops) AS max_hops_observed FROM itineraries
),
max_hops_counts AS
(
    SELECT count() AS max_hop_itinerary_count
    FROM itineraries
    CROSS JOIN max_hops
    WHERE Hops = max_hops_observed
)
SELECT
    TailNum AS `Aircraft ID`,
    FlightNum AS `Flight Number`,
    Carrier,
    FlightDate AS Date,
    Hops,
    Route,
    max_hops_observed AS `Maximum Hops Observed`,
    max_hop_itinerary_count AS `Maximum-Hop Itinerary Count`
FROM itineraries
CROSS JOIN max_hops
CROSS JOIN max_hops_counts
ORDER BY
    Hops DESC,
    Date DESC,
    Carrier,
    `Flight Number`,
    `Aircraft ID`
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

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.