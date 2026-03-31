# {{question_title}}

{{data_overview_md}}

## Key finding

The maximum observed hop count is **{{metric.max_hops}}** legs in a single day for one aircraft operating under one flight number. All {{metric.max_hops}}-hop itineraries belong to **{{metric.carrier}}** (Southwest Airlines), reflecting its characteristic point-to-point, multi-stop scheduling model.

## Most recent maximum-hop itinerary

The most recent {{metric.max_hops}}-hop itinerary was **{{metric.most_recent_carrier}} {{metric.most_recent_flight}}** on **{{metric.most_recent_date}}**, aircraft **{{metric.most_recent_tail}}**, flying the route:

> {{metric.most_recent_route}}

Departure times from each origin: {{metric.most_recent_schedule}}

The aircraft began its day at ISP (Long Island MacArthur) at 05:43 and reached SEA (Seattle-Tacoma) after 20:41, covering eight legs coast-to-coast over roughly 15 hours.

## Route repetition and clustering

Among the top 10 longest itineraries, strong route repetition is evident:

- **WN 3149** (CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN) appears **4 times** across January-February 2024, indicating a recurring weekly schedule pattern.
- **WN 2787** (MSY -> ATL -> CMH -> BWI -> RDU -> BNA -> DTW -> MDW -> LAX) appears **3 times** across September-October 2022.
- **WN 154** (ELP -> DAL -> LIT -> ATL -> RIC -> MDW -> MCI -> PHX -> SAN) appears **twice** in April 2023.
- **WN 366** appears once with a distinct route (ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA).

These are clearly **repeating weekly operating patterns**, not one-off itineraries. Southwest regularly assigns a single flight number to an aircraft that hops across the country through 8 or more cities in a single day.

## Top 10 results

{{result_table_md}}