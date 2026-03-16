# q002 Experiment Note

## Question

`q002` asks the model to use `default.ontime_v2` to determine the yearly flight-volume leader by `Reporting_Airline`, then identify the sharpest true leadership transition.

The prompt requires:

- completed flights only: `Cancelled = 0`
- yearly carrier ranking
- leader and runner-up share gap
- year-over-year leader share change
- top 5 carriers per year in the result
- transition logic that excludes the first year

## Why this question is useful

This is a good benchmark for structured analytical SQL because it combines:

- grouped aggregation
- ranking within partitions
- leader/runner-up extraction
- window functions over yearly leaders
- a final result shape that supports both reporting and visualization

It is also a good consistency check because multiple SQL formulations can return the same row count while still differing in output semantics.

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
- [q002-compare.md](/Users/bvt/work/ExploringDatabyLLMs/runs/q002-compare.md)

Run directories:

- Claude: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/claude/opus/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/claude/opus/run-001)
- Codex: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001)
- Gemini: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/gemini/gemini-3-flash-preview/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/gemini/gemini-3-flash-preview/run-001)

## Result summary

All three providers succeeded.

| runner | model | status | rows | duration_ms | read_rows | memory_usage |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| claude | opus | ok | 195 | 617 | 921,230,348 | 367,303,463 |
| codex | gpt-5.4 | ok | 195 | 603 | 921,230,348 | 360,158,400 |
| gemini | gemini-3-flash-preview | ok | 195 | 850 | 921,230,348 | 385,315,402 |

All three queries scanned the same amount of data and returned the same row count, which indicates the core analytical intent converged across models.

## Full SQL artifacts

Full generated SQL for each provider:

- Claude: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/claude/opus/run-001/query.sql)
- Codex: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/codex/gpt-5.4/run-001/query.sql)
- Gemini: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q002_top_carrier_by_flights_leadership/gemini/gemini-3-flash-preview/run-001/query.sql)

All three use the same broad pattern:

- aggregate completed flights by `Year, Reporting_Airline`
- compute yearly totals
- rank carriers within year
- extract leader and runner-up
- compute year-over-year leader share change
- emit top-5 rows per year

The important differences are semantic rather than structural.

## Real output differences

The three `result.json` files are not identical.

Normalized row hashes:

- Claude: `3138e1f00038e6338580228a1586c5ad06763c76f913f54a2c005e1c48500ad1`
- Codex: `b892018972daa0caf29f2462603b7840842eebbf0a3cdd3d08cb251e58a2a8ce`
- Gemini: `6be98385de3c56f03df52579401da0135e3bd9c2fdb7874589aa194faeae6740`

Pairwise row-level differences:

| pair | differing rows | main causes |
| --- | ---: | --- |
| Claude vs Codex | 157 / 195 | Claude repeats leader-transition fields on all top-5 rows; Codex only fills them on rank-1 rows. First-year prior leader is `''` in Claude vs `NULL` in Codex. |
| Claude vs Gemini | 195 / 195 | Gemini emits unrounded share values and treats the first year as a transition with a non-null share-change value. |
| Codex vs Gemini | 195 / 195 | Gemini repeats leader fields on all rows, leaves first-year transition logic wrong, and uses unrounded numeric values. |

Concrete examples from the actual outputs:

1. First-year handling (`1987`, leader row)
- Claude:
  - `PriorYearLeaderReportingAirline = ''`
  - `LeaderChanged = 0`
  - `LeaderShareChangePctPts = null`
- Codex:
  - `PriorYearLeaderReportingAirline = null`
  - `LeaderChanged = 0`
  - `LeaderShareChangePctPts = null`
- Gemini:
  - `PriorYearLeaderReportingAirline = ''`
  - `LeaderChanged = 1`
  - `LeaderShareChangePctPts = 14.221048631689575`

2. Non-leader rows in the same year
- Claude and Gemini repeat leader metadata on rank 2-5 rows.
- Codex sets leader-transition fields to `NULL` on non-leader rows.

3. Numeric normalization
- Claude and Codex round share values and share-gap values.
- Gemini returns full-precision floats, for example:
  - `SharePct = 14.221048631689575`
  - `LeaderShareGapPctPts = 1.5935567403247788`

## Data interpretation

Direct data checks show only three real leadership transitions in the dataset:

- `1990`: `DL` -> `US`
- `1992`: `US` -> `DL`
- `2000`: `DL` -> `WN`

The strongest transition by the revised question definition is `1990`, because the prompt now defines “most sharply” only over years where the annual leader actually changed.

## SQL comparison

### Claude

Strengths:
- clear CTE decomposition
- rounded percentage outputs
- correct first-year `LeaderChanged = 0`

Differences:
- repeats leader-transition fields on all top-5 rows
- uses empty string rather than `NULL` for first-year prior leader
- uses `argMax(... if(...))` patterns that are less direct than `maxIf`

### Codex

Strengths:
- cleanest match to the revised contract
- leader-transition fields populated only on rank-1 rows
- uses `NULL` for first-year prior leader
- explicitly casts nullable values in the final projection

Differences:
- more verbose final projection because it nulls non-leader rows explicitly
- uses `lagInFrame(... ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING)`, which is explicit but noisy

### Gemini

Strengths:
- overall query structure is still understandable and close to the other two
- correctly identifies leaders and runner-up rows

Real problems:
- first year is treated as a transition
- first year gets a non-null `LeaderShareChangePctPts`
- leader-transition fields are repeated on all top-5 rows
- numeric outputs are not normalized to the same rounded format as the others

This matters because equal row counts do not guarantee equal semantics.

## Execution stats

Execution metrics from `system.query_log`:

| runner | model | duration_ms | read_rows | memory_usage | relative note |
| --- | --- | ---: | ---: | ---: | --- |
| codex | gpt-5.4 | 603 | 921,230,348 | 360,158,400 | fastest and lowest memory |
| claude | opus | 617 | 921,230,348 | 367,303,463 | nearly identical to Codex |
| gemini | gemini-3-flash-preview | 850 | 921,230,348 | 385,315,402 | slowest and highest memory |

Interpretation:

- performance differences were small in absolute terms
- the database work was effectively identical across runs
- the meaningful differences came from SQL semantics and output shaping, not scan volume

## Prompt correction made during the experiment

The original q002 wording had an ambiguity:

- it asked where leadership changed most sharply
- but its ordering rule allowed non-transition years with large leader share swings to outrank real leadership changes

That was fixed so:

- the first year is explicitly not a transition
- “most sharply” applies only to years where the leader changed

This made the question more internally consistent without changing thresholds or reducing difficulty.

## Takeaway

`q002` is a good example of a benchmark question that is data-satisfiable and model-solvable, but still benefits from careful prompt tightening.

The main lesson from this run is:

- performance convergence is easy to observe
- semantic convergence needs stronger contracts

For this question, Codex, Claude, and Gemini all produced usable SQL, but Codex matched the revised output contract most closely.
