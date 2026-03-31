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

The maximum daily hops for a single aircraft on one flight number is **8 hops** (9 airports). All top results belong to **Southwest Airlines (WN)**, consistent with their point-to-point network model where a single aircraft cycles through many cities under one flight number.

The 10 most recent unique 8-hop routes are:

| # | Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |
|---|----------|--------|---------|------|------|------------|-------|
| 1 | N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |
| 2 | N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |
| 3 | N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |
| 4 | N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |
| 5 | N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |
| 6 | N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |
| 7 | N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |
| 8 | N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |
| 9 | N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |
| 10 | N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

The top 10 itineraries show a mix of recurring scheduled patterns and occasional/one-off routes.

**Recurrence counts for each route:**

1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA — **1 day** (one-off)
2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN — **5 days**
3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN — **2 days**
4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX — **5 days**
5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC — **40 days**
6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — **47 days**
7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS — **7 days**
8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK — **20 days**
9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK — **12 days**
10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — **5 days**

**Categories:**
- **High-recurrence scheduled patterns** (20+ days): 3 routes (#5, #6, #8) — these ran for weeks or months as stable daily scheduled rotations.
- **Moderate-recurrence patterns** (5–12 days): 5 routes (#2, #4, #7, #9, #10) — likely short-lived seasonal or weekly schedule patterns.
- **Low-recurrence / one-offs** (1–2 days): 2 routes (#1, #3) — appear to be irregular or ad-hoc itineraries, possibly driven by irregular operations or one-time schedule experiments.

Overall, **8 of the 10 routes** recurred on multiple days, indicating they represent recurring scheduled patterns rather than one-offs. The two highest-recurrence routes (WN 2884 at 47 days and WN 1956 at 40 days) were clearly long-running scheduled rotations.

- Rows returned: 10
- Columns: Route, recurrence_count, most_recent_date

| Route | recurrence_count | most_recent_date |
| --- | --- | --- |
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 | 2024-12-01T00:00:00Z |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

The top 10 eight-hop routes span **45 unique airports** across the entire continental United States, forming **coast-to-coast transcontinental zigzag patterns**.

**Key geographic observations:**

- **East-to-West sweep dominates:** Most routes originate on the East Coast or Gulf Coast (ISP, BWI, LGA, CLE, MSY, HOU) and terminate on the West Coast (SEA, OAK, SJC, DEN, SAN, LAX, LAS), reflecting Southwest's strategy of repositioning aircraft westward across the day.
- **Major Southwest hubs recur frequently:** DAL (Dallas Love), MDW (Chicago Midway), PHX (Phoenix), LAS (Las Vegas), OAK (Oakland), BWI (Baltimore), and BNA (Nashville) appear in multiple routes, anchoring the network's mid-continent and endpoint connections.
- **Mid-continent corridor:** Several routes thread through the Great Plains and Mountain West via ICT (Wichita), MCI (Kansas City), DEN (Denver), COS (Colorado Springs), SLC (Salt Lake City), and ABQ (Albuquerque) — reflecting Southwest's strong domestic interior presence.
- **Gulf Coast and Southeast representation:** MSY (New Orleans), HOU (Houston), PNS (Pensacola), VPS (Destin-Fort Walton Beach), MYR (Myrtle Beach), ATL (Atlanta), and FLL (Fort Lauderdale) appear, showing Sun Belt leisure destinations woven into operational routings.
- **No international or Hawaii airports:** All 45 airports are within the contiguous 48 states, consistent with Southwest's domestic point-to-point model.

- Rows returned: 45
- Columns: airport_code, DisplayAirportName, Latitude, Longitude

| airport_code | DisplayAirportName | Latitude | Longitude |
| --- | --- | --- | --- |
| SEA | Seattle/Tacoma International | 47.45 | -122.31166667 |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "id": "main",
      "answer_markdown": "The maximum daily hops for a single aircraft on one flight number is **8 hops** (9 airports). All top results belong to **Southwest Airlines (WN)**, consistent with their point-to-point network model where a single aircraft cycles through many cities under one flight number.\n\nThe 10 most recent unique 8-hop routes are:\n\n| # | Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|----------|--------|---------|------|------|------------|-------|\n| 1 | N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| 2 | N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| 3 | N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| 4 | N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| 5 | N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| 6 | N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| 7 | N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| 8 | N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| 9 | N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| 10 | N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH candidates AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate\n    FROM (\n        SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n               uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS hops\n        FROM ontime.fact_ontime\n        GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n        HAVING hops = (\n            SELECT max(h) FROM (\n                SELECT uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS h\n                FROM ontime.fact_ontime\n                GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n            )\n        )\n    )\n),\ndeduped_legs AS (\n    SELECT DISTINCT f.Tail_Number, f.FlightNum, f.Carrier, f.FlightDate,\n           f.OriginCode, f.DestCode, assumeNotNull(f.CRSDepTime) AS dep\n    FROM ontime.fact_ontime f\n    INNER JOIN candidates c\n        ON f.Tail_Number = c.Tail_Number AND f.FlightNum = c.FlightNum\n        AND f.Carrier = c.Carrier AND f.FlightDate = c.FlightDate\n),\nwith_route AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n           count() AS hops,\n           concat(\n               arrayStringConcat(groupArray(OriginCode), '-'),\n               '-',\n               arrayElement(groupArray(DestCode), count()::Int32)\n           ) AS Route\n    FROM (SELECT * FROM deduped_legs ORDER BY Tail_Number, FlightNum, Carrier, FlightDate, dep)\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\nroute_recurrence AS (\n    SELECT Route, uniqExact(FlightDate) AS recurrence_count\n    FROM with_route\n    GROUP BY Route\n),\nranked AS (\n    SELECT w.Tail_Number, w.FlightNum, w.Carrier, w.FlightDate, w.hops, w.Route,\n           r.recurrence_count,\n           ROW_NUMBER() OVER (PARTITION BY w.Route ORDER BY w.FlightDate DESC) AS rn\n    FROM with_route w\n    JOIN route_recurrence r ON w.Route = r.Route\n)\nSELECT\n    if(Tail_Number = '', 'unknown', Tail_Number) AS aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hops AS hop_count,\n    recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY FlightDate DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "answer_markdown": "The top 10 itineraries show a mix of recurring scheduled patterns and occasional/one-off routes.\n\n**Recurrence counts for each route:**\n\n1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA — **1 day** (one-off)\n2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN — **5 days**\n3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN — **2 days**\n4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX — **5 days**\n5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC — **40 days**\n6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — **47 days**\n7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS — **7 days**\n8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK — **20 days**\n9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK — **12 days**\n10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — **5 days**\n\n**Categories:**\n- **High-recurrence scheduled patterns** (20+ days): 3 routes (#5, #6, #8) — these ran for weeks or months as stable daily scheduled rotations.\n- **Moderate-recurrence patterns** (5–12 days): 5 routes (#2, #4, #7, #9, #10) — likely short-lived seasonal or weekly schedule patterns.\n- **Low-recurrence / one-offs** (1–2 days): 2 routes (#1, #3) — appear to be irregular or ad-hoc itineraries, possibly driven by irregular operations or one-time schedule experiments.\n\nOverall, **8 of the 10 routes** recurred on multiple days, indicating they represent recurring scheduled patterns rather than one-offs. The two highest-recurrence routes (WN 2884 at 47 days and WN 1956 at 40 days) were clearly long-running scheduled rotations.",
      "sql": "WITH candidates AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate\n    FROM (\n        SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n               uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS hops\n        FROM ontime.fact_ontime\n        GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n        HAVING hops = (\n            SELECT max(h) FROM (\n                SELECT uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS h\n                FROM ontime.fact_ontime\n                GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n            )\n        )\n    )\n),\ndeduped_legs AS (\n    SELECT DISTINCT f.Tail_Number, f.FlightNum, f.Carrier, f.FlightDate,\n           f.OriginCode, f.DestCode, assumeNotNull(f.CRSDepTime) AS dep\n    FROM ontime.fact_ontime f\n    INNER JOIN candidates c\n        ON f.Tail_Number = c.Tail_Number AND f.FlightNum = c.FlightNum\n        AND f.Carrier = c.Carrier AND f.FlightDate = c.FlightDate\n),\nwith_route AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n           count() AS hops,\n           concat(\n               arrayStringConcat(groupArray(OriginCode), '-'),\n               '-',\n               arrayElement(groupArray(DestCode), count()::Int32)\n           ) AS Route\n    FROM (SELECT * FROM deduped_legs ORDER BY Tail_Number, FlightNum, Carrier, FlightDate, dep)\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\nroute_recurrence AS (\n    SELECT Route, uniqExact(FlightDate) AS recurrence_count\n    FROM with_route\n    GROUP BY Route\n),\nranked AS (\n    SELECT w.Route, r.recurrence_count,\n           ROW_NUMBER() OVER (PARTITION BY w.Route ORDER BY w.FlightDate DESC) AS rn,\n           w.FlightDate\n    FROM with_route w\n    JOIN route_recurrence r ON w.Route = r.Route\n)\nSELECT Route, recurrence_count,\n       FlightDate AS most_recent_date\nFROM ranked\nWHERE rn = 1\nORDER BY FlightDate DESC\nLIMIT 10"
    },
    {
      "id": "q2",
      "answer_markdown": "The top 10 eight-hop routes span **45 unique airports** across the entire continental United States, forming **coast-to-coast transcontinental zigzag patterns**.\n\n**Key geographic observations:**\n\n- **East-to-West sweep dominates:** Most routes originate on the East Coast or Gulf Coast (ISP, BWI, LGA, CLE, MSY, HOU) and terminate on the West Coast (SEA, OAK, SJC, DEN, SAN, LAX, LAS), reflecting Southwest's strategy of repositioning aircraft westward across the day.\n- **Major Southwest hubs recur frequently:** DAL (Dallas Love), MDW (Chicago Midway), PHX (Phoenix), LAS (Las Vegas), OAK (Oakland), BWI (Baltimore), and BNA (Nashville) appear in multiple routes, anchoring the network's mid-continent and endpoint connections.\n- **Mid-continent corridor:** Several routes thread through the Great Plains and Mountain West via ICT (Wichita), MCI (Kansas City), DEN (Denver), COS (Colorado Springs), SLC (Salt Lake City), and ABQ (Albuquerque) — reflecting Southwest's strong domestic interior presence.\n- **Gulf Coast and Southeast representation:** MSY (New Orleans), HOU (Houston), PNS (Pensacola), VPS (Destin-Fort Walton Beach), MYR (Myrtle Beach), ATL (Atlanta), and FLL (Fort Lauderdale) appear, showing Sun Belt leisure destinations woven into operational routings.\n- **No international or Hawaii airports:** All 45 airports are within the contiguous 48 states, consistent with Southwest's domestic point-to-point model.",
      "sql": "WITH top_routes AS (\n    SELECT arrayJoin(splitByChar('-', Route)) AS airport_code\n    FROM (\n        SELECT 'ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA' AS Route\n        UNION ALL SELECT 'CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN'\n        UNION ALL SELECT 'ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN'\n        UNION ALL SELECT 'MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX'\n        UNION ALL SELECT 'MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC'\n        UNION ALL SELECT 'LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN'\n        UNION ALL SELECT 'SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS'\n        UNION ALL SELECT 'HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK'\n        UNION ALL SELECT 'BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK'\n        UNION ALL SELECT 'BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX'\n    )\n)\nSELECT DISTINCT\n    tr.airport_code,\n    d.DisplayAirportName,\n    d.Latitude,\n    d.Longitude\nFROM top_routes tr\nLEFT JOIN ontime.dim_airports d ON tr.airport_code = d.AirportCode\nORDER BY d.Longitude"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH candidates AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate\n    FROM (\n        SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n               uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS hops\n        FROM ontime.fact_ontime\n        GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n        HAVING hops = (\n            SELECT max(h) FROM (\n                SELECT uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS h\n                FROM ontime.fact_ontime\n                GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n            )\n        )\n    )\n),\ndeduped_legs AS (\n    SELECT DISTINCT f.Tail_Number, f.FlightNum, f.Carrier, f.FlightDate,\n           f.OriginCode, f.DestCode, assumeNotNull(f.CRSDepTime) AS dep\n    FROM ontime.fact_ontime f\n    INNER JOIN candidates c\n        ON f.Tail_Number = c.Tail_Number AND f.FlightNum = c.FlightNum\n        AND f.Carrier = c.Carrier AND f.FlightDate = c.FlightDate\n),\nwith_route AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n           count() AS hops,\n           concat(\n               arrayStringConcat(groupArray(OriginCode), '-'),\n               '-',\n               arrayElement(groupArray(DestCode), count()::Int32)\n           ) AS Route\n    FROM (SELECT * FROM deduped_legs ORDER BY Tail_Number, FlightNum, Carrier, FlightDate, dep)\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\nroute_recurrence AS (\n    SELECT Route, uniqExact(FlightDate) AS recurrence_count\n    FROM with_route\n    GROUP BY Route\n),\nranked AS (\n    SELECT w.Tail_Number, w.FlightNum, w.Carrier, w.FlightDate, w.hops, w.Route,\n           r.recurrence_count,\n           ROW_NUMBER() OVER (PARTITION BY w.Route ORDER BY w.FlightDate DESC) AS rn\n    FROM with_route w\n    JOIN route_recurrence r ON w.Route = r.Route\n)\nSELECT\n    if(Tail_Number = '', 'unknown', Tail_Number) AS aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hops AS hop_count,\n    recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY FlightDate DESC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\nDo not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.\nCount hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum daily hops for a single aircraft on one flight number is **8 hops** (9 airports). All top results belong to **Southwest Airlines (WN)**, consistent with their point-to-point network model where a single aircraft cycles through many cities under one flight number.\n\nThe 10 most recent unique 8-hop routes are:\n\n| # | Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|----------|--------|---------|------|------|------------|-------|\n| 1 | N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| 2 | N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| 3 | N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| 4 | N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| 5 | N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| 6 | N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| 7 | N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| 8 | N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| 9 | N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| 10 | N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH candidates AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate\n    FROM (\n        SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n               uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS hops\n        FROM ontime.fact_ontime\n        GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n        HAVING hops = (\n            SELECT max(h) FROM (\n                SELECT uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS h\n                FROM ontime.fact_ontime\n                GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n            )\n        )\n    )\n),\ndeduped_legs AS (\n    SELECT DISTINCT f.Tail_Number, f.FlightNum, f.Carrier, f.FlightDate,\n           f.OriginCode, f.DestCode, assumeNotNull(f.CRSDepTime) AS dep\n    FROM ontime.fact_ontime f\n    INNER JOIN candidates c\n        ON f.Tail_Number = c.Tail_Number AND f.FlightNum = c.FlightNum\n        AND f.Carrier = c.Carrier AND f.FlightDate = c.FlightDate\n),\nwith_route AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n           count() AS hops,\n           concat(\n               arrayStringConcat(groupArray(OriginCode), '-'),\n               '-',\n               arrayElement(groupArray(DestCode), count()::Int32)\n           ) AS Route\n    FROM (SELECT * FROM deduped_legs ORDER BY Tail_Number, FlightNum, Carrier, FlightDate, dep)\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\nroute_recurrence AS (\n    SELECT Route, uniqExact(FlightDate) AS recurrence_count\n    FROM with_route\n    GROUP BY Route\n),\nranked AS (\n    SELECT w.Tail_Number, w.FlightNum, w.Carrier, w.FlightDate, w.hops, w.Route,\n           r.recurrence_count,\n           ROW_NUMBER() OVER (PARTITION BY w.Route ORDER BY w.FlightDate DESC) AS rn\n    FROM with_route w\n    JOIN route_recurrence r ON w.Route = r.Route\n)\nSELECT\n    if(Tail_Number = '', 'unknown', Tail_Number) AS aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hops AS hop_count,\n    recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY FlightDate DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The top 10 itineraries show a mix of recurring scheduled patterns and occasional/one-off routes.\n\n**Recurrence counts for each route:**\n\n1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA — **1 day** (one-off)\n2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN — **5 days**\n3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN — **2 days**\n4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX — **5 days**\n5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC — **40 days**\n6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — **47 days**\n7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS — **7 days**\n8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK — **20 days**\n9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK — **12 days**\n10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — **5 days**\n\n**Categories:**\n- **High-recurrence scheduled patterns** (20+ days): 3 routes (#5, #6, #8) — these ran for weeks or months as stable daily scheduled rotations.\n- **Moderate-recurrence patterns** (5–12 days): 5 routes (#2, #4, #7, #9, #10) — likely short-lived seasonal or weekly schedule patterns.\n- **Low-recurrence / one-offs** (1–2 days): 2 routes (#1, #3) — appear to be irregular or ad-hoc itineraries, possibly driven by irregular operations or one-time schedule experiments.\n\nOverall, **8 of the 10 routes** recurred on multiple days, indicating they represent recurring scheduled patterns rather than one-offs. The two highest-recurrence routes (WN 2884 at 47 days and WN 1956 at 40 days) were clearly long-running scheduled rotations.",
      "sql": "WITH candidates AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate\n    FROM (\n        SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n               uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS hops\n        FROM ontime.fact_ontime\n        GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n        HAVING hops = (\n            SELECT max(h) FROM (\n                SELECT uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS h\n                FROM ontime.fact_ontime\n                GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n            )\n        )\n    )\n),\ndeduped_legs AS (\n    SELECT DISTINCT f.Tail_Number, f.FlightNum, f.Carrier, f.FlightDate,\n           f.OriginCode, f.DestCode, assumeNotNull(f.CRSDepTime) AS dep\n    FROM ontime.fact_ontime f\n    INNER JOIN candidates c\n        ON f.Tail_Number = c.Tail_Number AND f.FlightNum = c.FlightNum\n        AND f.Carrier = c.Carrier AND f.FlightDate = c.FlightDate\n),\nwith_route AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n           count() AS hops,\n           concat(\n               arrayStringConcat(groupArray(OriginCode), '-'),\n               '-',\n               arrayElement(groupArray(DestCode), count()::Int32)\n           ) AS Route\n    FROM (SELECT * FROM deduped_legs ORDER BY Tail_Number, FlightNum, Carrier, FlightDate, dep)\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\nroute_recurrence AS (\n    SELECT Route, uniqExact(FlightDate) AS recurrence_count\n    FROM with_route\n    GROUP BY Route\n),\nranked AS (\n    SELECT w.Route, r.recurrence_count,\n           ROW_NUMBER() OVER (PARTITION BY w.Route ORDER BY w.FlightDate DESC) AS rn,\n           w.FlightDate\n    FROM with_route w\n    JOIN route_recurrence r ON w.Route = r.Route\n)\nSELECT Route, recurrence_count,\n       FlightDate AS most_recent_date\nFROM ranked\nWHERE rn = 1\nORDER BY FlightDate DESC\nLIMIT 10"
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "The top 10 eight-hop routes span **45 unique airports** across the entire continental United States, forming **coast-to-coast transcontinental zigzag patterns**.\n\n**Key geographic observations:**\n\n- **East-to-West sweep dominates:** Most routes originate on the East Coast or Gulf Coast (ISP, BWI, LGA, CLE, MSY, HOU) and terminate on the West Coast (SEA, OAK, SJC, DEN, SAN, LAX, LAS), reflecting Southwest's strategy of repositioning aircraft westward across the day.\n- **Major Southwest hubs recur frequently:** DAL (Dallas Love), MDW (Chicago Midway), PHX (Phoenix), LAS (Las Vegas), OAK (Oakland), BWI (Baltimore), and BNA (Nashville) appear in multiple routes, anchoring the network's mid-continent and endpoint connections.\n- **Mid-continent corridor:** Several routes thread through the Great Plains and Mountain West via ICT (Wichita), MCI (Kansas City), DEN (Denver), COS (Colorado Springs), SLC (Salt Lake City), and ABQ (Albuquerque) — reflecting Southwest's strong domestic interior presence.\n- **Gulf Coast and Southeast representation:** MSY (New Orleans), HOU (Houston), PNS (Pensacola), VPS (Destin-Fort Walton Beach), MYR (Myrtle Beach), ATL (Atlanta), and FLL (Fort Lauderdale) appear, showing Sun Belt leisure destinations woven into operational routings.\n- **No international or Hawaii airports:** All 45 airports are within the contiguous 48 states, consistent with Southwest's domestic point-to-point model.",
      "sql": "WITH top_routes AS (\n    SELECT arrayJoin(splitByChar('-', Route)) AS airport_code\n    FROM (\n        SELECT 'ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA' AS Route\n        UNION ALL SELECT 'CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN'\n        UNION ALL SELECT 'ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN'\n        UNION ALL SELECT 'MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX'\n        UNION ALL SELECT 'MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC'\n        UNION ALL SELECT 'LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN'\n        UNION ALL SELECT 'SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS'\n        UNION ALL SELECT 'HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK'\n        UNION ALL SELECT 'BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK'\n        UNION ALL SELECT 'BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX'\n    )\n)\nSELECT DISTINCT\n    tr.airport_code,\n    d.DisplayAirportName,\n    d.Latitude,\n    d.Longitude\nFROM top_routes tr\nLEFT JOIN ontime.dim_airports d ON tr.airport_code = d.AirportCode\nORDER BY d.Longitude"
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
      "answer_markdown": "The maximum daily hops for a single aircraft on one flight number is **8 hops** (9 airports). All top results belong to **Southwest Airlines (WN)**, consistent with their point-to-point network model where a single aircraft cycles through many cities under one flight number.\n\nThe 10 most recent unique 8-hop routes are:\n\n| # | Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|----------|--------|---------|------|------|------------|-------|\n| 1 | N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| 2 | N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| 3 | N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| 4 | N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| 5 | N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| 6 | N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| 7 | N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| 8 | N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| 9 | N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| 10 | N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH candidates AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate\n    FROM (\n        SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n               uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS hops\n        FROM ontime.fact_ontime\n        GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n        HAVING hops = (\n            SELECT max(h) FROM (\n                SELECT uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS h\n                FROM ontime.fact_ontime\n                GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n            )\n        )\n    )\n),\ndeduped_legs AS (\n    SELECT DISTINCT f.Tail_Number, f.FlightNum, f.Carrier, f.FlightDate,\n           f.OriginCode, f.DestCode, assumeNotNull(f.CRSDepTime) AS dep\n    FROM ontime.fact_ontime f\n    INNER JOIN candidates c\n        ON f.Tail_Number = c.Tail_Number AND f.FlightNum = c.FlightNum\n        AND f.Carrier = c.Carrier AND f.FlightDate = c.FlightDate\n),\nwith_route AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n           count() AS hops,\n           concat(\n               arrayStringConcat(groupArray(OriginCode), '-'),\n               '-',\n               arrayElement(groupArray(DestCode), count()::Int32)\n           ) AS Route\n    FROM (SELECT * FROM deduped_legs ORDER BY Tail_Number, FlightNum, Carrier, FlightDate, dep)\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\nroute_recurrence AS (\n    SELECT Route, uniqExact(FlightDate) AS recurrence_count\n    FROM with_route\n    GROUP BY Route\n),\nranked AS (\n    SELECT w.Tail_Number, w.FlightNum, w.Carrier, w.FlightDate, w.hops, w.Route,\n           r.recurrence_count,\n           ROW_NUMBER() OVER (PARTITION BY w.Route ORDER BY w.FlightDate DESC) AS rn\n    FROM with_route w\n    JOIN route_recurrence r ON w.Route = r.Route\n)\nSELECT\n    if(Tail_Number = '', 'unknown', Tail_Number) AS aircraft_id,\n    FlightNum AS flight_number,\n    Carrier AS carrier,\n    FlightDate AS flight_date,\n    hops AS hop_count,\n    recurrence_count,\n    Route\nFROM ranked\nWHERE rn = 1\nORDER BY FlightDate DESC\nLIMIT 10",
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
      "answer_markdown": "The top 10 itineraries show a mix of recurring scheduled patterns and occasional/one-off routes.\n\n**Recurrence counts for each route:**\n\n1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA — **1 day** (one-off)\n2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN — **5 days**\n3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN — **2 days**\n4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX — **5 days**\n5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC — **40 days**\n6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — **47 days**\n7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS — **7 days**\n8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK — **20 days**\n9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK — **12 days**\n10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — **5 days**\n\n**Categories:**\n- **High-recurrence scheduled patterns** (20+ days): 3 routes (#5, #6, #8) — these ran for weeks or months as stable daily scheduled rotations.\n- **Moderate-recurrence patterns** (5–12 days): 5 routes (#2, #4, #7, #9, #10) — likely short-lived seasonal or weekly schedule patterns.\n- **Low-recurrence / one-offs** (1–2 days): 2 routes (#1, #3) — appear to be irregular or ad-hoc itineraries, possibly driven by irregular operations or one-time schedule experiments.\n\nOverall, **8 of the 10 routes** recurred on multiple days, indicating they represent recurring scheduled patterns rather than one-offs. The two highest-recurrence routes (WN 2884 at 47 days and WN 1956 at 40 days) were clearly long-running scheduled rotations.",
      "sql": "WITH candidates AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate\n    FROM (\n        SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n               uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS hops\n        FROM ontime.fact_ontime\n        GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n        HAVING hops = (\n            SELECT max(h) FROM (\n                SELECT uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS h\n                FROM ontime.fact_ontime\n                GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n            )\n        )\n    )\n),\ndeduped_legs AS (\n    SELECT DISTINCT f.Tail_Number, f.FlightNum, f.Carrier, f.FlightDate,\n           f.OriginCode, f.DestCode, assumeNotNull(f.CRSDepTime) AS dep\n    FROM ontime.fact_ontime f\n    INNER JOIN candidates c\n        ON f.Tail_Number = c.Tail_Number AND f.FlightNum = c.FlightNum\n        AND f.Carrier = c.Carrier AND f.FlightDate = c.FlightDate\n),\nwith_route AS (\n    SELECT Tail_Number, FlightNum, Carrier, FlightDate,\n           count() AS hops,\n           concat(\n               arrayStringConcat(groupArray(OriginCode), '-'),\n               '-',\n               arrayElement(groupArray(DestCode), count()::Int32)\n           ) AS Route\n    FROM (SELECT * FROM deduped_legs ORDER BY Tail_Number, FlightNum, Carrier, FlightDate, dep)\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\nroute_recurrence AS (\n    SELECT Route, uniqExact(FlightDate) AS recurrence_count\n    FROM with_route\n    GROUP BY Route\n),\nranked AS (\n    SELECT w.Route, r.recurrence_count,\n           ROW_NUMBER() OVER (PARTITION BY w.Route ORDER BY w.FlightDate DESC) AS rn,\n           w.FlightDate\n    FROM with_route w\n    JOIN route_recurrence r ON w.Route = r.Route\n)\nSELECT Route, recurrence_count,\n       FlightDate AS most_recent_date\nFROM ranked\nWHERE rn = 1\nORDER BY FlightDate DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "Route",
        "recurrence_count",
        "most_recent_date"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "most_recent_date": "2024-12-01T00:00:00Z",
        "recurrence_count": 1
      }
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "The top 10 eight-hop routes span **45 unique airports** across the entire continental United States, forming **coast-to-coast transcontinental zigzag patterns**.\n\n**Key geographic observations:**\n\n- **East-to-West sweep dominates:** Most routes originate on the East Coast or Gulf Coast (ISP, BWI, LGA, CLE, MSY, HOU) and terminate on the West Coast (SEA, OAK, SJC, DEN, SAN, LAX, LAS), reflecting Southwest's strategy of repositioning aircraft westward across the day.\n- **Major Southwest hubs recur frequently:** DAL (Dallas Love), MDW (Chicago Midway), PHX (Phoenix), LAS (Las Vegas), OAK (Oakland), BWI (Baltimore), and BNA (Nashville) appear in multiple routes, anchoring the network's mid-continent and endpoint connections.\n- **Mid-continent corridor:** Several routes thread through the Great Plains and Mountain West via ICT (Wichita), MCI (Kansas City), DEN (Denver), COS (Colorado Springs), SLC (Salt Lake City), and ABQ (Albuquerque) — reflecting Southwest's strong domestic interior presence.\n- **Gulf Coast and Southeast representation:** MSY (New Orleans), HOU (Houston), PNS (Pensacola), VPS (Destin-Fort Walton Beach), MYR (Myrtle Beach), ATL (Atlanta), and FLL (Fort Lauderdale) appear, showing Sun Belt leisure destinations woven into operational routings.\n- **No international or Hawaii airports:** All 45 airports are within the contiguous 48 states, consistent with Southwest's domestic point-to-point model.",
      "sql": "WITH top_routes AS (\n    SELECT arrayJoin(splitByChar('-', Route)) AS airport_code\n    FROM (\n        SELECT 'ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA' AS Route\n        UNION ALL SELECT 'CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN'\n        UNION ALL SELECT 'ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN'\n        UNION ALL SELECT 'MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX'\n        UNION ALL SELECT 'MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC'\n        UNION ALL SELECT 'LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN'\n        UNION ALL SELECT 'SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS'\n        UNION ALL SELECT 'HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK'\n        UNION ALL SELECT 'BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK'\n        UNION ALL SELECT 'BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX'\n    )\n)\nSELECT DISTINCT\n    tr.airport_code,\n    d.DisplayAirportName,\n    d.Latitude,\n    d.Longitude\nFROM top_routes tr\nLEFT JOIN ontime.dim_airports d ON tr.airport_code = d.AirportCode\nORDER BY d.Longitude",
      "row_count": 45,
      "result_columns": [
        "airport_code",
        "DisplayAirportName",
        "Latitude",
        "Longitude"
      ],
      "first_row": {
        "DisplayAirportName": "Seattle/Tacoma International",
        "Latitude": 47.45,
        "Longitude": -122.31166667,
        "airport_code": "SEA"
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