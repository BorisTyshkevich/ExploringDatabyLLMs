- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

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
- prefer a simple per-row itinerary representation from the primary query that the browser can reliably use for redraws
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis

### Data Source

SQL query for primary data source:

```sql
WITH ranked AS (
    SELECT
        FlightDate,
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        IATA_CODE_Reporting_Airline AS Carrier,
        count() AS hop_count,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),
                [argMax(DestCode, assumeNotNull(DepTime))]
            ),
            '-'
        ) AS Route,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),
                [argMax(DestCode, assumeNotNull(DepTime))]
            ),
            ','
        ) AS itinerary_sequence
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline
)
SELECT
    Tail_Number,
    FlightNum,
    Carrier,
    FlightDate,
    hop_count,
    Route,
    itinerary_sequence
FROM ranked
ORDER BY hop_count DESC, FlightDate DESC
LIMIT 10
```

Data example/snippet:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 4,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "The highest-hop example is Southwest Airlines (WN) flight 366 operated by tail **N957WN** on **2024-12-01**, with **8 hops** tracing **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA**. The aircraft started at Long Island MacArthur (ISP) early morning, stopped at Baltimore/Washington (BWI), Myrtle Beach (MYR), Nashville (BNA), Fort Walton Beach/Destin (VPS), Dallas Love Field (DAL), Las Vegas (LAS), and Oakland (OAK), finishing in Seattle/Tacoma (SEA) — a full cross-country marathon from the New York metro area to the Pacific Northwest. All other 8-hop itineraries in the top 10 also belong to Southwest.",
      "sql": "WITH ranked AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            '-'\n        ) AS Route,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            ','\n        ) AS itinerary_sequence\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n      AND Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n)\nSELECT\n    Tail_Number,\n    FlightNum,\n    Carrier,\n    FlightDate,\n    hop_count,\n    Route,\n    itinerary_sequence\nFROM ranked\nORDER BY hop_count DESC, FlightDate DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "itinerary_sequence"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "itinerary_sequence": "ISP,BWI,MYR,BNA,VPS,DAL,LAS,OAK,SEA"
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN flight 366 / tail N957WN on 2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA), which is also the lead itinerary. It is the only occurrence of that flight number in the top 10 and carries the latest date in the entire ranked set.",
      "sql": "WITH ranked AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            '-'\n        ) AS Route,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            ','\n        ) AS itinerary_sequence\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n      AND Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n    ORDER BY hop_count DESC, FlightDate DESC\n    LIMIT 10\n)\nSELECT Tail_Number, FlightNum, Carrier, FlightDate, hop_count, Route, itinerary_sequence\nFROM ranked\nORDER BY FlightDate DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "itinerary_sequence"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "itinerary_sequence": "ISP,BWI,MYR,BNA,VPS,DAL,LAS,OAK,SEA"
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top itineraries are strongly **recurring scheduled patterns**. Among the top 10 rows there are only 4 distinct routes. WN flight 3149 (CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN) ran on 4 consecutive Sundays in Jan–Feb 2024 with different tail numbers each time. WN flight 2787 (MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX) ran 3 times in Sep–Oct 2022. WN flight 154 (ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN) ran twice in Apr 2023. Only WN 366 appears once. The consistent weekly cadence and published flight numbers confirm these are regular Southwest scheduled turns, not anomalies.",
      "sql": "WITH ranked AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n      AND Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n    ORDER BY hop_count DESC, FlightDate DESC\n    LIMIT 10\n)\nSELECT\n    FlightNum,\n    Carrier,\n    Route,\n    count()            AS occurrences,\n    min(FlightDate)    AS first_seen,\n    max(FlightDate)    AS last_seen\nFROM ranked\nGROUP BY FlightNum, Carrier, Route\nORDER BY occurrences DESC",
      "row_count": 4,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "Route",
        "occurrences",
        "first_seen",
        "last_seen"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightNum": "3149",
        "Route": "CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN",
        "first_seen": "2024-01-14T00:00:00Z",
        "last_seen": "2024-02-18T00:00:00Z",
        "occurrences": 4
      }
    },
    {
      "id": "q4",
      "subquestion": "What geographic pattern do the top itineraries show?",
      "answer_markdown": "Every top itinerary is a Southwest Airlines **coast-to-coast diagonal sweep** across the continental US. Routes originate in the East or Southeast (Long Island NY, Cleveland OH, El Paso TX, New Orleans LA) and arc west through the South, Midwest, or Mid-Atlantic before terminating on the West Coast or Mountain West (Seattle, Denver, San Diego, Los Angeles). The lead route ISP→SEA traverses roughly 2,800 miles from the New York metro to the Pacific Northwest, passing through the Mid-Atlantic, Southeast Gulf Coast, South-Central, Southwest, and Pacific Coast in a single day. The pattern reflects Southwest's hub-and-spoke-free point-to-point network, where a single aircraft and flight number chain together many short segments to cover the full width of the country.",
      "sql": "SELECT\n    f.OriginCode,\n    f.DestCode,\n    f.DepTime,\n    o.DisplayAirportName  AS OriginName,\n    o.Latitude            AS OriginLat,\n    o.Longitude           AS OriginLon,\n    d.DisplayAirportName  AS DestName,\n    d.Latitude            AS DestLat,\n    d.Longitude           AS DestLon\nFROM ontime.fact_ontime f\nLEFT JOIN ontime.dim_airports o ON f.OriginCode = o.AirportCode\nLEFT JOIN ontime.dim_airports d ON f.DestCode   = d.AirportCode\nWHERE f.Cancelled = 0\n  AND f.Tail_Number = 'N957WN'\n  AND f.Flight_Number_Reporting_Airline = '366'\n  AND f.FlightDate = '2024-12-01'\nORDER BY assumeNotNull(f.DepTime)",
      "row_count": 8,
      "result_columns": [
        "OriginCode",
        "DestCode",
        "DepTime",
        "OriginName",
        "OriginLat",
        "OriginLon",
        "DestName",
        "DestLat",
        "DestLon"
      ],
      "first_row": {
        "DepTime": 543,
        "DestCode": "BWI",
        "DestLat": 39.17583333,
        "DestLon": -76.66888889,
        "DestName": "Baltimore/Washington International Thurgood Marshall",
        "OriginCode": "ISP",
        "OriginLat": 40.79611111,
        "OriginLon": -73.10055556,
        "OriginName": "Long Island MacArthur"
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
  "row_count": 4,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "The highest-hop example is Southwest Airlines (WN) flight 366 operated by tail **N957WN** on **2024-12-01**, with **8 hops** tracing **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA**. The aircraft started at Long Island MacArthur (ISP) early morning, stopped at Baltimore/Washington (BWI), Myrtle Beach (MYR), Nashville (BNA), Fort Walton Beach/Destin (VPS), Dallas Love Field (DAL), Las Vegas (LAS), and Oakland (OAK), finishing in Seattle/Tacoma (SEA) — a full cross-country marathon from the New York metro area to the Pacific Northwest. All other 8-hop itineraries in the top 10 also belong to Southwest.",
      "sql": "WITH ranked AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            '-'\n        ) AS Route,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            ','\n        ) AS itinerary_sequence\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n      AND Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n)\nSELECT\n    Tail_Number,\n    FlightNum,\n    Carrier,\n    FlightDate,\n    hop_count,\n    Route,\n    itinerary_sequence\nFROM ranked\nORDER BY hop_count DESC, FlightDate DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "itinerary_sequence"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "itinerary_sequence": "ISP,BWI,MYR,BNA,VPS,DAL,LAS,OAK,SEA"
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN flight 366 / tail N957WN on 2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA), which is also the lead itinerary. It is the only occurrence of that flight number in the top 10 and carries the latest date in the entire ranked set.",
      "sql": "WITH ranked AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            '-'\n        ) AS Route,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            ','\n        ) AS itinerary_sequence\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n      AND Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n    ORDER BY hop_count DESC, FlightDate DESC\n    LIMIT 10\n)\nSELECT Tail_Number, FlightNum, Carrier, FlightDate, hop_count, Route, itinerary_sequence\nFROM ranked\nORDER BY FlightDate DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "itinerary_sequence"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "itinerary_sequence": "ISP,BWI,MYR,BNA,VPS,DAL,LAS,OAK,SEA"
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top itineraries are strongly **recurring scheduled patterns**. Among the top 10 rows there are only 4 distinct routes. WN flight 3149 (CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN) ran on 4 consecutive Sundays in Jan–Feb 2024 with different tail numbers each time. WN flight 2787 (MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX) ran 3 times in Sep–Oct 2022. WN flight 154 (ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN) ran twice in Apr 2023. Only WN 366 appears once. The consistent weekly cadence and published flight numbers confirm these are regular Southwest scheduled turns, not anomalies.",
      "sql": "WITH ranked AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.2, arraySort(x -\u003e x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),\n                [argMax(DestCode, assumeNotNull(DepTime))]\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n      AND Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n    ORDER BY hop_count DESC, FlightDate DESC\n    LIMIT 10\n)\nSELECT\n    FlightNum,\n    Carrier,\n    Route,\n    count()            AS occurrences,\n    min(FlightDate)    AS first_seen,\n    max(FlightDate)    AS last_seen\nFROM ranked\nGROUP BY FlightNum, Carrier, Route\nORDER BY occurrences DESC",
      "row_count": 4,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "Route",
        "occurrences",
        "first_seen",
        "last_seen"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightNum": "3149",
        "Route": "CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN",
        "first_seen": "2024-01-14T00:00:00Z",
        "last_seen": "2024-02-18T00:00:00Z",
        "occurrences": 4
      }
    },
    {
      "id": "q4",
      "subquestion": "What geographic pattern do the top itineraries show?",
      "answer_markdown": "Every top itinerary is a Southwest Airlines **coast-to-coast diagonal sweep** across the continental US. Routes originate in the East or Southeast (Long Island NY, Cleveland OH, El Paso TX, New Orleans LA) and arc west through the South, Midwest, or Mid-Atlantic before terminating on the West Coast or Mountain West (Seattle, Denver, San Diego, Los Angeles). The lead route ISP→SEA traverses roughly 2,800 miles from the New York metro to the Pacific Northwest, passing through the Mid-Atlantic, Southeast Gulf Coast, South-Central, Southwest, and Pacific Coast in a single day. The pattern reflects Southwest's hub-and-spoke-free point-to-point network, where a single aircraft and flight number chain together many short segments to cover the full width of the country.",
      "sql": "SELECT\n    f.OriginCode,\n    f.DestCode,\n    f.DepTime,\n    o.DisplayAirportName  AS OriginName,\n    o.Latitude            AS OriginLat,\n    o.Longitude           AS OriginLon,\n    d.DisplayAirportName  AS DestName,\n    d.Latitude            AS DestLat,\n    d.Longitude           AS DestLon\nFROM ontime.fact_ontime f\nLEFT JOIN ontime.dim_airports o ON f.OriginCode = o.AirportCode\nLEFT JOIN ontime.dim_airports d ON f.DestCode   = d.AirportCode\nWHERE f.Cancelled = 0\n  AND f.Tail_Number = 'N957WN'\n  AND f.Flight_Number_Reporting_Airline = '366'\n  AND f.FlightDate = '2024-12-01'\nORDER BY assumeNotNull(f.DepTime)",
      "row_count": 8,
      "result_columns": [
        "OriginCode",
        "DestCode",
        "DepTime",
        "OriginName",
        "OriginLat",
        "OriginLon",
        "DestName",
        "DestLat",
        "DestLon"
      ],
      "first_row": {
        "DepTime": 543,
        "DestCode": "BWI",
        "DestLat": 39.17583333,
        "DestLon": -76.66888889,
        "DestName": "Baltimore/Washington International Thurgood Marshall",
        "OriginCode": "ISP",
        "OriginLat": 40.79611111,
        "OriginLon": -73.10055556,
        "OriginName": "Long Island MacArthur"
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