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

Determine which airline and origin-airport combinations perform worst in winter after applying a meaningful flight threshold.

Focus on winter departures only and evaluate completed flights at the `(carrier, origin airport)` level. Limit the analysis to combinations with enough winter traffic to be credible.

Use winter consistently as the business definition of the season for the full available history. You may apply a reasonable minimum-volume filter to remove noise, but do not invent a custom score or let delay-cause shares replace the primary performance ranking.

For each qualifying pair, quantify:

- winter flight volume
- departure on-time performance
- average departure delay
- how reported delay minutes split across weather and operational causes such as carrier, NAS, security, and late aircraft

If delay-cause fields are missing or unreported for some of the weakest pairs, surface that explicitly in both the proof query output and the report. Do not present a blanket weather-vs-operational conclusion as if it covers pairs whose cause data is unavailable.

Rank the worst-performing winter pairs by on-time performance, while using the delay-cause mix as context rather than as the primary ranking driver.

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked view of the weakest qualifying winter carrier-airport pairs
- a cause-composition view for the leading weak pairs that separates weather from operational causes

The proof query behind the concentration question must directly support both parts of that question. It should make carrier concentration and airport concentration separately inspectable from the returned result, rather than requiring the reader to infer one of them indirectly from lists or arrays.

The proof query behind the worst-pair question should preserve the ranked pair-level rows needed for the dashboard, not just a single worst pair or summary count.

## Dashboard Questions

- Which winter carrier-airport pair ranks worst overall?
- Are the worst pairs driven more by weather or by operational causes?
- Are the weakest pairs concentrated in a small number of carriers or airports?

In the report, answer those questions directly in prose. Name the worst winter pair, summarize whether the weakest pairs are driven more by weather or by operational causes, and state whether the weak set is concentrated in a small number of carriers or airports.

If some leading weak pairs lack delay-cause reporting, say so directly and limit the weather-vs-operational conclusion to the subset with measured cause data.

Do not use fallback phrases such as "the worst pair" or "the weakest pairs" when your verified query results let you name the actual carrier-airport combinations directly.

Keep the result business-readable and analytically sound. Exclude low-volume winter pairs before ranking them.

Generated report.md:

```md
# Worst winter carrier-origin pairs by departure performance

> Which winter carrier-airport pair ranks worst overall?

Using December through February and excluding cancelled or diverted flights, the worst qualifying pair is `DH` departing `Chicago O'Hare International (ORD)`. With a minimum threshold of 1,000 winter departures across the full history, it has 19,929 winter flights, the lowest departure on-time rate at 56.53%, and an average departure delay of 27.06 minutes.

- Rows returned: 25
- Columns: Carrier, OriginCode, DisplayAirportName, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause

| Carrier | OriginCode | DisplayAirportName | winter_flights | dep_otp_pct | avg_dep_delay_min | flights_with_any_reported_cause | flights_without_reported_cause |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DH | ORD | Chicago O'Hare International | 19929 | 56.53 | 27.06 | 5681 | 14248 |

> Are the worst pairs driven more by weather or by operational causes?

For the leading weak winter pairs that actually report delay causes, operational causes dominate reported delay minutes rather than weather. In the verified top-10 weak set, every pair with measured cause data is majority operational, with weather contributing only 1.5% to 26.8% of reported cause minutes. Cause reporting is incomplete, though: `PI-DFW`, `PI-LAX`, and `PI-DAY` show no reported cause minutes at all, so the weather-versus-operations conclusion applies only to the measured subset.

- Rows returned: 10
- Columns: Carrier, OriginCode, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause, weather_delay_min, carrier_delay_min, nas_delay_min, security_delay_min, late_aircraft_delay_min, total_reported_cause_min, weather_share_pct, operational_share_pct

| Carrier | OriginCode | winter_flights | dep_otp_pct | avg_dep_delay_min | flights_with_any_reported_cause | flights_without_reported_cause | weather_delay_min | carrier_delay_min | nas_delay_min | security_delay_min | late_aircraft_delay_min | total_reported_cause_min | weather_share_pct | operational_share_pct |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DH | ORD | 19929 | 56.53 | 27.06 | 5681 | 14248 | 55176 | 123407 | 41492 | 116 | 133239 | 353430 | 15.6 | 84.4 |

