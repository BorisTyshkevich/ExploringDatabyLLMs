Create `visual.html` using the `ontime-analyst-dashboard` skill.

The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.

General visual rules:

- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.
- Keep report logic out of `visual.html`; qforge already renders `report.md` separately.

Saved SQL to preserve in the final page contract:

```sql
{{saved_sql}}
```

Visual context:

- Question title: `{{question_title}}`
- Visual mode: `{{visual_mode}}`
- Visual type: `{{visual_type}}`

Question-specific visual guidance:

{{visual_prompt_md}}
