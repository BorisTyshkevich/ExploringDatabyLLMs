`visual.html` written. Here's what the dashboard delivers:

**Layout**
- Branded navy header + subtitle
- Hero narrative (anchored to Rank 1 — WN 366, 2024-12-01, N957WN, ISP→SEA)
- 5-card KPI strip: Tail Number, Flight Number, Date, Hop Count, Route Repetition — all anchored to Rank 1 regardless of table selection
- Two-column: Leaflet route map + Route Sequence/Detail panel
- Top-10 itinerary table with clickable rows and a distinct `active-row` outline state
- Query ledger with expand/collapse SQL (▶/▼)

**Map behavior**
- Present in initial HTML before enrichment
- Runs one airport-coordinate enrichment query (`dim_airports` lookup) after primary data loads; results cached for all subsequent row selections — no per-click queries
- On row click: redraws route polyline + markers, refreshes detail panel; KPIs stay anchored
- Degraded-state card (with explanation) if enrichment fails or selected route lacks enough coordinates

**Dynamic mode**
- JWE stored in `localStorage['OnTimeAnalystDashboard::auth::jwe']`; auto-runs on load if token present
- Footer contains token input (password field), Forget button, SQL textarea, Fetch button, status bar
- Endpoint: `https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=...`
- Ledger tracks both queries (Primary + Enrichment) with label, role, status, row count, expandable SQL
