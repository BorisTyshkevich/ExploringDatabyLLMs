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
    SELECT FlightDate, Tail_Number, carrier, flight_num,
        count() AS hop_count,
        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, arraySort(groupArray(tuple(dep_time, OriginCode)))), [argMax(DestCode, dep_time)]), '-') AS Route
    FROM legs_deduped
    GROUP BY FlightDate, Tail_Number, carrier, flight_num
),
top_itin AS (
    SELECT FlightDate, Tail_Number, carrier, flight_num, hop_count, Route
    FROM itineraries
    WHERE hop_count = (SELECT max(hop_count) FROM itineraries)
),
top10_routes AS (
    SELECT Route
    FROM top_itin
    GROUP BY Route
    ORDER BY max(FlightDate) DESC
    LIMIT 10
),
airport_list AS (
    SELECT DISTINCT arrayJoin(splitByChar('-', Route)) AS airport_code
    FROM top10_routes
)
SELECT
    al.airport_code,
    any(d.CityName) AS city,
    any(d.StateName) AS state,
    any(d.StateCode) AS state_code
FROM airport_list al
LEFT JOIN ontime.dim_airports d ON al.airport_code = toFixedString(d.AirportCode, 3)
GROUP BY al.airport_code
ORDER BY any(d.StateCode), al.airport_code
