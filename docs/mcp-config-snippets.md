# MCP Local Config Snippets

## Codex (`~/.codex/config.toml`)

```toml
[mcp_servers.altinity_ontime_demo]
url = "https://mcp.demo.altinity.cloud/<JWE_TOKEN>/http"
```

## Gemini CLI (`~/.gemini/settings.json`)

```json
{
  "mcpServers": {
    "demo-altinity": {
      "httpUrl": "https://mcp.demo.altinity.cloud/<JWE_TOKEN>/http",
      "timeout": 300000
    }
  }
}
```

## Security Note

For public articles, never publish reusable production tokens. Use redacted placeholders and short-lived examples.
