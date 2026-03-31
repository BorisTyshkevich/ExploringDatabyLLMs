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

Find American Airlines' worst network-wide month for departure delays, then identify which origins and routes contributed most to that peak.

Analyze completed American Airlines flights by month across the full network. Find the single month that stands out as the worst overall for departure delays.

Use the full available history unless the question explicitly asks for a narrower period. Do not invent a custom score for the peak month; identify it directly from the monthly delay metrics needed to answer the question.

For the monthly view, quantify:

- flight volume
- average departure delay
- the share of flights departing 15+ minutes late

Then drill into the selected peak month to show which origin airports and origin-destination routes contributed most to that bad month. Focus on contributors with enough flights in that month to be meaningful.

For the breadth-versus-concentration question, compute concentration against the full selected-month AA network. If you apply minimum-flight thresholds to keep origin or route tables business-meaningful, use those thresholds only for displayed contributor tables, not as the denominator for any network-wide concentration share.

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a monthly leaderboard showing how the network performed over time
- a drilldown into the selected peak month by origin and by route

The proof query behind the peak-month question should preserve the month-by-month network rows needed for the dashboard, not just the single worst month.

The proof query behind the breadth/concentration question must let the reviewer inspect full-network peak-month totals and the share captured by top origins and top routes against those full-network totals. Do not answer "across the network" using shares computed only from filtered subsets.

## Dashboard Questions

- Which month is the single worst American Airlines month for departure delays?
- Which origins contribute most to that peak month?
- Which routes contribute most to that peak month?
- Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

In the report, answer those questions directly in prose. Name the worst month, identify the leading origin and route contributors using the verified result, and summarize whether the peak looks broad or concentrated.

If you use thresholded contributor tables for readability, label them as display filters only and keep the network-wide breadth/concentration conclusion tied to full-network denominators.

Do not use fallback phrases such as "the peak month" or "the leading contributors" when your verified query results let you name the actual month, origins, and routes directly.

Keep the result business-readable and analytically sound.

Generated report.md:

