Find American Airlines' worst network-wide month for departure delays, then identify which origins and routes contributed most to that peak.

Definitions and filters:

- Filter to `IATA_CODE_Reporting_Airline = 'AA'`.
- Restrict to completed flights with `Cancelled = 0`.
- Aggregate months using `toStartOfMonth(FlightDate)`.

Monthly metrics:

- completed flights
- average `DepDelayMinutes`
- percentage of flights with `DepDel15 = 1`

Peak-month ranking:

- Rank months by average `DepDelayMinutes` descending, then `% DepDel15` descending, then completed flights descending, then `MonthStart` ascending so ties prefer the earlier month.
- Select the single worst month using those rules.

Contribution drilldowns for the selected month:

- At `Origin` level, compute completed flights, average `DepDelayMinutes`, and total departure delay minutes.
- At `(Origin, Dest)` route level, compute completed flights, average `DepDelayMinutes`, and total departure delay minutes.
- Only include origins or routes with at least `100` completed flights in the selected month.
- Rank contributors by total departure delay minutes descending, then average `DepDelayMinutes` descending, then completed flights descending.

Required output:

- Your SQL should return a single result set that supports both peak-month selection and peak-month drilldowns.
- Include these columns in this order:
  `RowType`,
  `MonthStart`,
  `Origin`,
  `Dest`,
  `CompletedFlights`,
  `AvgDepDelayMinutes`,
  `DepDel15Pct`,
  `TotalDepDelayMinutes`,
  `MonthRank`,
  `ContributionRank`,
  `IsPeakMonth`
- Use `RowType = 'month_summary'` for all monthly leaderboard rows.
- Use `RowType = 'peak_origin'` for top origin contributors in the selected peak month.
- Use `RowType = 'peak_route'` for top `(Origin, Dest)` contributors in the selected peak month.
- For `month_summary` rows, leave `Origin`, `Dest`, and `ContributionRank` empty.
- For drilldown rows, keep `IsPeakMonth = 1`.

Implementation expectations:

- Use one query that emits the monthly summary and the selected-month contribution rows together.
- Make the peak-month selection deterministic.

Report guidance:

Explain:

- which month is the single worst American Airlines month by the question's ranking rule,
- which origins contribute most to that peak month,
- which routes contribute most to that peak month,
- and whether the peak appears broad-based across the network or concentrated in a smaller set of origins and routes.
