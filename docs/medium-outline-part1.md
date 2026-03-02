# Part 1 Medium Outline

## Working title

From Prompt to Production: Building a Secure MCP Layer Between LLMs and ClickHouse

## Narrative Arc

1. Problem:
- BI is strong for known questions, weak for exploratory natural-language loops.

2. Architecture:
- `Chat UI -> LLM -> MCP Server -> ClickHouse`
- Why MCP is the safety/control plane.

3. Transport decision:
- Streamable HTTP first, SSE for compatibility, stdio for local dev.

4. Security model:
- Read-only DB users.
- Short-lived tokens.
- Bearer auth over URL tokens.
- Auditing/logging and rate limits.

5. Hands-on setup:
- Direct ClickHouse smoke checks.
- MCP health and query checks.
- Optional docker-compose local setup.

6. Production checklist:
- Auth, rotation, least privilege, data governance, observability.

7. Bridge to Part 2:
- Controlled benchmark of model SQL quality on Ontime.

## Must-fix content changes

- Remove or redact live long-lived token strings.
- Replace broad vendor support claims with sourced, date-stamped statements.
- Correct ChatGPT plan restrictions to current behavior with official links.
- Add explicit caveat that demo clusters differ in schema and data horizon.
