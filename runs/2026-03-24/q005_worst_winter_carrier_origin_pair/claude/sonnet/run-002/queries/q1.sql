SELECT
    IATA_CODE_Reporting_Airline AS carrier,
    OriginCode AS origin,
    count() AS winter_flights,
    round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct,
    round(avg(ifNull(DepDelayMinutes, 0)), 2) AS avg_dep_delay_min
FROM ontime.fact_ontime
WHERE toMonth(FlightDate) IN (12, 1, 2)
  AND Cancelled = 0
GROUP BY carrier, origin
HAVING winter_flights >= 1000
ORDER BY otp_pct ASC
LIMIT 20
