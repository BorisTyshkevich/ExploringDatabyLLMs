# q002 Experiment Note

## Question

**ID:** q002  
**Slug:** q002_top_carrier_by_flights_leadership  
**Title:** Yearly carrier leadership by completed flights

Determine which `Reporting_Airline` flew the most completed flights in each calendar year, identify leadership transitions, and find where leadership changed most sharply based on year-over-year share swing.

## Why this question is useful

This question tests an LLM's ability to:
- Combine aggregation with window functions across multiple analytical layers
- Compute year-over-year metrics using lag functions over a filtered subset (leaders only)
- Rank carriers within partitions and extract conditional values (leader vs runner-up)
- Produce a unified result set supporting both tabular analysis and visual dashboards (bump charts, share time series)

The business value lies in understanding market dominance patterns—which carriers held leadership, how often it changed hands, and which transitions represented the largest market-share swings.

## Experiment setup

- **Day:** 2026-03-17
- **Dataset:** `default.ontime_v2` (39 years: 1987–2025)
- **Runners:**
  - claude/opus (run-001)
  - codex/gpt-5.4 (run-001)
  - gemini/gemini-3.1-pro-preview (run-001)
- **Prompt:** Requested a single result set with `RowType` column distinguishing `carrier_year` rows (top 5 per year) and optional `year_summary` rows, plus leader-transition metrics on rank-1 rows.

## Result summary

| Runner | Model | Run | Status | Rows | Duration | Read Rows | Memory |
|--------|-------|-----|--------|-----:|----------|-----------|--------|
| claude | opus | run-001 | ok | 195 | 594 ms | 921 M | 346 MiB |
| codex | gpt-5.4 | run-001 | ok | 195 | 586 ms | 921 M | 350 MiB |
| gemini | gemini-3.1-pro-preview | run-001 | ok | 234 | 1,458 ms | 1,612 M | 585 MiB |

**Row count mismatch:** Claude and Codex return 195 rows (5 carriers × 39 years = 195 `carrier_year` rows). Gemini returns 234 rows (195 `carrier_year` + 39 `year_summary` rows).

**Fastest run:** codex/gpt-5.4 at 586 ms.  
**Lowest resource usage:** claude/opus with 921 M read rows and 346 MiB memory.

## Full SQL artifacts

### claude / opus
- [run-001 query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/query.sql)

### codex / gpt-5.4
- [run-001 query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/query.sql)

### gemini / gemini-3.1-pro-preview
- [run-001 query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/query.sql)

## Real output differences

### Row-type composition
- **Claude and Codex:** 195 rows, all `RowType = 'carrier_year'` (top 5 carriers per year).
- **Gemini:** 234 rows—195 `carrier_year` rows plus 39 `year_summary` rows (one per year). The prompt labeled `year_summary` rows as optional, so Gemini's inclusion is valid but diverges from the other two.

### Leader-transition field population
- **Claude:** Populates `LeaderReportingAirline`, `RunnerUpReportingAirline`, `LeaderShareGapPctPts`, `PriorYearLeaderReportingAirline`, `LeaderChanged`, and `LeaderShareChangePctPts` on **all 5 carrier rows per year** (broadcasting leader metrics to every row via a simple JOIN).
- **Codex:** Populates these fields **only on rank-1 rows**; ranks 2–5 have `NULL` for leader-transition columns. This more strictly follows the prompt instruction: "Repeat the leader-transition fields on the rank-1 row for each year."
- **Gemini:** Same as Codex for `carrier_year` rows (leader fields only on rank 1); additionally emits `year_summary` rows with full leader metrics.

### Numeric precision
- **Claude:** Rounds `SharePct` and `LeaderShareGapPctPts` to 4 decimal places (e.g., `14.221`, `1.5935`).
- **Codex and Gemini:** Retain full double precision (e.g., `14.221048631689575`, `1.5935567403247788`).

### NULL vs empty string for first-year prior leader
- **Claude:** Uses empty string `""` for `PriorYearLeaderReportingAirline` in 1987.
- **Codex and Gemini:** Use `NULL`.

### Core flight counts and rankings
All three runs agree on the underlying data: 1987 leader is DL with 183,756 completed flights; 1988 leader is DL with 749,514; the top-5 carrier sets and their `CompletedFlights` values match exactly across runs. The leadership-change events (e.g., 1990 DL→US transition) are consistently identified.

