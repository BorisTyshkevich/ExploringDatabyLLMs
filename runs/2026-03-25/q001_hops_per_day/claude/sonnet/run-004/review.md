# Analysis Review
Verdict: WARN

## Summary

The run correctly identifies the top 10 unique 8-hop routes by most-recent departure date, returns all seven required fields, and provides well-supported geographic analysis. One factual error exists in the q1 tier breakdown: the "Low recurrence (2–7 days)" bucket is stated as 4 routes when the data shows 5, causing the tier totals to sum to 9 instead of the required 10. The core data tables and SQL are correct throughout.

## Findings

### main — PASS

- `queries/main.sql` groups by `(Tail_Number, FlightNum, Carrier, FlightDate)` and sorts legs by `CRSDepTime` before concatenating into the Route string. Grain is correct.
- All 7 required fields are present in `results/main.json`: `Tail_Number`, `FlightNum`, `Carrier`, `FlightDate`, `hop_count`, `recurrence_count`, `Route`.
- 10 rows returned, all WN, all with `hop_count = 8` (9-airport itineraries).
- `recurrence_count` is computed via `countDistinct(FlightDate)` across all of `daily_routes` for matching route strings — correctly counts distinct calendar days across all history for that exact route, regardless of aircraft.
- Results are ordered by `most_recent_date DESC`; the most recent row (`ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA`, 2024-12-01) is verified in `results/main.json`.

### q1 — WARN: tier count error, totals sum to 9 not 10

The 10 routes and their recurrence values are correctly extracted from the main result and listed in `report.md` and `answer.raw.json`:

| Route | Recurrence |
|---|---|
| ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA | 1 |
| ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN | 2 |
| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 |
| MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX | 5 |
| BWI-MCO-MEM-MDW-IAD-ATL-MSY-DAL-LAX | 5 |
| SMF-SAN-PHX-COS-DEN-PSP-OAK-RNO-LAS | 7 |
| BWI-FLL-MSY-DAL-MAF-DEN-LAS-BUR-OAK | 12 |
| HOU-MSY-BNA-MYR-CMH-DAL-ABQ-LAS-OAK | 20 |
| MSY-TPA-BWI-ORD-DEN-SLC-LAS-BUR-SJC | 40 |
| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 |

However, the tier breakdown in `report.md` states:

> Low recurrence (2–7 days): **4 routes**

The routes with recurrence in 2–7 days from the data are: ELP (2), CLE (5), MSY-ATL (5), BWI-MCO (5), SMF-SAN (7) — that is **5 routes**, not 4. The tier total is therefore 1 + 4 + 2 + 2 = **9**, not 10. The report claims "totals add to 10" but they do not. This violates the explicit prompt requirement.

Additionally, `queries/q1.sql` returns all 719 max-hop routes ordered by recurrence (confirmed in `results/q1.json`, row_count=719), not just the 10 from main. The answer derives its 10-route table by cross-referencing main rather than from a q1 query scoped to those 10 routes. This is a process gap but does not affect the correctness of the route-level recurrence values shown.

### q2 — PASS

- `queries/q2.sql` hardcodes the 45 distinct airport codes appearing in the 10 main routes and joins `dim_airports` for longitude-based region classification.
- `results/q2.json` returns 3 rows: West (17 airports), Central (14), East (14), summing to 45. The airport lists in the result match the codes extracted from `results/main.json`.
- The geographic narrative in `report.md` is fully supported: all airports are CONUS, all routes span at least two geographic bands, and the coast-to-coast characterization is accurate.
- The longitude thresholds used (< −100 = West, < −85 = Central, else East) correctly classify all listed airports.

## Suggested Prompt Fixes

**Fix 1 — q1 tier counting:** The current prompt only says "make sure any category totals add up to 10," which was claimed but not verified before writing. Adding an explicit instruction to derive tier boundaries mechanically from the listed recurrences would reduce this class of error:

> After listing the 10 recurrence values, choose tier boundaries that include every route in exactly one tier. Show the boundary definition and recount each tier from the listed values before stating totals.

**Fix 2 — q1 query scope:** The q1 query returns all 719 max-hop routes rather than being scoped to the 10 from main. This makes it harder to verify that the q1 answer is derived from the correct 10 routes. Adding a constraint to the prompt would close this gap:

> For q1, write the query so that it operates only on the exact 10 Route strings returned by the main query (not the full population of maximum-hop routes).
