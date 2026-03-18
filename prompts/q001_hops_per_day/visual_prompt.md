The page must:

- use Leaflet with a slippy map and remote basemap tiles
- place the token input, SQL textarea, fetch action, and status text inside a literal `<footer>` element at the very end of the page
- keep the Lead Itinerary Map card and visible map container in the page even before airport-coordinate enrichment succeeds
- treat the first row returned by the primary query as the default selected itinerary on initial load and parse its `Route` string in JavaScript
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- run an explicit airport-coordinate enrichment query against `default.airports_bts` using airport codes parsed from the route strings
- label the map as airport-coordinate enrichment in the query ledger
- reuse the enrichment results for any itinerary selected from the primary result set without issuing a new per-click enrichment query
- use enrichment results to place airport markers and route lines for the currently selected itinerary, but treat that enrichment only as an upgrade to the map
- include KPI cards for tail number, flight number, date, hop count, and route repetition context, with the date shown as its own visible KPI value
- keep the KPI strip anchored to the top-ranked result even when the selected itinerary changes
- include a legend plus both a route sequence/detail panel and an itinerary table below the map
- make itinerary table rows clickable so selecting a row redraws the map and refreshes the route sequence/detail panel for that itinerary
- show a clear active-row state for the selected itinerary that is distinct from simple hover styling
- if enrichment fails or the selected itinerary lacks enough coordinates, keep the map card visible with degraded-state messaging for that selected itinerary, report the degraded map in the ledger, and continue rendering the non-map analysis
- if the user selects a different itinerary after the map has already been initialized, update the map or its degraded-state message for that newly selected itinerary instead of leaving stale map content or a blank map container
