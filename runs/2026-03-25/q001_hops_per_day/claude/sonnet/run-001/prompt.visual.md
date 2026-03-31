- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save.
- Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully.
- Do not write unchecked SQL.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `Highest daily hops for one aircraft on one flight number`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_map`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

- use the first dashboard-question proof query as the primary saved SQL already provided in the prompt
- use the other dashboard-question proof queries as supporting queries when they materially improve the narrative or supporting panels
- anchor the hero narrative and KPI strip to the top-ranked itinerary even when another itinerary is selected in the table
- show a lead-itinerary map that remains present even before airport-coordinate enrichment succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- include a narrative hero about the lead itinerary and the broader geographic pattern of the top itineraries
- label the map as airport-coordinate enrichment in the query ledger
- reuse the enrichment results for any itinerary selected from the primary result set without issuing a new per-click enrichment query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip anchored to the top-ranked result even when the selected itinerary changes
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- make the selected-row map behavior explicit: when the selected itinerary differs from Rank 1, the map title, plotted route, markers, bounds, and route detail panel must visibly update to that selected itinerary rather than leaving the lead route drawn
- keep the map/detail selection state separate from the anchored hero and KPI state
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- prefer the `Route` value from the primary query as the per-row itinerary representation for redraws
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis
- derive the ordered itinerary sequence for map redraws and the route detail panel by splitting `Route` on `-`

### Data Source

SQL query for primary data source:

```sql
WITH daily_itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        Reporting_Airline                AS Carrier,
        FlightDate,
        count()                          AS hop_count,
        arrayStringConcat(
            arrayConcat(
                [arrayElement(arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],
                arrayMap(x -> x.3, arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))
            ),
            '-'
        ) AS Route
    FROM ontime.fact_ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate
    HAVING hop_count >= 2
),
unique_routes AS (
    SELECT
        Route,
        max(hop_count)                   AS hop_count,
        argMax(Tail_Number, FlightDate)  AS Tail_Number,
        argMax(FlightNum,   FlightDate)  AS FlightNum,
        argMax(Carrier,     FlightDate)  AS Carrier,
        max(FlightDate)                  AS most_recent_date,
        count()                          AS occurrences
    FROM daily_itineraries
    GROUP BY Route
)
SELECT
    Tail_Number,
    FlightNum,
    Carrier,
    most_recent_date AS FlightDate,
    hop_count,
    Route,
    occurrences
FROM unique_routes
ORDER BY hop_count DESC, most_recent_date DESC
LIMIT 10
```

Data example/snippet:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "All top-10 unique routes tie at **8 hops** (9 airports). The leading example by recency is Southwest flight **WN 366**, tail **N957WN**, on **2024-12-01**: `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` — a coast-to-coast turn starting in Islip, NY and ending in Seattle.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "occurrences": 1
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN 366** (tail N957WN) on **2024-12-01**, route `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` with 8 hops. It is the only 8-hop unique route recorded in late 2024.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY most_recent_date DESC, hop_count DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2025-11-30T00:00:00Z",
        "FlightNum": "2179",
        "Route": "SAT-TPA-FLL-RDU-BNA-STL-AUS",
        "Tail_Number": "N925WN",
        "hop_count": 6,
        "occurrences": 1
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top 10 are a **mix of both**. Two routes are clearly recurring scheduled turn patterns: WN 2884 (`LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN`) appeared **46 times** and WN 1956 (`MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC`) appeared **40 times**, indicating they were regular daily assignments over many years. Others appear only 1–5 times, suggesting irregular or one-off operations. Southwest dominates all 10 entries, consistent with its high-utilization short-haul turn model.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "hop_count": 8,
        "occurrences": 1
      }
    }
  ]
}

### Multi-query additions

- The saved SQL shown below is the primary dashboard query for this page.
- The verified analysis package includes named supporting queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Use subquestion answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- The dashboard does not need to mirror `report.md`; it should combine narrative and interactive analysis.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "All top-10 unique routes tie at **8 hops** (9 airports). The leading example by recency is Southwest flight **WN 366**, tail **N957WN**, on **2024-12-01**: `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` — a coast-to-coast turn starting in Islip, NY and ending in Seattle.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "occurrences": 1
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN 366** (tail N957WN) on **2024-12-01**, route `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` with 8 hops. It is the only 8-hop unique route recorded in late 2024.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY most_recent_date DESC, hop_count DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2025-11-30T00:00:00Z",
        "FlightNum": "2179",
        "Route": "SAT-TPA-FLL-RDU-BNA-STL-AUS",
        "Tail_Number": "N925WN",
        "hop_count": 6,
        "occurrences": 1
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top 10 are a **mix of both**. Two routes are clearly recurring scheduled turn patterns: WN 2884 (`LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN`) appeared **46 times** and WN 1956 (`MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC`) appeared **40 times**, indicating they were regular daily assignments over many years. Others appear only 1–5 times, suggesting irregular or one-off operations. Southwest dominates all 10 entries, consistent with its high-utilization short-haul turn model.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "hop_count": 8,
        "occurrences": 1
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.

Create browser-ready HTML `visual.html`.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.