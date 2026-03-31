# Worst winter carrier-origin pairs by departure performance

> Which winter carrier-airport pair ranks worst overall?

Using December through February and excluding cancelled or diverted flights, the worst qualifying pair is `DH` departing `Chicago O'Hare International (ORD)`. With a minimum threshold of 1,000 winter departures across the full history, it has 19,929 winter flights, the lowest departure on-time rate at 56.53%, and an average departure delay of 27.06 minutes.

- Rows returned: 25
- Columns: Carrier, OriginCode, DisplayAirportName, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause

| Carrier | OriginCode | DisplayAirportName | winter_flights | dep_otp_pct | avg_dep_delay_min | flights_with_any_reported_cause | flights_without_reported_cause |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DH | ORD | Chicago O'Hare International | 19929 | 56.53 | 27.06 | 5681 | 14248 |

> Are the worst pairs driven more by weather or by operational causes?

For the leading weak winter pairs that actually report delay causes, operational causes dominate reported delay minutes rather than weather. In the verified top-10 weak set, every pair with measured cause data is majority operational, with weather contributing only 1.5% to 26.8% of reported cause minutes. Cause reporting is incomplete, though: `PI-DFW`, `PI-LAX`, and `PI-DAY` show no reported cause minutes at all, so the weather-versus-operations conclusion applies only to the measured subset.

- Rows returned: 10
- Columns: Carrier, OriginCode, winter_flights, dep_otp_pct, avg_dep_delay_min, flights_with_any_reported_cause, flights_without_reported_cause, weather_delay_min, carrier_delay_min, nas_delay_min, security_delay_min, late_aircraft_delay_min, total_reported_cause_min, weather_share_pct, operational_share_pct

| Carrier | OriginCode | winter_flights | dep_otp_pct | avg_dep_delay_min | flights_with_any_reported_cause | flights_without_reported_cause | weather_delay_min | carrier_delay_min | nas_delay_min | security_delay_min | late_aircraft_delay_min | total_reported_cause_min | weather_share_pct | operational_share_pct |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DH | ORD | 19929 | 56.53 | 27.06 | 5681 | 14248 | 55176 | 123407 | 41492 | 116 | 133239 | 353430 | 15.6 | 84.4 |

> Are the weakest pairs concentrated in a small number of carriers or airports?

The weak set is only mildly concentrated in carriers and not very concentrated in airports. Looking at the top 10 qualifying weak winter pairs, `PI` accounts for 3 pairs (30%), while `ORD` is the only airport that appears more than once with 2 pairs (20%). Every other carrier and airport appears once, so the concentration is more noticeable by carrier than by origin airport.

- Rows returned: 17
- Columns: concentration_type, entity, weak_pair_count, weak_pair_share_pct

| concentration_type | entity | weak_pair_count | weak_pair_share_pct |
| --- | --- | --- | --- |
| carrier | OH | 1 | 10 |
