You are running inside qforge.

Use the configured MCP server for all data access.
Do not construct raw OpenAPI URLs manually.
Stay within the configured dataset scope.
Do not reference tables outside the allowed dataset constraints.

Dataset constraints:

- Use `default.ontime_v2` as the primary fact table.
- Do not reference `default.ontime`.

You are generating presentation artifacts inside qforge.

The harness already executed the saved SQL and produced `result.json`.
Do not invent KPIs or data values.
The report is a Markdown template that qforge will render from `result.json`.

Return exactly these fenced sections:

```report
Use placeholders only where data is needed.
Allowed placeholders: {{row_count}}, {{generated_at}}, {{columns_csv}}, {{question_title}}, {{data_overview_md}}, {{result_table_md}}
```

```html
<!doctype html>
<html>...</html>
```

Report rules:

- The report must be Markdown.
- The report must be a template, not a data-filled summary.
- Prefer `{{data_overview_md}}` and `{{result_table_md}}` for JSON-derived sections.
- Keep the report concise and analytical.

Question title: `Yearly carrier leadership by completed flights`

Result columns: `RowType, Year, Reporting_Airline, RankInYear, CompletedFlights, SharePct, LeaderReportingAirline, RunnerUpReportingAirline, LeaderShareGapPctPts, PriorYearLeaderReportingAirline, LeaderChanged, LeaderShareChangePctPts`

Question-specific report guidance:

Explain:

- how often each carrier appears as the annual leader across the full time range,
- every true leadership transition with the prior leader, new leader, and share swing,
- which transition is the sharpest by the question's ranking rule,
- and whether the series contains long periods of stable dominance.

Create `visual.html` using the `ontime-analyst-dashboard` skill.

The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.

Saved SQL to embed directly in the final page:

```sql
WITH 
-- Step 1: Annual totals for completed flights
annual_totals AS (
    SELECT 
        Year,
        count() AS TotalCompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year
),

-- Step 2: Carrier-year aggregates for completed flights
carrier_year_stats AS (
    SELECT 
        Year,
        Reporting_Airline,
        count() AS CompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year, Reporting_Airline
),

-- Step 3: Ranked carrier-year with share calculation
ranked_carriers AS (
    SELECT 
        cy.Year,
        cy.Reporting_Airline,
        cy.CompletedFlights,
        round(100.0 * cy.CompletedFlights / at.TotalCompletedFlights, 4) AS SharePct,
        row_number() OVER (PARTITION BY cy.Year ORDER BY cy.CompletedFlights DESC, cy.Reporting_Airline ASC) AS RankInYear
    FROM carrier_year_stats cy
    JOIN annual_totals at ON cy.Year = at.Year
),

-- Step 4: Extract leaders (rank 1) and runners-up (rank 2) per year
year_leaders AS (
    SELECT 
        Year,
        argMax(Reporting_Airline, RankInYear = 1) AS LeaderReportingAirline,
        argMax(SharePct, RankInYear = 1) AS LeaderSharePct,
        argMax(Reporting_Airline, RankInYear = 2) AS RunnerUpReportingAirline,
        argMax(SharePct, RankInYear = 2) AS RunnerUpSharePct
    FROM ranked_carriers
    WHERE RankInYear IN (1, 2)
    GROUP BY Year
),

-- Step 5: Compute leader transitions with YoY share change
leader_transitions AS (
    SELECT 
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        round(LeaderSharePct - RunnerUpSharePct, 4) AS LeaderShareGapPctPts,
        LeaderSharePct,
        lagInFrame(LeaderReportingAirline) OVER (ORDER BY Year) AS PriorYearLeaderReportingAirline,
        lagInFrame(LeaderSharePct) OVER (ORDER BY Year) AS PriorYearLeaderSharePct,
        row_number() OVER (ORDER BY Year) AS YearSeq
    FROM year_leaders
),

-- Step 6: Final leader transition metrics
leader_metrics AS (
    SELECT 
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderShareGapPctPts,
        LeaderSharePct,
        PriorYearLeaderReportingAirline,
        CASE WHEN YearSeq = 1 THEN NULL ELSE round(LeaderSharePct - PriorYearLeaderSharePct, 4) END AS LeaderShareChangePctPts,
        CASE 
            WHEN YearSeq = 1 THEN 0
            WHEN LeaderReportingAirline != PriorYearLeaderReportingAirline THEN 1
            ELSE 0
        END AS LeaderChanged
    FROM leader_transitions
),

-- Step 7: Top 5 carriers per year for bump chart
top5_carriers AS (
    SELECT *
    FROM ranked_carriers
    WHERE RankInYear <= 5
)

-- Final output combining carrier_year rows
SELECT 
    'carrier_year' AS RowType,
    tc.Year,
    tc.Reporting_Airline,
    tc.RankInYear,
    tc.CompletedFlights,
    tc.SharePct,
    lm.LeaderReportingAirline,
    lm.RunnerUpReportingAirline,
    lm.LeaderShareGapPctPts,
    lm.PriorYearLeaderReportingAirline,
    lm.LeaderChanged,
    lm.LeaderShareChangePctPts
FROM top5_carriers tc
JOIN leader_metrics lm ON tc.Year = lm.Year
ORDER BY 
    tc.Year ASC,
    RowType ASC,
    tc.RankInYear ASC,
    tc.Reporting_Airline ASC
```

Visual context:

- Question title: `Yearly carrier leadership by completed flights`
- Visual type: `html_timeseries`

Question-specific visual guidance:

Build a dashboard that:

- shows KPI cards for total years analyzed, distinct annual leaders, largest leader share gap, and the sharpest leadership transition
- renders a bump chart for yearly carrier rank among the top carriers
- renders a time series of yearly completed-flight share for the leading carriers
- includes a compact table of all leadership-change years with prior leader, new leader, share swing, and share gap
- highlights true leadership transitions only; do not treat the first year as a transition
- makes the sharpest leadership transition visually distinct
- derives all shown carriers from fetched data instead of hardcoding airline names