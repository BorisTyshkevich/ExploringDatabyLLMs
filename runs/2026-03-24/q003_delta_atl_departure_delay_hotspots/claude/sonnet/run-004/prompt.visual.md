- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create browser-ready HTML `visual.html` using the proper  `*-analyst-dashboard` skill.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.

### Rules

- Question title: `Delta ATL departure delay hotspots by destination and time block`
- Visual mode: `dynamic`
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
    h.DestCode,
    h.DepTimeBlk,
    h.total_flights,
    h.avg_dep_delay,
    h.p90_dep_delay,
    h.pct_del15,
    m.credible_months
FROM (
    SELECT
        DestCode,
        DepTimeBlk,
        count() AS total_flights,
        round(avg(DepDelay), 1) AS avg_dep_delay,
        round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,
        round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15
    FROM ontime.fact_ontime
    WHERE
        OriginCode = 'ATL'
        AND IATA_CODE_Reporting_Airline = 'DL'
        AND Cancelled = 0
        AND DepDelay IS NOT NULL
    GROUP BY DestCode, DepTimeBlk
    HAVING total_flights >= 500
) h
JOIN (
    SELECT
        DestCode,
        DepTimeBlk,
        countIf(monthly_cnt >= 20) AS credible_months
    FROM (
        SELECT
            DestCode,
            DepTimeBlk,
            toStartOfMonth(FlightDate) AS ym,
            count() AS monthly_cnt
        FROM ontime.fact_ontime
        WHERE
            OriginCode = 'ATL'
            AND IATA_CODE_Reporting_Airline = 'DL'
            AND Cancelled = 0
            AND DepDelay IS NOT NULL
        GROUP BY DestCode, DepTimeBlk, ym
    )
    GROUP BY DestCode, DepTimeBlk
) m ON h.DestCode = m.DestCode AND h.DepTimeBlk = m.DepTimeBlk
ORDER BY h.avg_dep_delay DESC
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
      "id": "worst_hotspot",
      "subquestion": "Which destination and time block is the worst recurring hotspot?",
      "answer_markdown": "The worst recurring hotspot is **EWR (Newark) departing in the 1600–1659 block**, with an average departure delay of 27.1 minutes, a p90 delay of 82.9 minutes, and 39.3% of flights departing 15+ minutes late across 6,472 completed flights and 228 credible months of data. EWR also dominates the next two slots (1700–1759 at 25.1 min avg, 1800–1859 at 25.0 min avg), making the late-afternoon Newark leg the single most pressure-laden departure pattern Delta operates out of ATL.",
      "sql": "SELECT\n    h.DestCode,\n    h.DepTimeBlk,\n    h.total_flights,\n    h.avg_dep_delay,\n    h.p90_dep_delay,\n    h.pct_del15,\n    m.credible_months\nFROM (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        count() AS total_flights,\n        round(avg(DepDelay), 1) AS avg_dep_delay,\n        round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,\n        round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15\n    FROM ontime.fact_ontime\n    WHERE\n        OriginCode = 'ATL'\n        AND IATA_CODE_Reporting_Airline = 'DL'\n        AND Cancelled = 0\n        AND DepDelay IS NOT NULL\n    GROUP BY DestCode, DepTimeBlk\n    HAVING total_flights \u003e= 500\n) h\nJOIN (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        countIf(monthly_cnt \u003e= 20) AS credible_months\n    FROM (\n        SELECT\n            DestCode,\n            DepTimeBlk,\n            toStartOfMonth(FlightDate) AS ym,\n            count() AS monthly_cnt\n        FROM ontime.fact_ontime\n        WHERE\n            OriginCode = 'ATL'\n            AND IATA_CODE_Reporting_Airline = 'DL'\n            AND Cancelled = 0\n            AND DepDelay IS NOT NULL\n        GROUP BY DestCode, DepTimeBlk, ym\n    )\n    GROUP BY DestCode, DepTimeBlk\n) m ON h.DestCode = m.DestCode AND h.DepTimeBlk = m.DepTimeBlk\nORDER BY h.avg_dep_delay DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "total_flights",
        "avg_dep_delay",
        "p90_dep_delay",
        "pct_del15",
        "credible_months"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "avg_dep_delay": 27.1,
        "credible_months": 228,
        "p90_dep_delay": 82.9,
        "pct_del15": 39.3,
        "total_flights": 6472
      }
    },
    {
      "id": "persistence",
      "subquestion": "Is that hotspot consistently bad across time, or concentrated in a narrower period?",
      "answer_markdown": "The EWR 1600–1659 hotspot is **consistently bad across time**, not a narrow incident window. It accumulates 228 credible months (months with at least 20 departures) spanning from 1991 through the mid-2020s. Monthly data shows average delays persistently in the 15–50 minute range across decades, with peak months exceeding 45 minutes average delay and del15 rates above 60–70%. The breadth of credible months — more than 19 years of qualifying monthly presence — confirms this is a structural, recurring pressure point rather than a concentrated bad period.",
      "sql": "SELECT\n    DestCode,\n    DepTimeBlk,\n    toStartOfMonth(FlightDate) AS flight_month,\n    count() AS flights,\n    round(avg(DepDelay), 1) AS avg_dep_delay,\n    round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,\n    round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15\nFROM ontime.fact_ontime\nWHERE\n    OriginCode = 'ATL'\n    AND IATA_CODE_Reporting_Airline = 'DL'\n    AND Cancelled = 0\n    AND DepDelay IS NOT NULL\n    AND DestCode = 'EWR'\n    AND DepTimeBlk = '1600-1659'\nGROUP BY DestCode, DepTimeBlk, flight_month\nHAVING flights \u003e= 20\nORDER BY flight_month",
      "row_count": 228,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "flight_month",
        "flights",
        "avg_dep_delay",
        "p90_dep_delay",
        "pct_del15"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "avg_dep_delay": 18.8,
        "flight_month": "1991-04-01T00:00:00Z",
        "flights": 30,
        "p90_dep_delay": 39.4,
        "pct_del15": 36.7
      }
    },
    {
      "id": "pattern_summary",
      "subquestion": "What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?",
      "answer_markdown": "The top hotspots point overwhelmingly to **New York-area airports — EWR and JFK — in the afternoon and evening departure windows (1500–2059)**. EWR occupies four of the top ten slots by average delay (1500–1659, 1600–1659, 1700–1759, 1800–1859, 1900–1959), all with average delays above 22 minutes and del15 rates above 34%. JFK follows with multiple late-evening entries (1700–1759 at 23.9 min avg, 1900–1959 at 24.7 min avg). The pattern is clear: Delta's ATL-to-New York corridor concentrates departure pressure in the late afternoon, when cumulative schedule slippage from earlier wave departures compounds with congestion at the New York metro airports. The EWR 1600–1859 window is the single densest cluster of persistent hotspots, and any operational improvement there would have the highest impact on Delta's ATL departure performance.",
      "sql": "SELECT\n    h.DestCode,\n    h.DepTimeBlk,\n    h.total_flights,\n    h.avg_dep_delay,\n    h.p90_dep_delay,\n    h.pct_del15,\n    m.credible_months\nFROM (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        count() AS total_flights,\n        round(avg(DepDelay), 1) AS avg_dep_delay,\n        round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,\n        round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15\n    FROM ontime.fact_ontime\n    WHERE\n        OriginCode = 'ATL'\n        AND IATA_CODE_Reporting_Airline = 'DL'\n        AND Cancelled = 0\n        AND DepDelay IS NOT NULL\n    GROUP BY DestCode, DepTimeBlk\n    HAVING total_flights \u003e= 500\n) h\nJOIN (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        countIf(monthly_cnt \u003e= 20) AS credible_months\n    FROM (\n        SELECT\n            DestCode,\n            DepTimeBlk,\n            toStartOfMonth(FlightDate) AS ym,\n            count() AS monthly_cnt\n        FROM ontime.fact_ontime\n        WHERE\n            OriginCode = 'ATL'\n            AND IATA_CODE_Reporting_Airline = 'DL'\n            AND Cancelled = 0\n            AND DepDelay IS NOT NULL\n        GROUP BY DestCode, DepTimeBlk, ym\n    )\n    GROUP BY DestCode, DepTimeBlk\n) m ON h.DestCode = m.DestCode AND h.DepTimeBlk = m.DepTimeBlk\nWHERE m.credible_months \u003e= 24\nORDER BY h.avg_dep_delay DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "total_flights",
        "avg_dep_delay",
        "p90_dep_delay",
        "pct_del15",
        "credible_months"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "avg_dep_delay": 27.1,
        "credible_months": 228,
        "p90_dep_delay": 82.9,
        "pct_del15": 39.3,
        "total_flights": 6472
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
      "id": "worst_hotspot",
      "subquestion": "Which destination and time block is the worst recurring hotspot?",
      "answer_markdown": "The worst recurring hotspot is **EWR (Newark) departing in the 1600–1659 block**, with an average departure delay of 27.1 minutes, a p90 delay of 82.9 minutes, and 39.3% of flights departing 15+ minutes late across 6,472 completed flights and 228 credible months of data. EWR also dominates the next two slots (1700–1759 at 25.1 min avg, 1800–1859 at 25.0 min avg), making the late-afternoon Newark leg the single most pressure-laden departure pattern Delta operates out of ATL.",
      "sql": "SELECT\n    h.DestCode,\n    h.DepTimeBlk,\n    h.total_flights,\n    h.avg_dep_delay,\n    h.p90_dep_delay,\n    h.pct_del15,\n    m.credible_months\nFROM (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        count() AS total_flights,\n        round(avg(DepDelay), 1) AS avg_dep_delay,\n        round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,\n        round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15\n    FROM ontime.fact_ontime\n    WHERE\n        OriginCode = 'ATL'\n        AND IATA_CODE_Reporting_Airline = 'DL'\n        AND Cancelled = 0\n        AND DepDelay IS NOT NULL\n    GROUP BY DestCode, DepTimeBlk\n    HAVING total_flights \u003e= 500\n) h\nJOIN (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        countIf(monthly_cnt \u003e= 20) AS credible_months\n    FROM (\n        SELECT\n            DestCode,\n            DepTimeBlk,\n            toStartOfMonth(FlightDate) AS ym,\n            count() AS monthly_cnt\n        FROM ontime.fact_ontime\n        WHERE\n            OriginCode = 'ATL'\n            AND IATA_CODE_Reporting_Airline = 'DL'\n            AND Cancelled = 0\n            AND DepDelay IS NOT NULL\n        GROUP BY DestCode, DepTimeBlk, ym\n    )\n    GROUP BY DestCode, DepTimeBlk\n) m ON h.DestCode = m.DestCode AND h.DepTimeBlk = m.DepTimeBlk\nORDER BY h.avg_dep_delay DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "total_flights",
        "avg_dep_delay",
        "p90_dep_delay",
        "pct_del15",
        "credible_months"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "avg_dep_delay": 27.1,
        "credible_months": 228,
        "p90_dep_delay": 82.9,
        "pct_del15": 39.3,
        "total_flights": 6472
      }
    },
    {
      "id": "persistence",
      "subquestion": "Is that hotspot consistently bad across time, or concentrated in a narrower period?",
      "answer_markdown": "The EWR 1600–1659 hotspot is **consistently bad across time**, not a narrow incident window. It accumulates 228 credible months (months with at least 20 departures) spanning from 1991 through the mid-2020s. Monthly data shows average delays persistently in the 15–50 minute range across decades, with peak months exceeding 45 minutes average delay and del15 rates above 60–70%. The breadth of credible months — more than 19 years of qualifying monthly presence — confirms this is a structural, recurring pressure point rather than a concentrated bad period.",
      "sql": "SELECT\n    DestCode,\n    DepTimeBlk,\n    toStartOfMonth(FlightDate) AS flight_month,\n    count() AS flights,\n    round(avg(DepDelay), 1) AS avg_dep_delay,\n    round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,\n    round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15\nFROM ontime.fact_ontime\nWHERE\n    OriginCode = 'ATL'\n    AND IATA_CODE_Reporting_Airline = 'DL'\n    AND Cancelled = 0\n    AND DepDelay IS NOT NULL\n    AND DestCode = 'EWR'\n    AND DepTimeBlk = '1600-1659'\nGROUP BY DestCode, DepTimeBlk, flight_month\nHAVING flights \u003e= 20\nORDER BY flight_month",
      "row_count": 228,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "flight_month",
        "flights",
        "avg_dep_delay",
        "p90_dep_delay",
        "pct_del15"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "avg_dep_delay": 18.8,
        "flight_month": "1991-04-01T00:00:00Z",
        "flights": 30,
        "p90_dep_delay": 39.4,
        "pct_del15": 36.7
      }
    },
    {
      "id": "pattern_summary",
      "subquestion": "What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?",
      "answer_markdown": "The top hotspots point overwhelmingly to **New York-area airports — EWR and JFK — in the afternoon and evening departure windows (1500–2059)**. EWR occupies four of the top ten slots by average delay (1500–1659, 1600–1659, 1700–1759, 1800–1859, 1900–1959), all with average delays above 22 minutes and del15 rates above 34%. JFK follows with multiple late-evening entries (1700–1759 at 23.9 min avg, 1900–1959 at 24.7 min avg). The pattern is clear: Delta's ATL-to-New York corridor concentrates departure pressure in the late afternoon, when cumulative schedule slippage from earlier wave departures compounds with congestion at the New York metro airports. The EWR 1600–1859 window is the single densest cluster of persistent hotspots, and any operational improvement there would have the highest impact on Delta's ATL departure performance.",
      "sql": "SELECT\n    h.DestCode,\n    h.DepTimeBlk,\n    h.total_flights,\n    h.avg_dep_delay,\n    h.p90_dep_delay,\n    h.pct_del15,\n    m.credible_months\nFROM (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        count() AS total_flights,\n        round(avg(DepDelay), 1) AS avg_dep_delay,\n        round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,\n        round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15\n    FROM ontime.fact_ontime\n    WHERE\n        OriginCode = 'ATL'\n        AND IATA_CODE_Reporting_Airline = 'DL'\n        AND Cancelled = 0\n        AND DepDelay IS NOT NULL\n    GROUP BY DestCode, DepTimeBlk\n    HAVING total_flights \u003e= 500\n) h\nJOIN (\n    SELECT\n        DestCode,\n        DepTimeBlk,\n        countIf(monthly_cnt \u003e= 20) AS credible_months\n    FROM (\n        SELECT\n            DestCode,\n            DepTimeBlk,\n            toStartOfMonth(FlightDate) AS ym,\n            count() AS monthly_cnt\n        FROM ontime.fact_ontime\n        WHERE\n            OriginCode = 'ATL'\n            AND IATA_CODE_Reporting_Airline = 'DL'\n            AND Cancelled = 0\n            AND DepDelay IS NOT NULL\n        GROUP BY DestCode, DepTimeBlk, ym\n    )\n    GROUP BY DestCode, DepTimeBlk\n) m ON h.DestCode = m.DestCode AND h.DepTimeBlk = m.DepTimeBlk\nWHERE m.credible_months \u003e= 24\nORDER BY h.avg_dep_delay DESC\nLIMIT 15",
      "row_count": 15,
      "result_columns": [
        "DestCode",
        "DepTimeBlk",
        "total_flights",
        "avg_dep_delay",
        "p90_dep_delay",
        "pct_del15",
        "credible_months"
      ],
      "first_row": {
        "DepTimeBlk": "1600-1659",
        "DestCode": "EWR",
        "avg_dep_delay": 27.1,
        "credible_months": 228,
        "p90_dep_delay": 82.9,
        "pct_del15": 39.3,
        "total_flights": 6472
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.