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
- Do not suggest code changes to qforge in this review. Focus on the run and the prompt.
- Write Markdown only to `review.md`.

Question-specific guidance:

Find American Airlines' worst network-wide month for departure delays, then identify which origins and routes contributed most to that peak.

Analyze completed American Airlines flights by month across the full network. Find the single month that stands out as the worst overall for departure delays.

Use the full available history unless the question explicitly asks for a narrower period. Do not invent a custom score for the peak month; identify it directly from the monthly delay metrics needed to answer the question.

For the monthly view, quantify:

- flight volume
- average departure delay
- the share of flights departing 15+ minutes late

Then drill into the selected peak month to show which origin airports and origin-destination routes contributed most to that bad month. Focus on contributors with enough flights in that month to be meaningful.

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a monthly leaderboard showing how the network performed over time
- a drilldown into the selected peak month by origin and by route

The proof query behind the peak-month question should preserve the month-by-month network rows needed for the dashboard, not just the single worst month.

## Dashboard Questions

- Which month is the single worst American Airlines month for departure delays?
- Which origins contribute most to that peak month?
- Which routes contribute most to that peak month?
- Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

In the report, answer those questions directly in prose. Name the worst month, identify the leading origin and route contributors using the verified result, and summarize whether the peak looks broad or concentrated.

Do not use fallback phrases such as "the peak month" or "the leading contributors" when your verified query results let you name the actual month, origins, and routes directly.

Keep the result business-readable and analytically sound.

Generated report.md:

```md
# American Airlines peak network delay month and contributors

> Which month is the single worst American Airlines month for departure delays?

July 2024 is the single worst American Airlines month for departure delays in the full history. It had 86,083 completed AA flights, an average departure delay of 36.33 minutes, and 38.04% of flights departed at least 15 minutes late.

- Rows returned: 458
- Columns: month, flight_volume, avg_departure_delay_minutes, pct_departing_15_plus_late, worst_month_rank

| month | flight_volume | avg_departure_delay_minutes | pct_departing_15_plus_late | worst_month_rank |
| --- | --- | --- | --- | --- |
| 1987-10-01T00:00:00Z | 55871 | 4.27 | 8.8 | 402 |

> Which origins contribute most to that peak month?

In July 2024, DFW and CLT were the dominant origin contributors. DFW generated the most total departure delay minutes, while CLT was close behind and had the worse average delay rate; ORD, MIA, and PHL formed the next tier of meaningful contributors.

- Rows returned: 15
- Columns: OriginCode, flight_volume, total_departure_delay_minutes, avg_departure_delay_minutes, pct_departing_15_plus_late

| OriginCode | flight_volume | total_departure_delay_minutes | avg_departure_delay_minutes | pct_departing_15_plus_late |
| --- | --- | --- | --- | --- |
| DFW | 14962 | 593021 | 38.52 | 46.59 |

> Which routes contribute most to that peak month?

The largest July 2024 route contributor was DFW-LAX, followed very closely by CLT-MCO and DFW-SAT. CLT-RDU, RDU-CLT, and DFW-PHX also ranked among the strongest route-level contributors to the bad month.

- Rows returned: 20
- Columns: route, flight_volume, total_departure_delay_minutes, avg_departure_delay_minutes, pct_departing_15_plus_late

| route | flight_volume | total_departure_delay_minutes | avg_departure_delay_minutes | pct_departing_15_plus_late |
| --- | --- | --- | --- | --- |
| DFW-LAX | 443 | 18462 | 40.89 | 46.5 |

> Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

July 2024 looks broad across the network, but with clear hub concentration. AA had delays spread across 125 origins and 909 routes that month; the top two origins, DFW and CLT, accounted for 50.8% of meaningful origin delay minutes and the top five origins reached 73.71%, while the top ten routes made up only 15.82% of meaningful route delay minutes.

- Rows returned: 1
- Columns: total_origins_in_month, total_routes_in_month, top_2_origins_share_pct, top_5_origins_share_pct, top_10_routes_share_pct

| total_origins_in_month | total_routes_in_month | top_2_origins_share_pct | top_5_origins_share_pct | top_10_routes_share_pct |
| --- | --- | --- | --- | --- |
| 125 | 909 | 50.8 | 73.71 | 15.82 |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "subquestion": "Which month is the single worst American Airlines month for departure delays?",
      "answer_markdown": "July 2024 is the single worst American Airlines month for departure delays in the full history. It had 86,083 completed AA flights, an average departure delay of 36.33 minutes, and 38.04% of flights departed at least 15 minutes late.",
      "sql": "WITH monthly AS ( SELECT toStartOfMonth(FlightDate) AS month, count() AS flight_volume, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY month ) SELECT month, flight_volume, avg_departure_delay_minutes, pct_departing_15_plus_late, dense_rank() OVER (ORDER BY avg_departure_delay_minutes DESC, pct_departing_15_plus_late DESC, flight_volume DESC) AS worst_month_rank FROM monthly ORDER BY month"
    },
    {
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "In July 2024, DFW and CLT were the dominant origin contributors. DFW generated the most total departure delay minutes, while CLT was close behind and had the worse average delay rate; ORD, MIA, and PHL formed the next tier of meaningful contributors.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume >= 1000 ORDER BY total_departure_delay_minutes DESC, avg_departure_delay_minutes DESC LIMIT 20"
    },
    {
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "The largest July 2024 route contributor was DFW-LAX, followed very closely by CLT-MCO and DFW-SAT. CLT-RDU, RDU-CLT, and DFW-PHX also ranked among the strongest route-level contributors to the bad month.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY route HAVING flight_volume >= 200 ORDER BY total_departure_delay_minutes DESC, avg_departure_delay_minutes DESC LIMIT 20"
    },
    {
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "July 2024 looks broad across the network, but with clear hub concentration. AA had delays spread across 125 origins and 909 routes that month; the top two origins, DFW and CLT, accounted for 50.8% of meaningful origin delay minutes and the top five origins reached 73.71%, while the top ten routes made up only 15.82% of meaningful route delay minutes.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT ( SELECT uniqExact(OriginCode) FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ) AS total_origins_in_month, ( SELECT uniqExact(concat(OriginCode, '-', DestCode)) FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ) AS total_routes_in_month, ( SELECT round(sumIf(total_departure_delay_minutes, rn <= 2) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume >= 1000 ) ) ) AS top_2_origins_share_pct, ( SELECT round(sumIf(total_departure_delay_minutes, rn <= 5) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume >= 1000 ) ) ) AS top_5_origins_share_pct, ( SELECT round(sumIf(total_departure_delay_minutes, rn <= 10) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY route HAVING flight_volume >= 200 ) ) ) AS top_10_routes_share_pct"
    }
  ]
}
```

