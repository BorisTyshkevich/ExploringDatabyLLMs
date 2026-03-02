SELECT carrier, COUNT(*) AS flights
FROM ontime
WHERE year = 2019
GROUP BY carrier
ORDER BY flights DESC
LIMIT 1;
