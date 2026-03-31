# Highest daily hops for one aircraft on one flight number

> Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

The maximum distinct-hop count across all history is **8 hops** in a single day under one flight number. All 10 unique routes at that peak are operated by **Southwest Airlines (WN)**, reflecting Southwest's point-to-point milk-run scheduling model. The most recent unique 8-hop route was flown by tail **N957WN** on flight **WN 366** on **2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA, recurrence: 1 day). The highest-recurrence route among the top 10 is **LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN**, which appeared on 47 distinct days.

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
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

The 10 routes span a wide recurrence spectrum — from true one-offs to well-established scheduled patterns.

Recurrence count for each of the 10 routes (distinct calendar days flown, all history):

| Route | Recurrence |
|---|---|
| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |
| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |
| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |
| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |
| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |
| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |
| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |
| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |

**Tiers:** Two routes are effectively one-offs (recurrence ≤ 2 days): ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA (1) and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN (2) — total **2 routes**. Four routes show modest scheduling (5–7 days): CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS, BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX — total **4 routes**. Four routes are strongly recurring scheduled patterns (12–47 days): HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN — total **4 routes**. Overall, 8 of 10 routes appeared on more than 2 days, indicating that 8-hop itineraries are predominantly recurring scheduled rotations rather than one-off operations.

- Rows returned: 10
- Columns: Route, recurrence_count

| Route | recurrence_count |
| --- | --- |
| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

All 10 routes are exclusively **Southwest Airlines (WN)** itineraries and touch **45 distinct airports** across the contiguous United States, with no international, Hawaiian, or Alaskan stops.

The routes form continuous coast-to-coast or diagonal sweeps — classic Southwest milk runs. Each itinerary anchors at one of three regional clusters:

- **East/Gulf Coast origins:** ISP (Long Island), LGA (New York), CLE (Cleveland), BWI (Baltimore), FLL (Fort Lauderdale), HOU (Houston Hobby), MSY (New Orleans)
- **Sun Belt / Midwest bridges:** ATL, BNA, DAL, MEM, MDW, STL, DEN — high-frequency Southwest connecting points mid-country
- **West Coast termini:** OAK, LAX, BUR, SJC, SEA, SMF, LAS — all Bay Area, Southern California, or Pacific Northwest endpoints

A Southwest-specific fingerprint is evident throughout: Dallas Love Field (DAL) appears in place of DFW, Houston Hobby (HOU) in place of IAH, and Chicago Midway (MDW) in place of O'Hare (ORD) — the three airports where Southwest maintains dominance over legacy hubs. The Mountain West (PHX, DEN, LAS, SLC, COS, ABQ) consistently appears as the penultimate cluster bridging Midwest and Pacific stops. No route doubles back geographically; all flow directionally from one region to another, consistent with aircraft repositioning across Southwest's national network.

- Rows returned: 45
- Columns: airport_code, airport_name, lat, lon

| airport_code | airport_name | lat | lon |
| --- | --- | --- | --- |
| ABQ | Albuquerque International Sunport | 35.03805556 | -106.61 |
