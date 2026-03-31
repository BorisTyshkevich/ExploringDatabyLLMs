# Highest Daily Hops for One Aircraft on One Flight Number

{{data_overview_md}}

## Maximum Hop Count

The highest recorded daily hop count is **{{metric.max_hops}} hops** ({{metric.max_hops_leg_count}} legs in a single day), achieved by carrier **{{metric.top_carrier}}** on flight **{{metric.top_flight_number}}** operated by tail number **{{metric.top_tail_number}}** on **{{metric.top_date}}**.

The full route for this itinerary:

> {{metric.top_route}}

Departure times per leg: {{metric.top_departure_times}}

## Operating Pattern

{{metric.max_hops}} hops represents the ceiling across the entire dataset. The top-10 itineraries show that this is a **recurring scheduled pattern**, not a one-off event — the same flight number and identical route sequence appears on multiple dates with different tail numbers, indicating a fixed published multi-hop rotation rather than an irregular assignment.

## Top 10 Longest Itineraries

{{result_table_md}}

## Route Repetition and Clustering

All 10 results belong to a single carrier ({{metric.dominant_carrier}}). Several flight numbers recur across multiple dates with near-identical routing:

- **{{metric.repeated_flight_1}}** follows the same {{metric.max_hops}}-hop sequence on {{metric.repeated_flight_1_count}} separate dates.
- **{{metric.repeated_flight_2}}** and **{{metric.repeated_flight_3}}** also appear multiple times with consistent routing.

This clustering confirms that {{metric.max_hops}}-hop daily rotations are a deliberate scheduling strategy for point-to-point carriers that chain aircraft through multiple short segments in a single day.

---
*Generated: {{generated_at}}*