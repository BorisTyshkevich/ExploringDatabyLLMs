### main
Find American Airlines' worst network-wide month for departure delays, then identify which origins and routes contributed most to that peak.

Analyze completed American Airlines flights by month across the full network. Find the single month that stands out as the worst overall for departure delays.

Do not invent a custom score for the peak month; identify it directly from the monthly delay metrics needed to answer the question.

For the monthly view, quantify:

- flight volume
- average departure delay across all completed flights
- the share of flights departing 15+ minutes late

Then drill into the selected peak month to show which origin airports and origin-destination routes contributed most to that bad month. Focus on contributors with enough flights in that month to be meaningful.

For the breadth-versus-concentration question, compute concentration against the full selected-month AA network. If you apply minimum-flight thresholds to keep origin or route tables business-meaningful, use those thresholds only for displayed contributor tables, not as the denominator for any network-wide concentration share.

If you use thresholded contributor tables for readability, label them as display filters only and keep the network-wide breadth/concentration conclusion tied to full-network denominators.

Contributor tables must be one row per origin for the origin table and one row per route for the route table. Do not split the same origin or route into multiple rows because of name changes, city labels, or dimension-table history. If descriptive fields vary over time, still aggregate at the entity level and choose one representative display name for the output row.

### q1
Which month is the single worst American Airlines month for departure delays?

Preserve the month-by-month network rows needed for the dashboard, not just the single worst month.

### q2
Which origins contribute most to that peak month?

### q3
Which routes contribute most to that peak month?

### q4
Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

Let the reviewer inspect full-network peak-month totals and the share captured by top origins and top routes against those full-network totals. Do not answer "across the network" using shares computed only from filtered subsets.
