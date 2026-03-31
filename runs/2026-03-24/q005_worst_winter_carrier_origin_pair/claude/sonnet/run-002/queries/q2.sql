SELECT
    IATA_CODE_Reporting_Airline AS carrier,
    OriginCode AS origin,
    count() AS winter_flights,
    round(100.0 * countIf(DepDel15 = 0) / count(), 2) AS otp_pct,
    round(avg(ifNull(DepDelayMinutes, 0)), 2) AS avg_dep_delay_min,
    sum(ifNull(WeatherDelay, 0)) AS weather_min,
    sum(ifNull(CarrierDelay, 0)) AS carrier_min,
    sum(ifNull(NASDelay, 0)) AS nas_min,
    sum(ifNull(SecurityDelay, 0)) AS security_min,
    sum(ifNull(LateAircraftDelay, 0)) AS late_aircraft_min,
    sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)
        + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)) AS total_reported_delay,
    round(sum(ifNull(WeatherDelay, 0)) * 100.0
          / nullIf(sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)
                       + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)), 0), 1) AS weather_share_pct,
    round(sum(ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0) + ifNull(SecurityDelay, 0)
              + ifNull(LateAircraftDelay, 0)) * 100.0
          / nullIf(sum(ifNull(WeatherDelay, 0) + ifNull(CarrierDelay, 0) + ifNull(NASDelay, 0)
                       + ifNull(SecurityDelay, 0) + ifNull(LateAircraftDelay, 0)), 0), 1) AS operational_share_pct
FROM ontime.fact_ontime
WHERE toMonth(FlightDate) IN (12, 1, 2)
  AND Cancelled = 0
  AND (IATA_CODE_Reporting_Airline, OriginCode) IN (
      ('DH','ORD'), ('PI','DFW'), ('PI','LAX'), ('PI','DAY'),
      ('AS','DUT'), ('YV','ORD'), ('OO','OTH'), ('AA','EYW'),
      ('F9','PBI'), ('OH','RSW')
  )
GROUP BY carrier, origin
ORDER BY otp_pct ASC
