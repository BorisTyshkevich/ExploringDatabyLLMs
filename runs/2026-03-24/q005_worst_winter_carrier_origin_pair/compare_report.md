```markdown
# q005 Experiment Note

## Question

**Worst winter carrier-origin pairs by departure performance** -- identify which (carrier, origin airport) combinations have the worst winter departure on-time performance, determine whether delays are weather- or operationally-driven, and assess whether the weak set is concentrated in a few carriers or airports.

Prompt: [`report_prompt.md`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q005_worst_winter_carrier_origin_pair/report_prompt.md) | Visual: [`visual_prompt.md`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q005_worst_winter_carrier_origin_pair/visual_prompt.md)

## Why this question is useful

Winter weather is the default scapegoat for airline delays, but disaggregating by carrier-origin pair exposes whether poor on-time performance is genuinely weather-linked or stems from chronic operational issues at specific hubs. This question tests an LLM's ability to (a) define winter correctly, (b) apply a credible volume filter, (c) decompose delay causes from the five BTS cause columns, and (d) handle the null-cause-data caveat for older carriers like Piedmont (PI).

## Experiment setup

| Dimension | Value |
|---|---|
| Day | 2026-03-24 |
| Dataset | `ontime.fact_ontime` (BTS On-Time) |
| Runs | 3 total: claude/sonnet run-001 & run-002, codex/gpt-5.4 run-001 |
| Review | run-001 had no automated review; run-002 and codex run-001 both received review and scored **PASS** |
| Winter definition | Months 12, 1, 2 (consistent across all runs) |
| Volume filter | >= 1,000 winter flights (all runs) |

## Result summary

All three runs agree on the headline answer: **DH (Independence Air) at ORD (Chicago O'Hare) is the worst qualifying winter carrier-origin pair**, and **operational causes dominate over weather** for every measured pair. The runs differ in secondary details -- row counts (51 / 45 / 52), the set size used for concentration analysis (25 vs 30 vs 10 pairs), and minor numeric discrepancies in DH/ORD flight counts and average delay due to a diverted-flight exclusion difference.

Both reviews (claude/sonnet run-002 and codex/gpt-5.4 run-001) returned **PASS** with no substantive correctness defects.

## Full SQL artifacts

### claude / sonnet

- **run-001** -- three named query files:
  - `worst_pair.sql`, `cause_mix.sql`, `concentration_pattern.sql`
  - [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-24/q005_worst_winter_carrier_origin_pair/claude/sonnet/run-001/result.json) (51 manifest rows, 3 sub-question rows in result.json)
- **run-002** -- three query files (`q1.sql`, `q2.sql`, `q3.sql`):
  - [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-24/q005_worst_winter_carrier_origin_pair/claude/sonnet/run-002/result.json) (45 manifest rows, 3 sub-question rows in result.json)

### codex / gpt-5.4

- **run-001** -- three query files (`q1.sql`, `q2.sql`, `q3.sql`):
  - [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-24/q005_worst_winter_carrier_origin_pair/codex/gpt-5.4/run-001/result.json) (52 manifest rows, 3 sub-question rows in result.json)

## Real output differences

All runs return 3 sub-question result rows. The manifest row counts (51, 45, 52) reflect the sum of inner row counts across the three sub-questions. Verified differences:

### Sub-question 1 -- worst pair ranking

| Field | claude/sonnet run-001 | claude/sonnet run-002 | codex/gpt-5.4 run-001 |
|---|---|---|---|
| Top pair | DH / ORD | DH / ORD | DH / ORD |
| winter_flights | 19,986 | 19,986 | 19,929 |
| OTP % | 56.58 | 56.58 | 56.53 |
| avg_dep_delay_min | 28.35 | 28.35 | 27.06 |
| Rows returned | 25 | 20 | 25 |

The codex run reports 57 fewer flights (19,929 vs 19,986) and a lower average delay (27.06 vs 28.35 min). This is explained by codex additionally filtering `Diverted = 0`, while both claude runs filter only `Cancelled = 0`. The diverted-flight exclusion slightly changes the denominator and the average. The OTP ranking itself is unaffected -- DH/ORD is worst across all three.

### Sub-question 2 -- cause composition

All runs agree that operational causes dominate. The weather/operational share for DH/ORD is 15.6% / 84.4% in all three runs. The underlying cause-minute totals are identical (weather 55,176; carrier 123,407; NAS 41,492; security 116; late aircraft 133,239) despite the flight-count discrepancy, because cause-minute columns are null for diverted flights. The claude runs analyze 25 or 10 pairs; codex analyzes 10 pairs.

### Sub-question 3 -- concentration

| Run | Scope | Distinct carriers | Distinct airports |
|---|---|---|---|
| claude/sonnet run-001 | worst 25 pairs | 14 | 22 |
| claude/sonnet run-002 | worst 30 pairs | 15 | -- (carrier-grouped) |
| codex/gpt-5.4 run-001 | worst 10 pairs | -- (union of carrier/airport rows) | -- |

Run-001 aggregates the 25 worst into a single summary row. Run-002 expands to the worst 30 pairs grouped by carrier (B6 leads with 6 pairs). Codex limits to the top 10 and uses a UNION ALL to separate carrier and airport concentration into distinct rows (PI leads with 3 of 10 pairs). The qualitative conclusion is shared: moderate carrier concentration, weaker airport concentration.

## SQL comparison

All runs define winter as `Month IN (12, 1, 2)`, filter `Cancelled = 0`, group by `(carrier, origin)`, and apply `HAVING >= 1000` flights. Key structural differences:

| Aspect | claude/sonnet run-001 | claude/sonnet run-002 | codex/gpt-5.4 run-001 |
|---|---|---|---|
| Diverted exclusion | No | No | Yes (`Diverted = 0`) |
| Q1 LIMIT | 25 | 20 | 25 |
| Q1 delay column | `DepDelayMinutes` | `ifNull(DepDelayMinutes,0)` | `ifNull(DepDelay,0)` -- uses `DepDelay` not `DepDelayMinutes` |
| Q1 extra columns | weather_pct, operational_pct, rank | -- (OTP and delay only) | DisplayAirportName, cause-reporting counts |
| Q2 approach | CTE on worst 25 + share calc | Hardcoded 10-pair IN-list filter | CTE on worst 10 + share calc |
| Q2 null handling | `assumeNotNull()` | `ifNull(,0)` | `ifNull(toUInt64(),0)` |
| Q3 scope | worst 25 -> single summary row | worst 30 -> carrier-grouped rows | worst 10 -> UNION ALL carrier + airport rows |
| Airport dimension join | No | No | Yes (`dim_airports`) |

Codex is the only run that joins `dim_airports` to surface `DisplayAirportName` and uses `DepDelay` instead of `DepDelayMinutes`. Claude run-002 hardcodes the top-10 pair list into Q2 rather than deriving it from a CTE, which is less maintainable but functionally equivalent given the stable ranking.

## Presentation artifacts

### claude / sonnet

- **run-001**
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq005_worst_winter_carrier_origin_pair%2Fclaude%2Fsonnet%2Frun-001%2Freport.md) -- tables show only the first row per sub-question; full data in result.json
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-24/q005_worst_winter_carrier_origin_pair/claude/sonnet/run-001/visual.html)
- **run-002**
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq005_worst_winter_carrier_origin_pair%2Fclaude%2Fsonnet%2Frun-002%2Freport.md) -- same single-row table format; review noted this truncation as a minor gap
  - [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq005_worst_winter_carrier_origin_pair%2Fclaude%2Fsonnet%2Frun-002%2Freview.md) -- PASS; flagged PI null-cause data under-reporting and single-row table truncation
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-24/q005_worst_winter_carrier_origin_pair/claude/sonnet/run-002/visual.html)

### codex / gpt-5.4

- **run-001**
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq005_worst_winter_carrier_origin_pair%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md) -- explicitly flags incomplete cause reporting for PI pairs; uses code-formatted carrier names
  - [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-24%2Fq005_worst_winter_carrier_origin_pair%2Fcodex%2Fgpt-5.4%2Frun-001%2Freview.md) -- PASS; no substantive defects; no suggested prompt fixes
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-24/q005_worst_winter_carrier_origin_pair/codex/gpt-5.4/run-001/visual.html)

## Execution stats

| Provider / Model | Run | Status | Query time | Rows read | Bytes read | Peak memory | SQL gen (s) | Visual gen (s) | Total (s) |
|---|---|---|---|---|---|---|---:|---:|---:|
| claude / sonnet | run-001 | ok | n/a | n/a | n/a | n/a | 121.3 | 322.2 | 447 |
| claude / sonnet | run-002 | ok | n/a | n/a | n/a | n/a | 156.2 | 357.5 | 597 |
| codex / gpt-5.4 | run-001 | ok | n/a | n/a | n/a | n/a | 162.7 | 418.1 | 701 |

Query-log metrics (query time, rows read, bytes read, peak memory) were not captured for any run. Among the available timing data, claude/sonnet run-001 was the fastest end-to-end at 447 s (7.5 min), while codex/gpt-5.4 run-001 was the slowest at 701 s (11.7 min) -- a 57% spread driven mainly by longer visual generation (418 s vs 322 s).

## Takeaway

All three runs converge on the same headline: DH/ORD is the worst winter pair by OTP, and operational causes overwhelm weather for every measurable pair. The most meaningful divergence is the codex run's additional `Diverted = 0` filter, which slightly alters flight counts and average delay but does not change the ranking or conclusions. Codex also produced the most defensible handling of the cause-data gap -- explicitly flagging PI's missing cause minutes both in SQL (via `flights_with_any_reported_cause` / `flights_without_reported_cause` columns) and in report prose. The claude run-002 review correctly identified the same gap but as a minor issue in the report rather than a SQL-level fix. For this question, model choice had no impact on correctness; the differences are in analytical polish and presentation scope.
```
