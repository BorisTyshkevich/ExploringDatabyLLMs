# Highest Daily Hops for One Aircraft on One Flight Number

{{data_overview_md}}

## Maximum Hop Count

The highest recorded hop count is **{{metric.max_hops}}** legs flown by a single aircraft under a single flight number in one day. All top-10 entries belong to {{metric.dominant_carrier}}. This is not a one-off anomaly: the same 8-hop routing patterns recur across multiple weeks on the same flight numbers, indicating deliberate scheduled network itineraries.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Date | {{metric.most_recent_date}} |
| Carrier | {{metric.most_recent_carrier}} |
| Flight Number | {{metric.most_recent_flight_num}} |
| Aircraft (Tail) | {{metric.most_recent_aircraft}} |
| Route | {{metric.most_recent_route}} |

## Top 10 Longest and Most Recent Itineraries

{{result_table_md}}

*DepartureTimes column shows `AIRPORT@HHMM` for each origin leg in chronological order.*

## Route Repetition and Clustering

Three distinct 8-hop route patterns recur across different aircraft on the same flight number, confirming them as recurring scheduled itineraries. WN 3149 (CLE → DEN corridor) is the most frequently repeated, appearing on consecutive Sundays in early 2024 with consistent city pairs from the Midwest through the South and out to the West Coast. WN 2787 (MSY → LAX corridor) and WN 154 (ELP → SAN corridor) each recur across separate weeks as well. This tight clustering by flight number and consistent routing is characteristic of deliberate point-to-point network scheduling.