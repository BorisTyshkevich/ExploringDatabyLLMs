```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Highest Daily Hops — OnTime Analyst</title>
  <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"/>
  <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
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
      --border: rgba(22,50,68,0.10);
      --shadow: 0 18px 45px rgba(14,58,82,0.10);
      --radius-xl: 22px;
      --radius-lg: 16px;
      --radius-md: 12px;
      --radius-sm: 8px;
    }
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: linear-gradient(180deg, var(--bg-top) 0%, var(--bg-bottom) 100%);
      min-height: 100vh;
      font-family: "Segoe UI", "Helvetica Neue", Arial, sans-serif;
      color: var(--ink);
    }

    /* ── Header ── */
    header {
      background: linear-gradient(135deg, #0b2d42 0%, var(--navy) 45%, #1a5272 100%);
      color: #fff;
      padding: 2rem 2.5rem 1.75rem;
    }
    .header-inner { max-width: 1280px; margin: 0 auto; }
    header h1 {
      font-family: Georgia, ui-serif, serif;
      font-size: 1.8rem;
      font-weight: 700;
      letter-spacing: -0.02em;
      margin-bottom: 0.4rem;
    }
    header p { color: rgba(255,255,255,0.65); font-size: 0.92rem; }

    /* ── Wrapper ── */
    .main-wrap { max-width: 1280px; margin: 0 auto; padding: 1.75rem 2rem 2rem; }

    /* ── KPI strip ── */
    .kpi-strip {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(190px, 1fr));
      gap: 1rem;
      margin-bottom: 1.5rem;
    }
    .kpi-card {
      background: var(--panel);
      border-radius: var(--radius-lg);
      box-shadow: var(--shadow);
      padding: 1.1rem 1.4rem;
      border: 1px solid var(--border);
    }
    .kpi-label {
      font-size: 0.73rem;
      text-transform: uppercase;
      letter-spacing: 0.09em;
      color: var(--muted);
      margin-bottom: 0.4rem;
    }
    .kpi-value {
      font-size: 2rem;
      font-weight: 700;
      color: var(--navy);
      line-height: 1;
    }
    .kpi-sub { font-size: 0.8rem; color: var(--slate); margin-top: 0.35rem; }

    /* ── Panels ── */
    .panel {
      background: var(--panel);
      border-radius: var(--radius-lg);
      box-shadow: var(--shadow);
      border: 1px solid var(--border);
      overflow: hidden;
    }
    .panel-header {
      background: var(--panel-alt);
      padding: 0.7rem 1.25rem;
      border-bottom: 1px solid var(--border);
      font-weight: 600;
      font-size: 0.88rem;
      color: var(--navy);
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }
    .panel-body { padding: 1rem 1.25rem; }

    /* ── Content grid ── */
    .content-grid {
      display: grid;
      grid-template-columns: 1fr 360px;
      gap: 1.25rem;
      margin-bottom: 1.25rem;
      align-items: start;
    }
    @media (max-width: 900px) { .content-grid { grid-template-columns: 1fr; } }

    /* ── Map ── */
    #map-container { position: relative; }
    #map { height: 430px; width: 100%; }
    #map-degraded {
      height: 430px;
      display: none;
      align-items: center;
      justify-content: center;
      background: var(--panel-alt);
      color: var(--muted);
      font-size: 0.9rem;
      text-align: center;
      padding: 1.5rem;
      flex-direction: column;
      gap: 0.5rem;
    }
    .map-legend {
      display: flex;
      gap: 1rem;
      flex-wrap: wrap;
      align-items: center;
      padding: 0.55rem 1.1rem;
      background: var(--panel-alt);
      border-top: 1px solid var(--border);
      font-size: 0.78rem;
    }
    .legend-item { display: flex; align-items: center; gap: 0.35rem; }
    .legend-dot { width: 11px; height: 11px; border-radius: 50%; }
    .legend-line { width: 22px; height: 3px; border-radius: 2px; }

    /* ── Route detail ── */
    .route-panel {
      max-height: 520px;
      overflow-y: auto;
      display: flex;
      flex-direction: column;
    }
    .stop-item {
      display: flex;
      align-items: flex-start;
      padding: 0.6rem 1.25rem;
      border-bottom: 1px solid var(--border);
      gap: 0.75rem;
    }
    .stop-item:last-child { border-bottom: none; }
    .stop-dot {
      width: 28px; height: 28px;
      border-radius: 50%;
      display: flex; align-items: center; justify-content: center;
      font-size: 0.65rem; font-weight: 700;
      flex-shrink: 0; margin-top: 2px;
    }
    .stop-dot.origin { background: var(--teal); color: #fff; }
    .stop-dot.waypoint { background: var(--sky); color: #fff; }
    .stop-dot.dest { background: var(--red); color: #fff; }
    .stop-code { font-weight: 700; font-size: 1rem; color: var(--navy); }
    .stop-time { font-size: 0.79rem; color: var(--muted); margin-top: 1px; }
    .stop-name { font-size: 0.75rem; color: var(--slate); margin-top: 1px; }
    .stop-arrow {
      display: flex; align-items: center; justify-content: center;
      padding: 0.25rem 1.25rem;
      font-size: 0.7rem; color: var(--muted);
      border-bottom: 1px solid var(--border);
      background: var(--panel-alt);
    }

    /* ── Itinerary table ── */
    .table-wrap { overflow-x: auto; }
    .itin-table { width: 100%; border-collapse: collapse; font-size: 0.85rem; }
    .itin-table th {
      background: var(--navy); color: #fff;
      padding: 0.55rem 0.75rem;
      text-align: left; font-weight: 600;
      font-size: 0.75rem; text-transform: uppercase;
      letter-spacing: 0.05em;
      position: sticky; top: 0;
      white-space: nowrap;
    }
    .itin-table td {
      padding: 0.55rem 0.75rem;
      border-bottom: 1px solid var(--border);
      vertical-align: middle;
    }
    .itin-table tbody tr { cursor: pointer; transition: background 0.12s; }
    .itin-table tbody tr:hover td { background: var(--panel-alt); }
    .itin-table tbody tr.active td { background: rgba(60,136,181,0.13); }
    .itin-table tbody tr.active { border-left: 3px solid var(--sky); }
    .itin-table tbody tr.active td:first-child { font-weight: 700; }
    .badge {
      display: inline-block; padding: 0.15rem 0.5rem;
      border-radius: 99px; font-size: 0.72rem; font-weight: 700;
    }
    .badge-wn { background: #f0c040; color: #163244; }
    .hops-val { font-weight: 700; font-size: 1.05rem; color: var(--red); }
    .route-cell {
      font-size: 0.78rem; max-width: 300px;
      white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
    }
    .times-cell { font-size: 0.76rem; color: var(--muted); }

    /* ── Query ledger ── */
    .ledger-table { width: 100%; border-collapse: collapse; font-size: 0.82rem; }
    .ledger-table th {
      background: var(--panel-alt); padding: 0.45rem 0.75rem;
      text-align: left; font-weight: 600; font-size: 0.73rem;
      color: var(--muted); border-bottom: 1px solid var(--border);
    }
    .ledger-table td { padding: 0.45rem 0.75rem; border-bottom: 1px solid var(--border); vertical-align: middle; }
    .ledger-row { cursor: pointer; }
    .ledger-row:hover td { background: var(--panel-alt); }
    .ledger-sql-row td { padding: 0 0.75rem 0.5rem; background: var(--panel-alt); }
    .ledger-sql-block {
      background: #eef2f5; border-radius: var(--radius-sm);
      padding: 0.75rem 0.9rem; font-family: "Courier New", monospace;
      font-size: 0.76rem; white-space: pre-wrap; word-break: break-all;
      color: var(--ink); border: 1px solid var(--border);
    }
    .st-ok { color: var(--teal); font-weight: 600; }
    .st-pending { color: var(--amber); font-weight: 600; }
    .st-failed { color: var(--red); font-weight: 600; }
    .st-degraded { color: var(--amber); font-weight: 600; }

    /* ── Loading / empty states ── */
    .state-msg {
      text-align: center; padding: 3rem 2rem;
      color: var(--muted); font-size: 0.95rem;
    }
    #loading-state, #empty-state { display: none; }
    #dashboard-content { display: none; }
    .mb { margin-bottom: 1.25rem; }

    /* ── Footer ── */
    footer {
      background: linear-gradient(135deg, #0b2d42 0%, var(--navy) 100%);
      color: rgba(255,255,255,0.85);
      padding: 1.75rem 2.5rem 2rem;
      margin-top: 2.5rem;
    }
    footer h3 {
      font-size: 0.82rem; font-weight: 600;
      text-transform: uppercase; letter-spacing: 0.08em;
      color: rgba(255,255,255,0.5);
      margin-bottom: 1rem;
    }
    .footer-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1.25rem;
      align-items: start;
      max-width: 1280px;
    }
    @media (max-width: 700px) { .footer-grid { grid-template-columns: 1fr; } }
    .footer-grid label {
      display: block; font-size: 0.78rem;
      color: rgba(255,255,255,0.55); margin-bottom: 0.35rem;
    }
    .footer-grid input[type="password"],
    .footer-grid textarea {
      width: 100%;
      background: rgba(255,255,255,0.07);
      border: 1px solid rgba(255,255,255,0.18);
      border-radius: var(--radius-sm);
      color: #fff; padding: 0.6rem 0.85rem;
      font-size: 0.85rem; font-family: inherit; outline: none;
    }
    .footer-grid textarea {
      font-family: "Courier New", monospace;
      font-size: 0.76rem; resize: vertical; min-height: 140px;
    }
    .footer-actions {
      display: flex; gap: 0.75rem;
      align-items: center; margin-top: 1rem; flex-wrap: wrap;
    }
    .btn {
      padding: 0.55rem 1.3rem; border-radius: var(--radius-sm);
      font-size: 0.85rem; font-weight: 600;
      cursor: pointer; border: none; transition: opacity 0.15s;
    }
    .btn:disabled { opacity: 0.45; cursor: not-allowed; }
    .btn-primary { background: var(--sky); color: #fff; }
    .btn-secondary { background: rgba(255,255,255,0.12); color: rgba(255,255,255,0.8); }
    .btn-sm { padding: 0.35rem 0.85rem; font-size: 0.78rem; }
    #status-text { font-size: 0.84rem; color: rgba(255,255,255,0.6); }
  </style>
</head>
<body>

<header>
  <div class="header-inner">
    <h1>✈ Highest Daily Hops for One Aircraft on One Flight Number</h1>
    <p>Southwest Airlines point-to-point routing — top 10 itineraries chaining up to 8 legs under a single flight number in one calendar day</p>
  </div>
</header>

<div class="main-wrap">

  <!-- KPI strip — anchored to top-ranked result -->
  <div class="kpi-strip">
    <div class="kpi-card">
      <div class="kpi-label">Max Daily Hops</div>
      <div class="kpi-value" id="kpi-hops">—</div>
      <div class="kpi-sub">legs in one calendar day</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Tail Number</div>
      <div class="kpi-value" id="kpi-tail" style="font-size:1.35rem;letter-spacing:0.02em;">—</div>
      <div class="kpi-sub">aircraft on top-ranked day</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Flight Number</div>
      <div class="kpi-value" id="kpi-flight" style="font-size:1.5rem;">—</div>
      <div class="kpi-sub" id="kpi-date-sub">—</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Route Repeats</div>
      <div class="kpi-value" id="kpi-repeats" style="font-size:1.5rem;">—</div>
      <div class="kpi-sub" id="kpi-repeats-sub">distinct dates in top 10</div>
    </div>
  </div>

  <!-- Loading / empty states -->
  <div id="loading-state" class="state-msg">Loading itinerary data…</div>
  <div id="empty-state" class="state-msg">No results returned. Check your token and SQL.</div>

  <!-- Dashboard content -->
  <div id="dashboard-content">

    <!-- Map + route detail -->
    <div class="content-grid mb">

      <!-- Map panel -->
      <div class="panel">
        <div class="panel-header">
          <span>✈ Itinerary Map</span>
          <span id="map-subtitle" style="margin-left:auto;font-size:0.78rem;color:var(--muted);font-weight:400;">Select a row to plot route</span>
        </div>
        <div id="map-container">
          <div id="map"></div>
          <div id="map-degraded">
            <span style="font-size:1.5rem;">🗺</span>
            <span id="map-degraded-msg">Airport coordinates are loading…</span>
          </div>
        </div>
        <div class="map-legend">
          <div class="legend-item">
            <div class="legend-dot" style="background:var(--teal);"></div>
            <span>Origin</span>
          </div>
          <div class="legend-item">
            <div class="legend-dot" style="background:var(--sky);"></div>
            <span>Intermediate stop</span>
          </div>
          <div class="legend-item">
            <div class="legend-dot" style="background:var(--red);"></div>
            <span>Final destination</span>
          </div>
          <div class="legend-item">
            <div class="legend-line" style="background:var(--red);opacity:0.85;"></div>
            <span>Active route</span>
          </div>
        </div>
      </div>

      <!-- Route sequence panel -->
      <div class="panel route-panel">
        <div class="panel-header">🗒 Route Sequence &amp; Departure Times</div>
        <div id="route-detail">
          <div class="panel-body" style="color:var(--muted);font-size:0.88rem;">Select an itinerary row to see the stop sequence.</div>
        </div>
      </div>

    </div><!-- /content-grid -->

    <!-- Itinerary table -->
    <div class="panel mb">
      <div class="panel-header">
        <span>📋 Top 10 Longest Daily Itineraries</span>
        <button id="export-btn" class="btn btn-primary btn-sm" style="margin-left:auto;">Export CSV</button>
      </div>
      <div class="table-wrap">
        <table class="itin-table">
          <thead>
            <tr>
              <th style="width:2em;">#</th>
              <th>Tail</th>
              <th>Carrier</th>
              <th>Flight</th>
              <th>Date</th>
              <th style="width:4em;">Hops</th>
              <th>Route</th>
              <th>Departure Times</th>
            </tr>
          </thead>
          <tbody id="itin-tbody"></tbody>
        </table>
      </div>
    </div>

    <!-- Query ledger -->
    <div class="panel">
      <div class="panel-header">🔍 Query Ledger</div>
      <div class="table-wrap">
        <table class="ledger-table">
          <thead>
            <tr>
              <th style="width:1.6em;"></th>
              <th>Label</th>
              <th style="width:9em;">Role</th>
              <th style="width:7em;">Status</th>
              <th style="width:5em;text-align:right;">Rows</th>
            </tr>
          </thead>
          <tbody id="ledger-tbody"></tbody>
        </table>
      </div>
    </div>

  </div><!-- /dashboard-content -->
</div><!-- /main-wrap -->

<!-- Footer: JWE + SQL controls -->
<footer>
  <h3>Data Controls</h3>
  <div class="footer-grid">
    <div>
      <label for="jwe-input">JWE Token</label>
      <input type="password" id="jwe-input" placeholder="Paste your JWE token here…" autocomplete="off" spellcheck="false">
    </div>
    <div>
      <label for="sql-input">Analytical SQL (editable)</label>
      <textarea id="sql-input" spellcheck="false"></textarea>
    </div>
  </div>
  <div class="footer-actions">
    <button class="btn btn-primary" id="fetch-btn">Fetch Data</button>
    <button class="btn btn-secondary" id="forget-btn">Forget Token</button>
    <span id="status-text">Enter your JWE token and click Fetch Data.</span>
  </div>
</footer>

<script>
// ============================================================
// CONSTANTS
// ============================================================
const JWE_KEY = 'OnTimeAnalystDashboard::auth::jwe';
const BASE_URL = 'https://mcp.demo.altinity.cloud';

const SAVED_SQL = `WITH legs AS (
    SELECT
        Tail_Number,
        IATA_CODE_Reporting_Airline AS Carrier,
        Flight_Number_Reporting_Airline AS FlightNum,
        FlightDate,
        Origin,
        Dest,
        DepTime
    FROM ontime.ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND DepTime IS NOT NULL
),
grouped AS (
    SELECT
        Tail_Number,
        Carrier,
        FlightNum,
        FlightDate,
        count() AS Hops,
        arraySort(x -> tupleElement(x, 1), groupArray((DepTime, Origin, Dest))) AS sorted_legs
    FROM legs
    GROUP BY Tail_Number, Carrier, FlightNum, FlightDate
)
SELECT
    Tail_Number,
    Carrier,
    FlightNum,
    FlightDate,
    Hops,
    arrayStringConcat(
        arrayConcat(
            arrayMap(x -> tupleElement(x, 2), sorted_legs),
            [tupleElement(sorted_legs[length(sorted_legs)], 3)]
        ),
        ' -> '
    ) AS Route,
    arrayStringConcat(
        arrayMap(x -> concat(
            leftPad(toString(intDiv(tupleElement(x, 1), 100)), 2, '0'),
            ':',
            leftPad(toString(tupleElement(x, 1) % 100), 2, '0')
        ), sorted_legs),
        ', '
    ) AS DepartureTimes
