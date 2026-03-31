WITH airport_offsets AS
(
    SELECT
        airport_id,
        any(utc_local_time_variation) AS utc_local_time_variation
    FROM ontime.airports_latest
    GROUP BY airport_id
),
legs AS
(
    SELECT
        o.FlightDate,
        o.TailNum,
        o.FlightNum,
        o.Carrier,
        o.OriginAirportID,
        o.DestAirportID,
        o.Origin,
        o.Dest,
        o.DepTime,
        toDateTime(o.FlightDate)
            + toIntervalDay(if(o.DepTime = 2400, 1, 0))
            + toIntervalMinute(intDiv(if(o.DepTime = 2400, 0, o.DepTime), 100) * 60 + modulo(if(o.DepTime = 2400, 0, o.DepTime), 100)) AS dep_local_ts,
        (
            if(length(a1o.utc_local_time_variation) = 5,
                if(substring(a1o.utc_local_time_variation, 1, 1) = '-', -1, 1)
                * (toInt32(substring(a1o.utc_local_time_variation, 2, 2)) * 60 + toInt32(substring(a1o.utc_local_time_variation, 4, 2))),
                0
            )
        ) AS origin_offset_minutes
    FROM ontime.ontime AS o
    LEFT JOIN airport_offsets AS a1o ON o.OriginAirportID = a1o.airport_id
    WHERE o.Cancelled = 0
      AND o.Diverted = 0
      AND o.DepTime IS NOT NULL
      AND o.TailNum != ''
      AND o.FlightNum != ''
),
itineraries AS
(
    SELECT
        FlightDate,
        TailNum,
        FlightNum,
        Carrier,
        length(ordered_legs) AS Hops,
        arrayStringConcat(
            arrayMap(x -> concat(formatDateTime(x.2, '%H:%i'), ' ', x.3, '->', x.4), ordered_legs),
            ' | '
        ) AS Route
    FROM
    (
        SELECT
            FlightDate,
            TailNum,
            FlightNum,
            Carrier,
            arraySort(x -> x.1, groupArray((
                dep_local_ts - toIntervalMinute(origin_offset_minutes),
                dep_local_ts,
                Origin,
                Dest
            ))) AS ordered_legs
        FROM legs
        GROUP BY
            FlightDate,
            TailNum,
            FlightNum,
            Carrier
    )
),
max_hops AS
(
    SELECT max(Hops) AS max_hops_observed FROM itineraries
),
max_hops_counts AS
(
    SELECT count() AS max_hop_itinerary_count
    FROM itineraries
    CROSS JOIN max_hops
    WHERE Hops = max_hops_observed
)
SELECT
    TailNum AS `Aircraft ID`,
    FlightNum AS `Flight Number`,
    Carrier,
    FlightDate AS Date,
    Hops,
    Route,
    max_hops_observed AS `Maximum Hops Observed`,
    max_hop_itinerary_count AS `Maximum-Hop Itinerary Count`
FROM itineraries
CROSS JOIN max_hops
CROSS JOIN max_hops_counts
ORDER BY
    Hops DESC,
    Date DESC,
    Carrier,
    `Flight Number`,
    `Aircraft ID`
LIMIT 10
