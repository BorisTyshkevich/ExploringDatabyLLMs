
- Do not emit result rows or any data output.
- Write one JSON object containing the final verified SQL and a Markdown report template to `answer.raw.json`, not the debug query.

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
- The report must directly answer the business questions asked in the question-specific prompt.
- Do not make the report primarily about how to read the table, result shape, or row types.
- Do not replace analysis with instructions to the reader such as "use these rows to determine...".
- If the question-specific prompt asks multiple business questions, answer each of them explicitly in prose.
- Do not invent a custom scoring formula, analysis window, ranking rule, or business definition unless the question-specific prompt explicitly asks for it.
- If a reasonable guardrail or assumption is necessary to remove noise, keep it minimal and make it consistent with the question's business framing.
- Use the full available dataset history unless the question-specific prompt explicitly asks for a narrower time window.
- Allowed built-in placeholders: {{report_placeholders}}
- Allowed metric placeholders use this pattern only: `{{metric.<name>}}`
- Do not invent any placeholder outside the built-in list and `{{metric.<name>}}`.
- If a fact is needed in the report and is not covered by a built-in placeholder, put it in `metrics.named_values` and reference it via `{{metric.<name>}}`.
- Do not include TSV, JSON rows, HTML, or any other fenced blocks.

Invalid example:

`"report_markdown": "The key derived value is {{primary_value}}."`

Use `metrics.named_values.primary_value` and `{{metric.primary_value}}` .

Question title: `{{question_title}}`

Question-specific guidance:

{{question_prompt_md}}
