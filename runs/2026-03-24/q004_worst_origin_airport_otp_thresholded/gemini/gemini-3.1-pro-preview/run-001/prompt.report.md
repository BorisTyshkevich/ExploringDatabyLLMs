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

Question title: `Worst origin airports by departure on-time performance`

Question-specific guidance:

Identify which origin airports have the worst departure on-time performance after excluding low-volume airports.

Analyze completed departures at the origin-airport level across the full available history. Focus on airports with enough traffic to make the comparison meaningful, and rank the weakest performers by departure on-time performance.

You may apply a reasonable minimum-volume filter to remove noise, but do not invent a custom score or redefine on-time performance.

For each airport, quantify:

- completed departures
- departure on-time performance
- average departure delay
- a high-delay measure that reflects the worse end of the delay distribution
- the first and last dates represented in the data

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked view of the weakest qualifying origin airports
- a comparison view that makes the gap between the very worst airports and the middle of the ranked set easy to judge

The proof query behind the worst-airport question should preserve the ranked airport-level rows needed for the dashboard, not just a single top airport or summary statistic.

## Dashboard Questions

- Which airport ranks worst on departure on-time performance?
- How large is the spread between the worst airport and the middle of the ranked set?
- Are the weakest airports mostly major hubs, or is the bottom group more mixed?

In the report, answer those questions directly in prose. Name the worst airport, describe the spread between the bottom and the middle using the verified result, and summarize whether the weakest group is mostly hubs or more mixed.

Do not use fallback phrases such as "the worst airport" or "the bottom group" when your verified query results let you name the actual airport set directly.

Keep the result business-readable and analytically sound. Exclude low-volume airports before ranking them.