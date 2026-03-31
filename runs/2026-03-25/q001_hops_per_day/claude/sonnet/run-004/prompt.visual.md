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

- use main proof query as the primary saved SQL already provided in the prompt
- use the other section proof queries as supporting queries when they materially improve the narrative or supporting panels
- anchor the hero narrative to the top-ranked itinerary even when another itinerary is selected in the table
- show a lead-itinerary map that remains present even before airport-coordinate enrichment succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- include a narrative hero about the lead itinerary and the broader geographic pattern of the top itineraries
- label the map as airport-coordinate enrichment in the query ledger
- reuse the enrichment results for any itinerary selected from the primary result set without issuing a new per-click enrichment query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip synced to the currently selected itinerary
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- make the selected-row map behavior explicit: when the selected itinerary differs from Rank 1, the map title, plotted route, markers, bounds, and route detail panel must visibly update to that selected itinerary rather than leaving the lead route drawn
- keep the map/detail/KPI selection state separate from the anchored hero state
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- prefer the `Route` value from the primary query as the per-row itinerary representation for redraws
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis
- derive the ordered itinerary sequence for map redraws and the route detail panel by splitting `Route` on `-`

### Data Source

SQL query for primary data source:

```sql
WITH raw AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        IATA_CODE_Reporting_Airline AS Carrier,
        FlightDate,
        count() AS hop_count,
        arraySort(x -> x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs
    FROM ontime.fact_ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate
),
daily_routes AS (
    SELECT
        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.1, sorted_legs),
                [sorted_legs[length(sorted_legs)].2]
            ),
            '-'
        ) AS Route
    FROM raw
),
max_hops AS (
    SELECT max(hop_count) AS max_hop_count FROM daily_routes
),
top_route_days AS (
    SELECT * FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)
),
route_recurrence AS (
    SELECT Route, countDistinct(FlightDate) AS recurrence_count
    FROM daily_routes
    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)
    GROUP BY Route
),
recent_top AS (
    SELECT
        Route,
        argMax(Tail_Number, FlightDate) AS Tail_Number,
        argMax(FlightNum, FlightDate) AS FlightNum,
        argMax(Carrier, FlightDate) AS Carrier,
        max(FlightDate) AS most_recent_date,
        argMax(hop_count, FlightDate) AS hop_count
    FROM top_route_days
    GROUP BY Route
)
SELECT
    rt.Tail_Number,
    rt.FlightNum,
    rt.Carrier,
    rt.most_recent_date AS FlightDate,
    rt.hop_count,
    rr.recurrence_count,
    rt.Route
FROM recent_top rt
JOIN route_recurrence rr ON rt.Route = rr.Route
ORDER BY rt.most_recent_date DESC
LIMIT 10
```

