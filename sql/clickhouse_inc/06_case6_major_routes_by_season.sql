WITH routes AS (
    SELECT
        trimRight(toString(Origin)) AS origin,
        trimRight(toString(Dest)) AS dest,
        count() AS flights
    FROM ontime.ontime
    WHERE Year BETWEEN 2019 AND 2021
    GROUP BY origin, dest
    ORDER BY flights DESC
    LIMIT 50
)
SELECT
    concat(r.origin, '-', r.dest) AS route,
    multiIf(o.Month IN (12, 1, 2), 'Winter', o.Month IN (3, 4, 5), 'Spring', o.Month IN (6, 7, 8), 'Summer', 'Fall') AS season,
    round(avg(o.DepDelay), 2) AS avg_dep_delay,
    count() AS flights
FROM ontime.ontime o
INNER JOIN routes r ON trimRight(toString(o.Origin)) = r.origin AND trimRight(toString(o.Dest)) = r.dest
WHERE o.Year BETWEEN 2019 AND 2021
GROUP BY route, season
HAVING flights >= 1000
ORDER BY avg_dep_delay DESC
LIMIT 12;
