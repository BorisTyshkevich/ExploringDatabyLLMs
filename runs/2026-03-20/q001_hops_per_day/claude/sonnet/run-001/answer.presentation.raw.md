Here is the final browser-ready artifact:

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Highest Daily Hops — One Aircraft, One Flight Number</title>
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"/>
<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
<style>
:root {
  --bg-top:#eaf3f8;--bg-bottom:#f6fafc;
  --panel:#ffffff;--panel-alt:#f3f7fa;
  --ink:#163244;--muted:#5d7485;
  --navy:#0e3a52;--slate:#5c7080;
  --sky:#3c88b5;--teal:#1f8a70;
  --amber:#d48a1f;--red:#c54f36;
  --border:rgba(22,50,68,.10);--grid:rgba(22,50,68,.12);
  --shadow:0 18px 45px rgba(14,58,82,.10);
  --radius-xl:22px;--radius-lg:16px;--radius-md:12px;--radius-sm:8px;
}
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
body{font-family:"Segoe UI","Helvetica Neue",sans-serif;color:var(--ink);
  background:linear-gradient(175deg,var(--bg-top),var(--bg-bottom));min-height:100vh}
h1,h2,h3{font-family:Georgia,ui-serif,serif}

.header{background:linear-gradient(135deg,var(--navy),#1a5c7c);color:#fff;
  padding:2rem 0;text-align:center}
.header h1{font-size:1.6rem;font-weight:700;margin-bottom:.3rem}
.header p{color:rgba(255,255,255,.75);font-size:.92rem}

.container{max-width:1280px;margin:0 auto;padding:1.25rem 1rem}

/* KPI strip */
.kpi-strip{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));
  gap:.75rem;margin-bottom:1.25rem}
.kpi-card{background:var(--panel);border-radius:var(--radius-md);padding:1rem 1.1rem;
  box-shadow:var(--shadow);border-left:4px solid var(--sky)}
.kpi-card .label{font-size:.7rem;text-transform:uppercase;letter-spacing:.06em;color:var(--muted)}
.kpi-card .value{font-size:1.3rem;font-weight:700;color:var(--navy);margin-top:.2rem;word-break:break-word}
.kpi-card.primary{border-left-color:var(--red)}
.kpi-card.accent{border-left-color:var(--teal)}

/* Cards */
.card{background:var(--panel);border-radius:var(--radius-lg);box-shadow:var(--shadow);
  padding:1.25rem;margin-bottom:1.25rem}
.card h2{font-size:1.05rem;color:var(--navy);margin-bottom:.75rem;
  display:flex;align-items:center;gap:.4rem}

/* Map */
#map-container{height:420px;border-radius:var(--radius-md);overflow:hidden;
  background:var(--panel-alt);position:relative}
#map{width:100%;height:100%}
#map-degraded{display:none;position:absolute;inset:0;align-items:center;justify-content:center;
  background:var(--panel-alt);color:var(--muted);font-size:.9rem;text-align:center;padding:2rem;z-index:999}

/* Legend */
.legend{display:flex;gap:1.2rem;flex-wrap:wrap;margin-top:.6rem;font-size:.78rem;color:var(--muted);
  align-items:center}
.leg-dot{display:inline-block;width:10px;height:10px;border-radius:50%;
  margin-right:.3rem;vertical-align:middle}
.leg-line{display:inline-block;width:18px;height:3px;border-radius:2px;
  margin-right:.3rem;vertical-align:middle}

/* Route detail */
.route-detail{display:flex;flex-wrap:wrap;gap:.4rem;align-items:center;margin-top:.4rem}
.route-stop{background:var(--panel-alt);border:1px solid var(--border);border-radius:var(--radius-sm);
  padding:.4rem .65rem;font-size:.8rem;text-align:center;min-width:58px}
