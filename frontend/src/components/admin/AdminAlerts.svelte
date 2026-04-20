<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AlertEvent } from '../../lib/types';
  import { headlessEnabled } from '../../lib/config';

  let events = $state<AlertEvent[]>([]);
  let loadError = $state('');

  // Config state
  let discordUrl = $state('');
  let emailDest = $state('');
  let enabledDiscord = $state<Record<string, boolean>>({});
  let enabledEmail = $state<Record<string, boolean>>({});
  let saving = $state(false);
  let saveMsg = $state('');
  let saveMsgType = $state<'success' | 'error'>('success');
  let testResult = $state<Record<string, string>>({});

  let refreshInterval: ReturnType<typeof setInterval> | null = null;

  async function loadConfig() {
    try {
      const r = await adminApi.getAlertConfig();
      discordUrl = r.discordWebhookUrl ?? '';
      emailDest = r.emailDestination ?? '';
      enabledDiscord = {};
      enabledEmail = {};
      for (const s of r.subscriptions ?? []) {
        if (s.channel === 'discord') enabledDiscord[s.alertType] = true;
        if (s.channel === 'email') enabledEmail[s.alertType] = true;
      }
    } catch (e: any) {
      loadError = e?.message ?? 'Failed to load config';
    }
  }

  async function loadEvents() {
    try {
      const r = await adminApi.getRecentAlerts();
      events = (r.events ?? []).slice().reverse();
    } catch {}
  }

  onMount(() => {
    loadConfig();
    loadEvents();
    refreshInterval = setInterval(loadEvents, 15_000);
  });

  onDestroy(() => {
    if (refreshInterval) clearInterval(refreshInterval);
  });

  async function saveConfig() {
    saving = true;
    saveMsg = '';
    const subscriptions: Array<{ alertType: string; channel: 'discord' | 'email' }> = [];
    for (const t of allTypes) {
      if (enabledDiscord[t]) subscriptions.push({ alertType: t, channel: 'discord' });
      if (enabledEmail[t]) subscriptions.push({ alertType: t, channel: 'email' });
    }
    try {
      await adminApi.saveAlertConfig({ discordWebhookUrl: discordUrl, emailDestination: emailDest, subscriptions });
      saveMsg = 'Settings saved.';
      saveMsgType = 'success';
    } catch (e: any) {
      saveMsg = e?.message ?? 'Save failed';
      saveMsgType = 'error';
    } finally {
      saving = false;
      setTimeout(() => { saveMsg = ''; }, 4000);
    }
  }

  async function testChannel(channel: 'discord' | 'email') {
    testResult = { ...testResult, [channel]: '' };
    try {
      await adminApi.testAlert(channel);
      testResult = { ...testResult, [channel]: 'ok' };
    } catch {
      testResult = { ...testResult, [channel]: 'err' };
    }
    setTimeout(() => { testResult = { ...testResult, [channel]: '' }; }, 4000);
  }

  function relativeTime(ts: string): string {
    const diff = Date.now() - new Date(ts).getTime();
    const s = Math.floor(diff / 1000);
    if (s < 60) return `${s}s ago`;
    const m = Math.floor(s / 60);
    if (m < 60) return `${m}m ago`;
    const h = Math.floor(m / 60);
    if (h < 24) return `${h}h ago`;
    return `${Math.floor(h / 24)}d ago`;
  }

  const typeLabels: Record<string, string> = {
    failed_auth_spike: 'Failed Auth Spike',
    flagged_connection: 'Flagged Connection (Metadata Mismatch)',
    world_id_cross_account: 'WorldId Used by Multiple Accounts',
    registration_spike: 'Registration Spike (Same IP)',
    rate_limit_burst: 'Rate Limit Burst',
    admin_login: 'Admin Login',
    duplicate_connection_spike: 'Duplicate WS Connection Attempts',
    metadata_mismatch_spike: 'Metadata Mismatch Spike',
    execute_js_burst: 'execute-js Burst',
    password_reset_flood: 'Password Reset Flood',
    invalid_token_spike: 'Invalid Token Spike',
    new_ip_for_token: 'New IP for Known Browser Token',
    account_password_change: 'Account Password Change',
    client_disconnect_spike: 'Client Disconnect Spike',
    system_unhealthy: 'System Unhealthy',
    headless_session_flood: 'Headless Session Flood',
    stripe_payment_failed: 'Stripe Payment Failed',
    new_user_registration: 'New User Registration',
    new_subscription: 'New Subscription',
    subscription_cancelled: 'Subscription Cancelled',
    user_monthly_limit_approaching: 'User Monthly Limit Approaching (80%)',
  };

  const securityTypes = [
    'failed_auth_spike', 'flagged_connection', 'world_id_cross_account', 'registration_spike',
    'rate_limit_burst', 'admin_login', 'duplicate_connection_spike', 'metadata_mismatch_spike',
    'execute_js_burst', 'password_reset_flood', 'invalid_token_spike', 'new_ip_for_token',
    'account_password_change',
  ];
  let operationsTypes = $derived([
    'client_disconnect_spike',
    'system_unhealthy',
    ...($headlessEnabled ? ['headless_session_flood'] : []),
    'stripe_payment_failed',
  ]);
  const analyticsTypes = [
    'new_user_registration', 'new_subscription', 'subscription_cancelled',
    'user_monthly_limit_approaching',
  ];

  let allTypes = $derived([...securityTypes, ...operationsTypes, ...analyticsTypes]);

  let expandedEvents = $state(new Set<number>());
  function toggleExpand(i: number) {
    const next = new Set(expandedEvents);
    if (next.has(i)) next.delete(i); else next.add(i);
    expandedEvents = next;
  }
