#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const ROOT = '/Users/bvt/work/ExploringDatabyLLMs';
const RUNS_DIR = path.join(ROOT, 'runs');
const ALLOWED_DATES = new Set(process.argv.slice(2));

function walk(dir, out = []) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (dir === RUNS_DIR && ALLOWED_DATES.size > 0 && !ALLOWED_DATES.has(entry.name)) continue;
      walk(full, out);
    }
    else if (entry.isFile() && entry.name === 'visual.html') out.push(full);
  }
  return out;
}

function replaceOnce(text, searchValue, replaceValue, label, file) {
  if (!text.includes(searchValue)) {
    throw new Error(`${label} not found in ${file}`);
  }
  return text.replace(searchValue, replaceValue);
}

function replaceRegex(text, pattern, replaceValue, label, file) {
  if (!pattern.test(text)) {
    throw new Error(`${label} not found in ${file}`);
  }
  return text.replace(pattern, replaceValue);
}

function migrateVariantA(source, file) {
  let text = source;

  text = text.replace(/Enter your token to load the dashboard/g, 'Sign in to load the dashboard');
  text = text.replace(/Paste your JWE access token in the footer below, then click <strong>Run Query<\/strong>\./g, 'Use Google sign-in from the footer below, then click <strong>Run Query</strong>.');
  text = text.replace(/JWE Access Token/g, 'OAuth2 access');
  text = text.replace(/Forget Token/g, 'Logout');
  text = text.replace(/Enter token and click Run Query\./g, 'Sign in, then click Run Query.');
  text = text.replace(/placeholder="Paste token…"/g, 'placeholder="OAuth2 handled automatically"');

  text = replaceOnce(
    text,
    "const STORAGE_KEY = 'OnTimeAnalystDashboard::auth::jwe';",
    "const MCP_BASE = 'https://mcp.demo.altinity.cloud';\nconst OAUTH_STORAGE_KEY = 'OnTimeAnalystDashboard::auth::oauth2';",
    'variant A storage constant',
    file
  );

  text = replaceOnce(
    text,
    "const $tooltip = document.getElementById('tooltip');",
    "const $tooltip = document.getElementById('tooltip');\nlet oauthState = {};",
    'variant A oauth state insertion',
    file
  );

  text = replaceRegex(
    text,
    /\$sql\.value = SQL_PRIMARY;\s*const stored = localStorage\.getItem\(STORAGE_KEY\);\s*if \(stored\) \{ \$jwe\.value = stored; setTimeout\(runAll, 200\); \}/,
    '$sql.value = SQL_PRIMARY;',
    'variant A boot block',
    file
  );

  const helperBlock = `
/* ── OAuth helpers ─────────────────────────────────────── */
function loadAuthState() {
  try { oauthState = JSON.parse(localStorage.getItem(OAUTH_STORAGE_KEY)) || {}; } catch (err) { oauthState = {}; }
}
function saveAuthState() {
  localStorage.setItem(OAUTH_STORAGE_KEY, JSON.stringify(oauthState));
}
function clearAuthState() {
  localStorage.removeItem(OAUTH_STORAGE_KEY);
  oauthState = {};
}
function hasValidAccessToken() {
  return Boolean(oauthState.access_token && (!oauthState.expires_at || Date.now() < oauthState.expires_at));
}
function currentRedirectUri() {
  return window.location.origin + window.location.pathname;
}
function randomString(len) {
  const arr = new Uint8Array(len);
  crypto.getRandomValues(arr);
  return Array.from(arr, (b) => b.toString(16).padStart(2, '0')).join('');
}
async function sha256(plain) {
  const data = new TextEncoder().encode(plain);
  const hash = await crypto.subtle.digest('SHA-256', data);
  return btoa(String.fromCharCode(...new Uint8Array(hash))).replace(/\\+/g, '-').replace(/\\//g, '_').replace(/=+$/, '');
}
async function registerOAuthClient() {
  const resp = await fetch(MCP_BASE + '/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      client_name: 'OnTime Analyst Dashboard',
      redirect_uris: [currentRedirectUri()],
      grant_types: ['authorization_code'],
      response_types: ['code'],
      token_endpoint_auth_method: 'none'
    })
  });
  if (!resp.ok) throw new Error('OAuth registration failed: ' + await resp.text());
  return resp.json();
}
function ensureLoginButton() {
  let btn = document.getElementById('oauth-login-btn');
  if (!btn) {
    btn = document.createElement('button');
    btn.id = 'oauth-login-btn';
    btn.type = 'button';
    btn.className = 'btn btn-ghost';
    btn.textContent = 'Login with Google';
    $runBtn.parentNode.insertBefore(btn, $runBtn);
    btn.addEventListener('click', startLogin);
  }
  return btn;
}
function applyAuthCopy() {
  const loginBtn = ensureLoginButton();
  const label = document.querySelector('label[for="jwe-input"]');
  const formGroup = $jwe.closest('.form-group');
  if (label) label.textContent = 'OAuth2 redirect flow';
  if (formGroup) formGroup.style.display = 'none';
  if ($forget) $forget.textContent = 'Logout';
  const emptyTitle = document.querySelector('#empty-state h3');
  const emptyBody = document.querySelector('#empty-state p');
  if (emptyTitle) emptyTitle.textContent = 'Sign in to load the dashboard';
  if (emptyBody) emptyBody.innerHTML = 'Use Google sign-in from the footer below, then click <strong>Run Query</strong>.';
  if (hasValidAccessToken()) {
    $jwe.value = '__oauth_authenticated__';
    loginBtn.style.display = 'none';
    if ($forget) $forget.style.display = 'inline-flex';
  } else {
    $jwe.value = '';
    loginBtn.style.display = 'inline-flex';
    if ($forget) $forget.style.display = 'none';
  }
}
function updateAuthUi() {
  applyAuthCopy();
  if (hasValidAccessToken()) {
    setStatus('Authenticated with OAuth2. Query runs use the stored bearer token until it expires or you log out.', 'ok');
  } else if (oauthState.expires_at && Date.now() >= oauthState.expires_at) {
    setStatus('OAuth2 session expired. Sign in again, then rerun the dashboard.', 'err');
  } else {
    setStatus('Sign in with Google. This page must be registered as an OAuth redirect URI in Google Console and the Altinity MCP provider setup.', '');
  }
}
async function startLogin() {
  try {
    setStatus('Registering OAuth client and redirecting to Google sign-in.', '');
    const reg = await registerOAuthClient();
    const codeVerifier = randomString(32);
    const codeChallenge = await sha256(codeVerifier);
    const stateToken = randomString(16);
    oauthState.client_id = reg.client_id;
    oauthState.code_verifier = codeVerifier;
    oauthState.state = stateToken;
    saveAuthState();
    const params = new URLSearchParams({
      response_type: 'code',
      client_id: reg.client_id,
      redirect_uri: currentRedirectUri(),
      scope: 'openid email',
      state: stateToken,
      code_challenge: codeChallenge,
      code_challenge_method: 'S256'
    });
    window.location.href = MCP_BASE + '/authorize?' + params.toString();
  } catch (err) {
    setStatus('Login failed: ' + err.message, 'err');
    updateAuthUi();
  }
}
async function exchangeCode(code) {
  const resp = await fetch(MCP_BASE + '/token', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({
      grant_type: 'authorization_code',
      code,
      redirect_uri: currentRedirectUri(),
      client_id: oauthState.client_id,
      code_verifier: oauthState.code_verifier
    })
  });
  if (!resp.ok) throw new Error('Token exchange failed: ' + await resp.text());
  return resp.json();
}
async function initAuth() {
  loadAuthState();
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');
  const stateToken = params.get('state');
  if (code && stateToken) {
    window.history.replaceState({}, '', window.location.pathname);
    if (stateToken !== oauthState.state) {
      clearAuthState();
      updateAuthUi();
      setStatus('OAuth state mismatch. Please sign in again.', 'err');
      return false;
    }
    try {
      setStatus('Exchanging OAuth authorization code.', '');
      const token = await exchangeCode(code);
      oauthState.access_token = token.access_token;
      oauthState.token_type = token.token_type || 'Bearer';
      if (token.expires_in) oauthState.expires_at = Date.now() + Number(token.expires_in) * 1000;
      delete oauthState.code_verifier;
      delete oauthState.state;
      saveAuthState();
    } catch (err) {
      clearAuthState();
      updateAuthUi();
      setStatus('OAuth callback failed: ' + err.message, 'err');
      return false;
    }
  }
  if (oauthState.expires_at && Date.now() >= oauthState.expires_at) {
    clearAuthState();
  }
  updateAuthUi();
  return hasValidAccessToken();
}
function logout() {
  clearAuthState();
  updateAuthUi();
  setStatus('Logged out. Sign in again to rerun the dashboard.', '');
}

`;

  text = replaceOnce(text, "/* ── Helpers ────────────────────────────────────────────── */", helperBlock + "/* ── Helpers ────────────────────────────────────────────── */", 'variant A oauth helper insert', file);

  text = replaceRegex(
    text,
    /async function execQuery\(jwe, sql\) \{[\s\S]*?return \(data\.rows \?\? \[\]\)\.map\(r => \{ const o = \{\}; cols\.forEach\(\(c,i\) => \{ o\[c\] = r\[i\]; \}\); return o; \}\);\s*\}/,
    `async function execQuery(jwe, sql) {
  if (!hasValidAccessToken()) {
    updateAuthUi();
    throw new Error('OAuth login is required.');
  }
  const url = MCP_BASE + '/openapi/execute_query?query=' + encodeURIComponent(sql);
  const resp = await fetch(url, {
    headers: { Authorization: 'Bearer ' + oauthState.access_token }
  });
  if (!resp.ok) {
    const t = await resp.text();
    if (resp.status === 401 || resp.status === 403) {
      clearAuthState();
      updateAuthUi();
      throw new Error('Authentication expired. Please login again.');
    }
    throw new Error(\`HTTP \${resp.status}: \${t.slice(0,200)}\`);
  }
  const data = await resp.json();
  const cols = data.columns ?? [];
  return (data.rows ?? []).map(r => { const o = {}; cols.forEach((c,i) => { o[c] = r[i]; }); return o; });
}`,
    'variant A execQuery',
    file
  );

  text = text.replace(/localStorage\.setItem\(STORAGE_KEY, jwe\);\s*/g, '');

  text = replaceRegex(
    text,
    /\$forget\.addEventListener\('click', \(\) => \{\s*localStorage\.removeItem\(STORAGE_KEY\);\s*\$jwe\.value = '';\s*setStatus\('Token forgotten\.', ''\);\s*\}\);\s*\n\s*\n\}\)\(\);\s*\n<\/script>/,
    `$forget.addEventListener('click', logout);

(async function () {
  const isAuthenticated = await initAuth();
  if (isAuthenticated) {
    setTimeout(runAll, 200);
  }
})();

})();
</script>`,
    'variant A footer/init block',
    file
  );

  return text;
}

