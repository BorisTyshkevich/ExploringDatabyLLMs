`visual.html` written. Here's what the dashboard delivers:

**Queries (3 total, all run in-browser)**
| # | Label | Purpose |
|---|-------|---------|
| 1 | Primary Query | Top-10 itineraries by hop count |
| 2 | Airport Coordinate Enrichment | Single batch query for all unique airports across all 10 itineraries; cached for all row selections |
| 3 | Route Recurrence Pattern | How many times each flight number appeared in the top-10 |

**Layout**
- **Hero card** — narrative title anchored to Rank 1 (WN366 · N957WN · Dec 1 2024 — 8-hop marathon)
- **5 KPI cards** — Tail / Flight # / Date / Hop Count / Route Occurrences, always Rank 1
- **Leaflet map** — polyline + color-coded markers (green=origin, blue=intermediate, red=destination); lead route in orange `#e8612e`, others in sky blue; redraws on row click without re-querying
- **Route sequence panel** — numbered leg list + recurrence context, updates on row selection
- **Itinerary table** — all 10 rows, clickable; active row has blue left-border + tinted background; hover is lighter
- **Query ledger** — collapsed by default; each entry shows name / role / status / row count / expandable SQL

**Key behaviors**
- Hero + KPI strip stay anchored to Rank 1 regardless of table selection
- Map title, polyline, markers, and bounds all update to the selected row
- If airport coords are missing for a selected itinerary, the map card stays visible with a degraded-state message
- JWE auto-loaded from `localStorage['OnTimeAnalystDashboard::auth::jwe']`; Forget button clears it
