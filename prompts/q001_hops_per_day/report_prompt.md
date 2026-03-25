Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.

Return enough itinerary detail to support the dashboard, including:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`


## Dashboard Questions

- Which itinerary is the highest-hop example, and what does it look like?
- Which of the top-ranked itineraries is the most recent?
- Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?
