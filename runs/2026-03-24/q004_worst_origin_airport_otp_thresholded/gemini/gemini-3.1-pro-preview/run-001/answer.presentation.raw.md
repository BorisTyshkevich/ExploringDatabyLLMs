I have successfully generated the interactive dashboard and saved it as `visual.html` in the current project directory. 

The dashboard runs in **dynamic** mode, directly connecting to the configured ClickHouse endpoint using your token. It includes:
- **KPI Strip** summarizing the worst airport, its OTP, the highest average delay, and the total qualifying airports.
- **Narrative Insight** derived dynamically from secondary queries, confirming that the weakest performers are heavily trafficked major hubs.
- **Visualizations**: 
  - A ranked horizontal bar chart of the worst 25 airports with a specific highlight on the worst 5 performers and a pin annotating the absolute worst.
  - A scatter plot plotting `CompletedDepartures` (on a logarithmic scale to account for skewed volumes) vs `DepartureOtpPct`, highlighting the worst airport among its peers.
- **Detail Table** providing the full dataset for ranking validation.
- **Query Ledger** displaying the primary analysis SQL and the secondary enrichment queries run in the background.
