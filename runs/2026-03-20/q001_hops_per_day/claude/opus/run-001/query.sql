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
    WHERE Tail_Number != '' AND Cancelled = 0 AND DepTime IS NOT NULL
),
counts AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops
    FROM legs
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    ORDER BY Hops DESC, FlightDate DESC
    LIMIT 10
),
top_legs AS (
    SELECT l.*
    FROM legs l
    INNER JOIN counts c
        ON l.Tail_Number = c.Tail_Number
        AND l.Flight_Number_Reporting_Airline = c.Flight_Number_Reporting_Airline
        AND l.IATA_CODE_Reporting_Airline = c.IATA_CODE_Reporting_Airline
        AND l.FlightDate = c.FlightDate
    ORDER BY l.Tail_Number, l.Flight_Number_Reporting_Airline, l.FlightDate, l.DepTime
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    count() AS Hops,
    arrayStringConcat(
        arrayConcat(
            groupArray(Origin),
            [arrayElement(groupArray(Dest), toUInt32(count()))]
        ),
        '-'
    ) AS Route,
    arrayStringConcat(
        arrayMap(x -> concat(toString(intDiv(x, 100)), ':', leftPad(toString(x % 100), 2, '0')),
            groupArray(DepTime)
        ),
        ', '
    ) AS DepTimes
FROM top_legs
GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
ORDER BY Hops DESC, FlightDate DESC
