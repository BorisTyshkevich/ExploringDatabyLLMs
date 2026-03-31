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
WITH itineraries AS (
  SELECT
    Tail_Number,
    IATA_CODE_Reporting_Airline AS Carrier,
    Flight_Number_Reporting_Airline AS FlightNum,
    FlightDate,
    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,
    arrayStringConcat(
      arrayConcat(
        arrayMap(t -> t.2, arraySort(t -> t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),
        [argMax(DestCode, coalesce(CRSDepTime, 0))]
      ),
      '-'
    ) AS Route
  FROM ontime.fact_ontime
  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate
),
max_hop AS (
  SELECT max(hop_count) AS max_hop_count FROM itineraries
),
route_recurrence AS (
  SELECT Route, countDistinct(FlightDate) AS recurrence_count
  FROM itineraries GROUP BY Route
),
top_itineraries AS (
  SELECT i.Tail_Number, i.Carrier, i.FlightNum, i.FlightDate, i.hop_count, i.Route
  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count
),
top_with_recurrence AS (
  SELECT t.Tail_Number, t.Carrier, t.FlightNum, t.FlightDate, t.hop_count, r.recurrence_count, t.Route
  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route
)
SELECT
  argMax(Tail_Number, FlightDate) AS aircraft_id,
  argMax(FlightNum, FlightDate) AS flight_number,
  argMax(Carrier, FlightDate) AS carrier,
  max(FlightDate) AS flight_date,
  max(hop_count) AS hop_count,
  any(recurrence_count) AS recurrence_count,
  Route
FROM top_with_recurrence
GROUP BY Route
ORDER BY flight_date DESC
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
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum distinct-hop count across all history is **8 hops** in a single day under one flight number. All 10 unique routes at that peak are operated by **Southwest Airlines (WN)**, reflecting Southwest's point-to-point milk-run scheduling model. The most recent unique 8-hop route was flown by tail **N957WN** on flight **WN 366** on **2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA, recurrence: 1 day). The highest-recurrence route among the top 10 is **LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN**, which appeared on 47 distinct days.\n\n| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (\n  SELECT max(hop_count) AS max_hop_count FROM itineraries\n),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Tail_Number, i.Carrier, i.FlightNum, i.FlightDate, i.hop_count, i.Route\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop_with_recurrence AS (\n  SELECT t.Tail_Number, t.Carrier, t.FlightNum, t.FlightDate, t.hop_count, r.recurrence_count, t.Route\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n)\nSELECT\n  argMax(Tail_Number, FlightDate) AS aircraft_id,\n  argMax(FlightNum, FlightDate) AS flight_number,\n  argMax(Carrier, FlightDate) AS carrier,\n  max(FlightDate) AS flight_date,\n  max(hop_count) AS hop_count,\n  any(recurrence_count) AS recurrence_count,\n  Route\nFROM top_with_recurrence\nGROUP BY Route\nORDER BY flight_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "recurrence_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8,
        "recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes span a wide recurrence spectrum — from true one-offs to well-established scheduled patterns.\n\nRecurrence count for each of the 10 routes (distinct calendar days flown, all history):\n\n| Route | Recurrence |\n|---|---|\n| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |\n| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |\n| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |\n| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |\n| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |\n| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |\n| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |\n| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |\n| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |\n| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |\n\n**Tiers:** Two routes are effectively one-offs (recurrence ≤ 2 days): ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA (1) and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN (2) — total **2 routes**. Four routes show modest scheduling (5–7 days): CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS, BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — total **4 routes**. Four routes are strongly recurring scheduled patterns (12–47 days): HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — total **4 routes**. Overall, 8 of 10 routes appeared on more than 2 days, indicating that 8-hop itineraries are predominantly recurring scheduled rotations rather than one-off operations.",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Route, i.FlightDate, i.hop_count\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop10_routes AS (\n  SELECT t.Route, r.recurrence_count\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n  GROUP BY t.Route, r.recurrence_count\n  ORDER BY max(t.FlightDate) DESC\n  LIMIT 10\n)\nSELECT Route, recurrence_count\nFROM top10_routes\nORDER BY recurrence_count DESC",
      "row_count": 10,
      "result_columns": [
        "Route",
        "recurrence_count"
      ],
      "first_row": {
        "Route": "LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN",
        "recurrence_count": 47
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "All 10 routes are exclusively **Southwest Airlines (WN)** itineraries and touch **45 distinct airports** across the contiguous United States, with no international, Hawaiian, or Alaskan stops.\n\nThe routes form continuous coast-to-coast or diagonal sweeps — classic Southwest milk runs. Each itinerary anchors at one of three regional clusters:\n\n- **East/Gulf Coast origins:** ISP (Long Island), LGA (New York), CLE (Cleveland), BWI (Baltimore), FLL (Fort Lauderdale), HOU (Houston Hobby), MSY (New Orleans)\n- **Sun Belt / Midwest bridges:** ATL, BNA, DAL, MEM, MDW, STL, DEN — high-frequency Southwest connecting points mid-country\n- **West Coast termini:** OAK, LAX, BUR, SJC, SEA, SMF, LAS — all Bay Area, Southern California, or Pacific Northwest endpoints\n\nA Southwest-specific fingerprint is evident throughout: Dallas Love Field (DAL) appears in place of DFW, Houston Hobby (HOU) in place of IAH, and Chicago Midway (MDW) in place of O'Hare (ORD) — the three airports where Southwest maintains dominance over legacy hubs. The Mountain West (PHX, DEN, LAS, SLC, COS, ABQ) consistently appears as the penultimate cluster bridging Midwest and Pacific stops. No route doubles back geographically; all flow directionally from one region to another, consistent with aircraft repositioning across Southwest's national network.",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),\ntop_itineraries AS (\n  SELECT i.Route, i.FlightDate, i.hop_count\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop10_routes AS (\n  SELECT Route\n  FROM top_itineraries\n  GROUP BY Route\n  ORDER BY max(FlightDate) DESC\n  LIMIT 10\n),\nall_airports AS (\n  SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS ap\n  FROM top10_routes\n)\nSELECT\n  a.ap AS airport_code,\n  d.DisplayAirportName AS airport_name,\n  d.Latitude AS lat,\n  d.Longitude AS lon\nFROM all_airports a\nLEFT JOIN ontime.dim_airports d ON a.ap = d.AirportCode\nORDER BY airport_code",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "airport_name",
        "lat",
        "lon"
      ],
      "first_row": {
        "airport_code": "ABQ",
        "airport_name": "Albuquerque International Sunport",
        "lat": 35.03805556,
        "lon": -106.61
      }
    }
  ]
}