Data example/snippet:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified section answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum number of hops flown by a single aircraft on a single flight number in one day is **8 hops** (9-airport itineraries). All 10 unique routes are operated by Southwest Airlines (WN). The most recent top 10 unique routes by departure date are:\n\n| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -\u003e x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT * FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n),\nrecent_top AS (\n    SELECT\n        Route,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum, FlightDate) AS FlightNum,\n        argMax(Carrier, FlightDate) AS Carrier,\n        max(FlightDate) AS most_recent_date,\n        argMax(hop_count, FlightDate) AS hop_count\n    FROM top_route_days\n    GROUP BY Route\n)\nSELECT\n    rt.Tail_Number,\n    rt.FlightNum,\n    rt.Carrier,\n    rt.most_recent_date AS FlightDate,\n    rt.hop_count,\n    rr.recurrence_count,\n    rt.Route\nFROM recent_top rt\nJOIN route_recurrence rr ON rt.Route = rr.Route\nORDER BY rt.most_recent_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "recurrence_count",
        "Route"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes show a wide spread of recurrence, from a single occurrence to nearly 50 days:\n\n| Route | Recurrence |\n|---|---|\n| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |\n| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |\n| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |\n| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |\n| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |\n| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |\n| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |\n| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |\n| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |\n| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |\n\n**Tier breakdown (totals add to 10):**\n- **One-off (1 day):** 1 route — likely an irregular or seasonal assignment.\n- **Low recurrence (2–7 days):** 4 routes — occasional scheduled rotations, not daily.\n- **Moderate recurrence (12–20 days):** 2 routes — semi-regular scheduled patterns.\n- **High recurrence (40–47 days):** 2 routes — clearly stable, recurring scheduled rotations flown across many weeks.\n\nOverall, the majority of 8-hop itineraries are recurring scheduled patterns. Only one is a true one-off; 9 of 10 routes were flown on at least 2 distinct days, and 4 of 10 were flown on 12 or more days, indicating these represent genuine repeating flight schedules rather than one-time anomalies.",
      "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -\u003e x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT Route FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n)\nSELECT Route, recurrence_count\nFROM route_recurrence\nORDER BY recurrence_count ASC",
      "row_count": 719,
      "result_columns": [
        "Route",
        "recurrence_count"
      ],
      "first_row": {
        "Route": "EWR-MDW-ABQ-DEN-SMF-LAS-TUS-OAK-SLC",
        "recurrence_count": 1
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "All airports across the 10 routes are in the **continental United States**, and every route is operated by **Southwest Airlines (WN)**. The itineraries are true transcontinental or cross-country sweeps, connecting three geographic bands:\n\n- **West (17 airports):** SEA, OAK, SJC, SMF, RNO, LAX, BUR, SAN, PSP, LAS, PHX, SLC, ABQ, ELP, COS, DEN, MAF — dominated by Pacific Coast, Mountain West, and Desert Southwest airports.\n- **Central (14 airports):** STL, MEM, ORD, MDW, MCI, ICT, LIT, BNA, DAL, HOU, MSY, JAN, PNS, VPS — the Mississippi/Gulf corridor and Midwest hubs.\n- **East (14 airports):** LGA, ISP, BWI, IAD, RIC, RDU, MYR, CMH, CLE, DTW, ATL, TPA, MCO, FLL — stretching from the Northeast down the Atlantic Seaboard to Florida.\n\nThe dominant pattern is **coast-to-coast routing**: each itinerary typically originates on one coast or the Gulf South, threads through 2–3 inland stops, and terminates on the opposite coast. Several routes end or begin at West Coast leisure-travel airports (OAK, LAX, SJC, LAS). Recurring stops like BWI, BNA, DAL, LAS, MDW, and OAK appear across multiple routes, reflecting Southwest's key connecting cities. No route is confined to a single region; all 10 span at least two of the three geographic bands.",
      "sql": "SELECT\n    multiIf(d.Longitude \u003c -100, 'West', d.Longitude \u003c -85, 'Central', 'East') AS region,\n    count() AS airport_count,\n    groupArray(a.code) AS airports\nFROM (\n    SELECT DISTINCT arrayJoin([\n        'ISP','BWI','MYR','BNA','VPS','DAL','LAS','OAK','SEA',\n        'CLE','PNS','HOU','MCI','PHX','BUR','DEN',\n        'ELP','LIT','ATL','RIC','MDW','SAN',\n        'MSY','CMH','RDU','DTW','LAX',\n        'TPA','ORD','SLC','SJC',\n        'LGA','STL','ICT','COS','JAN',\n        'SMF','PSP','RNO',\n        'ABQ','MAF',\n        'FLL','MCO','MEM','IAD'\n    ]) AS code\n) a\nLEFT JOIN ontime.dim_airports d ON a.code = d.AirportCode\nGROUP BY region\nORDER BY min(d.Longitude)",
      "row_count": 3,
      "result_columns": [
        "region",
        "airport_count",
        "airports"
      ],
      "first_row": {
        "airport_count": 17,
        "airports": [
          "LAS",
          "OAK",
          "SEA",
          "PHX",
          "BUR",
          "DEN",
          "ELP",
          "SAN",
          "LAX",
          "SLC",
          "SJC",
          "COS",
          "SMF",
          "PSP",
          "RNO",
          "ABQ",
          "MAF"
        ],
        "region": "West"
      }
    }
  ]
}

### Multi-query additions

