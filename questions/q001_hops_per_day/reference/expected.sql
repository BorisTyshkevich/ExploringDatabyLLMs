WITH airport_offsets AS (
    SELECT
        IATA,
        toInt32(round(Timezone * 60)) AS utc_offset_minutes
    FROM default.airports
    WHERE IATA != ''
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
    FROM default.ontime_v2 AS o
    LEFT JOIN airport_offsets AS ao
        ON replaceAll(toString(o.Origin), '\0', '') = ao.IATA
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
