SELECT
    f.OriginCode,
    f.DestCode,
    f.DepTime,
    o.DisplayAirportName  AS OriginName,
    o.Latitude            AS OriginLat,
    o.Longitude           AS OriginLon,
    d.DisplayAirportName  AS DestName,
    d.Latitude            AS DestLat,
    d.Longitude           AS DestLon
FROM ontime.fact_ontime f
LEFT JOIN ontime.dim_airports o ON f.OriginCode = o.AirportCode
LEFT JOIN ontime.dim_airports d ON f.DestCode   = d.AirportCode
WHERE f.Cancelled = 0
  AND f.Tail_Number = 'N957WN'
  AND f.Flight_Number_Reporting_Airline = '366'
  AND f.FlightDate = '2024-12-01'
ORDER BY assumeNotNull(f.DepTime)
