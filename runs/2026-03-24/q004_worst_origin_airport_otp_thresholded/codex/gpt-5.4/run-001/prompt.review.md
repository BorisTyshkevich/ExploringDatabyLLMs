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

Identify which origin airports have the worst departure on-time performance after excluding low-volume airports.

Analyze completed departures at the origin-airport level across the full available history. Focus on airports with enough traffic to make the comparison meaningful, and rank the weakest performers by departure on-time performance.

You may apply a reasonable minimum-volume filter to remove noise, but do not invent a custom score or redefine on-time performance.

For each airport, quantify:

- completed departures
- departure on-time performance
- average departure delay
- a high-delay measure that reflects the worse end of the delay distribution
- the first and last dates represented in the data

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked view of the weakest qualifying origin airports
- a comparison view that makes the gap between the very worst airports and the middle of the ranked set easy to judge

The proof query behind the worst-airport question should preserve the ranked airport-level rows needed for the dashboard, not just a single top airport or summary statistic.

## Dashboard Questions

- Which airport ranks worst on departure on-time performance?
- How large is the spread between the worst airport and the middle of the ranked set?
- Are the weakest airports mostly major hubs, or is the bottom group more mixed?

In the report, answer those questions directly in prose. Name the worst airport, describe the spread between the bottom and the middle using the verified result, and summarize whether the weakest group is mostly hubs or more mixed.

Do not use fallback phrases such as "the worst airport" or "the bottom group" when your verified query results let you name the actual airport set directly.

Keep the result business-readable and analytically sound. Exclude low-volume airports before ranking them.

Generated report.md:

