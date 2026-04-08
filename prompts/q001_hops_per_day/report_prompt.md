### main
Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
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

### q1
Which airports or legs are the main operational stress points within the top 10 unique maximum-hop itineraries?

Return per-airport and per-leg: average departure delay, average arrival delay, rate of 15-plus-minute delays, diversion incidence, and the stop positions most associated with disruption.

### q2
Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?

Classify airport appearances by route position and return: airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

### q3
How geographically extreme is each of the top 10 unique maximum-hop itineraries?

Return: total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic. Summarize which routes are the most geographically expansive.
