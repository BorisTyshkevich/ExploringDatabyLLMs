Do not emit result rows or any data output.
Inspect the live schema first with `SHOW TABLES FROM ontime` and `DESCRIBE TABLE` for the tables you intend to use.
Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in source data reading subquery/CTE, and fix any errors internally.
Write one JSON object containing the final verified SQL and a Markdown report template to `answer.raw.json`, not the debug query.

`answer.raw.json` must contain plain JSON bytes only. Do not wrap the file contents in Markdown fences.

Your stdout response may contain a short status line, but qforge will ignore stdout and load only `answer.raw.json`.

Write exactly this JSON object shape to `answer.raw.json`:

{
  "sql": "-- one SQL statement",
  "report_markdown": "# {{question_title}}\\n\\n{{data_overview_md}}\\n\\nThe key derived value is {{metric.primary_value}}.",
  "metrics": {
    "summary_facts": [
      "Summarize the strongest derived fact from the query result."
    ],
    "named_values": {
      "primary_value": "example derived value"
    },
    "named_lists": {
      "example_list": [
        "Example ordered item"
      ]
    }
  }
}

Rules:

- Use one SQL statement only.
- JSON must contain exactly these top-level keys:
  - `sql`
  - `report_markdown`
  - `metrics`
- Write the artifact to `answer.raw.json`.
- The `answer.raw.json` file must contain raw JSON, not fenced Markdown.
- The report must be Markdown.
- The report must be a template, not a data-filled summary.
- Prefer `{{data_overview_md}}` and `{{result_table_md}}` for JSON-derived sections.
- Keep the report concise and analytical.
- Use placeholders only where data is needed.
- Derive SQL, metrics, and report claims only from the current question and the current query result shape.
- Do not rely on prior qforge runs, prior question ids, or previously observed values.
- Do not mention other question ids such as `q001` in the artifact unless the current prompt explicitly asks for cross-question comparison.
- Allowed built-in placeholders: {{report_placeholders}}
- Allowed metric placeholders use this pattern only: `{{metric.<name>}}`
- Do not invent any placeholder outside the built-in list and `{{metric.<name>}}`.
- If a fact is needed in the report and is not covered by a built-in placeholder, put it in `metrics.named_values` and reference it via `{{metric.<name>}}`.
- Do not include TSV, JSON rows, HTML, or any other fenced blocks.

Invalid example:

`"report_markdown": "The key derived value is {{primary_value}}."`

The placeholder `{{primary_value}}` is invalid. Use `metrics.named_values.primary_value` plus `{{metric.primary_value}}` instead.

Question title: `{{question_title}}`

Question-specific report guidance:

{{report_prompt_md}}
