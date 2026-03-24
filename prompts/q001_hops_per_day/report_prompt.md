Find the daily itineraries with the highest number of hops for a single aircraft using the same flight number.

Analyze completed flights across the full available history and focus on the unique top 10 longest itineraries, ordered so the most recent examples are easy to identify. Preserve enough itinerary detail for a dashboard map, route sequence panel, and clickable itinerary table.

Provide one proof query for each dashboard question below. The proof query behind the lead-itinerary question should preserve the ranked itinerary rows needed for the dashboard, not just a single top row.

Keep the proof queries efficient. Avoid rescanning the full fact table independently for every dashboard question when the same ranked itinerary candidate set can be reused or when a targeted follow-up on that candidate set is enough.

Return enough itinerary detail to support:

- aircraft id / tail number
- flight number
- carrier
- flight date
- hop count
- textual route in chronological leg order
- a simple browser-friendly itinerary-sequence field that lets the dashboard redraw the selected route and route detail panel without relying only on the textual `Route`

Keep these runtime constraints:

- Build each itinerary in chronological leg order using the actual departure timestamps.
- The textual `Route` must include every leg and the final destination airport.
- Delimiter in `Route` must be `-` without spaces.

In the answers, directly name the lead itinerary, directly name the most recent top itinerary, state whether the top itineraries look recurring or one-off, and summarize the geographic shape they trace.

## Dashboard Questions

- Which itinerary is the highest-hop example, and what does it look like?
- Which of the top-ranked itineraries is the most recent?
- Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?
- What geographic pattern do the top itineraries show?
