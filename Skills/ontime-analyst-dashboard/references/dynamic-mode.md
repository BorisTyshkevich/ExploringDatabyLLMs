# Dynamic Mode

Use dynamic mode for browser dashboards that fetch live data from the MCP OpenAPI interface.

## Endpoint and auth

- Base URL: `https://mcp.demo.altinity.cloud`
- Request URL:
  `https://mcp.demo.altinity.cloud/{JWE_TOKEN}/openapi/execute_query?query=...`
- JWE is stored in browser `localStorage`

## Storage keys

- auth key:
  `OnTimeAnalystDashboard::auth::jwe`
- layout key pattern:
  `OnTimeAnalystDashboard::<dashboardId>::layout`
- optional query key pattern:
  `OnTimeAnalystDashboard::<dashboardId>::query`

## Required UI

- Header remains visible before data loads
- Main dashboard content hidden before successful fetch
- The JWE/SQL control form must live in its own separate block at the very end of the page
- That control block is a footer-style utility panel, not part of the hero or empty-state layout
- After data loads, the control block must still remain at the bottom of the document, below the analytical content
- Do not place the JWE/SQL controls inside the hero, KPI strip, or main analytical card grid
- Provide:
  - JWE token input field (allow to enter new and show locally stored as ***)
  - forget stored token button
  - SQL textarea
  - fetch button
  - status text
  - empty-state hint

## Fetch flow

1. On page load, read the stored JWE token from `localStorage` and prefill the token field if present
2. Read token and SQL from inputs
3. Validate both are non-empty
4. Persist the JWE token to `localStorage` after every successful token entry so subsequent dashboards can reuse it
5. Optionally persist the current SQL by dashboard id
6. Call endpoint with `fetch`
7. Parse JSON payload with `columns` and `rows`
8. Treat empty results as valid when `count = 0`, even if `rows` is returned as `null`
9. Convert row arrays into objects
10. Run through the same normalization pipeline as static mode
11. Show content and render dashboard while keeping the control block at the bottom

Do not embed analytical result rows as CSV or JSON in dynamic mode. The browser should fetch the active dataset live after the user supplies JWE and SQL.

## Error handling

- On HTTP failure, show status with response code or API error text
- On empty result set, show a visible warning instead of a blank page
- On malformed payload, report that `columns`/`rows` were not usable
- Never print or echo the token in status messages
- If the result set is empty, keep KPI/chart containers stable and show a clear warning panel instead of a broken dashboard

## Payload rules

- Normal success payloads may include:
  - `columns`
  - `types`
  - `rows`
  - `count`
- Do not require `rows` to be an array in the empty-result case.
- If `columns` is present and `count = 0` and `rows` is `null`, treat it as an empty dataset, not as a malformed payload.

## Security

- Use `<input type="password">` for the token field
- Never hardcode a real token in examples
- Always provide a `Forget` control that removes the token from storage
- Forget must immediately clear the shared JWE storage key used by subsequent dashboards

## Compatibility note

Dynamic mode is for browser-opened dashboards. For benchmark-produced `visual.html`, prefer static mode unless live loading is explicitly requested.
