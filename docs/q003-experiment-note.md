# q003 Experiment Note

## Question

`q003` asks the model to use `default.ontime_v2` to find Delta departures out of ATL with the worst sustained departure delays at the `(Dest, DepTimeBlk)` level.

The prompt requires:

- completed flights only: `Cancelled = 0`
- Delta ATL departures only: `IATA_CODE_Reporting_Airline = 'DL'` and `Origin = 'ATL'`
- monthly qualification at `(MonthStart, Dest, DepTimeBlk)`
- hotspot qualification across qualifying monthly cells only
- hotspot ranking by:
  - `AvgDepDelayMinutes` descending
  - `P90DepDelayMinutes` descending
  - `DepDel15Pct` descending
  - `CompletedFlights` descending
- top 20 hotspot summary rows plus monthly trend rows for those hotspots

## Why this question is useful

This is a good benchmark for analytical SQL because it combines:

- thresholded monthly qualification
- rollup from qualified monthly cells to hotspot cells
- percentile aggregation
- rank-based top-N extraction
- a mixed result shape with summary rows and monthly trend rows

It is also a good consistency check because several SQL shapes can return the same row count while still computing hotspot metrics differently.

## Experiment setup

Date:
- `2026-03-15`

Providers:
- Claude `opus`
- Codex `gpt-5.4`
- Gemini `gemini-3-flash-preview`

Execution model:
- each provider generated SQL independently
- `qforge` executed the SQL itself against the MCP/OpenAPI-backed `default.ontime_v2`
- performance metrics were fetched later from `system.query_log`

Compare artifact:
- [q003-compare.md](/Users/bvt/work/ExploringDatabyLLMs/runs/q003-compare.md)

Run directories:

- Claude: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001)
- Codex: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001)
- Gemini: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3-flash-preview/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3-flash-preview/run-001)

## Result summary

All three providers succeeded.

| runner | model | status | rows | duration_ms | read_rows | memory_usage |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| claude | opus | ok | 832 | 1385 | 1,877,464,390 | 250,528,074 |
| codex | gpt-5.4 | ok | 832 | 900 | 1,146,680,615 | 283,105,876 |
| gemini | gemini-3-flash-preview | ok | 832 | 587 | 689,178,284 | 479,511,521 |

All three queries returned the same row count:

- `20` `hotspot_summary` rows
- `812` `monthly_trend` rows

But the three `result.json` files are not identical, and the differences are semantic rather than cosmetic.

## Full SQL artifacts

Full generated SQL for each provider:

- Claude: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql)
- Codex: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/query.sql)
- Gemini: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3-flash-preview/run-001/query.sql)

All three use the same broad pattern:

- qualify monthly `(MonthStart, Dest, DepTimeBlk)` cells with at least `40` completed flights
- keep hotspot cells with at least `1,000` flights across qualifying months
- rank hotspot cells and emit top-20 summaries
- join those hotspot cells back to qualified monthly rows for trend output

The important differences are in how the final hotspot metrics are computed and normalized.

## Real output differences

The three `result.json` files are not identical.

Normalized row hashes:

- Claude: `8ee922898317eedfdf6e4a970e8b23d47d49835b03b8c9127cbc2a1f538bd9c2`
- Codex: `260ce95cbf3a9dcf9edd7c776f800362c4780ef3ce04a450283c9a77210a9bea`
- Gemini: `9b7f89342717b98024e38ee0de6837e63a2dbd03f472ae9c50f5cd96188f0d55`

Pairwise row-level differences:

| pair | differing rows | main causes |
| --- | ---: | --- |
| Claude vs Codex | 826 / 832 | Claude rounds most metrics and uses `1970-01-01` as the `MonthStart` sentinel for `hotspot_summary`; Codex keeps full precision and uses `NULL`. Hotspot ordering is the same. |
| Claude vs Gemini | 819 / 832 | Gemini computes hotspot metrics differently, changes 4 hotspot ranks, uses `NULL` summary `MonthStart`, and keeps unrounded values. |
| Codex vs Gemini | 683 / 832 | Gemini changes 4 hotspot ranks and computes different `AvgDepDelayMinutes`, `P90DepDelayMinutes`, and `DepDel15Pct` values; Codex recomputes hotspot metrics from raw qualifying flights. |

Concrete examples from the actual outputs:

1. `hotspot_summary` `MonthStart`
- Claude:
  - `MonthStart = '1970-01-01'`
- Codex:
  - `MonthStart = NULL`
- Gemini:
  - `MonthStart = NULL`

2. Top-8 hotspot ordering
- Claude and Codex:
  - rank `6`: `SLC`, `1500-1559`
  - rank `7`: `LAX`, `1700-1759`
  - rank `8`: `MCO`, `2200-2259`
