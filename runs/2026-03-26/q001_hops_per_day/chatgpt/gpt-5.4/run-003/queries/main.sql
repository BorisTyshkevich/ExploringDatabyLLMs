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
        f.OriginCityMarketID AS origin_city_market_id,
        f.DestCityMarketID AS dest_city_market_id,
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
            origin_city_market_id,
            dest_city_market_id,
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
        r.FlightDate,
        r.carrier,
        r.flight_number,
        r.aircraft_id,
        r.first_dep_time,
        r.hop_count,
        rr.route_recurrence_count,
        r.Route,
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
)
SELECT
    if(aircraft_id = '', 'unknown', aircraft_id) AS aircraft_id,
    flight_number,
    carrier,
    FlightDate AS flight_date,
    hop_count,
    route_recurrence_count,
    Route
FROM ranked_unique
WHERE route_pick_rank = 1
ORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route
LIMIT 10
