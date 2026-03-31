# {{question_title}}

## Overview

{{data_overview_md}}

## Key Finding

The maximum number of hops observed for a single aircraft operating one flight number in a single day is **{{metric.max_hops}}**. This pattern is not a one-off anomaly — it appears repeatedly across multiple dates and flight numbers, all operated by **{{metric.dominant_carrier}}**, reflecting that carrier's point-to-point operating model where a single flight number can traverse the entire country in a chain of short-haul segments.

## Most Recent Maximum-Hop Itinerary

The most recent 8-hop itinerary is **{{metric.most_recent_carrier}} {{metric.most_recent_flight}}** on **{{metric.most_recent_date}}**, aircraft **{{metric.most_recent_tail}}**, flying the route **{{metric.most_recent_route}}** with departures at **{{metric.most_recent_dep_times}}**.

This itinerary spans coast-to-coast from the New York area (ISP) to the Pacific Northwest (SEA) in eight legs over roughly 15 hours.

## Route Repetition and Clustering

Among the top 10 longest itineraries, strong route repetition is evident:

- **WN 3149** appears four times (Jan–Feb 2024) on the identical route {{metric.wn3149_route}}, indicating a fixed weekly schedule.
- **WN 2787** appears three times (Sep–Oct 2022) repeating the route {{metric.wn2787_route}}.
- **WN 154** appears twice (Apr 2023) on the route {{metric.wn154_route}}.

All 10 entries are Southwest Airlines (WN), consistent with its unique operational model of multi-stop through-flights under a single flight number.

## Result Detail

{{result_table_md}}