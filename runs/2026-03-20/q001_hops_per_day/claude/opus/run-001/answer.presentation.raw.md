`visual.html` written. The dashboard includes:

- **KPI strip** anchored to the top-ranked result with max hops, tail number, flight, date, carrier, and route repetition context
- **Leaflet map** showing the selected itinerary route, present before enrichment completes with degraded-state messaging if coordinates fail
- **Route sequence panel** with airport codes, names (from enrichment), and departure times
- **Clickable itinerary table** with active-row styling distinct from hover; clicking redraws the map and route detail
- **Single enrichment query** for all airport coordinates, reused across itinerary selections without per-click re-fetch
- **Query ledger** tracking both the primary query and airport coordinate enrichment with expandable SQL
- **Dynamic mode** using `localStorage` JWE token and the `mcp.demo.altinity.cloud` endpoint
- **CSV export** for the result set
