Find the Delta departure-delay hotspots out of ATL that appear to be persistently problematic, not just noisy one-off periods.

Analyze completed Delta departures from ATL by destination and departure time block. Focus on combinations that have enough flight volume to be credible and enough repeated monthly presence to count as sustained hotspots.

For each hotspot, quantify:

- flight volume
- average departure delay
- a high-delay measure that captures the worse end of the distribution
- the share of flights departing 15+ minutes late
- how many months the hotspot meaningfully appears

Return one SQL query that supports two views from the same result:

- a ranked hotspot summary
- a monthly trend view for the leading hotspots

The output should let a BI dashboard answer:

- Which destination and time block is the worst recurring hotspot?
- Is that hotspot consistently bad across time, or concentrated in a narrower period?
- What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?

Keep the result business-readable and analytically sound. Exclude low-volume noise before identifying the leading hotspots.
