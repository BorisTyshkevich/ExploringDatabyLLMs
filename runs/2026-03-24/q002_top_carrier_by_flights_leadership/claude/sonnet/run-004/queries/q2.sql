WITH carrier_year AS (
    SELECT
        Year,
        Reporting_Airline AS carrier,
        countIf(Cancelled = 0) AS completed_flights
    FROM ontime.fact_ontime
    GROUP BY Year, carrier
),
ranked AS (
    SELECT
        Year,
        carrier,
        completed_flights,
        rank() OVER (PARTITION BY Year ORDER BY completed_flights DESC) AS rnk
    FROM carrier_year
),
leader_per_year AS (
    SELECT
        Year,
        maxIf(carrier, rnk = 1)          AS leader,
        maxIf(completed_flights, rnk = 1) AS leader_flights
    FROM ranked
    GROUP BY Year
),
transitions AS (
    SELECT
        Year,
        leader,
        leader_flights,
        lagInFrame(leader)          OVER (ORDER BY Year) AS prev_leader,
        lagInFrame(leader_flights)  OVER (ORDER BY Year) AS prev_leader_flights
    FROM leader_per_year
    ORDER BY Year
)
SELECT
    Year,
    leader          AS new_leader,
    leader_flights  AS new_leader_flights,
    prev_leader,
    prev_leader_flights,
    leader_flights - prev_leader_flights AS swing
FROM transitions
WHERE leader != prev_leader
  AND prev_leader != ''
ORDER BY Year
