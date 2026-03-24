Build a dynamic dashboard for Delta ATL departure delay hotspots.

Layout intent:

- headline and subtitle that clearly frame Delta departures from ATL
- narrative hero naming the single worst hotspot with its key metrics
- KPI strip derived from the primary fetched query
- hotspot ranking or heatmap from the primary query
- persistence chart or time-series panel using the `persistence` supporting query
- supporting pattern panel or table using the `pattern_summary` supporting query
- visible query ledger and footer controls following the dashboard skill contract
- concluding takeaway section that summarizes the broader operational pattern across the leading hotspots

Behavior:

- use `worst_hotspot` as the primary saved SQL already provided in the prompt
- use `persistence` and `pattern_summary` as supporting queries when they materially improve the dashboard
- if supporting queries are used, show them in the query ledger with their own status and SQL text
- the dashboard does not need to mirror `report.md`; use the narrative answers as framing and the fetched query results for charts and tables
- make the single worst hotspot visually prominent
- keep the page readable on mobile with stacked sections and horizontally scrollable tables or charts when needed
- preserve useful content when any supporting query fails or returns zero rows
