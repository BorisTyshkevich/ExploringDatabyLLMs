- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

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
# Highest daily hops for one aircraft on one flight number

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
- Allowed built-in placeholders: {{row_count}}, {{generated_at}}, {{columns_csv}}, {{question_title}}, {{data_overview_md}}, {{result_table_md}}
- Do not use `{{metric.<name>}}` placeholders in this mode.
- Do not invent any placeholder outside the built-in list.
- Do not write `answer.raw.json`.
- Do not include TSV, JSON rows, HTML, or any other fenced blocks outside the example file sections above.

Question title: `Highest daily hops for one aircraft on one flight number`

Question-specific guidance:

Find routes with the highest number of hops per day for a single aircraft using the same flight number.

- What does the itinerary look like? 
- Find the top 10 longest and most recent itineraries.
- Report Aircraft ID, Flight Number, Carrier, Date, Route, hop_count.

rules:

- return column names exactly as requested. Do not invent other column aliases.
- Build each itinerary in chronological leg order using the actual departure timestamps.
- The textual `Route` must include every leg and the final destination airport.
- Delimiter in `Route` must be `-` without spaces.