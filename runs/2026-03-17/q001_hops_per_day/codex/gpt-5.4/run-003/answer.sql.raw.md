```sql
WITH
    groupArray((
        addMinutes(addHours(toDateTime(FlightDate), intDiv(DepTime, 100)), DepTime % 100),
        concat(formatDateTime(addMinutes(addHours(toDateTime(FlightDate), intDiv(DepTime, 100)), DepTime % 100), '%H:%i'), ' ', Origin, '->', Dest)
    )) AS legs
SELECT
    Tail_Number AS `Aircraft ID`,
    Flight_Number_Reporting_Airline AS `Flight Number`,
    IATA_CODE_Reporting_Airline AS `Carrier`,
    FlightDate AS `Date`,
    arrayStringConcat(arrayMap(x -> x.2, arraySort(x -> x.1, legs)), ' | ') AS `Route`
FROM default.ontime_v2
WHERE Cancelled = 0
  AND Diverted = 0
  AND DepTime IS NOT NULL
  AND Tail_Number != ''
  AND Flight_Number_Reporting_Airline != ''
GROUP BY
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate
ORDER BY
    count() DESC,
    FlightDate DESC,
    max(addMinutes(addHours(toDateTime(FlightDate), intDiv(DepTime, 100)), DepTime % 100)) DESC
LIMIT 10
```