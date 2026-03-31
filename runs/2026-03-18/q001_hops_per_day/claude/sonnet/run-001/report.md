# Highest daily hops for one aircraft on one flight number

_Generated: 2026-03-18T21:33:25Z_

## Overview

- Rows returned: 10
- Generated at: 2026-03-18T21:33:25Z
- Columns: Tail_Number, Carrier, FlightNum, FlightDate, Hops, Route, DepartureTimes
- First row snapshot: Tail_Number=N957WN, Carrier=WN, FlightNum=366

## Top 10 Longest Daily Itineraries

| Tail_Number | Carrier | FlightNum | FlightDate | Hops | Route | DepartureTimes |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | WN | 366 | 2024-12-01T00:00:00Z | 8 | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA | 05:43, 08:10, 10:20, 11:42, 14:01, 16:43, 18:28, 20:41 |
| N7835A | WN | 3149 | 2024-02-18T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:21, 08:01, 10:07, 12:34, 15:14, 17:47, 19:02, 21:17 |
| N429WN | WN | 3149 | 2024-01-28T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:18, 08:00, 10:12, 12:39, 15:15, 17:56, 18:59, 20:59 |
| N228WN | WN | 3149 | 2024-01-21T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:20, 08:10, 10:17, 12:37, 15:10, 18:00, 19:03, 21:02 |
| N569WN | WN | 3149 | 2024-01-14T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:25, 08:05, 10:11, 12:41, 16:07, 19:09, 20:06, 21:56 |
| N7742B | WN | 154 | 2023-04-30T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | 06:27, 09:39, 11:24, 14:34, 16:56, 18:31, 20:40, 22:46 |
| N929WN | WN | 154 | 2023-04-16T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | 06:22, 09:41, 11:22, 14:33, 16:59, 18:32, 20:41, 22:26 |
| N8631A | WN | 2787 | 2022-10-23T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:08, 09:44, 11:56, 13:54, 15:53, 17:29, 20:45, 21:49 |
| N8809L | WN | 2787 | 2022-10-02T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:10, 09:17, 11:42, 13:55, 15:36, 17:16, 20:40, 21:47 |
| N8811L | WN | 2787 | 2022-09-25T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:01, 09:18, 11:41, 13:47, 15:39, 17:24, 21:12, 22:42 |

## Analysis

### Maximum Hop Count

The maximum observed is **{{max_hops}} hops** in a single calendar day under a single flight number, producing an itinerary of {{max_hops_plus_one}} airports. This is not a one-off anomaly: the top 10 table contains multiple distinct dates with the same hop count, confirming it reflects a repeating scheduled operation rather than an exceptional day.

### Most Recent Maximum-Hop Itinerary

The single most recent entry in the top 10 is:

- **Carrier / Flight:** {{most_recent_carrier}} {{most_recent_flight_num}}
- **Date:** {{most_recent_date}}
- **Tail Number:** {{most_recent_tail}}
- **Route:** {{most_recent_route}}
- **Leg departure times:** {{most_recent_dep_times}}

### Route Repetition and Clustering

Across the top 10 itineraries, distinct named flight numbers recur on multiple dates with identical or nearly identical routing:

| Tail_Number | Carrier | FlightNum | FlightDate | Hops | Route | DepartureTimes |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | WN | 366 | 2024-12-01T00:00:00Z | 8 | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA | 05:43, 08:10, 10:20, 11:42, 14:01, 16:43, 18:28, 20:41 |
| N7835A | WN | 3149 | 2024-02-18T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:21, 08:01, 10:07, 12:34, 15:14, 17:47, 19:02, 21:17 |
| N429WN | WN | 3149 | 2024-01-28T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:18, 08:00, 10:12, 12:39, 15:15, 17:56, 18:59, 20:59 |
| N228WN | WN | 3149 | 2024-01-21T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:20, 08:10, 10:17, 12:37, 15:10, 18:00, 19:03, 21:02 |
| N569WN | WN | 3149 | 2024-01-14T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:25, 08:05, 10:11, 12:41, 16:07, 19:09, 20:06, 21:56 |
| N7742B | WN | 154 | 2023-04-30T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | 06:27, 09:39, 11:24, 14:34, 16:56, 18:31, 20:40, 22:46 |
| N929WN | WN | 154 | 2023-04-16T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | 06:22, 09:41, 11:22, 14:33, 16:59, 18:32, 20:41, 22:26 |
| N8631A | WN | 2787 | 2022-10-23T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:08, 09:44, 11:56, 13:54, 15:53, 17:29, 20:45, 21:49 |
| N8809L | WN | 2787 | 2022-10-02T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:10, 09:17, 11:42, 13:55, 15:36, 17:16, 20:40, 21:47 |
| N8811L | WN | 2787 | 2022-09-25T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:01, 09:18, 11:41, 13:47, 15:39, 17:24, 21:12, 22:42 |

All carriers represented are Southwest Airlines (WN), consistent with WN's point-to-point network model in which a single flight number can chain many short segments across one day. Route patterns appear to repeat weekly or on a regular schedule, suggesting these are published schedule rotations rather than ad-hoc operations.