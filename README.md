# Ontime + LLMs Publication Kit

This project turns the two article drafts into a reproducible, testable package for Medium publication.

It includes:
- Validated SQL examples for both public ClickHouse demo environments.
- Scripts that execute all examples and store JSON outputs.
- MCP HTTP/OpenAPI smoke tests against the public Altinity MCP endpoint.
- Automated LLM benchmark harness (OpenAI + Anthropic batch CLIs) with reproducible scoring.
- A local `docker-compose` setup for running `altinity-mcp` in HTTP mode.
- Editorial critique and rewrite plan for both articles.

## Repository Layout

- `docs/review-critique.md`: publication critique, weak points, and improvements.
- `docs/medium-outline-part1.md`: Medium-ready structure for Part 1.
- `docs/medium-outline-part2.md`: Medium-ready structure for Part 2.
- `docs/methodology.md`: reproducible model benchmarking method.
- `sql/altinity/*.sql`: validated queries for `demo.demo.altinity.cloud` (`default.ontime`).
- `sql/clickhouse_inc/*.sql`: validated queries for `sql-clickhouse.clickhouse.com` (`ontime.ontime`).
- `sql/part2/*.sql`: as-written and corrected SQL used in the critique.
- `scripts/run_clickhouse_tests.sh`: runs SQL checks against both ClickHouse demos.
- `scripts/run_mcp_openapi_tests.sh`: runs MCP HTTP/OpenAPI checks.
- `scripts/run_all.sh`: executes all checks.
- `scripts/benchmark_llm_sql.py`: runs schema-aware SQL benchmark across 4 models.
- `scripts/summarize_benchmark_relaxed.py`: computes strict vs relaxed benchmark scoring.
- `docker-compose.yml`: local `altinity-mcp` service (HTTP transport).

## Quick Start

1. Copy environment file:

```bash
cp .env.example .env
```

2. Run all tests:

```bash
./scripts/run_all.sh
```

3. Inspect output artifacts:

- `results/clickhouse/*.json`
- `results/mcp/*.json`
- `results/benchmark/*`

Benchmark score outputs:

- `results/benchmark/benchmark_summary.json` (strict scoring)
- `results/benchmark/benchmark_relaxed_summary.json` (normalized business-answer scoring)

## Local Altinity MCP (HTTP) via Docker Compose

Start local MCP server:

```bash
docker compose up -d
```

Health check:

```bash
curl -s http://localhost:${LOCAL_MCP_PORT:-8080}/health | jq .
```

Stop:

```bash
docker compose down
```

## Notes

- The two public demo ClickHouse servers do not expose identical schema/database names for Ontime.
- The SQL in `sql/part2/01_*` and `sql/part2/02_*` intentionally reflects the current Part 2 article text and is expected to fail.
- Corrected versions are included and validated.
