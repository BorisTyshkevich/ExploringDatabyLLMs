`visual.html` written. Here's what the dashboard delivers:

**Layout:**
- **Header** — navy aviation gradient with Delta ATL framing
- **Hero card** — red card spotlighting EWR / 1600–1659 with avg delay, p90, del15%, and flight count pulled live from the primary query
- **KPI strip** — 4 cards: worst avg delay, NY metro slot count, total flights in top-15, mean del15+ rate
- **Heatmap** — destination × time-block grid, cells colored from light-blue → amber → red by avg delay; worst cell gets a white outline ring; hover tooltip shows all 4 metrics per cell; Export CSV button
- **Persistence chart** — SVG line chart of EWR/1600-1659 monthly avg delay 1991–2024 (enrichment query), with peak spike annotated
- **Pattern summary table** — destination rollup from the pattern_summary query, EWR/JFK rows highlighted
- **Takeaway** — live-derived narrative about NY metro concentration percentage
- **Query ledger** — all 3 queries listed with status, row count, and collapsible SQL (▶/▼ toggle)
- **Footer** — JWE token (password field), SQL textarea, Run/Forget controls, status message

**Dynamic behavior:** auto-loads on page open if a stored token is found; enrichment query failures degrade only their own panel; stale run protection via `runId` guard.
