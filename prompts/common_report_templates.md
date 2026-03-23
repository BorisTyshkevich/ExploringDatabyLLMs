
- Do not emit result rows or any data output.
- Write the final verified SQL to `query.sql`.
- Write the Markdown report template to `report.template.md`.

Write exactly these two files in the run directory:

`query.sql`

```sql
-- one SQL statement
```

`report.template.md`

```md
# {{question_title}}

{{data_overview_md}}

Add one short analytical takeaway grounded in the result set.
```

Rules:

- Use one SQL statement only in `query.sql`.
- `query.sql` must contain executable SQL only, not Markdown fences.
- `report.template.md` must contain Markdown only, not fenced Markdown.
- The report must be a template, not a data-filled summary.
- Prefer `{{data_overview_md}}` and `{{result_table_md}}` for JSON-derived sections.
- Keep the report concise and analytical.
- Use placeholders only where data is needed.
- Derive SQL and report claims only from the current question and the current query result shape.
- Do not rely on prior qforge runs, prior question ids, or previously observed values.
- Do not mention other question ids such as `q001` in the artifact unless the current prompt explicitly asks for cross-question comparison.
- Allowed built-in placeholders: {{report_placeholders}}
- Do not use `{{metric.<name>}}` placeholders in this mode.
- Do not invent any placeholder outside the built-in list.
- Do not write `answer.raw.json`.
- Do not include TSV, JSON rows, HTML, or any other fenced blocks outside the example file sections above.

Question title: `{{question_title}}`

Question-specific guidance:

{{question_prompt_md}}
