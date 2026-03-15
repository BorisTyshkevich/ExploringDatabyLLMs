# CLAUDE.md

## Project Overview

**qforge** is a Go-based benchmark harness for evaluating LLM-generated SQL against a real OnTime aviation dataset via an MCP (Model Context Protocol) OpenAPI interface. The core design principle is strict phase separation: models generate SQL only, the harness executes SQL directly, and comparison/validation happen after runs complete.

## Architecture

### Two-Phase Design

1. **SQL Phase** – LLM generates SQL; harness extracts and executes it against the MCP endpoint
2. **Presentation Phase** – LLM generates report/visual HTML from the canonical JSON result (deferred unless `--with-visual`)

Models never see execution results during SQL generation. The harness owns all data access.

### Directory Layout

```
cmd/qforge/        CLI entry point
internal/
  cli/             Command router and orchestration
  model/           Shared data structures and enums
  datasets/        Dataset config loading and MCP URL resolution
  questions/       Question metadata loading
  providers/       LLM provider adapters (Claude, Codex, Gemini)
  execute/         SQL execution against MCP OpenAPI
  extract/         Fenced code block extraction from LLM output
  prompts/         Prompt template building
  runs/            Run directory management and manifest handling
  render/          Report and table rendering
  compare/         Run comparison and reporting
  querylog/        system.query_log metric fetching
questions/         Per-question YAML + Markdown prompt files
datasets/          Dataset YAML configs (MCP endpoint, auth, tables)
Skills/            Claude skill definitions (HTML dashboard generation)
scripts/           Bash wrapper: ./scripts/qforge
```

### Run Directory Structure

`runs/YYYY-MM-DD/<question-slug>/<runner>/<model>/run-NNN/`

SQL-phase artifacts: `prompt.sql.md`, `answer.sql.raw.md`, `query.sql`, `result.json`, `manifest.json`, `stdout.log`, `stderr.log`

Presentation artifacts: `prompt.presentation.md`, `answer.presentation.raw.md`, `report.template.md`, `report.md`, `visual.html`

## Build & Run

```bash
# Run via wrapper (uses go run)
./scripts/qforge <subcommand> [flags]

# Or build directly
go build ./cmd/qforge
./qforge <subcommand> [flags]

# Run tests
go test ./...
```

## Commands

```bash
# List available questions
./scripts/qforge list-questions

# Generate SQL and execute (SQL phase only)
./scripts/qforge run --question q001 --runner claude --model claude-opus-4

# SQL phase + presentation in one step
./scripts/qforge run --question q001 --runner claude --model claude-opus-4 --with-visual

# Generate presentation for an existing run
./scripts/qforge process-visual --run runs/2026-01-15/q001_hops_per_day/claude/claude-opus-4/run-001

# Compare runs and fetch performance metrics
./scripts/qforge compare --question q001 --date 2026-01-15

# Inspect run manifest
./scripts/qforge inspect-run --run <run-dir>
```

All commands support `--verbose` for detailed logging.

## Environment

```bash
# Required: MCP authentication token
export MCP_JWE_TOKEN=<your-jwe-token>
```

## Dataset: ontime_v2

- **MCP endpoint:** `https://mcp.demo.altinity.cloud`
- **Primary table:** `default.ontime_v2`
- **Forbidden table:** `default.ontime` (legacy)
- **Auth:** Bearer JWE token via `MCP_JWE_TOKEN` env var

## LLM Providers

Three providers are supported, each invoked as a subprocess CLI:

| Provider | CLI binary | Notes |
|----------|------------|-------|
| Claude   | `claude`   | Uses MCP JSON config file |
| Codex    | `codex`    | Uses configuration flags |
| Gemini   | `gemini`   | Uses server name allowlist |

All providers must be installed and available on `PATH`.

## Questions

Six benchmark questions in `questions/`:

| ID   | Title | Visual Type |
|------|-------|-------------|
| q001 | Highest daily hops for one aircraft | `html_map` |
| q002 | Yearly carrier leadership | `html_timeseries` |
| q003 | Delta ATL departure delays | `html_heatmap` |
| q004 | Worst origin airport OTP | – |
| q005 | Worst winter carrier-origin pairs | `html_seasonal_dashboard` |
| q006 | American Airlines peak delay month | `html_contribution_dashboard` |

Each question has: `meta.yaml`, `prompt.md`, `visual_prompt.md`, optional `report_prompt.md`, optional `compare.yaml`. Question `q001` also includes reference outputs (`expected.sql`, `expected.tsv`, `expected_visual.html`).

## Key Design Decisions

- **Single external dependency:** `gopkg.in/yaml.v3`. Everything else is stdlib.
- **Canonical output is JSON.** Array-of-arrays from MCP is converted to array-of-objects by the harness.
- **Deferred metrics.** Query performance is fetched from `system.query_log` after execution, keyed by deterministic `log_comment` (SHA256 of query).
- **Validation by contract.** Optional `compare.yaml` per question defines comparison rules and tolerances.
- **Concurrent providers.** Multiple LLMs can run in parallel for the same question; results compared in a secondary pass.

## Skills

`Skills/ontime-analyst-dashboard/` defines a Claude skill for generating HTML dashboards. It supports:
- **Static mode:** Embedded CSV data (validator-safe)
- **Dynamic mode:** Live browser data fetching via MCP OpenAPI with browser-stored JWE
- **Map support:** Geographic visualization with Leaflet (`html_map`)
- **Edit mode:** Interactive layout controls

See `Skills/ontime-analyst-dashboard/SKILL.md` and `references/` for full documentation.
