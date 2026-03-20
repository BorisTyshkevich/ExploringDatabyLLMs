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

Report guidance:

Explain:

- the maximum hop count observed and whether it appears to be a repeated operating pattern or a one-off itinerary,
- the single most recent itinerary among the maximum-hop rows, including carrier, flight number, date, and full route,
- and any notable route repetition or clustering visible across the top 10 longest itineraries.
