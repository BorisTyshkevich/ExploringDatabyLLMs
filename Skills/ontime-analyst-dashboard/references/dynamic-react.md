# Dynamic React Mode

Use this when the presentation target is React and the visual mode is `dynamic`.

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
- The built page must keep working when `visual_build/index.html` is copied to sibling `visual.html`

## React guidance

- Use client-only React
- Do not use Next.js, routing, SSR, or any backend server
- Keep token, query text, run status, and fetched rows in explicit component state
- The footer token/query controls still need to render in the final page
- The query ledger still needs to render in the final page
- Keep the same dynamic fetch contract, endpoint template, localStorage key, and degradation behavior as HTML dynamic mode
