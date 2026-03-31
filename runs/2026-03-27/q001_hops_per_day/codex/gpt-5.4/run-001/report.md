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

Across the most recent five years in the OnTime data, the longest same-aircraft, same-flight-number daily itineraries reached 8 hops. The 10 most recent unique routes at that maximum include ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA on 2024-12-01 (route recurrence 1 day), CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN on 2024-02-18 (4 days), and ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN on 2023-04-30 (2 days).

- Rows returned: 10
- Columns: aircraft_id, flight_number, carrier, flight_date, hop_count, route_recurrence_count, Route

| aircraft_id | flight_number | carrier | flight_date | hop_count | route_recurrence_count | Route |
| --- | --- | --- | --- | --- | --- | --- |
| N957WN | 366 | WN | 2024-12-01T00:00:00Z | 8 | 1 | ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA |
