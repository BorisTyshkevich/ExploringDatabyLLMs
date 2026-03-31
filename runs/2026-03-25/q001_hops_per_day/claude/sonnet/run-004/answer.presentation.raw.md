`visual.html` written to `run-004/`. Here is what the dashboard includes:

**Structure**
- **Header** — branded aviation gradient with carrier badge
- **Hero** — narrative anchored permanently to rank-1 itinerary (N957WN · WN366 · ISP→SEA, Dec 1 2024), covering the coast-to-coast Southwest pattern
- **KPI strip** (5 cards) — Tail Number, Flight Number, Flight Date, Hop Count, Recurrence — synced to the selected row on every click
- **Leaflet map** — route polyline in red, numbered circle markers (teal=origin, navy=stops, red=final), bounds fitted, legend, overlay for pre-enrichment state
- **Stop sequence panel** — ordered airport list with city names (populated after enrichment), connector lines
- **Itinerary table** — 10 clickable rows; active row styled with inset left border + blue tint distinct from hover; recurrence mini-bar color-coded by tier (amber/sky/teal/navy)
- **Tier legend** below table
- **Query ledger** — primary query + coordinate enrichment query, both expandable, with status/row counts
- **Footer** — password JWE input, forget button, SQL textarea pre-filled, fetch button with disabled-during-run guard, status line

**Dynamic behavior**
- Auto-runs on load when a stored JWE token is found
- Enrichment query fetches all airport coordinates from `ontime.dim_airports` in one shot after the primary query; coordinates are cached and reused for every row selection without re-querying
- Selecting a row updates KPIs, map title, map route, and route detail panel while the hero stays anchored to rank 1
- If enrichment fails, the map shows a degraded message, the ledger records the failure, and all other panels remain functional
