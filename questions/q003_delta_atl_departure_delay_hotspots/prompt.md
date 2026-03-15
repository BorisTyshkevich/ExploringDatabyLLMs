Using `default.ontime_v2`, find which Delta departures out of ATL have the worst sustained departure delays at the `(Dest, DepTimeBlk)` level.

Definitions and filters:

- Use only `default.ontime_v2`.
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
- Use `RowType = 'monthly_trend'` for monthly rows belonging to those top 20 hotspot cells after monthly qualification.

Ordering:

- Sort by `RowType`, then `HotspotRank` ascending, then `MonthStart` ascending, then `Dest`, then `DepTimeBlk`.

Implementation expectations:

- Use CTEs for monthly qualification, final rollup, and extraction of monthly trend rows for the top-ranked hotspot cells.
- Exclude low-volume monthly cells before final ranking.
- Keep field names exactly from `default.ontime_v2`.
