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

Question title: `Delta ATL departure delay hotspots by destination and time block`

Result columns: `RowType, MonthStart, Dest, DepTimeBlk, QualifyingMonths, CompletedFlights, AvgDepDelayMinutes, P90DepDelayMinutes, DepDel15Pct, FirstQualifyingMonth, LastQualifyingMonth, HotspotRank`

Question-specific report guidance:

Explain:

- which destination and time block is the worst hotspot,
- whether the hotspot is persistent or concentrated in a narrower era,
- and what the top 5 hotspot cells suggest about ATL departure pressure for Delta.

Create `visual.html` using the `ontime-analyst-dashboard` skill.

The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.

Saved SQL to embed directly in the final page:

```sql
WITH
raw_flights AS (
    SELECT
        toStartOfMonth(FlightDate) AS MonthStart,
        Dest,
        DepTimeBlk,
        DepDelayMinutes,
        DepDel15
    FROM default.ontime_v2
    WHERE IATA_CODE_Reporting_Airline = 'DL'
      AND Origin = 'ATL'
      AND Cancelled = 0
),
qualifying_months AS (
    SELECT
        MonthStart,
        Dest,
        DepTimeBlk
    FROM raw_flights
    GROUP BY
        MonthStart,
        Dest,
        DepTimeBlk
    HAVING count() >= 40
),
qualifying_raw_flights AS (
    SELECT
        r.MonthStart,
        r.Dest,
        r.DepTimeBlk,
        r.DepDelayMinutes,
        r.DepDel15
    FROM raw_flights r
    INNER JOIN qualifying_months q
        ON r.MonthStart = q.MonthStart
       AND r.Dest = q.Dest
       AND r.DepTimeBlk = q.DepTimeBlk
),
hotspot_stats AS (
    SELECT
        Dest,
        DepTimeBlk,
        count(DISTINCT MonthStart) AS QualifyingMonths,
        count() AS CompletedFlights,
        round(avg(DepDelayMinutes), 2) AS AvgDepDelayMinutes,
        round(quantile(0.9)(DepDelayMinutes), 2) AS P90DepDelayMinutes,
        round(avg(DepDel15) * 100, 2) AS DepDel15Pct,
        min(MonthStart) AS FirstQualifyingMonth,
        max(MonthStart) AS LastQualifyingMonth
    FROM qualifying_raw_flights
    GROUP BY
        Dest,
        DepTimeBlk
    HAVING count() >= 1000
),
top_hotspots AS (
    SELECT
        Dest,
        DepTimeBlk,
        QualifyingMonths,
        CompletedFlights,
        AvgDepDelayMinutes,
        P90DepDelayMinutes,
        DepDel15Pct,
        FirstQualifyingMonth,
        LastQualifyingMonth,
        row_number() OVER (ORDER BY AvgDepDelayMinutes DESC, P90DepDelayMinutes DESC, DepDel15Pct DESC, CompletedFlights DESC) AS HotspotRank
    FROM hotspot_stats
    ORDER BY HotspotRank
    LIMIT 20
),
combined_results AS (
    SELECT
        'hotspot_summary' AS RowType,
        CAST(NULL AS Nullable(Date)) AS MonthStart,
        Dest,
        DepTimeBlk,
        QualifyingMonths,
        CompletedFlights,
        AvgDepDelayMinutes,
        P90DepDelayMinutes,
        DepDel15Pct,
        FirstQualifyingMonth,
        LastQualifyingMonth,
        HotspotRank
    FROM top_hotspots

    UNION ALL

    SELECT
        'monthly_trend' AS RowType,
        CAST(q.MonthStart AS Nullable(Date)) AS MonthStart,
        t.Dest,
        t.DepTimeBlk,
        t.QualifyingMonths,
        count() AS CompletedFlights,
        round(avg(q.DepDelayMinutes), 2) AS AvgDepDelayMinutes,
        round(quantile(0.9)(q.DepDelayMinutes), 2) AS P90DepDelayMinutes,
        round(avg(q.DepDel15) * 100, 2) AS DepDel15Pct,
        t.FirstQualifyingMonth,
        t.LastQualifyingMonth,
        t.HotspotRank
    FROM qualifying_raw_flights q
    INNER JOIN top_hotspots t
        ON q.Dest = t.Dest
       AND q.DepTimeBlk = t.DepTimeBlk
    GROUP BY
        q.MonthStart,
        t.Dest,
        t.DepTimeBlk,
        t.QualifyingMonths,
        t.FirstQualifyingMonth,
        t.LastQualifyingMonth,
        t.HotspotRank
)
SELECT
    RowType,
    MonthStart,
    Dest,
    DepTimeBlk,
    QualifyingMonths,
    CompletedFlights,
    AvgDepDelayMinutes,
    P90DepDelayMinutes,
    DepDel15Pct,
    FirstQualifyingMonth,
    LastQualifyingMonth,
    HotspotRank
FROM combined_results
ORDER BY
    RowType ASC,
    HotspotRank ASC,
    MonthStart ASC,
    Dest ASC,
    DepTimeBlk ASC
```

Visual context:

- Question title: `Delta ATL departure delay hotspots by destination and time block`
- Visual type: `html_heatmap`

Question-specific visual guidance:

Build a dynamic non-map dashboard for Delta ATL departure delay hotspots.

Layout intent:

- headline and subtitle that clearly frame Delta departures from ATL
- KPI strip for worst hotspot, worst average departure delay, p90 departure delay, and qualifying months
- primary heatmap with `Dest` on one axis and `DepTimeBlk` on the other, colored by average `DepDelayMinutes`
- supporting monthly trend view for the top 3 hotspot cells
- ranked table for the top 20 hotspot cells
- visible legend for the heatmap color scale

Behavior:

- derive the top 3 hotspot cells for the trend view from the fetched ranking rows, not from hardcoded labels
- make the single worst hotspot cell visually prominent
- keep the heatmap readable on mobile by allowing horizontal scrolling if necessary
- preserve useful content when the fetched result set is empty