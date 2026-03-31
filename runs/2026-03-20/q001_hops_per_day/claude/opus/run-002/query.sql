WITH legs AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Origin,
        Dest,
        DepTime
    FROM ontime.ontime
    WHERE Tail_Number != ''
      AND Cancelled = 0
      AND DepTime IS NOT NULL
),
counted AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.2,
                    arraySort(x -> x.1,
                        groupArray( (DepTime, Origin) )
                    )
                ),
                [ argMax(Dest, DepTime) ]
            ),
            '-'
        ) AS Route,
        arrayStringConcat(
            arrayMap(x -> lpad(toString(x.1), 4, '0'),
                arraySort(x -> x.1,
                    groupArray( (DepTime, Origin) )
                )
            ),
            ', '
        ) AS DepTimes
    FROM legs
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    HAVING Hops >= 2
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline AS Carrier,
    FlightDate,
    Hops,
    Route,
    DepTimes
FROM counted
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
