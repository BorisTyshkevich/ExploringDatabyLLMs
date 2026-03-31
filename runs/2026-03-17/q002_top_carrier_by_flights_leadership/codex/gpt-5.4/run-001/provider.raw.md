```report
# {{question_title}}

Generated from `result.json` at {{generated_at}} with {{row_count}} rows across `{{columns_csv}}`.

## Analytical focus
{{data_overview_md}}

Use the annual leader rows where `RankInYear = 1` to assess four questions: how often each carrier leads across the full span, which years represent true leadership changes (`LeaderChanged = 1` with a non-null prior leader), which of those transitions has the largest absolute `LeaderShareChangePctPts`, and whether repeated runs by the same `LeaderReportingAirline` indicate extended dominance.

## Result detail
{{result_table_md}}
```

```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Yearly carrier leadership by completed flights</title>
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
      font-family: "Segoe UI", "Helvetica Neue", sans-serif;
      color: var(--ink);
      background:
        radial-gradient(circle at top left, rgba(60,136,181,0.18), transparent 30%),
        radial-gradient(circle at top right, rgba(31,138,112,0.12), transparent 28%),
        linear-gradient(180deg, var(--bg-top), var(--bg-bottom));
    }
    .wrap {
      max-width: 1280px;
      margin: 0 auto;
      padding: 28px 20px 36px;
    }
    .hero {
      background: linear-gradient(140deg, var(--navy), #174f6e 58%, #2b7e98);
      color: #fff;
      border-radius: var(--radius-xl);
      padding: 28px 28px 22px;
      box-shadow: var(--shadow);
      position: relative;
      overflow: hidden;
    }
    .hero::after {
      content: "";
      position: absolute;
      inset: auto -40px -65px auto;
      width: 240px;
      height: 240px;
      background: radial-gradient(circle, rgba(255,255,255,0.18), transparent 62%);
      pointer-events: none;
    }
    .eyebrow {
      text-transform: uppercase;
      letter-spacing: 0.14em;
      font-size: 12px;
      opacity: 0.78;
      margin-bottom: 10px;
    }
    h1, h2, h3 {
      font-family: Georgia, ui-serif, serif;
      margin: 0;
      color: inherit;
    }
    h1 {
      font-size: clamp(2rem, 4vw, 3.1rem);
      line-height: 1.05;
      max-width: 760px;
    }
    .subtitle {
      margin-top: 14px;
      max-width: 860px;
      line-height: 1.55;
      color: rgba(255,255,255,0.9);
      font-size: 1rem;
    }
    .kpi-grid, .two-col, .single-col {
      display: grid;
      gap: 18px;
      margin-top: 22px;
    }
    .kpi-grid {
      grid-template-columns: repeat(auto-fit, minmax(190px, 1fr));
    }
    .two-col {
      grid-template-columns: repeat(auto-fit, minmax(330px, 1fr));
    }
    .single-col { grid-template-columns: 1fr; }
    .card {
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: var(--radius-lg);
      box-shadow: var(--shadow);
      padding: 18px 18px 16px;
    }
    .card h2 {
      font-size: 1.35rem;
      margin-bottom: 8px;
      color: var(--navy);
    }
    .card p {
      margin: 0;
      line-height: 1.5;
      color: var(--muted);
    }
    .kpi {
      padding: 16px 16px 14px;
      background: linear-gradient(180deg, #fff, #f8fbfd);
    }
    .kpi .label {
      color: var(--muted);
      font-size: 0.85rem;
      text-transform: uppercase;
      letter-spacing: 0.08em;
      margin-bottom: 10px;
    }
    .kpi .value {
      font-size: clamp(1.6rem, 3vw, 2.35rem);
      line-height: 1;
      color: var(--navy);
      font-weight: 700;
    }
    .kpi .meta {
      margin-top: 8px;
      color: var(--muted);
      font-size: 0.92rem;
      line-height: 1.35;
    }
    .leader-list {
      display: grid;
      gap: 10px;
      margin-top: 14px;
    }
    .leader-row {
      display: grid;
      grid-template-columns: 72px 1fr 84px;
      gap: 12px;
      align-items: center;
      background: var(--panel-alt);
      border-radius: var(--radius-md);
      padding: 10px 12px;
    }
    .leader-code {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      font-weight: 700;
      font-size: 1rem;
      color: #fff;
      border-radius: 999px;
      padding: 8px 0;
    }
    .leader-row strong { color: var(--ink); }
    .leader-row span { color: var(--muted); font-size: 0.9rem; }
    .leader-years {
      text-align: right;
      font-weight: 700;
      color: var(--navy);
    }
    .callout-strip {
      display: grid;
      gap: 12px;
      margin-top: 14px;
    }
    .callout {
      border-left: 4px solid var(--sky);
      background: linear-gradient(180deg, #fff, #f7fbfd);
      border-radius: var(--radius-md);
      padding: 12px 14px;
    }
    .callout.sharpest {
      border-left-color: var(--red);
      background: linear-gradient(180deg, rgba(197,79,54,0.08), rgba(255,255,255,0.96));
    }
    .callout .title {
      font-weight: 700;
      color: var(--navy);
      margin-bottom: 4px;
    }
    .callout .desc {
      color: var(--muted);
      line-height: 1.45;
      font-size: 0.95rem;
    }
    .legend {
      display: flex;
      flex-wrap: wrap;
      gap: 10px 14px;
      margin-top: 14px;
    }
    .legend-item {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      color: var(--muted);
      font-size: 0.9rem;
      padding: 4px 8px;
      background: var(--panel-alt);
      border-radius: 999px;
    }
    .swatch {
      width: 10px;
      height: 10px;
      border-radius: 999px;
      flex: 0 0 auto;
    }
    .chart-wrap {
      margin-top: 14px;
      background: linear-gradient(180deg, #fcfeff, #f4f8fb);
      border: 1px solid var(--border);
      border-radius: var(--radius-md);
      padding: 10px;
      overflow-x: auto;
    }
    svg {
      display: block;
      width: 100%;
      height: auto;
      min-width: 760px;
    }
    .axis text, .axis-label, .grid-label {
      fill: var(--muted);
      font-size: 12px;
      font-family: "Segoe UI", "Helvetica Neue", sans-serif;
    }
    .grid-line {
      stroke: var(--grid);
      stroke-width: 1;
    }
    .transition-note {
      margin-top: 12px;
      color: var(--muted);
      font-size: 0.92rem;
      line-height: 1.45;
    }
    .table-tools {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 12px;
      margin-top: 12px;
      color: var(--muted);
      font-size: 0.92rem;
    }
    button {
      border: 0;
      border-radius: 999px;
      background: var(--navy);
      color: #fff;
      padding: 9px 14px;
      cursor: pointer;
      font-weight: 600;
    }
    button.secondary {
      background: var(--panel-alt);
      color: var(--navy);
      border: 1px solid var(--border);
    }
    button:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 12px;
      font-size: 0.95rem;
    }
    thead th {
      position: sticky;
      top: 0;
      background: #eef5f9;
      color: var(--navy);
      text-align: left;
      padding: 10px 12px;
      border-bottom: 1px solid var(--border);
      font-size: 0.83rem;
      text-transform: uppercase;
      letter-spacing: 0.06em;
    }
    tbody td {
      padding: 11px 12px;
      border-bottom: 1px solid var(--border);
      color: var(--ink);
    }
    tbody tr:hover { background: var(--panel-alt); }
    tbody tr.sharpest { background: rgba(197,79,54,0.08); }
    .num { text-align: right; font-variant-numeric: tabular-nums; }
    .ledger {
      margin-top: 22px;
    }
    .ledger-entry {
      border: 1px solid var(--border);
      border-radius: var(--radius-md);
      overflow: hidden;
      background: var(--panel);
    }
    .ledger-head {
      width: 100%;
      display: grid;
      grid-template-columns: 1.4em 1fr 7em 6em 5em;
      gap: 10px;
      align-items: center;
      text-align: left;
      border-radius: 0;
      background: transparent;
      color: var(--ink);
      padding: 12px 14px;
    }
    .ledger-head:hover { background: var(--panel-alt); }
    .toggle-icon {
      font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
      color: var(--muted);
    }
    .ledger-role, .ledger-rows, .ledger-status {
      font-size: 0.85rem;
      color: var(--muted);
      text-transform: uppercase;
      letter-spacing: 0.06em;
    }
    .ledger-status.ok { color: var(--teal); font-weight: 700; }
    .ledger-sql {
      display: none;
      background: var(--panel-alt);
      border-top: 1px solid var(--border);
      padding: 14px;
    }
    .ledger-sql.open { display: block; }
    pre {
      margin: 0;
      white-space: pre-wrap;
      word-break: break-word;
      font-size: 0.84rem;
      line-height: 1.45;
      color: var(--ink);
      font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    }
    .footer-controls {
      margin-top: 26px;
      padding: 18px;
      background: var(--panel);
      border-radius: var(--radius-lg);
      border: 1px solid var(--border);
      box-shadow: var(--shadow);
    }
    .footer-grid {
      display: grid;
      grid-template-columns: 1fr auto auto;
      gap: 12px;
      align-items: end;
    }
    .field {
      display: grid;
      gap: 6px;
    }
    label {
      font-size: 0.83rem;
      color: var(--muted);
      text-transform: uppercase;
      letter-spacing: 0.06em;
    }
    input[type="password"], textarea {
      width: 100%;
      border: 1px solid var(--border);
      border-radius: var(--radius-sm);
      background: #fbfdfe;
      color: var(--ink);
      padding: 10px 12px;
      font: inherit;
    }
    textarea {
      min-height: 230px;
      margin-top: 12px;
      resize: vertical;
      font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
      font-size: 0.84rem;
      line-height: 1.45;
    }
    .status {
      margin-top: 10px;
      color: var(--muted);
      font-size: 0.92rem;
      line-height: 1.45;
    }
    .hint {
      margin-top: 8px;
      color: var(--muted);
      font-size: 0.88rem;
    }
    @media (max-width: 760px) {
      .footer-grid { grid-template-columns: 1fr; }
      .ledger-head { grid-template-columns: 1.4em 1fr; }
      .ledger-role, .ledger-rows, .ledger-status { display: none; }
      .leader-row { grid-template-columns: 60px 1fr 72px; }
    }
  </style>
</head>
<body>
  <div class="wrap">
    <section class="hero">
      <div class="eyebrow">OnTime Leadership Dashboard</div>
      <h1>Yearly carrier leadership by completed flights</h1>
      <div class="subtitle" id="subtitle">A static qforge artifact built from the saved leadership query. The dashboard isolates leader frequency, true handoffs, rank movement, and share dynamics without treating the first observed year as a transition.</div>
    </section>

    <section class="kpi-grid" id="kpis"></section>

    <section class="two-col">
      <article class="card">
        <h2>Leader Frequency Across the Full Span</h2>
        <p id="leaderSummary"></p>
        <div class="leader-list" id="leaderList"></div>
      </article>

      <article class="card">
        <h2>Leadership Transitions</h2>
        <p id="transitionSummary"></p>
        <div class="callout-strip" id="transitionCallouts"></div>
      </article>
    </section>

    <section class="single-col">
      <article class="card">
        <h2>Yearly Carrier Rank Bump Chart</h2>
        <p>The chart tracks rank positions for every carrier that appears in the fetched top-five set. Rank 1 is at the top; the sharpest leadership handoff is marked in red.</p>
        <div class="legend" id="bumpLegend"></div>
        <div class="chart-wrap"><svg id="bumpChart" viewBox="0 0 1100 420" aria-label="Bump chart of annual carrier ranks"></svg></div>
      </article>
    </section>

    <section class="single-col">
      <article class="card">
        <h2>Completed-Flight Share for Annual Leaders</h2>
        <p>The share trend focuses on carriers that actually hold the annual leadership position at least once in the series.</p>
        <div class="legend" id="shareLegend"></div>
        <div class="chart-wrap"><svg id="shareChart" viewBox="0 0 1100 420" aria-label="Time series of completed-flight share for annual leaders"></svg></div>
        <div class="transition-note" id="stabilityNote"></div>
      </article>
    </section>

    <section class="single-col">
      <article class="card">
        <h2>True Leadership-Change Years</h2>
        <p>Only years where the annual leader differs from the previous year are included. The first observed year is excluded by design.</p>
        <div class="table-tools">
          <span id="tableMeta"></span>
          <button id="exportTransitions" class="secondary">Export Transition CSV</button>
        </div>
        <div style="overflow:auto;">
          <table>
            <thead>
              <tr>
                <th>Year</th>
                <th>Prior leader</th>
                <th>New leader</th>
                <th class="num">Share swing</th>
                <th class="num">Leader gap</th>
              </tr>
            </thead>
            <tbody id="transitionTable"></tbody>
          </table>
        </div>
      </article>
    </section>

    <section class="card ledger" id="query-ledger">
      <h2>Query Ledger</h2>
      <p>This artifact is static in the browser. The qforge harness materialized the saved SQL into the result set used below.</p>
      <div class="ledger-entry">
        <button class="ledger-head" type="button" onclick="toggleLedgerEntry('ledger-primary')">
          <span class="toggle-icon" id="ledger-primary-toggle">▶</span>
          <span>Primary analytical query for yearly carrier leadership</span>
          <span class="ledger-role">primary</span>
          <span class="ledger-status ok">loaded</span>
          <span class="ledger-rows" id="ledgerRows">195 rows</span>
        </button>
        <div class="ledger-sql" id="ledger-primary">
<pre id="ledgerSql"></pre>
        </div>
      </div>
    </section>

    <section class="footer-controls">
      <div class="footer-grid">
        <div class="field">
          <label for="tokenInput">JWE Token</label>
          <input id="tokenInput" type="password" placeholder="Stored locally only">
        </div>
        <button id="fetchBtn" type="button">Save Token</button>
        <button id="forgetBtn" type="button" class="secondary">Forget</button>
      </div>
      <div class="field">
        <label for="sqlTextarea">Saved SQL</label>
        <textarea id="sqlTextarea" spellcheck="false"></textarea>
      </div>
      <div class="status" id="statusText">Static artifact sourced from the precomputed saved-query result. Footer controls remain available for token storage consistency across dashboards.</div>
      <div class="hint">Green signals improvement or clear leadership, amber denotes caution, and red marks the sharpest year-over-year leadership break.</div>
    </section>
  </div>

  <script>
    const dashboardId = "q002-yearly-carrier-leadership";
    const authKey = "OnTimeAnalystDashboard::auth::jwe";
    const savedSQL = `WITH annual_totals AS (
    SELECT
        Year,
        count() AS YearCompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY Year
),
annual_carrier_counts AS (
    SELECT
        Year,
        Reporting_Airline,
        count() AS CompletedFlights
    FROM default.ontime_v2
    WHERE Cancelled = 0
    GROUP BY
        Year,
        Reporting_Airline
),
ranked_carriers AS (
    SELECT
        acc.Year,
        acc.Reporting_Airline,
        row_number() OVER (
            PARTITION BY acc.Year
            ORDER BY acc.CompletedFlights DESC, acc.Reporting_Airline ASC
        ) AS RankInYear,
        acc.CompletedFlights,
        100.0 * acc.CompletedFlights / at.YearCompletedFlights AS SharePct
    FROM annual_carrier_counts AS acc
    INNER JOIN annual_totals AS at
        ON acc.Year = at.Year
),
leader_runner_up AS (
    SELECT
        Year,
        maxIf(Reporting_Airline, RankInYear = 1) AS LeaderReportingAirline,
        maxIf(Reporting_Airline, RankInYear = 2) AS RunnerUpReportingAirline,
        maxIf(SharePct, RankInYear = 1) AS LeaderSharePct,
        maxIf(SharePct, RankInYear = 2) AS RunnerUpSharePct
    FROM ranked_carriers
    WHERE RankInYear <= 2
    GROUP BY Year
),
leader_transitions AS (
    SELECT
        Year,
        LeaderReportingAirline,
        RunnerUpReportingAirline,
        LeaderSharePct - RunnerUpSharePct AS LeaderShareGapPctPts,
        prior_year_leader AS PriorYearLeaderReportingAirline,
        if(prior_year_leader IS NULL, 0, LeaderReportingAirline != prior_year_leader) AS LeaderChanged,
        if(
            prior_year_leader IS NULL,
            cast(NULL, 'Nullable(Float64)'),
            LeaderSharePct - prior_year_share_pct
        ) AS LeaderShareChangePctPts
    FROM (
        SELECT
            Year,
            LeaderReportingAirline,
            RunnerUpReportingAirline,
            LeaderSharePct,
            RunnerUpSharePct,
            lagInFrame(toNullable(LeaderReportingAirline)) OVER w AS prior_year_leader,
            lagInFrame(toNullable(LeaderSharePct)) OVER w AS prior_year_share_pct
        FROM leader_runner_up
        WINDOW w AS (ORDER BY Year ASC ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)
    )
)
SELECT
    'carrier_year' AS RowType,
    rc.Year,
    rc.Reporting_Airline,
    rc.RankInYear,
    rc.CompletedFlights,
    rc.SharePct,
    if(rc.RankInYear = 1, lt.LeaderReportingAirline, cast(NULL, 'Nullable(String)')) AS LeaderReportingAirline,
    if(rc.RankInYear = 1, lt.RunnerUpReportingAirline, cast(NULL, 'Nullable(String)')) AS RunnerUpReportingAirline,
    if(rc.RankInYear = 1, lt.LeaderShareGapPctPts, cast(NULL, 'Nullable(Float64)')) AS LeaderShareGapPctPts,
    if(rc.RankInYear = 1, lt.PriorYearLeaderReportingAirline, cast(NULL, 'Nullable(String)')) AS PriorYearLeaderReportingAirline,
    if(rc.RankInYear = 1, toUInt8(lt.LeaderChanged), cast(NULL, 'Nullable(UInt8)')) AS LeaderChanged,
    if(rc.RankInYear = 1, lt.LeaderShareChangePctPts, cast(NULL, 'Nullable(Float64)')) AS LeaderShareChangePctPts
