WITH airport_offsets AS (
    SELECT
        code,
        multiIf(
            utc_local_time_variation = '', 0,
            startsWith(utc_local_time_variation, '-'),
                -1 * ((toInt32OrZero(substring(utc_local_time_variation, 2, 2)) * 60) + toInt32OrZero(substring(utc_local_time_variation, 4, 2))),
            (toInt32OrZero(substring(utc_local_time_variation, 2, 2)) * 60) + toInt32OrZero(substring(utc_local_time_variation, 4, 2))
        ) AS utc_offset_minutes
    FROM ontime.airports_latest
),
legs AS (
    SELECT
        Tail_Number AS TailNum,
        Flight_Number_Reporting_Airline AS FlightNum,
        trimBoth(toString(Reporting_Airline)) AS Carrier,
        FlightDate,
        replaceAll(toString(Origin), '\0', '') AS Origin,
        replaceAll(toString(Dest), '\0', '') AS Dest,
        assumeNotNull(DepTime) AS DepTime,
        (toDateTime(FlightDate) + toIntervalHour(intDiv(DepTime, 100)) + toIntervalMinute(DepTime % 100))
            - toIntervalMinute(coalesce(ao.utc_offset_minutes, 0)) AS dep_ts_utc
    FROM ontime.ontime AS o
    LEFT JOIN airport_offsets AS ao
        ON replaceAll(toString(o.Origin), '\0', '') = ao.code
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
      AND Diverted = 0
      AND DepTime > 0
),
itineraries AS (
    SELECT
        TailNum,
        FlightNum,
        Carrier,
        FlightDate,
        count() AS hops,
        arraySort(x -> (x.1, x.2, x.3, x.4), groupArray((dep_ts_utc, Origin, Dest, DepTime))) AS route_legs
    FROM legs
    GROUP BY TailNum, FlightNum, Carrier, FlightDate
)
SELECT
    TailNum,
    FlightNum,
    Carrier,
    FlightDate,
    hops,
    arrayStringConcat(arrayMap(x -> concat(x.2, '-', x.3), route_legs), ' | ') AS Route
FROM itineraries
ORDER BY hops DESC, FlightDate DESC, TailNum, FlightNum
LIMIT 10
