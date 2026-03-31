# {{question_title}}

_Generated {{generated_at}} · {{row_count}} rows_

## Maximum Hop Count

The highest recorded hop count is **{{metric.max_hops}}** legs flown by a single aircraft under a single flight number in one calendar day. All top-10 entries are operated by **{{metric.dominant_carrier}}**. This is not a one-off anomaly: the same multi-hop routing patterns recur across multiple weeks on the same flight numbers, confirming these as deliberately scheduled point-to-point network itineraries.

## Most Recent Maximum-Hop Itinerary

| Field | Value |
|---|---|
| Date | {{metric.most_recent_date}} |
| Carrier | {{metric.most_recent_carrier}} |
| Flight Number | {{metric.most_recent_flight_num}} |
| Aircraft (Tail) | {{metric.most_recent_aircraft}} |
| Hops | {{metric.max_hops}} |

**Route (chronological, with airport names and departure times):**

> {{metric.most_recent_route}}

## Route Repetition and Clustering

Among the top 10 longest itineraries, strong route repetition is evident by flight number. The same {{metric.max_hops}}-hop city sequences recur on consecutive weeks, confirming them as recurring scheduled patterns rather than irregular one-off operations. Route clusters are anchored by specific hub sequences through the Midwest, South, and West Coast.

## Top 10 Longest and Most Recent Itineraries

{{result_table_md}}

_Route column format: `CODE (Airport Name) HH:MM -> ... -> CODE (Airport Name)` — each entry shows the origin code, full airport name from the semantic layer, and scheduled departure time. The final entry shows the terminal destination only._