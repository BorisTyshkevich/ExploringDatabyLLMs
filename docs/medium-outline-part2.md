# Part 2 Medium Outline

## Working title

How Good Are LLMs at ClickHouse SQL? A Reproducible Ontime Benchmark

## Narrative Arc

1. Benchmark setup:
- Two public ClickHouse environments.
- Same prompts, temperature, and tool policy.
- Scoring rubric.

2. Baseline failures first:
- Show common failure classes (wrong columns, missing filters, wrong DB/table).

3. Case-by-case results:
- Case 1..6 with generated SQL, execution status, and corrected SQL.

4. Cross-model analysis:
- Accuracy vs latency vs verbosity.
- Where schema-aware MCP tools materially help.

5. Operational takeaway:
- Human-in-the-loop thresholds by query risk.
- Safe default guardrails for production.

6. Conclusion:
- What is production-ready now vs what still needs supervision.

## Must-fix content changes

- Fix SQL casing/schema mismatches in examples.
- Fix year-window mismatch (Altinity demo max year is 2021).
- Replace placeholder links (`chatgpt.com/c/link-to-part1`) with final public links.
- Add a reproducibility appendix (SQL files + scripts).
