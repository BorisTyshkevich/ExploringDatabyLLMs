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
        count() OVER () AS total_qualifying,
        multiIf(
            s.completed_departures >= 2000000, 'Major hub (2M+ dep)',
            s.completed_departures >= 500000,  'Large airport (500k-2M dep)',
            s.completed_departures >= 100000,  'Mid-size (100k-500k dep)',
            'Regional/small (<100k dep)'
        ) AS volume_tier
    FROM airport_stats s
    LEFT JOIN ontime.dim_airports d ON s.OriginAirportID = d.AirportID
)
SELECT
    rank_worst,
    OriginCode,
    airport_name,
    completed_departures,
    volume_tier,
    otp_pct,
    avg_dep_delay_min,
    p90_dep_delay_min
FROM ranked
WHERE rank_worst <= 30
ORDER BY rank_worst ASC
