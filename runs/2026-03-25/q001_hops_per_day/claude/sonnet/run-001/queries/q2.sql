WITH daily_itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        Reporting_Airline               AS Carrier,
        FlightDate,
        count()                         AS hop_count,
        arrayStringConcat(
            arrayConcat(
                [arrayElement(arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))), 1).2],
                arrayMap(x -> x.3, arraySort(x -> x.1, groupArray((assumeNotNull(CRSDepTime), OriginCode, DestCode))))
            ),
            '-'
        ) AS Route
    FROM ontime.fact_ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Reporting_Airline, FlightDate
    HAVING hop_count >= 2
),
unique_routes AS (
    SELECT
        Route,
        max(hop_count)                  AS hop_count,
        argMax(Tail_Number, FlightDate) AS Tail_Number,
        argMax(FlightNum,   FlightDate) AS FlightNum,
        argMax(Carrier,     FlightDate) AS Carrier,
        max(FlightDate)                 AS most_recent_date,
        count()                         AS occurrences
    FROM daily_itineraries
    GROUP BY Route
)
SELECT Tail_Number, FlightNum, Carrier, most_recent_date AS FlightDate, hop_count, Route, occurrences
FROM unique_routes
ORDER BY most_recent_date DESC, hop_count DESC
LIMIT 1
