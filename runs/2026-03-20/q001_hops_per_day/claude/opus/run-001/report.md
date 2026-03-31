# Highest daily hops for one aircraft on one flight number

## Overview

- Rows returned: 10
- Generated at: 2026-03-20T10:09:50Z
- Columns: Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate, Hops, Route, DepTimes
- First row snapshot: Tail_Number=N957WN, Flight_Number_Reporting_Airline=366, IATA_CODE_Reporting_Airline=WN

## Key Finding

The maximum number of hops observed for a single aircraft operating one flight number in a single day is **8**. This pattern is not a one-off anomaly — it appears repeatedly across multiple dates and flight numbers, all operated by **Southwest Airlines (WN)**, reflecting that carrier's point-to-point operating model where a single flight number can traverse the entire country in a chain of short-haul segments.

## Most Recent Maximum-Hop Itinerary

The most recent 8-hop itinerary is **WN 366** on **2024-12-01**, aircraft **N957WN**, flying the route **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA** with departures at **5:43, 8:10, 10:20, 11:42, 14:01, 16:43, 18:28, 20:41**.

This itinerary spans coast-to-coast from the New York area (ISP) to the Pacific Northwest (SEA) in eight legs over roughly 15 hours.

## Route Repetition and Clustering

Among the top 10 longest itineraries, strong route repetition is evident:

- **WN 3149** appears four times (Jan–Feb 2024) on the identical route CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, indicating a fixed weekly schedule.
- **WN 2787** appears three times (Sep–Oct 2022) repeating the route MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX.
- **WN 154** appears twice (Apr 2023) on the route ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN.

All 10 entries are Southwest Airlines (WN), consistent with its unique operational model of multi-stop through-flights under a single flight number.

## Result Detail

| Tail_Number | Flight_Number_Reporting_Airline | IATA_CODE_Reporting_Airline | FlightDate | Hops | Route | DepTimes |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 5:43, 8:10, 10:20, 11:42, 14:01, 16:43, 18:28, 20:41 |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 6:21, 8:01, 10:07, 12:34, 15:14, 17:47, 19:02, 21:17 |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 6:18, 8:00, 10:12, 12:39, 15:15, 17:56, 18:59, 20:59 |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 6:20, 8:10, 10:17, 12:37, 15:10, 18:00, 19:03, 21:02 |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 6:25, 8:05, 10:11, 12:41, 16:07, 19:09, 20:06, 21:56 |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 6:27, 9:39, 11:24, 14:34, 16:56, 18:31, 20:40, 22:46 |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 6:22, 9:41, 11:22, 14:33, 16:59, 18:32, 20:41, 22:26 |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 6:08, 9:44, 11:56, 13:54, 15:53, 17:29, 20:45, 21:49 |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 6:10, 9:17, 11:42, 13:55, 15:36, 17:16, 20:40, 21:47 |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | 8 | MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 6:01, 9:18, 11:41, 13:47, 15:39, 17:24, 21:12, 22:42 |