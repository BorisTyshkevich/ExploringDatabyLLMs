- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `Yearly carrier leadership by completed flights`
- Visual mode: `dynamic`
- Presentation target: `html`
- Visual type: `html_timeseries`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

Build a dashboard that:

- uses `most_frequent_leader` as the primary saved SQL already provided in the prompt
- uses `leadership_swing`, `sharpest_transition`, and `stability_pattern` as supporting queries when they materially improve the dashboard
- shows KPI cards for total years analyzed, distinct annual leaders, largest leader share gap, and the sharpest leadership transition
- renders a bump chart for yearly carrier rank among the top carriers
- renders a time series of yearly completed-flight share for the leading carriers
- includes a compact table of all leadership-change years with prior leader, new leader, share swing, and share gap
- includes a narrative takeaway about whether the market shows long stable eras or frequent turnover
- highlights true leadership transitions only; do not treat the first year as a transition
- makes the sharpest leadership transition visually distinct
- derives all shown carriers from fetched data instead of hardcoding airline names
- shows supporting queries in the query ledger when used

### Data Source

SQL query for primary data source:

```sql
WITH carrier_year AS (
    SELECT
        Year,
        Reporting_Airline AS carrier,
        countIf(Cancelled = 0) AS completed_flights
    FROM ontime.fact_ontime
    GROUP BY Year, carrier
),
ranked AS (
    SELECT
        Year,
        carrier,
        completed_flights,
        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk
    FROM carrier_year
),
leaders AS (
    SELECT Year, carrier
    FROM ranked
    WHERE rnk = 1
)
SELECT
    carrier,
    count() AS years_led
FROM leaders
GROUP BY carrier
ORDER BY years_led DESC
```

Data example/snippet:

