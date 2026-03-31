SELECT
    OriginCode,
    count() AS flights,
    round(100.0 * count() / sum(count()) OVER (), 2) AS pct_of_total_flights,
    round(avg(DepDelay), 2) AS avg_dep_delay_min,
    round(100.0 * sum(DepDel15) / count(), 2) AS pct_dep_del15
FROM ontime.fact_ontime
WHERE Reporting_Airline = 'AA'
  AND Cancelled = 0
  AND DepDelay IS NOT NULL
  AND toStartOfMonth(FlightDate) = '2024-07-01'
GROUP BY OriginCode
ORDER BY flights DESC
LIMIT 20
