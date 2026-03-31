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

Question title: `American Airlines peak network delay month and contributors`

Question-specific guidance:

Find American Airlines' worst network-wide month for departure delays, then identify which origins and routes contributed most to that peak.

Analyze completed American Airlines flights by month across the full network. Find the single month that stands out as the worst overall for departure delays.

Use the full available history unless the question explicitly asks for a narrower period. Do not invent a custom score for the peak month; identify it directly from the monthly delay metrics needed to answer the question.

For the monthly view, quantify:

- flight volume
- average departure delay
- the share of flights departing 15+ minutes late

Then drill into the selected peak month to show which origin airports and origin-destination routes contributed most to that bad month. Focus on contributors with enough flights in that month to be meaningful.

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a monthly leaderboard showing how the network performed over time
- a drilldown into the selected peak month by origin and by route

The proof query behind the peak-month question should preserve the month-by-month network rows needed for the dashboard, not just the single worst month.

## Dashboard Questions

- Which month is the single worst American Airlines month for departure delays?
- Which origins contribute most to that peak month?
- Which routes contribute most to that peak month?
- Does the peak look broad across the network, or concentrated in a smaller set of origins and routes?

In the report, answer those questions directly in prose. Name the worst month, identify the leading origin and route contributors using the verified result, and summarize whether the peak looks broad or concentrated.

Do not use fallback phrases such as "the peak month" or "the leading contributors" when your verified query results let you name the actual month, origins, and routes directly.

Keep the result business-readable and analytically sound.