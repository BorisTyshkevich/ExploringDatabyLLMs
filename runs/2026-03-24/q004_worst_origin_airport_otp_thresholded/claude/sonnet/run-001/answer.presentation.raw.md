`visual.html` written. Here's what the dashboard includes:

**KPI Strip**
- Worst airport code + name (CEC / Jack McNamara Field)
- Worst OTP % with "~1 in N departures delayed" annotation
- Highest avg delay among ranked airports (with airport code)
- Total qualifying airport count

**Charts**
- **Ranked horizontal bar chart** — worst 25 airports, worst→least-worst top-to-bottom; rank #1 in darkest red, ranks 2–5 in red, ranks 6–25 in sky blue; "worst" label drawn on the #1 bar via a custom Chart.js plugin
- **Scatter plot** — completed departures (log scale) vs. OTP %; four datasets colored by volume tier (red=major hub, amber=large, sky=mid-size, teal=regional); worst airport annotated with label + emphasis ring; log scale keeps CEC's 11k alongside ORD's 11M readable

**Narrative** — derived live from fetched data: identifies regional airports dominating the bottom, names the major hubs in the bottom 30, quantifies the spread from worst to median using the `spread_to_middle` result

**Detail table** — all 30 rows; top-5 rows highlighted (red/amber row backgrounds); volume tier badges from `bottom_group_mix`; CSV export

**Query ledger** — all three queries (primary + 2 enrichments) with expand/collapse SQL, live status, and row counts

**Footer** — JWE token (password field), SQL textarea, Fetch / Forget buttons, status line; auto-loads on page open if token is stored
