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
    SELECT Year, maxIf(carrier, rnk = 1) AS leader
    FROM ranked
    GROUP BY Year
),
with_era_change AS (
    SELECT
        Year,
        leader,
        leader != lagInFrame(leader) OVER (ORDER BY Year) AS era_start
    FROM leader_per_year
    ORDER BY Year
),
with_era_id AS (
    SELECT
        Year,
        leader,
        sum(era_start) OVER (ORDER BY Year ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS era_id
    FROM with_era_change
)
SELECT
    leader,
    min(Year) AS era_start_year,
    max(Year) AS era_end_year,
    count()   AS years_in_era
FROM with_era_id
GROUP BY era_id, leader
ORDER BY era_start_year
