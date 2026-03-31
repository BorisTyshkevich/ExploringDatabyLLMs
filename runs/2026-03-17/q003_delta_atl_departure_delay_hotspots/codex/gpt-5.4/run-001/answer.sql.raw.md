```sql
WITH
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
qualifying_monthly_cells AS (
    SELECT
        MonthStart,
        Dest,
        DepTimeBlk,
        count() AS CompletedFlights
    FROM base_flights
    GROUP BY
        MonthStart,
        Dest,
        DepTimeBlk
    HAVING count() >= 40
),
qualified_flights AS (
    SELECT
        bf.MonthStart,
        bf.Dest,
        bf.DepTimeBlk,
        bf.DepDelayMinutes,
        bf.DepDel15
    FROM base_flights AS bf
    INNER JOIN qualifying_monthly_cells AS qmc
        ON bf.MonthStart = qmc.MonthStart
       AND bf.Dest = qmc.Dest
       AND bf.DepTimeBlk = qmc.DepTimeBlk
),
hotspot_rollup AS (
    SELECT
        Dest,
        DepTimeBlk,
        count() AS CompletedFlights,
        CAST(round(avg(DepDelayMinutes), 2), 'Nullable(Decimal(18,2))') AS AvgDepDelayMinutes,
        CAST(round(quantile(0.9)(DepDelayMinutes), 2), 'Nullable(Decimal(18,2))') AS P90DepDelayMinutes,
        CAST(round(100.0 * avg(toFloat64(DepDel15)), 2), 'Decimal(18,2)') AS DepDel15Pct,
        countDistinct(MonthStart) AS QualifyingMonths,
        min(MonthStart) AS FirstQualifyingMonth,
        max(MonthStart) AS LastQualifyingMonth
    FROM qualified_flights
    GROUP BY
        Dest,
        DepTimeBlk
    HAVING count() >= 1000
),
ranked_hotspots AS (
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
        row_number() OVER (
            ORDER BY
                AvgDepDelayMinutes DESC,
                P90DepDelayMinutes DESC,
                DepDel15Pct DESC,
                CompletedFlights DESC,
                Dest ASC,
                DepTimeBlk ASC
        ) AS HotspotRank
    FROM hotspot_rollup
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
        HotspotRank
    FROM ranked_hotspots
    WHERE HotspotRank <= 20
),
monthly_trends AS (
    SELECT
        qf.MonthStart,
        qf.Dest,
        qf.DepTimeBlk,
        th.QualifyingMonths,
        count() AS CompletedFlights,
        CAST(round(avg(qf.DepDelayMinutes), 2), 'Nullable(Decimal(18,2))') AS AvgDepDelayMinutes,
        CAST(round(quantile(0.9)(qf.DepDelayMinutes), 2), 'Nullable(Decimal(18,2))') AS P90DepDelayMinutes,
        CAST(round(100.0 * avg(toFloat64(qf.DepDel15)), 2), 'Decimal(18,2)') AS DepDel15Pct,
        th.FirstQualifyingMonth,
        th.LastQualifyingMonth,
        th.HotspotRank
    FROM qualified_flights AS qf
    INNER JOIN top_hotspots AS th
        ON qf.Dest = th.Dest
       AND qf.DepTimeBlk = th.DepTimeBlk
    GROUP BY
        qf.MonthStart,
        qf.Dest,
        qf.DepTimeBlk,
        th.QualifyingMonths,
        th.FirstQualifyingMonth,
        th.LastQualifyingMonth,
        th.HotspotRank
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
FROM (
    SELECT
        'hotspot_summary' AS RowType,
        CAST(NULL, 'Nullable(Date)') AS MonthStart,
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
    FROM monthly_trends
)
ORDER BY
    RowType,
    HotspotRank ASC,
    MonthStart ASC,
    Dest ASC,
    DepTimeBlk ASC
```