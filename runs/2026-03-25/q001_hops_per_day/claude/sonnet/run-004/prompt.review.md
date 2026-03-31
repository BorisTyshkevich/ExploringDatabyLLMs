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

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

The maximum number of hops flown by a single aircraft on a single flight number in one day is **8 hops** (9-airport itineraries). All 10 unique routes are operated by Southwest Airlines (WN). The most recent top 10 unique routes by departure date are:

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
- Columns: Tail_Number, FlightNum, Carrier, FlightDate, hop_count, recurrence_count, Route

| Tail_Number | FlightNum | Carrier | FlightDate | hop_count | recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

The 10 routes show a wide spread of recurrence, from a single occurrence to nearly 50 days:

| Route | Recurrence |
|---|---|
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |
| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |
| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |
| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |
| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |
| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |
| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |
| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |
| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |

**Tier breakdown (totals add to 10):**
- **One-off (1 day):** 1 route — likely an irregular or seasonal assignment.
- **Low recurrence (2–7 days):** 4 routes — occasional scheduled rotations, not daily.
- **Moderate recurrence (12–20 days):** 2 routes — semi-regular scheduled patterns.
- **High recurrence (40–47 days):** 2 routes — clearly stable, recurring scheduled rotations flown across many weeks.

Overall, the majority of 8-hop itineraries are recurring scheduled patterns. Only one is a true one-off; 9 of 10 routes were flown on at least 2 distinct days, and 4 of 10 were flown on 12 or more days, indicating these represent genuine repeating flight schedules rather than one-time anomalies.

- Rows returned: 719
- Columns: Route, recurrence_count

| Route | recurrence_count |
| --- | --- |
| EWR-MDW-ABQ-DEN-SMF-LAS-TUS-OAK-SLC | 1 |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

All airports across the 10 routes are in the **continental United States**, and every route is operated by **Southwest Airlines (WN)**. The itineraries are true transcontinental or cross-country sweeps, connecting three geographic bands:

- **West (17 airports):** SEA, OAK, SJC, SMF, RNO, LAX, BUR, SAN, PSP, LAS, PHX, SLC, ABQ, ELP, COS, DEN, MAF — dominated by Pacific Coast, Mountain West, and Desert Southwest airports.
- **Central (14 airports):** STL, MEM, ORD, MDW, MCI, ICT, LIT, BNA, DAL, HOU, MSY, JAN, PNS, VPS — the Mississippi/Gulf corridor and Midwest hubs.
- **East (14 airports):** LGA, ISP, BWI, IAD, RIC, RDU, MYR, CMH, CLE, DTW, ATL, TPA, MCO, FLL — stretching from the Northeast down the Atlantic Seaboard to Florida.

The dominant pattern is **coast-to-coast routing**: each itinerary typically originates on one coast or the Gulf South, threads through 2–3 inland stops, and terminates on the opposite coast. Several routes end or begin at West Coast leisure-travel airports (OAK, LAX, SJC, LAS). Recurring stops like BWI, BNA, DAL, LAS, MDW, and OAK appear across multiple routes, reflecting Southwest's key connecting cities. No route is confined to a single region; all 10 span at least two of the three geographic bands.

- Rows returned: 3
- Columns: region, airport_count, airports