## SQL comparison

| Aspect | Claude | Codex | Gemini |
|--------|--------|-------|--------|
| **CTE structure** | 7 CTEs | 4 CTEs | 8 CTEs + UNION ALL |
| **Leader extraction** | `argMax(col, RankInYear = 1)` | `maxIf(col, RankInYear = 1)` | `maxIf(col, RankInYear = 1)` |
| **Window function** | `lagInFrame(...) OVER (ORDER BY Year)` | `lagInFrame(...) OVER w` with explicit frame | `lag(...) OVER (ORDER BY Year)` |
| **Top-5 filter** | Separate `top5_carriers` CTE | `WHERE RankInYear <= 5` in final SELECT | `WHERE c.RankInYear <= 5` in FinalData |
| **year_summary rows** | Not emitted | Not emitted | Emitted via `UNION ALL` with SummaryData CTE |
| **Rounding** | `round(..., 4)` on SharePct and gap | No rounding | No rounding |
| **NULL handling** | Empty string for first-year prior leader | `toNullable()` casts, proper NULL | `CAST(NULL AS Nullable(...))` |

Gemini's query scans the base table twice (once for `YearlyCarrier`, once for `YearlyTotals` which re-aggregates), contributing to higher read volume (1.6B vs 921M rows). Claude and Codex both use a single aggregation pass for carrier-year stats and join to annual totals computed from the same CTE.

## Presentation artifacts

### claude / opus
- [run-001 report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq002_top_carrier_by_flights_leadership%2Fclaude%2Fopus%2Frun-001%2Freport.md)
- [run-001 visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q002_top_carrier_by_flights_leadership/claude/opus/run-001/visual.html)

### codex / gpt-5.4
- [run-001 report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq002_top_carrier_by_flights_leadership%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md)
- [run-001 visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/visual.html)

### gemini / gemini-3.1-pro-preview
- [run-001 report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq002_top_carrier_by_flights_leadership%2Fgemini%2Fgemini-3.1-pro-preview%2Frun-001%2Freport.md)
- [run-001 visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q002_top_carrier_by_flights_leadership/gemini/gemini-3.1-pro-preview/run-001/visual.html)

All three visuals use Chart.js for time-series and bump-chart rendering. Gemini's dashboard includes a dynamic data-fetch control panel with JWE token input; Claude and Codex embed static data or similar fetch mechanisms. Reports follow similar structures: metadata header, analysis summary, and a truncated result table showing 20 of N rows.

## Execution stats

### claude / opus
- **run-001:** 594 ms, 921,230,348 read rows, 2.46 GB read, 346.2 MiB memory, 34 threads

### codex / gpt-5.4
- **run-001:** 586 ms, 921,230,348 read rows, 2.46 GB read, 349.7 MiB memory, 34 threads

### gemini / gemini-3.1-pro-preview
- **run-001:** 1,458 ms, 1,612,153,109 read rows, 5.11 GB read, 584.9 MiB memory, 34 threads

Gemini's query reads 75% more rows and uses 69% more memory than Claude/Codex due to its `UNION ALL` approach and separate `YearlyTotals` aggregation that rescans the dataset.

## Takeaway

All three models produce correct answers to the core question—identifying annual flight-count leaders and leadership transitions. The differences are primarily structural:

1. **Row composition:** Gemini includes optional `year_summary` rows; Claude and Codex do not.
2. **Field broadcasting:** Claude broadcasts leader-transition fields to all top-5 rows per year; Codex and Gemini restrict them to rank-1 rows (closer to prompt intent).
3. **Precision:** Claude rounds to 4 decimals; others retain full precision.
4. **Efficiency:** Claude and Codex tie for fastest execution (~590 ms) and lowest resource usage. Gemini's more verbose query structure results in 2.5× longer runtime and higher I/O.

For dashboard consumption, Codex's approach is most precise (leader fields only on rank-1, no extra summary rows, full numeric precision). Claude's rounding and field broadcasting may simplify certain visualizations but deviates slightly from the prompt. Gemini's extra `year_summary` rows provide redundant annual aggregates that could be useful for some reporting contexts but cost significant additional compute.
