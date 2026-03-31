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
SELECT OriginCode, any(OriginCityName) AS OriginCityName, count() AS completed_departures, round(avg(DepDel15 = 0), 4) AS on_time_performance, round(avg(DepDelay), 2) AS avg_departure_delay, round(quantile(0.90)(DepDelay), 2) AS p90_departure_delay, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 GROUP BY OriginCode HAVING completed_departures >= 100000 ORDER BY on_time_performance ASC LIMIT 500
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
      "answer_markdown": "Chicago Midway International (MDW) ranks worst on departure on-time performance among qualifying airports. Over the entire historical dataset, MDW recorded an on-time performance of only 76.5%, significantly lower than its peers.",
      "sql": "SELECT OriginCode, any(OriginCityName) AS OriginCityName, count() AS completed_departures, round(avg(DepDel15 = 0), 4) AS on_time_performance, round(avg(DepDelay), 2) AS avg_departure_delay, round(quantile(0.90)(DepDelay), 2) AS p90_departure_delay, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 GROUP BY OriginCode HAVING completed_departures \u003e= 100000 ORDER BY on_time_performance ASC LIMIT 500",
      "row_count": 160,
      "result_columns": [
        "OriginCode",
        "OriginCityName",
        "completed_departures",
        "on_time_performance",
        "avg_departure_delay",
        "p90_departure_delay",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "OriginCityName": "Chicago, IL",
        "OriginCode": "MDW",
        "avg_departure_delay": 11.8,
        "completed_departures": 2512354,
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "on_time_performance": 0.7652,
        "p90_departure_delay": 38
      }
    },
    {
      "id": "q2",
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "There is an 8.5 percentage point gap between the worst airport and the median of the ranked set. While Chicago Midway (MDW) sits at the bottom with a 76.5% on-time performance, an average qualifying airport in the middle of the rankings (such as Islip, NY) achieves an on-time performance of approximately 85.0%.",
      "sql": "WITH airport_stats AS (SELECT OriginCode, any(OriginCityName) AS OriginCityName, count() AS completed_departures, round(avg(DepDel15 = 0), 4) AS on_time_performance, round(avg(DepDelay), 2) AS avg_departure_delay, round(quantile(0.90)(DepDelay), 2) AS p90_departure_delay, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 GROUP BY OriginCode HAVING completed_departures \u003e= 100000), ranked_airports AS (SELECT *, row_number() OVER (ORDER BY on_time_performance ASC) AS rnk, count() OVER () AS total_airports FROM airport_stats) SELECT * FROM ranked_airports WHERE rnk = 1 OR rnk = intDiv(total_airports, 2)",
      "row_count": 2,
      "result_columns": [
        "OriginCode",
        "OriginCityName",
        "completed_departures",
        "on_time_performance",
        "avg_departure_delay",
        "p90_departure_delay",
        "first_date",
        "last_date",
        "rnk",
        "total_airports"
      ],
      "first_row": {
        "OriginCityName": "Chicago, IL",
        "OriginCode": "MDW",
        "avg_departure_delay": 11.8,
        "completed_departures": 2512354,
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "on_time_performance": 0.7652,
        "p90_departure_delay": 36,
        "rnk": 1,
        "total_airports": 160
      }
    },
    {
      "id": "q3",
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The weakest performing airports are overwhelmingly major airline hubs. The bottom tier of the rankings is dominated by heavily trafficked hubs such as Chicago O'Hare (ORD), Newark (EWR), Houston (HOU), Dallas (DAL), JFK, and Baltimore (BWI), rather than being a mixed group of large and mid-sized regional airports.",
      "sql": "SELECT OriginCode, any(OriginCityName) AS OriginCityName, count() AS completed_departures, round(avg(DepDel15 = 0), 4) AS on_time_performance, round(avg(DepDelay), 2) AS avg_departure_delay, round(quantile(0.90)(DepDelay), 2) AS p90_departure_delay, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 GROUP BY OriginCode HAVING completed_departures \u003e= 100000 ORDER BY on_time_performance ASC LIMIT 20",
      "row_count": 20,
      "result_columns": [
        "OriginCode",
        "OriginCityName",
        "completed_departures",
        "on_time_performance",
        "avg_departure_delay",
        "p90_departure_delay",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "OriginCityName": "Chicago, IL",
        "OriginCode": "MDW",
        "avg_departure_delay": 11.8,
        "completed_departures": 2512354,
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "on_time_performance": 0.7652,
        "p90_departure_delay": 35
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
      "answer_markdown": "Chicago Midway International (MDW) ranks worst on departure on-time performance among qualifying airports. Over the entire historical dataset, MDW recorded an on-time performance of only 76.5%, significantly lower than its peers.",
      "sql": "SELECT OriginCode, any(OriginCityName) AS OriginCityName, count() AS completed_departures, round(avg(DepDel15 = 0), 4) AS on_time_performance, round(avg(DepDelay), 2) AS avg_departure_delay, round(quantile(0.90)(DepDelay), 2) AS p90_departure_delay, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 GROUP BY OriginCode HAVING completed_departures \u003e= 100000 ORDER BY on_time_performance ASC LIMIT 500",
      "row_count": 160,
      "result_columns": [
        "OriginCode",
        "OriginCityName",
        "completed_departures",
        "on_time_performance",
        "avg_departure_delay",
        "p90_departure_delay",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "OriginCityName": "Chicago, IL",
        "OriginCode": "MDW",
        "avg_departure_delay": 11.8,
        "completed_departures": 2512354,
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "on_time_performance": 0.7652,
        "p90_departure_delay": 38
      }
    },
    {
      "id": "q2",
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "There is an 8.5 percentage point gap between the worst airport and the median of the ranked set. While Chicago Midway (MDW) sits at the bottom with a 76.5% on-time performance, an average qualifying airport in the middle of the rankings (such as Islip, NY) achieves an on-time performance of approximately 85.0%.",
      "sql": "WITH airport_stats AS (SELECT OriginCode, any(OriginCityName) AS OriginCityName, count() AS completed_departures, round(avg(DepDel15 = 0), 4) AS on_time_performance, round(avg(DepDelay), 2) AS avg_departure_delay, round(quantile(0.90)(DepDelay), 2) AS p90_departure_delay, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 GROUP BY OriginCode HAVING completed_departures \u003e= 100000), ranked_airports AS (SELECT *, row_number() OVER (ORDER BY on_time_performance ASC) AS rnk, count() OVER () AS total_airports FROM airport_stats) SELECT * FROM ranked_airports WHERE rnk = 1 OR rnk = intDiv(total_airports, 2)",
      "row_count": 2,
      "result_columns": [
        "OriginCode",
        "OriginCityName",
        "completed_departures",
        "on_time_performance",
        "avg_departure_delay",
        "p90_departure_delay",
        "first_date",
        "last_date",
        "rnk",
        "total_airports"
      ],
      "first_row": {
        "OriginCityName": "Chicago, IL",
        "OriginCode": "MDW",
        "avg_departure_delay": 11.8,
        "completed_departures": 2512354,
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "on_time_performance": 0.7652,
        "p90_departure_delay": 36,
        "rnk": 1,
        "total_airports": 160
      }
    },
    {
      "id": "q3",
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The weakest performing airports are overwhelmingly major airline hubs. The bottom tier of the rankings is dominated by heavily trafficked hubs such as Chicago O'Hare (ORD), Newark (EWR), Houston (HOU), Dallas (DAL), JFK, and Baltimore (BWI), rather than being a mixed group of large and mid-sized regional airports.",
      "sql": "SELECT OriginCode, any(OriginCityName) AS OriginCityName, count() AS completed_departures, round(avg(DepDel15 = 0), 4) AS on_time_performance, round(avg(DepDelay), 2) AS avg_departure_delay, round(quantile(0.90)(DepDelay), 2) AS p90_departure_delay, min(FlightDate) AS first_date, max(FlightDate) AS last_date FROM ontime.fact_ontime WHERE Cancelled = 0 GROUP BY OriginCode HAVING completed_departures \u003e= 100000 ORDER BY on_time_performance ASC LIMIT 20",
      "row_count": 20,
      "result_columns": [
        "OriginCode",
        "OriginCityName",
        "completed_departures",
        "on_time_performance",
        "avg_departure_delay",
        "p90_departure_delay",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "OriginCityName": "Chicago, IL",
        "OriginCode": "MDW",
        "avg_departure_delay": 11.8,
        "completed_departures": 2512354,
        "first_date": "1987-10-01T00:00:00Z",
        "last_date": "2025-11-30T00:00:00Z",
        "on_time_performance": 0.7652,
        "p90_departure_delay": 35
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