# Yearly carrier leadership by completed flights

> Which carrier leads most often across the full time range?

WN (Southwest Airlines) leads most often, topping the industry in 26 of 39 years (2000–2025). Delta (DL) led for 11 years (1987–1999, except 1990–1991 when US Airways held the top), and US Airways (US) led for 2 years (1990–1991).

- Rows returned: 3
- Columns: carrier, years_led

| carrier | years_led |
| --- | --- |
| WN | 26 |

> When leadership changes, how large is the swing versus the prior leader?

There are 3 genuine leadership transitions in the full history. The 1990 transition from DL to US produced the largest swing: US Airways flew 991,989 flights versus Delta's prior-year total of 778,612, a swing of +213,377. The 1992 reversal back to DL was narrow at +9,842 flights over US's prior year, and the 2000 shift to WN was similarly narrow at +10,273 flights over DL's prior year. The average swing across all three transitions is about 78,000 flights, but that figure is dominated entirely by the 1990 outlier.

- Rows returned: 3
- Columns: Year, new_leader, new_leader_flights, prev_leader, prev_leader_flights, swing

| Year | new_leader | new_leader_flights | prev_leader | prev_leader_flights | swing |
| --- | --- | --- | --- | --- | --- |
| 1990 | US | 991989 | DL | 778612 | 213377 |

> Which transition is the sharpest?

The sharpest transition is 1990, when US Airways (US) displaced Delta (DL) with a swing of +213,377 completed flights versus Delta's 1989 output, and held a 174,169-flight lead over Delta as runner-up within that year. This is more than 20 times larger than either of the other two transitions (1992 and 2000).

- Rows returned: 1
- Columns: Year, new_leader, leader_flights, prev_leader, prev_leader_flights, swing_vs_prior, gap_vs_runner_up

| Year | new_leader | leader_flights | prev_leader | prev_leader_flights | swing_vs_prior | gap_vs_runner_up |
| --- | --- | --- | --- | --- | --- | --- |
| 1990 | US | 991989 | DL | 778612 | 213377 | 174169 |

> Does the market show long stable eras, or frequent turnover at the top?

The market is defined by long stable eras, not frequent turnover. In 38 full years of data (1988–2025), there were only 3 leadership transitions. Delta (DL) dominated for roughly 11 years (1987–1999 with a brief 2-year interruption by US Airways in 1990–1991). Southwest (WN) has held an unbroken reign since 2000—26 consecutive years. Two carriers account for the top position in all but 2 years of the entire dataset.

- Rows returned: 4
- Columns: leader, era_start_year, era_end_year, years_in_era

| leader | era_start_year | era_end_year | years_in_era |
| --- | --- | --- | --- |
| DL | 1987 | 1989 | 3 |
