```report
# {{question_title}}

**Report Generated:** {{generated_at}}
**Rows Analyzed:** {{row_count}}

## Data Overview
{{data_overview_md}}

## Analytical Summary

Based on the top 10 longest single-aircraft, single-flight-number itineraries:

*   **Maximum Hop Count & Pattern:** The result set captures the highest concentration of daily hops operated by a single tail and flight number. By assessing the `Route` sequences in the results, we can discern whether these extreme operational days are highly repetitive "milk-run" routes bouncing back and forth between a small subset of airports, or if they represent long linear point-to-point sequences stretching across the network.
*   **Lead Itinerary:** The foremost row highlights the single most intense, recent itinerary observed. This peak event outlines a specific Carrier and Aircraft operating a singular Flight Number over an extensive multi-stop route on a single `Date`. The sequence explicitly defines the entire operating chain from the initial origin to final destination. 
*   **Route Clustering:** Comparing the top 10 longest itineraries reveals whether extreme multi-hop schedules cluster around specific hub airports or if specific airlines employ this scheduling tactic more frequently than others. Repetition across dates for the same route suggests a structural routing strategy rather than a one-off disruption or aircraft routing anomaly.

## Result Details
{{result_table_md}}
```

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Highest daily hops for one aircraft on one flight number</title>
    <!-- Leaflet CSS & JS -->
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY=" crossorigin=""/>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js" integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo=" crossorigin=""></script>
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
            font-family: "Segoe UI", "Helvetica Neue", sans-serif;
            color: var(--ink);
        }
        body {
            margin: 0;
            padding: 0;
            background: linear-gradient(180deg, var(--bg-top) 0%, var(--bg-bottom) 100%);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }
        header {
            background: var(--navy);
            color: #fff;
            padding: 2rem;
            text-align: center;
        }
        header h1 {
            margin: 0;
            font-family: Georgia, ui-serif, serif;
            font-weight: normal;
        }
        header p {
            margin: 0.5rem 0 0;
            opacity: 0.8;
        }
        main {
            flex: 1;
            padding: 2rem;
            max-width: 1280px;
            margin: 0 auto;
            width: 100%;
            box-sizing: border-box;
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
        }
        .hidden {
            display: none !important;
        }
        .kpi-strip {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1.5rem;
        }
        .card {
            background: var(--panel);
            border-radius: var(--radius-lg);
            box-shadow: var(--shadow);
            padding: 1.5rem;
            position: relative;
            overflow: hidden;
        }
        .kpi-card h3 {
            margin: 0;
            font-size: 0.85rem;
            color: var(--muted);
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }
        .kpi-card .value {
            margin: 0.5rem 0 0;
            font-size: 2rem;
            font-weight: 600;
            color: var(--navy);
        }
        .kpi-card .sub {
            margin: 0.25rem 0 0;
            font-size: 0.85rem;
            color: var(--slate);
        }
        .grid-2 {
            display: grid;
            grid-template-columns: 2fr 1fr;
            gap: 1.5rem;
        }
        @media (max-width: 900px) {
            .grid-2 { grid-template-columns: 1fr; }
        }
        .map-container {
            height: 400px;
            width: 100%;
            background: var(--panel-alt);
            border-radius: var(--radius-md);
            border: 1px solid var(--border);
            position: relative;
        }
        .map-overlay {
            position: absolute;
            inset: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            background: rgba(255,255,255,0.8);
            z-index: 1000;
            border-radius: var(--radius-md);
            text-align: center;
            padding: 2rem;
            color: var(--muted);
            font-weight: bold;
        }
        .map-card h2 {
            margin: 0 0 1rem;
            font-family: Georgia, ui-serif, serif;
            font-weight: normal;
        }
        .list-card h2 {
            margin: 0 0 1rem;
            font-family: Georgia, ui-serif, serif;
            font-weight: normal;
        }
        .route-table {
            width: 100%;
            border-collapse: collapse;
            font-size: 0.9rem;
        }
        .route-table th, .route-table td {
            text-align: left;
            padding: 0.75rem 0.5rem;
            border-bottom: 1px solid var(--border);
        }
        .route-table th {
            color: var(--muted);
            font-weight: 600;
        }
        .route-table tr:hover {
            background: var(--panel-alt);
        }
        .ledger-panel {
            background: var(--panel);
            border-radius: var(--radius-lg);
            box-shadow: var(--shadow);
            padding: 1.5rem;
        }
        .ledger-panel h2 {
            margin: 0 0 1rem;
            font-family: Georgia, ui-serif, serif;
            font-weight: normal;
        }
        .ledger-row {
            display: grid;
            grid-template-columns: 1.5rem 1fr 8rem 6rem 5rem;
            gap: 0.5rem;
            align-items: center;
            padding: 0.5rem;
            border-bottom: 1px solid var(--border);
            font-size: 0.85rem;
            cursor: pointer;
        }
        .ledger-row:hover {
            background: var(--panel-alt);
        }
        .ledger-row .status-ok { color: var(--teal); font-weight: 600; }
        .ledger-row .status-pending { color: var(--amber); font-weight: 600; }
        .ledger-row .status-failed { color: var(--red); font-weight: 600; }
        .ledger-sql {
            display: none;
            padding: 1rem;
            background: var(--panel-alt);
            border-radius: var(--radius-sm);
            margin: 0.5rem;
            font-family: monospace;
            white-space: pre-wrap;
            font-size: 0.8rem;
            color: var(--navy);
            border: 1px solid var(--border);
        }
        .ledger-item.open .ledger-sql {
            display: block;
        }
        .controls-footer {
            background: var(--ink);
            color: #fff;
            padding: 2rem;
            margin-top: auto;
        }
        .controls-container {
            max-width: 1280px;
            margin: 0 auto;
            display: flex;
            flex-direction: column;
            gap: 1rem;
        }
        .controls-row {
            display: flex;
            gap: 1rem;
            align-items: center;
            flex-wrap: wrap;
        }
        .controls-footer input[type="password"] {
            padding: 0.5rem;
            border: 1px solid var(--slate);
            background: var(--navy);
            color: #fff;
            border-radius: var(--radius-sm);
            flex: 1;
            min-width: 200px;
        }
        .controls-footer button {
            padding: 0.5rem 1rem;
            border: none;
            border-radius: var(--radius-sm);
            cursor: pointer;
            font-weight: 600;
        }
        .btn-primary {
            background: var(--teal);
            color: #fff;
        }
        .btn-primary:hover { background: #1a7861; }
        .btn-secondary {
            background: var(--slate);
            color: #fff;
        }
        .btn-secondary:hover { background: #4a5a67; }
        .controls-footer textarea {
            width: 100%;
            height: 120px;
            padding: 0.5rem;
            border: 1px solid var(--slate);
            background: var(--navy);
            color: #fff;
            border-radius: var(--radius-sm);
            font-family: monospace;
            box-sizing: border-box;
            resize: vertical;
        }
        .error-message {
            background: #ffebe6;
            color: var(--red);
            padding: 1rem;
            border-radius: var(--radius-sm);
            border: 1px solid #ffbdad;
            margin-bottom: 1rem;
        }
        .badge {
            display: inline-block;
            padding: 0.15rem 0.5rem;
            border-radius: 12px;
            background: var(--sky);
            color: #fff;
            font-size: 0.75rem;
            font-weight: bold;
        }
    </style>
</head>
<body>

    <header>
        <h1>Highest Daily Hops Tracker</h1>
        <p>Analyzing the most intensive single-aircraft, single-flight-number operations in a day</p>
    </header>

    <main id="dashboard" class="hidden">
        <div id="error-container" class="hidden error-message"></div>

        <div class="kpi-strip">
            <div class="card kpi-card">
                <h3>Lead Tail & Flight</h3>
                <div class="value" id="kpi-tail">-</div>
                <div class="sub" id="kpi-carrier">-</div>
            </div>
            <div class="card kpi-card">
                <h3>Flight Date</h3>
                <div class="value" id="kpi-date">-</div>
                <div class="sub">Most recent max-hop sequence</div>
            </div>
            <div class="card kpi-card">
                <h3>Max Hop Count</h3>
                <div class="value" id="kpi-hops">-</div>
                <div class="sub" id="kpi-segments">Segments connected</div>
            </div>
            <div class="card kpi-card">
                <h3>Route Pattern</h3>
                <div class="value" id="kpi-pattern" style="font-size:1.5rem;">-</div>
                <div class="sub">Linear vs Milk-run</div>
            </div>
        </div>

        <div class="grid-2">
            <div class="card map-card">
                <h2>Lead Itinerary Map</h2>
                <div id="map" class="map-container">
                    <div id="map-overlay" class="map-overlay">
                        Awaiting map coordinates...
                    </div>
                </div>
            </div>
            <div class="card list-card">
                <h2>Route Sequence</h2>
                <table class="route-table">
                    <thead>
                        <tr>
                            <th>Stop</th>
                            <th>Airport</th>
                            <th>Departure</th>
                        </tr>
                    </thead>
                    <tbody id="sequence-tbody">
                        <!-- Filled by JS -->
                    </tbody>
                </table>
            </div>
        </div>

        <div class="card list-card">
            <h2>Top 10 Longest Itineraries</h2>
            <table class="route-table">
                <thead>
                    <tr>
                        <th>Date</th>
                        <th>Tail Number</th>
                        <th>Flight No.</th>
                        <th>Carrier</th>
                        <th>Hops</th>
                        <th>Full Route</th>
                    </tr>
                </thead>
                <tbody id="top-tbody">
                    <!-- Filled by JS -->
                </tbody>
            </table>
        </div>

        <div class="ledger-panel">
            <h2>Query Ledger</h2>
            <div id="ledger-list"></div>
        </div>
    </main>

    <div class="controls-footer">
        <div class="controls-container">
            <div class="controls-row">
                <label for="jwe-token">MCP Token (JWE):</label>
                <input type="password" id="jwe-token" placeholder="eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0...">
                <button type="button" class="btn-primary" id="btn-fetch">Run Query</button>
                <button type="button" class="btn-secondary" id="btn-forget">Forget Token</button>
                <span id="fetch-status">Ready</span>
            </div>
            <div class="controls-row" style="align-items: flex-start;">
                <textarea id="sql-query" spellcheck="false">
SELECT
    Tail_Number AS `Aircraft ID`,
    Flight_Number_Reporting_Airline AS `Flight Number`,
    Reporting_Airline AS `Carrier`,
    FlightDate AS `Date`,
    arrayStringConcat(
        arrayPushBack(
            arrayMap(x -> x.2, arraySort(groupArray(tuple(DepTime, Origin)))),
            argMax(Dest, DepTime)
        ),
        '-'
    ) AS Route,
    arraySort(groupArray(DepTime)) AS `Actual Departure Times`
FROM default.ontime_v2
WHERE length(Tail_Number) > 0 
  AND DepTime IS NOT NULL 
  AND length(toString(DepTime)) > 0
GROUP BY
    Tail_Number,
    Flight_Number_Reporting_Airline,
    Reporting_Airline,
    FlightDate
ORDER BY
    count() DESC,
    FlightDate DESC
LIMIT 10
</textarea>
            </div>
        </div>
    </div>

    <script>
        const AUTH_KEY = 'OnTimeAnalystDashboard::auth::jwe';
        const BASE_URL = 'https://mcp.demo.altinity.cloud';
        let mapInstance = null;
        let routeLayer = null;

        // Init
        document.addEventListener('DOMContentLoaded', () => {
            const stored = localStorage.getItem(AUTH_KEY);
            if (stored) document.getElementById('jwe-token').value = stored;

            document.getElementById('btn-fetch').addEventListener('click', runPrimaryQuery);
            document.getElementById('btn-forget').addEventListener('click', () => {
                localStorage.removeItem(AUTH_KEY);
                document.getElementById('jwe-token').value = '';
                document.getElementById('fetch-status').textContent = 'Token forgotten.';
            });
        });

        function normalizeDate(val) {
            return String(val ?? '').slice(0, 10);
        }

        async function fetchMCP(token, sql) {
            const url = `${BASE_URL}/${token}/openapi/execute_query?query=${encodeURIComponent(sql)}`;
            const res = await fetch(url);
            if (!res.ok) {
                const text = await res.text();
                throw new Error(`HTTP ${res.status}: ${text}`);
            }
            return await res.json();
        }

        // Ledger functions
        const ledger = [];
        function addLedgerEntry(label, role, sql) {
            const id = 'ledger-' + Date.now() + Math.random().toString(36).substring(2);
            const entry = { id, label, role, status: 'Pending', rows: '-', sql };
            ledger.push(entry);
            renderLedger();
            return id;
        }
        function updateLedgerEntry(id, status, rows) {
            const entry = ledger.find(e => e.id === id);
            if (entry) {
                entry.status = status;
                entry.rows = rows;
                renderLedger();
            }
        }
        function renderLedger() {
            const container = document.getElementById('ledger-list');
            container.innerHTML = ledger.map(e => `
                <div class="ledger-item" id="${e.id}">
                    <div class="ledger-row" onclick="this.parentElement.classList.toggle('open')">
                        <span>▶</span>
                        <span style="font-weight:600; color:var(--navy);">${e.label}</span>
                        <span><span class="badge">${e.role}</span></span>
                        <span class="status-${e.status.toLowerCase()}">${e.status}</span>
                        <span style="color:var(--muted);">${e.rows} rows</span>
                    </div>
                    <div class="ledger-sql">${e.sql}</div>
                </div>
            `).join('');
        }

        async function runPrimaryQuery() {
            const token = document.getElementById('jwe-token').value.trim();
            const sql = document.getElementById('sql-query').value.trim();
            const statusEl = document.getElementById('fetch-status');
            const errEl = document.getElementById('error-container');
            const dashboard = document.getElementById('dashboard');

            if (!token || !sql) {
                statusEl.textContent = 'Token and SQL required.';
                return;
            }

            localStorage.setItem(AUTH_KEY, token);
            statusEl.textContent = 'Fetching...';
            errEl.classList.add('hidden');
            dashboard.classList.remove('hidden');

            const ledgerId = addLedgerEntry('Lead Itineraries', 'Primary', sql);

            try {
                const data = await fetchMCP(token, sql);
                const count = data.count ?? 0;
                let rows = [];
                if (data.rows && data.columns) {
                    rows = data.rows.map(row => {
                        const obj = {};
                        data.columns.forEach((col, i) => obj[col] = row[i]);
                        return obj;
                    });
                }
                updateLedgerEntry(ledgerId, 'OK', count);

                if (count === 0 || rows.length === 0) {
                    showError("Query returned 0 rows.");
                    return;
                }

                statusEl.textContent = `Success (${count} rows)`;
                renderDashboard(rows, token);

            } catch (err) {
                updateLedgerEntry(ledgerId, 'Failed', 0);
                statusEl.textContent = 'Error';
                showError(err.message);
            }
        }

        function showError(msg) {
            const errEl = document.getElementById('error-container');
            errEl.textContent = msg;
            errEl.classList.remove('hidden');
        }

        function renderDashboard(rows, token) {
            // Render Top 10
            const tbody = document.getElementById('top-tbody');
            tbody.innerHTML = rows.map(r => {
                const routeArr = (r['Route'] || '').split('-');
                const hops = Math.max(0, routeArr.length - 1);
                return `
                <tr>
                    <td>${normalizeDate(r['Date'])}</td>
                    <td>${r['Aircraft ID']}</td>
                    <td>${r['Flight Number']}</td>
                    <td>${r['Carrier']}</td>
                    <td>${hops}</td>
                    <td>${r['Route']}</td>
                </tr>
            `}).join('');

            // Focus Lead Itinerary (Row 0)
            const lead = rows[0];
            const routeStr = lead['Route'] || '';
            const routeArr = routeStr.split('-');
            const depTimes = lead['Actual Departure Times'] || [];
            const hops = Math.max(0, routeArr.length - 1);
            
            // Format times
            const formattedTimes = depTimes.map(t => {
                const ts = String(t).padStart(4, '0');
                return `${ts.slice(0,2)}:${ts.slice(2,4)}`;
            });
            
            // Unique airports vs total stops
            const uniqueAirports = new Set(routeArr).size;
            let patternType = uniqueAirports === routeArr.length ? 'Linear Route' : 'Cyclic / Milk-run';

            // Update KPIs
            document.getElementById('kpi-tail').textContent = lead['Aircraft ID'] || 'N/A';
            document.getElementById('kpi-carrier').textContent = `Carrier: ${lead['Carrier']} | Flight: ${lead['Flight Number']}`;
            document.getElementById('kpi-date').textContent = normalizeDate(lead['Date']);
            document.getElementById('kpi-hops').textContent = hops;
            document.getElementById('kpi-segments').textContent = `${routeArr.length} airports visited`;
            document.getElementById('kpi-pattern').textContent = patternType;
            
            // Route Sequence Table
            const seqBody = document.getElementById('sequence-tbody');
            seqBody.innerHTML = routeArr.map((apt, i) => {
                const isLast = i === routeArr.length - 1;
                const time = isLast ? 'N/A (Arrival)' : (formattedTimes[i] || 'Unknown');
                return `<tr>
                    <td>${i+1}</td>
                    <td><strong>${apt}</strong></td>
                    <td>${time}</td>
                </tr>`;
            }).join('');

            // Enrich and Map
            enrichCoordinatesAndMap(routeArr, token);
        }

        async function enrichCoordinatesAndMap(airports, token) {
            const uniqueCodes = [...new Set(airports)].filter(Boolean);
            if (uniqueCodes.length === 0) return;

            const codesStr = uniqueCodes.map(c => `'${c.replace(/'/g, "''")}'`).join(',');
            const sql = `SELECT iata AS IATA, lat AS Lat, lon AS Lon, city AS City, name AS Name FROM default.airports_bts WHERE iata IN (${codesStr})`;
            
            const ledgerId = addLedgerEntry('Airport Coordinates', 'Enrichment', sql);

            try {
                const data = await fetchMCP(token, sql);
                let coords = [];
                if (data.rows && data.columns) {
                    coords = data.rows.map(row => {
                        const obj = {};
                        data.columns.forEach((col, i) => obj[col] = row[i]);
                        return obj;
                    });
                }
                updateLedgerEntry(ledgerId, 'OK', coords.length);

                if (coords.length > 0) {
                    drawMap(airports, coords);
                } else {
                    document.getElementById('map-overlay').textContent = 'Map degraded: Coordinate enrichment returned empty.';
                }
            } catch (err) {
                updateLedgerEntry(ledgerId, 'Failed', 0);
                document.getElementById('map-overlay').textContent = 'Map degraded: Enrichment query failed.';
                console.error("Enrichment failed:", err);
            }
        }

        function drawMap(routeArr, coordsData) {
            document.getElementById('map-overlay').classList.add('hidden');
            
            const coordMap = {};
            coordsData.forEach(c => {
                if (c.IATA && c.Lat && c.Lon) {
                    coordMap[c.IATA] = [parseFloat(c.Lat), parseFloat(c.Lon)];
                }
            });

            const latlngs = routeArr.map(apt => coordMap[apt]).filter(Boolean);
            if (latlngs.length === 0) return;

            if (!mapInstance) {
                mapInstance = L.map('map').setView(latlngs[0], 4);
                L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
                    attribution: '&copy; OpenStreetMap contributors &copy; CARTO',
                    subdomains: 'abcd',
                    maxZoom: 20
                }).addTo(mapInstance);
            }

            if (routeLayer) {
                mapInstance.removeLayer(routeLayer);
            }

            routeLayer = L.featureGroup().addTo(mapInstance);

            // Draw lines
            if (latlngs.length > 1) {
                L.polyline(latlngs, {
                    color: '#c54f36',
                    weight: 3,
                    opacity: 0.8,
                    dashArray: '5, 5'
                }).addTo(routeLayer);
            }

            // Draw markers
            latlngs.forEach((ll, i) => {
                const isStart = i === 0;
                const isEnd = i === latlngs.length - 1;
                const color = (isStart || isEnd) ? '#c54f36' : '#0e3a52';
                
                L.circleMarker(ll, {
                    radius: 6,
                    fillColor: color,
                    color: '#fff',
                    weight: 2,
                    opacity: 1,
                    fillOpacity: 1
                }).bindPopup(`<b>${routeArr[i]}</b><br>Stop ${i+1}`).addTo(routeLayer);
            });

            mapInstance.fitBounds(routeLayer.getBounds(), { padding: [40, 40] });
            mapInstance.invalidateSize();
        }
    </script>
</body>
</html>
```
