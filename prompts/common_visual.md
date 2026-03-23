Create browser-ready HTML `visual.html` using the proper  `*-analyst-dashboard` skill.

The returned `visual.html` must be final browser-ready HTML. qforge will not patch or rewrite it after generation.

Use file `query.sql` as the authoritative input for the primary data query.
Use file `visual_input.json` to understand the result shape before building visuals without querying dataset.
You may construct additional queries when needed for enrichment or drill-down, but do not rebuild the saved SQL.

Return exactly this fenced section:

```html
<!doctype html>
<html>...</html>
```

### Rules

- Question title: `{{question_title}}`
- Visual mode: `{{visual_mode}}`
- Visual type: `{{visual_type}}`
- Derive KPIs, chart values, table rows, filters, and highlights from the actual analytical data. Do not invent or hardcode them.
- Respect the declared visual mode and visual type shown below.
- Follow question-specific visual guidance after the shared contract. Put reusable runtime behavior in shared page code, not in prose comments.

{{visual_prompt_md}}

### Data Source

SQL query for primary data source:

```sql
{{saved_sql}}
```

Data example/snippet:

### Visual input requirements

{{visual_input_summary_json}}
