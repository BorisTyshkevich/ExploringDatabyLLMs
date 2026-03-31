# Highest daily hops for one aircraft on one flight number

**Report Generated:** 2026-03-17T16:35:50Z
**Rows Analyzed:** 10

## Data Overview
- Rows returned: 10
- Generated at: 2026-03-17T16:35:50Z
- Columns: Aircraft ID, Flight Number, Carrier, Date, Route, Actual Departure Times
- First row snapshot: Aircraft ID=N957WN, Flight Number=366, Carrier=WN

## Analytical Summary

Based on the top 10 longest single-aircraft, single-flight-number itineraries:

*   **Maximum Hop Count & Pattern:** The result set captures the highest concentration of daily hops operated by a single tail and flight number. By assessing the `Route` sequences in the results, we can discern whether these extreme operational days are highly repetitive "milk-run" routes bouncing back and forth between a small subset of airports, or if they represent long linear point-to-point sequences stretching across the network.
*   **Lead Itinerary:** The foremost row highlights the single most intense, recent itinerary observed. This peak event outlines a specific Carrier and Aircraft operating a singular Flight Number over an extensive multi-stop route on a single `Date`. The sequence explicitly defines the entire operating chain from the initial origin to final destination. 
*   **Route Clustering:** Comparing the top 10 longest itineraries reveals whether extreme multi-hop schedules cluster around specific hub airports or if specific airlines employ this scheduling tactic more frequently than others. Repetition across dates for the same route suggests a structural routing strategy rather than a one-off disruption or aircraft routing anomaly.

## Result Details
| Aircraft ID | Flight Number | Carrier | Date | Route | Actual Departure Times |
| --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | [543, 810, 1020, 1142, 1401, 1643, 1828, 2041] |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | [621, 801, 1007, 1234, 1514, 1747, 1902, 2117] |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | [618, 800, 1012, 1239, 1515, 1756, 1859, 2059] |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | [620, 810, 1017, 1237, 1510, 1800, 1903, 2102] |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | [625, 805, 1011, 1241, 1607, 1909, 2006, 2156] |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | [627, 939, 1124, 1434, 1656, 1831, 2040, 2246] |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | [622, 941, 1122, 1433, 1659, 1832, 2041, 2226] |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | [608, 944, 1156, 1354, 1553, 1729, 2045, 2149] |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | [610, 917, 1142, 1355, 1536, 1716, 2040, 2147] |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | [601, 918, 1141, 1347, 1539, 1724, 2112, 2242] |