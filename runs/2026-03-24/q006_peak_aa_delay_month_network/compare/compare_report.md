# q006 Experiment Note

## Question

**American Airlines peak network delay month and contributors** — Identify the single worst month for AA departure delays across the full ontime history, drill into which origins and routes drove that peak, and assess whether the disruption was broad or concentrated.

The question decomposes into four sub-questions:

1. Which month is the single worst AA month for departure delays?
2. Which origins contribute most to that peak month?
3. Which routes contribute most to that peak month?
4. Does the peak look broad across the network, or concentrated?

## Why this question is useful

This is a multi-step analytical question that tests whether a model can (a) rank hundreds of monthly aggregates to find a global maximum, (b) drill into the peak with appropriate contributor queries, and (c) synthesise a breadth-versus-concentration conclusion backed by correct share-of-total arithmetic. The last sub-question is especially tricky because it requires the denominator for concentration shares to be the full network, not a filtered subset.

## Experiment setup

- **Date**: 2026-03-24
- **Dataset**: Altinity `ontime.fact_ontime`
- **Runs**: 3 total across 2 provider/model combinations
  - **claude/sonnet** — 1 run (run-001), no automated review step
  - **codex/gpt-5.4** — 2 runs (run-001 failed review; run-002 passed review)
- Each run produced four SQL sub-queries (one per sub-question) and a Markdown report. Successful runs also produced an HTML visual dashboard.

## Result summary

All three runs agree on the headline answer: **July 2024** is the worst AA departure-delay month, with 86,083 completed flights, 36.33 min average departure delay, and 38.04% of flights departing 15+ minutes late. These figures are identical across every run's `result.json`.

The runs diverge on **how they rank origin and route contributors** and on **concentration methodology**:

| Aspect | claude/sonnet run-001 | codex/gpt-5.4 run-001 | codex/gpt-5.4 run-002 |
|---|---|---|---|
| Peak month | July 2024 (36.33 min) | July 2024 (36.33 min) | July 2024 (36.33 min) |
| Origin ranking metric | avg delay (CLT #1) | total delay minutes (DFW #1) | total delay minutes (DFW #1) |
| Top route | PHX-JAX by avg delay | DFW-LAX by total delay | DFW-LAX by total delay |
| Concentration denominator | flight count (no delay-minute shares) | filtered subset (HAVING >= 1000/200) | full network |
| Review verdict | n/a | FAIL | PASS |
| Status | ok | failed | ok |

## Full SQL artifacts

Each run produced four SQL files (one per sub-question) in its `queries/` directory.

### claude / sonnet

- **run-001**: `peak_month.sql`, `origin_contributors.sql`, `route_contributors.sql`, `concentration_pattern.sql` — uses `Reporting_Airline = 'AA'`; ranks origins and routes by `avg_dep_delay_min DESC`; origin threshold `HAVING flights >= 100`; route threshold `HAVING flights >= 30`; concentration query returns top-20 origins by flight volume with `pct_of_total_flights` via a window function but no delay-minute share columns.

### codex / gpt-5.4

- **run-001**: `q1.sql`–`q4.sql` — uses `Carrier = 'AA'`; ranks origins/routes by `total_departure_delay_minutes DESC`; origin threshold `HAVING flight_volume >= 1000`; route threshold `HAVING flight_volume >= 200`; `q4.sql` computes concentration shares but the denominator is the filtered subset (only origins with >= 1000 flights / routes with >= 200 flights), not the full network.
- **run-002**: `q1.sql`–`q4.sql` — same `Carrier = 'AA'` filter; same total-delay-minute ranking; origin threshold `flights >= 1000`; route threshold `flights >= 150`; critically, `q2.sql`/`q3.sql` use a `base` CTE to compute `network_delay_share_pct` against the unfiltered full-month network, and `q4.sql` similarly computes all shares against the full network denominator via `greatest(DepDelay, 0)`.

## Real output differences

All three runs return 4 meta-rows in `result.json` (one per sub-question) with 458 monthly rows for sub-question 1. The manifest row counts differ (508 / 494 / 656) because they reflect the sum of detailed rows across all sub-queries, which varies with different `LIMIT` and `HAVING` thresholds.

**Sub-question 1 (peak month)**: Identical across all runs — July 2024, 86,083 flights, 36.33 avg delay, 38.04% del15.

**Sub-question 2 (origin contributors)**: claude/sonnet ranks by average delay and surfaces CLT as #1 (52.0 min avg, 10,568 flights); codex runs rank by total delay minutes and surface DFW as #1 (593,021 total delay minutes, 14,962 flights). Both perspectives are valid but answer different notions of "contribute most."

**Sub-question 3 (route contributors)**: claude/sonnet ranks by average delay, surfacing low-volume routes like PHX-JAX (134.4 min avg, 31 flights); codex runs rank by total delay minutes, surfacing high-volume routes like DFW-LAX (18,462 total delay min, 443 flights). This is the largest qualitative divergence — the two approaches highlight completely different routes.

**Sub-question 4 (concentration)**: claude/sonnet describes concentration qualitatively using flight counts and average delays but does not compute delay-minute share percentages. codex/gpt-5.4 run-001 computes shares against a filtered denominator (top-2 origins = 50.8% of filtered total), which the review correctly flags as misleading. codex/gpt-5.4 run-002 fixes this: top-5 origins = 50.75% and top-10 routes = 4.88% of the full network delay minutes — a materially different and more defensible result.

## SQL comparison

**Airline filter column**: claude/sonnet uses `Reporting_Airline = 'AA'`; codex runs use `Carrier = 'AA'`. Both resolve to the same data in this dataset.

**Ranking philosophy**: The most consequential SQL difference is that claude/sonnet orders contributors by `avg(DepDelay) DESC`, while both codex runs order by `sum(DepDelayMinutes) DESC`. The average-delay approach highlights severity per flight; the total-delay approach highlights aggregate network impact. Neither is wrong, but they produce very different contributor lists.

**Volume thresholds**: claude/sonnet uses permissive thresholds (100 flights for origins, 30 for routes), yielding small-station results like EYW (119 flights). codex/gpt-5.4 uses stricter thresholds (1,000 for origins, 150–200 for routes), restricting output to high-volume airports and routes.

**Concentration denominator (the correctness-critical difference)**: codex run-001's `q4.sql` applies the `HAVING` filter before computing shares, so shares are relative to only "meaningful" entities — the review correctly identifies this as a denominator bug. codex run-002 fixes this with a `base` CTE that feeds an unfiltered `network` total into share calculations. claude/sonnet sidesteps the issue entirely by not computing delay-minute concentration shares at all.

**Additional columns**: codex run-002 enriches output with `origin_city`/`dest_city` names and `network_delay_share_pct`/`network_flight_share_pct` columns, making its output the most self-documenting of the three.

## Presentation artifacts

### claude / sonnet

- **run-001**
  - Report: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq006_peak_aa_delay_month_network%2Fclaude%2Fsonnet%2Frun-001%2Freport.md) — structured as four sub-question blocks with answer prose, row counts, and a sample first row. Emphasises CLT as the single dominant origin contributor.
  - Visual: [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-24/q006_peak_aa_delay_month_network/claude/sonnet/run-001/visual.html) (1,080 lines)

