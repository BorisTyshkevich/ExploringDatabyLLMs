SELECT
    IATA_CODE_Reporting_Airline AS carrier,
    count() AS flights
FROM ontime.ontime
WHERE Year = 2019
GROUP BY carrier
ORDER BY flights DESC
LIMIT 5;
