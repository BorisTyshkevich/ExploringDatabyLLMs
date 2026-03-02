SELECT trimRight(toString(Carrier)) AS carrier, count() AS flights FROM ontime WHERE Year = 2019 GROUP BY Carrier ORDER BY flights DESC LIMIT 1
