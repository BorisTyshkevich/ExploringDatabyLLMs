# Worst winter carrier-origin pairs by departure performance

> Which winter carrier-airport pair ranks worst overall?

DH (Independence Air) departing from ORD (Chicago O'Hare) ranks worst among all qualifying winter carrier-origin pairs with at least 1,000 winter flights. Across 19,986 winter departures, it achieved only 56.58% on-time performance and averaged 28.35 minutes of departure delay per completed flight — both the worst figures in the dataset by a meaningful margin. The next-worst qualifying pairs (PI/DFW at 61.5%, PI/LAX at 62.95%) trail DH/ORD by more than four percentage points.

- Rows returned: 20
- Columns: carrier, origin, winter_flights, otp_pct, avg_dep_delay_min

| carrier | origin | winter_flights | otp_pct | avg_dep_delay_min |
| --- | --- | --- | --- | --- |
| DH | ORD | 19986 | 56.58 | 28.35 |

> Are the worst pairs driven more by weather or by operational causes?

Operational causes dominate decisively across all worst pairs that have delay-cause data. For DH/ORD the breakdown is 15.6% weather vs 84.4% operational (carrier delay + NAS + late aircraft + security). YV/ORD is even more skewed: 5.0% weather vs 95.0% operational. AA/EYW sits at 5.7% weather vs 94.3% operational; F9/PBI at 1.5% vs 98.5%; OO/OTH at 2.9% vs 97.1%. OH/RSW is the most weather-affected at 26.8% weather vs 73.2% operational. In every case, operational causes account for the majority of reported delay minutes, making poor operational execution rather than weather the primary driver of chronic winter under-performance.

- Rows returned: 10
- Columns: carrier, origin, winter_flights, otp_pct, avg_dep_delay_min, weather_min, carrier_min, nas_min, security_min, late_aircraft_min, total_reported_delay, weather_share_pct, operational_share_pct

| carrier | origin | winter_flights | otp_pct | avg_dep_delay_min | weather_min | carrier_min | nas_min | security_min | late_aircraft_min | total_reported_delay | weather_share_pct | operational_share_pct |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DH | ORD | 19986 | 56.58 | 28.35 | 55176 | 123407 | 41492 | 116 | 133239 | 353430 | 15.6 | 84.4 |

> Are the weakest pairs concentrated in a small number of carriers or airports?

The weakest pairs are moderately spread rather than tightly concentrated. Among the 30 worst qualifying winter pairs, B6 (JetBlue) is the most frequent carrier with 6 pairs spread across MIA, RNO, SRQ, FLL, PHX, and PBI. PI, OO, and OH each contribute 3 pairs. YV, F9, EV, and FL each have 2 pairs, while DH, AS, AA, NW, G4, XE, and EA appear once each — 15 distinct carriers in total. On the airport side, ORD and EWR each host 2 pairs, FLL appears twice (B6 and F9), and no single airport dominates the list. The weakness is carrier-driven more than airport-driven: B6's Florida leisure-market airports and regional carriers at high-congestion or weather-prone hubs account for most of the bottom tier, but there is no single choke-point airport that concentrates failure.

- Rows returned: 15
- Columns: carrier, pairs_in_worst30, avg_otp_pct, origins

| carrier | pairs_in_worst30 | avg_otp_pct | origins |
| --- | --- | --- | --- |
| B6 | 6 | 68.43 | [MIA, RNO, SRQ, FLL, PHX, PBI] |
