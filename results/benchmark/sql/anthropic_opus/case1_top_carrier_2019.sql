SELECT trimRight(toString(Carrier)) AS carrier, count() AS flights FROM default.ontime WHERE Year = 2019 GROUP BY carrier ORDER BY flights DESC LIMIT 1
