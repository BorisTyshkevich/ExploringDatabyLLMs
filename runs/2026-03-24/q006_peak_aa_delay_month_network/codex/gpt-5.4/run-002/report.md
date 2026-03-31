# American Airlines peak network delay month and contributors

> Which month is the single worst American Airlines month for departure delays?

July 2024 is the single worst American Airlines month for departure delays in the full available history. It leads the AA monthly leaderboard on both average departure delay and the share of flights leaving 15 or more minutes late, with 86,083 completed flights, a 36.33 minute average departure delay, and 38.04% of flights departing 15+ minutes late.

- Rows returned: 458
- Columns: month, flights, avg_dep_delay, pct_dep_15_plus, worst_month_rank

| month | flights | avg_dep_delay | pct_dep_15_plus | worst_month_rank |
| --- | --- | --- | --- | --- |
| 1987-10-01T00:00:00Z | 55871 | 4.27 | 8.8 | 402 |

> Which origins contribute most to that peak month?

In July 2024, Dallas/Fort Worth (DFW) and Charlotte (CLT) are the dominant origin contributors. DFW produced 593,021 departure-delay minutes and 17.98% of the month's total AA delay minutes, while CLT added 560,409 minutes and 16.99%; ORD, MIA, and PHL are the next tier. The displayed origin table uses a 1,000-flight readability filter, but each airport's share is still measured against the full July 2024 AA network.

- Rows returned: 15
- Columns: OriginCode, origin_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, network_delay_share_pct, network_flight_share_pct

| OriginCode | origin_city | flights | avg_dep_delay | pct_dep_15_plus | total_dep_delay_minutes | network_delay_share_pct | network_flight_share_pct |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DFW | Dallas/Fort Worth, TX | 14962 | 38.52 | 46.59 | 593021 | 17.98 | 17.38 |

> Which routes contribute most to that peak month?

Among business-meaningful July 2024 AA routes, DFW-LAX is the largest route contributor by departure-delay minutes. CLT-MCO and DFW-SAT are close behind, followed by CLT-RDU and RDU-CLT. The displayed route table applies a 150-flight readability filter, but every route share is still computed against the full July 2024 AA network total.

- Rows returned: 182
- Columns: route, origin_city, dest_city, flights, avg_dep_delay, pct_dep_15_plus, total_dep_delay_minutes, network_delay_share_pct, network_flight_share_pct

| route | origin_city | dest_city | flights | avg_dep_delay | pct_dep_15_plus | total_dep_delay_minutes | network_delay_share_pct | network_flight_share_pct |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DFW-LAX | Dallas/Fort Worth, TX | Los Angeles, CA | 443 | 40.89 | 46.5 | 18462 | 0.56 | 0.51 |

> Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

July 2024 looks broad across the AA network, with notable hub concentration rather than a single-point failure. DFW alone accounts for 17.98% of delay minutes, and the top five origins account for 50.75%, but the top route contributes only 0.56% and the top ten routes together only 4.88% of delay minutes across 909 routes. That pattern points to network-wide disruption centered heavily in major hubs, especially DFW and CLT, not a narrow problem confined to a few routes.

- Rows returned: 1
- Columns: network_flights, network_delay_minutes, top_origin, top_origin_delay_minutes, top_origin_delay_share_pct, top_origin_flights, top_origin_flight_share_pct, top_route, top_route_delay_minutes, top_route_delay_share_pct, top_route_flights, top_route_flight_share_pct, top5_origin_delay_minutes, top5_origin_delay_share_pct, top5_origin_flights, top5_origin_flight_share_pct, top10_route_delay_minutes, top10_route_delay_share_pct, top10_route_flights, top10_route_flight_share_pct, origin_count, route_count

| network_flights | network_delay_minutes | top_origin | top_origin_delay_minutes | top_origin_delay_share_pct | top_origin_flights | top_origin_flight_share_pct | top_route | top_route_delay_minutes | top_route_delay_share_pct | top_route_flights | top_route_flight_share_pct | top5_origin_delay_minutes | top5_origin_delay_share_pct | top5_origin_flights | top5_origin_flight_share_pct | top10_route_delay_minutes | top10_route_delay_share_pct | top10_route_flights | top10_route_flight_share_pct | origin_count | route_count |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 86083 | 3.297804e+06 | DFW | 593021 | 17.98 | 14962 | 17.38 | DFW-LAX | 18462 | 0.56 | 443 | 0.51 | 1.673625e+06 | 50.75 | 39810 | 46.25 | 161038 | 4.88 | 3340 | 3.88 | 125 | 909 |
