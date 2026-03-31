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
WITH annual_totals AS (
    SELECT
        Year,
        count() AS YearCompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year
),
annual_carrier_counts AS (
    SELECT
        Year,
        Reporting_Airline,
        count() AS CompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY
        Year,
        Reporting_Airline
),
ranked_carriers AS (
    SELECT
        acc.Year,
        acc.Reporting_Airline,
        row_number() OVER (
            PARTITION BY acc.Year
            ORDER BY acc.CompletedFlights DESC, acc.Reporting_Airline ASC
        ) AS RankInYear,
        acc.CompletedFlights,
        100.0 * acc.CompletedFlights / at.YearCompletedFlights AS SharePct
    FROM annual_carrier_counts AS acc
    INNER JOIN annual_totals AS at
        ON acc.Year = at.Year
),
leader_runner_up AS (
    SELECT
        Year,
        maxIf(Reporting_Airline, RankInYear = 1) AS LeaderReportingAirline,
        maxIf(Reporting_Airline, RankInYear = 2) AS RunnerUpReportingAirline,
        maxIf(SharePct, RankInYear = 1) AS LeaderSharePct,
        maxIf(SharePct, RankInYear = 2) AS RunnerUpSharePct
    FROM ranked_carriers
    WHERE RankInYear <= 2
    GROUP BY Year
),
leader_transitions AS (
    SELECT
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderSharePct - RunnerUpSharePct AS LeaderShareGapPctPts,
        prior_year_leader AS PriorYearLeaderReportingAirline,
        if(prior_year_leader IS NULL, 0, LeaderReportingAirline != prior_year_leader) AS LeaderChanged,
        if(
            prior_year_leader IS NULL,
            cast(NULL, 'Nullable(Float64)'),
            LeaderSharePct - prior_year_share_pct
        ) AS LeaderShareChangePctPts
    FROM (
        SELECT
            Year,
            LeaderReportingAirline,
            RunnerUpReportingAirline,
            LeaderSharePct,
            RunnerUpSharePct,
            lagInFrame(toNullable(LeaderReportingAirline)) OVER w AS prior_year_leader,
            lagInFrame(toNullable(LeaderSharePct)) OVER w AS prior_year_share_pct
        FROM leader_runner_up
        WINDOW w AS (ORDER BY Year ASC ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)
    )
)
SELECT
    'carrier_year' AS RowType,
    rc.Year,
    rc.Reporting_Airline,
    rc.RankInYear,
    rc.CompletedFlights,
    rc.SharePct,
    if(rc.RankInYear = 1, lt.LeaderReportingAirline, cast(NULL, 'Nullable(String)')) AS LeaderReportingAirline,
    if(rc.RankInYear = 1, lt.RunnerUpReportingAirline, cast(NULL, 'Nullable(String)')) AS RunnerUpReportingAirline,
    if(rc.RankInYear = 1, lt.LeaderShareGapPctPts, cast(NULL, 'Nullable(Float64)')) AS LeaderShareGapPctPts,
    if(rc.RankInYear = 1, lt.PriorYearLeaderReportingAirline, cast(NULL, 'Nullable(String)')) AS PriorYearLeaderReportingAirline,
    if(rc.RankInYear = 1, toUInt8(lt.LeaderChanged), cast(NULL, 'Nullable(UInt8)')) AS LeaderChanged,
    if(rc.RankInYear = 1, lt.LeaderShareChangePctPts, cast(NULL, 'Nullable(Float64)')) AS LeaderShareChangePctPts
FROM ranked_carriers AS rc
INNER JOIN leader_transitions AS lt
    ON rc.Year = lt.Year
WHERE rc.RankInYear <= 5
ORDER BY
    rc.Year ASC,
    RowType,
    rc.RankInYear ASC,
    rc.Reporting_Airline ASC
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