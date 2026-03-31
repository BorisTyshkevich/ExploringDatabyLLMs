# Highest daily hops for one aircraft on one flight number

> Which itinerary is the highest-hop example, and what does it look like?

The highest-hop example is an 8-hop Southwest itinerary flown by aircraft N957WN as flight 366 on 2024-12-01, with route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.

- Rows returned: 1
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | Route |
| --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Which of the top-ranked itineraries is the most recent?

The most recent top-ranked itinerary is also the latest 8-hop unique route: N957WN operating Southwest flight 366 on 2024-12-01 over ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.

- Rows returned: 1
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | Route |
| --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

They look mostly like recurring patterns rather than one-offs. Among the 10 most recent unique 8-hop routes, 9 recur on multiple dates, averaging 14.1 occurrences each, although the single most recent example appears only once.

- Rows returned: 1
- Columns: route_count, one_off_routes, recurring_routes, avg_occurrences, max_occurrences

| route_count | one_off_routes | recurring_routes | avg_occurrences | max_occurrences |
| --- | --- | --- | --- | --- |
| 10 | 1 | 9 | 14.1 | 46 |
