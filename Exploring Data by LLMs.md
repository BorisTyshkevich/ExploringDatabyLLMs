# Exploring Data with LLMs + ClickHouse (Part 1)

How to use Altinity-MCP as a secure control plane between ChatGPT/Claude and ClickHouse.

## TL;DR

- LLMs are excellent for exploratory analytics, but only when they are connected through a constrained tool layer.
- MCP is that layer: `Chat UI -> LLM -> MCP server -> ClickHouse`.
- For production, the key controls are: read-only database users, short-lived tokens, HTTPS, and query guardrails.
- This article uses public demo environments so everything is reproducible.
- In Part 2 we benchmark real model-generated SQL on the Ontime dataset and compare OpenAI vs Anthropic models.

## Why LLMs for analytics at all?

BI dashboards are great for known metrics. They are weaker at fast, conversational exploration of unknown questions.

Typical friction with classic BI workflows:

1. You need predefined semantic models before non-technical users can self-serve.
2. Ad-hoc deep dives often fall back to SQL anyway.
3. Iteration cycles are slower when each new question requires dashboard or model updates.

LLMs help by reducing interface friction: users can ask a question in plain language, iterate quickly, and get both SQL and explanation.

That said, LLM-generated SQL without governance is risky. The architecture matters more than the model name.

## The architecture that actually works

Use MCP as the control and policy layer:

```text
Chat UI <-> LLM <-> MCP Server <-> ClickHouse
```

Why this pattern is important:

- The model gets a bounded tool surface (`list tables`, `describe table`, `execute query`).
- The MCP server enforces auth, limits, transport security, and logging.
- Database credentials stay behind the MCP boundary.

## Demo environments used in this series

Tested on **March 2, 2026**.

### ClickHouse demos

- Altinity demo: `demo.demo.altinity.cloud:9440` (`demo/demo`), table `default.ontime`
- ClickHouse Inc demo: `sql-clickhouse.clickhouse.com:9440` (`demo/<empty password>`), table `ontime.ontime`

Important caveat for reproducibility:

- These demos are not schema-identical and not time-range-identical.
- On March 2, 2026:
  - Altinity `default.ontime`: `1987..2021`, `201,575,308` rows
  - ClickHouse Inc `ontime.ontime`: `1987..2025`, `226,328,393` rows

If you compare model quality, use one environment per benchmark run or normalize schemas first.

## MCP transport choices

For shared team usage, prefer HTTP-based MCP deployment.

- `stdio`: simplest for local-only workflows.
- `HTTP/streamable HTTP`: best default for centralized/shared deployments.
- `SSE`: compatibility option for clients that still require it.

References:

- [Model Context Protocol transports](https://modelcontextprotocol.io/docs/concepts/transports)
- [Altinity MCP repository](https://github.com/Altinity/altinity-mcp)
- [OpenAI remote MCP tools guide](https://platform.openai.com/docs/guides/tools-remote-mcp)
- [Anthropic remote MCP announcement](https://www.anthropic.com/news/claude-code-remote-mcp)

## Deploying Altinity-MCP quickly (HTTP mode)

For local/dev reproducibility, Docker Compose is enough.

```yaml
services:
  altinity-mcp-local:
    image: ghcr.io/altinity/altinity-mcp:latest
    environment:
      MCP_TRANSPORT: http
      MCP_PORT: 8080
      MCP_OPENAPI: http
      CLICKHOUSE_HOST: demo.demo.altinity.cloud
      CLICKHOUSE_PORT: 8443
      CLICKHOUSE_PROTOCOL: http
      CLICKHOUSE_TLS: "true"
      CLICKHOUSE_DATABASE: default
      CLICKHOUSE_USERNAME: demo
      CLICKHOUSE_PASSWORD: demo
      CLICKHOUSE_READ_ONLY: "true"
    ports:
      - "8080:8080"
```

Then verify:

```bash
curl -s http://localhost:8080/health | jq .
```

For production, use Helm/Kubernetes and keep secrets in your secret manager.

## Token strategy (what to publish, what never to publish)

Use encrypted tokens (JWE/JWT) or bearer auth through your MCP layer, but do not publish reusable long-lived tokens in public docs.

### Good practice

- Short TTL tokens.
- Rotate regularly.
- Send token via Authorization header where possible.
- Log redaction for URLs and headers.

### Bad practice

- Hardcoding non-expiring tokens in article text, screenshots, public repos, or sample config.

When showing setup in Medium, always use placeholders:

```text
https://mcp.example.com/<JWE_TOKEN>/http
```

## Client integration patterns

### Chat clients

Use the client’s MCP connector UI and point it at your MCP HTTP endpoint.

### Codex CLI

Add to `~/.codex/config.toml`:

```toml
[mcp_servers.analytics]
url = "https://mcp.example.com/<JWE_TOKEN>/http"
```

### Claude Code CLI

```bash
claude mcp add --transport http analytics https://mcp.example.com/<JWE_TOKEN>/http
```

### Gemini CLI

`~/.gemini/settings.json`:

```json
{
  "mcpServers": {
    "analytics": {
      "httpUrl": "https://mcp.example.com/<JWE_TOKEN>/http",
      "timeout": 300000
    }
  }
}
```

Capability support changes quickly. Before publishing vendor-specific statements, verify current docs on the publication date.

## Security checklist for production

1. Read-only ClickHouse user for exploration workloads.
2. Query limits (`max_execution_time`, row/byte caps, and explicit `LIMIT` policy).
3. Strict allowlist for exposed MCP tools.
4. HTTPS everywhere (client <-> MCP and MCP <-> ClickHouse).
5. Token TTL + rotation + revocation process.
6. Auditing: who asked what, which SQL ran, and how much data was scanned.
7. Network controls: private networking/VPN where possible.
8. Human review gates for high-risk queries.

## What this architecture does not solve by itself

- Hallucinated but syntactically valid SQL.
- Silent semantic errors that return plausible but wrong numbers.
- Cost/performance issues from inefficient generated queries.

You still need validation loops and benchmarked model behavior.

## Reproducibility package used for this series

In this project folder:

- `sql/` contains tested queries for both demo environments.
- `scripts/run_all.sh` validates ClickHouse and MCP connectivity.
- `docker-compose.yml` runs Altinity-MCP locally in HTTP mode.

This keeps the article evidence-based instead of screenshot-based.

## Part 2 preview

In Part 2, we run a benchmark on model-generated SQL against Ontime using:

- OpenAI: `gpt-5`, `gpt-5.3-codex`
- Anthropic: `sonnet`, `opus`

We compare executability, strict output matching, and relaxed business-answer matching, with all SQL and outputs captured in `results/`.
