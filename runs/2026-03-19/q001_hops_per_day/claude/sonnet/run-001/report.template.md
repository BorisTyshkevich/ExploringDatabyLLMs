# Highest Daily Hops for One Aircraft on One Flight Number

## Overview

This report identifies the longest single-day itineraries flown by one aircraft tail under one flight number. A **hop** is one leg (one Origin→Dest segment); an itinerary of N hops contains N+1 distinct airports.

The maximum observed hop count across all carriers and dates is **{{metric.max_hops}}**. All top-10 itineraries reach this ceiling, meaning no aircraft exceeded 8 legs in a single day under a single flight number.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | {{metric.most_recent_carrier}} |
| Flight Number | {{metric.most_recent_flight}} |
| Tail Number | {{metric.most_recent_tail}} |
| Date | {{metric.most_recent_date}} |
| Route | {{metric.most_recent_route}} |
| Departure Times | {{metric.most_recent_dep_times}} |

## Repeated Operating Pattern vs. One-Off

The 8-hop ceiling is **not a one-off event** — it appears across multiple carriers, years, and flight numbers. The single most repeated pattern in the top 10 is **{{metric.top_recurring_flight}}**, which flew the identical {{metric.max_hops}}-hop route {{metric.top_recurring_count}} times:

> {{metric.top_recurring_route}}

This signals a scheduled, recurring turnaround pattern rather than an irregular re-routing.

## Top 10 Longest Itineraries

{{result_table_md}}

## Key Observations

- All top-10 itineraries are operated by **{{metric.dominant_carrier}}**, consistent with Southwest's point-to-point network model where a single aircraft and flight number chain many short segments in one day.
- Route clustering around certain city pairs (e.g., the recurring CLE–DEN spine) indicates structured schedule blocks, not ad-hoc assignments.
- Departure times span roughly 05:00–23:00 local, confirming full-day aircraft utilization on these itineraries.

---
_Generated: {{generated_at}}_