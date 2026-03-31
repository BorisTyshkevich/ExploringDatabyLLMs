You are running inside qforge.

Generate valid ClickHouse SQL.
Use the configured MCP server for all data access.
Do not construct raw OpenAPI URLs manually.
Stay within the configured dataset scope.
Do not reference tables outside the allowed dataset constraints.

Dataset constraints:

- Use `default.ontime_v2` as the primary fact table.
- Do not reference `default.ontime`.

Do not emit result rows or any data output.
Inspect the live schema first with `DESCRIBE TABLE default.ontime_v2`.
Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in source data reading subquery/CTE, and fix any errors internally.
Return only the final full SQL query, not the debug query.

Return exactly this fenced section:

```sql
-- one SQL statement
```

Rules:

- Use one SQL statement only.
- Emit only the final verified SQL.
- Do not include TSV, JSON rows, HTML, report prose, or any other fenced blocks.

Identify which origin airports have the worst departure on-time performance after excluding low-volume airports.

Definitions and rules:

- Restrict to completed flights with `Cancelled = 0`.
- Departure on-time performance is the share of flights with `DepDel15 = 0`.
- Aggregate at the `Origin` level.
- A qualifying airport must have at least `50,000` completed departures over the full table history.
- Rank airports by departure OTP ascending, then average `DepDelayMinutes` descending, then completed departures descending, then `Origin` ascending.

Required metrics:

- completed departures
- departure OTP percentage
- average `DepDelayMinutes`
- p90 `DepDelayMinutes`
- first flight date
- last flight date

Required output:

- Return the 25 worst qualifying origin airports.
- Include these columns in this order:
  `Origin`,
  `OriginCityName`,
  `OriginState`,
  `CompletedDepartures`,
  `DepartureOtpPct`,
  `AvgDepDelayMinutes`,
  `P90DepDelayMinutes`,
  `FirstFlightDate`,
  `LastFlightDate`

Ordering:

- Sort using the ranking rules above.

Implementation expectations:

- Use a threshold CTE before final ranking.
- Keep the ranking deterministic.
- Use `quantile` logic for p90 rather than approximating in prose.