FROM grouped
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10`;

// ============================================================
// STATE
// ============================================================
let rows = [];
let airportCoords = {};   // code -> { lat, lon, name }
let enrichmentDone = false;
let selectedIdx = 0;
let leafletMap = null;
let mapMarkers = [];
let mapLines = [];
let isRunning = false;
let activeRunId = 0;

// Query ledger
let ledger = [];
let ledgerSeq = 0;

// ============================================================
// BOOT
// ============================================================
document.addEventListener('DOMContentLoaded', () => {
  document.getElementById('sql-input').value = SAVED_SQL;

  const stored = localStorage.getItem(JWE_KEY);
  if (stored) {
    document.getElementById('jwe-input').value = stored;
    runFetch();
  }

  document.getElementById('fetch-btn').addEventListener('click', runFetch);
  document.getElementById('forget-btn').addEventListener('click', forgetToken);
  document.getElementById('export-btn').addEventListener('click', exportCsv);
});

// ============================================================
// FETCH FLOW
// ============================================================
async function runFetch() {
  if (isRunning) return;

  const jweEl = document.getElementById('jwe-input');
  const sqlEl = document.getElementById('sql-input');
  const fetchBtn = document.getElementById('fetch-btn');
  const statusEl = document.getElementById('status-text');

  const token = jweEl.value.trim();
  const sql = sqlEl.value.trim();
  if (!token) { statusEl.textContent = 'JWE token required.'; return; }
  if (!sql)   { statusEl.textContent = 'SQL required.'; return; }

  isRunning = true;
  activeRunId++;
  const myRun = activeRunId;
  fetchBtn.disabled = true;
  statusEl.textContent = 'Fetching…';

  // Reset UI
  rows = [];
  airportCoords = {};
  enrichmentDone = false;
  ledger = [];
  ledgerSeq = 0;
  document.getElementById('loading-state').style.display = 'block';
  document.getElementById('dashboard-content').style.display = 'none';
  document.getElementById('empty-state').style.display = 'none';
  renderLedger();

  const primaryId = addLedgerEntry('Primary Itinerary Query', 'Primary', 'Pending', null, sql);

  try {
    const url = buildUrl(token, sql);
    const resp = await fetch(url);
    if (myRun !== activeRunId) return;

    if (!resp.ok) {
      const txt = await resp.text();
      updateLedgerEntry(primaryId, 'Failed', null);
      statusEl.textContent = `HTTP ${resp.status}: ${txt.slice(0, 200)}`;
      document.getElementById('loading-state').style.display = 'none';
      return;
    }

    const payload = await resp.json();
    if (myRun !== activeRunId) return;

    localStorage.setItem(JWE_KEY, token);

    const raw = toObjects(payload);
    rows = raw.map(normalizeRow);
    updateLedgerEntry(primaryId, 'OK', rows.length);

    document.getElementById('loading-state').style.display = 'none';

    if (rows.length === 0) {
      document.getElementById('empty-state').style.display = 'block';
      statusEl.textContent = 'Query returned 0 rows.';
      return;
    }

    // Show dashboard shell
    document.getElementById('dashboard-content').style.display = 'block';
    renderKpis();
    renderTable();
    selectedIdx = 0;
    setActiveRow(0);
    renderRouteDetail(rows[0]);

    // Init Leaflet now that map container is visible
    if (!leafletMap) initLeaflet();
    showMapPending('Loading airport coordinates…');

    statusEl.textContent = `Loaded ${rows.length} rows. Fetching airport coordinates…`;

    // Run coordinate enrichment (once for all airports)
    await runEnrichment(token, myRun);
    if (myRun !== activeRunId) return;

    // Render map for initially selected row
    renderMap(rows[selectedIdx]);
    renderRouteDetail(rows[selectedIdx]); // refresh names
    statusEl.textContent = `Ready — ${rows.length} itineraries loaded.`;

  } catch (err) {
    if (myRun !== activeRunId) return;
    updateLedgerEntry(primaryId, 'Failed', null);
    statusEl.textContent = `Error: ${err.message}`;
    document.getElementById('loading-state').style.display = 'none';
  } finally {
    if (myRun === activeRunId) {
      isRunning = false;
      fetchBtn.disabled = false;
    }
  }
}

// ============================================================
// ENRICHMENT
// ============================================================
async function runEnrichment(token, myRun) {
  const codes = new Set();
  rows.forEach(r => r.stops.forEach(c => codes.add(c)));
  const list = [...codes].filter(Boolean);
  if (!list.length) { showMapDegraded('No airport codes found in route data.'); return; }

  const escapedList = list.map(c => `'${c.replace(/'/g, "''")}'`).join(', ');
  const enrichSql = `SELECT code, latitude, longitude, name FROM ontime.airports_latest WHERE code IN (${escapedList})`;
  const enrichId = addLedgerEntry('Airport Coordinate Lookup', 'Enrichment', 'Pending', null, enrichSql);

  try {
    const url = buildUrl(token, enrichSql);
    const resp = await fetch(url);
    if (myRun !== activeRunId) return;

    if (!resp.ok) {
      const txt = await resp.text();
      updateLedgerEntry(enrichId, 'Failed', null);
      showMapDegraded(`Coordinate lookup failed (HTTP ${resp.status}): ${txt.slice(0, 120)}`);
      addLedgerEntry('Map render', 'Map Render', 'Degraded', null, '-- enrichment failed, no route plotted --');
      return;
    }

    const payload = await resp.json();
    if (myRun !== activeRunId) return;

    const enrichRows = toObjects(payload);
    airportCoords = {};
    enrichRows.forEach(r => {
      const lat = parseFloat(r.latitude);
      const lon = parseFloat(r.longitude);
      if (!isNaN(lat) && !isNaN(lon)) {
        airportCoords[r.code] = { lat, lon, name: r.name || r.code };
      }
    });
    enrichmentDone = true;
    updateLedgerEntry(enrichId, 'OK', enrichRows.length);

  } catch (err) {
    if (myRun !== activeRunId) return;
    updateLedgerEntry(enrichId, 'Failed', null);
    showMapDegraded(`Coordinate lookup error: ${err.message}`);
  }
}

