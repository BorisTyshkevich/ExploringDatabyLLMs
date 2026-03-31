# Analysis Review
Verdict: FAIL

## Summary
The run correctly identifies July 2024 as the worst American Airlines month for departure delays, and the saved monthly evidence in `queries/q1.sql` and `results/q1.json` supports that answer. The failure is in the contributor and concentration metrics: the drilldown and network-share queries use `sum(DepDelay)` and describe it as delay minutes, which lets early departures subtract from the totals.

Because `report.md` uses those values to name the leading contributors and to justify the breadth-versus-concentration conclusion, the evidence support is materially wrong. The saved `review.md` verdict of PASS was not defensible.

## Findings
- `queries/q2.sql`, `queries/q3.sql`, and `queries/q4.sql` use `sum(DepDelay)` for totals that are described in `results/q2.json`, `results/q3.json`, `results/q4.json`, and `report.md` as departure delay minutes. That metric is not the same as delay minutes because negative `DepDelay` values from early departures reduce the total. For July 2024 AA flights, the full-network total is `3,127,538` by `sum(DepDelay)` but `3,297,804` by `sum(DepDelayMinutes)`, so the run understates delay totals by about 170k minutes.
- The origin-contributor prose in `report.md` is based on the wrong total-delay metric. `results/q2.json` reports DFW at `576,310` and CLT at `549,425`, but the corresponding delay-minute totals are `593,021` and `560,409`. The ranking still points to DFW and CLT, but the report claims specific delay-minute magnitudes that are not supported by the correct delay-minute measure.
- The route-contributor section has the same metric problem. `results/q3.json` reports DFW-LAX at `18,114`, CLT-MCO at `18,070`, and DFW-SAT at `17,701`, while the actual delay-minute totals are `18,462`, `18,376`, and `18,032`. `report.md` presents these as route delay minutes even though the SQL is summing raw signed departure delay.
- The breadth/concentration conclusion in `report.md` is also off because `queries/q4.sql` uses `sum(DepDelay)` in both the numerator and denominator. The reported shares in `results/q4.json` are `51.93%` for the top 5 origins and `5.02%` for the top 10 routes, but with true delay minutes the comparable shares are about `50.75%` and `4.88%`. The shape of the conclusion is similar, but the stated network-wide percentages are not supported by the correct metric.
- `queries/q1.sql` and `results/q1.json` do materially answer the month-selection question. The query preserves all 458 monthly rows, and July 2024 is correctly the highest month by average departure delay with `86,083` flights, `36.33` average departure delay minutes, and `0.3804` of flights departing 15+ minutes late.

## Suggested Prompt Fixes
- Add a metric guardrail to the contributor sections, for example: `When quantifying total delay minutes for origins, routes, or network concentration, use positive delay-minute fields such as DepDelayMinutes. Do not use signed raw DepDelay totals as a proxy for delay minutes.`
- Add a wording check for reported totals, for example: `If you describe a quantity as delay minutes in the report, the proof query must compute delay minutes directly rather than a signed net-delay measure that can be reduced by early departures.`
- Add a concentration-specific reminder, for example: `For the breadth/concentration question, compute network-wide shares from full-network delay-minute totals, then compare top origins and routes against those totals.`
