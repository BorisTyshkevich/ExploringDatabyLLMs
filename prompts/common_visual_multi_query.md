### Multi-query additions

- The dashboard should combine narrative and interactive analysis.
- The verified analysis package includes the main query and named supporting section queries that may be used for building visual panels for particular questions.
- Prefill editable SQL panels for all queries in the verified package.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.

### additional questions

- Saved supporting queries come from reviewed analysis artifacts; lookup queries are visual-only second-pass queries authored directly into the page for a concrete panel, lookup, or drill-down need.
- When the question-specific visual prompt explicitly asks an additional question, author a lookup query and use it for building a visual panel.
- Label that panel query in the query ledger with the header name.
- A lookup query uses the currently selected context from the primary/main query result set and fetches additional data to enrich it. It may also depend on current dashboard state when that dependency is explicit in the UI.
- A lookup query should use a filter to read only rows from the main set that need to be enriched with additional information.
- Example WHERE clause: `where (c1,c2,c3) in ( ('vc1','vc2', 'vc3'), ('wc1','wc2', 'wc3') )`
- Do not try to recalculate the main row set with a complicated SQL subquery/CTE.
- When the date selector changes, rerun the main query and every enrichment query to support the active date range.
- If a lookup query fails, degrade only the dependent panel with visible degraded-state messaging, keep the primary analysis visible, record the lookup failure in the ledger and visible status UI, and continue rendering the rest of the dashboard.
- If you run supporting or lookup queries, record them in the same visible query ledger as the main query.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{{visual_input_summary_json}}
