WITH airport_stats AS (
    SELECT
        OriginAirportID,
        OriginCode,
        count() AS completed_departures,
        round(countIf(DepDel15 = 0) * 100.0 / count(), 2) AS otp_pct,
        round(avg(DepDelay), 2) AS avg_dep_delay_min,
        round(quantile(0.90)(DepDelay), 2) AS p90_dep_delay_min,
        min(FlightDate) AS first_date,
        max(FlightDate) AS last_date
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
    GROUP BY OriginAirportID, OriginCode
    HAVING completed_departures >= 10000
),
ranked AS (
    SELECT
        s.*,
        d.DisplayAirportName AS airport_name,
        row_number() OVER (ORDER BY otp_pct ASC) AS rank_worst,
        count() OVER () AS total_qualifying
    FROM airport_stats s
    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID
)
SELECT
    rank_worst,
    total_qualifying,
    OriginCode,
    airport_name,
    completed_departures,
    otp_pct,
    avg_dep_delay_min,
    p90_dep_delay_min,
    first_date,
    last_date
FROM ranked
ORDER BY rank_worst ASC
LIMIT 30
