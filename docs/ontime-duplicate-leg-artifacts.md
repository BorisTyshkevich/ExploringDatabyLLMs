# Duplicate-Leg Artifacts in the OnTime Dataset

This note captures a useful data-quality finding from exploring the OnTime dataset for [`q001_hops_per_day`](/Users/bvt/work/ExploringDatabyLLMs/prompts/q001_hops_per_day).

The short version:

- exact duplicate flight-leg rows do exist in `ontime.fact_ontime`
- they are rare at whole-dataset scale
- they matter a lot for itinerary reconstruction queries
- missing `Tail_Number` makes the artifact much worse, because multiple same-day rows can collapse into one fake aircraft itinerary
- exact-row duplicates can be removed during loading, but broader leg-key collisions require a modeling decision

## Why this came up

The `q001` question asks for the highest number of same-day hops flown by one aircraft under one flight number.

An early query shape used:

- `count()` as `hop_count`
- `groupArray(...)` to collect legs
- sorting by `CRSDepTime`
- grouping by `(FlightDate, carrier, flight_number, Tail_Number)`

That query produced apparently extreme itineraries such as:

- `AA 1086` on `1988-11-12` with `9` hops
- `OO 6195` on `2003-02-28` with `10` hops

After checking raw rows, both turned out to be artifacts.

## Dataset-wide duplicate summary

Using this leg identity:

- `FlightDate`
- `IATA_CODE_Reporting_Airline`
- `Flight_Number_Reporting_Airline`
- `Tail_Number`
- `CRSDepTime`
- `OriginCode`
- `DestCode`

the whole table shows:

- total rows: `230,307,587`
- distinct grouped legs on that key: `230,306,507`
- duplicate leg groups: `1,073`
- excess duplicate rows: `1,080`
- worst single duplicate count for one leg: `5`

So exact duplicates are extremely rare globally, but not zero.

There are really two different duplicate problems here:

- exact duplicate rows, where every stored column matches
- duplicate leg-key groups, where rows share the same itinerary leg identity but differ in other operational columns

That distinction matters because the first kind can be removed safely during loading, while the second kind cannot be removed safely without choosing a survivor rule.

Most duplicate excess is concentrated in older years:

- `1989`: `438` excess rows
- `1988`: `404`
- `1987`: `135`
- `2003`: `34`
- later years: only scattered cases

## Exact-row duplicates vs leg-key collisions

The broad leg identity used above is analytically useful, but it is not a table primary key.

Two rows can match on:

- `FlightDate`
- `IATA_CODE_Reporting_Airline`
- `Flight_Number_Reporting_Airline`
- `Tail_Number`
- `CRSDepTime`
- `OriginCode`
- `DestCode`

and still differ in other columns such as:

- `DepTime`
- `ArrTime`
- `DepDelay`
- `ArrDelay`
- cancellation flags
- elapsed time fields

That means:

- exact-row deduplication is safe
- leg-key deduplication is policy-driven

In practice, some remaining leg-key collisions look like raw source duplication, but others may represent conflicting operational records or ambiguous same-day movements, especially when `Tail_Number = ''`.

## What the loader now does

The loader for `ontime.stage_ontime` now runs:

```sql
OPTIMIZE TABLE ontime.stage_ontime PARTITION <year> FINAL DEDUPLICATE
```

before publishing the year into `ontime.fact_ontime`.

That gives the dataset a narrower and safer contract:

- exact duplicate rows are removed during loading
- the published fact table should no longer retain byte-identical duplicate rows for a rebuilt year
- non-identical rows that collide on the leg key are preserved

This is deliberate. The goal is to remove clearly redundant source rows without silently deleting rows that may encode different operational facts.

## Why the remaining leg-key collisions are harder

After exact-row deduplication, some duplicate leg-key groups still remain.

For example, two rows may describe the same scheduled leg but disagree on actual departure or arrival metrics. In that situation, there is no universally correct “drop one” rule built into the raw dataset contract.

So there are two defensible layers:

- raw fact storage: keep non-identical leg-key collisions
- analytical itinerary layer: choose one row per leg key with an explicit winner rule when the question requires one-aircraft route reconstruction

This is especially important for blank-tail records:

- when `Tail_Number` is present, grouping often still corresponds to one physical aircraft
- when `Tail_Number = ''`, the grouping can mix true duplicates, conflicting records, and different aircraft flying the same day/flight number pattern

For that reason, the recommended place to “beat” the remaining duplicates is not the raw loader. It is an analytical view or query shape that states its survivor policy explicitly.

## Recommended handling for analysis

For exploratory or benchmark queries such as `q001`, use a deduplicated leg representation rather than trusting raw row counts.

Reasonable approaches are:

- count distinct legs rather than raw rows
- build deduplicated leg arrays before constructing routes
- create a dedicated analytical view that picks one row per leg key using an explicit rule, such as preferring:
  - non-empty `Tail_Number`
  - `Cancelled = 0`
  - rows with more populated operational timestamps and delay fields

This keeps `ontime.fact_ontime` close to the source while still giving itinerary-style analysis a stable contract.

