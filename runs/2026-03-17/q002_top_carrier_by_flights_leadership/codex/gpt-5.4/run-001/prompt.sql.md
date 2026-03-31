You are running inside qforge.

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

Determine which `Reporting_Airline` flew the most completed flights in each calendar year, and identify where leadership changed most sharply.

Definitions and rules:

- A completed flight is a row with `Cancelled = 0`.
- Group first by `Year` and `Reporting_Airline`.
- For each year, compute completed flights and carrier share of that year's completed flights.
- Rank carriers within each year by completed flights descending.
- The annual leader is rank 1; the runner-up is rank 2.
- Compute the leader share gap versus the runner-up in percentage points.
- Compute the leader's year-over-year share change using a window function over yearly leaders only.
- A leadership change means the annual leader differs from the previous year's leader.
- The first year has no prior year for comparison, so it is not a leadership change.
- “Most sharply” refers only to years where a leadership change occurred.
- Order leadership-change years by absolute value of the leader's year-over-year share change descending, then by larger leader share gap, then by later year.

Required output:

- Your SQL should return a single result set that supports both the annual summary and the dashboard.
- Include these columns in this order:
  `RowType`,
  `Year`,
  `Reporting_Airline`,
  `RankInYear`,
  `CompletedFlights`,
  `SharePct`,
  `LeaderReportingAirline`,
  `RunnerUpReportingAirline`,
  `LeaderShareGapPctPts`,
  `PriorYearLeaderReportingAirline`,
  `LeaderChanged`,
  `LeaderShareChangePctPts`
- Use `RowType = 'carrier_year'` for ranked carrier-year rows.
- Include at least the top 5 carriers per year as `carrier_year` rows so the visual can build a bump chart and share chart.
- Repeat the leader-transition fields on the rank-1 row for each year.
- For the first year, set `PriorYearLeaderReportingAirline` and `LeaderShareChangePctPts` to `NULL` or an empty value, and set `LeaderChanged = 0`.
- If you also emit a pure yearly summary row, label it `RowType = 'year_summary'`.

Ordering:

- Sort by `Year` ascending, then `RowType`, then `RankInYear` ascending, then `Reporting_Airline`.

Implementation expectations:

- Use CTEs to separate annual totals, annual carrier counts, per-year ranking, and yearly leader transitions.
- Make the result deterministic.