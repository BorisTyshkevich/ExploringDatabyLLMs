Static-mode requirements:

- Build a self-contained benchmark artifact. Do not require browser-side MCP access, tokens, localStorage, or live fetches to render the analytical content.
- Embed the analytical data needed by the page directly in the HTML using inline data blocks such as `<script type="application/json">` or `<script type="text/csv">`.
- Derive KPIs, filters, charts, and tables from the embedded data after parsing and normalization in browser JavaScript.
- Keep CSS and JavaScript inline.
- For non-map visuals, do not use remote `<script>` or `<link>` assets.
- For `html_map`, Leaflet and remote basemap assets are allowed, but the analytical dataset itself must still be embedded rather than fetched.
- Show clear empty-state or missing-field warnings instead of rendering broken visuals.
- Normalize temporal fields before grouping, filtering, or comparison logic. If a ClickHouse `Date` may appear as ISO datetime text, derive a stable `YYYY-MM-DD` key and reuse it consistently.
- If the page uses `<template>` cloning, bind behavior to scoped selectors or per-instance references rather than duplicated global `id` values.
