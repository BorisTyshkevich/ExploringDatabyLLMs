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

July 2024 is the single worst month for American Airlines departure delays. During this month, average departure delays reached a peak of 36.33 minutes, and 38.04% of the 86,083 completed flights departed 15 or more minutes late.

- Rows returned: 458
- Columns: Year, Month, FlightVolume, AvgDepDelay, PctDel15

| Year | Month | FlightVolume | AvgDepDelay | PctDel15 |
| --- | --- | --- | --- | --- |
| 2024 | 7 | 86083 | 36.33 | 0.3804 |

> Which origins contribute most to that peak month?

The delays in July 2024 were driven primarily by Dallas/Fort Worth (DFW) and Charlotte (CLT). DFW generated over 576,000 delay minutes (averaging 38.5 minutes per flight), while CLT was responsible for roughly 549,000 delay minutes with severe average delays reaching 52 minutes.

- Rows returned: 93
- Columns: OriginCode, FlightVolume, AvgDepDelay, PctDel15, TotalDelayMinutes, DelayedFlights

| OriginCode | FlightVolume | AvgDepDelay | PctDel15 | TotalDelayMinutes | DelayedFlights |
| --- | --- | --- | --- | --- | --- |
| DFW | 14962 | 38.52 | 0.4659 | 576310 | 6971 |

> Which routes contribute most to that peak month?

The most severely impacted routes were hub-to-hub and hub-to-spoke connections from DFW and CLT. The top route contributor was DFW to Los Angeles (LAX) with over 18,000 delay minutes, closely followed by CLT to Orlando (MCO) and DFW to San Antonio (SAT).

- Rows returned: 799
- Columns: OriginCode, DestCode, Route, FlightVolume, AvgDepDelay, PctDel15, TotalDelayMinutes, DelayedFlights

| OriginCode | DestCode | Route | FlightVolume | AvgDepDelay | PctDel15 | TotalDelayMinutes | DelayedFlights |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DFW | LAX | DFW-LAX | 443 | 40.89 | 0.465 | 18114 | 206 |

> Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

The peak is highly concentrated at the origin level but broadly distributed across routes. In July 2024, the top 5 origins (led overwhelmingly by DFW and CLT) accounted for nearly 52% of the total network departure delay minutes. In contrast, the network's delays were widely spread across individual routes, with the top 10 worst routes representing only 5% of all network delay minutes.

- Rows returned: 1
- Columns: network_delay, top_5_origins_delay, top_10_routes_delay, top_5_origins_pct, top_10_routes_pct

