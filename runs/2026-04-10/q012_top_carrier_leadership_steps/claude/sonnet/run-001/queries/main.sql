-- Yearly carrier leadership by completed flights
-- One row per year: leader, runner-up, share/gap metrics, and leadership transition flag
WITH
    -- Step 1: completed flights per carrier per year
    yearly_carrier AS (
        SELECT
            Year,
            Reporting_Airline AS carrier,
            countIf(Cancelled = 0) AS flights
        FROM ontime.fact_ontime
        GROUP BY Year, carrier
    ),
    -- Step 2: rank carriers within each year and compute year total
    with_rank AS (
        SELECT
            Year,
            carrier,
            flights,
            SUM(flights) OVER (PARTITION BY Year) AS year_total,
            rank() OVER (PARTITION BY Year ORDER BY flights DESC) AS rnk
        FROM yearly_carrier
    ),
    -- Step 3: pivot to one row per year with leader (rank 1) and runner-up (rank 2)
    pivoted AS (
        SELECT
            Year,
            maxIf(carrier, rnk = 1)            AS leader_carrier,
            maxIf(toInt64(flights), rnk = 1)   AS leader_flights,
            maxIf(carrier, rnk = 2)            AS runner_carrier,
            maxIf(toInt64(flights), rnk = 2)   AS runner_flights,
            max(year_total)                     AS year_total
        FROM with_rank
        WHERE rnk <= 2
        GROUP BY Year
    )
-- Step 4: compute share, gap, and leadership transition
SELECT
    Year,
    leader_carrier,
    leader_flights,
    round(leader_flights * 100.0 / year_total, 2)                    AS leader_share_pct,
    runner_carrier,
    runner_flights,
    round(runner_flights * 100.0 / year_total, 2)                    AS runner_share_pct,
    leader_flights - runner_flights                                    AS gap_flights,
    round((leader_flights - runner_flights) * 100.0 / year_total, 2) AS gap_share_pct_pts,
    lagInFrame(leader_carrier, 1, leader_carrier)
        OVER (ORDER BY Year ASC)                                      AS prev_leader,
    if(
        lagInFrame(leader_carrier, 1, leader_carrier)
            OVER (ORDER BY Year ASC) != leader_carrier,
        1, 0
    )                                                                 AS leadership_changed
FROM pivoted
ORDER BY Year
