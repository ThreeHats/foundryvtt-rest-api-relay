<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { adminFetch, adminGet, adminMutate } from '../../lib/adminAuth';

  // ---------------------------------------------------------------------------
  // Terminal / process runner state
  // Only UI-driving values are $state. Terminal output is managed directly on
  // the DOM element to avoid O(n²) re-renders as lines stream in.
  // ---------------------------------------------------------------------------
  let running = $state(false);
  let operationKind = $state<'test' | 'docs' | null>(null);
  let runId = $state<string | null>(null);
  let status = $state<string>('idle');
  let eventSource: EventSource | null = null;

  // Terminal DOM — plain let, not $state. We write to it directly.
  let terminalEl: HTMLElement | null = null;
  let pendingHtml: string[] = [];
  let flushTimer: ReturnType<typeof setTimeout> | null = null;
  let wasConnected = false; // tracks whether the current EventSource connected at least once

  onDestroy(() => {
    eventSource?.close();
    if (flushTimer) clearTimeout(flushTimer);
  });

  // ---------------------------------------------------------------------------
  // ANSI → HTML converter
  // Handles the SGR codes that Jest/pnpm emit: bold, dim, reverse, and the
  // standard 8-color + bright foreground palette (30-37, 39, 90-97).
  // ---------------------------------------------------------------------------
  const ANSI_FG: Record<number, string> = {
    30: '#555', 31: '#e06c75', 32: '#98c379', 33: '#e5c07b',
    34: '#61afef', 35: '#c678dd', 36: '#56b6c2', 37: '#abb2bf',
    90: '#666', 91: '#ff7b72', 92: '#7ee787', 93: '#ffa657',
    94: '#79c0ff', 95: '#d2a8ff', 96: '#76e3ea', 97: '#f0f6fc',
  };

  type AnsiState = { color: string; bold: boolean; dim: boolean; invert: boolean };
  const ANSI_RESET: AnsiState = { color: '', bold: false, dim: false, invert: false };

  function applyAnsiCode(n: number, s: AnsiState): AnsiState {
    s = { ...s };
    if (n === 0)  return { ...ANSI_RESET };
    if (n === 1)  s.bold = true;
    if (n === 2)  s.dim = true;
    if (n === 22) { s.bold = false; s.dim = false; }
    if (n === 7)  s.invert = true;
    if (n === 27) s.invert = false;
    if (n === 39) s.color = '';
    if (ANSI_FG[n] !== undefined) s.color = ANSI_FG[n];
    return s;
  }

  function stateToCSS(s: AnsiState): string {
    const p: string[] = [];
    if (s.color)  p.push(`color:${s.color}`);
    if (s.bold)   p.push('font-weight:700');
    if (s.dim)    p.push('opacity:0.55');
    if (s.invert) p.push('background:#d0d0d0;color:#1a1a1a;border-radius:2px;padding:0 2px');
    return p.join(';');
  }

  function ansiToHtml(raw: string): string {
    const esc = (s: string) =>
      s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');

    let cur: AnsiState = { ...ANSI_RESET };
    let html = '';
    let spanOpen = false;

    for (const chunk of raw.split(/(\u001b\[[0-9;]*m)/)) {
      const m = chunk.match(/^\u001b\[([0-9;]*)m$/);
      if (m) {
        const next = (m[1] || '0').split(';').reduce(
          (s, c) => applyAnsiCode(parseInt(c || '0', 10), s),
          cur
        );
        if (JSON.stringify(next) !== JSON.stringify(cur)) {
          if (spanOpen) { html += '</span>'; spanOpen = false; }
          cur = next;
          const css = stateToCSS(cur);
          if (css) { html += `<span style="${css}">`; spanOpen = true; }
        }
      } else if (chunk) {
        html += esc(chunk);
      }
    }
    if (spanOpen) html += '</span>';
    return html;
  }

  // ---------------------------------------------------------------------------
  // Terminal DOM helpers — O(1) per line, no reactive state involved
  // ---------------------------------------------------------------------------
  function clearTerminal() {
    if (terminalEl) terminalEl.innerHTML = '';
    pendingHtml = [];
    if (flushTimer) { clearTimeout(flushTimer); flushTimer = null; }
  }

  const MAX_TERMINAL_LINES = 5000;

  function flushPending() {
    flushTimer = null;
    if (!pendingHtml.length || !terminalEl) return;
    const html = pendingHtml.map(l => `<div>${l || '<br>'}</div>`).join('');
    terminalEl.insertAdjacentHTML('beforeend', html);
    pendingHtml = [];
    // Trim oldest lines once we exceed the cap to keep the DOM bounded
    while (terminalEl.childElementCount > MAX_TERMINAL_LINES) {
      terminalEl.removeChild(terminalEl.firstChild!);
    }
    terminalEl.scrollTop = terminalEl.scrollHeight;
  }

  function queueHtmlLine(html: string) {
    pendingHtml.push(html);
    if (!flushTimer) flushTimer = setTimeout(flushPending, 50);
  }

  // ---------------------------------------------------------------------------
  // Process controls
  // ---------------------------------------------------------------------------
  async function startOperation(kind: 'test' | 'docs') {
    clearTerminal();
    running = true;
    operationKind = kind;
    status = 'starting';
    const endpoint = kind === 'docs' ? '/admin/api/tests/docs/run' : '/admin/api/tests/run';
    try {
      const r = await adminMutate<{ runId: string; status: string }>('POST', endpoint, {});
      runId = r.runId;
      status = r.status;
      streamRun(r.runId);
    } catch (e: any) {
      queueHtmlLine(ansiToHtml(`Error: ${e?.message ?? e}`));
      flushPending();
      running = false;
      operationKind = null;
    }
  }

  function streamRun(id: string) {
    wasConnected = false;
    eventSource = new EventSource(`/admin/api/tests/run/${id}/stream`);

    const appendLine = (ev: MessageEvent, prefix = '') => {
      let text: string;
      try {
        const parsed = JSON.parse(ev.data);
        text = prefix + parsed.data;
      } catch {
        text = prefix + ev.data;
      }
      queueHtmlLine(ansiToHtml(text));
    };

    eventSource.addEventListener('stdout', (ev: MessageEvent) => appendLine(ev));
    eventSource.addEventListener('stderr', (ev: MessageEvent) => appendLine(ev));
    eventSource.addEventListener('error',  (ev: MessageEvent) => appendLine(ev, 'Error: '));

    eventSource.addEventListener('complete', (ev: MessageEvent) => {
      try { status = JSON.parse(ev.data).data; } catch { status = 'completed'; }
      flushPending();
      running = false;
      operationKind = null;
      eventSource?.close();
    });

    eventSource.onopen = () => {
      if (wasConnected) {
        queueHtmlLine(ansiToHtml('\u001b[32m[reconnected]\u001b[0m'));
        flushPending();
      }
      wasConnected = true;
    };

    eventSource.onerror = () => {
      if (!eventSource || eventSource.readyState === EventSource.CLOSED) {
        // Browser gave up (HTTP error, auth failure, etc.) — reset UI
        flushPending();
        queueHtmlLine(ansiToHtml('\u001b[31m[stream closed]\u001b[0m'));
        flushPending();
        running = false;
        operationKind = null;
      } else {
        // readyState === CONNECTING — browser is auto-retrying with Last-Event-ID
        queueHtmlLine(ansiToHtml('\u001b[33m[reconnecting...]\u001b[0m'));
        flushPending();
      }
    };
  }

  async function cancel() {
    if (!runId) return;
    await adminFetch(`/admin/api/tests/run/${runId}/cancel`, { method: 'POST' });
  }

  function openReport() {
    window.open('/admin/api/tests/report', '_blank');
  }

  // ---------------------------------------------------------------------------
  // Env var editor
  // ---------------------------------------------------------------------------
  type EnvVars = Record<string, string>;

  const defaultEnv: EnvVars = {
    TEST_BASE_URL: '', TEST_API_KEY: '', TEST_FOUNDRY_VERSIONS: '',
    TEST_DEFAULT_WORLD: '', USE_EXISTING_SESSION: '',
    FOUNDRY_USERNAME: '', FOUNDRY_PASSWORD: '',
    TEST_USER_EMAIL: '', TEST_USER_PASSWORD: '',
    TEST_ADMIN_EMAIL: '', TEST_ADMIN_PASSWORD: '',
    CAPTURE_BROWSER_CONSOLE: '',
  };

  let envVars = $state<EnvVars>({ ...defaultEnv });
  let envSaveStatus = $state<'idle' | 'saving' | 'saved' | 'error'>('idle');
  let envSaveMessage = $state('');
  // Track which password fields are revealed
  let showPassword = $state<Record<string, boolean>>({});

  onMount(async () => {
    try {
      const data = await adminGet<EnvVars>('/admin/api/tests/env');
      envVars = { ...defaultEnv, ...data };
    } catch {
      // If .env.test doesn't exist yet, leave defaults empty
    }
  });

  async function saveEnv() {
    envSaveStatus = 'saving';
    try {
      await adminMutate('POST', '/admin/api/tests/env', envVars);
      envSaveStatus = 'saved';
      envSaveMessage = 'Saved';
    } catch (e: any) {
      envSaveStatus = 'error';
      envSaveMessage = e?.message ?? 'Save failed';
    }
    setTimeout(() => { envSaveStatus = 'idle'; envSaveMessage = ''; }, 2500);
  }

  function toggleShow(key: string) {
    showPassword = { ...showPassword, [key]: !showPassword[key] };
  }

  function getEnvVal(key: string): string {
    return envVars[key] ?? '';
  }

  function setEnvVal(key: string, val: string) {
    envVars[key] = val;
  }

  // Static (non-version-specific) field definitions
  type FieldDef = { key: string; label: string; secret?: boolean; hint?: string; section?: string };

  const STATIC_FIELDS: FieldDef[] = [
    { key: 'TEST_BASE_URL',           label: 'Relay Server URL',     section: 'Relay' },
    { key: 'TEST_API_KEY',            label: 'API Key',              section: 'Relay', secret: true, hint: 'leave blank for ephemeral mode' },
    { key: 'TEST_FOUNDRY_VERSIONS',   label: 'Test Versions',        section: 'Foundry', hint: 'e.g. 12,13,14' },
    { key: 'TEST_DEFAULT_WORLD',      label: 'Default World',        section: 'Foundry', hint: 'fallback if version-specific not set' },
    { key: 'USE_EXISTING_SESSION',    label: 'Use Existing Session', section: 'Session', hint: 'true or false' },
    { key: 'FOUNDRY_USERNAME',        label: 'Foundry Username',     section: 'Session' },
    { key: 'FOUNDRY_PASSWORD',        label: 'Foundry Password',     section: 'Session', secret: true },
    { key: 'TEST_USER_EMAIL',         label: 'Test User Email',      section: 'Auth' },
    { key: 'TEST_USER_PASSWORD',      label: 'Test User Password',   section: 'Auth', secret: true },
    { key: 'TEST_ADMIN_EMAIL',        label: 'Admin Email',          section: 'Auth' },
    { key: 'TEST_ADMIN_PASSWORD',     label: 'Admin Password',       section: 'Auth', secret: true },
    { key: 'CAPTURE_BROWSER_CONSOLE', label: 'Browser Console Log',  section: 'Misc', hint: 'error / warn / debug' },
  ];

  // Per-version field templates — keys use V{N} placeholder filled at render time
  function versionFields(v: string): FieldDef[] {
    return [
      { key: `FOUNDRY_V${v}_URL`,         label: 'URL' },
      { key: `FOUNDRY_V${v}_WORLD`,        label: 'World' },
      { key: `TEST_CLIENT_ID_V${v}`,       label: 'Client ID',     hint: 'when USE_EXISTING_SESSION=true' },
      { key: `TEST_PLAYER_USER_ID_V${v}`,  label: 'Player User ID', hint: 'non-GM player for permission tests' },
    ];
  }

  // Reactive list of versions parsed from TEST_FOUNDRY_VERSIONS
  let parsedVersions = $derived(
    (envVars.TEST_FOUNDRY_VERSIONS ?? '')
      .split(',')
      .map(v => v.trim())
      .filter(Boolean)
  );
</script>

<div class="admin-page">
  <h1>Tests & Docs</h1>

  <!-- ── Action bar ── -->
  <div class="actions">
    <button onclick={() => startOperation('test')} disabled={running}>
      Run Full Test Suite
    </button>
    <button onclick={() => startOperation('docs')} disabled={running}>
      Regenerate Docs
    </button>
    {#if running}
      <button class="danger" onclick={cancel}>Cancel</button>
    {/if}
    <button class="secondary" onclick={openReport}>View Last Report</button>
    <span class="status">
      {#if operationKind === 'test'}
        Status: {status} <span class="kind-badge">tests</span>
      {:else if operationKind === 'docs'}
        Status: {status} <span class="kind-badge">docs</span>
      {:else}
        Status: {status}
      {/if}
    </span>
  </div>

  <!-- ── Terminal ── -->
  <div class="terminal" bind:this={terminalEl}></div>

  <!-- ── Test Environment ── -->
  <section class="env-section">
    <h2>Test Environment <span class="env-file">.env.test</span></h2>
    <p class="hint">Changes are saved to disk and take effect on the next test run.</p>

    <div class="env-grid">
      {#each STATIC_FIELDS as field, i}
        {#if field.section && (i === 0 || STATIC_FIELDS[i - 1].section !== field.section)}
          <div class="env-section-header">{field.section}</div><div></div>
        {/if}
        <label class="env-label" for={field.key}>{field.label}</label>
        <div class="env-input-wrap">
          <input
            id={field.key}
            type={field.secret && !showPassword[field.key] ? 'password' : 'text'}
            placeholder={field.hint ?? ''}
            value={getEnvVal(field.key)}
            oninput={(e) => setEnvVal(field.key, (e.target as HTMLInputElement).value)}
          />
          {#if field.secret}
            <button
              type="button"
              class="toggle-show"
              onclick={() => toggleShow(field.key)}
              title={showPassword[field.key] ? 'Hide' : 'Show'}
            >{showPassword[field.key] ? '🙈' : '👁'}</button>
          {/if}
        </div>
      {/each}

      {#each parsedVersions as v}
        <div class="env-section-header">Foundry v{v}</div><div></div>
        {#each versionFields(v) as field}
          <label class="env-label" for={field.key}>{field.label}</label>
          <div class="env-input-wrap">
            <input
              id={field.key}
              type="text"
              placeholder={field.hint ?? ''}
              value={getEnvVal(field.key)}
              oninput={(e) => setEnvVal(field.key, (e.target as HTMLInputElement).value)}
            />
          </div>
        {/each}
      {/each}
    </div>

    <div class="env-footer">
      <button onclick={saveEnv} disabled={envSaveStatus === 'saving'}>
        {envSaveStatus === 'saving' ? 'Saving…' : 'Save'}
      </button>
      {#if envSaveMessage}
        <span class="save-msg" class:error={envSaveStatus === 'error'}>{envSaveMessage}</span>
      {/if}
    </div>
  </section>
</div>

<style>
  .admin-page { padding: 1.5rem; max-width: 900px; }
  h1 { margin-bottom: 0.25rem; }
  h2 { margin: 0 0 0.25rem; font-size: 1rem; }

  /* ── Actions ── */
  .actions { display: flex; gap: 0.75rem; align-items: center; margin-bottom: 1rem; flex-wrap: wrap; }
  .status { color: var(--color-text-muted); margin-left: 0.25rem; }
  .kind-badge {
    display: inline-block;
    font-size: 0.7rem;
    padding: 0 5px;
    border-radius: 3px;
    background: var(--color-border, #444);
    color: var(--color-text-muted);
    vertical-align: middle;
    margin-left: 4px;
  }
  button.danger    { background: var(--color-error, #e53935); color: #fff; border: none; }
  button.secondary { background: transparent; border: 1px solid var(--color-border, #444); color: var(--color-text-muted); }
  button.secondary:hover { border-color: var(--color-text-muted); color: var(--color-text, #eee); }

  /* ── Terminal ── */
  .terminal {
    background: #1a1a1a;
    color: #d0d0d0;
    padding: 1rem 1.25rem;
    border-radius: var(--radius-md, 6px);
    font-family: 'SF Mono', 'Menlo', 'Consolas', 'Liberation Mono', monospace;
    font-size: 0.78rem;
    line-height: 1.55;
    min-height: 6rem;
    max-height: 60vh;
    overflow-y: auto;
    white-space: pre-wrap;
    word-break: break-all;
    margin-bottom: 2rem;
    scrollbar-width: thin;
    scrollbar-color: #444 #1a1a1a;
  }
  .terminal::-webkit-scrollbar        { width: 6px; }
  .terminal::-webkit-scrollbar-track  { background: #1a1a1a; }
  .terminal::-webkit-scrollbar-thumb  { background: #444; border-radius: 3px; }
  .terminal :global(span) { display: inline; }
  .terminal :global(div)  { margin: 0; }

  /* ── Env section ── */
  .env-section {
    border-top: 1px solid var(--color-border, #333);
    padding-top: 1.5rem;
  }
  .env-file {
    font-family: monospace;
    font-size: 0.78rem;
    color: var(--color-text-muted);
    font-weight: normal;
    margin-left: 0.4rem;
  }
  .hint { color: var(--color-text-muted); font-size: 0.82rem; margin: 0 0 1.25rem; }

  .env-grid {
    display: grid;
    grid-template-columns: 180px 1fr;
    gap: 0.5rem 1rem;
    align-items: center;
    margin-bottom: 1.25rem;
  }
  .env-section-header {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--color-text-muted);
    padding-top: 0.75rem;
    border-top: 1px solid var(--color-border, #333);
    margin-top: 0.25rem;
  }
  .env-section-header:first-child { border-top: none; padding-top: 0; }
  .env-label {
    font-size: 0.85rem;
    color: var(--color-text-muted);
    text-align: right;
    padding-right: 0.25rem;
  }
  .env-input-wrap {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }
  .env-input-wrap input {
    flex: 1;
    min-width: 0;
    font-family: 'SF Mono', 'Menlo', monospace;
    font-size: 0.82rem;
    padding: 0.35rem 0.6rem;
    background: var(--color-surface, #1f1f1f);
    border: 1px solid var(--color-border, #444);
    border-radius: var(--radius-sm, 4px);
    color: var(--color-text, #eee);
  }
  .env-input-wrap input:focus {
    outline: none;
    border-color: var(--color-primary, #61afef);
  }
  .toggle-show {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    padding: 0 0.2rem;
    opacity: 0.6;
    flex-shrink: 0;
  }
  .toggle-show:hover { opacity: 1; }

  .env-footer { display: flex; align-items: center; gap: 1rem; }
  .save-msg { font-size: 0.85rem; color: var(--color-success, #98c379); }
  .save-msg.error { color: var(--color-error, #e06c75); }
</style>
