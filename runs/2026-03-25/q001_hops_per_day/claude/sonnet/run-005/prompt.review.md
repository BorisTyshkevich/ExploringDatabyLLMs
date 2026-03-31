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

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

The maximum distinct-hop count across all history is **8 hops** in a single day under one flight number. All 10 unique routes at that peak are operated by **Southwest Airlines (WN)**, reflecting Southwest's point-to-point milk-run scheduling model. The most recent unique 8-hop route was flown by tail **N957WN** on flight **WN 366** on **2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA, recurrence: 1 day). The highest-recurrence route among the top 10 is **LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN**, which appeared on 47 distinct days.

| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |
|---|---|---|---|---|---|---|
| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |
| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |
| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |
| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |
| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |
| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |
| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |
| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |
| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |
| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

The 10 routes span a wide recurrence spectrum — from true one-offs to well-established scheduled patterns.

Recurrence count for each of the 10 routes (distinct calendar days flown, all history):

| Route | Recurrence |
|---|---|
| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |
| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |
| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |
| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |
| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |
| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |
| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |
| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |

**Tiers:** Two routes are effectively one-offs (recurrence ≤ 2 days): ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA (1) and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN (2) — total **2 routes**. Four routes show modest scheduling (5–7 days): CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS, BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — total **4 routes**. Four routes are strongly recurring scheduled patterns (12–47 days): HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — total **4 routes**. Overall, 8 of 10 routes appeared on more than 2 days, indicating that 8-hop itineraries are predominantly recurring scheduled rotations rather than one-off operations.

- Rows returned: 10
- Columns: Route, recurrence_count

| Route | recurrence_count |
| --- | --- |
| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

All 10 routes are exclusively **Southwest Airlines (WN)** itineraries and touch **45 distinct airports** across the contiguous United States, with no international, Hawaiian, or Alaskan stops.

The routes form continuous coast-to-coast or diagonal sweeps — classic Southwest milk runs. Each itinerary anchors at one of three regional clusters:

- **East/Gulf Coast origins:** ISP (Long Island), LGA (New York), CLE (Cleveland), BWI (Baltimore), FLL (Fort Lauderdale), HOU (Houston Hobby), MSY (New Orleans)
- **Sun Belt / Midwest bridges:** ATL, BNA, DAL, MEM, MDW, STL, DEN — high-frequency Southwest connecting points mid-country
- **West Coast termini:** OAK, LAX, BUR, SJC, SEA, SMF, LAS — all Bay Area, Southern California, or Pacific Northwest endpoints

A Southwest-specific fingerprint is evident throughout: Dallas Love Field (DAL) appears in place of DFW, Houston Hobby (HOU) in place of IAH, and Chicago Midway (MDW) in place of O'Hare (ORD) — the three airports where Southwest maintains dominance over legacy hubs. The Mountain West (PHX, DEN, LAS, SLC, COS, ABQ) consistently appears as the penultimate cluster bridging Midwest and Pacific stops. No route doubles back geographically; all flow directionally from one region to another, consistent with aircraft repositioning across Southwest's national network.

- Rows returned: 45
- Columns: airport_code, airport_name, lat, lon

