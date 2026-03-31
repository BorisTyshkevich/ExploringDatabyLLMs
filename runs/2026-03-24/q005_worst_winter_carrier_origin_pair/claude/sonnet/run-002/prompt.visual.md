- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `Worst winter carrier-origin pairs by departure performance`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_seasonal_dashboard`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

Build a dashboard that:

- uses `worst_pair` as the primary saved SQL already provided in the prompt
- uses `cause_mix` and `concentration_pattern` as supporting queries when they materially improve the dashboard
- shows KPI cards for worst winter pair, worst OTP, average delay of the worst pair, and total qualifying pairs
- renders a ranked chart for the worst winter `(Reporting_Airline, OriginCode)` pairs
- renders a stacked bar chart of delay-cause shares for the top 10 pairs
- renders a compact table of the full ranked result
- includes a narrative takeaway about whether the weakest winter pairs cluster in a small set of carriers or origins
- derives the top 10 pairs for the cause-share chart from fetched ranking data, not from hardcoded labels
- makes winter framing explicit in titles and copy
- clearly separates ranking severity from cause composition
- keeps carrier-airport labels readable without truncating meaning
- shows supporting queries in the query ledger when used

### Data Source

SQL query for primary data source:

```sql
SELECT
    IATA_CODE_Reporting_Airline AS carrier,
    OriginCode AS origin,
    count() AS winter_flights,
    round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct,
    round(avg(ifNull(DepDelayMinutes, 0)), 2) AS avg_dep_delay_min
FROM ontime.fact_ontime
WHERE toMonth(FlightDate) IN (12, 1, 2)
  AND Cancelled = 0
