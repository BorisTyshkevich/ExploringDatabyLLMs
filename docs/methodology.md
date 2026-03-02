# Reproducible LLM SQL Evaluation Method

## Goal

Measure how reliably each model converts natural language prompts into correct, executable, and efficient ClickHouse SQL.

## Test Matrix

- Models: GPT, Claude, and at least two open-weight models.
- Environments:
  - Altinity demo: `default.ontime`
  - ClickHouse Inc demo: `ontime.ontime`
- Cases:
  - Basic aggregation
  - Multi-filter query
  - Group + ranking
  - Seasonal logic
  - Time series trend
  - Route-level analysis

## Scoring (0-5 per dimension)

- Correctness: query answers the business question.
- Executability: runs without manual fixes.
- Schema alignment: uses existing tables/columns.
- Efficiency: avoids unnecessary scans/functions.
- Robustness: still works when prompt wording changes.

Total score per case: `25`.

## Execution Rules

- Run each generated SQL unchanged first.
- Capture errors and classify (`UNKNOWN_TABLE`, `UNKNOWN_IDENTIFIER`, filter omissions, logic error).
- If corrected manually, store corrected SQL separately and mark as `human-assisted`.
- Record query output rows and execution time.

## Output Artifacts

- JSON results in `results/`.
- A table per case with: generated SQL, status, error class, corrected SQL (if any), score.
- Final leaderboard across models and both environments.
