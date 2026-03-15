Create `visual.html` using the `ontime-analyst-dashboard` skill in dynamic mode.

Dynamic-mode requirements:

- The dashboard must fetch data live in the browser from MCP OpenAPI using browser-stored JWE.
- Do not embed CSV data, embedded `result.json`, or any hardcoded metric values.
- Include the dynamic-mode auth and query controls required by the skill.
- Place the JWE token input, SQL textarea, fetch button, forget button, and status text together in a separate control block at the very end of the page.
- Load any previously stored JWE token from browser `localStorage` on page init so subsequent dashboards can reuse it.
- Persist the current JWE token back to `localStorage` after successful entry.
- Keep the control block at the bottom of the page after data loads; do not leave it in the empty-state or hero area.
- Normalize browser-fetched `columns` + `rows` into row objects before deriving KPIs, charts, or tables.
- Handle empty results and malformed payloads with visible warnings.
- Treat `count = 0` with `rows = null` as an empty result, not as a malformed payload.
- The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.
- Prefill the dashboard SQL textarea with the saved query shown below.

Saved SQL to embed directly in the final page:

```sql
{{saved_sql}}
```

Visual context:

- Question title: `{{question_title}}`
- Visual type: `{{visual_type}}`

Question-specific visual guidance:

{{visual_prompt_md}}
