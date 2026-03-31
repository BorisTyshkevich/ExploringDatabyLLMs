# Worst origin airports by departure on-time performance

> Which airport ranks worst on departure on-time performance?

Jack McNamara Field (CEC, Crescent City, CA) ranks worst among the 304 qualifying origin airports (those with at least 10,000 completed departures). It has an OTP of 69.24%, meaning nearly 1 in 3 departures is delayed 15 or more minutes. Its average departure delay is 19.5 minutes and the 90th-percentile delay is 84 minutes — the highest in the dataset. The next-worst airports are Aspen Pitkin County Sardy Field (ASE, 75.09%) and Nantucket Memorial (ACK, 75.77%). The first major hub to appear in the ranking is Chicago Midway International (MDW, rank 4, 76.52%), followed by Chicago O'Hare (ORD, rank 9, 78.61%) and Newark Liberty (EWR, rank 11, 78.69%).

- Rows returned: 30
- Columns: rank_worst, total_qualifying, OriginCode, airport_name, completed_departures, otp_pct, avg_dep_delay_min, p90_dep_delay_min, first_date, last_date

| rank_worst | total_qualifying | OriginCode | airport_name | completed_departures | otp_pct | avg_dep_delay_min | p90_dep_delay_min | first_date | last_date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 304 | CEC | Jack McNamara Field | 11891 | 69.24 | 19.5 | 83 | 2003-01-01T00:00:00Z | 2015-04-06T00:00:00Z |

> How large is the spread between the worst airport and the middle of the ranked set?

Among 304 qualifying airports, the middle-ranked airport (rank 152) is Pensacola International (PNS) with an OTP of 85.31%. The worst airport, Jack McNamara Field (CEC), has an OTP of 69.24%. The spread is 16.07 percentage points — meaning CEC departs on time 16 points less often than a typical mid-tier airport. This is a substantial gap, reflecting that CEC's performance is genuinely extreme rather than marginally below average.

- Rows returned: 2
- Columns: rank_worst, total_qualifying, OriginCode, airport_name, completed_departures, otp_pct, avg_dep_delay_min, p90_dep_delay_min, spread_from_worst_pct

| rank_worst | total_qualifying | OriginCode | airport_name | completed_departures | otp_pct | avg_dep_delay_min | p90_dep_delay_min | spread_from_worst_pct |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 304 | CEC | Jack McNamara Field | 11891 | 69.24 | 19.5 | 84 | 0 |

> Are the weakest airports mostly major hubs, or is the bottom group more mixed?

The bottom group is genuinely mixed. Among the 30 worst-performing qualifying airports, small and mid-size regional airports dominate the very bottom — Jack McNamara Field (CEC, 11k departures), Aspen Pitkin County Sardy Field (ASE, 99k), Nantucket Memorial (ACK, 15k), Modesto City-County (MOD, 18k), and Trenton Mercer (TTN, 28k) all rank in the worst 10. However, some of the busiest hubs in the country also appear throughout the bottom 30: Chicago Midway (MDW, 2.5M, rank 4), Chicago O'Hare (ORD, 11.1M, rank 9), Newark Liberty (EWR, 4.5M, rank 11), JFK (3.1M, rank 17), Baltimore/Washington (BWI, 3.3M, rank 18), Las Vegas (LAS, 5.2M, rank 20), Denver (DEN, 7.4M, rank 22), Fort Lauderdale (FLL, 2.3M, rank 23), San Francisco (SFO, 5.1M, rank 26), Dallas/Fort Worth (DFW, 10M, rank 27), Miami (MIA, 2.8M, rank 29), and Philadelphia (PHL, 3.6M, rank 30). The conclusion is that poor departure OTP is not limited to major hubs — the weakest performers span the full volume spectrum.

- Rows returned: 30
- Columns: rank_worst, OriginCode, airport_name, completed_departures, volume_tier, otp_pct, avg_dep_delay_min, p90_dep_delay_min

| rank_worst | OriginCode | airport_name | completed_departures | volume_tier | otp_pct | avg_dep_delay_min | p90_dep_delay_min |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | CEC | Jack McNamara Field | 11891 | Regional/small (<100k dep) | 69.24 | 19.5 | 83 |
