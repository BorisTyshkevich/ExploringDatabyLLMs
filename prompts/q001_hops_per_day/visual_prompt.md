Create a `visual.html` artifact for this question using the ontime-analyst-dashboard skill in dynamic mode

The page must:

- be a complete browser-openable HTML document
- use a scalable slippy-map approach with Leaflet and realistic remote basemap tiles
- include a map container plus JavaScript that initializes the map
- treat the saved SQL as the primary analytical query and keep it visible in the query controls
- show the lead itinerary returned by the primary query and parse its `Route` string in JavaScript for downstream rendering
- derive hop count, stop sequence, and repeated-route comparisons in JavaScript from the saved result set
- because the primary query lacks coordinates, run an explicit airport-coordinate enrichment query against `default.airports_bts` using airport codes parsed from the lead/top route strings
- show a visible query ledger that lists the primary query and the airport enrichment query with purpose, role, status, target table, and row count
- label the map as airport-coordinate enrichment rather than implying the primary query already contained map geometry
- use the enrichment query results to place airport markers and route lines for the lead itinerary when coordinates are available
- render the primary visual from fields present in the fetched data; do not invent coordinates or geocode city names
- include a headline, a short subtitle, and KPI cards for tail number, flight number, date, hop count, and route repetition context
- include a legend and a route table or route sequence panel below the map
- if airport-coordinate enrichment fails or returns insufficient rows, keep the route-centric overview, report the degraded map in the ledger, and continue rendering the non-map analysis
- remote JS/CSS libraries and remote tile URLs are allowed for this map artifact
