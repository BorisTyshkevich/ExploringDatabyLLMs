#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const files = [
  '/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-25/q001_hops_per_day/claude/sonnet/run-002/visual.html',
  '/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-25/q001_hops_per_day/claude/sonnet/run-004/visual.html',
  '/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-25/q001_hops_per_day/claude/sonnet/run-005/visual.html',
  '/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-26/q001_hops_per_day/chatgpt/gpt-5.4/run-002/visual.html',
  '/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-26/q001_hops_per_day/chatgpt/gpt-5.4/run-004/visual.html',
  '/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-26/q001_hops_per_day/claude/opus/run-002/visual.html',
  '/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-26/q001_hops_per_day/claude/sonnet/run-001/visual.html',
  '/Users/bvt/work/ExploringDatabyLLMs/runs/2026-03-26/q001_hops_per_day/codex/gpt-5.4/run-001/visual.html'
];

const block = String.raw`
const __oauthRecentMigration = (() => {
  if (window.__oauthRecentMigrationLoaded) return null;
  window.__oauthRecentMigrationLoaded = true;

  const MCP_BASE = 'https://mcp.demo.altinity.cloud';
  const OAUTH_STORAGE_KEY = 'OnTimeAnalystDashboard::auth::oauth2';
  const LEGACY_JWE_KEY = 'OnTimeAnalystDashboard::auth::jwe';
  let oauthState = {};

  try { localStorage.removeItem(LEGACY_JWE_KEY); } catch (err) {}

  function q(selector) {
    return document.querySelector(selector);
  }

  function tokenInput() {
    return q('#jwe-input, #tok, #tokenInput, #jweToken, #jwe-inp');
  }

  function runButton() {
    return q('#run-btn, #runBtn');
  }

  function forgetButton() {
    return q('#forget-btn, #btn-forget, #forgetBtn, #forgetToken');
  }

  function statusNodes() {
    return [
      q('#status-msg'),
      q('#status-line'),
      q('#status-area'),
      q('#runStatusText'),
      q('#globalStatus'),
      q('.status-line'),
      q('.status-msg')
    ].filter(Boolean);
  }

  function setVisibleStatus(message) {
    statusNodes().forEach((node) => {
      node.textContent = message;
    });
  }

  function loadAuthState() {
    try {
      oauthState = JSON.parse(localStorage.getItem(OAUTH_STORAGE_KEY)) || {};
    } catch (err) {
      oauthState = {};
    }
  }

  function saveAuthState() {
    localStorage.setItem(OAUTH_STORAGE_KEY, JSON.stringify(oauthState));
  }

  function clearAuthState() {
    localStorage.removeItem(OAUTH_STORAGE_KEY);
    oauthState = {};
    try { localStorage.removeItem(LEGACY_JWE_KEY); } catch (err) {}
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
    return btoa(String.fromCharCode(...new Uint8Array(hash)))
      .replace(/\+/g, '-')
      .replace(/\//g, '_')
      .replace(/=+$/, '');
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

  function updateTextCopy() {
    const candidates = [
      q('#empty-state p'),
      q('#empty-state .empty-text'),
      q('#empty-state .status-msg'),
      q('#empty-state .small-note')
    ].filter(Boolean);
    candidates.forEach((node) => {
      node.innerHTML = 'Use Google sign-in from the footer controls, then run the dashboard.';
    });

    const titles = [
      q('#empty-state h3'),
      q('#empty-state .empty-title')
    ].filter(Boolean);
    titles.forEach((node) => {
      node.textContent = 'Sign in to load the dashboard';
    });

    document.querySelectorAll('label, .ctrl-lbl, .f-label').forEach((node) => {
      if (/JWE/i.test(node.textContent)) node.textContent = 'OAuth2 redirect flow';
    });

    document.querySelectorAll('p, div, span').forEach((node) => {
      if (!node.firstElementChild && /JWE token|stored token|live fetch/i.test(node.textContent)) {
        node.textContent = node.textContent
          .replace(/JWE token/gi, 'OAuth2 login')
          .replace(/stored token/gi, 'OAuth session')
          .replace(/live fetch/gi, 'authenticated fetch');
      }
    });
  }

  function applyInputState(authenticated) {
    const input = tokenInput();
    if (!input) return;
    input.value = authenticated ? '__oauth_authenticated__' : '';
    try {
      if (authenticated) input.dataset.rawToken = '__oauth_authenticated__';
      else delete input.dataset.rawToken;
    } catch (err) {}
    input.setAttribute('type', 'text');
    input.setAttribute('readonly', 'readonly');
    input.setAttribute('placeholder', 'OAuth2 handled automatically');
    input.style.display = 'none';

    const label = document.querySelector('label[for="' + input.id + '"]');
    if (label) label.style.display = 'none';
  }

  function ensureLoginButton() {
    let btn = q('#oauth-login-btn');
    if (btn) return btn;
    const run = runButton();
    if (!run || !run.parentNode) return null;
    btn = document.createElement('button');
    btn.id = 'oauth-login-btn';
    btn.type = 'button';
    btn.textContent = 'Login with Google';
    btn.className = run.className || '';
    if (!btn.className.includes('secondary') && !btn.className.includes('ghost')) {
      btn.className = (btn.className + ' secondary').trim();
    }
    run.parentNode.insertBefore(btn, run);
    btn.addEventListener('click', startLogin);
    return btn;
  }

  function applyUi() {
    updateTextCopy();
    applyInputState(hasValidAccessToken());

    const login = ensureLoginButton();
    const forget = forgetButton();
    if (forget) {
      forget.textContent = 'Logout';
      forget.style.display = hasValidAccessToken() ? '' : 'none';
    }
    if (login) {
      login.style.display = hasValidAccessToken() ? 'none' : '';
    }
  }

  function triggerRun() {
    if (!hasValidAccessToken()) return;
    if (typeof window.runAll === 'function') return window.runAll();
    if (typeof window.startRun === 'function') return window.startRun();
    if (typeof window.runDashboard === 'function') return window.runDashboard({ auto: true });
    if (typeof window.runQuery === 'function' && window.runQuery.length === 0) return window.runQuery();
    const run = runButton();
    if (run) run.click();
  }

  async function startLogin() {
    try {
      setVisibleStatus('Registering OAuth client and redirecting to Google sign-in.');
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
      setVisibleStatus('Login failed: ' + err.message);
      applyUi();
    }
  }

  function logout() {
    clearAuthState();
    applyUi();
    setVisibleStatus('Logged out. Sign in again to rerun the dashboard.');
  }

  function wrapFetch() {
    if (window.__oauthRecentFetchWrapped) return;
    const nativeFetch = window.fetch.bind(window);
    window.fetch = async function patchedFetch(input, init) {
      let url = typeof input === 'string' ? input : (input && input.url) || '';
      let normalizedUrl = url;

      const prefixedMatch = url.match(/^https:\/\/mcp\.demo\.altinity\.cloud\/[^/]+\/openapi\/execute_query\?query=(.*)$/);
      if (prefixedMatch) {
        normalizedUrl = MCP_BASE + '/openapi/execute_query?query=' + prefixedMatch[1];
      }

      const directQuery = normalizedUrl.startsWith(MCP_BASE + '/openapi/execute_query?query=');
      if (directQuery) {
        if (!hasValidAccessToken()) {
          applyUi();
          throw new Error('OAuth login is required.');
        }
        const headers = new Headers((init && init.headers) || (input instanceof Request ? input.headers : undefined) || {});
        headers.set('Authorization', 'Bearer ' + oauthState.access_token);
        if (typeof input === 'string') {
          input = normalizedUrl;
          init = Object.assign({}, init || {}, { headers });
        } else {
          input = new Request(normalizedUrl, Object.assign({}, input, { headers }));
          init = undefined;
        }
      }

      const response = await nativeFetch(input, init);
      if (directQuery && (response.status === 401 || response.status === 403)) {
        clearAuthState();
        applyUi();
        setVisibleStatus('Authentication expired. Sign in again to run queries.');
      }
      return response;
    };
    window.__oauthRecentFetchWrapped = true;
  }

  async function initOAuthRecent() {
    wrapFetch();
    loadAuthState();

    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');
    const stateToken = params.get('state');
    if (code && stateToken) {
      window.history.replaceState({}, '', window.location.pathname);
      if (stateToken !== oauthState.state) {
        clearAuthState();
        applyUi();
        setVisibleStatus('OAuth state mismatch. Please sign in again.');
        return;
      }
      try {
        setVisibleStatus('Exchanging OAuth authorization code.');
        const token = await exchangeCode(code);
        oauthState.access_token = token.access_token;
        oauthState.token_type = token.token_type || 'Bearer';
        if (token.expires_in) oauthState.expires_at = Date.now() + Number(token.expires_in) * 1000;
        delete oauthState.code_verifier;
        delete oauthState.state;
        saveAuthState();
      } catch (err) {
        clearAuthState();
        applyUi();
        setVisibleStatus('OAuth callback failed: ' + err.message);
        return;
      }
    }

    if (oauthState.expires_at && Date.now() >= oauthState.expires_at) {
      clearAuthState();
    }

    applyUi();

    const forget = forgetButton();
    if (forget && !forget.dataset.oauthBound) {
      forget.addEventListener('click', () => logout());
      forget.dataset.oauthBound = '1';
    }

    if (hasValidAccessToken()) {
      setVisibleStatus('Authenticated with OAuth2. Running the dashboard queries.');
      setTimeout(() => triggerRun(), 50);
    } else {
      setVisibleStatus('Sign in with Google to run this dashboard. The page URL must be registered as an OAuth redirect URI.');
    }
  }

  queueMicrotask(() => {
    initOAuthRecent();
  });

  return { startLogin, logout };
})();
`;

for (const file of files) {
  let text = fs.readFileSync(file, 'utf8');
  if (text.includes('__oauthRecentMigration')) continue;
  const idx = text.indexOf('<script>');
  if (idx === -1) throw new Error(`No inline <script> tag found in ${file}`);
  text = text.slice(0, idx + '<script>'.length) + '\n' + block + '\n' + text.slice(idx + '<script>'.length);
  fs.writeFileSync(file, text);
  console.log(`Injected OAuth wrapper into ${file}`);
}
