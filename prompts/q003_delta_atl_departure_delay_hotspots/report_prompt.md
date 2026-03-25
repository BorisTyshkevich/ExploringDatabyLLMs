### main
Find the Delta departure-delay hotspots out of ATL that appear to be persistently problematic, not just noisy one-off periods.

Analyze completed Delta departures from ATL by destination and departure time block across the full available history. Focus on combinations that have enough flight volume to be credible and enough repeated monthly presence to count as sustained hotspots.

You may apply reasonable minimum-volume filters to remove noise, but do not invent a custom hotspot score or a narrower analysis window.

For each hotspot, quantify:

- flight volume
- average departure delay
- a high-delay measure that captures the worse end of the distribution
- the share of flights departing 15+ minutes late
- how many months the hotspot meaningfully appears

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked hotspot summary view
- a monthly trend view for the leading hotspots

The monthly trend view should include only months that are credible for interpretation, not thin low-volume months that would add noise.

### q1
Which destination and time block is the worst recurring hotspot?

### q2
Is that hotspot consistently bad across time, or concentrated in a narrower period?

### q3
What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?
