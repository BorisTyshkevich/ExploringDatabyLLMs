Both files written. Here's a summary of what was produced:

**`main.sql`** — groups every non-cancelled flight by `(Tail_Number, FlightNum, Carrier, FlightDate)`, builds the chronologically-ordered `Route` string via `arraySort + groupArray`, then deduplicates by unique `Route` string, retaining the most recent occurrence and occurrence count. Results are ordered by `hop_count DESC, most_recent_date DESC LIMIT 10`.

**Key findings from the query results:**

| Q | Finding |
|---|---------|
| Highest-hop example | WN 366 / N957WN, 2024-12-01, `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` — 8 hops, 9 airports |
| Most recent | Same flight (WN 366, 2024-12-01) is both the most recent and tied for highest hop count |
| Recurring vs one-offs | Mixed — WN 2884 recurred 46×, WN 1956 recurred 40× (clear scheduled patterns); others appeared 1–5× (one-offs). All 10 are Southwest, consistent with their high-turn short-haul model |
