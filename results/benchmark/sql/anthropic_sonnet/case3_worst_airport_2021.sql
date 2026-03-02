SELECT
    Origin AS airport,
    round(countIf(ArrDel15 = 1) * 100.0 / count(), 2) AS delayed_arrival_pct,
    count() AS flights
FROM ontime
WHERE Year = 2021
GROUP BY Origin
HAVING flights >= 10000
ORDER BY delayed_arrival_pct DESC
LIMIT 1
