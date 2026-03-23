Create browser-ready HTML `visual.html` using the proper  `*-analyst-dashboard` skill.

Write the file or provide a download link. Do not include the HTML source in the response. Do not open the artifact view frame.

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

{{visual_input_summary_json}}

