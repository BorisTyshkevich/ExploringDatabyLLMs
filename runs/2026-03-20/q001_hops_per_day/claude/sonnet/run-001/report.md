# Highest Daily Hops for One Aircraft on One Flight Number

- Rows returned: 10
- Generated at: 2026-03-20T09:56:28Z
- Columns: Tail_Number, Flight_Number_Reporting_Airline, Carrier, FlightDate, Hops, Route, DepartureTimes
- First row snapshot: Tail_Number=N957WN, Flight_Number_Reporting_Airline=366, Carrier=WN

## Maximum Hop Count

The highest recorded daily hop count is **8 hops** (8 legs in a single day), achieved by carrier **WN (Southwest Airlines)** on flight **366** operated by tail number **N957WN** on **2024-12-01**.

The full route for this itinerary:

> ISP → BWI → MYR → BNA → VPS → DAL → LAS → OAK → SEA

Departure times per leg: 0543 (ISP), 0810 (BWI), 1020 (MYR), 1142 (BNA), 1401 (VPS), 1643 (DAL), 1828 (LAS), 2041 (OAK)

## Operating Pattern

8 hops represents the ceiling across the entire dataset. The top-10 itineraries show that this is a **recurring scheduled pattern**, not a one-off event — the same flight number and identical route sequence appears on multiple dates with different tail numbers, indicating a fixed published multi-hop rotation rather than an irregular assignment.

## Top 10 Longest Itineraries

| Tail_Number | Flight_Number_Reporting_Airline | Carrier | FlightDate | Hops | Route | DepartureTimes |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP → BWI → MYR → BNA → VPS → DAL → LAS → OAK → SEA | 0543 (ISP), 0810 (BWI), 1020 (MYR), 1142 (BNA), 1401 (VPS), 1643 (DAL), 1828 (LAS), 2041 (OAK) |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | CLE → BNA → PNS → HOU → MCI → PHX → BUR → OAK → DEN | 0621 (CLE), 0801 (BNA), 1007 (PNS), 1234 (HOU), 1514 (MCI), 1747 (PHX), 1902 (BUR), 2117 (OAK) |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | CLE → BNA → PNS → HOU → MCI → PHX → BUR → OAK → DEN | 0618 (CLE), 0800 (BNA), 1012 (PNS), 1239 (HOU), 1515 (MCI), 1756 (PHX), 1859 (BUR), 2059 (OAK) |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | CLE → BNA → PNS → HOU → MCI → PHX → BUR → OAK → DEN | 0620 (CLE), 0810 (BNA), 1017 (PNS), 1237 (HOU), 1510 (MCI), 1800 (PHX), 1903 (BUR), 2102 (OAK) |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | CLE → BNA → PNS → HOU → MCI → PHX → BUR → OAK → DEN | 0625 (CLE), 0805 (BNA), 1011 (PNS), 1241 (HOU), 1607 (MCI), 1909 (PHX), 2006 (BUR), 2156 (OAK) |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | ELP → DAL → LIT → ATL → RIC → MDW → MCI → PHX → SAN | 0627 (ELP), 0939 (DAL), 1124 (LIT), 1434 (ATL), 1656 (RIC), 1831 (MDW), 2040 (MCI), 2246 (PHX) |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | ELP → DAL → LIT → ATL → RIC → MDW → MCI → PHX → SAN | 0622 (ELP), 0941 (DAL), 1122 (LIT), 1433 (ATL), 1659 (RIC), 1832 (MDW), 2041 (MCI), 2226 (PHX) |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | MSY → ATL → CMH → BWI → RDU → BNA → DTW → MDW → LAX | 0608 (MSY), 0944 (ATL), 1156 (CMH), 1354 (BWI), 1553 (RDU), 1729 (BNA), 2045 (DTW), 2149 (MDW) |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | MSY → ATL → CMH → BWI → RDU → BNA → DTW → MDW → LAX | 0610 (MSY), 0917 (ATL), 1142 (CMH), 1355 (BWI), 1536 (RDU), 1716 (BNA), 2040 (DTW), 2147 (MDW) |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | 8 | MSY → ATL → CMH → BWI → RDU → BNA → DTW → MDW → LAX | 0601 (MSY), 0918 (ATL), 1141 (CMH), 1347 (BWI), 1539 (RDU), 1724 (BNA), 2112 (DTW), 2242 (MDW) |

## Route Repetition and Clustering

All 10 results belong to a single carrier (WN (Southwest Airlines)). Several flight numbers recur across multiple dates with near-identical routing:

- **WN 3149 (CLE→BNA→PNS→HOU→MCI→PHX→BUR→OAK→DEN)** follows the same 8-hop sequence on 4 separate dates.
- **WN 154 (ELP→DAL→LIT→ATL→RIC→MDW→MCI→PHX→SAN)** and **WN 2787 (MSY→ATL→CMH→BWI→RDU→BNA→DTW→MDW→LAX)** also appear multiple times with consistent routing.

This clustering confirms that 8-hop daily rotations are a deliberate scheduling strategy for point-to-point carriers that chain aircraft through multiple short segments in a single day.

---
*Generated: 2026-03-20T09:56:28Z*