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
        ) AS Route
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline
    ORDER BY hop_count DESC, FlightDate DESC
    LIMIT 10
)
SELECT
    FlightNum,
    Carrier,
    Route,
    count()            AS occurrences,
    min(FlightDate)    AS first_seen,
    max(FlightDate)    AS last_seen
FROM ranked
GROUP BY FlightNum, Carrier, Route
ORDER BY occurrences DESC
