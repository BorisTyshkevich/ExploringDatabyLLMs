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
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis

### Data Source

SQL query for primary data source:

```sql
WITH itineraries AS (
    SELECT
        FlightDate,
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        IATA_CODE_Reporting_Airline AS Carrier,
        count() AS HopCount,
        arrayStringConcat(
            arrayMap(x -> x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))
            || [arrayElement(arrayMap(x -> x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],
            '-'
        ) AS Route,
        arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode, DestCode))) AS LegSeq
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
        AND Tail_Number != ''
        AND Flight_Number_Reporting_Airline != ''
    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline
),
ranked AS (
    SELECT
        row_number() OVER (ORDER BY HopCount DESC, FlightDate DESC) AS Rank,
        FlightDate,
        Tail_Number,
        FlightNum,
        Carrier,
        HopCount,
        Route,
        LegSeq
    FROM itineraries
)
SELECT *
FROM ranked
WHERE Rank <= 10
ORDER BY Rank
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
      "answer_markdown": "The highest-hop example is **WN flight 366 on 2024-12-01**, operated by tail **N957WN**, covering **8 hops** (9 airports) along the route **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA**. The aircraft started at Long Island MacArthur (ISP), made stops in Baltimore (BWI), Myrtle Beach (MYR), Nashville (BNA), Fort Walton Beach (VPS), Dallas Love Field (DAL), Las Vegas (LAS), and Oakland (OAK), before finishing in Seattle (SEA). All top-10 itineraries share the same 8-hop maximum and are Southwest Airlines operations.",
      "sql": "WITH itineraries AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS HopCount,\n        arrayStringConcat(\n            arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))\n            || [arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],\n            '-'\n        ) AS Route,\n        arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode, DestCode))) AS LegSeq\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n        AND Tail_Number != ''\n        AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n),\nranked AS (\n    SELECT\n        row_number() OVER (ORDER BY HopCount DESC, FlightDate DESC) AS Rank,\n        FlightDate,\n        Tail_Number,\n        FlightNum,\n        Carrier,\n        HopCount,\n        Route,\n        LegSeq\n    FROM itineraries\n)\nSELECT *\nFROM ranked\nWHERE Rank \u003c= 10\nORDER BY Rank",
      "row_count": 10,
      "result_columns": [
        "Rank",
        "FlightDate",
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "HopCount",
        "Route",
        "LegSeq"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "HopCount": 8,
        "LegSeq": [
          {
            "": "BWI"
          },
          {
            "": "MYR"
          },
          {
            "": "BNA"
          },
          {
            "": "VPS"
          },
          {
            "": "DAL"
          },
          {
            "": "LAS"
          },
          {
            "": "OAK"
          },
          {
            "": "SEA"
          }
        ],
        "Rank": 1,
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN"
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN flight 366 on 2024-12-01** (tail N957WN), route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA. It is both Rank 1 (the lead itinerary) and the most recent date in the top-10, ahead of the next most recent (WN 3149 on 2024-02-18).",
      "sql": "WITH itineraries AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS HopCount,\n        arrayStringConcat(\n            arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))\n            || [arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n        AND Tail_Number != ''\n        AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n)\nSELECT FlightDate, Tail_Number, FlightNum, Carrier, HopCount, Route\nFROM itineraries\nORDER BY HopCount DESC, FlightDate DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "FlightDate",
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "HopCount",
        "Route"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "HopCount": 8,
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN"
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top itineraries are overwhelmingly **recurring scheduled patterns**, not one-offs. The most frequent 8-hop route (WN 2215: CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO) ran **61 times** between August and October 2009. WN 1923 ran 60 times, WN 2558 ran 55 times. These routes operated nearly every day of a seasonal schedule window, indicating deliberate long-day rotations assigned to a flight number — not anomalies.",
      "sql": "WITH itineraries AS (\n    SELECT\n        FlightDate,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS HopCount,\n        arrayStringConcat(\n            arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))\n            || [arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n        AND Tail_Number != ''\n        AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n    HAVING HopCount = 8\n)\nSELECT FlightNum, Carrier, Route, count() AS DateCount,\n       min(FlightDate) AS FirstSeen, max(FlightDate) AS LastSeen\nFROM itineraries\nGROUP BY FlightNum, Carrier, Route\nORDER BY DateCount DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "Route",
        "DateCount",
        "FirstSeen",
        "LastSeen"
      ],
      "first_row": {
        "Carrier": "WN",
        "DateCount": 61,
        "FirstSeen": "2009-08-17T00:00:00Z",
        "FlightNum": "2215",
        "LastSeen": "2009-10-30T00:00:00Z",
        "Route": "CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO"
      }
    },
    {
      "id": "q4",
      "subquestion": "What geographic pattern do the top itineraries show?",
      "answer_markdown": "All 8-hop itineraries trace **long east-to-west (or southeast-to-northwest) transcontinental sweeps** across the continental United States. Routes consistently originate in the Eastern US (e.g., Columbus OH, Pittsburgh PA, Hartford CT, Fort Lauderdale FL, Richmond VA) and terminate at West Coast or Pacific Northwest airports (e.g., Reno, San Jose, Oakland, Portland, Ontario CA, Seattle). Along the way they pass through Midwest hubs (MDW, MCI, STL), Southern cities (DAL, HOU, BNA), and Desert Southwest airports (PHX, ELP, ABQ, LAS). The pattern reflects Southwest Airlines routing an aircraft through its network in a single long operational day, effectively a transcontinental relay covering 2,000+ miles west to east.",
      "sql": "WITH itineraries AS (\n    SELECT\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS HopCount,\n        arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode)))), 1) AS StartAirport,\n        arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1) AS EndAirport,\n        arrayStringConcat(\n            arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))\n            || [arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0 AND Tail_Number != '' AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n    HAVING HopCount = 8\n),\nroute_summary AS (\n    SELECT FlightNum, Carrier, Route, StartAirport, EndAirport, count() AS DateCount\n    FROM itineraries\n    GROUP BY FlightNum, Carrier, Route, StartAirport, EndAirport\n    ORDER BY DateCount DESC\n    LIMIT 10\n)\nSELECT\n    rs.FlightNum, rs.Carrier, rs.Route, rs.DateCount,\n    a_start.DisplayAirportName AS StartName, a_start.Latitude AS StartLat, a_start.Longitude AS StartLon,\n    a_end.DisplayAirportName AS EndName, a_end.Latitude AS EndLat, a_end.Longitude AS EndLon\nFROM route_summary rs\nLEFT JOIN ontime.dim_airports a_start ON a_start.AirportCode = rs.StartAirport\nLEFT JOIN ontime.dim_airports a_end ON a_end.AirportCode = rs.EndAirport\nORDER BY rs.DateCount DESC",
      "row_count": 10,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "Route",
        "DateCount",
        "StartName",
        "StartLat",
        "StartLon",
        "EndName",
        "EndLat",
        "EndLon"
      ],
      "first_row": {
        "Carrier": "WN",
        "DateCount": 61,
        "EndLat": 39.49916667,
        "EndLon": -119.76805556,
        "EndName": "Reno/Tahoe International",
        "FlightNum": "2215",
        "Route": "CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO",
        "StartLat": 39.99694444,
        "StartLon": -82.89222222,
        "StartName": "John Glenn Columbus International"
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
      "answer_markdown": "The highest-hop example is **WN flight 366 on 2024-12-01**, operated by tail **N957WN**, covering **8 hops** (9 airports) along the route **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA**. The aircraft started at Long Island MacArthur (ISP), made stops in Baltimore (BWI), Myrtle Beach (MYR), Nashville (BNA), Fort Walton Beach (VPS), Dallas Love Field (DAL), Las Vegas (LAS), and Oakland (OAK), before finishing in Seattle (SEA). All top-10 itineraries share the same 8-hop maximum and are Southwest Airlines operations.",
      "sql": "WITH itineraries AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS HopCount,\n        arrayStringConcat(\n            arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))\n            || [arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],\n            '-'\n        ) AS Route,\n        arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode, DestCode))) AS LegSeq\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n        AND Tail_Number != ''\n        AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n),\nranked AS (\n    SELECT\n        row_number() OVER (ORDER BY HopCount DESC, FlightDate DESC) AS Rank,\n        FlightDate,\n        Tail_Number,\n        FlightNum,\n        Carrier,\n        HopCount,\n        Route,\n        LegSeq\n    FROM itineraries\n)\nSELECT *\nFROM ranked\nWHERE Rank \u003c= 10\nORDER BY Rank",
      "row_count": 10,
      "result_columns": [
        "Rank",
        "FlightDate",
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "HopCount",
        "Route",
        "LegSeq"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "HopCount": 8,
        "LegSeq": [
          {
            "": "BWI"
          },
          {
            "": "MYR"
          },
          {
            "": "BNA"
          },
          {
            "": "VPS"
          },
          {
            "": "DAL"
          },
          {
            "": "LAS"
          },
          {
            "": "OAK"
          },
          {
            "": "SEA"
          }
        ],
        "Rank": 1,
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN"
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN flight 366 on 2024-12-01** (tail N957WN), route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA. It is both Rank 1 (the lead itinerary) and the most recent date in the top-10, ahead of the next most recent (WN 3149 on 2024-02-18).",
      "sql": "WITH itineraries AS (\n    SELECT\n        FlightDate,\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS HopCount,\n        arrayStringConcat(\n            arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))\n            || [arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n        AND Tail_Number != ''\n        AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n)\nSELECT FlightDate, Tail_Number, FlightNum, Carrier, HopCount, Route\nFROM itineraries\nORDER BY HopCount DESC, FlightDate DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "FlightDate",
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "HopCount",
        "Route"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "HopCount": 8,
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN"
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top itineraries are overwhelmingly **recurring scheduled patterns**, not one-offs. The most frequent 8-hop route (WN 2215: CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO) ran **61 times** between August and October 2009. WN 1923 ran 60 times, WN 2558 ran 55 times. These routes operated nearly every day of a seasonal schedule window, indicating deliberate long-day rotations assigned to a flight number — not anomalies.",
      "sql": "WITH itineraries AS (\n    SELECT\n        FlightDate,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS HopCount,\n        arrayStringConcat(\n            arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))\n            || [arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n        AND Tail_Number != ''\n        AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n    HAVING HopCount = 8\n)\nSELECT FlightNum, Carrier, Route, count() AS DateCount,\n       min(FlightDate) AS FirstSeen, max(FlightDate) AS LastSeen\nFROM itineraries\nGROUP BY FlightNum, Carrier, Route\nORDER BY DateCount DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "Route",
        "DateCount",
        "FirstSeen",
        "LastSeen"
      ],
      "first_row": {
        "Carrier": "WN",
        "DateCount": 61,
        "FirstSeen": "2009-08-17T00:00:00Z",
        "FlightNum": "2215",
        "LastSeen": "2009-10-30T00:00:00Z",
        "Route": "CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO"
      }
    },
    {
      "id": "q4",
      "subquestion": "What geographic pattern do the top itineraries show?",
      "answer_markdown": "All 8-hop itineraries trace **long east-to-west (or southeast-to-northwest) transcontinental sweeps** across the continental United States. Routes consistently originate in the Eastern US (e.g., Columbus OH, Pittsburgh PA, Hartford CT, Fort Lauderdale FL, Richmond VA) and terminate at West Coast or Pacific Northwest airports (e.g., Reno, San Jose, Oakland, Portland, Ontario CA, Seattle). Along the way they pass through Midwest hubs (MDW, MCI, STL), Southern cities (DAL, HOU, BNA), and Desert Southwest airports (PHX, ELP, ABQ, LAS). The pattern reflects Southwest Airlines routing an aircraft through its network in a single long operational day, effectively a transcontinental relay covering 2,000+ miles west to east.",
      "sql": "WITH itineraries AS (\n    SELECT\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        count() AS HopCount,\n        arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode)))), 1) AS StartAirport,\n        arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1) AS EndAirport,\n        arrayStringConcat(\n            arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))\n            || [arrayElement(arrayMap(x -\u003e x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0 AND Tail_Number != '' AND Flight_Number_Reporting_Airline != ''\n    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline\n    HAVING HopCount = 8\n),\nroute_summary AS (\n    SELECT FlightNum, Carrier, Route, StartAirport, EndAirport, count() AS DateCount\n    FROM itineraries\n    GROUP BY FlightNum, Carrier, Route, StartAirport, EndAirport\n    ORDER BY DateCount DESC\n    LIMIT 10\n)\nSELECT\n    rs.FlightNum, rs.Carrier, rs.Route, rs.DateCount,\n    a_start.DisplayAirportName AS StartName, a_start.Latitude AS StartLat, a_start.Longitude AS StartLon,\n    a_end.DisplayAirportName AS EndName, a_end.Latitude AS EndLat, a_end.Longitude AS EndLon\nFROM route_summary rs\nLEFT JOIN ontime.dim_airports a_start ON a_start.AirportCode = rs.StartAirport\nLEFT JOIN ontime.dim_airports a_end ON a_end.AirportCode = rs.EndAirport\nORDER BY rs.DateCount DESC",
      "row_count": 10,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "Route",
        "DateCount",
        "StartName",
        "StartLat",
        "StartLon",
        "EndName",
        "EndLat",
        "EndLon"
      ],
      "first_row": {
        "Carrier": "WN",
        "DateCount": 61,
        "EndLat": 39.49916667,
        "EndLon": -119.76805556,
        "EndName": "Reno/Tahoe International",
        "FlightNum": "2215",
        "Route": "CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO",
        "StartLat": 39.99694444,
        "StartLon": -82.89222222,
        "StartName": "John Glenn Columbus International"
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