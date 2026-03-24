Build a dashboard that:

- uses `peak_month` as the primary saved SQL already provided in the prompt
- uses `origin_contributors`, `route_contributors`, and `concentration_pattern` as supporting queries when they materially improve the dashboard
- shows KPI cards for peak month, peak average departure delay, peak `% DepDel15`, and completed flights in the peak month
- renders a monthly time-series chart of average `DepDelayMinutes` across all AA months
- visually highlights the peak month on that chart
- renders a bar chart of top origin contributors within the peak month
- renders a route contribution table for the peak month
- includes a narrative takeaway about whether the peak month was broad across the network or concentrated in a smaller set of origins and routes
- derives the peak month from fetched monthly data instead of hardcoding it
- clearly separates the network-wide trend from the peak-month drilldown
- annotates the peak month on the time series
- makes the contribution logic easy to read without external narrative
- shows supporting queries in the query ledger when used
