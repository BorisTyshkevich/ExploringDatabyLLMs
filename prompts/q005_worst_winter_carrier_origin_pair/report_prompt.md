### main
Determine which airline and origin-airport combinations perform worst in winter after applying a meaningful flight threshold.

Focus on winter departures only and evaluate completed flights at the `(carrier, origin airport)` level. Limit the analysis to combinations with enough winter traffic to be credible.

Use winter consistently as the business definition of the season across the analyzed window. Default to the most recent 5 years unless the question explicitly asks for a different window. You may apply a reasonable minimum-volume filter to remove noise, but do not invent a custom score or let delay-cause shares replace the primary performance ranking.

For each qualifying pair, quantify:

- winter flight volume
- departure on-time performance
- average departure delay
- how reported delay minutes split across weather and operational causes such as carrier, NAS, security, and late aircraft

Rank the worst-performing winter pairs by on-time performance, while using the delay-cause mix as context rather than as the primary ranking driver.

If some leading weak pairs lack delay-cause reporting, say so directly and limit the weather-vs-operational conclusion to the subset with measured cause data.

### q1
Which winter carrier-airport pair ranks worst overall?

Preserve the ranked pair-level rows needed for the dashboard, not just a single worst pair or summary count.

### q2
Are the worst pairs driven more by weather or by operational causes?

Return the cause-composition rows needed to separate weather from operational causes for the leading weak pairs.

### q3
Are the weakest pairs concentrated in a small number of carriers or airports?

Make carrier concentration and airport concentration separately inspectable from the returned result, rather than requiring the reader to infer one of them indirectly from lists or arrays.
