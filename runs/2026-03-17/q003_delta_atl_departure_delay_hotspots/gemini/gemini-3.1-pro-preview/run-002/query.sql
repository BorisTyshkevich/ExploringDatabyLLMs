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
