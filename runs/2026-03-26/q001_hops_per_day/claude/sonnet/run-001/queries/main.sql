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
    argMax(Tail_Number, FlightDate) AS aircraft_id,
    argMax(flight_num, FlightDate) AS flight_number,
    argMax(carrier, FlightDate) AS carrier,
    max(FlightDate) AS most_recent_date,
    any(hop_count) AS hop_count,
    countDistinct(FlightDate) AS recurrence_count,
    Route
FROM top_itin
GROUP BY Route
ORDER BY most_recent_date DESC
LIMIT 10