// ============================================================
// KPIs — anchored to top-ranked row (row 0), never change on selection
// ============================================================
function renderKpis() {
  const top = rows[0];
  if (!top) return;

  document.getElementById('kpi-hops').textContent = top.hops;
  document.getElementById('kpi-tail').textContent = top.tail;
  document.getElementById('kpi-flight').textContent = `${top.carrier} ${top.flightNum}`;
  document.getElementById('kpi-date-sub').textContent = `Date: ${top.dateKey}`;

  // Count distinct dates across top 10
  const distinctDates = new Set(rows.map(r => r.dateKey)).size;
  document.getElementById('kpi-repeats').textContent = distinctDates;
  document.getElementById('kpi-repeats-sub').textContent =
    `distinct dates in top ${rows.length} — repeating schedule`;
}

// ============================================================
// ITINERARY TABLE
// ============================================================
function renderTable() {
  const tbody = document.getElementById('itin-tbody');
  tbody.innerHTML = '';
  rows.forEach((r, i) => {
    const tr = document.createElement('tr');
    tr.dataset.idx = i;
    tr.innerHTML = `
      <td>${i + 1}</td>
      <td><code style="font-size:0.82rem;">${esc(r.tail)}</code></td>
      <td><span class="badge badge-wn">${esc(r.carrier)}</span></td>
      <td><strong>${esc(r.flightNum)}</strong></td>
      <td>${esc(r.dateKey)}</td>
      <td class="hops-val">${r.hops}</td>
      <td class="route-cell" title="${esc(r.route)}">${esc(r.route)}</td>
      <td class="times-cell">${esc(r.depTimes)}</td>
    `;
    tr.addEventListener('click', () => selectRow(i));
    tbody.appendChild(tr);
  });
}

