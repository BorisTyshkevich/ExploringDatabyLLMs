WITH ranked AS (
    SELECT
        FlightDate,
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        IATA_CODE_Reporting_Airline AS Carrier,
        count() AS hop_count,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),
                [argMax(DestCode, assumeNotNull(DepTime))]
            ),
            '-'
        ) AS Route,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),
                [argMax(DestCode, assumeNotNull(DepTime))]
            ),
            ','
        ) AS itinerary_sequence
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline
)
SELECT
    Tail_Number,
    FlightNum,
    Carrier,
    FlightDate,
    hop_count,
    Route,
    itinerary_sequence
FROM ranked
ORDER BY hop_count DESC, FlightDate DESC
LIMIT 10
