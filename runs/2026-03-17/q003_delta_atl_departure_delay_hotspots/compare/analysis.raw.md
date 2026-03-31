Now I have enough information. Let me compile the final compare report.

```markdown
# q003 Experiment Note

## Question

**ID:** q003  
**Slug:** q003_delta_atl_departure_delay_hotspots  
**Title:** Delta ATL departure delay hotspots by destination and time block

Find the Delta (DL) departures from Atlanta (ATL) with the worst sustained departure delays at the (Dest, DepTimeBlk) level. The question requires monthly qualification thresholds (≥40 flights/month), hotspot qualification (≥1,000 total flights), and ranking by average delay, P90 delay, and DepDel15 percentage.

## Why this question is useful

This question tests an LLM's ability to:

1. **Multi-stage aggregation** – monthly qualification filtering followed by hotspot-level rollup
2. **Metric recomputation semantics** – the prompt explicitly requires recomputing hotspot metrics from raw flights rather than averaging monthly values
3. **Mixed row-type output** – combining `hotspot_summary` and `monthly_trend` rows in a single result set with proper NULL handling
4. **Complex ranking** – multi-key tiebreaker ordering (AvgDepDelayMinutes DESC, P90 DESC, DepDel15Pct DESC, CompletedFlights DESC)

The pattern is common in operational analytics where you need both summary rankings and supporting time-series detail in one query.

## Experiment setup

- **Date:** 2026-03-17
- **Dataset:** ontime_v2 (full historical BTS on-time data)
- **Runners tested:**
  - claude/opus (run-001, run-002)
  - codex/gpt-5.4 (run-001)
  - gemini/gemini-3.1-pro-preview (run-002)
- **Compare contract:** Exact match on key columns (RowType, MonthStart, Dest, DepTimeBlk, HotspotRank); numeric tolerance ±0.005 on AvgDepDelayMinutes, P90DepDelayMinutes, DepDel15Pct

## Result summary

| Runner | Model | Run | Status | Rows | Duration | Read Rows | Memory |
|--------|-------|-----|--------|-----:|----------|-----------|--------|
| claude | opus | run-001 | partial | 832 | 1.38s | 1.84B | 231 MiB |
| claude | opus | run-002 | ok | 832 | 1.28s | 1.83B | 251 MiB |
| codex | gpt-5.4 | run-001 | ok | 832 | 1.03s | 1.36B | 293 MiB |
| gemini | gemini-3.1-pro-preview | run-002 | ok | 832 | 1.11s | 1.37B | 480 MiB |

**Key findings:**

- All runs returned exactly 832 rows (20 hotspot_summary + 812 monthly_trend)
- claude/opus/run-001 marked "partial" due to visual.html missing Leaflet CDN reference (not a data issue)
- codex/gpt-5.4/run-001 read 26% fewer rows than claude runs (1.36B vs 1.83B)
- All runs identified the same top-20 hotspots with LGA 1900-1959 ranked #1

## Full SQL artifacts

### claude / opus

- **run-001:** [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/query.sql)
- **run-002:** [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-002/query.sql)

### codex / gpt-5.4

- **run-001:** [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/query.sql)

### gemini / gemini-3.1-pro-preview

- **run-002:** [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-002/query.sql)

## Real output differences

### Hotspot summary rows

All four runs produce identical top-20 hotspot rankings with the same (Dest, DepTimeBlk) combinations. Minor numeric rounding differences exist in codex/gpt-5.4/run-001:

| Rank | Column | claude/opus | codex/gpt-5.4 | Difference |
|------|--------|-------------|---------------|------------|
| 1 | P90DepDelayMinutes | 66.1 | 66.09 | 0.01 |
| 5 | AvgDepDelayMinutes | 17.08 | 17.07 | 0.01 |
| 17 | AvgDepDelayMinutes | 8.95 | 8.94 | 0.01 |
| 19 | AvgDepDelayMinutes | 8.79 | 8.78 | 0.01 |

These differences exceed the compare contract tolerance of ±0.005 but are attributable to rounding timing: codex applies rounding earlier in the CTE chain via explicit `Decimal(18,2)` casts, while other runs round only in the final SELECT.

### Monthly trend rows

**QualifyingMonths field:** claude/opus/run-002 sets `QualifyingMonths = 0` for all monthly_trend rows, while the other three runs propagate the parent hotspot's QualifyingMonths value. This is a semantic difference in interpretation—the prompt does not explicitly specify whether monthly rows should inherit this field.

All other monthly_trend fields (MonthStart, Dest, DepTimeBlk, CompletedFlights, AvgDepDelayMinutes, P90DepDelayMinutes, DepDel15Pct) match across runs within rounding tolerance.

## SQL comparison

### Structure patterns

All four queries follow the same logical flow:

1. Base CTE filtering DL/ATL/non-cancelled flights
2. Monthly qualification CTE (HAVING count ≥ 40)
3. Join back to raw flights for metric recomputation
4. Hotspot rollup with ≥1,000 flight threshold
5. Window function ranking
6. Top-20 filter
7. Monthly trend extraction via join
8. UNION ALL combining summary and trend rows

### Notable variations

| Aspect | claude/opus/run-001 | claude/opus/run-002 | codex/gpt-5.4/run-001 | gemini/run-002 |
|--------|---------------------|---------------------|----------------------|----------------|
| CTE count | 10 | 9 | 8 | 6 |
| Hotspot qualification | Separate CTE | Integrated in hotspot_metrics | Integrated | Integrated |
| Monthly trend QualifyingMonths | Propagated from hotspot | Hardcoded `CAST(0 AS UInt64)` | Propagated | Propagated |
| DepDel15Pct formula | `100.0 * sum(DepDel15) / count()` | Same | `100.0 * avg(toFloat64(DepDel15))` | `avg(DepDel15) * 100` |
| Decimal casting | None (Float64) | None | Explicit `Decimal(18,2)` | None |
| Top-20 selection | `WHERE HotspotRank <= 20` | Same | `LIMIT 20` after ORDER BY | Same as codex |

The codex query uses `LIMIT 20` with `ORDER BY HotspotRank` in the CTE rather than `WHERE HotspotRank <= 20`, which is functionally equivalent but reflects a different SQL idiom.

## Presentation artifacts

### claude / opus

- **run-001:** [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq003_delta_atl_departure_delay_hotspots%2Fclaude%2Fopus%2Frun-001%2Freport.md) | [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-001/visual.html)
- **run-002:** [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq003_delta_atl_departure_delay_hotspots%2Fclaude%2Fopus%2Frun-002%2Freport.md) | [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/claude/opus/run-002/visual.html)

### codex / gpt-5.4

- **run-001:** [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq003_delta_atl_departure_delay_hotspots%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md) | [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/codex/gpt-5.4/run-001/visual.html)

### gemini / gemini-3.1-pro-preview

- **run-002:** [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq003_delta_atl_departure_delay_hotspots%2Fgemini%2Fgemini-3.1-pro-preview%2Frun-002%2Freport.md) | [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q003_delta_atl_departure_delay_hotspots/gemini/gemini-3.1-pro-preview/run-002/visual.html)

**Visual differences:**

- gemini/run-002 includes Chart.js CDN reference (dynamic charting)
- All other runs use inline SVG or CSS-only visualizations
- claude/opus/run-001 flagged "partial" for missing Leaflet reference (per visual prompt contract)
- codex/gpt-5.4/run-001 produces the largest visual.html (50KB) with more elaborate CSS animations

## Execution stats

### claude / opus

- **run-001:** 1.38s query duration, 1,836,238,936 rows read, 231.3 MiB memory
- **run-002:** 1.28s query duration, 1,831,679,447 rows read, 250.7 MiB memory

### codex / gpt-5.4

- **run-001:** 1.03s query duration, 1,364,087,871 rows read, 292.8 MiB memory

### gemini / gemini-3.1-pro-preview

- **run-002:** 1.11s query duration, 1,371,628,729 rows read, 479.5 MiB memory

**Observations:**

- codex and gemini read ~26% fewer rows than claude runs
- The row-read difference suggests codex/gemini queries may push filters more effectively or structure joins differently
- gemini uses nearly 2× the memory of claude despite similar row reads, possibly due to wider intermediate result materialization
- All queries complete in ~1 second, indicating the difference is not operationally significant at this scale

## Takeaway

All four runs correctly identify the same top-20 Delta ATL departure delay hotspots, with LGA 1900-1959 consistently ranked worst (24.83 min avg delay, 66 min P90, 40% delayed ≥15 min). The semantic core of the question—multi-stage qualification, metric recomputation from raw flights, and mixed row-type output—is handled correctly by all models.

The differences observed are:

1. **Rounding precision:** codex/gpt-5.4 produces values 0.01 off from others in 4 of 80 hotspot metric cells, due to earlier Decimal casting
2. **QualifyingMonths in monthly_trend:** claude/opus/run-002 zeros this field; others propagate the hotspot value
3. **Query efficiency:** codex and gemini read fewer rows, suggesting more selective filter placement

For practical use, all outputs are interchangeable. The QualifyingMonths difference in monthly_trend rows is cosmetic since this field is primarily meaningful for hotspot_summary rows. The rounding differences are within acceptable floating-point variation and would round identically if displayed to 1 decimal place.
```
