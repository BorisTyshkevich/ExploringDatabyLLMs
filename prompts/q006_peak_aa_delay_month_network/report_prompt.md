Find American Airlines' worst network-wide month for departure delays, then identify which origins and routes contributed most to that peak.

Analyze completed American Airlines flights by month across the full network. Find the single month that stands out as the worst overall for departure delays.

Use the full available history unless the question explicitly asks for a narrower period. Do not invent a custom score for the peak month; identify it directly from the monthly delay metrics needed to answer the question.

For the monthly view, quantify:

- flight volume
- average departure delay
- the share of flights departing 15+ minutes late

Then drill into the selected peak month to show which origin airports and origin-destination routes contributed most to that bad month. Focus on contributors with enough flights in that month to be meaningful.

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a monthly leaderboard showing how the network performed over time
- a drilldown into the selected peak month by origin and by route

The proof query behind the peak-month question should preserve the month-by-month network rows needed for the dashboard, not just the single worst month.

The output should let a BI dashboard answer:

- Which month is the single worst American Airlines month for departure delays?
- Which origins contribute most to that peak month?
- Which routes contribute most to that peak month?
- Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

In the report, answer those questions directly in prose. Name the worst month, identify the leading origin and route contributors using the verified result, and summarize whether the peak looks broad or concentrated.

Do not use fallback phrases such as "the peak month" or "the leading contributors" when your verified query results let you name the actual month, origins, and routes directly.

Keep the result business-readable and analytically sound.
