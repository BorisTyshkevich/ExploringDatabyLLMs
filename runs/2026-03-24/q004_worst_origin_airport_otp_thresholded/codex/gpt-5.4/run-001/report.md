# Worst origin airports by departure on-time performance

> Which airport ranks worst on departure on-time performance?

Among origin airports with at least 250,000 completed departures, Chicago Midway International (MDW) ranks worst on departure on-time performance. It logged 2,512,354 completed departures with 76.52% departure OTP, an average departure delay of 11.80 minutes, a 90th-percentile departure delay of 36.95 minutes, and coverage from 1987-10-01 through 2025-11-30.

- Rows returned: 15
- Columns: otp_rank, origin_code, DisplayAirportName, completed_departures, dep_otp_pct, avg_dep_delay_min, p90_dep_delay_min, first_date, last_date

| otp_rank | origin_code | DisplayAirportName | completed_departures | dep_otp_pct | avg_dep_delay_min | p90_dep_delay_min | first_date | last_date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | MDW | Chicago Midway International | 2.512354e+06 | 76.52 | 11.8 | 36.91 | 1987-10-01T00:00:00Z | 2025-11-30T00:00:00Z |

> How large is the spread between the worst airport and the middle of the ranked set?

Using the same 250,000-departure cutoff, Chicago Midway International (MDW) sits 8.35 percentage points below Richmond International (RIC), the median-ranked qualifying airport. MDW also runs 3.71 minutes worse on average departure delay and 8.55 minutes worse at the 90th percentile of departure delay.

- Rows returned: 1
- Columns: worst_code, worst_airport, worst_dep_otp_pct, middle_code, middle_airport, middle_dep_otp_pct, otp_gap_pct_points, worst_avg_dep_delay_min, middle_avg_dep_delay_min, avg_delay_gap_min, worst_p90_dep_delay_min, middle_p90_dep_delay_min, p90_delay_gap_min

| worst_code | worst_airport | worst_dep_otp_pct | middle_code | middle_airport | middle_dep_otp_pct | otp_gap_pct_points | worst_avg_dep_delay_min | middle_avg_dep_delay_min | avg_delay_gap_min | worst_p90_dep_delay_min | middle_p90_dep_delay_min | p90_delay_gap_min |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| MDW | Chicago Midway International | 76.52 | RIC | Richmond International | 84.87 | 8.35 | 11.8 | 8.09 | 3.71 | 36.93 | 28.39 | 8.54 |

> Are the weakest airports mostly major hubs, or is the bottom group more mixed?

The weakest qualifying airports are more mixed than purely hub-driven. The bottom dozen includes giant hubs such as ORD, EWR, JFK, DEN, SFO, and DFW, but also large secondary or focus airports such as MDW, HOU, DAL, BWI, FLL, and LAS, so the laggards are not confined to one airport type.

- Rows returned: 12
- Columns: otp_rank, origin_code, DisplayAirportName, completed_departures, traffic_rank, dep_otp_pct, avg_dep_delay_min, p90_dep_delay_min, first_date, last_date

| otp_rank | origin_code | DisplayAirportName | completed_departures | traffic_rank | dep_otp_pct | avg_dep_delay_min | p90_dep_delay_min | first_date | last_date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | MDW | Chicago Midway International | 2.512354e+06 | 27 | 76.52 | 11.8 | 36.9 | 1987-10-01T00:00:00Z | 2025-11-30T00:00:00Z |
