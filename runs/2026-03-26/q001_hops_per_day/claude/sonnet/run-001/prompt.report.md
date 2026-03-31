- Connect to clickhouse server though MCP connection
- Do not use direct HTTP by any tools like curl.
- Use the `ontime` database to answer analytical questions
- Use `ontime-semantic-layer` skill for schema inspection, join guidance, and dimension semantics.
- write correct and efficient ClickHouse SQL 
- Before writing any SQL artifact, self-verify every SQL statement you intend to save.
- Run a cheap debug execution for each query first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both applied inside the main data-reading subquery or CTE.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, window, join, or unknown-column errors in a loop until every saved query runs successfully.
- Do not write unchecked SQL.

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
- Use the full available dataset history unless the question-specific prompt explicitly asks for a narrower time window.

## Do not

- Do not merge several section questions into one unioned or row-typed SQL result unless the question-specific guidance explicitly requires that.
- Do not collapse a section to a final aggregate when the presentation needs the underlying ranked, time-series, or drilldown rows.
- Do not emit result rows or any data output.
- Do not output fenced Markdown artifacts.

## Question title 
`Highest daily hops for one aircraft on one flight number`

## Question-specific guidance

Read every top-level `###` section in the question-specific guidance below.

### main
Find the longest itineraries with the highest number of hops for a single aircraft using the same flight number.
Define uniqueness by the full textual `Route` string and output the most recent top 10 unique routes by departure time.
Do not exclude rows solely because `Tail_Number` is empty. If an itinerary qualifies but the aircraft id is missing in the source data, keep it in the result and surface the aircraft id as empty / unknown rather than filtering it out.
Count hops from distinct same-day legs, not raw source rows; do not let duplicate or conflicting same-time rows inflate hop count or create artifact routes.

Return:

- aircraft id
- flight number
- carrier
- flight date
- hop count
- route recurrence count: total number of days across all history on which this exact Route string was flown by any aircraft
- textual `Route` in chronological order, including every origin and the final destination, using `-` as the delimiter throughout, for example `SMF-SAN-PHX-COS-DEN`

### q1
Do the top itineraries appear to be recurring scheduled patterns or mostly one-offs?

List the recurrence count for each of the 10 routes explicitly before summarizing any tiers or categories, and make sure any category totals add up to 10.

### q2
What geographic pattern do the top itineraries show?

Base the geographic answer only on the airports appearing in the 10 routes returned by `main`, not on the broader population of all maximum-hop flights in history.