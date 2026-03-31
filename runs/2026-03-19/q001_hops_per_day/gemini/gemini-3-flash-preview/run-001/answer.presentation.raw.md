I will begin by reading the `query.sql`, `report.md`, and `result.json` files to understand the data structure and content for the visualization.
I will now generate the `visual.html` file, incorporating the Leaflet library for mapping, dynamic data fetching from the Altinity MCP endpoint, and the interactive UI components as specified.

```html
<!doctype html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Highest daily hops for one aircraft on one flight number</title>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <style>
        :root {
            --primary: #003366;
            --secondary: #00509d;
            --accent: #ffcc00;
            --bg: #f4f7f6;
            --text: #333;
            --card-bg: #fff;
            --border: #ddd;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg);
            color: var(--text);
            margin: 0;
            display: flex;
            flex-direction: column;
            min-height: 100vh;
        }
        header {
            background-color: var(--primary);
            color: white;
            padding: 1rem 2rem;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        main {
            flex: 1;
            padding: 1.5rem;
            max-width: 1400px;
            margin: 0 auto;
            width: 100%;
            box-sizing: border-box;
        }
        .kpi-container {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
            gap: 1rem;
            margin-bottom: 1.5rem;
        }
        .kpi-card {
            background: var(--card-bg);
            padding: 1rem;
            border-radius: 8px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            border-top: 4px solid var(--secondary);
            text-align: center;
        }
        .kpi-label {
            font-size: 0.75rem;
            text-transform: uppercase;
            color: #666;
            margin-bottom: 0.5rem;
        }
        .kpi-value {
            font-size: 1.25rem;
            font-weight: bold;
            color: var(--primary);
        }
        .map-section {
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
            margin-bottom: 1.5rem;
        }
        #map-card {
            background: var(--card-bg);
            border-radius: 8px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            overflow: hidden;
            height: 500px;
            position: relative;
        }
        #map {
            height: 100%;
            width: 100%;
        }
        #map-overlay {
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(255,255,255,0.7);
            display: none;
            align-items: center;
            justify-content: center;
            z-index: 1000;
            font-weight: bold;
            color: var(--primary);
        }
        .details-container {
            display: grid;
            grid-template-columns: 1fr 2fr;
            gap: 1.5rem;
        }
        .detail-panel, .table-panel {
            background: var(--card-bg);
            border-radius: 8px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            padding: 1rem;
        }
        h2 {
            margin-top: 0;
            font-size: 1.1rem;
            border-bottom: 1px solid var(--border);
            padding-bottom: 0.5rem;
            margin-bottom: 1rem;
        }
        .route-sequence {
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
        }
        .route-step {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.5rem;
            background: #f9f9f9;
            border-radius: 4px;
            font-size: 0.9rem;
        }
        .step-num {
            background: var(--secondary);
            color: white;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 0.75rem;
            flex-shrink: 0;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 0.85rem;
        }
        th, td {
            text-align: left;
            padding: 0.75rem;
            border-bottom: 1px solid var(--border);
        }
        th {
            background: #f8f9fa;
            font-weight: 600;
        }
        tr.active {
            background-color: #e3f2fd;
            border-left: 4px solid var(--secondary);
        }
        tr.clickable:hover {
            background-color: #f1f1f1;
            cursor: pointer;
        }
        .ledger-panel {
            margin-top: 2rem;
            background: #2d2d2d;
            color: #ccc;
            padding: 1rem;
            border-radius: 8px;
            font-family: monospace;
            font-size: 0.8rem;
            max-height: 300px;
            overflow-y: auto;
        }
        .ledger-entry {
            margin-bottom: 1rem;
            border-left: 2px solid #555;
            padding-left: 0.5rem;
        }
        .ledger-sql {
            color: #88c0d0;
            word-break: break-all;
        }
        footer {
            background: #eee;
            padding: 1rem 2rem;
            border-top: 1px solid var(--border);
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 0.8rem;
        }
        .controls {
            display: flex;
            gap: 1rem;
        }
        input[type="password"] {
            padding: 0.25rem;
            border: 1px solid var(--border);
            border-radius: 4px;
        }
        .loading-spinner {
            display: inline-block;
            width: 1rem;
            height: 1rem;
            border: 2px solid rgba(0,0,0,0.1);
            border-radius: 50%;
            border-top-color: var(--secondary);
            animation: spin 1s linear infinite;
        }
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        .degraded-msg {
            color: #d32f2f;
            font-style: italic;
            font-size: 0.8rem;
            margin-top: 0.5rem;
        }
    </style>
</head>
<body>

<header>
    <h1>Highest daily hops for one aircraft on one flight number</h1>
</header>

<main>
    <div class="kpi-container" id="kpi-strip">
        <div class="kpi-card">
            <div class="kpi-label">Top Aircraft</div>
            <div class="kpi-value" id="kpi-tail">-</div>
        </div>
        <div class="kpi-card">
            <div class="kpi-label">Flight Number</div>
            <div class="kpi-value" id="kpi-flight">-</div>
        </div>
        <div class="kpi-card">
            <div class="kpi-label">Date</div>
            <div class="kpi-value" id="kpi-date">-</div>
        </div>
        <div class="kpi-card">
            <div class="kpi-label">Max Daily Hops</div>
            <div class="kpi-value" id="kpi-hops">-</div>
        </div>
        <div class="kpi-card">
            <div class="kpi-label">Route Repetition</div>
            <div class="kpi-value" id="kpi-rep">-</div>
        </div>
    </div>

    <div class="map-section">
        <div id="map-card">
            <div id="map"></div>
            <div id="map-overlay">Loading enrichment data...</div>
        </div>
    </div>

    <div class="details-container">
        <div class="detail-panel">
            <h2>Route Sequence</h2>
            <div id="route-details" class="route-sequence">
                <p>Select an itinerary to see details.</p>
            </div>
            <div id="enrichment-status"></div>
        </div>
        <div class="table-panel">
            <h2>Top 10 Itineraries</h2>
            <div style="overflow-x: auto;">
                <table id="itinerary-table">
                    <thead>
                        <tr>
                            <th>Aircraft</th>
                            <th>Flight</th>
                            <th>Carrier</th>
                            <th>Date</th>
                            <th>Hops</th>
                            <th>Route</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr><td colspan="6" style="text-align:center;">Loading data...</td></tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>

    <div class="ledger-panel">
        <h3>Query Ledger</h3>
        <div id="query-ledger"></div>
    </div>
</main>

<footer>
    <div>Altinity On-Time Analyst Dashboard</div>
    <div class="controls">
        <label>JWE: <input type="password" id="jwe-input" placeholder="Enter JWE token"></label>
        <button onclick="saveJWE()">Save</button>
    </div>
</footer>

<script>
    const STORAGE_KEY = 'OnTimeAnalystDashboard::auth::jwe';
    const PRIMARY_SQL = `SELECT Tail_Number AS AircraftID, Flight_Number_Reporting_Airline AS FlightNumber, Reporting_Airline AS Carrier, FlightDate AS Date, count() AS Hops, arrayStringConcat(arrayPushBack(arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((DepTime, Origin, Dest)))), (arraySort(x -> x.1, groupArray((DepTime, Origin, Dest)))[-1]).3), ' -> ') AS Route, arrayStringConcat(arrayMap(x -> lpad(toString(x.1), 4, '0'), arraySort(x -> x.1, groupArray((DepTime, Origin, Dest)))), ', ') AS DepTimes FROM ontime.ontime WHERE Cancelled = 0 AND Tail_Number != '' GROUP BY AircraftID, FlightNumber, Carrier, Date ORDER BY Hops DESC, Date DESC LIMIT 10`;
    
    let map, routeLayer, markerLayer;
    let mainData = [];
    let airportCoords = {};
    let selectedIdx = 0;

    function initMap() {
        map = L.map('map').setView([39.8283, -98.5795], 4);
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; OpenStreetMap contributors'
        }).addTo(map);
        routeLayer = L.layerGroup().addTo(map);
        markerLayer = L.layerGroup().addTo(map);
    }

    function saveJWE() {
        const val = document.getElementById('jwe-input').value;
        if (val) {
            localStorage.setItem(STORAGE_KEY, val);
            location.reload();
        }
    }

    function getJWE() {
        return localStorage.getItem(STORAGE_KEY);
    }

    async function executeQuery(sql, label) {
        const jwe = getJWE();
        if (!jwe) {
            addToLedger(label, sql, "Error: JWE missing. Please provide it in the footer.");
            return null;
        }

        const url = `https://mcp.demo.altinity.cloud/${jwe}/openapi/execute_query?query=${encodeURIComponent(sql)}`;
        
        try {
            const start = Date.now();
            const response = await fetch(url);
            const duration = Date.now() - start;
            
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            
            const result = await response.json();
            addToLedger(label, sql, `Success (${duration}ms, ${result.rows ? result.rows.length : 0} rows)`);
            return result;
        } catch (err) {
            addToLedger(label, sql, `Error: ${err.message}`);
            return null;
        }
    }

    function addToLedger(label, sql, status) {
        const ledger = document.getElementById('query-ledger');
        const div = document.createElement('div');
        div.className = 'ledger-entry';
        div.innerHTML = `
            <strong>[${new Date().toLocaleTimeString()}] ${label}</strong><br>
            <div class="ledger-sql">${sql}</div>
            <div style="color: ${status.startsWith('Error') ? '#ff6b6b' : '#a3be8c'}">${status}</div>
        `;
        ledger.prepend(div);
    }

    async function loadData() {
        document.getElementById('jwe-input').value = getJWE() || '';
        
        const result = await executeQuery(PRIMARY_SQL, "Primary Analytical Query");
        if (!result || !result.rows) return;

        mainData = result.rows;
        updateKPIs(mainData[0]);
        renderTable();
        
        // Enrich airport data
        await enrichAirports();
        
        // Initial selection
        selectItinerary(0);
    }

    function updateKPIs(topRow) {
        document.getElementById('kpi-tail').textContent = topRow.AircraftID;
        document.getElementById('kpi-flight').textContent = topRow.FlightNumber;
        document.getElementById('kpi-date').textContent = new Date(topRow.Date).toLocaleDateString();
        document.getElementById('kpi-hops').textContent = topRow.Hops;
        
        // Count repetitions of the top route
        const repCount = mainData.filter(r => r.Route === topRow.Route).length;
        document.getElementById('kpi-rep').textContent = `${repCount} occurrences in Top 10`;
    }

    async function enrichAirports() {
        const airports = new Set();
        mainData.forEach(row => {
            row.Route.split(' -> ').forEach(code => airports.add(code.trim()));
        });

        const codes = Array.from(airports).map(c => `'${c}'`).join(',');
        const sql = `SELECT code, latitude, longitude, name FROM ontime.airports_latest WHERE code IN (${codes})`;
        
        const result = await executeQuery(sql, "Airport Coordinate Enrichment");
        if (result && result.rows) {
            result.rows.forEach(r => {
                airportCoords[r.code] = { lat: r.latitude, lon: r.longitude, name: r.name };
            });
        }
    }

    function renderTable() {
        const tbody = document.querySelector('#itinerary-table tbody');
        tbody.innerHTML = '';
        
        mainData.forEach((row, idx) => {
            const tr = document.createElement('tr');
            tr.className = 'clickable' + (idx === selectedIdx ? ' active' : '');
            tr.onclick = () => selectItinerary(idx);
            
            tr.innerHTML = `
                <td>${row.AircraftID}</td>
                <td>${row.FlightNumber}</td>
                <td>${row.Carrier}</td>
                <td>${new Date(row.Date).toLocaleDateString()}</td>
                <td>${row.Hops}</td>
                <td title="${row.Route}">${row.Route.length > 40 ? row.Route.substring(0, 37) + '...' : row.Route}</td>
            `;
            tbody.appendChild(tr);
        });
    }

    function selectItinerary(idx) {
        selectedIdx = idx;
        const row = mainData[idx];
        
        // Update Table UI
        document.querySelectorAll('#itinerary-table tr').forEach((tr, i) => {
            if (i > 0) tr.classList.toggle('active', (i-1) === idx);
        });

        // Update Details Panel
        const codes = row.Route.split(' -> ');
        const times = row.DepTimes.split(', ');
        const container = document.getElementById('route-details');
        container.innerHTML = '';
        
        let missingCoords = [];
        
        codes.forEach((code, i) => {
            const isLast = i === codes.length - 1;
            const step = document.createElement('div');
            step.className = 'route-step';
            const info = airportCoords[code] ? airportCoords[code].name : 'Unknown Airport';
            const time = isLast ? '' : ` (Dep: ${times[i]})`;
            step.innerHTML = `
                <div class="step-num">${i + 1}</div>
                <div><strong>${code}</strong> - ${info}${time}</div>
            `;
            container.appendChild(step);
            
            if (!airportCoords[code]) missingCoords.push(code);
        });

        const statusDiv = document.getElementById('enrichment-status');
        if (missingCoords.length > 0) {
            statusDiv.innerHTML = `<div class="degraded-msg">Missing coordinates for: ${missingCoords.join(', ')}</div>`;
            addToLedger("Map Visualization", "N/A", `Degraded: Missing coords for ${missingCoords.length} airports in itinerary ${idx+1}`);
        } else {
            statusDiv.innerHTML = '';
        }

        // Update Map
        drawRoute(codes);
    }

    function drawRoute(codes) {
        routeLayer.clearLayers();
        markerLayer.clearLayers();
        
        const path = [];
        codes.forEach(code => {
            const coord = airportCoords[code];
            if (coord && coord.lat && coord.lon) {
                const latlng = [coord.lat, coord.lon];
                path.push(latlng);
                
                L.marker(latlng).addTo(markerLayer)
                    .bindPopup(`<b>${code}</b><br>${coord.name}`);
            }
        });

        if (path.length > 1) {
            const polyline = L.polyline(path, {
                color: 'var(--secondary)',
                weight: 3,
                opacity: 0.7,
                dashArray: '10, 10'
            }).addTo(routeLayer);
            
            map.fitBounds(polyline.getBounds(), { padding: [50, 50] });
        } else if (path.length === 1) {
            map.setView(path[0], 6);
        }
    }

    window.onload = () => {
        initMap();
        loadData();
    };
</script>

</body>
</html>
```
