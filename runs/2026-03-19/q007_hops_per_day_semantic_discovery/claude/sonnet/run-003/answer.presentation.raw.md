```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>Highest Daily Hops — OnTime Dashboard</title>
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
  --grid: rgba(22,50,68,0.12);
  --border: rgba(22,50,68,0.10);
  --shadow: 0 18px 45px rgba(14,58,82,0.10);
  --radius-xl: 22px;
  --radius-lg: 16px;
  --radius-md: 12px;
  --radius-sm: 8px;
}
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
body {
  font-family: "Segoe UI", "Helvetica Neue", Arial, sans-serif;
  background: linear-gradient(160deg, var(--bg-top) 0%, var(--bg-bottom) 100%);
  color: var(--ink);
  min-height: 100vh;
}
a { color: var(--sky); }

/* ── Header ── */
header {
  background: linear-gradient(135deg, var(--navy) 0%, #1a5272 60%, #245f80 100%);
  color: #fff;
  padding: 28px 40px 24px;
}
.header-inner { max-width: 1280px; margin: 0 auto; }
.header-badge {
  display: inline-block;
  font-size: 11px;
  letter-spacing: .08em;
  text-transform: uppercase;
  background: rgba(255,255,255,.15);
  border-radius: 4px;
  padding: 3px 8px;
  margin-bottom: 10px;
  font-family: "Segoe UI", sans-serif;
}
header h1 {
  font-family: Georgia, ui-serif, serif;
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
  margin-bottom: 6px;
}
header p.subtitle {
  font-size: 14px;
  opacity: .78;
  max-width: 700px;
}

/* ── Layout ── */
main { max-width: 1280px; margin: 0 auto; padding: 28px 24px 48px; }
section { margin-bottom: 24px; }

/* ── KPI Strip ── */
.kpi-strip {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
}
.kpi-card {
  background: var(--panel);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow);
  padding: 18px 20px;
  border-top: 3px solid var(--sky);
}
.kpi-card.highlight { border-top-color: var(--red); }
.kpi-label {
  font-size: 11px;
  letter-spacing: .06em;
  text-transform: uppercase;
  color: var(--muted);
  margin-bottom: 6px;
}
.kpi-value {
  font-size: 26px;
  font-weight: 700;
  color: var(--navy);
  font-family: Georgia, ui-serif, serif;
  line-height: 1;
}
.kpi-sub {
  font-size: 12px;
  color: var(--muted);
  margin-top: 4px;
}

/* ── Two-column map + detail ── */
.map-detail-row {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 16px;
  align-items: start;
}
@media (max-width: 900px) {
  .map-detail-row { grid-template-columns: 1fr; }
}

/* ── Panel card ── */
.panel {
  background: var(--panel);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow);
  overflow: hidden;
}
.panel-header {
  background: var(--panel-alt);
  border-bottom: 1px solid var(--border);
  padding: 12px 18px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: .05em;
  text-transform: uppercase;
  color: var(--slate);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.panel-body { padding: 16px 18px; }

/* ── Map ── */
#map {
  width: 100%;
  height: 420px;
  border-radius: 0 0 var(--radius-lg) var(--radius-lg);
}
.map-degraded {
  display: none;
  height: 420px;
  align-items: center;
  justify-content: center;
  color: var(--muted);
  font-size: 14px;
  text-align: center;
  padding: 24px;
}

/* ── Route Detail Panel ── */
.route-detail-panel { height: 100%; min-height: 460px; }
.stop-list { list-style: none; padding: 4px 0; }
.stop-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid var(--border);
}
.stop-item:last-child { border-bottom: none; }
.stop-num {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--sky);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 1px;
}
.stop-num.final { background: var(--red); }
.stop-code { font-weight: 700; font-size: 15px; color: var(--navy); }
.stop-name { font-size: 12px; color: var(--muted); margin-top: 2px; }
.stop-time {
  margin-left: auto;
  font-size: 12px;
  color: var(--slate);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
  padding-top: 3px;
}
.route-meta {
  font-size: 12px;
  color: var(--muted);
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--border);
}

/* ── Itinerary Table ── */
.table-wrap { overflow-x: auto; }
table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}
thead th {
  background: var(--panel-alt);
  border-bottom: 2px solid var(--grid);
  padding: 10px 12px;
  text-align: left;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: .05em;
  text-transform: uppercase;
  color: var(--slate);
  position: sticky;
  top: 0;
  white-space: nowrap;
}
tbody tr {
  cursor: pointer;
  transition: background .15s;
  border-bottom: 1px solid var(--border);
}
tbody tr:hover { background: var(--panel-alt); }
tbody tr.active {
  background: #ddeef8;
  border-left: 3px solid var(--sky);
}
tbody tr.active td:first-child { padding-left: 9px; }
tbody td {
  padding: 10px 12px;
  vertical-align: top;
}
.badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
  background: var(--sky);
  color: #fff;
}
.badge.hops8 { background: var(--red); }
.route-cell {
  font-size: 11.5px;
  color: var(--muted);
  font-variant-numeric: tabular-nums;
  max-width: 400px;
  word-break: break-all;
}
.route-cell strong { color: var(--ink); font-size: 12px; }

/* ── Legend ── */
.legend {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 12px;
  color: var(--muted);
  align-items: center;
}
.legend-item { display: flex; align-items: center; gap: 6px; }
.legend-dot {
  width: 12px; height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}
.legend-line {
  width: 24px; height: 3px;
  border-radius: 2px;
  flex-shrink: 0;
}

/* ── Footer / SQL ── */
footer {
  background: var(--navy);
  color: rgba(255,255,255,.7);
  padding: 24px 40px;
  font-size: 12px;
}
.footer-inner { max-width: 1280px; margin: 0 auto; }
.sql-toggle {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255,255,255,.55);
  font-size: 12px;
  margin-bottom: 8px;
  user-select: none;
}
.sql-toggle:hover { color: rgba(255,255,255,.85); }
.sql-block {
  display: none;
  background: rgba(0,0,0,.25);
  border-radius: var(--radius-sm);
  padding: 12px 16px;
  margin-top: 8px;
}
.sql-block pre {
  font-family: "Cascadia Code", "Fira Code", "Menlo", monospace;
  font-size: 11.5px;
  white-space: pre-wrap;
  word-break: break-word;
  color: #b8d4e8;
  line-height: 1.6;
}
.footer-meta {
  margin-top: 14px;
  opacity: .5;
  font-size: 11px;
}
</style>
</head>
<body>

<!-- ══ Embedded Data ══════════════════════════════════════════════════ -->
<script type="application/json" id="result-data">
{
  "columns": ["Tail_Number","Flight_Number_Reporting_Airline","IATA_CODE_Reporting_Airline","FlightDate","Hops","Route"],
  "rows": [
    {"Tail_Number":"N957WN","Flight_Number_Reporting_Airline":"366","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2024-12-01T00:00:00Z","Hops":8,"Route":"ISP(543) -> BWI(810) -> MYR(1020) -> BNA(1142) -> VPS(1401) -> DAL(1643) -> LAS(1828) -> OAK(2041) -> SEA"},
    {"Tail_Number":"N7835A","Flight_Number_Reporting_Airline":"3149","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2024-02-18T00:00:00Z","Hops":8,"Route":"CLE(621) -> BNA(801) -> PNS(1007) -> HOU(1234) -> MCI(1514) -> PHX(1747) -> BUR(1902) -> OAK(2117) -> DEN"},
    {"Tail_Number":"N263WN","Flight_Number_Reporting_Airline":"3149","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2024-02-04T00:00:00Z","Hops":8,"Route":"BUR(0) -> OAK(0) -> PHX(0) -> CLE(614) -> BNA(810) -> PNS(1010) -> HOU(1239) -> MCI(1514) -> PHX"},
    {"Tail_Number":"N429WN","Flight_Number_Reporting_Airline":"3149","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2024-01-28T00:00:00Z","Hops":8,"Route":"CLE(618) -> BNA(800) -> PNS(1012) -> HOU(1239) -> MCI(1515) -> PHX(1756) -> BUR(1859) -> OAK(2059) -> DEN"},
    {"Tail_Number":"N228WN","Flight_Number_Reporting_Airline":"3149","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2024-01-21T00:00:00Z","Hops":8,"Route":"CLE(620) -> BNA(810) -> PNS(1017) -> HOU(1237) -> MCI(1510) -> PHX(1800) -> BUR(1903) -> OAK(2102) -> DEN"},
    {"Tail_Number":"N569WN","Flight_Number_Reporting_Airline":"3149","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2024-01-14T00:00:00Z","Hops":8,"Route":"CLE(625) -> BNA(805) -> PNS(1011) -> HOU(1241) -> MCI(1607) -> PHX(1909) -> BUR(2006) -> OAK(2156) -> DEN"},
    {"Tail_Number":"N7742B","Flight_Number_Reporting_Airline":"154","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2023-04-30T00:00:00Z","Hops":8,"Route":"ELP(627) -> DAL(939) -> LIT(1124) -> ATL(1434) -> RIC(1656) -> MDW(1831) -> MCI(2040) -> PHX(2246) -> SAN"},
    {"Tail_Number":"N929WN","Flight_Number_Reporting_Airline":"154","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2023-04-16T00:00:00Z","Hops":8,"Route":"ELP(622) -> DAL(941) -> LIT(1122) -> ATL(1433) -> RIC(1659) -> MDW(1832) -> MCI(2041) -> PHX(2226) -> SAN"},
    {"Tail_Number":"N8631A","Flight_Number_Reporting_Airline":"2787","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2022-10-23T00:00:00Z","Hops":8,"Route":"MSY(608) -> ATL(944) -> CMH(1156) -> BWI(1354) -> RDU(1553) -> BNA(1729) -> DTW(2045) -> MDW(2149) -> LAX"},
    {"Tail_Number":"N8809L","Flight_Number_Reporting_Airline":"2787","IATA_CODE_Reporting_Airline":"WN","FlightDate":"2022-10-02T00:00:00Z","Hops":8,"Route":"MSY(610) -> ATL(917) -> CMH(1142) -> BWI(1355) -> RDU(1536) -> BNA(1716) -> DTW(2040) -> MDW(2147) -> LAX"}
  ],
  "generated_at": "2026-03-19T12:23:06.840438Z"
}
</script>

<script type="application/json" id="airport-coords">
{
  "ISP": {"name":"Long Island MacArthur","lat":40.7952,"lon":-73.1002},
  "BWI": {"name":"Baltimore/Washington Intl","lat":39.1754,"lon":-76.6683},
  "MYR": {"name":"Myrtle Beach Intl","lat":33.6797,"lon":-78.9283},
  "BNA": {"name":"Nashville Intl","lat":36.1245,"lon":-86.6782},
  "VPS": {"name":"Destin-Fort Walton Beach","lat":30.4832,"lon":-86.5254},
  "DAL": {"name":"Dallas Love Field","lat":32.8481,"lon":-96.8518},
  "LAS": {"name":"Harry Reid Intl (Las Vegas)","lat":36.0840,"lon":-115.1537},
  "OAK": {"name":"Oakland Intl","lat":37.7213,"lon":-122.2208},
  "SEA": {"name":"Seattle-Tacoma Intl","lat":47.4502,"lon":-122.3088},
  "CLE": {"name":"Cleveland Hopkins Intl","lat":41.4117,"lon":-81.8498},
  "PNS": {"name":"Pensacola Intl","lat":30.4734,"lon":-87.1866},
  "HOU": {"name":"William P. Hobby Airport","lat":29.6454,"lon":-95.2789},
  "MCI": {"name":"Kansas City Intl","lat":39.2976,"lon":-94.7139},
  "PHX": {"name":"Phoenix Sky Harbor Intl","lat":33.4373,"lon":-112.0078},
  "BUR": {"name":"Hollywood Burbank Airport","lat":34.2007,"lon":-118.3585},
  "DEN": {"name":"Denver Intl","lat":39.8561,"lon":-104.6737},
  "ELP": {"name":"El Paso Intl","lat":31.8072,"lon":-106.3779},
  "LIT": {"name":"Bill and Hillary Clinton Natl (Little Rock)","lat":34.7294,"lon":-92.2243},
  "ATL": {"name":"Hartsfield-Jackson Atlanta Intl","lat":33.6407,"lon":-84.4277},
  "RIC": {"name":"Richmond Intl","lat":37.5052,"lon":-77.3197},
  "MDW": {"name":"Chicago Midway Intl","lat":41.7868,"lon":-87.7522},
  "SAN": {"name":"San Diego Intl","lat":32.7338,"lon":-117.1933},
  "MSY": {"name":"Louis Armstrong New Orleans Intl","lat":29.9934,"lon":-90.2580},
  "CMH": {"name":"John Glenn Columbus Intl","lat":39.9980,"lon":-82.8919},
  "RDU": {"name":"Raleigh-Durham Intl","lat":35.8776,"lon":-78.7875},
  "DTW": {"name":"Detroit Metro Wayne County","lat":42.2124,"lon":-83.3534},
  "LAX": {"name":"Los Angeles Intl","lat":33.9425,"lon":-118.4081}
}
</script>

<!-- ══ Header ══════════════════════════════════════════════════════════ -->
<header>
  <div class="header-inner">
    <div class="header-badge">OnTime Analytics · Southwest Airlines</div>
    <h1>Highest Daily Hops for One Aircraft on One Flight Number</h1>
    <p class="subtitle">Maximum consecutive legs flown by a single tail number under one flight number on a single calendar day — all top-10 itineraries belong to Southwest Airlines (WN) with 8 hops each.</p>
  </div>
</header>

<!-- ══ Main ══════════════════════════════════════════════════════════ -->
<main>

  <!-- KPI Strip (anchored to top result) -->
  <section class="kpi-strip" id="kpi-strip">
    <div class="kpi-card highlight">
      <div class="kpi-label">Max Daily Hops</div>
      <div class="kpi-value" id="kpi-hops">—</div>
      <div class="kpi-sub">All top-10 itineraries</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Tail Number (Top)</div>
      <div class="kpi-value" id="kpi-tail" style="font-size:20px">—</div>
      <div class="kpi-sub">Most recent max-hop flight</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Flight (Top)</div>
      <div class="kpi-value" id="kpi-flight" style="font-size:20px">—</div>
      <div class="kpi-sub" id="kpi-carrier-sub">WN · Southwest Airlines</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Date (Top)</div>
      <div class="kpi-value" id="kpi-date" style="font-size:20px">—</div>
      <div class="kpi-sub" id="kpi-repeat-sub">—</div>
    </div>
  </section>

  <!-- Map + Detail row -->
  <section>
    <div class="map-detail-row">

      <!-- Map card -->
      <div class="panel">
        <div class="panel-header">
          <span id="map-title">Route Map</span>
          <span id="map-itinerary-label" style="font-size:11px;color:var(--muted);font-weight:400"></span>
        </div>
        <div id="map-degraded" class="map-degraded">
          <div>
            <div style="font-size:24px;margin-bottom:8px">✈</div>
            <div id="map-degraded-msg">Coordinates unavailable for one or more airports in this itinerary.</div>
          </div>
        </div>
        <div id="map"></div>
      </div>

      <!-- Route detail -->
      <div class="panel route-detail-panel">
        <div class="panel-header">
          <span>Stop Sequence</span>
          <span id="detail-hops-badge" style="background:var(--red);color:#fff;padding:2px 8px;border-radius:4px;font-size:11px;font-weight:700"></span>
        </div>
        <div class="panel-body" style="padding:10px 14px;overflow-y:auto;max-height:420px;">
          <ul class="stop-list" id="stop-list"></ul>
          <div class="route-meta" id="route-meta"></div>
        </div>
      </div>

    </div>
  </section>

  <!-- Legend -->
  <section>
    <div class="panel">
      <div class="panel-body" style="padding:12px 18px;">
        <div class="legend">
          <div class="legend-item">
            <div class="legend-dot" style="background:var(--navy)"></div>
            <span>Origin / stop airport</span>
          </div>
          <div class="legend-item">
            <div class="legend-dot" style="background:var(--red)"></div>
            <span>Final destination</span>
          </div>
          <div class="legend-item">
            <div class="legend-line" style="background:var(--sky)"></div>
            <span>Flight leg</span>
          </div>
          <div class="legend-item" style="margin-left:auto;font-style:italic">
            Click any table row to re-draw the map for that itinerary
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Itinerary Table -->
  <section>
    <div class="panel">
      <div class="panel-header">
        Top 10 Longest Daily Itineraries
        <span id="table-count" style="font-size:11px;color:var(--muted);font-weight:400"></span>
      </div>
      <div class="table-wrap">
        <table id="itinerary-table">
          <thead>
            <tr>
              <th>#</th>
              <th>Carrier</th>
              <th>Flight</th>
              <th>Tail</th>
              <th>Date</th>
              <th>Hops</th>
              <th>Route</th>
            </tr>
          </thead>
          <tbody id="table-body"></tbody>
        </table>
      </div>
    </div>
  </section>

</main>

<!-- ══ Footer / SQL ════════════════════════════════════════════════ -->
<footer>
  <div class="footer-inner">
    <div class="sql-toggle" id="sql-toggle" onclick="toggleSql()">
      <span id="sql-arrow">▶</span>
      <span>Primary Query · daily_itineraries ranked by hop count (LIMIT 10)</span>
    </div>
    <div class="sql-block" id="sql-block">
      <pre>WITH legs AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        Origin,
        Dest,
        assumeNotNull(DepTime) AS DepTime
    FROM ontime.ontime
    WHERE DepTime IS NOT NULL
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
),
daily_itineraries AS (
    SELECT
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate,
        count() AS Hops,
        arrayStringConcat(
            arrayConcat(
                arrayMap(
                    x -> concat(x.2, '(', toString(x.1), ')'),
                    arraySort(x -> x.1, groupArray((DepTime, Origin)))
                ),
                [argMax(Dest, DepTime)]
            ),
            ' -> '
        ) AS Route
    FROM legs
    GROUP BY
        Tail_Number,
        Flight_Number_Reporting_Airline,
        IATA_CODE_Reporting_Airline,
        FlightDate
)
SELECT
    Tail_Number,
    Flight_Number_Reporting_Airline,
    IATA_CODE_Reporting_Airline,
    FlightDate,
    Hops,
    Route
FROM daily_itineraries
ORDER BY Hops DESC, FlightDate DESC
LIMIT 10</pre>
    </div>
    <div class="footer-meta">
      Static artifact · Data embedded · Source: ontime.ontime via ClickHouse MCP ·
      Generated: 2026-03-19T12:23:06Z · q007_hops_per_day_semantic_discovery
    </div>
  </div>
</footer>

<!-- ══ Application Logic ═══════════════════════════════════════════ -->
<script>
(function () {
  'use strict';

  /* ── Parse embedded data ── */
  const resultData = JSON.parse(document.getElementById('result-data').textContent);
  const airportCoords = JSON.parse(document.getElementById('airport-coords').textContent);
  const rows = resultData.rows || [];

  /* ── State ── */
  let selectedIdx = 0;
  let map = null;
  let routeLayer = null;
  let markersLayer = null;

  /* ── Helpers ── */
  function normalizeDate(raw) {
    if (!raw) return '—';
    return raw.slice(0, 10); // YYYY-MM-DD from ISO timestamp
  }

  function formatTime(t) {
    if (t === null || t === undefined || t === 0) return null;
    const s = String(t).padStart(4, '0');
    return s.slice(0, 2) + ':' + s.slice(2);
  }

  function formatDateDisplay(raw) {
    const d = normalizeDate(raw);
    if (d === '—') return d;
    const parts = d.split('-');
    const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
    return `${months[parseInt(parts[1],10)-1]} ${parseInt(parts[2],10)}, ${parts[0]}`;
  }

  /* ── Parse route string into stops ── */
  function parseRoute(routeStr) {
    return routeStr.split(' -> ').map((seg, i, arr) => {
      const m = seg.match(/^([A-Z]{2,4})\((\d+)\)$/);
      if (m) return { code: m[1], time: parseInt(m[2], 10), isFinal: false };
      const m2 = seg.match(/^([A-Z]{2,4})$/);
      if (m2) return { code: m2[1], time: null, isFinal: true };
      return { code: seg, time: null, isFinal: i === arr.length - 1 };
    });
  }

  /* ── KPI strip (anchored to row 0) ── */
  function renderKPIs() {
    const top = rows[0];
    if (!top) return;
    document.getElementById('kpi-hops').textContent = top.Hops;
    document.getElementById('kpi-tail').textContent = top.Tail_Number;
    document.getElementById('kpi-flight').textContent = `WN ${top.Flight_Number_Reporting_Airline}`;
    document.getElementById('kpi-date').textContent = formatDateDisplay(top.FlightDate);
    const fn3149 = rows.filter(r => r.Flight_Number_Reporting_Airline === '3149').length;
    document.getElementById('kpi-repeat-sub').textContent =
      `WN3149 repeats ${fn3149}× in top 10`;
  }

  /* ── Route detail panel ── */
  function renderRouteDetail(row) {
    const stops = parseRoute(row.Route);
    const list = document.getElementById('stop-list');
    list.innerHTML = '';
    stops.forEach((stop, i) => {
      const li = document.createElement('li');
      li.className = 'stop-item';
      const numDiv = document.createElement('div');
      numDiv.className = 'stop-num' + (stop.isFinal ? ' final' : '');
      numDiv.textContent = i + 1;
      const infoDiv = document.createElement('div');
      infoDiv.style.flex = '1';
      const codeDiv = document.createElement('div');
      codeDiv.className = 'stop-code';
      codeDiv.textContent = stop.code;
      const nameDiv = document.createElement('div');
      nameDiv.className = 'stop-name';
      const coords = airportCoords[stop.code];
      nameDiv.textContent = coords ? coords.name : 'Unknown airport';
      infoDiv.appendChild(codeDiv);
      infoDiv.appendChild(nameDiv);
      const timeDiv = document.createElement('div');
      timeDiv.className = 'stop-time';
      const ft = stop.time !== null ? formatTime(stop.time) : null;
      timeDiv.textContent = stop.isFinal ? 'Arrival' : (ft ? `Dep ${ft}` : '—');
      li.appendChild(numDiv);
      li.appendChild(infoDiv);
      li.appendChild(timeDiv);
      list.appendChild(li);
    });
    const badge = document.getElementById('detail-hops-badge');
    badge.textContent = `${row.Hops} hops`;
    const meta = document.getElementById('route-meta');
    meta.textContent = `WN ${row.Flight_Number_Reporting_Airline} · ${row.Tail_Number} · ${formatDateDisplay(row.FlightDate)}`;
  }

  /* ── Map rendering ── */
  function initMap() {
    map = L.map('map', { zoomControl: true, scrollWheelZoom: false });
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
      maxZoom: 18
    }).addTo(map);
    routeLayer = L.layerGroup().addTo(map);
    markersLayer = L.layerGroup().addTo(map);
  }

  function getMarkerIcon(isFinal) {
    const color = isFinal ? '#c54f36' : '#0e3a52';
    const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 22 22">
      <circle cx="11" cy="11" r="9" fill="${color}" stroke="#fff" stroke-width="2"/>
    </svg>`;
    return L.divIcon({
      html: svg,
      className: '',
      iconSize: [22, 22],
      iconAnchor: [11, 11],
      popupAnchor: [0, -14]
    });
  }

  function renderMapRoute(row) {
    routeLayer.clearLayers();
    markersLayer.clearLayers();

    const stops = parseRoute(row.Route);
    const points = [];
    const missingCodes = [];

    stops.forEach(stop => {
      const c = airportCoords[stop.code];
      if (c) {
        points.push({ lat: c.lat, lon: c.lon, code: stop.code, name: c.name, time: stop.time, isFinal: stop.isFinal });
      } else {
        missingCodes.push(stop.code);
      }
    });

    const mapEl = document.getElementById('map');
    const degradedEl = document.getElementById('map-degraded');
    const degradedMsg = document.getElementById('map-degraded-msg');

    if (points.length < 2) {
      mapEl.style.display = 'none';
      degradedEl.style.display = 'flex';
      degradedMsg.textContent = missingCodes.length > 0
        ? `Coordinates unavailable for: ${missingCodes.join(', ')}`
        : 'Not enough coordinate data to render this route.';
      return;
    }

    mapEl.style.display = 'block';
    degradedEl.style.display = 'none';

    for (let i = 0; i < points.length - 1; i++) {
      const p1 = points[i], p2 = points[i + 1];
      L.polyline([[p1.lat, p1.lon], [p2.lat, p2.lon]], {
        color: '#3c88b5',
        weight: 2.5,
        opacity: 0.7
      }).addTo(routeLayer);
    }

    points.forEach((pt, idx) => {
      const icon = getMarkerIcon(pt.isFinal);
      const ft = pt.time !== null ? formatTime(pt.time) : null;
      const timeStr = pt.isFinal ? 'Final destination' : (ft ? `Dep ${ft}` : 'Dep —');
      const popup = `<div style="font-family:'Segoe UI',sans-serif;font-size:13px;min-width:140px">
        <strong style="font-size:15px;color:#0e3a52">${pt.code}</strong><br/>
        <span style="color:#5d7485;font-size:12px">${pt.name}</span><br/>
        <span style="color:#5c7080;font-size:12px">Stop ${idx + 1} · ${timeStr}</span>
      </div>`;
      L.marker([pt.lat, pt.lon], { icon }).addTo(markersLayer).bindPopup(popup);
    });

    const latLngs = points.map(p => [p.lat, p.lon]);
    try {
      map.fitBounds(L.latLngBounds(latLngs), { padding: [30, 30] });
    } catch (e) {
      map.setView([38, -96], 4);
    }

    if (missingCodes.length > 0) {
      const notice = L.control({ position: 'bottomleft' });
      notice.onAdd = function() {
        const d = L.DomUtil.create('div');
        d.style.cssText = 'background:rgba(255,255,255,.85);padding:6px 10px;border-radius:6px;font-size:11px;color:#d48a1f;max-width:220px';
        d.textContent = `Note: ${missingCodes.join(', ')} not mapped (no coordinates)`;
        return d;
      };
      notice.addTo(map);
    }

    map.invalidateSize();
  }

  function updateMapHeader(row) {
    document.getElementById('map-title').textContent = 'Route Map';
    document.getElementById('map-itinerary-label').textContent =
      `WN ${row.Flight_Number_Reporting_Airline} · ${row.Tail_Number} · ${formatDateDisplay(row.FlightDate)}`;
  }

  function selectRow(idx) {
    selectedIdx = idx;
    const row = rows[idx];
    if (!row) return;

    const trs = document.querySelectorAll('#table-body tr');
    trs.forEach((tr, i) => {
      tr.classList.toggle('active', i === idx);
    });

    renderMapRoute(row);
    updateMapHeader(row);
    renderRouteDetail(row);
  }

  function renderTable() {
    const tbody = document.getElementById('table-body');
    tbody.innerHTML = '';

    const flightCounts = {};
    rows.forEach(r => {
      const k = `WN${r.Flight_Number_Reporting_Airline}`;
      flightCounts[k] = (flightCounts[k] || 0) + 1;
    });

    rows.forEach((row, idx) => {
      const tr = document.createElement('tr');
      if (idx === selectedIdx) tr.classList.add('active');
      tr.onclick = () => selectRow(idx);

      const date = formatDateDisplay(row.FlightDate);
      const flightKey = `WN${row.Flight_Number_Reporting_Airline}`;
      const repeatCount = flightCounts[flightKey];

      tr.innerHTML = `
        <td>${idx + 1}</td>
        <td>${row.IATA_CODE_Reporting_Airline}</td>
        <td><span class="badge hops8">WN ${row.Flight_Number_Reporting_Airline}</span></td>
        <td style="font-variant-numeric:tabular-nums;white-space:nowrap">${row.Tail_Number}</td>
        <td style="white-space:nowrap">${date}</td>
        <td><span class="badge hops8">${row.Hops}</span></td>
        <td class="route-cell"><strong>${row.Route.split(' -> ')[0]}&nbsp;→&nbsp;${row.Route.split(' -> ').pop()}</strong>
          ${repeatCount > 1 ? `<br/><span style="color:var(--amber);font-size:11px">✕${repeatCount} occurrences in top 10</span>` : ''}
        </td>
      `;
      tbody.appendChild(tr);
    });

    document.getElementById('table-count').textContent = `${rows.length} itineraries · all 8 hops`;
  }

  window.toggleSql = function() {
    const block = document.getElementById('sql-block');
    const arrow = document.getElementById('sql-arrow');
    const visible = block.style.display === 'block';
    block.style.display = visible ? 'none' : 'block';
    arrow.textContent = visible ? '▶' : '▼';
  };

  function init() {
    if (!rows.length) {
      document.querySelector('main').innerHTML =
        '<div style="padding:40px;color:var(--red);font-size:16px">No data available to render.</div>';
      return;
    }

    renderKPIs();
    renderTable();
    initMap();
    selectRow(0);
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
</script>

</body>
</html>
```