</script>

<div class="admin-page">
  <h1>Alerts</h1>
  {#if loadError}<p class="error">{loadError}</p>{/if}

  <!-- ── Notification Settings ── -->
  <section>
    <h2>Notification Settings</h2>
    <div class="dest-inputs">
      <div class="dest-row">
        <label class="dest-label" for="discord-url">Discord Webhook URL</label>
        <input
          class="dest-input"
          id="discord-url"
          type="url"
          bind:value={discordUrl}
          placeholder="https://discord.com/api/webhooks/…"
        />
        <button
          class="btn-test"
          class:ok={testResult['discord'] === 'ok'}
          class:err={testResult['discord'] === 'err'}
          onclick={() => testChannel('discord')}
          disabled={!discordUrl}
        >
          {testResult['discord'] === 'ok' ? 'Sent ✓' : testResult['discord'] === 'err' ? 'Failed ✗' : 'Test'}
        </button>
      </div>
      <div class="dest-row">
        <label class="dest-label" for="email-dest">Email Destination</label>
        <input
          class="dest-input"
          id="email-dest"
          type="email"
          bind:value={emailDest}
          placeholder="admin@example.com"
        />
        <button
          class="btn-test"
          class:ok={testResult['email'] === 'ok'}
          class:err={testResult['email'] === 'err'}
          onclick={() => testChannel('email')}
          disabled={!emailDest}
        >
          {testResult['email'] === 'ok' ? 'Sent ✓' : testResult['email'] === 'err' ? 'Failed ✗' : 'Test'}
        </button>
      </div>
    </div>

    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Alert Type</th>
            <th class="ch-col">Discord</th>
            <th class="ch-col">Email</th>
          </tr>
        </thead>
        <tbody>
          <!-- Security -->
          <tr class="section-header"><td colspan="3">Security</td></tr>
          {#each securityTypes as t (t)}
            <tr>
              <td>{typeLabels[t]}</td>
              <td class="ch-col"><input type="checkbox" bind:checked={enabledDiscord[t]} /></td>
              <td class="ch-col"><input type="checkbox" bind:checked={enabledEmail[t]} /></td>
            </tr>
          {/each}

          <!-- Operations -->
          <tr class="section-header"><td colspan="3">Operations</td></tr>
          {#each operationsTypes as t (t)}
            <tr>
              <td>{typeLabels[t]}</td>
              <td class="ch-col"><input type="checkbox" bind:checked={enabledDiscord[t]} /></td>
              <td class="ch-col"><input type="checkbox" bind:checked={enabledEmail[t]} /></td>
            </tr>
          {/each}

          <!-- Analytics -->
          <tr class="section-header"><td colspan="3">Analytics</td></tr>
          {#each analyticsTypes as t (t)}
            <tr>
              <td>{typeLabels[t]}</td>
              <td class="ch-col"><input type="checkbox" bind:checked={enabledDiscord[t]} /></td>
              <td class="ch-col"><input type="checkbox" bind:checked={enabledEmail[t]} /></td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <div class="save-row">
      <button class="btn-save" onclick={saveConfig} disabled={saving}>
        {saving ? 'Saving…' : 'Save Settings'}
      </button>
      {#if saveMsg}
        <span class="save-msg" class:success={saveMsgType === 'success'} class:error={saveMsgType === 'error'}>
          {saveMsg}
        </span>
      {/if}
    </div>
  </section>

  <!-- ── Recent Events ── -->
  <section>
    <h2>Recent Events <span class="sub">(auto-refreshes every 15s)</span></h2>
    {#if events.length === 0}
      <p class="muted">No events yet.</p>
    {:else}
      <ul class="event-list">
        {#each events as ev, i (i)}
          <li class="event-item">
            <div class="event-header" onclick={() => toggleExpand(i)}>
              <span class="badge sev-{ev.severity}">{ev.severity}</span>
              <span class="event-type">{typeLabels[ev.type] ?? ev.type}</span>
              <span class="event-msg">{ev.message}</span>
              <span class="event-time">{relativeTime(ev.timestamp)}</span>
              {#if ev.details && Object.keys(ev.details).length > 0}
                <button class="expand-btn">{expandedEvents.has(i) ? '▾' : '▸'}</button>
              {/if}
            </div>
            {#if expandedEvents.has(i) && ev.details}
              <dl class="event-details">
                {#each Object.entries(ev.details) as [k, v]}
                  <dt>{k}</dt><dd>{String(v)}</dd>
                {/each}
              </dl>
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</div>

<style>
  .admin-page { padding: 1.5rem; }
  h1 { margin-top: 0; }
  h2 { font-size: 1rem; font-weight: 600; margin-bottom: 0.5rem; }
  .sub { font-size: 0.75rem; font-weight: 400; color: var(--color-text-muted, #888); }
  .muted { color: var(--color-text-muted, #888); font-size: 0.875rem; }
  .error { color: var(--color-error, #e53935); font-size: 0.875rem; }
  .success { color: #43a047; }
  section { margin-bottom: 2rem; }

  /* Destination inputs */
  .dest-inputs { display: flex; flex-direction: column; gap: 0.5rem; margin-bottom: 1rem; }
  .dest-row { display: flex; align-items: center; gap: 0.5rem; }
  .dest-label { font-size: 0.875rem; font-weight: 500; min-width: 160px; white-space: nowrap; }
  .dest-input { flex: 1; padding: 0.35rem 0.5rem; border: 1px solid var(--color-border); border-radius: var(--radius-sm, 4px); background: var(--color-bg, white); font-size: 0.875rem; }
  .btn-test { padding: 0.3rem 0.6rem; font-size: 0.75rem; border: 1px solid var(--color-border); border-radius: var(--radius-sm, 4px); background: var(--color-bg-elevated); cursor: pointer; white-space: nowrap; }
  .btn-test:disabled { opacity: 0.5; cursor: default; }
  .btn-test.ok { background: #43a047; color: white; border-color: #43a047; }
  .btn-test.err { background: var(--color-error, #e53935); color: white; border-color: var(--color-error, #e53935); }

  /* Matrix table */
  .table-wrap { overflow-x: auto; border: 1px solid var(--color-border); border-radius: var(--radius-md); margin-bottom: 0.75rem; }
  table { width: 100%; border-collapse: collapse; }
  th, td { text-align: left; padding: 0.4rem 0.75rem; border-bottom: 1px solid var(--color-border-light); font-size: 0.875rem; }
  th { background: var(--color-bg-elevated); font-weight: 600; }
  .ch-col { text-align: center; width: 80px; }
  .ch-col input[type="checkbox"] { cursor: pointer; width: 15px; height: 15px; }
  tr.section-header td { background: var(--color-bg-elevated); font-weight: 600; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.04em; color: var(--color-text-muted, #888); padding: 0.3rem 0.75rem; }
  tr:last-child td { border-bottom: none; }

  /* Save */
  .save-row { display: flex; align-items: center; gap: 0.75rem; }
  .btn-save { padding: 0.45rem 1rem; font-size: 0.875rem; font-weight: 500; border: none; border-radius: var(--radius-sm, 4px); background: var(--color-primary, #1e88e5); color: white; cursor: pointer; }
  .btn-save:disabled { opacity: 0.6; cursor: default; }
  .save-msg { font-size: 0.875rem; }

  /* Event feed */
  .event-list { list-style: none; padding: 0; margin: 0; border: 1px solid var(--color-border); border-radius: var(--radius-md); overflow: hidden; }
  .event-item { border-bottom: 1px solid var(--color-border-light); }
  .event-item:last-child { border-bottom: none; }
  .event-header { display: flex; align-items: center; gap: 0.5rem; padding: 0.5rem 0.75rem; cursor: pointer; font-size: 0.875rem; }
  .event-header:hover { background: var(--color-bg-elevated); }
  .badge { padding: 0.1rem 0.4rem; border-radius: 3px; font-size: 0.7rem; font-weight: 700; text-transform: uppercase; white-space: nowrap; }
  .sev-critical { background: #e53935; color: white; }
  .sev-warning  { background: #fb8c00; color: white; }
  .sev-info     { background: #1e88e5; color: white; }
  .event-type { font-weight: 600; white-space: nowrap; }
  .event-msg { flex: 1; color: var(--color-text-muted, #666); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .event-time { white-space: nowrap; color: var(--color-text-muted, #888); font-size: 0.75rem; }
  .expand-btn { background: none; border: none; cursor: pointer; padding: 0 0.25rem; font-size: 0.8rem; }
  .event-details { display: grid; grid-template-columns: max-content 1fr; gap: 0.1rem 0.75rem; margin: 0; padding: 0.5rem 0.75rem 0.75rem 2.5rem; background: var(--color-bg-elevated); font-size: 0.8rem; }
  .event-details dt { font-weight: 600; color: var(--color-text-muted, #666); }
  .event-details dd { margin: 0; font-family: monospace; }
</style>
