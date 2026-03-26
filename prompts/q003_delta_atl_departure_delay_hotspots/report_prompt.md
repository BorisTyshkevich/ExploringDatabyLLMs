### main
Find the Delta departure-delay hotspots out of ATL that appear to be persistently problematic, not just noisy one-off periods.

Analyze completed Delta departures from ATL by destination and departure time block across the most recent 5 years unless the question explicitly asks for a different window. Focus on combinations that have enough flight volume to be credible and enough repeated monthly presence to count as sustained hotspots.

You may apply reasonable minimum-volume filters to remove noise, but do not invent a custom hotspot score or a different analysis window.

For each hotspot, quantify:

- flight volume
- average departure delay
- a high-delay measure that captures the worse end of the distribution
- the share of flights departing 15+ minutes late
- how many months the hotspot meaningfully appears

### q1
Which destination and time block is the worst recurring hotspot?

### q2
Is that hotspot consistently bad across time, or concentrated in a narrower period?

Preserve the monthly rows needed for that trend view, and include only months that are credible for interpretation rather than thin low-volume months that would add noise.

### q3
What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?
