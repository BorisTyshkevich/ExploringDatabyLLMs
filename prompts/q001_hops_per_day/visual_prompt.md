The page must:

- use Leaflet with a slippy map and remote basemap tiles
- keep the Lead Itinerary Map card and visible map container in the page even before airport-coordinate enrichment succeeds
- show the lead itinerary returned by the primary query and parse its `Route` string in JavaScript
- derive hop count, stop sequence, and repeated-route comparisons from the result set
- run an explicit airport-coordinate enrichment query against `default.airports_bts` using airport codes parsed from the route strings
- label the map as airport-coordinate enrichment in the query ledger
- use enrichment results to place airport markers and route lines for the lead itinerary, but treat that enrichment only as an upgrade to the map
- include KPI cards for tail number, flight number, date, hop count, and route repetition context
- include a legend and a route table or route sequence panel below the map
- if enrichment fails, keep the map card visible with degraded-state messaging, report the degraded map in the ledger, and continue rendering the non-map analysis
