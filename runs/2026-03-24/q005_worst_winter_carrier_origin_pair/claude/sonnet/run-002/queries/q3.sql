WITH worst30 AS (
    SELECT
        IATA_CODE_Reporting_Airline AS carrier,
        OriginCode AS origin,
        count() AS winter_flights,
        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct
    FROM ontime.fact_ontime
    WHERE toMonth(FlightDate) IN (12, 1, 2)
      AND Cancelled = 0
    GROUP BY carrier, origin
    HAVING winter_flights >= 1000
    ORDER BY otp_pct ASC
    LIMIT 30
)
SELECT
    carrier,
    count() AS pairs_in_worst30,
    round(avg(otp_pct), 2) AS avg_otp_pct,
    groupArray(origin) AS origins
FROM worst30
GROUP BY carrier
ORDER BY pairs_in_worst30 DESC, avg_otp_pct ASC
