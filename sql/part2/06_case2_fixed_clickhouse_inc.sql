SELECT avg(DepDelay) AS avg_delay
FROM ontime.ontime
WHERE IATA_CODE_Reporting_Airline = 'DL'
  AND trimRight(toString(Origin)) = 'ATL'
  AND Year = 2022
  AND Month = 7;
