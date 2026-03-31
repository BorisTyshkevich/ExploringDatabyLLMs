All 6 source files are in place. Here's a summary of what was built:

---

**`visual_src/` source tree**

| File | Purpose |
|---|---|
| `package.json` | React 18, Leaflet 1.9, Vite 5 |
| `index.html` | Entry HTML |
| `vite.config.js` | `base: './'`, output → `../visual_build`, assets → `visual_assets/` |
| `src/main.jsx` | React root mount |
| `src/App.jsx` | Full dashboard (550 lines) |
| `src/styles.css` | Aviation theme with all layout rules |

**Dashboard behavior**

- **Auto-loads** on mount if a JWE is stored in `localStorage[OnTimeAnalystDashboard::auth::jwe]`
- **Hero + KPI strip** permanently anchored to Rank #1 (WN 366 / N957WN), even when another row is selected
- **Map** uses Leaflet with `--red` for Rank #1 route, `--sky` for others; always visible, degrades gracefully if enrichment fails
- **Row selection** updates the map title, Leaflet route/markers/bounds, and the route sequence panel — independently of the hero/KPI anchor
- **Coordinate enrichment** (q4 pattern) is cached per `tail|flightNum|date` key — no re-query on re-selection
- **Query ledger** lists primary query + one enrichment entry per itinerary viewed; SQL collapses/expands with ▶/▼
- **Footer controls**: JWE password input, Forget button, SQL textarea, Run Query button, status line

**To build:**
```bash
cd visual_src
npm install
npm run build   # outputs to ../visual_build/
```
