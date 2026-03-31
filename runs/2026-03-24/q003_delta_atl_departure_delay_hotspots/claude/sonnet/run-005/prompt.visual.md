- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `Delta ATL departure delay hotspots by destination and time block`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_heatmap`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

Build a dynamic dashboard for Delta ATL departure delay hotspots.

Layout intent:

- headline and subtitle that clearly frame Delta departures from ATL
- narrative hero naming the single worst hotspot with its key metrics
- KPI strip derived from the primary fetched query
- hotspot ranking or heatmap from the primary query
- persistence chart or time-series panel using the `persistence` supporting query
- supporting pattern panel or table using the `pattern_summary` supporting query
- visible query ledger and footer controls following the dashboard skill contract
- concluding takeaway section that summarizes the broader operational pattern across the leading hotspots

Behavior:

- use `worst_hotspot` as the primary saved SQL already provided in the prompt
- use `persistence` and `pattern_summary` as supporting queries when they materially improve the dashboard
- if supporting queries are used, show them in the query ledger with their own status and SQL text
- the dashboard does not need to mirror `report.md`; use the narrative answers as framing and the fetched query results for charts and tables
- make the single worst hotspot visually prominent
- keep the page readable on mobile with stacked sections and horizontally scrollable tables or charts when needed
- preserve useful content when any supporting query fails or returns zero rows

### Data Source

SQL query for primary data source:

```sql
SELECT
    DestCode,
    DepTimeBlk,
    count() AS flight_count,
    round(avg(DepDelay), 1) AS avg_dep_delay_min,
    round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay_min,
    round(avg(DepDel15) * 100, 1) AS pct_dep_del15,
    countDistinct(toStartOfMonth(FlightDate)) AS active_months
FROM ontime.fact_ontime
WHERE
    IATA_CODE_Reporting_Airline = 'DL'
    AND OriginCode = 'ATL'
    AND Cancelled = 0
    AND DepDelay IS NOT NULL
