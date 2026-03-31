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

The maximum daily hops for a single aircraft on one flight number is **8 hops** (9 airports). All top results belong to **Southwest Airlines (WN)**, consistent with their point-to-point network model where a single aircraft cycles through many cities under one flight number.

The 10 most recent unique 8-hop routes are:

| # | Aircraft | Flight | Carrier | Date | Hops | Recurrence | Route |
|---|----------|--------|---------|------|------|------------|-------|
| 1 | N957WN | 366 | WN | 2024-12-01 | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |
| 2 | N7835A | 3149 | WN | 2024-02-18 | 8 | 5 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN |
| 3 | N7742B | 154 | WN | 2023-04-30 | 8 | 2 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN |
| 4 | N8631A | 2787 | WN | 2022-10-23 | 8 | 5 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX |
| 5 | N416WN | 1956 | WN | 2022-09-01 | 8 | 40 | MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC |
| 6 | N7713A | 2884 | WN | 2022-08-31 | 8 | 47 | LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN |
| 7 | N219WN | 3378 | WN | 2021-10-31 | 8 | 7 | SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS |
| 8 | N262WN | 904 | WN | 2021-08-27 | 8 | 20 | HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK |
| 9 | N484WN | 2294 | WN | 2021-08-25 | 8 | 12 | BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK |
| 10 | N225WN | 3530 | WN | 2021-08-08 | 8 | 5 | BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX |

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

The top 10 itineraries show a mix of recurring scheduled patterns and occasional/one-off routes.

**Recurrence counts for each route:**

1. ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA — **1 day** (one-off)
2. CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN — **5 days**
3. ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN — **2 days**
4. MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX — **5 days**
5. MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC — **40 days**
6. LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — **47 days**
7. SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS — **7 days**
8. HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK — **20 days**
9. BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK — **12 days**
10. BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — **5 days**

**Categories:**
- **High-recurrence scheduled patterns** (20+ days): 3 routes (#5, #6, #8) — these ran for weeks or months as stable daily scheduled rotations.
- **Moderate-recurrence patterns** (5–12 days): 5 routes (#2, #4, #7, #9, #10) — likely short-lived seasonal or weekly schedule patterns.
- **Low-recurrence / one-offs** (1–2 days): 2 routes (#1, #3) — appear to be irregular or ad-hoc itineraries, possibly driven by irregular operations or one-time schedule experiments.

Overall, **8 of the 10 routes** recurred on multiple days, indicating they represent recurring scheduled patterns rather than one-offs. The two highest-recurrence routes (WN 2884 at 47 days and WN 1956 at 40 days) were clearly long-running scheduled rotations.

- Rows returned: 10
- Columns: Route, recurrence_count, most_recent_date

| Route | recurrence_count | most_recent_date |
| --- | --- | --- |
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 | 2024-12-01T00:00:00Z |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

The top 10 eight-hop routes span **45 unique airports** across the entire continental United States, forming **coast-to-coast transcontinental zigzag patterns**.

**Key geographic observations:**

- **East-to-West sweep dominates:** Most routes originate on the East Coast or Gulf Coast (ISP, BWI, LGA, CLE, MSY, HOU) and terminate on the West Coast (SEA, OAK, SJC, DEN, SAN, LAX, LAS), reflecting Southwest's strategy of repositioning aircraft westward across the day.
- **Major Southwest hubs recur frequently:** DAL (Dallas Love), MDW (Chicago Midway), PHX (Phoenix), LAS (Las Vegas), OAK (Oakland), BWI (Baltimore), and BNA (Nashville) appear in multiple routes, anchoring the network's mid-continent and endpoint connections.
- **Mid-continent corridor:** Several routes thread through the Great Plains and Mountain West via ICT (Wichita), MCI (Kansas City), DEN (Denver), COS (Colorado Springs), SLC (Salt Lake City), and ABQ (Albuquerque) — reflecting Southwest's strong domestic interior presence.
- **Gulf Coast and Southeast representation:** MSY (New Orleans), HOU (Houston), PNS (Pensacola), VPS (Destin-Fort Walton Beach), MYR (Myrtle Beach), ATL (Atlanta), and FLL (Fort Lauderdale) appear, showing Sun Belt leisure destinations woven into operational routings.
- **No international or Hawaii airports:** All 45 airports are within the contiguous 48 states, consistent with Southwest's domestic point-to-point model.

- Rows returned: 45
- Columns: airport_code, DisplayAirportName, Latitude, Longitude

| airport_code | DisplayAirportName | Latitude | Longitude |
| --- | --- | --- | --- |
| SEA | Seattle/Tacoma International | 47.45 | -122.31166667 |
