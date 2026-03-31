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
    SELECT Route FROM daily_routes WHERE hop_count = (SELECT max_hop_count FROM max_hops)
),
route_recurrence AS (
    SELECT Route, countDistinct(FlightDate) AS recurrence_count
    FROM daily_routes
    WHERE Route IN (SELECT DISTINCT Route FROM top_route_days)
    GROUP BY Route
)
SELECT Route, recurrence_count
FROM route_recurrence
ORDER BY recurrence_count ASC
