# Analysis Review
Verdict: WARN

## Summary

The main query correctly answers the core question: it groups by (Tail_Number, FlightNumber, Carrier, FlightDate), constructs the chronological Route string, deduplicates by Route, and returns the top 10 most-recent unique itineraries with all seven required columns. The results (10 rows, all 8-hop WN routes, verified from `results/main.json`) are plausible and well-formed. Dashboard subquestions are addressed with supporting queries. However, several issues reduce confidence in the supporting evidence: a redundant/misleading `num_flights` metric, grain inconsistency in `queries/q1.sql`, a hardcoded hop-count filter in `queries/q2.sql`, an unexplained occurrences discrepancy across subqueries, and one prose claim with no supporting query result.

## Findings

**1. `num_flights` is redundant with `hop_count` (main.sql, results/main.json)**

`num_flights` is computed as `argMax(hops, FlightDate)` — the hop count on the most recent date. `hop_count` is `max(hops)`. For every row in `results/main.json`, both equal 8. The column name implies route recurrence (how many times the itinerary was flown), but it just re-expresses the per-day segment count. The prompt's "number of flights" is ambiguous, but a same-valued duplicate column produces no additional information and could mislead readers.

**2. `queries/q1.sql` drops `Tail_Number` from the GROUP BY**

`main.sql` defines one itinerary as (Tail_Number, FlightNumber, Carrier, FlightDate). `queries/q1.sql` groups only by (FlightNumber, Carrier, FlightDate), omitting Tail_Number. If two aircraft fly the same flight number on the same day, the CTE merges their segments into one pseudo-itinerary, potentially producing a garbled Route string or inflated `occurrences` count. This explains why the same route `FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC` shows 114 occurrences in `results/q1.json` but only 99 in `results/q2.json` (which does include Tail_Number). The report cites the q1 figure ("114 times") without acknowledging the grain difference.

**3. `queries/q2.sql` hardcodes `HAVING hops = 8`**

The geographic analysis filters `HAVING hops = 8`, silently assuming the top routes are always exactly 8-hop. This happens to be true for the current top 10, but it is a fragile assumption that would silently exclude any future top-ranked route with a different hop count. The query should derive the maximum hop threshold from the data rather than hard-coding it.

**4. Unsupported prose claim about WN 3149 recurrence (`report.md`)**

The report states: "CLE-BNA-PNS-HOU-MCI-PHX-BUR-OAK-DEN (WN flight 3149) appears across multiple weeks in 2024." This route appears in `results/main.json` with a single flight_date of 2024-02-18. It does not appear in `results/q1.json` (not in the top 10 by occurrences) or `results/q2.json`. No executed query result supports the "multiple weeks in 2024" claim for this specific route.

**5. Occurrences discrepancy between q1 and q2 not explained (`report.md`)**

The same route `FLL-JAX-IND-MDW-MCI-DAL-ABQ-LAX-SJC` is shown with 114 occurrences in the first dashboard table and 99 in the second. These figures come from queries with different grains (issue 2 above). The report presents both numbers without comment, which is confusing.

## Suggested Prompt Fixes

**For the `num_flights` ambiguity:**

The prompt phrase "number of flights" is unclear. Add a definition to close the gap:

> Return "number of flights" as the count of individual segment-legs in the itinerary on that day (i.e., the same value as hop count), **or** as the total number of calendar days across all history on which this exact Route was flown — whichever is intended. State explicitly which meaning is required.

If recurrence is the intended meaning, replace with:
> - route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft

**For the hardcoded `HAVING hops = 8` in q2:**

Add a prompt instruction such as:
> For dashboard subquestions, do not hard-code any numeric threshold derived from the main result (e.g., hop count). Filter or join dynamically against the main result set instead.

**For unsupported claims:**

> All prose claims about a specific route's recurrence frequency or date range must be directly traceable to a row in an executed query result. Do not infer recurrence from a route appearing in the top-10 main result alone.
