```report
# {{question_title}}

{{data_overview_md}}

## Leadership Analysis

This report examines yearly carrier leadership based on completed flights. The data details how often each carrier appears as the annual leader across the full time range. We observe periods of stable dominance and note every true leadership transition, including the prior leader, new leader, and the resulting share swing. 

### Key Transitions

The transition analysis identifies the sharpest leadership transition by ranking rule, highlighting the most significant shifts in market share among top carriers. The table below outlines these exact points of transition.

{{result_table_md}}
```

```html
<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Yearly carrier leadership by completed flights</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
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
body {
    font-family: "Segoe UI", "Helvetica Neue", sans-serif;
    background: linear-gradient(180deg, var(--bg-top) 0%, var(--bg-bottom) 100%);
    color: var(--ink);
    margin: 0;
    padding: 0;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
}
h1, h2, h3 { font-family: Georgia, ui-serif, serif; color: var(--navy); }
.header { background: linear-gradient(90deg, var(--navy), var(--sky)); color: white; padding: 20px 40px; }
.header h1 { color: white; margin: 0; font-size: 24px; }
.header p { margin: 5px 0 0 0; opacity: 0.8; }
.container { max-width: 1280px; margin: 0 auto; padding: 20px; flex-grow: 1; width: 100%; box-sizing: border-box; }
.kpi-strip { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 20px; }
.kpi-card { background: var(--panel); border-radius: var(--radius-md); padding: 20px; box-shadow: var(--shadow); border: 1px solid var(--border); }
.kpi-label { font-size: 14px; color: var(--slate); margin-bottom: 5px; }
.kpi-value { font-size: 28px; font-weight: bold; color: var(--navy); }
.kpi-sub { font-size: 12px; color: var(--muted); margin-top: 5px; }
.dashboard-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }
@media (max-width: 900px) { .dashboard-grid { grid-template-columns: 1fr; } }
.card { background: var(--panel); border-radius: var(--radius-lg); padding: 20px; box-shadow: var(--shadow); border: 1px solid var(--border); }
.card h2 { margin-top: 0; font-size: 18px; margin-bottom: 15px; }
table { width: 100%; border-collapse: collapse; font-size: 14px; }
th, td { padding: 10px; text-align: left; border-bottom: 1px solid var(--border); }
th { background: var(--panel-alt); color: var(--navy); font-weight: bold; position: sticky; top: 0; }
tr:hover { background: var(--panel-alt); }
.sharpest { border-left: 4px solid var(--red); background: #fff5f3; }
.sharpest:hover { background: #ffebeb; }
.controls { background: var(--panel-alt); padding: 20px; border-top: 1px solid var(--border); margin-top: auto; }
.controls-inner { max-width: 1280px; margin: 0 auto; display: flex; flex-direction: column; gap: 10px; }
.control-row { display: flex; gap: 10px; align-items: center; }
input[type="password"] { padding: 8px; border: 1px solid var(--border); border-radius: var(--radius-sm); flex-grow: 1; max-width: 300px; }
button { padding: 8px 16px; background: var(--sky); color: white; border: none; border-radius: var(--radius-sm); cursor: pointer; }
button:hover { background: var(--navy); }
textarea { width: 100%; height: 100px; padding: 10px; border: 1px solid var(--border); border-radius: var(--radius-sm); font-family: monospace; resize: vertical; box-sizing: border-box; }
.ledger { margin-top: 20px; border-top: 1px solid var(--border); padding-top: 10px; }
.ledger-entry { margin-bottom: 5px; background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-sm); }
.ledger-header { display: flex; padding: 10px; cursor: pointer; background: var(--panel-alt); align-items: center; }
.ledger-header:hover { background: #eaf3f8; }
.ledger-toggle { width: 20px; font-family: monospace; }
.ledger-label { flex-grow: 1; font-weight: bold; }
.ledger-role { width: 100px; color: var(--muted); font-size: 12px; }
.ledger-status { width: 80px; font-weight: bold; }
.ledger-status.ok { color: var(--teal); }
.ledger-status.pending { color: var(--amber); }
.ledger-status.failed { color: var(--red); }
.ledger-rows { width: 60px; text-align: right; font-size: 12px; color: var(--muted); }
.ledger-sql { display: none; padding: 10px; background: var(--panel); font-family: monospace; font-size: 12px; white-space: pre-wrap; word-break: break-all; border-top: 1px solid var(--border); }
.ledger-entry.expanded .ledger-sql { display: block; }
.ledger-entry.expanded .ledger-toggle::before { content: "▼"; }
.ledger-entry:not(.expanded) .ledger-toggle::before { content: "▶"; }
.status-msg { margin-top: 10px; font-weight: bold; }
#dashboard-content { display: none; }
</style>
</head>
<body>

<div class="header">
    <h1>Yearly carrier leadership by completed flights</h1>
    <p>Analysis of market leadership and transitions over time</p>
</div>

<div class="container" id="loading-state">
    <div class="card">
        <h2>Awaiting Data</h2>
        <p>Please enter your JWE token below and fetch data to view the dashboard.</p>
    </div>
</div>

<div class="container" id="dashboard-content">
    <div class="kpi-strip">
        <div class="kpi-card">
            <div class="kpi-label">Total Years Analyzed</div>
            <div class="kpi-value" id="kpi-years">-</div>
        </div>
        <div class="kpi-card">
            <div class="kpi-label">Distinct Annual Leaders</div>
            <div class="kpi-value" id="kpi-leaders">-</div>
        </div>
        <div class="kpi-card">
            <div class="kpi-label">Largest Leader Share Gap</div>
            <div class="kpi-value" id="kpi-gap">-</div>
            <div class="kpi-sub" id="kpi-gap-sub"></div>
        </div>
        <div class="kpi-card">
            <div class="kpi-label">Sharpest Transition</div>
            <div class="kpi-value" id="kpi-transition">-</div>
            <div class="kpi-sub" id="kpi-transition-sub"></div>
        </div>
    </div>

    <div class="dashboard-grid">
        <div class="card" style="grid-column: 1 / -1;">
            <h2>Top Carriers Share (%)</h2>
            <canvas id="timeSeriesChart" height="80"></canvas>
        </div>
        <div class="card">
            <h2>Yearly Carrier Rank (Top 5)</h2>
            <canvas id="bumpChart"></canvas>
        </div>
        <div class="card">
            <h2>Leadership Transitions</h2>
            <div style="overflow-x: auto;">
                <table id="transitionsTable">
                    <thead>
                        <tr>
                            <th>Year</th>
                            <th>Prior Leader</th>
                            <th>New Leader</th>
                            <th>Share Swing (pts)</th>
                            <th>Share Gap (pts)</th>
                        </tr>
                    </thead>
                    <tbody></tbody>
                </table>
            </div>
        </div>
    </div>
</div>

<div class="controls">
    <div class="controls-inner">
        <h3>Data Connection</h3>
        <div class="control-row">
            <input type="password" id="jwe-token" placeholder="Enter JWE Token">
            <button id="btn-fetch">Fetch Data</button>
            <button id="btn-forget" style="background: var(--muted);">Forget Token</button>
        </div>
        <textarea id="sql-query">WITH YearlyCarrier AS (
    SELECT
        Year,
        Reporting_Airline,
        count() AS CompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year, Reporting_Airline
),
YearlyTotals AS (
    SELECT
        Year,
        sum(CompletedFlights) AS TotalCompletedFlights
    FROM YearlyCarrier
    GROUP BY Year
),
CarrierShares AS (
    SELECT
        c.Year,
        c.Reporting_Airline,
        c.CompletedFlights,
        (c.CompletedFlights * 100.0) / t.TotalCompletedFlights AS SharePct,
        row_number() OVER (PARTITION BY c.Year ORDER BY c.CompletedFlights DESC, c.Reporting_Airline ASC) AS RankInYear
    FROM YearlyCarrier c
    JOIN YearlyTotals t ON c.Year = t.Year
),
YearlyLeaders AS (
    SELECT
        Year,
        maxIf(Reporting_Airline, RankInYear = 1) AS LeaderReportingAirline,
        maxIf(Reporting_Airline, RankInYear = 2) AS RunnerUpReportingAirline,
        maxIf(SharePct, RankInYear = 1) - maxIf(SharePct, RankInYear = 2) AS LeaderShareGapPctPts,
        maxIf(SharePct, RankInYear = 1) AS LeaderSharePct
    FROM CarrierShares
    GROUP BY Year
),
LeaderTransitions AS (
    SELECT
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderShareGapPctPts,
        LeaderSharePct,
        lag(LeaderReportingAirline) OVER (ORDER BY Year) AS RawPriorLeader,
        lag(LeaderSharePct) OVER (ORDER BY Year) AS RawPriorShare,
        row_number() OVER (ORDER BY Year) AS YearNum
    FROM YearlyLeaders
),
ProcessedTransitions AS (
    SELECT
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderShareGapPctPts,
        CASE WHEN YearNum = 1 THEN CAST(NULL AS Nullable(String)) ELSE CAST(RawPriorLeader AS Nullable(String)) END AS PriorYearLeaderReportingAirline,
        CASE WHEN YearNum = 1 THEN 0 WHEN LeaderReportingAirline != RawPriorLeader THEN 1 ELSE 0 END AS LeaderChanged,
        CASE WHEN YearNum = 1 THEN CAST(NULL AS Nullable(Float64)) ELSE CAST(LeaderSharePct - RawPriorShare AS Nullable(Float64)) END AS LeaderShareChangePctPts
    FROM LeaderTransitions
),
FinalData AS (
    SELECT
        CAST('carrier_year' AS String) AS RowType,
        c.Year AS Year,
        CAST(c.Reporting_Airline AS String) AS Reporting_Airline,
        CAST(c.RankInYear AS UInt32) AS RankInYear,
        CAST(c.CompletedFlights AS UInt64) AS CompletedFlights,
        CAST(c.SharePct AS Float64) AS SharePct,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.LeaderReportingAirline AS String) ELSE CAST('' AS String) END AS LeaderReportingAirline,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.RunnerUpReportingAirline AS String) ELSE CAST('' AS String) END AS RunnerUpReportingAirline,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.LeaderShareGapPctPts AS Nullable(Float64)) ELSE CAST(NULL AS Nullable(Float64)) END AS LeaderShareGapPctPts,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.PriorYearLeaderReportingAirline AS Nullable(String)) ELSE CAST(NULL AS Nullable(String)) END AS PriorYearLeaderReportingAirline,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.LeaderChanged AS UInt8) ELSE CAST(0 AS UInt8) END AS LeaderChanged,
        CASE WHEN c.RankInYear = 1 THEN CAST(t.LeaderShareChangePctPts AS Nullable(Float64)) ELSE CAST(NULL AS Nullable(Float64)) END AS LeaderShareChangePctPts
    FROM CarrierShares c
    JOIN ProcessedTransitions t ON c.Year = t.Year
    WHERE c.RankInYear <= 5
),
SummaryData AS (
    SELECT
        CAST('year_summary' AS String) AS RowType,
        t.Year AS Year,
        CAST('' AS String) AS Reporting_Airline,
        CAST(0 AS UInt32) AS RankInYear,
        CAST(tot.TotalCompletedFlights AS UInt64) AS CompletedFlights,
        CAST(100.0 AS Float64) AS SharePct,
        CAST(t.LeaderReportingAirline AS String) AS LeaderReportingAirline,
        CAST(t.RunnerUpReportingAirline AS String) AS RunnerUpReportingAirline,
        CAST(t.LeaderShareGapPctPts AS Nullable(Float64)) AS LeaderShareGapPctPts,
        CAST(t.PriorYearLeaderReportingAirline AS Nullable(String)) AS PriorYearLeaderReportingAirline,
        CAST(t.LeaderChanged AS UInt8) AS LeaderChanged,
        CAST(t.LeaderShareChangePctPts AS Nullable(Float64)) AS LeaderShareChangePctPts
    FROM ProcessedTransitions t
    JOIN YearlyTotals tot ON t.Year = tot.Year
)
SELECT * FROM FinalData
UNION ALL
SELECT * FROM SummaryData
ORDER BY Year ASC, RowType ASC, RankInYear ASC, Reporting_Airline ASC</textarea>
        <div id="status-text" class="status-msg"></div>

        <div class="ledger" id="query-ledger">
            <h3>Query Ledger</h3>
        </div>
    </div>
</div>

<script>
const AUTH_KEY = 'OnTimeAnalystDashboard::auth::jwe';
let timeSeriesChartInstance = null;
let bumpChartInstance = null;

document.addEventListener('DOMContentLoaded', () => {
    const storedToken = localStorage.getItem(AUTH_KEY);
    if (storedToken) {
        document.getElementById('jwe-token').value = storedToken;
    }

    document.getElementById('btn-fetch').addEventListener('click', fetchData);
    document.getElementById('btn-forget').addEventListener('click', () => {
        localStorage.removeItem(AUTH_KEY);
        document.getElementById('jwe-token').value = '';
        document.getElementById('status-text').innerText = 'Token forgotten.';
        document.getElementById('status-text').style.color = 'var(--slate)';
    });
});

function addLedgerEntry(id, label, role, sql) {
    const ledger = document.getElementById('query-ledger');
    const entry = document.createElement('div');
    entry.className = 'ledger-entry';
    entry.id = `ledger-${id}`;
    entry.innerHTML = `
        <div class="ledger-header" onclick="this.parentElement.classList.toggle('expanded')">
            <div class="ledger-toggle"></div>
            <div class="ledger-label">${label}</div>
            <div class="ledger-role">${role}</div>
            <div class="ledger-status pending" id="ledger-status-${id}">Pending</div>
            <div class="ledger-rows" id="ledger-rows-${id}">-</div>
        </div>
        <div class="ledger-sql">${sql.replace(/</g, '&lt;').replace(/>/g, '&gt;')}</div>
    `;
    ledger.appendChild(entry);
}

function updateLedgerEntry(id, status, rows) {
    const statusEl = document.getElementById(`ledger-status-${id}`);
    const rowsEl = document.getElementById(`ledger-rows-${id}`);
    if (statusEl) {
        statusEl.className = `ledger-status ${status}`;
        statusEl.innerText = status.toUpperCase();
    }
    if (rowsEl) {
        rowsEl.innerText = rows !== null ? `${rows} rows` : '-';
    }
}

async function fetchData() {
    const token = document.getElementById('jwe-token').value.trim();
    const sql = document.getElementById('sql-query').value.trim();
    const statusText = document.getElementById('status-text');

    if (!token || !sql) {
        statusText.innerText = 'Error: Token and SQL are required.';
        statusText.style.color = 'var(--red)';
        return;
    }

    localStorage.setItem(AUTH_KEY, token);
    statusText.innerText = 'Fetching data...';
    statusText.style.color = 'var(--amber)';
    
    document.getElementById('query-ledger').innerHTML = '<h3>Query Ledger</h3>';
    addLedgerEntry('primary', 'Primary Analytical Query', 'source', sql);

    try {
        const url = `https://mcp.demo.altinity.cloud/${encodeURIComponent(token)}/openapi/execute_query?query=${encodeURIComponent(sql)}`;
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const payload = await response.json();
        
        if (payload.count === 0 && !payload.rows) {
            payload.rows = [];
        }

        if (!payload.columns || !Array.isArray(payload.rows)) {
            throw new Error('Malformed payload: columns or rows missing.');
        }

        const data = payload.rows.map(row => {
            let obj = {};
            payload.columns.forEach((col, i) => {
                obj[col.name] = row[i];
            });
            return obj;
        });

        updateLedgerEntry('primary', 'ok', payload.count);
        statusText.innerText = 'Data loaded successfully.';
        statusText.style.color = 'var(--teal)';

        renderDashboard(data);
    } catch (err) {
        updateLedgerEntry('primary', 'failed', null);
        statusText.innerText = `Error: ${err.message}`;
        statusText.style.color = 'var(--red)';
        console.error(err);
    }
}

