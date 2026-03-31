`visual.html` written. Here's what the dashboard delivers:

**KPI strip (4 cards)**
- Worst Winter Pair (carrier/origin + flight count)
- Worst OTP % (on-time departure rate)
- Avg Delay — Worst Pair (minutes)
- Total Qualifying Pairs

**Charts**
- Horizontal ranked bar chart — all 25 pairs sorted worst-first, color-coded: red (rank 1), amber (ranks 2–5), blue (6–12), teal (rest). Tooltip shows OTP, avg delay, and flight count.
- Stacked bar chart — delay-cause shares (Late Aircraft / Carrier / NAS / Weather / Security) for the top 10 pairs. Labels derived from live ranking data, not hardcoded.

**Narrative** — auto-generated from fetched data: distinct carrier/airport counts, which carriers repeat, average operational vs. weather split.

**Full table** — all 25 rows with OTP badge coloring (red < 65%, amber < 72%, teal otherwise) + CSV export.

**Query ledger** — `worst_pair` (Primary) and `cause_mix` (Enrichment) both registered with expandable SQL, live status, and row counts.

**Dynamic behavior** — auto-loads on page open if a stored JWE token is found; cause-mix chart degrades gracefully without blocking the primary ranking view.
