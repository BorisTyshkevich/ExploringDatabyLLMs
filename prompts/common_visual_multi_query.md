### Multi-query additions

- The saved SQL shown below is the primary section query for this page.
- The verified analysis package includes named supporting section queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Treat the supporting queries in the verified analysis package as first-class runtime queries, not just narrative context.
- Prefill editable SQL panels for the primary query and each supporting query from the verified package.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- When the question-specific visual prompt explicitly asks for a lookup query, author that query only in the visual pass; do not treat it as part of the reviewed analysis package.
- Do not assume auxiliary lookup or enrichment schema details from memory. Use only columns you have checked against the live endpoint or semantic-layer guidance.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- If you add a lookup query in the visual pass, record it in that same unified query ledger with its own label and status.
- When the date selector changes, rerun the primary query and every supporting query that supports the active date range.
- The dashboard does not need to mirror `report.md`; it should combine narrative and interactive analysis.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{{visual_input_summary_json}}