```md
# Worst origin airports by departure on-time performance

> Which airport ranks worst on departure on-time performance?

Among origin airports with at least 250,000 completed departures, Chicago Midway International (MDW) ranks worst on departure on-time performance. It logged 2,512,354 completed departures with 76.52% departure OTP, an average departure delay of 11.80 minutes, a 90th-percentile departure delay of 36.95 minutes, and coverage from 1987-10-01 through 2025-11-30.

- Rows returned: 15
- Columns: otp_rank, origin_code, DisplayAirportName, completed_departures, dep_otp_pct, avg_dep_delay_min, p90_dep_delay_min, first_date, last_date

| otp_rank | origin_code | DisplayAirportName | completed_departures | dep_otp_pct | avg_dep_delay_min | p90_dep_delay_min | first_date | last_date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | MDW | Chicago Midway International | 2.512354e+06 | 76.52 | 11.8 | 36.91 | 1987-10-01T00:00:00Z | 2025-11-30T00:00:00Z |

> How large is the spread between the worst airport and the middle of the ranked set?

Using the same 250,000-departure cutoff, Chicago Midway International (MDW) sits 8.35 percentage points below Richmond International (RIC), the median-ranked qualifying airport. MDW also runs 3.71 minutes worse on average departure delay and 8.55 minutes worse at the 90th percentile of departure delay.

- Rows returned: 1
- Columns: worst_code, worst_airport, worst_dep_otp_pct, middle_code, middle_airport, middle_dep_otp_pct, otp_gap_pct_points, worst_avg_dep_delay_min, middle_avg_dep_delay_min, avg_delay_gap_min, worst_p90_dep_delay_min, middle_p90_dep_delay_min, p90_delay_gap_min

| worst_code | worst_airport | worst_dep_otp_pct | middle_code | middle_airport | middle_dep_otp_pct | otp_gap_pct_points | worst_avg_dep_delay_min | middle_avg_dep_delay_min | avg_delay_gap_min | worst_p90_dep_delay_min | middle_p90_dep_delay_min | p90_delay_gap_min |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| MDW | Chicago Midway International | 76.52 | RIC | Richmond International | 84.87 | 8.35 | 11.8 | 8.09 | 3.71 | 36.93 | 28.39 | 8.54 |

> Are the weakest airports mostly major hubs, or is the bottom group more mixed?

The weakest qualifying airports are more mixed than purely hub-driven. The bottom dozen includes giant hubs such as ORD, EWR, JFK, DEN, SFO, and DFW, but also large secondary or focus airports such as MDW, HOU, DAL, BWI, FLL, and LAS, so the laggards are not confined to one airport type.

- Rows returned: 12
- Columns: otp_rank, origin_code, DisplayAirportName, completed_departures, traffic_rank, dep_otp_pct, avg_dep_delay_min, p90_dep_delay_min, first_date, last_date

| otp_rank | origin_code | DisplayAirportName | completed_departures | traffic_rank | dep_otp_pct | avg_dep_delay_min | p90_dep_delay_min | first_date | last_date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | MDW | Chicago Midway International | 2.512354e+06 | 27 | 76.52 | 11.8 | 36.9 | 1987-10-01T00:00:00Z | 2025-11-30T00:00:00Z |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "subquestion": "Which airport ranks worst on departure on-time performance?",
      "answer_markdown": "Among origin airports with at least 250,000 completed departures, Chicago Midway International (MDW) ranks worst on departure on-time performance. It logged 2,512,354 completed departures with 76.52% departure OTP, an average departure delay of 11.80 minutes, a 90th-percentile departure delay of 36.95 minutes, and coverage from 1987-10-01 through 2025-11-30.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() >= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank FROM airport_stats) SELECT r.otp_rank, r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min, r.first_date, r.last_date FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 ORDER BY r.otp_rank LIMIT 15"
    },
    {
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "Using the same 250,000-departure cutoff, Chicago Midway International (MDW) sits 8.35 percentage points below Richmond International (RIC), the median-ranked qualifying airport. MDW also runs 3.71 minutes worse on average departure delay and 8.55 minutes worse at the 90th percentile of departure delay.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() >= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank, count() OVER () AS airport_count FROM airport_stats), worst AS (SELECT r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank = 1), middle AS (SELECT r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank = intDiv(r.airport_count + 1, 2)) SELECT worst.origin_code AS worst_code, worst.DisplayAirportName AS worst_airport, worst.dep_otp_pct AS worst_dep_otp_pct, middle.origin_code AS middle_code, middle.DisplayAirportName AS middle_airport, middle.dep_otp_pct AS middle_dep_otp_pct, round(middle.dep_otp_pct - worst.dep_otp_pct, 2) AS otp_gap_pct_points, worst.avg_dep_delay_min AS worst_avg_dep_delay_min, middle.avg_dep_delay_min AS middle_avg_dep_delay_min, round(worst.avg_dep_delay_min - middle.avg_dep_delay_min, 2) AS avg_delay_gap_min, worst.p90_dep_delay_min AS worst_p90_dep_delay_min, middle.p90_dep_delay_min AS middle_p90_dep_delay_min, round(worst.p90_dep_delay_min - middle.p90_dep_delay_min, 2) AS p90_delay_gap_min FROM worst CROSS JOIN middle"
    },
    {
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The weakest qualifying airports are more mixed than purely hub-driven. The bottom dozen includes giant hubs such as ORD, EWR, JFK, DEN, SFO, and DFW, but also large secondary or focus airports such as MDW, HOU, DAL, BWI, FLL, and LAS, so the laggards are not confined to one airport type.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() >= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank, dense_rank() OVER (ORDER BY completed_departures DESC) AS traffic_rank FROM airport_stats) SELECT r.otp_rank, r.origin_code, a.DisplayAirportName, r.completed_departures, r.traffic_rank, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min, r.first_date, r.last_date FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank <= 12 ORDER BY r.otp_rank"
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
      "subquestion": "Which airport ranks worst on departure on-time performance?",
      "answer_markdown": "Among origin airports with at least 250,000 completed departures, Chicago Midway International (MDW) ranks worst on departure on-time performance. It logged 2,512,354 completed departures with 76.52% departure OTP, an average departure delay of 11.80 minutes, a 90th-percentile departure delay of 36.95 minutes, and coverage from 1987-10-01 through 2025-11-30.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() \u003e= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank FROM airport_stats) SELECT r.otp_rank, r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min, r.first_date, r.last_date FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 ORDER BY r.otp_rank LIMIT 15"
    },
    {
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "Using the same 250,000-departure cutoff, Chicago Midway International (MDW) sits 8.35 percentage points below Richmond International (RIC), the median-ranked qualifying airport. MDW also runs 3.71 minutes worse on average departure delay and 8.55 minutes worse at the 90th percentile of departure delay.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() \u003e= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank, count() OVER () AS airport_count FROM airport_stats), worst AS (SELECT r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank = 1), middle AS (SELECT r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank = intDiv(r.airport_count + 1, 2)) SELECT worst.origin_code AS worst_code, worst.DisplayAirportName AS worst_airport, worst.dep_otp_pct AS worst_dep_otp_pct, middle.origin_code AS middle_code, middle.DisplayAirportName AS middle_airport, middle.dep_otp_pct AS middle_dep_otp_pct, round(middle.dep_otp_pct - worst.dep_otp_pct, 2) AS otp_gap_pct_points, worst.avg_dep_delay_min AS worst_avg_dep_delay_min, middle.avg_dep_delay_min AS middle_avg_dep_delay_min, round(worst.avg_dep_delay_min - middle.avg_dep_delay_min, 2) AS avg_delay_gap_min, worst.p90_dep_delay_min AS worst_p90_dep_delay_min, middle.p90_dep_delay_min AS middle_p90_dep_delay_min, round(worst.p90_dep_delay_min - middle.p90_dep_delay_min, 2) AS p90_delay_gap_min FROM worst CROSS JOIN middle"
    },
    {
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The weakest qualifying airports are more mixed than purely hub-driven. The bottom dozen includes giant hubs such as ORD, EWR, JFK, DEN, SFO, and DFW, but also large secondary or focus airports such as MDW, HOU, DAL, BWI, FLL, and LAS, so the laggards are not confined to one airport type.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() \u003e= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank, dense_rank() OVER (ORDER BY completed_departures DESC) AS traffic_rank FROM airport_stats) SELECT r.otp_rank, r.origin_code, a.DisplayAirportName, r.completed_departures, r.traffic_rank, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min, r.first_date, r.last_date FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank \u003c= 12 ORDER BY r.otp_rank"
    }
  ]
}
```

