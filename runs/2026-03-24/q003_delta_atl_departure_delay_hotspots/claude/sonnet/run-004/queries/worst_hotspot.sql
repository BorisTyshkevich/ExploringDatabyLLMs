SELECT
    h.DestCode,
    h.DepTimeBlk,
    h.total_flights,
    h.avg_dep_delay,
    h.p90_dep_delay,
    h.pct_del15,
    m.credible_months
FROM (
    SELECT
        DestCode,
        DepTimeBlk,
        count() AS total_flights,
        round(avg(DepDelay), 1) AS avg_dep_delay,
        round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,
        round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15
    FROM ontime.fact_ontime
    WHERE
        OriginCode = 'ATL'
        AND IATA_CODE_Reporting_Airline = 'DL'
        AND Cancelled = 0
        AND DepDelay IS NOT NULL
    GROUP BY DestCode, DepTimeBlk
    HAVING total_flights >= 500
) h
JOIN (
    SELECT
        DestCode,
        DepTimeBlk,
        countIf(monthly_cnt >= 20) AS credible_months
    FROM (
        SELECT
            DestCode,
            DepTimeBlk,
            toStartOfMonth(FlightDate) AS ym,
            count() AS monthly_cnt
        FROM ontime.fact_ontime
        WHERE
            OriginCode = 'ATL'
            AND IATA_CODE_Reporting_Airline = 'DL'
            AND Cancelled = 0
            AND DepDelay IS NOT NULL
        GROUP BY DestCode, DepTimeBlk, ym
    )
    GROUP BY DestCode, DepTimeBlk
) m ON h.DestCode = m.DestCode AND h.DepTimeBlk = m.DepTimeBlk
ORDER BY h.avg_dep_delay DESC
LIMIT 15
