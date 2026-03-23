Identify which origin airports have the worst departure on-time performance after excluding low-volume airports.

Analyze completed departures at the origin-airport level. Focus on airports with enough traffic to make the comparison meaningful, and rank the weakest performers by departure on-time performance.

For each airport, quantify:

- completed departures
- departure on-time performance
- average departure delay
- a high-delay measure that reflects the worse end of the delay distribution
- the first and last dates represented in the data

Return one SQL query that produces a ranked view of the weakest qualifying origin airports.

The output should let a BI dashboard answer:

- Which airport ranks worst on departure on-time performance?
- How large is the spread between the worst airport and the middle of the ranked set?
- Are the weakest airports mostly major hubs, or is the bottom group more mixed?

Keep the result business-readable and analytically sound. Exclude low-volume airports before ranking them.
