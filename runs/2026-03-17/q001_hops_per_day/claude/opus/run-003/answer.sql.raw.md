The query works correctly. Now let me finalize it with proper column aliases as requested (Aircraft ID instead of Tail_Number, etc.).

```sql
WITH flight_legs AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Origin,
        Dest,
        DepTime
    FROM default.ontime_v2
    WHERE Tail_Number != ''
      AND Cancelled = 0
      AND DepTime IS NOT NULL
),
grouped AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arraySort(x -> x.1, groupArray(tuple(DepTime, Origin, Dest))) AS legs_sorted
    FROM flight_legs
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate
    HAVING Hops >= 2
)
SELECT
    Tail_Number AS `Aircraft ID`,
    Flight_Number_Reporting_Airline AS `Flight Number`,
    IATA_CODE_Reporting_Airline AS Carrier,
    FlightDate AS `Date`,
    Hops,
    arrayStringConcat(
        arrayMap(x -> concat(
            lpad(toString(x.1 DIV 100), 2, '0'), ':', lpad(toString(x.1 MOD 100), 2, '0'),
            ' ', x.2
        ), legs_sorted),
        ' → '
    ) || ' → ' || legs_sorted[length(legs_sorted)].3 AS Route
FROM grouped
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10
```
