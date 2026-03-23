Determine which airline and origin-airport combinations perform worst in winter after applying a meaningful flight threshold.

Focus on winter departures only and evaluate completed flights at the `(carrier, origin airport)` level. Limit the analysis to combinations with enough winter traffic to be credible.

For each qualifying pair, quantify:

- winter flight volume
- departure on-time performance
- average departure delay
- how reported delay minutes split across weather and operational causes such as carrier, NAS, security, and late aircraft

Rank the worst-performing winter pairs by on-time performance, while using the delay-cause mix as context rather than as the primary ranking driver.

Return one SQL query that produces a ranked view of the weakest qualifying winter carrier-airport pairs.

The output should let a BI dashboard answer:

- Which winter carrier-airport pair ranks worst overall?
- Are the worst pairs driven more by weather or by operational causes?
- Are the weakest pairs concentrated in a small number of carriers or airports?

Keep the result business-readable and analytically sound. Exclude low-volume winter pairs before ranking them.
