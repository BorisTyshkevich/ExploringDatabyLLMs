# Highest Daily Hops for One Aircraft on One Flight Number

- Rows returned: 10
- Generated at: 2026-03-19T09:12:22Z
- Columns: FlightDate, Tail_Number, Carrier, FlightNum, Hops, Route, DepartureTimes
- First row snapshot: FlightDate=2024-12-01T00:00:00Z, Tail_Number=N957WN, Carrier=WN

## Maximum Hop Count

The highest recorded hop count is **8** legs flown by a single aircraft under a single flight number in one day. All top-10 entries belong to Southwest Airlines (WN). This is not a one-off anomaly: the same 8-hop routing patterns recur across multiple weeks on the same flight numbers, indicating deliberate scheduled network itineraries.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Date | 2024-12-01 |
| Carrier | WN |
| Flight Number | 366 |
| Aircraft (Tail) | N957WN |
| Route | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA |

## Top 10 Longest and Most Recent Itineraries

| FlightDate | Tail_Number | Carrier | FlightNum | Hops | Route | DepartureTimes |
| --- | --- | --- | --- | --- | --- | --- |
| 2024-12-01T00:00:00Z | N957WN | WN | 366 | 8 | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA | ISP@543, BWI@810, MYR@1020, BNA@1142, VPS@1401, DAL@1643, LAS@1828, OAK@2041 |
| 2024-02-18T00:00:00Z | N7835A | WN | 3149 | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE@621, BNA@801, PNS@1007, HOU@1234, MCI@1514, PHX@1747, BUR@1902, OAK@2117 |
| 2024-01-28T00:00:00Z | N429WN | WN | 3149 | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE@618, BNA@800, PNS@1012, HOU@1239, MCI@1515, PHX@1756, BUR@1859, OAK@2059 |
| 2024-01-21T00:00:00Z | N228WN | WN | 3149 | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE@620, BNA@810, PNS@1017, HOU@1237, MCI@1510, PHX@1800, BUR@1903, OAK@2102 |
| 2024-01-14T00:00:00Z | N569WN | WN | 3149 | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE@625, BNA@805, PNS@1011, HOU@1241, MCI@1607, PHX@1909, BUR@2006, OAK@2156 |
| 2023-04-30T00:00:00Z | N7742B | WN | 154 | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | ELP@627, DAL@939, LIT@1124, ATL@1434, RIC@1656, MDW@1831, MCI@2040, PHX@2246 |
| 2023-04-16T00:00:00Z | N929WN | WN | 154 | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | ELP@622, DAL@941, LIT@1122, ATL@1433, RIC@1659, MDW@1832, MCI@2041, PHX@2226 |
| 2022-10-23T00:00:00Z | N8631A | WN | 2787 | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY@608, ATL@944, CMH@1156, BWI@1354, RDU@1553, BNA@1729, DTW@2045, MDW@2149 |
| 2022-10-02T00:00:00Z | N8809L | WN | 2787 | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY@610, ATL@917, CMH@1142, BWI@1355, RDU@1536, BNA@1716, DTW@2040, MDW@2147 |
| 2022-09-25T00:00:00Z | N8811L | WN | 2787 | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY@601, ATL@918, CMH@1141, BWI@1347, RDU@1539, BNA@1724, DTW@2112, MDW@2242 |

*DepartureTimes column shows `AIRPORT@HHMM` for each origin leg in chronological order.*

## Route Repetition and Clustering

Three distinct 8-hop route patterns recur across different aircraft on the same flight number, confirming them as recurring scheduled itineraries. WN 3149 (CLE → DEN corridor) is the most frequently repeated, appearing on consecutive Sundays in early 2024 with consistent city pairs from the Midwest through the South and out to the West Coast. WN 2787 (MSY → LAX corridor) and WN 154 (ELP → SAN corridor) each recur across separate weeks as well. This tight clustering by flight number and consistent routing is characteristic of deliberate point-to-point network scheduling.