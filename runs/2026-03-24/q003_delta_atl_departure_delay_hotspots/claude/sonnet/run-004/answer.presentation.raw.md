`visual.html` written. Here's what the dynamic dashboard includes:

**Layout (all content derived from live browser queries):**
- **Hero** — EWR 1600–1659 prominently badged as worst hotspot with key metrics
- **KPI strip** — worst avg delay, P90 delay, del-15 rate, credible months (all from primary query result)
- **Heatmap** — destination × time-block grid, cells colored teal→amber→red by avg delay; EWR 1600–1659 outlined as worst cell
- **Persistence chart** — Chart.js dual-axis line chart (avg delay + del-15 rate) over all 228+ qualifying months for EWR 1600–1659
- **Pattern table** — top 15 hotspots with ≥24 credible months, EWR/JFK badges color-coded, CSV export
- **Query ledger** — collapsible entries for all three queries (worst_hotspot primary + persistence + pattern_summary supporting), each showing status and row count
- **Takeaway** — narrative summary derived from live query results

**Dynamic behavior:**
- Auto-runs on page load if a stored JWE token exists
- Runs `persistence` and `pattern_summary` supporting queries concurrently after primary succeeds
- Each supporting section degrades independently on failure without breaking the primary dashboard
- JWE stored in `localStorage['OnTimeAnalystDashboard::auth::jwe']`; Forget button clears it
