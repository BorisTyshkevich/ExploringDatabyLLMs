Created [visual.html](/Users/bvt/work/ExploringDatabyLLMs-runs/2026-03-23/q001_hops_per_day/codex/gpt-5.4/run-001/visual.html).

The page is a dynamic `html_map` dashboard wired to the saved primary SQL and a single cached airport-coordinate enrichment query against `ontime.dim_airports`, with lead KPIs anchored to the top-ranked itinerary, clickable itinerary rows, a persistent map card with degraded-state handling, a legend, route sequence/detail panels, export, and the required footer controls and query ledger.

Verification: I debug-ran the primary itinerary query and the airport enrichment query through ClickHouse MCP successfully, then syntax-checked the embedded browser script with Node. I did not open the artifact view frame.