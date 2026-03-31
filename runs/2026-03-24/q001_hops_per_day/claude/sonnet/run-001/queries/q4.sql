WITH itineraries AS (
    SELECT
        Flight_Number_Reporting_Airline AS FlightNum,
        IATA_CODE_Reporting_Airline AS Carrier,
        count() AS HopCount,
        arrayElement(arrayMap(x -> x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode)))), 1) AS StartAirport,
        arrayElement(arrayMap(x -> x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1) AS EndAirport,
        arrayStringConcat(
            arrayMap(x -> x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))
            || [arrayElement(arrayMap(x -> x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],
            '-'
        ) AS Route
    FROM ontime.fact_ontime
    WHERE Cancelled = 0 AND Tail_Number != '' AND Flight_Number_Reporting_Airline != ''
    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline
    HAVING HopCount = 8
),
route_summary AS (
    SELECT FlightNum, Carrier, Route, StartAirport, EndAirport, count() AS DateCount
    FROM itineraries
    GROUP BY FlightNum, Carrier, Route, StartAirport, EndAirport
    ORDER BY DateCount DESC
    LIMIT 10
)
SELECT
    rs.FlightNum, rs.Carrier, rs.Route, rs.DateCount,
    a_start.DisplayAirportName AS StartName, a_start.Latitude AS StartLat, a_start.Longitude AS StartLon,
    a_end.DisplayAirportName AS EndName, a_end.Latitude AS EndLat, a_end.Longitude AS EndLon
FROM route_summary rs
LEFT JOIN ontime.dim_airports a_start ON a_start.AirportCode = rs.StartAirport
LEFT JOIN ontime.dim_airports a_end ON a_end.AirportCode = rs.EndAirport
ORDER BY rs.DateCount DESC