| region | airport_count | airports |
| --- | --- | --- |
| West | 17 | [LAS, OAK, SEA, PHX, BUR, DEN, ELP, SAN, LAX, SLC, SJC, COS, SMF, PSP, RNO, ABQ, MAF] |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "id": "main",
      "answer_markdown": "The maximum number of hops flown by a single aircraft on a single flight number in one day is **8 hops** (9-airport itineraries). All 10 unique routes are operated by Southwest Airlines (WN). The most recent top 10 unique routes by departure date are:\n\n| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -> x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -> x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT * FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n),\nrecent_top AS (\n    SELECT\n        Route,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum, FlightDate) AS FlightNum,\n        argMax(Carrier, FlightDate) AS Carrier,\n        max(FlightDate) AS most_recent_date,\n        argMax(hop_count, FlightDate) AS hop_count\n    FROM top_route_days\n    GROUP BY Route\n)\nSELECT\n    rt.Tail_Number,\n    rt.FlightNum,\n    rt.Carrier,\n    rt.most_recent_date AS FlightDate,\n    rt.hop_count,\n    rr.recurrence_count,\n    rt.Route\nFROM recent_top rt\nJOIN route_recurrence rr ON rt.Route = rr.Route\nORDER BY rt.most_recent_date DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "answer_markdown": "The 10 routes show a wide spread of recurrence, from a single occurrence to nearly 50 days:\n\n| Route | Recurrence |\n|---|---|\n| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |\n| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |\n| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |\n| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |\n| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |\n| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |\n| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |\n| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |\n| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |\n| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |\n\n**Tier breakdown (totals add to 10):**\n- **One-off (1 day):** 1 route — likely an irregular or seasonal assignment.\n- **Low recurrence (2–7 days):** 4 routes — occasional scheduled rotations, not daily.\n- **Moderate recurrence (12–20 days):** 2 routes — semi-regular scheduled patterns.\n- **High recurrence (40–47 days):** 2 routes — clearly stable, recurring scheduled rotations flown across many weeks.\n\nOverall, the majority of 8-hop itineraries are recurring scheduled patterns. Only one is a true one-off; 9 of 10 routes were flown on at least 2 distinct days, and 4 of 10 were flown on 12 or more days, indicating these represent genuine repeating flight schedules rather than one-time anomalies.",
      "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -> x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -> x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT Route FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n)\nSELECT Route, recurrence_count\nFROM route_recurrence\nORDER BY recurrence_count ASC"
    },
    {
      "id": "q2",
      "answer_markdown": "All airports across the 10 routes are in the **continental United States**, and every route is operated by **Southwest Airlines (WN)**. The itineraries are true transcontinental or cross-country sweeps, connecting three geographic bands:\n\n- **West (17 airports):** SEA, OAK, SJC, SMF, RNO, LAX, BUR, SAN, PSP, LAS, PHX, SLC, ABQ, ELP, COS, DEN, MAF — dominated by Pacific Coast, Mountain West, and Desert Southwest airports.\n- **Central (14 airports):** STL, MEM, ORD, MDW, MCI, ICT, LIT, BNA, DAL, HOU, MSY, JAN, PNS, VPS — the Mississippi/Gulf corridor and Midwest hubs.\n- **East (14 airports):** LGA, ISP, BWI, IAD, RIC, RDU, MYR, CMH, CLE, DTW, ATL, TPA, MCO, FLL — stretching from the Northeast down the Atlantic Seaboard to Florida.\n\nThe dominant pattern is **coast-to-coast routing**: each itinerary typically originates on one coast or the Gulf South, threads through 2–3 inland stops, and terminates on the opposite coast. Several routes end or begin at West Coast leisure-travel airports (OAK, LAX, SJC, LAS). Recurring stops like BWI, BNA, DAL, LAS, MDW, and OAK appear across multiple routes, reflecting Southwest's key connecting cities. No route is confined to a single region; all 10 span at least two of the three geographic bands.",
      "sql": "SELECT\n    multiIf(d.Longitude < -100, 'West', d.Longitude < -85, 'Central', 'East') AS region,\n    count() AS airport_count,\n    groupArray(a.code) AS airports\nFROM (\n    SELECT DISTINCT arrayJoin([\n        'ISP','BWI','MYR','BNA','VPS','DAL','LAS','OAK','SEA',\n        'CLE','PNS','HOU','MCI','PHX','BUR','DEN',\n        'ELP','LIT','ATL','RIC','MDW','SAN',\n        'MSY','CMH','RDU','DTW','LAX',\n        'TPA','ORD','SLC','SJC',\n        'LGA','STL','ICT','COS','JAN',\n        'SMF','PSP','RNO',\n        'ABQ','MAF',\n        'FLL','MCO','MEM','IAD'\n    ]) AS code\n) a\nLEFT JOIN ontime.dim_airports d ON a.code = d.AirportCode\nGROUP BY region\nORDER BY min(d.Longitude)"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -\u003e x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT * FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n),\nrecent_top AS (\n    SELECT\n        Route,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum, FlightDate) AS FlightNum,\n        argMax(Carrier, FlightDate) AS Carrier,\n        max(FlightDate) AS most_recent_date,\n        argMax(hop_count, FlightDate) AS hop_count\n    FROM top_route_days\n    GROUP BY Route\n)\nSELECT\n    rt.Tail_Number,\n    rt.FlightNum,\n    rt.Carrier,\n    rt.most_recent_date AS FlightDate,\n    rt.hop_count,\n    rr.recurrence_count,\n    rt.Route\nFROM recent_top rt\nJOIN route_recurrence rr ON rt.Route = rr.Route\nORDER BY rt.most_recent_date DESC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "id": "main",
      "subquestion": "Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.\nDefine uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.\n\nReturn:\n\n- aircraft id\n- flight number\n- carrier\n- flight date\n- hop count\n- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft\n- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`",
      "answer_markdown": "The maximum number of hops flown by a single aircraft on a single flight number in one day is **8 hops** (9-airport itineraries). All 10 unique routes are operated by Southwest Airlines (WN). The most recent top 10 unique routes by departure date are:\n\n| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |\n|---|---|---|---|---|---|---|\n| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |\n| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |\n| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |\n| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |\n| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |\n| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |\n| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |\n| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |\n| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |\n| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |",
      "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -\u003e x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT * FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n),\nrecent_top AS (\n    SELECT\n        Route,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum, FlightDate) AS FlightNum,\n        argMax(Carrier, FlightDate) AS Carrier,\n        max(FlightDate) AS most_recent_date,\n        argMax(hop_count, FlightDate) AS hop_count\n    FROM top_route_days\n    GROUP BY Route\n)\nSELECT\n    rt.Tail_Number,\n    rt.FlightNum,\n    rt.Carrier,\n    rt.most_recent_date AS FlightDate,\n    rt.hop_count,\n    rr.recurrence_count,\n    rt.Route\nFROM recent_top rt\nJOIN route_recurrence rr ON rt.Route = rr.Route\nORDER BY rt.most_recent_date DESC\nLIMIT 10"
    },
    {
      "id": "q1",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?\n\nList the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.",
      "answer_markdown": "The 10 routes show a wide spread of recurrence, from a single occurrence to nearly 50 days:\n\n| Route | Recurrence |\n|---|---|\n| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |\n| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |\n| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |\n| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |\n| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |\n| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |\n| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |\n| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |\n| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |\n| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |\n\n**Tier breakdown (totals add to 10):**\n- **One-off (1 day):** 1 route — likely an irregular or seasonal assignment.\n- **Low recurrence (2–7 days):** 4 routes — occasional scheduled rotations, not daily.\n- **Moderate recurrence (12–20 days):** 2 routes — semi-regular scheduled patterns.\n- **High recurrence (40–47 days):** 2 routes — clearly stable, recurring scheduled rotations flown across many weeks.\n\nOverall, the majority of 8-hop itineraries are recurring scheduled patterns. Only one is a true one-off; 9 of 10 routes were flown on at least 2 distinct days, and 4 of 10 were flown on 12 or more days, indicating these represent genuine repeating flight schedules rather than one-time anomalies.",
      "sql": "WITH raw AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        IATA_CODE_Reporting_Airline AS Carrier,\n        FlightDate,\n        count() AS hop_count,\n        arraySort(x -\u003e x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate\n),\ndaily_routes AS (\n    SELECT\n        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                arrayMap(x -\u003e x.1, sorted_legs),\n                [sorted_legs[length(sorted_legs)].2]\n            ),\n            '-'\n        ) AS Route\n    FROM raw\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count FROM daily_routes\n),\ntop_route_days AS (\n    SELECT Route FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n),\nroute_recurrence AS (\n    SELECT Route, countDistinct(FlightDate) AS recurrence_count\n    FROM daily_routes\n    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)\n    GROUP BY Route\n)\nSELECT Route, recurrence_count\nFROM route_recurrence\nORDER BY recurrence_count ASC"
    },
    {
      "id": "q2",
      "subquestion": "What geographic pattern do the top itineraries show?\n\nBase the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.",
      "answer_markdown": "All airports across the 10 routes are in the **continental United States**, and every route is operated by **Southwest Airlines (WN)**. The itineraries are true transcontinental or cross-country sweeps, connecting three geographic bands:\n\n- **West (17 airports):** SEA, OAK, SJC, SMF, RNO, LAX, BUR, SAN, PSP, LAS, PHX, SLC, ABQ, ELP, COS, DEN, MAF — dominated by Pacific Coast, Mountain West, and Desert Southwest airports.\n- **Central (14 airports):** STL, MEM, ORD, MDW, MCI, ICT, LIT, BNA, DAL, HOU, MSY, JAN, PNS, VPS — the Mississippi/Gulf corridor and Midwest hubs.\n- **East (14 airports):** LGA, ISP, BWI, IAD, RIC, RDU, MYR, CMH, CLE, DTW, ATL, TPA, MCO, FLL — stretching from the Northeast down the Atlantic Seaboard to Florida.\n\nThe dominant pattern is **coast-to-coast routing**: each itinerary typically originates on one coast or the Gulf South, threads through 2–3 inland stops, and terminates on the opposite coast. Several routes end or begin at West Coast leisure-travel airports (OAK, LAX, SJC, LAS). Recurring stops like BWI, BNA, DAL, LAS, MDW, and OAK appear across multiple routes, reflecting Southwest's key connecting cities. No route is confined to a single region; all 10 span at least two of the three geographic bands.",
      "sql": "SELECT\n    multiIf(d.Longitude \u003c -100, 'West', d.Longitude \u003c -85, 'Central', 'East') AS region,\n    count() AS airport_count,\n    groupArray(a.code) AS airports\nFROM (\n    SELECT DISTINCT arrayJoin([\n        'ISP','BWI','MYR','BNA','VPS','DAL','LAS','OAK','SEA',\n        'CLE','PNS','HOU','MCI','PHX','BUR','DEN',\n        'ELP','LIT','ATL','RIC','MDW','SAN',\n        'MSY','CMH','RDU','DTW','LAX',\n        'TPA','ORD','SLC','SJC',\n        'LGA','STL','ICT','COS','JAN',\n        'SMF','PSP','RNO',\n        'ABQ','MAF',\n        'FLL','MCO','MEM','IAD'\n    ]) AS code\n) a\nLEFT JOIN ontime.dim_airports d ON a.code = d.AirportCode\nGROUP BY region\nORDER BY min(d.Longitude)"
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
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `queries/main.sql`
- `queries/q1.sql`
- `queries/q2.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/main.json`
- `results/q1.json`
- `results/q2.json`