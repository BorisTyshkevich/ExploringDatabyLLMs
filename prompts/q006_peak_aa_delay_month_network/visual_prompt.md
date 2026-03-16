Use the `ontime-analyst-dashboard` skill in dynamic mode.

Build a final `visual.html` page for this question that:

- shows KPI cards for peak month, peak average departure delay, peak `% DepDel15`, and completed flights in the peak month
- renders a monthly time-series chart of average `DepDelayMinutes` across all AA months
- visually highlights the peak month on that chart
- renders a bar chart of top origin contributors within the peak month
- renders a route contribution table for the peak month
- derives the peak month from fetched monthly data instead of hardcoding it
- clearly separates the network-wide trend from the peak-month drilldown
- annotates the peak month on the time series
- makes the contribution logic easy to read without external narrative
