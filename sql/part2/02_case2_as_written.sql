SELECT AVG(depdelay) AS avg_delay
FROM ontime
WHERE carrier = 'DL'
  AND origin = 'ATL'
  AND year = 2022
  AND month = 7;
