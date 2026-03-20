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

Report guidance:

Explain:

- which winter `(Reporting_Airline, Origin)` pair ranks worst overall,
- whether the worst-ranked pairs are driven more by weather or by operational causes such as carrier, NAS, or late-aircraft delay,
- and whether the worst pairs are concentrated within a small number of carriers or airports.
