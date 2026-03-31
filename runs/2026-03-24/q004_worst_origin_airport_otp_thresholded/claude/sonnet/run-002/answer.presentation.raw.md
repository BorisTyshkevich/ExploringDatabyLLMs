`visual.html` written. Here's what the dashboard delivers:

**KPI Strip (4 cards)**
- Worst Airport — code + full name
- Worst Departure OTP — rate and % delayed
- Highest Avg Delay When Late — across all ranked airports (not just worst OTP)
- Qualifying Airports — count + median OTP + spread from `spread_to_middle`

**Charts**
- **Horizontal bar chart** (top 25 worst, worst at top): worst 5 in red, ranks 6–25 in sky blue, teal dashed median line annotated with value; worst airport labelled "WORST" inline
- **Scatter plot** (all qualifying airports, log x-axis): three-layer dataset — others (sky), worst 2–5 (red), worst airport (red star); teal median reference line; worst airport code annotated; log scale handles the volume skew from 100K to 10M+ departures

**Supporting detail**
- Narrative takeaway derived from fetched data (bottom-10 codes, avg volume, spread, count below 85%)
- Full ranked table — sortable by any column, filterable, CSV export, worst-5 rows highlighted in red with star on rank 1
- Query ledger with expandable SQL for both `worst_airport` (primary) and `spread_to_middle` (supporting)

**Dynamic mode**
- JWE stored in `localStorage['OnTimeAnalystDashboard::auth::jwe']`; auto-runs on page load if token present
- Both queries run in parallel via `Promise.allSettled`; spread failure degrades only KPI/median lines, not the main dashboard
