# Highest daily hops for one aircraft on one flight number

- Rows returned: 10
- Generated at: 2026-03-19T18:38:30Z
- Columns: AircraftID, FlightNumber, Carrier, Date, Hops, Route, DepTimes
- First row snapshot: AircraftID=N957WN, FlightNumber=366, Carrier=WN

The highest number of hops recorded for a single aircraft using a consistent flight number in one day is **8**. This pattern is observed across multiple dates, notably with Southwest Airlines (WN) flight 3149, which frequently operates an 8-hop route such as `CLE -> BNA -> PNS -> HOU -> MCI -> PHX -> BUR -> OAK -> DEN`.

## Longest Recent Itinerary
The most recent instance of an 8-hop itinerary was on **2024-12-01** by carrier **WN** (Flight **366**), using aircraft **N957WN**. The full route was:
`ISP -> BWI -> MYR -> BNA -> VPS -> DAL -> LAS -> OAK -> SEA`

Departure times from each origin were: `0543, 0810, 1020, 1142, 1401, 1643, 1828, 2041`.

## Operational Analysis
The top 10 results show that 8 hops is the operational limit for a single flight number/aircraft combination in this dataset. Most of these occurrences belong to Southwest Airlines, indicating a high-utilization point-to-point routing model. The repetition of specific routes (e.g., flight 3149 and flight 2787) suggests these are scheduled multi-stop service patterns rather than one-off occurrences.