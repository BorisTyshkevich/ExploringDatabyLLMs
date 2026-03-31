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

Question title: `Highest daily hops for one aircraft on one flight number`

Question-specific guidance:

Find the daily itineraries with the highest number of hops for a single aircraft using the same flight number.

Analyze completed flights across the full available history and focus on the unique top 10 longest itineraries, ordered so the most recent examples are easy to identify. Preserve enough itinerary detail for a dashboard map, route sequence panel, and clickable itinerary table.

Provide one proof query for each dashboard question below. The proof query behind the lead-itinerary question should preserve the ranked itinerary rows needed for the dashboard, not just a single top row.

Keep the proof queries efficient. Avoid rescanning the full fact table independently for every dashboard question when the same ranked itinerary candidate set can be reused or when a targeted follow-up on that candidate set is enough.

Return enough itinerary detail to support:

- aircraft id / tail number
- flight number
- carrier
- flight date
- hop count
- textual route in chronological leg order
- a simple browser-friendly itinerary-sequence field that lets the dashboard redraw the selected route and route detail panel without relying only on the textual `Route`

Keep these runtime constraints:

- Build each itinerary in chronological leg order using the actual departure timestamps.
- The textual `Route` must include every leg and the final destination airport.
- Delimiter in `Route` must be `-` without spaces.

In the answers, directly name the lead itinerary, directly name the most recent top itinerary, state whether the top itineraries look recurring or one-off, and summarize the geographic shape they trace.

## Dashboard Questions

- Which itinerary is the highest-hop example, and what does it look like?
- Which of the top-ranked itineraries is the most recent?
- Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?
- What geographic pattern do the top itineraries show?