- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create browser-ready HTML `visual.html` using the proper  `*-analyst-dashboard` skill.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.

### Rules

- Question title: `Highest daily hops for one aircraft on one flight number`
- Visual mode: `dynamic`
- Visual type: `html_map`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

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

### Data Source

SQL query for primary data source:

```sql
WITH leg_rows AS (
    SELECT
        TailNum,
        FlightNum,
        Carrier,
        FlightDate,
        OriginCode,
        DestCode,
        addMinutes(toDateTime(FlightDate), intDiv(ifNull(DepTime, CRSDepTime), 100) * 60 + modulo(ifNull(DepTime, CRSDepTime), 100)) AS dep_ts
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND Diverted = 0
      AND TailNum != ''
      AND FlightNum != ''
      AND ifNull(DepTime, CRSDepTime) IS NOT NULL
), itineraries AS (
    SELECT
        TailNum,
        FlightNum,
        Carrier,
        FlightDate,
        arraySort(x -> x.1, groupArray((dep_ts, OriginCode, DestCode))) AS ordered_legs
    FROM leg_rows
    GROUP BY
        TailNum,
        FlightNum,
        Carrier,
        FlightDate
    HAVING length(ordered_legs) > 1
)
SELECT
    TailNum AS `Aircraft ID`,
    FlightNum AS `Flight Number`,
    Carrier,
    FlightDate AS `Date`,
    concat(
        arrayStringConcat(arrayMap(x -> x.2, ordered_legs), '→'),
        '→',
        tupleElement(arrayElement(ordered_legs, length(ordered_legs)), 3)
    ) AS Route,
    length(ordered_legs) AS hop_count
FROM itineraries
ORDER BY hop_count DESC, `Date` DESC, `Aircraft ID`, `Flight Number`
LIMIT 10
```

Data example/snippet:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": [
    "Aircraft ID",
    "Flight Number",
    "Carrier",
    "Date",
    "Route",
    "hop_count"
  ],
  "row_count": 10,
  "sample_rows": [
    {
      "Aircraft ID": "N957WN",
      "Carrier": "WN",
      "Date": "2024-12-01T00:00:00Z",
      "Flight Number": "366",
      "Route": "ISP→BWI→MYR→BNA→VPS→DAL→LAS→OAK→SEA",
      "hop_count": 8
    },
    {
      "Aircraft ID": "N7835A",
      "Carrier": "WN",
      "Date": "2024-02-18T00:00:00Z",
      "Flight Number": "3149",
      "Route": "CLE→BNA→PNS→HOU→MCI→PHX→BUR→OAK→DEN",
      "hop_count": 8
    }
  ],
  "field_shape_notes": {
    "Date": "ISO-like timestamp string"
  },
  "mode_hint": "Dynamic mode still fetches live data in the browser via query.sql and the configured endpoint."
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.