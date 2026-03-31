The dashboard is ready at `visual.html`. Key features:

- **Dynamic mode** — fetches live data via the JWE-authenticated ClickHouse endpoint
- **Leaflet map** — always present; shows the selected itinerary's route with numbered segments, origin/destination highlighted in red, intermediate stops in navy
- **Airport coordinate enrichment** — single bulk query against `ontime.dim_airports` for all airports across all 10 rows; cached and reused on row selection (no per-click queries)
- **KPI strip** — anchored to the top-ranked result: Aircraft, Flight Number, Date, Hop Count, Route Repetition context
- **Route sequence panel** — visual stop-by-stop display with airport code badges and arrows; updates on row click
- **Clickable itinerary table** — active-row highlighting distinct from hover; clicking redraws the map and route panel
- **Degraded map handling** — if enrichment fails or a route lacks coordinates, the map card stays visible with an explanatory message
- **Query ledger** — primary query + enrichment query, each with expand/collapse SQL, status, and row count
- **CSV export** and **footer controls** (JWE input, forget, SQL textarea, fetch button)
- Auto-runs on load if a stored JWE token is found
