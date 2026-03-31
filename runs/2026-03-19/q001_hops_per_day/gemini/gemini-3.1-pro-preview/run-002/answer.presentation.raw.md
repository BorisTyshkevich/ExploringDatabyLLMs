```html
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Highest daily hops for one aircraft on one flight number</title>
    <!-- Tailwind CSS -->
    <script src="https://cdn.tailwindcss.com"></script>
    <!-- Leaflet CSS & JS -->
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY=" crossorigin=""/>
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js" integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo=" crossorigin=""></script>
    <!-- FontAwesome for icons -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    
    <style>
        #map { height: 450px; width: 100%; border-radius: 0.5rem; z-index: 1; }
        .ledger-entry { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace; font-size: 0.75rem; border-bottom: 1px solid #374151; padding: 0.5rem 0; word-break: break-all; }
        .ledger-entry:last-child { border-bottom: none; }
        .active-row { background-color: #eff6ff !important; border-left: 4px solid #3b82f6; }
        td { transition: background-color 0.2s; }
        
        .degraded-overlay {
            position: absolute; top: 0; left: 0; right: 0; bottom: 0;
            background: rgba(255, 255, 255, 0.85); z-index: 1000;
            display: flex; align-items: center; justify-content: center;
            flex-direction: column; text-align: center; padding: 1rem;
            display: none; border-radius: 0.5rem;
            backdrop-filter: blur(2px);
        }
        
        /* Custom map icon styling */
        .route-stop-icon {
            color: white; border-radius: 50%; width: 24px; height: 24px; 
            display: flex; align-items: center; justify-content: center; 
            font-size: 11px; font-weight: bold; border: 2px solid white; 
            box-shadow: 0 1px 4px rgba(0,0,0,0.4);
        }
    </style>
</head>
<body class="bg-gray-100 text-gray-800 p-4 md:p-8 font-sans flex flex-col min-h-screen">

<div class="max-w-7xl mx-auto w-full flex-grow flex flex-col gap-6">
    <header class="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
        <h1 class="text-2xl md:text-3xl font-bold text-gray-900">Highest daily hops for one aircraft on one flight number</h1>
        <p class="text-gray-500 mt-2">Interactive analysis of the longest multi-leg itineraries flown by a single aircraft in a single day.</p>
    </header>

    <!-- Error container -->
    <div id="error-message" class="hidden bg-red-50 border-l-4 border-red-500 text-red-700 p-4 rounded shadow-sm"></div>

    <!-- KPI Strip -->
    <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
        <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200 flex flex-col justify-center">
            <div class="text-xs text-gray-500 font-bold uppercase tracking-wider mb-1">Top Date</div>
            <div class="text-xl font-bold text-gray-900" id="kpi-date">--</div>
        </div>
        <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200 flex flex-col justify-center">
            <div class="text-xs text-gray-500 font-bold uppercase tracking-wider mb-1">Aircraft Tail</div>
            <div class="text-xl font-bold text-blue-600" id="kpi-tail">--</div>
        </div>
        <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200 flex flex-col justify-center">
            <div class="text-xs text-gray-500 font-bold uppercase tracking-wider mb-1">Flight Number</div>
            <div class="text-xl font-bold text-gray-900" id="kpi-flight">--</div>
        </div>
        <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200 flex flex-col justify-center">
            <div class="text-xs text-gray-500 font-bold uppercase tracking-wider mb-1">Max Hops</div>
            <div class="text-xl font-bold text-green-600" id="kpi-hops">--</div>
        </div>
        <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200 flex flex-col justify-center col-span-2 md:col-span-1">
            <div class="text-xs text-gray-500 font-bold uppercase tracking-wider mb-1">Route Context</div>
            <div class="text-sm font-bold text-gray-800 leading-tight" id="kpi-repetition">--</div>
        </div>
    </div>
    
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Map Area -->
        <div class="lg:col-span-2 bg-white p-4 rounded-lg shadow-sm border border-gray-200 flex flex-col">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-lg font-bold text-gray-800"><i class="fas fa-map-marker-alt text-blue-500 mr-2"></i>Itinerary Map</h2>
                <div class="text-xs bg-gray-100 px-2 py-1 rounded text-gray-600 border border-gray-200">Selected Route Viewer</div>
            </div>
            
            <div class="relative flex-grow rounded overflow-hidden border border-gray-300">
                <div id="map"></div>
                <div id="map-degraded-overlay" class="degraded-overlay">
                    <i class="fas fa-satellite-dish text-4xl text-orange-400 mb-3"></i>
                    <h3 class="font-bold text-xl text-gray-800 mb-1">Map Degraded</h3>
                    <p class="text-gray-600 max-w-sm mt-1 text-sm">Airport coordinates are missing or enrichment failed for this route. Cannot draw the full path.</p>
                </div>
            </div>
            
            <!-- Map Legend -->
            <div class="mt-4 flex flex-wrap gap-4 text-xs font-medium text-gray-600 bg-gray-50 p-3 rounded border border-gray-200">
                <div class="flex items-center"><span class="w-3 h-3 rounded-full bg-green-500 mr-2 shadow-sm"></span> Origin</div>
                <div class="flex items-center"><span class="w-3 h-3 rounded-full bg-blue-500 mr-2 shadow-sm"></span> Connection</div>
                <div class="flex items-center"><span class="w-3 h-3 rounded-full bg-red-500 mr-2 shadow-sm"></span> Destination</div>
                <div class="ml-auto flex items-center italic text-gray-500"><i class="fas fa-mouse-pointer mr-1"></i> Click table rows below to update map</div>
            </div>
        </div>
        
        <!-- Route Sequence Panel -->
        <div class="bg-white flex flex-col rounded-lg shadow-sm border border-gray-200 h-[565px]">
            <div class="p-4 border-b border-gray-200 bg-gray-50 rounded-t-lg">
                <h2 class="text-lg font-bold text-gray-800"><i class="fas fa-list-ol text-blue-500 mr-2"></i>Route Sequence</h2>
            </div>
            <div class="p-4 overflow-y-auto flex-grow bg-white rounded-b-lg" id="route-detail-panel">
                <div class="flex flex-col items-center justify-center h-full text-gray-400">
                    <i class="fas fa-route text-3xl mb-3 text-gray-300"></i>
                    <p class="italic text-sm">Loading route details...</p>
                </div>
            </div>
        </div>
    </div>

    <!-- Data Table -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-4 border-b border-gray-200 bg-gray-50 flex justify-between items-center">
            <h2 class="text-lg font-bold text-gray-800"><i class="fas fa-table text-blue-500 mr-2"></i>Top 10 Longest Daily Itineraries</h2>
            <span class="text-xs font-semibold text-gray-500 bg-gray-200 px-2 py-1 rounded">Click row to select</span>
        </div>
        <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 text-sm">
                <thead class="bg-white">
                    <tr>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Date</th>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Carrier</th>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Flight</th>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Tail Number</th>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Hops</th>
                        <th scope="col" class="px-6 py-3 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Route String</th>
                    </tr>
                </thead>
                <tbody id="data-table-body" class="bg-white divide-y divide-gray-200 cursor-pointer">
                    <tr>
                        <td colspan="6" class="px-6 py-8 text-center text-gray-500 italic">Executing analytical query...</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>

    <!-- Control Panel & Ledger -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-2 border-t border-gray-300 pt-6">
        <div class="lg:col-span-1 bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <h3 class="text-sm font-bold text-gray-800 mb-3 uppercase tracking-wide"><i class="fas fa-key text-gray-400 mr-2"></i>Authentication</h3>
            <div class="flex flex-col gap-3">
                <label class="text-xs font-semibold text-gray-600" for="jwe-input">JWE Token (Required for API access)</label>
                <div class="flex flex-col sm:flex-row gap-2">
                    <input type="password" id="jwe-input" class="border border-gray-300 p-2 rounded flex-grow text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none" placeholder="Paste Altinity MCP JWE Token">
                    <button onclick="saveJWE()" class="bg-gray-800 text-white px-4 py-2 rounded text-sm hover:bg-gray-700 transition-colors whitespace-nowrap font-medium">Save & Reload</button>
                </div>
                <p class="text-xs text-gray-500 mt-1">Stored securely in browser localStorage.</p>
            </div>
        </div>
        
        <div class="lg:col-span-2 bg-gray-900 rounded-lg shadow-sm overflow-hidden flex flex-col border border-gray-800">
            <div class="bg-gray-800 p-2 px-4 flex justify-between items-center border-b border-gray-700">
                <h3 class="text-xs font-bold text-gray-300 uppercase tracking-wider"><i class="fas fa-terminal text-gray-400 mr-2"></i>Query Ledger</h3>
                <span class="flex h-2 w-2 relative">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
                </span>
            </div>
            <div id="query-ledger" class="p-3 h-40 overflow-y-auto font-mono text-xs text-gray-300 bg-gray-900 flex-grow">
                <!-- Ledger entries injected here -->
                <div class="text-gray-500 italic">Awaiting operations...</div>
            </div>
        </div>
    </div>
</div>

<script>
    const JWE_KEY = 'OnTimeAnalystDashboard::auth::jwe';
    const PRIMARY_SQL = `SELECT Date, Carrier, AircraftID, FlightNumber, Hops, arrayStringConcat(arrayPushBack(arrayMap(x -> x.2, sorted_legs), arrayMap(x -> x.3, sorted_legs)[-1]), ' - ') AS Route, arrayStringConcat(arrayMap(x -> toString(x.1), sorted_legs), ', ') AS DepartureTimes FROM ( SELECT FlightDate AS Date, Reporting_Airline AS Carrier, Tail_Number AS AircraftID, Flight_Number_Reporting_Airline AS FlightNumber, count() AS Hops, arraySort(x -> x.1, groupArray((DepTime, Origin, Dest))) AS sorted_legs FROM ontime.ontime WHERE Cancelled = 0 AND Tail_Number != '' AND Flight_Number_Reporting_Airline != '' GROUP BY Date, Carrier, AircraftID, FlightNumber ) ORDER BY Hops DESC, Date DESC LIMIT 10`;

    let primaryData = [];
    let airportDict = {}; 
    let mapInstance = null;
    let polylineLayer = null;
    let markerLayer = null;
    let currentSelectedIndex = -1;

    // Initialize Leaflet Map
    function initMap() {
        // Default to a view of the US
        mapInstance = L.map('map', { zoomControl: false }).setView([39.8283, -98.5795], 4);
        L.control.zoom({ position: 'bottomright' }).addTo(mapInstance);
        
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
            maxZoom: 18,
            className: 'map-tiles'
        }).addTo(mapInstance);
        
        polylineLayer = L.featureGroup().addTo(mapInstance);
        markerLayer = L.featureGroup().addTo(mapInstance);
    }

    function appendToLedger(title, query, statusId, isError = false) {
        const ledger = document.getElementById('query-ledger');
        if (ledger.innerHTML.includes('Awaiting operations...')) {
            ledger.innerHTML = '';
        }
        
        const entry = document.createElement('div');
        entry.className = 'ledger-entry';
        
        const queryText = query ? `<div class="text-gray-400 mt-1 mb-1 truncate cursor-help" title="${query.replace(/"/g, '&quot;')}">${query}</div>` : '';
        const statusClass = isError ? 'text-red-400' : 'text-blue-300';
        
        entry.innerHTML = `
            <div class="font-bold text-blue-400 flex items-center"><i class="fas fa-play text-[10px] mr-2"></i> ${title}</div>
            ${queryText}
            <div id="${statusId}" class="${statusClass} font-semibold">Running...</div>
        `;
        ledger.appendChild(entry);
        ledger.scrollTop = ledger.scrollHeight;
    }

    function updateLedgerStatus(statusId, message, isSuccess) {
        const el = document.getElementById(statusId);
        if (el) {
            el.innerText = message;
            el.className = isSuccess ? "text-green-400 font-semibold" : "text-red-400 font-semibold";
        }
    }

    async function executeQuery(query, title) {
        const jwe = localStorage.getItem(JWE_KEY);
        if (!jwe) throw new Error("Authentication missing.");
        
        const statusId = 'status-' + Date.now() + Math.floor(Math.random() * 1000);
        appendToLedger(title, query, statusId);
        
        const url = `https://mcp.demo.altinity.cloud/${jwe}/openapi/execute_query?query=${encodeURIComponent(query)}`;
        
        try {
            const response = await fetch(url);
            if (!response.ok) {
                if (response.status === 401 || response.status === 403) {
                    throw new Error(`Auth Error (${response.status}). Please check your JWE Token.`);
                }
                const errText = await response.text();
                throw new Error(`HTTP ${response.status}: ${errText}`);
            }
            
            const data = await response.json();
            updateLedgerStatus(statusId, `✓ Success (${data.rows} rows)`, true);
            return data.data;
        } catch (e) {
            updateLedgerStatus(statusId, `✗ Error: ${e.message}`, false);
            throw e;
        }
    }

    function showError(msg) {
        const el = document.getElementById('error-message');
        el.innerHTML = `<div class="flex items-center"><i class="fas fa-exclamation-circle mr-2"></i><span>${msg}</span></div>`;
        el.classList.remove('hidden');
    }

    function renderKPIs(topRow) {
        document.getElementById('kpi-date').innerText = topRow.Date || '--';
        document.getElementById('kpi-tail').innerText = topRow.AircraftID || '--';
        document.getElementById('kpi-flight').innerText = topRow.Carrier ? (topRow.Carrier + ' ' + topRow.FlightNumber) : '--';
        document.getElementById('kpi-hops').innerText = topRow.Hops || '--';
        
        if (topRow.Route && primaryData.length > 0) {
            const sameRouteCount = primaryData.filter(r => r.Route === topRow.Route).length;
            document.getElementById('kpi-repetition').innerText = `${sameRouteCount} of Top 10 share this exact route`;
        } else {
            document.getElementById('kpi-repetition').innerText = '--';
        }
    }

    function renderTable(data) {
        const tbody = document.getElementById('data-table-body');
        tbody.innerHTML = '';
        
        if (data.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="px-6 py-4 text-center text-gray-500">No data found.</td></tr>';
            return;
        }
        
        data.forEach((row, i) => {
            const tr = document.createElement('tr');
            tr.id = `row-${i}`;
            tr.className = "hover:bg-gray-50 transition-colors duration-150 ease-in-out border-l-4 border-transparent";
            tr.onclick = () => selectItinerary(i);
            
            tr.innerHTML = `
                <td class="px-6 py-3 whitespace-nowrap font-medium text-gray-900">${row.Date || ''}</td>
                <td class="px-6 py-3 whitespace-nowrap text-gray-700">${row.Carrier || ''}</td>
                <td class="px-6 py-3 whitespace-nowrap font-semibold text-gray-800">${row.FlightNumber || ''}</td>
                <td class="px-6 py-3 whitespace-nowrap text-blue-600 font-mono text-xs">${row.AircraftID || ''}</td>
                <td class="px-6 py-3 whitespace-nowrap font-bold text-green-600">${row.Hops || ''}</td>
                <td class="px-6 py-3 text-gray-500 font-mono text-[10px] sm:text-xs truncate max-w-xs" title="${row.Route}">${row.Route || ''}</td>
            `;
            tbody.appendChild(tr);
        });
    }

    function selectItinerary(index) {
        if (index < 0 || index >= primaryData.length) return;
        
        // Update table active row styling
        if (currentSelectedIndex >= 0) {
            const oldRow = document.getElementById(`row-${currentSelectedIndex}`);
            if (oldRow) {
                oldRow.classList.remove('active-row');
                oldRow.classList.remove('border-blue-500');
                oldRow.classList.add('border-transparent');
            }
        }
        currentSelectedIndex = index;
        const newRow = document.getElementById(`row-${index}`);
        if (newRow) {
            newRow.classList.add('active-row');
            newRow.classList.remove('border-transparent');
            newRow.classList.add('border-blue-500');
        }
        
        const row = primaryData[index];
        const routeParts = row.Route ? row.Route.split(' - ') : [];
        const depTimes = row.DepartureTimes ? row.DepartureTimes.split(', ') : [];
        
        // Render Route Sequence Panel
        let detailHtml = `
            <div class="mb-4 pb-3 border-b border-gray-100">
                <div class="text-sm text-gray-500 uppercase tracking-wide font-bold">Flight ${row.Carrier} ${row.FlightNumber}</div>
                <div class="text-xs text-gray-400 mt-1"><i class="far fa-calendar-alt mr-1"></i> ${row.Date} &nbsp;|&nbsp; <i class="fas fa-plane mr-1"></i> Tail: ${row.AircraftID}</div>
            </div>
            <div class="relative border-l-2 border-gray-200 ml-3 mt-4 space-y-6">
        `;
        
        for (let i = 0; i < routeParts.length; i++) {
            const apCode = routeParts[i].trim();
            const apInfo = airportDict[apCode];
            const name = apInfo && apInfo.name ? apInfo.name : 'Unknown Airport';
            
            const isOrigin = i === 0;
            const isDest = i === routeParts.length - 1;
            
            let timeStr = '';
            let timeIcon = '';
            
            if (isDest) {
                timeStr = 'Final Destination';
                timeIcon = 'fa-plane-arrival';
            } else {
                const depTimeRaw = depTimes[i] || '--';
                // format time if possible (e.g. 1430 to 14:30)
                const formattedTime = (depTimeRaw.length === 4 && !isNaN(depTimeRaw)) 
                    ? depTimeRaw.substring(0,2) + ':' + depTimeRaw.substring(2,4) 
                    : depTimeRaw;
                timeStr = `Departs: <span class="font-mono bg-gray-100 px-1 rounded">${formattedTime}</span>`;
                timeIcon = 'fa-plane-departure';
            }
            
            let colorClass = 'bg-blue-500 border-blue-200';
            if (isOrigin) colorClass = 'bg-green-500 border-green-200';
            else if (isDest) colorClass = 'bg-red-500 border-red-200';
            
            detailHtml += `
                <div class="relative pl-6">
                    <div class="absolute -left-[9px] top-1 w-4 h-4 rounded-full ${colorClass} border-4 shadow-sm"></div>
                    <div class="flex flex-col">
                        <div class="font-bold text-gray-800 flex items-center gap-2">
                            <span class="bg-gray-800 text-white text-[10px] px-1.5 py-0.5 rounded font-mono">${apCode}</span>
                            <span class="text-sm truncate" title="${name}">${name}</span>
                        </div>
                        <div class="text-xs text-gray-500 mt-1 flex items-center">
                            <i class="fas ${timeIcon} w-4 text-gray-400"></i> ${timeStr}
                        </div>
                    </div>
                </div>
            `;
        }
        detailHtml += `</div>`;
        document.getElementById('route-detail-panel').innerHTML = detailHtml;
        
        drawMap(routeParts);
    }

    function drawMap(routeParts) {
        if (!mapInstance) return;
        
        polylineLayer.clearLayers();
        markerLayer.clearLayers();
        
        const latlngs = [];
        let missingCoords = false;
        
        for (const codeStr of routeParts) {
            const code = codeStr.trim();
            const ap = airportDict[code];
            if (ap && ap.latitude && ap.longitude) {
                latlngs.push([ap.latitude, ap.longitude]);
            } else {
                missingCoords = true;
            }
        }
        
        const overlay = document.getElementById('map-degraded-overlay');
        if (missingCoords || latlngs.length < 2) {
            overlay.style.display = 'flex';
            const statusId = 'status-' + Date.now() + Math.floor(Math.random() * 1000);
            appendToLedger('Map Degraded', 'Missing coordinate data for one or more airports in the sequence.', statusId, true);
            updateLedgerStatus(statusId, '✓ Logged', false); // false for red color styling
        } else {
            overlay.style.display = 'none';
        }
        
        if (latlngs.length > 0) {
            // Draw markers
            for (let i=0; i < latlngs.length; i++) {
                const ll = latlngs[i];
                const code = routeParts[i].trim();
                const apName = airportDict[code] ? airportDict[code].name : '';
                
                const isStart = i === 0;
                const isEnd = i === latlngs.length - 1;
                
                let bgColor = '#3b82f6'; // blue
                if (isStart) bgColor = '#22c55e'; // green
                else if (isEnd) bgColor = '#ef4444'; // red
                
                const icon = L.divIcon({
                    className: 'custom-div-icon',
                    html: `<div class="route-stop-icon" style="background-color: ${bgColor};">${i+1}</div>`,
                    iconSize: [24, 24],
                    iconAnchor: [12, 12]
                });
                
                const popupContent = `
                    <div class="text-center font-sans">
                        <div class="font-bold text-lg border-b pb-1 mb-1">${code}</div>
                        <div class="text-xs text-gray-600 mb-2">${apName}</div>
                        <div class="text-xs bg-gray-100 rounded px-2 py-1 inline-block">Stop ${i+1} of ${latlngs.length}</div>
                    </div>
                `;
                
                L.marker(ll, { icon: icon }).bindPopup(popupContent).addTo(markerLayer);
            }
        }
        
        if (latlngs.length > 1) {
            // Draw dashed polyline
            const polyline = L.polyline(latlngs, {
                color: '#3b82f6', 
                weight: 3, 
                opacity: 0.8, 
                dashArray: '8, 8',
                lineJoin: 'round'
            });
            polylineLayer.addLayer(polyline);
            
            // Add subtle glow/background line
            const bgLine = L.polyline(latlngs, {
                color: '#93c5fd',
                weight: 6,
                opacity: 0.3,
                lineJoin: 'round'
            });
            polylineLayer.addLayer(bgLine);
            
            // Fit bounds
            mapInstance.fitBounds(polylineLayer.getBounds(), { padding: [40, 40], maxZoom: 6 });
        } else if (latlngs.length === 1) {
            mapInstance.setView(latlngs[0], 5);
        }
    }

    function saveJWE() {
        const val = document.getElementById('jwe-input').value.trim();
        if (val) {
            localStorage.setItem(JWE_KEY, val);
            location.reload();
        } else {
            showError("Please enter a valid JWE token.");
        }
    }

    async function initApp() {
        initMap();
        
        const token = localStorage.getItem(JWE_KEY);
        if (!token) {
            document.getElementById('jwe-input').value = '';
            document.getElementById('data-table-body').innerHTML = `
                <tr><td colspan="6" class="px-6 py-8 text-center text-red-500 font-bold"><i class="fas fa-lock mr-2"></i>Authentication Required. Please enter JWE token below.</td></tr>
            `;
            const statusId = 'status-init';
            appendToLedger('System Ready', 'Waiting for JWE token...', statusId, true);
            updateLedgerStatus(statusId, 'Action Required', false);
            return;
        }
        
        document.getElementById('jwe-input').value = '******** (Loaded from storage)';
        
        try {
            // 1. Run Primary Analytical Query
            primaryData = await executeQuery(PRIMARY_SQL, 'Primary Analytical Query');
            
            if (!primaryData || primaryData.length === 0) {
                renderTable([]);
                document.getElementById('route-detail-panel').innerHTML = `<div class="p-4 text-gray-500 text-center">No results found.</div>`;
                return;
            }
            
            // 2. Extract distinct airports for enrichment
            const airports = new Set();
            primaryData.forEach(row => {
                if (row.Route) {
                    row.Route.split(' - ').forEach(ap => airports.add(ap.trim()));
                }
            });
            
            // 3. Run Enrichment Query
            if (airports.size > 0) {
                const apList = Array.from(airports).map(ap => `'${ap.replace(/'/g, "''")}'`).join(', ');
                const enrSql = `SELECT code, name, latitude, longitude FROM ontime.airports_latest WHERE code IN (${apList})`;
                
                try {
                    const enrData = await executeQuery(enrSql, 'Airport-Coordinate Enrichment');
                    enrData.forEach(ap => {
                        airportDict[ap.code] = ap;
                    });
                } catch (e) {
                    console.error("Enrichment failed:", e);
                    // Map will fall back to degraded state during render
                }
            }
            
            // 4. Render UI
            renderKPIs(primaryData[0]);
            renderTable(primaryData);
            
            // Select top row automatically
            selectItinerary(0);
            
        } catch (e) {
            showError(e.message);
            document.getElementById('data-table-body').innerHTML = `
                <tr><td colspan="6" class="px-6 py-8 text-center text-red-500"><i class="fas fa-exclamation-triangle mr-2"></i>Failed to load data: ${e.message}</td></tr>
            `;
        }
    }

    // Start application on load
    document.addEventListener('DOMContentLoaded', initApp);
</script>
</body>
</html>
```