Saved visual_input.json:

```json
{
  "question_title": "Worst origin airports by departure on-time performance",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which airport ranks worst on departure on-time performance?",
      "answer_markdown": "Among origin airports with at least 250,000 completed departures, Chicago Midway International (MDW) ranks worst on departure on-time performance. It logged 2,512,354 completed departures with 76.52% departure OTP, an average departure delay of 11.80 minutes, a 90th-percentile departure delay of 36.95 minutes, and coverage from 1987-10-01 through 2025-11-30.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() \u003e= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank FROM airport_stats) SELECT r.otp_rank, r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min, r.first_date, r.last_date FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 ORDER BY r.otp_rank LIMIT 15",
      "row_count": 15,
      "result_columns": [
        "otp_rank",
        "origin_code",
        "DisplayAirportName",
        "completed_departures",
        "dep_otp_pct",
        "avg_dep_delay_min",
        "p90_dep_delay_min",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "DisplayAirportName": "Chicago Midway International",
        "avg_dep_delay_min": 11.8,
        "completed_departures": 2512354,
        "dep_otp_pct": 76.52,
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "origin_code": "MDW",
        "otp_rank": 1,
        "p90_dep_delay_min": 36.91
      }
    },
    {
      "id": "q2",
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "Using the same 250,000-departure cutoff, Chicago Midway International (MDW) sits 8.35 percentage points below Richmond International (RIC), the median-ranked qualifying airport. MDW also runs 3.71 minutes worse on average departure delay and 8.55 minutes worse at the 90th percentile of departure delay.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() \u003e= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank, count() OVER () AS airport_count FROM airport_stats), worst AS (SELECT r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank = 1), middle AS (SELECT r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank = intDiv(r.airport_count + 1, 2)) SELECT worst.origin_code AS worst_code, worst.DisplayAirportName AS worst_airport, worst.dep_otp_pct AS worst_dep_otp_pct, middle.origin_code AS middle_code, middle.DisplayAirportName AS middle_airport, middle.dep_otp_pct AS middle_dep_otp_pct, round(middle.dep_otp_pct - worst.dep_otp_pct, 2) AS otp_gap_pct_points, worst.avg_dep_delay_min AS worst_avg_dep_delay_min, middle.avg_dep_delay_min AS middle_avg_dep_delay_min, round(worst.avg_dep_delay_min - middle.avg_dep_delay_min, 2) AS avg_delay_gap_min, worst.p90_dep_delay_min AS worst_p90_dep_delay_min, middle.p90_dep_delay_min AS middle_p90_dep_delay_min, round(worst.p90_dep_delay_min - middle.p90_dep_delay_min, 2) AS p90_delay_gap_min FROM worst CROSS JOIN middle",
      "row_count": 1,
      "result_columns": [
        "worst_code",
        "worst_airport",
        "worst_dep_otp_pct",
        "middle_code",
        "middle_airport",
        "middle_dep_otp_pct",
        "otp_gap_pct_points",
        "worst_avg_dep_delay_min",
        "middle_avg_dep_delay_min",
        "avg_delay_gap_min",
        "worst_p90_dep_delay_min",
        "middle_p90_dep_delay_min",
        "p90_delay_gap_min"
      ],
      "first_row": {
        "avg_delay_gap_min": 3.71,
        "middle_airport": "Richmond International",
        "middle_avg_dep_delay_min": 8.09,
        "middle_code": "RIC",
        "middle_dep_otp_pct": 84.87,
        "middle_p90_dep_delay_min": 28.39,
        "otp_gap_pct_points": 8.35,
        "p90_delay_gap_min": 8.54,
        "worst_airport": "Chicago Midway International",
        "worst_avg_dep_delay_min": 11.8,
        "worst_code": "MDW",
        "worst_dep_otp_pct": 76.52,
        "worst_p90_dep_delay_min": 36.93
      }
    },
    {
      "id": "q3",
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The weakest qualifying airports are more mixed than purely hub-driven. The bottom dozen includes giant hubs such as ORD, EWR, JFK, DEN, SFO, and DFW, but also large secondary or focus airports such as MDW, HOU, DAL, BWI, FLL, and LAS, so the laggards are not confined to one airport type.",
      "sql": "WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() \u003e= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank, dense_rank() OVER (ORDER BY completed_departures DESC) AS traffic_rank FROM airport_stats) SELECT r.otp_rank, r.origin_code, a.DisplayAirportName, r.completed_departures, r.traffic_rank, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min, r.first_date, r.last_date FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank \u003c= 12 ORDER BY r.otp_rank",
      "row_count": 12,
      "result_columns": [
        "otp_rank",
        "origin_code",
        "DisplayAirportName",
        "completed_departures",
        "traffic_rank",
        "dep_otp_pct",
        "avg_dep_delay_min",
        "p90_dep_delay_min",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "DisplayAirportName": "Chicago Midway International",
        "avg_dep_delay_min": 11.8,
        "completed_departures": 2512354,
        "dep_otp_pct": 76.52,
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "origin_code": "MDW",
        "otp_rank": 1,
        "p90_dep_delay_min": 36.9,
        "traffic_rank": 27
      }
    }
  ]
}
```

