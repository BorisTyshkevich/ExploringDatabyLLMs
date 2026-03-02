SELECT Year, Month, avg(DepDelay) AS avg_dep_delay FROM ontime WHERE Year BETWEEN 2018 AND 2021 AND trimRight(toString(Carrier)) = 'AA' GROUP BY Year, Month ORDER BY avg_dep_delay DESC LIMIT 1
