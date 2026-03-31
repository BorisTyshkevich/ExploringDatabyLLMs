```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Highest Daily Hops — OnTime Analyst</title>
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" crossorigin=""/>
<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js" crossorigin=""></script>
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
  --grid: rgba(22,50,68,0.12);
  --border: rgba(22,50,68,0.10);
  --shadow: 0 18px 45px rgba(14,58,82,0.10);
  --radius-xl: 22px;
  --radius-lg: 16px;
  --radius-md: 12px;
  --radius-sm: 8px;
}
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:"Segoe UI","Helvetica Neue",sans-serif;background:linear-gradient(160deg,var(--bg-top) 0%,var(--bg-bottom) 100%);min-height:100vh;color:var(--ink)}
header{background:linear-gradient(135deg,var(--navy) 0%,#1a5a7a 100%);color:#fff;padding:28px 40px 24px}
.header-inner{max-width:1280px;margin:0 auto}
.header-inner h1{font-family:Georgia,ui-serif,serif;font-size:1.9rem;font-weight:700;letter-spacing:-0.02em;display:flex;align-items:center;gap:10px}
.header-inner p{margin-top:6px;color:rgba(255,255,255,0.72);font-size:0.92rem}
.main{max-width:1280px;margin:0 auto;padding:28px 40px 40px}
.pre-load{background:var(--panel);border-radius:var(--radius-xl);padding:60px;text-align:center;box-shadow:var(--shadow);margin-bottom:24px}
.pre-load h2{font-family:Georgia,ui-serif,serif;color:var(--navy);margin-bottom:10px}
.pre-load p{color:var(--muted)}
#dashboard{display:none}
.kpi-strip{display:grid;grid-template-columns:repeat(auto-fit,minmax(155px,1fr));gap:16px;margin-bottom:24px}
.kpi-card{background:var(--panel);border-radius:var(--radius-lg);padding:18px 20px;box-shadow:var(--shadow);border:1px solid var(--border)}
.kpi-label{font-size:0.7rem;text-transform:uppercase;letter-spacing:0.08em;color:var(--muted);margin-bottom:6px}
.kpi-value{font-family:Georgia,ui-serif,serif;font-size:1.65rem;font-weight:700;color:var(--navy);line-height:1.1}
.kpi-sub{font-size:0.73rem;color:var(--slate);margin-top:4px}
.map-detail-grid{display:grid;grid-template-columns:1fr 340px;gap:20px;margin-bottom:24px}
@media(max-width:900px){.map-detail-grid{grid-template-columns:1fr}}
.card{background:var(--panel);border-radius:var(--radius-xl);padding:24px;box-shadow:var(--shadow);border:1px solid var(--border)}
.card-title{font-family:Georgia,ui-serif,serif;font-size:1rem;font-weight:700;color:var(--navy);margin-bottom:16px;display:flex;align-items:center;gap:8px;flex-wrap:wrap}
.badge{background:var(--sky);color:#fff;font-family:"Segoe UI",sans-serif;font-size:0.68rem;font-weight:600;padding:2px 8px;border-radius:20px;letter-spacing:0.05em}
#map-container{height:400px;border-radius:var(--radius-lg);overflow:hidden;background:#dce9f0;position:relative}
.map-legend{margin-top:12px;display:flex;gap:18px;font-size:0.78rem;color:var(--muted);flex-wrap:wrap}
.legend-item{display:flex;align-items:center;gap:6px}
.legend-dot{width:12px;height:12px;border-radius:50%;flex-shrink:0}
.legend-line{width:24px;height:3px;border-radius:2px;flex-shrink:0}
.map-degraded{display:none;background:#fff8e6;border:1px solid var(--amber);border-radius:var(--radius-md);padding:14px 16px;margin-top:12px;font-size:0.83rem;color:#7a5000}
.route-header{font-size:0.78rem;color:var(--muted);text-transform:uppercase;letter-spacing:0.07em;margin-bottom:12px}
.stop-list{list-style:none;max-height:360px;overflow-y:auto}
.stop-item{display:flex;align-items:flex-start;gap:10px;padding:8px 0;border-bottom:1px solid var(--border)}
.stop-item:last-child{border-bottom:none}
.stop-num{flex-shrink:0;width:22px;height:22px;border-radius:50%;background:var(--navy);color:#fff;font-size:0.68rem;font-weight:700;display:flex;align-items:center;justify-content:center;margin-top:1px}
.stop-num.stop-first{background:var(--teal)}
.stop-num.stop-last{background:var(--red)}
.stop-info{flex:1;min-width:0}
.stop-code{font-weight:700;font-size:0.95rem;color:var(--navy)}
.stop-name{font-size:0.72rem;color:var(--muted);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.stop-time{font-size:0.78rem;color:var(--sky);font-weight:600;margin-top:1px}
.table-card{margin-bottom:24px}
.export-row{display:flex;justify-content:flex-end;margin-bottom:10px}
.btn-export{background:var(--panel-alt);border:1px solid var(--border);border-radius:var(--radius-sm);padding:6px 14px;font-size:0.78rem;color:var(--navy);cursor:pointer}
.btn-export:hover{background:#dce9f0}
table{width:100%;border-collapse:collapse;font-size:0.85rem}
th{text-align:left;font-size:0.69rem;text-transform:uppercase;letter-spacing:0.07em;color:var(--muted);padding:10px 12px;border-bottom:2px solid var(--border);background:var(--panel-alt);white-space:nowrap}
td{padding:9px 12px;border-bottom:1px solid var(--border);color:var(--ink);vertical-align:top}
tr.data-row{cursor:pointer;transition:background 0.12s}
tr.data-row:hover td{background:#f0f7fc}
tr.data-row.active td{background:#e3f0fa;border-left:3px solid var(--sky)}
tr.data-row.active td:first-child{padding-left:9px}
.route-cell{font-size:0.78rem;color:var(--slate);max-width:300px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.hops-badge{display:inline-block;background:var(--navy);color:#fff;font-size:0.75rem;font-weight:700;padding:2px 8px;border-radius:10px}
.carrier-badge{display:inline-block;background:var(--sky);color:#fff;font-size:0.72rem;font-weight:600;padding:2px 7px;border-radius:8px}
.ledger-title{font-family:Georgia,ui-serif,serif;font-size:1rem;color:var(--navy);margin-bottom:12px;font-weight:700}
.ledger-table{width:100%;border-collapse:collapse;font-size:0.82rem}
.ledger-table th{font-size:0.68rem;text-transform:uppercase;letter-spacing:0.07em;color:var(--muted);padding:8px 10px;border-bottom:2px solid var(--border);background:var(--panel-alt);text-align:left}
.ledger-row{cursor:pointer}
.ledger-row:hover td{background:var(--panel-alt)}
.ledger-row td{padding:9px 10px;border-bottom:1px solid var(--border);vertical-align:top}
.ledger-toggle{font-family:monospace;font-size:0.82rem;color:var(--slate);user-select:none}
.ledger-sql{display:none;padding:10px;background:var(--panel-alt);border-radius:var(--radius-sm);margin-top:6px}
.ledger-sql pre{font-family:monospace;font-size:0.73rem;color:var(--ink);white-space:pre-wrap;word-break:break-all}
.status-ok{color:var(--teal);font-weight:600}
.status-pending{color:var(--amber);font-weight:600}
.status-failed{color:var(--red);font-weight:600}
footer{background:var(--navy);color:rgba(255,255,255,0.85);padding:28px 40px;margin-top:40px}
.footer-inner{max-width:1280px;margin:0 auto}
.footer-inner h3{font-family:Georgia,ui-serif,serif;font-size:1rem;color:#fff;margin-bottom:16px}
.footer-grid{display:grid;grid-template-columns:1fr 1fr;gap:20px}
@media(max-width:700px){.footer-grid{grid-template-columns:1fr}}
.footer-field label{display:block;font-size:0.7rem;text-transform:uppercase;letter-spacing:0.07em;color:rgba(255,255,255,0.6);margin-bottom:6px}
.footer-field input,.footer-field textarea{width:100%;background:rgba(255,255,255,0.08);border:1px solid rgba(255,255,255,0.18);border-radius:var(--radius-sm);color:#fff;padding:9px 12px;font-family:monospace;font-size:0.78rem;resize:vertical}
.footer-field textarea{min-height:130px}
.footer-actions{display:flex;gap:10px;margin-top:14px;align-items:center;flex-wrap:wrap}
.btn-primary{background:var(--sky);color:#fff;border:none;border-radius:var(--radius-sm);padding:9px 20px;font-size:0.85rem;font-weight:600;cursor:pointer;transition:opacity 0.15s}
.btn-primary:disabled{opacity:0.45;cursor:not-allowed}
.btn-secondary{background:transparent;color:rgba(255,255,255,0.7);border:1px solid rgba(255,255,255,0.3);border-radius:var(--radius-sm);padding:9px 16px;font-size:0.85rem;cursor:pointer}
.btn-secondary:hover{background:rgba(255,255,255,0.08)}
#status-text{font-size:0.82rem;color:rgba(255,255,255,0.7)}
.alert-warn{background:#fff8e6;border:1px solid var(--amber);border-radius:var(--radius-md);padding:14px 18px;margin-bottom:16px;font-size:0.85rem;color:#7a5000}
</style>
</head>
<body>

<header>
  <div class="header-inner">
    <h1>&#9992; Highest Daily Hops — One Aircraft, One Flight Number</h1>
    <p>Southwest Airlines dominates the top-10 with recurring 8-hop itineraries spanning the continental US in a single day</p>
  </div>
</header>

<div class="main">

  <div id="pre-load-state" class="pre-load">
    <h2>Enter your API token to load data</h2>
    <p>Supply your JWE token in the controls below and click Fetch Data.</p>
  </div>

  <div id="dashboard">

    <!-- KPI strip — anchored to the top-ranked row, never changes on row selection -->
    <div class="kpi-strip">
      <div class="kpi-card">
        <div class="kpi-label">Max Hops</div>
        <div class="kpi-value" id="kpi-hops">—</div>
        <div class="kpi-sub">legs flown in one day</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">Top Tail Number</div>
        <div class="kpi-value" id="kpi-tail" style="font-size:1.2rem">—</div>
        <div class="kpi-sub" id="kpi-tail-sub">—</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">Flight Number</div>
        <div class="kpi-value" id="kpi-fnum" style="font-size:1.35rem">—</div>
        <div class="kpi-sub">highest-hop flight</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">Date</div>
        <div class="kpi-value" id="kpi-date" style="font-size:1.1rem">—</div>
        <div class="kpi-sub">most recent max-hop day</div>
      </div>
      <div class="kpi-card">
        <div class="kpi-label">Recurring Routes</div>
        <div class="kpi-value" id="kpi-repeats">—</div>
        <div class="kpi-sub" id="kpi-repeats-sub">flight numbers seen 2+ days</div>
      </div>
    </div>

    <!-- Map + Route detail -->
    <div class="map-detail-grid">
      <div class="card">
        <div class="card-title">
          Route Map
          <span class="badge" id="map-badge">—</span>
        </div>
        <div id="map-container"></div>
        <div class="map-legend">
          <div class="legend-item"><div class="legend-dot" style="background:var(--teal)"></div>Origin</div>
          <div class="legend-item"><div class="legend-dot" style="background:var(--navy)"></div>Intermediate</div>
          <div class="legend-item"><div class="legend-dot" style="background:var(--red)"></div>Final destination</div>
          <div class="legend-item"><div class="legend-line" style="background:var(--sky);opacity:0.8"></div>Route legs</div>
        </div>
        <div class="map-degraded" id="map-degraded">Airport coordinates unavailable for this itinerary.</div>
      </div>

      <div class="card">
        <div class="card-title">Route Sequence</div>
        <div class="route-header" id="route-header">—</div>
        <ul class="stop-list" id="stop-list"></ul>
      </div>
    </div>

    <!-- Itinerary table -->
    <div class="card table-card">
      <div class="card-title">Top 10 Max-Hop Itineraries</div>
      <div class="export-row">
        <button class="btn-export" onclick="exportCsv()">Export CSV</button>
      </div>
      <div style="overflow-x:auto">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>Date</th>
              <th>Carrier</th>
              <th>Flight</th>
              <th>Tail</th>
              <th>Hops</th>
              <th>Route</th>
            </tr>
          </thead>
          <tbody id="itinerary-tbody"></tbody>
        </table>
      </div>
    </div>

    <!-- Query ledger -->
    <div class="card">
      <div class="ledger-title">Query Ledger</div>
      <table class="ledger-table">
        <thead>
          <tr>
            <th style="width:1.6em"></th>
            <th>Label</th>
            <th style="width:7em">Role</th>
            <th style="width:5.5em">Status</th>
            <th style="width:4em">Rows</th>
          </tr>
        </thead>
        <tbody id="ledger-tbody"></tbody>
      </table>
    </div>

  </div><!-- #dashboard -->

</div><!-- .main -->

<footer>
  <div class="footer-inner">
    <h3>Data Controls</h3>
    <div class="footer-grid">
      <div class="footer-field">
        <label>JWE Token</label>
        <input type="password" id="jwe-input" placeholder="Paste JWE token…" autocomplete="off">
      </div>
      <div class="footer-field">
        <label>SQL Query</label>
        <textarea id="sql-input" spellcheck="false"></textarea>
      </div>
    </div>
    <div class="footer-actions">
      <button class="btn-primary" id="fetch-btn" onclick="runPrimaryQuery()">Fetch Data</button>
      <button class="btn-secondary" onclick="forgetToken()">Forget Token</button>
      <span id="status-text"></span>
    </div>
  </div>
</footer>

<script>
// ── Constants ──────────────────────────────────────────────
const STORAGE_KEY = 'OnTimeAnalystDashboard::auth::jwe';
const BASE_URL    = 'https://mcp.demo.altinity.cloud';

const PRIMARY_SQL = `WITH itineraries AS (
    SELECT
        FlightDate,
        Tail_Number,
        IATA_CODE_Reporting_Airline AS Carrier,
        Flight_Number_Reporting_Airline AS FlightNum,
        count() AS Hops,
        arrayStringConcat(
            arrayConcat(
                arrayMap(
                    x -> x.2,
                    arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Origin)))
                ),
                [arrayElement(
                    arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Dest)))),
                    -1
                )]
            ),
            ' -> '
        ) AS Route,
        arrayStringConcat(
            arrayMap(
                x -> x.2 || '@' || toString(x.1),
                arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), Origin)))
            ),
            ', '
        ) AS DepartureTimes
    FROM ontime.ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND DepTime IS NOT NULL
    GROUP BY FlightDate, Tail_Number, Carrier, FlightNum
)
SELECT
    FlightDate,
    Tail_Number,
    Carrier,
    FlightNum,
    Hops,
    Route,
    DepartureTimes
