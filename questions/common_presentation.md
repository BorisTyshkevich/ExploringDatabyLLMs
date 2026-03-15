You are generating presentation artifacts inside qforge.

The harness already executed the saved SQL and produced `result.json`.
Do not invent KPIs or data values.
The report is a Markdown template that qforge will render from `result.json`.

Return exactly these fenced sections:

```report
Use placeholders only where data is needed.
Allowed placeholders: {{report_placeholders}}
```

```html
<!doctype html>
<html>...</html>
```

Report rules:

- The report must be Markdown.
- The report must be a template, not a data-filled summary.
- Prefer `{{data_overview_md}}` and `{{result_table_md}}` for JSON-derived sections.
- Keep the report concise and analytical.

Question title: `{{question_title}}`

Result columns: `{{result_columns_csv}}`

Question-specific report guidance:

{{report_prompt_md}}
