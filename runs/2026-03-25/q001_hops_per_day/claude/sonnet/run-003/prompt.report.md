- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save.
- Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully.
- Do not write unchecked SQL.

- Write one JSON object to `answer.raw.json` file with shape:

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

- answer in prose to the main question.
- Provide proof query for the main question.
- Provide one proof query for each dashboard question.
- Read every bullet under `## Dashboard Questions` in the question-specific guidance below.
- Return one object in `subquestions` for every listed dashboard question, in the same order.
- Preserve the required `subquestion` text exactly.
- Each `answer_markdown` must directly answer that subquestion in concise prose.
- All prose claims about a specific route's recurrence frequency or date range must be directly traceable to a row in an executed query result.
- Each `sql` must be one executable SQL statement only and should serve as the proof query for that subquestion.
- Use one proof query per dashboard question. Do not merge several dashboard questions into one unioned or row-typed SQL result unless the question-specific guidance explicitly requires that.
- Derive SQL and answers only from the current question and the current query result shape.
- Use the full available dataset history unless the question-specific prompt explicitly asks for a narrower time window.
- For dashboard subquestions, do not hard-code any numeric threshold derived from the main result (e.g., hop count). Filter or join dynamically against the main result set instead.

## Do not

- Do not invent extra dashboard questions, custom scoring formulas, analysis windows, ranking rules, or business definitions unless the question-specific prompt explicitly asks for them.
- Do not emit result rows or any data output.
- Do not output fenced Markdown artifacts.

Question title: `Highest daily hops for one aircraft on one flight number`

Question-specific guidance:

Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

## Dashboard Questions

- Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?
- What geographic pattern do the top itineraries show?