Proof queries:

queries/q1.sql:
```sql
WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() >= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank FROM airport_stats) SELECT r.otp_rank, r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min, r.first_date, r.last_date FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 ORDER BY r.otp_rank LIMIT 15
```

queries/q2.sql:
```sql
WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() >= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank, count() OVER () AS airport_count FROM airport_stats), worst AS (SELECT r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank = 1), middle AS (SELECT r.origin_code, a.DisplayAirportName, r.completed_departures, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank = intDiv(r.airport_count + 1, 2)) SELECT worst.origin_code AS worst_code, worst.DisplayAirportName AS worst_airport, worst.dep_otp_pct AS worst_dep_otp_pct, middle.origin_code AS middle_code, middle.DisplayAirportName AS middle_airport, middle.dep_otp_pct AS middle_dep_otp_pct, round(middle.dep_otp_pct - worst.dep_otp_pct, 2) AS otp_gap_pct_points, worst.avg_dep_delay_min AS worst_avg_dep_delay_min, middle.avg_dep_delay_min AS middle_avg_dep_delay_min, round(worst.avg_dep_delay_min - middle.avg_dep_delay_min, 2) AS avg_delay_gap_min, worst.p90_dep_delay_min AS worst_p90_dep_delay_min, middle.p90_dep_delay_min AS middle_p90_dep_delay_min, round(worst.p90_dep_delay_min - middle.p90_dep_delay_min, 2) AS p90_delay_gap_min FROM worst CROSS JOIN middle
```

