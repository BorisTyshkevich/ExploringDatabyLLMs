### Multi-query additions

- The saved SQL shown below is the primary section query for this page.
- The verified analysis package includes named supporting section queries that may be used for enrichment, drill-down, or secondary visuals when the question-specific prompt calls for them.
- Use section answers as narrative framing, but derive displayed KPIs, charts, tables, and interactions from live browser execution of the primary saved SQL and any supporting queries you actually run.
- Do not assume auxiliary lookup or enrichment schema details from memory. Use only columns you have checked against the live endpoint or semantic-layer guidance.
- If you run supporting queries, record them in the same visible query ledger as the primary query.
- The dashboard does not need to mirror `report.md`; it should combine narrative and interactive analysis.

### Verified Analysis Package

Use this JSON package as the supporting context for the visual:

{{visual_input_summary_json}}
