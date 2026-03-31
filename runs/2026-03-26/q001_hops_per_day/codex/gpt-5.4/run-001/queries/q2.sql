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
), max_hops AS (
    SELECT max(hop_count) AS max_hop_count FROM itineraries
), top_unique AS (
    SELECT
        i.FlightDate,
        i.first_dep_sort,
        i.legs,
        arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-') AS route,
        row_number() OVER (
            PARTITION BY arrayStringConcat(arrayConcat(arrayMap(x -> x.3, i.legs), [arrayElement(i.legs, length(i.legs)).4]), '-')
            ORDER BY i.FlightDate DESC, i.first_dep_sort DESC, i.carrier DESC, i.flight_number DESC, i.aircraft_id DESC
        ) AS route_rank
    FROM itineraries i
    CROSS JOIN max_hops m
    WHERE i.hop_count = m.max_hop_count
), top10 AS (
    SELECT *
    FROM top_unique
    WHERE route_rank = 1
    ORDER BY FlightDate DESC, first_dep_sort DESC
    LIMIT 10
), airport_mentions AS (
    SELECT arrayJoin(arrayConcat(arrayMap(x -> x.3, legs), [arrayElement(legs, length(legs)).4])) AS airport_code
    FROM top10
)
SELECT
    a.AirportCode AS airport_code,
    a.StateCode AS state_code,
    a.CountryCodeISO AS country_code,
    multiIf(a.Longitude <= -115, 'far_west', a.Longitude <= -95, 'interior_west_central', 'east') AS geo_bucket,
    count() AS airport_mentions
FROM airport_mentions am
LEFT JOIN ontime.dim_airports a
    ON am.airport_code = toString(a.AirportCode)
   AND a.IsLatest = 1
GROUP BY airport_code, state_code, country_code, geo_bucket
ORDER BY airport_mentions DESC, airport_code
LIMIT 100