GROUP BY carrier, origin
HAVING winter_flights >= 1000
ORDER BY otp_pct ASC
LIMIT 20
```

Data example/snippet:

{
  "question_title": "Worst winter carrier-origin pairs by departure performance",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which winter carrier-airport pair ranks worst overall?",
      "answer_markdown": "DH (Independence Air) departing from ORD (Chicago O'Hare) ranks worst among all qualifying winter carrier-origin pairs with at least 1,000 winter flights. Across 19,986 winter departures, it achieved only 56.58% on-time performance and averaged 28.35 minutes of departure delay per completed flight — both the worst figures in the dataset by a meaningful margin. The next-worst qualifying pairs (PI/DFW at 61.5%, PI/LAX at 62.95%) trail DH/ORD by more than four percentage points.",
      "sql": "SELECT\n    IATA_CODE_Reporting_Airline AS carrier,\n    OriginCode AS origin,\n    count() AS winter_flights,\n    round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct,\n    round(avg(ifNull(DepDelayMinutes, 0)), 2) AS avg_dep_delay_min\nFROM ontime.fact_ontime\nWHERE toMonth(FlightDate) IN (12, 1, 2)\n  AND Cancelled = 0\nGROUP BY carrier, origin\nHAVING winter_flights \u003e= 1000\nORDER BY otp_pct ASC\nLIMIT 20",
      "row_count": 20,
      "result_columns": [
        "carrier",
        "origin",
        "winter_flights",
        "otp_pct",
        "avg_dep_delay_min"
      ],
      "first_row": {
        "avg_dep_delay_min": 28.35,
        "carrier": "DH",
        "origin": "ORD",
        "otp_pct": 56.58,
        "winter_flights": 19986
      }
    },
    {
      "id": "q2",
      "subquestion": "Are the worst pairs driven more by weather or by operational causes?",
      "answer_markdown": "Operational causes dominate decisively across all worst pairs that have delay-cause data. For DH/ORD the breakdown is 15.6% weather vs 84.4% operational (carrier delay + NAS + late aircraft + security). YV/ORD is even more skewed: 5.0% weather vs 95.0% operational. AA/EYW sits at 5.7% weather vs 94.3% operational; F9/PBI at 1.5% vs 98.5%; OO/OTH at 2.9% vs 97.1%. OH/RSW is the most weather-affected at 26.8% weather vs 73.2% operational. In every case, operational causes account for the majority of reported delay minutes, making poor operational execution rather than weather the primary driver of chronic winter under-performance.",
      "sql": "SELECT\n    IATA_CODE_Reporting_Airline AS carrier,\n    OriginCode AS origin,\n    count() AS winter_flights,\n    round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct,\n    round(avg(ifNull(DepDelayMinutes, 0)), 2) AS avg_dep_delay_min,\n    sum(ifNull(WeatherDelay, 0)) AS weather_min,\n    sum(ifNull(CarrierDelay, 0)) AS carrier_min,\n    sum(ifNull(NASDelay, 0)) AS nas_min,\n    sum(ifNull(SecurityDelay, 0)) AS security_min,\n    sum(ifNull(LateAircraftDelay, 0)) AS late_aircraft_min,\n    sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)\n        + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)) AS total_reported_delay,\n    round(sum(ifNull(WeatherDelay, 0)) * 100.0\n          / nullIf(sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)\n                       + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)), 0), 1) AS weather_share_pct,\n    round(sum(ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0) + ifNull(SecurityDelay, 0)\n              + ifNull(LateAircraftDelay, 0)) * 100.0\n          / nullIf(sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)\n                       + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)), 0), 1) AS operational_share_pct\nFROM ontime.fact_ontime\nWHERE toMonth(FlightDate) IN (12, 1, 2)\n  AND Cancelled = 0\n  AND (IATA_CODE_Reporting_Airline, OriginCode) IN (\n      ('DH','ORD'), ('PI','DFW'), ('PI','LAX'), ('PI','DAY'),\n      ('AS','DUT'), ('YV','ORD'), ('OO','OTH'), ('AA','EYW'),\n      ('F9','PBI'), ('OH','RSW')\n  )\nGROUP BY carrier, origin\nORDER BY otp_pct ASC",
      "row_count": 10,
      "result_columns": [
        "carrier",
        "origin",
        "winter_flights",
        "otp_pct",
        "avg_dep_delay_min",
        "weather_min",
        "carrier_min",
        "nas_min",
        "security_min",
        "late_aircraft_min",
        "total_reported_delay",
        "weather_share_pct",
        "operational_share_pct"
      ],
      "first_row": {
        "avg_dep_delay_min": 28.35,
        "carrier": "DH",
        "carrier_min": 123407,
        "late_aircraft_min": 133239,
        "nas_min": 41492,
        "operational_share_pct": 84.4,
        "origin": "ORD",
        "otp_pct": 56.58,
        "security_min": 116,
        "total_reported_delay": 353430,
        "weather_min": 55176,
        "weather_share_pct": 15.6,
        "winter_flights": 19986
      }
    },
    {
      "id": "q3",
      "subquestion": "Are the weakest pairs concentrated in a small number of carriers or airports?",
      "answer_markdown": "The weakest pairs are moderately spread rather than tightly concentrated. Among the 30 worst qualifying winter pairs, B6 (JetBlue) is the most frequent carrier with 6 pairs spread across MIA, RNO, SRQ, FLL, PHX, and PBI. PI, OO, and OH each contribute 3 pairs. YV, F9, EV, and FL each have 2 pairs, while DH, AS, AA, NW, G4, XE, and EA appear once each — 15 distinct carriers in total. On the airport side, ORD and EWR each host 2 pairs, FLL appears twice (B6 and F9), and no single airport dominates the list. The weakness is carrier-driven more than airport-driven: B6's Florida leisure-market airports and regional carriers at high-congestion or weather-prone hubs account for most of the bottom tier, but there is no single choke-point airport that concentrates failure.",
      "sql": "WITH worst30 AS (\n    SELECT\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode AS origin,\n        count() AS winter_flights,\n        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct\n    FROM ontime.fact_ontime\n    WHERE toMonth(FlightDate) IN (12, 1, 2)\n      AND Cancelled = 0\n    GROUP BY carrier, origin\n    HAVING winter_flights \u003e= 1000\n    ORDER BY otp_pct ASC\n    LIMIT 30\n)\nSELECT\n    carrier,\n    count() AS pairs_in_worst30,\n    round(avg(otp_pct), 2) AS avg_otp_pct,\n    groupArray(origin) AS origins\nFROM worst30\nGROUP BY carrier\nORDER BY pairs_in_worst30 DESC, avg_otp_pct ASC",
      "row_count": 15,
      "result_columns": [
        "carrier",
        "pairs_in_worst30",
        "avg_otp_pct",
        "origins"
      ],
      "first_row": {
        "avg_otp_pct": 68.43,
        "carrier": "B6",
        "origins": [
          "MIA",
          "RNO",
          "SRQ",
          "FLL",
          "PHX",
          "PBI"
        ],
        "pairs_in_worst30": 6
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
  "question_title": "Worst winter carrier-origin pairs by departure performance",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which winter carrier-airport pair ranks worst overall?",
      "answer_markdown": "DH (Independence Air) departing from ORD (Chicago O'Hare) ranks worst among all qualifying winter carrier-origin pairs with at least 1,000 winter flights. Across 19,986 winter departures, it achieved only 56.58% on-time performance and averaged 28.35 minutes of departure delay per completed flight — both the worst figures in the dataset by a meaningful margin. The next-worst qualifying pairs (PI/DFW at 61.5%, PI/LAX at 62.95%) trail DH/ORD by more than four percentage points.",
      "sql": "SELECT\n    IATA_CODE_Reporting_Airline AS carrier,\n    OriginCode AS origin,\n    count() AS winter_flights,\n    round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct,\n    round(avg(ifNull(DepDelayMinutes, 0)), 2) AS avg_dep_delay_min\nFROM ontime.fact_ontime\nWHERE toMonth(FlightDate) IN (12, 1, 2)\n  AND Cancelled = 0\nGROUP BY carrier, origin\nHAVING winter_flights \u003e= 1000\nORDER BY otp_pct ASC\nLIMIT 20",
      "row_count": 20,
      "result_columns": [
        "carrier",
        "origin",
        "winter_flights",
        "otp_pct",
        "avg_dep_delay_min"
      ],
      "first_row": {
        "avg_dep_delay_min": 28.35,
        "carrier": "DH",
        "origin": "ORD",
        "otp_pct": 56.58,
        "winter_flights": 19986
      }
    },
    {
      "id": "q2",
      "subquestion": "Are the worst pairs driven more by weather or by operational causes?",
      "answer_markdown": "Operational causes dominate decisively across all worst pairs that have delay-cause data. For DH/ORD the breakdown is 15.6% weather vs 84.4% operational (carrier delay + NAS + late aircraft + security). YV/ORD is even more skewed: 5.0% weather vs 95.0% operational. AA/EYW sits at 5.7% weather vs 94.3% operational; F9/PBI at 1.5% vs 98.5%; OO/OTH at 2.9% vs 97.1%. OH/RSW is the most weather-affected at 26.8% weather vs 73.2% operational. In every case, operational causes account for the majority of reported delay minutes, making poor operational execution rather than weather the primary driver of chronic winter under-performance.",
      "sql": "SELECT\n    IATA_CODE_Reporting_Airline AS carrier,\n    OriginCode AS origin,\n    count() AS winter_flights,\n    round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct,\n    round(avg(ifNull(DepDelayMinutes, 0)), 2) AS avg_dep_delay_min,\n    sum(ifNull(WeatherDelay, 0)) AS weather_min,\n    sum(ifNull(CarrierDelay, 0)) AS carrier_min,\n    sum(ifNull(NASDelay, 0)) AS nas_min,\n    sum(ifNull(SecurityDelay, 0)) AS security_min,\n    sum(ifNull(LateAircraftDelay, 0)) AS late_aircraft_min,\n    sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)\n        + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)) AS total_reported_delay,\n    round(sum(ifNull(WeatherDelay, 0)) * 100.0\n          / nullIf(sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)\n                       + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)), 0), 1) AS weather_share_pct,\n    round(sum(ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0) + ifNull(SecurityDelay, 0)\n              + ifNull(LateAircraftDelay, 0)) * 100.0\n          / nullIf(sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)\n                       + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)), 0), 1) AS operational_share_pct\nFROM ontime.fact_ontime\nWHERE toMonth(FlightDate) IN (12, 1, 2)\n  AND Cancelled = 0\n  AND (IATA_CODE_Reporting_Airline, OriginCode) IN (\n      ('DH','ORD'), ('PI','DFW'), ('PI','LAX'), ('PI','DAY'),\n      ('AS','DUT'), ('YV','ORD'), ('OO','OTH'), ('AA','EYW'),\n      ('F9','PBI'), ('OH','RSW')\n  )\nGROUP BY carrier, origin\nORDER BY otp_pct ASC",
      "row_count": 10,
      "result_columns": [
        "carrier",
        "origin",
        "winter_flights",
        "otp_pct",
        "avg_dep_delay_min",
        "weather_min",
        "carrier_min",
        "nas_min",
        "security_min",
        "late_aircraft_min",
        "total_reported_delay",
        "weather_share_pct",
        "operational_share_pct"
      ],
      "first_row": {
        "avg_dep_delay_min": 28.35,
        "carrier": "DH",
        "carrier_min": 123407,
        "late_aircraft_min": 133239,
        "nas_min": 41492,
        "operational_share_pct": 84.4,
        "origin": "ORD",
        "otp_pct": 56.58,
        "security_min": 116,
        "total_reported_delay": 353430,
        "weather_min": 55176,
        "weather_share_pct": 15.6,
        "winter_flights": 19986
      }
    },
    {
      "id": "q3",
      "subquestion": "Are the weakest pairs concentrated in a small number of carriers or airports?",
      "answer_markdown": "The weakest pairs are moderately spread rather than tightly concentrated. Among the 30 worst qualifying winter pairs, B6 (JetBlue) is the most frequent carrier with 6 pairs spread across MIA, RNO, SRQ, FLL, PHX, and PBI. PI, OO, and OH each contribute 3 pairs. YV, F9, EV, and FL each have 2 pairs, while DH, AS, AA, NW, G4, XE, and EA appear once each — 15 distinct carriers in total. On the airport side, ORD and EWR each host 2 pairs, FLL appears twice (B6 and F9), and no single airport dominates the list. The weakness is carrier-driven more than airport-driven: B6's Florida leisure-market airports and regional carriers at high-congestion or weather-prone hubs account for most of the bottom tier, but there is no single choke-point airport that concentrates failure.",
      "sql": "WITH worst30 AS (\n    SELECT\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode AS origin,\n        count() AS winter_flights,\n        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct\n    FROM ontime.fact_ontime\n    WHERE toMonth(FlightDate) IN (12, 1, 2)\n      AND Cancelled = 0\n    GROUP BY carrier, origin\n    HAVING winter_flights \u003e= 1000\n    ORDER BY otp_pct ASC\n    LIMIT 30\n)\nSELECT\n    carrier,\n    count() AS pairs_in_worst30,\n    round(avg(otp_pct), 2) AS avg_otp_pct,\n    groupArray(origin) AS origins\nFROM worst30\nGROUP BY carrier\nORDER BY pairs_in_worst30 DESC, avg_otp_pct ASC",
      "row_count": 15,
      "result_columns": [
        "carrier",
        "pairs_in_worst30",
        "avg_otp_pct",
        "origins"
      ],
      "first_row": {
        "avg_otp_pct": 68.43,
        "carrier": "B6",
        "origins": [
          "MIA",
          "RNO",
          "SRQ",
          "FLL",
          "PHX",
          "PBI"
        ],
        "pairs_in_worst30": 6
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