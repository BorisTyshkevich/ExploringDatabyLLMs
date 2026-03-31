Now I have all the data I need. Let me write the report.

```markdown
# q001 Experiment Note

**Question:** Highest daily hops for one aircraft on one flight number
**Day:** 2026-03-26
**Runs:** claude/opus/run-002 · claude/sonnet/run-001 · codex/gpt-5.4/run-001

---

## Question

Find the longest single-day itineraries — measured by hop count — for one aircraft operating under the same flight number. Define uniqueness by the full textual route string, return the 10 most recent unique routes at the maximum hop count, include all required fields (aircraft id, flight number, carrier, date, hop count, recurrence count, route), and deduplicate same-time legs before counting.

Full prompt: [report_prompt.md](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/report_prompt.md) · [visual_prompt.md](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/visual_prompt.md)

---

## Why this question is useful

This question probes several capabilities at once: window-function or aggregation-based deduplication of raw flight legs, correct route-string construction in chronological leg order, a global-maximum filter that requires a correlated subquery or two-pass aggregation, and a recurrence join that scans all history rather than just the top rows. It also exercises the edge-case instruction to retain empty tail numbers rather than silently filtering them — a detail that separates careful prompt-following from loose interpretation.

---

## Experiment setup

| Field | Value |
|---|---|
| Dataset | `ontime.fact_ontime` (BTS On-Time Performance) |
| Question ID | `q001` |
| Subqueries | 3 per run (main, q1, q2) |
| Runners | claude/opus, claude/sonnet, codex/gpt-5.4 |
| Presentation target | HTML visual |

All three runs use the multi-query harness: `main.sql` finds the top-10 routes, `q1.sql` analyses recurrence tiers, and `q2.sql` examines the geographic pattern of the airports in those routes.

---

## Result summary

All three runs agree on the headline answer: **8 hops** (9 airports) is the global maximum, achieved exclusively by Southwest Airlines (WN). All runs return the same 10 distinct route strings in the same order. The first-row result is identical across all three:

- Aircraft **N957WN**, flight **WN 366**, date **2024-12-01**, route **ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA**

Two recurrence counts differ by ±1 between opus and the other two runs (details in §7). These small differences trace to a design choice: opus does not filter cancelled flights in `main.sql`, while sonnet and codex both add `WHERE Cancelled = 0`.

claude/sonnet/run-001 finished with status `partial` / verdict `WARN` due to two factual errors in the q2 prose narrative (state count and California-endpoint count), documented in the review. The SQL and data outputs for sonnet are otherwise sound.

codex/gpt-5.4/run-001 skipped presentation generation entirely (visual.html was not produced by the model; the phase is marked `skipped` in the manifest).

---

## Full SQL artifacts

### claude / opus

- **run-002**
  - [main.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/main.sql)
  - [q1.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/q1.sql)
  - [q2.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/claude/opus/run-002/queries/q2.sql)
  - [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/claude/opus/run-002/result.json)

### claude / sonnet

- **run-001**
  - [main.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/main.sql)
  - [q1.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/q1.sql)
  - [q2.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/queries/q2.sql)
  - [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/result.json)

### codex / gpt-5.4

- **run-001**
  - [main.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/main.sql)
  - [q1.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/q1.sql)
  - [q2.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/queries/q2.sql)
  - [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/result.json)

---

## Real output differences

The 10 route strings and their ordering are **identical** across all three runs. All three agree on hop_count = 8, carrier = WN, and aircraft IDs for every row. The only numeric differences are in two recurrence counts:

| Route | opus/run-002 | sonnet/run-001 | codex/run-001 |
|---|---:|---:|---:|
| CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN | 5 | 4 | 4 |
| LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN | 47 | 46 | 46 |
| All other 8 routes | agree | agree | agree |

The two discrepancies (±1) align exactly with the cancelled-flight filter: opus/main.sql does not apply `WHERE Cancelled = 0`, while sonnet and codex both do. Each off-by-one is consistent with one cancelled occurrence of that exact route string existing in the data.

The sonnet run carries two additional prose errors in q2 (not SQL errors), noted in its review: the state count was reported as 27 (plus DC) instead of the correct 24, and the California-endpoint count was reported as 7/10 instead of the correct 6/10. These are narrative arithmetic mistakes on the already-executed result set; the underlying q2 SQL and data are not affected.

---

## SQL comparison

All three runs share the same CTE pipeline shape: a deduplication step → itinerary aggregation with route-string construction → global-max filter → recurrence join → window-ranked LIMIT 10. Key divergences:

**Cancelled-flight filtering.** opus/main.sql scans all rows with no `WHERE Cancelled = 0` guard. sonnet/main.sql and codex/main.sql both add this filter in their first CTE. This is the likely source of the ±1 recurrence discrepancy above.

**Hop-counting method.** opus uses `uniqExact(OriginCode, DestCode, CRSDepTime)` in a subquery candidates CTE to identify max-hop days, then does a `SELECT DISTINCT` join to expand legs. sonnet groups by `(FlightDate, Tail_Number, carrier, flight_num, OriginCode, DestCode)` with `min(CRSDepTime)` to deduplicate, then counts rows per itinerary. codex similarly groups by those six columns plus `dep_sort` (with a coalesce fallback: `coalesce(CRSDepTime, DepTime, 9999)`) and tracks both `dep_sort` and `arr_sort` for sorting.

**Route-string construction.** All three use ClickHouse array functions to concatenate origins and append the final destination. opus sorts within a subquery before `groupArray`. sonnet uses `arraySort(groupArray(tuple(dep_time, OriginCode)))` inside the aggregate. codex uses `arraySort(groupArray((dep_sort, arr_sort, OriginCode, DestCode)))` with a 4-tuple that also captures arrival time for tiebreaking.

**Tie-breaking in the final ORDER BY.** The prompt requires "most recent top 10 unique routes by departure time" but does not specify a strict tiebreaker beyond date. opus orders by `FlightDate DESC`. sonnet orders by `most_recent_date DESC`. codex adds secondary tiebreakers (`first_dep_sort DESC, carrier DESC, flight_number DESC, aircraft_id DESC`). All produce the same ranked output, suggesting no ties exist at this date boundary.

**q2 dimension join.** opus joins to an airport dimension and returns `DisplayAirportName`, `Latitude`, `Longitude`. sonnet joins a different dimension view returning `city`, `state`, `state_code`. codex builds its own airport list from `fact_ontime` using a frequency count per airport code, returning `geo_bucket` and `airport_mentions` without a dimension join.

---

## Presentation artifacts

### claude / opus

- **run-002**
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-002%2Freport.md): Full narrative with labelled table of all 10 routes, per-route recurrence list, and geographic observations. All claims verified against result.json by the review.
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-26/q001_hops_per_day/claude/opus/run-002/visual.html): Generated in 241 s.
  - [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-002%2Freview.md): Verdict PASS — no defects found.

### claude / sonnet

- **run-001**
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fsonnet%2Frun-001%2Freport.md): Narrative with the same 10-route table and tiered recurrence analysis. Contains two prose errors in the q2 section (state count: 27 stated, 24 actual; CA endpoints: 7 stated, 6 actual).
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/visual.html): Generated in 679 s.
  - [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-26%2Fq001_hops_per_day%2Fclaude%2Fsonnet%2Frun-001%2Freview.md): Verdict WARN — main and q1 pass; q2 flagged for two miscounts in prose.

### codex / gpt-5.4

- **run-001**
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-26%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md): Concise narrative presenting the main result in prose and table form; q1 recurrence list presented inline as a single sentence rather than a numbered list. No factual errors found by the review.
  - visual.html: Presentation generation was **skipped** for this run; no HTML visual was produced.
  - [review.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-26%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-001%2Freview.md): Verdict PASS — no substantive defects found; MCP spot-check of itinerary continuity passed.

---

## Execution stats

| Provider | Model | Run | Status | Query time | Rows read | Bytes read | Peak memory | SQL gen time | Visual gen time | Total duration |
|---|---|---|---|---:|---:|---:|---:|---:|---:|---:|
| claude | opus | run-002 | ok / PASS | 122.4 s | 3,075,342,183 | 38.0 GiB | 49.1 GiB | 654 s | 241 s | 1,164 s |
| claude | sonnet | run-001 | partial / WARN | 214.4 s | 3,720,509,274 | 43.9 GiB | 59.2 GiB | 676 s | 679 s | 1,633 s |
| codex | gpt-5.4 | run-001 | ok / PASS | 306.3 s | 3,386,975,101 | 52.8 GiB | 57.0 GiB | n/a | skipped | 482 s |

Query time shows a 2.5× spread between the fastest (opus, 122 s) and slowest (codex, 306 s). The difference is consistent with architectural choices: opus's two-pass candidate-then-expand pattern reads the fewest rows and bytes; codex's 4-tuple sort with coalesce fallback and route-recurrence pre-computation scans roughly 310 M more rows than opus. Despite the longer query time, codex's total wall-clock duration (482 s) is the shortest of the three because it skipped presentation generation entirely and does not report a separate SQL-generation timer; the claude runs spent 895–1355 s in model-side generation phases alone.

---

## Takeaway

All three runs converge on the correct headline answer (8-hop max, WN-only, same 10 routes) and pass or near-pass review. The main differentiator is depth of correctness at the prose level: opus produced a fully verified narrative with no reviewer flags; codex produced a compact but factually clean output; sonnet introduced two countable errors in the q2 geographic summary that the automated review caught. The ±1 recurrence discrepancy between opus and the other two runs is a real, small, and explainable data difference tied to the cancelled-flight filter — neither approach is strictly wrong given the ambiguous prompt, but including cancelled flights in recurrence counting is arguably undesirable. On efficiency, opus is the clear winner on query resources; on total wall-clock time, codex is fastest by virtue of skipping the visual phase.
```
