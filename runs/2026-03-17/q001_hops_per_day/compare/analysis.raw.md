Now I have all the artifacts I need to write the comparison report. Let me compile the findings:

```markdown
# q001 Experiment Note

**Question ID:** `q001_hops_per_day`  
**Day:** `2026-03-17`  
**Generated:** `2026-03-17T20:20:04Z`

---

## Question

> Find the highest number of hops per day for a single aircraft using the same flight number.  
> Report Aircraft ID, Flight Number, Carrier, Date, Route.  
> For the longest trip, show the actual departure time from each origin.  
> What does the itinerary look like? Find the top 10 longest and most recent itineraries.

---

## Why this question is useful

This question tests an LLM's ability to:

1. **Aggregate flight legs** into multi-hop itineraries by grouping on tail number, flight number, carrier, and date
2. **Order leg sequences chronologically** using departure times
3. **Construct human-readable route strings** that show airport progression with timing
4. **Rank and filter** by hop count while using recency as a tiebreaker

It reveals differences in how models handle array aggregation, tuple sorting, and string concatenation in ClickHouse.

---

## Experiment setup

Three provider/model combinations were evaluated:

| Provider | Model | Run ID | Status |
|----------|-------|--------|--------|
| claude | opus | run-003 | ok |
| codex | gpt-5.4 | run-003 | ok |
| gemini | gemini-3.1-pro-preview | run-003 | partial (presentation_render failed) |

All runs queried the `default.ontime_v2` dataset with a limit of 10 rows. The compare.json artifact was generated at `2026-03-17T20:20:04Z`.

---

## Result summary

**Data agreement:** All three runs returned the **same 10 itineraries** with identical Aircraft IDs, Flight Numbers, Carriers, and Dates. The core analytical result is consistent across models.

**Column differences:**
- **claude/opus** returned 6 columns including an explicit `Hops` count
- **codex/gpt-5.4** returned 5 columns (no explicit `Hops` column)
- **gemini/gemini-3.1-pro-preview** returned 6 columns including `Actual Departure Times` as a raw integer array

**Route formatting differences:**
- claude/opus: `05:43 ISP → 08:10 BWI → ... → SEA` (times + arrows + final destination)
- codex/gpt-5.4: `05:43 ISP->BWI | 08:10 BWI->MYR | ...` (pipe-delimited leg segments)
- gemini/gemini-3.1-pro-preview: `ISP-BWI-MYR-BNA-...` (dash-separated airport codes only)

**Max hop count:** All runs found **8 hops** as the maximum, with the lead itinerary being **N957WN** on flight **WN 366** on **2024-12-01**.

---

## Full SQL artifacts

### claude / opus
- **run-003**: [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/query.sql) | [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/result.json)

### codex / gpt-5.4
- **run-003**: [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/query.sql) | [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/result.json)

### gemini / gemini-3.1-pro-preview
- **run-003**: [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/query.sql) | [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/result.json)

---

## Real output differences

**Row-level agreement:** The 10 rows are identical across all runs for the core fields (Aircraft ID, Flight Number, Carrier, Date). The same itineraries appear in the same rank order.

**Route string differences (Row 1 example):**

| Model | Route representation |
|-------|---------------------|
| claude/opus | `05:43 ISP → 08:10 BWI → 10:20 MYR → 11:42 BNA → 14:01 VPS → 16:43 DAL → 18:28 LAS → 20:41 OAK → SEA` |
| codex/gpt-5.4 | `05:43 ISP->BWI \| 08:10 BWI->MYR \| 10:20 MYR->BNA \| 11:42 BNA->VPS \| 14:01 VPS->DAL \| 16:43 DAL->LAS \| 18:28 LAS->OAK \| 20:41 OAK->SEA` |
| gemini | `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` (times in separate array: `[543, 810, 1020, 1142, 1401, 1643, 1828, 2041]`) |

**Column schema differences:**

| Model | Columns |
|-------|---------|
| claude/opus | Aircraft ID, Flight Number, Carrier, Date, **Hops**, Route |
| codex/gpt-5.4 | Aircraft ID, Flight Number, Carrier, Date, Route |
| gemini | Aircraft ID, Flight Number, Carrier, Date, Route, **Actual Departure Times** |

Only claude/opus provides an explicit `Hops` column. Gemini separates departure times into a raw integer array rather than embedding them in the route string.

---

## SQL comparison

**Structural approach:**

| Model | CTE usage | Filter conditions | Route construction |
|-------|-----------|-------------------|-------------------|
| claude/opus | 2 CTEs (`flight_legs`, `grouped`) | `Cancelled=0`, `DepTime IS NOT NULL`, `Tail_Number != ''` | `arrayStringConcat` with time formatting via `DIV/MOD` |
| codex/gpt-5.4 | Inline CTE via `groupArray` in SELECT | `Cancelled=0`, `Diverted=0`, `DepTime IS NOT NULL`, `Tail_Number != ''`, `Flight_Number != ''` | `arrayStringConcat` with `formatDateTime` for time |
| gemini | No CTE, direct aggregation | `length(Tail_Number) > 0`, `DepTime IS NOT NULL` | `arrayPushBack` + `argMax` for final destination |

**Key differences:**

1. **Diverted flights:** codex/gpt-5.4 explicitly excludes diverted flights (`Diverted=0`); others do not
2. **Carrier column:** gemini uses `Reporting_Airline` while others use `IATA_CODE_Reporting_Airline`
3. **Hop count:** claude/opus calculates `count() AS Hops` in the CTE; codex uses `count()` only for ordering; gemini derives hops implicitly
4. **Time formatting:** claude uses integer arithmetic (`DIV/MOD`); codex uses `formatDateTime`; gemini outputs raw integers
5. **Route terminator:** claude appends final destination via array indexing; codex includes origin-dest pairs; gemini uses `argMax(Dest, DepTime)`

---

## Presentation artifacts

### claude / opus
- **run-003**: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-003%2Freport.md) | [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/visual.html)
  - Report includes analytical commentary on route repetition, milk-run patterns, and carrier clustering
  - Visual is a dynamic dashboard with Leaflet map, KPI cards, query ledger with expandable SQL, and CSV export

### codex / gpt-5.4
- **run-003**: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-003%2Freport.md) | [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/visual.html)
  - Report is more concise, focusing on route pattern interpretation and table rendering
  - Visual includes route clustering analysis, carrier filtering, and repeat-pattern detection

### gemini / gemini-3.1-pro-preview
- **run-003**: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-17%2Fq001_hops_per_day%2Fgemini%2Fgemini-3.1-pro-preview%2Frun-003%2Freport.md) | [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/visual.html)
  - Report includes analytical summary sections; departure times shown as raw integer arrays
  - Visual has simpler layout but includes map and sequence table; **presentation_render phase failed** (partial status)

---

## Execution stats

### claude / opus
- **run-003**
  - Duration: **18.63 s**
  - Read rows: **193,061,941**
  - Read bytes: 2.65 GB
  - Memory: **35.1 GiB** (lowest)

### codex / gpt-5.4
- **run-003**
  - Duration: **7.88 s** (fastest)
  - Read rows: **193,061,941**
  - Read bytes: 2.65 GB
  - Memory: **45.3 GiB**

### gemini / gemini-3.1-pro-preview
- **run-003**
  - Duration: **15.53 s**
  - Read rows: **230,307,587** (+19% more rows scanned)
  - Read bytes: 2.76 GB
  - Memory: **48.9 GiB** (highest)

**Performance notes:**
- codex/gpt-5.4 was **2.4× faster** than claude/opus despite identical row reads
- gemini scanned **37M additional rows** (19% more), likely due to missing the `Cancelled=0` filter
- claude/opus had the lowest memory footprint, possibly due to CTE structure

---

## Takeaway

All three models correctly identified the same top-10 itineraries with the same max hop count of 8. The **core analytical result is identical**.

Differences are primarily cosmetic:
- **Route formatting:** claude embeds times inline; codex uses pipe-delimited segments; gemini separates times into an array
- **Column schema:** only claude provides an explicit `Hops` column
- **Filter strictness:** codex excludes diverted flights; gemini's looser filters caused higher row reads

For **query efficiency**, codex/gpt-5.4 produced the fastest execution despite comparable I/O. For **memory efficiency**, claude/opus used ~27% less memory than gemini. The gemini run had a partial failure in the presentation rendering phase, resulting in a degraded visual artifact.

All three approaches are valid interpretations of the prompt; the choice depends on whether the consumer prefers inline timing, explicit hop counts, or structured array output.
```
