Build a dashboard that:

- uses `worst_pair` as the primary saved SQL already provided in the prompt
- uses `cause_mix` and `concentration_pattern` as supporting queries when they materially improve the dashboard
- shows KPI cards for worst winter pair, worst OTP, average delay of the worst pair, and total qualifying pairs
- renders a ranked chart for the worst winter `(Reporting_Airline, OriginCode)` pairs
- renders a stacked bar chart of delay-cause shares for the top 10 pairs
- renders a compact table of the full ranked result
- includes a narrative takeaway about whether the weakest winter pairs cluster in a small set of carriers or origins
- derives the top 10 pairs for the cause-share chart from fetched ranking data, not from hardcoded labels
- makes winter framing explicit in titles and copy
- clearly separates ranking severity from cause composition
- keeps carrier-airport labels readable without truncating meaning
- shows supporting queries in the query ledger when used
