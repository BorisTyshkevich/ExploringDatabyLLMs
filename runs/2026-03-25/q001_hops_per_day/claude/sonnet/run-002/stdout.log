Both proof queries execute successfully. Both artifacts are written and verified.

**Summary:**

- **`main.sql`**: Groups flights by `(Tail_Number, FlightNum, Carrier, FlightDate)`, builds the chronologically-ordered route string using `arraySort` on `CRSDepTime`, then deduplicates by unique `Route` string — keeping the most recent aircraft/date for each unique route — ordered by max hop count descending, top 10.

- **`answer.raw.json`** answers both dashboard questions:
  1. **Recurring patterns**: Top itineraries are definitively recurring scheduled rotations — the highest-frequency 8-hop routes each appear 54–114 times over multi-month windows, all on Southwest Airlines (WN).
  2. **Geographic pattern**: All top routes are full-width continental US daisy-chains on Southwest, typically running east-to-west through 6–8 states, touching Florida/Gulf Coast → Midwest → Texas/Southwest → California in a single day's rotation.