.route-stop .code{font-weight:700;color:var(--navy);font-size:.9rem}
.route-stop .dep{color:var(--muted);font-size:.68rem;margin-top:.1rem}
.route-stop.terminus{border-color:var(--teal);background:#f0faf6}
.route-arrow{color:var(--sky);font-weight:700;font-size:1.1rem;padding:0 .1rem}

/* Table */
.scroll-wrap{overflow-x:auto}
.itinerary-table{width:100%;border-collapse:collapse;font-size:.8rem}
.itinerary-table thead{background:var(--navy);color:#fff;position:sticky;top:0;z-index:1}
.itinerary-table th{padding:.5rem .6rem;text-align:left;font-weight:600;font-size:.7rem;
  text-transform:uppercase;letter-spacing:.04em;white-space:nowrap}
.itinerary-table td{padding:.5rem .6rem;border-bottom:1px solid var(--border);vertical-align:top}
.itinerary-table tbody tr{cursor:pointer;transition:background .15s}
.itinerary-table tbody tr:hover{background:#f0f8ff}
.itinerary-table tbody tr.active{background:#dbeef8;outline:2px solid var(--sky);outline-offset:-2px}
.itinerary-table .hops-badge{background:var(--red);color:#fff;border-radius:999px;
  padding:.15rem .55rem;font-weight:700;font-size:.78rem;display:inline-block}
.route-cell{max-width:320px;word-break:break-word;font-size:.76rem;color:var(--slate)}
.dep-cell{max-width:260px;font-size:.72rem;color:var(--muted)}

/* Ledger */
.ledger-table{width:100%;border-collapse:collapse;font-size:.78rem}
.ledger-table th{text-align:left;padding:.4rem .5rem;border-bottom:2px solid var(--border);
  font-size:.68rem;text-transform:uppercase;color:var(--muted);letter-spacing:.04em}
.ledger-table td{padding:.4rem .5rem;border-bottom:1px solid var(--border);vertical-align:top}
.ledger-table .ledger-row{cursor:pointer}
.ledger-table .ledger-row:hover{background:var(--panel-alt)}
.ledger-toggle{font-family:monospace;width:1.2em;display:inline-block;font-size:.8rem}
.ledger-sql-row td{padding:0}
.ledger-sql-row .sql-block{display:none;background:var(--panel-alt);border-radius:var(--radius-sm);
  padding:.6rem .75rem;margin:.2rem .5rem .4rem;font-family:monospace;font-size:.72rem;
  white-space:pre-wrap;word-break:break-all}
.status-ok{color:var(--teal);font-weight:600}
.status-pending{color:var(--amber);font-weight:600}
.status-failed{color:var(--red);font-weight:600}

/* Empty state */
.empty-state{text-align:center;padding:3rem 1rem;color:var(--muted)}
.empty-state h2{color:var(--navy);margin-bottom:.5rem;font-size:1.2rem}

/* Export bar */
.export-bar{display:flex;justify-content:flex-end;margin-bottom:.5rem}

/* Footer controls */
.footer-controls{background:var(--panel);border-radius:var(--radius-lg);box-shadow:var(--shadow);
  padding:1.25rem;margin-top:1.25rem}
.footer-controls h3{font-size:.95rem;color:var(--navy);margin-bottom:.75rem;font-family:Georgia,serif}
.ctrl-row{display:flex;gap:.5rem;margin-bottom:.5rem;align-items:flex-start;flex-wrap:wrap}
.ctrl-row label{font-size:.75rem;color:var(--muted);min-width:65px;padding-top:.45rem}
.ctrl-row input[type=password]{flex:1;min-width:220px;padding:.4rem .6rem;border:1px solid var(--border);
  border-radius:var(--radius-sm);font-size:.8rem}
.ctrl-row textarea{flex:1;min-width:220px;height:130px;padding:.5rem .6rem;border:1px solid var(--border);
  border-radius:var(--radius-sm);font-family:monospace;font-size:.74rem;resize:vertical}
.btn{padding:.42rem .95rem;border:none;border-radius:var(--radius-sm);font-size:.8rem;
  font-weight:600;cursor:pointer;transition:opacity .15s}
.btn-primary{background:var(--navy);color:#fff}
.btn-primary:hover{opacity:.85}
.btn-primary:disabled{opacity:.5;cursor:not-allowed}
.btn-secondary{background:var(--panel-alt);color:var(--navy);border:1px solid var(--border)}
.btn-secondary:hover{opacity:.8}
.status-text{font-size:.78rem;color:var(--muted);margin-top:.35rem}

/* Responsive */
@media(max-width:768px){
  .kpi-strip{grid-template-columns:1fr 1fr}
  #map-container{height:300px}
}

/* Leaflet tooltip override */
.airport-label{background:transparent;border:none;box-shadow:none;
  font-weight:700;font-size:.72rem;color:var(--navy);white-space:nowrap}
.airport-label::before{display:none}
</style>
</head>
<body>

<div class="header">
  <h1>Highest Daily Hops for One Aircraft on One Flight Number</h1>
  <p>Multi-stop itineraries — Southwest Airlines point-to-point network &mdash; dynamic analysis</p>
</div>

<div class="container">
  <!-- Empty state -->
  <div id="empty-state" class="empty-state">
    <h2>Awaiting data</h2>
    <p>Enter your JWE token below and click <strong>Fetch</strong> to load the dashboard.</p>
  </div>

  <!-- Dashboard -->
  <div id="dashboard" style="display:none">

    <!-- KPI strip — anchored to top-ranked row, never changes on selection -->
    <div class="kpi-strip" id="kpi-strip">
      <div class="kpi-card primary">
        <div class="label">Max Hops</div>
        <div class="value" id="kpi-hops">—</div>
      </div>
      <div class="kpi-card">
        <div class="label">Tail Number</div>
        <div class="value" id="kpi-tail">—</div>
      </div>
      <div class="kpi-card">
        <div class="label">Flight</div>
        <div class="value" id="kpi-flight">—</div>
      </div>
      <div class="kpi-card">
        <div class="label">Date (Top Record)</div>
        <div class="value" id="kpi-date">—</div>
      </div>
      <div class="kpi-card accent">
        <div class="label">Route Repeats in Top 10</div>
        <div class="value" id="kpi-repetition">—</div>
      </div>
    </div>

    <!-- Map card -->
    <div class="card">
      <h2 id="map-title">&#x1F5FA; Selected Itinerary Route Map</h2>
      <div id="map-container">
        <div id="map"></div>
        <div id="map-degraded">Awaiting airport coordinate enrichment&hellip;</div>
      </div>
      <div class="legend">
        <span><span class="leg-dot" style="background:var(--red)"></span>Origin</span>
        <span><span class="leg-dot" style="background:var(--navy)"></span>Stop</span>
        <span><span class="leg-dot" style="background:var(--teal)"></span>Final destination</span>
        <span><span class="leg-line" style="background:var(--sky)"></span>Route arc</span>
      </div>
    </div>

    <!-- Route detail panel -->
    <div class="card">
      <h2>Route Sequence &amp; Detail</h2>
      <div id="route-detail" class="route-detail"></div>
    </div>

    <!-- Itinerary table -->
    <div class="card">
      <h2>Top 10 Itineraries by Hop Count</h2>
      <div class="export-bar">
        <button class="btn btn-secondary" id="btn-export">Export CSV</button>
      </div>
      <div class="scroll-wrap">
        <table class="itinerary-table">
          <thead>
            <tr>
              <th>#</th>
              <th>Tail</th>
              <th>Flight</th>
              <th>Carrier</th>
              <th>Date</th>
              <th>Hops</th>
              <th>Route</th>
              <th>Departure Schedule</th>
            </tr>
          </thead>
          <tbody id="table-body"></tbody>
        </table>
      </div>
    </div>

    <!-- Query ledger -->
    <div class="card">
      <h2>Query Ledger</h2>
      <table class="ledger-table">
        <thead>
          <tr><th></th><th>Label</th><th>Role</th><th>Status</th><th>Rows</th></tr>
        </thead>
        <tbody id="ledger-body"></tbody>
      </table>
    </div>

  </div><!-- /dashboard -->
</div><!-- /container -->

<!-- Footer control block -->
<footer>
  <div class="container">
    <div class="footer-controls">
      <h3>Data Connection</h3>
      <div class="ctrl-row">
        <label>JWE Token</label>
        <input type="password" id="jwe-input" placeholder="Paste JWE token here">
        <button class="btn btn-secondary" id="btn-forget">Forget</button>
      </div>
      <div class="ctrl-row">
        <label>SQL</label>
        <textarea id="sql-input"></textarea>
      </div>
      <div class="ctrl-row">
        <button class="btn btn-primary" id="btn-fetch">Fetch</button>
      </div>
      <div class="status-text" id="status-text"></div>
    </div>
  </div>
</footer>

<script>
(function () {
  "use strict";

  /* ─── Constants ─── */
  const STORAGE_KEY = "OnTimeAnalystDashboard::auth::jwe";
  const BASE_URL    = "https://mcp.demo.altinity.cloud";

  const SAVED_SQL = `WITH leg_data AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline AS Carrier,
        FlightDate,
        arraySort(
            x -> x.1,
            groupArray((
                coalesce(DepTime, CRSDepTime, 0),
                Origin,
                Dest
            ))
        ) AS sorted_legs
    FROM ontime.ontime
    WHERE Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
      AND Cancelled = 0
    GROUP BY Tail_Number, Flight_Number_Reporting_Airline, Carrier, FlightDate
    HAVING length(sorted_legs) >= 2
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    Carrier,
    FlightDate,
    length(sorted_legs) AS Hops,
    arrayStringConcat(
        arrayPushBack(
            arrayMap(x -> x.2, sorted_legs),
            arrayElement(sorted_legs, -1).3
        ),
        ' \u2192 '
    ) AS Route,
    arrayStringConcat(
        arrayMap(x -> concat(leftPad(toString(x.1), 4, '0'), ' (', x.2, ')'), sorted_legs),
        ', '
    ) AS DepartureTimes
FROM leg_data
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10`;

  /* ─── State ─── */
  let rows           = [];
  let selectedIdx    = 0;
  let airportCoords  = {};   // code → {lat, lon, name}
  let enrichmentDone = false;
  let leafletMap     = null;
  let routeLayer     = null;
  let isRunning      = false;
  const ledger       = [];

  /* ─── DOM refs ─── */
  const $jwe      = document.getElementById("jwe-input");
  const $sql      = document.getElementById("sql-input");
  const $fetch    = document.getElementById("btn-fetch");
  const $forget   = document.getElementById("btn-forget");
  const $status   = document.getElementById("status-text");
  const $empty    = document.getElementById("empty-state");
  const $dash     = document.getElementById("dashboard");
  const $tbody    = document.getElementById("table-body");
  const $detail   = document.getElementById("route-detail");
  const $mapEl    = document.getElementById("map");
  const $mapDeg   = document.getElementById("map-degraded");
  const $mapTitle = document.getElementById("map-title");
  const $export   = document.getElementById("btn-export");

  /* ─── Init ─── */
  $sql.value = SAVED_SQL;
  const savedToken = localStorage.getItem(STORAGE_KEY);
  if (savedToken) $jwe.value = savedToken;

  $fetch.addEventListener("click", run);
  $forget.addEventListener("click", () => {
    localStorage.removeItem(STORAGE_KEY);
    $jwe.value = "";
    setStatus("Token cleared.");
  });
  $export.addEventListener("click", exportCSV);

  if (savedToken) run();

  /* ─── Utilities ─── */
  function setStatus(msg) { $status.textContent = msg; }

  function esc(s) {
    const d = document.createElement("div");
    d.textContent = String(s ?? "");
    return d.innerHTML;
  }

  function normDate(v) { return String(v ?? "").slice(0, 10); }

  /** Parse "0543 (ISP), 0810 (BWI)" → {ISP:"05:43", BWI:"08:10"} */
  function parseDepartureTimes(dtStr) {
    const map = {};
    const re  = /(\d{4})\s*\((\w+)\)/g;
    let m;
    while ((m = re.exec(dtStr ?? "")) !== null) {
      const hhmm = m[1].padStart(4, "0");
      map[m[2]] = hhmm.slice(0, 2) + ":" + hhmm.slice(2);
    }
    return map;
  }

  /** All unique airport codes across all rows */
  function allAirportCodes(dataRows) {
    const codes = new Set();
    dataRows.forEach(r => {
      (r.Route ?? "").split(/\s*\u2192\s*/).forEach(c => { if (c.trim()) codes.add(c.trim()); });
    });
    return [...codes];
  }

  /** How many top-10 rows share the same flight number AND route as the top-ranked row */
  function computeRepetition(topRow, allRows) {
    if (!topRow) return "—";
    const fn    = topRow.Flight_Number_Reporting_Airline;
    const route = topRow.Route;
    const count = allRows.filter(r =>
      r.Flight_Number_Reporting_Airline === fn && r.Route === route
    ).length;
    return count <= 1 ? "Unique in top 10" : `${count}× (WN\u00a0${fn})`;
  }

  /* ─── Ledger ─── */
  function addLedger(label, role, sql) {
    const entry = { label, role, sql: sql || "", status: "Pending", rows: "—" };
    ledger.push(entry);
    renderLedger();
    return entry;
  }

  function updateLedger(entry, status, rowCount) {
    entry.status = status;
    if (rowCount != null) entry.rows = rowCount;
    renderLedger();
  }

  function renderLedger() {
    const $lb = document.getElementById("ledger-body");
    $lb.innerHTML = "";
    ledger.forEach((e, i) => {
      const cls = e.status === "OK" ? "status-ok"
        : e.status === "Pending"    ? "status-pending" : "status-failed";
      const tr = document.createElement("tr");
      tr.className = "ledger-row";
      tr.innerHTML = `
        <td><span class="ledger-toggle" data-i="${i}">\u25b6</span></td>
        <td>${esc(e.label)}</td>
        <td>${esc(e.role)}</td>
        <td class="${cls}">${esc(e.status)}</td>
        <td>${esc(String(e.rows))}</td>`;
      tr.addEventListener("click", () => toggleSQL(i));
      $lb.appendChild(tr);

      const sqlTr = document.createElement("tr");
      sqlTr.className = "ledger-sql-row";
      sqlTr.id = "lsql-" + i;
      sqlTr.innerHTML = `<td colspan="5"><div class="sql-block">${esc(e.sql)}</div></td>`;
      $lb.appendChild(sqlTr);
    });
  }

  function toggleSQL(i) {
    const row = document.getElementById("lsql-" + i);
    if (!row) return;
    const blk = row.querySelector(".sql-block");
    const tog = document.querySelector(`.ledger-toggle[data-i="${i}"]`);
    const open = blk.style.display !== "block";
    blk.style.display = open ? "block" : "none";
    if (tog) tog.textContent = open ? "\u25bc" : "\u25b6";
  }

  /* ─── Query execution ─── */
  async function executeQuery(sql) {
    const token = $jwe.value.trim();
    const url   = `${BASE_URL}/${encodeURIComponent(token)}/openapi/execute_query?query=${encodeURIComponent(sql)}`;
    const resp  = await fetch(url);
    if (!resp.ok) {
      const txt = await resp.text().catch(() => "");
      throw new Error(`HTTP ${resp.status}: ${txt.slice(0, 300)}`);
    }
    const payload = await resp.json();
    const cols    = payload.columns ?? [];
    const raw     = payload.rows    ?? [];
    return raw.map(r => {
      if (!Array.isArray(r)) return r;
      const obj = {};
      cols.forEach((c, i) => { obj[c] = r[i]; });
      return obj;
    });
  }

  /* ─── Main fetch flow ─── */
  async function run() {
    if (isRunning) return;
    isRunning = true;
    $fetch.disabled = true;

    const token = $jwe.value.trim();
    const sql   = $sql.value.trim();
    if (!token || !sql) {
      setStatus("JWE token and SQL are required.");
      isRunning = false; $fetch.disabled = false; return;
    }
    localStorage.setItem(STORAGE_KEY, token);
    ledger.length  = 0;
    airportCoords  = {};
    enrichmentDone = false;

    setStatus("Running primary query\u2026");
    const primaryEntry = addLedger(
      "Primary: Highest hop itineraries (top 10)",
      "primary",
      sql
    );

    try {
      rows = await executeQuery(sql);
      updateLedger(primaryEntry, "OK", rows.length);
    } catch (err) {
      updateLedger(primaryEntry, "Failed", 0);
      setStatus("Primary query failed: " + err.message);
      isRunning = false; $fetch.disabled = false; return;
    }

    if (!rows.length) {
      setStatus("Query returned no rows.");
      $empty.innerHTML = "<h2>No data</h2><p>Zero rows returned.</p>";
      isRunning = false; $fetch.disabled = false; return;
    }

    rows = rows.map(r => ({ ...r, _date: normDate(r.FlightDate) }));

    $empty.style.display = "none";
    $dash.style.display  = "";

    renderKPIs(rows);
    renderTable();
    initMap();
    selectRow(0, /* skipMapRender */ true);
    showMapDegraded("Airport coordinate enrichment in progress\u2026");

    setStatus("Running airport coordinate enrichment\u2026");
    await runEnrichment();

    renderMapRoute(rows[selectedIdx]);

    setStatus("Dashboard ready.");
    isRunning = false;
    $fetch.disabled = false;
  }

  /* ─── KPIs — anchored to top row only ─── */
  function renderKPIs(dataRows) {
    const top = dataRows[0];
    document.getElementById("kpi-hops").textContent       = top?.Hops ?? "—";
    document.getElementById("kpi-tail").textContent       = top?.Tail_Number ?? "—";
    document.getElementById("kpi-flight").textContent     =
      (top?.Carrier ? top.Carrier + "\u00a0" : "") + (top?.Flight_Number_Reporting_Airline ?? "—");
    document.getElementById("kpi-date").textContent       = top?._date ?? "—";
    document.getElementById("kpi-repetition").textContent = computeRepetition(top, dataRows);
  }

  /* ─── Table ─── */
  function renderTable() {
    $tbody.innerHTML = "";
    rows.forEach((r, i) => {
      const tr = document.createElement("tr");
      tr.dataset.idx = i;
      tr.innerHTML = `
        <td>${i + 1}</td>
        <td>${esc(r.Tail_Number)}</td>
        <td>${esc(r.Flight_Number_Reporting_Airline)}</td>
        <td>${esc(r.Carrier)}</td>
        <td>${esc(r._date)}</td>
        <td><span class="hops-badge">${esc(r.Hops)}</span></td>
        <td class="route-cell">${esc(r.Route)}</td>
        <td class="dep-cell">${esc(r.DepartureTimes)}</td>`;
      tr.addEventListener("click", () => selectRow(i));
      $tbody.appendChild(tr);
    });
  }

  /* ─── Row selection ─── */
  function selectRow(idx, skipMapRender) {
    selectedIdx = idx;
    $tbody.querySelectorAll("tr").forEach((tr, i) => {
      tr.classList.toggle("active", i === idx);
    });
    const row = rows[idx];
    $mapTitle.textContent =
      `Route Map — ${row?.Carrier ?? ""}\u00a0${row?.Flight_Number_Reporting_Airline ?? ""} · ${row?._date ?? ""}`;
    renderRouteDetail(row);
    if (!skipMapRender) renderMapRoute(row);
  }

  /* ─── Route detail panel ─── */
  function renderRouteDetail(row) {
    $detail.innerHTML = "";
    if (!row) return;
    const stops   = (row.Route ?? "").split(/\s*\u2192\s*/).filter(Boolean);
    const timeMap = parseDepartureTimes(row.DepartureTimes);
    stops.forEach((code, i) => {
      if (i > 0) {
        const arrow = document.createElement("div");
        arrow.className = "route-arrow";
        arrow.textContent = "\u2192";
        $detail.appendChild(arrow);
      }
      const el = document.createElement("div");
      const isEndpoint = i === 0 || i === stops.length - 1;
      el.className = "route-stop" + (isEndpoint ? " terminus" : "");
      const depHtml = timeMap[code]
        ? `<div class="dep">${esc(timeMap[code])}</div>`
        : i === stops.length - 1 ? `<div class="dep">ARR</div>` : "";
      el.innerHTML = `<div class="code">${esc(code)}</div>${depHtml}`;
      $detail.appendChild(el);
    });
  }

  /* ─── Map ─── */
  function initMap() {
    if (leafletMap) { try { leafletMap.remove(); } catch (_) {} leafletMap = null; }
    try {
      leafletMap = L.map($mapEl, { scrollWheelZoom: true }).setView([38, -96], 4);
      L.tileLayer("https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png", {
        attribution: "&copy; OpenStreetMap &copy; CARTO", maxZoom: 18
      }).addTo(leafletMap);
      routeLayer = L.layerGroup().addTo(leafletMap);
      setTimeout(() => leafletMap?.invalidateSize(), 200);
    } catch (e) {
      showMapDegraded("Map initialization failed: " + e.message);
    }
  }

  function showMapDegraded(msg) {
    $mapDeg.textContent = msg;
    $mapDeg.style.display = "flex";
  }

  function hideMapDegraded() {
    $mapDeg.style.display = "none";
  }

  function renderMapRoute(row) {
    if (!leafletMap || !routeLayer) return;
    routeLayer.clearLayers();
    if (!row) return;

    const stops  = (row.Route ?? "").split(/\s*\u2192\s*/).filter(Boolean);
    const coords = stops.map(c => airportCoords[c.trim()]).filter(Boolean);

    if (coords.length < 2) {
      if (enrichmentDone) {
        const missing = stops.filter(c => !airportCoords[c.trim()]);
        showMapDegraded(
          `Insufficient coordinates for this itinerary` +
          (missing.length ? ` (missing: ${missing.join(", ")})` : "") + "."
        );
      }
      return;
    }

    hideMapDegraded();
    const timeMap = parseDepartureTimes(row.DepartureTimes);
    const latlngs = coords.map(c => [c.lat, c.lon]);

    L.polyline(latlngs, { color: "#3c88b5", weight: 2.5, opacity: 0.75 }).addTo(routeLayer);

    stops.forEach((code, i) => {
      const c = airportCoords[code.trim()];
      if (!c) return;
      const isOrigin = i === 0;
      const isDest   = i === stops.length - 1;
      const color    = isOrigin ? "#c54f36" : isDest ? "#1f8a70" : "#0e3a52";
      const marker   = L.circleMarker([c.lat, c.lon], {
        radius: (isOrigin || isDest) ? 8 : 6,
        fillColor: color, color: "#fff", weight: 2, fillOpacity: 0.9
      }).addTo(routeLayer);
      const depStr = timeMap[code]
        ? `Dep: <strong>${timeMap[code]}</strong>`
        : isDest ? "(Final destination)" : "";
      marker.bindPopup(
        `<strong>${esc(code)}</strong>${c.name ? " — " + esc(c.name) : ""}<br>${depStr}`
      );
      marker.bindTooltip(code, {
        permanent: true, direction: "top", offset: [0, -9], className: "airport-label"
      });
    });

    try { leafletMap.fitBounds(L.latLngBounds(latlngs).pad(0.15)); } catch (_) {}
  }

  /* ─── Airport enrichment (once, reused for all rows) ─── */
  async function runEnrichment() {
    const codes = allAirportCodes(rows);
    if (!codes.length) { enrichmentDone = true; return; }
    const codeList  = codes.map(c => `'${c.replace(/'/g, "''")}'`).join(",");
    const enrichSQL = `SELECT code, name, latitude, longitude FROM ontime.airports_latest WHERE code IN (${codeList})`;
    const entry     = addLedger("Airport coordinate enrichment", "enrichment", enrichSQL);

    try {
      const eRows = await executeQuery(enrichSQL);
      updateLedger(entry, "OK", eRows.length);
      eRows.forEach(r => {
        const lat = parseFloat(r.latitude);
        const lon = parseFloat(r.longitude);
        if (!isNaN(lat) && !isNaN(lon)) {
          airportCoords[r.code] = { lat, lon, name: r.name ?? "" };
        }
      });
    } catch (err) {
      updateLedger(entry, "Failed", 0);
      showMapDegraded("Airport coordinate enrichment failed: " + err.message);
    }
    enrichmentDone = true;
  }

  /* ─── CSV export ─── */
  function exportCSV() {
    if (!rows.length) return;
    const headers = ["Tail_Number","Flight_Number_Reporting_Airline","Carrier",
      "FlightDate","Hops","Route","DepartureTimes"];
    const lines = [headers.join(",")];
    rows.forEach(r => {
      lines.push(headers.map(h => {
        const v = h === "FlightDate" ? r._date : (r[h] ?? "");
        return '"' + String(v).replace(/"/g, '""') + '"';
      }).join(","));
    });
    const blob = new Blob([lines.join("\n")], { type: "text/csv" });
    const a    = document.createElement("a");
    a.href     = URL.createObjectURL(blob);
    a.download = "hops_per_day_top10.csv";
    a.click();
  }

})();
</script>
</body>
</html>
```
