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

Within the last five years of OnTime data, the highest observed same-aircraft, same-flight-number itineraries reached 8 hops. The most recent unique 8-hop route was ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01 by WN flight 366 with aircraft N957WN. Among the top 10 unique routes, recurrence ranges from 1 day to 46 days, led by LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN at 46 days and MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC at 40 days.

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, route_recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | route_recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

The key connector airports in the top 10 routes are DAL, DEN, LAS, BWI, MSY, and OAK, each appearing in 5 of 10 itineraries. DAL is entirely an intermediate connector with 5 intermediate appearances, while DEN and LAS split between connector and terminal roles. The most common final destinations are OAK and LAX, with 2 terminal appearances each.

- Rows returned: 45
- Columns: airport_code, airport_name, city_name, state_code, total_appearances, origin_appearances, intermediate_stop_appearances, final_destination_appearances, share_of_itineraries

| airport_code | airport_name | city_name | state_code | total_appearances | origin_appearances | intermediate_stop_appearances | final_destination_appearances | share_of_itineraries |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DAL | Dallas Love Field | Dallas, TX | TX | 5 | 0 | 5 | 0 | 0.5 |

> How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.

The most geographically expansive route by flown distance is BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX at 5,169 miles, spanning 9 airports, 8 city markets, 9 states, and 3 local-time offsets. The widest time-zone spread is 4 offsets, reached by CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN, ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN, HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK, MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC, and BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK. All top 10 routes are entirely domestic.

- Rows returned: 10
- Columns: t.Route, aircraft_id, flight_number, carrier, flight_date, hop_count, total_flown_distance, unique_airports, unique_city_markets, unique_states, unique_local_time_offsets, entirely_domestic

| t.Route | aircraft_id | flight_number | carrier | flight_date | hop_count | total_flown_distance | unique_airports | unique_city_markets | unique_states | unique_local_time_offsets | entirely_domestic |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | N225WN | 3530 | WN | 2021-08-08T00:00:00Z | 8 | 5169 | 9 | 8 | 9 | 3 | 1 |

> Which airports or legs are the main operational stress points within the top 10 unique maximum-hop itineraries?
Return per-airport and per-leg average departure delay, average arrival delay, rate of 15-plus-minute delays, and diversion incidence, and identify the stop positions most associated with disruption.

The main operational stress points are late-route segments, especially leg 8, which averages 11.7 minutes of departure delay, 5.8 minutes of arrival delay, and a 50% 15-plus-minute delay rate. At the airport level, RNO, OAK, COS, and DAL show the highest outbound stress in this top-10 set, while legs MDW-LAX, OAK-RNO, DAL-LAX, MSY-ATL, PHX-SAN, RNO-LAS, COS-DEN, and BNA-DTW each show a 100% 15-plus-minute delay rate. No diversions appear in these representative itineraries.

- Rows returned: 124
- Columns: entity_type, entity_key, entity_label, stop_position, avg_departure_delay, avg_arrival_delay, delay_15_plus_rate, diversion_incidence

| entity_type | entity_key | entity_label | stop_position | avg_departure_delay | avg_arrival_delay | delay_15_plus_rate | diversion_incidence |
| --- | --- | --- | --- | --- | --- | --- | --- |
| airport | RNO | departures from RNO |  | 30 | 20 | 1 | 0 |
