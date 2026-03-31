`visual.html` written. Here's what the dashboard delivers:

**KPI strip (4 cards)**
- Worst winter pair (carrier / origin) — red accent
- Lowest winter OTP % — red, with flight count sub-label
- Average departure delay for the worst pair — amber
- Total qualifying pairs count — navy

**Charts (side-by-side)**
- Horizontal ranked bar chart of all 20 worst pairs, worst at top, color-coded red/amber/sky by severity
- Stacked horizontal bar chart of delay-cause shares (weather / carrier / NAS / late aircraft / security) for the top 10 pairs — labels derived live from the primary fetch result, not hardcoded

**Full ranked table** with inline OTP progress bars, color-coded values, rank badges (red for #1, amber for top-3), and CSV export

**Narrative block** (dark navy panel) — derives carrier/airport concentration stats from live data and surfaces the key finding: weakness is carrier-driven, not airport-driven

**Query ledger** — two entries (`worst_pair` primary + `cause_mix` enrichment), each with expand/collapse SQL reveal, status badges, and row counts

**Dynamic behavior**
- Auto-runs if a stored JWE token is found in `localStorage`
- Cause-mix SQL is built dynamically in the browser from the top-10 pairs returned by the primary query
- Enrichment failure degrades only the cause chart; ranking, table, and KPIs remain visible
- Token stored/forgotten via footer controls; never echoed in UI
