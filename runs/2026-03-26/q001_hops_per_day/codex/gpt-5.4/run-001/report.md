# Highest daily hops for one aircraft on one flight number

> Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.
Count hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

The maximum observed same-aircraft, same-flight-number itinerary length is 8 hops. The 10 most recent unique max-hop routes are all Southwest (WN) itineraries, with the newest on 2024-12-01: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA flown by N957WN on flight 366. Across these 10 rows, route recurrence ranges from 1 day to 46 days, led by LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN (46), MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC (40), HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK (20), and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK (12).

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, route_recurrence_count, route

| aircraft_id | flight_number | carrier | flight_date | hop_count | route_recurrence_count | route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

Recurrence counts for the 10 routes are: ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA = 1; CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN = 4; ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN = 2; MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX = 5; MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC = 40; LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN = 46; SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS = 7; HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK = 20; BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK = 12; BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX = 5. Overall, these top itineraries skew toward recurring patterns rather than pure one-offs: 4 are recurring (10+ days), 5 are occasional (2-9 days), and 1 is a one-off, which sums to all 10 routes.

- Rows returned: 10
- Columns: route, route_recurrence_count, recurrence_tier, tier_count

| route | route_recurrence_count | recurrence_tier | tier_count |
| --- | --- | --- | --- |
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 | one-off | 1 |

> What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.

The geographic pattern is entirely domestic U.S. and strongly tied to Southwest/West Coast plus Sun Belt and Mid-Atlantic nodes. Every airport in the 10 routes is a US airport, and the most repeated points are BWI, DAL, DEN, LAS, MSY, and OAK with 5 mentions each, followed by BNA and PHX with 4 each. California is especially prominent through OAK, BUR, LAX, SAN, PSP, SJC, and SMF, so these longest same-flight-number days look like long domestic zig-zag chains linking eastern and southeastern stops into western endpoints and western connection banks.

- Rows returned: 45
- Columns: airport_code, state_code, country_code, geo_bucket, airport_mentions

| airport_code | state_code | country_code | geo_bucket | airport_mentions |
| --- | --- | --- | --- | --- |
| BWI | MD | US | east | 5 |
