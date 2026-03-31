`visual.html` written. Here's what the dashboard includes:

**KPI strip (4 cards)**
- Years Analyzed — derived from yearly ranks data with min/max year range
- Distinct Annual Leaders — count from primary query, sub-label shows top leader + years led
- Largest Leader Gap — year with the widest gap between #1 and #2 (flights), derived live from yearly ranks
- Sharpest Transition — year + carrier pair + swing magnitude from the sharpest_transition query

**Charts (side-by-side)**
- Bump chart — rank 1–5 per year per carrier, Y-axis inverted so rank 1 is at top; red dashed vertical line marks the sharpest transition year
- Flight share time series — percentage of all industry completed flights per year per carrier; same red dashed line annotation

**Leadership transitions table**
- Columns: Year, Prior Leader, New Leader, New Leader Flights, Prior Leader Flights, Swing vs Prior, Gap vs #2 (same year)
- Sharpest transition row highlighted in amber with a ★ marker

**Narrative takeaway**
- Derives era count, longest era, and transition frequency from live data
- Era chips strip (one per stability era) colored by carrier

**Query ledger**
- 5 entries: Primary (most_frequent_leader), Yearly Carrier Ranks, Leadership Swing, Sharpest Transition, Stability Pattern
- Each row expandable to show full SQL; status updates as queries complete

**Dynamic behavior**
- Auto-loads if JWE is already stored; guarded against concurrent runs
- Enrichment query failures degrade only their component; primary dashboard still renders
