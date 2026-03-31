# Highest daily hops for one aircraft on one flight number

> Which itinerary is the highest-hop example, and what does it look like?

The highest-hop example is **WN flight 366 on 2024-12-01**, operated by tail **N957WN**, covering **8 hops** (9 airports) along the route **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA**. The aircraft started at Long Island MacArthur (ISP), made stops in Baltimore (BWI), Myrtle Beach (MYR), Nashville (BNA), Fort Walton Beach (VPS), Dallas Love Field (DAL), Las Vegas (LAS), and Oakland (OAK), before finishing in Seattle (SEA). All top-10 itineraries share the same 8-hop maximum and are Southwest Airlines operations.

- Rows returned: 10
- Columns: Rank, FlightDate, Tail_Number, FlightNum, Carrier, HopCount, Route, LegSeq

| Rank | FlightDate | Tail_Number | FlightNum | Carrier | HopCount | Route | LegSeq |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 2024-12-01T00:00:00Z | N957WN | 366 | WN | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | [{=BWI}, {=MYR}, {=BNA}, {=VPS}, {=DAL}, {=LAS}, {=OAK}, {=SEA}] |

> Which of the top-ranked itineraries is the most recent?

The most recent top-ranked itinerary is **WN flight 366 on 2024-12-01** (tail N957WN), route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA. It is both Rank 1 (the lead itinerary) and the most recent date in the top-10, ahead of the next most recent (WN 3149 on 2024-02-18).

- Rows returned: 1
- Columns: FlightDate, Tail_Number, FlightNum, Carrier, HopCount, Route

| FlightDate | Tail_Number | FlightNum | Carrier | HopCount | Route |
| --- | --- | --- | --- | --- | --- |
| 2024-12-01T00:00:00Z | N957WN | 366 | WN | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

The top itineraries are overwhelmingly **recurring scheduled patterns**, not one-offs. The most frequent 8-hop route (WN 2215: CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO) ran **61 times** between August and October 2009. WN 1923 ran 60 times, WN 2558 ran 55 times. These routes operated nearly every day of a seasonal schedule window, indicating deliberate long-day rotations assigned to a flight number — not anomalies.

- Rows returned: 10
- Columns: FlightNum, Carrier, Route, DateCount, FirstSeen, LastSeen

| FlightNum | Carrier | Route | DateCount | FirstSeen | LastSeen |
| --- | --- | --- | --- | --- | --- |
| 2215 | WN | CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO | 61 | 2009-08-17T00:00:00Z | 2009-10-30T00:00:00Z |

> What geographic pattern do the top itineraries show?

All 8-hop itineraries trace **long east-to-west (or southeast-to-northwest) transcontinental sweeps** across the continental United States. Routes consistently originate in the Eastern US (e.g., Columbus OH, Pittsburgh PA, Hartford CT, Fort Lauderdale FL, Richmond VA) and terminate at West Coast or Pacific Northwest airports (e.g., Reno, San Jose, Oakland, Portland, Ontario CA, Seattle). Along the way they pass through Midwest hubs (MDW, MCI, STL), Southern cities (DAL, HOU, BNA), and Desert Southwest airports (PHX, ELP, ABQ, LAS). The pattern reflects Southwest Airlines routing an aircraft through its network in a single long operational day, effectively a transcontinental relay covering 2,000+ miles west to east.

- Rows returned: 10
- Columns: FlightNum, Carrier, Route, DateCount, StartName, StartLat, StartLon, EndName, EndLat, EndLon

| FlightNum | Carrier | Route | DateCount | StartName | StartLat | StartLon | EndName | EndLat | EndLon |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 2215 | WN | CMH-MDW-MCI-DAL-ELP-PHX-SNA-SJC-RNO | 61 | John Glenn Columbus International | 39.99694444 | -82.89222222 | Reno/Tahoe International | 39.49916667 | -119.76805556 |