- The saved SQL shown below is the primary section query for this page.
- The verified analysis package includes named supporting section queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- The dashboard does not need to mirror `report.md`; it should combine narrative and interactive analysis.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{
  "question_title": "Highest daily hops for one aircraft on one flight number",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified section answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum number of hops flown by a single aircraft on a single flight number in one day is **8 hops** (9-airport itineraries). All 10 unique routes are operated by Southwest Airlines (WN). The most recent top 10 unique routes by departure date are:\n\n| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -\u003e x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT * FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n),\nrecent_top AS (\n    SELECT\n        Route,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum, FlightDate) AS FlightNum,\n        argMax(Carrier, FlightDate) AS Carrier,\n        max(FlightDate) AS most_recent_date,\n        argMax(hop_count, FlightDate) AS hop_count\n    FROM top_route_days\n    GROUP BY Route\n)\nSELECT\n    rt.Tail_Number,\n    rt.FlightNum,\n    rt.Carrier,\n    rt.most_recent_date AS FlightDate,\n    rt.hop_count,\n    rr.recurrence_count,\n    rt.Route\nFROM recent_top rt\nJOIN route_recurrence rr ON rt.Route = rr.Route\nORDER BY rt.most_recent_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "recurrence_count",
        "Route"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes show a wide spread of recurrence, from a single occurrence to nearly 50 days:\n\n| Route | Recurrence |\n|---|---|\n| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |\n| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |\n| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |\n| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |\n| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |\n| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |\n| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |\n| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |\n| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |\n| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |\n\n**Tier breakdown (totals add to 10):**\n- **One-off (1 day):** 1 route — likely an irregular or seasonal assignment.\n- **Low recurrence (2–7 days):** 4 routes — occasional scheduled rotations, not daily.\n- **Moderate recurrence (12–20 days):** 2 routes — semi-regular scheduled patterns.\n- **High recurrence (40–47 days):** 2 routes — clearly stable, recurring scheduled rotations flown across many weeks.\n\nOverall, the majority of 8-hop itineraries are recurring scheduled patterns. Only one is a true one-off; 9 of 10 routes were flown on at least 2 distinct days, and 4 of 10 were flown on 12 or more days, indicating these represent genuine repeating flight schedules rather than one-time anomalies.",
      "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -\u003e x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT Route FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n)\nSELECT Route, recurrence_count\nFROM route_recurrence\nORDER BY recurrence_count ASC",
      "row_count": 719,
      "result_columns": [
        "Route",
        "recurrence_count"
      ],
      "first_row": {
        "Route": "EWR-MDW-ABQ-DEN-SMF-LAS-TUS-OAK-SLC",
        "recurrence_count": 1
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "All airports across the 10 routes are in the **continental United States**, and every route is operated by **Southwest Airlines (WN)**. The itineraries are true transcontinental or cross-country sweeps, connecting three geographic bands:\n\n- **West (17 airports):** SEA, OAK, SJC, SMF, RNO, LAX, BUR, SAN, PSP, LAS, PHX, SLC, ABQ, ELP, COS, DEN, MAF — dominated by Pacific Coast, Mountain West, and Desert Southwest airports.\n- **Central (14 airports):** STL, MEM, ORD, MDW, MCI, ICT, LIT, BNA, DAL, HOU, MSY, JAN, PNS, VPS — the Mississippi/Gulf corridor and Midwest hubs.\n- **East (14 airports):** LGA, ISP, BWI, IAD, RIC, RDU, MYR, CMH, CLE, DTW, ATL, TPA, MCO, FLL — stretching from the Northeast down the Atlantic Seaboard to Florida.\n\nThe dominant pattern is **coast-to-coast routing**: each itinerary typically originates on one coast or the Gulf South, threads through 2–3 inland stops, and terminates on the opposite coast. Several routes end or begin at West Coast leisure-travel airports (OAK, LAX, SJC, LAS). Recurring stops like BWI, BNA, DAL, LAS, MDW, and OAK appear across multiple routes, reflecting Southwest's key connecting cities. No route is confined to a single region; all 10 span at least two of the three geographic bands.",
      "sql": "SELECT\n    multiIf(d.Longitude \u003c -100, 'West', d.Longitude \u003c -85, 'Central', 'East') AS region,\n    count() AS airport_count,\n    groupArray(a.code) AS airports\nFROM (\n    SELECT DISTINCT arrayJoin([\n        'ISP','BWI','MYR','BNA','VPS','DAL','LAS','OAK','SEA',\n        'CLE','PNS','HOU','MCI','PHX','BUR','DEN',\n        'ELP','LIT','ATL','RIC','MDW','SAN',\n        'MSY','CMH','RDU','DTW','LAX',\n        'TPA','ORD','SLC','SJC',\n        'LGA','STL','ICT','COS','JAN',\n        'SMF','PSP','RNO',\n        'ABQ','MAF',\n        'FLL','MCO','MEM','IAD'\n    ]) AS code\n) a\nLEFT JOIN ontime.dim_airports d ON a.code = d.AirportCode\nGROUP BY region\nORDER BY min(d.Longitude)",
      "row_count": 3,
      "result_columns": [
        "region",
        "airport_count",
        "airports"
      ],
      "first_row": {
        "airport_count": 17,
        "airports": [
          "LAS",
          "OAK",
          "SEA",
          "PHX",
          "BUR",
          "DEN",
          "ELP",
          "SAN",
          "LAX",
          "SLC",
          "SJC",
          "COS",
          "SMF",
          "PSP",
          "RNO",
          "ABQ",
          "MAF"
        ],
        "region": "West"
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