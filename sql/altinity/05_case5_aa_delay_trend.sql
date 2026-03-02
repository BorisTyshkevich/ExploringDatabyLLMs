SELECT
    Year,
    Month,
    round(avg(DepDelay), 2) AS avg_dep_delay,
    count() AS flights
FROM default.ontime
WHERE trimRight(toString(Carrier)) = 'AA'
  AND Year BETWEEN 2018 AND 2021
GROUP BY Year, Month
ORDER BY Year, Month;
