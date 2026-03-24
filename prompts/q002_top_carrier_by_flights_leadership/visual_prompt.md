Build a dashboard that:

- uses `most_frequent_leader` as the primary saved SQL already provided in the prompt
- uses `leadership_swing`, `sharpest_transition`, and `stability_pattern` as supporting queries when they materially improve the dashboard
- shows KPI cards for total years analyzed, distinct annual leaders, largest leader share gap, and the sharpest leadership transition
- renders a bump chart for yearly carrier rank among the top carriers
- renders a time series of yearly completed-flight share for the leading carriers
- includes a compact table of all leadership-change years with prior leader, new leader, share swing, and share gap
- includes a narrative takeaway about whether the market shows long stable eras or frequent turnover
- highlights true leadership transitions only; do not treat the first year as a transition
- makes the sharpest leadership transition visually distinct
- derives all shown carriers from fetched data instead of hardcoding airline names
- shows supporting queries in the query ledger when used
