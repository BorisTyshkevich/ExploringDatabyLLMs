- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save.
- Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully.
- Do not write unchecked SQL.

You are reviewing a qforge analysis run after SQL generation, harness execution, and report rendering.

Your job is to judge whether the analysis artifacts actually answer the question correctly and whether the SQL evidence supports the written claims.

Return the final review by writing `review.md` in the run directory.

`review.md` must use exactly this top-level structure:

```md
# Analysis Review
Verdict: PASS

## Summary
...

## Findings
...

## Suggested Prompt Fixes
...
```

Rules:

- Set `Verdict: PASS` only when the analysis is materially aligned with the question.
- Set `Verdict: WARN` when the analysis is materially aligned overall but has limited evidence-support, wording, or minor artifact issues that should not block downstream use.
- Set `Verdict: FAIL` when any substantive correctness, grain, metric, or evidence-support problem exists that makes the run unreliable.
- Base your judgment only on the provided question prompt and the verified run artifacts.
- When artifact file paths are listed, read the files directly from the run directory instead of relying only on prompt excerpts.
- Check whether all required `### main` / `### qN` sections were answered directly.
- Check whether SQL grain and returned metrics match the prompt.
- Check whether any explicitly requested returned fields are actually present in the executed result, and only treat missing identifiers as defects when the prompt clearly requires them to be non-empty.
- Check whether prose claims are supported by the executed query results.
- Check for duplicated entities, inconsistent counts, missing requested metrics, unsupported inference, and proof-query/result mismatches.
- Treat `report.md` as a monitoring artifact. Do not treat its one-row example table or abbreviated preview as a defect when the underlying query result preserves the required ranked or detailed rows.
- In `## Findings`, cite concrete artifact names such as `queries/main.sql`, `queries/q1.sql`, `results/q1.json`, or `report.md`.
- In `## Suggested Prompt Fixes`, propose prompt-level changes only when they would reduce the observed failure mode.
- When proposing a prompt fix, be concrete. Prefer 1-3 short replacement or insertion snippets that could be added to the question prompt, not vague advice.
- Tie each prompt fix to the specific failure you found. Explain what ambiguity in the current prompt allowed the bad result and how the new wording would close that gap.
- If the prompt is already sufficiently specific and the failure is not plausibly prompt-driven, say `None.` instead of inventing weak prompt edits.
- Do not suggest code changes to qforge in this review. Focus on the run and the prompt.
- Write Markdown only to `review.md`.

Question-specific guidance:

### main
Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.
Count hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

### q1
Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

### q2
What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

Generated report.md:

