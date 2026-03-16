Create `visual.html` using the `ontime-analyst-dashboard` skill in dynamic mode.

Dynamic-mode requirements:

- The dashboard must execute the saved SQL in the browser via MCP OpenAPI using browser-stored JWE as the primary analytical query.
- Additional browser queries are allowed for explicit enrichment or drill-down purposes when they materially improve the visualization.
- Do not invent hidden enrichment requests, embedded analytical CSV, embedded `result.json`, or hardcoded metric values.
- Include the auth/query controls required by the skill and keep them in a separate control block at the very end of the page.
- Load any previously stored JWE token from browser `localStorage` on page init and persist the current JWE token after successful entry.
- Normalize browser-fetched `columns` + `rows` into row objects before deriving KPIs, charts, tables, filtering, formatting, and highlights.
- Handle empty results and malformed payloads with visible warnings, and treat `count = 0` with `rows = null` as an empty result.
- The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.
- Prefill the dashboard SQL textarea with the saved query shown below and treat it as the primary query source.
- If additional queries are used, show them in a visible query ledger with purpose, role, status, and row count where available.
- Separate “primary query” and “enrichment queries” clearly in both code and UI.
- Prefer dataset-native tables rather than inferred or geocoded coordinates when enrichment is needed.
- If the page uses HTML templates for cards, avoid duplicated fixed `id` values inside cloned template content; use scoped selectors, `data-role` hooks, or unique IDs so JS always binds to the visible rendered card.
- For Leaflet maps, prefer initializing the map only after the visible container is in layout; if initialization happens before reveal, call `invalidateSize()` after the container becomes visible.

Saved SQL to embed directly in the final page:

```sql
{{saved_sql}}
```

Visual context:

- Question title: `{{question_title}}`
- Visual type: `{{visual_type}}`

Question-specific visual guidance:

{{visual_prompt_md}}
