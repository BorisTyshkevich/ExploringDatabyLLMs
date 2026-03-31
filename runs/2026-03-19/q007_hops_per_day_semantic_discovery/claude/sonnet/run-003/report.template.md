# Highest Daily Hops for One Aircraft on One Flight Number

{{data_overview_md}}

## Maximum Hop Count

The highest number of hops recorded for a single aircraft operating under a single flight number on one calendar day is **{{metric.max_hops}}**. This means the aircraft completed {{metric.max_hops}} individual flight legs under the same flight number within a single day.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Carrier | {{metric.top_carrier}} |
| Flight Number | {{metric.top_flight_number}} |
| Aircraft (Tail) | {{metric.top_tail_number}} |
| Date | {{metric.top_date}} |
| Route | {{metric.top_route}} |

The route above shows each origin airport followed by its actual departure time (HHMM) in parentheses, ending at the final destination.

## Operating Pattern Analysis

{{metric.pattern_summary}}

## Top 10 Longest Itineraries

{{result_table_md}}

### Notable Observations

- All top-10 itineraries belong to **{{metric.dominant_carrier}}**, reflecting that carrier's hub-and-spoke multi-stop routing strategy.
- The most repeated flight number in the top 10 is **{{metric.most_repeated_flight}}**, suggesting a recurrent scheduled itinerary rather than an ad-hoc assignment.
- Route clustering across dates for the same flight number indicates a **{{metric.recurrence_pattern}}** operating schedule.

_Generated: {{generated_at}}_