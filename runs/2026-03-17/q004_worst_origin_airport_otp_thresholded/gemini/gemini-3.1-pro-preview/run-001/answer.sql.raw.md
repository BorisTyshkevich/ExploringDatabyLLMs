```sql
WITH OriginStats AS (
    SELECT
        Origin,
        any(OriginCityName) AS OriginCityName,
        any(OriginState) AS OriginState,
        count() AS CompletedDepartures,
        countIf(DepDel15 = 0) / count() AS DepartureOtpPct,
        avg(DepDelayMinutes) AS AvgDepDelayMinutes,
        quantile(0.90)(DepDelayMinutes) AS P90DepDelayMinutes,
        min(FlightDate) AS FirstFlightDate,
        max(FlightDate) AS LastFlightDate
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Origin
    HAVING CompletedDepartures >= 50000
)
SELECT
    Origin,
    OriginCityName,
    OriginState,
    CompletedDepartures,
    DepartureOtpPct,
    AvgDepDelayMinutes,
    P90DepDelayMinutes,
    FirstFlightDate,
    LastFlightDate
FROM OriginStats
ORDER BY
    DepartureOtpPct ASC,
    AvgDepDelayMinutes DESC,
    CompletedDepartures DESC,
    Origin ASC
LIMIT 25;
```
