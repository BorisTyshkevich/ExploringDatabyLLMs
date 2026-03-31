WITH winter_pairs AS (
    SELECT
        IATA_CODE_Reporting_Airline AS carrier,
        OriginCode AS origin,
        count() AS winter_flights,
        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime,
        round(avg(DepDelayMinutes), 2) AS avg_dep_delay_min,
        round(sum(assumeNotNull(WeatherDelay)) / nullIf(sum(assumeNotNull(WeatherDelay) + assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)), 0) * 100, 2) AS weather_pct,
        round(sum(assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)) / nullIf(sum(assumeNotNull(WeatherDelay) + assumeNotNull(CarrierDelay) + assumeNotNull(NASDelay) + assumeNotNull(SecurityDelay) + assumeNotNull(LateAircraftDelay)), 0) * 100, 2) AS operational_pct
    FROM ontime.fact_ontime
    WHERE toMonth(FlightDate) IN (12, 1, 2)
      AND Cancelled = 0
    GROUP BY carrier, origin
    HAVING winter_flights >= 1000
)
SELECT
    carrier,
    origin,
    winter_flights,
    pct_ontime,
    avg_dep_delay_min,
    weather_pct,
    operational_pct,
    row_number() OVER (ORDER BY pct_ontime ASC) AS rank
FROM winter_pairs
ORDER BY pct_ontime ASC
LIMIT 25
