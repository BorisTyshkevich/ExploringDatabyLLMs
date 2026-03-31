The report has been written to `compare_report.md`. Here's what it covers:

- **All 5 runs agree on the answer**: 8 hops max, same 10 Southwest Airlines itineraries, same order
- **Differences are cosmetic only**: column naming, route separators (`-` vs `→` vs `->`), departure time formatting
- **SQL strategies diverge significantly**: from 1-CTE single-pass (sonnet) to 5-CTE with window functions and airport JOINs (codex), with read volumes ranging 193M to 772M rows
- **Performance spread**: 7.3x on query time (4.37 s to 32.10 s), 4x on rows read
- **Presentation quality gap**: Claude reports contain genuine analysis; codex reports have unresolved template placeholders. Both codex runs failed visual rendering.
