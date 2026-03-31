# Highest daily hops for one aircraft on one flight number

> Which itinerary is the highest-hop example, and what does it look like?

The highest-hop example is Southwest Airlines (WN) flight 366 operated by tail **N957WN** on **2024-12-01**, with **8 hops** tracing **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA**. The aircraft started at Long Island MacArthur (ISP) early morning, stopped at Baltimore/Washington (BWI), Myrtle Beach (MYR), Nashville (BNA), Fort Walton Beach/Destin (VPS), Dallas Love Field (DAL), Las Vegas (LAS), and Oakland (OAK), finishing in Seattle/Tacoma (SEA) — a full cross-country marathon from the New York metro area to the Pacific Northwest. All other 8-hop itineraries in the top 10 also belong to Southwest.

- Rows returned: 10
- Columns: Tail_Number, FlightNum, Carrier, FlightDate, hop_count, Route, itinerary_sequence

| Tail_Number | FlightNum | Carrier | FlightDate | hop_count | Route | itinerary_sequence |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | ISP,BWI,MYR,BNA,VPS,DAL,LAS,OAK,SEA |

> Which of the top-ranked itineraries is the most recent?

The most recent top-ranked itinerary is **WN flight 366 / tail N957WN on 2024-12-01** (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA), which is also the lead itinerary. It is the only occurrence of that flight number in the top 10 and carries the latest date in the entire ranked set.

- Rows returned: 1
- Columns: Tail_Number, FlightNum, Carrier, FlightDate, hop_count, Route, itinerary_sequence

| Tail_Number | FlightNum | Carrier | FlightDate | hop_count | Route | itinerary_sequence |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | ISP,BWI,MYR,BNA,VPS,DAL,LAS,OAK,SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

The top itineraries are strongly **recurring scheduled patterns**. Among the top 10 rows there are only 4 distinct routes. WN flight 3149 (CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN) ran on 4 consecutive Sundays in Jan–Feb 2024 with different tail numbers each time. WN flight 2787 (MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX) ran 3 times in Sep–Oct 2022. WN flight 154 (ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN) ran twice in Apr 2023. Only WN 366 appears once. The consistent weekly cadence and published flight numbers confirm these are regular Southwest scheduled turns, not anomalies.

- Rows returned: 4
- Columns: FlightNum, Carrier, Route, occurrences, first_seen, last_seen

| FlightNum | Carrier | Route | occurrences | first_seen | last_seen |
| --- | --- | --- | --- | --- | --- |
| 3149 | WN | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 4 | 2024-01-14T00:00:00Z | 2024-02-18T00:00:00Z |

> What geographic pattern do the top itineraries show?

Every top itinerary is a Southwest Airlines **coast-to-coast diagonal sweep** across the continental US. Routes originate in the East or Southeast (Long Island NY, Cleveland OH, El Paso TX, New Orleans LA) and arc west through the South, Midwest, or Mid-Atlantic before terminating on the West Coast or Mountain West (Seattle, Denver, San Diego, Los Angeles). The lead route ISP→SEA traverses roughly 2,800 miles from the New York metro to the Pacific Northwest, passing through the Mid-Atlantic, Southeast Gulf Coast, South-Central, Southwest, and Pacific Coast in a single day. The pattern reflects Southwest's hub-and-spoke-free point-to-point network, where a single aircraft and flight number chain together many short segments to cover the full width of the country.

- Rows returned: 8
- Columns: OriginCode, DestCode, DepTime, OriginName, OriginLat, OriginLon, DestName, DestLat, DestLon

| OriginCode | DestCode | DepTime | OriginName | OriginLat | OriginLon | DestName | DestLat | DestLon |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| ISP | BWI | 543 | Long Island MacArthur | 40.79611111 | -73.10055556 | Baltimore/Washington International Thurgood Marshall | 39.17583333 | -76.66888889 |
