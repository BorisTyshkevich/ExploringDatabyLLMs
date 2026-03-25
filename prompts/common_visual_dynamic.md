### Dynamic-mode additions

- Use this endpoint template for every browser query: `{{dynamic_query_endpoint_template}}`
- Keep JWE in `localStorage['OnTimeAnalystDashboard::auth::jwe']`.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Before writing `visual.html`, self-verify every browser-side SQL statement you intend to ship, including primary, supporting, enrichment, drill-down, and lookup queries.
- For each query, run a cheap live check against the real endpoint and schema first, usually with a small `LIMIT`, a narrow `WHERE` filter, or both when that preserves the query shape.
- Treat successful execution as mandatory. Fix any syntax, type, aggregate, join, or unknown-column errors in a loop until every shipped browser query runs successfully.
