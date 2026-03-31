# Analysis Review
Verdict: FAIL

## Summary
The run is unreliable for the asked question. The core issue is scope drift in [queries/main.sql](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/queries/main.sql): it hard-codes `FlightDate >= addYears(today(), -5)`, so [results/main.json](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/results/main.json) only reports 8-hop routes from the last five years instead of answering the all-history question. Live MCP validation against `ontime.fact_ontime` shows the all-history maximum is 9 hops, not 8, so the `main` answer and all downstream sections are based on the wrong itinerary set.

The saved SQL also does not enforce leg-to-leg continuity when building `Route`; it only sorts legs by departure/arrival time and concatenates origins/destinations. That conflicts with the prompt’s requirement to avoid artifact routes. Live validation shows the all-history 9-hop candidate produced by this logic is `DFW-OKC-BNA-OMA-DSM-BNA-RSW-SRQ-DFW-OKC`, which is not a contiguous chain because the first leg ends at `OKC` but the next leg starts at `TPA`. That makes the method itself unsafe for this question, especially for empty `Tail_Number` cases the prompt explicitly called out.

## Findings
`1.` [queries/main.sql](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/queries/main.sql), [queries/q1.sql](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/queries/q1.sql), [queries/q2.sql](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/queries/q2.sql), and [queries/q3.sql](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/queries/q3.sql) all impose `FlightDate >= addYears(today(), -5)` even though the prompt did not request a time window. As a result, [report.md](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/report.md) claims “the highest observed ... itineraries reached 8 hops,” but live MCP validation shows the all-history maximum is 9 hops. This is a substantive correctness failure, not a wording issue.

`2.` [results/main.json](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/results/main.json) returns 10 rows of 8-hop routes, and [results/q1.json](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/results/q1.json), [results/q2.json](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/results/q2.json), and [results/q3.json](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/results/q3.json) all analyze that same 10-route set. Because the underlying `main` population is wrong, every downstream section is answering the wrong question population as well.

`3.` [queries/main.sql](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/queries/main.sql) constructs `Route` by sorting legs on `(dep_hhmm, arr_hhmm)` and concatenating the first leg’s origin plus all destinations, but never checks that each leg’s destination matches the next leg’s origin. The prompt explicitly warned against artifact routes. Live MCP inspection of the all-history 9-hop candidate shows the generated route is non-contiguous, so this SQL method is not robust for empty/unknown aircraft-id cases.

`4.` The prose in [report.md](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-002/report.md) is internally supported by the saved 5-year result files for q1-q3, but that does not rescue the run: those sections inherit the wrong itinerary universe from `main.sql`. The defect is therefore material and blocking.

## Suggested Prompt Fixes
Add an explicit scope guard to prevent invented date windows:

```md
Do not apply any date filter unless the prompt explicitly asks for one. Use the full available history.
```

Add an explicit continuity requirement so time-sorted concatenation is not treated as sufficient:

```md
Only keep itineraries whose ordered legs form a contiguous chain, where each leg's destination equals the next leg's origin. Discard any non-contiguous sequence as an artifact, including cases caused by empty or ambiguous aircraft identifiers.
```

Add an explicit instruction for downstream sections to reuse the exact verified `main` population:

```md
For q1-q3, analyze exactly the itinerary set returned by `main` after all filtering and validation; do not recompute a broader or narrower candidate set.
```