function selectRow(idx) {
  selectedIdx = idx;
  setActiveRow(idx);
  renderRouteDetail(rows[idx]);
  renderMap(rows[idx]);
}

function setActiveRow(idx) {
  document.querySelectorAll('#itin-tbody tr').forEach((tr, i) => {
    tr.classList.toggle('active', i === idx);
  });
}

// ============================================================
// ROUTE DETAIL PANEL
// ============================================================
function renderRouteDetail(row) {
  const container = document.getElementById('route-detail');
  container.innerHTML = '';

  const times = row.depTimes.split(', ');

  row.stops.forEach((code, i) => {
    const isOrigin = i === 0;
    const isDest = i === row.stops.length - 1;
    const hasDepart = !isDest;

    const dotClass = isOrigin ? 'origin' : (isDest ? 'dest' : 'waypoint');
    const dotLabel = isOrigin ? 'O' : (isDest ? '●' : String(i));
    const depTime = hasDepart ? (times[i] ?? '') : '';
    const coordInfo = airportCoords[code];
    const nameStr = coordInfo ? `<div class="stop-name">${esc(coordInfo.name)}</div>` : '';
    const typeTag = isOrigin ? ' <small style="color:var(--teal);font-size:0.7rem;">(origin)</small>'
                  : isDest   ? ' <small style="color:var(--red);font-size:0.7rem;">(destination)</small>'
                  : '';

    const div = document.createElement('div');
    div.className = 'stop-item';
    div.innerHTML = `
      <div class="stop-dot ${dotClass}">${dotLabel}</div>
      <div>
        <div class="stop-code">${esc(code)}${typeTag}</div>
        ${depTime ? `<div class="stop-time">Dep: ${esc(depTime)}</div>` : ''}
        ${nameStr}
      </div>
    `;
    container.appendChild(div);

    if (!isDest) {
      const arrow = document.createElement('div');
      arrow.className = 'stop-arrow';
      arrow.textContent = '↓';
      container.appendChild(arrow);
    }
  });
}

