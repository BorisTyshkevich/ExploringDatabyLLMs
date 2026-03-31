SELECT
    DestCode,
    DepTimeBlk,
    toStartOfMonth(FlightDate) AS flight_month,
    count() AS monthly_flights,
    round(avg(DepDelay), 1) AS avg_dep_delay_min,
    round(avg(DepDel15) * 100, 1) AS pct_dep_del15
FROM ontime.fact_ontime
WHERE
    IATA_CODE_Reporting_Airline = 'DL'
    AND OriginCode = 'ATL'
    AND DestCode = 'EWR'
    AND DepTimeBlk = '1600-1659'
    AND Cancelled = 0
    AND DepDelay IS NOT NULL
GROUP BY DestCode, DepTimeBlk, flight_month
HAVING monthly_flights >= 30
ORDER BY flight_month