```md
# American Airlines peak network delay month and contributors

> Which month is the single worst American Airlines month for departure delays?

July 2024 is the single worst American Airlines month for departure delays in the full available history. It leads the AA monthly leaderboard on both average departure delay and the share of flights leaving 15 or more minutes late, with 86,083 completed flights, a 36.33 minute average departure delay, and 38.04% of flights departing 15+ minutes late.

- Rows returned: 458
- Columns: month, flights, avg_dep_delay, pct_dep_15_plus, worst_month_rank

| month | flights | avg_dep_delay | pct_dep_15_plus | worst_month_rank |
| --- | --- | --- | --- | --- |
| 1987-10-01T00:00:00Z | 55871 | 4.27 | 8.8 | 402 |

> Which origins contribute most to that peak month?

In July 2024, Dallas/Fort Worth (DFW) and Charlotte (CLT) are the dominant origin contributors. DFW produced 593,021 departure-delay minutes and 17.98% of the month's total AA delay minutes, while CLT added 560,409 minutes and 16.99%; ORD, MIA, and PHL are the next tier. The displayed origin table uses a 1,000-flight readability filter, but each airport's share is still measured against the full July 2024 AA network.

- Rows returned: 15
- Columns: OriginCode, origin_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, network_delay_share_pct, network_flight_share_pct

| OriginCode | origin_city | flights | avg_dep_delay | pct_dep_15_plus | total_dep_delay_minutes | network_delay_share_pct | network_flight_share_pct |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DFW | Dallas/Fort Worth, TX | 14962 | 38.52 | 46.59 | 593021 | 17.98 | 17.38 |

> Which routes contribute most to that peak month?

Among business-meaningful July 2024 AA routes, DFW-LAX is the largest route contributor by departure-delay minutes. CLT-MCO and DFW-SAT are close behind, followed by CLT-RDU and RDU-CLT. The displayed route table applies a 150-flight readability filter, but every route share is still computed against the full July 2024 AA network total.

- Rows returned: 182
- Columns: route, origin_city, dest_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, network_delay_share_pct, network_flight_share_pct

| route | origin_city | dest_city | flights | avg_dep_delay | pct_dep_15_plus | total_dep_delay_minutes | network_delay_share_pct | network_flight_share_pct |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DFW-LAX | Dallas/Fort Worth, TX | Los Angeles, CA | 443 | 40.89 | 46.5 | 18462 | 0.56 | 0.51 |

> Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

July 2024 looks broad across the AA network, with notable hub concentration rather than a single-point failure. DFW alone accounts for 17.98% of delay minutes, and the top five origins account for 50.75%, but the top route contributes only 0.56% and the top ten routes together only 4.88% of delay minutes across 909 routes. That pattern points to network-wide disruption centered heavily in major hubs, especially DFW and CLT, not a narrow problem confined to a few routes.

- Rows returned: 1
- Columns: network_flights, network_delay_minutes, top_origin, top_origin_delay_minutes, top_origin_delay_share_pct, top_origin_flights, top_origin_flight_share_pct, top_route, top_route_delay_minutes, top_route_delay_share_pct, top_route_flights, top_route_flight_share_pct, top5_origin_delay_minutes, top5_origin_delay_share_pct, top5_origin_flights, top5_origin_flight_share_pct, top10_route_delay_minutes, top10_route_delay_share_pct, top10_route_flights, top10_route_flight_share_pct, origin_count, route_count

| network_flights | network_delay_minutes | top_origin | top_origin_delay_minutes | top_origin_delay_share_pct | top_origin_flights | top_origin_flight_share_pct | top_route | top_route_delay_minutes | top_route_delay_share_pct | top_route_flights | top_route_flight_share_pct | top5_origin_delay_minutes | top5_origin_delay_share_pct | top5_origin_flights | top5_origin_flight_share_pct | top10_route_delay_minutes | top10_route_delay_share_pct | top10_route_flights | top10_route_flight_share_pct | origin_count | route_count |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 86083 | 3.297804e+06 | DFW | 593021 | 17.98 | 14962 | 17.38 | DFW-LAX | 18462 | 0.56 | 443 | 0.51 | 1.673625e+06 | 50.75 | 39810 | 46.25 | 161038 | 4.88 | 3340 | 3.88 | 125 | 909 |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "subquestion": "Which month is the single worst American Airlines month for departure delays?",
      "answer_markdown": "July 2024 is the single worst American Airlines month for departure delays in the full available history. It leads the AA monthly leaderboard on both average departure delay and the share of flights leaving 15 or more minutes late, with 86,083 completed flights, a 36.33 minute average departure delay, and 38.04% of flights departing 15+ minutes late.",
      "sql": "WITH monthly AS ( SELECT toStartOfMonth(FlightDate) AS month, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY month ) SELECT month, flights, avg_dep_delay, pct_dep_15_plus, row_number() OVER (ORDER BY avg_dep_delay DESC, pct_dep_15_plus DESC, flights DESC, month DESC) AS worst_month_rank FROM monthly ORDER BY month"
    },
    {
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "In July 2024, Dallas/Fort Worth (DFW) and Charlotte (CLT) are the dominant origin contributors. DFW produced 593,021 departure-delay minutes and 17.98% of the month's total AA delay minutes, while CLT added 560,409 minutes and 16.99%; ORD, MIA, and PHL are the next tier. The displayed origin table uses a 1,000-flight readability filter, but each airport's share is still measured against the full July 2024 AA network.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, OriginCityName, DepDelay, DepDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), origins AS ( SELECT OriginCode, any(OriginCityName) AS origin_city, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus, sum(greatest(DepDelay, 0)) AS total_dep_delay_minutes FROM base GROUP BY OriginCode ) SELECT OriginCode, origin_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, round(total_dep_delay_minutes / network_delay_minutes * 100, 2) AS network_delay_share_pct, round(flights / network_flights * 100, 2) AS network_flight_share_pct FROM origins CROSS JOIN network WHERE flights >= 1000 ORDER BY total_dep_delay_minutes DESC, flights DESC"
    },
    {
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "Among business-meaningful July 2024 AA routes, DFW-LAX is the largest route contributor by departure-delay minutes. CLT-MCO and DFW-SAT are close behind, followed by CLT-RDU and RDU-CLT. The displayed route table applies a 150-flight readability filter, but every route share is still computed against the full July 2024 AA network total.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, OriginCityName, DestCode, DestCityName, DepDelay, DepDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), routes AS ( SELECT concat(OriginCode, '-', DestCode) AS route, any(OriginCityName) AS origin_city, any(DestCityName) AS dest_city, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus, sum(greatest(DepDelay, 0)) AS total_dep_delay_minutes FROM base GROUP BY OriginCode, DestCode ) SELECT route, origin_city, dest_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, round(total_dep_delay_minutes / network_delay_minutes * 100, 2) AS network_delay_share_pct, round(flights / network_flights * 100, 2) AS network_flight_share_pct FROM routes CROSS JOIN network WHERE flights >= 150 ORDER BY total_dep_delay_minutes DESC, flights DESC"
    },
    {
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "July 2024 looks broad across the AA network, with notable hub concentration rather than a single-point failure. DFW alone accounts for 17.98% of delay minutes, and the top five origins account for 50.75%, but the top route contributes only 0.56% and the top ten routes together only 4.88% of delay minutes across 909 routes. That pattern points to network-wide disruption centered heavily in major hubs, especially DFW and CLT, not a narrow problem confined to a few routes.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, DestCode, DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), origin_ranked AS ( SELECT OriginCode, count() AS flights, sum(greatest(DepDelay, 0)) AS delay_minutes FROM base GROUP BY OriginCode ), route_ranked AS ( SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flights, sum(greatest(DepDelay, 0)) AS delay_minutes FROM base GROUP BY OriginCode, DestCode ), origin_top1 AS ( SELECT OriginCode AS top_origin, flights AS top_origin_flights, delay_minutes AS top_origin_delay_minutes FROM origin_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 1 ), route_top1 AS ( SELECT route AS top_route, flights AS top_route_flights, delay_minutes AS top_route_delay_minutes FROM route_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 1 ), top_origins AS ( SELECT sum(delay_minutes) AS top5_origin_delay_minutes, sum(flights) AS top5_origin_flights FROM (SELECT * FROM origin_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 5) ), top_routes AS ( SELECT sum(delay_minutes) AS top10_route_delay_minutes, sum(flights) AS top10_route_flights FROM (SELECT * FROM route_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 10) ) SELECT network_flights, network_delay_minutes, top_origin, top_origin_delay_minutes, round(top_origin_delay_minutes / network_delay_minutes * 100, 2) AS top_origin_delay_share_pct, top_origin_flights, round(top_origin_flights / network_flights * 100, 2) AS top_origin_flight_share_pct, top_route, top_route_delay_minutes, round(top_route_delay_minutes / network_delay_minutes * 100, 2) AS top_route_delay_share_pct, top_route_flights, round(top_route_flights / network_flights * 100, 2) AS top_route_flight_share_pct, top5_origin_delay_minutes, round(top5_origin_delay_minutes / network_delay_minutes * 100, 2) AS top5_origin_delay_share_pct, top5_origin_flights, round(top5_origin_flights / network_flights * 100, 2) AS top5_origin_flight_share_pct, top10_route_delay_minutes, round(top10_route_delay_minutes / network_delay_minutes * 100, 2) AS top10_route_delay_share_pct, top10_route_flights, round(top10_route_flights / network_flights * 100, 2) AS top10_route_flight_share_pct, (SELECT count() FROM origin_ranked) AS origin_count, (SELECT count() FROM route_ranked) AS route_count FROM network CROSS JOIN origin_top1 CROSS JOIN route_top1 CROSS JOIN top_origins CROSS JOIN top_routes"
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
      "answer_markdown": "July 2024 is the single worst American Airlines month for departure delays in the full available history. It leads the AA monthly leaderboard on both average departure delay and the share of flights leaving 15 or more minutes late, with 86,083 completed flights, a 36.33 minute average departure delay, and 38.04% of flights departing 15+ minutes late.",
      "sql": "WITH monthly AS ( SELECT toStartOfMonth(FlightDate) AS month, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY month ) SELECT month, flights, avg_dep_delay, pct_dep_15_plus, row_number() OVER (ORDER BY avg_dep_delay DESC, pct_dep_15_plus DESC, flights DESC, month DESC) AS worst_month_rank FROM monthly ORDER BY month"
    },
    {
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "In July 2024, Dallas/Fort Worth (DFW) and Charlotte (CLT) are the dominant origin contributors. DFW produced 593,021 departure-delay minutes and 17.98% of the month's total AA delay minutes, while CLT added 560,409 minutes and 16.99%; ORD, MIA, and PHL are the next tier. The displayed origin table uses a 1,000-flight readability filter, but each airport's share is still measured against the full July 2024 AA network.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, OriginCityName, DepDelay, DepDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), origins AS ( SELECT OriginCode, any(OriginCityName) AS origin_city, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus, sum(greatest(DepDelay, 0)) AS total_dep_delay_minutes FROM base GROUP BY OriginCode ) SELECT OriginCode, origin_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, round(total_dep_delay_minutes / network_delay_minutes * 100, 2) AS network_delay_share_pct, round(flights / network_flights * 100, 2) AS network_flight_share_pct FROM origins CROSS JOIN network WHERE flights \u003e= 1000 ORDER BY total_dep_delay_minutes DESC, flights DESC"
    },
    {
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "Among business-meaningful July 2024 AA routes, DFW-LAX is the largest route contributor by departure-delay minutes. CLT-MCO and DFW-SAT are close behind, followed by CLT-RDU and RDU-CLT. The displayed route table applies a 150-flight readability filter, but every route share is still computed against the full July 2024 AA network total.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, OriginCityName, DestCode, DestCityName, DepDelay, DepDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), routes AS ( SELECT concat(OriginCode, '-', DestCode) AS route, any(OriginCityName) AS origin_city, any(DestCityName) AS dest_city, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus, sum(greatest(DepDelay, 0)) AS total_dep_delay_minutes FROM base GROUP BY OriginCode, DestCode ) SELECT route, origin_city, dest_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, round(total_dep_delay_minutes / network_delay_minutes * 100, 2) AS network_delay_share_pct, round(flights / network_flights * 100, 2) AS network_flight_share_pct FROM routes CROSS JOIN network WHERE flights \u003e= 150 ORDER BY total_dep_delay_minutes DESC, flights DESC"
    },
    {
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "July 2024 looks broad across the AA network, with notable hub concentration rather than a single-point failure. DFW alone accounts for 17.98% of delay minutes, and the top five origins account for 50.75%, but the top route contributes only 0.56% and the top ten routes together only 4.88% of delay minutes across 909 routes. That pattern points to network-wide disruption centered heavily in major hubs, especially DFW and CLT, not a narrow problem confined to a few routes.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, DestCode, DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), origin_ranked AS ( SELECT OriginCode, count() AS flights, sum(greatest(DepDelay, 0)) AS delay_minutes FROM base GROUP BY OriginCode ), route_ranked AS ( SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flights, sum(greatest(DepDelay, 0)) AS delay_minutes FROM base GROUP BY OriginCode, DestCode ), origin_top1 AS ( SELECT OriginCode AS top_origin, flights AS top_origin_flights, delay_minutes AS top_origin_delay_minutes FROM origin_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 1 ), route_top1 AS ( SELECT route AS top_route, flights AS top_route_flights, delay_minutes AS top_route_delay_minutes FROM route_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 1 ), top_origins AS ( SELECT sum(delay_minutes) AS top5_origin_delay_minutes, sum(flights) AS top5_origin_flights FROM (SELECT * FROM origin_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 5) ), top_routes AS ( SELECT sum(delay_minutes) AS top10_route_delay_minutes, sum(flights) AS top10_route_flights FROM (SELECT * FROM route_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 10) ) SELECT network_flights, network_delay_minutes, top_origin, top_origin_delay_minutes, round(top_origin_delay_minutes / network_delay_minutes * 100, 2) AS top_origin_delay_share_pct, top_origin_flights, round(top_origin_flights / network_flights * 100, 2) AS top_origin_flight_share_pct, top_route, top_route_delay_minutes, round(top_route_delay_minutes / network_delay_minutes * 100, 2) AS top_route_delay_share_pct, top_route_flights, round(top_route_flights / network_flights * 100, 2) AS top_route_flight_share_pct, top5_origin_delay_minutes, round(top5_origin_delay_minutes / network_delay_minutes * 100, 2) AS top5_origin_delay_share_pct, top5_origin_flights, round(top5_origin_flights / network_flights * 100, 2) AS top5_origin_flight_share_pct, top10_route_delay_minutes, round(top10_route_delay_minutes / network_delay_minutes * 100, 2) AS top10_route_delay_share_pct, top10_route_flights, round(top10_route_flights / network_flights * 100, 2) AS top10_route_flight_share_pct, (SELECT count() FROM origin_ranked) AS origin_count, (SELECT count() FROM route_ranked) AS route_count FROM network CROSS JOIN origin_top1 CROSS JOIN route_top1 CROSS JOIN top_origins CROSS JOIN top_routes"
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
      "answer_markdown": "July 2024 is the single worst American Airlines month for departure delays in the full available history. It leads the AA monthly leaderboard on both average departure delay and the share of flights leaving 15 or more minutes late, with 86,083 completed flights, a 36.33 minute average departure delay, and 38.04% of flights departing 15+ minutes late.",
      "sql": "WITH monthly AS ( SELECT toStartOfMonth(FlightDate) AS month, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY month ) SELECT month, flights, avg_dep_delay, pct_dep_15_plus, row_number() OVER (ORDER BY avg_dep_delay DESC, pct_dep_15_plus DESC, flights DESC, month DESC) AS worst_month_rank FROM monthly ORDER BY month",
      "row_count": 458,
      "result_columns": [
        "month",
        "flights",
        "avg_dep_delay",
        "pct_dep_15_plus",
        "worst_month_rank"
      ],
      "first_row": {
        "avg_dep_delay": 4.27,
        "flights": 55871,
        "month": "1987-10-01T00:00:00Z",
        "pct_dep_15_plus": 8.8,
        "worst_month_rank": 402
      }
    },
    {
      "id": "q2",
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "In July 2024, Dallas/Fort Worth (DFW) and Charlotte (CLT) are the dominant origin contributors. DFW produced 593,021 departure-delay minutes and 17.98% of the month's total AA delay minutes, while CLT added 560,409 minutes and 16.99%; ORD, MIA, and PHL are the next tier. The displayed origin table uses a 1,000-flight readability filter, but each airport's share is still measured against the full July 2024 AA network.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, OriginCityName, DepDelay, DepDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), origins AS ( SELECT OriginCode, any(OriginCityName) AS origin_city, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus, sum(greatest(DepDelay, 0)) AS total_dep_delay_minutes FROM base GROUP BY OriginCode ) SELECT OriginCode, origin_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, round(total_dep_delay_minutes / network_delay_minutes * 100, 2) AS network_delay_share_pct, round(flights / network_flights * 100, 2) AS network_flight_share_pct FROM origins CROSS JOIN network WHERE flights \u003e= 1000 ORDER BY total_dep_delay_minutes DESC, flights DESC",
      "row_count": 15,
      "result_columns": [
        "OriginCode",
        "origin_city",
        "flights",
        "avg_dep_delay",
        "pct_dep_15_plus",
        "total_dep_delay_minutes",
        "network_delay_share_pct",
        "network_flight_share_pct"
      ],
      "first_row": {
        "OriginCode": "DFW",
        "avg_dep_delay": 38.52,
        "flights": 14962,
        "network_delay_share_pct": 17.98,
        "network_flight_share_pct": 17.38,
        "origin_city": "Dallas/Fort Worth, TX",
        "pct_dep_15_plus": 46.59,
        "total_dep_delay_minutes": 593021
      }
    },
    {
      "id": "q3",
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "Among business-meaningful July 2024 AA routes, DFW-LAX is the largest route contributor by departure-delay minutes. CLT-MCO and DFW-SAT are close behind, followed by CLT-RDU and RDU-CLT. The displayed route table applies a 150-flight readability filter, but every route share is still computed against the full July 2024 AA network total.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, OriginCityName, DestCode, DestCityName, DepDelay, DepDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), routes AS ( SELECT concat(OriginCode, '-', DestCode) AS route, any(OriginCityName) AS origin_city, any(DestCityName) AS dest_city, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus, sum(greatest(DepDelay, 0)) AS total_dep_delay_minutes FROM base GROUP BY OriginCode, DestCode ) SELECT route, origin_city, dest_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, round(total_dep_delay_minutes / network_delay_minutes * 100, 2) AS network_delay_share_pct, round(flights / network_flights * 100, 2) AS network_flight_share_pct FROM routes CROSS JOIN network WHERE flights \u003e= 150 ORDER BY total_dep_delay_minutes DESC, flights DESC",
      "row_count": 182,
      "result_columns": [
        "route",
        "origin_city",
        "dest_city",
        "flights",
        "avg_dep_delay",
        "pct_dep_15_plus",
        "total_dep_delay_minutes",
        "network_delay_share_pct",
        "network_flight_share_pct"
      ],
      "first_row": {
        "avg_dep_delay": 40.89,
        "dest_city": "Los Angeles, CA",
        "flights": 443,
        "network_delay_share_pct": 0.56,
        "network_flight_share_pct": 0.51,
        "origin_city": "Dallas/Fort Worth, TX",
        "pct_dep_15_plus": 46.5,
        "route": "DFW-LAX",
        "total_dep_delay_minutes": 18462
      }
    },
    {
      "id": "q4",
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "July 2024 looks broad across the AA network, with notable hub concentration rather than a single-point failure. DFW alone accounts for 17.98% of delay minutes, and the top five origins account for 50.75%, but the top route contributes only 0.56% and the top ten routes together only 4.88% of delay minutes across 909 routes. That pattern points to network-wide disruption centered heavily in major hubs, especially DFW and CLT, not a narrow problem confined to a few routes.",
      "sql": "WITH toDate('2024-07-01') AS peak_month, base AS ( SELECT OriginCode, DestCode, DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL AND toStartOfMonth(FlightDate) = peak_month ), network AS ( SELECT count() AS network_flights, sum(greatest(DepDelay, 0)) AS network_delay_minutes FROM base ), origin_ranked AS ( SELECT OriginCode, count() AS flights, sum(greatest(DepDelay, 0)) AS delay_minutes FROM base GROUP BY OriginCode ), route_ranked AS ( SELECT concat(OriginCode, '-', DestCode) AS route, count() AS flights, sum(greatest(DepDelay, 0)) AS delay_minutes FROM base GROUP BY OriginCode, DestCode ), origin_top1 AS ( SELECT OriginCode AS top_origin, flights AS top_origin_flights, delay_minutes AS top_origin_delay_minutes FROM origin_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 1 ), route_top1 AS ( SELECT route AS top_route, flights AS top_route_flights, delay_minutes AS top_route_delay_minutes FROM route_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 1 ), top_origins AS ( SELECT sum(delay_minutes) AS top5_origin_delay_minutes, sum(flights) AS top5_origin_flights FROM (SELECT * FROM origin_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 5) ), top_routes AS ( SELECT sum(delay_minutes) AS top10_route_delay_minutes, sum(flights) AS top10_route_flights FROM (SELECT * FROM route_ranked ORDER BY delay_minutes DESC, flights DESC LIMIT 10) ) SELECT network_flights, network_delay_minutes, top_origin, top_origin_delay_minutes, round(top_origin_delay_minutes / network_delay_minutes * 100, 2) AS top_origin_delay_share_pct, top_origin_flights, round(top_origin_flights / network_flights * 100, 2) AS top_origin_flight_share_pct, top_route, top_route_delay_minutes, round(top_route_delay_minutes / network_delay_minutes * 100, 2) AS top_route_delay_share_pct, top_route_flights, round(top_route_flights / network_flights * 100, 2) AS top_route_flight_share_pct, top5_origin_delay_minutes, round(top5_origin_delay_minutes / network_delay_minutes * 100, 2) AS top5_origin_delay_share_pct, top5_origin_flights, round(top5_origin_flights / network_flights * 100, 2) AS top5_origin_flight_share_pct, top10_route_delay_minutes, round(top10_route_delay_minutes / network_delay_minutes * 100, 2) AS top10_route_delay_share_pct, top10_route_flights, round(top10_route_flights / network_flights * 100, 2) AS top10_route_flight_share_pct, (SELECT count() FROM origin_ranked) AS origin_count, (SELECT count() FROM route_ranked) AS route_count FROM network CROSS JOIN origin_top1 CROSS JOIN route_top1 CROSS JOIN top_origins CROSS JOIN top_routes",
      "row_count": 1,
      "result_columns": [
        "network_flights",
        "network_delay_minutes",
        "top_origin",
        "top_origin_delay_minutes",
        "top_origin_delay_share_pct",
        "top_origin_flights",
        "top_origin_flight_share_pct",
        "top_route",
        "top_route_delay_minutes",
        "top_route_delay_share_pct",
        "top_route_flights",
        "top_route_flight_share_pct",
        "top5_origin_delay_minutes",
        "top5_origin_delay_share_pct",
        "top5_origin_flights",
        "top5_origin_flight_share_pct",
        "top10_route_delay_minutes",
        "top10_route_delay_share_pct",
        "top10_route_flights",
        "top10_route_flight_share_pct",
        "origin_count",
        "route_count"
      ],
      "first_row": {
        "network_delay_minutes": 3297804,
        "network_flights": 86083,
        "origin_count": 125,
        "route_count": 909,
        "top10_route_delay_minutes": 161038,
        "top10_route_delay_share_pct": 4.88,
        "top10_route_flight_share_pct": 3.88,
        "top10_route_flights": 3340,
        "top5_origin_delay_minutes": 1673625,
        "top5_origin_delay_share_pct": 50.75,
        "top5_origin_flight_share_pct": 46.25,
        "top5_origin_flights": 39810,
        "top_origin": "DFW",
        "top_origin_delay_minutes": 593021,
        "top_origin_delay_share_pct": 17.98,
        "top_origin_flight_share_pct": 17.38,
        "top_origin_flights": 14962,
        "top_route": "DFW-LAX",
        "top_route_delay_minutes": 18462,
        "top_route_delay_share_pct": 0.56,
        "top_route_flight_share_pct": 0.51,
        "top_route_flights": 443
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