### Multi-query additions

- The saved SQL shown below is the primary section query for this page.
- The verified analysis package includes named supporting section queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- Do not assume auxiliary lookup or enrichment schema details from memory. Use only columns you have checked against the live endpoint or semantic-layer guidance.
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
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum distinct-hop count across all history is **8 hops** in a single day under one flight number. All 10 unique routes at that peak are operated by **Southwest Airlines (WN)**, reflecting Southwest's point-to-point milk-run scheduling model. The most recent unique 8-hop route was flown by tail **N957WN** on flight **WN 366** on **2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA, recurrence: 1 day). The highest-recurrence route among the top 10 is **LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN**, which appeared on 47 distinct days.\n\n| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (\n  SELECT max(hop_count) AS max_hop_count FROM itineraries\n),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Tail_Number, i.Carrier, i.FlightNum, i.FlightDate, i.hop_count, i.Route\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop_with_recurrence AS (\n  SELECT t.Tail_Number, t.Carrier, t.FlightNum, t.FlightDate, t.hop_count, r.recurrence_count, t.Route\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n)\nSELECT\n  argMax(Tail_Number, FlightDate) AS aircraft_id,\n  argMax(FlightNum, FlightDate) AS flight_number,\n  argMax(Carrier, FlightDate) AS carrier,\n  max(FlightDate) AS flight_date,\n  max(hop_count) AS hop_count,\n  any(recurrence_count) AS recurrence_count,\n  Route\nFROM top_with_recurrence\nGROUP BY Route\nORDER BY flight_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "recurrence_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8,
        "recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes span a wide recurrence spectrum — from true one-offs to well-established scheduled patterns.\n\nRecurrence count for each of the 10 routes (distinct calendar days flown, all history):\n\n| Route | Recurrence |\n|---|---|\n| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |\n| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |\n| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |\n| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |\n| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |\n| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |\n| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |\n| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |\n| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |\n| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |\n\n**Tiers:** Two routes are effectively one-offs (recurrence ≤ 2 days): ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA (1) and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN (2) — total **2 routes**. Four routes show modest scheduling (5–7 days): CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS, BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — total **4 routes**. Four routes are strongly recurring scheduled patterns (12–47 days): HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — total **4 routes**. Overall, 8 of 10 routes appeared on more than 2 days, indicating that 8-hop itineraries are predominantly recurring scheduled rotations rather than one-off operations.",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Route, i.FlightDate, i.hop_count\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop10_routes AS (\n  SELECT t.Route, r.recurrence_count\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n  GROUP BY t.Route, r.recurrence_count\n  ORDER BY max(t.FlightDate) DESC\n  LIMIT 10\n)\nSELECT Route, recurrence_count\nFROM top10_routes\nORDER BY recurrence_count DESC",
      "row_count": 10,
      "result_columns": [
        "Route",
        "recurrence_count"
      ],
      "first_row": {
        "Route": "LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN",
        "recurrence_count": 47
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "All 10 routes are exclusively **Southwest Airlines (WN)** itineraries and touch **45 distinct airports** across the contiguous United States, with no international, Hawaiian, or Alaskan stops.\n\nThe routes form continuous coast-to-coast or diagonal sweeps — classic Southwest milk runs. Each itinerary anchors at one of three regional clusters:\n\n- **East/Gulf Coast origins:** ISP (Long Island), LGA (New York), CLE (Cleveland), BWI (Baltimore), FLL (Fort Lauderdale), HOU (Houston Hobby), MSY (New Orleans)\n- **Sun Belt / Midwest bridges:** ATL, BNA, DAL, MEM, MDW, STL, DEN — high-frequency Southwest connecting points mid-country\n- **West Coast termini:** OAK, LAX, BUR, SJC, SEA, SMF, LAS — all Bay Area, Southern California, or Pacific Northwest endpoints\n\nA Southwest-specific fingerprint is evident throughout: Dallas Love Field (DAL) appears in place of DFW, Houston Hobby (HOU) in place of IAH, and Chicago Midway (MDW) in place of O'Hare (ORD) — the three airports where Southwest maintains dominance over legacy hubs. The Mountain West (PHX, DEN, LAS, SLC, COS, ABQ) consistently appears as the penultimate cluster bridging Midwest and Pacific stops. No route doubles back geographically; all flow directionally from one region to another, consistent with aircraft repositioning across Southwest's national network.",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),\ntop_itineraries AS (\n  SELECT i.Route, i.FlightDate, i.hop_count\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop10_routes AS (\n  SELECT Route\n  FROM top_itineraries\n  GROUP BY Route\n  ORDER BY max(FlightDate) DESC\n  LIMIT 10\n),\nall_airports AS (\n  SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS ap\n  FROM top10_routes\n)\nSELECT\n  a.ap AS airport_code,\n  d.DisplayAirportName AS airport_name,\n  d.Latitude AS lat,\n  d.Longitude AS lon\nFROM all_airports a\nLEFT JOIN ontime.dim_airports d ON a.ap = d.AirportCode\nORDER BY airport_code",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "airport_name",
        "lat",
        "lon"
      ],
      "first_row": {
        "airport_code": "ABQ",
        "airport_name": "Albuquerque International Sunport",
        "lat": 35.03805556,
        "lon": -106.61
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Before writing `visual.html`, self-verify every browser-side SQL statement you intend to ship, including primary, supporting, enrichment, drill-down, and lookup queries.
- For each query, run a cheap live check against the real endpoint and schema first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both when that preserves the query shape.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, join, or unknown-column errors in a loop until every shipped browser query runs successfully.

Create browser-ready HTML `visual.html`.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.