| network_delay | top_5_origins_delay | top_10_routes_delay | top_5_origins_pct | top_10_routes_pct |
| --- | --- | --- | --- | --- |
| 3.127538e+06 | 1.624257e+06 | 156919 | 51.93 | 5.02 |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "subquestion": "Which month is the single worst American Airlines month for departure delays?",
      "answer_markdown": "July 2024 is the single worst month for American Airlines departure delays. During this month, average departure delays reached a peak of 36.33 minutes, and 38.04% of the 86,083 completed flights departed 15 or more minutes late.",
      "sql": "SELECT Year, Month, count() AS FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 GROUP BY Year, Month ORDER BY AvgDepDelay DESC;"
    },
    {
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "The delays in July 2024 were driven primarily by Dallas/Fort Worth (DFW) and Charlotte (CLT). DFW generated over 576,000 delay minutes (averaging 38.5 minutes per flight), while CLT was responsible for roughly 549,000 delay minutes with severe average delays reaching 52 minutes.",
      "sql": "SELECT OriginCode, count() as FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15, sum(DepDelay) as TotalDelayMinutes, sum(DepDel15) as DelayedFlights FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode HAVING FlightVolume >= 100 ORDER BY TotalDelayMinutes DESC;"
    },
    {
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "The most severely impacted routes were hub-to-hub and hub-to-spoke connections from DFW and CLT. The top route contributor was DFW to Los Angeles (LAX) with over 18,000 delay minutes, closely followed by CLT to Orlando (MCO) and DFW to San Antonio (SAT).",
      "sql": "SELECT OriginCode, DestCode, concat(OriginCode, '-', DestCode) as Route, count() as FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15, sum(DepDelay) as TotalDelayMinutes, sum(DepDel15) as DelayedFlights FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode, DestCode HAVING FlightVolume >= 30 ORDER BY TotalDelayMinutes DESC;"
    },
    {
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "The peak is highly concentrated at the origin level but broadly distributed across routes. In July 2024, the top 5 origins (led overwhelmingly by DFW and CLT) accounted for nearly 52% of the total network departure delay minutes. In contrast, the network's delays were widely spread across individual routes, with the top 10 worst routes representing only 5% of all network delay minutes.",
      "sql": "WITH network_total AS (SELECT sum(DepDelay) as network_delay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7), top_origins AS (SELECT sum(DepDelay) as origin_delay FROM (SELECT sum(DepDelay) as DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode ORDER BY DepDelay DESC LIMIT 5)), top_routes AS (SELECT sum(DepDelay) as route_delay FROM (SELECT sum(DepDelay) as DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode, DestCode ORDER BY DepDelay DESC LIMIT 10)) SELECT n.network_delay, o.origin_delay as top_5_origins_delay, r.route_delay as top_10_routes_delay, round(o.origin_delay / n.network_delay * 100, 2) as top_5_origins_pct, round(r.route_delay / n.network_delay * 100, 2) as top_10_routes_pct FROM network_total n, top_origins o, top_routes r;"
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
      "answer_markdown": "July 2024 is the single worst month for American Airlines departure delays. During this month, average departure delays reached a peak of 36.33 minutes, and 38.04% of the 86,083 completed flights departed 15 or more minutes late.",
      "sql": "SELECT Year, Month, count() AS FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 GROUP BY Year, Month ORDER BY AvgDepDelay DESC;"
    },
    {
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "The delays in July 2024 were driven primarily by Dallas/Fort Worth (DFW) and Charlotte (CLT). DFW generated over 576,000 delay minutes (averaging 38.5 minutes per flight), while CLT was responsible for roughly 549,000 delay minutes with severe average delays reaching 52 minutes.",
      "sql": "SELECT OriginCode, count() as FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15, sum(DepDelay) as TotalDelayMinutes, sum(DepDel15) as DelayedFlights FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode HAVING FlightVolume \u003e= 100 ORDER BY TotalDelayMinutes DESC;"
    },
    {
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "The most severely impacted routes were hub-to-hub and hub-to-spoke connections from DFW and CLT. The top route contributor was DFW to Los Angeles (LAX) with over 18,000 delay minutes, closely followed by CLT to Orlando (MCO) and DFW to San Antonio (SAT).",
      "sql": "SELECT OriginCode, DestCode, concat(OriginCode, '-', DestCode) as Route, count() as FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15, sum(DepDelay) as TotalDelayMinutes, sum(DepDel15) as DelayedFlights FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode, DestCode HAVING FlightVolume \u003e= 30 ORDER BY TotalDelayMinutes DESC;"
    },
    {
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "The peak is highly concentrated at the origin level but broadly distributed across routes. In July 2024, the top 5 origins (led overwhelmingly by DFW and CLT) accounted for nearly 52% of the total network departure delay minutes. In contrast, the network's delays were widely spread across individual routes, with the top 10 worst routes representing only 5% of all network delay minutes.",
      "sql": "WITH network_total AS (SELECT sum(DepDelay) as network_delay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7), top_origins AS (SELECT sum(DepDelay) as origin_delay FROM (SELECT sum(DepDelay) as DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode ORDER BY DepDelay DESC LIMIT 5)), top_routes AS (SELECT sum(DepDelay) as route_delay FROM (SELECT sum(DepDelay) as DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode, DestCode ORDER BY DepDelay DESC LIMIT 10)) SELECT n.network_delay, o.origin_delay as top_5_origins_delay, r.route_delay as top_10_routes_delay, round(o.origin_delay / n.network_delay * 100, 2) as top_5_origins_pct, round(r.route_delay / n.network_delay * 100, 2) as top_10_routes_pct FROM network_total n, top_origins o, top_routes r;"
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
      "answer_markdown": "July 2024 is the single worst month for American Airlines departure delays. During this month, average departure delays reached a peak of 36.33 minutes, and 38.04% of the 86,083 completed flights departed 15 or more minutes late.",
      "sql": "SELECT Year, Month, count() AS FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 GROUP BY Year, Month ORDER BY AvgDepDelay DESC;",
      "row_count": 458,
      "result_columns": [
        "Year",
        "Month",
        "FlightVolume",
        "AvgDepDelay",
        "PctDel15"
      ],
      "first_row": {
        "AvgDepDelay": 36.33,
        "FlightVolume": 86083,
        "Month": 7,
        "PctDel15": 0.3804,
        "Year": 2024
      }
    },
    {
      "id": "q2",
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "The delays in July 2024 were driven primarily by Dallas/Fort Worth (DFW) and Charlotte (CLT). DFW generated over 576,000 delay minutes (averaging 38.5 minutes per flight), while CLT was responsible for roughly 549,000 delay minutes with severe average delays reaching 52 minutes.",
      "sql": "SELECT OriginCode, count() as FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15, sum(DepDelay) as TotalDelayMinutes, sum(DepDel15) as DelayedFlights FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode HAVING FlightVolume \u003e= 100 ORDER BY TotalDelayMinutes DESC;",
      "row_count": 93,
      "result_columns": [
        "OriginCode",
        "FlightVolume",
        "AvgDepDelay",
        "PctDel15",
        "TotalDelayMinutes",
        "DelayedFlights"
      ],
      "first_row": {
        "AvgDepDelay": 38.52,
        "DelayedFlights": 6971,
        "FlightVolume": 14962,
        "OriginCode": "DFW",
        "PctDel15": 0.4659,
        "TotalDelayMinutes": 576310
      }
    },
    {
      "id": "q3",
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "The most severely impacted routes were hub-to-hub and hub-to-spoke connections from DFW and CLT. The top route contributor was DFW to Los Angeles (LAX) with over 18,000 delay minutes, closely followed by CLT to Orlando (MCO) and DFW to San Antonio (SAT).",
      "sql": "SELECT OriginCode, DestCode, concat(OriginCode, '-', DestCode) as Route, count() as FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15, sum(DepDelay) as TotalDelayMinutes, sum(DepDel15) as DelayedFlights FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode, DestCode HAVING FlightVolume \u003e= 30 ORDER BY TotalDelayMinutes DESC;",
      "row_count": 799,
      "result_columns": [
        "OriginCode",
        "DestCode",
        "Route",
        "FlightVolume",
        "AvgDepDelay",
        "PctDel15",
        "TotalDelayMinutes",
        "DelayedFlights"
      ],
      "first_row": {
        "AvgDepDelay": 40.89,
        "DelayedFlights": 206,
        "DestCode": "LAX",
        "FlightVolume": 443,
        "OriginCode": "DFW",
        "PctDel15": 0.465,
        "Route": "DFW-LAX",
        "TotalDelayMinutes": 18114
      }
    },
    {
      "id": "q4",
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "The peak is highly concentrated at the origin level but broadly distributed across routes. In July 2024, the top 5 origins (led overwhelmingly by DFW and CLT) accounted for nearly 52% of the total network departure delay minutes. In contrast, the network's delays were widely spread across individual routes, with the top 10 worst routes representing only 5% of all network delay minutes.",
      "sql": "WITH network_total AS (SELECT sum(DepDelay) as network_delay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7), top_origins AS (SELECT sum(DepDelay) as origin_delay FROM (SELECT sum(DepDelay) as DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode ORDER BY DepDelay DESC LIMIT 5)), top_routes AS (SELECT sum(DepDelay) as route_delay FROM (SELECT sum(DepDelay) as DepDelay FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND Year = 2024 AND Month = 7 GROUP BY OriginCode, DestCode ORDER BY DepDelay DESC LIMIT 10)) SELECT n.network_delay, o.origin_delay as top_5_origins_delay, r.route_delay as top_10_routes_delay, round(o.origin_delay / n.network_delay * 100, 2) as top_5_origins_pct, round(r.route_delay / n.network_delay * 100, 2) as top_10_routes_pct FROM network_total n, top_origins o, top_routes r;",
      "row_count": 1,
      "result_columns": [
        "network_delay",
        "top_5_origins_delay",
        "top_10_routes_delay",
        "top_5_origins_pct",
        "top_10_routes_pct"
      ],
      "first_row": {
        "network_delay": 3127538,
        "top_10_routes_delay": 156919,
        "top_10_routes_pct": 5.02,
        "top_5_origins_delay": 1624257,
        "top_5_origins_pct": 51.93
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