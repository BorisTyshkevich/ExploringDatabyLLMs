# Highest Daily Hops for One Aircraft on One Flight Number

## Overview

- Rows returned: 10
- Generated at: 2026-03-19T11:04:55Z
- Columns: Tail_Number, Flight_Number_Reporting_Airline, Carrier, FlightDate, Hops, Route, DepartureSchedule
- First row snapshot: Tail_Number=N957WN, Flight_Number_Reporting_Airline=366, Carrier=WN

The maximum number of hops flown by a single aircraft under the same flight number in a single day is **8**, achieved exclusively by carrier **WN (Southwest Airlines)**. All 10 entries in the top-10 list share this hop count, indicating that 8-hop itineraries are a structured, repeating feature of Southwest Airlines scheduling rather than a one-off occurrence.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | WN |
| Flight Number | 366 |
| Aircraft (Tail) | N957WN |
| Date | 2024-12-01 |
| Hops | 8 |
| Route | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA |
| Departure Schedule | ISP 05:43, BWI 08:10, MYR 10:20, BNA 11:42, VPS 14:01, DAL 16:43, LAS 18:28, OAK 20:41 |

The route spans 9 airports across the country, with departures beginning before 06:00 local time and the last segment departing after 20:00, covering roughly 16 hours of continuous flying operations.

## Top 10 Longest and Most Recent Itineraries

| Tail_Number | Flight_Number_Reporting_Airline | Carrier | FlightDate | Hops | Route | DepartureSchedule |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA | ISP 05:43, BWI 08:10, MYR 10:20, BNA 11:42, VPS 14:01, DAL 16:43, LAS 18:28, OAK 20:41 |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE 06:21, BNA 08:01, PNS 10:07, HOU 12:34, MCI 15:14, PHX 17:47, BUR 19:02, OAK 21:17 |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE 06:18, BNA 08:00, PNS 10:12, HOU 12:39, MCI 15:15, PHX 17:56, BUR 18:59, OAK 20:59 |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE 06:20, BNA 08:10, PNS 10:17, HOU 12:37, MCI 15:10, PHX 18:00, BUR 19:03, OAK 21:02 |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE 06:25, BNA 08:05, PNS 10:11, HOU 12:41, MCI 16:07, PHX 19:09, BUR 20:06, OAK 21:56 |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | ELP 06:27, DAL 09:39, LIT 11:24, ATL 14:34, RIC 16:56, MDW 18:31, MCI 20:40, PHX 22:46 |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | ELP 06:22, DAL 09:41, LIT 11:22, ATL 14:33, RIC 16:59, MDW 18:32, MCI 20:41, PHX 22:26 |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY 06:08, ATL 09:44, CMH 11:56, BWI 13:54, RDU 15:53, BNA 17:29, DTW 20:45, MDW 21:49 |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY 06:10, ATL 09:17, CMH 11:42, BWI 13:55, RDU 15:36, BNA 17:16, DTW 20:40, MDW 21:47 |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY 06:01, ATL 09:18, CMH 11:41, BWI 13:47, RDU 15:39, BNA 17:24, DTW 21:12, MDW 22:42 |

## Route Patterns and Clustering

The top-10 itineraries reveal strong route clustering by flight number:

- **Flight 3149 (CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN)** operates the identical 8-stop sequence 3 times across different dates and tail numbers, confirming it as a regularly scheduled transcontinental multi-hop service.
- **Flight 154 (ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN)** and **Flight 2787 (MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX)** each repeat their respective fixed routes on multiple dates, further demonstrating that 8-hop itineraries are codified schedule patterns.
- All top-10 entries are operated by WN (Southwest Airlines), which historically employs a point-to-point multi-stop model suited to high daily utilization of individual aircraft.