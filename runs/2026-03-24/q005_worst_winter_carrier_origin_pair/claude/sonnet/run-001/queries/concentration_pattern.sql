WITH worst_pairs AS (
    SELECT
        IATA_CODE_Reporting_Airline AS carrier,
        OriginCode AS origin,
        count() AS winter_flights,
        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime
    FROM ontime.fact_ontime
    WHERE toMonth(FlightDate) IN (12, 1, 2)
      AND Cancelled = 0
    GROUP BY carrier, origin
    HAVING winter_flights >= 1000
    ORDER BY pct_ontime ASC
    LIMIT 25
)
SELECT
    uniq(carrier) AS distinct_carriers,
    uniq(origin) AS distinct_airports,
    count() AS total_pairs,
    groupArray(carrier) AS carriers,
    groupArray(origin) AS airports
FROM worst_pairs
