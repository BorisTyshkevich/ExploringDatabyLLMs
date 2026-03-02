SELECT round(avg(DepDelay), 2) AS avg_dep_delay FROM default.ontime WHERE Year = 2021 AND Month = 7 AND trimRight(toString(Carrier)) = 'DL' AND trimRight(toString(Origin)) = 'ATL'