// ============================================================
// MAP
// ============================================================
function initLeaflet() {
  leafletMap = L.map('map', { zoomControl: true }).setView([39.5, -98.35], 4);
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors',
    maxZoom: 18
  }).addTo(leafletMap);
  leafletMap.invalidateSize();
}

function clearMapLayers() {
  mapMarkers.forEach(m => leafletMap?.removeLayer(m));
  mapLines.forEach(l => leafletMap?.removeLayer(l));
  mapMarkers = [];
  mapLines = [];
}

function showMapPending(msg) {
  document.getElementById('map').style.display = 'none';
  const d = document.getElementById('map-degraded');
  d.style.display = 'flex';
  document.getElementById('map-degraded-msg').textContent = msg;
}

function showMapDegraded(msg) {
  document.getElementById('map').style.display = 'none';
  const d = document.getElementById('map-degraded');
  d.style.display = 'flex';
  document.getElementById('map-degraded-msg').textContent = msg;
}

function showMapActive() {
  document.getElementById('map').style.display = 'block';
  document.getElementById('map-degraded').style.display = 'none';
  leafletMap?.invalidateSize();
}

function renderMap(row) {
  if (!leafletMap) return;
  clearMapLayers();

  if (!enrichmentDone) {
    showMapPending('Airport coordinates loading…');
    return;
  }

  const validCount = row.stops.filter(c => airportCoords[c]).length;

  if (validCount < 2) {
    showMapDegraded(`Not enough coordinates for route: ${row.route}`);
    const existing = ledger.find(e => e.role === 'Map Render');
    if (!existing) {
      addLedgerEntry(`Map degraded — missing coords for: ${row.route}`, 'Map Render', 'Degraded', null, '-- insufficient coordinate data --');
    } else {
      existing.status = 'Degraded';
      existing.label = `Map degraded — missing coords for: ${row.route}`;
      renderLedger();
    }
    return;
  }

  showMapActive();
  document.getElementById('map-subtitle').textContent =
    `${row.carrier}${row.flightNum} · ${row.dateKey} · ${row.hops} hops`;

  // Draw leg polylines
  row.stops.forEach((code, i) => {
    if (i === row.stops.length - 1) return;
    const c1 = airportCoords[row.stops[i]];
    const c2 = airportCoords[row.stops[i + 1]];
    if (!c1 || !c2) return;
    const line = L.polyline([[c1.lat, c1.lon], [c2.lat, c2.lon]], {
      color: '#c54f36',
      weight: 2.8,
      opacity: 0.82
    }).addTo(leafletMap);
    mapLines.push(line);
  });

  // Draw airport markers
  const times = row.depTimes.split(', ');
  row.stops.forEach((code, i) => {
    const c = airportCoords[code];
    if (!c) return;
    const isOrigin = i === 0;
    const isDest = i === row.stops.length - 1;
    const color = isOrigin ? '#1f8a70' : (isDest ? '#c54f36' : '#3c88b5');
    const radius = isOrigin || isDest ? 8 : 6;
    const depTime = times[i] ? ` — Dep: ${times[i]}` : '';

    const marker = L.circleMarker([c.lat, c.lon], {
      radius, fillColor: color,
      color: '#fff', weight: 2,
      opacity: 1, fillOpacity: 0.95
    }).addTo(leafletMap);

    const roleLabel = isOrigin ? 'Origin' : (isDest ? 'Destination' : `Stop ${i}`);
    marker.bindPopup(`<strong>${esc(code)}</strong> <small>(${roleLabel})</small><br>${esc(c.name)}${depTime}`);
    mapMarkers.push(marker);
  });

  // Fit bounds
  const pts = row.stops.filter(c => airportCoords[c]).map(c => [airportCoords[c].lat, airportCoords[c].lon]);
  if (pts.length >= 2) {
    try { leafletMap.fitBounds(L.latLngBounds(pts), { padding: [45, 45] }); }
    catch (_) { /* ignore */ }
  }

  // Update ledger
  const existing = ledger.find(e => e.role === 'Map Render');
  if (!existing) {
    addLedgerEntry(`Map rendered — ${row.carrier}${row.flightNum} ${row.dateKey}`, 'Map Render', 'OK', null, '-- Leaflet render, no query --');
  } else {
    existing.status = 'OK';
    existing.label = `Map rendered — ${row.carrier}${row.flightNum} ${row.dateKey}`;
    renderLedger();
  }
}