```md
# Highest daily hops for one aircraft on one flight number

> Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.
Count hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

The maximum single-day hop count is **8 legs** for one aircraft under the same flight number, achieved exclusively by **Southwest Airlines (WN)**. Across all history, there are multiple distinct 8-leg route strings. The 10 most recently flown unique routes (one row per distinct route, showing the most recent occurrence) are listed below with their aircraft id, flight number, flight date, hop count, recurrence count (distinct days that exact route was ever flown), and the full textual route:

| Aircraft | Flight# | Carrier | Date | Hops | Recurrence | Route |
|---|---|---|---|---|---|---|
| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |
| N7835A | 3149 | WN | 2024-02-18 | 8 | 4 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |
| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |
| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |
| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |
| N7713A | 2884 | WN | 2022-08-31 | 8 | 46 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |
| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |
| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |
| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |
| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, most_recent_date, hop_count, recurrence_count, Route

| aircraft_id | flight_number | carrier | most_recent_date | hop_count | recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

The 10 routes show a wide spectrum from one-offs to established scheduled patterns. Recurrence counts for each route (exact days flown in full history):

1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA: **1**
2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN: **4**
3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN: **2**
4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX: **5**
5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC: **40**
6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN: **46**
7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS: **7**
8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK: **20**
9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK: **12**
10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX: **5**

Tiered summary (totals add to 10):
- **One-off (1 day):** 1 route — route #1 was flown exactly once.
- **Rare (2–7 days):** 5 routes — routes #2, #3, #4, #7, #10 were flown on only a handful of occasions.
- **Recurring (12–20 days):** 2 routes — routes #8 and #9 appear periodically, suggesting semi-regular scheduling.
- **Highly recurring (40–46 days):** 2 routes — routes #5 and #6 were each flown on more than 40 distinct days, indicating firmly established scheduled patterns.

Overall, 6 of 10 routes are rare-to-one-off while 4 show meaningful recurrence, meaning most 8-hop itineraries are ad-hoc operational assignments, but a small subset are genuine recurring scheduled turns.

- Rows returned: 10
- Columns: Route, recurrence_count, tier

| Route | recurrence_count | tier |
| --- | --- | --- |
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 | one-off |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

The 10 routes collectively touch **45 unique airports across 27 US states** (plus DC). All airports are in the contiguous United States — no Hawaii, Alaska, or international airports appear, consistent with Southwest Airlines' domestic network during the covered periods.

Key geographic patterns:

1. **Coast-to-coast sweep**: Every route crosses multiple time zones in a single day. Routes span from the Northeast (LGA, ISP, BWI) or Southeast (MSY, TPA, ATL) all the way to the West Coast (OAK, BUR, LAX, SJC, SEA, SMF) or terminate in the Mountain West (DEN, LAS, PHX). This is characteristic of maximum aircraft utilization.

2. **California dominates as a terminus**: 7 California airports appear (BUR, LAX, OAK, PSP, SAN, SJC, SMF) — more than any other state — and California is the endpoint of 7 of the 10 routes, confirming that the West Coast is the preferred turn-around anchor for these ultra-long turns.

3. **Texas and the South serve as through-hubs**: Dallas Love Field (DAL), Houston (HOU), and other Texas airports appear as interior waypoints in multiple routes, reflecting Southwest's historical strength in the South-Central corridor.

4. **Florida and the Southeast as origin clusters**: 5 Florida airports (FLL, MCO, PNS, TPA, VPS) and several Southeast cities (ATL, MYR, RDU, RIC, BNA, MEM) frequently appear early in routes, suggesting morning departures from the East before heading westward.

5. **No large hub concentration**: The routes do not funnel through any single mega-hub. The 45 airports are spread across the network, consistent with Southwest's point-to-point model rather than hub-and-spoke.

- Rows returned: 45
- Columns: airport_code, city, state, state_code

| airport_code | city | state | state_code |
| --- | --- | --- | --- |
| LIT | Little Rock, AR | Arkansas | AR |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "id": "main",
      "answer_markdown": "The maximum single-day hop count is **8 legs** for one aircraft under the same flight number, achieved exclusively by **Southwest Airlines (WN)**. Across all history, there are multiple distinct 8-leg route strings. The 10 most recently flown unique routes (one row per distinct route, showing the most recent occurrence) are listed below with their aircraft id, flight number, flight date, hop count, recurrence count (distinct days that exact route was ever flown), and the full textual route:\n\n| Aircraft | Flight# | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 4 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 46 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    argMax(Tail_Number, FlightDate) AS aircraft_id,\n    argMax(flight_num, FlightDate) AS flight_number,\n    argMax(carrier, FlightDate) AS carrier,\n    max(FlightDate) AS most_recent_date,\n    any(hop_count) AS hop_count,\n    countDistinct(FlightDate) AS recurrence_count,\n    Route\nFROM top_itin\nGROUP BY Route\nORDER BY most_recent_date DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "answer_markdown": "The 10 routes show a wide spectrum from one-offs to established scheduled patterns. Recurrence counts for each route (exact days flown in full history):\n\n1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA: **1**\n2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN: **4**\n3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN: **2**\n4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX: **5**\n5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC: **40**\n6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN: **46**\n7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS: **7**\n8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK: **20**\n9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK: **12**\n10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX: **5**\n\nTiered summary (totals add to 10):\n- **One-off (1 day):** 1 route — route #1 was flown exactly once.\n- **Rare (2–7 days):** 5 routes — routes #2, #3, #4, #7, #10 were flown on only a handful of occasions.\n- **Recurring (12–20 days):** 2 routes — routes #8 and #9 appear periodically, suggesting semi-regular scheduling.\n- **Highly recurring (40–46 days):** 2 routes — routes #5 and #6 were each flown on more than 40 distinct days, indicating firmly established scheduled patterns.\n\nOverall, 6 of 10 routes are rare-to-one-off while 4 show meaningful recurrence, meaning most 8-hop itineraries are ad-hoc operational assignments, but a small subset are genuine recurring scheduled turns.",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    Route,\n    countDistinct(FlightDate) AS recurrence_count,\n    multiIf(\n        countDistinct(FlightDate) = 1, 'one-off',\n        countDistinct(FlightDate) <= 7, 'rare (2-7)',\n        countDistinct(FlightDate) <= 20, 'recurring (8-20)',\n        'highly recurring (40+)'\n    ) AS tier\nFROM top_itin\nGROUP BY Route\nORDER BY max(FlightDate) DESC\nLIMIT 10"
    },
    {
      "id": "q2",
      "answer_markdown": "The 10 routes collectively touch **45 unique airports across 27 US states** (plus DC). All airports are in the contiguous United States — no Hawaii, Alaska, or international airports appear, consistent with Southwest Airlines' domestic network during the covered periods.\n\nKey geographic patterns:\n\n1. **Coast-to-coast sweep**: Every route crosses multiple time zones in a single day. Routes span from the Northeast (LGA, ISP, BWI) or Southeast (MSY, TPA, ATL) all the way to the West Coast (OAK, BUR, LAX, SJC, SEA, SMF) or terminate in the Mountain West (DEN, LAS, PHX). This is characteristic of maximum aircraft utilization.\n\n2. **California dominates as a terminus**: 7 California airports appear (BUR, LAX, OAK, PSP, SAN, SJC, SMF) — more than any other state — and California is the endpoint of 7 of the 10 routes, confirming that the West Coast is the preferred turn-around anchor for these ultra-long turns.\n\n3. **Texas and the South serve as through-hubs**: Dallas Love Field (DAL), Houston (HOU), and other Texas airports appear as interior waypoints in multiple routes, reflecting Southwest's historical strength in the South-Central corridor.\n\n4. **Florida and the Southeast as origin clusters**: 5 Florida airports (FLL, MCO, PNS, TPA, VPS) and several Southeast cities (ATL, MYR, RDU, RIC, BNA, MEM) frequently appear early in routes, suggesting morning departures from the East before heading westward.\n\n5. **No large hub concentration**: The routes do not funnel through any single mega-hub. The 45 airports are spread across the network, consistent with Southwest's point-to-point model rather than hub-and-spoke.",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n),\ntop10_routes AS (\n    SELECT Route\n    FROM top_itin\n    GROUP BY Route\n    ORDER BY max(FlightDate) DESC\n    LIMIT 10\n),\nairport_list AS (\n    SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS airport_code\n    FROM top10_routes\n)\nSELECT\n    al.airport_code,\n    any(d.CityName) AS city,\n    any(d.StateName) AS state,\n    any(d.StateCode) AS state_code\nFROM airport_list al\nLEFT JOIN ontime.dim_airports d ON al.airport_code = toFixedString(d.AirportCode, 3)\nGROUP BY al.airport_code\nORDER BY any(d.StateCode), al.airport_code"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    argMax(Tail_Number, FlightDate) AS aircraft_id,\n    argMax(flight_num, FlightDate) AS flight_number,\n    argMax(carrier, FlightDate) AS carrier,\n    max(FlightDate) AS most_recent_date,\n    any(hop_count) AS hop_count,\n    countDistinct(FlightDate) AS recurrence_count,\n    Route\nFROM top_itin\nGROUP BY Route\nORDER BY most_recent_date DESC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum single-day hop count is **8 legs** for one aircraft under the same flight number, achieved exclusively by **Southwest Airlines (WN)**. Across all history, there are multiple distinct 8-leg route strings. The 10 most recently flown unique routes (one row per distinct route, showing the most recent occurrence) are listed below with their aircraft id, flight number, flight date, hop count, recurrence count (distinct days that exact route was ever flown), and the full textual route:\n\n| Aircraft | Flight# | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 4 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 46 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    argMax(Tail_Number, FlightDate) AS aircraft_id,\n    argMax(flight_num, FlightDate) AS flight_number,\n    argMax(carrier, FlightDate) AS carrier,\n    max(FlightDate) AS most_recent_date,\n    any(hop_count) AS hop_count,\n    countDistinct(FlightDate) AS recurrence_count,\n    Route\nFROM top_itin\nGROUP BY Route\nORDER BY most_recent_date DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes show a wide spectrum from one-offs to established scheduled patterns. Recurrence counts for each route (exact days flown in full history):\n\n1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA: **1**\n2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN: **4**\n3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN: **2**\n4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX: **5**\n5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC: **40**\n6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN: **46**\n7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS: **7**\n8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK: **20**\n9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK: **12**\n10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX: **5**\n\nTiered summary (totals add to 10):\n- **One-off (1 day):** 1 route — route #1 was flown exactly once.\n- **Rare (2–7 days):** 5 routes — routes #2, #3, #4, #7, #10 were flown on only a handful of occasions.\n- **Recurring (12–20 days):** 2 routes — routes #8 and #9 appear periodically, suggesting semi-regular scheduling.\n- **Highly recurring (40–46 days):** 2 routes — routes #5 and #6 were each flown on more than 40 distinct days, indicating firmly established scheduled patterns.\n\nOverall, 6 of 10 routes are rare-to-one-off while 4 show meaningful recurrence, meaning most 8-hop itineraries are ad-hoc operational assignments, but a small subset are genuine recurring scheduled turns.",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT\n        FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n)\nSELECT\n    Route,\n    countDistinct(FlightDate) AS recurrence_count,\n    multiIf(\n        countDistinct(FlightDate) = 1, 'one-off',\n        countDistinct(FlightDate) \u003c= 7, 'rare (2-7)',\n        countDistinct(FlightDate) \u003c= 20, 'recurring (8-20)',\n        'highly recurring (40+)'\n    ) AS tier\nFROM top_itin\nGROUP BY Route\nORDER BY max(FlightDate) DESC\nLIMIT 10"
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "The 10 routes collectively touch **45 unique airports across 27 US states** (plus DC). All airports are in the contiguous United States — no Hawaii, Alaska, or international airports appear, consistent with Southwest Airlines' domestic network during the covered periods.\n\nKey geographic patterns:\n\n1. **Coast-to-coast sweep**: Every route crosses multiple time zones in a single day. Routes span from the Northeast (LGA, ISP, BWI) or Southeast (MSY, TPA, ATL) all the way to the West Coast (OAK, BUR, LAX, SJC, SEA, SMF) or terminate in the Mountain West (DEN, LAS, PHX). This is characteristic of maximum aircraft utilization.\n\n2. **California dominates as a terminus**: 7 California airports appear (BUR, LAX, OAK, PSP, SAN, SJC, SMF) — more than any other state — and California is the endpoint of 7 of the 10 routes, confirming that the West Coast is the preferred turn-around anchor for these ultra-long turns.\n\n3. **Texas and the South serve as through-hubs**: Dallas Love Field (DAL), Houston (HOU), and other Texas airports appear as interior waypoints in multiple routes, reflecting Southwest's historical strength in the South-Central corridor.\n\n4. **Florida and the Southeast as origin clusters**: 5 Florida airports (FLL, MCO, PNS, TPA, VPS) and several Southeast cities (ATL, MYR, RDU, RIC, BNA, MEM) frequently appear early in routes, suggesting morning departures from the East before heading westward.\n\n5. **No large hub concentration**: The routes do not funnel through any single mega-hub. The 45 airports are spread across the network, consistent with Southwest's point-to-point model rather than hub-and-spoke.",
      "sql": "WITH\nlegs_deduped AS (\n    SELECT\n        FlightDate, Tail_Number,\n        IATA_CODE_Reporting_Airline AS carrier,\n        Flight_Number_Reporting_Airline AS flight_num,\n        OriginCode, DestCode,\n        min(ifNull(CRSDepTime, 0)) AS dep_time\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode\n),\nitineraries AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num,\n        count() AS hop_count,\n        arrayStringConcat(arrayConcat(arrayMap(x -\u003e x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route\n    FROM legs_deduped\n    GROUP BY FlightDate, Tail_Number, carrier, flight_num\n),\ntop_itin AS (\n    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route\n    FROM itineraries\n    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)\n),\ntop10_routes AS (\n    SELECT Route\n    FROM top_itin\n    GROUP BY Route\n    ORDER BY max(FlightDate) DESC\n    LIMIT 10\n),\nairport_list AS (\n    SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS airport_code\n    FROM top10_routes\n)\nSELECT\n    al.airport_code,\n    any(d.CityName) AS city,\n    any(d.StateName) AS state,\n    any(d.StateCode) AS state_code\nFROM airport_list al\nLEFT JOIN ontime.dim_airports d ON al.airport_code = toFixedString(d.AirportCode, 3)\nGROUP BY al.airport_code\nORDER BY any(d.StateCode), al.airport_code"
    }
  ]
}
```

Saved visual_input.json:

```json
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
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `queries/main.sql`
- `queries/q1.sql`
- `queries/q2.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/main.json`
- `results/q1.json`
- `results/q2.json`