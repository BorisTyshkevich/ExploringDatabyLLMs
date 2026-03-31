WITH sorted AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arraySort(x -> x.1, groupArray(tuple(DepTime, Origin, Dest))) AS legs
    FROM ontime.ontime
    WHERE Tail_Number != '' AND DepTime IS NOT NULL
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    HAVING Hops >= 6
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    Hops,
    arrayStringConcat(arrayMap(x -> x.2, legs), ' -> ')
        || ' -> ' || legs[length(legs)].3 AS Route,
    arrayStringConcat(
        arrayMap(
            x -> concat(x.2, ' ', lpad(toString(x.1 DIV 100), 2, '0'), ':', lpad(toString(x.1 MOD 100), 2, '0')),
            legs
        ), ', '
    ) AS DepSchedule
FROM sorted
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
