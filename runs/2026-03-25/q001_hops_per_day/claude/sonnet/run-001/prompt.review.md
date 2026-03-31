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

All top-10 unique routes tie at **8 hops** (9 airports). The leading example by recency is Southwest flight **WN 366**, tail **N957WN**, on **2024-12-01**: `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` — a coast-to-coast turn starting in Islip, NY and ending in Seattle.

- Rows returned: 1
- Columns: Tail_Number, FlightNum, Carrier, FlightDate, hop_count, Route, occurrences

| Tail_Number | FlightNum | Carrier | FlightDate | hop_count | Route | occurrences |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |

> Which of the top-ranked itineraries is the most recent?

The most recent top-ranked itinerary is **WN 366** (tail N957WN) on **2024-12-01**, route `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` with 8 hops. It is the only 8-hop unique route recorded in late 2024.

- Rows returned: 1
- Columns: Tail_Number, FlightNum, Carrier, FlightDate, hop_count, Route, occurrences

| Tail_Number | FlightNum | Carrier | FlightDate | hop_count | Route | occurrences |
| --- | --- | --- | --- | --- | --- | --- |
| N925WN | 2179 | WN | 2025-11-30T00:00:00Z | 6 | SAT-TPA-FLL-RDU-BNA-STL-AUS | 1 |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

The top 10 are a **mix of both**. Two routes are clearly recurring scheduled turn patterns: WN 2884 (`LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN`) appeared **46 times** and WN 1956 (`MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC`) appeared **40 times**, indicating they were regular daily assignments over many years. Others appear only 1–5 times, suggesting irregular or one-off operations. Southwest dominates all 10 entries, consistent with its high-utilization short-haul turn model.

- Rows returned: 10
- Columns: FlightNum, Carrier, FlightDate, hop_count, Route, occurrences

| FlightNum | Carrier | FlightDate | hop_count | Route | occurrences |
| --- | --- | --- | --- | --- | --- |
| 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |
```

Saved main.sql:

```sql
WITH daily_itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        Reporting_Airline                AS Carrier,
        FlightDate,
        count()                          AS hop_count,
        arrayStringConcat(
            arrayConcat(
                [arrayElement(arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],
                arrayMap(x -> x.3, arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))
            ),
            '-'
        ) AS Route
    FROM ontime.fact_ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate
    HAVING hop_count >= 2
),
unique_routes AS (
    SELECT
        Route,
        max(hop_count)                   AS hop_count,
        argMax(Tail_Number, FlightDate)  AS Tail_Number,
        argMax(FlightNum,   FlightDate)  AS FlightNum,
        argMax(Carrier,     FlightDate)  AS Carrier,
        max(FlightDate)                  AS most_recent_date,
        count()                          AS occurrences
    FROM daily_itineraries
    GROUP BY Route
)
SELECT
    Tail_Number,
    FlightNum,
    Carrier,
    most_recent_date AS FlightDate,
    hop_count,
    Route,
    occurrences
