# Worst origin airports by departure on-time performance

> Which airport ranks worst on departure on-time performance?

Chicago Midway International (MDW) ranks worst on departure on-time performance among qualifying airports. Over the entire historical dataset, MDW recorded an on-time performance of only 76.5%, significantly lower than its peers.

- Rows returned: 160
- Columns: OriginCode, OriginCityName, completed_departures, on_time_performance, avg_departure_delay, p90_departure_delay, first_date, last_date

| OriginCode | OriginCityName | completed_departures | on_time_performance | avg_departure_delay | p90_departure_delay | first_date | last_date |
| --- | --- | --- | --- | --- | --- | --- | --- |
| MDW | Chicago, IL | 2.512354e+06 | 0.7652 | 11.8 | 38 | 1987-10-01T00:00:00Z | 2025-11-30T00:00:00Z |

> How large is the spread between the worst airport and the middle of the ranked set?

There is an 8.5 percentage point gap between the worst airport and the median of the ranked set. While Chicago Midway (MDW) sits at the bottom with a 76.5% on-time performance, an average qualifying airport in the middle of the rankings (such as Islip, NY) achieves an on-time performance of approximately 85.0%.

- Rows returned: 2
- Columns: OriginCode, OriginCityName, completed_departures, on_time_performance, avg_departure_delay, p90_departure_delay, first_date, last_date, rnk, total_airports

| OriginCode | OriginCityName | completed_departures | on_time_performance | avg_departure_delay | p90_departure_delay | first_date | last_date | rnk | total_airports |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| MDW | Chicago, IL | 2.512354e+06 | 0.7652 | 11.8 | 36 | 1987-10-01T00:00:00Z | 2025-11-30T00:00:00Z | 1 | 160 |

> Are the weakest airports mostly major hubs, or is the bottom group more mixed?

The weakest performing airports are overwhelmingly major airline hubs. The bottom tier of the rankings is dominated by heavily trafficked hubs such as Chicago O'Hare (ORD), Newark (EWR), Houston (HOU), Dallas (DAL), JFK, and Baltimore (BWI), rather than being a mixed group of large and mid-sized regional airports.

- Rows returned: 20
- Columns: OriginCode, OriginCityName, completed_departures, on_time_performance, avg_departure_delay, p90_departure_delay, first_date, last_date

| OriginCode | OriginCityName | completed_departures | on_time_performance | avg_departure_delay | p90_departure_delay | first_date | last_date |
| --- | --- | --- | --- | --- | --- | --- | --- |
| MDW | Chicago, IL | 2.512354e+06 | 0.7652 | 11.8 | 35 | 1987-10-01T00:00:00Z | 2025-11-30T00:00:00Z |
