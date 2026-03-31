WITH itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline AS Carrier,
        FlightDate,
        count() AS Hops,
        groupArray((toUInt32(coalesce(DepTime, CRSDepTime, 0)), Origin, Dest)) AS legs_raw
    FROM ontime.ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
),
itineraries_sorted AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        Carrier,
        FlightDate,
        Hops,
        arraySort(x -> x.1, legs_raw) AS legs_sorted
    FROM itineraries
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    Carrier,
    FlightDate,
    Hops,
    arrayStringConcat(
        arrayConcat(
            arrayMap(x -> x.2, legs_sorted),
            [arrayElement(arrayMap(x -> x.3, legs_sorted), -1)]
        ),
        ' -> '
    ) AS Route,
    arrayStringConcat(
        arrayMap(
            x -> concat(x.2, ' ', printf('%02d:%02d', intDiv(x.1, 100), x.1 % 100)),
            legs_sorted
        ),
        ', '
    ) AS DepartureSchedule
FROM itineraries_sorted
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
