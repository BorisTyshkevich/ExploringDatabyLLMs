# Highest daily hops for one aircraft on one flight number

- Rows returned: 10
- Generated at: 2026-03-20T12:01:47Z
- Columns: Aircraft ID, Flight Number, Carrier, Date, Hops, Route, Departure Times From Origin, Max Hops Overall, Max Hop Itinerary Count, Same-Hops Route Count, Route Frequency In Top 10
- First row snapshot: Aircraft ID=N957WN, Flight Number=366, Carrier=WN

The key derived value is the maximum observed daily hop count in `Hops`.

| Aircraft ID | Flight Number | Carrier | Date | Hops | Route | Departure Times From Origin | Max Hops Overall | Max Hop Itinerary Count | Same-Hops Route Count | Route Frequency In Top 10 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA | 05:43 ISP \| 08:10 BWI \| 10:20 MYR \| 11:42 BNA \| 14:01 VPS \| 16:43 DAL \| 18:28 LAS \| 20:41 OAK | 8 | 9859 | 1 | 1 |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:21 CLE \| 08:01 BNA \| 10:07 PNS \| 12:34 HOU \| 15:14 MCI \| 17:47 PHX \| 19:02 BUR \| 21:17 OAK | 8 | 9859 | 4 | 4 |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:18 CLE \| 08:00 BNA \| 10:12 PNS \| 12:39 HOU \| 15:15 MCI \| 17:56 PHX \| 18:59 BUR \| 20:59 OAK | 8 | 9859 | 4 | 4 |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:20 CLE \| 08:10 BNA \| 10:17 PNS \| 12:37 HOU \| 15:10 MCI \| 18:00 PHX \| 19:03 BUR \| 21:02 OAK | 8 | 9859 | 4 | 4 |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN | 06:25 CLE \| 08:05 BNA \| 10:11 PNS \| 12:41 HOU \| 16:07 MCI \| 19:09 PHX \| 20:06 BUR \| 21:56 OAK | 8 | 9859 | 4 | 4 |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | 06:27 ELP \| 09:39 DAL \| 11:24 LIT \| 14:34 ATL \| 16:56 RIC \| 18:31 MDW \| 20:40 MCI \| 22:46 PHX | 8 | 9859 | 2 | 2 |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN | 06:22 ELP \| 09:41 DAL \| 11:22 LIT \| 14:33 ATL \| 16:59 RIC \| 18:32 MDW \| 20:41 MCI \| 22:26 PHX | 8 | 9859 | 2 | 2 |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:08 MSY \| 09:44 ATL \| 11:56 CMH \| 13:54 BWI \| 15:53 RDU \| 17:29 BNA \| 20:45 DTW \| 21:49 MDW | 8 | 9859 | 5 | 3 |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:10 MSY \| 09:17 ATL \| 11:42 CMH \| 13:55 BWI \| 15:36 RDU \| 17:16 BNA \| 20:40 DTW \| 21:47 MDW | 8 | 9859 | 5 | 3 |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | 8 | MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX | 06:01 MSY \| 09:18 ATL \| 11:41 CMH \| 13:47 BWI \| 15:39 RDU \| 17:24 BNA \| 21:12 DTW \| 22:42 MDW | 8 | 9859 | 5 | 3 |

The maximum-hop result should be interpreted as whether the maximum-hop result is repeated or one-off based on `Max Hop Itinerary Count` and the repeated routes visible among the returned rows.

The single most recent itinerary among the maximum-hop rows is the carrier, flight number, date, aircraft id, route, and departure timeline from the most recent row where `Hops` equals `Max Hops Overall`.

Across the top 10 longest itineraries, notable route repetition or clustering is the strongest route repetition or clustering pattern visible across the top 10 rows, using `Route Frequency In Top 10` and `Same-Hops Route Count`.