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

Report guidance:

Explain:

- which airport ranks worst on departure on-time performance and how poor its OTP is,
- how wide the spread is between the worst airport and the median airport within the ranked set,
- and whether the bottom 25 are dominated by major hubs or show a more mixed airport profile.
