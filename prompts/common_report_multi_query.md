## Output

Write one JSON object to `answer.raw.json` file with shape:

{
  "subquestions": [
    {
      "id": "Section id such as main, q1, q2, or q0.",
      "answer_markdown": "A concise prose answer to that subquestion.",
      "sql": "-- one SQL statement that proves the answer"
    }
  ]
}

## Rules:

- Provide one proof sql query for each parsed section question.
- Return one object in `subquestions` for every parsed section id.
- Preserve each required `id` exactly.
- Each `answer_markdown` must directly answer that section's question in concise business-readable prose.
- All prose claims must be directly traceable to a row in an executed query result.
- When a conclusion depends on shares, concentration, or other denominator-based comparisons, keep the denominator-bearing totals inspectable in the proof query result.
- Use the most recent 5 years by default unless the question-specific prompt explicitly asks for a narrower or wider time window.

## Do not

- Do not merge several section questions into one unioned or row-typed SQL result unless the question-specific guidance explicitly requires that.
- Do not collapse a section to a final aggregate when the presentation needs the underlying ranked, time-series, or drilldown rows.
- Do not emit result rows or any data output.
- Do not output fenced Markdown artifacts.

## Question title 
`{{question_title}}`

## Question-specific guidance

Read every top-level `###` section in the question-specific guidance below.

{{question_prompt_md}}
