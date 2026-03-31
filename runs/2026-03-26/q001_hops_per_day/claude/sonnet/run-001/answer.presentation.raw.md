`visual.html` written. Here's what the dashboard delivers:

**Layout**
- **KPI strip** (5 cards): Aircraft tail, Flight #, Date, Hop count, Recurrence — all synced to the selected itinerary
- **Hero narrative** — anchored to Rank 1 (N957WN / WN 366 / Dec 1 2024), stays fixed even when another row is selected
- **Leaflet route map** — polyline + airport markers; overlay shows loading/degraded state before/during enrichment
- **Route sequence panel** — numbered stop list with airport names from enrichment cache
- **Itinerary table** — 10 clickable rows; active row has left-border highlight distinct from hover; clicking redraws map, detail, and KPI strip
- **Query ledger** — two entries (primary + enrichment) with ▶/▼ expand toggle
- **CSV export**

**Interactions**
- Selecting any row updates the map title, polyline, markers, bounds, and route detail panel for that itinerary; the hero text stays on Rank 1
- Enrichment runs once; coordinates are cached and reused for all subsequent row selections without additional queries
- If enrichment fails, the map card stays visible with a degraded-state message; all other panels continue rendering
- Auto-runs on load if a stored JWE token is found