Saved analysis.json:

```json
{
  "sql": "",
  "report_markdown": "",
  "subquestions": [
    {
      "subquestion": "Which month is the single worst American Airlines month for departure delays?",
      "answer_markdown": "July 2024 is the single worst American Airlines month for departure delays in the full history. It had 86,083 completed AA flights, an average departure delay of 36.33 minutes, and 38.04% of flights departed at least 15 minutes late.",
      "sql": "WITH monthly AS ( SELECT toStartOfMonth(FlightDate) AS month, count() AS flight_volume, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY month ) SELECT month, flight_volume, avg_departure_delay_minutes, pct_departing_15_plus_late, dense_rank() OVER (ORDER BY avg_departure_delay_minutes DESC, pct_departing_15_plus_late DESC, flight_volume DESC) AS worst_month_rank FROM monthly ORDER BY month"
    },
    {
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "In July 2024, DFW and CLT were the dominant origin contributors. DFW generated the most total departure delay minutes, while CLT was close behind and had the worse average delay rate; ORD, MIA, and PHL formed the next tier of meaningful contributors.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume \u003e= 1000 ORDER BY total_departure_delay_minutes DESC, avg_departure_delay_minutes DESC LIMIT 20"
    },
    {
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "The largest July 2024 route contributor was DFW-LAX, followed very closely by CLT-MCO and DFW-SAT. CLT-RDU, RDU-CLT, and DFW-PHX also ranked among the strongest route-level contributors to the bad month.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY route HAVING flight_volume \u003e= 200 ORDER BY total_departure_delay_minutes DESC, avg_departure_delay_minutes DESC LIMIT 20"
    },
    {
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "July 2024 looks broad across the network, but with clear hub concentration. AA had delays spread across 125 origins and 909 routes that month; the top two origins, DFW and CLT, accounted for 50.8% of meaningful origin delay minutes and the top five origins reached 73.71%, while the top ten routes made up only 15.82% of meaningful route delay minutes.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT ( SELECT uniqExact(OriginCode) FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ) AS total_origins_in_month, ( SELECT uniqExact(concat(OriginCode, '-', DestCode)) FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ) AS total_routes_in_month, ( SELECT round(sumIf(total_departure_delay_minutes, rn \u003c= 2) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume \u003e= 1000 ) ) ) AS top_2_origins_share_pct, ( SELECT round(sumIf(total_departure_delay_minutes, rn \u003c= 5) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume \u003e= 1000 ) ) ) AS top_5_origins_share_pct, ( SELECT round(sumIf(total_departure_delay_minutes, rn \u003c= 10) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY route HAVING flight_volume \u003e= 200 ) ) ) AS top_10_routes_share_pct"
    }
  ]
}
```

Saved visual_input.json:

