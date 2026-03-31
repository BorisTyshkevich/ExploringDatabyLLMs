# Highest daily hops for one aircraft on one flight number

> Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

The maximum number of hops flown by a single aircraft on a single flight number in one day is **8 hops** (9-airport itineraries). All 10 unique routes are operated by Southwest Airlines (WN). The most recent top 10 unique routes by departure date are:

| Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |
|---|---|---|---|---|---|---|
| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |
| N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |
| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |
| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |
| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |
| N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |
| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |
| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |
| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |
| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |

- Rows returned: 10
- Columns: Tail_Number, FlightNum, Carrier, FlightDate, hop_count, recurrence_count, Route

| Tail_Number | FlightNum | Carrier | FlightDate | hop_count | recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

The 10 routes show a wide spread of recurrence, from a single occurrence to nearly 50 days:

| Route | Recurrence |
|---|---|
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |
| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |
| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |
| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |
| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |
| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |
| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |
| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |
| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |

**Tier breakdown (totals add to 10):**
- **One-off (1 day):** 1 route — likely an irregular or seasonal assignment.
- **Low recurrence (2–7 days):** 4 routes — occasional scheduled rotations, not daily.
- **Moderate recurrence (12–20 days):** 2 routes — semi-regular scheduled patterns.
- **High recurrence (40–47 days):** 2 routes — clearly stable, recurring scheduled rotations flown across many weeks.

Overall, the majority of 8-hop itineraries are recurring scheduled patterns. Only one is a true one-off; 9 of 10 routes were flown on at least 2 distinct days, and 4 of 10 were flown on 12 or more days, indicating these represent genuine repeating flight schedules rather than one-time anomalies.

- Rows returned: 719
- Columns: Route, recurrence_count

| Route | recurrence_count |
| --- | --- |
| EWR-MDW-ABQ-DEN-SMF-LAS-TUS-OAK-SLC | 1 |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

All airports across the 10 routes are in the **continental United States**, and every route is operated by **Southwest Airlines (WN)**. The itineraries are true transcontinental or cross-country sweeps, connecting three geographic bands:

- **West (17 airports):** SEA, OAK, SJC, SMF, RNO, LAX, BUR, SAN, PSP, LAS, PHX, SLC, ABQ, ELP, COS, DEN, MAF — dominated by Pacific Coast, Mountain West, and Desert Southwest airports.
- **Central (14 airports):** STL, MEM, ORD, MDW, MCI, ICT, LIT, BNA, DAL, HOU, MSY, JAN, PNS, VPS — the Mississippi/Gulf corridor and Midwest hubs.
- **East (14 airports):** LGA, ISP, BWI, IAD, RIC, RDU, MYR, CMH, CLE, DTW, ATL, TPA, MCO, FLL — stretching from the Northeast down the Atlantic Seaboard to Florida.

The dominant pattern is **coast-to-coast routing**: each itinerary typically originates on one coast or the Gulf South, threads through 2–3 inland stops, and terminates on the opposite coast. Several routes end or begin at West Coast leisure-travel airports (OAK, LAX, SJC, LAS). Recurring stops like BWI, BNA, DAL, LAS, MDW, and OAK appear across multiple routes, reflecting Southwest's key connecting cities. No route is confined to a single region; all 10 span at least two of the three geographic bands.

- Rows returned: 3
- Columns: region, airport_count, airports

| region | airport_count | airports |
| --- | --- | --- |
| West | 17 | [LAS, OAK, SEA, PHX, BUR, DEN, ELP, SAN, LAX, SLC, SJC, COS, SMF, PSP, RNO, ABQ, MAF] |
