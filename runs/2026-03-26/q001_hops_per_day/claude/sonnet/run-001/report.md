# Highest daily hops for one aircraft on one flight number

> Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.
Count hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

The maximum single-day hop count is **8 legs** for one aircraft under the same flight number, achieved exclusively by **Southwest Airlines (WN)**. Across all history, there are multiple distinct 8-leg route strings. The 10 most recently flown unique routes (one row per distinct route, showing the most recent occurrence) are listed below with their aircraft id, flight number, flight date, hop count, recurrence count (distinct days that exact route was ever flown), and the full textual route:

| Aircraft | Flight# | Carrier | Date | Hops | Recurrence | Route |
|---|---|---|---|---|---|---|
| N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |
| N7835A | 3149 | WN | 2024-02-18 | 8 | 4 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |
| N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |
| N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |
| N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |
| N7713A | 2884 | WN | 2022-08-31 | 8 | 46 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |
| N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |
| N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |
| N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |
| N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, most_recent_date, hop_count, recurrence_count, Route

| aircraft_id | flight_number | carrier | most_recent_date | hop_count | recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

The 10 routes show a wide spectrum from one-offs to established scheduled patterns. Recurrence counts for each route (exact days flown in full history):

1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA: **1**
2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN: **4**
3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN: **2**
4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX: **5**
5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC: **40**
6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN: **46**
7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS: **7**
8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK: **20**
9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK: **12**
10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX: **5**

Tiered summary (totals add to 10):
- **One-off (1 day):** 1 route — route #1 was flown exactly once.
- **Rare (2–7 days):** 5 routes — routes #2, #3, #4, #7, #10 were flown on only a handful of occasions.
- **Recurring (12–20 days):** 2 routes — routes #8 and #9 appear periodically, suggesting semi-regular scheduling.
- **Highly recurring (40–46 days):** 2 routes — routes #5 and #6 were each flown on more than 40 distinct days, indicating firmly established scheduled patterns.

Overall, 6 of 10 routes are rare-to-one-off while 4 show meaningful recurrence, meaning most 8-hop itineraries are ad-hoc operational assignments, but a small subset are genuine recurring scheduled turns.

- Rows returned: 10
- Columns: Route, recurrence_count, tier

| Route | recurrence_count | tier |
| --- | --- | --- |
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 | one-off |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

The 10 routes collectively touch **45 unique airports across 27 US states** (plus DC). All airports are in the contiguous United States — no Hawaii, Alaska, or international airports appear, consistent with Southwest Airlines' domestic network during the covered periods.

Key geographic patterns:

1. **Coast-to-coast sweep**: Every route crosses multiple time zones in a single day. Routes span from the Northeast (LGA, ISP, BWI) or Southeast (MSY, TPA, ATL) all the way to the West Coast (OAK, BUR, LAX, SJC, SEA, SMF) or terminate in the Mountain West (DEN, LAS, PHX). This is characteristic of maximum aircraft utilization.

2. **California dominates as a terminus**: 7 California airports appear (BUR, LAX, OAK, PSP, SAN, SJC, SMF) — more than any other state — and California is the endpoint of 7 of the 10 routes, confirming that the West Coast is the preferred turn-around anchor for these ultra-long turns.

3. **Texas and the South serve as through-hubs**: Dallas Love Field (DAL), Houston (HOU), and other Texas airports appear as interior waypoints in multiple routes, reflecting Southwest's historical strength in the South-Central corridor.

4. **Florida and the Southeast as origin clusters**: 5 Florida airports (FLL, MCO, PNS, TPA, VPS) and several Southeast cities (ATL, MYR, RDU, RIC, BNA, MEM) frequently appear early in routes, suggesting morning departures from the East before heading westward.

5. **No large hub concentration**: The routes do not funnel through any single mega-hub. The 45 airports are spread across the network, consistent with Southwest's point-to-point model rather than hub-and-spoke.

- Rows returned: 45
- Columns: airport_code, city, state, state_code

| airport_code | city | state | state_code |
| --- | --- | --- | --- |
| LIT | Little Rock, AR | Arkansas | AR |
