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
WITH monthly AS ( SELECT toStartOfMonth(FlightDate) AS month, count() AS flights, round(avg(DepDelay), 2) AS avg_dep_delay, round(avg(DepDel15) * 100, 2) AS pct_dep_15_plus FROM ontime.fact_ontime WHERE Carrier = 'AA' AND Cancelled = 0 AND DepDelay IS NOT NULL GROUP BY month ) SELECT month, flights, avg_dep_delay, pct_dep_15_plus, row_number() OVER (ORDER BY avg_dep_delay DESC, pct_dep_15_plus DESC, flights DESC, month DESC) AS worst_month_rank FROM monthly ORDER BY month
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

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.

Create browser-ready HTML `visual.html`.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.