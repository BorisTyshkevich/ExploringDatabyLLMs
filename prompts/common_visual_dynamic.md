### Dynamic-mode additions

- Use this endpoint template for every browser query: `{{dynamic_query_endpoint_template}}`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Provide a visible start/end date selector for the analytical range.
- Drive the date selector through SQL reruns, not client-side filtering alone.
- Default the visual to the most recent 5 years relative to the latest available analytical date when a usable date field exists.
- If no safe date field can be detected for a query, keep the selector visible but disable it with a clear note for that query or view.
- Expose editable SQL controls for the primary query and every supporting query the page uses.
- Provide individual run buttons for editable queries and a `Run all` path when the page uses multiple queries.
- Keep one unified query ledger that records each execution with query id, role, effective date range, status, rows, and expandable SQL text.
- Prefer an explicit SQL wrapping or parameter-insertion strategy for date predicates instead of brittle string replacement.
- If a supporting query cannot be safely date-parameterized, keep it editable and manually runnable, and surface that limitation in the UI.
- Before writing `visual.html`, self-verify every browser-side SQL statement you intend to ship, including primary, supporting, enrichment, drill-down, and lookup queries.
- For each query, run a cheap live check against the real endpoint and schema first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both when that preserves the query shape.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, join, or unknown-column errors in a loop until every shipped browser query runs successfully.
