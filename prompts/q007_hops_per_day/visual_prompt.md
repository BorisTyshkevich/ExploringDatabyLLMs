The page must:

- show a lead-itinerary map using embedded airport coordinates from the static artifact
- treat the first row returned by the primary query as the default selected itinerary on initial load
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- reuse embedded airport metadata for any itinerary selected from the primary result set without issuing browser-side data fetches
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip anchored to the top-ranked result even when the selected itinerary changes
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- if the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary and continue rendering the non-map analysis
