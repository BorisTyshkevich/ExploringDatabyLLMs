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
- Keep auth state, query text, run status, and fetched rows in explicit component state
- Support both dynamic auth variants from HTML dynamic mode:
  - `OAuth2` by default
  - `JWE` when the prompt explicitly asks for token-entry footer controls
- The footer SQL/query controls still need to render in the final page
- The query ledger still needs to render in the final page
- Keep the same auth-mode contract, query execution contract, endpoint template, and degradation behavior as HTML dynamic mode
- For `JWE`, keep token field state, masked stored-token behavior, and `Forget` behavior explicit in component state and storage helpers
- For `OAuth2`, keep callback parsing, redirect-state recovery, login/logout state, and token-expiry handling explicit in component state and effects
- If using React effects for OAuth2 callback handling, ensure the callback is consumed once and the cleaned URL state does not trigger duplicate exchanges