GROUP BY DestCode, DepTimeBlk
HAVING flight_count >= 500
ORDER BY avg_dep_delay_min DESC
LIMIT 15
```

Data example/snippet:

{
  "question_title": "Delta ATL departure delay hotspots by destination and time block",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which destination and time block is the worst recurring hotspot?",
      "answer_markdown": "EWR in the 1600–1659 block is the worst recurring hotspot out of ATL for Delta. Across 6,472 completed departures spanning 263 credible months, it averages 27.1 minutes of departure delay, with a 90th-percentile delay of 82.9 minutes and 39.3% of flights leaving 15 or more minutes late. The next-worst credible combinations are also EWR — in the 1700–1759 and 1800–1859 windows — confirming that the EWR afternoon-to-evening corridor is structurally the most problematic slot in Delta's ATL network.",
      "sql": "SELECT\n    DestCode,\n    DepTimeBlk,\n    count() AS flight_count,\n    round(avg(DepDelay), 1) AS avg_dep_delay_min,\n    round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay_min,\n    round(avg(DepDel15) * 100, 1) AS pct_dep_del15,\n    countDistinct(toStartOfMonth(FlightDate)) AS active_months\nFROM ontime.fact_ontime\nWHERE\n    IATA_CODE_Reporting_Airline = 'DL'\n    AND OriginCode = 'ATL'\n    AND Cancelled = 0\n    AND DepDelay IS NOT NULL\nGROUP BY DestCode, DepTimeBlk\nHAVING flight_count \u003e= 500\nORDER BY avg_dep_delay_min DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "flight_count",
        "avg_dep_delay_min",
        "p90_dep_delay_min",
        "pct_dep_del15",
        "active_months"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "active_months": 263,
        "avg_dep_delay_min": 27.1,
        "flight_count": 6472,
        "p90_dep_delay_min": 82.9,
        "pct_dep_del15": 39.3
      }
    },
    {
      "id": "q2",
      "subquestion": "Is that hotspot consistently bad across time, or concentrated in a narrower period?",
      "answer_markdown": "The EWR / 1600–1659 hotspot is persistently bad across a multi-decade span, not concentrated in a narrow window. Monthly data with at least 30 flights shows elevated delays recurring from 1991 through 2024, with severe spikes appearing in many different years: average monthly delays exceed 40 minutes in months spread across 1992, 2011, 2014, 2016, 2017, 2018, 2019, and 2024. There is no single concentrated period — the hotspot reappears whenever Delta operates meaningful volume on this slot, across more than three decades of data.",
      "sql": "SELECT\n    DestCode,\n    DepTimeBlk,\n    toStartOfMonth(FlightDate) AS flight_month,\n    count() AS monthly_flights,\n    round(avg(DepDelay), 1) AS avg_dep_delay_min,\n    round(avg(DepDel15) * 100, 1) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE\n    IATA_CODE_Reporting_Airline = 'DL'\n    AND OriginCode = 'ATL'\n    AND DestCode = 'EWR'\n    AND DepTimeBlk = '1600-1659'\n    AND Cancelled = 0\n    AND DepDelay IS NOT NULL\nGROUP BY DestCode, DepTimeBlk, flight_month\nHAVING monthly_flights \u003e= 30\nORDER BY flight_month",
      "row_count": 49,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "flight_month",
        "monthly_flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "avg_dep_delay_min": 18.8,
        "flight_month": "1991-04-01T00:00:00Z",
        "monthly_flights": 30,
        "pct_dep_del15": 36.7
      }
    },
    {
      "id": "q3",
      "subquestion": "What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?",
      "answer_markdown": "The top hotspots reveal that Delta's departure pressure out of ATL is overwhelmingly concentrated on New York metro routes. Of the 15 worst destination/time-block combinations (each with at least 500 flights), five are EWR in the 1500–1900 window and three more are JFK in the late-afternoon through evening window. Together, EWR and JFK account for the majority of the highest-average-delay slots. The pattern points to downstream congestion at New York-area airports compounding throughout the afternoon: delays build progressively from the 1500 block onward, reflecting both EWR/JFK slot pressure and the cumulative cascade of late-arriving aircraft. Outside the New York metro, OKC in the 2000–2059 block and IAH in the 2000–2059 block also appear, suggesting late-evening departure slots to mid-continent destinations carry elevated risk, likely driven by late-arriving equipment.",
      "sql": "SELECT\n    DestCode,\n    count() AS hotspot_slots,\n    sum(flight_count) AS total_flights,\n    round(avg(avg_dep_delay_min), 1) AS mean_avg_delay,\n    round(avg(pct_dep_del15), 1) AS mean_pct_del15\nFROM (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        count() AS flight_count,\n        round(avg(DepDelay), 1) AS avg_dep_delay_min,\n        round(avg(DepDel15) * 100, 1) AS pct_dep_del15\n    FROM ontime.fact_ontime\n    WHERE\n        IATA_CODE_Reporting_Airline = 'DL'\n        AND OriginCode = 'ATL'\n        AND Cancelled = 0\n        AND DepDelay IS NOT NULL\n    GROUP BY DestCode, DepTimeBlk\n    HAVING flight_count \u003e= 500\n    ORDER BY avg_dep_delay_min DESC\n    LIMIT 15\n)\nGROUP BY DestCode\nORDER BY mean_avg_delay DESC",
      "row_count": 9,
      "result_columns": [
        "DestCode",
        "hotspot_slots",
        "total_flights",
        "mean_avg_delay",
        "mean_pct_del15"
      ],
      "first_row": {
        "DestCode": "OKC",
        "hotspot_slots": 1,
        "mean_avg_delay": 24.7,
        "mean_pct_del15": 21.9,
        "total_flights": 4143
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
  "question_title": "Delta ATL departure delay hotspots by destination and time block",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which destination and time block is the worst recurring hotspot?",
      "answer_markdown": "EWR in the 1600–1659 block is the worst recurring hotspot out of ATL for Delta. Across 6,472 completed departures spanning 263 credible months, it averages 27.1 minutes of departure delay, with a 90th-percentile delay of 82.9 minutes and 39.3% of flights leaving 15 or more minutes late. The next-worst credible combinations are also EWR — in the 1700–1759 and 1800–1859 windows — confirming that the EWR afternoon-to-evening corridor is structurally the most problematic slot in Delta's ATL network.",
      "sql": "SELECT\n    DestCode,\n    DepTimeBlk,\n    count() AS flight_count,\n    round(avg(DepDelay), 1) AS avg_dep_delay_min,\n    round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay_min,\n    round(avg(DepDel15) * 100, 1) AS pct_dep_del15,\n    countDistinct(toStartOfMonth(FlightDate)) AS active_months\nFROM ontime.fact_ontime\nWHERE\n    IATA_CODE_Reporting_Airline = 'DL'\n    AND OriginCode = 'ATL'\n    AND Cancelled = 0\n    AND DepDelay IS NOT NULL\nGROUP BY DestCode, DepTimeBlk\nHAVING flight_count \u003e= 500\nORDER BY avg_dep_delay_min DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "flight_count",
        "avg_dep_delay_min",
        "p90_dep_delay_min",
        "pct_dep_del15",
        "active_months"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "active_months": 263,
        "avg_dep_delay_min": 27.1,
        "flight_count": 6472,
        "p90_dep_delay_min": 82.9,
        "pct_dep_del15": 39.3
      }
    },
    {
      "id": "q2",
      "subquestion": "Is that hotspot consistently bad across time, or concentrated in a narrower period?",
      "answer_markdown": "The EWR / 1600–1659 hotspot is persistently bad across a multi-decade span, not concentrated in a narrow window. Monthly data with at least 30 flights shows elevated delays recurring from 1991 through 2024, with severe spikes appearing in many different years: average monthly delays exceed 40 minutes in months spread across 1992, 2011, 2014, 2016, 2017, 2018, 2019, and 2024. There is no single concentrated period — the hotspot reappears whenever Delta operates meaningful volume on this slot, across more than three decades of data.",
      "sql": "SELECT\n    DestCode,\n    DepTimeBlk,\n    toStartOfMonth(FlightDate) AS flight_month,\n    count() AS monthly_flights,\n    round(avg(DepDelay), 1) AS avg_dep_delay_min,\n    round(avg(DepDel15) * 100, 1) AS pct_dep_del15\nFROM ontime.fact_ontime\nWHERE\n    IATA_CODE_Reporting_Airline = 'DL'\n    AND OriginCode = 'ATL'\n    AND DestCode = 'EWR'\n    AND DepTimeBlk = '1600-1659'\n    AND Cancelled = 0\n    AND DepDelay IS NOT NULL\nGROUP BY DestCode, DepTimeBlk, flight_month\nHAVING monthly_flights \u003e= 30\nORDER BY flight_month",
      "row_count": 49,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "flight_month",
        "monthly_flights",
        "avg_dep_delay_min",
        "pct_dep_del15"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "avg_dep_delay_min": 18.8,
        "flight_month": "1991-04-01T00:00:00Z",
        "monthly_flights": 30,
        "pct_dep_del15": 36.7
      }
    },
    {
      "id": "q3",
      "subquestion": "What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?",
      "answer_markdown": "The top hotspots reveal that Delta's departure pressure out of ATL is overwhelmingly concentrated on New York metro routes. Of the 15 worst destination/time-block combinations (each with at least 500 flights), five are EWR in the 1500–1900 window and three more are JFK in the late-afternoon through evening window. Together, EWR and JFK account for the majority of the highest-average-delay slots. The pattern points to downstream congestion at New York-area airports compounding throughout the afternoon: delays build progressively from the 1500 block onward, reflecting both EWR/JFK slot pressure and the cumulative cascade of late-arriving aircraft. Outside the New York metro, OKC in the 2000–2059 block and IAH in the 2000–2059 block also appear, suggesting late-evening departure slots to mid-continent destinations carry elevated risk, likely driven by late-arriving equipment.",
      "sql": "SELECT\n    DestCode,\n    count() AS hotspot_slots,\n    sum(flight_count) AS total_flights,\n    round(avg(avg_dep_delay_min), 1) AS mean_avg_delay,\n    round(avg(pct_dep_del15), 1) AS mean_pct_del15\nFROM (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        count() AS flight_count,\n        round(avg(DepDelay), 1) AS avg_dep_delay_min,\n        round(avg(DepDel15) * 100, 1) AS pct_dep_del15\n    FROM ontime.fact_ontime\n    WHERE\n        IATA_CODE_Reporting_Airline = 'DL'\n        AND OriginCode = 'ATL'\n        AND Cancelled = 0\n        AND DepDelay IS NOT NULL\n    GROUP BY DestCode, DepTimeBlk\n    HAVING flight_count \u003e= 500\n    ORDER BY avg_dep_delay_min DESC\n    LIMIT 15\n)\nGROUP BY DestCode\nORDER BY mean_avg_delay DESC",
      "row_count": 9,
      "result_columns": [
        "DestCode",
        "hotspot_slots",
        "total_flights",
        "mean_avg_delay",
        "mean_pct_del15"
      ],
      "first_row": {
        "DestCode": "OKC",
        "hotspot_slots": 1,
        "mean_avg_delay": 24.7,
        "mean_pct_del15": 21.9,
        "total_flights": 4143
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