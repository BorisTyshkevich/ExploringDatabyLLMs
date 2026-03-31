WITH max_dt AS (
    SELECT max(FlightDate) AS mx
    FROM ontime.fact_ontime
),
base AS (
    SELECT
        f.FlightDate,
        f.IATA_CODE_Reporting_Airline AS carrier,
        f.Flight_Number_Reporting_Airline AS flight_number,
        ifNull(nullIf(f.Tail_Number, ''), '') AS aircraft_id,
        f.OriginCode AS origin_code,
        f.DestCode AS dest_code,
        f.OriginAirportID AS origin_airport_id,
        f.DestAirportID AS dest_airport_id,
        coalesce(f.CRSDepTime, f.DepTime, 9999) AS dep_sort_time,
        coalesce(f.Distance, 0) AS distance,
        row_number() OVER (
            PARTITION BY
                f.FlightDate,
                f.IATA_CODE_Reporting_Airline,
                f.Flight_Number_Reporting_Airline,
                ifNull(nullIf(f.Tail_Number, ''), ''),
                f.OriginCode,
                f.DestCode,
                coalesce(f.CRSDepTime, f.DepTime, 9999)
            ORDER BY
                f.DepTime DESC NULLS LAST,
                f.ArrTime DESC NULLS LAST,
                coalesce(f.Distance, 0) DESC
        ) AS rn
    FROM ontime.fact_ontime AS f
    CROSS JOIN max_dt
    WHERE f.Cancelled = 0
      AND f.FlightDate >= addYears(mx, -5)
),
legs AS (
    SELECT *
    FROM base
    WHERE rn = 1
),
itineraries AS (
    SELECT
        FlightDate,
        carrier,
        flight_number,
        aircraft_id,
        arraySort(groupArray((
            dep_sort_time,
            origin_code,
            dest_code,
            origin_airport_id,
            dest_airport_id,
            distance
        ))) AS legs_arr,
        min(dep_sort_time) AS first_dep_time
    FROM legs
    GROUP BY FlightDate, carrier, flight_number, aircraft_id
),
routes AS (
    SELECT
        FlightDate,
        carrier,
        flight_number,
        aircraft_id,
        first_dep_time,
        length(legs_arr) AS hop_count,
        legs_arr,
        arrayStringConcat(
            arrayConcat(
                [tupleElement(arrayElement(legs_arr, 1), 2)],
                arrayMap(x -> tupleElement(x, 3), legs_arr)
            ),
            '-'
        ) AS Route
    FROM itineraries
),
route_recurrence AS (
    SELECT
        Route,
        uniqExact(FlightDate) AS route_recurrence_count
    FROM routes
    GROUP BY Route
),
ranked_unique AS (
    SELECT
        r.*,
        rr.route_recurrence_count,
        row_number() OVER (
            PARTITION BY r.Route
            ORDER BY
                r.hop_count DESC,
                r.FlightDate DESC,
                r.first_dep_time DESC,
                r.carrier,
                r.flight_number,
                r.aircraft_id
        ) AS route_pick_rank
    FROM routes AS r
    INNER JOIN route_recurrence AS rr
        ON r.Route = rr.Route
),
top_routes AS (
    SELECT *
    FROM ranked_unique
    WHERE route_pick_rank = 1
    ORDER BY hop_count DESC, FlightDate DESC, first_dep_time DESC, Route
    LIMIT 10
),
route_airports AS (
    SELECT
        Route,
        arrayJoin(
            arrayZip(
                arrayConcat(
                    [tupleElement(arrayElement(legs_arr, 1), 2)],
                    arrayMap(x -> tupleElement(x, 3), legs_arr)
                ),
                arrayEnumerate(
                    arrayConcat(
                        [tupleElement(arrayElement(legs_arr, 1), 2)],
                        arrayMap(x -> tupleElement(x, 3), legs_arr)
                    )
                ),
                arrayMap(
                    x -> length(legs_arr) + 1,
                    arrayConcat(
                        [tupleElement(arrayElement(legs_arr, 1), 2)],
                        arrayMap(x -> tupleElement(x, 3), legs_arr)
                    )
                )
            )
        ) AS airport_pos
    FROM top_routes
)
SELECT
    trimBoth(toString(airport_pos.1)) AS airport_code,
    any(a.DisplayAirportName) AS airport_name,
    any(concat(a.CityName, ', ', a.StateCode)) AS city_state,
    count() AS total_appearances,
    countIf(airport_pos.2 = 1) AS origin_appearances,
    countIf(airport_pos.2 > 1 AND airport_pos.2 < airport_pos.3) AS intermediate_stop_appearances,
    countIf(airport_pos.2 = airport_pos.3) AS final_destination_appearances,
    round(uniqExact(Route) / 10.0, 3) AS share_of_itineraries_containing_airport
FROM route_airports AS ra
LEFT JOIN ontime.dim_airports AS a
    ON a.AirportCode = ra.airport_pos.1
   AND a.IsLatest = 1
GROUP BY airport_code
ORDER BY total_appearances DESC, share_of_itineraries_containing_airport DESC, airport_code
