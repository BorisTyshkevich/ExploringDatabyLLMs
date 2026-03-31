- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create browser-ready HTML `visual.html` using the proper  `*-analyst-dashboard` skill.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.

### Rules

- Question title: `Worst winter carrier-origin pairs by departure performance`
- Visual mode: `dynamic`
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
WITH winter_pairs AS (
    SELECT
        IATA_CODE_Reporting_Airline AS carrier,
        OriginCode AS origin,
        count() AS winter_flights,
        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime,
        round(avg(DepDelayMinutes), 2) AS avg_dep_delay_min,
        round(sum(assumeNotNull(WeatherDelay)) / nullIf(sum(assumeNotNull(WeatherDelay) + assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)), 0) * 100, 2) AS weather_pct,
        round(sum(assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)) / nullIf(sum(assumeNotNull(WeatherDelay) + assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)), 0) * 100, 2) AS operational_pct
    FROM ontime.fact_ontime
    WHERE toMonth(FlightDate) IN (12, 1, 2)
      AND Cancelled = 0
    GROUP BY carrier, origin
    HAVING winter_flights >= 1000
)
SELECT
    carrier,
    origin,
    winter_flights,
    pct_ontime,
    avg_dep_delay_min,
    weather_pct,
    operational_pct,
    row_number() OVER (ORDER BY pct_ontime ASC) AS rank
