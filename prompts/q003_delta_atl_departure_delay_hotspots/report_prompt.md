Find which Delta departures out of ATL have the worst sustained departure delays at the `(Dest, DepTimeBlk)` level.

Definitions and filters:

- Filter to `IATA_CODE_Reporting_Airline = 'DL'`.
- Filter to `Origin = 'ATL'`.
- Restrict to completed flights with `Cancelled = 0`.
- Use `FlightDate` truncated to month as the monthly grain.

Metrics to compute:

- `CompletedFlights`
- average `DepDelayMinutes`
- p90 `DepDelayMinutes`
- percentage of flights delayed 15+ minutes using `DepDel15`
- number of qualifying months

Metric semantics:

- Use `quantile(0.9)(DepDelayMinutes)` for p90.
- After monthly qualification, recompute final hotspot-level metrics from the raw flights that belong to qualifying monthly cells.
- Do not compute hotspot metrics by averaging monthly averages, monthly p90 values, or monthly delayed-15 percentages.

Threshold rules:

- A monthly cell `(month, Dest, DepTimeBlk)` qualifies only if it has at least `40` completed flights.
- A `(Dest, DepTimeBlk)` hotspot qualifies only if it has at least `1,000` completed flights across all qualifying monthly cells.
- “Worst sustained” means rank by average `DepDelayMinutes` descending, then p90 `DepDelayMinutes` descending, then `% DepDel15` descending, then completed flights descending.

Return one SQL query that produces rows supporting both hotspot ranking and monthly trend analysis.

Include these columns in this order:

- `RowType`
- `MonthStart`
- `Dest`
- `DepTimeBlk`
- `QualifyingMonths`
- `CompletedFlights`
- `AvgDepDelayMinutes`
- `P90DepDelayMinutes`
- `DepDel15Pct`
- `FirstQualifyingMonth`
- `LastQualifyingMonth`
- `HotspotRank`

Row rules:

- Use `RowType = 'hotspot_summary'` for the top 20 final hotspot cells.
- For `hotspot_summary` rows, set `MonthStart = NULL`.
- Use `RowType = 'monthly_trend'` for monthly rows belonging to those top 20 hotspot cells after monthly qualification.
- Populate `FirstQualifyingMonth` and `LastQualifyingMonth` for both row types.

Numeric normalization:

- Round `AvgDepDelayMinutes`, `P90DepDelayMinutes`, and `DepDel15Pct` to exactly 2 decimal places in the final output.

Ordering:

- Sort by `RowType`, then `HotspotRank` ascending, then `MonthStart` ascending, then `Dest`, then `DepTimeBlk`.

Implementation expectations:

- Use CTEs for monthly qualification, final rollup, and extraction of monthly trend rows for the top-ranked hotspot cells.
- Exclude low-volume monthly cells before final ranking.
- Verify that the final top-20 `(Dest, DepTimeBlk)` ranking is based on hotspot metrics recomputed over all qualifying raw flights.

Report guidance:

Explain:

- which destination and time block is the worst hotspot,
- whether the hotspot is persistent or concentrated in a narrower era,
- and what the top 5 hotspot cells suggest about ATL departure pressure for Delta.
