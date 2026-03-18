Dynamic-mode additions:

- Build the page in dynamic mode using the `ontime-analyst-dashboard` skill contract.
- Execute the embedded saved SQL in the browser as the primary query.
- Keep the embedded saved SQL authoritative for the artifact.
- Surface every browser query in a visible query ledger.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Keep additional browser queries limited to explicit enrichment or drill-down that materially improves the visualization.
- Use literal semantic HTML where the page contract calls for it; do not substitute a footer-styled `<div>` for a real `<footer>`.
- Put the token input, SQL textarea, fetch button, and status text inside a real `<footer>` element at the end of the page.
- Keep query-ledger SQL hidden by default and reveal it through an explicit expand/collapse interaction.
- Surface empty, failed, and degraded states in visible page UI, not only in console output or logs.
