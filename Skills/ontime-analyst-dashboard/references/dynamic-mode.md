## Dynamic mode

- The saved SQL runs in the browser as the primary query for the page
- Additional browser queries are allowed only for explicit saved supporting panels, lookup queries, or drill-down
- Do not embed the primary analytical dataset or examples as static payloads
- Default analysis to the most recent 5 years unless the question explicitly asks for a different window
- Dynamic dashboards support two auth modes:
  - `OAuth2` is the default
  - `JWE` remains a supported explicit alternative

## Auth mode selection

- Default to `OAuth2` unless the prompt explicitly asks for manual token entry or the legacy JWE footer controls
- Use `OAuth2` for redirect-based PKCE login where the dashboard page is also the redirect/callback page
- Use `JWE` when the dashboard must expose a password-style token field and `Forget` control in the footer
- Keep the auth choice explicit in the generated page code and prompts; do not blend both modes unless the prompt explicitly asks for that hybrid behavior

## Required UI

- Header remains visible before data loads
- Main dashboard content hidden before successful fetch
- Provide a visible date-range selector with start and end controls
- The primary SQL editor is the authoritative analytical query source for the page
- If the page uses supporting queries, expose editable SQL controls for them too
- Visual-only lookup queries may be authored directly into the page for a concrete panel or lookup need even though they are not part of reviewed analysis artifacts
- The SQL control form must live in a real `<footer>` block at the very end of the page
- That control block is a footer-style utility panel, not part of the hero or empty-state layout
- After data loads, the control block must still remain at the bottom of the document, below the analytical content
- Do not place the footer SQL controls inside the hero, KPI strip, or main analytical card grid
- Provide a visible query ledger or provenance section that lists the primary query, any saved supporting queries, and any lookup queries
- Provide:
  - start date input
  - end date input
  - at least one SQL editor
  - per-query run button
  - `Run all` button when multiple queries are present
  - status text
  - empty-state hint
- Auth-specific controls:
  - `JWE`: put token input and `Forget` in the footer control block with the SQL controls
  - `OAuth2`: provide obvious login/logout state; auth controls may live in the header or footer, but they must remain visually separate from the analytical content

## Data shape

Use provided JSON data as an example that can be received by executing the SQL query.

## Shared fetch flow

1. On page load, initialize auth state, query state, and date-range state before any auto-run decision
2. If the page auto-runs when stored auth is available, use the same guarded run path as the manual fetch action
3. While a run is active, disable the fetch button, show it in an inactive state, and ignore additional manual clicks
4. Read the current SQL from the active query editor
5. Validate that SQL is non-empty and auth is available for the selected auth mode
6. Call the endpoint with `fetch` using the current SQL shown in the active query editor
7. If `response.ok` is false, read the response text and surface that API error directly instead of trying to parse it as JSON first
8. Parse JSON payload with `columns` and `rows` only for successful responses
9. Treat empty results as valid when `count = 0`, even if `rows` is returned as `null`
10. Convert row arrays into objects
11. Run through the same normalization pipeline as static mode
12. Normalize temporal fields before UI formatting, grouping, filtering, or comparison logic
13. Apply the selected date range through SQL reruns for the primary query and any supporting queries that support date filtering
14. If needed, run explicit supporting, lookup, or drill-down queries with a concrete purpose and record them in the query ledger
15. If a query template uses placeholder filter variables such as `__START_DATE__`, `__END_DATE__`, or selection-context tokens, resolve them through one shared helper before execution
16. Use that same substitution path for manual query runs, `Run all`, auto-load, and selection-driven lookup reruns; never send raw placeholder tokens to the SQL endpoint
17. Re-enable run buttons only after the active run has finished or failed
18. Show content and render dashboard while keeping the control block at the bottom

## JWE auth contract

