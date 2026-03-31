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

Question title: `Worst winter carrier-origin pairs by departure performance`

Question-specific guidance:

Determine which airline and origin-airport combinations perform worst in winter after applying a meaningful flight threshold.

Focus on winter departures only and evaluate completed flights at the `(carrier, origin airport)` level. Limit the analysis to combinations with enough winter traffic to be credible.

Use winter consistently as the business definition of the season for the full available history. You may apply a reasonable minimum-volume filter to remove noise, but do not invent a custom score or let delay-cause shares replace the primary performance ranking.

For each qualifying pair, quantify:

- winter flight volume
- departure on-time performance
- average departure delay
- how reported delay minutes split across weather and operational causes such as carrier, NAS, security, and late aircraft

If delay-cause fields are missing or unreported for some of the weakest pairs, surface that explicitly in both the proof query output and the report. Do not present a blanket weather-vs-operational conclusion as if it covers pairs whose cause data is unavailable.

Rank the worst-performing winter pairs by on-time performance, while using the delay-cause mix as context rather than as the primary ranking driver.

Provide one proof query for each required business question. Across those proof queries, include enough evidence to support both:

- a ranked view of the weakest qualifying winter carrier-airport pairs
- a cause-composition view for the leading weak pairs that separates weather from operational causes

The proof query behind the concentration question must directly support both parts of that question. It should make carrier concentration and airport concentration separately inspectable from the returned result, rather than requiring the reader to infer one of them indirectly from lists or arrays.

The proof query behind the worst-pair question should preserve the ranked pair-level rows needed for the dashboard, not just a single worst pair or summary count.

## Dashboard Questions

- Which winter carrier-airport pair ranks worst overall?
- Are the worst pairs driven more by weather or by operational causes?
- Are the weakest pairs concentrated in a small number of carriers or airports?

In the report, answer those questions directly in prose. Name the worst winter pair, summarize whether the weakest pairs are driven more by weather or by operational causes, and state whether the weak set is concentrated in a small number of carriers or airports.

If some leading weak pairs lack delay-cause reporting, say so directly and limit the weather-vs-operational conclusion to the subset with measured cause data.

Do not use fallback phrases such as "the worst pair" or "the weakest pairs" when your verified query results let you name the actual carrier-airport combinations directly.

Keep the result business-readable and analytically sound. Exclude low-volume winter pairs before ranking them.