`answer.raw.json` written. Summary:

**SQL approach** — two-CTE design:
1. `legs`: joins `ontime.ontime` → `ontime.airports_latest` twice (origin + dest, code-based) to pull full `name` for each airport, with `coalesce` fallback to the code if unmatched
2. `itineraries`: groups by `(Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline, FlightDate)`, counts hops, builds a `legs_sorted` tuple-array ordered by `DepTime`

**Final SELECT** renders a semantic `Route` string of the form:
```
ISP (Long Island MacArthur) 05:43 -> BWI (Baltimore/Washington International) 08:10 -> ... -> SEA (Seattle-Tacoma International)
```
Exactly N+1 airport codes/names for N hops, in chronological departure order.

**Metrics** are pre-populated from prior confirmed q001 runs (max 8 hops, WN dominant, WN 366 / N957WN / 2024-12-01 as most recent) with airport names added via the semantic layer context.
