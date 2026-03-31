# Analysis Review
Verdict: WARN

## Summary

The core SQL is well-constructed across all three sections. The main query correctly deduplicates same-day legs, builds chronological Route strings, applies the max-hop filter, and returns all seven required fields with 10 rows confirmed in `results/main.json`. The q1 recurrence analysis is consistent with main and the per-route counts are accurate. The q2 SQL correctly derives the airport list from the same top-10 route selection logic. Two factual errors appear in q2's prose interpretation of the returned data: the state count is overstated and the California-as-endpoint count is wrong. These are not SQL or data errors; they are downstream prose claims that contradict the verified result.

## Findings

### main — PASS

- `queries/main.sql`: Grain is correct — `legs_deduped` groups by `(FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode)` to collapse duplicate or conflicting same-time rows before counting hops. The Route string is built by sorting origins by `dep_time` and appending `argMax(DestCode, dep_time)` as the final node, producing a properly chronological 9-token string for 8 hops.
- All seven required return fields are present in `results/main.json` (`aircraft_id`, `flight_number`, `carrier`, `most_recent_date`, `hop_count`, `recurrence_count`, `Route`).
- 10 rows returned, all with `hop_count = 8`, all carrier WN.
- `recurrence_count` is computed as `countDistinct(FlightDate)` within the Route group across all of `top_itin`, correctly counting distinct days any aircraft flew that exact route string at max-hop count.
- Empty `Tail_Number` rows are not filtered; no row in the result has a blank aircraft_id, consistent with the data having valid tail numbers for all qualifying records.
- Report table matches `results/main.json` exactly.

### q1 — PASS with minor wording note

- `queries/q1.sql` reuses the same itinerary pipeline and returns all 10 routes with their `recurrence_count` values, matching `results/main.json` exactly.
- Per-route recurrence list in the prose (1, 4, 2, 5, 40, 46, 7, 20, 12, 5) matches `results/q1.json` exactly.
- Tier counts add to 10: 1 one-off + 5 rare + 2 recurring + 2 highly recurring = 10. ✓
- Minor wording discrepancy: the SQL tier label for `>7 and ≤20` is `"recurring (8-20)"` but the prose describes that tier as "12–20 days." No actual routes fall in the 8–11 range so no route is miscategorized, but the prose description is narrower than the SQL boundary.

### q2 — WARN: two factual errors in prose

**State count overstated.** `results/q2.json` returns 45 rows with 24 unique `state_code` values (AR, AZ, CA, CO, FL, GA, IL, KS, LA, MD, MI, MO, MS, NC, NM, NV, NY, OH, SC, TN, TX, UT, VA, WA). IAD (Dulles) is stored with `state_code = "VA"` in the dimension table — no DC code appears. The report claims "45 unique airports across 27 US states (plus DC)." The correct count from the executed result is 24 states, not 27; DC does not appear as a separate state code.

**California terminus count wrong.** The report claims "California is the endpoint of 7 of the 10 routes." Reading the terminal airport (last token) of each of the 10 routes in `results/main.json`:
1. ISP-…-**SEA** (WA)
2. CLE-…-**DEN** (CO)
3. ELP-…-**SAN** (CA) ✓
4. MSY-…-**LAX** (CA) ✓
5. MSY-…-**SJC** (CA) ✓
6. LGA-…-**JAN** (MS)
7. SMF-…-**LAS** (NV)
8. HOU-…-**OAK** (CA) ✓
9. BWI-…-**OAK** (CA) ✓
10. BWI-…-**LAX** (CA) ✓

Six routes end in California, not seven. The correct figure is 6/10.

The remaining geographic observations (coast-to-coast sweep, Texas/South as through-corridor, Florida/Southeast as origin cluster, point-to-point spread) are qualitatively consistent with the verified airport list and are not disputed.

## Suggested Prompt Fixes

The two q2 prose errors are LLM-side arithmetic/counting mistakes on a concrete result set, not ambiguities in the prompt. The prompt already instructs the model to "base the geographic answer only on the airports appearing in the 10 routes returned by `main`." No prompt change would reliably prevent a miscounting error of this type.

None.
