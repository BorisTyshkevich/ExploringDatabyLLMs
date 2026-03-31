- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

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
- Set `Verdict: FAIL` when any substantive correctness, grain, metric, or evidence-support problem exists.
- Base your judgment only on the provided question prompt and the verified run artifacts.
- When artifact file paths are listed, read the files directly from the run directory instead of relying only on prompt excerpts.
- Check whether all dashboard questions were answered directly.
- Check whether SQL grain and returned metrics match the prompt.
- Check whether prose claims are supported by the executed query results.
- Check for duplicated entities, inconsistent counts, missing requested metrics, unsupported inference, and proof-query/result mismatches.
- Treat `report.md` as a monitoring artifact. Do not treat its one-row example table or abbreviated preview as a defect when the underlying query result preserves the required ranked or detailed rows.
- In `## Findings`, cite concrete artifact names such as `queries/q1.sql`, `results/q1.json`, or `report.md`.
- In `## Suggested Prompt Fixes`, propose prompt-level changes only when they would reduce the observed failure mode.
- When proposing a prompt fix, be concrete. Prefer 1-3 short replacement or insertion snippets that could be added to the question prompt, not vague advice.
- Tie each prompt fix to the specific failure you found. Explain what ambiguity in the current prompt allowed the bad result and how the new wording would close that gap.
- If the prompt is already sufficiently specific and the failure is not plausibly prompt-driven, say `None.` instead of inventing weak prompt edits.
- Do not suggest code changes to qforge in this review. Focus on the run and the prompt.
- Write Markdown only to `review.md`.

Question-specific guidance:

Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.

Return enough itinerary detail to support the dashboard, including:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`


## Dashboard Questions

- Which itinerary is the highest-hop example, and what does it look like?
- Which of the top-ranked itineraries is the most recent?
- Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

Generated report.md:

```md
# Highest daily hops for one aircraft on one flight number

> Which itinerary is the highest-hop example, and what does it look like?

The highest-hop example is an 8-hop Southwest itinerary flown by aircraft N957WN as flight 366 on 2024-12-01, with route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.

- Rows returned: 1
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | Route |
| --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Which of the top-ranked itineraries is the most recent?

The most recent top-ranked itinerary is also the latest 8-hop unique route: N957WN operating Southwest flight 366 on 2024-12-01 over ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.

- Rows returned: 1
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | Route |
| --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

They look mostly like recurring patterns rather than one-offs. Among the 10 most recent unique 8-hop routes, 9 recur on multiple dates, averaging 14.1 occurrences each, although the single most recent example appears only once.

- Rows returned: 1
- Columns: route_count, one_off_routes, recurring_routes, avg_occurrences, max_occurrences

