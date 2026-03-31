WITH leg_rows AS (
    SELECT
        TailNum,
        FlightNum,
        Carrier,
        FlightDate,
        OriginCode,
        DestCode,
        addMinutes(toDateTime(FlightDate), intDiv(ifNull(DepTime, CRSDepTime), 100) * 60 + modulo(ifNull(DepTime, CRSDepTime), 100)) AS dep_ts
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND Diverted = 0
      AND TailNum != ''
      AND FlightNum != ''
      AND ifNull(DepTime, CRSDepTime) IS NOT NULL
), itineraries AS (
    SELECT
        TailNum,
        FlightNum,
        Carrier,
        FlightDate,
        arraySort(x -> x.1, groupArray((dep_ts, OriginCode, DestCode))) AS ordered_legs
    FROM leg_rows
    GROUP BY
        TailNum,
        FlightNum,
        Carrier,
        FlightDate
    HAVING length(ordered_legs) > 1
)
SELECT
    TailNum AS `Aircraft ID`,
    FlightNum AS `Flight Number`,
    Carrier,
    FlightDate AS `Date`,
    concat(
        arrayStringConcat(arrayMap(x -> x.2, ordered_legs), '→'),
        '→',
        tupleElement(arrayElement(ordered_legs, length(ordered_legs)), 3)
    ) AS Route,
    length(ordered_legs) AS hop_count
FROM itineraries
ORDER BY hop_count DESC, `Date` DESC, `Aircraft ID`, `Flight Number`
LIMIT 10
