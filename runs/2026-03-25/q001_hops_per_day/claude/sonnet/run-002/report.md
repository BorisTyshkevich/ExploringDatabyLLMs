# Highest daily hops for one aircraft on one flight number

> Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

The top itineraries are overwhelmingly recurring scheduled patterns. The highest-frequency 8-hop routes each appear dozens to over a hundred times across different dates, spanning multi-month windows. For example, the route `FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC` appeared 114 times between March and August 2008, and `CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN` (WN flight 3149) appears across multiple weeks in 2024. These are classic Southwest Airlines through-plane rotation schedules that repeat the same airport sequence day after day on fixed flight numbers.

- Rows returned: 10
- Columns: Route, max_hops, occurrences, first_seen, last_seen

| Route | max_hops | occurrences | first_seen | last_seen |
| --- | --- | --- | --- | --- |
| FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC | 8 | 114 | 2008-03-10T00:00:00Z | 2008-08-21T00:00:00Z |

> What geographic pattern do the top itineraries show?

All top itineraries are Southwest Airlines (WN) cross-country daisy-chains that traverse the full width of the continental United States. Routes typically originate on the East Coast or Gulf Coast (e.g., FLL, PVD, BDL, LGA, HOU, MSY), pass through Midwest hubs (MDW, STL, MCI, CMH), continue through Texas/Southwest hubs (DAL, ABQ, PHX, ELP), and terminate on the West Coast (SJC, OAK, ONT, SEA, RNO). Each 8-hop route covers 6–7 different states. The pattern reflects Southwest's hub-and-spoke-free model where a single aircraft flies a continuous east-to-west (or west-to-east) rotation through its point-to-point network in a single day.

- Rows returned: 10
- Columns: Route, occurrences, state_count, states_sample

| Route | occurrences | state_count | states_sample |
| --- | --- | --- | --- |
| FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC | 99 | 7 | [NM, TX, FL, IN, CA, MO, IL] |
