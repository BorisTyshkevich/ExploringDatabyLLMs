# Highest daily hops for one aircraft on one flight number

- Rows returned: 10
- Generated at: 2026-03-19T22:38:47Z
- Columns: Date, Carrier, AircraftID, FlightNumber, Hops, Route, DepartureTimes
- First row snapshot: Date=2024-12-01T00:00:00Z, Carrier=WN, AircraftID=N957WN

## Summary

The maximum observed daily hop count for a single aircraft on the same flight number is **8**. This does not appear to be a one-off anomaly but rather a deliberate and repeated operating pattern, exclusively operated by carrier **WN** (Southwest Airlines) among the top 10 longest itineraries.

## Most Recent Longest Itinerary

The most recent 8-hop itinerary was flown by **WN** on flight **366** using aircraft **N957WN** on **2024-12-01**. 

The full 8-leg route traversed: **ISP - BWI - MYR - BNA - VPS - DAL - LAS - OAK - SEA**.

## Route Repetition and Clustering

Looking at the top 10 longest itineraries, there is a clear pattern of route repetition and clustering:
- **Consistent Carrier**: All top 10 records belong to Southwest Airlines (WN).
- **Identical Repeated Routes**: Southwest repeatedly flies the exact same 8-hop route using the same flight number on different dates. For example, flight 3149 flew the exact same 8-leg cross-country route (CLE - BNA - PNS - HOU - MCI - PHX - BUR - OAK - DEN) on at least four separate dates in early 2024.
- **Transcontinental Sweeps**: These high-hop flights generally sweep across the country, serving as a "milk run" connecting multiple regional airports and major hubs from the East Coast/Midwest to the West Coast.

| Date | Carrier | AircraftID | FlightNumber | Hops | Route | DepartureTimes |
| --- | --- | --- | --- | --- | --- | --- |
| 2024-12-01T00:00:00Z | WN | N957WN | 366 | 8 | ISP - BWI - MYR - BNA - VPS - DAL - LAS - OAK - SEA | 543, 810, 1020, 1142, 1401, 1643, 1828, 2041 |
| 2024-02-18T00:00:00Z | WN | N7835A | 3149 | 8 | CLE - BNA - PNS - HOU - MCI - PHX - BUR - OAK - DEN | 621, 801, 1007, 1234, 1514, 1747, 1902, 2117 |
| 2024-01-28T00:00:00Z | WN | N429WN | 3149 | 8 | CLE - BNA - PNS - HOU - MCI - PHX - BUR - OAK - DEN | 618, 800, 1012, 1239, 1515, 1756, 1859, 2059 |
| 2024-01-21T00:00:00Z | WN | N228WN | 3149 | 8 | CLE - BNA - PNS - HOU - MCI - PHX - BUR - OAK - DEN | 620, 810, 1017, 1237, 1510, 1800, 1903, 2102 |
| 2024-01-14T00:00:00Z | WN | N569WN | 3149 | 8 | CLE - BNA - PNS - HOU - MCI - PHX - BUR - OAK - DEN | 625, 805, 1011, 1241, 1607, 1909, 2006, 2156 |
| 2023-04-30T00:00:00Z | WN | N7742B | 154 | 8 | ELP - DAL - LIT - ATL - RIC - MDW - MCI - PHX - SAN | 627, 939, 1124, 1434, 1656, 1831, 2040, 2246 |
| 2023-04-16T00:00:00Z | WN | N929WN | 154 | 8 | ELP - DAL - LIT - ATL - RIC - MDW - MCI - PHX - SAN | 622, 941, 1122, 1433, 1659, 1832, 2041, 2226 |
| 2022-10-23T00:00:00Z | WN | N8631A | 2787 | 8 | MSY - ATL - CMH - BWI - RDU - BNA - DTW - MDW - LAX | 608, 944, 1156, 1354, 1553, 1729, 2045, 2149 |
| 2022-10-02T00:00:00Z | WN | N8809L | 2787 | 8 | MSY - ATL - CMH - BWI - RDU - BNA - DTW - MDW - LAX | 610, 917, 1142, 1355, 1536, 1716, 2040, 2147 |
| 2022-09-25T00:00:00Z | WN | N8811L | 2787 | 8 | MSY - ATL - CMH - BWI - RDU - BNA - DTW - MDW - LAX | 601, 918, 1141, 1347, 1539, 1724, 2112, 2242 |