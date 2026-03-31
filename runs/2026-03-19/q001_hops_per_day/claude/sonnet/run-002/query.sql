WITH itineraries AS (
    SELECT
        FlightDate,
        Tail_Number,
        IATA_CODE_Reporting_Airline AS Carrier,
        Flight_Number_Reporting_Airline AS FlightNum,
        count() AS Hops,
        arrayStringConcat(
            arrayConcat(
                arrayMap(
                    x -> x.2,
                    arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Origin)))
                ),
                [arrayElement(
                    arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Dest)))),
                    -1
                )]
            ),
            ' -> '
        ) AS Route,
        arrayStringConcat(
            arrayMap(
                x -> x.2 || '@' || toString(x.1),
                arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Origin)))
            ),
            ', '
        ) AS DepartureTimes
    FROM ontime.ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND DepTime IS NOT NULL
    GROUP BY FlightDate, Tail_Number, Carrier, FlightNum
)
SELECT
    FlightDate,
    Tail_Number,
    Carrier,
    FlightNum,
    Hops,
    Route,
    DepartureTimes
FROM itineraries
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
