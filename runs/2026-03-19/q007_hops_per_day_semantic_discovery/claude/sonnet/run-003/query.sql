WITH legs AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Origin,
        Dest,
        assumeNotNull(DepTime) AS DepTime
    FROM ontime.ontime
    WHERE DepTime IS NOT NULL
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
),
daily_itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arrayStringConcat(
            arrayConcat(
                arrayMap(
                    x -> concat(x.2, '(', toString(x.1), ')'),
                    arraySort(x -> x.1, groupArray((DepTime, Origin)))
                ),
                [argMax(Dest, DepTime)]
            ),
            ' -> '
        ) AS Route
    FROM legs
    GROUP BY
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    Hops,
    Route
FROM daily_itineraries
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
