The query produces 20 hotspot summary rows and 812 monthly trend rows. The query is verified and working correctly.

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