// ============================================================
// QUERY LEDGER
// ============================================================
function addLedgerEntry(label, role, status, rowCount, sql) {
  const id = ++ledgerSeq;
  ledger.push({ id, label, role, status, rowCount, sql: sql ?? '', expanded: false });
  renderLedger();
  return id;
}

function updateLedgerEntry(id, status, rowCount) {
  const e = ledger.find(x => x.id === id);
  if (!e) return;
  e.status = status;
  if (rowCount !== null) e.rowCount = rowCount;
  renderLedger();
}

function renderLedger() {
  const tbody = document.getElementById('ledger-tbody');
  tbody.innerHTML = '';
  ledger.forEach(entry => {
    const stCls = entry.status === 'OK' ? 'st-ok'
                : entry.status === 'Pending' ? 'st-pending'
                : entry.status === 'Degraded' ? 'st-degraded'
                : 'st-failed';
    const rowsStr = entry.rowCount !== null && entry.rowCount !== undefined ? entry.rowCount : '—';

    const tr = document.createElement('tr');
    tr.className = 'ledger-row';
    tr.innerHTML = `
      <td style="font-family:monospace;color:var(--muted);user-select:none;">${entry.expanded ? '▼' : '▶'}</td>
      <td>${esc(entry.label)}</td>
      <td style="color:var(--slate);font-size:0.77rem;">${esc(entry.role)}</td>
      <td class="${stCls}">${esc(entry.status)}</td>
      <td style="text-align:right;">${rowsStr}</td>
    `;
    tr.addEventListener('click', () => { entry.expanded = !entry.expanded; renderLedger(); });
    tbody.appendChild(tr);

    if (entry.expanded) {
      const sqlTr = document.createElement('tr');
      sqlTr.className = 'ledger-sql-row';
      sqlTr.innerHTML = `<td colspan="5"><pre class="ledger-sql-block">${esc(entry.sql)}</pre></td>`;
      tbody.appendChild(sqlTr);
    }
  });
}

