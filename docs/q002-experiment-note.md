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
- [q002-compare.md](/Users/bvt/work/OnTime-LLM/runs/q002-compare.md)

## Result summary

All three providers succeeded.

| runner | model | status | rows | duration_ms | read_rows | memory_usage |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| claude | opus | ok | 195 | 617 | 921,230,348 | 367,303,463 |
| codex | gpt-5.4 | ok | 195 | 603 | 921,230,348 | 360,158,400 |
| gemini | gemini-3-flash-preview | ok | 195 | 850 | 921,230,348 | 385,315,402 |

All three queries scanned the same amount of data and returned the same row count, which indicates the core analytical intent converged across models.

## Data interpretation

Direct data checks show only three real leadership transitions in the dataset:

- `1990`: `DL` -> `US`
- `1992`: `US` -> `DL`
- `2000`: `DL` -> `WN`

The strongest transition by the revised question definition is `1990`, because the prompt now defines “most sharply” only over years where the annual leader actually changed.

## What differed across models

The SQLs were not identical.

Codex produced the cleanest match to the revised contract:
- leader-transition fields populated only on the rank-1 row
- first-year prior-leader fields treated as null
- clear separation of yearly totals, ranking, and transition logic

Claude and Gemini were still logically valid, but looser:
- they repeated leader-transition fields on all top-5 rows
- Claude used empty strings for first-year prior-leader handling
- Gemini left some numeric fields less normalized

This matters because equal row counts do not guarantee equal semantics.

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
