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
      "id": "example_id",
      "subquestion": "Question text copied exactly from the required subquestion contract.",
      "answer_markdown": "A concise prose answer to that subquestion.",
      "sql": "-- one SQL statement that proves the answer"
    }
  ]
}

Rules:

- Write the artifact to `answer.raw.json`.
- The `answer.raw.json` file must contain raw JSON, not fenced Markdown.
- Return one object in `subquestions` for every required subquestion listed below.
- Preserve the required subquestion `id` values exactly.
- Preserve the required `subquestion` text exactly.
- Each `answer_markdown` must directly answer that subquestion in concise prose.
- Each `sql` must be one executable SQL statement only and should serve as the proof query for that subquestion.
- Use one proof query per required subquestion. Do not merge several subquestions into one unioned or row-typed SQL result unless the question-specific guidance explicitly requires that.
- Derive SQL and answers only from the current question and the current query result shape.
- Do not rely on prior qforge runs, prior question ids, or previously observed values.
- Do not invent extra required subquestions, custom scoring formulas, analysis windows, ranking rules, or business definitions unless the question-specific prompt explicitly asks for them.
- Use the full available dataset history unless the question-specific prompt explicitly asks for a narrower time window.

Required subquestion contract:

- `worst_hotspot`: Which destination and time block is the worst recurring hotspot?
- `persistence`: Is that hotspot consistently bad across time, or concentrated in a narrower period?
- `pattern_summary`: What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?

Question title: `Delta ATL departure delay hotspots by destination and time block`

Question-specific guidance:

Find the Delta departure-delay hotspots out of ATL that appear to be persistently problematic, not just noisy one-off periods.

Analyze completed Delta departures from ATL by destination and departure time block across the full available history. Focus on combinations that have enough flight volume to be credible and enough repeated monthly presence to count as sustained hotspots.

You may apply reasonable minimum-volume filters to remove noise, but do not invent a custom hotspot score or a narrower analysis window.

For each hotspot, quantify:

- flight volume
- average departure delay
- a high-delay measure that captures the worse end of the distribution
- the share of flights departing 15+ minutes late
- how many months the hotspot meaningfully appears

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked hotspot summary view
- a monthly trend view for the leading hotspots

The monthly trend view should include only months that are credible for interpretation, not thin low-volume months that would add noise.

The output should let a BI dashboard answer:

- Which destination and time block is the worst recurring hotspot?
- Is that hotspot consistently bad across time, or concentrated in a narrower period?
- What do the top hotspots suggest about where Delta faces the most departure-pressure out of ATL?

In the report, answer those questions directly in prose before the table. Name the actual worst hotspot using the verified result, state whether it looks persistent or concentrated, and summarize the pattern across the leading hotspots.

Do not use fallback phrases such as "the top-ranked hotspot" or "add one short takeaway" when your verified query results let you name the destination, departure time block, and key pattern directly.

Keep the result business-readable and analytically sound. Exclude low-volume noise before identifying the leading hotspots.