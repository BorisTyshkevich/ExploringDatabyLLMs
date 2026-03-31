# Highest Daily Hops for One Aircraft on One Flight Number

- Rows returned: 10
- Generated at: 2026-03-19T12:23:06Z
- Columns: Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate, Hops, Route
- First row snapshot: Tail_Number=N957WN, Flight_Number_Reporting_Airline=366, IATA_CODE_Reporting_Airline=WN

## Maximum Hop Count

The highest number of hops recorded for a single aircraft operating under a single flight number on one calendar day is **8**. This means the aircraft completed 8 individual flight legs under the same flight number within a single day.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | WN (Southwest Airlines) |
| Flight Number | 366 |
| Aircraft (Tail) | N957WN |
| Date | 2024-12-01 |
| Route | ISP(543) -> BWI(810) -> MYR(1020) -> BNA(1142) -> VPS(1401) -> DAL(1643) -> LAS(1828) -> OAK(2041) -> SEA |

The route above shows each origin airport followed by its actual departure time (HHMM) in parentheses, ending at the final destination.

## Operating Pattern Analysis

The 8-hop maximum appears to be a repeated scheduled pattern, not a one-off. Flight WN3149 ran the same 8-leg CLE→BNA→PNS→HOU→MCI→PHX→BUR→OAK→DEN itinerary on four consecutive Sundays in January–February 2024, operated by different aircraft each week. Flight WN154 similarly repeated an 8-hop ELP-based itinerary across two Sundays in April 2023. This strongly indicates these are scheduled multi-stop routes, not exceptional occurrences.

## Top 10 Longest Itineraries

| Tail_Number | Flight_Number_Reporting_Airline | IATA_CODE_Reporting_Airline | FlightDate | Hops | Route |
| --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP(543) -> BWI(810) -> MYR(1020) -> BNA(1142) -> VPS(1401) -> DAL(1643) -> LAS(1828) -> OAK(2041) -> SEA |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | CLE(621) -> BNA(801) -> PNS(1007) -> HOU(1234) -> MCI(1514) -> PHX(1747) -> BUR(1902) -> OAK(2117) -> DEN |
| N263WN | 3149 | WN | 2024-02-04T00:00:00Z | 8 | BUR(0) -> OAK(0) -> PHX(0) -> CLE(614) -> BNA(810) -> PNS(1010) -> HOU(1239) -> MCI(1514) -> PHX |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | CLE(618) -> BNA(800) -> PNS(1012) -> HOU(1239) -> MCI(1515) -> PHX(1756) -> BUR(1859) -> OAK(2059) -> DEN |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | CLE(620) -> BNA(810) -> PNS(1017) -> HOU(1237) -> MCI(1510) -> PHX(1800) -> BUR(1903) -> OAK(2102) -> DEN |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | CLE(625) -> BNA(805) -> PNS(1011) -> HOU(1241) -> MCI(1607) -> PHX(1909) -> BUR(2006) -> OAK(2156) -> DEN |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | ELP(627) -> DAL(939) -> LIT(1124) -> ATL(1434) -> RIC(1656) -> MDW(1831) -> MCI(2040) -> PHX(2246) -> SAN |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | ELP(622) -> DAL(941) -> LIT(1122) -> ATL(1433) -> RIC(1659) -> MDW(1832) -> MCI(2041) -> PHX(2226) -> SAN |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | MSY(608) -> ATL(944) -> CMH(1156) -> BWI(1354) -> RDU(1553) -> BNA(1729) -> DTW(2045) -> MDW(2149) -> LAX |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | MSY(610) -> ATL(917) -> CMH(1142) -> BWI(1355) -> RDU(1536) -> BNA(1716) -> DTW(2040) -> MDW(2147) -> LAX |

### Notable Observations

- All top-10 itineraries belong to **WN (Southwest Airlines)**, reflecting that carrier's hub-and-spoke multi-stop routing strategy.
- The most repeated flight number in the top 10 is **WN3149 (4 occurrences in top 10)**, suggesting a recurrent scheduled itinerary rather than an ad-hoc assignment.
- Route clustering across dates for the same flight number indicates a **weekly (approximately every 7 days for WN3149 in Jan–Feb 2024)** operating schedule.

_Generated: 2026-03-19T12:23:06Z_