> Are the weakest pairs concentrated in a small number of carriers or airports?

The weak set is only mildly concentrated in carriers and not very concentrated in airports. Looking at the top 10 qualifying weak winter pairs, `PI` accounts for 3 pairs (30%), while `ORD` is the only airport that appears more than once with 2 pairs (20%). Every other carrier and airport appears once, so the concentration is more noticeable by carrier than by origin airport.

- Rows returned: 17
- Columns: concentration_type, entity, weak_pair_count, weak_pair_share_pct

| concentration_type | entity | weak_pair_count | weak_pair_share_pct |
| --- | --- | --- | --- |
| carrier | OH | 1 | 10 |
```

Saved answer.raw.json:

```json
{
  "subquestions": [
    {
      "subquestion": "Which winter carrier-airport pair ranks worst overall?",
      "answer_markdown": "Using December through February and excluding cancelled or diverted flights, the worst qualifying pair is `DH` departing `Chicago O'Hare International (ORD)`. With a minimum threshold of 1,000 winter departures across the full history, it has 19,929 winter flights, the lowest departure on-time rate at 56.53%, and an average departure delay of 27.06 minutes.",
      "sql": "WITH ranked_pairs AS (SELECT f.Carrier, f.OriginCode, a.DisplayAirportName, count() AS winter_flights, round(avg(if(f.DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(f.DepDelay), 0.0)), 2) AS avg_dep_delay_min, countIf((f.CarrierDelay IS NOT NULL) OR (f.WeatherDelay IS NOT NULL) OR (f.NASDelay IS NOT NULL) OR (f.SecurityDelay IS NOT NULL) OR (f.LateAircraftDelay IS NOT NULL)) AS flights_with_any_reported_cause, countIf((f.CarrierDelay IS NULL) AND (f.WeatherDelay IS NULL) AND (f.NASDelay IS NULL) AND (f.SecurityDelay IS NULL) AND (f.LateAircraftDelay IS NULL)) AS flights_without_reported_cause FROM ontime.fact_ontime AS f LEFT JOIN ontime.dim_airports AS a ON f.OriginAirportID = a.AirportID WHERE f.Month IN (12, 1, 2) AND f.Cancelled = 0 AND f.Diverted = 0 GROUP BY f.Carrier, f.OriginCode, a.DisplayAirportName HAVING winter_flights >= 1000) SELECT Carrier, OriginCode, DisplayAirportName, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause FROM ranked_pairs ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, winter_flights DESC LIMIT 25"
    },
    {
      "subquestion": "Are the worst pairs driven more by weather or by operational causes?",
      "answer_markdown": "For the leading weak winter pairs that actually report delay causes, operational causes dominate reported delay minutes rather than weather. In the verified top-10 weak set, every pair with measured cause data is majority operational, with weather contributing only 1.5% to 26.8% of reported cause minutes. Cause reporting is incomplete, though: `PI-DFW`, `PI-LAX`, and `PI-DAY` show no reported cause minutes at all, so the weather-versus-operations conclusion applies only to the measured subset.",
      "sql": "WITH weak_pairs AS (SELECT f.Carrier, f.OriginCode, count() AS winter_flights, round(avg(if(f.DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(f.DepDelay), 0.0)), 2) AS avg_dep_delay_min, sum(ifNull(toUInt64(f.WeatherDelay), 0)) AS weather_delay_min, sum(ifNull(toUInt64(f.CarrierDelay), 0)) AS carrier_delay_min, sum(ifNull(toUInt64(f.NASDelay), 0)) AS nas_delay_min, sum(ifNull(toUInt64(f.SecurityDelay), 0)) AS security_delay_min, sum(ifNull(toUInt64(f.LateAircraftDelay), 0)) AS late_aircraft_delay_min, countIf((f.CarrierDelay IS NOT NULL) OR (f.WeatherDelay IS NOT NULL) OR (f.NASDelay IS NOT NULL) OR (f.SecurityDelay IS NOT NULL) OR (f.LateAircraftDelay IS NOT NULL)) AS flights_with_any_reported_cause, countIf((f.CarrierDelay IS NULL) AND (f.WeatherDelay IS NULL) AND (f.NASDelay IS NULL) AND (f.SecurityDelay IS NULL) AND (f.LateAircraftDelay IS NULL)) AS flights_without_reported_cause FROM ontime.fact_ontime AS f WHERE f.Month IN (12, 1, 2) AND f.Cancelled = 0 AND f.Diverted = 0 GROUP BY f.Carrier, f.OriginCode HAVING winter_flights >= 1000 ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC LIMIT 10) SELECT Carrier, OriginCode, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause, weather_delay_min, carrier_delay_min, nas_delay_min, security_delay_min, late_aircraft_delay_min, weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min AS total_reported_cause_min, round(weather_delay_min * 100.0 / nullIf(weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min, 0), 1) AS weather_share_pct, round((carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min) * 100.0 / nullIf(weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min, 0), 1) AS operational_share_pct FROM weak_pairs ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, winter_flights DESC"
    },
    {
      "subquestion": "Are the weakest pairs concentrated in a small number of carriers or airports?",
      "answer_markdown": "The weak set is only mildly concentrated in carriers and not very concentrated in airports. Looking at the top 10 qualifying weak winter pairs, `PI` accounts for 3 pairs (30%), while `ORD` is the only airport that appears more than once with 2 pairs (20%). Every other carrier and airport appears once, so the concentration is more noticeable by carrier than by origin airport.",
      "sql": "WITH weak_pairs AS (SELECT Carrier, OriginCode, count() AS winter_flights, round(avg(if(DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(DepDelay), 0.0)), 2) AS avg_dep_delay_min FROM ontime.fact_ontime WHERE Month IN (12, 1, 2) AND Cancelled = 0 AND Diverted = 0 GROUP BY Carrier, OriginCode HAVING winter_flights >= 1000 ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC LIMIT 10) SELECT 'carrier' AS concentration_type, Carrier AS entity, count() AS weak_pair_count, round(count() * 100.0 / (SELECT count() FROM weak_pairs), 1) AS weak_pair_share_pct FROM weak_pairs GROUP BY Carrier UNION ALL SELECT 'airport' AS concentration_type, OriginCode AS entity, count() AS weak_pair_count, round(count() * 100.0 / (SELECT count() FROM weak_pairs), 1) AS weak_pair_share_pct FROM weak_pairs GROUP BY OriginCode ORDER BY concentration_type, weak_pair_count DESC, entity ASC"
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
      "subquestion": "Which winter carrier-airport pair ranks worst overall?",
      "answer_markdown": "Using December through February and excluding cancelled or diverted flights, the worst qualifying pair is `DH` departing `Chicago O'Hare International (ORD)`. With a minimum threshold of 1,000 winter departures across the full history, it has 19,929 winter flights, the lowest departure on-time rate at 56.53%, and an average departure delay of 27.06 minutes.",
      "sql": "WITH ranked_pairs AS (SELECT f.Carrier, f.OriginCode, a.DisplayAirportName, count() AS winter_flights, round(avg(if(f.DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(f.DepDelay), 0.0)), 2) AS avg_dep_delay_min, countIf((f.CarrierDelay IS NOT NULL) OR (f.WeatherDelay IS NOT NULL) OR (f.NASDelay IS NOT NULL) OR (f.SecurityDelay IS NOT NULL) OR (f.LateAircraftDelay IS NOT NULL)) AS flights_with_any_reported_cause, countIf((f.CarrierDelay IS NULL) AND (f.WeatherDelay IS NULL) AND (f.NASDelay IS NULL) AND (f.SecurityDelay IS NULL) AND (f.LateAircraftDelay IS NULL)) AS flights_without_reported_cause FROM ontime.fact_ontime AS f LEFT JOIN ontime.dim_airports AS a ON f.OriginAirportID = a.AirportID WHERE f.Month IN (12, 1, 2) AND f.Cancelled = 0 AND f.Diverted = 0 GROUP BY f.Carrier, f.OriginCode, a.DisplayAirportName HAVING winter_flights \u003e= 1000) SELECT Carrier, OriginCode, DisplayAirportName, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause FROM ranked_pairs ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, winter_flights DESC LIMIT 25"
    },
    {
      "subquestion": "Are the worst pairs driven more by weather or by operational causes?",
      "answer_markdown": "For the leading weak winter pairs that actually report delay causes, operational causes dominate reported delay minutes rather than weather. In the verified top-10 weak set, every pair with measured cause data is majority operational, with weather contributing only 1.5% to 26.8% of reported cause minutes. Cause reporting is incomplete, though: `PI-DFW`, `PI-LAX`, and `PI-DAY` show no reported cause minutes at all, so the weather-versus-operations conclusion applies only to the measured subset.",
      "sql": "WITH weak_pairs AS (SELECT f.Carrier, f.OriginCode, count() AS winter_flights, round(avg(if(f.DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(f.DepDelay), 0.0)), 2) AS avg_dep_delay_min, sum(ifNull(toUInt64(f.WeatherDelay), 0)) AS weather_delay_min, sum(ifNull(toUInt64(f.CarrierDelay), 0)) AS carrier_delay_min, sum(ifNull(toUInt64(f.NASDelay), 0)) AS nas_delay_min, sum(ifNull(toUInt64(f.SecurityDelay), 0)) AS security_delay_min, sum(ifNull(toUInt64(f.LateAircraftDelay), 0)) AS late_aircraft_delay_min, countIf((f.CarrierDelay IS NOT NULL) OR (f.WeatherDelay IS NOT NULL) OR (f.NASDelay IS NOT NULL) OR (f.SecurityDelay IS NOT NULL) OR (f.LateAircraftDelay IS NOT NULL)) AS flights_with_any_reported_cause, countIf((f.CarrierDelay IS NULL) AND (f.WeatherDelay IS NULL) AND (f.NASDelay IS NULL) AND (f.SecurityDelay IS NULL) AND (f.LateAircraftDelay IS NULL)) AS flights_without_reported_cause FROM ontime.fact_ontime AS f WHERE f.Month IN (12, 1, 2) AND f.Cancelled = 0 AND f.Diverted = 0 GROUP BY f.Carrier, f.OriginCode HAVING winter_flights \u003e= 1000 ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC LIMIT 10) SELECT Carrier, OriginCode, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause, weather_delay_min, carrier_delay_min, nas_delay_min, security_delay_min, late_aircraft_delay_min, weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min AS total_reported_cause_min, round(weather_delay_min * 100.0 / nullIf(weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min, 0), 1) AS weather_share_pct, round((carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min) * 100.0 / nullIf(weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min, 0), 1) AS operational_share_pct FROM weak_pairs ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, winter_flights DESC"
    },
    {
      "subquestion": "Are the weakest pairs concentrated in a small number of carriers or airports?",
      "answer_markdown": "The weak set is only mildly concentrated in carriers and not very concentrated in airports. Looking at the top 10 qualifying weak winter pairs, `PI` accounts for 3 pairs (30%), while `ORD` is the only airport that appears more than once with 2 pairs (20%). Every other carrier and airport appears once, so the concentration is more noticeable by carrier than by origin airport.",
      "sql": "WITH weak_pairs AS (SELECT Carrier, OriginCode, count() AS winter_flights, round(avg(if(DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(DepDelay), 0.0)), 2) AS avg_dep_delay_min FROM ontime.fact_ontime WHERE Month IN (12, 1, 2) AND Cancelled = 0 AND Diverted = 0 GROUP BY Carrier, OriginCode HAVING winter_flights \u003e= 1000 ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC LIMIT 10) SELECT 'carrier' AS concentration_type, Carrier AS entity, count() AS weak_pair_count, round(count() * 100.0 / (SELECT count() FROM weak_pairs), 1) AS weak_pair_share_pct FROM weak_pairs GROUP BY Carrier UNION ALL SELECT 'airport' AS concentration_type, OriginCode AS entity, count() AS weak_pair_count, round(count() * 100.0 / (SELECT count() FROM weak_pairs), 1) AS weak_pair_share_pct FROM weak_pairs GROUP BY OriginCode ORDER BY concentration_type, weak_pair_count DESC, entity ASC"
    }
  ]
}
```

Saved visual_input.json:

```json
{
  "question_title": "Worst winter carrier-origin pairs by departure performance",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which winter carrier-airport pair ranks worst overall?",
      "answer_markdown": "Using December through February and excluding cancelled or diverted flights, the worst qualifying pair is `DH` departing `Chicago O'Hare International (ORD)`. With a minimum threshold of 1,000 winter departures across the full history, it has 19,929 winter flights, the lowest departure on-time rate at 56.53%, and an average departure delay of 27.06 minutes.",
      "sql": "WITH ranked_pairs AS (SELECT f.Carrier, f.OriginCode, a.DisplayAirportName, count() AS winter_flights, round(avg(if(f.DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(f.DepDelay), 0.0)), 2) AS avg_dep_delay_min, countIf((f.CarrierDelay IS NOT NULL) OR (f.WeatherDelay IS NOT NULL) OR (f.NASDelay IS NOT NULL) OR (f.SecurityDelay IS NOT NULL) OR (f.LateAircraftDelay IS NOT NULL)) AS flights_with_any_reported_cause, countIf((f.CarrierDelay IS NULL) AND (f.WeatherDelay IS NULL) AND (f.NASDelay IS NULL) AND (f.SecurityDelay IS NULL) AND (f.LateAircraftDelay IS NULL)) AS flights_without_reported_cause FROM ontime.fact_ontime AS f LEFT JOIN ontime.dim_airports AS a ON f.OriginAirportID = a.AirportID WHERE f.Month IN (12, 1, 2) AND f.Cancelled = 0 AND f.Diverted = 0 GROUP BY f.Carrier, f.OriginCode, a.DisplayAirportName HAVING winter_flights \u003e= 1000) SELECT Carrier, OriginCode, DisplayAirportName, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause FROM ranked_pairs ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, winter_flights DESC LIMIT 25",
      "row_count": 25,
      "result_columns": [
        "Carrier",
        "OriginCode",
        "DisplayAirportName",
        "winter_flights",
        "dep_otp_pct",
        "avg_dep_delay_min",
        "flights_with_any_reported_cause",
        "flights_without_reported_cause"
      ],
      "first_row": {
        "Carrier": "DH",
        "DisplayAirportName": "Chicago O'Hare International",
        "OriginCode": "ORD",
        "avg_dep_delay_min": 27.06,
        "dep_otp_pct": 56.53,
        "flights_with_any_reported_cause": 5681,
        "flights_without_reported_cause": 14248,
        "winter_flights": 19929
      }
    },
    {
      "id": "q2",
      "subquestion": "Are the worst pairs driven more by weather or by operational causes?",
      "answer_markdown": "For the leading weak winter pairs that actually report delay causes, operational causes dominate reported delay minutes rather than weather. In the verified top-10 weak set, every pair with measured cause data is majority operational, with weather contributing only 1.5% to 26.8% of reported cause minutes. Cause reporting is incomplete, though: `PI-DFW`, `PI-LAX`, and `PI-DAY` show no reported cause minutes at all, so the weather-versus-operations conclusion applies only to the measured subset.",
      "sql": "WITH weak_pairs AS (SELECT f.Carrier, f.OriginCode, count() AS winter_flights, round(avg(if(f.DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(f.DepDelay), 0.0)), 2) AS avg_dep_delay_min, sum(ifNull(toUInt64(f.WeatherDelay), 0)) AS weather_delay_min, sum(ifNull(toUInt64(f.CarrierDelay), 0)) AS carrier_delay_min, sum(ifNull(toUInt64(f.NASDelay), 0)) AS nas_delay_min, sum(ifNull(toUInt64(f.SecurityDelay), 0)) AS security_delay_min, sum(ifNull(toUInt64(f.LateAircraftDelay), 0)) AS late_aircraft_delay_min, countIf((f.CarrierDelay IS NOT NULL) OR (f.WeatherDelay IS NOT NULL) OR (f.NASDelay IS NOT NULL) OR (f.SecurityDelay IS NOT NULL) OR (f.LateAircraftDelay IS NOT NULL)) AS flights_with_any_reported_cause, countIf((f.CarrierDelay IS NULL) AND (f.WeatherDelay IS NULL) AND (f.NASDelay IS NULL) AND (f.SecurityDelay IS NULL) AND (f.LateAircraftDelay IS NULL)) AS flights_without_reported_cause FROM ontime.fact_ontime AS f WHERE f.Month IN (12, 1, 2) AND f.Cancelled = 0 AND f.Diverted = 0 GROUP BY f.Carrier, f.OriginCode HAVING winter_flights \u003e= 1000 ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC LIMIT 10) SELECT Carrier, OriginCode, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause, weather_delay_min, carrier_delay_min, nas_delay_min, security_delay_min, late_aircraft_delay_min, weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min AS total_reported_cause_min, round(weather_delay_min * 100.0 / nullIf(weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min, 0), 1) AS weather_share_pct, round((carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min) * 100.0 / nullIf(weather_delay_min + carrier_delay_min + nas_delay_min + security_delay_min + late_aircraft_delay_min, 0), 1) AS operational_share_pct FROM weak_pairs ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC, winter_flights DESC",
      "row_count": 10,
      "result_columns": [
        "Carrier",
        "OriginCode",
        "winter_flights",
        "dep_otp_pct",
        "avg_dep_delay_min",
        "flights_with_any_reported_cause",
        "flights_without_reported_cause",
        "weather_delay_min",
        "carrier_delay_min",
        "nas_delay_min",
        "security_delay_min",
        "late_aircraft_delay_min",
        "total_reported_cause_min",
        "weather_share_pct",
        "operational_share_pct"
      ],
      "first_row": {
        "Carrier": "DH",
        "OriginCode": "ORD",
        "avg_dep_delay_min": 27.06,
        "carrier_delay_min": 123407,
        "dep_otp_pct": 56.53,
        "flights_with_any_reported_cause": 5681,
        "flights_without_reported_cause": 14248,
        "late_aircraft_delay_min": 133239,
        "nas_delay_min": 41492,
        "operational_share_pct": 84.4,
        "security_delay_min": 116,
        "total_reported_cause_min": 353430,
        "weather_delay_min": 55176,
        "weather_share_pct": 15.6,
        "winter_flights": 19929
      }
    },
    {
      "id": "q3",
      "subquestion": "Are the weakest pairs concentrated in a small number of carriers or airports?",
      "answer_markdown": "The weak set is only mildly concentrated in carriers and not very concentrated in airports. Looking at the top 10 qualifying weak winter pairs, `PI` accounts for 3 pairs (30%), while `ORD` is the only airport that appears more than once with 2 pairs (20%). Every other carrier and airport appears once, so the concentration is more noticeable by carrier than by origin airport.",
      "sql": "WITH weak_pairs AS (SELECT Carrier, OriginCode, count() AS winter_flights, round(avg(if(DepDel15 = 0, 1.0, 0.0)) * 100, 2) AS dep_otp_pct, round(avg(ifNull(toFloat64(DepDelay), 0.0)), 2) AS avg_dep_delay_min FROM ontime.fact_ontime WHERE Month IN (12, 1, 2) AND Cancelled = 0 AND Diverted = 0 GROUP BY Carrier, OriginCode HAVING winter_flights \u003e= 1000 ORDER BY dep_otp_pct ASC, avg_dep_delay_min DESC LIMIT 10) SELECT 'carrier' AS concentration_type, Carrier AS entity, count() AS weak_pair_count, round(count() * 100.0 / (SELECT count() FROM weak_pairs), 1) AS weak_pair_share_pct FROM weak_pairs GROUP BY Carrier UNION ALL SELECT 'airport' AS concentration_type, OriginCode AS entity, count() AS weak_pair_count, round(count() * 100.0 / (SELECT count() FROM weak_pairs), 1) AS weak_pair_share_pct FROM weak_pairs GROUP BY OriginCode ORDER BY concentration_type, weak_pair_count DESC, entity ASC",
      "row_count": 17,
      "result_columns": [
        "concentration_type",
        "entity",
        "weak_pair_count",
        "weak_pair_share_pct"
      ],
      "first_row": {
        "concentration_type": "carrier",
        "entity": "OH",
        "weak_pair_count": 1,
        "weak_pair_share_pct": 10
      }
    }
  ]
}
```

Proof queries are saved as files in the run directory. Read the SQL files you need to verify grain, filters, metrics, and ranking logic:

- `queries/q1.sql`
- `queries/q2.sql`
- `queries/q3.sql`

Executed query results are saved as files in the run directory. Read the result files you need for verification instead of assuming the report summary is complete:

- `results/q1.json`
- `results/q2.json`
- `results/q3.json`