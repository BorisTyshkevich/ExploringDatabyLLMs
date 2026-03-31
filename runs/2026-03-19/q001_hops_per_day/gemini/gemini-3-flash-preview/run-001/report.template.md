# Highest daily hops for one aircraft on one flight number

{{data_overview_md}}

The highest number of hops recorded for a single aircraft using a consistent flight number in one day is **{{metric.primary_value}}**. This pattern is observed across multiple dates, notably with Southwest Airlines (WN) flight 3149, which frequently operates an 8-hop route such as `CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN`.

## Longest Recent Itinerary
The most recent instance of an 8-hop itinerary was on **{{metric.most_recent_date}}** by carrier **{{metric.most_recent_carrier}}** (Flight **{{metric.most_recent_flight}}**), using aircraft **{{metric.most_recent_tail}}**. The full route was:
`{{metric.most_recent_route}}`

Departure times from each origin were: `{{metric.most_recent_dep_times}}`.

## Operational Analysis
The top 10 results show that 8 hops is the operational limit for a single flight number/aircraft combination in this dataset. Most of these occurrences belong to Southwest Airlines, indicating a high-utilization point-to-point routing model. The repetition of specific routes (e.g., flight 3149 and flight 2787) suggests these are scheduled multi-stop service patterns rather than one-off occurrences.