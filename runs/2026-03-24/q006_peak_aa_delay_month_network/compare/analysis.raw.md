The `compare_report.md` has been written. Here's a brief summary of the key findings:

- **All 3 runs agree** on July 2024 as the worst AA month (36.33 min avg delay, 38.04% del15, 86,083 flights)
- **Ranking divergence**: claude/sonnet ranks contributors by average delay (CLT #1), codex ranks by total delay minutes (DFW #1) — both valid but highlighting different perspectives
- **The critical bug**: codex run-001 computed concentration shares against a filtered denominator (HAVING >= 1000), inflating the top-2 origin share to 50.8% vs the correct ~35% against the full network. The automated review caught this as a FAIL.
- **run-002 fixed it**: codex's retry computed all shares against the unfiltered network, passed review, and produced the most self-documenting output with city names and explicit share columns
- **Performance**: claude/sonnet generated SQL in ~83s vs ~188s for codex, with total durations of 381s vs 674s for complete runs