function renderDashboard(data) {
    document.getElementById('loading-state').style.display = 'none';
    document.getElementById('dashboard-content').style.display = 'block';

    const summaryData = data.filter(d => d.RowType === 'year_summary');
    const carrierData = data.filter(d => d.RowType === 'carrier_year');

    // KPIs
    const years = summaryData.length;
    const leaders = new Set(summaryData.map(d => d.LeaderReportingAirline).filter(Boolean)).size;
    
    let maxGap = 0;
    let maxGapYear = null;
    let sharpestTransitionVal = 0;
    let sharpestTransitionYear = null;
    let sharpestTransitionNode = null;

    summaryData.forEach(d => {
        if (d.LeaderShareGapPctPts > maxGap) {
            maxGap = d.LeaderShareGapPctPts;
            maxGapYear = d.Year;
        }
        if (d.LeaderChanged == 1 && Math.abs(d.LeaderShareChangePctPts) > Math.abs(sharpestTransitionVal)) {
            sharpestTransitionVal = d.LeaderShareChangePctPts;
            sharpestTransitionYear = d.Year;
            sharpestTransitionNode = d;
        }
    });

    document.getElementById('kpi-years').innerText = years;
    document.getElementById('kpi-leaders').innerText = leaders;
    document.getElementById('kpi-gap').innerText = maxGap > 0 ? maxGap.toFixed(1) + ' pts' : '-';
    document.getElementById('kpi-gap-sub').innerText = maxGapYear ? `in ${maxGapYear}` : '';
    
    if (sharpestTransitionNode) {
        const sign = sharpestTransitionNode.LeaderShareChangePctPts > 0 ? '+' : '';
        document.getElementById('kpi-transition').innerText = `${sign}${sharpestTransitionNode.LeaderShareChangePctPts.toFixed(1)} pts`;
        document.getElementById('kpi-transition-sub').innerText = `${sharpestTransitionNode.PriorYearLeaderReportingAirline} \u2192 ${sharpestTransitionNode.LeaderReportingAirline} (${sharpestTransitionYear})`;
    } else {
        document.getElementById('kpi-transition').innerText = 'None';
        document.getElementById('kpi-transition-sub').innerText = '';
    }

    // Chart Colors
    const colors = [
        '#1f8a70', '#3c88b5', '#c54f36', '#d48a1f', '#0e3a52', 
        '#5c7080', '#8b5a2b', '#483d8b', '#2e8b57', '#b22222'
    ];
    const carrierColors = {};
    let colorIdx = 0;
    const allCarriers = [...new Set(carrierData.map(d => d.Reporting_Airline))];
    allCarriers.forEach(c => {
        carrierColors[c] = colors[colorIdx % colors.length];
        colorIdx++;
    });

    const yearsLabels = summaryData.map(d => d.Year);

    // Time Series Chart
    const timeSeriesDatasets = [];
    allCarriers.forEach(carrier => {
        const carrierPts = [];
        yearsLabels.forEach(year => {
            const row = carrierData.find(d => d.Year === year && d.Reporting_Airline === carrier);
            carrierPts.push(row ? row.SharePct : null);
        });
        if (carrierPts.some(v => v !== null)) {
            timeSeriesDatasets.push({
                label: carrier,
                data: carrierPts,
                borderColor: carrierColors[carrier],
                backgroundColor: carrierColors[carrier],
                tension: 0.3,
                fill: false,
                spanGaps: true
            });
        }
    });

    const ctxTime = document.getElementById('timeSeriesChart').getContext('2d');
    if (timeSeriesChartInstance) timeSeriesChartInstance.destroy();
    timeSeriesChartInstance = new Chart(ctxTime, {
        type: 'line',
        data: { labels: yearsLabels, datasets: timeSeriesDatasets },
        options: {
            responsive: true,
            plugins: { legend: { position: 'right' } },
            scales: {
                y: { title: { display: true, text: 'Share of Flights (%)' } },
                x: { title: { display: true, text: 'Year' } }
            }
        }
    });

    // Bump Chart
    const bumpDatasets = [];
    allCarriers.forEach(carrier => {
        const rankPts = [];
        yearsLabels.forEach(year => {
            const row = carrierData.find(d => d.Year === year && d.Reporting_Airline === carrier);
            rankPts.push(row ? row.RankInYear : null);
        });
        if (rankPts.some(v => v !== null)) {
            bumpDatasets.push({
                label: carrier,
                data: rankPts,
                borderColor: carrierColors[carrier],
                backgroundColor: carrierColors[carrier],
                tension: 0.1,
                fill: false,
                spanGaps: false, // Break line when outside top 5
                borderWidth: 3,
                pointRadius: 4
            });
        }
    });

    const ctxBump = document.getElementById('bumpChart').getContext('2d');
    if (bumpChartInstance) bumpChartInstance.destroy();
    bumpChartInstance = new Chart(ctxBump, {
        type: 'line',
        data: { labels: yearsLabels, datasets: bumpDatasets },
        options: {
            responsive: true,
            plugins: { legend: { display: false }, tooltip: { mode: 'index', intersect: false } },
            scales: {
                y: { 
                    reverse: true, 
                    min: 1, 
                    max: 5, 
                    ticks: { stepSize: 1 },
                    title: { display: true, text: 'Rank' } 
                },
                x: { title: { display: true, text: 'Year' } }
            }
        }
    });

    // Transitions Table
    const tbody = document.querySelector('#transitionsTable tbody');
    tbody.innerHTML = '';
    const transitions = summaryData.filter(d => d.LeaderChanged == 1 && d.PriorYearLeaderReportingAirline);
    
    transitions.forEach(t => {
        const isSharpest = t.Year === sharpestTransitionYear;
        const tr = document.createElement('tr');
        if (isSharpest) tr.className = 'sharpest';
        
        tr.innerHTML = `
            <td>${t.Year}${isSharpest ? ' <strong>(Sharpest)</strong>' : ''}</td>
            <td>${t.PriorYearLeaderReportingAirline || '-'}</td>
            <td><strong>${t.LeaderReportingAirline}</strong></td>
            <td style="color: ${t.LeaderShareChangePctPts > 0 ? 'var(--teal)' : 'var(--red)'}">
                ${t.LeaderShareChangePctPts > 0 ? '+' : ''}${t.LeaderShareChangePctPts ? t.LeaderShareChangePctPts.toFixed(2) : '-'}
            </td>
            <td>${t.LeaderShareGapPctPts ? t.LeaderShareGapPctPts.toFixed(2) : '-'}</td>
        `;
        tbody.appendChild(tr);
    });

    if (transitions.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; color: var(--muted);">No true leadership transitions found.</td></tr>';
    }
}
</script>
</body>
</html>
```
