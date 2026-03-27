### main query

- use main query as the primary source of information for visualizing
- use the other/supporting queries when they materially improve the narrative or supporting panels
- anchor the hero narrative to the top-ranked itinerary even when another itinerary is selected in the table
- show a lead-itinerary map only after airport-coordinate lookup succeeds
- treat the first row returned by the main query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- include a narrative hero about the lead itinerary and the broader geographic pattern of the top itineraries
- label the map as airport-coordinate lookup in the query ledger
- reuse the lookup results for any itinerary selected from the main query result set without issuing a new per-click lookup query
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip synced to the currently selected itinerary
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- make the selected-row map behavior explicit: when the selected itinerary differs from Rank 1, the map title, plotted route, markers, bounds, and route detail panel must visibly update to that selected itinerary rather than leaving the lead route drawn
- keep the map/detail/KPI selection state separate from the anchored hero state
- place selection-driven itinerary detail right after KPI cards, but before map panel.
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- prefer the `Route` value from the main query as the per-row itinerary representation for redraws
- if airport-coordinate lookup fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis
- derive the ordered itinerary sequence for map redraws and the route detail panel by splitting `Route` on `-`

### operational-stress

Which airports or legs are the main operational stress points within the top 10 unique maximum-hop itineraries?
Return per-airport and per-leg average departure delay, average arrival delay, rate of 15-plus-minute delays, and diversion incidence, and identify the stop positions most associated with disruption.

### key connectors

Which airports act as the key connectors, origins, and termini within the top 10 unique maximum-hop itineraries?
Classify airport appearances by route position and return airport code, airport name, city/state, total appearances, origin appearances, intermediate-stop appearances, final-destination appearances, and share of itineraries containing that airport.

### geographically extreme

How geographically extreme is each of the top 10 unique maximum-hop itineraries?
Return total flown distance, unique airports, unique city markets, unique states, unique local-time offsets, and whether the route is entirely domestic, then summarize which routes are the most geographically expansive.
