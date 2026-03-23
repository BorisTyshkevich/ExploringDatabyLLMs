Find routes with the highest number of hops per day for a single aircraft using the same flight number.

- What does the itinerary look like? 
- Find the top 10 longest and most recent itineraries.
- Report Aircraft ID, Flight Number, Carrier, Date, Route, hop_count.

rules:

- return column names exactly as requested. Do not invent other column aliases.
- Build each itinerary in chronological leg order using the actual departure timestamps.
- The textual `Route` must include every leg and the final destination airport.
- Delimiter in `Route` must be `-` without spaces. 
