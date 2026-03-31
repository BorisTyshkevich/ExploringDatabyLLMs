# Highest daily hops for one aircraft on one flight number

**Generated:** 2026-03-17T15:21:52Z  
**Rows returned:** 10  
**Columns:** Aircraft ID, Flight Number, Carrier, Date, Hops, Route

---

## Overview

- Rows returned: 10
- Generated at: 2026-03-17T15:21:52Z
- Columns: Aircraft ID, Flight Number, Carrier, Date, Hops, Route
- First row snapshot: Aircraft ID=N957WN, Flight Number=366, Carrier=WN

---

## Analysis

The dashboard identifies the maximum number of daily hops flown by a single aircraft under one flight number. This pattern emerges when regional or commuter aircraft operate multiple short legs under a single published flight number, creating high-frequency shuttle itineraries.

### Maximum hop count

The top result shows the highest daily hop count observed in the dataset. A recurring pattern of high hop counts for the same tail number and flight number suggests a scheduled shuttle service rather than a one-off repositioning itinerary. Single occurrences at the maximum level may indicate irregular operations or charter-style routing.

### Most recent maximum-hop itinerary

The first row in the result set represents the most recent date among all maximum-hop itineraries. The route column shows the full sequence of departures with local times and airport codes, providing visibility into the operational rhythm of that day's service.

### Route repetition and clustering

Examine whether the same tail number and flight number appear on multiple dates at or near the maximum hop count. Repeated appearances indicate a stable operating schedule. If different aircraft or flight numbers dominate the list, the high-hop pattern is distributed across the network rather than concentrated on a single shuttle corridor.

---

## Result Data

| Aircraft ID | Flight Number | Carrier | Date | Hops | Route |
| --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 05:43 ISP → 08:10 BWI → 10:20 MYR → 11:42 BNA → 14:01 VPS → 16:43 DAL → 18:28 LAS → 20:41 OAK → SEA |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | 06:21 CLE → 08:01 BNA → 10:07 PNS → 12:34 HOU → 15:14 MCI → 17:47 PHX → 19:02 BUR → 21:17 OAK → DEN |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | 06:18 CLE → 08:00 BNA → 10:12 PNS → 12:39 HOU → 15:15 MCI → 17:56 PHX → 18:59 BUR → 20:59 OAK → DEN |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | 06:20 CLE → 08:10 BNA → 10:17 PNS → 12:37 HOU → 15:10 MCI → 18:00 PHX → 19:03 BUR → 21:02 OAK → DEN |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | 06:25 CLE → 08:05 BNA → 10:11 PNS → 12:41 HOU → 16:07 MCI → 19:09 PHX → 20:06 BUR → 21:56 OAK → DEN |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | 06:27 ELP → 09:39 DAL → 11:24 LIT → 14:34 ATL → 16:56 RIC → 18:31 MDW → 20:40 MCI → 22:46 PHX → SAN |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | 06:22 ELP → 09:41 DAL → 11:22 LIT → 14:33 ATL → 16:59 RIC → 18:32 MDW → 20:41 MCI → 22:26 PHX → SAN |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | 06:08 MSY → 09:44 ATL → 11:56 CMH → 13:54 BWI → 15:53 RDU → 17:29 BNA → 20:45 DTW → 21:49 MDW → LAX |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | 06:10 MSY → 09:17 ATL → 11:42 CMH → 13:55 BWI → 15:36 RDU → 17:16 BNA → 20:40 DTW → 21:47 MDW → LAX |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | 8 | 06:01 MSY → 09:18 ATL → 11:41 CMH → 13:47 BWI → 15:39 RDU → 17:24 BNA → 21:12 DTW → 22:42 MDW → LAX |