{
  "question_title": "Yearly carrier leadership by completed flights",
  "result_columns": null,
  "row_count": 4,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which carrier leads most often across the full time range?",
      "answer_markdown": "WN (Southwest Airlines) leads most often, topping the industry in 26 of 39 years (2000–2025). Delta (DL) led for 11 years (1987–1999, except 1990–1991 when US Airways held the top), and US Airways (US) led for 2 years (1990–1991).",
      "sql": "WITH carrier_year AS (\n    SELECT\n        Year,\n        Reporting_Airline AS carrier,\n        countIf(Cancelled = 0) AS completed_flights\n    FROM ontime.fact_ontime\n    GROUP BY Year, carrier\n),\nranked AS (\n    SELECT\n        Year,\n        carrier,\n        completed_flights,\n        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk\n    FROM carrier_year\n),\nleaders AS (\n    SELECT Year, carrier\n    FROM ranked\n    WHERE rnk = 1\n)\nSELECT\n    carrier,\n    count() AS years_led\nFROM leaders\nGROUP BY carrier\nORDER BY years_led DESC",
      "row_count": 3,
      "result_columns": [
        "carrier",
        "years_led"
      ],
      "first_row": {
        "carrier": "WN",
        "years_led": 26
      }
    },
    {
      "id": "q2",
      "subquestion": "When leadership changes, how large is the swing versus the prior leader?",
      "answer_markdown": "There are 3 genuine leadership transitions in the full history. The 1990 transition from DL to US produced the largest swing: US Airways flew 991,989 flights versus Delta's prior-year total of 778,612, a swing of +213,377. The 1992 reversal back to DL was narrow at +9,842 flights over US's prior year, and the 2000 shift to WN was similarly narrow at +10,273 flights over DL's prior year. The average swing across all three transitions is about 78,000 flights, but that figure is dominated entirely by the 1990 outlier.",
      "sql": "WITH carrier_year AS (\n    SELECT\n        Year,\n        Reporting_Airline AS carrier,\n        countIf(Cancelled = 0) AS completed_flights\n    FROM ontime.fact_ontime\n    GROUP BY Year, carrier\n),\nranked AS (\n    SELECT\n        Year,\n        carrier,\n        completed_flights,\n        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk\n    FROM carrier_year\n),\nleader_per_year AS (\n    SELECT\n        Year,\n        maxIf(carrier, rnk = 1)          AS leader,\n        maxIf(completed_flights, rnk = 1) AS leader_flights\n    FROM ranked\n    GROUP BY Year\n),\ntransitions AS (\n    SELECT\n        Year,\n        leader,\n        leader_flights,\n        lagInFrame(leader)          OVER (ORDER BY Year) AS prev_leader,\n        lagInFrame(leader_flights)  OVER (ORDER BY Year) AS prev_leader_flights\n    FROM leader_per_year\n    ORDER BY Year\n)\nSELECT\n    Year,\n    leader          AS new_leader,\n    leader_flights  AS new_leader_flights,\n    prev_leader,\n    prev_leader_flights,\n    leader_flights - prev_leader_flights AS swing\nFROM transitions\nWHERE leader != prev_leader\n  AND prev_leader != ''\nORDER BY Year",
      "row_count": 3,
      "result_columns": [
        "Year",
        "new_leader",
        "new_leader_flights",
        "prev_leader",
        "prev_leader_flights",
        "swing"
      ],
      "first_row": {
        "Year": 1990,
        "new_leader": "US",
        "new_leader_flights": 991989,
        "prev_leader": "DL",
        "prev_leader_flights": 778612,
        "swing": 213377
      }
    },
    {
      "id": "q3",
      "subquestion": "Which transition is the sharpest?",
      "answer_markdown": "The sharpest transition is 1990, when US Airways (US) displaced Delta (DL) with a swing of +213,377 completed flights versus Delta's 1989 output, and held a 174,169-flight lead over Delta as runner-up within that year. This is more than 20 times larger than either of the other two transitions (1992 and 2000).",
      "sql": "WITH carrier_year AS (\n    SELECT\n        Year,\n        Reporting_Airline AS carrier,\n        countIf(Cancelled = 0) AS completed_flights\n    FROM ontime.fact_ontime\n    GROUP BY Year, carrier\n),\nranked AS (\n    SELECT\n        Year,\n        carrier,\n        completed_flights,\n        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk\n    FROM carrier_year\n),\nleader_runner AS (\n    SELECT\n        Year,\n        maxIf(carrier, rnk = 1)          AS leader,\n        maxIf(completed_flights, rnk = 1) AS leader_flights,\n        maxIf(carrier, rnk = 2)          AS runner_up,\n        maxIf(completed_flights, rnk = 2) AS runner_up_flights\n    FROM ranked\n    GROUP BY Year\n),\nwith_lag AS (\n    SELECT\n        Year,\n        leader,\n        leader_flights,\n        runner_up,\n        runner_up_flights,\n        leader_flights - runner_up_flights                       AS gap_vs_runner_up,\n        lagInFrame(leader)         OVER (ORDER BY Year)          AS prev_leader,\n        lagInFrame(leader_flights) OVER (ORDER BY Year)          AS prev_leader_flights,\n        leader_flights - lagInFrame(leader_flights) OVER (ORDER BY Year) AS swing_vs_prior\n    FROM leader_runner\n    ORDER BY Year\n)\nSELECT\n    Year,\n    leader          AS new_leader,\n    leader_flights,\n    prev_leader,\n    prev_leader_flights,\n    swing_vs_prior,\n    gap_vs_runner_up\nFROM with_lag\nWHERE leader != prev_leader\n  AND prev_leader != ''\nORDER BY swing_vs_prior DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Year",
        "new_leader",
        "leader_flights",
        "prev_leader",
        "prev_leader_flights",
        "swing_vs_prior",
        "gap_vs_runner_up"
      ],
      "first_row": {
        "Year": 1990,
        "gap_vs_runner_up": 174169,
        "leader_flights": 991989,
        "new_leader": "US",
        "prev_leader": "DL",
        "prev_leader_flights": 778612,
        "swing_vs_prior": 213377
      }
    },
    {
      "id": "q4",
      "subquestion": "Does the market show long stable eras, or frequent turnover at the top?",
      "answer_markdown": "The market is defined by long stable eras, not frequent turnover. In 38 full years of data (1988–2025), there were only 3 leadership transitions. Delta (DL) dominated for roughly 11 years (1987–1999 with a brief 2-year interruption by US Airways in 1990–1991). Southwest (WN) has held an unbroken reign since 2000—26 consecutive years. Two carriers account for the top position in all but 2 years of the entire dataset.",
      "sql": "WITH carrier_year AS (\n    SELECT\n        Year,\n        Reporting_Airline AS carrier,\n        countIf(Cancelled = 0) AS completed_flights\n    FROM ontime.fact_ontime\n    GROUP BY Year, carrier\n),\nranked AS (\n    SELECT\n        Year,\n        carrier,\n        completed_flights,\n        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk\n    FROM carrier_year\n),\nleader_per_year AS (\n    SELECT Year, maxIf(carrier, rnk = 1) AS leader\n    FROM ranked\n    GROUP BY Year\n),\nwith_era_change AS (\n    SELECT\n        Year,\n        leader,\n        leader != lagInFrame(leader) OVER (ORDER BY Year) AS era_start\n    FROM leader_per_year\n    ORDER BY Year\n),\nwith_era_id AS (\n    SELECT\n        Year,\n        leader,\n        sum(era_start) OVER (ORDER BY Year ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS era_id\n    FROM with_era_change\n)\nSELECT\n    leader,\n    min(Year) AS era_start_year,\n    max(Year) AS era_end_year,\n    count()   AS years_in_era\nFROM with_era_id\nGROUP BY era_id, leader\nORDER BY era_start_year",
      "row_count": 4,
      "result_columns": [
        "leader",
        "era_start_year",
        "era_end_year",
        "years_in_era"
      ],
      "first_row": {
        "era_end_year": 1989,
        "era_start_year": 1987,
        "leader": "DL",
        "years_in_era": 3
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
  "question_title": "Yearly carrier leadership by completed flights",
  "result_columns": null,
  "row_count": 4,
  "mode_hint": "This visual pass receives only verified subquestion answers plus proof-query previews: row count, column names, and the first result row for each query.",
  "query_summaries": [
    {
      "id": "q1",
      "subquestion": "Which carrier leads most often across the full time range?",
      "answer_markdown": "WN (Southwest Airlines) leads most often, topping the industry in 26 of 39 years (2000–2025). Delta (DL) led for 11 years (1987–1999, except 1990–1991 when US Airways held the top), and US Airways (US) led for 2 years (1990–1991).",
      "sql": "WITH carrier_year AS (\n    SELECT\n        Year,\n        Reporting_Airline AS carrier,\n        countIf(Cancelled = 0) AS completed_flights\n    FROM ontime.fact_ontime\n    GROUP BY Year, carrier\n),\nranked AS (\n    SELECT\n        Year,\n        carrier,\n        completed_flights,\n        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk\n    FROM carrier_year\n),\nleaders AS (\n    SELECT Year, carrier\n    FROM ranked\n    WHERE rnk = 1\n)\nSELECT\n    carrier,\n    count() AS years_led\nFROM leaders\nGROUP BY carrier\nORDER BY years_led DESC",
      "row_count": 3,
      "result_columns": [
        "carrier",
        "years_led"
      ],
      "first_row": {
        "carrier": "WN",
        "years_led": 26
      }
    },
    {
      "id": "q2",
      "subquestion": "When leadership changes, how large is the swing versus the prior leader?",
      "answer_markdown": "There are 3 genuine leadership transitions in the full history. The 1990 transition from DL to US produced the largest swing: US Airways flew 991,989 flights versus Delta's prior-year total of 778,612, a swing of +213,377. The 1992 reversal back to DL was narrow at +9,842 flights over US's prior year, and the 2000 shift to WN was similarly narrow at +10,273 flights over DL's prior year. The average swing across all three transitions is about 78,000 flights, but that figure is dominated entirely by the 1990 outlier.",
      "sql": "WITH carrier_year AS (\n    SELECT\n        Year,\n        Reporting_Airline AS carrier,\n        countIf(Cancelled = 0) AS completed_flights\n    FROM ontime.fact_ontime\n    GROUP BY Year, carrier\n),\nranked AS (\n    SELECT\n        Year,\n        carrier,\n        completed_flights,\n        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk\n    FROM carrier_year\n),\nleader_per_year AS (\n    SELECT\n        Year,\n        maxIf(carrier, rnk = 1)          AS leader,\n        maxIf(completed_flights, rnk = 1) AS leader_flights\n    FROM ranked\n    GROUP BY Year\n),\ntransitions AS (\n    SELECT\n        Year,\n        leader,\n        leader_flights,\n        lagInFrame(leader)          OVER (ORDER BY Year) AS prev_leader,\n        lagInFrame(leader_flights)  OVER (ORDER BY Year) AS prev_leader_flights\n    FROM leader_per_year\n    ORDER BY Year\n)\nSELECT\n    Year,\n    leader          AS new_leader,\n    leader_flights  AS new_leader_flights,\n    prev_leader,\n    prev_leader_flights,\n    leader_flights - prev_leader_flights AS swing\nFROM transitions\nWHERE leader != prev_leader\n  AND prev_leader != ''\nORDER BY Year",
      "row_count": 3,
      "result_columns": [
        "Year",
        "new_leader",
        "new_leader_flights",
        "prev_leader",
        "prev_leader_flights",
        "swing"
      ],
      "first_row": {
        "Year": 1990,
        "new_leader": "US",
        "new_leader_flights": 991989,
        "prev_leader": "DL",
        "prev_leader_flights": 778612,
        "swing": 213377
      }
    },
    {
      "id": "q3",
      "subquestion": "Which transition is the sharpest?",
      "answer_markdown": "The sharpest transition is 1990, when US Airways (US) displaced Delta (DL) with a swing of +213,377 completed flights versus Delta's 1989 output, and held a 174,169-flight lead over Delta as runner-up within that year. This is more than 20 times larger than either of the other two transitions (1992 and 2000).",
      "sql": "WITH carrier_year AS (\n    SELECT\n        Year,\n        Reporting_Airline AS carrier,\n        countIf(Cancelled = 0) AS completed_flights\n    FROM ontime.fact_ontime\n    GROUP BY Year, carrier\n),\nranked AS (\n    SELECT\n        Year,\n        carrier,\n        completed_flights,\n        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk\n    FROM carrier_year\n),\nleader_runner AS (\n    SELECT\n        Year,\n        maxIf(carrier, rnk = 1)          AS leader,\n        maxIf(completed_flights, rnk = 1) AS leader_flights,\n        maxIf(carrier, rnk = 2)          AS runner_up,\n        maxIf(completed_flights, rnk = 2) AS runner_up_flights\n    FROM ranked\n    GROUP BY Year\n),\nwith_lag AS (\n    SELECT\n        Year,\n        leader,\n        leader_flights,\n        runner_up,\n        runner_up_flights,\n        leader_flights - runner_up_flights                       AS gap_vs_runner_up,\n        lagInFrame(leader)         OVER (ORDER BY Year)          AS prev_leader,\n        lagInFrame(leader_flights) OVER (ORDER BY Year)          AS prev_leader_flights,\n        leader_flights - lagInFrame(leader_flights) OVER (ORDER BY Year) AS swing_vs_prior\n    FROM leader_runner\n    ORDER BY Year\n)\nSELECT\n    Year,\n    leader          AS new_leader,\n    leader_flights,\n    prev_leader,\n    prev_leader_flights,\n    swing_vs_prior,\n    gap_vs_runner_up\nFROM with_lag\nWHERE leader != prev_leader\n  AND prev_leader != ''\nORDER BY swing_vs_prior DESC\nLIMIT 1",
      "row_count": 1,
      "result_columns": [
        "Year",
        "new_leader",
        "leader_flights",
        "prev_leader",
        "prev_leader_flights",
        "swing_vs_prior",
        "gap_vs_runner_up"
      ],
      "first_row": {
        "Year": 1990,
        "gap_vs_runner_up": 174169,
        "leader_flights": 991989,
        "new_leader": "US",
        "prev_leader": "DL",
        "prev_leader_flights": 778612,
        "swing_vs_prior": 213377
      }
    },
    {
      "id": "q4",
      "subquestion": "Does the market show long stable eras, or frequent turnover at the top?",
      "answer_markdown": "The market is defined by long stable eras, not frequent turnover. In 38 full years of data (1988–2025), there were only 3 leadership transitions. Delta (DL) dominated for roughly 11 years (1987–1999 with a brief 2-year interruption by US Airways in 1990–1991). Southwest (WN) has held an unbroken reign since 2000—26 consecutive years. Two carriers account for the top position in all but 2 years of the entire dataset.",
      "sql": "WITH carrier_year AS (\n    SELECT\n        Year,\n        Reporting_Airline AS carrier,\n        countIf(Cancelled = 0) AS completed_flights\n    FROM ontime.fact_ontime\n    GROUP BY Year, carrier\n),\nranked AS (\n    SELECT\n        Year,\n        carrier,\n        completed_flights,\n        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk\n    FROM carrier_year\n),\nleader_per_year AS (\n    SELECT Year, maxIf(carrier, rnk = 1) AS leader\n    FROM ranked\n    GROUP BY Year\n),\nwith_era_change AS (\n    SELECT\n        Year,\n        leader,\n        leader != lagInFrame(leader) OVER (ORDER BY Year) AS era_start\n    FROM leader_per_year\n    ORDER BY Year\n),\nwith_era_id AS (\n    SELECT\n        Year,\n        leader,\n        sum(era_start) OVER (ORDER BY Year ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS era_id\n    FROM with_era_change\n)\nSELECT\n    leader,\n    min(Year) AS era_start_year,\n    max(Year) AS era_end_year,\n    count()   AS years_in_era\nFROM with_era_id\nGROUP BY era_id, leader\nORDER BY era_start_year",
      "row_count": 4,
      "result_columns": [
        "leader",
        "era_start_year",
        "era_end_year",
        "years_in_era"
      ],
      "first_row": {
        "era_end_year": 1989,
        "era_start_year": 1987,
        "leader": "DL",
        "years_in_era": 3
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