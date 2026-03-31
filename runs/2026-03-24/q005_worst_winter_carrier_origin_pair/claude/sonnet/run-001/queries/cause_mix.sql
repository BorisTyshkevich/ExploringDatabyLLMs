WITH worst_pairs AS (
    SELECT
        IATA_CODE_Reporting_Airline AS carrier,
        OriginCode AS origin,
        count() AS winter_flights,
        round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS pct_ontime,
        round(avg(DepDelayMinutes), 2) AS avg_dep_delay_min,
        round(sum(assumeNotNull(CarrierDelay)), 0) AS total_carrier_delay,
        round(sum(assumeNotNull(WeatherDelay)), 0) AS total_weather_delay,
        round(sum(assumeNotNull(NASDelay)), 0) AS total_nas_delay,
        round(sum(assumeNotNull(SecurityDelay)), 0) AS total_security_delay,
        round(sum(assumeNotNull(LateAircraftDelay)), 0) AS total_late_aircraft_delay
    FROM ontime.fact_ontime
    WHERE toMonth(FlightDate) IN (12, 1, 2)
      AND Cancelled = 0
    GROUP BY carrier, origin
    HAVING winter_flights >= 1000
    ORDER BY pct_ontime ASC
    LIMIT 25
)
SELECT
    carrier,
    origin,
    pct_ontime,
    total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay AS total_delay_min,
    round(100.0 * total_weather_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS weather_share_pct,
    round(100.0 * total_carrier_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS carrier_share_pct,
    round(100.0 * total_nas_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS nas_share_pct,
    round(100.0 * total_late_aircraft_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS late_aircraft_share_pct,
    round(100.0 * total_security_delay / nullIf(total_carrier_delay + total_weather_delay + total_nas_delay + total_security_delay + total_late_aircraft_delay, 0), 2) AS security_share_pct
FROM worst_pairs
ORDER BY pct_ontime ASC
