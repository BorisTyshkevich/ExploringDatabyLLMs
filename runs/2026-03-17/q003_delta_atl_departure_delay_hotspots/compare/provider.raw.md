```markdown
# q003 Experiment Note

## Question
`q003` asks for Delta (`DL`) departures from ATL, qualified first at monthly `(MonthStart, Dest, DepTimeBlk)` grain, then re-aggregated from raw qualifying flights to rank the top 20 hotspot cells and emit their monthly trend rows. The source prompt is [prompt.md](/Users/bvt/work/ExploringDatabyLLMs/prompts/q003_delta_atl_departure_delay_hotspots/prompt.md), with reporting guidance in [report_prompt.md](/Users/bvt/work/ExploringDatabyLLMs/prompts/q003_delta_atl_departure_delay_hotspots/report_prompt.md) and comparison rules in [compare.yaml](/Users/bvt/work/ExploringDatabyLLMs/prompts/q003_delta_atl_departure_delay_hotspots/compare.yaml).

## Why this question is useful
This benchmark is useful because it forces a model to preserve aggregation semantics across two grains. It is easy to get superficially plausible SQL that averages monthly aggregates instead of recomputing hotspot metrics from raw qualifying flights, or that mishandles the mixed output shape (`hotspot_summary` plus `monthly_trend`).

## Experiment setup
The verified structured summary is [compare.json](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/compare/compare.json). I verified the five saved SQL artifacts, the four available `result.json` files, and the successful `report.md` and `visual.html` artifacts directly from the run directories.

## Result summary
Three runs finished with non-empty results: `codex/gpt-5.4/run-001`, `claude/opus/run-002`, and `gemini/gemini-3.1-pro-preview/run-002`. Each returned 832 rows with the same key set: 20 `hotspot_summary` rows and 812 `monthly_trend` rows.

Across those successful runs, the hotspot ranking is stable. Rank 1 is `LGA` at `1900-1959`, with 1,630 completed flights, 29 qualifying months, average departure delay `24.83`, `P90DepDelayMinutes` about `66.1`, `DepDel15Pct` `39.75`, and qualifying span `1998-11-01` to `2001-03-01`. The top five cells are identical by rank and key: `LGA 1900-1959`, `LGA 1500-1559`, `MCO 1700-1759`, `DFW 1700-1759`, `ORD 1500-1559`.

`gemini/.../run-001` did not produce a result set: [compare.json](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/compare/compare.json) marks `sql_execution` failed, row count `0`, and warns that `result.json` is missing.

## Full SQL artifacts
- [claude/opus/run-002/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-002/query.sql)
- [codex/gpt-5.4/run-001/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/query.sql)
- [gemini/gemini-3.1-pro-preview/run-002/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-002/query.sql)
- [claude/opus/run-001/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql)
- [gemini/gemini-3.1-pro-preview/run-001/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-001/query.sql)

## Real output differences
The successful outputs are not byte-identical, but the differences are narrow and verifiable.

- `claude/opus/run-001` and `gemini/.../run-002` match exactly at the row-content level.
- `codex/gpt-5.4/run-001` differs from those two on 135 of 832 rows, only by `0.01` in already-rounded numeric fields: `AvgDepDelayMinutes` on 61 rows, `P90DepDelayMinutes` on 46 rows, and `DepDel15Pct` on 35 rows. Keys, row counts, hotspot membership, and ranking stay the same.
- `claude/opus/run-002` has the same 0.01-level numeric differences as above, plus one substantive schema/output difference: all 812 `monthly_trend` rows set `QualifyingMonths = 0`. The compare contract marks `QualifyingMonths` as an exact-compare field in [compare.yaml](/Users/bvt/work/ExploringDatabyLLMs/prompts/q003_delta_atl_departure_delay_hotspots/compare.yaml), so this is a real output discrepancy rather than formatting noise.

## SQL comparison
All successful SQL files follow the intended broad shape: filter ATL Delta completed flights, qualify monthly cells at `>= 40`, rejoin raw qualifying flights, rank hotspots at `>= 1000` flights, and emit both summary and trend rows.

The main verified SQL differences are:
- [codex/gpt-5.4/run-001/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/query.sql) and [gemini/.../run-002/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-002/query.sql) carry hotspot `QualifyingMonths` through to monthly rows.
- [claude/opus/run-002/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-002/query.sql) explicitly writes `CAST(0 AS UInt64) AS QualifyingMonths` in `monthly_trend_rows`, which explains the 812-row mismatch there.
- [codex/gpt-5.4/run-001/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/query.sql) adds `Dest` and `DepTimeBlk` as extra tie-breakers in the `row_number()` ordering; the others do not. That did not change the observed top-20 ranking on this day.
- [claude/opus/run-001/query.sql](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql) uses a placeholder date in the summary CTE and nulls it in the final `SELECT`; the others emit `NULL` directly.

## Presentation artifacts
The successful report artifacts all surface the same top-20 table, but differ in narrative specificity.
- [codex/gpt-5.4/run-001/report.md](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/report.md) stays mostly template-like and does not name the worst hotspot in prose.
- [claude/opus/run-002/report.md](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-002/report.md) adds headings and framing, but the prose remains generic.
- [gemini/gemini-3.1-pro-preview/run-002/report.md](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-002/report.md) is the most concrete: it names `LGA 1900-1959` and summarizes the top-5 pattern directly from the result.
- The successful visuals all expose a heatmap, a top-20 ranking table, and a top-3 monthly trend view with embedded saved SQL: [codex/gpt-5.4/run-001/visual.html](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/visual.html), [claude/opus/run-001/visual.html](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual.html), [claude/opus/run-002/visual.html](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-002/visual.html), and [gemini/gemini-3.1-pro-preview/run-002/visual.html](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-002/visual.html).

## Execution stats
From [compare.json](/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/compare/compare.json):
- Fastest successful query: `codex/gpt-5.4/run-001` at `1.03 s`.
- Lowest read volume: `codex/gpt-5.4/run-001` at `1,364,087,871` rows.
- Lowest memory among successful runs: `claude/opus/run-002` at `250.7 MiB`.
- `gemini/.../run-001` failed before producing a result, with `query_duration_ms = 8`, `read_rows = 0`, and only the warning about missing `result.json`.

## Takeaway
This question cleanly separates semantic correctness from superficial similarity. On `2026-03-17`, all successful runs found the same hotspot ranking and returned the same row keys, but they did not produce identical outputs: Codex differed from the Claude-1/Gemini-2 pair only by 0.01-level rounded metrics, while Claude run-002 introduced a real contract-visible discrepancy by zeroing `QualifyingMonths` on every monthly trend row.
```