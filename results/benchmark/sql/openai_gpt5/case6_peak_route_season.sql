SELECT route, season, avg_dep_delay, flights
FROM
(
  SELECT route, season, avg(DepDelay) AS avg_dep_delay, count() AS flights
  FROM
  (
    SELECT
      concat(trimRight(toString(Origin)),'-',trimRight(toString(Dest))) AS route,
      multiIf(Month IN (12,1,2), 'Winter', Month IN (3,4,5), 'Spring', Month IN (6,7,8), 'Summer', 'Fall') AS season,
      DepDelay
    FROM default.ontime
    WHERE Year BETWEEN 2019 AND 2021
  )
  WHERE route IN
  (
    SELECT route FROM
    (
      SELECT
        concat(trimRight(toString(Origin)),'-',trimRight(toString(Dest))) AS route,
        count() AS flights
      FROM default.ontime
      WHERE Year BETWEEN 2019 AND 2021
      GROUP BY route
      ORDER BY flights DESC
      LIMIT 50
    )
  )
  GROUP BY route, season
)
ORDER BY avg_dep_delay DESC, flights DESC
LIMIT 1
