WITH sorted AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS hop_count,
        arraySort(t -> t.3, groupArray(tuple(OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS legs
    FROM ontime.fact_ontime
    WHERE Tail_Number != ''
      AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
)
SELECT
    Tail_Number                    AS `Aircraft ID`,
    Flight_Number_Reporting_Airline AS `Flight Number`,
    IATA_CODE_Reporting_Airline     AS Carrier,
    FlightDate                      AS Date,
    arrayStringConcat(
        arrayConcat(
            [legs[1].1],
            arrayMap(t -> t.2, legs)
        ),
        ' → '
    )                               AS Route,
    hop_count
FROM sorted
ORDER BY hop_count DESC, FlightDate DESC
LIMIT 10