FROM itineraries
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10`;

// ── State ──────────────────────────────────────────────────
let rows = [];
let airportCoords = {};   // code -> {lat, lon, name}
let selectedIdx   = 0;
let leafletMap    = null;
let routeLayerGroup = null;
let isRunning     = false;
let activeRunId   = 0;
let ledgerEntries = [];

// ── Helpers ────────────────────────────────────────────────
function normDate(v) { return String(v ?? '').slice(0, 10); }

function fmtDate(v) {
  const d = normDate(v);
  if (!d) return '—';
  const [y, m, day] = d.split('-');
  const mo = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
  return `${mo[parseInt(m,10)-1]} ${parseInt(day,10)}, ${y}`;
}

function fmtTime(raw) {
  const s = String(raw ?? '').padStart(4, '0');
  return `${s.slice(0,2)}:${s.slice(2)}`;
}

function parseRoute(route) {
  if (!route) return [];
  return route.split(' -> ').map(s => s.trim()).filter(Boolean);
}

function parseDepartureTimes(dt) {
  const map = {};
  if (!dt) return map;
  dt.split(',').forEach(part => {
    const idx = part.indexOf('@');
    if (idx < 0) return;
    map[part.slice(0, idx).trim()] = part.slice(idx + 1).trim();
  });
  return map;
}

function esc(s) {
  return String(s ?? '')
    .replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

function setStatus(msg) { document.getElementById('status-text').textContent = msg; }

// ── MCP fetch ─────────────────────────────────────────────
async function execQuery(jwe, sql) {
  const url = `${BASE_URL}/${encodeURIComponent(jwe)}/openapi/execute_query?query=${encodeURIComponent(sql)}`;
  const resp = await fetch(url);
  if (!resp.ok) {
    const txt = await resp.text();
    throw new Error(`HTTP ${resp.status}: ${txt.slice(0,250)}`);
  }
  const payload = await resp.json();
  const { columns = [], rows: rawRows, count = 0 } = payload;
  if (!count || !rawRows) return [];
  return rawRows.map(row => {
    if (Array.isArray(row)) {
      const obj = {};
      columns.forEach((col, i) => { obj[col] = row[i]; });
      return obj;
    }
    return row;
  });
}

// ── Query ledger ──────────────────────────────────────────
function ledgerAdd(label, role, sql) {
  const id = ledgerEntries.length;
  ledgerEntries.push({ id, label, role, sql, status: 'Pending', rowCount: '…' });
  renderLedger();
  return id;
}

function ledgerUpdate(id, status, rowCount) {
  if (!ledgerEntries[id]) return;
  ledgerEntries[id].status   = status;
  ledgerEntries[id].rowCount = rowCount;
  renderLedger();
}

function renderLedger() {
  const tbody = document.getElementById('ledger-tbody');
  tbody.innerHTML = '';
  ledgerEntries.forEach(e => {
    const tr = document.createElement('tr');
    tr.className = 'ledger-row';
    const sc = e.status === 'OK' ? 'status-ok' : e.status === 'Pending' ? 'status-pending' : 'status-failed';
    tr.innerHTML = `
      <td><span class="ledger-toggle">&#9658;</span></td>
      <td>
        <div>${esc(e.label)}</div>
        <div class="ledger-sql" id="lsql-${e.id}"><pre>${esc(e.sql)}</pre></div>
      </td>
      <td>${esc(e.role)}</td>
      <td><span class="${sc}">${esc(e.status)}</span></td>
      <td>${esc(String(e.rowCount))}</td>
    `;
    function toggle() {
      const blk = document.getElementById(`lsql-${e.id}`);
      const tog = tr.querySelector('.ledger-toggle');
      const open = blk.style.display === 'block';
      blk.style.display = open ? 'none' : 'block';
      tog.innerHTML = open ? '&#9658;' : '&#9660;';
    }
    tr.addEventListener('click', toggle);
    tbody.appendChild(tr);
  });
}

// ── Map ───────────────────────────────────────────────────
function initMap() {
  if (leafletMap) return;
  leafletMap = L.map('map-container', { zoomControl: true });
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors',
    maxZoom: 18
  }).addTo(leafletMap);
}

function drawRoute(row) {
  if (!leafletMap) initMap();

  // Clear previous route layers
  if (routeLayerGroup) {
    leafletMap.removeLayer(routeLayerGroup);
    routeLayerGroup = null;
  }

  const degradedEl = document.getElementById('map-degraded');
  const airports     = parseRoute(row.Route);
  const depTimes     = parseDepartureTimes(row.DepartureTimes);

  const missing = airports.filter(c => !airportCoords[c]);
  if (missing.length > 0 && airports.filter(c => airportCoords[c]).length < 2) {
    degradedEl.style.display = 'block';
    degradedEl.textContent = `Airport coordinates unavailable (missing: ${missing.join(', ')}). Map cannot be rendered for this itinerary.`;
    return;
  }
  degradedEl.style.display = 'none';

  routeLayerGroup = L.layerGroup().addTo(leafletMap);

  const latlngs = [];
  airports.forEach((code, i) => {
    const c = airportCoords[code];
    if (!c) return;
    const ll = [c.lat, c.lon];
    latlngs.push(ll);

    const color = i === 0 ? '#1f8a70' : (i === airports.length - 1 ? '#c54f36' : '#0e3a52');
    const radius = (i === 0 || i === airports.length - 1) ? 9 : 7;
    const dt = depTimes[code];
    const timeStr = dt ? ` · Dep ${fmtTime(dt)}` : (i === airports.length - 1 ? ' · Arrival' : '');
    L.circleMarker(ll, {
      radius, fillColor: color, color: '#fff', weight: 2, fillOpacity: 0.95
    }).bindPopup(`<strong>${esc(code)}</strong><br>${esc(c.name || code)}${timeStr}<br>Stop ${i+1} of ${airports.length}`)
      .addTo(routeLayerGroup);

    // Airport label
    const icon = L.divIcon({
      className: '',
      html: `<div style="background:${color};color:#fff;font-size:10px;font-weight:700;padding:2px 5px;border-radius:4px;white-space:nowrap;box-shadow:0 1px 3px rgba(0,0,0,.3)">${esc(code)}</div>`,
      iconAnchor: [-12, 10]
    });
    L.marker(ll, { icon, interactive: false }).addTo(routeLayerGroup);
  });

  if (latlngs.length >= 2) {
    L.polyline(latlngs, { color: '#3c88b5', weight: 3, opacity: 0.75 }).addTo(routeLayerGroup);
    try {
      leafletMap.fitBounds(L.latLngBounds(latlngs).pad(0.12));
    } catch (_) {}
  }

  setTimeout(() => { try { leafletMap.invalidateSize(); } catch(_) {} }, 120);
}

// ── Enrichment query ──────────────────────────────────────
async function runEnrichmentQuery(jwe) {
  const allCodes = new Set();
  rows.forEach(r => parseRoute(r.Route).forEach(c => allCodes.add(c)));
  if (!allCodes.size) return;

  const list = [...allCodes].map(c => `'${c}'`).join(', ');
  const sql = `SELECT code, name, latitude AS lat, longitude AS lon\nFROM ontime.airports_latest\nWHERE code IN (${list})`;

  const lid = ledgerAdd('Airport coordinate enrichment', 'Enrichment', sql);
  try {
    const data = await execQuery(jwe, sql);
    data.forEach(r => {
      airportCoords[r.code] = { lat: parseFloat(r.lat), lon: parseFloat(r.lon), name: r.name };
    });
    ledgerUpdate(lid, 'OK', data.length);
  } catch (e) {
    ledgerUpdate(lid, 'Failed', 0);
    const degradedEl = document.getElementById('map-degraded');
    degradedEl.style.display = 'block';
    degradedEl.textContent = `Airport coordinate enrichment failed: ${e.message.replace(/jwe|token/gi,'***')}. Map may be unavailable.`;
  }
}

// ── Route sequence panel ──────────────────────────────────
function renderRoutePanel(row) {
  const airports  = parseRoute(row.Route);
  const depTimes  = parseDepartureTimes(row.DepartureTimes);
  document.getElementById('route-header').textContent =
    `${row.Carrier}${row.FlightNum} · ${row.Tail_Number} · ${fmtDate(row.FlightDate)}`;
  const ul = document.getElementById('stop-list');
  ul.innerHTML = '';
  airports.forEach((code, i) => {
    const li   = document.createElement('li');
    li.className = 'stop-item';
    const numCls = i === 0 ? 'stop-num stop-first' : (i === airports.length - 1 ? 'stop-num stop-last' : 'stop-num');
    const numLbl = i === airports.length - 1 ? '&#10003;' : String(i + 1);
    const info   = airportCoords[code];
    const nameHtml = info ? `<div class="stop-name">${esc(info.name)}</div>` : '';
    const dt = depTimes[code];
    const timeHtml = dt
      ? `<div class="stop-time">Dep ${fmtTime(dt)}</div>`
      : (i === airports.length - 1 ? `<div class="stop-time" style="color:var(--red)">Arrival</div>` : '');
    li.innerHTML = `
      <div class="${numCls}">${numLbl}</div>
      <div class="stop-info">
        <div class="stop-code">${esc(code)}</div>
        ${nameHtml}
        ${timeHtml}
      </div>`;
    ul.appendChild(li);
  });
}

// ── Itinerary table ───────────────────────────────────────
function renderTable() {
  const tbody = document.getElementById('itinerary-tbody');
  tbody.innerHTML = '';
  rows.forEach((row, i) => {
    const tr = document.createElement('tr');
    tr.className = 'data-row' + (i === selectedIdx ? ' active' : '');
    tr.dataset.idx = i;
    tr.innerHTML = `
      <td>${i + 1}</td>
      <td style="white-space:nowrap">${esc(fmtDate(row.FlightDate))}</td>
      <td><span class="carrier-badge">${esc(row.Carrier)}</span></td>
      <td style="font-weight:600">${esc(String(row.FlightNum))}</td>
      <td style="font-family:monospace;font-size:0.82rem">${esc(row.Tail_Number)}</td>
      <td><span class="hops-badge">${row.Hops}</span></td>
      <td class="route-cell" title="${esc(row.Route)}">${esc(row.Route)}</td>
    `;
    tr.addEventListener('click', () => selectItinerary(i));
    tbody.appendChild(tr);
  });
}

function updateTableActive() {
  document.querySelectorAll('#itinerary-tbody .data-row').forEach((tr, i) => {
    tr.classList.toggle('active', i === selectedIdx);
  });
}

// ── Selection ─────────────────────────────────────────────
function selectItinerary(idx) {
  selectedIdx = idx;
  updateTableActive();
  const row = rows[idx];
  document.getElementById('map-badge').textContent =
    `${row.Carrier}${row.FlightNum} · ${normDate(row.FlightDate)}`;
  renderRoutePanel(row);
  try { drawRoute(row); } catch (e) {
    const degradedEl = document.getElementById('map-degraded');
    degradedEl.style.display = 'block';
    degradedEl.textContent = `Map render error: ${e.message}`;
  }
}

// ── KPI strip ─────────────────────────────────────────────
function renderKPIs() {
  if (!rows.length) return;
  const top = rows[0]; // anchored to top-ranked result

  document.getElementById('kpi-hops').textContent  = top.Hops;
  document.getElementById('kpi-tail').textContent  = top.Tail_Number;
  document.getElementById('kpi-tail-sub').textContent = `${top.Carrier} · ${fmtDate(top.FlightDate)}`;
  document.getElementById('kpi-fnum').textContent  = `${top.Carrier}${top.FlightNum}`;
  document.getElementById('kpi-date').textContent  = fmtDate(top.FlightDate);

  // Route repetition: count flight numbers appearing ≥2 times
  const seen = {};
  rows.forEach(r => { const k = `${r.Carrier}${r.FlightNum}`; seen[k] = (seen[k] ?? 0) + 1; });
  const recurring = Object.values(seen).filter(v => v >= 2).length;
  document.getElementById('kpi-repeats').textContent = recurring;
  document.getElementById('kpi-repeats-sub').textContent = 'flight numbers seen 2+ days';
}

// ── Primary query flow ────────────────────────────────────
async function runPrimaryQuery() {
  if (isRunning) return;
  const jwe = document.getElementById('jwe-input').value.trim();
  const sql = document.getElementById('sql-input').value.trim();
  if (!jwe || !sql) { setStatus('Please enter a JWE token and SQL.'); return; }

  isRunning = true;
  activeRunId++;
  const myRun = activeRunId;
  document.getElementById('fetch-btn').disabled = true;
  setStatus('Fetching…');

  const primaryLid = ledgerAdd('Highest daily hops per aircraft per flight number', 'Primary', sql);

  try { localStorage.setItem(STORAGE_KEY, jwe); } catch (_) {}

  try {
    const data = await execQuery(jwe, sql);
    if (myRun !== activeRunId) return;

    ledgerUpdate(primaryLid, 'OK', data.length);

    if (!data.length) {
      setStatus('Query returned no rows.');
      document.getElementById('pre-load-state').innerHTML =
        '<h2 style="font-family:Georgia,serif;color:var(--amber)">No data returned</h2><p>The query executed successfully but returned an empty result set.</p>';
      return;
    }

    rows       = data;
    selectedIdx = 0;

    // Show dashboard
    document.getElementById('pre-load-state').style.display = 'none';
    document.getElementById('dashboard').style.display = 'block';

    renderKPIs();
    renderTable();

    // Init map (container now visible)
    initMap();

    setStatus('Loading airport coordinates…');
    await runEnrichmentQuery(jwe);
    if (myRun !== activeRunId) return;

    // Select row 0 — redraws map with (now enriched) coords
    selectItinerary(0);
    setStatus(`Loaded ${data.length} itineraries.`);

  } catch (e) {
    if (myRun !== activeRunId) return;
    ledgerUpdate(primaryLid, 'Failed', 0);
    const safe = e.message.replace(/jwe|token/gi, '***');
    setStatus(`Error: ${safe}`);
    document.getElementById('pre-load-state').style.display = 'block';
    document.getElementById('pre-load-state').innerHTML =
      `<h2 style="font-family:Georgia,serif;color:var(--red)">Query failed</h2><p style="color:var(--muted)">${esc(safe)}</p>`;
  } finally {
    if (myRun === activeRunId) {
      isRunning = false;
      document.getElementById('fetch-btn').disabled = false;
    }
  }
}

function forgetToken() {
  try { localStorage.removeItem(STORAGE_KEY); } catch (_) {}
  document.getElementById('jwe-input').value = '';
  setStatus('Token forgotten.');
}

// ── Export CSV ────────────────────────────────────────────
function exportCsv() {
  if (!rows.length) return;
  const cols = ['FlightDate','Tail_Number','Carrier','FlightNum','Hops','Route','DepartureTimes'];
  const lines = [cols.join(',')];
  rows.forEach(r => {
    lines.push(cols.map(c => `"${String(r[c] ?? '').replace(/"/g,'""')}"`).join(','));
  });
  const a = document.createElement('a');
  a.href = URL.createObjectURL(new Blob([lines.join('\n')], { type: 'text/csv' }));
  a.download = 'hops_per_day_top10.csv';
  a.click();
}

// ── Init ──────────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', () => {
  document.getElementById('sql-input').value = PRIMARY_SQL;

  let storedJwe = '';
  try { storedJwe = localStorage.getItem(STORAGE_KEY) ?? ''; } catch (_) {}

  if (storedJwe) {
    document.getElementById('jwe-input').value = storedJwe;
    runPrimaryQuery();
  } else {
    document.getElementById('pre-load-state').style.display = 'block';
  }
});
</script>
</body>
</html>
```
