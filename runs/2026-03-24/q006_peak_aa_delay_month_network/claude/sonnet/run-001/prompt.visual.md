- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create browser-ready HTML `visual.html` using the proper  `*-analyst-dashboard` skill.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.

### Rules

- Question title: `American Airlines peak network delay month and contributors`
- Visual mode: `dynamic`
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
SELECT
    toStartOfMonth(FlightDate) AS flight_month,
    count() AS total_flights,
    round(avg(DepDelay), 2) AS avg_dep_delay_min,
    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15
FROM ontime.fact_ontime
WHERE Reporting_Airline = 'AA'
  AND Cancelled = 0
  AND DepDelay IS NOT NULL
GROUP BY flight_month
ORDER BY flight_month
```

Data example/snippet:

{
  "question_title": "American Airlines peak network delay month and contributors",
  "result_columns": null,
  "row_count": 4,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "peak_month",
      "subquestion": "Which month is the single worst American Airlines month for departure delays?",
      "answer_markdown": "**July 2024** is the single worst month for American Airlines departure delays across the full dataset history. It averaged **36.33 minutes** of departure delay per flight and had **38.04%** of its 86,083 completed flights depart 15 or more minutes late — both figures higher than any other calendar month on record. The next closest months are May 2024 (32.74 min avg) and July 2025 (31.35 min avg), making July 2024 a clear standout.",
      "sql": "SELECT\n    toStartOfMonth(FlightDate) AS flight_month,\n    count() AS total_flights,\n    round(avg(DepDelay), 2) AS avg_dep_delay_min,\n    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE Reporting_Airline = 'AA'\n  AND Cancelled = 0\n  AND DepDelay IS NOT NULL\nGROUP BY flight_month\nORDER BY flight_month",
      "row_count": 458,
      "result_columns": [
        "flight_month",
        "total_flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "avg_dep_delay_min": 4.27,
        "flight_month": "1987-10-01T00:00:00Z",
        "pct_dep_del15": 8.8,
        "total_flights": 55871
      }
    },
    {
      "id": "origin_contributors",
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "**Charlotte (CLT)** is by far the dominant origin contributor in July 2024, with 10,568 departures — roughly 12% of all AA July 2024 flights — averaging **52.0 minutes** of departure delay and 52.2% of flights delayed 15+ minutes. No other origin comes close in combined volume and severity. Among smaller but notably delayed origins: IAH (632 flights, 47.1 min avg), EWR (615 flights, 44.9 min avg), RSW (354 flights, 45.3 min avg), and CVG (312 flights, 47.4 min avg) also stand out. CLT's massive scale makes it the single largest driver of AA's July 2024 network-wide delay figure.",
      "sql": "SELECT\n    OriginCode,\n    count() AS flights,\n    round(avg(DepDelay), 2) AS avg_dep_delay_min,\n    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE Reporting_Airline = 'AA'\n  AND Cancelled = 0\n  AND DepDelay IS NOT NULL\n  AND toStartOfMonth(FlightDate) = '2024-07-01'\nGROUP BY OriginCode\nHAVING flights \u003e= 100\nORDER BY avg_dep_delay_min DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "OriginCode",
        "flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "OriginCode": "EYW",
        "avg_dep_delay_min": 70,
        "flights": 119,
        "pct_dep_del15": 45.38
      }
    },
    {
      "id": "route_contributors",
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "The worst routes out of CLT led the July 2024 ranking: **CLT→SDF** averaged 102.7 minutes of departure delay (69.9% del15, 83 flights), **CLT→ROC** averaged 91.9 minutes (69.1% del15, 55 flights), and **CLT→TYS** averaged 90.0 minutes (69.0% del15, 58 flights). Other severe routes include **DFW→CID** (95.1 min, 60 flights), **JAC→DFW** (93.9 min, 63 flights), and **DFW→ICT** (93.4 min, 59 flights). CLT-origin routes dominate the worst performers, reinforcing its role as the network's delay epicenter that month. The PHX→JAX route had the single highest average at 134.4 minutes but on only 31 flights.",
      "sql": "SELECT\n    OriginCode,\n    DestCode,\n    count() AS flights,\n    round(avg(DepDelay), 2) AS avg_dep_delay_min,\n    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE Reporting_Airline = 'AA'\n  AND Cancelled = 0\n  AND DepDelay IS NOT NULL\n  AND toStartOfMonth(FlightDate) = '2024-07-01'\nGROUP BY OriginCode, DestCode\nHAVING flights \u003e= 30\nORDER BY avg_dep_delay_min DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "OriginCode",
        "DestCode",
        "flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "DestCode": "JAX",
        "OriginCode": "PHX",
        "avg_dep_delay_min": 134.42,
        "flights": 31,
        "pct_dep_del15": 45.16
      }
    },
    {
      "id": "concentration_pattern",
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "The July 2024 peak is **concentrated, not broad**. CLT alone accounted for ~10,568 of 86,083 AA departures (~12.3%) and averaged over 52 minutes of delay — well above the 36.3-minute network average — meaning CLT's underperformance pulled the entire network number up significantly. The top 5 origins by average delay (EYW, JAC, ALB, PNS, ECP) were all small stations with under 250 flights each. At the route level, the worst-performing routes are CLT-hub spokes and a handful of DFW feeder routes, not a broad spread across the full network. The concentration around CLT as AA's largest East Coast hub — where cascading delays compound across hundreds of daily departures — explains why July 2024 stands so far above other months in the historical record.",
      "sql": "SELECT\n    OriginCode,\n    count() AS flights,\n    round(100.0 * count() / sum(count()) OVER (), 2) AS pct_of_total_flights,\n    round(avg(DepDelay), 2) AS avg_dep_delay_min,\n    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE Reporting_Airline = 'AA'\n  AND Cancelled = 0\n  AND DepDelay IS NOT NULL\n  AND toStartOfMonth(FlightDate) = '2024-07-01'\nGROUP BY OriginCode\nORDER BY flights DESC\nLIMIT 20",
      "row_count": 20,
      "result_columns": [
        "OriginCode",
        "flights",
        "pct_of_total_flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "OriginCode": "DFW",
        "avg_dep_delay_min": 38.52,
        "flights": 14962,
        "pct_dep_del15": 46.59,
        "pct_of_total_flights": 17.38
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
      "id": "peak_month",
      "subquestion": "Which month is the single worst American Airlines month for departure delays?",
      "answer_markdown": "**July 2024** is the single worst month for American Airlines departure delays across the full dataset history. It averaged **36.33 minutes** of departure delay per flight and had **38.04%** of its 86,083 completed flights depart 15 or more minutes late — both figures higher than any other calendar month on record. The next closest months are May 2024 (32.74 min avg) and July 2025 (31.35 min avg), making July 2024 a clear standout.",
      "sql": "SELECT\n    toStartOfMonth(FlightDate) AS flight_month,\n    count() AS total_flights,\n    round(avg(DepDelay), 2) AS avg_dep_delay_min,\n    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE Reporting_Airline = 'AA'\n  AND Cancelled = 0\n  AND DepDelay IS NOT NULL\nGROUP BY flight_month\nORDER BY flight_month",
      "row_count": 458,
      "result_columns": [
        "flight_month",
        "total_flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "avg_dep_delay_min": 4.27,
        "flight_month": "1987-10-01T00:00:00Z",
        "pct_dep_del15": 8.8,
        "total_flights": 55871
      }
    },
    {
      "id": "origin_contributors",
      "subquestion": "Which origins contribute most to that peak month?",
      "answer_markdown": "**Charlotte (CLT)** is by far the dominant origin contributor in July 2024, with 10,568 departures — roughly 12% of all AA July 2024 flights — averaging **52.0 minutes** of departure delay and 52.2% of flights delayed 15+ minutes. No other origin comes close in combined volume and severity. Among smaller but notably delayed origins: IAH (632 flights, 47.1 min avg), EWR (615 flights, 44.9 min avg), RSW (354 flights, 45.3 min avg), and CVG (312 flights, 47.4 min avg) also stand out. CLT's massive scale makes it the single largest driver of AA's July 2024 network-wide delay figure.",
      "sql": "SELECT\n    OriginCode,\n    count() AS flights,\n    round(avg(DepDelay), 2) AS avg_dep_delay_min,\n    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE Reporting_Airline = 'AA'\n  AND Cancelled = 0\n  AND DepDelay IS NOT NULL\n  AND toStartOfMonth(FlightDate) = '2024-07-01'\nGROUP BY OriginCode\nHAVING flights \u003e= 100\nORDER BY avg_dep_delay_min DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "OriginCode",
        "flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "OriginCode": "EYW",
        "avg_dep_delay_min": 70,
        "flights": 119,
        "pct_dep_del15": 45.38
      }
    },
    {
      "id": "route_contributors",
      "subquestion": "Which routes contribute most to that peak month?",
      "answer_markdown": "The worst routes out of CLT led the July 2024 ranking: **CLT→SDF** averaged 102.7 minutes of departure delay (69.9% del15, 83 flights), **CLT→ROC** averaged 91.9 minutes (69.1% del15, 55 flights), and **CLT→TYS** averaged 90.0 minutes (69.0% del15, 58 flights). Other severe routes include **DFW→CID** (95.1 min, 60 flights), **JAC→DFW** (93.9 min, 63 flights), and **DFW→ICT** (93.4 min, 59 flights). CLT-origin routes dominate the worst performers, reinforcing its role as the network's delay epicenter that month. The PHX→JAX route had the single highest average at 134.4 minutes but on only 31 flights.",
      "sql": "SELECT\n    OriginCode,\n    DestCode,\n    count() AS flights,\n    round(avg(DepDelay), 2) AS avg_dep_delay_min,\n    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE Reporting_Airline = 'AA'\n  AND Cancelled = 0\n  AND DepDelay IS NOT NULL\n  AND toStartOfMonth(FlightDate) = '2024-07-01'\nGROUP BY OriginCode, DestCode\nHAVING flights \u003e= 30\nORDER BY avg_dep_delay_min DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "OriginCode",
        "DestCode",
        "flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "DestCode": "JAX",
        "OriginCode": "PHX",
        "avg_dep_delay_min": 134.42,
        "flights": 31,
        "pct_dep_del15": 45.16
      }
    },
    {
      "id": "concentration_pattern",
      "subquestion": "Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?",
      "answer_markdown": "The July 2024 peak is **concentrated, not broad**. CLT alone accounted for ~10,568 of 86,083 AA departures (~12.3%) and averaged over 52 minutes of delay — well above the 36.3-minute network average — meaning CLT's underperformance pulled the entire network number up significantly. The top 5 origins by average delay (EYW, JAC, ALB, PNS, ECP) were all small stations with under 250 flights each. At the route level, the worst-performing routes are CLT-hub spokes and a handful of DFW feeder routes, not a broad spread across the full network. The concentration around CLT as AA's largest East Coast hub — where cascading delays compound across hundreds of daily departures — explains why July 2024 stands so far above other months in the historical record.",
      "sql": "SELECT\n    OriginCode,\n    count() AS flights,\n    round(100.0 * count() / sum(count()) OVER (), 2) AS pct_of_total_flights,\n    round(avg(DepDelay), 2) AS avg_dep_delay_min,\n    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE Reporting_Airline = 'AA'\n  AND Cancelled = 0\n  AND DepDelay IS NOT NULL\n  AND toStartOfMonth(FlightDate) = '2024-07-01'\nGROUP BY OriginCode\nORDER BY flights DESC\nLIMIT 20",
      "row_count": 20,
      "result_columns": [
        "OriginCode",
        "flights",
        "pct_of_total_flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "OriginCode": "DFW",
        "avg_dep_delay_min": 38.52,
        "flights": 14962,
        "pct_dep_del15": 46.59,
        "pct_of_total_flights": 17.38
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.