FROM unique_routes
ORDER BY hop_count DESC, most_recent_date DESC
LIMIT 10
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "All top-10 unique routes tie at **8 hops** (9 airports). The leading example by recency is Southwest flight **WN 366**, tail **N957WN**, on **2024-12-01**: `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` — a coast-to-coast turn starting in Islip, NY and ending in Seattle.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -> x.3, arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count >= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 1"
    },
    {
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN 366** (tail N957WN) on **2024-12-01**, route `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` with 8 hops. It is the only 8-hop unique route recorded in late 2024.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -> x.3, arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count >= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY most_recent_date DESC, hop_count DESC\nLIMIT 1"
    },
    {
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top 10 are a **mix of both**. Two routes are clearly recurring scheduled turn patterns: WN 2884 (`LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN`) appeared **46 times** and WN 1956 (`MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC`) appeared **40 times**, indicating they were regular daily assignments over many years. Others appear only 1–5 times, suggesting irregular or one-off operations. Southwest dominates all 10 entries, consistent with its high-utilization short-haul turn model.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -> x.3, arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count >= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 10"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline                AS Carrier,\n        FlightDate,\n        count()                          AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                   AS hop_count,\n        argMax(Tail_Number, FlightDate)  AS Tail_Number,\n        argMax(FlightNum,   FlightDate)  AS FlightNum,\n        argMax(Carrier,     FlightDate)  AS Carrier,\n        max(FlightDate)                  AS most_recent_date,\n        count()                          AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT\n    Tail_Number,\n    FlightNum,\n    Carrier,\n    most_recent_date AS FlightDate,\n    hop_count,\n    Route,\n    occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 10",
  "report_markdown": "",
  "subquestions": [
    {
      "subquestion": "Which itinerary is the highest-hop example, and what does it look like?",
      "answer_markdown": "All top-10 unique routes tie at **8 hops** (9 airports). The leading example by recency is Southwest flight **WN 366**, tail **N957WN**, on **2024-12-01**: `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` — a coast-to-coast turn starting in Islip, NY and ending in Seattle.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 1"
    },
    {
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN 366** (tail N957WN) on **2024-12-01**, route `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` with 8 hops. It is the only 8-hop unique route recorded in late 2024.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY most_recent_date DESC, hop_count DESC\nLIMIT 1"
    },
    {
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top 10 are a **mix of both**. Two routes are clearly recurring scheduled turn patterns: WN 2884 (`LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN`) appeared **46 times** and WN 1956 (`MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC`) appeared **40 times**, indicating they were regular daily assignments over many years. Others appear only 1–5 times, suggesting irregular or one-off operations. Southwest dominates all 10 entries, consistent with its high-utilization short-haul turn model.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 10"
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
      "answer_markdown": "All top-10 unique routes tie at **8 hops** (9 airports). The leading example by recency is Southwest flight **WN 366**, tail **N957WN**, on **2024-12-01**: `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` — a coast-to-coast turn starting in Islip, NY and ending in Seattle.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "Tail_Number": "N957WN",
        "hop_count": 8,
        "occurrences": 1
      }
    },
    {
      "id": "q2",
      "subquestion": "Which of the top-ranked itineraries is the most recent?",
      "answer_markdown": "The most recent top-ranked itinerary is **WN 366** (tail N957WN) on **2024-12-01**, route `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` with 8 hops. It is the only 8-hop unique route recorded in late 2024.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY most_recent_date DESC, hop_count DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Tail_Number",
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2025-11-30T00:00:00Z",
        "FlightNum": "2179",
        "Route": "SAT-TPA-FLL-RDU-BNA-STL-AUS",
        "Tail_Number": "N925WN",
        "hop_count": 6,
        "occurrences": 1
      }
    },
    {
      "id": "q3",
      "subquestion": "Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?",
      "answer_markdown": "The top 10 are a **mix of both**. Two routes are clearly recurring scheduled turn patterns: WN 2884 (`LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN`) appeared **46 times** and WN 1956 (`MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC`) appeared **40 times**, indicating they were regular daily assignments over many years. Others appear only 1–5 times, suggesting irregular or one-off operations. Southwest dominates all 10 entries, consistent with its high-utilization short-haul turn model.",
      "sql": "WITH daily_itineraries AS (\n    SELECT\n        Tail_Number,\n        Flight_Number_Reporting_Airline AS FlightNum,\n        Reporting_Airline               AS Carrier,\n        FlightDate,\n        count()                         AS hop_count,\n        arrayStringConcat(\n            arrayConcat(\n                [arrayElement(arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],\n                arrayMap(x -\u003e x.3, arraySort(x -\u003e x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))\n            ),\n            '-'\n        ) AS Route\n    FROM ontime.fact_ontime\n    WHERE Tail_Number != ''\n      AND Flight_Number_Reporting_Airline != ''\n      AND Cancelled = 0\n    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate\n    HAVING hop_count \u003e= 2\n),\nunique_routes AS (\n    SELECT\n        Route,\n        max(hop_count)                  AS hop_count,\n        argMax(Tail_Number, FlightDate) AS Tail_Number,\n        argMax(FlightNum,   FlightDate) AS FlightNum,\n        argMax(Carrier,     FlightDate) AS Carrier,\n        max(FlightDate)                 AS most_recent_date,\n        count()                         AS occurrences\n    FROM daily_itineraries\n    GROUP BY Route\n)\nSELECT FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences\nFROM unique_routes\nORDER BY hop_count DESC, most_recent_date DESC\nLIMIT 10",
      "row_count": 10,
      "result_columns": [
        "FlightNum",
        "Carrier",
        "FlightDate",
        "hop_count",
        "Route",
        "occurrences"
      ],
      "first_row": {
        "Carrier": "WN",
        "FlightDate": "2024-12-01T00:00:00Z",
        "FlightNum": "366",
        "Route": "ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA",
        "hop_count": 8,
        "occurrences": 1
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