- Read the stored JWE token from `localStorage` and prefill the token field if present
- Validate both token and SQL before running
- Persist the JWE token to `localStorage` after every successful token entry so subsequent dashboards can reuse it
- Use `<input type="password">` for the token field
- Allow the stored token state to render as masked `***` rather than exposing the raw token
- Always provide a `Forget` control that removes the token from storage
- `Forget` must immediately clear the shared JWE storage key used by subsequent dashboards
- Keep the token field, `Forget`, SQL controls, and status in the footer utility panel

## OAuth2 auth contract

- Model the browser auth flow on [`app.html`](/Users/bvt/work/ExploringDatabyLLMs/app.html)
- Use redirect-based OAuth2 with PKCE; the dashboard page is also the redirect/callback page
- Register the OAuth client dynamically against the MCP server before redirecting to the authorization endpoint
- Generate and persist a PKCE verifier and challenge before redirect
- Build `redirect_uri` from the dashboard page URL so the page can complete the callback itself
- Persist enough OAuth state to validate the callback and resume after redirect
- Parse `code` and `state` from the callback URL on page load
- Validate callback state before token exchange
- Exchange the authorization code for an access token on the callback page
- Persist the resulting access token and expiry metadata in browser storage
- Clear callback parameters from the visible URL after parsing them
- Show an authenticated state before exposing query execution controls
- Provide explicit login and logout controls
- On `401` or `403`, clear invalid auth state and return the dashboard to a login-required state
- Do not require manual bearer-token entry in normal OAuth2 usage
- Generated dashboards should assume the redirect page URL has been registered in Google Console or the relevant provider configuration; the page cannot automate that external setup

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
- Never print or echo secrets or tokens in status messages
- If the result set is empty, keep KPI/chart containers stable and show a clear warning panel instead of a broken dashboard
- If a lookup or enrichment query fails, report the failed query in the ledger, explain which visual degraded, and continue rendering the rest of the page
- Surface empty, failed, and degraded states in visible page UI, not only in console output or logs
- Keep primary-query failures separate from secondary-render failures; a component error must degrade that component instead of being reported as a primary-query failure
- Never expose stored auth material in visible UI text

## Query ledger contract

- Every query (primary, saved supporting, and lookup) must appear in a single unified ledger
- Each ledger entry must include: query id, label, role, effective date range, status, rows, and the full SQL text
- SQL query text is hidden by default with a clickable row to expand/reveal
- Use ▶ toggle icon to expand and show query text
- Use ▼ toggle icon to collapse query text
- Ledger entries should be added immediately (Pending status) and updated on completion
- On status update, also update the SQL field if it wasn't known at registration time

## Query execution contract

- Dynamic dashboards default to one primary query: the saved SQL prefilled into the page.
- The embedded saved SQL is authoritative for the artifact; browser storage must not silently replace it.
- Additional browser queries are allowed for enrichment or drill-down when they materially improve the visualization and remain within dataset policy.
- Supporting queries that are shipped with the dashboard must be editable and individually executable.
- Lookup queries are authored in the visual-only second pass and are not part of reviewed `answer.raw.json` artifacts.
- When multiple shipped queries are present, provide a `Run all` path that reruns the date-compatible query set against the selected range.
- If query templates contain placeholder variables, route every execution path through one shared substitution helper so date and selection filters are applied consistently.
- Primary-query success must be enough to render the main dashboard shell and any visuals driven directly by the primary result set.
- Enrichment and drill-down queries upgrade dependent visuals or details; they must not gate whether the dashboard shell renders at all.
- If a secondary query fails, degrade only the dependent component, keep the primary-query analysis visible, and record the failure in both status text and the query ledger.
- Do not generate hidden follow-up SQL or alternate result shapes without surfacing them to the user.
- Prefer explicit query templates, wrappers, or parameters for date filtering instead of brittle SQL string surgery.
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

- Never hardcode a real token in examples
- Never expose bearer tokens, JWE tokens, PKCE verifiers, or access tokens in status text or visible UI
- Keep auth storage scoped to the selected auth mode; do not reuse a JWE key as an OAuth state store
