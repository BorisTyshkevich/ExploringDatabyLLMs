# {{question_title}}

*Generated: {{generated_at}} | Rows: {{row_count}}*

## Overview

{{data_overview_md}}

This analysis identifies departure delay hotspots for Delta Air Lines (DL) flights originating from Atlanta Hartsfield–Jackson (ATL). A hotspot is defined as a destination–time block combination that has had at least 1,000 completed flights across months where the monthly count exceeded 40.

## Key Findings

### Worst Hotspot

The top-ranked hotspot represents the destination and departure time block with the highest average departure delay. The ranking considers average delay, 90th-percentile delay, and percentage of flights delayed 15+ minutes.

### Persistence vs. Concentration

The `QualifyingMonths` column indicates how many months the hotspot met the monthly activity threshold. A high count suggests a persistent structural issue, while a low count with a narrow date range may indicate a temporary bottleneck.

### Top 5 Implications

The top 5 hotspot cells reveal which combinations of destination and departure window generate the most consistent delay pressure at ATL for Delta. These patterns may reflect:

- Hub connection timing conflicts
- Downstream airport constraints
- Seasonal demand spikes on specific routes
- Gate and runway congestion during peak time blocks

## Hotspot Summary

{{result_table_md}}

---

*Data filtered to DL carrier, ATL origin, non-cancelled flights. Hotspots require ≥1,000 flights across months with ≥40 monthly flights.*