// ============================================================
// EXPORT
// ============================================================
function exportCsv() {
  if (!rows.length) return;
  const hdr = 'Tail_Number,Carrier,FlightNum,FlightDate,Hops,Route,DepartureTimes';
  const lines = rows.map(r =>
    [r.tail, r.carrier, r.flightNum, r.dateKey, r.hops,
     `"${r.route.replace(/"/g,'""')}"`,
     `"${r.depTimes.replace(/"/g,'""')}"`].join(',')
  );
  const blob = new Blob([[hdr, ...lines].join('\n')], { type: 'text/csv' });
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = 'hops_per_day_top10.csv';
  a.click();
}

// ============================================================
// FORGET TOKEN
// ============================================================
function forgetToken() {
  localStorage.removeItem(JWE_KEY);
  document.getElementById('jwe-input').value = '';
  document.getElementById('status-text').textContent = 'Token removed from local storage.';
}

// ============================================================
// HELPERS
// ============================================================
function buildUrl(token, sql) {
  return `${BASE_URL}/${encodeURIComponent(token)}/openapi/execute_query?query=${encodeURIComponent(sql)}`;
}

function toObjects(payload) {
  const cols = payload?.columns ?? [];
  const rawRows = payload?.rows;
  if (!cols.length || !rawRows || rawRows.length === 0) return [];
  return rawRows.map(r => {
    if (Array.isArray(r)) {
      const obj = {};
      cols.forEach((c, i) => { obj[c] = r[i]; });
      return obj;
    }
    return r;
  });
}

function normalizeRow(r) {
  const rawDate = r.FlightDate ?? r.flightdate ?? r.flight_date ?? '';
  const dateKey = String(rawDate).slice(0, 10);
  const route = r.Route ?? r.route ?? '';
  const stops = route.split(' -> ').map(s => s.trim()).filter(Boolean);
  return {
    tail:      String(r.Tail_Number ?? r.tail_number ?? ''),
    carrier:   String(r.Carrier ?? r.carrier ?? ''),
    flightNum: String(r.FlightNum ?? r.flightnum ?? r.flight_num ?? ''),
    dateKey,
    hops:      parseInt(r.Hops ?? r.hops ?? 0, 10),
    route,
    stops,
    depTimes:  String(r.DepartureTimes ?? r.departuretimes ?? r.departure_times ?? '')
  };
}

function esc(s) {
  return String(s ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
</script>
</body>
</html>
```
