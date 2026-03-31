# American Airlines peak network delay month and contributors

> Which month is the single worst American Airlines month for departure delays?

July 2024 is the single worst American Airlines month for departure delays in the full history. It had 86,083 completed AA flights, an average departure delay of 36.33 minutes, and 38.04% of flights departed at least 15 minutes late.

- Rows returned: 458
- Columns: month, flight_volume, avg_departure_delay_minutes, pct_departing_15_plus_late, worst_month_rank

| month | flight_volume | avg_departure_delay_minutes | pct_departing_15_plus_late | worst_month_rank |
| --- | --- | --- | --- | --- |
| 1987-10-01T00:00:00Z | 55871 | 4.27 | 8.8 | 402 |

> Which origins contribute most to that peak month?

In July 2024, DFW and CLT were the dominant origin contributors. DFW generated the most total departure delay minutes, while CLT was close behind and had the worse average delay rate; ORD, MIA, and PHL formed the next tier of meaningful contributors.

- Rows returned: 15
- Columns: OriginCode, flight_volume, total_departure_delay_minutes, avg_departure_delay_minutes, pct_departing_15_plus_late

| OriginCode | flight_volume | total_departure_delay_minutes | avg_departure_delay_minutes | pct_departing_15_plus_late |
| --- | --- | --- | --- | --- |
| DFW | 14962 | 593021 | 38.52 | 46.59 |

> Which routes contribute most to that peak month?

The largest July 2024 route contributor was DFW-LAX, followed very closely by CLT-MCO and DFW-SAT. CLT-RDU, RDU-CLT, and DFW-PHX also ranked among the strongest route-level contributors to the bad month.

- Rows returned: 20
- Columns: route, flight_volume, total_departure_delay_minutes, avg_departure_delay_minutes, pct_departing_15_plus_late

| route | flight_volume | total_departure_delay_minutes | avg_departure_delay_minutes | pct_departing_15_plus_late |
| --- | --- | --- | --- | --- |
| DFW-LAX | 443 | 18462 | 40.89 | 46.5 |

> Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

July 2024 looks broad across the network, but with clear hub concentration. AA had delays spread across 125 origins and 909 routes that month; the top two origins, DFW and CLT, accounted for 50.8% of meaningful origin delay minutes and the top five origins reached 73.71%, while the top ten routes made up only 15.82% of meaningful route delay minutes.

- Rows returned: 1
- Columns: total_origins_in_month, total_routes_in_month, top_2_origins_share_pct, top_5_origins_share_pct, top_10_routes_share_pct

| total_origins_in_month | total_routes_in_month | top_2_origins_share_pct | top_5_origins_share_pct | top_10_routes_share_pct |
| --- | --- | --- | --- | --- |
| 125 | 909 | 50.8 | 73.71 | 15.82 |
