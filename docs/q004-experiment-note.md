# q004 Experiment Note

## Question

`q004` asks the model to use `default.ontime_v2` to identify the 25 worst origin airports for departure on-time performance after excluding low-volume airports.

The prompt requires:

- completed flights only: `Cancelled = 0`
- aggregation at `Origin`
- a minimum threshold of `50,000` completed departures over full history
- ranking by:
  - `DepartureOtpPct` ascending
  - `AvgDepDelayMinutes` descending
  - `CompletedDepartures` descending
  - `Origin` ascending
- required metrics:
  - `CompletedDepartures`
  - `DepartureOtpPct`
  - `AvgDepDelayMinutes`
  - `P90DepDelayMinutes`
  - `FirstFlightDate`
  - `LastFlightDate`

## Why this question is useful

This is a good benchmark for a compact analytical ranking query because it combines:

- threshold filtering
- stable multi-column ranking
- percentile calculation
- date range aggregation
- a result that is easy to compare across models

It also exposes an important benchmark-design issue: percentile functions are not interchangeable, and prompt wording that says only “use quantile logic” leaves room for materially different answers.

## Experiment setup

Date:
- `2026-03-15`

Providers:
- Claude `opus`
- Codex `gpt-5.4`
- Gemini `gemini-3-flash-preview`

Execution model:
- each provider generated SQL independently
- `qforge` executed the SQL against `default.ontime_v2`
- query performance came from deferred `system.query_log` lookup

Compare artifact:
- [q004-compare.md](/Users/bvt/work/ExploringDatabyLLMs/runs/q004-compare.md)

Run directories:

- Claude: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/claude/opus/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/claude/opus/run-001)
- Codex: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/codex/gpt-5.4/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/codex/gpt-5.4/run-001)
- Gemini: [/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/gemini/gemini-3-flash-preview/run-001](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/gemini/gemini-3-flash-preview/run-001)

## Result summary

All three providers succeeded and returned `25` rows.

| runner | model | status | rows | duration_ms | read_rows | memory_usage |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| claude | opus | ok | 25 | 343 | 230,307,587 | 287,751,941 |
| codex | gpt-5.4 | ok | 25 | 527 | 460,615,174 | 555,270,949 |
| gemini | gemini-3-flash-preview | ok | 25 | 572 | 460,615,174 | 374,894,860 |

The top of the ranked result is stable across all three runs:

1. `ASE`
2. `MDW`
3. `ACV`
4. `SFB`
5. `ORD`

## Full SQL artifacts

Full generated SQL for each provider:

- Claude: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/claude/opus/run-001/query.sql)
- Codex: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/codex/gpt-5.4/run-001/query.sql)
- Gemini: [query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-15/q004_worst_origin_airport_otp_thresholded/gemini/gemini-3-flash-preview/run-001/query.sql)

All three queries use the same broad plan:

- filter to completed flights
- identify qualifying airports above the threshold
- aggregate airport-level metrics
- order by the prompt’s ranking rules
- `LIMIT 25`

The main SQL-level difference is how the threshold step is organized:

- Claude computes all airport stats first, then filters with a `qualifying` CTE.
- Codex creates a `completed_departures` base CTE, then a separate `qualifying_origins` CTE, then joins back.
- Gemini uses a `HAVING` threshold CTE and joins it back to the base table.

## Real output differences

The three `result.json` files are not identical.

Normalized row hashes:

- Claude: `c6633c381cd66fd32946896dcf7f7ad282556654e1df1fd402fd19a85dade652`
- Codex: `05c4b82285d0c5a3ab677b3607a532aeb16666748f6d523bf61f6a7d174038f6`
- Gemini: `4fd360ac11976a9c2852d7bd62474bfedc7a15564d6f483236dbdb724d8afd67`

But the divergence is narrow and easy to localize.

Pairwise row-level differences:

| pair | differing rows | differing field |
| --- | ---: | --- |
| Claude vs Codex | 23 / 25 | `P90DepDelayMinutes` only |
| Claude vs Gemini | 19 / 25 | `P90DepDelayMinutes` only |
| Codex vs Gemini | 20 / 25 | `P90DepDelayMinutes` only |

When `P90DepDelayMinutes` is ignored, all three results are identical row-for-row:

- same airport ordering
- same `CompletedDepartures`
- same `DepartureOtpPct`
- same `AvgDepDelayMinutes`
- same `FirstFlightDate`
- same `LastFlightDate`

Concrete examples from the actual outputs:

| airport | Claude p90 | Codex p90 | Gemini p90 |
| --- | ---: | ---: | ---: |
| ASE | 63 | 68 | 67 |
| MDW | 36 | 39 | 36 |
| ACV | 64 | 63 | 62 |
| SFB | 45.9 | 48 | 47 |
| ORD | 44 | 42 | 45 |

## What a direct data check shows

A direct validation query against the current dataset confirms the ranking itself:

- `ASE` is the worst qualifying airport by departure OTP
- followed by `MDW`, `ACV`, `SFB`, and `ORD`

But it also shows that percentile values depend strongly on the exact quantile function:

For a spot check on the top 5 airports:

| airport | `quantile(0.9)` | `quantileExact(0.9)` | `quantileTDigest(0.9)` |
| --- | ---: | ---: | ---: |
| ASE | 67 | 68 | 68.1457 |
| MDW | 36 | 37 | 36.9260 |
| ACV | 63 | 62 | 62.4139 |
| SFB | 48 | 48 | 47.6566 |
| ORD | 39 | 43 | 42.7800 |

That means q004 currently leaves too much freedom in the percentile definition. All three providers used acceptable “quantile logic,” but they did not converge on the same p90 output.

## SQL comparison

### Claude

Strengths:
- shortest and clearest query
- clean threshold CTE
- rounded output values

Observed issue:
- `P90DepDelayMinutes` differs from the other providers and from direct spot checks on `quantile(0.9)`

### Codex

Strengths:
- strongest separation between base rows, thresholding, and final aggregation
- explicit `toFloat64(ifNull(..., 0))` for delay minutes

Observed issue:
- zero-filling `DepDelayMinutes` before the percentile step changes the p90 definition
- this likely explains part of the percentile drift
- it also scanned roughly twice as many rows as Claude

### Gemini

Strengths:
- concise threshold CTE using `HAVING`
- simplest join-back structure

Observed issue:
- leaves `P90DepDelayMinutes` unrounded
- percentile outputs still differ from the other two and from direct spot checks

## Execution stats

| runner | model | duration_ms | read_rows | memory_usage | relative note |
| --- | --- | ---: | ---: | ---: | --- |
| claude | opus | 343 | 230,307,587 | 287,751,941 | fastest and most efficient |
| codex | gpt-5.4 | 527 | 460,615,174 | 555,270,949 | slowest scan footprint and highest memory |
| gemini | gemini-3-flash-preview | 572 | 460,615,174 | 374,894,860 | similar scan volume to Codex, lower memory |

Interpretation:

- Claude produced the most efficient execution plan for this question.
- Codex and Gemini scanned roughly double Claude’s `read_rows`.
- The execution gap is real here, unlike q002 where database work was nearly identical.

## Takeaway

`q004` is a strong example of a benchmark question where:

- the business result is stable across models
- the execution cost is not stable
- the percentile metric is underspecified

The key lesson is that q004 needs a stricter percentile contract if it is going to be used for semantic comparison. Right now the models agree on the airport ranking, but not on the p90 metric definition.
