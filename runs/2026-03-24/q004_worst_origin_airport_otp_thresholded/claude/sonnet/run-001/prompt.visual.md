- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create browser-ready HTML `visual.html` using the proper  `*-analyst-dashboard` skill.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.

### Rules

- Question title: `Worst origin airports by departure on-time performance`
- Visual mode: `dynamic`
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
WITH airport_stats AS (
    SELECT
        OriginAirportID,
        OriginCode,
        count() AS completed_departures,
        round(countIf(DepDel15 = 0) * 100.0 / count(), 2) AS otp_pct,
        round(avg(DepDelay), 2) AS avg_dep_delay_min,
        round(quantile(0.90)(DepDelay), 2) AS p90_dep_delay_min,
        min(FlightDate) AS first_date,
        max(FlightDate) AS last_date
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
    GROUP BY OriginAirportID, OriginCode
    HAVING completed_departures >= 10000
),
ranked AS (
    SELECT
        s.*,
        d.DisplayAirportName AS airport_name,
        row_number() OVER (ORDER BY otp_pct ASC) AS rank_worst,
        count() OVER () AS total_qualifying
    FROM airport_stats s
    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID
)
SELECT
    rank_worst,
    total_qualifying,
    OriginCode,
    airport_name,
    completed_departures,
    otp_pct,
    avg_dep_delay_min,
    p90_dep_delay_min,
    first_date,
    last_date
