SELECT
    toStartOfMonth(FlightDate) AS flight_month,
    count() AS total_flights,
    round(avg(DepDelay), 2) AS avg_dep_delay_min,
    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15
FROM ontime.fact_ontime
WHERE Reporting_Airline = 'AA'
  AND Cancelled = 0
  AND DepDelay IS NOT NULL
GROUP BY flight_month
ORDER BY flight_month
