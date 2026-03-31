# Highest Daily Hops for One Aircraft on One Flight Number

## Overview

{{data_overview_md}}

The maximum number of hops flown by a single aircraft under the same flight number in a single day is **{{metric.max_hops}}**, achieved exclusively by carrier **{{metric.max_hops_carrier}}**. All 10 entries in the top-10 list share this hop count, indicating that {{metric.max_hops}}-hop itineraries are a structured, repeating feature of Southwest Airlines scheduling rather than a one-off occurrence.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | {{metric.most_recent_carrier}} |
| Flight Number | {{metric.most_recent_flight_number}} |
| Aircraft (Tail) | {{metric.most_recent_tail}} |
| Date | {{metric.most_recent_date}} |
| Hops | {{metric.max_hops}} |
| Route | {{metric.most_recent_route}} |
| Departure Schedule | {{metric.most_recent_departure_schedule}} |

The route spans {{metric.most_recent_airport_count}} airports across the country, with departures beginning before 06:00 local time and the last segment departing after 20:00, covering roughly 16 hours of continuous flying operations.

## Top 10 Longest and Most Recent Itineraries

{{result_table_md}}

## Route Patterns and Clustering

The top-10 itineraries reveal strong route clustering by flight number:

- **Flight {{metric.flight_3149_label}}** operates the identical {{metric.max_hops}}-stop sequence {{metric.flight_3149_occurrences}} times across different dates and tail numbers, confirming it as a regularly scheduled transcontinental multi-hop service.
- **Flight {{metric.flight_154_label}}** and **Flight {{metric.flight_2787_label}}** each repeat their respective fixed routes on multiple dates, further demonstrating that {{metric.max_hops}}-hop itineraries are codified schedule patterns.
- All top-10 entries are operated by {{metric.max_hops_carrier}}, which historically employs a point-to-point multi-stop model suited to high daily utilization of individual aircraft.