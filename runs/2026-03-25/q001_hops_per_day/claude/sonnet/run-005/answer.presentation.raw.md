`visual.html` written. Here's what's in the dashboard:

**Layout (top to bottom)**
- **Header** — navy gradient with title and subtitle
- **Hero** — always anchored to Rank 1 (N957WN / WN 366 / ISP→SEA / 2024-12-01); includes Southwest context narrative
- **KPI strip** — 5 cards (Tail Number, Flight Number, Date, Hop Count, Route Recurrence); syncs to selected table row
- **Map + Route sequence** (side by side):
  - Leaflet map initialized immediately after primary query; route drawn after enrichment; map title updates with selection
  - Route sequence panel shows numbered airport list (green=origin, red=dest, blue=intermediate) with airport names populated from enrichment cache
- **Itinerary table** — 10 rows, clickable; active row distinctly highlighted in blue; selecting Rank 2–10 updates map, KPIs, route detail, and map title — hero stays locked to Rank 1
- **Query ledger** — expandable rows for Primary query + Airport Coordinate Enrichment (role, status, row count, full SQL on expand)
- **Footer** — JWE password field, Forget button, SQL textarea (pre-filled with primary query), Run button, status line

**Key behaviors**
- Auto-runs on load if JWE stored in `localStorage['OnTimeAnalystDashboard::auth::jwe']`
- Enrichment airport coordinates cached once in `airportCoords{}`, reused for all row selections with no new queries
- Map degraded state (message card) shown if enrichment fails or selected route lacks coordinates; rest of dashboard continues rendering
- Recurrence tiers color-coded: amber=one-off (≤2), sky=scheduled (5–7), teal=recurring (12+)
