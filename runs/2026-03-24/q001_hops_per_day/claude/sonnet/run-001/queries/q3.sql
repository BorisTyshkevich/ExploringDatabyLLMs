WITH itineraries AS (
    SELECT
        FlightDate,
        Flight_Number_Reporting_Airline AS FlightNum,
        IATA_CODE_Reporting_Airline AS Carrier,
        count() AS HopCount,
        arrayStringConcat(
            arrayMap(x -> x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), OriginCode))))
            || [arrayElement(arrayMap(x -> x.2, arraySort(groupArray((coalesce(DepTime, CRSDepTime, 0), DestCode)))), -1)],
            '-'
        ) AS Route
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
        AND Tail_Number != ''
        AND Flight_Number_Reporting_Airline != ''
    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline
    HAVING HopCount = 8
)
SELECT FlightNum, Carrier, Route, count() AS DateCount,
       min(FlightDate) AS FirstSeen, max(FlightDate) AS LastSeen
FROM itineraries
GROUP BY FlightNum, Carrier, Route
ORDER BY DateCount DESC
LIMIT 10
