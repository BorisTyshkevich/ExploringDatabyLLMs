# Highest daily hops for one aircraft on one flight number

- Rows returned: 10
- Generated at: 2026-03-23T13:19:43Z
- Columns: Aircraft ID, Flight Number, Carrier, Date, Route, hop_count
- First row snapshot: Aircraft ID=N957WN, Flight Number=366, Carrier=WN

| Aircraft ID | Flight Number | Carrier | Date | Route | hop_count |
| --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | ISP→BWI→MYR→BNA→VPS→DAL→LAS→OAK→SEA | 8 |
| N7835A | 3149 | WN | 2024-02-18T00:00:00Z | CLE→BNA→PNS→HOU→MCI→PHX→BUR→OAK→DEN | 8 |
| N429WN | 3149 | WN | 2024-01-28T00:00:00Z | CLE→BNA→PNS→HOU→MCI→PHX→BUR→OAK→DEN | 8 |
| N228WN | 3149 | WN | 2024-01-21T00:00:00Z | CLE→BNA→PNS→HOU→MCI→PHX→BUR→OAK→DEN | 8 |
| N569WN | 3149 | WN | 2024-01-14T00:00:00Z | CLE→BNA→PNS→HOU→MCI→PHX→BUR→OAK→DEN | 8 |
| N7742B | 154 | WN | 2023-04-30T00:00:00Z | ELP→DAL→LIT→ATL→RIC→MDW→MCI→PHX→SAN | 8 |
| N929WN | 154 | WN | 2023-04-16T00:00:00Z | ELP→DAL→LIT→ATL→RIC→MDW→MCI→PHX→SAN | 8 |
| N8631A | 2787 | WN | 2022-10-23T00:00:00Z | MSY→ATL→CMH→BWI→RDU→BNA→DTW→MDW→LAX | 8 |
| N8809L | 2787 | WN | 2022-10-02T00:00:00Z | MSY→ATL→CMH→BWI→RDU→BNA→DTW→MDW→LAX | 8 |
| N8811L | 2787 | WN | 2022-09-25T00:00:00Z | MSY→ATL→CMH→BWI→RDU→BNA→DTW→MDW→LAX | 8 |

The result set surfaces the most recent examples among the highest-hop same-flight-number aircraft-day itineraries, which helps distinguish repeat operating patterns from one-off long sequences.