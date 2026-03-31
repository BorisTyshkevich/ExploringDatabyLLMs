WITH legs AS (
    SELECT
        Tail_Number,
        IATA_CODE_Reporting_Airline AS Carrier,
        Flight_Number_Reporting_Airline AS FlightNum,
        FlightDate,
        Origin,
        Dest,
        DepTime
    FROM ontime.ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND DepTime IS NOT NULL
),
grouped AS (
    SELECT
        Tail_Number,
        Carrier,
        FlightNum,
        FlightDate,
        count() AS Hops,
        arraySort(x -> tupleElement(x, 1), groupArray((DepTime, Origin, Dest))) AS sorted_legs
    FROM legs
    GROUP BY Tail_Number, Carrier, FlightNum, FlightDate
)
SELECT
    Tail_Number,
    Carrier,
    FlightNum,
    FlightDate,
    Hops,
    arrayStringConcat(
        arrayConcat(
            arrayMap(x -> tupleElement(x, 2), sorted_legs),
            [tupleElement(sorted_legs[length(sorted_legs)], 3)]
        ),
        ' -> '
    ) AS Route,
    arrayStringConcat(
        arrayMap(x -> concat(
            leftPad(toString(intDiv(tupleElement(x, 1), 100)), 2, '0'),
            ':',
            leftPad(toString(tupleElement(x, 1) % 100), 2, '0')
        ), sorted_legs),
        ', '
    ) AS DepartureTimes
FROM grouped
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
