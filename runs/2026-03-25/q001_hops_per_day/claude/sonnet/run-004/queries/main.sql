WITH raw AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        IATA_CODE_Reporting_Airline AS Carrier,
        FlightDate,
        count() AS hop_count,
        arraySort(x -> x.3, groupArray((OriginCode, DestCode, assumeNotNull(CRSDepTime)))) AS sorted_legs
    FROM ontime.fact_ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate
),
daily_routes AS (
    SELECT
        Tail_Number, FlightNum, Carrier, FlightDate, hop_count,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.1, sorted_legs),
                [sorted_legs[length(sorted_legs)].2]
            ),
            '-'
        ) AS Route
    FROM raw
),
max_hops AS (
    SELECT max(hop_count) AS max_hop_count FROM daily_routes
),
top_route_days AS (
    SELECT * FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)
),
route_recurrence AS (
    SELECT Route, countDistinct(FlightDate) AS recurrence_count
    FROM daily_routes
    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)
    GROUP BY Route
),
recent_top AS (
    SELECT
        Route,
        argMax(Tail_Number, FlightDate) AS Tail_Number,
        argMax(FlightNum, FlightDate) AS FlightNum,
        argMax(Carrier, FlightDate) AS Carrier,
        max(FlightDate) AS most_recent_date,
        argMax(hop_count, FlightDate) AS hop_count
    FROM top_route_days
    GROUP BY Route
)
SELECT
    rt.Tail_Number,
    rt.FlightNum,
    rt.Carrier,
    rt.most_recent_date AS FlightDate,
    rt.hop_count,
    rr.recurrence_count,
    rt.Route
FROM recent_top rt
JOIN route_recurrence rr ON rt.Route = rr.Route
ORDER BY rt.most_recent_date DESC
LIMIT 10
