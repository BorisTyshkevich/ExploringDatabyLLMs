Do not emit result rows or any data output.
Inspect the live schema first with `SHOW TABLES FROM ontime` and `DESCRIBE TABLE` for the tables you intend to use.
Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in source data reading subquery/CTE, and fix any errors internally.
Return the final verified SQL and a Markdown report template, not the debug query.

Return exactly these fenced sections:

```sql
-- one SQL statement
```

```report
Use placeholders only where data is needed.
Allowed placeholders: {{report_placeholders}}
```

Rules:

- Use one SQL statement only.
- Emit only the final verified SQL.
- The report must be Markdown.
- The report must be a template, not a data-filled summary.
- Prefer `{{data_overview_md}}` and `{{result_table_md}}` for JSON-derived sections.
- Keep the report concise and analytical.
- Do not include TSV, JSON rows, HTML, or any other fenced blocks.

Question title: `{{question_title}}`

Question-specific report guidance:

{{report_prompt_md}}
