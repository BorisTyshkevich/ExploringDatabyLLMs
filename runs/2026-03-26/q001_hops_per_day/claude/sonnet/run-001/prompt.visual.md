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
WITH
legs_deduped AS (
    SELECT
        FlightDate, Tail_Number,
        IATA_CODE_Reporting_Airline AS carrier,
        Flight_Number_Reporting_Airline AS flight_num,
        OriginCode, DestCode,
        min(ifNull(CRSDepTime, 0)) AS dep_time
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode
),
itineraries AS (
    SELECT
        FlightDate, Tail_Number, carrier, flight_num,
        count() AS hop_count,
        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route
    FROM legs_deduped
    GROUP BY FlightDate, Tail_Number, carrier, flight_num
),
top_itin AS (
    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route
    FROM itineraries
    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)
)
SELECT
    argMax(Tail_Number, FlightDate) AS aircraft_id,
    argMax(flight_num, FlightDate) AS flight_number,
    argMax(carrier, FlightDate) AS carrier,
    max(FlightDate) AS most_recent_date,
    any(hop_count) AS hop_count,
    countDistinct(FlightDate) AS recurrence_count,
    Route
FROM top_itin
GROUP BY Route
ORDER BY most_recent_date DESC
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
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum single-day hop count is **8 legs** for one aircraft under the same flight number, achieved exclusively by **Southwest Airlines (WN)**. Across all history, there are multiple distinct 8-leg route strings. The 10 most recently flown unique routes (one row per distinct route, showing the most recent occurrence) are listed below with their aircraft id, flight number, flight date, hop count, recurrence count (distinct days that exact route was ever flown), and the full textual route:\n\n| Aircraft | Flight# | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 4 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 46 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    argMax(Tail_Number, FlightDate) AS aircraft_id,\n    argMax(flight_num, FlightDate) AS flight_number,\n    argMax(carrier, FlightDate) AS carrier,\n    max(FlightDate) AS most_recent_date,\n    any(hop_count) AS hop_count,\n    countDistinct(FlightDate) AS recurrence_count,\n    Route\nFROM top_itin\nGROUP BY Route\nORDER BY most_recent_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "most_recent_date",
        "hop_count",
        "recurrence_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_number": "366",
        "hop_count": 8,
        "most_recent_date": "2024-12-01T00:00:00Z",
        "recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes show a wide spectrum from one-offs to established scheduled patterns. Recurrence counts for each route (exact days flown in full history):\n\n1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA: **1**\n2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN: **4**\n3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN: **2**\n4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX: **5**\n5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC: **40**\n6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN: **46**\n7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS: **7**\n8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK: **20**\n9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK: **12**\n10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX: **5**\n\nTiered summary (totals add to 10):\n- **One-off (1 day):** 1 route — route #1 was flown exactly once.\n- **Rare (2–7 days):** 5 routes — routes #2, #3, #4, #7, #10 were flown on only a handful of occasions.\n- **Recurring (12–20 days):** 2 routes — routes #8 and #9 appear periodically, suggesting semi-regular scheduling.\n- **Highly recurring (40–46 days):** 2 routes — routes #5 and #6 were each flown on more than 40 distinct days, indicating firmly established scheduled patterns.\n\nOverall, 6 of 10 routes are rare-to-one-off while 4 show meaningful recurrence, meaning most 8-hop itineraries are ad-hoc operational assignments, but a small subset are genuine recurring scheduled turns.",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    Route,\n    countDistinct(FlightDate) AS recurrence_count,\n    multiIf(\n        countDistinct(FlightDate) = 1, 'one-off',\n        countDistinct(FlightDate) \u003c= 7, 'rare (2-7)',\n        countDistinct(FlightDate) \u003c= 20, 'recurring (8-20)',\n        'highly recurring (40+)'\n    ) AS tier\nFROM top_itin\nGROUP BY Route\nORDER BY max(FlightDate) DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Route",
        "recurrence_count",
        "tier"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "recurrence_count": 1,
        "tier": "one-off"
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "The 10 routes collectively touch **45 unique airports across 27 US states** (plus DC). All airports are in the contiguous United States — no Hawaii, Alaska, or international airports appear, consistent with Southwest Airlines' domestic network during the covered periods.\n\nKey geographic patterns:\n\n1. **Coast-to-coast sweep**: Every route crosses multiple time zones in a single day. Routes span from the Northeast (LGA, ISP, BWI) or Southeast (MSY, TPA, ATL) all the way to the West Coast (OAK, BUR, LAX, SJC, SEA, SMF) or terminate in the Mountain West (DEN, LAS, PHX). This is characteristic of maximum aircraft utilization.\n\n2. **California dominates as a terminus**: 7 California airports appear (BUR, LAX, OAK, PSP, SAN, SJC, SMF) — more than any other state — and California is the endpoint of 7 of the 10 routes, confirming that the West Coast is the preferred turn-around anchor for these ultra-long turns.\n\n3. **Texas and the South serve as through-hubs**: Dallas Love Field (DAL), Houston (HOU), and other Texas airports appear as interior waypoints in multiple routes, reflecting Southwest's historical strength in the South-Central corridor.\n\n4. **Florida and the Southeast as origin clusters**: 5 Florida airports (FLL, MCO, PNS, TPA, VPS) and several Southeast cities (ATL, MYR, RDU, RIC, BNA, MEM) frequently appear early in routes, suggesting morning departures from the East before heading westward.\n\n5. **No large hub concentration**: The routes do not funnel through any single mega-hub. The 45 airports are spread across the network, consistent with Southwest's point-to-point model rather than hub-and-spoke.",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n),\ntop10_routes AS (\n    SELECT Route\n    FROM top_itin\n    GROUP BY Route\n    ORDER BY max(FlightDate) DESC\n    LIMIT 10\n),\nairport_list AS (\n    SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS airport_code\n    FROM top10_routes\n)\nSELECT\n    al.airport_code,\n    any(d.CityName) AS city,\n    any(d.StateName) AS state,\n    any(d.StateCode) AS state_code\nFROM airport_list al\nLEFT JOIN ontime.dim_airports d ON al.airport_code = toFixedString(d.AirportCode, 3)\nGROUP BY al.airport_code\nORDER BY any(d.StateCode), al.airport_code",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "city",
        "state",
        "state_code"
      ],
      "first_row": {
        "airport_code": "LIT",
        "city": "Little Rock, AR",
        "state": "Arkansas",
        "state_code": "AR"
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
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum single-day hop count is **8 legs** for one aircraft under the same flight number, achieved exclusively by **Southwest Airlines (WN)**. Across all history, there are multiple distinct 8-leg route strings. The 10 most recently flown unique routes (one row per distinct route, showing the most recent occurrence) are listed below with their aircraft id, flight number, flight date, hop count, recurrence count (distinct days that exact route was ever flown), and the full textual route:\n\n| Aircraft | Flight# | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 4 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 46 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    argMax(Tail_Number, FlightDate) AS aircraft_id,\n    argMax(flight_num, FlightDate) AS flight_number,\n    argMax(carrier, FlightDate) AS carrier,\n    max(FlightDate) AS most_recent_date,\n    any(hop_count) AS hop_count,\n    countDistinct(FlightDate) AS recurrence_count,\n    Route\nFROM top_itin\nGROUP BY Route\nORDER BY most_recent_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "most_recent_date",
        "hop_count",
        "recurrence_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_number": "366",
        "hop_count": 8,
        "most_recent_date": "2024-12-01T00:00:00Z",
        "recurrence_count": 1
      }
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes show a wide spectrum from one-offs to established scheduled patterns. Recurrence counts for each route (exact days flown in full history):\n\n1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA: **1**\n2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN: **4**\n3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN: **2**\n4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX: **5**\n5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC: **40**\n6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN: **46**\n7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS: **7**\n8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK: **20**\n9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK: **12**\n10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX: **5**\n\nTiered summary (totals add to 10):\n- **One-off (1 day):** 1 route — route #1 was flown exactly once.\n- **Rare (2–7 days):** 5 routes — routes #2, #3, #4, #7, #10 were flown on only a handful of occasions.\n- **Recurring (12–20 days):** 2 routes — routes #8 and #9 appear periodically, suggesting semi-regular scheduling.\n- **Highly recurring (40–46 days):** 2 routes — routes #5 and #6 were each flown on more than 40 distinct days, indicating firmly established scheduled patterns.\n\nOverall, 6 of 10 routes are rare-to-one-off while 4 show meaningful recurrence, meaning most 8-hop itineraries are ad-hoc operational assignments, but a small subset are genuine recurring scheduled turns.",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    Route,\n    countDistinct(FlightDate) AS recurrence_count,\n    multiIf(\n        countDistinct(FlightDate) = 1, 'one-off',\n        countDistinct(FlightDate) \u003c= 7, 'rare (2-7)',\n        countDistinct(FlightDate) \u003c= 20, 'recurring (8-20)',\n        'highly recurring (40+)'\n    ) AS tier\nFROM top_itin\nGROUP BY Route\nORDER BY max(FlightDate) DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Route",
        "recurrence_count",
        "tier"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "recurrence_count": 1,
        "tier": "one-off"
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "The 10 routes collectively touch **45 unique airports across 27 US states** (plus DC). All airports are in the contiguous United States — no Hawaii, Alaska, or international airports appear, consistent with Southwest Airlines' domestic network during the covered periods.\n\nKey geographic patterns:\n\n1. **Coast-to-coast sweep**: Every route crosses multiple time zones in a single day. Routes span from the Northeast (LGA, ISP, BWI) or Southeast (MSY, TPA, ATL) all the way to the West Coast (OAK, BUR, LAX, SJC, SEA, SMF) or terminate in the Mountain West (DEN, LAS, PHX). This is characteristic of maximum aircraft utilization.\n\n2. **California dominates as a terminus**: 7 California airports appear (BUR, LAX, OAK, PSP, SAN, SJC, SMF) — more than any other state — and California is the endpoint of 7 of the 10 routes, confirming that the West Coast is the preferred turn-around anchor for these ultra-long turns.\n\n3. **Texas and the South serve as through-hubs**: Dallas Love Field (DAL), Houston (HOU), and other Texas airports appear as interior waypoints in multiple routes, reflecting Southwest's historical strength in the South-Central corridor.\n\n4. **Florida and the Southeast as origin clusters**: 5 Florida airports (FLL, MCO, PNS, TPA, VPS) and several Southeast cities (ATL, MYR, RDU, RIC, BNA, MEM) frequently appear early in routes, suggesting morning departures from the East before heading westward.\n\n5. **No large hub concentration**: The routes do not funnel through any single mega-hub. The 45 airports are spread across the network, consistent with Southwest's point-to-point model rather than hub-and-spoke.",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n),\ntop10_routes AS (\n    SELECT Route\n    FROM top_itin\n    GROUP BY Route\n    ORDER BY max(FlightDate) DESC\n    LIMIT 10\n),\nairport_list AS (\n    SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS airport_code\n    FROM top10_routes\n)\nSELECT\n    al.airport_code,\n    any(d.CityName) AS city,\n    any(d.StateName) AS state,\n    any(d.StateCode) AS state_code\nFROM airport_list al\nLEFT JOIN ontime.dim_airports d ON al.airport_code = toFixedString(d.AirportCode, 3)\nGROUP BY al.airport_code\nORDER BY any(d.StateCode), al.airport_code",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "city",
        "state",
        "state_code"
      ],
      "first_row": {
        "airport_code": "LIT",
        "city": "Little Rock, AR",
        "state": "Arkansas",
        "state_code": "AR"
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