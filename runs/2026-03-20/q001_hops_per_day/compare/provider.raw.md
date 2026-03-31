```markdown
# q001 Experiment Note

## Question

`q001_hops_per_day` asks for the highest number of same-day hops flown by a single aircraft on a single flight number, with the top 10 longest and most recent itineraries, full route strings, and per-leg departure times.

## Why this question is useful

This prompt tests several failure-prone behaviors at once:

- reconstructing a multi-leg itinerary in chronological order from flight-leg rows,
- keeping route text consistent with hop count,
- ranking by both hop count and recency,
- and presenting repeated-route patterns rather than just a single maximum row.

It is also useful because the compare contract tolerates header aliases but still requires exact agreement on the business result: key columns plus exact `Hops` and `Route`.

## Experiment setup

The note is based on verified local artifacts from `compare.json`, the three `query.sql` files, the three `result.json` files, and the three presentation artifacts (`report.md`, `visual.html`).

The question prompt requires chronological leg ordering, a route string with every leg plus final destination, and the top 10 longest and most recent itineraries. The compare contract normalizes header aliases across runs to `Aircraft ID`, `Flight Number`, `Carrier`, and `Date`, then compares `Hops` and `Route` exactly.

## Result summary

All three runs returned 10 rows, and after normalizing the alias/header differences required by the compare contract, all three runs produced the same 10 business-result rows in the same order.

The shared result is:

- maximum observed hop count: `8`
- all top-10 rows are carrier `WN`
- the most recent max-hop itinerary is `WN 366`, aircraft `N957WN`, on `2024-12-01`, route `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA`
- repeated routes in the top 10 are:
  - `CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN`: 4 rows
  - `MSY-ATL-CMH-BWI-RDU-BNA-DTW-MDW-LAX`: 3 rows
  - `ELP-DAL-LIT-ATL-RIC-MDW-MCI-PHX-SAN`: 2 rows
  - `ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA`: 1 row

Execution status was not identical: `claude/opus` and `claude/sonnet` finished cleanly, while `codex/gpt-5.4/run-001` is marked `partial` in `compare.json` because `presentation_render` failed even though `query.sql`, `result.json`, `report.md`, and `visual.html` artifacts were produced.

## Full SQL artifacts

### claude / opus

- `run-001`: [query.sql](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-20/q001_hops_per_day/claude/opus/run-001/query.sql), [result.json](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-20/q001_hops_per_day/claude/opus/run-001/result.json)

### claude / sonnet

- `run-001`: [query.sql](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-20/q001_hops_per_day/claude/sonnet/run-001/query.sql), [result.json](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-20/q001_hops_per_day/claude/sonnet/run-001/result.json)

### codex / gpt-5.4

- `run-001`: [query.sql](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-001/query.sql), [result.json](https://github.com/boristyshkevich/ExploringDatabyLLMs-runs/blob/main/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-001/result.json)

## Real output differences

There is no business-result difference across the three `result.json` files after contract normalization.

Verified equivalence:

- the same 10 `(aircraft, flight number, carrier, date)` keys appear in all three runs,
- `Hops` is `8` for every row in every run,
- normalized `Route` strings match exactly for all 10 rows,
- normalized departure-time sequences also match exactly for all 10 rows.

The observed differences are presentation/schema differences, not answer differences:

- `claude/opus` returns prompt-aligned columns with aliases such as `Tail_Number`, `Flight_Number_Reporting_Airline`, `IATA_CODE_Reporting_Airline`, and `DepTimes`.
- `claude/sonnet` uses `Carrier` and `DepartureTimes`, and formats routes with `→`.
- `codex/gpt-5.4` uses compare-contract-friendly aliases (`Aircraft ID`, `Flight Number`, `Date`) and adds four extra derived columns:
  - `Max Hops Overall`
  - `Max Hop Itinerary Count`
  - `Same-Hops Route Count`
  - `Route Frequency In Top 10`

Those extra `codex` fields do not change the underlying top-10 itinerary set.

## SQL comparison

The SQL strategies differ materially even though the outputs match.

- `claude/opus` uses a two-step pattern: build `legs`, compute top 10 groups in `counts`, join back in `top_legs`, then aggregate ordered legs into route and departure-time strings.
- `claude/sonnet` does the itinerary construction in one grouped pass with `groupArray` plus `arraySort`, then orders the grouped itineraries and limits to 10.
- `codex/gpt-5.4` builds timestamped legs, aggregates itineraries, computes `max_hops`, then adds window-based repetition metrics before selecting the top 10.

Other verified SQL-shape differences:

- `claude/opus` filters `Cancelled = 0` and `DepTime IS NOT NULL`.
- `claude/sonnet` filters `Cancelled = 0`, allows `coalesce(DepTime, CRSDepTime, 0)` inside the grouped sort key, and requires `length(sorted_legs) >= 2`.
- `codex/gpt-5.4` adds `Diverted = 0` and constructs explicit `DateTime` values for sorting.
- Route formatting differs by SQL output string:
  - `claude/opus`: hyphen-delimited airport chain
  - `claude/sonnet`: arrow-delimited chain
  - `codex/gpt-5.4`: `->` chain plus extra repetition metrics

## Presentation artifacts

### claude / opus

- `run-001`: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-20%2Fq001_hops_per_day%2Fclaude%2Fopus%2Frun-001%2Freport.md) gives a direct narrative answer to the report prompt: maximum hop count, the most recent itinerary, and route repetition counts.
- `run-001`: [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-20/q001_hops_per_day/claude/opus/run-001/visual.html) includes the required anchored KPI strip, selected-itinerary map, route-sequence panel, itinerary table, and query ledger. The HTML also contains degraded-map messaging for pending or failed airport enrichment.

### claude / sonnet

- `run-001`: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-20%2Fq001_hops_per_day%2Fclaude%2Fsonnet%2Frun-001%2Freport.md) also answers the report prompt directly, with the most recent itinerary called out explicitly and the repeated-route patterns summarized from the top 10.
- `run-001`: [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-20/q001_hops_per_day/claude/sonnet/run-001/visual.html) includes the same core required pieces as opus: anchored KPIs, selected-itinerary map, route-sequence/detail panel, clickable itinerary table, query ledger, and explicit degraded-state handling around enrichment.

### codex / gpt-5.4

- `run-001`: [report.md](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/md.html?file=2026-03-20%2Fq001_hops_per_day%2Fcodex%2Fgpt-5.4%2Frun-001%2Freport.md) is more templated: it includes the full table and points to the derived columns, but it does not turn those values into a fully written narrative the way the two Claude reports do.
- `run-001`: [visual.html](https://boristyshkevich.github.io/ExploringDatabyLLMs-runs/2026-03-20/q001_hops_per_day/codex/gpt-5.4/run-001/visual.html) is the richest HTML artifact by feature count. Verified elements include a lead-itinerary map, legend, route sequence/detail panel, itinerary table, query ledger, degraded-map messaging, and an additional carrier filter/reset control. Even so, the run is still `partial` because `compare.json` records `presentation_render: failed`.

## Execution stats

### claude / opus

- `run-001`: status `ok`; 10 rows; query time `4.37 s`; read rows `386,123,882`; memory `24.9 GiB`; fastest successful run and lowest memory use.

### claude / sonnet

- `run-001`: status `ok`; 10 rows; query time `7.36 s`; read rows `193,061,941`; memory `41.2 GiB`; lowest read volume.

### codex / gpt-5.4

- `run-001`: status `partial`; 10 rows; query time `19.50 s`; read rows `579,185,823`; memory `50.1 GiB`; slowest query and highest resource use of the three.

## Takeaway

For `q001` on `2026-03-20`, the models converged on the same answer set despite noticeably different SQL plans and output schemas. The main benchmark separation here is not correctness of the returned itineraries, but efficiency and presentation behavior: `claude/opus` was fastest and most memory-efficient, `claude/sonnet` read the fewest rows, and `codex/gpt-5.4` returned the same top 10 plus extra context columns but incurred the heaviest query cost and did not finish presentation rendering cleanly.
```