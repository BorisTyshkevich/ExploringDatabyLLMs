Determine which airline and origin-airport combinations perform worst in winter after applying a meaningful flight threshold.

Focus on winter departures only and evaluate completed flights at the `(carrier, origin airport)` level. Limit the analysis to combinations with enough winter traffic to be credible.

Use winter consistently as the business definition of the season for the full available history. You may apply a reasonable minimum-volume filter to remove noise, but do not invent a custom score or let delay-cause shares replace the primary performance ranking.

For each qualifying pair, quantify:

- winter flight volume
- departure on-time performance
- average departure delay
- how reported delay minutes split across weather and operational causes such as carrier, NAS, security, and late aircraft

Rank the worst-performing winter pairs by on-time performance, while using the delay-cause mix as context rather than as the primary ranking driver.

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked view of the weakest qualifying winter carrier-airport pairs
- a cause-composition view for the leading weak pairs that separates weather from operational causes

The proof query behind the worst-pair question should preserve the ranked pair-level rows needed for the dashboard, not just a single worst pair or summary count.

## Dashboard Questions

- Which winter carrier-airport pair ranks worst overall?
- Are the worst pairs driven more by weather or by operational causes?
- Are the weakest pairs concentrated in a small number of carriers or airports?

In the report, answer those questions directly in prose. Name the worst winter pair, summarize whether the weakest pairs are driven more by weather or by operational causes, and state whether the weak set is concentrated in a small number of carriers or airports.

Do not use fallback phrases such as "the worst pair" or "the weakest pairs" when your verified query results let you name the actual carrier-airport combinations directly.

Keep the result business-readable and analytically sound. Exclude low-volume winter pairs before ranking them.
