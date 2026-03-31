WITH itineraries AS (
    SELECT
        FlightDate AS flight_date,
        TailNum AS aircraft_id,
        FlightNum AS flight_number,
        Carrier AS carrier,
        arraySort(
            x -> (x.1, x.2, x.3),
            groupArray((coalesce(DepTime, CRSDepTime, toUInt16(0)), OriginCode, DestCode))
        ) AS legs
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND TailNum != ''
      AND FlightNum != ''
    GROUP BY
        flight_date,
        aircraft_id,
        flight_number,
        carrier
    HAVING count() > 1
),
ranked AS (
    SELECT
        flight_date,
        aircraft_id,
        flight_number,
        carrier,
        length(legs) AS hop_count,
        arrayStringConcat(arrayMap(x -> x.2, legs), '-') || '-' || arrayElement(legs, -1).3 AS Route
    FROM itineraries
),
max_hops AS (
    SELECT max(hop_count) AS max_hop_count
    FROM ranked
),
unique_routes AS (
    SELECT
        flight_date,
        aircraft_id,
        flight_number,
        carrier,
        hop_count,
        Route,
        row_number() OVER (
            PARTITION BY Route
            ORDER BY flight_date DESC, aircraft_id DESC, flight_number DESC, carrier DESC
        ) AS route_rank
    FROM ranked
    WHERE hop_count = (SELECT max_hop_count FROM max_hops)
)
SELECT
    aircraft_id,
    flight_number,
    carrier,
    flight_date,
    hop_count,
    Route
FROM unique_routes
WHERE route_rank = 1
ORDER BY flight_date DESC, Route ASC
LIMIT 10
