SELECT
    DestCode,
    count() AS hotspot_slots,
    sum(flight_count) AS total_flights,
    round(avg(avg_dep_delay_min), 1) AS mean_avg_delay,
    round(avg(pct_dep_del15), 1) AS mean_pct_del15
FROM (
    SELECT
        DestCode,
        DepTimeBlk,
        count() AS flight_count,
        round(avg(DepDelay), 1) AS avg_dep_delay_min,
        round(avg(DepDel15) * 100, 1) AS pct_dep_del15
    FROM ontime.fact_ontime
    WHERE
        IATA_CODE_Reporting_Airline = 'DL'
        AND OriginCode = 'ATL'
        AND Cancelled = 0
        AND DepDelay IS NOT NULL
    GROUP BY DestCode, DepTimeBlk
    HAVING flight_count >= 500
    ORDER BY avg_dep_delay_min DESC
    LIMIT 15
)
GROUP BY DestCode
ORDER BY mean_avg_delay DESC