## Example 1: AA 1086 on 1988-11-12

Raw grouped rows for:

- `FlightDate = 1988-11-12`
- `carrier = AA`
- `flight_number = 1086`
- `Tail_Number = ''`

show:

- total rows: `9`
- distinct legs by `(CRSDepTime, OriginCode, DestCode)`: `8`
- duplicate rows: `1`

The duplicated leg is:

- `2150 DFW -> OKC`
- appears `2` times

The raw-row itinerary query therefore produced:

- route: `TPA-BNA-OMA-DSM-BNA-RSW-SRQ-DFW-OKC-OKC`
- raw hop count: `9`

But this is not a genuine 9-leg aircraft rotation. It is an 8-leg itinerary with one duplicated source row at the final leg.

Important distinction:

- `9` airports in a route string does not mean `9` hops
- `hop_count` is the number of legs
- `airport_count = hop_count + 1`

## Example 2: OO 6195 on 2003-02-28

This case is more severe.

Raw grouped rows for:

- `FlightDate = 2003-02-28`
- `carrier = OO`
- `flight_number = 6195`
- `Tail_Number = ''`

show:

- total rows: `10`
- distinct legs: `3`
- excess duplicate rows: `7`

Grouped raw legs:

- `2030 LAX -> IPL` appears `2` times
- `2030 LAX -> SAN` appears `3` times
- `2135 SAN -> IPL` appears `5` times

The raw-row itinerary query produced:

- route: `LAX-IPL-IPL-SAN-SAN-SAN-IPL-IPL-IPL-IPL-IPL`
- raw hop count: `10`

This is not a plausible one-aircraft itinerary. It is a merged artifact built from repeated rows and conflicting same-time movements.

## Why empty `Tail_Number` makes this worse

Duplicate rows alone are not the whole story.

The early query grouped by:

- `FlightDate`
- `carrier`
- `flight_number`
- `Tail_Number`

When `Tail_Number` is present, that usually keeps physical aircraft separate.

When `Tail_Number = ''`, every row for the same day and flight number collapses into the same bucket. That means:

- duplicate legs can inflate hop count
- conflicting same-time departures can get merged together
- separate physical aircraft can be mistaken for one itinerary

This is why the blank-tail cases are especially dangerous for itinerary reconstruction.

## Query-design lesson

For `q001`, the unsafe pattern is:

```sql
count() AS hop_count
```

plus:

```sql
groupArray((CRSDepTime, OriginCode, DestCode))
```

That pattern treats raw rows as trustworthy legs.

The safer pattern is:

- count distinct legs, not raw rows
- use deduplicated leg arrays when constructing the route
- treat repeated adjacent airports like `...-OKC-OKC` or `IPL-IPL-...` as artifact signals
- be especially skeptical when `Tail_Number` is empty

In practice, the stable q001 result comes from deduplicating legs rather than trusting raw row count.

## Practical implication for exploratory analysis

This is a good example of why “interesting top records” in operational datasets need row-level validation.

The artifact routes were not random hallucinations by the model. They were reachable from real source rows using a plausible but naive query pattern.

That makes them analytically interesting:

- the dataset is mostly clean at scale
- older years contain a small but real duplicate-row seam
- missing aircraft ids can turn a rare seam into a dramatic false top result
- exact-row cleanup improves the raw table, but itinerary correctness still depends on leg-level query design

For a blog post, this is a strong story because it shows the value of iterative analytical debugging:

1. ask an eye-catching question
2. get a surprising answer
3. inspect the raw rows
4. discover a data-quality seam
5. refine the query contract so the benchmark measures the intended phenomenon rather than a source artifact

## Related Altinity posts

These Altinity blog posts use the OnTime dataset in ways that are relevant background for a longer article:

- [ClickHouse Aggregation Fun, Part 1: Internals and Handy Tools](https://altinity.com/blog/clickhouse-aggregation-fun-part-1-internals-and-handy-tools)
  - uses the US airline OnTime dataset to explain aggregation behavior and performance in ClickHouse
- [Creating Beautiful Grafana Dashboards on ClickHouse: a Tutorial](https://altinity.com/blog/2019-12-28-creating-beautiful-grafana-dashboards-on-clickhouse-a-tutorial)
  - builds visual analysis over the OnTime dataset, including airline and airport views
- [Using the Altinity Grafana Plugin for ClickHouse in Grafana Cloud](https://altinity.com/blog/using-altinitys-grafana-clickhouse-plugin-in-grafana-cloud)
  - shows a richer dashboard workflow on OnTime data, including map-style airport analysis
- [MySQL to ClickHouse Data Migration and Replication](https://altinity.com/blog/2017-12-1-mysql-to-clickhouse-data-migration-and-replication-airlineontime-dataset)
  - uses a subset of airline OnTime data as a migration and replication example
- [What’s Up with Parquet Performance in ClickHouse?](https://altinity.com/blog/whats-up-with-parquet-performance-in-clickhouse)
  - includes OnTime-based benchmarks and is useful for framing scale/performance context around the dataset
