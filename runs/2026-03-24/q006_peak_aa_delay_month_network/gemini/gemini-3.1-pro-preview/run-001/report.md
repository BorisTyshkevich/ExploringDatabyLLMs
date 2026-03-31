# American Airlines peak network delay month and contributors

> Which month is the single worst American Airlines month for departure delays?

July 2024 is the single worst month for American Airlines departure delays. During this month, average departure delays reached a peak of 36.33 minutes, and 38.04% of the 86,083 completed flights departed 15 or more minutes late.

- Rows returned: 458
- Columns: Year, Month, FlightVolume, AvgDepDelay, PctDel15

| Year | Month | FlightVolume | AvgDepDelay | PctDel15 |
| --- | --- | --- | --- | --- |
| 2024 | 7 | 86083 | 36.33 | 0.3804 |

> Which origins contribute most to that peak month?

The delays in July 2024 were driven primarily by Dallas/Fort Worth (DFW) and Charlotte (CLT). DFW generated over 576,000 delay minutes (averaging 38.5 minutes per flight), while CLT was responsible for roughly 549,000 delay minutes with severe average delays reaching 52 minutes.

- Rows returned: 93
- Columns: OriginCode, FlightVolume, AvgDepDelay, PctDel15, TotalDelayMinutes, DelayedFlights

| OriginCode | FlightVolume | AvgDepDelay | PctDel15 | TotalDelayMinutes | DelayedFlights |
| --- | --- | --- | --- | --- | --- |
| DFW | 14962 | 38.52 | 0.4659 | 576310 | 6971 |

> Which routes contribute most to that peak month?

The most severely impacted routes were hub-to-hub and hub-to-spoke connections from DFW and CLT. The top route contributor was DFW to Los Angeles (LAX) with over 18,000 delay minutes, closely followed by CLT to Orlando (MCO) and DFW to San Antonio (SAT).

- Rows returned: 799
- Columns: OriginCode, DestCode, Route, FlightVolume, AvgDepDelay, PctDel15, TotalDelayMinutes, DelayedFlights

| OriginCode | DestCode | Route | FlightVolume | AvgDepDelay | PctDel15 | TotalDelayMinutes | DelayedFlights |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DFW | LAX | DFW-LAX | 443 | 40.89 | 0.465 | 18114 | 206 |

> Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

The peak is highly concentrated at the origin level but broadly distributed across routes. In July 2024, the top 5 origins (led overwhelmingly by DFW and CLT) accounted for nearly 52% of the total network departure delay minutes. In contrast, the network's delays were widely spread across individual routes, with the top 10 worst routes representing only 5% of all network delay minutes.

- Rows returned: 1
- Columns: network_delay, top_5_origins_delay, top_10_routes_delay, top_5_origins_pct, top_10_routes_pct

| network_delay | top_5_origins_delay | top_10_routes_delay | top_5_origins_pct | top_10_routes_pct |
| --- | --- | --- | --- | --- |
| 3.127538e+06 | 1.624257e+06 | 156919 | 51.93 | 5.02 |
