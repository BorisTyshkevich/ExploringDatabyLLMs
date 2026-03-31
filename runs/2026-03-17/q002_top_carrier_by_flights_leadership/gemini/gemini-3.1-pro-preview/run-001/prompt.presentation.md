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
WITH YearlyCarrier AS (
    SELECT
        Year,
        Reporting_Airline,
        count() AS CompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year, Reporting_Airline
),
YearlyTotals AS (
    SELECT
        Year,
        sum(CompletedFlights) AS TotalCompletedFlights
    FROM YearlyCarrier
    GROUP BY Year
),
CarrierShares AS (
    SELECT
        c.Year,
        c.Reporting_Airline,
        c.CompletedFlights,
        (c.CompletedFlights * 100.0) / t.TotalCompletedFlights AS SharePct,
        row_number() OVER (PARTITION BY c.Year ORDER BY c.CompletedFlights DESC, c.Reporting_Airline ASC) AS RankInYear
    FROM YearlyCarrier c
    JOIN YearlyTotals t ON c.Year = t.Year
),
YearlyLeaders AS (
    SELECT
        Year,
        maxIf(Reporting_Airline, RankInYear = 1) AS LeaderReportingAirline,
        maxIf(Reporting_Airline, RankInYear = 2) AS RunnerUpReportingAirline,
        maxIf(SharePct, RankInYear = 1) - maxIf(SharePct, RankInYear = 2) AS LeaderShareGapPctPts,
        maxIf(SharePct, RankInYear = 1) AS LeaderSharePct
    FROM CarrierShares
    GROUP BY Year
),
LeaderTransitions AS (
    SELECT
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderShareGapPctPts,
        LeaderSharePct,
        lag(LeaderReportingAirline) OVER (ORDER BY Year) AS RawPriorLeader,
        lag(LeaderSharePct) OVER (ORDER BY Year) AS RawPriorShare,
        row_number() OVER (ORDER BY Year) AS YearNum
    FROM YearlyLeaders
),
ProcessedTransitions AS (
    SELECT
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderShareGapPctPts,
        CASE WHEN YearNum = 1 THEN CAST(NULL AS Nullable(String)) ELSE CAST(RawPriorLeader AS Nullable(String)) END AS PriorYearLeaderReportingAirline,
        CASE WHEN YearNum = 1 THEN 0 WHEN LeaderReportingAirline != RawPriorLeader THEN 1 ELSE 0 END AS LeaderChanged,
        CASE WHEN YearNum = 1 THEN CAST(NULL AS Nullable(Float64)) ELSE CAST(LeaderSharePct - RawPriorShare AS Nullable(Float64)) END AS LeaderShareChangePctPts
    FROM LeaderTransitions
),
FinalData AS (
    SELECT
        CAST('carrier_year' AS String) AS RowType,
        c.Year AS Year,
        CAST(c.Reporting_Airline AS String) AS Reporting_Airline,
        CAST(c.RankInYear AS UInt32) AS RankInYear,
        CAST(c.CompletedFlights AS UInt64) AS CompletedFlights,
        CAST(c.SharePct AS Float64) AS SharePct,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.LeaderReportingAirline AS String) ELSE CAST('' AS String) END AS LeaderReportingAirline,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.RunnerUpReportingAirline AS String) ELSE CAST('' AS String) END AS RunnerUpReportingAirline,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.LeaderShareGapPctPts AS Nullable(Float64)) ELSE CAST(NULL AS Nullable(Float64)) END AS LeaderShareGapPctPts,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.PriorYearLeaderReportingAirline AS Nullable(String)) ELSE CAST(NULL AS Nullable(String)) END AS PriorYearLeaderReportingAirline,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.LeaderChanged AS UInt8) ELSE CAST(0 AS UInt8) END AS LeaderChanged,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.LeaderShareChangePctPts AS Nullable(Float64)) ELSE CAST(NULL AS Nullable(Float64)) END AS LeaderShareChangePctPts
    FROM CarrierShares c
    JOIN ProcessedTransitions t ON c.Year = t.Year
    WHERE c.RankInYear <= 5
),
SummaryData AS (
    SELECT
        CAST('year_summary' AS String) AS RowType,
        t.Year AS Year,
        CAST('' AS String) AS Reporting_Airline,
        CAST(0 AS UInt32) AS RankInYear,
        CAST(tot.TotalCompletedFlights AS UInt64) AS CompletedFlights,
        CAST(100.0 AS Float64) AS SharePct,
        CAST(t.LeaderReportingAirline AS String) AS LeaderReportingAirline,
        CAST(t.RunnerUpReportingAirline AS String) AS RunnerUpReportingAirline,
        CAST(t.LeaderShareGapPctPts AS Nullable(Float64)) AS LeaderShareGapPctPts,
        CAST(t.PriorYearLeaderReportingAirline AS Nullable(String)) AS PriorYearLeaderReportingAirline,
        CAST(t.LeaderChanged AS UInt8) AS LeaderChanged,
        CAST(t.LeaderShareChangePctPts AS Nullable(Float64)) AS LeaderShareChangePctPts
    FROM ProcessedTransitions t
    JOIN YearlyTotals tot ON t.Year = tot.Year
)
SELECT * FROM FinalData
UNION ALL
SELECT * FROM SummaryData
ORDER BY Year ASC, RowType ASC, RankInYear ASC, Reporting_Airline ASC
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