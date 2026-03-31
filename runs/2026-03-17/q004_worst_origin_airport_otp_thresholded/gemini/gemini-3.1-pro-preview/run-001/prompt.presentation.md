You are running inside qforge.

Generate valid ClickHouse SQL.
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

Question title: `Worst origin airports by departure on-time performance`

Result columns: `Origin, OriginCityName, OriginState, CompletedDepartures, DepartureOtpPct, AvgDepDelayMinutes, P90DepDelayMinutes, FirstFlightDate, LastFlightDate`

Question-specific report guidance:

Explain:

- which airport ranks worst on departure on-time performance and how poor its OTP is,
- how wide the spread is between the worst airport and the median airport within the ranked set,
- and whether the bottom 25 are dominated by major hubs or show a more mixed airport profile.

Create `visual.html` using the `ontime-analyst-dashboard` skill.

The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.

General visual rules:

- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.
- Keep report logic out of `visual.html`; qforge already renders `report.md` separately.

Saved SQL to preserve in the final page contract:

```sql
WITH OriginStats AS (
    SELECT
        Origin,
        any(OriginCityName) AS OriginCityName,
        any(OriginState) AS OriginState,
        count() AS CompletedDepartures,
        countIf(DepDel15 = 0) / count() AS DepartureOtpPct,
        avg(DepDelayMinutes) AS AvgDepDelayMinutes,
        quantile(0.90)(DepDelayMinutes) AS P90DepDelayMinutes,
        min(FlightDate) AS FirstFlightDate,
        max(FlightDate) AS LastFlightDate
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Origin
    HAVING CompletedDepartures >= 50000
)
SELECT
    Origin,
    OriginCityName,
    OriginState,
    CompletedDepartures,
    DepartureOtpPct,
    AvgDepDelayMinutes,
    P90DepDelayMinutes,
    FirstFlightDate,
    LastFlightDate
FROM OriginStats
ORDER BY
    DepartureOtpPct ASC,
    AvgDepDelayMinutes DESC,
    CompletedDepartures DESC,
    Origin ASC
LIMIT 25;
```

Visual context:

- Question title: `Worst origin airports by departure on-time performance`
- Visual mode: `dynamic`
- Visual type: `html_ranked_dashboard`

Question-specific visual guidance:

Build a dashboard that:

- shows KPI cards for worst airport, worst OTP, highest average departure delay among ranked airports, and qualifying airport count
- renders a ranked horizontal bar or lollipop chart for the worst 25 airports by departure OTP
- renders a scatter plot of `CompletedDepartures` vs `DepartureOtpPct`
- renders a detail table with the full ranked result
- derives chart extents and highlighted airports from fetched data instead of hardcoding them
- uses one accent treatment for the worst 5 airports
- annotates the single worst airport in both charts
- keeps the scatter plot readable despite skewed volume differences

Dynamic-mode additions:

- Build the page in dynamic mode using the `ontime-analyst-dashboard` skill contract.
- Execute the embedded saved SQL in the browser as the primary query.
- Keep the embedded saved SQL authoritative for the artifact.
- Surface every browser query in a visible query ledger.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Keep additional browser queries limited to explicit enrichment or drill-down that materially improves the visualization.

Build a dashboard that:

- shows KPI cards for worst airport, worst OTP, highest average departure delay among ranked airports, and qualifying airport count
- renders a ranked horizontal bar or lollipop chart for the worst 25 airports by departure OTP
- renders a scatter plot of `CompletedDepartures` vs `DepartureOtpPct`
- renders a detail table with the full ranked result
- derives chart extents and highlighted airports from fetched data instead of hardcoding them
- uses one accent treatment for the worst 5 airports
- annotates the single worst airport in both charts
- keeps the scatter plot readable despite skewed volume differences