SELECT
    avg(DepDelay) AS avg_dep_delay
FROM default.ontime
WHERE trimRight(toString(Carrier)) = 'DL'
  AND trimRight(toString(Origin)) = 'ATL'
  AND Year = 2021
  AND Month = 7;
