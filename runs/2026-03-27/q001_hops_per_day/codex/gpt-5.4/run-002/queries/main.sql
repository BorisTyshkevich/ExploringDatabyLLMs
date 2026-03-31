WITH
    (SELECT max(FlightDate) FROM ontime.fact_ontime) AS max_fd,
    addYears(max_fd, -5) AS start_fd,
    base AS (
        SELECT
            FlightDate,
            Carrier,
            FlightNum,
            if(Tail_Number = '', '', Tail_Number) AS aircraft_id,
            OriginCode,
            DestCode,
            coalesce(CRSDepTime, DepTime, 0) AS dep_sort,
            coalesce(CRSArrTime, ArrTime, 0) AS arr_sort
        FROM ontime.fact_ontime
        WHERE FlightDate > start_fd
          AND Cancelled = 0
          AND FlightNum != ''
          AND Carrier != ''
    ),
    dedup AS (
        SELECT
            FlightDate,
            Carrier,
            FlightNum,
            aircraft_id,
            dep_sort,
            min(arr_sort) AS arr_sort,
            min(OriginCode) AS OriginCode,
            min(DestCode) AS DestCode
        FROM base
        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, dep_sort
    ),
    itineraries AS (
        SELECT
            FlightDate,
            Carrier,
            FlightNum,
            aircraft_id,
            min(dep_sort) AS first_dep_time,
            arraySort(x -> (x.1, x.2, x.3, x.4), groupArray((dep_sort, arr_sort, OriginCode, DestCode))) AS legs,
            length(legs) AS hop_count,
            arrayStringConcat(arrayConcat(arrayMap(x -> x.3, legs), [legs[length(legs)].4]), '-') AS Route
        FROM dedup
        GROUP BY FlightDate, Carrier, FlightNum, aircraft_id
        HAVING hop_count > 1
    ),
    route_days AS (
        SELECT
            Route,
            countDistinct(FlightDate) AS route_recurrence_count
        FROM itineraries
        GROUP BY Route
    ),
    ranked AS (
        SELECT
            i.aircraft_id,
            i.FlightNum,
            i.Carrier,
            i.FlightDate,
            i.hop_count,
            r.route_recurrence_count,
            i.Route,
            i.first_dep_time,
            row_number() OVER (
                PARTITION BY i.Route
                ORDER BY i.FlightDate DESC, i.first_dep_time DESC, i.hop_count DESC, i.Carrier, i.FlightNum, i.aircraft_id
            ) AS rn
        FROM itineraries AS i
        INNER JOIN route_days AS r ON i.Route = r.Route
    )
SELECT
    aircraft_id,
    FlightNum AS flight_number,
    Carrier AS carrier,
    FlightDate AS flight_date,
    hop_count,
    route_recurrence_count,
    Route
FROM ranked
WHERE rn = 1
ORDER BY hop_count DESC, flight_date DESC, first_dep_time DESC, Route
LIMIT 10
