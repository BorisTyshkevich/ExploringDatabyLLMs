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
leader_runner AS (
    SELECT
        Year,
        maxIf(carrier, rnk = 1)          AS leader,
        maxIf(completed_flights, rnk = 1) AS leader_flights,
        maxIf(carrier, rnk = 2)          AS runner_up,
        maxIf(completed_flights, rnk = 2) AS runner_up_flights
    FROM ranked
    GROUP BY Year
),
with_lag AS (
    SELECT
        Year,
        leader,
        leader_flights,
        runner_up,
        runner_up_flights,
        leader_flights - runner_up_flights                       AS gap_vs_runner_up,
        lagInFrame(leader)         OVER (ORDER BY Year)          AS prev_leader,
        lagInFrame(leader_flights) OVER (ORDER BY Year)          AS prev_leader_flights,
        leader_flights - lagInFrame(leader_flights) OVER (ORDER BY Year) AS swing_vs_prior
    FROM leader_runner
    ORDER BY Year
)
SELECT
    Year,
    leader          AS new_leader,
    leader_flights,
    prev_leader,
    prev_leader_flights,
    swing_vs_prior,
    gap_vs_runner_up
FROM with_lag
WHERE leader != prev_leader
  AND prev_leader != ''
ORDER BY swing_vs_prior DESC
LIMIT 1
