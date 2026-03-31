`visual.html` written. Here's a summary of what was built:

---

**Dashboard: Highest Daily Hops per Aircraft**

**Queries (3 total, all visible in the ledger):**
| Query | Role | Purpose |
|---|---|---|
| Top Itineraries by Hop Count | Primary | The provided SQL, LIMIT 10 |
| Airport Coordinate Enrichment | Enrichment | `dim_airports` lat/lon for all airports across all 10 routes — cached once, reused on row selection |
| Route Occurrence Counts | Supporting | Counts how many times each of the 10 exact routes was flown — powers the Route Occurrences KPI |

**Layout:**
- **Hero narrative** — anchored to Rank 1, describes the lead itinerary's cross-country daisy-chain pattern
- **KPI strip (5 cards)** — Tail Number, Flight Number + carrier name, Flight Date, Max Daily Hops, Route Occurrences (loads async from supporting query) — always anchored to Rank 1 regardless of selected row
- **Map + Route detail panel** — Leaflet map with route polyline and airport markers; updates on row click; degrades gracefully if enrichment fails; map card present in DOM from initial load
- **Itinerary table** — clickable rows with clear selected-row state (separate from KPI anchor); Occurrences column fills in after supporting query completes
- **Query ledger** — expand/collapse per entry, status color-coded
- **Footer** — JWE token (password field), SQL textarea, Fetch button, Forget token
