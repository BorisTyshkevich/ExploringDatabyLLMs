WITH main_top10 AS (
    WITH dedup_legs AS (
        SELECT
            FlightDate,
            ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,
            Flight_Number_Reporting_Airline AS flight_number,
            IATA_CODE_Reporting_Airline AS carrier,
            OriginCode,
            DestCode,
            coalesce(CRSDepTime, DepTime, toUInt16(9999)) AS dep_sort,
            min(coalesce(CRSArrTime, ArrTime, toUInt16(9999))) AS arr_sort
        FROM ontime.fact_ontime
        WHERE Cancelled = 0
        GROUP BY FlightDate, aircraft_id, flight_number, carrier, OriginCode, DestCode, dep_sort
    ), itineraries AS (
        SELECT
            FlightDate,
            aircraft_id,
            flight_number,
            carrier,
            count() AS hop_count,
            min(dep_sort) AS first_dep_sort,
            arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs
        FROM dedup_legs
        GROUP BY FlightDate, aircraft_id, flight_number, carrier
        HAVING count() > 1
    ), route_recurrence AS (
        SELECT
            arrayStringConcat(arrayConcat(arrayMap(x -> x.3, legs), [arrayElement(legs, length(legs)).4]), '-') AS route,
            uniqExact(FlightDate) AS route_recurrence_count
        FROM itineraries
        GROUP BY route
    ), max_hops AS (
        SELECT max(hop_count) AS max_hop_count FROM itineraries
    ), top_unique AS (
        SELECT
            i.FlightDate,
            i.first_dep_sort,
            arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,
            rr.route_recurrence_count,
            row_number() OVER (
                PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')
                ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC
            ) AS route_rank
        FROM itineraries i
        CROSS JOIN max_hops m
        INNER JOIN route_recurrence rr
            ON arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') = rr.route
        WHERE i.hop_count = m.max_hop_count
    )
    SELECT route, route_recurrence_count, FlightDate, first_dep_sort
    FROM top_unique
    WHERE route_rank = 1
    ORDER BY FlightDate DESC, first_dep_sort DESC
    LIMIT 10
)
SELECT
    route,
    route_recurrence_count,
    multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count < 10, 'occasional', 'recurring') AS recurrence_tier,
    count() OVER (PARTITION BY multiIf(route_recurrence_count = 1, 'one-off', route_recurrence_count < 10, 'occasional', 'recurring')) AS tier_count
FROM main_top10
ORDER BY FlightDate DESC, first_dep_sort DESC
