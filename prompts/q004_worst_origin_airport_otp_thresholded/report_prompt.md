### main
Identify which origin airports have the worst departure on-time performance after excluding low-volume airports.

Analyze completed departures at the origin-airport level across the full available history. Focus on airports with enough traffic to make the comparison meaningful, and rank the weakest performers by departure on-time performance.

You may apply a reasonable minimum-volume filter to remove noise, but do not invent a custom score or redefine on-time performance.

For each airport, quantify:

- completed departures
- departure on-time performance
- average departure delay
- a high-delay measure that reflects the worse end of the delay distribution
- the first and last dates represented in the data

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked view of the weakest qualifying origin airports
- a comparison view that makes the gap between the very worst airports and the middle of the ranked set easy to judge

The proof query behind the worst-airport question should preserve the ranked airport-level rows needed for the dashboard, not just a single top airport or summary statistic.

### q1
Which airport ranks worst on departure on-time performance?

### q2
How large is the spread between the worst airport and the middle of the ranked set?

### q3
Are the weakest airports mostly major hubs, or is the bottom group more mixed?
