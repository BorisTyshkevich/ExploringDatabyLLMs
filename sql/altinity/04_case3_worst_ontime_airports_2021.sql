SELECT
    trimRight(toString(Origin)) AS airport,
    count() AS flights,
    round(100.0 * avg(ArrDel15), 2) AS delayed_arrival_pct
FROM default.ontime
WHERE Year = 2021
GROUP BY airport
HAVING flights >= 10000
ORDER BY delayed_arrival_pct DESC
LIMIT 10;
