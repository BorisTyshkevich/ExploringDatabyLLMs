Build a dynamic non-map dashboard for Delta ATL departure delay hotspots.

Layout intent:

- headline and subtitle that clearly frame Delta departures from ATL
- KPI strip for worst hotspot, worst average departure delay, p90 departure delay, and qualifying months
- primary heatmap with `Dest` on one axis and `DepTimeBlk` on the other, colored by average `DepDelayMinutes`
- supporting monthly trend view for the top 3 hotspot cells
- ranked table for the top 20 hotspot cells
- visible legend for the heatmap color scale

Behavior:

- derive the top 3 hotspot cells for the trend view from the fetched ranking rows, not from hardcoded labels
- make the single worst hotspot cell visually prominent
- keep the heatmap readable on mobile by allowing horizontal scrolling if necessary
- preserve useful content when the fetched result set is empty
