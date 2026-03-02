# LLM SQL Benchmark Results

Generated at: 2026-03-02T18:02:35Z

## Aggregate Scores

| Model | Vendor | Execution OK | Correct | Total | Correct Rate | Avg Runtime (s) |
|---|---|---:|---:|---:|---:|---:|
| openai_gpt5 | openai | 6 | 1 | 6 | 16.67% | 1.322 |
| openai_gpt53codex | openai | 6 | 1 | 6 | 16.67% | 0.9647 |
| anthropic_sonnet | anthropic | 6 | 1 | 6 | 16.67% | 0.9415 |
| anthropic_opus | anthropic | 6 | 2 | 6 | 33.33% | 1.8451 |

## Per-Case Outcomes

| Model | Case | Status | Error |
|---|---|---|---|
| openai_gpt5 | case1_top_carrier_2019 | correct |  |
| openai_gpt5 | case2_avg_depdelay_delta_atl_2021_07 | runs_wrong |  |
| openai_gpt5 | case3_worst_airport_2021 | runs_wrong |  |
| openai_gpt5 | case4_worst_winter_carrier_airport | runs_wrong |  |
| openai_gpt5 | case5_peak_aa_month | runs_wrong |  |
| openai_gpt5 | case6_peak_route_season | runs_wrong |  |
| openai_gpt53codex | case1_top_carrier_2019 | correct |  |
| openai_gpt53codex | case2_avg_depdelay_delta_atl_2021_07 | runs_wrong |  |
| openai_gpt53codex | case3_worst_airport_2021 | runs_wrong |  |
| openai_gpt53codex | case4_worst_winter_carrier_airport | runs_wrong |  |
| openai_gpt53codex | case5_peak_aa_month | runs_wrong |  |
| openai_gpt53codex | case6_peak_route_season | runs_wrong |  |
| anthropic_sonnet | case1_top_carrier_2019 | correct |  |
| anthropic_sonnet | case2_avg_depdelay_delta_atl_2021_07 | runs_wrong |  |
| anthropic_sonnet | case3_worst_airport_2021 | runs_wrong |  |
| anthropic_sonnet | case4_worst_winter_carrier_airport | runs_wrong |  |
| anthropic_sonnet | case5_peak_aa_month | runs_wrong |  |
| anthropic_sonnet | case6_peak_route_season | runs_wrong |  |
| anthropic_opus | case1_top_carrier_2019 | correct |  |
| anthropic_opus | case2_avg_depdelay_delta_atl_2021_07 | runs_wrong |  |
| anthropic_opus | case3_worst_airport_2021 | runs_wrong |  |
| anthropic_opus | case4_worst_winter_carrier_airport | runs_wrong |  |
| anthropic_opus | case5_peak_aa_month | runs_wrong |  |
| anthropic_opus | case6_peak_route_season | correct |  |
