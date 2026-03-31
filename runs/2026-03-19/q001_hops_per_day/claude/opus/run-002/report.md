# Highest daily hops for one aircraft on one flight number

- Rows returned: 10
- Generated at: 2026-03-19T08:34:12Z
- Columns: Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate, Hops, Route, DepSchedule
- First row snapshot: Tail_Number=N957WN, Flight_Number_Reporting_Airline=366, IATA_CODE_Reporting_Airline=WN

## Key finding

The maximum observed hop count is **8** legs in a single day for one aircraft operating under one flight number. All 8-hop itineraries belong to **WN (Southwest Airlines)** (Southwest Airlines), reflecting its characteristic point-to-point, multi-stop scheduling model.

## Most recent maximum-hop itinerary

The most recent 8-hop itinerary was **WN 366** on **2024-12-01**, aircraft **N957WN**, flying the route:

> ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA

Departure times from each origin: ISP 05:43, BWI 08:10, MYR 10:20, BNA 11:42, VPS 14:01, DAL 16:43, LAS 18:28, OAK 20:41

The aircraft began its day at ISP (Long Island MacArthur) at 05:43 and reached SEA (Seattle-Tacoma) after 20:41, covering eight legs coast-to-coast over roughly 15 hours.

## Route repetition and clustering

Among the top 10 longest itineraries, strong route repetition is evident:

- **WN 3149** (CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN) appears **4 times** across January-February 2024, indicating a recurring weekly schedule pattern.
- **WN 2787** (MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX) appears **3 times** across September-October 2022.
- **WN 154** (ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN) appears **twice** in April 2023.
- **WN 366** appears once with a distinct route (ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA).

These are clearly **repeating weekly operating patterns**, not one-off itineraries. Southwest regularly assigns a single flight number to an aircraft that hops across the country through 8 or more cities in a single day.

## Top 10 results

| Tail_Number | Flight_Number_Reporting_Airline | IATA_CODE_Reporting_Airline | FlightDate | Hops | Route | DepSchedule |
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