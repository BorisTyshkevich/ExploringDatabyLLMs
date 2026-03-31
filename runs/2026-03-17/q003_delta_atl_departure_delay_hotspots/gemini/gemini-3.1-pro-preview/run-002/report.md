# Delta ATL departure delay hotspots by destination and time block

- Rows returned: 832
- Generated at: 2026-03-17T13:08:38Z
- Columns: RowType, MonthStart, Dest, DepTimeBlk, QualifyingMonths, CompletedFlights, AvgDepDelayMinutes, P90DepDelayMinutes, DepDel15Pct, FirstQualifyingMonth, LastQualifyingMonth, HotspotRank
- First row snapshot: RowType=hotspot_summary, Dest=LGA, DepTimeBlk=1900-1959

## Executive Summary

Based on the departure delay statistics, the worst hotspot for Delta departures from ATL is **LGA at 1900-1959**. This time block experiences the highest average departure delay (~24.8 minutes) and 90th percentile departure delay (~66.1 minutes), indicating severe late-evening congestion. 

This hotspot appears significantly concentrated in a narrower era, primarily showing qualifying flights in the late 1990s through 2001, which suggests historical structural challenges on the ATL-LGA route during evening peak hours rather than a currently persistent issue across all decades.

Analyzing the top 5 hotspot cells reveals a strong concentration of pressure during the **1500-1959 (afternoon to evening) window** to major hub destinations:
1. **LGA (1900-1959)**
2. **LGA (1500-1559)**
3. **MCO (1700-1759)**
4. **DFW (1700-1759)**
5. **ORD (1500-1559)**

This pattern indicates that Delta's ATL departure pressure is heavily driven by connecting bank structures feeding into East Coast and Central business hubs (New York, Orlando, Dallas, Chicago) during the late afternoon and evening periods, where cascading network delays are most likely to accumulate.

## Hotspot Rankings

| RowType | MonthStart | Dest | DepTimeBlk | QualifyingMonths | CompletedFlights | AvgDepDelayMinutes | P90DepDelayMinutes | DepDel15Pct | FirstQualifyingMonth | LastQualifyingMonth | HotspotRank |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| hotspot_summary |  | LGA | 1900-1959 | 29 | 1630 | 24.83 | 66.1 | 39.75 | 1998-11-01T00:00:00Z | 2001-03-01T00:00:00Z | 1 |
| hotspot_summary |  | LGA | 1500-1559 | 35 | 1793 | 21.31 | 56 | 34.02 | 1998-11-01T00:00:00Z | 2001-10-01T00:00:00Z | 2 |
| hotspot_summary |  | MCO | 1700-1759 | 42 | 2516 | 18.18 | 47 | 32.83 | 1993-06-01T00:00:00Z | 2013-12-01T00:00:00Z | 3 |
| hotspot_summary |  | DFW | 1700-1759 | 91 | 5339 | 17.42 | 42 | 34.07 | 1994-06-01T00:00:00Z | 2005-01-01T00:00:00Z | 4 |
| hotspot_summary |  | ORD | 1500-1559 | 25 | 1482 | 17.08 | 44 | 31.51 | 1995-05-01T00:00:00Z | 1997-05-01T00:00:00Z | 5 |
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
| hotspot_summary |  | DFW | 0900-0959 | 18 | 1070 | 8.95 | 26 | 19.72 | 2003-08-01T00:00:00Z | 2005-01-01T00:00:00Z | 17 |
| hotspot_summary |  | SLC | 0800-0859 | 56 | 3132 | 8.91 | 22 | 15.01 | 1997-04-01T00:00:00Z | 2005-09-01T00:00:00Z | 18 |
| hotspot_summary |  | DCA | 1300-1359 | 42 | 2283 | 8.79 | 19 | 15.72 | 1992-06-01T00:00:00Z | 1995-11-01T00:00:00Z | 19 |
| hotspot_summary |  | MCO | 1300-1359 | 42 | 2483 | 7.94 | 20 | 14.14 | 1993-06-01T00:00:00Z | 2015-01-01T00:00:00Z | 20 |

_Showing 20 of 832 rows._