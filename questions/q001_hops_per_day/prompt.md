Using `default.ontime_v2`, find the highest number of hops per day for a single aircraft using the same flight number.

Report Aircraft ID, Flight Number, Carrier, Date, Route.
For the longest trip, show the actual departure time from each origin.
What does the itinerary look like? Find the top 10 longest and most recent itineraries.

Definitions and rules:

- Keep column names exactly as they exist in `default.ontime_v2`.
- Do not invent column aliases 
- Build each itinerary in chronological leg order using the actual departure timestamps.
- The textual `Route` must include every leg and the final destination airport.
- If `Hops = N`, the route must contain exactly `N` legs or `N + 1` airport codes.
