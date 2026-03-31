The dashboard has been written to `visual.html`. It includes:

- **Header** with title and analytical subtitle
- **Hero narrative** anchored to the rank-1 itinerary (N957WN, WN 366, ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA)
- **5 KPI cards** — Aircraft, Flight Number, Flight Date, Hop Count, Route Recurrence — all synced to the currently selected row
- **Leaflet map** showing the selected itinerary's route with origin (teal), destination (red), and intermediate stops (navy), plus dashed leg lines
- **Route sequence panel** alongside the map listing each stop with airport names from enrichment
- **Clickable itinerary table** — selecting a row redraws the map, updates KPIs, and refreshes the route detail panel; rank-1 has distinct styling and selected row has a highlighted state
- **Export CSV** button
- **Query ledger** with expand/collapse SQL for both the primary query and the airport coordinate enrichment query
- **Footer controls** with JWE token input, forget button, read-only SQL textarea, and fetch button

The map degrades gracefully with a message if coordinates are unavailable. The hero narrative remains anchored to rank 1 regardless of selection. Auto-loads if a stored JWE token exists.
