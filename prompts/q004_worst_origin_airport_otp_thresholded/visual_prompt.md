Build a dashboard that:

- uses `worst_airport` as the primary saved SQL already provided in the prompt
- uses `spread_to_middle` and `bottom_group_mix` as supporting queries when they materially improve the dashboard
- shows KPI cards for worst airport, worst OTP, highest average departure delay among ranked airports, and qualifying airport count
- renders a ranked horizontal bar or lollipop chart for the worst 25 airports by departure OTP
- renders a scatter plot of `CompletedDepartures` vs `DepartureOtpPct`
- renders a detail table with the full ranked result
- includes a narrative takeaway about whether the weakest airports cluster in large hubs or a more mixed set
- derives chart extents and highlighted airports from fetched data instead of hardcoding them
- uses one accent treatment for the worst 5 airports
- annotates the single worst airport in both charts
- keeps the scatter plot readable despite skewed volume differences
- shows supporting queries in the query ledger when used
