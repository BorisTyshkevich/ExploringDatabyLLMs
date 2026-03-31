`review.md` written. Key findings:

**Verdict: WARN**

- **main** — fully correct. All 7 required fields present, route strings built chronologically via `CRSDepTime` sort, recurrence counted as distinct dates across all history, correct 10-row result ordered by most-recent date.

- **q1** — one factual error: the "Low recurrence (2–7 days)" tier is reported as 4 routes, but the data shows 5 routes in that band (recurrences: 2, 5, 5, 5, 7). Tier totals sum to 9, not 10 — violating the explicit prompt requirement. The 10 route/recurrence values in the table are individually correct. Additionally, `queries/q1.sql` returns 719 rows (all max-hop routes) rather than being scoped to the 10 from main.

- **q2** — correct. Hardcoded airport list matches the 10 main routes, longitude-based region classification produces West (17) / Central (14) / East (14), geographic narrative is well-supported.
