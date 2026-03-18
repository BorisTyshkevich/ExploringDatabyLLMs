Dynamic-mode additions:

- Build the page in dynamic mode using the `ontime-analyst-dashboard` skill contract.
- Execute the embedded saved SQL in the browser as the primary query.
- Keep the embedded saved SQL authoritative for the artifact.
- Surface every browser query in a visible query ledger.
- Do not embed the primary analytical dataset as `result.json` payloads or CSV snapshots.
- Keep additional browser queries limited to explicit enrichment or drill-down that materially improves the visualization.
