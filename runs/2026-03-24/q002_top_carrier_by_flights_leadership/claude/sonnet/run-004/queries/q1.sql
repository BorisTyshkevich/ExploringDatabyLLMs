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
leaders AS (
    SELECT Year, carrier
    FROM ranked
    WHERE rnk = 1
)
SELECT
    carrier,
    count() AS years_led
FROM leaders
GROUP BY carrier
ORDER BY years_led DESC
