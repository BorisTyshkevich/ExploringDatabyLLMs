Create a `visual.html` artifact for this question using the ontime-analyst-dashboard skill in dynamic mode

The page must:

- be a complete browser-openable HTML document
- use a scalable slippy-map approach with Leaflet and realistic remote basemap tiles
- include a map container plus JavaScript that initializes the map
- show the full itinerary for the highest-hop itinerary returned by the query
- highlight the single longest leg within that itinerary as the primary emphasized route
- render the map in a realistic atlas-like style using the basemap for state borders and labels
- place leg numbers (`1`, `2`, ...) directly on the map for each segment
- include a headline, a short subtitle, and KPI cards for tail number, flight number, date, hop count, and longest-leg distance
- include a legend and a route table below the map
- embed itinerary-leg data and airport-coordinate data as CSV inside `<script type="text/csv">` tags
- parse the embedded CSV with a small inline parser written in plain JavaScript
- derive the longest leg in JavaScript from the embedded data, not by hardcoding the highlighted segment
- add airport markers with tooltips or popups that show airport code, city, and departure context
- remote JS/CSS libraries and remote tile URLs are allowed for this map artifact
