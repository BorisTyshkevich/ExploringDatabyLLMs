# Highest daily hops for one aircraft on one flight number

- Rows returned: 10
- Generated at: 2026-03-20T14:04:31Z
- Columns: Aircraft ID, Flight Number, Carrier, Date, Hops, Route, Maximum Hops Observed, Maximum-Hop Itinerary Count
- First row snapshot: Aircraft ID=N957WN, Flight Number=366, Carrier=WN

The key derived value is maximum value in `Maximum Hops Observed`.

| Aircraft ID | Flight Number | Carrier | Date | Hops | Route | Maximum Hops Observed | Maximum-Hop Itinerary Count |
| --- | --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 05:43 ISP->BWI \| 08:10 BWI->MYR \| 10:20 MYR->BNA \| 11:42 BNA->VPS \| 14:01 VPS->DAL \| 16:43 DAL->LAS \| 18:28 LAS->OAK \| 20:41 OAK->SEA | 8 | 9859 |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | 8 | 06:21 CLE->BNA \| 08:01 BNA->PNS \| 10:07 PNS->HOU \| 12:34 HOU->MCI \| 15:14 MCI->PHX \| 17:47 PHX->BUR \| 19:02 BUR->OAK \| 21:17 OAK->DEN | 8 | 9859 |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | 8 | 06:18 CLE->BNA \| 08:00 BNA->PNS \| 10:12 PNS->HOU \| 12:39 HOU->MCI \| 15:15 MCI->PHX \| 17:56 PHX->BUR \| 18:59 BUR->OAK \| 20:59 OAK->DEN | 8 | 9859 |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | 8 | 06:20 CLE->BNA \| 08:10 BNA->PNS \| 10:17 PNS->HOU \| 12:37 HOU->MCI \| 15:10 MCI->PHX \| 18:00 PHX->BUR \| 19:03 BUR->OAK \| 21:02 OAK->DEN | 8 | 9859 |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | 8 | 06:25 CLE->BNA \| 08:05 BNA->PNS \| 10:11 PNS->HOU \| 12:41 HOU->MCI \| 16:07 MCI->PHX \| 19:09 PHX->BUR \| 20:06 BUR->OAK \| 21:56 OAK->DEN | 8 | 9859 |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | 8 | 06:27 ELP->DAL \| 09:39 DAL->LIT \| 11:24 LIT->ATL \| 14:34 ATL->RIC \| 16:56 RIC->MDW \| 18:31 MDW->MCI \| 20:40 MCI->PHX \| 22:46 PHX->SAN | 8 | 9859 |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | 8 | 06:22 ELP->DAL \| 09:41 DAL->LIT \| 11:22 LIT->ATL \| 14:33 ATL->RIC \| 16:59 RIC->MDW \| 18:32 MDW->MCI \| 20:41 MCI->PHX \| 22:26 PHX->SAN | 8 | 9859 |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | 8 | 06:08 MSY->ATL \| 09:44 ATL->CMH \| 11:56 CMH->BWI \| 13:54 BWI->RDU \| 15:53 RDU->BNA \| 17:29 BNA->DTW \| 20:45 DTW->MDW \| 21:49 MDW->LAX | 8 | 9859 |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | 8 | 06:10 MSY->ATL \| 09:17 ATL->CMH \| 11:42 CMH->BWI \| 13:55 BWI->RDU \| 15:36 RDU->BNA \| 17:16 BNA->DTW \| 20:40 DTW->MDW \| 21:47 MDW->LAX | 8 | 9859 |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | 8 | 06:01 MSY->ATL \| 09:18 ATL->CMH \| 11:41 CMH->BWI \| 13:47 BWI->RDU \| 15:39 RDU->BNA \| 17:24 BNA->DTW \| 21:12 DTW->MDW \| 22:42 MDW->LAX | 8 | 9859 |

The maximum-hop pattern should be characterized as interpretation of `Maximum-Hop Itinerary Count` as a recurring pattern or a one-off itinerary. The single most recent itinerary among the maximum-hop rows is carrier, flight number, date, and route from the first row in the result set. Across the top 10 longest itineraries, the main route clustering signal is short description of the dominant repeated route or flight-number cluster visible across the top 10 rows.