FROM ranked
ORDER BY rank_worst ASC
LIMIT 30
```

Data example/snippet:

{
  "question_title": "Worst origin airports by departure on-time performance",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "worst_airport",
      "subquestion": "Which airport ranks worst on departure on-time performance?",
      "answer_markdown": "Jack McNamara Field (CEC, Crescent City, CA) ranks worst among the 304 qualifying origin airports (those with at least 10,000 completed departures). It has an OTP of 69.24%, meaning nearly 1 in 3 departures is delayed 15 or more minutes. Its average departure delay is 19.5 minutes and the 90th-percentile delay is 84 minutes — the highest in the dataset. The next-worst airports are Aspen Pitkin County Sardy Field (ASE, 75.09%) and Nantucket Memorial (ACK, 75.77%). The first major hub to appear in the ranking is Chicago Midway International (MDW, rank 4, 76.52%), followed by Chicago O'Hare (ORD, rank 9, 78.61%) and Newark Liberty (EWR, rank 11, 78.69%).",
      "sql": "WITH airport_stats AS (\n    SELECT\n        OriginAirportID,\n        OriginCode,\n        count() AS completed_departures,\n        round(countIf(DepDel15 = 0) * 100.0 / count(), 2) AS otp_pct,\n        round(avg(DepDelay), 2) AS avg_dep_delay_min,\n        round(quantile(0.90)(DepDelay), 2) AS p90_dep_delay_min,\n        min(FlightDate) AS first_date,\n        max(FlightDate) AS last_date\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY OriginAirportID, OriginCode\n    HAVING completed_departures \u003e= 10000\n),\nranked AS (\n    SELECT\n        s.*,\n        d.DisplayAirportName AS airport_name,\n        row_number() OVER (ORDER BY otp_pct ASC) AS rank_worst,\n        count() OVER () AS total_qualifying\n    FROM airport_stats s\n    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID\n)\nSELECT\n    rank_worst,\n    total_qualifying,\n    OriginCode,\n    airport_name,\n    completed_departures,\n    otp_pct,\n    avg_dep_delay_min,\n    p90_dep_delay_min,\n    first_date,\n    last_date\nFROM ranked\nORDER BY rank_worst ASC\nLIMIT 30",
      "row_count": 30,
      "result_columns": [
        "rank_worst",
        "total_qualifying",
        "OriginCode",
        "airport_name",
        "completed_departures",
        "otp_pct",
        "avg_dep_delay_min",
        "p90_dep_delay_min",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "OriginCode": "CEC",
        "airport_name": "Jack McNamara Field",
        "avg_dep_delay_min": 19.5,
        "completed_departures": 11891,
        "first_date": "2003-01-01T00:00:00Z",
        "last_date": "2015-04-06T00:00:00Z",
        "otp_pct": 69.24,
        "p90_dep_delay_min": 83,
        "rank_worst": 1,
        "total_qualifying": 304
      }
    },
    {
      "id": "spread_to_middle",
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "Among 304 qualifying airports, the middle-ranked airport (rank 152) is Pensacola International (PNS) with an OTP of 85.31%. The worst airport, Jack McNamara Field (CEC), has an OTP of 69.24%. The spread is 16.07 percentage points — meaning CEC departs on time 16 points less often than a typical mid-tier airport. This is a substantial gap, reflecting that CEC's performance is genuinely extreme rather than marginally below average.",
      "sql": "WITH airport_stats AS (\n    SELECT\n        OriginAirportID,\n        OriginCode,\n        count() AS completed_departures,\n        round(countIf(DepDel15 = 0) * 100.0 / count(), 2) AS otp_pct,\n        round(avg(DepDelay), 2) AS avg_dep_delay_min,\n        round(quantile(0.90)(DepDelay), 2) AS p90_dep_delay_min,\n        min(FlightDate) AS first_date,\n        max(FlightDate) AS last_date\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY OriginAirportID, OriginCode\n    HAVING completed_departures \u003e= 10000\n),\nranked AS (\n    SELECT\n        s.*,\n        d.DisplayAirportName AS airport_name,\n        row_number() OVER (ORDER BY otp_pct ASC) AS rank_worst,\n        count() OVER () AS total_qualifying\n    FROM airport_stats s\n    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID\n)\nSELECT\n    rank_worst,\n    total_qualifying,\n    OriginCode,\n    airport_name,\n    completed_departures,\n    otp_pct,\n    avg_dep_delay_min,\n    p90_dep_delay_min,\n    round(otp_pct - min(otp_pct) OVER (), 2) AS spread_from_worst_pct\nFROM ranked\nWHERE rank_worst = 1\n   OR rank_worst = toUInt64(ceil(total_qualifying / 2.0))\nORDER BY rank_worst ASC",
      "row_count": 2,
      "result_columns": [
        "rank_worst",
        "total_qualifying",
        "OriginCode",
        "airport_name",
        "completed_departures",
        "otp_pct",
        "avg_dep_delay_min",
        "p90_dep_delay_min",
        "spread_from_worst_pct"
      ],
      "first_row": {
        "OriginCode": "CEC",
        "airport_name": "Jack McNamara Field",
        "avg_dep_delay_min": 19.5,
        "completed_departures": 11891,
        "otp_pct": 69.24,
        "p90_dep_delay_min": 84,
        "rank_worst": 1,
        "spread_from_worst_pct": 0,
        "total_qualifying": 304
      }
    },
    {
      "id": "bottom_group_mix",
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The bottom group is genuinely mixed. Among the 30 worst-performing qualifying airports, small and mid-size regional airports dominate the very bottom — Jack McNamara Field (CEC, 11k departures), Aspen Pitkin County Sardy Field (ASE, 99k), Nantucket Memorial (ACK, 15k), Modesto City-County (MOD, 18k), and Trenton Mercer (TTN, 28k) all rank in the worst 10. However, some of the busiest hubs in the country also appear throughout the bottom 30: Chicago Midway (MDW, 2.5M, rank 4), Chicago O'Hare (ORD, 11.1M, rank 9), Newark Liberty (EWR, 4.5M, rank 11), JFK (3.1M, rank 17), Baltimore/Washington (BWI, 3.3M, rank 18), Las Vegas (LAS, 5.2M, rank 20), Denver (DEN, 7.4M, rank 22), Fort Lauderdale (FLL, 2.3M, rank 23), San Francisco (SFO, 5.1M, rank 26), Dallas/Fort Worth (DFW, 10M, rank 27), Miami (MIA, 2.8M, rank 29), and Philadelphia (PHL, 3.6M, rank 30). The conclusion is that poor departure OTP is not limited to major hubs — the weakest performers span the full volume spectrum.",
      "sql": "WITH airport_stats AS (\n    SELECT\n        OriginAirportID,\n        OriginCode,\n        count() AS completed_departures,\n        round(countIf(DepDel15 = 0) * 100.0 / count(), 2) AS otp_pct,\n        round(avg(DepDelay), 2) AS avg_dep_delay_min,\n        round(quantile(0.90)(DepDelay), 2) AS p90_dep_delay_min,\n        min(FlightDate) AS first_date,\n        max(FlightDate) AS last_date\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY OriginAirportID, OriginCode\n    HAVING completed_departures \u003e= 10000\n),\nranked AS (\n    SELECT\n        s.*,\n        d.DisplayAirportName AS airport_name,\n        row_number() OVER (ORDER BY otp_pct ASC) AS rank_worst,\n        count() OVER () AS total_qualifying,\n        multiIf(\n            s.completed_departures \u003e= 2000000, 'Major hub (2M+ dep)',\n            s.completed_departures \u003e= 500000,  'Large airport (500k-2M dep)',\n            s.completed_departures \u003e= 100000,  'Mid-size (100k-500k dep)',\n            'Regional/small (\u003c100k dep)'\n        ) AS volume_tier\n    FROM airport_stats s\n    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID\n)\nSELECT\n    rank_worst,\n    OriginCode,\n    airport_name,\n    completed_departures,\n    volume_tier,\n    otp_pct,\n    avg_dep_delay_min,\n    p90_dep_delay_min\nFROM ranked\nWHERE rank_worst \u003c= 30\nORDER BY rank_worst ASC",
      "row_count": 30,
      "result_columns": [
        "rank_worst",
        "OriginCode",
        "airport_name",
        "completed_departures",
        "volume_tier",
        "otp_pct",
        "avg_dep_delay_min",
        "p90_dep_delay_min"
      ],
      "first_row": {
        "OriginCode": "CEC",
        "airport_name": "Jack McNamara Field",
        "avg_dep_delay_min": 19.5,
        "completed_departures": 11891,
        "otp_pct": 69.24,
        "p90_dep_delay_min": 83,
        "rank_worst": 1,
        "volume_tier": "Regional/small (\u003c100k dep)"
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
      "id": "worst_airport",
      "subquestion": "Which airport ranks worst on departure on-time performance?",
      "answer_markdown": "Jack McNamara Field (CEC, Crescent City, CA) ranks worst among the 304 qualifying origin airports (those with at least 10,000 completed departures). It has an OTP of 69.24%, meaning nearly 1 in 3 departures is delayed 15 or more minutes. Its average departure delay is 19.5 minutes and the 90th-percentile delay is 84 minutes — the highest in the dataset. The next-worst airports are Aspen Pitkin County Sardy Field (ASE, 75.09%) and Nantucket Memorial (ACK, 75.77%). The first major hub to appear in the ranking is Chicago Midway International (MDW, rank 4, 76.52%), followed by Chicago O'Hare (ORD, rank 9, 78.61%) and Newark Liberty (EWR, rank 11, 78.69%).",
      "sql": "WITH airport_stats AS (\n    SELECT\n        OriginAirportID,\n        OriginCode,\n        count() AS completed_departures,\n        round(countIf(DepDel15 = 0) * 100.0 / count(), 2) AS otp_pct,\n        round(avg(DepDelay), 2) AS avg_dep_delay_min,\n        round(quantile(0.90)(DepDelay), 2) AS p90_dep_delay_min,\n        min(FlightDate) AS first_date,\n        max(FlightDate) AS last_date\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY OriginAirportID, OriginCode\n    HAVING completed_departures \u003e= 10000\n),\nranked AS (\n    SELECT\n        s.*,\n        d.DisplayAirportName AS airport_name,\n        row_number() OVER (ORDER BY otp_pct ASC) AS rank_worst,\n        count() OVER () AS total_qualifying\n    FROM airport_stats s\n    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID\n)\nSELECT\n    rank_worst,\n    total_qualifying,\n    OriginCode,\n    airport_name,\n    completed_departures,\n    otp_pct,\n    avg_dep_delay_min,\n    p90_dep_delay_min,\n    first_date,\n    last_date\nFROM ranked\nORDER BY rank_worst ASC\nLIMIT 30",
      "row_count": 30,
      "result_columns": [
        "rank_worst",
        "total_qualifying",
        "OriginCode",
        "airport_name",
        "completed_departures",
        "otp_pct",
        "avg_dep_delay_min",
        "p90_dep_delay_min",
        "first_date",
        "last_date"
      ],
      "first_row": {
        "OriginCode": "CEC",
        "airport_name": "Jack McNamara Field",
        "avg_dep_delay_min": 19.5,
        "completed_departures": 11891,
        "first_date": "2003-01-01T00:00:00Z",
        "last_date": "2015-04-06T00:00:00Z",
        "otp_pct": 69.24,
        "p90_dep_delay_min": 83,
        "rank_worst": 1,
        "total_qualifying": 304
      }
    },
    {
      "id": "spread_to_middle",
      "subquestion": "How large is the spread between the worst airport and the middle of the ranked set?",
      "answer_markdown": "Among 304 qualifying airports, the middle-ranked airport (rank 152) is Pensacola International (PNS) with an OTP of 85.31%. The worst airport, Jack McNamara Field (CEC), has an OTP of 69.24%. The spread is 16.07 percentage points — meaning CEC departs on time 16 points less often than a typical mid-tier airport. This is a substantial gap, reflecting that CEC's performance is genuinely extreme rather than marginally below average.",
      "sql": "WITH airport_stats AS (\n    SELECT\n        OriginAirportID,\n        OriginCode,\n        count() AS completed_departures,\n        round(countIf(DepDel15 = 0) * 100.0 / count(), 2) AS otp_pct,\n        round(avg(DepDelay), 2) AS avg_dep_delay_min,\n        round(quantile(0.90)(DepDelay), 2) AS p90_dep_delay_min,\n        min(FlightDate) AS first_date,\n        max(FlightDate) AS last_date\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY OriginAirportID, OriginCode\n    HAVING completed_departures \u003e= 10000\n),\nranked AS (\n    SELECT\n        s.*,\n        d.DisplayAirportName AS airport_name,\n        row_number() OVER (ORDER BY otp_pct ASC) AS rank_worst,\n        count() OVER () AS total_qualifying\n    FROM airport_stats s\n    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID\n)\nSELECT\n    rank_worst,\n    total_qualifying,\n    OriginCode,\n    airport_name,\n    completed_departures,\n    otp_pct,\n    avg_dep_delay_min,\n    p90_dep_delay_min,\n    round(otp_pct - min(otp_pct) OVER (), 2) AS spread_from_worst_pct\nFROM ranked\nWHERE rank_worst = 1\n   OR rank_worst = toUInt64(ceil(total_qualifying / 2.0))\nORDER BY rank_worst ASC",
      "row_count": 2,
      "result_columns": [
        "rank_worst",
        "total_qualifying",
        "OriginCode",
        "airport_name",
        "completed_departures",
        "otp_pct",
        "avg_dep_delay_min",
        "p90_dep_delay_min",
        "spread_from_worst_pct"
      ],
      "first_row": {
        "OriginCode": "CEC",
        "airport_name": "Jack McNamara Field",
        "avg_dep_delay_min": 19.5,
        "completed_departures": 11891,
        "otp_pct": 69.24,
        "p90_dep_delay_min": 84,
        "rank_worst": 1,
        "spread_from_worst_pct": 0,
        "total_qualifying": 304
      }
    },
    {
      "id": "bottom_group_mix",
      "subquestion": "Are the weakest airports mostly major hubs, or is the bottom group more mixed?",
      "answer_markdown": "The bottom group is genuinely mixed. Among the 30 worst-performing qualifying airports, small and mid-size regional airports dominate the very bottom — Jack McNamara Field (CEC, 11k departures), Aspen Pitkin County Sardy Field (ASE, 99k), Nantucket Memorial (ACK, 15k), Modesto City-County (MOD, 18k), and Trenton Mercer (TTN, 28k) all rank in the worst 10. However, some of the busiest hubs in the country also appear throughout the bottom 30: Chicago Midway (MDW, 2.5M, rank 4), Chicago O'Hare (ORD, 11.1M, rank 9), Newark Liberty (EWR, 4.5M, rank 11), JFK (3.1M, rank 17), Baltimore/Washington (BWI, 3.3M, rank 18), Las Vegas (LAS, 5.2M, rank 20), Denver (DEN, 7.4M, rank 22), Fort Lauderdale (FLL, 2.3M, rank 23), San Francisco (SFO, 5.1M, rank 26), Dallas/Fort Worth (DFW, 10M, rank 27), Miami (MIA, 2.8M, rank 29), and Philadelphia (PHL, 3.6M, rank 30). The conclusion is that poor departure OTP is not limited to major hubs — the weakest performers span the full volume spectrum.",
      "sql": "WITH airport_stats AS (\n    SELECT\n        OriginAirportID,\n        OriginCode,\n        count() AS completed_departures,\n        round(countIf(DepDel15 = 0) * 100.0 / count(), 2) AS otp_pct,\n        round(avg(DepDelay), 2) AS avg_dep_delay_min,\n        round(quantile(0.90)(DepDelay), 2) AS p90_dep_delay_min,\n        min(FlightDate) AS first_date,\n        max(FlightDate) AS last_date\n    FROM ontime.fact_ontime\n    WHERE Cancelled = 0\n    GROUP BY OriginAirportID, OriginCode\n    HAVING completed_departures \u003e= 10000\n),\nranked AS (\n    SELECT\n        s.*,\n        d.DisplayAirportName AS airport_name,\n        row_number() OVER (ORDER BY otp_pct ASC) AS rank_worst,\n        count() OVER () AS total_qualifying,\n        multiIf(\n            s.completed_departures \u003e= 2000000, 'Major hub (2M+ dep)',\n            s.completed_departures \u003e= 500000,  'Large airport (500k-2M dep)',\n            s.completed_departures \u003e= 100000,  'Mid-size (100k-500k dep)',\n            'Regional/small (\u003c100k dep)'\n        ) AS volume_tier\n    FROM airport_stats s\n    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID\n)\nSELECT\n    rank_worst,\n    OriginCode,\n    airport_name,\n    completed_departures,\n    volume_tier,\n    otp_pct,\n    avg_dep_delay_min,\n    p90_dep_delay_min\nFROM ranked\nWHERE rank_worst \u003c= 30\nORDER BY rank_worst ASC",
      "row_count": 30,
      "result_columns": [
        "rank_worst",
        "OriginCode",
        "airport_name",
        "completed_departures",
        "volume_tier",
        "otp_pct",
        "avg_dep_delay_min",
        "p90_dep_delay_min"
      ],
      "first_row": {
        "OriginCode": "CEC",
        "airport_name": "Jack McNamara Field",
        "avg_dep_delay_min": 19.5,
        "completed_departures": 11891,
        "otp_pct": 69.24,
        "p90_dep_delay_min": 83,
        "rank_worst": 1,
        "volume_tier": "Regional/small (\u003c100k dep)"
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.