function migrateVariantB(source, file) {
  let text = source;

  text = text.replace(/JWE token/g, 'OAuth2 access');
  text = text.replace(/Paste JWE token/g, 'OAuth2 handled automatically');
  text = text.replace(/Forget stored token/g, 'Logout');
  text = text.replace(/stored token/gi, 'OAuth session');
  text = text.replace(/JWE, date range, and SQL editors/g, 'OAuth2, date range, and SQL editors');
  text = text.replace(/Enter the MCP JWE token/g, 'Sign in with OAuth2');
  text = text.replace(/Enter a JWE token and run the queries from the footer controls\./g, 'Sign in and run the queries from the footer controls.');

  text = replaceRegex(
    text,
    /const STORAGE_KEY = 'OnTimeAnalystDashboard::auth::jwe';\s*const ENDPOINT_TEMPLATE = 'https:\/\/mcp\.demo\.altinity\.cloud\/\{JWE\}\/openapi\/execute_query\?query=';/,
    "const MCP_BASE = 'https://mcp.demo.altinity.cloud';\n    const OAUTH_STORAGE_KEY = 'OnTimeAnalystDashboard::auth::oauth2';",
    'variant B constants',
    file
  );

  text = replaceOnce(
    text,
    "    const state = {",
    "    let oauthState = {};\n\n    const state = {",
    'variant B oauth state insert',
    file
  );

  const helperBlock = `
    function loadAuthState() {
      try { oauthState = JSON.parse(localStorage.getItem(OAUTH_STORAGE_KEY)) || {}; } catch (err) { oauthState = {}; }
    }

    function saveAuthState() {
      localStorage.setItem(OAUTH_STORAGE_KEY, JSON.stringify(oauthState));
    }

    function clearAuthState() {
      localStorage.removeItem(OAUTH_STORAGE_KEY);
      oauthState = {};
    }

    function hasValidAccessToken() {
      return Boolean(oauthState.access_token && (!oauthState.expires_at || Date.now() < oauthState.expires_at));
    }

    function currentRedirectUri() {
      return window.location.origin + window.location.pathname;
    }

    function randomString(len) {
      const arr = new Uint8Array(len);
      crypto.getRandomValues(arr);
      return Array.from(arr, (b) => b.toString(16).padStart(2, '0')).join('');
    }

    async function sha256(plain) {
      const data = new TextEncoder().encode(plain);
      const hash = await crypto.subtle.digest('SHA-256', data);
      return btoa(String.fromCharCode(...new Uint8Array(hash))).replace(/\\+/g, '-').replace(/\\//g, '_').replace(/=+$/, '');
    }

    async function registerOAuthClient() {
      const response = await fetch(MCP_BASE + '/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          client_name: 'OnTime Analyst Dashboard',
          redirect_uris: [currentRedirectUri()],
          grant_types: ['authorization_code'],
          response_types: ['code'],
          token_endpoint_auth_method: 'none'
        })
      });
      if (!response.ok) throw new Error('OAuth registration failed: ' + await response.text());
      return response.json();
    }

    function ensureLoginButton() {
      let btn = document.getElementById('oauth-login-btn');
      if (!btn) {
        btn = document.createElement('button');
        btn.id = 'oauth-login-btn';
        btn.type = 'button';
        btn.className = 'secondary';
        btn.textContent = 'Login with Google';
        els.runBtn.parentNode.insertBefore(btn, els.runBtn);
        btn.addEventListener('click', startLogin);
      }
      return btn;
    }

    function applyAuthCopy() {
      const label = document.querySelector('label[for="tokenInput"], label[for="token-input"]');
      const tokenGroup = els.tokenInput?.closest('div');
      const loginBtn = ensureLoginButton();
      if (label) label.textContent = 'OAuth2 redirect flow';
      if (tokenGroup) tokenGroup.style.display = 'none';
      if (els.forgetBtn) els.forgetBtn.textContent = 'Logout';
      if (hasValidAccessToken()) {
        els.tokenInput.value = '__oauth_authenticated__';
        loginBtn.style.display = 'inline-flex';
        loginBtn.style.display = 'none';
        if (els.forgetBtn) els.forgetBtn.style.display = 'inline-flex';
      } else {
        els.tokenInput.value = '';
        loginBtn.style.display = 'inline-flex';
        if (els.forgetBtn) els.forgetBtn.style.display = 'none';
      }
    }

    function updateAuthUi() {
      applyAuthCopy();
      if (hasValidAccessToken()) {
        setStatus('Authenticated with OAuth2. Query runs now use a stored bearer token.', 'ok');
      } else if (oauthState.expires_at && Date.now() >= oauthState.expires_at) {
        setStatus('OAuth2 session expired. Sign in again to rerun the dashboard.', 'warn');
      } else {
        setStatus('Sign in with Google. This page must be registered as an OAuth redirect URI in Google Console and the Altinity MCP provider setup.', '');
      }
    }

    async function startLogin() {
      try {
        setStatus('Registering OAuth client and redirecting to Google sign-in.', '');
        const registration = await registerOAuthClient();
        const codeVerifier = randomString(32);
        const codeChallenge = await sha256(codeVerifier);
        const stateToken = randomString(16);
        oauthState.client_id = registration.client_id;
        oauthState.code_verifier = codeVerifier;
        oauthState.state = stateToken;
        saveAuthState();
        const params = new URLSearchParams({
          response_type: 'code',
          client_id: registration.client_id,
          redirect_uri: currentRedirectUri(),
          scope: 'openid email',
          state: stateToken,
          code_challenge: codeChallenge,
          code_challenge_method: 'S256'
        });
        window.location.href = MCP_BASE + '/authorize?' + params.toString();
      } catch (err) {
        setStatus('Login failed: ' + err.message, 'error');
        updateAuthUi();
      }
    }

    async function exchangeCode(code) {
      const response = await fetch(MCP_BASE + '/token', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({
          grant_type: 'authorization_code',
          code,
          redirect_uri: currentRedirectUri(),
          client_id: oauthState.client_id,
          code_verifier: oauthState.code_verifier
        })
      });
      if (!response.ok) throw new Error('Token exchange failed: ' + await response.text());
      return response.json();
    }

    async function initAuth() {
      loadAuthState();
      const params = new URLSearchParams(window.location.search);
      const code = params.get('code');
      const stateToken = params.get('state');
      if (code && stateToken) {
        window.history.replaceState({}, '', window.location.pathname);
        if (stateToken !== oauthState.state) {
          clearAuthState();
          updateAuthUi();
          setStatus('OAuth state mismatch. Please sign in again.', 'error');
          return false;
        }
        try {
          setStatus('Exchanging OAuth authorization code.', 'warn');
          const token = await exchangeCode(code);
          oauthState.access_token = token.access_token;
          oauthState.token_type = token.token_type || 'Bearer';
          if (token.expires_in) oauthState.expires_at = Date.now() + Number(token.expires_in) * 1000;
          delete oauthState.code_verifier;
          delete oauthState.state;
          saveAuthState();
        } catch (err) {
          clearAuthState();
          updateAuthUi();
          setStatus('OAuth callback failed: ' + err.message, 'error');
          return false;
        }
      }
      if (oauthState.expires_at && Date.now() >= oauthState.expires_at) {
        clearAuthState();
      }
      updateAuthUi();
      return hasValidAccessToken();
    }

    function logout() {
      clearAuthState();
      updateAuthUi();
      setStatus('Logged out. Sign in again to rerun the dashboard.', '');
    }

`;

  text = replaceOnce(text, "    function setStatus(text, tone = '') {", helperBlock + "    function setStatus(text, tone = '') {", 'variant B helper insert', file);

  text = replaceRegex(
    text,
    /async function executeBrowserQuery\(\{ label, role, sql, runId \}\) \{[\s\S]*?\n    \}\n\n    function primaryPeak/,
    `async function executeBrowserQuery({ label, role, sql, runId }) {
      const ledgerId = registerLedgerEntry(label, role, sql);

      try {
        if (!hasValidAccessToken()) {
          updateAuthUi();
          throw new Error('OAuth login is required.');
        }
        const endpoint = MCP_BASE + '/openapi/execute_query?query=' + encodeURIComponent(sql);
        const response = await fetch(endpoint, {
          headers: {
            Authorization: 'Bearer ' + oauthState.access_token
          }
        });
        if (!response.ok) {
          const text = await response.text();
          if (response.status === 401 || response.status === 403) {
            clearAuthState();
            updateAuthUi();
            throw new Error('Authentication expired. Please login again.');
          }
          throw new Error(text || \`HTTP \${response.status}\`);
        }
        const payload = await response.json();
        const rows = convertResult(payload);
        if (runId !== state.activeRunId) {
          return null;
        }
        updateLedgerEntry(ledgerId, {
          status: 'OK',
          rows: String(payload.count ?? rows.length),
          sql
        });
        return { payload, rows };
      } catch (error) {
        updateLedgerEntry(ledgerId, {
          status: 'Failed',
          rows: '0',
          sql
        });
        throw error;
      }
    }

    function primaryPeak`,
    'variant B executeBrowserQuery',
    file
  );

  text = text.replace(/localStorage\.setItem\(STORAGE_KEY, token\);\s*/g, '');

  text = replaceRegex(
    text,
    /function initFromStorage\(\) \{[\s\S]*?\n    \}\n\n    function wireEvents\(\) \{/,
    `function initFromStorage() {
    }

    function wireEvents() {`,
    'variant B initFromStorage',
    file
  );

  text = replaceRegex(
    text,
    /els\.forgetBtn\.addEventListener\('click', \(\) => \{\s*localStorage\.removeItem\(STORAGE_KEY\);\s*els\.tokenInput\.value = '';\s*setStatus\('Stored token removed\.', ''\);\s*\}\);/,
    "els.forgetBtn.addEventListener('click', logout);",
    'variant B forget handler',
    file
  );

  text = replaceRegex(
    text,
    /setDefaultSqlEditors\(\);\s*initFromStorage\(\);\s*wireEvents\(\);\s*\n\s*if \(getToken\(\)\) \{\s*runAll\(\);\s*\}\s*/m,
    `async function init() {
      setDefaultSqlEditors();
      initFromStorage();
      wireEvents();
      const isAuthenticated = await initAuth();
      if (isAuthenticated) {
        setStatus('Authenticated. Running the dashboard queries.', 'ok');
        runAll();
      }
    }

    init();
`,
    'variant B init tail A',
    file
  );

  text = replaceRegex(
    text,
    /setDefaultSqlEditors\(\);\s*wireEvents\(\);\s*\n\s*if \(hasValidAccessToken\(\)\) \{\s*runAll\(\);\s*\}\s*/m,
    `async function init() {
      setDefaultSqlEditors();
      wireEvents();
      const isAuthenticated = await initAuth();
      if (isAuthenticated) {
        setStatus('Authenticated. Running the dashboard queries.', 'ok');
        runAll();
      }
    }

    init();
`,
    'variant B init tail B',
    file
  );

  return text;
}

function migrateFile(file) {
  const source = fs.readFileSync(file, 'utf8');
  if (source.includes('OAUTH_STORAGE_KEY') && source.includes('registerOAuthClient')) {
    return { file, status: 'skipped', variant: 'already_oauth' };
  }

  let output;
  let variant;
  if (source.includes('$jwe') || source.includes('id="jwe-input"')) {
    output = migrateVariantA(source, file);
    variant = 'variant_a';
  } else if (source.includes('els.tokenInput') || source.includes('id="tokenInput"') || source.includes('id="token-input"')) {
    output = migrateVariantB(source, file);
    variant = 'variant_b';
  } else {
    return { file, status: 'skipped', variant: 'unmatched' };
  }

  if (output === source) {
    return { file, status: 'unchanged', variant };
  }

  fs.writeFileSync(file, output);
  return { file, status: 'updated', variant };
}

function main() {
  const files = walk(RUNS_DIR);
  const results = files.map((file) => {
    try {
      return migrateFile(file);
    } catch (error) {
      return { file, status: 'error', error: error.message };
    }
  });

  for (const result of results) {
    if (result.status === 'updated') console.log(`UPDATED  ${result.variant}  ${result.file}`);
    else if (result.status === 'skipped') console.log(`SKIPPED  ${result.variant}  ${result.file}`);
    else if (result.status === 'unchanged') console.log(`UNCHANGED ${result.variant}  ${result.file}`);
    else console.log(`ERROR    ${result.file}\n  ${result.error}`);
  }

  const summary = results.reduce((acc, result) => {
    acc[result.status] = (acc[result.status] || 0) + 1;
    return acc;
  }, {});
  console.log('\nSummary:', JSON.stringify(summary, null, 2));

  if (summary.error) process.exitCode = 1;
}

main();
