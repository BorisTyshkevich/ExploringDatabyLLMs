`visual.html` written. Here's what the dashboard includes:

**KPI strip** (anchored to the top-ranked result — N957WN / WN366 / 2024-12-01 / 8 hops): tail number, flight number, date, hop count, and route repetition context derived from result-set grouping.

**Lead itinerary map** (Leaflet):
- Map card is always in the DOM; shows a spinner placeholder before enrichment completes
- Numbered stop markers — red for origin/destination, navy for intermediate stops
- Polylines in `--sky` connecting consecutive stops
- `invalidateSize()` called on init and after each route redraw
- One enrichment query fetches all airport coordinates for all 10 rows at once and caches them; no re-query on row clicks

**Route sequence / detail panel** — stop badges with sequence numbers, leg-by-leg table, updates on selection.

**Itinerary table** — click any row to redraw the map and route panel; active row gets a distinct `--sky` left-border highlight separate from hover state.

**Query ledger** — three entries (Primary, Airport Coordinate Enrichment, Map Rendering), each collapsible to show SQL, with OK / Pending / Failed / Degraded status badges.

**Degraded-state handling** — if enrichment fails or a selected route lacks enough coordinates, the map card stays visible with an amber warning; the ledger entry is updated to Degraded; all other panels continue rendering normally.
