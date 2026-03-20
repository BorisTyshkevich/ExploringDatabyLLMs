## Dynamic mode 

- The saved SQL runs in the browser as the primary query for the page
- Additional browser queries are allowed only for explicit enrichment or drill-down
- Do not embed the primary analytical dataset or examples as stati payloads 

## Required UI

- Header remains visible before data loads
- Main dashboard content hidden before successful fetch
- The saved SQL textarea is the primary analytical query source for the page
- The JWE/SQL control form must live in a real `<footer>` block at the very end of the page
- That control block is a footer-style utility panel, not part of the hero or empty-state layout
- After data loads, the control block must still remain at the bottom of the document, below the analytical content
- Do not place the JWE/SQL controls inside the hero, KPI strip, or main analytical card grid
- Provide a visible query ledger or provenance section that lists the primary query and any enrichment queries
- Provide:
  - JWE token input field (allow to enter new and show locally stored as ***)
  - forget stored token button
  - SQL textarea
  - fetch button
  - status text
  - empty-state hint

## Data Shape

use provided JSON data as an example that can be received by executing the SQL query.

# Fetch flow

1. On page load, read the stored JWE token from `localStorage` and prefill the token field if present
2. If the page auto-runs when a stored token is present, use the same guarded run path as the manual fetch action
3. While a run is active, disable the fetch button, show it in an inactive state, and ignore additional manual clicks
4. Read token and SQL from inputs
5. Validate both are non-empty
6. Persist the JWE token to `localStorage` after every successful token entry so subsequent dashboards can reuse it
7. Call endpoint with `fetch` using the saved SQL shown in the textarea
8. If `response.ok` is false, read the response text and surface that API error directly instead of trying to parse it as JSON first
9. Parse JSON payload with `columns` and `rows` only for successful responses
10. Treat empty results as valid when `count = 0`, even if `rows` is returned as `null`
11. Convert row arrays into objects
12. Run through the same normalization pipeline as static mode
13. Normalize temporal fields before UI formatting, grouping, filtering, or comparison logic
14. If needed, run explicit enrichment or drill-down queries with a concrete purpose and record them in the query ledger
15. Re-enable the fetch button only after the active run has finished or failed
16. Show content and render dashboard while keeping the control block at the bottom

## Response-shape contract

- Browser query responses are tabular payloads
- Use `columns` and `rows` as the primary response fields unless the prompt explicitly defines a different shape
- Convert row arrays into object rows before deriving KPIs, charts, tables, or filters
- Do not assume wrappers such as `meta`, `data`, or `items`
- Treat empty tabular payloads as valid when the response indicates zero rows

## Error handling

- On HTTP failure, show status with response code or API error text from the plain-text response body when present
- On empty result set, show a visible warning instead of a blank page
- On malformed payload, report that `columns`/`rows` were not usable
- Never print or echo the token in status messages
- If the result set is empty, keep KPI/chart containers stable and show a clear warning panel instead of a broken dashboard
- If an enrichment query fails, report the failed query in the ledger, explain which visual degraded, and continue rendering the rest of the page
- Surface empty, failed, and degraded states in visible page UI, not only in console output or logs
- Keep primary-query failures separate from secondary-render failures; a component error must degrade that component instead of being reported as a primary-query failure
- Never expose the token in status messages or visible UI text

## Query ledger contract

- Every query (primary and enrichment) must appear in a single unified ledger
- Each ledger entry must include: label, role, status, rows, and the full SQL text
- SQL query text is hidden by default with a clickable row to expand/reveal
- Use ▶ toggle icon to expand and show query text
- Use ▼ toggle icon to collapse query text
- Ledger entries should be added immediately (Pending status) and updated on completion
- On status update, also update the SQL field if it wasn't known at registration time

## Query execution contract

- Dynamic dashboards default to one primary query: the saved SQL prefilled into the page.
- The embedded saved SQL is authoritative for the artifact; browser storage must not silently replace it.
- Additional browser queries are allowed for enrichment or drill-down when they materially improve the visualization and remain within dataset policy.
- Primary-query success must be enough to render the main dashboard shell and any visuals driven directly by the primary result set.
- Enrichment and drill-down queries upgrade dependent visuals or details; they must not gate whether the dashboard shell renders at all.
- If a secondary query fails, degrade only the dependent component, keep the primary-query analysis visible, and record the failure in both status text and the query ledger.
- Do not generate hidden follow-up SQL or alternate result shapes without surfacing them to the user.
- Prefer dataset-native dimensions and lookup tables rather than inferred or geocoded data when enrichment is needed.
- Prefer `data-role` selectors or stored element references for card internals so map/chart initialization always targets the rendered node rather than inert template content.
- For Leaflet maps, prefer delayed initialization after the dashboard or map card becomes visible. If the map must be created before final layout settles, call `invalidateSize()` after reveal.
- When interactive drill-down is present, keep a clear selected-row or selected-item state and update only the dependent detail/map panels from that selection.
- Keep explicit run state such as `isRunning` or `activeRunId` so auto-load, manual fetch, and async completions cannot overlap or duplicate ledger entries.
- Ignore or cancel stale async completions from earlier runs once a newer run has started.

## Temporal normalization rules

- Add a small helper that normalizes date-like values once, for example by deriving `flightDateKey = String(value ?? '').slice(0, 10)` when the source may be ISO-like.
- Use the normalized key for filtering, grouping, comparisons, and display helpers that expect a date-only value.
- Keep the original raw value only when the exact source timestamp is analytically meaningful.
- Never append `T00:00:00` blindly to a value that may already contain a time component.

## Security

- Use `<input type="password">` for the token field
- Never hardcode a real token in examples
- Always provide a `Forget` control that removes the token from storage
- Forget must immediately clear the shared JWE storage key used by subsequent dashboards
