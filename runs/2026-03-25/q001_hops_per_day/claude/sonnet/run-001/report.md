# Highest daily hops for one aircraft on one flight number

> Which itinerary is the highest-hop example, and what does it look like?

All top-10 unique routes tie at **8 hops** (9 airports). The leading example by recency is Southwest flight **WN 366**, tail **N957WN**, on **2024-12-01**: `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` — a coast-to-coast turn starting in Islip, NY and ending in Seattle.

- Rows returned: 1
- Columns: Tail_Number, FlightNum, Carrier, FlightDate, hop_count, Route, occurrences

| Tail_Number | FlightNum | Carrier | FlightDate | hop_count | Route | occurrences |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |

> Which of the top-ranked itineraries is the most recent?

The most recent top-ranked itinerary is **WN 366** (tail N957WN) on **2024-12-01**, route `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` with 8 hops. It is the only 8-hop unique route recorded in late 2024.

- Rows returned: 1
- Columns: Tail_Number, FlightNum, Carrier, FlightDate, hop_count, Route, occurrences

| Tail_Number | FlightNum | Carrier | FlightDate | hop_count | Route | occurrences |
| --- | --- | --- | --- | --- | --- | --- |
| N925WN | 2179 | WN | 2025-11-30T00:00:00Z | 6 | SAT-TPA-FLL-RDU-BNA-STL-AUS | 1 |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

The top 10 are a **mix of both**. Two routes are clearly recurring scheduled turn patterns: WN 2884 (`LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN`) appeared **46 times** and WN 1956 (`MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC`) appeared **40 times**, indicating they were regular daily assignments over many years. Others appear only 1–5 times, suggesting irregular or one-off operations. Southwest dominates all 10 entries, consistent with its high-utilization short-haul turn model.

- Rows returned: 10
- Columns: FlightNum, Carrier, FlightDate, hop_count, Route, occurrences

| FlightNum | Carrier | FlightDate | hop_count | Route | occurrences |
| --- | --- | --- | --- | --- | --- |
| 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |
