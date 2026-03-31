# q001 Experiment Note

**Day:** 2026-03-20 | **Question:** q001 — Highest daily hops for one aircraft on one flight number

## Question

Find the aircraft+flight-number combinations that logged the most legs ("hops") in a single calendar day across the entire `ontime.ontime` dataset. Return the top 10 by hop count (descending), with the full route chain and departure times.

Prompt: [`report_prompt.md`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/report_prompt.md) | Visual prompt: [`visual_prompt.md`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/visual_prompt.md) | Compare contract: [`compare.yaml`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/compare.yaml)

## Why this question is useful

It exercises several SQL skills at once: grouping by a composite key (tail number + flight number + date), counting legs, then reconstructing ordered route chains with `groupArray` / `arraySort`. It also tests whether a model adds sensible filters (cancelled flights, null departure times, diverted flights) and whether it can keep read volume low through single-pass aggregation versus multi-pass JOINs.

## Experiment setup

Five runs across three model variants, all querying the same ClickHouse `ontime.ontime` table via the qforge harness:

| Provider | Model | Run | Notes |
|----------|-------|-----|-------|
| claude | opus | run-001 | Full pipeline (SQL + report + visual) |
| claude | opus | run-002 | SQL + report only (visual skipped) |
| claude | sonnet | run-001 | Full pipeline |
| codex | gpt-5.4 | run-001 | Full pipeline, visual render failed |
| codex | gpt-5.4 | run-002 | Full pipeline, visual render failed |

## Result summary

**All five runs agree on the answer.** Every run returns 10 rows, all with Hops = 8, all Southwest Airlines (WN). The same 10 (tail number, flight number, date) tuples appear in every result set in the same order. The top entry across all runs is N957WN / WN 366 on 2024-12-01 flying ISP through SEA in 8 legs.

The differences between runs are limited to cosmetic formatting of routes and departure times, extra analytical columns in the codex runs, and column naming conventions. No run disagrees on any factual value.

## Full SQL artifacts

### claude / opus

- **run-001:** [`query.sql`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/claude/opus/run-001/query.sql) | [`result.json`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/claude/opus/run-001/result.json)
- **run-002:** [`query.sql`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/claude/opus/run-002/query.sql) | [`result.json`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/claude/opus/run-002/result.json)

### claude / sonnet

- **run-001:** [`query.sql`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/claude/sonnet/run-001/query.sql) | [`result.json`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/claude/sonnet/run-001/result.json)

### codex / gpt-5.4

- **run-001:** [`query.sql`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-001/query.sql) | [`result.json`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-001/result.json)
- **run-002:** [`query.sql`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-002/query.sql) | [`result.json`](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-002/result.json)

## Real output differences

All five runs return **identical factual data** for the core fields: tail number, flight number, carrier (WN), flight date, and hop count (8 for every row). The same 10 itineraries appear in the same order. Verified differences are strictly cosmetic or structural:

**Column naming.** Claude runs use raw schema names (`Tail_Number`, `Flight_Number_Reporting_Airline`); opus/run-002 and sonnet alias the carrier column to `Carrier`. Codex runs rename all columns to human-readable labels (`Aircraft ID`, `Flight Number`, `Carrier`, `Date`).

**Route formatting.** Dash-separated in opus runs (`ISP-BWI-MYR-...`), arrow-separated in sonnet (`ISP → BWI → MYR → ...`), and leg-pair formatted in codex/run-002 (`05:43 ISP->BWI | 08:10 BWI->MYR | ...`). Codex/run-001 uses `->` between airports with departure times in a separate column.

**Departure time formatting.** Opus/run-001 uses `H:MM` (`5:43, 8:10, ...`), opus/run-002 uses zero-padded HHMM (`0543, 0810, ...`), sonnet annotates each time with its origin airport (`0543 (ISP), 0810 (BWI), ...`), codex/run-001 uses `HH:MM ORIGIN` pipe-delimited (`05:43 ISP | 08:10 BWI | ...`), and codex/run-002 folds departure times into the Route column itself.

**Extra columns.** Both codex runs add `Maximum Hops Observed` (8) and `Maximum-Hop Itinerary Count` (9,859). Codex/run-001 further adds `Same-Hops Route Count` and `Route Frequency In Top 10` via window functions.

## SQL comparison

All five queries share the same high-level pattern: filter cancelled flights, group by (tail number, flight number, carrier, date), count legs, reconstruct ordered route chains, order by hops descending, limit 10. They diverge on structure, filtering, and route-building strategy.

**claude/opus/run-001 (3 CTEs, two-pass).** A `legs` CTE filters rows, a `counts` CTE gets the top-10 groups by hop count, then a `top_legs` CTE JOINs back to `legs` to recover per-leg origin/dest for route assembly. The JOIN causes a double scan, reading 386M rows — the highest among the three Claude runs. Route is built with `groupArray(Origin)` plus the last destination appended via `arrayElement`.

**claude/opus/run-002 (2 CTEs, single-pass).** A `legs` CTE feeds a `counted` CTE that does grouping, counting, and route assembly in one pass using `arraySort(x -> x.1, groupArray((DepTime, Origin)))` with a `HAVING Hops >= 2` pre-filter. This halves the read volume to 193M rows. Route uses `argMax(Dest, DepTime)` for the final destination.

