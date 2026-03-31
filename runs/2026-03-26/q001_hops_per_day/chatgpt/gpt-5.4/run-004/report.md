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
- route recurrence count: total number of days across the analyzed window on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

Across the most recent five years ending 2025-11-30, the longest same-aircraft, same-flight-number itineraries all reached 8 hops. The most recent top 10 unique routes are all Southwest (WN) examples, led by ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01; recurrence across the window ranges from 1 day to 46 days, with LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN the most recurrent of the top 10.

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, FlightDate, hop_count, route_recurrence_count, Route

| aircraft_id | flight_number | carrier | FlightDate | hop_count | route_recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

Dallas Love Field, Denver, Las Vegas, Baltimore/Washington, New Orleans, and Oakland are the key nodes, each appearing in 5 of the 10 itineraries. DAL is the strongest pure connector with 5 intermediate-stop appearances; BWI and MSY are the most common origins at 2 each; OAK and LAX are the most common final destinations at 2 each.

- Rows returned: 45
- Columns: airport_code, airport_name, city_state, total_appearances, origin_appearances, intermediate_stop_appearances, final_destination_appearances, share_of_itineraries_containing_airport

| airport_code | airport_name | city_state | total_appearances | origin_appearances | intermediate_stop_appearances | final_destination_appearances | share_of_itineraries_containing_airport |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DAL | Dallas Love Field | Dallas, TX | 5 | 0 | 5 | 0 | 0.5 |

> How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.

All 10 top routes are entirely domestic. The most geographically expansive by flown distance is BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX at 5,169 miles; the broadest state coverage is 9 states on BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX, MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX, and ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA; five routes span 4 local-time offsets, while every route touches 9 unique airports.

- Rows returned: 10
- Columns: Route, total_flown_distance, unique_airports, unique_city_markets, unique_states, unique_local_time_offsets, entirely_domestic

| Route | total_flown_distance | unique_airports | unique_city_markets | unique_states | unique_local_time_offsets | entirely_domestic |
| --- | --- | --- | --- | --- | --- | --- |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5169 | 9 | 8 | 9 | 3 | Yes |
