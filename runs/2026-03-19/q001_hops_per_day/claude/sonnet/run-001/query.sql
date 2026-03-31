WITH leg_data AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS hops,
        groupArray((toUInt32(ifNull(DepTime, ifNull(CRSDepTime, 9999))), Origin, Dest)) AS legs_raw
    FROM ontime.ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
    GROUP BY
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate
),
itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        hops,
        arraySort(x -> x.1, legs_raw) AS legs_sorted
    FROM leg_data
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    hops,
    arrayStringConcat(
        arrayConcat(
            [legs_sorted[1].2],
            arrayMap(x -> x.3, legs_sorted)
        ),
        ' -> '
    ) AS Route,
    arrayStringConcat(
        arrayMap(
            x -> concat(x.2, ' ', lpad(toString(intDiv(x.1, 100)), 2, '0'), ':', lpad(toString(x.1 % 100), 2, '0')),
            legs_sorted
        ),
        ' | '
    ) AS DepartureTimes
FROM itineraries
ORDER BY hops DESC, FlightDate DESC
LIMIT 10
