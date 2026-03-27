### Multi-query additions

- The dashboard should combine narrative and interactive analysis.
- The verified analysis package includes the main query and named supporting section queries that may be used for building visual panels for particular questions.
- Prefill editable SQL panels for all queries in the verified package.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.

### additional questions

- When the question-specific visual prompt explicitly asks an additional question, author a lookup query and use it for building a visual panel
- label that panel query in the query ledger with the header name
- additional lookup query use the currently selected context from the primary/main query result set and fetches additional data to enrich
- lookup query should use a filter to read only rows from the main set that need to be enriched with additional information. 
- example WHERE clause: `where (c1,c2,c3) in ( ('vc1','vc2', 'vc3'), ('wc1','wc2', 'wc3') )`
- Don't try to calculate the main row set again by a complicated SQL subquery/CTE.
- When the date selector changes, rerun the main query and every enrichment query to support the active date range.
- if the enrichment lookup fails, keep the panel visible with degraded-state messaging, report that degraded panel in the ledger, and continue rendering the rest of the dashboard
- If you run supporting or lookup queries, record them in the same visible query ledger as the main query.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{{visual_input_summary_json}}
