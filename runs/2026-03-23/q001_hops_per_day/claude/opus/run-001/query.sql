WITH legs AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        OriginCode,
        DestCode,
        assumeNotNull(CRSDepTime) AS dep
    FROM ontime.fact_ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
),
grouped AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        groupArray((dep, OriginCode, DestCode)) AS legs_arr,
        count() AS hop_count
    FROM legs
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    HAVING hop_count >= 2
)
SELECT
    Tail_Number AS "Aircraft ID",
    Flight_Number_Reporting_Airline AS "Flight Number",
    IATA_CODE_Reporting_Airline AS "Carrier",
    FlightDate AS "Date",
    arrayStringConcat(
        arrayPushBack(
            arrayMap(x -> x.2, arraySort(x -> x.1, legs_arr)),
            arrayElement(arraySort(x -> x.1, legs_arr), -1).3
        ),
        '-'
    ) AS "Route",
    hop_count
FROM grouped
ORDER BY hop_count DESC, FlightDate DESC
LIMIT 10
