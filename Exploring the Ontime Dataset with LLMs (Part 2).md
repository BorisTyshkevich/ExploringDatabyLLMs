# Exploring the Ontime Dataset with LLMs (Part 2)

A reproducible benchmark of model-generated ClickHouse SQL on 200M+ flight rows.

## TL;DR

- We benchmarked 4 models in batch mode on 6 real SQL tasks over Ontime.
- Environment: Altinity demo `default.ontime` (as of March 2, 2026: `1987..2021`, `201,575,308` rows).
- **Executability:** 100% for all tested models.
- **Strict exact-match score:** low (format/precision differences penalized heavily).
- **Relaxed business-answer score:** high (all models produced the same business answers after normalization).
- Main risk is no longer syntax; it is subtle semantic drift, precision policy, and output-shape inconsistency.

## Why this benchmark matters

Part 1 covered architecture and security. This part answers a practical question:

> If we ask different LLMs the same analytics questions, how good is their generated ClickHouse SQL in real conditions?

Too many comparisons rely on screenshots or one-off prompts. Here we used scripts, saved SQL, and executed every query against a live ClickHouse dataset.

## Benchmark setup

### Date and data

Tested on **March 2, 2026**.

Dataset profile on Altinity demo:

```sql
SELECT min(Year) AS min_year, max(Year) AS max_year, count() AS rows
FROM default.ontime;
```

Observed:

- `min_year = 1987`
- `max_year = 2021`
- `rows = 201,575,308`

### Models tested

- OpenAI `gpt-5` (Codex CLI batch mode)
- OpenAI `gpt-5.3-codex` (Codex CLI batch mode)
- Anthropic `sonnet` (Claude Code batch mode)
- Anthropic `opus` (Claude Code batch mode)

### Cases (6)

1. Top carrier by flights in 2019.
2. Average departure delay for Delta out of ATL in July 2021.
3. Worst on-time origin airport in 2021 (minimum flight threshold).
4. Worst winter carrier-airport pair (2019-2021, minimum flight threshold).
5. Peak AA delay month across 2018-2021.
6. Highest-delay route-season among top 50 routes (2019-2021).

### Prompting mode

Schema-aware prompt with explicit case-sensitive column hints:

- `Year, Month, Carrier, Origin, Dest, DepDelay, ArrDel15`
- instruction to output one SQL query only

### Scoring

We used two score layers:

1. **Execution score**: query runs successfully.
2. **Strict match**: generated output equals expected output exactly.

Because exact match can be too harsh (rounding, aliasing, fixed-string padding), we also inspected:

3. **Relaxed business-answer match**: same business outcome after normalization.

## Results

### Aggregate leaderboard

| Model | Execution OK | Strict Match | Relaxed Match | Total Cases |
|---|---:|---:|---:|---:|
| openai_gpt5 | 6 | 1 | 6 | 6 |
| openai_gpt53codex | 6 | 1 | 6 | 6 |
| anthropic_sonnet | 6 | 1 | 6 | 6 |
| anthropic_opus | 6 | 2 | 6 | 6 |

### Key interpretation

- All models were operationally usable for these tasks (100% execution success).
- Strict exact-match underestimates practical quality because many "mismatches" were harmless formatting differences:
  - `avg(DepDelay)` vs `round(avg(DepDelay), 2)`
  - fixed-width string padding (`DAL\x00\x00`) vs trimmed strings
  - equivalent route representation (`Origin + Dest` vs prebuilt `route` alias)
- After normalization, all four models converged to the same business answers for these 6 tasks.
- Full strict/relaxed scoring artifacts are published in `results/benchmark/benchmark_summary.json` and `results/benchmark/benchmark_relaxed_summary.json`.

## What each model got right and wrong

## OpenAI `gpt-5`

Pros:

- Stable SQL generation.
- Strong instruction following on output fields.
- Correct ranking/grouping logic in complex route-season query.

Cons:

- Often returned unrounded metrics where rounded output was requested.
- Could pass execution while still violating output-format constraints.

## OpenAI `gpt-5.3-codex`

Pros:

- Fastest OpenAI model in this run.
- Consistent query structure and deterministic ordering.

Cons:

- Same precision/format drift as `gpt-5`.
- More likely to optimize for runnable SQL than strict reporting format.

## Anthropic `sonnet`

Pros:

- Solid execution reliability.
- Efficient SQL with reasonable plans for this workload.

Cons:

- More frequent formatting deviations (integer rounding in one aggregate case).
- Returned padded fixed strings when not explicitly trimmed.

## Anthropic `opus`

Pros:

- Highest strict score in this run.
- Best adherence on the most complex route-season task.

Cons:

- Slower on some cases.
- Still exhibited output-format drift on multiple tasks.

## Hard critique: weaker points in the overall approach

1. **Schema-aware prompting can hide real-world brittleness.**
   If production users ask free-form questions without schema context, failure rates will rise.

2. **Execution success is not correctness.**
   A query that runs can still encode the wrong business logic.

3. **Strict equality can be too strict; relaxed matching can be too lenient.**
   You need both views plus manual review for high-impact metrics.

4. **Single dataset benchmarks can overfit.**
   Ontime is useful, but production workloads include joins, slowly changing dimensions, and messy source semantics.

## Pro/Con summary of the core idea

### Using LLMs as SQL copilots for ClickHouse

Pros:

- High velocity for exploration.
- Lower barrier for non-SQL users.
- Good enough operational reliability with guardrails.

Cons:

- Semantics and formatting drift are still common.
- Results can look correct while violating reporting policy.
- Requires validation, not just generation.

## Recommended production workflow

1. Use MCP with read-only credentials and query limits.
2. Generate SQL from LLM.
3. Auto-validate:
   - schema checks
   - forbidden-pattern checks
   - cost checks
4. Execute in sandboxed/read-only mode.
5. Post-validate output shape and KPI tolerances.
6. Require human approval for high-impact decisions.

## Corrected SQL examples for the two common failure classes

### 1) Case sensitivity and fixed strings

```sql
SELECT trimRight(toString(Carrier)) AS carrier, count() AS flights
FROM default.ontime
WHERE Year = 2019
GROUP BY carrier
ORDER BY flights DESC
LIMIT 1;
```

### 2) Time-window alignment with actual dataset horizon

```sql
SELECT round(avg(DepDelay), 4) AS avg_dep_delay
FROM default.ontime
WHERE trimRight(toString(Carrier)) = 'DL'
  AND trimRight(toString(Origin)) = 'ATL'
  AND Year = 2021
  AND Month = 7;
```

(Using `Year = 2022` on this specific dataset returns no rows because max year is 2021.)

## Reproducibility assets

All benchmark artifacts are included in this project:

- `scripts/benchmark_llm_sql.py`
- `results/benchmark/benchmark_summary.json`
- `results/benchmark/benchmark_report.md`
- generated SQL per model in `results/benchmark/sql/`

## Conclusion

For this Ontime workload, the tested OpenAI and Anthropic models were all good at generating executable ClickHouse SQL when given schema-aware prompts. The differentiator is not "can it run?" but "does it follow reporting semantics exactly?"

If your use case is exploratory analysis, LLM + MCP is ready now.
If your use case is audited business reporting, add strict validation gates and human review.

---

If you want a follow-up Part 3, the next logical step is a governance-focused benchmark:

- blind prompts (no schema hints)
- cost-aware query scoring
- semantic diff checks against approved KPI definitions
- pass/fail gates for production promotion
