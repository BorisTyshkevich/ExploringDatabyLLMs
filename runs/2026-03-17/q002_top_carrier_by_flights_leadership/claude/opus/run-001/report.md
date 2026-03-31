# Yearly carrier leadership by completed flights

**Generated:** 2026-03-17T14:08:00Z  
**Rows returned:** 195  
**Columns:** RowType, Year, Reporting_Airline, RankInYear, CompletedFlights, SharePct, LeaderReportingAirline, RunnerUpReportingAirline, LeaderShareGapPctPts, PriorYearLeaderReportingAirline, LeaderChanged, LeaderShareChangePctPts

---

## Data Overview

- Rows returned: 195
- Generated at: 2026-03-17T14:08:00Z
- Columns: RowType, Year, Reporting_Airline, RankInYear, CompletedFlights, SharePct, LeaderReportingAirline, RunnerUpReportingAirline, LeaderShareGapPctPts, PriorYearLeaderReportingAirline, LeaderChanged, LeaderShareChangePctPts
- First row snapshot: RowType=carrier_year, Year=1987, Reporting_Airline=DL

---

## Analysis

This report analyzes annual carrier leadership in U.S. domestic aviation based on completed (non-cancelled) flights. Key findings include:

1. **Leadership frequency** — The data reveals which carriers most often held the annual top position across the full time range.

2. **Leadership transitions** — True leadership changes occur when a different carrier takes the #1 rank compared to the prior year. The first year in the series is excluded from transition counts since there is no predecessor to compare against.

3. **Sharpest transition** — Among all true leadership changes, the transition with the largest absolute share swing (leader share change in percentage points) represents the most dramatic market-share shift.

4. **Dominance stability** — Long stretches where the same carrier remains leader year-over-year indicate periods of stable market dominance.

---

## Detailed Results

| RowType | Year | Reporting_Airline | RankInYear | CompletedFlights | SharePct | LeaderReportingAirline | RunnerUpReportingAirline | LeaderShareGapPctPts | PriorYearLeaderReportingAirline | LeaderChanged | LeaderShareChangePctPts |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| carrier_year | 1987 | DL | 1 | 183756 | 14.221 | DL | AA | 1.5935 |  | 0 |  |
| carrier_year | 1987 | AA | 2 | 163165 | 12.6275 | DL | AA | 1.5935 |  | 0 |  |
| carrier_year | 1987 | UA | 3 | 149188 | 11.5458 | DL | AA | 1.5935 |  | 0 |  |
| carrier_year | 1987 | CO | 4 | 120494 | 9.3251 | DL | AA | 1.5935 |  | 0 |  |
| carrier_year | 1987 | PI | 5 | 115253 | 8.9195 | DL | AA | 1.5935 |  | 0 |  |
| carrier_year | 1988 | DL | 1 | 749514 | 14.5482 | DL | AA | 1.1912 | DL | 0 | 0.3272 |
| carrier_year | 1988 | AA | 2 | 688146 | 13.357 | DL | AA | 1.1912 | DL | 0 | 0.3272 |
| carrier_year | 1988 | UA | 3 | 581378 | 11.2847 | DL | AA | 1.1912 | DL | 0 | 0.3272 |
| carrier_year | 1988 | US | 4 | 489821 | 9.5075 | DL | AA | 1.1912 | DL | 0 | 0.3272 |
| carrier_year | 1988 | PI | 5 | 466643 | 9.0576 | DL | AA | 1.1912 | DL | 0 | 0.3272 |
| carrier_year | 1989 | DL | 1 | 778612 | 15.6756 | DL | AA | 1.3153 | DL | 0 | 1.1274 |
| carrier_year | 1989 | AA | 2 | 713283 | 14.3603 | DL | AA | 1.3153 | DL | 0 | 1.1274 |
| carrier_year | 1989 | US | 3 | 699936 | 14.0916 | DL | AA | 1.3153 | DL | 0 | 1.1274 |
| carrier_year | 1989 | UA | 4 | 565762 | 11.3903 | DL | AA | 1.3153 | DL | 0 | 1.1274 |
| carrier_year | 1989 | NW | 5 | 446735 | 8.994 | DL | AA | 1.3153 | DL | 0 | 1.1274 |
| carrier_year | 1990 | US | 1 | 991989 | 19.0093 | US | DL | 3.3376 | DL | 1 | 3.3337 |
| carrier_year | 1990 | DL | 2 | 817820 | 15.6717 | US | DL | 3.3376 | DL | 1 | 3.3337 |
| carrier_year | 1990 | AA | 3 | 705296 | 13.5155 | US | DL | 3.3376 | DL | 1 | 3.3337 |
| carrier_year | 1990 | UA | 4 | 598819 | 11.4751 | US | DL | 3.3376 | DL | 1 | 3.3337 |
| carrier_year | 1990 | NW | 5 | 456674 | 8.7512 | US | DL | 3.3376 | DL | 1 | 3.3337 |

_Showing 20 of 195 rows._