- Gemini:
  - rank `6`: `SLC`, `1500-1559`
  - rank `7`: `MCO`, `2200-2259`
  - rank `8`: `LAX`, `1700-1759`

3. Mid-table ordering
- Claude and Codex:
  - rank `14`: `EWR`, `1400-1459`
  - rank `15`: `DFW`, `1500-1559`
- Gemini:
  - rank `14`: `DFW`, `1500-1559`
  - rank `15`: `EWR`, `1400-1459`

4. Metric computation differences for the top hotspot (`LGA`, `1900-1959`)
- Claude:
  - `AvgDepDelayMinutes = 24.83`
  - `P90DepDelayMinutes = 66.1`
  - `DepDel15Pct = 39.75`
- Codex:
  - `AvgDepDelayMinutes = 24.82883435582822`
  - `P90DepDelayMinutes = 66.380005`
  - `DepDel15Pct = 39.75460122699386`
- Gemini:
  - `AvgDepDelayMinutes = 24.905597355192945`
  - `P90DepDelayMinutes = 64.4896551724138`
  - `DepDel15Pct = 39.936996164113474`

The Gemini values are not just formatting differences. They come from a different aggregation strategy.

## SQL comparison

### Claude

Strengths:
- clear staged CTE flow
- enforces monthly qualification before hotspot ranking
- recomputes hotspot aggregates from raw flights restricted to qualifying monthly cells
- rounded result shape is presentation-friendly

Differences:
- uses `toDate('1970-01-01')` instead of `NULL` for summary `MonthStart`
- rounds metrics in the final output, which creates many row-level diffs relative to Codex even when the ranking matches
- uses exact `quantile(0.9)` rather than an approximate percentile aggregate

### Codex

Strengths:
- closest match to the prompt’s intended hotspot semantics
- computes monthly qualification first, then recomputes hotspot metrics from raw qualifying flights
- uses `NULL` for summary `MonthStart`
- keeps full-precision numeric outputs

Differences:
- uses `quantileTDigest(0.9)`, so percentile values are approximate rather than exact
- reads more data than Gemini because it makes a second pass over qualified raw flights

### Gemini

Strengths:
- understandable structure and valid result shape
- applies both the monthly and hotspot flight-count thresholds
- fastest runtime in this run set

Real problems:
- computes hotspot metrics by averaging monthly aggregates instead of recomputing from raw qualifying flights
- averages monthly p90 values, which is not the same as the p90 over all qualifying flights in a hotspot
- changes hotspot ordering in 4 of the top-20 positions
- highest memory usage of the three runs

This matters because equal row counts do not guarantee equal hotspot semantics.

## Execution stats

Execution metrics from `system.query_log`:

| runner | model | duration_ms | read_rows | memory_usage | relative note |
| --- | --- | ---: | ---: | ---: | --- |
| gemini | gemini-3-flash-preview | 587 | 689,178,284 | 479,511,521 | fastest, least data read, highest memory |
| codex | gpt-5.4 | 900 | 1,146,680,615 | 283,105,876 | middle runtime, moderate read volume, moderate memory |
| claude | opus | 1385 | 1,877,464,390 | 250,528,074 | slowest, highest read volume, lowest memory |

Interpretation:

- performance differences were driven by query shape, not by result size
- Gemini appears faster because it performs less raw-flight recomputation, but that simplification also changes the hotspot metrics
- Claude and Codex converge on the same hotspot ordering even though one rounds aggressively and the other keeps full precision

The final point is an inference from the SQL structure plus the observed `read_rows` differences.

## Takeaway

`q003` is a good benchmark for thresholded, sustained-delay analysis because it separates:

- monthly qualification logic
- hotspot rollup semantics
- percentile handling
- ranking stability under slightly different aggregation strategies

The main lesson from this run is:

- row-count convergence is not enough
- hotspot semantics depend on whether the model recomputes final metrics from qualifying raw flights or averages monthly summaries

For this question, Codex matched the intended semantics most closely, Claude was very close but normalized the output differently, and Gemini produced a faster but materially different hotspot ranking.

## Prompt tightening after this run

After this comparison, the `q003` prompt and compare contract were tightened to reduce avoidable output drift across providers.

The revised contract now makes three things explicit:

- final hotspot metrics must be recomputed from raw flights belonging to qualifying monthly cells
- `hotspot_summary` rows must use `MonthStart = NULL`
- `AvgDepDelayMinutes`, `P90DepDelayMinutes`, and `DepDel15Pct` must be rounded to 2 decimals

This keeps the benchmark focused on sustained-hotspot semantics rather than allowing multiple incompatible rollup interpretations to pass as equivalent.
