```markdown
# q001 Experiment Note

## Question

`q001` asks for the highest number of same-day hops flown by a single aircraft on a single flight number, with `Aircraft ID`, `Flight Number`, `Carrier`, `Date`, `Route`, and the actual departure-time sequence for the longest trip. The SQL prompt also requires the top 10 longest and most recent itineraries, with routes built in chronological leg order.

## Why this question is useful

This benchmark is a good stress test for itinerary reconstruction rather than simple aggregation. A successful answer has to group by aircraft, flight number, carrier, and date; sort legs by actual departure time; preserve the full airport chain; and still rank the longest itineraries correctly. It also exposes whether a model returns the exact requested schema or only an approximation of it.

## Experiment setup

The note is based on the verified local artifacts for `2026-03-17`:

- Compare summary: [`compare.json`](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/compare/compare.json)
- Question prompt: [`prompt.md`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/prompt.md)
- Compare contract: [`compare.yaml`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/compare.yaml)
- Report prompt: [`report_prompt.md`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/report_prompt.md)
- Visual prompt: [`visual_prompt.md`](https://github.com/boristyshkevich/ExploringDatabyLLMs/blob/main/prompts/q001_hops_per_day/visual_prompt.md)

Verified runs:

- `claude / opus / run-003`
- `codex / gpt-5.4 / run-003`
- `gemini / gemini-3.1-pro-preview / run-003`

## Result summary

All three runs returned 10 rows and the same 10 keyed itineraries. The top row is the same in every `result.json`: aircraft `N957WN`, flight `366`, carrier `WN`, date `2024-12-01T00:00:00Z`, with route content corresponding to `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA`.

The maximum observed hop count in the saved outputs is 8. Claude reports `Hops = 8` explicitly on every row; Codex’s route strings contain 8 timed leg segments per row; Gemini’s `Actual Departure Times` arrays contain 8 departure times per row.

After normalizing route formatting, the top 10 collapse to 4 airport sequences:

- `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA` appears 1 time
- `CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN` appears 4 times
- `ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN` appears 2 times
- `MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX` appears 3 times

## Full SQL artifacts

### claude / opus
- `run-003`: [query.sql](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/claude/opus/run-003/query.sql), [result.json](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/claude/opus/run-003/result.json)

### codex / gpt-5.4
- `run-003`: [query.sql](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/query.sql), [result.json](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/result.json)

### gemini / gemini-3.1-pro-preview
- `run-003`: [query.sql](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/query.sql), [result.json](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/result.json)

## Real output differences

The outputs are identical at the itinerary-content level, but not byte-identical in representation.

- Claude and Codex encode the same departure times for every keyed row.
- Gemini encodes the same airport sequence as Claude/Codex, and its `Actual Departure Times` arrays match the departure times embedded in Claude/Codex route strings for all 10 rows.
- Claude is the only run that returns the requested `Hops` column explicitly.
- Codex omits `Hops` and returns only `Aircraft ID`, `Flight Number`, `Carrier`, `Date`, and `Route`.
- Gemini also omits `Hops`, and instead adds `Actual Departure Times`.

The raw `Route` field differs by format:

- Claude: `05:43 ISP → 08:10 BWI → ... → SEA`
- Codex: `05:43 ISP->BWI | 08:10 BWI->MYR | ... | 20:41 OAK->SEA`
- Gemini: `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA`

One additional metadata inconsistency exists in Gemini’s saved output: its `result.json` `log_comment` says `run=run-004` even though the directory and compare artifact identify `run-003`.

## SQL comparison

Claude’s SQL is the closest match to the prompt shape. It uses two CTEs, filters `Cancelled = 0`, groups by aircraft/flight/carrier/date, computes `count() AS Hops`, sorts legs by `DepTime`, and emits both `Hops` and a time-ordered route ending with the final destination.

Codex’s SQL is more compact and was the fastest run, but it does not project `Hops`. It also adds `Diverted = 0` and `Flight_Number_Reporting_Airline != ''`, ranks by `count() DESC`, and formats `Route` as timed `Origin->Dest` segments separated by `|`.

Gemini’s SQL also does not project `Hops`. It groups by `Reporting_Airline` rather than `IATA_CODE_Reporting_Airline`, does not filter `Cancelled = 0` or `Diverted = 0`, emits an airport-only `Route`, and adds `Actual Departure Times` as a separate array.

## Presentation artifacts

### claude / opus
- `run-003`: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-003%2Freport.md), [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q001_hops_per_day/claude/opus/run-003/visual.html)
- The report includes the full result table with `Hops`.
- The visual includes Leaflet, KPI cards, a lead-itinerary map, a route-sequence panel, and a result table with `Hops` and `Route`.

### codex / gpt-5.4
- `run-003`: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-003%2Freport.md), [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q001_hops_per_day/codex/gpt-5.4/run-003/visual.html)
- The report mirrors the saved result schema and therefore omits `Hops`.
- The visual is the richest of the three in terms of explicit degraded-state messaging: it keeps the map visible, describes airport-coordinate enrichment from `default.airports_bts`, includes route clustering and route-sequence sections, and renders a leaderboard with a derived hop-count display.

### gemini / gemini-3.1-pro-preview
- `run-003`: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-17%2Fq001_hops_per_day%2Fgemini%2Fgemini-3.1-pro-preview%2Frun-003%2Freport.md), [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-17/q001_hops_per_day/gemini/gemini-3.1-pro-preview/run-003/visual.html)
- The report omits `Hops` and instead shows `Actual Departure Times`.
- The visual includes Leaflet, KPI cards, a lead-itinerary map, a route-sequence section, and degraded-map messaging for enrichment failure.
- `compare.json` marks this run `partial` because `presentation_render` failed, even though both `report.md` and `visual.html` were written.

## Execution stats

### claude / opus
- `run-003`: status `ok`; query time `18.63 s`; read rows `193,061,941`; read bytes `2,650,424,488`; memory `35.1 GiB`; lowest read volume and lowest memory usage of the three.

### codex / gpt-5.4
- `run-003`: status `ok`; query time `7.88 s`; read rows `193,061,941`; read bytes `2,654,458,009`; memory `45.3 GiB`; fastest successful run, about 2.36x faster than Claude and 1.97x faster than Gemini.

### gemini / gemini-3.1-pro-preview
- `run-003`: status `partial`; query time `15.53 s`; read rows `230,307,587`; read bytes `2,763,997,654`; memory `48.9 GiB`; highest read volume and highest memory usage.

## Takeaway

For `q001`, the three runs agree on the actual top-10 itineraries after normalizing route formatting, so this benchmark is more about output fidelity and execution profile than about answer disagreement. Claude is the only run that returns the requested `Hops` column directly; Codex matches the itinerary content with the best query time but leaves hop count implicit in the route string; Gemini also preserves the itinerary content, adds explicit departure-time arrays, but its run is only `partial` because the presentation render did not finish cleanly.
```