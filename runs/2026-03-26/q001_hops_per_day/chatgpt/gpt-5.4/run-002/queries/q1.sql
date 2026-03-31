WITH recent_legs AS (
    SELECT FlightDate, Carrier, FlightNum, Tail_Number AS aircraft_id, OriginCode, DestCode,
           coalesce(DepTime, CRSDepTime) AS dep_hhmm, coalesce(ArrTime, CRSArrTime) AS arr_hhmm
    FROM ontime.fact_ontime
    WHERE FlightDate >= addYears(today(), -5)
      AND Cancelled = 0 AND Carrier != '' AND FlightNum != ''
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, OriginCode, DestCode, dep_hhmm, arr_hhmm
),
recent_itineraries AS (
    SELECT FlightDate, Carrier, FlightNum, aircraft_id, count() AS hop_count,
           arraySort(groupArray((ifNull(dep_hhmm, 9999), ifNull(arr_hhmm, 9999), OriginCode, DestCode))) AS legs_sorted
    FROM recent_legs
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id
),
recent_routes AS (
    SELECT FlightDate, Carrier, FlightNum, aircraft_id, hop_count,
           arrayStringConcat(arrayConcat([tupleElement(legs_sorted[1], 3)], arrayMap(x -> tupleElement(x, 4), legs_sorted)), '-') AS Route
    FROM recent_itineraries
    WHERE hop_count >= 2
),
top_routes AS (
    SELECT *, max(FlightDate) OVER (PARTITION BY Route) AS most_recent_flight_date,
           row_number() OVER (PARTITION BY Route ORDER BY hop_count DESC, FlightDate DESC, Carrier, FlightNum, aircraft_id) AS rn
    FROM recent_routes
),
top10 AS (
    SELECT Route
    FROM top_routes
    WHERE rn = 1
    ORDER BY hop_count DESC, most_recent_flight_date DESC, Route
    LIMIT 10
),
airports AS (
    SELECT
        Route,
        airport_code,
        pos,
        length(airport_list) AS route_len
    FROM (
        SELECT Route, splitByChar('-', Route) AS airport_list
        FROM top10
    )
    ARRAY JOIN airport_list AS airport_code, arrayEnumerate(airport_list) AS pos
)
SELECT
    a.airport_code,
    any(d.DisplayAirportName) AS airport_name,
    any(d.CityName) AS city_name,
    any(d.StateCode) AS state_code,
    count() AS total_appearances,
    countIf(pos = 1) AS origin_appearances,
    countIf(pos > 1 AND pos < route_len) AS intermediate_stop_appearances,
    countIf(pos = route_len) AS final_destination_appearances,
    round(countDistinct(Route) / 10.0, 3) AS share_of_itineraries
FROM airports a
LEFT JOIN ontime.dim_airports d
    ON a.airport_code = d.AirportCode AND d.IsLatest = 1
GROUP BY a.airport_code
ORDER BY total_appearances DESC, intermediate_stop_appearances DESC, a.airport_code
