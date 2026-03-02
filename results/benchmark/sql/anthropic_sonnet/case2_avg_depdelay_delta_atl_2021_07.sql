SELECT round(avg(DepDelay)) AS avg_dep_delay
FROM ontime
WHERE Year = 2021
  AND Month = 7
  AND trimRight(toString(Carrier)) = 'DL'
  AND trimRight(toString(Origin)) = 'ATL'