queries/q3.sql:
```sql
WITH airport_stats AS (SELECT OriginAirportID, any(OriginCode) AS origin_code, count() AS completed_departures, round(100.0 * avg(if(DepDel15 = 0, 1, 0)), 2) AS dep_otp_pct, round(avg(toFloat64(ifNull(DepDelay, 0))), 2) AS avg_dep_delay_min, round(quantileTDigest(0.9)(toFloat64(ifNull(DepDelay, 0))), 2) AS p90_dep_delay_min, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepTime IS NOT NULL GROUP BY OriginAirportID HAVING count() >= 250000), ranked AS (SELECT *, row_number() OVER (ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, completed_departures DESC, origin_code ASC) AS otp_rank, dense_rank() OVER (ORDER BY completed_departures DESC) AS traffic_rank FROM airport_stats) SELECT r.otp_rank, r.origin_code, a.DisplayAirportName, r.completed_departures, r.traffic_rank, r.dep_otp_pct, r.avg_dep_delay_min, r.p90_dep_delay_min, r.first_date, r.last_date FROM ranked r LEFT JOIN ontime.dim_airports_bts_full a ON a.AirportID = r.OriginAirportID AND a.IsLatest = 1 WHERE r.otp_rank <= 12 ORDER BY r.otp_rank
```

Executed query results:

results/q1.json:
```json
{
  "columns": [
    "otp_rank",
    "origin_code",
    "DisplayAirportName",
    "completed_departures",
    "dep_otp_pct",
    "avg_dep_delay_min",
    "p90_dep_delay_min",
    "first_date",
    "last_date"
  ],
  "rows": [
    {
      "DisplayAirportName": "Chicago Midway International",
      "avg_dep_delay_min": 11.8,
      "completed_departures": 2512354,
      "dep_otp_pct": 76.52,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "MDW",
      "otp_rank": 1,
      "p90_dep_delay_min": 36.91
    },
    {
      "DisplayAirportName": "Chicago O'Hare International",
      "avg_dep_delay_min": 12.16,
      "completed_departures": 11143843,
      "dep_otp_pct": 78.61,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "ORD",
      "otp_rank": 2,
      "p90_dep_delay_min": 42.77
    },
    {
      "DisplayAirportName": "Newark Liberty International",
      "avg_dep_delay_min": 12.62,
      "completed_departures": 4548113,
      "dep_otp_pct": 78.69,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "EWR",
      "otp_rank": 3,
      "p90_dep_delay_min": 44.31
    },
    {
      "DisplayAirportName": "William P Hobby",
      "avg_dep_delay_min": 10.31,
      "completed_departures": 2084605,
      "dep_otp_pct": 78.92,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "HOU",
      "otp_rank": 4,
      "p90_dep_delay_min": 32.81
    },
    {
      "DisplayAirportName": "Dallas Love Field",
      "avg_dep_delay_min": 10.34,
      "completed_departures": 1931453,
      "dep_otp_pct": 78.99,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "DAL",
      "otp_rank": 5,
      "p90_dep_delay_min": 32.11
    },
    {
      "DisplayAirportName": "John F. Kennedy International",
      "avg_dep_delay_min": 11.99,
      "completed_departures": 3072552,
      "dep_otp_pct": 79.52,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "JFK",
      "otp_rank": 6,
      "p90_dep_delay_min": 41.29
    },
    {
      "DisplayAirportName": "Baltimore/Washington International Thurgood Marshall",
      "avg_dep_delay_min": 10.39,
      "completed_departures": 3285203,
      "dep_otp_pct": 79.63,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "BWI",
      "otp_rank": 7,
      "p90_dep_delay_min": 32.91
    },
    {
      "DisplayAirportName": "Harry Reid International",
      "avg_dep_delay_min": 9.99,
      "completed_departures": 5189435,
      "dep_otp_pct": 79.75,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "LAS",
      "otp_rank": 8,
      "p90_dep_delay_min": 33.52
    },
    {
      "DisplayAirportName": "Denver International",
      "avg_dep_delay_min": 10.65,
      "completed_departures": 7385210,
      "dep_otp_pct": 79.93,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "DEN",
      "otp_rank": 9,
      "p90_dep_delay_min": 36.06
    },
    {
      "DisplayAirportName": "Fort Lauderdale-Hollywood International",
      "avg_dep_delay_min": 11.25,
      "completed_departures": 2312776,
      "dep_otp_pct": 79.99,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "FLL",
      "otp_rank": 10,
      "p90_dep_delay_min": 38.18
    },
    {
      "DisplayAirportName": "San Francisco International",
      "avg_dep_delay_min": 10.58,
      "completed_departures": 5139849,
      "dep_otp_pct": 80.4,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "SFO",
      "otp_rank": 11,
      "p90_dep_delay_min": 38.99
    },
    {
      "DisplayAirportName": "Dallas/Fort Worth International",
      "avg_dep_delay_min": 10.64,
      "completed_departures": 10023965,
      "dep_otp_pct": 80.49,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "DFW",
      "otp_rank": 12,
      "p90_dep_delay_min": 36.1
    },
    {
      "DisplayAirportName": "Miami International",
      "avg_dep_delay_min": 11.12,
      "completed_departures": 2818035,
      "dep_otp_pct": 80.57,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "MIA",
      "otp_rank": 13,
      "p90_dep_delay_min": 37.83
    },
    {
      "DisplayAirportName": "Philadelphia International",
      "avg_dep_delay_min": 10.43,
      "completed_departures": 3556123,
      "dep_otp_pct": 80.66,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "PHL",
      "otp_rank": 14,
      "p90_dep_delay_min": 35.42
    },
    {
      "DisplayAirportName": "Hartsfield-Jackson Atlanta International",
      "avg_dep_delay_min": 9.97,
      "completed_departures": 12072831,
      "dep_otp_pct": 81,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "ATL",
      "otp_rank": 15,
      "p90_dep_delay_min": 32.46
    }
  ],
  "row_count": 15,
  "generated_at": "2026-03-24T14:18:36.998836Z",
  "source_query_sha256": "04071e28a4f8db9f25003a01c7380f31e502ae4b389a40982edc236368a4df2c",
  "log_comment": "qforge|question=q004|run=run-001|runner=codex|model=gpt-5.4|phase=full|subquestion=q1"
}
```

