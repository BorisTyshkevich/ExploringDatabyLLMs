SELECT 
    Origin,
    Dest,
    CASE 
        WHEN Month IN (12,1,2) THEN 'Winter'
        WHEN Month IN (3,4,5) THEN 'Spring'
        WHEN Month IN (6,7,8) THEN 'Summer'
        WHEN Month IN (9,10,11) THEN 'Fall'
    END AS season,
    avg(DepDelay) AS avg_dep_delay,
    count() AS flights
FROM ontime
WHERE Year BETWEEN 2019 AND 2021
    AND (Origin, Dest) IN (
        SELECT Origin, Dest
        FROM ontime
        WHERE Year BETWEEN 2019 AND 2021
        GROUP BY Origin, Dest
        ORDER BY count() DESC
        LIMIT 50
    )
GROUP BY Origin, Dest, season
ORDER BY avg_dep_delay DESC
LIMIT 1
