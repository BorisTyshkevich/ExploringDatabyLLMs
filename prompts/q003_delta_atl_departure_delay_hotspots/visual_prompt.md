Build a browser-ready analytical brief for Delta ATL departure delay hotspots.

Layout intent:

- headline and subtitle that clearly frame Delta departures from ATL
- hero section naming the single worst hotspot with its key metrics
- compact KPI strip for worst average delay, p90 departure delay, late-15-plus share, and qualifying months
- one section for each required business question with clear prose and an evidence card
- an evidence card should show the proof-query row count, column names, and the first preview row in a readable compact table
- a concluding takeaway section that summarizes the broader operational pattern across the leading hotspots

Behavior:

- make the single worst hotspot visually prominent
- rely only on the verified analysis package; do not imply broader coverage than the provided query previews support
- do not invent synthetic trend lines, inferred heatmap cells, or approximate values beyond the embedded proof-query previews
- if the provided evidence is too thin for a richer chart, prefer a textual evidence panel over a fabricated visualization
- keep the page readable on mobile with stacked sections and horizontally scrollable evidence tables if needed
- preserve useful content when any proof query returns zero rows
