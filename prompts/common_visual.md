Create the presentation artifact using the proper `*-analyst-dashboard` skill.

### Rules

- Question title: `{{question_title}}`
- Visual mode: `{{visual_mode}}`
- Presentation target: `{{presentation_target}}`
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

{{visual_input_summary_json}}
