SELECT
    DestCode,
    DepTimeBlk,
    count() AS flight_count,
    round(avg(DepDelay), 1) AS avg_dep_delay_min,
    round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay_min,
    round(avg(DepDel15) * 100, 1) AS pct_dep_del15,
    countDistinct(toStartOfMonth(FlightDate)) AS active_months
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
