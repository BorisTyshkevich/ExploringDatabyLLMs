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
WITH legs AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        OriginCode,
        DestCode,
        assumeNotNull(CRSDepTime) AS dep
    FROM ontime.fact_ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
),
grouped AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        groupArray((dep, OriginCode, DestCode)) AS legs_arr,
        count() AS hop_count
    FROM legs
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    HAVING hop_count >= 2
)
SELECT
    Tail_Number AS "Aircraft ID",
    Flight_Number_Reporting_Airline AS "Flight Number",
    IATA_CODE_Reporting_Airline AS "Carrier",
    FlightDate AS "Date",
    arrayStringConcat(
        arrayPushBack(
            arrayMap(x -> x.2, arraySort(x -> x.1, legs_arr)),
            arrayElement(arraySort(x -> x.1, legs_arr), -1).3
        ),
        '-'
    ) AS "Route",
    hop_count
FROM grouped
ORDER BY hop_count DESC, FlightDate DESC
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
      "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
      "hop_count": 8
    },
    {
      "Aircraft ID": "N7835A",
      "Carrier": "WN",
      "Date": "2024-02-18T00:00:00Z",
      "Flight Number": "3149",
      "Route": "CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN",
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