results/q2.json:
```json
{
  "columns": [
    "worst_code",
    "worst_airport",
    "worst_dep_otp_pct",
    "middle_code",
    "middle_airport",
    "middle_dep_otp_pct",
    "otp_gap_pct_points",
    "worst_avg_dep_delay_min",
    "middle_avg_dep_delay_min",
    "avg_delay_gap_min",
    "worst_p90_dep_delay_min",
    "middle_p90_dep_delay_min",
    "p90_delay_gap_min"
  ],
  "rows": [
    {
      "avg_delay_gap_min": 3.71,
      "middle_airport": "Richmond International",
      "middle_avg_dep_delay_min": 8.09,
      "middle_code": "RIC",
      "middle_dep_otp_pct": 84.87,
      "middle_p90_dep_delay_min": 28.39,
      "otp_gap_pct_points": 8.35,
      "p90_delay_gap_min": 8.54,
      "worst_airport": "Chicago Midway International",
      "worst_avg_dep_delay_min": 11.8,
      "worst_code": "MDW",
      "worst_dep_otp_pct": 76.52,
      "worst_p90_dep_delay_min": 36.93
    }
  ],
  "row_count": 1,
  "generated_at": "2026-03-24T14:18:38.012339Z",
  "source_query_sha256": "2c1b482a586d9adb93a4f648f83c67bafd8ff25e90cd777ecdc1870256f0ebf2",
  "log_comment": "qforge|question=q004|run=run-001|runner=codex|model=gpt-5.4|phase=full|subquestion=q2"
}
```

