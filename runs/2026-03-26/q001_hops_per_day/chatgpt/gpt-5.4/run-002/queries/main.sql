WITH recent_legs AS (
    SELECT
        FlightDate,
        Carrier,
        FlightNum,
        Tail_Number AS aircraft_id,
        OriginCode,
        DestCode,
        coalesce(DepTime, CRSDepTime) AS dep_hhmm,
        coalesce(ArrTime, CRSArrTime) AS arr_hhmm
    FROM ontime.fact_ontime
    WHERE FlightDate >= addYears(today(), -5)
      AND Cancelled = 0
      AND Carrier != ''
      AND FlightNum != ''
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, OriginCode, DestCode, dep_hhmm, arr_hhmm
),
recent_itineraries AS (
    SELECT
        FlightDate,
        Carrier,
        FlightNum,
        aircraft_id,
        count() AS hop_count,
        arraySort(groupArray((ifNull(dep_hhmm, 9999), ifNull(arr_hhmm, 9999), OriginCode, DestCode))) AS legs_sorted
    FROM recent_legs
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id
),
recent_routes AS (
    SELECT
        FlightDate,
        Carrier,
        FlightNum,
        aircraft_id,
        hop_count,
        arrayStringConcat(arrayConcat([tupleElement(legs_sorted[1], 3)], arrayMap(x -> tupleElement(x, 4), legs_sorted)), '-') AS Route
    FROM recent_itineraries
    WHERE hop_count >= 2
),
top_routes AS (
    SELECT
        if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,
        FlightNum AS flight_number,
        Carrier AS carrier,
        FlightDate AS flight_date,
        hop_count,
        Route,
        max(FlightDate) OVER (PARTITION BY Route) AS most_recent_flight_date,
        row_number() OVER (PARTITION BY Route ORDER BY hop_count DESC, FlightDate DESC, Carrier, FlightNum, aircraft_id) AS rn
    FROM recent_routes
),
top10 AS (
    SELECT aircraft_id, flight_number, carrier, flight_date, hop_count, Route, most_recent_flight_date
    FROM top_routes
    WHERE rn = 1
    ORDER BY hop_count DESC, most_recent_flight_date DESC, Route
    LIMIT 10
),
history_legs AS (
    SELECT
        FlightDate,
        Carrier,
        FlightNum,
        Tail_Number AS aircraft_id,
        OriginCode,
        DestCode,
        coalesce(DepTime, CRSDepTime) AS dep_hhmm,
        coalesce(ArrTime, CRSArrTime) AS arr_hhmm
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND Carrier != ''
      AND FlightNum != ''
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, OriginCode, DestCode, dep_hhmm, arr_hhmm
),
history_routes AS (
    SELECT
        FlightDate,
        arrayStringConcat(arrayConcat([tupleElement(legs_sorted[1], 3)], arrayMap(x -> tupleElement(x, 4), legs_sorted)), '-') AS Route
    FROM (
        SELECT
            FlightDate,
            Carrier,
            FlightNum,
            aircraft_id,
            count() AS hop_count,
            arraySort(groupArray((ifNull(dep_hhmm, 9999), ifNull(arr_hhmm, 9999), OriginCode, DestCode))) AS legs_sorted
        FROM history_legs
        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id
        HAVING hop_count >= 2
    )
    WHERE Route IN (SELECT Route FROM top10)
),
recurrence AS (
    SELECT Route, countDistinct(FlightDate) AS route_recurrence_count
    FROM history_routes
    GROUP BY Route
)
SELECT
    t.aircraft_id,
    t.flight_number,
    t.carrier,
    t.flight_date,
    t.hop_count,
    r.route_recurrence_count,
    t.Route
FROM top10 t
LEFT JOIN recurrence r USING (Route)
ORDER BY t.hop_count DESC, t.most_recent_flight_date DESC, t.Route