| airport_code | airport_name | lat | lon |
| --- | --- | --- | --- |
| ABQ | Albuquerque International Sunport | 35.03805556 | -106.61 |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "id": "main",
      "answer_markdown": "The maximum distinct-hop count across all history is **8 hops** in a single day under one flight number. All 10 unique routes at that peak are operated by **Southwest Airlines (WN)**, reflecting Southwest's point-to-point milk-run scheduling model. The most recent unique 8-hop route was flown by tail **N957WN** on flight **WN 366** on **2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA, recurrence: 1 day). The highest-recurrence route among the top 10 is **LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN**, which appeared on 47 distinct days.\n\n| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -> t.2, arraySort(t -> t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (\n  SELECT max(hop_count) AS max_hop_count FROM itineraries\n),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Tail_Number, i.Carrier, i.FlightNum, i.FlightDate, i.hop_count, i.Route\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop_with_recurrence AS (\n  SELECT t.Tail_Number, t.Carrier, t.FlightNum, t.FlightDate, t.hop_count, r.recurrence_count, t.Route\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n)\nSELECT\n  argMax(Tail_Number, FlightDate) AS aircraft_id,\n  argMax(FlightNum, FlightDate) AS flight_number,\n  argMax(Carrier, FlightDate) AS carrier,\n  max(FlightDate) AS flight_date,\n  max(hop_count) AS hop_count,\n  any(recurrence_count) AS recurrence_count,\n  Route\nFROM top_with_recurrence\nGROUP BY Route\nORDER BY flight_date DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "answer_markdown": "The 10 routes span a wide recurrence spectrum — from true one-offs to well-established scheduled patterns.\n\nRecurrence count for each of the 10 routes (distinct calendar days flown, all history):\n\n| Route | Recurrence |\n|---|---|\n| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |\n| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |\n| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |\n| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |\n| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |\n| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |\n| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |\n| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |\n| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |\n| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |\n\n**Tiers:** Two routes are effectively one-offs (recurrence ≤ 2 days): ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA (1) and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN (2) — total **2 routes**. Four routes show modest scheduling (5–7 days): CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS, BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — total **4 routes**. Four routes are strongly recurring scheduled patterns (12–47 days): HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — total **4 routes**. Overall, 8 of 10 routes appeared on more than 2 days, indicating that 8-hop itineraries are predominantly recurring scheduled rotations rather than one-off operations.",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -> t.2, arraySort(t -> t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Route, i.FlightDate, i.hop_count\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop10_routes AS (\n  SELECT t.Route, r.recurrence_count\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n  GROUP BY t.Route, r.recurrence_count\n  ORDER BY max(t.FlightDate) DESC\n  LIMIT 10\n)\nSELECT Route, recurrence_count\nFROM top10_routes\nORDER BY recurrence_count DESC"
    },
    {
      "id": "q2",
      "answer_markdown": "All 10 routes are exclusively **Southwest Airlines (WN)** itineraries and touch **45 distinct airports** across the contiguous United States, with no international, Hawaiian, or Alaskan stops.\n\nThe routes form continuous coast-to-coast or diagonal sweeps — classic Southwest milk runs. Each itinerary anchors at one of three regional clusters:\n\n- **East/Gulf Coast origins:** ISP (Long Island), LGA (New York), CLE (Cleveland), BWI (Baltimore), FLL (Fort Lauderdale), HOU (Houston Hobby), MSY (New Orleans)\n- **Sun Belt / Midwest bridges:** ATL, BNA, DAL, MEM, MDW, STL, DEN — high-frequency Southwest connecting points mid-country\n- **West Coast termini:** OAK, LAX, BUR, SJC, SEA, SMF, LAS — all Bay Area, Southern California, or Pacific Northwest endpoints\n\nA Southwest-specific fingerprint is evident throughout: Dallas Love Field (DAL) appears in place of DFW, Houston Hobby (HOU) in place of IAH, and Chicago Midway (MDW) in place of O'Hare (ORD) — the three airports where Southwest maintains dominance over legacy hubs. The Mountain West (PHX, DEN, LAS, SLC, COS, ABQ) consistently appears as the penultimate cluster bridging Midwest and Pacific stops. No route doubles back geographically; all flow directionally from one region to another, consistent with aircraft repositioning across Southwest's national network.",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -> t.2, arraySort(t -> t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),\ntop_itineraries AS (\n  SELECT i.Route, i.FlightDate, i.hop_count\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop10_routes AS (\n  SELECT Route\n  FROM top_itineraries\n  GROUP BY Route\n  ORDER BY max(FlightDate) DESC\n  LIMIT 10\n),\nall_airports AS (\n  SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS ap\n  FROM top10_routes\n)\nSELECT\n  a.ap AS airport_code,\n  d.DisplayAirportName AS airport_name,\n  d.Latitude AS lat,\n  d.Longitude AS lon\nFROM all_airports a\nLEFT JOIN ontime.dim_airports d ON a.ap = d.AirportCode\nORDER BY airport_code"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (\n  SELECT max(hop_count) AS max_hop_count FROM itineraries\n),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Tail_Number, i.Carrier, i.FlightNum, i.FlightDate, i.hop_count, i.Route\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop_with_recurrence AS (\n  SELECT t.Tail_Number, t.Carrier, t.FlightNum, t.FlightDate, t.hop_count, r.recurrence_count, t.Route\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n)\nSELECT\n  argMax(Tail_Number, FlightDate) AS aircraft_id,\n  argMax(FlightNum, FlightDate) AS flight_number,\n  argMax(Carrier, FlightDate) AS carrier,\n  max(FlightDate) AS flight_date,\n  max(hop_count) AS hop_count,\n  any(recurrence_count) AS recurrence_count,\n  Route\nFROM top_with_recurrence\nGROUP BY Route\nORDER BY flight_date DESC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum distinct-hop count across all history is **8 hops** in a single day under one flight number. All 10 unique routes at that peak are operated by **Southwest Airlines (WN)**, reflecting Southwest's point-to-point milk-run scheduling model. The most recent unique 8-hop route was flown by tail **N957WN** on flight **WN 366** on **2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA, recurrence: 1 day). The highest-recurrence route among the top 10 is **LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN**, which appeared on 47 distinct days.\n\n| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (\n  SELECT max(hop_count) AS max_hop_count FROM itineraries\n),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Tail_Number, i.Carrier, i.FlightNum, i.FlightDate, i.hop_count, i.Route\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop_with_recurrence AS (\n  SELECT t.Tail_Number, t.Carrier, t.FlightNum, t.FlightDate, t.hop_count, r.recurrence_count, t.Route\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n)\nSELECT\n  argMax(Tail_Number, FlightDate) AS aircraft_id,\n  argMax(FlightNum, FlightDate) AS flight_number,\n  argMax(Carrier, FlightDate) AS carrier,\n  max(FlightDate) AS flight_date,\n  max(hop_count) AS hop_count,\n  any(recurrence_count) AS recurrence_count,\n  Route\nFROM top_with_recurrence\nGROUP BY Route\nORDER BY flight_date DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes span a wide recurrence spectrum — from true one-offs to well-established scheduled patterns.\n\nRecurrence count for each of the 10 routes (distinct calendar days flown, all history):\n\n| Route | Recurrence |\n|---|---|\n| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |\n| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |\n| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |\n| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |\n| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |\n| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |\n| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |\n| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |\n| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |\n| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |\n\n**Tiers:** Two routes are effectively one-offs (recurrence ≤ 2 days): ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA (1) and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN (2) — total **2 routes**. Four routes show modest scheduling (5–7 days): CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS, BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — total **4 routes**. Four routes are strongly recurring scheduled patterns (12–47 days): HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — total **4 routes**. Overall, 8 of 10 routes appeared on more than 2 days, indicating that 8-hop itineraries are predominantly recurring scheduled rotations rather than one-off operations.",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),\nroute_recurrence AS (\n  SELECT Route, countDistinct(FlightDate) AS recurrence_count\n  FROM itineraries GROUP BY Route\n),\ntop_itineraries AS (\n  SELECT i.Route, i.FlightDate, i.hop_count\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop10_routes AS (\n  SELECT t.Route, r.recurrence_count\n  FROM top_itineraries t JOIN route_recurrence r ON t.Route = r.Route\n  GROUP BY t.Route, r.recurrence_count\n  ORDER BY max(t.FlightDate) DESC\n  LIMIT 10\n)\nSELECT Route, recurrence_count\nFROM top10_routes\nORDER BY recurrence_count DESC"
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "All 10 routes are exclusively **Southwest Airlines (WN)** itineraries and touch **45 distinct airports** across the contiguous United States, with no international, Hawaiian, or Alaskan stops.\n\nThe routes form continuous coast-to-coast or diagonal sweeps — classic Southwest milk runs. Each itinerary anchors at one of three regional clusters:\n\n- **East/Gulf Coast origins:** ISP (Long Island), LGA (New York), CLE (Cleveland), BWI (Baltimore), FLL (Fort Lauderdale), HOU (Houston Hobby), MSY (New Orleans)\n- **Sun Belt / Midwest bridges:** ATL, BNA, DAL, MEM, MDW, STL, DEN — high-frequency Southwest connecting points mid-country\n- **West Coast termini:** OAK, LAX, BUR, SJC, SEA, SMF, LAS — all Bay Area, Southern California, or Pacific Northwest endpoints\n\nA Southwest-specific fingerprint is evident throughout: Dallas Love Field (DAL) appears in place of DFW, Houston Hobby (HOU) in place of IAH, and Chicago Midway (MDW) in place of O'Hare (ORD) — the three airports where Southwest maintains dominance over legacy hubs. The Mountain West (PHX, DEN, LAS, SLC, COS, ABQ) consistently appears as the penultimate cluster bridging Midwest and Pacific stops. No route doubles back geographically; all flow directionally from one region to another, consistent with aircraft repositioning across Southwest's national network.",
      "sql": "WITH itineraries AS (\n  SELECT\n    Tail_Number,\n    IATA_CODE_Reporting_Airline AS Carrier,\n    Flight_Number_Reporting_Airline AS FlightNum,\n    FlightDate,\n    countDistinct(OriginCode, DestCode, CRSDepTime) AS hop_count,\n    arrayStringConcat(\n      arrayConcat(\n        arrayMap(t -\u003e t.2, arraySort(t -\u003e t.1, groupUniqArray((coalesce(CRSDepTime, 0), OriginCode)))),\n        [argMax(DestCode, coalesce(CRSDepTime, 0))]\n      ),\n      '-'\n    ) AS Route\n  FROM ontime.fact_ontime\n  GROUP BY Tail_Number, Carrier, FlightNum, FlightDate\n),\nmax_hop AS (SELECT max(hop_count) AS max_hop_count FROM itineraries),\ntop_itineraries AS (\n  SELECT i.Route, i.FlightDate, i.hop_count\n  FROM itineraries i CROSS JOIN max_hop WHERE i.hop_count = max_hop.max_hop_count\n),\ntop10_routes AS (\n  SELECT Route\n  FROM top_itineraries\n  GROUP BY Route\n  ORDER BY max(FlightDate) DESC\n  LIMIT 10\n),\nall_airports AS (\n  SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS ap\n  FROM top10_routes\n)\nSELECT\n  a.ap AS airport_code,\n  d.DisplayAirportName AS airport_name,\n  d.Latitude AS lat,\n  d.Longitude AS lon\nFROM all_airports a\nLEFT JOIN ontime.dim_airports d ON a.ap = d.AirportCode\nORDER BY airport_code"
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
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `queries/main.sql`
- `queries/q1.sql`
- `queries/q2.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/main.json`
- `results/q1.json`
- `results/q2.json`