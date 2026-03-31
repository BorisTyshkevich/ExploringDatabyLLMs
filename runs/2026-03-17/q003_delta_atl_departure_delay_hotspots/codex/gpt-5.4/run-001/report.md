# Delta ATL departure delay hotspots by destination and time block

- Rows returned: 832
- Generated at: 2026-03-17T12:30:15Z
- Columns: RowType, MonthStart, Dest, DepTimeBlk, QualifyingMonths, CompletedFlights, AvgDepDelayMinutes, P90DepDelayMinutes, DepDel15Pct, FirstQualifyingMonth, LastQualifyingMonth, HotspotRank
- First row snapshot: RowType=hotspot_summary, Dest=LGA, DepTimeBlk=1900-1959

The worst hotspot is the `hotspot_summary` cell with `HotspotRank = 1` in the rendered result set. Its destination, departure time block, average delay, tail-risk delay (`P90DepDelayMinutes`), and `DepDel15Pct` define the highest-pressure ATL departure cell for Delta in this output.

Persistence versus concentration is indicated by `QualifyingMonths` together with the span from `FirstQualifyingMonth` to `LastQualifyingMonth`. A long span with many qualifying months points to a structural hotspot, while a shorter span suggests a narrower disruption era.

The top 5 hotspot cells show whether ATL departure pressure is concentrated in a small set of Delta destination-time banks or distributed across multiple banks. Compare rank order, delay severity, and qualifying-month coverage across those five cells to judge whether the bottleneck is sustained, bank-specific, or destination-specific.

## Result Detail

| RowType | MonthStart | Dest | DepTimeBlk | QualifyingMonths | CompletedFlights | AvgDepDelayMinutes | P90DepDelayMinutes | DepDel15Pct | FirstQualifyingMonth | LastQualifyingMonth | HotspotRank |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| hotspot_summary |  | LGA | 1900-1959 | 29 | 1630 | 24.83 | 66.09 | 39.75 | 1998-11-01T00:00:00Z | 2001-03-01T00:00:00Z | 1 |
| hotspot_summary |  | LGA | 1500-1559 | 35 | 1793 | 21.31 | 56 | 34.02 | 1998-11-01T00:00:00Z | 2001-10-01T00:00:00Z | 2 |
| hotspot_summary |  | MCO | 1700-1759 | 42 | 2516 | 18.18 | 47 | 32.83 | 1993-06-01T00:00:00Z | 2013-12-01T00:00:00Z | 3 |
| hotspot_summary |  | DFW | 1700-1759 | 91 | 5339 | 17.42 | 42 | 34.07 | 1994-06-01T00:00:00Z | 2005-01-01T00:00:00Z | 4 |
| hotspot_summary |  | ORD | 1500-1559 | 25 | 1482 | 17.07 | 44 | 31.51 | 1995-05-01T00:00:00Z | 1997-05-01T00:00:00Z | 5 |
| hotspot_summary |  | SLC | 1500-1559 | 28 | 1647 | 15.3 | 41 | 28.9 | 1997-04-01T00:00:00Z | 2000-03-01T00:00:00Z | 6 |
| hotspot_summary |  | LAX | 1700-1759 | 21 | 1229 | 15.05 | 35 | 30.02 | 2003-05-01T00:00:00Z | 2005-01-01T00:00:00Z | 7 |
| hotspot_summary |  | MCO | 2200-2259 | 21 | 1219 | 15.02 | 52 | 25.43 | 1999-07-01T00:00:00Z | 2001-03-01T00:00:00Z | 8 |
| hotspot_summary |  | TPA | 1600-1659 | 29 | 1714 | 14.41 | 40 | 27.77 | 1998-12-01T00:00:00Z | 2002-11-01T00:00:00Z | 9 |
| hotspot_summary |  | DFW | 1300-1359 | 101 | 5822 | 14.02 | 32 | 27.4 | 1992-11-01T00:00:00Z | 2001-03-01T00:00:00Z | 10 |
| hotspot_summary |  | LGA | 1200-1259 | 29 | 1472 | 13.53 | 38 | 22.83 | 1998-11-01T00:00:00Z | 2001-03-01T00:00:00Z | 11 |
| hotspot_summary |  | TPA | 1700-1759 | 17 | 1010 | 13.18 | 33 | 21.98 | 2002-12-01T00:00:00Z | 2005-01-01T00:00:00Z | 12 |
| hotspot_summary |  | PBI | 1700-1759 | 18 | 1074 | 12.86 | 35 | 24.77 | 1996-12-01T00:00:00Z | 1998-05-01T00:00:00Z | 13 |
| hotspot_summary |  | EWR | 1400-1459 | 35 | 1828 | 12.06 | 40 | 21.12 | 2002-02-01T00:00:00Z | 2005-01-01T00:00:00Z | 14 |
| hotspot_summary |  | DFW | 1500-1559 | 85 | 4702 | 12.01 | 31 | 23.14 | 1995-12-01T00:00:00Z | 2004-11-01T00:00:00Z | 15 |
| hotspot_summary |  | DFW | 0800-0859 | 48 | 2648 | 10.98 | 24.3 | 16.88 | 1997-04-01T00:00:00Z | 2001-03-01T00:00:00Z | 16 |
| hotspot_summary |  | DFW | 0900-0959 | 18 | 1070 | 8.94 | 26 | 19.72 | 2003-08-01T00:00:00Z | 2005-01-01T00:00:00Z | 17 |
| hotspot_summary |  | SLC | 0800-0859 | 56 | 3132 | 8.91 | 22 | 15.01 | 1997-04-01T00:00:00Z | 2005-09-01T00:00:00Z | 18 |
| hotspot_summary |  | DCA | 1300-1359 | 42 | 2283 | 8.78 | 19 | 15.72 | 1992-06-01T00:00:00Z | 1995-11-01T00:00:00Z | 19 |
| hotspot_summary |  | MCO | 1300-1359 | 42 | 2483 | 7.94 | 20 | 14.14 | 1993-06-01T00:00:00Z | 2015-01-01T00:00:00Z | 20 |

_Showing 20 of 832 rows._