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
-- Base filtered flights
base_flights AS (
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

-- Monthly aggregation with qualification filter (>=40 flights)
monthly_qualified AS (
    SELECT
        MonthStart,
        Dest,
        DepTimeBlk,
        count() AS MonthlyFlights
    FROM base_flights
    GROUP BY MonthStart, Dest, DepTimeBlk
    HAVING MonthlyFlights >= 40
),

-- Hotspots with >=1000 total flights across qualifying months
hotspot_qualified AS (
    SELECT
        Dest,
        DepTimeBlk,
        count() AS QualifyingMonths,
        sum(MonthlyFlights) AS TotalQualifyingFlights,
        min(MonthStart) AS FirstQualifyingMonth,
        max(MonthStart) AS LastQualifyingMonth
    FROM monthly_qualified
    GROUP BY Dest, DepTimeBlk
    HAVING TotalQualifyingFlights >= 1000
),

-- Recompute hotspot metrics from raw flights belonging to qualifying monthly cells
hotspot_metrics AS (
    SELECT
        bf.Dest AS Dest,
        bf.DepTimeBlk AS DepTimeBlk,
        hq.QualifyingMonths AS QualifyingMonths,
        hq.FirstQualifyingMonth AS FirstQualifyingMonth,
        hq.LastQualifyingMonth AS LastQualifyingMonth,
        count() AS CompletedFlights,
        avg(bf.DepDelayMinutes) AS AvgDepDelayMinutes,
        quantile(0.9)(bf.DepDelayMinutes) AS P90DepDelayMinutes,
        100.0 * sum(bf.DepDel15) / count() AS DepDel15Pct
    FROM base_flights bf
    INNER JOIN monthly_qualified mq
        ON bf.MonthStart = mq.MonthStart
       AND bf.Dest = mq.Dest
       AND bf.DepTimeBlk = mq.DepTimeBlk
    INNER JOIN hotspot_qualified hq
        ON bf.Dest = hq.Dest
       AND bf.DepTimeBlk = hq.DepTimeBlk
    GROUP BY bf.Dest, bf.DepTimeBlk, hq.QualifyingMonths, hq.FirstQualifyingMonth, hq.LastQualifyingMonth
),

-- Rank hotspots
hotspot_ranked AS (
    SELECT
        Dest,
        DepTimeBlk,
        QualifyingMonths,
        FirstQualifyingMonth,
        LastQualifyingMonth,
        CompletedFlights,
        AvgDepDelayMinutes,
        P90DepDelayMinutes,
        DepDel15Pct,
        row_number() OVER (
            ORDER BY AvgDepDelayMinutes DESC, P90DepDelayMinutes DESC, DepDel15Pct DESC, CompletedFlights DESC
        ) AS HotspotRank
    FROM hotspot_metrics
),

-- Top 20 hotspots
top20_hotspots AS (
    SELECT * FROM hotspot_ranked WHERE HotspotRank <= 20
),

-- Monthly trend metrics for top 20 hotspots (recomputed from raw flights)
monthly_trend_metrics AS (
    SELECT
        bf.MonthStart AS MonthStart,
        bf.Dest AS Dest,
        bf.DepTimeBlk AS DepTimeBlk,
        t20.HotspotRank AS HotspotRank,
        t20.FirstQualifyingMonth AS FirstQualifyingMonth,
        t20.LastQualifyingMonth AS LastQualifyingMonth,
        count() AS CompletedFlights,
        avg(bf.DepDelayMinutes) AS AvgDepDelayMinutes,
        quantile(0.9)(bf.DepDelayMinutes) AS P90DepDelayMinutes,
        100.0 * sum(bf.DepDel15) / count() AS DepDel15Pct
    FROM base_flights bf
    INNER JOIN monthly_qualified mq
        ON bf.MonthStart = mq.MonthStart
       AND bf.Dest = mq.Dest
       AND bf.DepTimeBlk = mq.DepTimeBlk
    INNER JOIN top20_hotspots t20
        ON bf.Dest = t20.Dest
       AND bf.DepTimeBlk = t20.DepTimeBlk
    GROUP BY bf.MonthStart, bf.Dest, bf.DepTimeBlk, t20.HotspotRank, t20.FirstQualifyingMonth, t20.LastQualifyingMonth
),

-- Hotspot summary rows
hotspot_summary_rows AS (
    SELECT
        'hotspot_summary' AS RowType,
        CAST(NULL AS Nullable(Date)) AS MonthStart,
        Dest,
        DepTimeBlk,
        QualifyingMonths,
        CompletedFlights,
        round(AvgDepDelayMinutes, 2) AS AvgDepDelayMinutes,
        round(P90DepDelayMinutes, 2) AS P90DepDelayMinutes,
        round(DepDel15Pct, 2) AS DepDel15Pct,
        FirstQualifyingMonth,
        LastQualifyingMonth,
        HotspotRank
    FROM top20_hotspots
),

-- Monthly trend rows
monthly_trend_rows AS (
    SELECT
        'monthly_trend' AS RowType,
        CAST(MonthStart AS Nullable(Date)) AS MonthStart,
        Dest,
        DepTimeBlk,
        CAST(0 AS UInt64) AS QualifyingMonths,
        CompletedFlights,
        round(AvgDepDelayMinutes, 2) AS AvgDepDelayMinutes,
        round(P90DepDelayMinutes, 2) AS P90DepDelayMinutes,
        round(DepDel15Pct, 2) AS DepDel15Pct,
        FirstQualifyingMonth,
        LastQualifyingMonth,
        HotspotRank
    FROM monthly_trend_metrics
)

-- Final union
SELECT * FROM hotspot_summary_rows
UNION ALL
SELECT * FROM monthly_trend_rows
ORDER BY RowType, HotspotRank ASC, MonthStart ASC, Dest, DepTimeBlk
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