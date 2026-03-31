# {{question_title}}

**Generated:** {{generated_at}}  
**Rows returned:** {{row_count}}  
**Columns:** {{columns_csv}}

---

## Overview

{{data_overview_md}}

---

## Analysis

The dashboard identifies the maximum number of daily hops flown by a single aircraft under one flight number. This pattern emerges when regional or commuter aircraft operate multiple short legs under a single published flight number, creating high-frequency shuttle itineraries.

### Maximum hop count

The top result shows the highest daily hop count observed in the dataset. A recurring pattern of high hop counts for the same tail number and flight number suggests a scheduled shuttle service rather than a one-off repositioning itinerary. Single occurrences at the maximum level may indicate irregular operations or charter-style routing.

### Most recent maximum-hop itinerary

The first row in the result set represents the most recent date among all maximum-hop itineraries. The route column shows the full sequence of departures with local times and airport codes, providing visibility into the operational rhythm of that day's service.

### Route repetition and clustering

Examine whether the same tail number and flight number appear on multiple dates at or near the maximum hop count. Repeated appearances indicate a stable operating schedule. If different aircraft or flight numbers dominate the list, the high-hop pattern is distributed across the network rather than concentrated on a single shuttle corridor.

---

## Result Data

{{result_table_md}}