```json
{
  "question_title": "American Airlines peak network delay month and contributors",
  "result_columns": null,
  "row_count": 4,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which month is the single worst American Airlines month for departure delays?",
      "answer_markdown": "July 2024 is the single worst American Airlines month for departure delays in the full history. It had 86,083 completed AA flights, an average departure delay of 36.33 minutes, and 38.04% of flights departed at least 15 minutes late.",
      "sql": "WITH monthly AS ( SELECT toStartOfMonth(FlightDate) AS month, count() AS flight_volume, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY month ) SELECT month, flight_volume, avg_departure_delay_minutes, pct_departing_15_plus_late, dense_rank() OVER (ORDER BY avg_departure_delay_minutes DESC, pct_departing_15_plus_late DESC, flight_volume DESC) AS worst_month_rank FROM monthly ORDER BY month",
      "row_count": 458,
      "result_columns": [
        "month",
        "flight_volume",
        "avg_departure_delay_minutes",
        "pct_departing_15_plus_late",
        "worst_month_rank"
      ],
      "first_row": {
        "avg_departure_delay_minutes": 4.27,
        "flight_volume": 55871,
        "month": "1987-10-01T00:00:00Z",
        "pct_departing_15_plus_late": 8.8,
        "worst_month_rank": 402
      }
    },
    {
      "id": "q2",
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "In July 2024, DFW and CLT were the dominant origin contributors. DFW generated the most total departure delay minutes, while CLT was close behind and had the worse average delay rate; ORD, MIA, and PHL formed the next tier of meaningful contributors.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume \u003e= 1000 ORDER BY total_departure_delay_minutes DESC, avg_departure_delay_minutes DESC LIMIT 20",
      "row_count": 15,
      "result_columns": [
        "OriginCode",
        "flight_volume",
        "total_departure_delay_minutes",
        "avg_departure_delay_minutes",
        "pct_departing_15_plus_late"
      ],
      "first_row": {
        "OriginCode": "DFW",
        "avg_departure_delay_minutes": 38.52,
        "flight_volume": 14962,
        "pct_departing_15_plus_late": 46.59,
        "total_departure_delay_minutes": 593021
      }
    },
    {
      "id": "q3",
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "The largest July 2024 route contributor was DFW-LAX, followed very closely by CLT-MCO and DFW-SAT. CLT-RDU, RDU-CLT, and DFW-PHX also ranked among the strongest route-level contributors to the bad month.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes, round(avg(DepDelay), 2) AS avg_departure_delay_minutes, round(avg(DepDel15) * 100, 2) AS pct_departing_15_plus_late FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY route HAVING flight_volume \u003e= 200 ORDER BY total_departure_delay_minutes DESC, avg_departure_delay_minutes DESC LIMIT 20",
      "row_count": 20,
      "result_columns": [
        "route",
        "flight_volume",
        "total_departure_delay_minutes",
        "avg_departure_delay_minutes",
        "pct_departing_15_plus_late"
      ],
      "first_row": {
        "avg_departure_delay_minutes": 40.89,
        "flight_volume": 443,
        "pct_departing_15_plus_late": 46.5,
        "route": "DFW-LAX",
        "total_departure_delay_minutes": 18462
      }
    },
    {
      "id": "q4",
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "July 2024 looks broad across the network, but with clear hub concentration. AA had delays spread across 125 origins and 909 routes that month; the top two origins, DFW and CLT, accounted for 50.8% of meaningful origin delay minutes and the top five origins reached 73.71%, while the top ten routes made up only 15.82% of meaningful route delay minutes.",
      "sql": "WITH toDate('2024-07-01') AS peak_month SELECT ( SELECT uniqExact(OriginCode) FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ) AS total_origins_in_month, ( SELECT uniqExact(concat(OriginCode, '-', DestCode)) FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ) AS total_routes_in_month, ( SELECT round(sumIf(total_departure_delay_minutes, rn \u003c= 2) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume \u003e= 1000 ) ) ) AS top_2_origins_share_pct, ( SELECT round(sumIf(total_departure_delay_minutes, rn \u003c= 5) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT OriginCode, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY OriginCode HAVING flight_volume \u003e= 1000 ) ) ) AS top_5_origins_share_pct, ( SELECT round(sumIf(total_departure_delay_minutes, rn \u003c= 10) * 100.0 / max(total_delay_minutes), 2) FROM ( SELECT total_departure_delay_minutes, row_number() OVER (ORDER BY total_departure_delay_minutes DESC) AS rn, sum(total_departure_delay_minutes) OVER () AS total_delay_minutes FROM ( SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flight_volume, sum(DepDelayMinutes) AS total_departure_delay_minutes FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month GROUP BY route HAVING flight_volume \u003e= 200 ) ) ) AS top_10_routes_share_pct",
      "row_count": 1,
      "result_columns": [
        "total_origins_in_month",
        "total_routes_in_month",
        "top_2_origins_share_pct",
        "top_5_origins_share_pct",
        "top_10_routes_share_pct"
      ],
      "first_row": {
        "top_10_routes_share_pct": 15.82,
        "top_2_origins_share_pct": 50.8,
        "top_5_origins_share_pct": 73.71,
        "total_origins_in_month": 125,
        "total_routes_in_month": 909
      }
    }
  ]
}
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `queries/q1.sql`
- `queries/q2.sql`
- `queries/q3.sql`
- `queries/q4.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/q1.json`
- `results/q2.json`
- `results/q3.json`
- `results/q4.json`