- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `American Airlines peak network delay month and contributors`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_contribution_dashboard`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

Build a dashboard that:

- uses `peak_month` as the primary saved SQL already provided in the prompt
- uses `origin_contributors`, `route_contributors`, and `concentration_pattern` as supporting queries when they materially improve the dashboard
- shows KPI cards for peak month, peak average departure delay, peak `% DepDel15`, and completed flights in the peak month
- renders a monthly time-series chart of average `DepDelayMinutes` across all AA months
- visually highlights the peak month on that chart
- renders a bar chart of top origin contributors within the peak month
- renders a route contribution table for the peak month
- includes a narrative takeaway about whether the peak month was broad across the network or concentrated in a smaller set of origins and routes
- derives the peak month from fetched monthly data instead of hardcoding it
- clearly separates the network-wide trend from the peak-month drilldown
- annotates the peak month on the time series
- makes the contribution logic easy to read without external narrative
- shows supporting queries in the query ledger when used

### Data Source

SQL query for primary data source:

```sql
SELECT Year, Month, count() AS FlightVolume, round(avg(DepDelay), 2) AS AvgDepDelay, round(sum(DepDel15) / count(), 4) AS PctDel15 FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 GROUP BY Year, Month ORDER BY AvgDepDelay DESC;
```

Data example/snippet:

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

### Multi-query additions

- The saved SQL shown below is the primary dashboard query for this page.
- The verified analysis package includes named supporting queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Use subquestion answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- The dashboard does not need to mirror `report.md`; it should combine narrative and interactive analysis.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

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

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.

Create browser-ready HTML `visual.html`.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.