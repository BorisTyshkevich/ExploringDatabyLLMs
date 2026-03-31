# Yearly carrier leadership by completed flights

Generated from `result.json` at 2026-03-17T14:08:29Z with 195 rows across `RowType, Year, Reporting_Airline, RankInYear, CompletedFlights, SharePct, LeaderReportingAirline, RunnerUpReportingAirline, LeaderShareGapPctPts, PriorYearLeaderReportingAirline, LeaderChanged, LeaderShareChangePctPts`.

## Analytical focus
- Rows returned: 195
- Generated at: 2026-03-17T14:08:29Z
- Columns: RowType, Year, Reporting_Airline, RankInYear, CompletedFlights, SharePct, LeaderReportingAirline, RunnerUpReportingAirline, LeaderShareGapPctPts, PriorYearLeaderReportingAirline, LeaderChanged, LeaderShareChangePctPts
- First row snapshot: RowType=carrier_year, Year=1987, Reporting_Airline=DL

Use the annual leader rows where `RankInYear = 1` to assess four questions: how often each carrier leads across the full span, which years represent true leadership changes (`LeaderChanged = 1` with a non-null prior leader), which of those transitions has the largest absolute `LeaderShareChangePctPts`, and whether repeated runs by the same `LeaderReportingAirline` indicate extended dominance.

## Result detail
| RowType | Year | Reporting_Airline | RankInYear | CompletedFlights | SharePct | LeaderReportingAirline | RunnerUpReportingAirline | LeaderShareGapPctPts | PriorYearLeaderReportingAirline | LeaderChanged | LeaderShareChangePctPts |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| carrier_year | 1987 | DL | 1 | 183756 | 14.221048631689575 | DL | AA | 1.5935567403247788 |  | 0 |  |
| carrier_year | 1987 | AA | 2 | 163165 | 12.627491891364796 |  |  |  |  |  |  |
| carrier_year | 1987 | UA | 3 | 149188 | 11.545798794404016 |  |  |  |  |  |  |
| carrier_year | 1987 | CO | 4 | 120494 | 9.325143308663684 |  |  |  |  |  |  |
| carrier_year | 1987 | PI | 5 | 115253 | 8.919537418903975 |  |  |  |  |  |  |
| carrier_year | 1988 | DL | 1 | 749514 | 14.548209380828515 | DL | AA | 1.1911645590111508 | DL | 0 | 0.32716074913894033 |
| carrier_year | 1988 | AA | 2 | 688146 | 13.357044821817365 |  |  |  |  |  |  |
| carrier_year | 1988 | UA | 3 | 581378 | 11.284657622682593 |  |  |  |  |  |  |
| carrier_year | 1988 | US | 4 | 489821 | 9.507518828369857 |  |  |  |  |  |  |
| carrier_year | 1988 | PI | 5 | 466643 | 9.057629437339344 |  |  |  |  |  |  |
| carrier_year | 1989 | DL | 1 | 778612 | 15.675589159327446 | DL | AA | 1.3152514528284982 | DL | 0 | 1.1273797784989306 |
| carrier_year | 1989 | AA | 2 | 713283 | 14.360337706498948 |  |  |  |  |  |  |
| carrier_year | 1989 | US | 3 | 699936 | 14.091626090816755 |  |  |  |  |  |  |
| carrier_year | 1989 | UA | 4 | 565762 | 11.390336488468472 |  |  |  |  |  |  |
| carrier_year | 1989 | NW | 5 | 446735 | 8.993997425023178 |  |  |  |  |  |  |
| carrier_year | 1990 | US | 1 | 991989 | 19.00931984397621 | US | DL | 3.3375715132985277 | DL | 1 | 3.3337306846487653 |
| carrier_year | 1990 | DL | 2 | 817820 | 15.671748330677683 |  |  |  |  |  |  |
| carrier_year | 1990 | AA | 3 | 705296 | 13.515469676253513 |  |  |  |  |  |  |
| carrier_year | 1990 | UA | 4 | 598819 | 11.475068674803845 |  |  |  |  |  |  |
| carrier_year | 1990 | NW | 5 | 456674 | 8.751167735154313 |  |  |  |  |  |  |

_Showing 20 of 195 rows._