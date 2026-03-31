import { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

// ── Constants ────────────────────────────────────────────────────────────────

const LS_JWE_KEY = 'OnTimeAnalystDashboard::auth::jwe';

const PRIMARY_SQL = `WITH ranked AS (
    SELECT
        FlightDate,
        Tail_Number,
        Flight_Number_Reporting_Airline AS FlightNum,
        IATA_CODE_Reporting_Airline AS Carrier,
        count() AS hop_count,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),
                [argMax(DestCode, assumeNotNull(DepTime))]
            ),
            '-'
        ) AS Route,
        arrayStringConcat(
            arrayConcat(
                arrayMap(x -> x.2, arraySort(x -> x.1, groupArray((assumeNotNull(DepTime), OriginCode)))),
                [argMax(DestCode, assumeNotNull(DepTime))]
            ),
            ','
        ) AS itinerary_sequence
    FROM ontime.fact_ontime
    WHERE Cancelled = 0
      AND Tail_Number != ''
      AND Flight_Number_Reporting_Airline != ''
    GROUP BY FlightDate, Tail_Number, Flight_Number_Reporting_Airline, IATA_CODE_Reporting_Airline
)
SELECT
    Tail_Number,
    FlightNum,
    Carrier,
    FlightDate,
    hop_count,
    Route,
    itinerary_sequence
FROM ranked
ORDER BY hop_count DESC, FlightDate DESC
LIMIT 10`;

function makeEnrichmentSQL(tail, flightNum, date) {
  const esc = s => String(s).replace(/'/g, "''");
  return `SELECT\n    f.OriginCode,\n    f.DestCode,\n    f.DepTime,\n    o.DisplayAirportName AS OriginName,\n    o.Latitude           AS OriginLat,\n    o.Longitude          AS OriginLon,\n    d.DisplayAirportName AS DestName,\n    d.Latitude           AS DestLat,\n    d.Longitude          AS DestLon\nFROM ontime.fact_ontime f\nLEFT JOIN ontime.dim_airports o ON f.OriginCode = o.AirportCode\nLEFT JOIN ontime.dim_airports d ON f.DestCode   = d.AirportCode\nWHERE f.Cancelled = 0\n  AND f.Tail_Number = '${esc(tail)}'\n  AND f.Flight_Number_Reporting_Airline = '${esc(flightNum)}'\n  AND f.FlightDate = '${esc(date)}'\nORDER BY assumeNotNull(f.DepTime)`;
}

// ── Utilities ────────────────────────────────────────────────────────────────

function buildUrl(jwe, sql) {
  return `https://mcp.demo.altinity.cloud/${jwe}/openapi/execute_query?query=${encodeURIComponent(sql)}`;
}

function normalizeDate(val) {
  return String(val ?? '').slice(0, 10);
}

function parseRows(data) {
  const cols = (data.columns || []).map(c => (typeof c === 'string' ? c : c.name));
  return (data.rows || []).map(r =>
    Object.fromEntries(cols.map((c, i) => [c, r[i]]))
  );
}

function rowCacheKey(row) {
  return `${row.Tail_Number}|${row.FlightNum}|${normalizeDate(row.FlightDate)}`;
}

function formatDepTime(t) {
  if (t == null || t === '') return '—';
  const n = parseInt(t, 10);
  if (isNaN(n)) return String(t);
  return `${String(Math.floor(n / 100)).padStart(2, '0')}:${String(n % 100).padStart(2, '0')}`;
}

function carrierName(code) {
  const names = { WN: 'Southwest Airlines', AA: 'American Airlines', DL: 'Delta Air Lines', UA: 'United Airlines' };
  return names[code] || code;
}

function buildHeroNarrative(rows) {
  if (!rows.length) return '';
  const lead = rows[0];
  const airports = (lead.itinerary_sequence || lead.Route?.replace(/-/g, ',')).split(',');
  const origin = airports[0] || '';
  const dest = airports[airports.length - 1] || '';
  const carriers = [...new Set(rows.map(r => r.Carrier))];
  const singleCarrier = carriers.length === 1;
  const routes = [...new Set(rows.map(r => r.Route))].length;
  const carrier = carrierName(lead.Carrier);

  return `${carrier} (${lead.Carrier}) flight ${lead.FlightNum}, tail ${lead.Tail_Number}, ` +
    `on ${normalizeDate(lead.FlightDate)} completed ${lead.hop_count} hops in a single day — ` +
    `departing ${origin} and arriving at ${dest} while making ${lead.hop_count - 1} intermediate stops. ` +
    (singleCarrier
      ? `All ${rows.length} top itineraries in this ranking belong to ${carrier}, `
      : `The top ${rows.length} itineraries span ${carriers.join(', ')}, `) +
    `covering ${routes} distinct routes. These are recurring scheduled patterns, not anomalies — ` +
    `a single flight number chains many short segments to sweep the full width of the country in one day.`;
}

// ── Map helpers ──────────────────────────────────────────────────────────────

const COLORS = {
  routeRank1: '#c54f36',
  routeOther: '#3c88b5',
  markerOrigin: '#1f8a70',
  markerDest: '#c54f36',
  markerStop: '#0e3a52',
};

function drawRoute(map, segs, isRank1) {
  if (!segs || segs.length === 0) return null;

  const points = [];
  segs.forEach((seg, i) => {
    const oLat = parseFloat(seg.OriginLat);
    const oLon = parseFloat(seg.OriginLon);
    const dLat = parseFloat(seg.DestLat);
    const dLon = parseFloat(seg.DestLon);
    if (i === 0 && !isNaN(oLat) && !isNaN(oLon)) {
      points.push({ code: seg.OriginCode, name: seg.OriginName, lat: oLat, lon: oLon, role: 'origin' });
    }
    if (!isNaN(dLat) && !isNaN(dLon)) {
      points.push({ code: seg.DestCode, name: seg.DestName, lat: dLat, lon: dLon, role: i === segs.length - 1 ? 'dest' : 'stop' });
    }
  });

  if (points.length < 2) return null;

  const group = L.layerGroup();
  const latlngs = points.map(p => [p.lat, p.lon]);

  L.polyline(latlngs, {
    color: isRank1 ? COLORS.routeRank1 : COLORS.routeOther,
    weight: 3.5,
    opacity: 0.85,
  }).addTo(group);

  points.forEach(p => {
    const fillColor =
      p.role === 'origin' ? COLORS.markerOrigin :
      p.role === 'dest'   ? COLORS.markerDest   :
                            COLORS.markerStop;
    const radius = p.role === 'origin' || p.role === 'dest' ? 9 : 6;

    L.circleMarker([p.lat, p.lon], {
      radius,
      color: '#fff',
      weight: 2,
      fillColor,
      fillOpacity: 0.92,
    })
      .bindTooltip(p.code, { permanent: true, direction: 'top', className: 'airport-label' })
      .bindPopup(`<b>${p.code}</b><br/>${p.name || ''}`)
      .addTo(group);
  });

  return { group, latlngs };
}

// ── App ──────────────────────────────────────────────────────────────────────

export default function App() {
  const [jwe, setJwe]               = useState(() => localStorage.getItem(LS_JWE_KEY) || '');
  const [sql, setSql]               = useState(PRIMARY_SQL);
  const [rows, setRows]             = useState([]);
  const [isRunning, setIsRunning]   = useState(false);
  const [statusMsg, setStatusMsg]   = useState('');
  const [hasData, setHasData]       = useState(false);
  const [selectedIdx, setSelectedIdx] = useState(0);
  const [coordCache, setCoordCache] = useState({});     // key → segs[] | null
  const [mapStatus, setMapStatus]   = useState('idle'); // idle | loading | ok | degraded
  const [ledger, setLedger]         = useState([]);

  const runIdRef        = useRef(0);
  const mapRef          = useRef(null);
  const mapContainerRef = useRef(null);
  const routeLayerRef   = useRef(null);

  // ── Ledger ────────────────────────────────────────────────────────────────

  const addLedgerEntry = useCallback((id, label, role, sqlText) => {
    setLedger(prev => {
      if (prev.find(e => e.id === id)) return prev;
      return [...prev, { id, label, role, sql: sqlText, status: 'Pending', rows: null, expanded: false }];
    });
  }, []);

  const updateLedgerEntry = useCallback((id, updates) => {
    setLedger(prev => prev.map(e => e.id === id ? { ...e, ...updates } : e));
  }, []);

  const toggleLedgerEntry = useCallback((id) => {
    setLedger(prev => prev.map(e => e.id === id ? { ...e, expanded: !e.expanded } : e));
  }, []);

  // ── Primary fetch ─────────────────────────────────────────────────────────

  const runPrimaryQuery = useCallback(async (token, sqlText) => {
    const thisRunId = ++runIdRef.current;
    setIsRunning(true);
    setStatusMsg('Fetching data…');
    setLedger([]);
    setHasData(false);
    setRows([]);
    setCoordCache({});
    setSelectedIdx(0);
    setMapStatus('idle');

    const entryId = 'primary';
    addLedgerEntry(entryId, 'Top 10 highest-hop itineraries', 'Primary', sqlText);

    try {
      const res = await fetch(buildUrl(token, sqlText));
      if (thisRunId !== runIdRef.current) return;

      if (!res.ok) {
        const text = await res.text();
        updateLedgerEntry(entryId, { status: 'Failed', rows: 0 });
        setStatusMsg(`Error ${res.status}: ${text.slice(0, 300)}`);
        return;
      }

      const data = await res.json();
      if (thisRunId !== runIdRef.current) return;

      const parsed = parseRows(data);
      updateLedgerEntry(entryId, { status: 'OK', rows: parsed.length });
      setRows(parsed);
      setHasData(true);
      setStatusMsg(`Loaded ${parsed.length} itinerary rows.`);
      localStorage.setItem(LS_JWE_KEY, token);
    } catch (err) {
      if (thisRunId !== runIdRef.current) return;
      updateLedgerEntry(entryId, { status: 'Failed', rows: 0 });
      setStatusMsg(`Fetch failed: ${err.message}`);
    } finally {
      if (thisRunId === runIdRef.current) setIsRunning(false);
    }
  }, [addLedgerEntry, updateLedgerEntry]);

  // ── Enrichment fetch ──────────────────────────────────────────────────────

  const runEnrichment = useCallback(async (row, token) => {
    const key     = rowCacheKey(row);
    const date    = normalizeDate(row.FlightDate);
    const enrichSql = makeEnrichmentSQL(row.Tail_Number, row.FlightNum, date);
    const entryId = `enrich_${key}`;

    setMapStatus('loading');
    addLedgerEntry(entryId, `Airport-coordinate enrichment — ${row.Route}`, 'Enrichment', enrichSql);

    try {
      const res = await fetch(buildUrl(token, enrichSql));
      if (!res.ok) {
        const text = await res.text();
        updateLedgerEntry(entryId, { status: 'Failed', rows: 0 });
        setCoordCache(prev => ({ ...prev, [key]: null }));
        setMapStatus('degraded');
        return;
      }
      const data = await res.json();
      const segs = parseRows(data);
      updateLedgerEntry(entryId, { status: 'OK', rows: segs.length });
      setCoordCache(prev => ({ ...prev, [key]: segs }));
      setMapStatus(segs.length > 0 ? 'ok' : 'degraded');
    } catch (err) {
      updateLedgerEntry(entryId, { status: 'Failed', rows: 0 });
      setCoordCache(prev => ({ ...prev, [key]: null }));
      setMapStatus('degraded');
    }
  }, [addLedgerEntry, updateLedgerEntry]);

  // ── Auto-load on mount ────────────────────────────────────────────────────

  useEffect(() => {
    const stored = localStorage.getItem(LS_JWE_KEY);
    if (stored) {
      setJwe(stored);
      runPrimaryQuery(stored, PRIMARY_SQL);
    }
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  // ── Enrichment trigger on selection change ────────────────────────────────

  useEffect(() => {
    if (!hasData || rows.length === 0) return;
    const row = rows[selectedIdx];
    if (!row) return;
    const key = rowCacheKey(row);

    if (Object.prototype.hasOwnProperty.call(coordCache, key)) {
      const segs = coordCache[key];
      setMapStatus(segs && segs.length > 0 ? 'ok' : 'degraded');
      return;
    }

    const token = jwe.trim() || localStorage.getItem(LS_JWE_KEY) || '';
    if (!token) {
      setMapStatus('degraded');
      return;
    }
    runEnrichment(row, token);
  }, [selectedIdx, hasData]); // eslint-disable-line react-hooks/exhaustive-deps

  // ── Map initialization ────────────────────────────────────────────────────

  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;
    const map = L.map(mapContainerRef.current, {
      center: [39, -97],
      zoom: 4,
      scrollWheelZoom: false,
    });
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
      maxZoom: 18,
    }).addTo(map);
    mapRef.current = map;
  }, []);

  // ── Map redraw ────────────────────────────────────────────────────────────

  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;

    if (routeLayerRef.current) {
      map.removeLayer(routeLayerRef.current);
      routeLayerRef.current = null;
    }

    if (mapStatus === 'loading' || mapStatus === 'idle' || !hasData) return;

    const row = rows[selectedIdx];
    if (!row) return;

    const key  = rowCacheKey(row);
    const segs = coordCache[key];
    if (!segs || segs.length === 0) return;

    const result = drawRoute(map, segs, selectedIdx === 0);
    if (!result) return;

    result.group.addTo(map);
    routeLayerRef.current = result.group;

    try {
      map.fitBounds(L.latLngBounds(result.latlngs).pad(0.12));
    } catch (_) { /* ignore */ }

    setTimeout(() => map.invalidateSize(), 50);
  }, [mapStatus, coordCache, selectedIdx, hasData, rows]);

  // Invalidate map after data loads (layout may shift)
  useEffect(() => {
    if (hasData && mapRef.current) {
      setTimeout(() => mapRef.current?.invalidateSize(), 150);
    }
  }, [hasData]);

  // ── Derived values ────────────────────────────────────────────────────────

  const leadRow     = rows[0] ?? null;
  const selectedRow = rows[selectedIdx] ?? null;

  const routeOccurrences = useMemo(() => {
    const counts = {};
    rows.forEach(r => {
      const k = `${r.FlightNum}|${r.Route}`;
      counts[k] = (counts[k] || 0) + 1;
    });
    return counts;
  }, [rows]);

  const leadOccurrences = leadRow
    ? (routeOccurrences[`${leadRow.FlightNum}|${leadRow.Route}`] ?? 1)
    : 1;

  const heroNarrative = useMemo(() => buildHeroNarrative(rows), [rows]);

  const selectedSegs = selectedRow
    ? (coordCache[rowCacheKey(selectedRow)] ?? [])
    : [];

  const selectedAirports = selectedRow
    ? ((selectedRow.itinerary_sequence || selectedRow.Route?.replace(/-/g, ',')).split(','))
    : [];

  // ── Handlers ──────────────────────────────────────────────────────────────

  function handleFetch() {
    if (isRunning) return;
    const token = jwe.trim();
    if (!token) { setStatusMsg('Please enter a JWE token.'); return; }
    if (!sql.trim()) { setStatusMsg('SQL is empty.'); return; }
    runPrimaryQuery(token, sql);
  }

  function handleForget() {
    localStorage.removeItem(LS_JWE_KEY);
    setJwe('');
    setStatusMsg('Token cleared from storage.');
  }

  function handleRowClick(idx) {
    if (idx === selectedIdx) return;
    setSelectedIdx(idx);
  }

  // ── Render ────────────────────────────────────────────────────────────────

  const mapTitle = selectedRow
    ? selectedIdx === 0
      ? `${selectedRow.Route} (Rank #1)`
      : selectedRow.Route
    : 'Route Map';

  return (
    <div className="app">
      {/* ── Header ── */}
      <header className="page-header">
        <div className="header-inner">
          <div className="header-badge">✈ OnTime Analytics</div>
          <h1 className="header-title">Highest Daily Hops for One Aircraft</h1>
          <p className="header-sub">Top itineraries by hop count — one flight number, one tail, one day</p>
        </div>
      </header>

      {/* ── Main ── */}
      <main className="page-main">

        {/* Empty state */}
        {!hasData && !isRunning && (
          <div className="empty-state">
            <div className="empty-icon">✈</div>
            <h2>Ready to load data</h2>
            <p>Enter your JWE token in the controls below and click <strong>Run Query</strong> to fetch the top itineraries.</p>
          </div>
        )}

        {/* Loading state */}
        {isRunning && !hasData && (
          <div className="loading-state">
            <div className="spinner" />
            <p>Loading itinerary data…</p>
          </div>
        )}

        {/* Dashboard */}
        {hasData && leadRow && (
          <>
            {/* ── Hero (anchored to rank 1) ── */}
            <section className="hero-section">
              <div className="hero-card">
                <div className="hero-eyebrow">Champion Itinerary — Rank #1</div>
                <div className="hero-route">{leadRow.Route}</div>
                <p className="hero-narrative">{heroNarrative}</p>
              </div>
            </section>

            {/* ── KPI strip (anchored to rank 1) ── */}
            <section className="kpi-strip">
              <div className="kpi-card">
                <div className="kpi-label">Tail Number</div>
                <div className="kpi-value">{leadRow.Tail_Number}</div>
                <div className="kpi-sub">Aircraft identifier</div>
              </div>
              <div className="kpi-card">
                <div className="kpi-label">Flight Number</div>
                <div className="kpi-value">{leadRow.Carrier} {leadRow.FlightNum}</div>
                <div className="kpi-sub">Carrier + flight #</div>
              </div>
              <div className="kpi-card">
                <div className="kpi-label">Date</div>
                <div className="kpi-value">{normalizeDate(leadRow.FlightDate)}</div>
                <div className="kpi-sub">Flight date</div>
              </div>
              <div className="kpi-card kpi-highlight">
                <div className="kpi-label">Max Hops</div>
                <div className="kpi-value kpi-big">{leadRow.hop_count}</div>
                <div className="kpi-sub">Segments in one day</div>
              </div>
              <div className="kpi-card">
                <div className="kpi-label">Route Pattern</div>
                <div className="kpi-value">{leadOccurrences}×</div>
                <div className="kpi-sub">Times in top 10</div>
              </div>
            </section>

            {/* ── Map ── */}
            <section className="map-section">
              <div className="panel map-panel">
                <div className="panel-header">
                  <div>
                    <h3 className="panel-title">Route Map</h3>
                    <div className="panel-sub map-subtitle">{mapTitle}</div>
                  </div>
                  <div className="map-legend">
                    <span className="legend-item"><span className="legend-dot dot-origin" /> Origin</span>
                    <span className="legend-item"><span className="legend-dot dot-stop" /> Stop</span>
                    <span className="legend-item"><span className="legend-dot dot-dest" /> Destination</span>
                    <span className="legend-item"><span className="legend-line line-rank1" /> Rank #1</span>
                    <span className="legend-item"><span className="legend-line line-other" /> Other</span>
                  </div>
                </div>

                <div className="map-wrap">
                  {mapStatus === 'loading' && (
                    <div className="map-overlay">
                      <div className="spinner spinner-sm" />
                      <span>Loading coordinates…</span>
                    </div>
                  )}
                  {mapStatus === 'degraded' && (
                    <div className="map-overlay map-degraded">
                      <span>⚠ Coordinate enrichment unavailable for this itinerary — non-map analysis continues below.</span>
                    </div>
                  )}
                  <div ref={mapContainerRef} className="leaflet-map" />
                </div>
              </div>
            </section>

            {/* ── Detail + Table ── */}
            <section className="detail-table-row">

              {/* Route sequence panel */}
              <div className="panel route-detail-panel">
                <div className="panel-header">
                  <div>
                    <h3 className="panel-title">Route Sequence</h3>
                    {selectedRow && (
                      <div className="panel-sub">
                        {selectedRow.Route} · {selectedRow.hop_count} hops · {normalizeDate(selectedRow.FlightDate)}
                        {selectedIdx !== 0 && <span className="selection-badge">Selected</span>}
                      </div>
                    )}
                  </div>
                </div>

                {selectedSegs.length > 0 ? (
                  <div className="table-scroll">
                    <table className="seq-table">
                      <thead>
                        <tr>
                          <th>#</th>
                          <th>Orig</th>
                          <th>Dest</th>
                          <th>Dep</th>
                          <th>Origin Airport</th>
                          <th>Destination Airport</th>
                        </tr>
                      </thead>
                      <tbody>
                        {selectedSegs.map((seg, i) => (
                          <tr key={i} className={i === 0 ? 'seq-first' : i === selectedSegs.length - 1 ? 'seq-last' : ''}>
                            <td className="td-num">{i + 1}</td>
                            <td className="td-code">{seg.OriginCode}</td>
                            <td className="td-code">{seg.DestCode}</td>
                            <td className="td-time">{formatDepTime(seg.DepTime)}</td>
                            <td className="td-name">{seg.OriginName || '—'}</td>
                            <td className="td-name">{seg.DestName || '—'}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                ) : selectedAirports.length > 0 ? (
                  <div className="seq-chips">
                    {selectedAirports.map((ap, i) => (
                      <span key={i}>
                        <span className={`airport-chip chip-${
                          i === 0 ? 'origin' : i === selectedAirports.length - 1 ? 'dest' : 'stop'
                        }`}>{ap}</span>
                        {i < selectedAirports.length - 1 && <span className="chip-arrow">→</span>}
                      </span>
                    ))}
                    {mapStatus === 'loading' && (
                      <p className="muted-note">Loading segment detail…</p>
                    )}
                  </div>
                ) : (
                  <p className="muted-note">Select an itinerary row to view the route detail.</p>
                )}
              </div>

              {/* Itinerary table */}
              <div className="panel itinerary-panel">
                <div className="panel-header">
                  <div>
                    <h3 className="panel-title">Top 10 Itineraries</h3>
                    <div className="panel-sub">Click a row to update the map and route detail</div>
                  </div>
                </div>
                <div className="table-scroll">
                  <table className="itinerary-table">
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
                    <tbody>
                      {rows.map((row, i) => (
                        <tr
                          key={i}
                          className={[
                            'itinerary-row',
                            i === 0       ? 'row-rank1'    : '',
                            i === selectedIdx ? 'row-selected' : '',
                          ].join(' ')}
                          onClick={() => handleRowClick(i)}
                          title={`Click to view route: ${row.Route}`}
                        >
                          <td className="td-rank">{i + 1}{i === 0 && <span className="rank-star">★</span>}</td>
                          <td>{row.Carrier}</td>
                          <td className="td-flight">{row.FlightNum}</td>
                          <td className="td-tail">{row.Tail_Number}</td>
                          <td className="td-date">{normalizeDate(row.FlightDate)}</td>
                          <td className="td-hops">{row.hop_count}</td>
                          <td className="td-route">{row.Route}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

            </section>

            {/* ── Query ledger ── */}
            <section className="ledger-section">
              <div className="panel">
                <div className="panel-header">
                  <h3 className="panel-title">Query Ledger</h3>
                  <div className="panel-sub">Click a row to expand SQL</div>
                </div>
                <div className="ledger-list">
                  {ledger.map(entry => (
                    <div key={entry.id} className="ledger-entry" onClick={() => toggleLedgerEntry(entry.id)}>
                      <div className="ledger-row">
                        <span className="ledger-toggle">{entry.expanded ? '▼' : '▶'}</span>
                        <span className="ledger-label">{entry.label}</span>
                        <span className="ledger-role">{entry.role}</span>
                        <span className={`ledger-status ls-${entry.status.toLowerCase()}`}>{entry.status}</span>
                        <span className="ledger-rows">{entry.rows != null ? `${entry.rows} rows` : '—'}</span>
                      </div>
                      {entry.expanded && (
                        <div className="ledger-sql">
                          <pre>{entry.sql}</pre>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </section>
          </>
        )}
      </main>

      {/* ── Footer controls ── */}
      <footer className="page-footer">
        <div className="footer-inner">
          <h4 className="footer-title">Data Controls</h4>
          <div className="footer-fields">
            <div className="field-group">
              <label className="field-label">JWE Token</label>
              <div className="token-row">
                <input
                  type="password"
                  className="token-input"
                  placeholder="Paste JWE token…"
                  value={jwe}
                  onChange={e => setJwe(e.target.value)}
                  autoComplete="off"
                />
                <button className="btn btn-ghost" onClick={handleForget} title="Clear stored token">
                  Forget
                </button>
              </div>
            </div>
            <div className="field-group">
              <label className="field-label">SQL Query</label>
              <textarea
                className="sql-textarea"
                value={sql}
                onChange={e => setSql(e.target.value)}
                rows={9}
                spellCheck={false}
                autoComplete="off"
              />
            </div>
            <div className="footer-actions">
              <button
                className="btn btn-primary"
                onClick={handleFetch}
                disabled={isRunning}
              >
                {isRunning ? 'Loading…' : 'Run Query'}
              </button>
              {statusMsg && <span className="status-msg">{statusMsg}</span>}
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
