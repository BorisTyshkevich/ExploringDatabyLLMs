SELECT
    DestCode,
    DepTimeBlk,
    toStartOfMonth(FlightDate) AS flight_month,
    count() AS flights,
    round(avg(DepDelay), 1) AS avg_dep_delay,
    round(quantile(0.90)(DepDelay), 1) AS p90_dep_delay,
    round(100.0 * countIf(DepDel15 = 1) / count(), 1) AS pct_del15
FROM ontime.fact_ontime
WHERE
    OriginCode = 'ATL'
    AND IATA_CODE_Reporting_Airline = 'DL'
    AND Cancelled = 0
    AND DepDelay IS NOT NULL
    AND DestCode = 'EWR'
    AND DepTimeBlk = '1600-1659'
GROUP BY DestCode, DepTimeBlk, flight_month
HAVING flights >= 20
ORDER BY flight_month
