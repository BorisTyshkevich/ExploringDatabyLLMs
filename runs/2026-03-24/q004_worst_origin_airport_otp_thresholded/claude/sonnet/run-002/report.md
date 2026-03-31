# Worst origin airports by departure on-time performance

> Which airport ranks worst on departure on-time performance?

Chicago Midway International (MDW) ranks worst among all qualifying origin airports, with an on-time departure rate of 76.52% — meaning nearly 1 in 4 departures was delayed 15 or more minutes. MDW logged 2.51 million completed departures across the full dataset (1987–2025) and, when delayed, averaged 26.8 minutes late with a 90th-percentile delay of 37 minutes. The minimum-volume threshold applied was 100,000 completed departures, leaving 160 qualifying airports in the ranked set.

- Rows returned: 161
- Columns: rank, OriginCode, airport_name, city, total_departures, otp_rate, pct_delayed, avg_dep_delay_when_late, p90_delay_min, first_date, last_date

| rank | OriginCode | airport_name | city | total_departures | otp_rate | pct_delayed | avg_dep_delay_when_late | p90_delay_min | first_date | last_date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | MDW | Chicago Midway International | Chicago, IL | 2.512354e+06 | 76.52 | 23.48 | 26.82 | 37 | 1987-10-01T00:00:00Z | 2025-11-30T00:00:00Z |

> How large is the spread between the worst airport and the middle of the ranked set?

The spread between the worst airport and the median of the 160 qualifying airports is 8.53 percentage points. MDW sits at 76.52% OTP while the median airport in the ranked set sits at 85.05% OTP. This gap is substantial: MDW's delay rate (23.48%) is roughly 1.6x the median delay rate (14.95%), indicating that the very bottom of the distribution is a clear outlier rather than a gradual tail.

- Rows returned: 1
- Columns: worst_otp, median_otp, spread_pp, qualifying_airports

| worst_otp | median_otp | spread_pp | qualifying_airports |
| --- | --- | --- | --- |
| 76.52 | 85.05 | 8.53 | 160 |

> Are the weakest airports mostly major hubs, or is the bottom group more mixed?

The bottom 10 airports are overwhelmingly major hubs and high-traffic gateways. MDW (Chicago Midway), ORD (Chicago O'Hare), EWR (Newark Liberty), JFK (John F. Kennedy), DEN (Denver), and LAS (Las Vegas) are among the nation's busiest airports by operations. HOU (Houston Hobby) and DAL (Dallas Love Field) are primary Southwest Airlines intra-Texas hubs. BWI (Baltimore/Washington) and FLL (Fort Lauderdale) are large leisure and low-cost-carrier gateways. Every airport in the bottom 10 carries well over 1 million historical departures, and none are regional or small-market airports. The pattern strongly suggests that operational complexity, congestion, and constrained airspace — not obscurity — drive the worst OTP outcomes.

- Rows returned: 10
- Columns: OriginCode, airport_name, city, total_departures, otp_rate, pct_delayed, avg_dep_delay_when_late, p90_delay_min

| OriginCode | airport_name | city | total_departures | otp_rate | pct_delayed | avg_dep_delay_when_late | p90_delay_min |
| --- | --- | --- | --- | --- | --- | --- | --- |
| MDW | Chicago Midway International | Chicago, IL | 2.512354e+06 | 76.52 | 23.48 | 26.82 | 37 |