### codex / gpt-5.4

- **run-001**
  - Report: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq006_peak_aa_delay_month_network%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md) — same four-block structure. Identifies DFW and CLT as co-dominant. Uses filtered concentration shares (50.8% for top-2 origins).
  - Review: [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq006_peak_aa_delay_month_network%2Fcodex%2Fgpt-5.4%2Frun-001%2Freview.md) — verdict **FAIL**; the concentration query's denominator is a filtered subset, not the full network.
  - Visual: not generated (run failed review).
- **run-002**
  - Report: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq006_peak_aa_delay_month_network%2Fcodex%2Fgpt-5.4%2Frun-002%2Freport.md) — richer than run-001; includes city names, explicit network share percentages, and a note that display filters do not affect the denominator.
  - Review: [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq006_peak_aa_delay_month_network%2Fcodex%2Fgpt-5.4%2Frun-002%2Freview.md) — verdict **PASS**; all concentration shares computed against full network.
  - Visual: [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-24/q006_peak_aa_delay_month_network/codex/gpt-5.4/run-002/visual.html) (1,205 lines)

## Execution stats

| Provider/Model | Run | Status | SQL Gen (s) | Visual Gen (s) | Total Duration (s) | Manifest Rows | Query Time | Rows Read | Peak Memory |
|---|---|---|---:|---:|---:|---:|---|---|---|
| claude/sonnet | run-001 | ok | 83.2 | 294.0 | 381 | 508 | n/a | n/a | n/a |
| codex/gpt-5.4 | run-001 | failed | 187.7 | n/a | 298 | 494 | n/a | n/a | n/a |
| codex/gpt-5.4 | run-002 | ok | 188.2 | 382.1 | 674 | 656 | n/a | n/a | n/a |

Query-log metrics (query time, rows read, peak memory) were not captured for any run. claude/sonnet generated SQL in roughly half the wall-clock time of either codex run (83 s vs ~188 s). codex run-002's total duration (674 s) is nearly double claude/sonnet's (381 s), driven by both slower SQL generation and a longer visual-generation phase. codex run-001 finished in 298 s only because the failed review skipped visual generation entirely.

## Takeaway

All three runs correctly identify July 2024 as AA's worst departure-delay month and agree on the headline numbers. The meaningful differences are analytical rather than factual:

1. **Ranking metric matters**: claude/sonnet's average-delay ranking highlights severity at small stations (CLT, EYW); codex's total-delay-minutes ranking highlights aggregate hub impact (DFW, CLT). For a "network contributor" question, total delay minutes is arguably more aligned with the prompt's intent, but neither is incorrect.

2. **Concentration denominator is the hardest part**: codex/gpt-5.4 run-001 failed its own automated review because its concentration shares used a filtered denominator — a subtle but substantive bug. run-002 fixed this by computing all shares against the full network. claude/sonnet avoided the trap entirely by not computing delay-minute shares, but at the cost of a less quantitative concentration answer.

3. **The review loop adds value**: The codex pipeline's self-review caught the denominator bug in run-001 and produced a materially improved run-002 with correct full-network shares, city-name enrichment, and explicit denominator documentation. This is a concrete example of automated review improving analytical correctness.
