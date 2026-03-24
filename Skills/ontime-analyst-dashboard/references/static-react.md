# Static React Mode

Use this when the presentation target is React and the visual mode is `static`.

## Required source tree

- `visual_src/package.json`
- `visual_src/index.html`
- `visual_src/src/main.jsx`
- `visual_src/src/App.jsx`
- optional `visual_src/src/styles.css`
- optional `visual_src/vite.config.js`

## Build contract

- Build for static hosting only
- Use relative asset paths
- Final build output must land in `../visual_build`
- Built assets must land under `visual_assets/`
- The final built page must be self-contained with respect to analytical data and must not require MCP access, tokens, or localStorage

## React guidance

- Keep the same static semantics as HTML static mode
- Bake the analytical data needed by the page into the React source or HTML entry so the built artifact is self-contained
- Do not use remote runtime dependencies for non-map dashboards
- Do not introduce server-side rendering, routing, or a backend server
