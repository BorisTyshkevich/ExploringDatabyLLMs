WITH leg_data AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline AS Carrier,
        FlightDate,
        arraySort(
            x -> x.1,
            groupArray((
                coalesce(DepTime, CRSDepTime, 0),
                Origin,
                Dest
            ))
        ) AS sorted_legs
    FROM ontime.ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Carrier, FlightDate
    HAVING length(sorted_legs) >= 2
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    Carrier,
    FlightDate,
    length(sorted_legs) AS Hops,
    arrayStringConcat(
        arrayPushBack(
            arrayMap(x -> x.2, sorted_legs),
            arrayElement(sorted_legs, -1).3
        ),
        ' → '
    ) AS Route,
    arrayStringConcat(
        arrayMap(x -> concat(leftPad(toString(x.1), 4, '0'), ' (', x.2, ')'), sorted_legs),
        ', '
    ) AS DepartureTimes
FROM leg_data
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
