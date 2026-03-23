
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

Brief findings that directly answer the question.

{{data_overview_md}}

{{result_table_md}}
```

Rules:

- Use one SQL statement only in `query.sql`.
- `query.sql` must contain executable SQL only, not Markdown fences.
- `report.template.md` must contain Markdown only, not fenced Markdown.
- `report.template.md` may include concrete findings derived from your own verified query results when the question asks for a direct answer.
- Prefer `{{data_overview_md}}` and `{{result_table_md}}` for JSON-derived sections.
- Keep the report concise and analytical.
- Use placeholders only where data is needed.
- Derive SQL and report claims only from the current question and the current query result shape.
- Do not rely on prior qforge runs, prior question ids, or previously observed values.
- Do not mention other question ids such as `q001` in the artifact unless the current prompt explicitly asks for cross-question comparison.
- The report must directly answer the business questions asked in the question-specific prompt.
- Do not make the report primarily about how to read the table, result shape, or row types.
- Do not replace analysis with instructions to the reader such as "use these rows to determine...".
- If the question-specific prompt asks multiple business questions, answer each of them explicitly in prose.
- Do not invent a custom scoring formula, analysis window, ranking rule, or business definition unless the question-specific prompt explicitly asks for it.
- If a reasonable guardrail or assumption is necessary to remove noise, keep it minimal and make it consistent with the question's business framing.
- Use the full available dataset history unless the question-specific prompt explicitly asks for a narrower time window.
- Do not leave report prose as generic placeholders such as "top-ranked hotspot" or "add one takeaway" when your verified query result lets you name the actual entity or finding.
- When the question asks for the top item, a transition, a worst month, or a worst hotspot, name it directly in the report using the verified result.
- Allowed built-in placeholders: {{report_placeholders}}
- Do not invent any placeholder outside the built-in list.
- Do not include TSV, JSON rows, HTML, or any other fenced blocks outside the example file sections above.

Question title: `{{question_title}}`

Question-specific guidance:

{{question_prompt_md}}
