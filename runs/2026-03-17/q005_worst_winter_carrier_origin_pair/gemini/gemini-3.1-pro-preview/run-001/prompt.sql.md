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

Determine which `(Reporting_Airline, Origin)` pairs perform worst in winter after applying a minimum flight threshold.

Definitions and filters:

- Winter is `Month IN (12, 1, 2)`.
- Restrict to completed flights with `Cancelled = 0`.
- Aggregate at `(Reporting_Airline, Origin)`.
- A qualifying pair must have at least `5,000` completed winter departures over the full table history.

Ranking metrics:

- primary: departure OTP percentage, where OTP is the share with `DepDel15 = 0`
- secondary: average `DepDelayMinutes`
- tertiary: completed winter departures

Cause context:

- For each qualifying pair, compute total non-null minutes for `CarrierDelay`, `WeatherDelay`, `NASDelay`, `SecurityDelay`, and `LateAircraftDelay`.
- Convert those totals into percentage shares of total reported delay-cause minutes.
- If total reported delay-cause minutes are zero for a pair, emit zero shares rather than `NULL` or divide-by-zero results.
- Ranking must still be based on OTP and average departure delay, not cause shares.

Required output:

- Return the 20 worst qualifying winter pairs.
- Include these columns in this order:
  `Reporting_Airline`,
  `Origin`,
  `OriginCityName`,
  `CompletedWinterDepartures`,
  `DepartureOtpPct`,
  `AvgDepDelayMinutes`,
  `CarrierDelaySharePct`,
  `WeatherDelaySharePct`,
  `NASDelaySharePct`,
  `SecurityDelaySharePct`,
  `LateAircraftDelaySharePct`

Ordering:

- Sort by departure OTP ascending, then average `DepDelayMinutes` descending, then completed winter departures descending, then `Reporting_Airline`, then `Origin`.

Implementation expectations:

- Use separate CTEs for winter performance and delay-cause decomposition.
- Handle missing cause minutes safely.
- Keep season logic explicit in SQL.