```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Highest daily hops for one aircraft on one flight number</title>
  <link
    rel="stylesheet"
    href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"
    integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
    crossorigin=""
  >
  <style>
    :root {
      --bg-top: #eaf3f8;
      --bg-bottom: #f6fafc;
      --panel: #ffffff;
      --panel-alt: #f3f7fa;
      --ink: #163244;
      --muted: #5d7485;
      --navy: #0e3a52;
      --slate: #5c7080;
      --sky: #3c88b5;
      --teal: #1f8a70;
      --amber: #d48a1f;
      --red: #c54f36;
      --grid: rgba(22, 50, 68, 0.12);
      --border: rgba(22, 50, 68, 0.10);
      --shadow: 0 18px 45px rgba(14, 58, 82, 0.10);
      --radius-xl: 22px;
      --radius-lg: 16px;
      --radius-md: 12px;
      --radius-sm: 8px;
    }

    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      font-family: "Segoe UI", "Helvetica Neue", sans-serif;
      color: var(--ink);
      background:
        radial-gradient(circle at top left, rgba(60, 136, 181, 0.16), transparent 30%),
        radial-gradient(circle at top right, rgba(31, 138, 112, 0.14), transparent 28%),
        linear-gradient(180deg, var(--bg-top) 0%, var(--bg-bottom) 100%);
    }
    h1, h2, h3, h4 {
      margin: 0;
      font-family: Georgia, ui-serif, serif;
      font-weight: 700;
      letter-spacing: 0.01em;
      color: var(--navy);
    }
    p { margin: 0; }
    button, input, select, textarea {
      font: inherit;
    }
    .page {
      max-width: 1280px;
      margin: 0 auto;
      padding: 24px 20px 36px;
    }
    .hero {
      padding: 28px 30px;
      border-radius: var(--radius-xl);
      background:
        linear-gradient(135deg, rgba(14, 58, 82, 0.96), rgba(60, 136, 181, 0.92)),
        linear-gradient(180deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0));
      color: #f4fbff;
      box-shadow: var(--shadow);
    }
    .hero h1 {
      color: #f8fcff;
      font-size: clamp(2rem, 3.3vw, 3rem);
      line-height: 1.05;
    }
    .hero p {
      margin-top: 10px;
      max-width: 860px;
      color: rgba(244, 251, 255, 0.86);
      font-size: 1rem;
      line-height: 1.5;
    }
    .hero-meta {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
      margin-top: 18px;
    }
    .pill {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      padding: 8px 12px;
      border-radius: 999px;
      background: rgba(255, 255, 255, 0.14);
      color: #f2fbff;
      font-size: 0.9rem;
    }
    .status-bar {
      margin-top: 16px;
      padding: 14px 16px;
      border-radius: var(--radius-md);
      background: rgba(255, 255, 255, 0.18);
      color: #f7fbff;
      font-size: 0.95rem;
    }
    .status-bar[data-tone="warning"] {
      background: rgba(212, 138, 31, 0.22);
    }
    .status-bar[data-tone="error"] {
      background: rgba(197, 79, 54, 0.24);
    }
    .status-bar[data-tone="ok"] {
      background: rgba(31, 138, 112, 0.24);
    }
    .content {
      display: none;
      margin-top: 22px;
    }
    .content.visible {
      display: block;
    }
    .warning-panel {
      display: none;
      margin-top: 18px;
      padding: 16px 18px;
      border-radius: var(--radius-md);
      background: rgba(212, 138, 31, 0.12);
      border: 1px solid rgba(212, 138, 31, 0.22);
      color: var(--ink);
    }
    .warning-panel.visible {
      display: block;
    }
    .kpi-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      gap: 14px;
      margin-bottom: 18px;
    }
    .kpi-card, .panel {
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: var(--radius-lg);
      box-shadow: var(--shadow);
    }
    .kpi-card {
      padding: 16px 18px;
      min-height: 112px;
    }
    .kpi-label {
      font-size: 0.82rem;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      color: var(--muted);
    }
    .kpi-value {
      margin-top: 10px;
      color: var(--navy);
      font-size: clamp(1.35rem, 2.3vw, 2rem);
      line-height: 1.05;
      font-weight: 700;
    }
    .kpi-sub {
      margin-top: 10px;
      color: var(--muted);
      font-size: 0.92rem;
      line-height: 1.35;
    }
    .main-grid {
      display: grid;
      grid-template-columns: minmax(0, 1.65fr) minmax(320px, 0.95fr);
      gap: 18px;
      align-items: start;
    }
    .panel {
      padding: 18px;
    }
    .panel-head {
      display: flex;
      align-items: start;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 12px;
    }
    .panel-title {
      font-size: 1.28rem;
    }
    .panel-note {
      color: var(--muted);
      font-size: 0.92rem;
      line-height: 1.4;
    }
    .badge {
      flex: 0 0 auto;
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 7px 10px;
      border-radius: 999px;
      background: var(--panel-alt);
      color: var(--navy);
      font-size: 0.84rem;
    }
    .map-shell {
      position: relative;
      height: 440px;
      border-radius: 14px;
      overflow: hidden;
      border: 1px solid var(--border);
      background:
        linear-gradient(180deg, rgba(60, 136, 181, 0.18), rgba(255, 255, 255, 0.08)),
        linear-gradient(135deg, #dbe9f2, #f5f9fc);
    }
    #map {
      position: absolute;
      inset: 0;
    }
    .map-message {
      position: absolute;
      left: 16px;
      right: 16px;
      bottom: 16px;
      z-index: 450;
      padding: 12px 14px;
      border-radius: var(--radius-md);
      background: rgba(255, 255, 255, 0.92);
      color: var(--ink);
      border: 1px solid rgba(22, 50, 68, 0.14);
      box-shadow: 0 10px 24px rgba(14, 58, 82, 0.10);
      font-size: 0.92rem;
      line-height: 1.35;
    }
    .legend {
      display: flex;
      flex-wrap: wrap;
      gap: 12px 18px;
      margin-top: 14px;
      font-size: 0.92rem;
      color: var(--muted);
    }
    .legend-item {
      display: inline-flex;
      align-items: center;
      gap: 8px;
    }
    .legend-swatch {
      width: 14px;
      height: 14px;
      border-radius: 999px;
      display: inline-block;
    }
    .legend-line {
      width: 20px;
      height: 0;
      border-top: 3px solid var(--sky);
      display: inline-block;
    }
    .legend-line.highlight {
      border-top-color: var(--red);
    }
    .detail-grid {
      display: grid;
      gap: 14px;
    }
    .detail-topline {
      display: flex;
      justify-content: space-between;
      gap: 10px;
      flex-wrap: wrap;
      padding: 12px 14px;
      border-radius: var(--radius-md);
      background: var(--panel-alt);
    }
    .detail-topline strong {
      color: var(--navy);
    }
    .sequence-list {
      list-style: none;
      margin: 0;
      padding: 0;
      display: grid;
      gap: 10px;
    }
    .sequence-item {
      padding: 12px 14px;
      border-radius: var(--radius-md);
      background: var(--panel-alt);
      border: 1px solid var(--border);
    }
    .sequence-item strong {
      color: var(--navy);
    }
    .sequence-meta {
      margin-top: 5px;
      color: var(--muted);
      font-size: 0.9rem;
    }
    .detail-stats {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
      gap: 10px;
    }
    .mini-stat {
      padding: 12px 14px;
      border-radius: var(--radius-md);
      background: linear-gradient(180deg, #f6fafc, #eef4f8);
      border: 1px solid var(--border);
    }
    .mini-stat .label {
      color: var(--muted);
      font-size: 0.82rem;
      text-transform: uppercase;
      letter-spacing: 0.06em;
    }
    .mini-stat .value {
      margin-top: 8px;
      color: var(--navy);
      font-size: 1.16rem;
      font-weight: 700;
    }
    .lower-grid {
      display: grid;
      grid-template-columns: minmax(0, 1fr);
      gap: 18px;
      margin-top: 18px;
    }
    .toolbar {
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 12px;
    }
    .toolbar-group {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
    }
    .control, .button, .sql-textarea {
      border-radius: var(--radius-sm);
      border: 1px solid rgba(22, 50, 68, 0.16);
      background: #fff;
      color: var(--ink);
    }
    .control {
      min-height: 40px;
      padding: 9px 12px;
    }
    .button {
      min-height: 40px;
      padding: 9px 14px;
      cursor: pointer;
      background: linear-gradient(180deg, #fff, #f3f8fb);
    }
    .button:hover:not(:disabled) {
      background: linear-gradient(180deg, #f8fcff, #edf5f9);
    }
    .button:disabled {
      cursor: not-allowed;
      opacity: 0.55;
    }
    .button.primary {
      background: linear-gradient(180deg, var(--navy), #124966);
      color: #fff;
      border-color: rgba(14, 58, 82, 0.4);
    }
    .button.primary:hover:not(:disabled) {
      background: linear-gradient(180deg, #10435f, #14506f);
    }
    .table-wrap {
      overflow: auto;
      border: 1px solid var(--border);
      border-radius: var(--radius-md);
    }
    table {
      width: 100%;
      border-collapse: collapse;
      min-width: 860px;
    }
    th, td {
      padding: 11px 12px;
      border-bottom: 1px solid var(--border);
      vertical-align: top;
      text-align: left;
      font-size: 0.94rem;
    }
    thead th {
      position: sticky;
      top: 0;
      background: #f2f7fb;
      color: var(--navy);
      z-index: 1;
    }
    tbody tr {
      cursor: pointer;
      transition: background 120ms ease;
    }
    tbody tr:hover {
      background: #f7fbfd;
    }
    tbody tr.active-row {
      background: rgba(60, 136, 181, 0.16);
      box-shadow: inset 4px 0 0 var(--red);
    }
    tbody tr:last-child td {
      border-bottom: 0;
    }
    .route-cell {
      min-width: 420px;
      color: var(--muted);
      font-size: 0.9rem;
      line-height: 1.35;
    }
    .mono {
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    }
    .ledger-list {
      display: grid;
      gap: 10px;
    }
    .ledger-entry {
      border: 1px solid var(--border);
      border-radius: var(--radius-md);
      overflow: hidden;
      background: var(--panel);
    }
    .ledger-summary {
      display: grid;
      grid-template-columns: 1em minmax(0, 1fr) 6em 5.5em 4.5em;
      gap: 12px;
      align-items: center;
      width: 100%;
      padding: 12px 14px;
      border: 0;
      background: transparent;
      color: inherit;
      cursor: pointer;
      text-align: left;
    }
    .ledger-summary:hover {
      background: var(--panel-alt);
    }
    .ledger-toggle {
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      color: var(--navy);
    }
    .ledger-label {
      font-weight: 600;
      color: var(--navy);
    }
    .ledger-meta {
      color: var(--muted);
      font-size: 0.88rem;
    }
    .ledger-status.ok { color: var(--teal); font-weight: 700; }
    .ledger-status.pending { color: var(--amber); font-weight: 700; }
    .ledger-status.failed { color: var(--red); font-weight: 700; }
    .ledger-sql {
      display: none;
      padding: 0 14px 14px;
    }
    .ledger-entry.open .ledger-sql {
      display: block;
    }
    .ledger-sql pre {
      margin: 0;
      padding: 14px;
      border-radius: var(--radius-sm);
      background: var(--panel-alt);
      white-space: pre-wrap;
      word-break: break-word;
      font-size: 0.84rem;
      line-height: 1.45;
      color: var(--ink);
    }
    footer {
      margin-top: 24px;
      padding: 18px;
      border-radius: var(--radius-lg);
      background: rgba(255, 255, 255, 0.88);
      border: 1px solid var(--border);
      box-shadow: var(--shadow);
    }
    .footer-grid {
      display: grid;
      gap: 14px;
    }
    .footer-controls {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto auto;
      gap: 10px;
      align-items: start;
    }
    .token-wrap {
      display: grid;
      gap: 8px;
    }
    .sql-textarea {
      width: 100%;
      min-height: 260px;
      padding: 12px 14px;
      resize: vertical;
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 0.84rem;
      line-height: 1.45;
    }
    .footer-hint {
      color: var(--muted);
      font-size: 0.92rem;
      line-height: 1.4;
    }
    .muted {
      color: var(--muted);
    }
    .empty-state {
      padding: 14px 0 2px;
      color: var(--muted);
      font-size: 0.94rem;
      line-height: 1.4;
    }
    @media (max-width: 980px) {
      .main-grid {
        grid-template-columns: 1fr;
      }
      .footer-controls {
        grid-template-columns: 1fr;
      }
      .map-shell {
        height: 360px;
      }
    }
    @media (max-width: 640px) {
      .page {
        padding: 16px 14px 28px;
      }
      .hero {
        padding: 22px 18px;
      }
      .panel, .kpi-card, footer {
        padding: 16px;
      }
      .ledger-summary {
        grid-template-columns: 1em minmax(0, 1fr);
      }
      .ledger-summary > :nth-child(3),
      .ledger-summary > :nth-child(4),
      .ledger-summary > :nth-child(5) {
        display: none;
      }
    }
  </style>
</head>
<body>
  <div class="page">
    <header class="hero">
      <h1>Highest daily hops for one aircraft on one flight number</h1>
      <p>
        Dynamic route dashboard for the top-ranked aircraft-day itineraries. The lead itinerary anchors the KPI strip,
        while row selection redraws the itinerary map and route sequence without changing the headline ranking context.
      </p>
      <div class="hero-meta">
        <span class="pill">Visual mode: dynamic</span>
        <span class="pill">Visual type: html_map</span>
        <span class="pill">Dataset scope: <span class="mono">ontime.ontime</span> + <span class="mono">ontime.airports_latest</span></span>
      </div>
      <div class="status-bar" id="statusBar" data-tone="warning">Awaiting token and query execution.</div>
    </header>

    <div class="warning-panel" id="warningPanel"></div>

    <main class="content" id="dashboardContent">
      <section class="kpi-grid" id="kpiGrid"></section>

      <section class="main-grid">
        <article class="panel">
          <div class="panel-head">
            <div>
              <h2 class="panel-title">Lead-Itinerary Map</h2>
              <p class="panel-note">
                Airport-coordinate enrichment upgrades the plotted route when lookup rows are available. The map card remains present even in degraded state.
              </p>
            </div>
            <div class="badge" id="selectedBadge">Selected itinerary</div>
          </div>
          <div class="map-shell">
            <div id="map" aria-label="Lead itinerary map"></div>
            <div class="map-message" id="mapMessage">Load the primary result set to render the lead itinerary and queue airport-coordinate enrichment.</div>
          </div>
          <div class="legend">
            <span class="legend-item"><span class="legend-swatch" style="background: var(--navy);"></span> Airport stop</span>
            <span class="legend-item"><span class="legend-line highlight"></span> Selected itinerary path</span>
            <span class="legend-item"><span class="legend-line"></span> Context route direction</span>
            <span class="legend-item"><span class="legend-swatch" style="background: var(--amber);"></span> Degraded map keeps non-map analysis active</span>
          </div>
        </article>

        <article class="panel">
          <div class="panel-head">
            <div>
              <h2 class="panel-title">Route Sequence</h2>
              <p class="panel-note">Selection-sensitive detail panel for hop order, airport names, and repetition context inside the loaded result set.</p>
            </div>
          </div>
          <div class="detail-grid">
            <div class="detail-topline" id="detailTopline"></div>
            <div class="detail-stats" id="detailStats"></div>
            <ol class="sequence-list" id="sequenceList"></ol>
          </div>
        </article>
      </section>

      <section class="lower-grid">
        <article class="panel">
          <div class="toolbar">
            <div class="toolbar-group">
              <select class="control" id="flightFilter" aria-label="Filter flight number">
                <option value="all">All flight numbers</option>
              </select>
              <select class="control" id="yearFilter" aria-label="Filter year">
                <option value="all">All years</option>
              </select>
            </div>
            <div class="toolbar-group">
              <button class="button" id="exportButton" type="button">Export filtered rows</button>
            </div>
          </div>
          <div class="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Aircraft ID</th>
                  <th>Flight Number</th>
                  <th>Carrier</th>
                  <th>Date</th>
                  <th>Hops</th>
                  <th>Route</th>
                </tr>
              </thead>
              <tbody id="itineraryTableBody"></tbody>
            </table>
          </div>
        </article>

        <article class="panel">
          <div class="panel-head">
            <div>
              <h2 class="panel-title">Query Ledger</h2>
              <p class="panel-note">
                Unified provenance for the primary analytical query and airport-coordinate enrichment query. Expand any row to inspect the exact SQL.
              </p>
            </div>
          </div>
          <div class="ledger-list" id="ledgerList"></div>
        </article>
      </section>
    </main>

    <footer>
      <div class="footer-grid">
        <div>
          <h2 class="panel-title">Browser Query Controls</h2>
          <p class="footer-hint">
            This dashboard executes the saved SQL through the configured endpoint template. The token is stored only in local browser storage and is never shown in status text.
          </p>
        </div>

        <div class="footer-controls">
          <div class="token-wrap">
            <label for="tokenInput">JWE token</label>
            <input class="control" id="tokenInput" type="password" placeholder="Paste JWE token">
          </div>
          <button class="button" id="forgetButton" type="button">Forget stored token</button>
          <button class="button primary" id="fetchButton" type="button">Run saved query</button>
        </div>

        <div>
          <label for="sqlInput">Saved SQL</label>
          <textarea class="sql-textarea" id="sqlInput" spellcheck="false">WITH airport_offsets AS
(
    SELECT
        airport_id,
        any(utc_local_time_variation) AS utc_local_time_variation
    FROM ontime.airports_latest
    GROUP BY airport_id
),
legs AS
(
    SELECT
        o.FlightDate,
        o.TailNum,
        o.FlightNum,
        o.Carrier,
        o.OriginAirportID,
        o.DestAirportID,
        o.Origin,
        o.Dest,
        o.DepTime,
        toDateTime(o.FlightDate)
            + toIntervalDay(if(o.DepTime = 2400, 1, 0))
            + toIntervalMinute(intDiv(if(o.DepTime = 2400, 0, o.DepTime), 100) * 60 + modulo(if(o.DepTime = 2400, 0, o.DepTime), 100)) AS dep_local_ts,
        (
            if(length(a1o.utc_local_time_variation) = 5,
                if(substring(a1o.utc_local_time_variation, 1, 1) = '-', -1, 1)
                * (toInt32(substring(a1o.utc_local_time_variation, 2, 2)) * 60 + toInt32(substring(a1o.utc_local_time_variation, 4, 2))),
                0
            )
        ) AS origin_offset_minutes
    FROM ontime.ontime AS o
    LEFT JOIN airport_offsets AS a1o ON o.OriginAirportID = a1o.airport_id
    WHERE o.Cancelled = 0
      AND o.Diverted = 0
      AND o.DepTime IS NOT NULL
      AND o.TailNum != ''
      AND o.FlightNum != ''
),
itineraries AS
(
    SELECT
        FlightDate,
        TailNum,
        FlightNum,
        Carrier,
        length(ordered_legs) AS Hops,
        arrayStringConcat(
            arrayMap(x -> concat(formatDateTime(x.2, '%H:%i'), ' ', x.3, '->', x.4), ordered_legs),
            ' | '
        ) AS Route
    FROM
    (
        SELECT
            FlightDate,
            TailNum,
            FlightNum,
            Carrier,
            arraySort(x -> x.1, groupArray((
                dep_local_ts - toIntervalMinute(origin_offset_minutes),
                dep_local_ts,
                Origin,
                Dest
            ))) AS ordered_legs
        FROM legs
        GROUP BY
            FlightDate,
            TailNum,
            FlightNum,
            Carrier
    )
),
max_hops AS
(
    SELECT max(Hops) AS max_hops_observed FROM itineraries
),
max_hops_counts AS
(
    SELECT count() AS max_hop_itinerary_count
    FROM itineraries
    CROSS JOIN max_hops
    WHERE Hops = max_hops_observed
)
SELECT
    TailNum AS `Aircraft ID`,
    FlightNum AS `Flight Number`,
    Carrier,
    FlightDate AS Date,
    Hops,
    Route,
    max_hops_observed AS `Maximum Hops Observed`,
    max_hop_itinerary_count AS `Maximum-Hop Itinerary Count`
FROM itineraries
CROSS JOIN max_hops
CROSS JOIN max_hops_counts
ORDER BY
    Hops DESC,
    Date DESC,
    Carrier,
    `Flight Number`,
    `Aircraft ID`
LIMIT 10</textarea>
        </div>

        <div class="empty-state" id="emptyState">
          Enter a valid token or rely on the one stored in <span class="mono">localStorage['OnTimeAnalystDashboard::auth::jwe']</span>.
          The dashboard shell stays visible, and the analytical panels appear after a successful primary query.
        </div>
      </div>
    </footer>
  </div>

  <script
    src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"
    integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo="
    crossorigin=""
  ></script>
  <script>
    (() => {
      const STORAGE_KEY = 'OnTimeAnalystDashboard::auth::jwe';
      const ENDPOINT_TEMPLATE = 'https://mcp.demo.altinity.cloud/{JWE}/openapi/execute_query?query=';

      const state = {
        isRunning: false,
        activeRunId: 0,
        primaryRows: [],
        filteredRows: [],
        selectedKey: null,
        leadKey: null,
        airportLookup: new Map(),
        mapReady: false,
        map: null,
        markerLayer: null,
        routeLayer: null,
        ledger: [],
        enrichmentStatus: 'idle',
        currentFilterFlight: 'all',
        currentFilterYear: 'all'
      };

      const els = {
        statusBar: document.getElementById('statusBar'),
        dashboardContent: document.getElementById('dashboardContent'),
        warningPanel: document.getElementById('warningPanel'),
        kpiGrid: document.getElementById('kpiGrid'),
        selectedBadge: document.getElementById('selectedBadge'),
        detailTopline: document.getElementById('detailTopline'),
        detailStats: document.getElementById('detailStats'),
        sequenceList: document.getElementById('sequenceList'),
        itineraryTableBody: document.getElementById('itineraryTableBody'),
        ledgerList: document.getElementById('ledgerList'),
        tokenInput: document.getElementById('tokenInput'),
        sqlInput: document.getElementById('sqlInput'),
        fetchButton: document.getElementById('fetchButton'),
        forgetButton: document.getElementById('forgetButton'),
        emptyState: document.getElementById('emptyState'),
        mapMessage: document.getElementById('mapMessage'),
        flightFilter: document.getElementById('flightFilter'),
        yearFilter: document.getElementById('yearFilter'),
        exportButton: document.getElementById('exportButton')
      };

      function escapeHtml(value) {
        return String(value ?? '')
          .replaceAll('&', '&amp;')
          .replaceAll('<', '&lt;')
          .replaceAll('>', '&gt;')
          .replaceAll('"', '&quot;')
          .replaceAll("'", '&#39;');
      }

      function formatDateKey(raw) {
        return String(raw ?? '').slice(0, 10);
      }

      function formatDateLabel(raw) {
        const key = formatDateKey(raw);
        if (!key) return 'Unknown date';
        const dt = new Date(key + 'T00:00:00Z');
        if (Number.isNaN(dt.getTime())) return key;
        return new Intl.DateTimeFormat('en-US', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
          timeZone: 'UTC'
        }).format(dt);
      }

      function routeToLegs(route) {
        return String(route ?? '')
          .split('|')
          .map(part => part.trim())
          .filter(Boolean)
          .map((part, index) => {
            const match = part.match(/^(\d{2}:\d{2})\s+([A-Z0-9]{3})\-\>([A-Z0-9]{3})$/);
            if (!match) {
              return {
                seq: index + 1,
                time: '',
                origin: '',
                dest: '',
                raw: part
              };
            }
            return {
              seq: index + 1,
              time: match[1],
              origin: match[2],
              dest: match[3],
              raw: part
            };
          });
      }

      function rowKey(row) {
        return [
          row['Aircraft ID'] ?? '',
          row['Flight Number'] ?? '',
          row.Carrier ?? '',
          formatDateKey(row.Date),
          row.Route ?? ''
        ].join('||');
      }

      function normalizeRows(columns, rows) {
        return (rows ?? []).map(row => {
          const obj = {};
          columns.forEach((column, index) => {
            obj[column] = row?.[index] ?? null;
          });
          obj._dateKey = formatDateKey(obj.Date);
          obj._legs = routeToLegs(obj.Route);
          obj._key = rowKey(obj);
          return obj;
        });
      }

      function parseResponse(payload) {
        if (!payload || !Array.isArray(payload.columns) || !Array.isArray(payload.rows ?? [])) {
          throw new Error('Malformed payload: columns/rows were not usable.');
        }
        return normalizeRows(payload.columns, payload.rows ?? []);
      }

      function setStatus(message, tone = 'warning') {
        els.statusBar.textContent = message;
        els.statusBar.dataset.tone = tone;
      }

      function showWarning(message = '') {
        els.warningPanel.textContent = message;
        els.warningPanel.classList.toggle('visible', Boolean(message));
      }

      function maskToken(token) {
        return token ? '***' : '';
      }

      function getStoredToken() {
        return localStorage.getItem(STORAGE_KEY) ?? '';
      }

      function prefillToken() {
        const token = getStoredToken();
        if (token) {
          els.tokenInput.value = token;
          els.tokenInput.placeholder = maskToken(token);
        }
      }

      function setRunning(isRunning) {
        state.isRunning = isRunning;
        els.fetchButton.disabled = isRunning;
        els.forgetButton.disabled = isRunning;
      }

      function ensureMap() {
        if (state.mapReady || typeof L === 'undefined') return;
        state.map = L.map('map', {
          zoomControl: true,
          attributionControl: true
        }).setView([39.5, -98.35], 4);
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
          maxZoom: 18,
          attribution: '&copy; OpenStreetMap contributors'
        }).addTo(state.map);
        state.markerLayer = L.layerGroup().addTo(state.map);
        state.routeLayer = L.layerGroup().addTo(state.map);
        state.mapReady = true;
        setTimeout(() => state.map?.invalidateSize(), 0);
      }

      function addLedgerEntry({ label, role, status, rows, sql }) {
        const id = `${Date.now()}-${Math.random().toString(16).slice(2)}`;
        state.ledger.push({ id, label, role, status, rows, sql, open: false });
        renderLedger();
        return id;
      }

      function updateLedgerEntry(id, patch) {
        const entry = state.ledger.find(item => item.id === id);
        if (!entry) return;
        Object.assign(entry, patch);
        renderLedger();
      }

      function renderLedger() {
        els.ledgerList.innerHTML = state.ledger.map(entry => `
          <div class="ledger-entry ${entry.open ? 'open' : ''}" data-ledger-id="${escapeHtml(entry.id)}">
            <button class="ledger-summary" type="button" data-ledger-toggle="${escapeHtml(entry.id)}">
              <span class="ledger-toggle">${entry.open ? '▼' : '▶'}</span>
              <span class="ledger-label">${escapeHtml(entry.label)}</span>
              <span class="ledger-meta">${escapeHtml(entry.role)}</span>
              <span class="ledger-status ${escapeHtml(String(entry.status).toLowerCase())}">${escapeHtml(entry.status)}</span>
              <span class="ledger-meta">${escapeHtml(entry.rows ?? 'n/a')}</span>
            </button>
            <div class="ledger-sql">
              <pre>${escapeHtml(entry.sql ?? '')}</pre>
            </div>
          </div>
        `).join('');
      }

      els.ledgerList.addEventListener('click', event => {
        const button = event.target.closest('[data-ledger-toggle]');
        if (!button) return;
        const id = button.getAttribute('data-ledger-toggle');
        const entry = state.ledger.find(item => item.id === id);
        if (!entry) return;
        entry.open = !entry.open;
        renderLedger();
      });

      async function executeSql({ token, sql, ledgerLabel, role, runId }) {
        const ledgerId = addLedgerEntry({
          label: ledgerLabel,
          role,
          status: 'Pending',
          rows: '...',
          sql
        });
        try {
          const endpoint = ENDPOINT_TEMPLATE.replace('{JWE}', encodeURIComponent(token)) + encodeURIComponent(sql);
          const response = await fetch(endpoint);
          if (runId !== state.activeRunId) {
            throw new Error('Stale run ignored.');
          }
          if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || `HTTP ${response.status}`);
          }
          const payload = await response.json();
          const rowCount = payload?.count ?? (Array.isArray(payload?.rows) ? payload.rows.length : 0);
          updateLedgerEntry(ledgerId, {
            status: 'OK',
            rows: rowCount,
            sql
          });
          return payload;
        } catch (error) {
          if (String(error?.message || '') === 'Stale run ignored.') {
            return null;
          }
          updateLedgerEntry(ledgerId, {
            status: 'Failed',
            rows: '0',
            sql
          });
          throw error;
        }
      }

      function renderKpis() {
        const lead = state.primaryRows[0];
        if (!lead) {
          els.kpiGrid.innerHTML = '';
          return;
        }
        const sameRouteCount = state.primaryRows.filter(row => row.Route === lead.Route).length;
        const sameFlightCount = state.primaryRows.filter(row => row['Flight Number'] === lead['Flight Number']).length;
        const cards = [
          {
            label: 'Tail Number',
            value: lead['Aircraft ID'],
            sub: `${lead.Carrier} carrier record in the ranked result`
          },
          {
            label: 'Flight Number',
            value: lead['Flight Number'],
            sub: `${sameFlightCount} of ${state.primaryRows.length} top rows share this flight number`
          },
          {
            label: 'Date',
            value: formatDateLabel(lead.Date),
            sub: 'Visible KPI fixed to the top-ranked itinerary'
          },
          {
            label: 'Hop Count',
            value: String(lead.Hops ?? '0'),
            sub: `Maximum observed = ${lead['Maximum Hops Observed'] ?? 'n/a'}`
          },
          {
            label: 'Route Repetition',
            value: `${sameRouteCount} of ${state.primaryRows.length}`,
            sub: `${lead['Maximum-Hop Itinerary Count'] ?? 'n/a'} maximum-hop itineraries in the full dataset`
          }
        ];
        els.kpiGrid.innerHTML = cards.map(card => `
          <div class="kpi-card">
            <div class="kpi-label">${escapeHtml(card.label)}</div>
            <div class="kpi-value">${escapeHtml(card.value)}</div>
            <div class="kpi-sub">${escapeHtml(card.sub)}</div>
          </div>
        `).join('');
      }

      function populateFilters() {
        const flights = [...new Set(state.primaryRows.map(row => row['Flight Number']).filter(Boolean))].sort((a, b) => String(a).localeCompare(String(b), undefined, { numeric: true }));
        const years = [...new Set(state.primaryRows.map(row => row._dateKey.slice(0, 4)).filter(Boolean))].sort((a, b) => b.localeCompare(a));
        els.flightFilter.innerHTML = '<option value="all">All flight numbers</option>' + flights.map(value => `<option value="${escapeHtml(value)}">${escapeHtml(value)}</option>`).join('');
        els.yearFilter.innerHTML = '<option value="all">All years</option>' + years.map(value => `<option value="${escapeHtml(value)}">${escapeHtml(value)}</option>`).join('');
        els.flightFilter.value = 'all';
        els.yearFilter.value = 'all';
        state.currentFilterFlight = 'all';
        state.currentFilterYear = 'all';
      }

      function applyFilters() {
        state.filteredRows = state.primaryRows.filter(row => {
          const flightOk = state.currentFilterFlight === 'all' || row['Flight Number'] === state.currentFilterFlight;
          const yearOk = state.currentFilterYear === 'all' || row._dateKey.slice(0, 4) === state.currentFilterYear;
          return flightOk && yearOk;
        });
        const visibleKeys = new Set(state.filteredRows.map(row => row._key));
        if (!visibleKeys.has(state.selectedKey)) {
          state.selectedKey = state.filteredRows[0]?._key ?? state.primaryRows[0]?._key ?? null;
        }
      }

      function selectedRow() {
        return state.primaryRows.find(row => row._key === state.selectedKey) ?? null;
      }

      function lookupAirport(code) {
        return state.airportLookup.get(code) ?? null;
      }

      function routeComparison(row) {
        return {
          sameRouteCount: state.primaryRows.filter(candidate => candidate.Route === row.Route).length,
          sameFlightCount: state.primaryRows.filter(candidate => candidate['Flight Number'] === row['Flight Number']).length,
          sameHopCount: state.primaryRows.filter(candidate => candidate.Hops === row.Hops).length
        };
      }

      function renderDetailPanel() {
        const row = selectedRow();
        if (!row) {
          els.detailTopline.innerHTML = '<span>No itinerary selected.</span>';
          els.detailStats.innerHTML = '';
          els.sequenceList.innerHTML = '';
          return;
        }

        const comparison = routeComparison(row);
        const firstLeg = row._legs[0];
        const lastLeg = row._legs[row._legs.length - 1];

        els.selectedBadge.textContent = `${row['Aircraft ID']} · Flight ${row['Flight Number']} · ${formatDateLabel(row.Date)}`;
        els.detailTopline.innerHTML = `
          <span><strong>${escapeHtml(row['Aircraft ID'])}</strong> on flight <strong>${escapeHtml(row['Flight Number'])}</strong> (${escapeHtml(row.Carrier ?? '')})</span>
          <span>${escapeHtml(firstLeg?.origin ?? 'n/a')} to ${escapeHtml(lastLeg?.dest ?? 'n/a')} across ${escapeHtml(String(row.Hops ?? 0))} hops</span>
        `;
        els.detailStats.innerHTML = `
          <div class="mini-stat">
            <div class="label">Same route in top 10</div>
            <div class="value">${escapeHtml(String(comparison.sameRouteCount))}</div>
          </div>
          <div class="mini-stat">
            <div class="label">Same flight number</div>
            <div class="value">${escapeHtml(String(comparison.sameFlightCount))}</div>
          </div>
          <div class="mini-stat">
            <div class="label">Rows with same hops</div>
            <div class="value">${escapeHtml(String(comparison.sameHopCount))}</div>
          </div>
          <div class="mini-stat">
            <div class="label">Sequence count</div>
            <div class="value">${escapeHtml(String(row._legs.length))}</div>
          </div>
        `;
        els.sequenceList.innerHTML = row._legs.map(leg => {
          const originInfo = lookupAirport(leg.origin);
          const destInfo = lookupAirport(leg.dest);
          const originLabel = originInfo?.name ? `${leg.origin} · ${originInfo.name}` : leg.origin;
          const destLabel = destInfo?.name ? `${leg.dest} · ${destInfo.name}` : leg.dest;
          const metaParts = [];
          if (originInfo?.utc_local_time_variation) metaParts.push(`Origin UTC ${originInfo.utc_local_time_variation}`);
          if (destInfo?.utc_local_time_variation) metaParts.push(`Destination UTC ${destInfo.utc_local_time_variation}`);
          return `
            <li class="sequence-item">
              <strong>Hop ${escapeHtml(String(leg.seq))}</strong> · ${escapeHtml(leg.time || 'Unknown time')} · ${escapeHtml(originLabel)} to ${escapeHtml(destLabel)}
              <div class="sequence-meta">${escapeHtml(metaParts.join(' · ') || 'Airport enrichment pending or unavailable for one or both stops.')}</div>
            </li>
          `;
        }).join('');
      }

      function renderTable() {
        els.itineraryTableBody.innerHTML = state.filteredRows.map(row => `
          <tr data-key="${escapeHtml(row._key)}" class="${row._key === state.selectedKey ? 'active-row' : ''}">
            <td class="mono">${escapeHtml(row['Aircraft ID'])}</td>
            <td class="mono">${escapeHtml(row['Flight Number'])}</td>
            <td>${escapeHtml(row.Carrier ?? '')}</td>
            <td>${escapeHtml(formatDateLabel(row.Date))}</td>
            <td>${escapeHtml(String(row.Hops ?? ''))}</td>
            <td class="route-cell">${escapeHtml(row.Route ?? '')}</td>
          </tr>
        `).join('');
      }

      function renderEmptyResult() {
        els.dashboardContent.classList.add('visible');
        els.kpiGrid.innerHTML = '';
        els.detailTopline.innerHTML = '<span>No rows returned by the primary query.</span>';
        els.detailStats.innerHTML = '';
        els.sequenceList.innerHTML = '';
        els.itineraryTableBody.innerHTML = '';
        showWarning('The primary query returned zero rows. KPI, map, and table containers remain stable, but no itineraries were available to analyze.');
        ensureMap();
        state.markerLayer?.clearLayers();
        state.routeLayer?.clearLayers();
        els.mapMessage.textContent = 'No itinerary rows were returned, so no route could be plotted.';
      }

      function buildAirportEnrichmentSql(rows) {
        const codes = [...new Set(rows.flatMap(row => row._legs.flatMap(leg => [leg.origin, leg.dest])).filter(Boolean))].sort();
        if (!codes.length) return '';
        const quoted = codes.map(code => `'${String(code).replaceAll("'", "''")}'`).join(', ');
        return `SELECT
    code,
    any(name) AS name,
    any(latitude) AS latitude,
    any(longitude) AS longitude,
    any(utc_local_time_variation) AS utc_local_time_variation
FROM ontime.airports_latest
WHERE code IN (${quoted})
GROUP BY code
ORDER BY code`;
      }

      function storeAirportLookup(payload) {
        const columns = payload?.columns ?? [];
        const rows = payload?.rows ?? [];
        const codeIndex = columns.indexOf('code');
        const nameIndex = columns.indexOf('name');
        const latIndex = columns.indexOf('latitude');
        const lonIndex = columns.indexOf('longitude');
        const tzIndex = columns.indexOf('utc_local_time_variation');
        state.airportLookup = new Map();
        rows.forEach(row => {
          const code = row?.[codeIndex];
          if (!code) return;
          state.airportLookup.set(code, {
            code,
            name: row?.[nameIndex] ?? '',
            latitude: Number(row?.[latIndex]),
            longitude: Number(row?.[lonIndex]),
            utc_local_time_variation: row?.[tzIndex] ?? ''
          });
        });
      }

      function polylinePointsForRow(row) {
        const stops = [];
        row._legs.forEach((leg, index) => {
          const origin = lookupAirport(leg.origin);
          const dest = lookupAirport(leg.dest);
          if (index === 0 && origin && Number.isFinite(origin.latitude) && Number.isFinite(origin.longitude)) {
            stops.push({
              code: leg.origin,
              name: origin.name,
              lat: origin.latitude,
              lon: origin.longitude,
              label: `${leg.origin}${origin.name ? ' · ' + origin.name : ''}`
            });
          }
          if (dest && Number.isFinite(dest.latitude) && Number.isFinite(dest.longitude)) {
            stops.push({
              code: leg.dest,
              name: dest.name,
              lat: dest.latitude,
              lon: dest.longitude,
              label: `${leg.dest}${dest.name ? ' · ' + dest.name : ''}`
            });
          }
        });
        return stops;
      }

      function renderMapForSelection(source = 'pending') {
        ensureMap();
        const row = selectedRow();
        state.markerLayer?.clearLayers();
        state.routeLayer?.clearLayers();
        if (!row) {
          els.mapMessage.textContent = 'No selected itinerary is available to plot.';
          return;
        }

        const points = polylinePointsForRow(row);
        const enoughPoints = points.length >= 2;
        if (!enoughPoints) {
          const reason = source === 'failed'
            ? 'Airport-coordinate enrichment failed for this run. The map remains visible in degraded mode while route details and the table continue rendering.'
            : 'Airport-coordinate enrichment is pending or incomplete for the selected itinerary. The map card stays visible and will upgrade when coordinates are available.';
          els.mapMessage.textContent = reason;
          state.map?.setView([39.5, -98.35], 4);
          state.map?.invalidateSize();
          return;
        }

        const latLngs = points.map(point => [point.lat, point.lon]);
        L.polyline(latLngs, {
          color: getComputedStyle(document.documentElement).getPropertyValue('--red').trim() || '#c54f36',
          weight: 4,
          opacity: 0.88
        }).addTo(state.routeLayer);

        points.forEach((point, index) => {
          const marker = L.circleMarker([point.lat, point.lon], {
            radius: index === 0 || index === points.length - 1 ? 7 : 5.5,
            color: getComputedStyle(document.documentElement).getPropertyValue('--navy').trim() || '#0e3a52',
            weight: 2,
            fillColor: index === 0 || index === points.length - 1
              ? (getComputedStyle(document.documentElement).getPropertyValue('--amber').trim() || '#d48a1f')
              : (getComputedStyle(document.documentElement).getPropertyValue('--sky').trim() || '#3c88b5'),
            fillOpacity: 0.94
          }).addTo(state.markerLayer);
          marker.bindPopup(`${escapeHtml(point.code)}${point.name ? `<br>${escapeHtml(point.name)}` : ''}`);
        });

        state.map.fitBounds(latLngs, { padding: [28, 28] });
        state.map.invalidateSize();
        els.mapMessage.textContent = `${formatDateLabel(row.Date)} · ${row['Aircraft ID']} · flight ${row['Flight Number']} plotted with ${points.length} mapped stops.`;
      }

      function renderAll() {
        renderKpis();
        applyFilters();
        renderDetailPanel();
        renderTable();
        renderMapForSelection(state.enrichmentStatus === 'failed' ? 'failed' : (state.airportLookup.size ? 'ok' : 'pending'));
      }

      async function runPrimaryQuery(autoRun = false) {
        if (state.isRunning) return;
        const token = els.tokenInput.value.trim();
        const sql = els.sqlInput.value.trim();
        if (!token) {
          setStatus('Enter a JWE token before executing the saved query.', 'warning');
          return;
        }
        if (!sql) {
          setStatus('Saved SQL is empty and cannot be executed.', 'error');
          return;
        }

        const runId = Date.now();
        state.activeRunId = runId;
        state.ledger = [];
        renderLedger();
        state.airportLookup = new Map();
        state.enrichmentStatus = 'pending';
        setRunning(true);
        showWarning('');
        setStatus(autoRun ? 'Running saved query from stored token...' : 'Running saved query...', 'warning');

        try {
          const payload = await executeSql({
            token,
            sql,
            ledgerLabel: 'Primary analytical query',
            role: 'primary',
            runId
          });
          if (!payload || runId !== state.activeRunId) return;
          localStorage.setItem(STORAGE_KEY, token);
          els.tokenInput.placeholder = maskToken(token);

          const rows = parseResponse(payload);
          state.primaryRows = rows;
          state.leadKey = rows[0]?._key ?? null;
          state.selectedKey = state.leadKey;
          populateFilters();

          els.dashboardContent.classList.add('visible');
          els.emptyState.textContent = `Last successful fetch returned ${rows.length} row${rows.length === 1 ? '' : 's'}.`;

          if (!rows.length) {
            renderEmptyResult();
            setStatus('Primary query completed with zero rows.', 'warning');
            return;
          }

          renderAll();
          setStatus(`Primary query completed with ${rows.length} rows. Running airport-coordinate enrichment for mapped stops...`, 'warning');

          const enrichmentSql = buildAirportEnrichmentSql(rows);
          if (!enrichmentSql) {
            state.enrichmentStatus = 'failed';
            showWarning('No airport codes could be parsed from the returned routes. The table and route detail remain available, but the map is degraded.');
            renderMapForSelection('failed');
            setStatus('Primary query completed, but no airport codes were available for map enrichment.', 'warning');
            return;
          }

          try {
            const enrichmentPayload = await executeSql({
              token,
              sql: enrichmentSql,
              ledgerLabel: 'Airport-coordinate enrichment',
              role: 'enrichment',
              runId
            });
            if (!enrichmentPayload || runId !== state.activeRunId) return;
            storeAirportLookup(enrichmentPayload);
            state.enrichmentStatus = 'ok';
            renderDetailPanel();
            renderMapForSelection('ok');
            setStatus(`Primary query completed with ${rows.length} rows and airport-coordinate enrichment returned ${enrichmentPayload?.count ?? enrichmentPayload?.rows?.length ?? 0} lookup rows.`, 'ok');
          } catch (error) {
            state.enrichmentStatus = 'failed';
            showWarning(`Airport-coordinate enrichment failed. The selected itinerary map is degraded, but the route sequence and itinerary table remain available. ${String(error?.message || '')}`.trim());
            renderDetailPanel();
            renderMapForSelection('failed');
            setStatus(`Primary query completed, but airport-coordinate enrichment failed: ${String(error?.message || 'Unknown error')}`, 'warning');
          }
        } catch (error) {
          els.dashboardContent.classList.remove('visible');
          showWarning('');
          setStatus(`Primary query failed: ${String(error?.message || 'Unknown error')}`, 'error');
        } finally {
          if (runId === state.activeRunId) {
            setRunning(false);
          }
        }
      }

      function exportFilteredRows() {
        const rows = state.filteredRows;
        if (!rows.length) {
          setStatus('No filtered rows are available to export.', 'warning');
          return;
        }
        const columns = ['Aircraft ID', 'Flight Number', 'Carrier', 'Date', 'Hops', 'Route', 'Maximum Hops Observed', 'Maximum-Hop Itinerary Count'];
        const csv = [
          columns.join(','),
          ...rows.map(row => columns.map(column => {
            const value = row[column] ?? '';
            const escaped = String(value).replaceAll('"', '""');
            return `"${escaped}"`;
          }).join(','))
        ].join('\n');
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = 'highest_daily_hops_itineraries.csv';
        document.body.appendChild(link);
        link.click();
        link.remove();
        URL.revokeObjectURL(url);
      }

      els.fetchButton.addEventListener('click', () => runPrimaryQuery(false));
      els.forgetButton.addEventListener('click', () => {
        localStorage.removeItem(STORAGE_KEY);
        els.tokenInput.value = '';
        els.tokenInput.placeholder = '';
        setStatus('Stored token removed from local browser storage.', 'warning');
      });
      els.itineraryTableBody.addEventListener('click', event => {
        const row = event.target.closest('tr[data-key]');
        if (!row) return;
        state.selectedKey = row.getAttribute('data-key');
        renderTable();
        renderDetailPanel();
        renderMapForSelection(state.airportLookup.size ? 'ok' : 'pending');
      });
      els.flightFilter.addEventListener('change', event => {
        state.currentFilterFlight = event.target.value;
        applyFilters();
        renderTable();
        renderDetailPanel();
        renderMapForSelection(state.enrichmentStatus === 'failed' ? 'failed' : (state.airportLookup.size ? 'ok' : 'pending'));
      });
      els.yearFilter.addEventListener('change', event => {
        state.currentFilterYear = event.target.value;
        applyFilters();
        renderTable();
        renderDetailPanel();
        renderMapForSelection(state.enrichmentStatus === 'failed' ? 'failed' : (state.airportLookup.size ? 'ok' : 'pending'));
      });
      els.exportButton.addEventListener('click', exportFilteredRows);

      prefillToken();
      ensureMap();
      const storedToken = getStoredToken();
      if (storedToken) {
        runPrimaryQuery(true);
      }
    })();
  </script>
</body>
</html>
```