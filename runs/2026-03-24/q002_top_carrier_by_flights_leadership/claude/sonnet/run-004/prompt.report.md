- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in a data reading subquery or CTE. Fix any errors in a loop until done.

- Do not emit result rows or any data output.
- Write one JSON object to `answer.raw.json`.

Write exactly this JSON object shape to `answer.raw.json`:

{
  "subquestions": [
    {
      "subquestion": "Question text copied exactly from the required subquestion contract.",
      "answer_markdown": "A concise prose answer to that subquestion.",
      "sql": "-- one SQL statement that proves the answer"
    }
  ]
}

Rules:

- Write the artifact to `answer.raw.json`.
- The `answer.raw.json` file must contain raw JSON, not fenced Markdown.
- Read every bullet under `## Dashboard Questions` in the question-specific guidance below.
- Return one object in `subquestions` for every listed dashboard question, in the same order.
- Preserve the required `subquestion` text exactly.
- Each `answer_markdown` must directly answer that subquestion in concise prose.
- Each `sql` must be one executable SQL statement only and should serve as the proof query for that subquestion.
- Use one proof query per dashboard question. Do not merge several dashboard questions into one unioned or row-typed SQL result unless the question-specific guidance explicitly requires that.
- Derive SQL and answers only from the current question and the current query result shape.
- Do not rely on prior qforge runs, prior question ids, or previously observed values.
- Do not invent extra dashboard questions, custom scoring formulas, analysis windows, ranking rules, or business definitions unless the question-specific prompt explicitly asks for them.
- Use the full available dataset history unless the question-specific prompt explicitly asks for a narrower time window.

Question title: `Yearly carrier leadership by completed flights`

Question-specific guidance:

Determine which carrier led the industry in completed flights each calendar year, and identify where leadership changed most sharply.

Analyze completed flights by year and carrier across the full available history. For each year, show the leading carriers, each carrier's share of total completed flights, and the gap between the leader and the runner-up.

Pay special attention to true leadership transitions, where the top carrier changes from one year to the next. The result should make it easy to see both long periods of stable dominance and the years when leadership shifted most dramatically.

Do not invent a custom scoring formula or a narrower analysis window. Rank and compare leadership changes directly from the yearly leadership metrics needed to answer the question.

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- an annual leadership view showing the top carriers by year
- a transition view highlighting years when the leading carrier changed

The proof query behind the leading-carrier question should preserve the annual carrier-by-year leadership rows needed for the dashboard, not just a final overall count of years led.

## Dashboard Questions

- Which carrier leads most often across the full time range?
- When leadership changes, how large is the swing versus the prior leader?
- Which transition is the sharpest?
- Does the market show long stable eras, or frequent turnover at the top?

In the report, answer those questions directly in prose. Name the carrier that leads most often, identify the sharpest leadership change using the verified result, and summarize whether the market is defined more by long stable eras or by frequent turnover.

Do not use fallback phrases such as "the leading carrier" or "the sharpest transition" when your verified query results let you name the actual carrier and year directly.

Keep the result business-readable and analytically sound.