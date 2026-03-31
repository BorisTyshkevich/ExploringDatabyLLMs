`answer.raw.json` written. Here's a summary of findings:

**main**: The global maximum is **8 distinct hops** in one day under one flight number. All 10 unique routes are **Southwest Airlines (WN)** milk-run itineraries, the most recent being N957WN / WN 366 on 2024-12-01 (ISP-BWI-MYR-BNA-VPS-DAL-LAS-OAK-SEA). The highest-recurrence route is LGA-STL-ICT-DEN-COS-PHX-ELP-HOU-JAN at 47 days.

> Note: A 10-hop record existed in the raw data (OO/6195/2003-02-28) but was caused by duplicate rows with identical origin, destination, and departure times — not a real 10-segment itinerary. Using `countDistinct(OriginCode, DestCode, CRSDepTime)` eliminates that artifact and yields the correct global max of 8.

**q1**: Split into three tiers — 2 one-offs (≤2 days), 4 modest patterns (5–7 days), and 4 strongly recurring schedules (12–47 days). 8 of 10 routes are recurring.

**q2**: All airports lie in the contiguous US. Routes sweep coast-to-coast through Southwest's own network of secondary airports (DAL over DFW, HOU over IAH, MDW over ORD), bridging East/Gulf Coast origins through Sun Belt/Mountain West transitions to West Coast termini.
