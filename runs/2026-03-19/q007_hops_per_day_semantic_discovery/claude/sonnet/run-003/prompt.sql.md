Generate only the artifacts requested in this prompt.
Use the configured MCP server for all data access.
Do not construct raw OpenAPI URLs manually.
Stay within the configured dataset scope.

Dataset semantic layer:

Use `ontime.ontime` as the primary fact table for flight operations.

Use `ontime.airports_latest` as the semantic airport dimension for current airport reference data, including:

- `code`
- `name`
- `latitude`
- `longitude`
- `utc_local_time_variation`

`ontime.airports_latest` can be used directly to get coordinates and other columns for enrichment in application code (such as JavaScript or python)
or used for SQL JOINs.

Preferred sql joins:

- `ontime.ontime.OriginAirportID = ontime.airports_latest.airport_id`
- `ontime.ontime.DestAirportID = ontime.airports_latest.airport_id`

For `ontime.ontime`, use `OriginAirportID` and `DestAirportID` when enriching airport names, coordinates, or other airport attributes inside SQL.
Do not join `ontime.ontime` legs directly to `ontime.airports_latest` by airport code when the airport id columns are already available.

Fallback sql joins:

- use `ontime.ontime.Origin = ontime.airports_latest.code`
- use `ontime.ontime.Dest = ontime.airports_latest.code`

Use the `airport_id` joins when those columns are available.
Use code-based joins only when the analytical result exposes route strings or airport codes but not airport IDs.

Do not emit result rows or any data output.
Inspect the live schema first with `SHOW TABLES FROM ontime` and `DESCRIBE TABLE` for the tables you intend to use.
Before finalizing your answer, self-verify the query with a quick debug execution, usually with a small `LIMIT` or `WHERE` filter in source data reading subquery/CTE, and fix any errors internally.
Write one JSON object containing the final verified SQL and a Markdown report template to `answer.raw.json`, not the debug query.

`answer.raw.json` must contain plain JSON bytes only. Do not wrap the file contents in Markdown fences.

Your stdout response may contain a short status line, but qforge will ignore stdout and load only `answer.raw.json`.

Write exactly this JSON object shape to `answer.raw.json`:

{
  "sql": "-- one SQL statement",
  "report_markdown": "# Highest daily hops for one aircraft on one flight number\\n\\n{{data_overview_md}}\\n\\nThe key derived value is {{metric.primary_value}}.",
  "metrics": {
    "summary_facts": [
      "Summarize the strongest derived fact from the query result."
    ],
    "named_values": {
      "primary_value": "example derived value"
    },
    "named_lists": {
      "example_list": [
        "Example ordered item"
      ]
    }
  }
}

Rules:

- Use one SQL statement only.
- JSON must contain exactly these top-level keys:
  - `sql`
  - `report_markdown`
  - `metrics`
- Write the artifact to `answer.raw.json`.
- The `answer.raw.json` file must contain raw JSON, not fenced Markdown.
- The report must be Markdown.
- The report must be a template, not a data-filled summary.
- Prefer `{{data_overview_md}}` and `{{result_table_md}}` for JSON-derived sections.
- Keep the report concise and analytical.
- Use placeholders only where data is needed.
- Derive SQL, metrics, and report claims only from the current question and the current query result shape.
- Do not rely on prior qforge runs, prior question ids, or previously observed values.
- Do not mention other question ids such as `q001` in the artifact unless the current prompt explicitly asks for cross-question comparison.
- Allowed built-in placeholders: {{row_count}}, {{generated_at}}, {{columns_csv}}, {{question_title}}, {{data_overview_md}}, {{result_table_md}}
- Allowed metric placeholders use this pattern only: `{{metric.<name>}}`
- Do not invent any placeholder outside the built-in list and `{{metric.<name>}}`.
- If a fact is needed in the report and is not covered by a built-in placeholder, put it in `metrics.named_values` and reference it via `{{metric.<name>}}`.
- Do not include TSV, JSON rows, HTML, or any other fenced blocks.

Invalid example:

`"report_markdown": "The key derived value is {{primary_value}}."`

The placeholder `{{primary_value}}` is invalid. Use `metrics.named_values.primary_value` plus `{{metric.primary_value}}` instead.

Question title: `Highest daily hops for one aircraft on one flight number`

Question-specific report guidance:

Explain:

- the maximum hop count observed and whether it appears to be a repeated operating pattern or a one-off itinerary,
- the single most recent itinerary among the maximum-hop rows, including carrier, flight number, date, and full route,
- and any notable route repetition or clustering visible across the top 10 longest itineraries.

Find the highest number of hops per day for a single aircraft using the same flight number.

Report Aircraft ID, Flight Number, Carrier, Date, Route.
For the longest trip, show the actual departure time from each origin.
What does the itinerary look like? Find the top 10 longest and most recent itineraries.

Definitions and rules:

- Do not invent column aliases 
- Build each itinerary in chronological leg order using the actual departure timestamps.
- Use actual `DepTime` values only. Exclude rows where actual departure time is missing; do not fall back to `CRSDepTime`.
- The textual `Route` must include every leg and the final destination airport.
- If `Hops = N`, the route must contain exactly `N` legs or `N + 1` airport codes.