**claude/sonnet/run-001 (1 CTE, single-pass).** Groups directly in a `leg_data` CTE using `groupArray` of tuples sorted by `coalesce(DepTime, CRSDepTime, 0)`. The coalesce fallback is a defensive choice not seen in other runs. Adds `Flight_Number_Reporting_Airline != ''` as an extra filter. Also reads 193M rows. Uses `arrayElement(sorted_legs, -1).3` for the last destination.

**codex/gpt-5.4/run-001 (5 CTEs, scalar subqueries + window).** Converts `DepTime` to a full `DateTime`, adds a `Diverted = 0` filter not used by Claude runs, and builds route chains with 4-element tuples. Adds three scalar/window-function columns (`Max Hops Overall`, `Max Hop Itinerary Count`, `Same-Hops Route Count`, `Route Frequency In Top 10`). The extra computed columns and lack of early `HAVING` filter push read volume to 579M rows.

**codex/gpt-5.4/run-002 (4 CTEs, JOIN to airports_latest).** The most complex query. Joins `ontime.airports_latest` to convert departure times to UTC for sorting, using parsed UTC-offset strings. Uses `TailNum`/`FlightNum` column aliases (legacy naming). The airport JOIN and UTC conversion push read volume to 772M rows — the highest of all runs.

## Presentation artifacts

### claude / opus

- **run-001:** [`report.md`](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-20%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-001%2Freport.md) | [`visual.html`](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-20/q001_hops_per_day/claude/opus/run-001/visual.html) — Structured report with Overview, Key Finding, Most Recent Maximum-Hop Itinerary, and Route Repetition sections. Highlights the Southwest point-to-point model. Visual rendered successfully.
- **run-002:** [`report.md`](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-20%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-002%2Freport.md) — Similar structure and narrative to run-001 but more concise. Visual generation was skipped for this run.

### claude / sonnet

- **run-001:** [`report.md`](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-20%2Fq001_hops_per_day%2Fclaude%2Fsonnet%2Frun-001%2Freport.md) | [`visual.html`](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-20/q001_hops_per_day/claude/sonnet/run-001/visual.html) — Well-structured report with Maximum Hop Count, Operating Pattern, and Route Repetition sections. Uses blockquote for the top route. Departure times include origin airport codes inline. Visual rendered successfully.

### codex / gpt-5.4

- **run-001:** [`report.md`](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-20%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md) | [`visual.html`](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-001/visual.html) — Report is template-like, with placeholder language ("should be interpreted as", "notable route repetition or clustering is") rather than filled-in analysis. Visual render failed (status: partial).
- **run-002:** [`report.md`](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-20%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-002%2Freport.md) | [`visual.html`](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-002/visual.html) — Same template-like quality as run-001, with unresolved placeholder phrases. Visual render failed (status: partial).

## Execution stats

| Provider/Model | Run | Status | Query time | Rows read | Bytes read | Peak memory | SQL gen time | Visual gen time | Total duration |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|
| claude/opus | run-001 | ok | 4.37 s | 386,123,882 | 4.3 GiB | 24.9 GiB | 87.05 s | 313.33 s | 240 s |
| claude/opus | run-002 | ok | 14.45 s | 193,061,941 | 2.5 GiB | 48.9 GiB | n/a | n/a | 14 s |
| claude/sonnet | run-001 | ok | 7.36 s | 193,061,941 | 3.0 GiB | 41.2 GiB | 104.05 s | 352.80 s | 331 s |
| codex/gpt-5.4 | run-001 | partial | 19.50 s | 579,185,823 | 5.4 GiB | 50.1 GiB | 297.21 s | 645.69 s | 982 s |
| codex/gpt-5.4 | run-002 | partial | 32.10 s | 772,328,680 | 12.8 GiB | 43.9 GiB | 425.21 s | 636.04 s | 1,099 s |

Query execution time spans a 7.3x range: claude/opus/run-001 at 4.37 s versus codex/gpt-5.4/run-002 at 32.10 s. The read-volume spread is 4x (193M to 772M rows), driven primarily by whether the query uses single-pass aggregation with early `HAVING` (opus/run-002, sonnet) versus multi-CTE JOINs and extra table lookups (codex/run-002). SQL generation time shows the widest gap: Claude models produced SQL in 87–104 s while codex took 297–425 s. Both codex runs failed at visual rendering, ending in `partial` status.

## Takeaway

All five runs converge on the same answer: the maximum daily hops for one aircraft on one flight number is **8**, seen exclusively on Southwest Airlines multi-stop through-flights. The factual agreement is complete — every run identifies the same 10 itineraries in the same order.

The key differentiators are efficiency and presentation quality. The Claude runs produced leaner SQL (193M–386M rows read, 4–14 s query time) and fully rendered reports with genuine analytical narrative. The codex/gpt-5.4 runs generated more complex queries that read 2–4x more data (579M–772M rows), took 3–5x longer to generate SQL, and produced template-like reports with unresolved placeholder language. Both codex runs also failed at visual rendering.

Within the Claude family, opus/run-002's single-pass `arraySort` + `HAVING` approach achieved the lowest read volume (tied with sonnet at 193M rows), while opus/run-001's two-pass JOIN strategy was the fastest wall-clock query (4.37 s) despite reading more data — likely benefiting from ClickHouse parallel execution on the JOIN pattern.
