# q003 Experiment Note

## Question

`q003` asks the model to use `default.ontime_v2` to find Delta departures out of ATL with the worst sustained departure delays at the `(Dest, DepTimeBlk)` level.

The prompt requires:

- completed flights only: `Cancelled = 0`
- Delta ATL departures only: `IATA_CODE_Reporting_Airline = 'DL'` and `Origin = 'ATL'`
- monthly qualification at `(MonthStart, Dest, DepTimeBlk)` with at least `40` flights
- hotspot qualification across qualifying monthly cells only, with at least `1,000` flights
- hotspot ranking by:
  - `AvgDepDelayMinutes` descending
  - `P90DepDelayMinutes` descending
  - `DepDel15Pct` descending
  - `CompletedFlights` descending
- top 20 `hotspot_summary` rows plus `monthly_trend` rows for those hotspots

## Why this question is useful

This is a strong benchmark for analytical SQL because it combines:

- thresholded monthly qualification
- recomputation from qualifying raw flights
- percentile aggregation
- top-N ranking
- a mixed result shape with both summary rows and monthly trend rows

It is also a good regression test for prompt quality. Earlier versions of this question allowed materially different hotspot semantics to pass with the same row count. The tightened prompt now tests whether providers converge on the same sustained-hotspot logic.

## Experiment setup

Date:
- `2026-03-16`

Providers:
- Claude `opus`
- Codex `gpt-5.4`
- Gemini `gemini-3.1-pro-preview`

Execution model:
- each provider generated SQL independently
- `qforge` executed the SQL itself against the MCP/OpenAPI-backed `default.ontime_v2`
- query performance metrics were fetched later from `system.query_log`

Compare artifacts:
- [q003-2026-03-16-compare.md](/Users/bvt/work/ExploringDatabyLLMs/runs/q003-2026-03-16-compare.md)
- [q003-2026-03-16-compare.json](/Users/bvt/work/ExploringDatabyLLMs/runs/q003-2026-03-16-compare.json)

Run directories:

- Claude: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001)
- Codex: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001)
- Gemini: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001)

## Result summary

All three providers succeeded.

| runner | model | status | rows | duration_ms | read_rows | memory_usage |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| claude | opus | ok | 832 | 1007 | 1,372,563,373 | 227,604,873 |
| codex | gpt-5.4 | ok | 832 | 1078 | 1,370,586,836 | 226,107,330 |
| gemini | gemini-3.1-pro-preview | ok | 832 | 879 | 1,143,250,214 | 309,491,491 |

All three queries returned the same row count:

- `20` `hotspot_summary` rows
- `812` `monthly_trend` rows

More importantly, all three `result.json` payloads are identical after normalized row comparison.

Normalized row hashes:

- Claude: `2fc9a380082b506e66663e8d3a94eea589a5a0eaf3b61fa4de390cf5cbb14e7a`
- Codex: `2fc9a380082b506e66663e8d3a94eea589a5a0eaf3b61fa4de390cf5cbb14e7a`
- Gemini: `2fc9a380082b506e66663e8d3a94eea589a5a0eaf3b61fa4de390cf5cbb14e7a`

## Full SQL artifacts

Full generated SQL for each provider:

- Claude: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql)
- Codex: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/query.sql)
- Gemini: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-16/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/query.sql)

All three follow the same intended logical shape:

- filter to completed Delta departures from ATL
- qualify monthly `(MonthStart, Dest, DepTimeBlk)` cells at `>= 40` flights
- join back to raw flights from those qualifying cells
- recompute hotspot metrics from qualifying raw flights
- rank hotspot cells and keep the top 20
- emit `hotspot_summary` rows with `MonthStart = NULL`
- emit `monthly_trend` rows for the top hotspots

## Real output differences

There are no verified row-level output differences in this run set.

Pairwise normalized row differences:

| pair | differing rows |
| --- | ---: |
| Claude vs Codex | 0 / 832 |
| Claude vs Gemini | 0 / 832 |
| Codex vs Gemini | 0 / 832 |

This is the key result of the 2026-03-16 run. After the prompt tightening done after the earlier `q003` comparison, all three providers now converge not only on row count but also on the exact final output.

## SQL comparison

The SQL text is still not identical.

Common behavior across all three:

- monthly qualification is explicit
- final hotspot metrics are recomputed from qualifying raw flights rather than averaged from monthly summaries
- `hotspot_summary` rows use `NULL` for `MonthStart`
- final numeric outputs are rounded to 2 decimals

Concrete implementation differences:

- Claude uses a straightforward staged CTE flow with `quantile(0.9)` and computes `DepDel15Pct` as `100.0 * sum(DepDel15) / count()`.
- Codex uses a similar staged flow, also with `quantile(0.9)`, and computes delayed-15 share as `100.0 * avg(toFloat64(DepDel15))`.
- Gemini rounds hotspot-level values earlier in the CTE chain and uses monthly aggregates only for the `monthly_trend` rows, while still recomputing hotspot ranking metrics from qualifying raw flights.

These SQL differences are real, but for this question and this dataset they do not change the final output.

## Execution stats

Execution metrics from `system.query_log`:

| runner | model | duration_ms | read_rows | memory_usage | relative note |
| --- | --- | ---: | ---: | ---: | --- |
| gemini | gemini-3.1-pro-preview | 879 | 1,143,250,214 | 309,491,491 | fastest and lowest read volume, highest memory |
| claude | opus | 1007 | 1,372,563,373 | 227,604,873 | middle runtime, highest read volume, moderate memory |
| codex | gpt-5.4 | 1078 | 1,370,586,836 | 226,107,330 | slightly slower than Claude, nearly identical read volume, lowest memory |

Interpretation:

- performance still differs by query shape even when results are identical
- Gemini achieved the lowest runtime and lowest read volume in this run set
- Codex used the least memory
- Claude and Codex were extremely close on read volume, which is consistent with their similarly structured SQL

The performance interpretation is an inference from the observed query-log metrics plus the SQL structure.

## Takeaway

The main result from the 2026-03-16 `q003` run is convergence.

All three providers:

- returned `832` rows
- satisfied the tightened hotspot semantics
- produced identical normalized outputs

That makes `q003` a much cleaner benchmark than it was in the earlier comparison. It is still useful because execution efficiency differs across providers, but it no longer suffers from the semantic ambiguity that previously allowed materially different hotspot rankings to pass with the same row count.

In short:

- the prompt correction worked
- output agreement is now exact across Claude, Codex, and Gemini
- remaining variation is primarily in SQL shape and execution cost, not in final analytical meaning
