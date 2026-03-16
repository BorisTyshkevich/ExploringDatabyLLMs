Use the `ontime-analyst-dashboard` skill in dynamic mode.

Build a final `visual.html` page for this question that:

- shows KPI cards for worst winter pair, worst OTP, average delay of the worst pair, and total qualifying pairs
- renders a ranked chart for the worst winter `(Reporting_Airline, Origin)` pairs
- renders a stacked bar chart of delay-cause shares for the top 10 pairs
- renders a compact table of the full ranked result
- derives the top 10 pairs for the cause-share chart from fetched ranking data, not from hardcoded labels
- makes winter framing explicit in titles and copy
- clearly separates ranking severity from cause composition
- keeps carrier-airport labels readable without truncating meaning
