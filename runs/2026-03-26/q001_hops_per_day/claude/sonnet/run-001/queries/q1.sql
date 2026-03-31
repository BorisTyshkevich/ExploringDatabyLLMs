WITH
legs_deduped AS (
    SELECT
        FlightDate, Tail_Number,
        IATA_CODE_Reporting_Airline AS carrier,
        Flight_Number_Reporting_Airline AS flight_num,
        OriginCode, DestCode,
        min(ifNull(CRSDepTime, 0)) AS dep_time
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
    GROUP BY FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode
),
itineraries AS (
    SELECT
        FlightDate, Tail_Number, carrier, flight_num,
        count() AS hop_count,
        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route
    FROM legs_deduped
    GROUP BY FlightDate, Tail_Number, carrier, flight_num
),
top_itin AS (
    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route
    FROM itineraries
    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)
)
SELECT
    Route,
    countDistinct(FlightDate) AS recurrence_count,
    multiIf(
        countDistinct(FlightDate) = 1, 'one-off',
        countDistinct(FlightDate) <= 7, 'rare (2-7)',
        countDistinct(FlightDate) <= 20, 'recurring (8-20)',
        'highly recurring (40+)'
    ) AS tier
FROM top_itin
GROUP BY Route
ORDER BY max(FlightDate) DESC
LIMIT 10
