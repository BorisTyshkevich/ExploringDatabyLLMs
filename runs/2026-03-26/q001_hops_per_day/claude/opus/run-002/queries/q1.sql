WITH candidates AS (
    SELECT Tail_Number, FlightNum, Carrier, FlightDate
    FROM (
        SELECT Tail_Number, FlightNum, Carrier, FlightDate,
               uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS hops
        FROM ontime.fact_ontime
        GROUP BY Tail_Number, FlightNum, Carrier, FlightDate
        HAVING hops = (
            SELECT max(h) FROM (
                SELECT uniqExact(OriginCode, DestCode, assumeNotNull(CRSDepTime)) AS h
                FROM ontime.fact_ontime
                GROUP BY Tail_Number, FlightNum, Carrier, FlightDate
            )
        )
    )
),
deduped_legs AS (
    SELECT DISTINCT f.Tail_Number, f.FlightNum, f.Carrier, f.FlightDate,
           f.OriginCode, f.DestCode, assumeNotNull(f.CRSDepTime) AS dep
    FROM ontime.fact_ontime f
    INNER JOIN candidates c
        ON f.Tail_Number = c.Tail_Number AND f.FlightNum = c.FlightNum
        AND f.Carrier = c.Carrier AND f.FlightDate = c.FlightDate
),
with_route AS (
    SELECT Tail_Number, FlightNum, Carrier, FlightDate,
           count() AS hops,
           concat(
               arrayStringConcat(groupArray(OriginCode), '-'),
               '-',
               arrayElement(groupArray(DestCode), count()::Int32)
           ) AS Route
    FROM (SELECT * FROM deduped_legs ORDER BY Tail_Number, FlightNum, Carrier, FlightDate, dep)
    GROUP BY Tail_Number, FlightNum, Carrier, FlightDate
),
route_recurrence AS (
    SELECT Route, uniqExact(FlightDate) AS recurrence_count
    FROM with_route
    GROUP BY Route
),
ranked AS (
    SELECT w.Route, r.recurrence_count,
           ROW_NUMBER() OVER (PARTITION BY w.Route ORDER BY w.FlightDate DESC) AS rn,
           w.FlightDate
    FROM with_route w
    JOIN route_recurrence r ON w.Route = r.Route
)
SELECT Route, recurrence_count,
       FlightDate AS most_recent_date
FROM ranked
WHERE rn = 1
ORDER BY FlightDate DESC
LIMIT 10
