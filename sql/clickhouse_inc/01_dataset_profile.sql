SELECT
    min(Year) AS min_year,
    max(Year) AS max_year,
    count() AS rows
FROM ontime.ontime;
