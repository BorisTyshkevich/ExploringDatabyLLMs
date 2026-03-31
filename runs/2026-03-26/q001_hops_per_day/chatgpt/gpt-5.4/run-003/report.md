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

Within the most recent 5 years of available data, the highest observed same-aircraft, same-flight-number itineraries reached 8 hops after deduplicating same-day leg rows. The most recent top route was WN 366 on 2024-12-01 by aircraft N957WN, flying ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA. Among the top 10 unique 8-hop routes, recurrence ranged from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN recurring most often at 46 days and MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC next at 40 days.

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, route_recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | route_recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

Across the top 10 unique maximum-hop itineraries, the strongest connector airports were DAL, DEN, LAS, MSY, OAK, and BWI, each appearing in 5 of the 10 itineraries. DAL was purely an intermediate stop in all 5 appearances, while BWI and MSY were the most common origins at 2 appearances each. The most common termini were OAK and LAX with 2 final-destination appearances each. BWI, DAL, DEN, LAS, MSY, and OAK each appeared in 50% of the itineraries.

- Rows returned: 45
- Columns: airport_code, airport_name, city_state, total_appearances, origin_appearances, intermediate_stop_appearances, final_destination_appearances, share_of_itineraries_containing_airport

| airport_code | airport_name | city_state | total_appearances | origin_appearances | intermediate_stop_appearances | final_destination_appearances | share_of_itineraries_containing_airport |
| --- | --- | --- | --- | --- | --- | --- | --- |
| BWI | Baltimore/Washington International Thurgood Marshall | Baltimore, MD, MD | 5 | 2 | 3 | 0 | 0.5 |

> How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.

All top 10 unique maximum-hop itineraries were entirely domestic. The most geographically expansive by total flown distance was BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK at 5,169 miles, spanning 9 airports, 8 city markets, 9 states, and 3 local-time offsets. Several routes were broader in time-zone spread, with CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN, HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK each touching 4 distinct local-time offsets.

- Rows returned: 10
- Columns: Route, aircraft_id, flight_number, carrier, flight_date, hop_count, total_flown_distance, unique_airports, unique_city_markets, unique_states, unique_local_time_offsets, entirely_domestic

| Route | aircraft_id | flight_number | carrier | flight_date | hop_count | total_flown_distance | unique_airports | unique_city_markets | unique_states | unique_local_time_offsets | entirely_domestic |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | N225WN | 3530 | WN | 2021-08-08T00:00:00Z | 8 | 5169 | 9 | 8 | 9 | 3 | 1 |
