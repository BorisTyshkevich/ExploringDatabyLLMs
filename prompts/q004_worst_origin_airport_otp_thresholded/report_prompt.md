### main
Identify which origin airports have the worst departure on-time performance after excluding low-volume airports.

Analyze completed departures at the origin-airport level. Focus on airports with enough traffic to make the comparison meaningful, and rank the weakest performers by departure on-time performance.

You may apply a reasonable minimum-volume filter to remove noise, but do not invent a custom score or redefine on-time performance.

For each airport, quantify:

- completed departures
- departure on-time performance
- average departure delay across all completed departures
- a high-delay measure that reflects the worse end of the delay distribution for delayed departures
- the first and last dates represented in the data

The ranking must be one row per origin airport. Do not split the same airport into multiple ranked rows because of airport-name changes, city labels, or dimension-table history. If descriptive fields vary over time, still aggregate at the airport level and choose one representative display name for the output row.

### q1
Which airport ranks worst on departure on-time performance?

Preserve the ranked airport-level rows needed for the dashboard, not just a single top airport or summary statistic.

### q2
How large is the spread between the worst airport and the middle of the ranked set?

Return enough ranked rows to make the gap between the very worst airports and the middle of the qualifying set easy to judge.

### q3
Are the weakest airports mostly major hubs, or is the bottom group more mixed?

In the report, answer that question using only evidence returned by the proof query. If you make an interpretation such as "mostly major hubs," make it explicit that it is an inference from the returned airport set and traffic levels, not a directly queried classification.
