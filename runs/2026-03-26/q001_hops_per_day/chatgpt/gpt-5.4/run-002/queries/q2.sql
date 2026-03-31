WITH recent_legs AS (
    SELECT
        FlightDate,
        Carrier,
        FlightNum,
        Tail_Number AS aircraft_id,
        OriginAirportID,
        DestAirportID,
        OriginCode,
        DestCode,
        coalesce(DepTime, CRSDepTime) AS dep_hhmm,
        coalesce(ArrTime, CRSArrTime) AS arr_hhmm,
        min(Distance) AS Distance
    FROM ontime.fact_ontime
    WHERE FlightDate >= addYears(today(), -5)
      AND Cancelled = 0 AND Carrier != '' AND FlightNum != ''
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, OriginAirportID, DestAirportID, OriginCode, DestCode, dep_hhmm, arr_hhmm
),
recent_itineraries AS (
    SELECT
        FlightDate, Carrier, FlightNum, aircraft_id,
        count() AS hop_count,
        arraySort(groupArray((ifNull(dep_hhmm, 9999), ifNull(arr_hhmm, 9999), OriginCode, DestCode, OriginAirportID, DestAirportID, ifNull(Distance, 0)))) AS legs_sorted
    FROM recent_legs
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id
),
recent_routes AS (
    SELECT
        FlightDate, Carrier, FlightNum, aircraft_id, hop_count, legs_sorted,
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
    SELECT aircraft_id, FlightDate, Carrier, FlightNum, hop_count, Route, legs_sorted, most_recent_flight_date
    FROM top_routes
    WHERE rn = 1
    ORDER BY hop_count DESC, most_recent_flight_date DESC, Route
    LIMIT 10
),
route_airports AS (
    SELECT
        Route,
        arrayJoin(arrayConcat([tupleElement(legs_sorted[1], 5)], arrayMap(x -> tupleElement(x, 6), legs_sorted))) AS airport_id
    FROM top10
)
SELECT
    t.Route,
    if(t.aircraft_id = '', 'unknown', t.aircraft_id) AS aircraft_id,
    t.FlightNum AS flight_number,
    t.Carrier AS carrier,
    t.FlightDate AS flight_date,
    t.hop_count,
    arraySum(arrayMap(x -> toUInt64(tupleElement(x, 7)), t.legs_sorted)) AS total_flown_distance,
    uniqExact(ra.airport_id) AS unique_airports,
    uniqExact(d.CityMarketID) AS unique_city_markets,
    uniqExact(d.StateCode) AS unique_states,
    uniqExact(d.UtcLocalTimeVariation) AS unique_local_time_offsets,
    min(d.CountryCodeISO = 'US') AS entirely_domestic
FROM top10 t
LEFT JOIN route_airports ra ON t.Route = ra.Route
LEFT JOIN ontime.dim_airports d ON ra.airport_id = d.AirportID AND d.IsLatest = 1
GROUP BY t.Route, aircraft_id, flight_number, carrier, flight_date, t.hop_count, t.legs_sorted
ORDER BY total_flown_distance DESC, unique_states DESC, t.Route