results/q3.json:
```json
{
  "columns": [
    "otp_rank",
    "origin_code",
    "DisplayAirportName",
    "completed_departures",
    "traffic_rank",
    "dep_otp_pct",
    "avg_dep_delay_min",
    "p90_dep_delay_min",
    "first_date",
    "last_date"
  ],
  "rows": [
    {
      "DisplayAirportName": "Chicago Midway International",
      "avg_dep_delay_min": 11.8,
      "completed_departures": 2512354,
      "dep_otp_pct": 76.52,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "MDW",
      "otp_rank": 1,
      "p90_dep_delay_min": 36.9,
      "traffic_rank": 27
    },
    {
      "DisplayAirportName": "Chicago O'Hare International",
      "avg_dep_delay_min": 12.16,
      "completed_departures": 11143843,
      "dep_otp_pct": 78.61,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "ORD",
      "otp_rank": 2,
      "p90_dep_delay_min": 42.75,
      "traffic_rank": 2
    },
    {
      "DisplayAirportName": "Newark Liberty International",
      "avg_dep_delay_min": 12.62,
      "completed_departures": 4548113,
      "dep_otp_pct": 78.69,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "EWR",
      "otp_rank": 3,
      "p90_dep_delay_min": 44.31,
      "traffic_rank": 13
    },
    {
      "DisplayAirportName": "William P Hobby",
      "avg_dep_delay_min": 10.31,
      "completed_departures": 2084605,
      "dep_otp_pct": 78.92,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "HOU",
      "otp_rank": 4,
      "p90_dep_delay_min": 32.8,
      "traffic_rank": 34
    },
    {
      "DisplayAirportName": "Dallas Love Field",
      "avg_dep_delay_min": 10.34,
      "completed_departures": 1931453,
      "dep_otp_pct": 78.99,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "DAL",
      "otp_rank": 5,
      "p90_dep_delay_min": 32.1,
      "traffic_rank": 37
    },
    {
      "DisplayAirportName": "John F. Kennedy International",
      "avg_dep_delay_min": 11.99,
      "completed_departures": 3072552,
      "dep_otp_pct": 79.52,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "JFK",
      "otp_rank": 6,
      "p90_dep_delay_min": 41.28,
      "traffic_rank": 23
    },
    {
      "DisplayAirportName": "Baltimore/Washington International Thurgood Marshall",
      "avg_dep_delay_min": 10.39,
      "completed_departures": 3285203,
      "dep_otp_pct": 79.63,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "BWI",
      "otp_rank": 7,
      "p90_dep_delay_min": 32.89,
      "traffic_rank": 22
    },
    {
      "DisplayAirportName": "Harry Reid International",
      "avg_dep_delay_min": 9.99,
      "completed_departures": 5189435,
      "dep_otp_pct": 79.75,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "LAS",
      "otp_rank": 8,
      "p90_dep_delay_min": 33.52,
      "traffic_rank": 8
    },
    {
      "DisplayAirportName": "Denver International",
      "avg_dep_delay_min": 10.65,
      "completed_departures": 7385210,
      "dep_otp_pct": 79.93,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "DEN",
      "otp_rank": 9,
      "p90_dep_delay_min": 36.05,
      "traffic_rank": 4
    },
    {
      "DisplayAirportName": "Fort Lauderdale-Hollywood International",
      "avg_dep_delay_min": 11.25,
      "completed_departures": 2312776,
      "dep_otp_pct": 79.99,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "FLL",
      "otp_rank": 10,
      "p90_dep_delay_min": 38.17,
      "traffic_rank": 30
    },
    {
      "DisplayAirportName": "San Francisco International",
      "avg_dep_delay_min": 10.58,
      "completed_departures": 5139849,
      "dep_otp_pct": 80.4,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "SFO",
      "otp_rank": 11,
      "p90_dep_delay_min": 38.97,
      "traffic_rank": 10
    },
    {
      "DisplayAirportName": "Dallas/Fort Worth International",
      "avg_dep_delay_min": 10.64,
      "completed_departures": 10023965,
      "dep_otp_pct": 80.49,
      "first_date": "1987-10-01T00:00:00Z",
      "last_date": "2025-11-30T00:00:00Z",
      "origin_code": "DFW",
      "otp_rank": 12,
      "p90_dep_delay_min": 36.04,
      "traffic_rank": 3
    }
  ],
  "row_count": 12,
  "generated_at": "2026-03-24T14:18:38.674267Z",
  "source_query_sha256": "e2e66da5217b2e593d4f581ac1fee68c18255c67be2a4e11f0b0e976c71ba02a",
  "log_comment": "qforge|question=q004|run=run-001|runner=codex|model=gpt-5.4|phase=full|subquestion=q3"
}
```