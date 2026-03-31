# Highest Daily Hops for One Aircraft on One Flight Number

## Overview

This report identifies the longest single-day itineraries flown by one aircraft tail under one flight number. A **hop** is one leg (one Origin→Dest segment); an itinerary of N hops contains N+1 distinct airports.

The maximum observed hop count across all carriers and dates is **8**. All top-10 itineraries reach this ceiling, meaning no aircraft exceeded 8 legs in a single day under a single flight number.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | WN (Southwest Airlines) |
| Flight Number | 366 |
| Tail Number | N957WN |
| Date | 2024-12-01 |
| Route | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA |
| Departure Times | ISP 05:43 | BWI 08:10 | MYR 10:20 | BNA 11:42 | VPS 14:01 | DAL 16:43 | LAS 18:28 | OAK 20:41 |

## Repeated Operating Pattern vs. One-Off

The 8-hop ceiling is **not a one-off event** — it appears across multiple carriers, years, and flight numbers. The single most repeated pattern in the top 10 is **WN 3149**, which flew the identical 8-hop route 4 times:

> CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN

This signals a scheduled, recurring turnaround pattern rather than an irregular re-routing.

## Top 10 Longest Itineraries

| Tail_Number | Flight_Number_Reporting_Airline | IATA_CODE_Reporting_Airline | FlightDate | hops | Route | DepartureTimes |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA | ISP 05:43 \| BWI 08:10 \| MYR 10:20 \| BNA 11:42 \| VPS 14:01 \| DAL 16:43 \| LAS 18:28 \| OAK 20:41 |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE 06:21 \| BNA 08:01 \| PNS 10:07 \| HOU 12:34 \| MCI 15:14 \| PHX 17:47 \| BUR 19:02 \| OAK 21:17 |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE 06:18 \| BNA 08:00 \| PNS 10:12 \| HOU 12:39 \| MCI 15:15 \| PHX 17:56 \| BUR 18:59 \| OAK 20:59 |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE 06:20 \| BNA 08:10 \| PNS 10:17 \| HOU 12:37 \| MCI 15:10 \| PHX 18:00 \| BUR 19:03 \| OAK 21:02 |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | CLE 06:25 \| BNA 08:05 \| PNS 10:11 \| HOU 12:41 \| MCI 16:07 \| PHX 19:09 \| BUR 20:06 \| OAK 21:56 |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | ELP 06:27 \| DAL 09:39 \| LIT 11:24 \| ATL 14:34 \| RIC 16:56 \| MDW 18:31 \| MCI 20:40 \| PHX 22:46 |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | ELP 06:22 \| DAL 09:41 \| LIT 11:22 \| ATL 14:33 \| RIC 16:59 \| MDW 18:32 \| MCI 20:41 \| PHX 22:26 |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY 06:08 \| ATL 09:44 \| CMH 11:56 \| BWI 13:54 \| RDU 15:53 \| BNA 17:29 \| DTW 20:45 \| MDW 21:49 |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY 06:10 \| ATL 09:17 \| CMH 11:42 \| BWI 13:55 \| RDU 15:36 \| BNA 17:16 \| DTW 20:40 \| MDW 21:47 |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | MSY 06:01 \| ATL 09:18 \| CMH 11:41 \| BWI 13:47 \| RDU 15:39 \| BNA 17:24 \| DTW 21:12 \| MDW 22:42 |

## Key Observations

- All top-10 itineraries are operated by **WN (Southwest Airlines)**, consistent with Southwest's point-to-point network model where a single aircraft and flight number chain many short segments in one day.
- Route clustering around certain city pairs (e.g., the recurring CLE–DEN spine) indicates structured schedule blocks, not ad-hoc assignments.
- Departure times span roughly 05:00–23:00 local, confirming full-day aircraft utilization on these itineraries.

---
_Generated: 2026-03-19T07:35:52Z_