- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `Worst origin airports by departure on-time performance`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_ranked_dashboard`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

Build a dashboard that:

- uses `worst_airport` as the primary saved SQL already provided in the prompt
- uses `spread_to_middle` and `bottom_group_mix` as supporting queries when they materially improve the dashboard
- shows KPI cards for worst airport, worst OTP, highest average departure delay among ranked airports, and qualifying airport count
- renders a ranked horizontal bar or lollipop chart for the worst 25 airports by departure OTP
- renders a scatter plot of `CompletedDepartures` vs `DepartureOtpPct`
- renders a detail table with the full ranked result
- includes a narrative takeaway about whether the weakest airports cluster in large hubs or a more mixed set
- derives chart extents and highlighted airports from fetched data instead of hardcoding them
- uses one accent treatment for the worst 5 airports
- annotates the single worst airport in both charts
- keeps the scatter plot readable despite skewed volume differences
- shows supporting queries in the query ledger when used

### Data Source

SQL query for primary data source:

```sql
SELECT row_number() OVER (ORDER BY round(100 - avg(f.DepDel15) * 100, 2) ASC) AS rank, f.OriginCode, coalesce(d.DisplayAirportName, f.OriginCode) AS airport_name, f.OriginCityName AS city, count() AS total_departures, round(100 - avg(f.DepDel15) * 100, 2) AS otp_rate, round(avg(f.DepDel15) * 100, 2) AS pct_delayed, round(avgIf(f.DepDelay, f.DepDelay > 0), 2) AS avg_dep_delay_when_late, quantileExact(0.90)(f.DepDelayMinutes) AS p90_delay_min, min(f.FlightDate) AS first_date, max(f.FlightDate) AS last_date FROM ontime.fact_ontime f LEFT JOIN ontime.dim_airports d ON f.OriginAirportID = d.AirportID WHERE f.Cancelled = 0 AND f.DepDelay IS NOT NULL GROUP BY f.OriginCode, d.DisplayAirportName, f.OriginCityName HAVING count() >= 100000 ORDER BY otp_rate ASC
```

Data example/snippet:

{
  "question_title": "Worst origin airports by departure on-time performance",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which airport ranks worst on departure on-time performance?",
      "answer_markdown": "Chicago Midway International (MDW) ranks worst among all qualifying origin airports, with an on-time departure rate of 76.52% — meaning nearly 1 in 4 departures was delayed 15 or more minutes. MDW logged 2.51 million completed departures across the full dataset (1987–2025) and, when delayed, averaged 26.8 minutes late with a 90th-percentile delay of 37 minutes. The minimum-volume threshold applied was 100,000 completed departures, leaving 160 qualifying airports in the ranked set.",
      "sql": "SELECT row_number() OVER (ORDER BY round(100 - avg(f.DepDel15) * 100, 2) ASC) AS rank, f.OriginCode, coalesce(d.DisplayAirportName, f.OriginCode) AS airport_name, f.OriginCityName AS city, count() AS total_departures, round(100 - avg(f.DepDel15) * 100, 2) AS otp_rate, round(avg(f.DepDel15) * 100, 2) AS pct_delayed, round(avgIf(f.DepDelay, f.DepDelay \u003e 0), 2) AS avg_dep_delay_when_late, quantileExact(0.90)(f.DepDelayMinutes) AS p90_delay_min, min(f.FlightDate) AS first_date, max(f.FlightDate) AS last_date FROM ontime.fact_ontime f LEFT JOIN ontime.dim_airports d ON f.OriginAirportID = d.AirportID WHERE f.Cancelled = 0 AND f.DepDelay IS NOT NULL GROUP BY f.OriginCode, d.DisplayAirportName, f.OriginCityName HAVING count() \u003e= 100000 ORDER BY otp_rate ASC",
      "row_count": 161,
      "result_columns": [
        "rank",
        "OriginCode",
        "airport_name",
        "city",
        "total_departures",
        "otp_rate",
        "pct_delayed",
        "avg_dep_delay_when_late",
        "p90_delay_min",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "OriginCode": "MDW",
        "airport_name": "Chicago Midway International",
        "avg_dep_delay_when_late": 26.82,
        "city": "Chicago, IL",
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "otp_rate": 76.52,
        "p90_delay_min": 37,
        "pct_delayed": 23.48,
        "rank": 1,
        "total_departures": 2512354
      }
    },
    {
      "id": "q2",
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "The spread between the worst airport and the median of the 160 qualifying airports is 8.53 percentage points. MDW sits at 76.52% OTP while the median airport in the ranked set sits at 85.05% OTP. This gap is substantial: MDW's delay rate (23.48%) is roughly 1.6x the median delay rate (14.95%), indicating that the very bottom of the distribution is a clear outlier rather than a gradual tail.",
      "sql": "WITH ranked AS (SELECT OriginCode, round(100 - avg(DepDel15) * 100, 2) AS otp_rate, count() AS total_departures FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY OriginCode HAVING total_departures \u003e= 100000) SELECT min(otp_rate) AS worst_otp, round(median(otp_rate), 2) AS median_otp, round(median(otp_rate) - min(otp_rate), 2) AS spread_pp, count() AS qualifying_airports FROM ranked",
      "row_count": 1,
      "result_columns": [
        "worst_otp",
        "median_otp",
        "spread_pp",
        "qualifying_airports"
      ],
      "first_row": {
        "median_otp": 85.05,
        "qualifying_airports": 160,
        "spread_pp": 8.53,
        "worst_otp": 76.52
      }
    },
    {
      "id": "q3",
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The bottom 10 airports are overwhelmingly major hubs and high-traffic gateways. MDW (Chicago Midway), ORD (Chicago O'Hare), EWR (Newark Liberty), JFK (John F. Kennedy), DEN (Denver), and LAS (Las Vegas) are among the nation's busiest airports by operations. HOU (Houston Hobby) and DAL (Dallas Love Field) are primary Southwest Airlines intra-Texas hubs. BWI (Baltimore/Washington) and FLL (Fort Lauderdale) are large leisure and low-cost-carrier gateways. Every airport in the bottom 10 carries well over 1 million historical departures, and none are regional or small-market airports. The pattern strongly suggests that operational complexity, congestion, and constrained airspace — not obscurity — drive the worst OTP outcomes.",
      "sql": "SELECT f.OriginCode, coalesce(d.DisplayAirportName, f.OriginCode) AS airport_name, f.OriginCityName AS city, count() AS total_departures, round(100 - avg(f.DepDel15) * 100, 2) AS otp_rate, round(avg(f.DepDel15) * 100, 2) AS pct_delayed, round(avgIf(f.DepDelay, f.DepDelay \u003e 0), 2) AS avg_dep_delay_when_late, quantileExact(0.90)(f.DepDelayMinutes) AS p90_delay_min FROM ontime.fact_ontime f LEFT JOIN ontime.dim_airports d ON f.OriginAirportID = d.AirportID WHERE f.Cancelled = 0 AND f.DepDelay IS NOT NULL GROUP BY f.OriginCode, d.DisplayAirportName, f.OriginCityName HAVING count() \u003e= 100000 ORDER BY otp_rate ASC LIMIT 10",
      "row_count": 10,
      "result_columns": [
        "OriginCode",
        "airport_name",
        "city",
        "total_departures",
        "otp_rate",
        "pct_delayed",
        "avg_dep_delay_when_late",
        "p90_delay_min"
      ],
      "first_row": {
        "OriginCode": "MDW",
        "airport_name": "Chicago Midway International",
        "avg_dep_delay_when_late": 26.82,
        "city": "Chicago, IL",
        "otp_rate": 76.52,
        "p90_delay_min": 37,
        "pct_delayed": 23.48,
        "total_departures": 2512354
      }
    }
  ]
}

### Multi-query additions

- The saved SQL shown below is the primary dashboard query for this page.
- The verified analysis package includes named supporting queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Use subquestion answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- The dashboard does not need to mirror `report.md`; it should combine narrative and interactive analysis.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{
  "question_title": "Worst origin airports by departure on-time performance",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which airport ranks worst on departure on-time performance?",
      "answer_markdown": "Chicago Midway International (MDW) ranks worst among all qualifying origin airports, with an on-time departure rate of 76.52% — meaning nearly 1 in 4 departures was delayed 15 or more minutes. MDW logged 2.51 million completed departures across the full dataset (1987–2025) and, when delayed, averaged 26.8 minutes late with a 90th-percentile delay of 37 minutes. The minimum-volume threshold applied was 100,000 completed departures, leaving 160 qualifying airports in the ranked set.",
      "sql": "SELECT row_number() OVER (ORDER BY round(100 - avg(f.DepDel15) * 100, 2) ASC) AS rank, f.OriginCode, coalesce(d.DisplayAirportName, f.OriginCode) AS airport_name, f.OriginCityName AS city, count() AS total_departures, round(100 - avg(f.DepDel15) * 100, 2) AS otp_rate, round(avg(f.DepDel15) * 100, 2) AS pct_delayed, round(avgIf(f.DepDelay, f.DepDelay \u003e 0), 2) AS avg_dep_delay_when_late, quantileExact(0.90)(f.DepDelayMinutes) AS p90_delay_min, min(f.FlightDate) AS first_date, max(f.FlightDate) AS last_date FROM ontime.fact_ontime f LEFT JOIN ontime.dim_airports d ON f.OriginAirportID = d.AirportID WHERE f.Cancelled = 0 AND f.DepDelay IS NOT NULL GROUP BY f.OriginCode, d.DisplayAirportName, f.OriginCityName HAVING count() \u003e= 100000 ORDER BY otp_rate ASC",
      "row_count": 161,
      "result_columns": [
        "rank",
        "OriginCode",
        "airport_name",
        "city",
        "total_departures",
        "otp_rate",
        "pct_delayed",
        "avg_dep_delay_when_late",
        "p90_delay_min",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "OriginCode": "MDW",
        "airport_name": "Chicago Midway International",
        "avg_dep_delay_when_late": 26.82,
        "city": "Chicago, IL",
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "otp_rate": 76.52,
        "p90_delay_min": 37,
        "pct_delayed": 23.48,
        "rank": 1,
        "total_departures": 2512354
      }
    },
    {
      "id": "q2",
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "The spread between the worst airport and the median of the 160 qualifying airports is 8.53 percentage points. MDW sits at 76.52% OTP while the median airport in the ranked set sits at 85.05% OTP. This gap is substantial: MDW's delay rate (23.48%) is roughly 1.6x the median delay rate (14.95%), indicating that the very bottom of the distribution is a clear outlier rather than a gradual tail.",
      "sql": "WITH ranked AS (SELECT OriginCode, round(100 - avg(DepDel15) * 100, 2) AS otp_rate, count() AS total_departures FROM ontime.fact_ontime WHERE Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY OriginCode HAVING total_departures \u003e= 100000) SELECT min(otp_rate) AS worst_otp, round(median(otp_rate), 2) AS median_otp, round(median(otp_rate) - min(otp_rate), 2) AS spread_pp, count() AS qualifying_airports FROM ranked",
      "row_count": 1,
      "result_columns": [
        "worst_otp",
        "median_otp",
        "spread_pp",
        "qualifying_airports"
      ],
      "first_row": {
        "median_otp": 85.05,
        "qualifying_airports": 160,
        "spread_pp": 8.53,
        "worst_otp": 76.52
      }
    },
    {
      "id": "q3",
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The bottom 10 airports are overwhelmingly major hubs and high-traffic gateways. MDW (Chicago Midway), ORD (Chicago O'Hare), EWR (Newark Liberty), JFK (John F. Kennedy), DEN (Denver), and LAS (Las Vegas) are among the nation's busiest airports by operations. HOU (Houston Hobby) and DAL (Dallas Love Field) are primary Southwest Airlines intra-Texas hubs. BWI (Baltimore/Washington) and FLL (Fort Lauderdale) are large leisure and low-cost-carrier gateways. Every airport in the bottom 10 carries well over 1 million historical departures, and none are regional or small-market airports. The pattern strongly suggests that operational complexity, congestion, and constrained airspace — not obscurity — drive the worst OTP outcomes.",
      "sql": "SELECT f.OriginCode, coalesce(d.DisplayAirportName, f.OriginCode) AS airport_name, f.OriginCityName AS city, count() AS total_departures, round(100 - avg(f.DepDel15) * 100, 2) AS otp_rate, round(avg(f.DepDel15) * 100, 2) AS pct_delayed, round(avgIf(f.DepDelay, f.DepDelay \u003e 0), 2) AS avg_dep_delay_when_late, quantileExact(0.90)(f.DepDelayMinutes) AS p90_delay_min FROM ontime.fact_ontime f LEFT JOIN ontime.dim_airports d ON f.OriginAirportID = d.AirportID WHERE f.Cancelled = 0 AND f.DepDelay IS NOT NULL GROUP BY f.OriginCode, d.DisplayAirportName, f.OriginCityName HAVING count() \u003e= 100000 ORDER BY otp_rate ASC LIMIT 10",
      "row_count": 10,
      "result_columns": [
        "OriginCode",
        "airport_name",
        "city",
        "total_departures",
        "otp_rate",
        "pct_delayed",
        "avg_dep_delay_when_late",
        "p90_delay_min"
      ],
      "first_row": {
        "OriginCode": "MDW",
        "airport_name": "Chicago Midway International",
        "avg_dep_delay_when_late": 26.82,
        "city": "Chicago, IL",
        "otp_rate": 76.52,
        "p90_delay_min": 37,
        "pct_delayed": 23.48,
        "total_departures": 2512354
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.

Create browser-ready HTML `visual.html`.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.