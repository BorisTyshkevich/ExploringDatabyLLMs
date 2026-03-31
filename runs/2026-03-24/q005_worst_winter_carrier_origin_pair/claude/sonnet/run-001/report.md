# Worst winter carrier-origin pairs by departure performance

> Which winter carrier-airport pair ranks worst overall?

The worst winter carrier-airport pair is **DH (Independence Air) departing from ORD (Chicago O'Hare)**, with only 56.58% of winter flights departing on time and an average departure delay of 28.35 minutes across 19,986 qualifying winter flights. The second-worst pair is PI (Piedmont Airlines) at DFW (61.5% on-time), followed by PI at LAX (62.95%) and PI at DAY (63.15%). The ranked query below returns the full set of worst-performing winter pairs (minimum 1,000 winter flights) sorted by on-time percentage ascending.

- Rows returned: 25
- Columns: carrier, origin, winter_flights, pct_ontime, avg_dep_delay_min, weather_pct, operational_pct, rank

| carrier | origin | winter_flights | pct_ontime | avg_dep_delay_min | weather_pct | operational_pct | rank |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DH | ORD | 19986 | 56.58 | 28.35 | 15.61 | 84.39 | 1 |

> Are the worst pairs driven more by weather or by operational causes?

The worst winter pairs are overwhelmingly driven by **operational causes**, not weather. Across the 22 worst pairs with delay-cause data, operational delays (carrier delays + late aircraft + NAS + security) account for 70–99% of reported delay minutes in nearly every case. Late aircraft propagation alone is the single largest component for most pairs — for example, DH/ORD (37.7% late aircraft), YV/ORD (45.2% late aircraft), AA/EYW (61.1% late aircraft), and AS/DUT (66.0% late aircraft). Carrier-attributed delays are the second-largest factor. Weather is a minor contributor for most pairs (typically 1–15%), with only two notable exceptions: OH/EWR (41.4% weather) and OH/ATL (28.7% weather). Three PI pairs show no delay-cause data, consistent with older data predating structured cause reporting.

- Rows returned: 25
- Columns: carrier, origin, pct_ontime, total_delay_min, weather_share_pct, carrier_share_pct, nas_share_pct, late_aircraft_share_pct, security_share_pct

| carrier | origin | pct_ontime | total_delay_min | weather_share_pct | carrier_share_pct | nas_share_pct | late_aircraft_share_pct | security_share_pct |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DH | ORD | 56.58 | 353430 | 15.61 | 34.92 | 11.74 | 37.7 | 0.03 |

> Are the weakest pairs concentrated in a small number of carriers or airports?

The 25 worst winter pairs span **14 distinct carriers and 22 distinct airports**, so the weak set is fairly dispersed across airports but shows moderate carrier concentration. Four carriers account for the majority of repeating entries: **B6 (JetBlue)** appears 4 times (MIA, RNO, SRQ, FLL), **OH (Comair)** 3 times (RSW, EWR, ATL), **OO (SkyWest)** 3 times (OTH, ASE, CEC), and **PI (Piedmont)** 3 times (DFW, LAX, DAY). On the airport side, **EWR** appears in 3 pairs (OH, FL, XE) and **ORD** in 2 pairs (DH, YV). The pattern suggests that chronic winter departure underperformance clusters around a handful of regional and low-cost carriers operating systemically poor schedules, rather than being uniformly distributed across the industry.

- Rows returned: 1
- Columns: distinct_carriers, distinct_airports, total_pairs, carriers, airports

| distinct_carriers | distinct_airports | total_pairs | carriers | airports |
| --- | --- | --- | --- | --- |
| 14 | 22 | 25 | [DH, PI, PI, PI, AS, YV, OO, AA, F9, OH, B6, OO, EV, NW, EV, B6, OH, G4, B6, FL, OH, FL, XE, OO, B6] | [ORD, DFW, LAX, DAY, DUT, ORD, OTH, EYW, PBI, RSW, MIA, ASE, SLC, SJU, DEN, RNO, EWR, SCK, SRQ, SFO, ATL, EWR, EWR, CEC, FLL] |
