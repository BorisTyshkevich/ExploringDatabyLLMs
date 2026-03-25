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

- Provide one proof query for each parsed section question.
- Return one object in `subquestions` for every parsed section id.
- Preserve each required `id` exactly.
- Each `answer_markdown` must directly answer that section's question in concise prose.
- All prose claims about a specific route's recurrence frequency or date range must be directly traceable to a row in an executed query result.
- Each `sql` must be one executable SQL statement only and should serve as the proof query for that section.
- Use one proof query per section question. Do not merge several section questions into one unioned or row-typed SQL result unless the question-specific guidance explicitly requires that.
- Derive SQL and answers only from the current question and the current query result shape.
- Use the full available dataset history unless the question-specific prompt explicitly asks for a narrower time window.
- For additional section questions, do not hard-code any numeric threshold derived from another section result (e.g., hop count). Filter or join dynamically instead.

## Do not

- Do not invent extra section ids, custom scoring formulas, analysis windows, ranking rules, or business definitions unless the question-specific prompt explicitly asks for them.
- Do not emit result rows or any data output.
- Do not output fenced Markdown artifacts.

## Question title 
`{{question_title}}`

## Question-specific guidance

Read every top-level `###` section in the question-specific guidance below.

{{question_prompt_md}}