FROM ranked_carriers AS rc
INNER JOIN leader_transitions AS lt
    ON rc.Year = lt.Year
WHERE rc.RankInYear <= 5
ORDER BY
    rc.Year ASC,
    RowType,
    rc.RankInYear ASC,
    rc.Reporting_Airline ASC`;

    const result = {
      columns: ["RowType","Year","Reporting_Airline","RankInYear","CompletedFlights","SharePct","LeaderReportingAirline","RunnerUpReportingAirline","LeaderShareGapPctPts","PriorYearLeaderReportingAirline","LeaderChanged","LeaderShareChangePctPts"],
      rows: [["carrier_year",1987,"DL",1,183756,14.221048631689575,"DL","AA",1.5935567403247788,null,0,null],["carrier_year",1987,"AA",2,163165,12.627491891364796,null,null,null,null,null,null],["carrier_year",1987,"UA",3,149188,11.545798794404016,null,null,null,null,null,null],["carrier_year",1987,"CO",4,120494,9.325143308663684,null,null,null,null,null,null],["carrier_year",1987,"PI",5,115253,8.919537418903975,null,null,null,null,null,null],["carrier_year",1988,"DL",1,749514,14.548209380828515,"DL","AA",1.1911645590111508,"DL",0,0.32716074913894033],["carrier_year",1988,"AA",2,688146,13.357044821817365,null,null,null,null,null,null],["carrier_year",1988,"UA",3,581378,11.284657622682593,null,null,null,null,null,null],["carrier_year",1988,"US",4,489821,9.507518828369857,null,null,null,null,null,null],["carrier_year",1988,"PI",5,466643,9.057629437339344,null,null,null,null,null,null],["carrier_year",1989,"DL",1,778612,15.675589159327446,"DL","AA",1.3152514528284982,"DL",0,1.1273797784989306],["carrier_year",1989,"AA",2,713283,14.360337706498948,null,null,null,null,null,null],["carrier_year",1989,"US",3,699936,14.091626090816755,null,null,null,null,null,null],["carrier_year",1989,"UA",4,565762,11.390336488468472,null,null,null,null,null,null],["carrier_year",1989,"NW",5,446735,8.993997425023178,null,null,null,null,null,null],["carrier_year",1990,"US",1,991989,19.00931984397621,"US","DL",3.3375715132985277,"DL",1,3.3337306846487653],["carrier_year",1990,"DL",2,817820,15.671748330677683,null,null,null,null,null,null],["carrier_year",1990,"AA",3,705296,13.515469676253513,null,null,null,null,null,null],["carrier_year",1990,"UA",4,598819,11.475068674803845,null,null,null,null,null,null],["carrier_year",1990,"NW",5,456674,8.751167735154313,null,null,null,null,null,null],["carrier_year",1991,"US",1,900252,17.885493362365946,"US","DL",0.6248435457402728,"US",0,-1.1238264816102657],["carrier_year",1991,"DL",2,868801,17.260649816625673,null,null,null,null,null,null],["carrier_year",1991,"AA",3,720240,14.309157590664002,null,null,null,null,null,null],["carrier_year",1991,"UA",4,621361,12.344707971915716,null,null,null,null,null,null],["carrier_year",1991,"NW",5,457489,9.089028930627684,null,null,null,null,null,null],["carrier_year",1992,"DL",1,910094,18.05985369854391,"DL","US",0.8488445169498036,"US",1,0.17436033617796554],["carrier_year",1992,"US",2,867318,17.211009181594108,null,null,null,null,null,null],["carrier_year",1992,"AA",3,777340,15.425490854819529,null,null,null,null,null,null],["carrier_year",1992,"UA",4,631353,12.52853310991699,null,null,null,null,null,null],["carrier_year",1992,"NW",5,471001,9.346517120064389,null,null,null,null,null,null],["carrier_year",1993,"DL",1,890515,17.772423411225994,"DL","US",1.4358598953909407,"DL",0,-0.28743028731791753],["carrier_year",1993,"US",2,818569,16.336563515835053,null,null,null,null,null,null],["carrier_year",1993,"AA",3,777427,15.515473423040815,null,null,null,null,null,null],["carrier_year",1993,"UA",4,642173,12.816146229156422,null,null,null,null,null,null],["carrier_year",1993,"WN",5,486067,9.70066594074708,null,null,null,null,null,null],["carrier_year",1994,"DL",1,864990,16.916446261402598,"DL","US",0.5230078062968246,"DL",0,-0.8559771498233957],["carrier_year",1994,"US",2,838247,16.393438455105773,null,null,null,null,null,null],["carrier_year",1994,"AA",3,714082,13.96516697214406,null,null,null,null,null,null],["carrier_year",1994,"UA",4,630722,12.334911176874149,null,null,null,null,null,null],["carrier_year",1994,"WN",5,562961,11.009722082065075,null,null,null,null,null,null],["carrier_year",1995,"DL",1,871929,16.65407322658833,"DL","US",2.043575340032433,"DL",0,-0.2623730348142672],["carrier_year",1995,"US",2,764937,14.610497886555898,null,null,null,null,null,null],["carrier_year",1995,"UA",3,711481,13.589474227060107,null,null,null,null,null,null],["carrier_year",1995,"WN",4,684961,13.08293525201842,null,null,null,null,null,null],["carrier_year",1995,"AA",5,677176,12.934239704480731,null,null,null,null,null,null],["carrier_year",1996,"DL",1,871187,16.678392639955952,"DL","WN",2.3591126702348078,"DL",0,0.024319413367621223],["carrier_year",1996,"WN",2,747960,14.319279969721144,null,null,null,null,null,null],["carrier_year",1996,"UA",3,715447,13.696836590856574,null,null,null,null,null,null],["carrier_year",1996,"US",4,708343,13.560834445147046,null,null,null,null,null,null],["carrier_year",1996,"AA",5,641045,12.272451505681977,null,null,null,null,null,null],["carrier_year",1997,"DL",1,908107,17.08869644416343,"DL","WN",2.303258513232768,"DL",0,0.4103038042074765],["carrier_year",1997,"WN",2,785710,14.78543793093066,null,null,null,null,null,null],["carrier_year",1997,"UA",3,726980,13.680260741276006,null,null,null,null,null,null],["carrier_year",1997,"US",4,705973,13.28495242826604,null,null,null,null,null,null],["carrier_year",1997,"AA",5,650475,12.240594797217957,null,null,null,null,null,null],["carrier_year",1998,"DL",1,898170,17.13995540638432,"DL","WN",1.7564365716501538,"DL",0,0.0512589622208921],["carrier_year",1998,"WN",2,806129,15.383518834734167,null,null,null,null,null,null],["carrier_year",1998,"UA",3,727775,13.888273985861641,null,null,null,null,null,null],["carrier_year",1998,"US",4,678816,12.95397972448443,null,null,null,null,null,null],["carrier_year",1998,"AA",5,639257,12.199067518642375,null,null,null,null,null,null],["carrier_year",1999,"DL",1,892387,16.60695779139876,"DL","WN",0.8618101959348099,"DL",0,-0.5329976149855611],["carrier_year",1999,"WN",2,846077,15.74514759546395,null,null,null,null,null,null],["carrier_year",1999,"UA",3,747011,13.901569774896517,null,null,null,null,null,null],["carrier_year",1999,"US",4,682121,12.693993363447374,null,null,null,null,null,null],["carrier_year",1999,"AA",5,663693,12.351055805885581,null,null,null,null,null,null],["carrier_year",2000,"WN",1,902660,16.425268630641078,"WN","DL",0.4767487626822895,"DL",1,-0.18168916075768138],["carrier_year",2000,"DL",2,876460,15.948519867958789,null,null,null,null,null,null],["carrier_year",2000,"UA",3,732400,13.327129533912577,null,null,null,null,null,null],["carrier_year",2000,"US",4,720569,13.1118465334815,null,null,null,null,null,null],["carrier_year",2000,"AA",5,712588,12.966620126039999,null,null,null,null,null,null],["carrier_year",2001,"WN",1,939453,16.376528741330638,"WN","DL",2.3889312486076193,"WN",0,-0.04873988931043982],["carrier_year",2001,"DL",2,802410,13.987597492723019,null,null,null,null,null,null],["carrier_year",2001,"AA",3,684755,11.936637530850252,null,null,null,null,null,null],["carrier_year",2001,"UA",4,670455,11.687360173706224,null,null,null,null,null,null],["carrier_year",2001,"US",5,658132,11.472545846986934,null,null,null,null,null,null],["carrier_year",2002,"WN",1,946507,18.180325211247478,"WN","AA",1.9916384567985652,"WN",0,1.8037964699168398],["carrier_year",2002,"AA",2,842818,16.188686754448913,null,null,null,null,null,null],["carrier_year",2002,"DL",3,720693,13.842933139923506,null,null,null,null,null,null],["carrier_year",2002,"UA",4,583265,11.203242431739291,null,null,null,null,null,null],["carrier_year",2002,"US",5,506032,9.719765756933635,null,null,null,null,null,null],["carrier_year",2003,"WN",1,948848,14.855842365692006,"WN","AA",3.2673528065287254,"WN",0,-3.3244828455554725],["carrier_year",2003,"AA",2,740161,11.58848955916328,null,null,null,null,null,null],["carrier_year",2003,"DL",3,653674,10.23438728073554,null,null,null,null,null,null],["carrier_year",2003,"UA",4,538038,8.423907427482796,null,null,null,null,null,null],["carrier_year",2003,"NW",5,493046,7.719480522733863,null,null,null,null,null,null],["carrier_year",2004,"WN",1,980301,14.001349715595643,"WN","AA",4.201329005673804,"WN",0,-0.854492650096363],["carrier_year",2004,"AA",2,686146,9.800020709921839,null,null,null,null,null,null],["carrier_year",2004,"DL",3,676921,9.668262758918656,null,null,null,null,null,null],["carrier_year",2004,"UA",4,549264,7.844975523014793,null,null,null,null,null,null],["carrier_year",2004,"NW",5,501276,7.159577089113365,null,null,null,null,null,null],["carrier_year",2005,"WN",1,1027275,14.661012389449237,"WN","AA",5.187181855924111,"WN",0,0.6596626738535942],["carrier_year",2005,"AA",2,663817,9.473830533525126,null,null,null,null,null,null],["carrier_year",2005,"DL",3,640571,9.142069423788067,null,null,null,null,null,null],["carrier_year",2005,"MQ",4,515286,7.354033175254669,null,null,null,null,null,null],["carrier_year",2005,"OO",5,507815,7.247408928035983,null,null,null,null,null,null],["carrier_year",2006,"WN",1,1090370,15.53236273338359,"WN","AA",6.508558134287409,"WN",0,0.8713503439343526],["carrier_year",2006,"AA",2,633470,9.02380459909618,null,null,null,null,null,null],["carrier_year",2006,"OO",3,535265,7.62487058382436,null,null,null,null,null,null],["carrier_year",2006,"MQ",4,530098,7.551266469401372,null,null,null,null,null,null],["carrier_year",2006,"US",5,499333,7.113017857010582,null,null,null,null,null,null],["carrier_year",2007,"WN",1,1158878,15.886686254540829,"WN","AA",7.4430586036422035,"WN",0,0.3543235211572391],["carrier_year",2007,"AA",2,615933,8.443627650898625,null,null,null,null,null,null],["carrier_year",2007,"OO",3,583694,8.001673555506235,null,null,null,null,null,null],["carrier_year",2007,"MQ",4,517702,7.097010425038957,null,null,null,null,null,null],["carrier_year",2007,"UA",5,478073,6.553749193415612,null,null,null,null,null,null],["carrier_year",2008,"WN",1,1189365,17.30666644936902,"WN","AA",8.758647403618065,"WN",0,1.4199801948281898],["carrier_year",2008,"AA",2,587445,8.548019045750953,null,null,null,null,null,null],["carrier_year",2008,"OO",3,554723,8.071875271925212,null,null,null,null,null,null],["carrier_year",2008,"MQ",4,472362,6.8734253802296585,null,null,null,null,null,null],["carrier_year",2008,"US",5,447007,6.504480163392311,null,null,null,null,null,null],["carrier_year",2009,"WN",1,1123684,17.665465370667206,"WN","AA",9.13918264499345,"WN",0,0.35879892129818813],["carrier_year",2009,"AA",2,542349,8.526282725673756,null,null,null,null,null,null],["carrier_year",2009,"OO",3,539792,8.486084062212502,null,null,null,null,null,null],["carrier_year",2009,"MQ",4,425890,6.695427759684623,null,null,null,null,null,null],["carrier_year",2009,"DL",5,423217,6.653405457208311,null,null,null,null,null,null],["carrier_year",2010,"WN",1,1112890,17.562162470951712,"WN","DL",6.229802700453318,"WN",0,-0.10330289971549433],["carrier_year",2010,"DL",2,718116,11.332359770498394,null,null,null,null,null,null],["carrier_year",2010,"OO",3,587688,9.274117062987958,null,null,null,null,null,null],["carrier_year",2010,"AA",4,531817,8.392434615113917,null,null,null,null,null,null],["carrier_year",2010,"MQ",5,424901,6.705227287575459,null,null,null,null,null,null],["carrier_year",2011,"WN",1,1143570,19.15751302957816,"WN","DL",7.059902973596747,"WN",0,1.5953505586264463],["carrier_year",2011,"DL",2,722143,12.097610055981411,null,null,null,null,null,null],["carrier_year",2011,"OO",3,574684,9.62732164877541,null,null,null,null,null,null],["carrier_year",2011,"AA",4,524652,8.789166842426997,null,null,null,null,null,null],["carrier_year",2011,"MQ",5,430520,7.212232315900198,null,null,null,null,null,null],["carrier_year",2012,"WN",1,1130955,18.793183668721646,"WN","EV",6.736103956529686,"WN",0,-0.36432936085651235],["carrier_year",2012,"EV",2,725583,12.05707971219196,null,null,null,null,null,null],["carrier_year",2012,"DL",3,723166,12.016916200003324,null,null,null,null,null,null],["carrier_year",2012,"OO",4,606589,10.079745426145333,null,null,null,null,null,null],["carrier_year",2012,"UA",5,523646,8.701473936090663,null,null,null,null,null,null],["carrier_year",2013,"WN",1,1122264,17.889047050515902,"WN","DL",5.898314648830711,"WN",0,-0.9041366182057438],["carrier_year",2013,"DL",2,752235,11.990732401685191,null,null,null,null,null,null],["carrier_year",2013,"EV",3,727389,11.594683643980126,null,null,null,null,null,null],["carrier_year",2013,"OO",4,612433,9.762268728470847,null,null,null,null,null,null],["carrier_year",2013,"AA",5,528142,8.418658254522617,null,null,null,null,null,null],["carrier_year",2014,"WN",1,1159468,20.36717434062198,"WN","DL",6.422028984896256,"WN",0,2.4781272901060767],["carrier_year",2014,"DL",2,793873,13.945145355725723,null,null,null,null,null,null],["carrier_year",2014,"EV",3,651893,11.451129640862089,null,null,null,null,null,null],["carrier_year",2014,"OO",4,594722,10.44686585417052,null,null,null,null,null,null],["carrier_year",2014,"AA",5,529240,9.296611332120229,null,null,null,null,null,null],["carrier_year",2015,"WN",1,1245812,21.744974643034492,"WN","DL",6.523691373744478,"WN",0,1.3778003024125134],["carrier_year",2015,"DL",2,872057,15.221283269290014,null,null,null,null,null,null],["carrier_year",2015,"AA",3,715065,12.481072820876232,null,null,null,null,null,null],["carrier_year",2015,"OO",4,578393,10.095536982071652,null,null,null,null,null,null],["carrier_year",2015,"EV",5,556746,9.717700305191219,null,null,null,null,null,null],["carrier_year",2016,"WN",1,1283578,23.12004563567436,"WN","DL",6.57351124329654,"WN",0,1.3750709926398663],["carrier_year",2016,"DL",2,918630,16.54653439237782,null,null,null,null,null,null],["carrier_year",2016,"AA",3,903628,16.276315578541507,null,null,null,null,null,null],["carrier_year",2016,"OO",4,597487,10.762046955247103,null,null,null,null,null,null],["carrier_year",2016,"UA",5,539597,9.719321509774222,null,null,null,null,null,null],["carrier_year",2017,"WN",1,1311398,23.451629465119986,"WN","DL",7.095908083206382,"WN",0,0.33158382944562703],["carrier_year",2017,"DL",2,914600,16.355721381913604,null,null,null,null,null,null],["carrier_year",2017,"AA",3,884210,15.812259351740465,null,null,null,null,null,null],["carrier_year",2017,"OO",4,696069,12.447748334339844,null,null,null,null,null,null],["carrier_year",2017,"UA",5,578522,10.345664383673105,null,null,null,null,null,null],["carrier_year",2018,"WN",1,1334277,18.820052764196177,"WN","DL",5.480124846678034,"WN",0,-4.631576700923809],["carrier_year",2018,"DL",2,945755,13.339927917518143,null,null,null,null,null,null],["carrier_year",2018,"AA",3,901873,12.720969818563834,null,null,null,null,null,null],["carrier_year",2018,"OO",4,763521,10.769507011341593,null,null,null,null,null,null],["carrier_year",2018,"UA",5,616662,8.69805248660866,null,null,null,null,null,null],["carrier_year",2019,"WN",1,1330324,18.25584926591622,"WN","DL",4.668242325387935,"WN",0,-0.5642034982799586],["carrier_year",2019,"DL",2,990144,13.587606940528284,null,null,null,null,null,null],["carrier_year",2019,"AA",3,926625,12.715944631555633,null,null,null,null,null,null],["carrier_year",2019,"OO",4,818992,11.238912101105637,null,null,null,null,null,null],["carrier_year",2019,"UA",5,620526,8.515391078851414,null,null,null,null,null,null],["carrier_year",2020,"WN",1,883434,20.044698365446575,"WN","OO",7.064792209324487,"WN",0,1.7888490995303563],["carrier_year",2020,"OO",2,572066,12.979906156122087,null,null,null,null,null,null],["carrier_year",2020,"DL",3,552332,12.532151057785685,null,null,null,null,null,null],["carrier_year",2020,"AA",4,535544,12.151239301888676,null,null,null,null,null,null],["carrier_year",2020,"UA",5,285624,6.480673062087618,null,null,null,null,null,null],["carrier_year",2021,"WN",1,1041135,17.66932605607629,"WN","DL",5.036513569335051,"WN",0,-2.3753723093702845],["carrier_year",2021,"DL",2,744367,12.632812486741239,null,null,null,null,null,null],["carrier_year",2021,"OO",3,740005,12.558784046378937,null,null,null,null,null,null],["carrier_year",2021,"AA",4,719921,12.217934161868056,null,null,null,null,null,null],["carrier_year",2021,"UA",5,440730,7.479723640732953,null,null,null,null,null,null],["carrier_year",2022,"WN",1,1264494,19.311626219219228,"WN","DL",5.889562325763098,"WN",0,1.642300163142938],["carrier_year",2022,"DL",2,878855,13.42206389345613,null,null,null,null,null,null],["carrier_year",2022,"AA",3,848251,12.954672977553813,null,null,null,null,null,null],["carrier_year",2022,"OO",4,720245,10.99973762331933,null,null,null,null,null,null],["carrier_year",2022,"UA",5,614587,9.386105764986855,null,null,null,null,null,null],["carrier_year",2023,"WN",1,1424140,21.06730312070171,"WN","DL",6.644571841761053,"WN",0,1.755676901482481],["carrier_year",2023,"DL",2,974970,14.422731278940656,null,null,null,null,null,null],["carrier_year",2023,"AA",3,930553,13.765670594799905,null,null,null,null,null,null],["carrier_year",2023,"UA",4,721942,10.679688057048908,null,null,null,null,null,null],["carrier_year",2023,"OO",5,666911,9.865614470157638,null,null,null,null,null,null],["carrier_year",2024,"WN",1,1407647,20.158931744044533,"WN","DL",5.837531538452064,"WN",0,-0.9083713766571755],["carrier_year",2024,"DL",2,1000027,14.32140020559247,null,null,null,null,null,null],["carrier_year",2024,"AA",3,969054,13.877835453272967,null,null,null,null,null,null],["carrier_year",2024,"UA",4,747973,10.711731459228218,null,null,null,null,null,null],["carrier_year",2024,"OO",5,736131,10.542142016908535,null,null,null,null,null,null],["carrier_year",2025,"WN",1,1262665,19.961694364344616,"WN","DL",5.22090162897759,"WN",0,-0.19723737969991717],["carrier_year",2025,"DL",2,932420,14.740792735367027,null,null,null,null,null,null],["carrier_year",2025,"AA",3,876508,13.85687003591845,null,null,null,null,null,null],["carrier_year",2025,"OO",4,760491,12.02273675823342,null,null,null,null,null,null],["carrier_year",2025,"UA",5,723139,11.43223238225325,null,null,null,null,null,null]],
      count: 195
    };

    const palette = ["#0e3a52","#3c88b5","#1f8a70","#d48a1f","#c54f36","#7c6a0a","#5c7080","#9d5d8f","#2e786f","#87543f","#335f9b","#7b8c36"];

    function rowObjects(payload) {
      return payload.rows.map(row => Object.fromEntries(payload.columns.map((col, idx) => [col, row[idx]])));
    }
    function fmtInt(value) {
      return Number(value).toLocaleString("en-US");
    }
    function fmtPct(value) {
      return `${Number(value).toFixed(2)} pts`;
    }
    function fmtShare(value) {
      return `${Number(value).toFixed(2)}%`;
    }
    function escapeHTML(value) {
      return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;");
    }
    function unique(values) {
      return [...new Set(values)];
    }
    function slug(value) {
      return String(value).toLowerCase().replace(/[^a-z0-9]+/g, "-");
    }
    function seriesColorMap(keys) {
      const map = new Map();
      keys.forEach((key, index) => map.set(key, palette[index % palette.length]));
      return map;
    }
    function scaleLinear(domainMin, domainMax, rangeMin, rangeMax) {
      const span = domainMax - domainMin || 1;
      return value => rangeMin + ((value - domainMin) / span) * (rangeMax - rangeMin);
    }
    function buildPath(points) {
      return points.map((point, index) => `${index === 0 ? "M" : "L"} ${point[0].toFixed(2)} ${point[1].toFixed(2)}`).join(" ");
    }
    function longestStreak(leaders) {
      let best = null;
      let current = null;
      leaders.forEach(row => {
        if (!current || current.airline !== row.Reporting_Airline || row.Year !== current.end + 1) {
          current = { airline: row.Reporting_Airline, start: row.Year, end: row.Year, years: 1 };
        } else {
          current.end = row.Year;
          current.years += 1;
        }
        if (!best || current.years > best.years) {
          best = { ...current };
        }
      });
      return best;
    }
    function renderKpis(metrics) {
      const items = [
        { label: "Total years analyzed", value: fmtInt(metrics.totalYears), meta: `${metrics.startYear} to ${metrics.endYear}` },
        { label: "Distinct annual leaders", value: fmtInt(metrics.distinctLeaders), meta: `${metrics.leaderCounts.map(item => `${item.airline} ${item.count}`).join(" • ")}` },
        { label: "Largest leader share gap", value: fmtPct(metrics.largestGap.LeaderShareGapPctPts), meta: `${metrics.largestGap.Year} • ${metrics.largestGap.LeaderReportingAirline} over ${metrics.largestGap.RunnerUpReportingAirline}` },
        { label: "Sharpest leadership transition", value: fmtPct(Math.abs(metrics.sharpest.LeaderShareChangePctPts)), meta: `${metrics.sharpest.Year} • ${metrics.sharpest.PriorYearLeaderReportingAirline} to ${metrics.sharpest.Reporting_Airline}` }
      ];
      document.getElementById("kpis").innerHTML = items.map(item => `
        <article class="card kpi">
          <div class="label">${escapeHTML(item.label)}</div>
          <div class="value">${escapeHTML(item.value)}</div>
          <div class="meta">${escapeHTML(item.meta)}</div>
        </article>
      `).join("");
    }
    function renderLeaderFrequency(metrics, colors) {
      document.getElementById("leaderSummary").textContent =
        `Only ${metrics.distinctLeaders} carriers lead across ${metrics.totalYears} observed years. ${metrics.longest.airline} owns the longest uninterrupted run, holding the top spot from ${metrics.longest.start} through ${metrics.longest.end}.`;
      document.getElementById("leaderList").innerHTML = metrics.leaderCounts.map(item => `
        <div class="leader-row">
          <div class="leader-code" style="background:${colors.get(item.airline)}">${escapeHTML(item.airline)}</div>
          <div><strong>${escapeHTML(item.airline)}</strong><br><span>${item.count === 1 ? "1 annual lead" : `${item.count} annual leads`}</span></div>
          <div class="leader-years">${escapeHTML(String(item.count))}</div>
        </div>
      `).join("");
    }
    function renderTransitions(metrics) {
      document.getElementById("transitionSummary").textContent =
        `${metrics.transitions.length} true handoffs occur in the full series. The largest swing happens in ${metrics.sharpest.Year}, when ${metrics.sharpest.PriorYearLeaderReportingAirline} gives way to ${metrics.sharpest.Reporting_Airline}.`;
      document.getElementById("transitionCallouts").innerHTML = metrics.transitions.map(row => {
        const sharpest = row.Year === metrics.sharpest.Year;
        return `
          <div class="callout ${sharpest ? "sharpest" : ""}">
            <div class="title">${row.Year}: ${row.PriorYearLeaderReportingAirline} → ${row.Reporting_Airline}${sharpest ? " • sharpest handoff" : ""}</div>
            <div class="desc">Leader share swing ${fmtPct(row.LeaderShareChangePctPts)} and leader-over-runner-up gap ${fmtPct(row.LeaderShareGapPctPts)}.</div>
          </div>
        `;
      }).join("");
    }
    function renderLegend(targetId, carriers, colors) {
      document.getElementById(targetId).innerHTML = carriers.map(carrier => `
        <span class="legend-item"><span class="swatch" style="background:${colors.get(carrier)}"></span>${escapeHTML(carrier)}</span>
      `).join("");
    }
    function renderBumpChart(rows, metrics, colors) {
      const svg = document.getElementById("bumpChart");
      const width = 1100;
      const height = 420;
      const margin = { top: 28, right: 70, bottom: 46, left: 60 };
      const years = unique(rows.map(r => r.Year)).sort((a, b) => a - b);
      const carriers = unique(rows.map(r => r.Reporting_Airline)).sort((a, b) => a.localeCompare(b));
      const x = scaleLinear(years[0], years[years.length - 1], margin.left, width - margin.right);
      const y = rank => margin.top + ((rank - 1) / 4) * (height - margin.top - margin.bottom);
      const bandX = x(metrics.sharpest.Year);
      let out = "";

      for (let rank = 1; rank <= 5; rank += 1) {
        const yy = y(rank);
        out += `<line class="grid-line" x1="${margin.left}" y1="${yy}" x2="${width - margin.right}" y2="${yy}"></line>`;
        out += `<text class="grid-label" x="${margin.left - 12}" y="${yy + 4}" text-anchor="end">Rank ${rank}</text>`;
      }
      years.forEach(year => {
        const xx = x(year);
        out += `<line class="grid-line" x1="${xx}" y1="${margin.top}" x2="${xx}" y2="${height - margin.bottom}"></line>`;
        out += `<text class="axis-label" x="${xx}" y="${height - 18}" text-anchor="middle">${year}</text>`;
      });
      out += `<rect x="${bandX - 7}" y="${margin.top}" width="14" height="${height - margin.top - margin.bottom}" fill="rgba(197,79,54,0.12)"></rect>`;

      carriers.forEach(carrier => {
        const points = rows
          .filter(r => r.Reporting_Airline === carrier)
          .sort((a, b) => a.Year - b.Year)
          .map(r => [x(r.Year), y(r.RankInYear)]);
        if (points.length < 2) return;
        const prominent = metrics.leaderSet.has(carrier);
        out += `<path d="${buildPath(points)}" fill="none" stroke="${colors.get(carrier)}" stroke-width="${prominent ? 3.5 : 2}" stroke-linecap="round" stroke-linejoin="round" opacity="${prominent ? 0.95 : 0.72}"></path>`;
        const last = points[points.length - 1];
        out += `<text class="axis-label" x="${last[0] + 8}" y="${last[1] + 4}" fill="${colors.get(carrier)}">${carrier}</text>`;
      });

      out += `<text class="axis-label" x="${bandX + 10}" y="${margin.top + 16}" fill="var(--red)">Sharpest transition</text>`;
      svg.innerHTML = out;
    }
    function renderShareChart(rows, metrics, colors) {
      const svg = document.getElementById("shareChart");
      const width = 1100;
      const height = 420;
      const margin = { top: 24, right: 40, bottom: 46, left: 64 };
      const years = unique(rows.map(r => r.Year)).sort((a, b) => a - b);
      const leadersOnly = rows.filter(r => metrics.leaderSet.has(r.Reporting_Airline));
      const maxShare = Math.max(...leadersOnly.map(r => Number(r.SharePct)));
      const minShare = Math.min(...leadersOnly.map(r => Number(r.SharePct)));
      const x = scaleLinear(years[0], years[years.length - 1], margin.left, width - margin.right);
      const y = value => margin.top + ((maxShare - value) / (maxShare - minShare || 1)) * (height - margin.top - margin.bottom);
      const bandX = x(metrics.sharpest.Year);
      let out = "";

      const yTicks = 5;
      for (let i = 0; i <= yTicks; i += 1) {
        const value = minShare + ((maxShare - minShare) / yTicks) * i;
        const yy = y(value);
        out += `<line class="grid-line" x1="${margin.left}" y1="${yy}" x2="${width - margin.right}" y2="${yy}"></line>`;
        out += `<text class="grid-label" x="${margin.left - 12}" y="${yy + 4}" text-anchor="end">${value.toFixed(1)}%</text>`;
      }
      years.forEach(year => {
        const xx = x(year);
        out += `<line class="grid-line" x1="${xx}" y1="${margin.top}" x2="${xx}" y2="${height - margin.bottom}"></line>`;
        out += `<text class="axis-label" x="${xx}" y="${height - 18}" text-anchor="middle">${year}</text>`;
      });
      out += `<rect x="${bandX - 7}" y="${margin.top}" width="14" height="${height - margin.top - margin.bottom}" fill="rgba(197,79,54,0.12)"></rect>`;

      [...metrics.leaderSet].forEach(carrier => {
        const points = rows
          .filter(r => r.Reporting_Airline === carrier)
          .sort((a, b) => a.Year - b.Year)
          .map(r => [x(r.Year), y(Number(r.SharePct)), r]);
        if (points.length < 2) return;
        out += `<path d="${buildPath(points.map(point => [point[0], point[1]]))}" fill="none" stroke="${colors.get(carrier)}" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round"></path>`;
        points.forEach(point => {
          const isSharpestYear = point[2].Year === metrics.sharpest.Year && point[2].Reporting_Airline === metrics.sharpest.Reporting_Airline;
          out += `<circle cx="${point[0]}" cy="${point[1]}" r="${isSharpestYear ? 5.5 : 3.5}" fill="${isSharpestYear ? "var(--red)" : colors.get(carrier)}" stroke="#fff" stroke-width="1.5"></circle>`;
        });
      });
      out += `<text class="axis-label" x="${bandX + 10}" y="${margin.top + 16}" fill="var(--red)">${metrics.sharpest.Year} handoff</text>`;
      svg.innerHTML = out;
    }
    function renderTransitionTable(metrics) {
      document.getElementById("tableMeta").textContent = `${metrics.transitions.length} change years detected`;
      document.getElementById("transitionTable").innerHTML = metrics.transitions.map(row => `
        <tr class="${row.Year === metrics.sharpest.Year ? "sharpest" : ""}">
          <td>${row.Year}</td>
          <td>${escapeHTML(row.PriorYearLeaderReportingAirline)}</td>
          <td>${escapeHTML(row.Reporting_Airline)}</td>
          <td class="num">${fmtPct(row.LeaderShareChangePctPts)}</td>
          <td class="num">${fmtPct(row.LeaderShareGapPctPts)}</td>
        </tr>
      `).join("");
    }
    function exportTransitions(metrics) {
      const header = ["Year","PriorLeader","NewLeader","LeaderShareChangePctPts","LeaderShareGapPctPts"];
      const lines = [header.join(",")].concat(metrics.transitions.map(row => [
        row.Year,
        row.PriorYearLeaderReportingAirline,
        row.Reporting_Airline,
        Number(row.LeaderShareChangePctPts).toFixed(6),
        Number(row.LeaderShareGapPctPts).toFixed(6)
      ].join(",")));
      const blob = new Blob([lines.join("\n")], { type: "text/csv;charset=utf-8" });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `${dashboardId}-transitions.csv`;
      a.click();
      URL.revokeObjectURL(url);
    }
    function toggleLedgerEntry(id) {
      const panel = document.getElementById(id);
      const toggle = document.getElementById(`${id}-toggle`);
      const open = panel.classList.toggle("open");
      toggle.textContent = open ? "▼" : "▶";
    }
    function initFooter() {
      const tokenInput = document.getElementById("tokenInput");
      const sqlTextarea = document.getElementById("sqlTextarea");
      const statusText = document.getElementById("statusText");
      const stored = localStorage.getItem(authKey);
      if (stored) tokenInput.value = stored;
      sqlTextarea.value = savedSQL;
      document.getElementById("ledgerSql").textContent = savedSQL;
      document.getElementById("fetchBtn").addEventListener("click", () => {
        if (tokenInput.value.trim()) {
          localStorage.setItem(authKey, tokenInput.value.trim());
          statusText.textContent = "Token saved locally. This page remains a static artifact and does not rerun queries in-browser.";
        } else {
          statusText.textContent = "Enter a token to store it locally, or use Forget to clear the existing token.";
        }
      });
      document.getElementById("forgetBtn").addEventListener("click", () => {
        localStorage.removeItem(authKey);
        tokenInput.value = "";
        statusText.textContent = "Stored token cleared from localStorage.";
      });
    }

    (function main() {
      const rows = rowObjects(result);
      const allRows = rows.slice().sort((a, b) => a.Year - b.Year || a.RankInYear - b.RankInYear || a.Reporting_Airline.localeCompare(b.Reporting_Airline));
      const leaderRows = allRows.filter(row => Number(row.RankInYear) === 1);
      const transitions = leaderRows.filter(row => Number(row.LeaderChanged) === 1 && row.PriorYearLeaderReportingAirline);
      const years = leaderRows.map(row => Number(row.Year));
      const leaderCounts = Array.from(leaderRows.reduce((map, row) => {
        map.set(row.Reporting_Airline, (map.get(row.Reporting_Airline) || 0) + 1);
        return map;
      }, new Map()).entries())
        .map(([airline, count]) => ({ airline, count }))
        .sort((a, b) => b.count - a.count || a.airline.localeCompare(b.airline));
      const largestGap = leaderRows.reduce((best, row) => !best || Number(row.LeaderShareGapPctPts) > Number(best.LeaderShareGapPctPts) ? row : best, null);
      const sharpest = transitions.reduce((best, row) => !best || Math.abs(Number(row.LeaderShareChangePctPts)) > Math.abs(Number(best.LeaderShareChangePctPts)) ? row : best, null);
      const longest = longestStreak(leaderRows);
      const colors = seriesColorMap(unique(allRows.map(row => row.Reporting_Airline)).sort((a, b) => a.localeCompare(b)));

      const metrics = {
        totalYears: leaderRows.length,
        startYear: Math.min(...years),
        endYear: Math.max(...years),
        distinctLeaders: leaderCounts.length,
        leaderCounts,
        largestGap,
        sharpest,
        transitions,
        longest,
        leaderSet: new Set(leaderCounts.map(item => item.airline))
      };

      renderKpis(metrics);
      renderLeaderFrequency(metrics, colors);
      renderTransitions(metrics);
      renderLegend("bumpLegend", unique(allRows.map(row => row.Reporting_Airline)).sort((a, b) => a.localeCompare(b)), colors);
      renderLegend("shareLegend", [...metrics.leaderSet], colors);
      renderBumpChart(allRows, metrics, colors);
      renderShareChart(allRows, metrics, colors);
      renderTransitionTable(metrics);
      document.getElementById("subtitle").textContent =
        `From ${metrics.startYear} to ${metrics.endYear}, only ${metrics.distinctLeaders} carriers top the annual completed-flight ranking. ${metrics.longest.airline} then sustains the longest hold at ${metrics.longest.years} consecutive years (${metrics.longest.start}–${metrics.longest.end}).`;
      document.getElementById("stabilityNote").textContent =
        `Stable dominance is present. The longest uninterrupted leadership spell is ${metrics.longest.years} years by ${metrics.longest.airline}, spanning ${metrics.longest.start} through ${metrics.longest.end}. Before that, the series changes hands only in ${metrics.transitions.map(row => row.Year).join(", ")}.`;
      document.getElementById("ledgerRows").textContent = `${fmtInt(result.count)} rows`;
      document.getElementById("exportTransitions").addEventListener("click", () => exportTransitions(metrics));
      initFooter();
    })();
  </script>
</body>
</html>
```