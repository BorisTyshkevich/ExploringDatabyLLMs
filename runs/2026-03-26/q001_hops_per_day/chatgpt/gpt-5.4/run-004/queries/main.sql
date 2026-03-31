WITH params AS (
    SELECT addYears(toDate((SELECT max(FlightDate) FROM ontime.fact_ontime)), -5) AS start_date,
           (SELECT max(FlightDate) FROM ontime.fact_ontime) AS end_date
),
legs_raw AS (
    SELECT
        FlightDate,
        ifNull(nullIf(Tail_Number, ''), '') AS aircraft_id,
        Flight_Number_Reporting_Airline AS flight_number,
        IATA_CODE_Reporting_Airline AS carrier,
        OriginAirportID,
        DestAirportID,
        OriginCode,
        DestCode,
        coalesce(CRSDepTime, DepTime, 0) AS dep_hhmm,
        coalesce(CRSArrTime, ArrTime, 0) AS arr_hhmm,
        coalesce(DepTime, CRSDepTime, 0) AS dep_sort_hhmm,
        Distance,
        row_number() OVER (
            PARTITION BY FlightDate, ifNull(nullIf(Tail_Number, ''), ''), Flight_Number_Reporting_Airline,
                         IATA_CODE_Reporting_Airline, OriginAirportID, DestAirportID,
                         coalesce(CRSDepTime, DepTime, 0), coalesce(CRSArrTime, ArrTime, 0)
            ORDER BY Cancelled ASC, Diverted ASC, coalesce(ActualElapsedTime, CRSElapsedTime, 1000000) ASC
        ) AS rn
    FROM ontime.fact_ontime
    WHERE FlightDate >= (SELECT start_date FROM params)
      AND FlightDate <= (SELECT end_date FROM params)
      AND Cancelled = 0
),
legs AS (
    SELECT *,
           toDateTime(FlightDate) + toIntervalMinute(intDiv(dep_sort_hhmm,100)*60 + (dep_sort_hhmm % 100)) AS dep_ts
    FROM legs_raw
    WHERE rn = 1
),
itineraries AS (
    SELECT
        FlightDate,
        aircraft_id,
        flight_number,
        carrier,
        arraySort(groupArray((dep_ts, OriginCode, DestCode, OriginAirportID, DestAirportID, toUInt32(ifNull(Distance,0))))) AS legs_sorted
    FROM legs
    GROUP BY FlightDate, aircraft_id, flight_number, carrier
    HAVING length(legs_sorted) > 1
       AND arrayAll(i -> legs_sorted[i].5 = legs_sorted[i + 1].4, range(1, length(legs_sorted)))
),
routes AS (
    SELECT
        FlightDate,
        aircraft_id,
        flight_number,
        carrier,
        length(legs_sorted) AS hop_count,
        legs_sorted[1].1 AS first_dep_ts,
        arrayStringConcat(arrayConcat(arrayMap(x -> x.2, legs_sorted), [legs_sorted[length(legs_sorted)].3]), '-') AS Route
    FROM itineraries
),
route_days AS (
    SELECT Route, uniqExact(FlightDate) AS route_recurrence_count
    FROM routes
    GROUP BY Route
),
unique_ranked AS (
    SELECT r.*, rd.route_recurrence_count,
           row_number() OVER (PARTITION BY Route ORDER BY first_dep_ts DESC, FlightDate DESC, aircraft_id DESC, flight_number DESC, carrier DESC) AS route_rn
    FROM routes r
    INNER JOIN route_days rd USING (Route)
)
SELECT aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route
FROM unique_ranked
WHERE route_rn = 1
ORDER BY hop_count DESC, first_dep_ts DESC
LIMIT 10
