# Highest daily hops for one aircraft on one flight number

- Rows returned: 10
- Generated at: 2026-03-20T16:51:15Z
- Columns: Tail_Number, Flight_Number_Reporting_Airline, Carrier, FlightDate, Hops, Route, DepTimes
- First row snapshot: Tail_Number=N957WN, Flight_Number_Reporting_Airline=366, Carrier=WN

## Key finding

The maximum number of hops observed for a single aircraft operating one flight number in a single day is **8**. This represents 8 legs flown under a single flight number from early morning to late evening, covering a transcontinental itinerary.

## Most recent maximum-hop itinerary

The most recent 8-hop itinerary is **WN 366** on **2024-12-01**, operated by tail **N957WN**. The full route was **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA**, with actual departure times at each origin: 0543, 0810, 1020, 1142, 1401, 1643, 1828, 2041.

## Route repetition and clustering

The top 10 longest itineraries are dominated by **Southwest Airlines (WN)**, which is the only carrier appearing at this hop count. Several flight numbers show strong repetition: the route pattern is not a one-off but a recurring weekly operating schedule where the same flight number hops through 8 cities. Notable clusters include flight WN 3149 (CLE through DEN, appearing multiple weeks) and WN 154 (ELP through SAN), each repeating with different tail numbers on successive weeks. This confirms a scheduled multi-leg pattern rather than irregular repositioning.

## Top 10 itineraries

| Tail_Number | Flight_Number_Reporting_Airline | Carrier | FlightDate | Hops | Route | DepTimes |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 0543, 0810, 1020, 1142, 1401, 1643, 1828, 2041 |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 0621, 0801, 1007, 1234, 1514, 1747, 1902, 2117 |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 0618, 0800, 1012, 1239, 1515, 1756, 1859, 2059 |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 0620, 0810, 1017, 1237, 1510, 1800, 1903, 2102 |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 0625, 0805, 1011, 1241, 1607, 1909, 2006, 2156 |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 0627, 0939, 1124, 1434, 1656, 1831, 2040, 2246 |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 0622, 0941, 1122, 1433, 1659, 1832, 2041, 2226 |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 0608, 0944, 1156, 1354, 1553, 1729, 2045, 2149 |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 0610, 0917, 1142, 1355, 1536, 1716, 2040, 2147 |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | 8 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 0601, 0918, 1141, 1347, 1539, 1724, 2112, 2242 |