FROM winter_pairs
ORDER BY pct_ontime ASC
LIMIT 25
```

Data example/snippet:

{
  "question_title": "Worst winter carrier-origin pairs by departure performance",
  "result_columns": null,
  "row_count": 3,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "worst_pair",
      "subquestion": "Which winter carrier-airport pair ranks worst overall?",
      "answer_markdown": "The worst winter carrier-airport pair is **DH (Independence Air) departing from ORD (Chicago O'Hare)**, with only 56.58% of winter flights departing on time and an average departure delay of 28.35 minutes across 19,986 qualifying winter flights. The second-worst pair is PI (Piedmont Airlines) at DFW (61.5% on-time), followed by PI at LAX (62.95%) and PI at DAY (63.15%). The ranked query below returns the full set of worst-performing winter pairs (minimum 1,000 winter flights) sorted by on-time percentage ascending.",
      "sql": "WITH winter_pairs AS (\n    SELECT\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode AS origin,\n        count() AS winter_flights,\n        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime,\n        round(avg(DepDelayMinutes), 2) AS avg_dep_delay_min,\n        round(sum(assumeNotNull(WeatherDelay)) / nullIf(sum(assumeNotNull(WeatherDelay) + assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)), 0) * 100, 2) AS weather_pct,\n        round(sum(assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)) / nullIf(sum(assumeNotNull(WeatherDelay) + assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)), 0) * 100, 2) AS operational_pct\n    FROM ontime.fact_ontime\n    WHERE toMonth(FlightDate) IN (12, 1, 2)\n      AND Cancelled = 0\n    GROUP BY carrier, origin\n    HAVING winter_flights \u003e= 1000\n)\nSELECT\n    carrier,\n    origin,\n    winter_flights,\n    pct_ontime,\n    avg_dep_delay_min,\n    weather_pct,\n    operational_pct,\n    row_number() OVER (ORDER BY pct_ontime ASC) AS rank\nFROM winter_pairs\nORDER BY pct_ontime ASC\nLIMIT 25",
      "row_count": 25,
      "result_columns": [
        "carrier",
        "origin",
        "winter_flights",
        "pct_ontime",
        "avg_dep_delay_min",
        "weather_pct",
        "operational_pct",
        "rank"
      ],
      "first_row": {
        "avg_dep_delay_min": 28.35,
        "carrier": "DH",
        "operational_pct": 84.39,
        "origin": "ORD",
        "pct_ontime": 56.58,
        "rank": 1,
        "weather_pct": 15.61,
        "winter_flights": 19986
      }
    },
    {
      "id": "cause_mix",
      "subquestion": "Are the worst pairs driven more by weather or by operational causes?",
      "answer_markdown": "The worst winter pairs are overwhelmingly driven by **operational causes**, not weather. Across the 22 worst pairs with delay-cause data, operational delays (carrier delays + late aircraft + NAS + security) account for 70–99% of reported delay minutes in nearly every case. Late aircraft propagation alone is the single largest component for most pairs — for example, DH/ORD (37.7% late aircraft), YV/ORD (45.2% late aircraft), AA/EYW (61.1% late aircraft), and AS/DUT (66.0% late aircraft). Carrier-attributed delays are the second-largest factor. Weather is a minor contributor for most pairs (typically 1–15%), with only two notable exceptions: OH/EWR (41.4% weather) and OH/ATL (28.7% weather). Three PI pairs show no delay-cause data, consistent with older data predating structured cause reporting.",
      "sql": "WITH worst_pairs AS (\n    SELECT\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode AS origin,\n        count() AS winter_flights,\n        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime,\n        round(avg(DepDelayMinutes), 2) AS avg_dep_delay_min,\n        round(sum(assumeNotNull(CarrierDelay)), 0) AS total_carrier_delay,\n        round(sum(assumeNotNull(WeatherDelay)), 0) AS total_weather_delay,\n        round(sum(assumeNotNull(NASDelay)), 0) AS total_nas_delay,\n        round(sum(assumeNotNull(SecurityDelay)), 0) AS total_security_delay,\n        round(sum(assumeNotNull(LateAircraftDelay)), 0) AS total_late_aircraft_delay\n    FROM ontime.fact_ontime\n    WHERE toMonth(FlightDate) IN (12, 1, 2)\n      AND Cancelled = 0\n    GROUP BY carrier, origin\n    HAVING winter_flights \u003e= 1000\n    ORDER BY pct_ontime ASC\n    LIMIT 25\n)\nSELECT\n    carrier,\n    origin,\n    pct_ontime,\n    total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay AS total_delay_min,\n    round(100.0 * total_weather_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS weather_share_pct,\n    round(100.0 * total_carrier_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS carrier_share_pct,\n    round(100.0 * total_nas_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS nas_share_pct,\n    round(100.0 * total_late_aircraft_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS late_aircraft_share_pct,\n    round(100.0 * total_security_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS security_share_pct\nFROM worst_pairs\nORDER BY pct_ontime ASC",
      "row_count": 25,
      "result_columns": [
        "carrier",
        "origin",
        "pct_ontime",
        "total_delay_min",
        "weather_share_pct",
        "carrier_share_pct",
        "nas_share_pct",
        "late_aircraft_share_pct",
        "security_share_pct"
      ],
      "first_row": {
        "carrier": "DH",
        "carrier_share_pct": 34.92,
        "late_aircraft_share_pct": 37.7,
        "nas_share_pct": 11.74,
        "origin": "ORD",
        "pct_ontime": 56.58,
        "security_share_pct": 0.03,
        "total_delay_min": 353430,
        "weather_share_pct": 15.61
      }
    },
    {
      "id": "concentration_pattern",
      "subquestion": "Are the weakest pairs concentrated in a small number of carriers or airports?",
      "answer_markdown": "The 25 worst winter pairs span **14 distinct carriers and 22 distinct airports**, so the weak set is fairly dispersed across airports but shows moderate carrier concentration. Four carriers account for the majority of repeating entries: **B6 (JetBlue)** appears 4 times (MIA, RNO, SRQ, FLL), **OH (Comair)** 3 times (RSW, EWR, ATL), **OO (SkyWest)** 3 times (OTH, ASE, CEC), and **PI (Piedmont)** 3 times (DFW, LAX, DAY). On the airport side, **EWR** appears in 3 pairs (OH, FL, XE) and **ORD** in 2 pairs (DH, YV). The pattern suggests that chronic winter departure underperformance clusters around a handful of regional and low-cost carriers operating systemically poor schedules, rather than being uniformly distributed across the industry.",
      "sql": "WITH worst_pairs AS (\n    SELECT\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode AS origin,\n        count() AS winter_flights,\n        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime\n    FROM ontime.fact_ontime\n    WHERE toMonth(FlightDate) IN (12, 1, 2)\n      AND Cancelled = 0\n    GROUP BY carrier, origin\n    HAVING winter_flights \u003e= 1000\n    ORDER BY pct_ontime ASC\n    LIMIT 25\n)\nSELECT\n    uniq(carrier) AS distinct_carriers,\n    uniq(origin) AS distinct_airports,\n    count() AS total_pairs,\n    groupArray(carrier) AS carriers,\n    groupArray(origin) AS airports\nFROM worst_pairs",
      "row_count": 1,
      "result_columns": [
        "distinct_carriers",
        "distinct_airports",
        "total_pairs",
        "carriers",
        "airports"
      ],
      "first_row": {
        "airports": [
          "ORD",
          "DFW",
          "LAX",
          "DAY",
          "DUT",
          "ORD",
          "OTH",
          "EYW",
          "PBI",
          "RSW",
          "MIA",
          "ASE",
          "SLC",
          "SJU",
          "DEN",
          "RNO",
          "EWR",
          "SCK",
          "SRQ",
          "SFO",
          "ATL",
          "EWR",
          "EWR",
          "CEC",
          "FLL"
        ],
        "carriers": [
          "DH",
          "PI",
          "PI",
          "PI",
          "AS",
          "YV",
          "OO",
          "AA",
          "F9",
          "OH",
          "B6",
          "OO",
          "EV",
          "NW",
          "EV",
          "B6",
          "OH",
          "G4",
          "B6",
          "FL",
          "OH",
          "FL",
          "XE",
          "OO",
          "B6"
        ],
        "distinct_airports": 22,
        "distinct_carriers": 14,
        "total_pairs": 25
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
      "id": "worst_pair",
      "subquestion": "Which winter carrier-airport pair ranks worst overall?",
      "answer_markdown": "The worst winter carrier-airport pair is **DH (Independence Air) departing from ORD (Chicago O'Hare)**, with only 56.58% of winter flights departing on time and an average departure delay of 28.35 minutes across 19,986 qualifying winter flights. The second-worst pair is PI (Piedmont Airlines) at DFW (61.5% on-time), followed by PI at LAX (62.95%) and PI at DAY (63.15%). The ranked query below returns the full set of worst-performing winter pairs (minimum 1,000 winter flights) sorted by on-time percentage ascending.",
      "sql": "WITH winter_pairs AS (\n    SELECT\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode AS origin,\n        count() AS winter_flights,\n        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime,\n        round(avg(DepDelayMinutes), 2) AS avg_dep_delay_min,\n        round(sum(assumeNotNull(WeatherDelay)) / nullIf(sum(assumeNotNull(WeatherDelay) + assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)), 0) * 100, 2) AS weather_pct,\n        round(sum(assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)) / nullIf(sum(assumeNotNull(WeatherDelay) + assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)), 0) * 100, 2) AS operational_pct\n    FROM ontime.fact_ontime\n    WHERE toMonth(FlightDate) IN (12, 1, 2)\n      AND Cancelled = 0\n    GROUP BY carrier, origin\n    HAVING winter_flights \u003e= 1000\n)\nSELECT\n    carrier,\n    origin,\n    winter_flights,\n    pct_ontime,\n    avg_dep_delay_min,\n    weather_pct,\n    operational_pct,\n    row_number() OVER (ORDER BY pct_ontime ASC) AS rank\nFROM winter_pairs\nORDER BY pct_ontime ASC\nLIMIT 25",
      "row_count": 25,
      "result_columns": [
        "carrier",
        "origin",
        "winter_flights",
        "pct_ontime",
        "avg_dep_delay_min",
        "weather_pct",
        "operational_pct",
        "rank"
      ],
      "first_row": {
        "avg_dep_delay_min": 28.35,
        "carrier": "DH",
        "operational_pct": 84.39,
        "origin": "ORD",
        "pct_ontime": 56.58,
        "rank": 1,
        "weather_pct": 15.61,
        "winter_flights": 19986
      }
    },
    {
      "id": "cause_mix",
      "subquestion": "Are the worst pairs driven more by weather or by operational causes?",
      "answer_markdown": "The worst winter pairs are overwhelmingly driven by **operational causes**, not weather. Across the 22 worst pairs with delay-cause data, operational delays (carrier delays + late aircraft + NAS + security) account for 70–99% of reported delay minutes in nearly every case. Late aircraft propagation alone is the single largest component for most pairs — for example, DH/ORD (37.7% late aircraft), YV/ORD (45.2% late aircraft), AA/EYW (61.1% late aircraft), and AS/DUT (66.0% late aircraft). Carrier-attributed delays are the second-largest factor. Weather is a minor contributor for most pairs (typically 1–15%), with only two notable exceptions: OH/EWR (41.4% weather) and OH/ATL (28.7% weather). Three PI pairs show no delay-cause data, consistent with older data predating structured cause reporting.",
      "sql": "WITH worst_pairs AS (\n    SELECT\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode AS origin,\n        count() AS winter_flights,\n        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime,\n        round(avg(DepDelayMinutes), 2) AS avg_dep_delay_min,\n        round(sum(assumeNotNull(CarrierDelay)), 0) AS total_carrier_delay,\n        round(sum(assumeNotNull(WeatherDelay)), 0) AS total_weather_delay,\n        round(sum(assumeNotNull(NASDelay)), 0) AS total_nas_delay,\n        round(sum(assumeNotNull(SecurityDelay)), 0) AS total_security_delay,\n        round(sum(assumeNotNull(LateAircraftDelay)), 0) AS total_late_aircraft_delay\n    FROM ontime.fact_ontime\n    WHERE toMonth(FlightDate) IN (12, 1, 2)\n      AND Cancelled = 0\n    GROUP BY carrier, origin\n    HAVING winter_flights \u003e= 1000\n    ORDER BY pct_ontime ASC\n    LIMIT 25\n)\nSELECT\n    carrier,\n    origin,\n    pct_ontime,\n    total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay AS total_delay_min,\n    round(100.0 * total_weather_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS weather_share_pct,\n    round(100.0 * total_carrier_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS carrier_share_pct,\n    round(100.0 * total_nas_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS nas_share_pct,\n    round(100.0 * total_late_aircraft_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS late_aircraft_share_pct,\n    round(100.0 * total_security_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS security_share_pct\nFROM worst_pairs\nORDER BY pct_ontime ASC",
      "row_count": 25,
      "result_columns": [
        "carrier",
        "origin",
        "pct_ontime",
        "total_delay_min",
        "weather_share_pct",
        "carrier_share_pct",
        "nas_share_pct",
        "late_aircraft_share_pct",
        "security_share_pct"
      ],
      "first_row": {
        "carrier": "DH",
        "carrier_share_pct": 34.92,
        "late_aircraft_share_pct": 37.7,
        "nas_share_pct": 11.74,
        "origin": "ORD",
        "pct_ontime": 56.58,
        "security_share_pct": 0.03,
        "total_delay_min": 353430,
        "weather_share_pct": 15.61
      }
    },
    {
      "id": "concentration_pattern",
      "subquestion": "Are the weakest pairs concentrated in a small number of carriers or airports?",
      "answer_markdown": "The 25 worst winter pairs span **14 distinct carriers and 22 distinct airports**, so the weak set is fairly dispersed across airports but shows moderate carrier concentration. Four carriers account for the majority of repeating entries: **B6 (JetBlue)** appears 4 times (MIA, RNO, SRQ, FLL), **OH (Comair)** 3 times (RSW, EWR, ATL), **OO (SkyWest)** 3 times (OTH, ASE, CEC), and **PI (Piedmont)** 3 times (DFW, LAX, DAY). On the airport side, **EWR** appears in 3 pairs (OH, FL, XE) and **ORD** in 2 pairs (DH, YV). The pattern suggests that chronic winter departure underperformance clusters around a handful of regional and low-cost carriers operating systemically poor schedules, rather than being uniformly distributed across the industry.",
      "sql": "WITH worst_pairs AS (\n    SELECT\n        IATA_CODE_Reporting_Airline AS carrier,\n        OriginCode AS origin,\n        count() AS winter_flights,\n        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime\n    FROM ontime.fact_ontime\n    WHERE toMonth(FlightDate) IN (12, 1, 2)\n      AND Cancelled = 0\n    GROUP BY carrier, origin\n    HAVING winter_flights \u003e= 1000\n    ORDER BY pct_ontime ASC\n    LIMIT 25\n)\nSELECT\n    uniq(carrier) AS distinct_carriers,\n    uniq(origin) AS distinct_airports,\n    count() AS total_pairs,\n    groupArray(carrier) AS carriers,\n    groupArray(origin) AS airports\nFROM worst_pairs",
      "row_count": 1,
      "result_columns": [
        "distinct_carriers",
        "distinct_airports",
        "total_pairs",
        "carriers",
        "airports"
      ],
      "first_row": {
        "airports": [
          "ORD",
          "DFW",
          "LAX",
          "DAY",
          "DUT",
          "ORD",
          "OTH",
          "EYW",
          "PBI",
          "RSW",
          "MIA",
          "ASE",
          "SLC",
          "SJU",
          "DEN",
          "RNO",
          "EWR",
          "SCK",
          "SRQ",
          "SFO",
          "ATL",
          "EWR",
          "EWR",
          "CEC",
          "FLL"
        ],
        "carriers": [
          "DH",
          "PI",
          "PI",
          "PI",
          "AS",
          "YV",
          "OO",
          "AA",
          "F9",
          "OH",
          "B6",
          "OO",
          "EV",
          "NW",
          "EV",
          "B6",
          "OH",
          "G4",
          "B6",
          "FL",
          "OH",
          "FL",
          "XE",
          "OO",
          "B6"
        ],
        "distinct_airports": 22,
        "distinct_carriers": 14,
        "total_pairs": 25
      }
    }
  ]
}

### Dynamic-mode additions

- Use this endpoint template for every browser query: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.