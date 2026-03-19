Find the highest number of hops per day for a single aircraft using the same flight number.

Report Aircraft ID, Flight Number, Carrier, Date, Route.
For the longest trip, show the actual departure time from each origin.
What does the itinerary look like? Find the top 10 longest and most recent itineraries.

Definitions and rules:

- Do not invent column aliases 
- Build each itinerary in chronological leg order using the actual departure timestamps.
- Use actual `DepTime` values only. Exclude rows where actual departure time is missing; do not fall back to `CRSDepTime`.
- The textual `Route` must include every leg and the final destination airport.
- If `Hops = N`, the route must contain exactly `N` legs or `N + 1` airport codes.
