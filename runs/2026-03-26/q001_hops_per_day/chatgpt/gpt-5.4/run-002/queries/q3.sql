WITH recent_legs AS (
    SELECT
        FlightDate,
        Carrier,
        FlightNum,
        Tail_Number AS aircraft_id,
        OriginAirportID,
        DestAirportID,
        OriginCode,
        DestCode,
        coalesce(DepTime, CRSDepTime) AS dep_hhmm,
        coalesce(ArrTime, CRSArrTime) AS arr_hhmm,
        min(DepDelay) AS DepDelay,
        min(ArrDelay) AS ArrDelay,
        max(DepDel15) AS DepDel15,
        max(ArrDel15) AS ArrDel15,
        max(Diverted) AS Diverted
    FROM ontime.fact_ontime
    WHERE FlightDate >= addYears(today(), -5)
      AND Cancelled = 0 AND Carrier != '' AND FlightNum != ''
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id, OriginAirportID, DestAirportID, OriginCode, DestCode, dep_hhmm, arr_hhmm
),
recent_itineraries AS (
    SELECT
        FlightDate, Carrier, FlightNum, aircraft_id,
        count() AS hop_count,
        arraySort(groupArray((ifNull(dep_hhmm, 9999), ifNull(arr_hhmm, 9999), OriginCode, DestCode, OriginAirportID, DestAirportID, DepDelay, ArrDelay, DepDel15, ArrDel15, Diverted))) AS legs_sorted
    FROM recent_legs
    GROUP BY FlightDate, Carrier, FlightNum, aircraft_id
),
recent_routes AS (
    SELECT
        FlightDate, Carrier, FlightNum, aircraft_id, hop_count, legs_sorted,
        arrayStringConcat(arrayConcat([tupleElement(legs_sorted[1], 3)], arrayMap(x -> tupleElement(x, 4), legs_sorted)), '-') AS Route
    FROM recent_itineraries
    WHERE hop_count >= 2
),
top_routes AS (
    SELECT *, max(FlightDate) OVER (PARTITION BY Route) AS most_recent_flight_date,
           row_number() OVER (PARTITION BY Route ORDER BY hop_count DESC, FlightDate DESC, Carrier, FlightNum, aircraft_id) AS rn
    FROM recent_routes
),
top10 AS (
    SELECT Route, legs_sorted
    FROM top_routes
    WHERE rn = 1
    ORDER BY hop_count DESC, most_recent_flight_date DESC, Route
    LIMIT 10
),
leg_rows AS (
    SELECT
        Route,
        pos AS leg_position,
        tupleElement(leg, 3) AS origin_code,
        tupleElement(leg, 4) AS dest_code,
        tupleElement(leg, 7) AS dep_delay,
        tupleElement(leg, 8) AS arr_delay,
        tupleElement(leg, 9) AS dep_del15,
        tupleElement(leg, 10) AS arr_del15,
        tupleElement(leg, 11) AS diverted
    FROM top10
    ARRAY JOIN legs_sorted AS leg, arrayEnumerate(legs_sorted) AS pos
)
SELECT *
FROM (
    SELECT
        'airport' AS entity_type,
        origin_code AS entity_key,
        concat('departures from ', origin_code) AS entity_label,
        CAST(NULL AS Nullable(UInt64)) AS stop_position,
        round(avg(toFloat64(dep_delay)), 2) AS avg_departure_delay,
        round(avg(toFloat64(arr_delay)), 2) AS avg_arrival_delay,
        round(avg(greatest(toUInt8(dep_del15), toUInt8(arr_del15))), 3) AS delay_15_plus_rate,
        round(avg(toFloat64(diverted)), 3) AS diversion_incidence
    FROM leg_rows
    GROUP BY origin_code

    UNION ALL

    SELECT
        'leg' AS entity_type,
        concat(origin_code, '-', dest_code) AS entity_key,
        concat(origin_code, ' to ', dest_code) AS entity_label,
        CAST(NULL AS Nullable(UInt64)) AS stop_position,
        round(avg(toFloat64(dep_delay)), 2) AS avg_departure_delay,
        round(avg(toFloat64(arr_delay)), 2) AS avg_arrival_delay,
        round(avg(greatest(toUInt8(dep_del15), toUInt8(arr_del15))), 3) AS delay_15_plus_rate,
        round(avg(toFloat64(diverted)), 3) AS diversion_incidence
    FROM leg_rows
    GROUP BY origin_code, dest_code

    UNION ALL

    SELECT
        'stop_position' AS entity_type,
        toString(leg_position) AS entity_key,
        concat('leg ', toString(leg_position)) AS entity_label,
        toUInt64(leg_position) AS stop_position,
        round(avg(toFloat64(dep_delay)), 2) AS avg_departure_delay,
        round(avg(toFloat64(arr_delay)), 2) AS avg_arrival_delay,
        round(avg(greatest(toUInt8(dep_del15), toUInt8(arr_del15))), 3) AS delay_15_plus_rate,
        round(avg(toFloat64(diverted)), 3) AS diversion_incidence
    FROM leg_rows
    GROUP BY leg_position
)
ORDER BY entity_type, delay_15_plus_rate DESC, avg_arrival_delay DESC, entity_key
