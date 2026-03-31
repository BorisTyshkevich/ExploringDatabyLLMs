Both files written. The query:

- Collects all non-cancelled legs with a valid tail number and flight number
- Groups by aircraft, flight number, carrier, and date, assembling legs into an array of tuples
- Sorts legs by scheduled departure time and builds the route string (all origins plus the final destination, joined by `-`)
- Returns the top 10 by hop count (descending), breaking ties by most recent date

The results show **8 hops** as the maximum, exclusively by **Southwest Airlines (WN)**, reflecting their multi-stop point-to-point routing model.
