The `compare_report.md` has been written. Here's the report:

```markdown
# q001 Experiment Note

## Question

**q001 — Highest daily hops for one aircraft on one flight number**

Find the top 10 cases where a single aircraft (tail number) flew the most legs in one day under the same flight number. Return the aircraft ID, flight number, carrier, date, reconstructed route, and hop count. The Route column must list every leg and the final destination, delimited by `-` without spaces.

## Why this question is useful

This question stress-tests an LLM's ability to aggregate multi-leg flight segments into ordered itineraries. It requires grouping by (tail number, flight number, date), sorting legs chronologically by departure time, concatenating airport codes into a readable route string, and ranking the results by hop count and recency. The route-construction step is especially revealing — it forces the model to choose between `groupArray` with sorting versus self-joins, and to correctly include both origin and destination airports without duplication. The explicit delimiter rule (`-` without spaces) further tests whether the model follows formatting constraints.

## Experiment setup

- **Date**: 2026-03-23
- **Dataset**: Altinity OnTime (`ontime.fact_ontime`, ClickHouse)
- **Prompt**: [report_prompt.md](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/report_prompt.md)
- **Visual prompt**: [visual_prompt.md](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/visual_prompt.md)
- **Compare contract**: [compare.yaml](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/compare.yaml)
- **Runners (4 runs)**:
  - claude / opus — run-001
  - claude / sonnet — run-001
  - codex / gpt-5.4 — run-001
  - gemini / gemini-3.1-pro-preview — run-002

## Result summary

All four runs returned **10 rows with identical values** for Aircraft ID, Flight Number, Carrier, Date, and hop_count. Every row is a Southwest Airlines (WN) flight with 8 hops. The top result across all runs is N957WN / WN 366 on 2024-12-01, route ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA.

The **only difference** across the four result sets is the Route delimiter:

| Run | Delimiter | Prompt-compliant? |
| --- | --- | --- |
| claude/opus/run-001 | `-` (hyphen) | Yes |
| claude/sonnet/run-001 | ` → ` (spaced arrow) | No |
| codex/gpt-5.4/run-001 | `→` (compact arrow) | No |
| gemini/gemini-3.1-pro-preview/run-002 | `-` (hyphen) | Yes |

The gemini run finished with **partial** status because its `presentation_render` phase failed; all other phases completed successfully.

## Full SQL artifacts

### claude / opus
- **run-001**: [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-23/q001_hops_per_day/claude/opus/run-001/query.sql) · [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-23/q001_hops_per_day/claude/opus/run-001/result.json)

### claude / sonnet
- **run-001**: [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-23/q001_hops_per_day/claude/sonnet/run-001/query.sql) · [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-23/q001_hops_per_day/claude/sonnet/run-001/result.json)

### codex / gpt-5.4
- **run-001**: [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-23/q001_hops_per_day/codex/gpt-5.4/run-001/query.sql) · [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-23/q001_hops_per_day/codex/gpt-5.4/run-001/result.json)

### gemini / gemini-3.1-pro-preview
- **run-002**: [query.sql](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-23/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-002/query.sql) · [result.json](https://github.com/BorisTyshkevich/ExploringDatabyLLMs/blob/main/runs/2026-03-23/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-002/result.json)

## Real output differences

The 10 result rows are **data-identical** across all four runs. Every row matches on Aircraft ID, Flight Number, Carrier, Date, and hop_count. The only difference is in the `Route` string delimiter:

| Run | Route example (row 1) | Delimiter |
| --- | --- | --- |
| claude/opus/run-001 | `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` | `-` |
| claude/sonnet/run-001 | `ISP → BWI → MYR → BNA → VPS → DAL → LAS → OAK → SEA` | ` → ` |
| codex/gpt-5.4/run-001 | `ISP→BWI→MYR→BNA→VPS→DAL→LAS→OAK→SEA` | `→` |
| gemini/gemini-3.1-pro-preview/run-002 | `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` | `-` |

The airport sequences are identical in all 10 rows across all four runs. No differences exist in any other column.

## SQL comparison

All four queries solve the problem with the same high-level strategy: group legs by tail number + flight number + date, sort legs by departure time, then assemble the route string. They differ in structure, column naming, filtering, and route assembly.

| Aspect | claude/opus | claude/sonnet | codex/gpt-5.4 | gemini |
| --- | --- | --- | --- | --- |
| CTEs | 2 (`legs`, `grouped`) | 1 (`sorted`) | 2 (`leg_rows`, `itineraries`) | 0 (single SELECT) |
| Column names | `Tail_Number`, `Flight_Number_Reporting_Airline`, `IATA_CODE_Reporting_Airline` | same as opus | `TailNum`, `FlightNum`, `Carrier` (old-style) | `Tail_Number`, `Flight_Number_Reporting_Airline`, `Reporting_Airline` |
| Cancel/divert filter | `Cancelled = 0` | `Cancelled = 0` | `Cancelled = 0 AND Diverted = 0` | `Cancelled = 0` |
| Empty-tail filter | `Tail_Number != ''` + `Flight_Number != ''` | `Tail_Number != ''` | `TailNum != ''` + `FlightNum != ''` | `Tail_Number != ''` |
| Dep-time sort key | `assumeNotNull(CRSDepTime)` | `assumeNotNull(CRSDepTime)` | `ifNull(DepTime, CRSDepTime)` → DateTime | `DepTime` (actual) |
| HAVING threshold | `hop_count >= 2` | none | `length(ordered_legs) > 1` | none |
| Route construction | `arrayPushBack` origins + last dest | `arrayConcat` first origin + dests | `arrayStringConcat` origins + `concat` last dest | `arrayPushBack` sorted origins + `argMax(DestCode, DepTime)` |
| Route delimiter | `-` | ` → ` | `→` | `-` |
| ORDER BY tiebreakers | `hop_count DESC, FlightDate DESC` | `hop_count DESC, FlightDate DESC` | `hop_count DESC, Date DESC, Aircraft ID, Flight Number` | `hop_count DESC, Date DESC` |

Despite these differences in filtering, column aliases, and sort tiebreakers, all four queries arrive at the same top-10 result set because the highest-hop rows (8 hops each) are unambiguous. Gemini's single-SELECT approach is the most compact but uses `Reporting_Airline` (same values as `IATA_CODE_Reporting_Airline` for this data) and raw `DepTime` instead of CRSDepTime. Codex's query is the most defensive (excludes diverted flights, requires non-empty flight number, null-safe departure time) but reads ~22% more bytes due to additional column access.

## Presentation artifacts

### claude / opus
- **run-001**:
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-23%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-001%2Freport.md) — Flat table with a one-paragraph analytical note about Southwest's point-to-point model.
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-23/q001_hops_per_day/claude/opus/run-001/visual.html) — Leaflet map dashboard.

### claude / sonnet
- **run-001**:
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-23%2Fq001_hops_per_day%2Fclaude%2Fsonnet%2Frun-001%2Freport.md) — Titled "Top 10 Longest Itineraries" section with a full Markdown table and an analytical takeaway on hub-and-spoke routing.
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-23/q001_hops_per_day/claude/sonnet/run-001/visual.html) — Leaflet map dashboard.

### codex / gpt-5.4
- **run-001**:
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-23%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md) — Compact table with a one-sentence summary about repeat operating patterns.
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-23/q001_hops_per_day/codex/gpt-5.4/run-001/visual.html) — Leaflet map dashboard with SRI integrity attributes on Leaflet CSS.

### gemini / gemini-3.1-pro-preview
- **run-002**:
  - [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs/md.html?file=runs%2F2026-03-23%2Fq001_hops_per_day%2Fgemini%2Fgemini-3.1-pro-preview%2Frun-002%2Freport.md) — Table with commentary on Southwest's through-flight network model.
  - [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs/runs/2026-03-23/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-002/visual.html) — HTML was generated but the **presentation_render phase failed** (status: partial), so the dashboard may not render correctly.

All dashboards follow the shared visual prompt's color palette and layout conventions but differ in header treatment and CSS layering.

## Execution stats

| Provider / Model | Run | Status | Query time | Rows read | Bytes read | Peak memory | SQL gen | Visual gen | Total duration |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| claude / opus | run-001 | ok | 7.54 s | 193,061,941 | 2.47 GiB | 42.8 GiB | 87.3 s | 207.7 s | 313 s |
| claude / sonnet | run-001 | ok | 7.42 s | 193,061,941 | 2.47 GiB | 42.5 GiB | 96.3 s | 546.8 s | 661 s |
| codex / gpt-5.4 | run-001 | ok | 7.55 s | 193,061,941 | 3.01 GiB | 41.3 GiB | 84.2 s | 438.6 s | 541 s |
| gemini / gemini-3.1-pro-preview | run-002 | partial | 14.59 s | 193,061,941 | 2.47 GiB | 48.9 GiB | 355.1 s | 265.8 s | 656 s |

Query execution times cluster tightly at 7.4–7.6 s for the three fully successful runs, while gemini's query took roughly 2x longer (14.6 s) and consumed 14% more peak memory (48.9 GiB vs. 41–43 GiB). The largest performance spread is in SQL generation time: gemini took 355 s — 4.2x slower than codex's 84 s — which accounts for most of its longer total duration. Codex read ~0.5 GiB more data than the others (3.01 GiB vs. 2.47 GiB), likely due to pulling additional `DepTime` and `Diverted` columns.

## Takeaway

All four models correctly solved this multi-leg itinerary reconstruction problem and returned identical top-10 rows. The only output difference is a cosmetic route delimiter: claude/opus and gemini followed the prompt's explicit `-` rule, while claude/sonnet and codex/gpt-5.4 substituted arrow characters — a recurring pattern where LLMs inject "prettier" formatting defaults over explicit instructions. SQL approaches ranged from gemini's zero-CTE single SELECT to codex's defensive two-CTE strategy with diversion filtering, but none of these structural differences changed the result. On performance, claude/opus achieved the fastest end-to-end completion (313 s) driven by its 208 s visual generation time, while gemini was slowest across the board with 2x query latency, 4x SQL generation time, and a failed visual render phase.
```