| route_count | one_off_routes | recurring_routes | avg_occurrences | max_occurrences |
| --- | --- | --- | --- | --- |
| 10 | 1 | 9 | 14.1 | 46 |
```

Saved main.sql:

```sql
WITH itineraries AS (
    SELECT
        FlightDate AS flight_date,
        TailNum AS aircraft_id,
        FlightNum AS flight_number,
        Carrier AS carrier,
        arraySort(
            x -> (x.1, x.2, x.3),
            groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))
        ) AS legs
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND TailNum != ''
      AND FlightNum != ''
    GROUP BY
        flight_date,
        aircraft_id,
        flight_number,
        carrier
    HAVING count() > 1
),
ranked AS (
    SELECT
        flight_date,
        aircraft_id,
        flight_number,
        carrier,
        length(legs) AS hop_count,
        arrayStringConcat(arrayMap(x -> x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route
    FROM itineraries
),
max_hops AS (
    SELECT max(hop_count) AS max_hop_count
    FROM ranked
),
unique_routes AS (
    SELECT
        flight_date,
        aircraft_id,
        flight_number,
        carrier,
        hop_count,
        Route,
        row_number() OVER (
            PARTITION BY Route
            ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC
        ) AS route_rank
    FROM ranked
    WHERE hop_count = (SELECT max_hop_count FROM max_hops)
)
SELECT
    aircraft_id,
    flight_number,
    carrier,
    flight_date,
    hop_count,
    Route
FROM unique_routes
WHERE route_rank = 1
ORDER BY flight_date DESC, Route ASC
LIMIT 10
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "The highest-hop example is an 8-hop Southwest itinerary flown by aircraft N957WN as flight 366 on 2024-12-01, with route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -> (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() > 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -> x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1"
    },
    {
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is also the latest 8-hop unique route: N957WN operating Southwest flight 366 on 2024-12-01 over ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -> (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() > 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -> x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1"
    },
    {
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "They look mostly like recurring patterns rather than one-offs. Among the 10 most recent unique 8-hop routes, 9 recur on multiple dates, averaging 14.1 occurrences each, although the single most recent example appears only once.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -> (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() > 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -> x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ), top10 AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 10 ), route_occurrences AS ( SELECT t.Route, count() AS occurrences FROM top10 t INNER JOIN ranked r ON r.Route = t.Route AND r.hop_count = t.hop_count GROUP BY t.Route ) SELECT count() AS route_count, countIf(occurrences = 1) AS one_off_routes, countIf(occurrences > 1) AS recurring_routes, round(avg(occurrences), 2) AS avg_occurrences, max(occurrences) AS max_occurrences FROM route_occurrences"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH itineraries AS (\n    SELECT\n        FlightDate AS flight_date,\n        TailNum AS aircraft_id,\n        FlightNum AS flight_number,\n        Carrier AS carrier,\n        arraySort(\n            x -\u003e (x.1, x.2, x.3),\n            groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))\n        ) AS legs\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n      AND TailNum != ''\n      AND FlightNum != ''\n    GROUP BY\n        flight_date,\n        aircraft_id,\n        flight_number,\n        carrier\n    HAVING count() \u003e 1\n),\nranked AS (\n    SELECT\n        flight_date,\n        aircraft_id,\n        flight_number,\n        carrier,\n        length(legs) AS hop_count,\n        arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route\n    FROM itineraries\n),\nmax_hops AS (\n    SELECT max(hop_count) AS max_hop_count\n    FROM ranked\n),\nunique_routes AS (\n    SELECT\n        flight_date,\n        aircraft_id,\n        flight_number,\n        carrier,\n        hop_count,\n        Route,\n        row_number() OVER (\n            PARTITION BY Route\n            ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC\n        ) AS route_rank\n    FROM ranked\n    WHERE hop_count = (SELECT max_hop_count FROM max_hops)\n)\nSELECT\n    aircraft_id,\n    flight_number,\n    carrier,\n    flight_date,\n    hop_count,\n    Route\nFROM unique_routes\nWHERE route_rank = 1\nORDER BY flight_date DESC, Route ASC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "The highest-hop example is an 8-hop Southwest itinerary flown by aircraft N957WN as flight 366 on 2024-12-01, with route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1"
    },
    {
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is also the latest 8-hop unique route: N957WN operating Southwest flight 366 on 2024-12-01 over ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1"
    },
    {
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "They look mostly like recurring patterns rather than one-offs. Among the 10 most recent unique 8-hop routes, 9 recur on multiple dates, averaging 14.1 occurrences each, although the single most recent example appears only once.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ), top10 AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 10 ), route_occurrences AS ( SELECT t.Route, count() AS occurrences FROM top10 t INNER JOIN ranked r ON r.Route = t.Route AND r.hop_count = t.hop_count GROUP BY t.Route ) SELECT count() AS route_count, countIf(occurrences = 1) AS one_off_routes, countIf(occurrences \u003e 1) AS recurring_routes, round(avg(occurrences), 2) AS avg_occurrences, max(occurrences) AS max_occurrences FROM route_occurrences"
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
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "The highest-hop example is an 8-hop Southwest itinerary flown by aircraft N957WN as flight 366 on 2024-12-01, with route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1",
      "row_count": 1,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is also the latest 8-hop unique route: N957WN operating Southwest flight 366 on 2024-12-01 over ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ) SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 1",
      "row_count": 1,
      "result_columns": [
        "aircraft_id",
        "flight_number",
        "carrier",
        "flight_date",
        "hop_count",
        "Route"
      ],
      "first_row": {
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "aircraft_id": "N957WN",
        "carrier": "WN",
        "flight_date": "2024-12-01T00:00:00Z",
        "flight_number": "366",
        "hop_count": 8
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "They look mostly like recurring patterns rather than one-offs. Among the 10 most recent unique 8-hop routes, 9 recur on multiple dates, averaging 14.1 occurrences each, although the single most recent example appears only once.",
      "sql": "WITH itineraries AS ( SELECT FlightDate AS flight_date, TailNum AS aircraft_id, FlightNum AS flight_number, Carrier AS carrier, arraySort(x -\u003e (x.1, x.2, x.3), groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))) AS legs FROM ontime.fact_ontime WHERE Cancelled = 0 AND TailNum != '' AND FlightNum != '' GROUP BY flight_date, aircraft_id, flight_number, carrier HAVING count() \u003e 1 ), ranked AS ( SELECT flight_date, aircraft_id, flight_number, carrier, length(legs) AS hop_count, arrayStringConcat(arrayMap(x -\u003e x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route FROM itineraries ), max_hops AS ( SELECT max(hop_count) AS max_hop_count FROM ranked ), unique_routes AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route, row_number() OVER (PARTITION BY Route ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rank FROM ranked WHERE hop_count = (SELECT max_hop_count FROM max_hops) ), top10 AS ( SELECT flight_date, aircraft_id, flight_number, carrier, hop_count, Route FROM unique_routes WHERE route_rank = 1 ORDER BY flight_date DESC, Route ASC LIMIT 10 ), route_occurrences AS ( SELECT t.Route, count() AS occurrences FROM top10 t INNER JOIN ranked r ON r.Route = t.Route AND r.hop_count = t.hop_count GROUP BY t.Route ) SELECT count() AS route_count, countIf(occurrences = 1) AS one_off_routes, countIf(occurrences \u003e 1) AS recurring_routes, round(avg(occurrences), 2) AS avg_occurrences, max(occurrences) AS max_occurrences FROM route_occurrences",
      "row_count": 1,
      "result_columns": [
        "route_count",
        "one_off_routes",
        "recurring_routes",
        "avg_occurrences",
        "max_occurrences"
      ],
      "first_row": {
        "avg_occurrences": 14.1,
        "max_occurrences": 46,
        "one_off_routes": 1,
        "recurring_routes": 9,
        "route_count": 10
      }
    }
  ]
}
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `main.sql`
- `queries/q1.sql`
- `queries/q2.sql`
- `queries/q3.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/main.json`
- `results/q1.json`
- `results/q2.json`
- `results/q3.json`