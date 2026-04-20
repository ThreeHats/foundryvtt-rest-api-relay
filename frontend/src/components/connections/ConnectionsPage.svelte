<script lang="ts">
  import { onMount } from 'svelte';
  import {
    fetchKnownClients, deleteKnownClient, generatePairingCode,
    setKnownClientAutoStart, setKnownClientCredential, fetchCredentials,
    fetchConnectionTokens, deleteConnectionToken,
    fetchConnectionLogs, renameConnectionToken,
    updateKnownClientCrossWorldSettings, updateConnectionTokenPermissions,
    fetchWorldNotificationSettings, updateWorldNotificationSettings, testWorldNotification,
  } from '../../lib/api';
  import type { KnownClient, KnownClientToken, ConnectionToken, Credential, ConnectionLog, KnownClientNotificationSettings } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';
  import { headlessEnabled } from '../../lib/config';

  // ── Data ──────────────────────────────────────────────────────────────────
  let clients     = $state<KnownClient[]>([]);
  let tokens      = $state<ConnectionToken[]>([]);
  let credentials = $state<Credential[]>([]);
  let logs        = $state<ConnectionLog[]>([]);
  let logsTotal   = $state(0);
  let logsOffset  = $state(0);
  const logsLimit = 20;

  let loading     = $state(true);
  let logsLoading = $state(false);
  let message     = $state('');
  let messageType = $state<'success' | 'error'>('error');

  // ── Pairing code (new world) ──────────────────────────────────────────────
  let pairingCode      = $state<string | null>(null);
  let pairingExpiresAt = $state<string | null>(null);
  let pairingCountdown = $state('');
  let pairingForWorld  = $state<KnownClient | null>(null); // null = new-world code
  let countdownInterval: ReturnType<typeof setInterval> | null = null;
  let pollInterval: ReturnType<typeof setInterval> | null = null;
  let refreshInterval: ReturnType<typeof setInterval> | null = null;
  let preCodeClientsCount = 0;
  let preCodeTokensCount = 0;

  // ── Inline token rename ───────────────────────────────────────────────────
  let editingTokenId   = $state<number | null>(null);
  let editingTokenName = $state('');

  function startEditToken(token: KnownClientToken) {
    editingTokenId = token.id;
    editingTokenName = token.name;
  }

  function cancelEditToken() {
    editingTokenId = null;
    editingTokenName = '';
  }

  async function saveTokenName(tokenId: number) {
    const newName = editingTokenName.trim();
    editingTokenId = null;
    editingTokenName = '';
    if (!newName) return;
    const r = await renameConnectionToken(tokenId, newName);
    if (r.ok) {
      tokens = tokens.map(t => t.id === tokenId ? { ...t, name: newName } : t);
    } else { message = r.error; messageType = 'error'; }
  }

  // ── Confirm modal ─────────────────────────────────────────────────────────
  let modal = $state<{ open: boolean; title: string; message: string; confirmLabel: string; action: () => Promise<void> }>({
    open: false, title: '', message: '', confirmLabel: 'Confirm', action: async () => {},
  });

  function ask(title: string, message: string, confirmLabel: string, action: () => Promise<void>) {
    modal = { open: true, title, message, confirmLabel, action };
  }

  async function runModal() {
    modal.open = false;
    await modal.action();
  }

  // ── Browser token selection (bulk delete) ─────────────────────────────────
  let selectedTokenIds = $state(new Set<number>());

  function toggleTokenSelect(id: number) {
    const next = new Set(selectedTokenIds);
    if (next.has(id)) next.delete(id); else next.add(id);
    selectedTokenIds = next;
  }

  function handleBulkDeleteTokens() {
    const n = selectedTokenIds.size;
    ask(
      `Delete ${n} browser${n === 1 ? '' : 's'}`,
      `Permanently delete ${n} browser connection${n === 1 ? '' : 's'}? This cannot be undone.`,
      'Delete all',
      async () => {
        await Promise.all([...selectedTokenIds].map(id => deleteConnectionToken(id)));
        selectedTokenIds = new Set();
        await Promise.all([loadTokens(), loadClients()]);
      },
    );
  }

  // ── Inactive browser collapse ─────────────────────────────────────────────
  let expandedInactive = $state(new Set<number>());

  function toggleInactive(clientId: number) {
    const next = new Set(expandedInactive);
    if (next.has(clientId)) next.delete(clientId); else next.add(clientId);
    expandedInactive = next;
  }

  // ── World notification settings (per-world, lazy-loaded) ─────────────────
  let worldNotifSettings = $state(new Map<number, KnownClientNotificationSettings>());
  let worldNotifOpen     = $state(new Set<number>());
  let worldNotifSaving   = $state(new Set<number>());
  let worldNotifTesting  = $state(new Set<number>());
  let worldNotifMsg      = $state(new Map<number, { text: string; type: 'success' | 'error' }>());

  async function toggleWorldNotif(id: number) {
    const next = new Set(worldNotifOpen);
    if (next.has(id)) {
      next.delete(id);
    } else {
      next.add(id);
      if (!worldNotifSettings.has(id)) {
        const r = await fetchWorldNotificationSettings(id);
        if (r.ok) {
          const m = new Map(worldNotifSettings);
          m.set(id, r.data);
          worldNotifSettings = m;
        }
      }
    }
    worldNotifOpen = next;
  }

  function updateWorldNotifField<K extends keyof KnownClientNotificationSettings>(id: number, field: K, value: KnownClientNotificationSettings[K]) {
    const current = worldNotifSettings.get(id);
    if (!current) return;
    const m = new Map(worldNotifSettings);
    m.set(id, { ...current, [field]: value });
    worldNotifSettings = m;
  }

  async function saveWorldNotif(id: number) {
    const ns = worldNotifSettings.get(id);
    if (!ns) return;
    const saving = new Set(worldNotifSaving); saving.add(id); worldNotifSaving = saving;
    const r = await updateWorldNotificationSettings(id, ns);
    const done = new Set(worldNotifSaving); done.delete(id); worldNotifSaving = done;
    const msgs = new Map(worldNotifMsg);
    if (r.ok) {
      const m = new Map(worldNotifSettings); m.set(id, r.data); worldNotifSettings = m;
      msgs.set(id, { text: 'Saved.', type: 'success' });
    } else {
      msgs.set(id, { text: r.error, type: 'error' });
    }
    worldNotifMsg = msgs;
    setTimeout(() => { const m = new Map(worldNotifMsg); m.delete(id); worldNotifMsg = m; }, 3000);
  }

  async function sendWorldNotifTest(id: number) {
    const testing = new Set(worldNotifTesting); testing.add(id); worldNotifTesting = testing;
    const r = await testWorldNotification(id);
    const done = new Set(worldNotifTesting); done.delete(id); worldNotifTesting = done;
    const msgs = new Map(worldNotifMsg);
    msgs.set(id, r.ok ? { text: 'Test sent!', type: 'success' } : { text: r.error, type: 'error' });
    worldNotifMsg = msgs;
    setTimeout(() => { const m = new Map(worldNotifMsg); m.delete(id); worldNotifMsg = m; }, 3000);
  }

  // ── Token IP allowlist modal ──────────────────────────────────────────────
  let editIPToken    = $state<KnownClientToken | null>(null);
  let editIPValue    = $state('');
  let editIPLoading  = $state(false);
  let editIPError    = $state('');

  function openIPModal(token: KnownClientToken) {
    editIPToken   = token;
    editIPValue   = token.allowedIps ?? '';
    editIPLoading = false;
    editIPError   = '';
  }

  async function handleSaveIPAllowlist() {
    if (!editIPToken) return;
    editIPLoading = true;
    editIPError   = '';
    const r = await updateConnectionTokenPermissions(editIPToken.id, { allowedIps: editIPValue });
    editIPLoading = false;
    if (r.ok) {
      // Update the token in the local client list so the row reflects the change immediately.
      clients = clients.map(c => ({
        ...c,
        tokens: c.tokens.map(t =>
          t.id === editIPToken!.id ? { ...t, allowedIps: editIPValue } : t
        ),
      }));
      editIPToken = null;
    } else {
      editIPError = r.error ?? 'Failed to save';
    }
  }

  // ── World cross-world settings modal ─────────────────────────────────────
  let editWorldClient = $state<KnownClient | null>(null);
  let editTargets     = $state<string[]>([]);  // selected target clientIds
  let editScopes      = $state<string[]>([]);
  let editRateLimit   = $state(0);             // 0 = unlimited
  let editLoading     = $state(false);

  // Scopes that can be enforced via cross-world remote-request tunneling.
  // These are the unique values from ActionToScopeRequired in scopes.go —
  // any action not in that map is denied by default regardless of this list.
  const ALL_SCOPES = [
    'entity:read', 'entity:write',
    'roll:read', 'roll:execute',
    'chat:read', 'chat:write',
    'encounter:read', 'encounter:manage',
    'macro:list', 'macro:execute', 'macro:write',
    'scene:read', 'scene:write',
    'user:read', 'user:write',
    'file:read', 'file:write',
    'structure:read', 'structure:write',
    'clients:read',
    'sheet:read',
    'playlist:control',
    'world:info',
    'search',
    'execute-js',
  ];

  onMount(() => {
    loadAll();
    refreshInterval = setInterval(() => {
      loadClients();
      loadTokens();
    }, 10_000);
    return () => {
      if (countdownInterval) clearInterval(countdownInterval);
      if (pollInterval) clearInterval(pollInterval);
      if (refreshInterval) clearInterval(refreshInterval);
    };
  });

  async function loadAll() {
    loading = true;
    const promises: Promise<void>[] = [loadClients(), loadTokens(), loadLogs()];
    if ($headlessEnabled) promises.push(loadCredentials());
    await Promise.all(promises);
    loading = false;
  }

  async function loadClients() {
    const r = await fetchKnownClients();
    if (r.ok) clients = r.data.clients || [];
  }

  async function loadTokens() {
    const r = await fetchConnectionTokens();
    if (r.ok) tokens = r.data.tokens || [];
  }

  async function loadCredentials() {
    const r = await fetchCredentials();
    if (r.ok) credentials = r.data.credentials || [];
  }

  async function loadLogs() {
    logsLoading = true;
    const r = await fetchConnectionLogs(logsLimit, logsOffset);
    logsLoading = false;
    if (r.ok) { logs = r.data.logs || []; logsTotal = r.data.total || 0; }
  }

  // ── Token grouping ────────────────────────────────────────────────────────
  // Tokens for a specific world: match by clientId, or by activeTokenId for
  // legacy tokens that were created before clientId was stored on the token.
  function tokensForClient(clientId: string, activeTokenId: number): ConnectionToken[] {
    return tokens.filter(t =>
      t.clientId === clientId ||
      (t.id === activeTokenId && activeTokenId !== 0)
    );
  }


  // ── Pairing ────────────────────────────────────────────────────────────────
  async function handleNewWorldCode() {
    message = '';
    const r = await generatePairingCode();
    if (!r.ok) { message = r.error; messageType = 'error'; return; }
    pairingForWorld = null;
    pairingCode = r.data.code;
    pairingExpiresAt = r.data.expiresAt;
    startCountdown();
    startCodePolling();
  }

  async function handleAddBrowser(client: KnownClient) {
    message = '';
    const r = await generatePairingCode({ clientId: client.clientId });
    if (!r.ok) { message = r.error; messageType = 'error'; return; }
    pairingForWorld = client;
    pairingCode = r.data.code;
    pairingExpiresAt = r.data.expiresAt;
    startCountdown();
    startCodePolling();
  }

  function startCountdown() {
    if (countdownInterval) clearInterval(countdownInterval);
    updateCountdown();
    countdownInterval = setInterval(updateCountdown, 1000);
  }

  function startCodePolling() {
    preCodeClientsCount = clients.length;
    preCodeTokensCount = tokens.length;
    if (pollInterval) clearInterval(pollInterval);
    pollInterval = setInterval(async () => {
      await Promise.all([loadClients(), loadTokens()]);
      if (clients.length > preCodeClientsCount || tokens.length > preCodeTokensCount) {
        dismissCode();
      }
    }, 3000);
  }

  function stopCodePolling() {
    if (pollInterval) { clearInterval(pollInterval); pollInterval = null; }
  }

  function updateCountdown() {
    if (!pairingExpiresAt) return;
    const remaining = new Date(pairingExpiresAt).getTime() - Date.now();
    if (remaining <= 0) {
      pairingCountdown = 'Expired';
      pairingCode = null;
      if (countdownInterval) clearInterval(countdownInterval);
      stopCodePolling();
      return;
    }
    const m = Math.floor(remaining / 60000);
    const s = Math.floor((remaining % 60000) / 1000);
    pairingCountdown = `${m}:${s.toString().padStart(2, '0')}`;
  }

  async function copyCode() {
    if (!pairingCode) return;
    await navigator.clipboard.writeText(pairingCode);
    message = 'Code copied!'; messageType = 'success';
    setTimeout(() => { message = ''; }, 2000);
  }

  function dismissCode() {
    pairingCode = null;
    pairingForWorld = null;
    if (countdownInterval) clearInterval(countdownInterval);
    stopCodePolling();
  }

  // ── Known client actions ──────────────────────────────────────────────────
  function handleDeleteClient(client: KnownClient) {
    const label = client.customName || client.worldTitle || client.clientId;
    ask(
      `Remove world`,
      `Remove "${label}"? This will also delete all paired browsers and disconnect any active session.`,
      'Remove',
      async () => {
        const r = await deleteKnownClient(client.id);
        if (r.ok) { loadClients(); } else { message = r.error; messageType = 'error'; }
      },
    );
  }

  async function handleCredentialChange(client: KnownClient, e: Event) {
    const select = e.target as HTMLSelectElement;
    const value = select.value === '' ? null : Number(select.value);
    const r = await setKnownClientCredential(client.id, value);
    if (r.ok) {
      clients = clients.map(c => c.id === client.id ? { ...c, credentialId: value } : c);
    } else {
      message = r.error; messageType = 'error';
      select.value = client.credentialId === null ? '' : String(client.credentialId);
    }
  }

  async function handleToggleAutoStart(client: KnownClient) {
    const next = !client.autoStartOnRemoteRequest;
    const r = await setKnownClientAutoStart(client.id, next);
    if (r.ok) {
      clients = clients.map(c => c.id === client.id ? { ...c, autoStartOnRemoteRequest: next } : c);
    } else { message = r.error; messageType = 'error'; }
  }

  function autoStartDisabled(client: KnownClient): boolean {
    return credentials.length === 0 || client.credentialId === null;
  }

  function autoStartTitle(client: KnownClient): string {
    if (credentials.length === 0) return 'Add a stored credential to enable auto-start headless sessions';
    if (client.credentialId === null) return 'Select a specific credential above to enable auto-start';
    return 'Allow the relay to spawn a headless session for this world when a sibling client sends a remote-request while this world is offline';
  }

  // ── Token actions ─────────────────────────────────────────────────────────
  function handleDeleteToken(token: KnownClientToken) {
    ask(
      'Delete browser',
      `Delete "${token.name}"? This cannot be undone.`,
      'Delete',
      async () => {
        const r = await deleteConnectionToken(token.id);
        if (r.ok) {
          selectedTokenIds.delete(token.id);
          await Promise.all([loadTokens(), loadClients()]);
        } else { message = r.error; messageType = 'error'; }
      },
    );
  }

  function openCrossWorldModal(client: KnownClient) {
    editWorldClient = client;
    editTargets = [...(client.allowedTargetClients || [])];
    editScopes = [...(client.remoteScopes || [])];
    editRateLimit = client.remoteRequestsPerHour ?? 0;
    editLoading = false;
  }

  function toggleScope(scope: string) {
    editScopes = editScopes.includes(scope)
      ? editScopes.filter(s => s !== scope)
      : [...editScopes, scope];
  }

  function toggleTarget(clientId: string) {
    editTargets = editTargets.includes(clientId)
      ? editTargets.filter(id => id !== clientId)
      : [...editTargets, clientId];
  }

  async function handleSaveCrossWorld() {
    if (!editWorldClient) return;
    editLoading = true;
    const r = await updateKnownClientCrossWorldSettings(editWorldClient.id, {
      allowedTargetClients: editTargets,
      remoteScopes: editScopes,
      remoteRequestsPerHour: editRateLimit,
    });
    editLoading = false;
    if (r.ok) {
      clients = clients.map(c => c.id === editWorldClient!.id ? {
        ...c,
        allowedTargetClients: [...editTargets],
        remoteScopes: [...editScopes],
        remoteRequestsPerHour: editRateLimit,
      } : c);
      editWorldClient = null;
    } else { message = r.error; messageType = 'error'; }
  }

  // ── Audit log pagination ──────────────────────────────────────────────────
  function nextPage() {
    if (logsOffset + logsLimit < logsTotal) { logsOffset += logsLimit; loadLogs(); }
  }
  function prevPage() {
    if (logsOffset > 0) { logsOffset = Math.max(0, logsOffset - logsLimit); loadLogs(); }
  }

  function formatLastSeen(v: string | null) {
    return v ? new Date(v).toLocaleString() : 'Never';
  }
</script>

<ConfirmModal
  open={modal.open}
  title={modal.title}
  message={modal.message}
  confirmLabel={modal.confirmLabel}
  dangerous={true}
  onConfirm={runModal}
  onCancel={() => modal.open = false}
/>

<h2 class="page-title">Connections</h2>

{#if message}
  <div class="alert mb-2" class:alert-success={messageType === 'success'} class:alert-error={messageType === 'error'}>
    {message}
  </div>
{/if}

<!-- ── Pairing Code ──────────────────────────────────────────────────────── -->
<div class="card">
  <h3 class="card-title mb-1">Pair a Foundry World</h3>
  <p class="text-muted mb-2" style="font-size: 0.85rem;">
    Generate a one-time 6-character code. Enter it in the Foundry REST API module's Connection dialog to link a world (and this browser) to your relay account.
  </p>

  {#if pairingCode}
    <p class="text-muted mb-1" style="font-size: 0.85rem;">
      {#if pairingForWorld}
        Code for adding a browser to <strong>{pairingForWorld.customName || pairingForWorld.worldTitle}</strong>:
      {:else}
        Code for a <strong>new world</strong>:
      {/if}
    </p>
    <div class="pairing-display">
      <code class="pairing-code">{pairingCode}</code>
      <div class="pairing-meta">
        <span class="text-muted" style="font-size: 0.8rem;">Expires in {pairingCountdown}</span>
        <div class="flex gap-1">
          <button class="btn btn-sm btn-primary" onclick={copyCode}>Copy</button>
          <button class="btn btn-sm btn-secondary" onclick={dismissCode}>Dismiss</button>
        </div>
      </div>
    </div>
  {:else}
    <button class="btn btn-primary" onclick={handleNewWorldCode}>Generate Pairing Code</button>
  {/if}
</div>

<!-- ── Worlds ────────────────────────────────────────────────────────────── -->
<div class="card">
  <div class="card-header">
    <h3 class="card-title">Worlds</h3>
    <div class="flex gap-1 items-center">
      {#if selectedTokenIds.size > 0}
        <button class="btn btn-sm btn-danger" onclick={handleBulkDeleteTokens}>
          Delete {selectedTokenIds.size} browser{selectedTokenIds.size === 1 ? '' : 's'}
        </button>
      {/if}
      <button class="btn btn-sm btn-secondary" onclick={loadAll} disabled={loading}>&#8635; Refresh</button>
    </div>
  </div>

  {#if loading}
    <p class="text-muted">Loading...</p>
  {:else if clients.length === 0}
    <p class="text-muted" style="font-size: 0.875rem;">No worlds paired yet. Generate a pairing code above and use it in your Foundry module.</p>
  {:else}
    {#each clients as client (client.id)}
      {@const allTokens = client.tokens || []}
      {@const activeTokens = allTokens.filter(t => t.id === client.activeTokenId && client.activeTokenId !== 0)}
      {@const inactiveTokens = allTokens.filter(t => !(t.id === client.activeTokenId && client.activeTokenId !== 0))}
      <div class="world-block">
        <div class="world-row">
          <!-- World info -->
          <div class="world-info">
            <div class="world-name">
              {#if client.isOnline}
                <span class="status-dot online" title="Online"></span>
              {:else}
                <span class="status-dot offline" title="Offline"></span>
              {/if}
              {#if client.customName}
                <strong>{client.customName}</strong>
                <span class="text-muted" style="font-size: 0.8rem;">({client.worldTitle})</span>
              {:else}
                <strong>{client.worldTitle || 'Unknown'}</strong>
              {/if}
              <span class="badge {client.isOnline ? 'badge-active' : 'badge-disabled'}" style="font-size: 0.7rem;">{client.isOnline ? 'online' : 'offline'}</span>
            </div>
            <div class="world-meta text-muted">
              {client.systemTitle || client.systemId}{client.systemVersion ? ` v${client.systemVersion}` : ''} · Foundry {client.foundryVersion || '?'} · Last seen: {formatLastSeen(client.lastSeenAt)}
            </div>
            <div class="world-meta text-muted">
              <code class="inline-code" style="font-size: 0.7rem;">{client.clientId}</code>
            </div>
          </div>
          <!-- World actions -->
          <div class="world-actions">
            {#if $headlessEnabled}
              {#if credentials.length > 0}
                <label class="text-muted" style="font-size: 0.75rem;">
                  Credential:
                  <select onchange={(e) => handleCredentialChange(client, e)} style="font-size: 0.75rem;">
                    <option value="" selected={client.credentialId === null}>(none)</option>
                    {#each credentials as cred}
                      <option value={cred.id} selected={client.credentialId === cred.id}>{cred.name}</option>
                    {/each}
                  </select>
                </label>
              {/if}
              <label class="toggle-label text-muted" style="font-size: 0.75rem;" title={autoStartTitle(client)}>
                <input type="checkbox" checked={client.autoStartOnRemoteRequest}
                  onchange={() => handleToggleAutoStart(client)}
                  disabled={autoStartDisabled(client)} />
                Auto-start
              </label>
            {/if}
            <button class="btn btn-sm btn-secondary" onclick={() => handleAddBrowser(client)} title="Generate a code so another GM's browser can pair to this world">
              + Add Browser
            </button>
            <button class="btn btn-sm btn-secondary" onclick={() => openCrossWorldModal(client)} title="Configure cross-world tunneling for this world">
              Cross-world{#if client.allowedTargetClients?.length} ({client.allowedTargetClients.length}){/if}
            </button>
            <button class="btn btn-sm btn-secondary" onclick={() => toggleWorldNotif(client.id)} title="Per-world notification settings">
              Notifications{worldNotifOpen.has(client.id) ? ' ▾' : ' ▸'}
            </button>
            <button class="btn btn-sm btn-danger" onclick={() => handleDeleteClient(client)}>Remove</button>
          </div>
        </div>

        <!-- Per-world notification settings -->
        {#if worldNotifOpen.has(client.id)}
          {@const ns = worldNotifSettings.get(client.id)}
          {@const msg = worldNotifMsg.get(client.id)}
          <div class="world-notif-panel">
            {#if ns}
              <div class="world-notif-form">
                <div class="form-group-row">
                  <div class="form-group flex-1">
                    <label class="form-label">Discord Webhook URL</label>
                    <input class="form-input" type="url" placeholder="https://discord.com/api/webhooks/..."
                      value={ns.discordWebhookUrl}
                      oninput={(e) => updateWorldNotifField(client.id, 'discordWebhookUrl', (e.target as HTMLInputElement).value)} />
                  </div>
                  {#if ns.smtpAvailable}
                    <div class="form-group flex-1">
                      <label class="form-label">Notify Email</label>
                      <input class="form-input" type="email" placeholder="you@example.com"
                        value={ns.notifyEmail}
                        oninput={(e) => updateWorldNotifField(client.id, 'notifyEmail', (e.target as HTMLInputElement).value)} />
                    </div>
                  {/if}
                </div>
                <div class="notif-toggles">
                  <label class="toggle-label">
                    <input type="checkbox" checked={ns.notifyOnConnect}
                      onchange={(e) => updateWorldNotifField(client.id, 'notifyOnConnect', (e.target as HTMLInputElement).checked)} />
                    Connect
                  </label>
                  <label class="toggle-label">
                    <input type="checkbox" checked={ns.notifyOnDisconnect}
                      onchange={(e) => updateWorldNotifField(client.id, 'notifyOnDisconnect', (e.target as HTMLInputElement).checked)} />
                    Disconnect
                  </label>
                  <label class="toggle-label">
                    <input type="checkbox" checked={ns.notifyOnExecuteJs}
                      onchange={(e) => updateWorldNotifField(client.id, 'notifyOnExecuteJs', (e.target as HTMLInputElement).checked)} />
                    Execute JS
                  </label>
                  <label class="toggle-label">
                    <input type="checkbox" checked={ns.notifyOnMacroExecute}
                      onchange={(e) => updateWorldNotifField(client.id, 'notifyOnMacroExecute', (e.target as HTMLInputElement).checked)} />
                    Macro Execute
                  </label>
                </div>
                <div class="flex gap-1 items-center mt-1">
                  <button class="btn btn-sm btn-primary" onclick={() => saveWorldNotif(client.id)} disabled={worldNotifSaving.has(client.id)}>
                    {worldNotifSaving.has(client.id) ? 'Saving…' : 'Save'}
                  </button>
                  <button class="btn btn-sm btn-secondary" onclick={() => sendWorldNotifTest(client.id)} disabled={worldNotifTesting.has(client.id)}>
                    {worldNotifTesting.has(client.id) ? 'Sending…' : 'Test'}
                  </button>
                  {#if msg}
                    <span class="notif-msg" class:notif-msg-success={msg.type === 'success'} class:notif-msg-error={msg.type === 'error'}>{msg.text}</span>
                  {/if}
                </div>
              </div>
            {:else}
              <p class="text-muted" style="font-size: 0.8rem;">Loading…</p>
            {/if}
          </div>
        {/if}

        <!-- Browsers for this world -->

        {#snippet tokenRow(token: KnownClientToken)}
          <div class="token-row" class:token-row-active={token.id === client.activeTokenId} class:token-row-selected={selectedTokenIds.has(token.id)}>
            <div class="token-info">
              <input type="checkbox" class="token-checkbox" checked={selectedTokenIds.has(token.id)} onchange={() => toggleTokenSelect(token.id)} />
              <span class="token-status-dot" class:token-dot-online={token.id === client.activeTokenId} title={token.id === client.activeTokenId ? 'Connected' : 'Not connected'}></span>
              {#if editingTokenId === token.id}
                <input class="token-name-input" bind:value={editingTokenName}
                  onkeydown={(e) => { if (e.key === 'Enter') saveTokenName(token.id); if (e.key === 'Escape') cancelEditToken(); }}
                  autofocus />
                <button class="btn btn-xs btn-primary" onclick={() => saveTokenName(token.id)}>✓</button>
                <button class="btn btn-xs btn-secondary" onclick={cancelEditToken}>✗</button>
              {:else}
                <span class="token-name-text" onclick={() => startEditToken(token)} title="Click to rename">{token.name}</span>
              {/if}
              {#if token.source === 'headless'}
                <span class="badge badge-disabled" style="font-size: 0.65rem;">headless</span>
              {/if}
              {#if token.allowedIps}
                <span class="badge badge-disabled" style="font-size: 0.65rem;" title="IP allowlist: {token.allowedIps}">🔒 IP restricted</span>
              {/if}
              <span class="text-muted" style="font-size: 0.7rem;">Last used: {token.lastUsedAt ? new Date(token.lastUsedAt).toLocaleString() : 'Never'}</span>
            </div>
            <div class="token-actions">
              <button class="btn btn-sm btn-secondary" onclick={() => openIPModal(token)}>IP Allowlist</button>
              <button class="btn btn-sm btn-danger" onclick={() => handleDeleteToken(token)}>Delete</button>
            </div>
          </div>
        {/snippet}

        {#each activeTokens as token (token.id)}
          {@render tokenRow(token)}
        {/each}

        {#if inactiveTokens.length > 0}
          <button class="inactive-toggle" onclick={() => toggleInactive(client.id)}>
            {inactiveTokens.length} inactive {inactiveTokens.length === 1 ? 'browser' : 'browsers'}
            {expandedInactive.has(client.id) ? ' ▾' : ' ▸'}
          </button>
          {#if expandedInactive.has(client.id)}
            {#each inactiveTokens as token (token.id)}
              {@render tokenRow(token)}
            {/each}
          {/if}
        {/if}
      </div>
    {/each}
  {/if}
</div>

<!-- ── Connection Audit Log ────────────────────────────────────────────────── -->
<div class="card">
  <h3 class="card-title mb-1">Connection Audit Log</h3>
  {#if logsLoading}
    <p class="text-muted">Loading logs...</p>
  {:else if logs.length === 0}
    <p class="text-muted" style="font-size: 0.875rem;">No connection logs yet.</p>
  {:else}
    <div class="table-wrapper">
      <table class="table">
        <thead>
          <tr><th>Timestamp</th><th>Browser</th><th>IP Address</th><th>Status</th></tr>
        </thead>
        <tbody>
          {#each logs as log (log.id)}
            <tr>
              <td>{log.createdAt ? new Date(log.createdAt).toLocaleString() : '—'}</td>
              <td>
                {#if log.tokenName}
                  <span style="font-size: 0.85rem;">{log.tokenName}</span>
                {:else}
                  <code class="inline-code" style="font-size: 0.75rem;">{log.clientId}</code>
                {/if}
              </td>
              <td><code class="inline-code">{log.ipAddress || '—'}</code></td>
              <td>
                {#if log.flagged}
                  <span class="badge badge-expired" title={log.flagReason || ''}>⚠ {log.flagReason || 'flagged'}</span>
                {:else}
                  <span class="badge badge-active">ok</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
    <div class="flex items-center justify-between mt-1">
      <span class="text-muted" style="font-size: 0.8rem;">Showing {logsOffset + 1}–{logsOffset + logs.length} of {logsTotal}</span>
      <div class="flex gap-1">
        <button class="btn btn-sm btn-secondary" onclick={prevPage} disabled={logsOffset === 0}>Previous</button>
        <button class="btn btn-sm btn-secondary" onclick={nextPage} disabled={logs.length < logsLimit}>Next</button>
      </div>
    </div>
  {/if}
</div>

<!-- ── IP Allowlist Modal (per token) ────────────────────────────────────── -->
{#if editIPToken}
  <div class="modal-backdrop" role="dialog" aria-modal="true">
    <div class="modal">
      <h3>IP Allowlist — {editIPToken.name}</h3>
      <p class="text-muted" style="font-size: 0.85rem; margin-bottom: 0.75rem;">
        Restrict this browser's relay connection to specific IP addresses or CIDR ranges. Leave blank to allow connections from any IP.
      </p>
      <label class="form-label">
        Allowed IPs / CIDRs
        <input class="input" type="text" bind:value={editIPValue}
          placeholder="e.g. 1.2.3.4, 10.0.0.0/8"
          onkeydown={(e) => { if (e.key === 'Enter') handleSaveIPAllowlist(); if (e.key === 'Escape') editIPToken = null; }} />
        <span class="text-muted" style="font-size: 0.75rem;">Comma-separated. Accepts individual IPs and CIDR ranges. Clear to remove the restriction.</span>
      </label>
      {#if editIPError}
        <p class="text-error" style="font-size: 0.85rem; margin-top: 0.5rem;">{editIPError}</p>
      {/if}
      <div class="modal-actions">
        <button class="btn btn-secondary" onclick={() => editIPToken = null} disabled={editIPLoading}>Cancel</button>
        <button class="btn btn-primary" onclick={handleSaveIPAllowlist} disabled={editIPLoading}>
          {editIPLoading ? 'Saving…' : 'Save'}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- ── Cross-world Settings Modal (per world) ────────────────────────────── -->
{#if editWorldClient}
  <div class="modal-backdrop" role="dialog" aria-modal="true">
    <div class="modal">
      <h3>Cross-world settings — {editWorldClient.customName || editWorldClient.worldTitle}</h3>
      <p class="text-muted" style="font-size: 0.85rem; margin-bottom: 0.75rem;">
        These settings apply to all browsers paired to this world. Allowed target worlds and scopes control which remote-request operations this world may invoke on other worlds.
      </p>

      <fieldset class="form-fieldset">
        <legend>Allowed scopes</legend>
        <p class="text-muted" style="font-size: 0.8rem;">Operations this world may invoke on allowed target worlds via remote-request.</p>
        <div class="scope-grid">
          {#each ALL_SCOPES as scope}
            <label class="scope-label" class:dangerous={scope === 'execute-js' || scope === 'macro:execute' || scope === 'macro:write'}>
              <input type="checkbox" checked={editScopes.includes(scope)} onchange={() => toggleScope(scope)} />
              {scope}
              {#if scope === 'execute-js' || scope === 'macro:execute' || scope === 'macro:write'}
                <span class="badge badge-expired" style="margin-left: 0.25rem; font-size: 0.65rem;">dangerous</span>
              {/if}
            </label>
          {/each}
        </div>
      </fieldset>

      <fieldset class="form-fieldset">
        <legend>Allowed target worlds</legend>
        {#if clients.filter(c => c.clientId !== editWorldClient?.clientId).length === 0}
          <p class="text-muted" style="font-size: 0.8rem;">No other worlds paired yet. Pair another world first to enable cross-world access.</p>
        {:else}
          <p class="text-muted" style="font-size: 0.8rem;">This world may send remote-request operations to the selected worlds.</p>
          <div class="scope-grid">
            {#each clients.filter(c => c.clientId !== editWorldClient?.clientId) as target}
              <label class="scope-label">
                <input type="checkbox"
                  checked={editTargets.includes(target.clientId)}
                  onchange={() => toggleTarget(target.clientId)} />
                <span>
                  {target.customName || target.worldTitle || target.clientId}
                  {#if target.isOnline}<span class="text-muted" style="font-size: 0.7rem;"> (online)</span>{/if}
                </span>
              </label>
            {/each}
          </div>
        {/if}
      </fieldset>

      <label class="form-label">
        Rate limit (per hour)
        <input class="input" type="number" min="0" bind:value={editRateLimit} placeholder="0 = unlimited" />
        <span class="text-muted" style="font-size: 0.75rem;">Max cross-world requests per hour from this world. 0 = no limit.</span>
      </label>

      <div class="modal-actions">
        <button class="btn btn-secondary" onclick={() => editWorldClient = null} disabled={editLoading}>Cancel</button>
        <button class="btn btn-primary" onclick={handleSaveCrossWorld} disabled={editLoading}>
          {editLoading ? 'Saving…' : 'Save'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .world-block {
    border: 1px solid var(--color-border, #e5e7eb);
    border-radius: var(--radius-md);
    margin-bottom: 0.75rem;
    overflow: hidden;
  }
  .world-row {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.75rem 1rem;
    background: var(--color-bg-raised);
  }
  .world-info { flex: 1; min-width: 0; }
  .world-name {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    flex-wrap: wrap;
    font-size: 0.95rem;
    margin-bottom: 0.2rem;
  }
  .world-meta { font-size: 0.8rem; margin-bottom: 0.1rem; }
  .world-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
    flex-shrink: 0;
  }
  .token-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.4rem 1rem 0.4rem 1.75rem;
    border-top: 1px solid var(--color-border, #e5e7eb);
    background: var(--color-bg-sunken);
    font-size: 0.85rem;
  }
  .token-info {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
    flex: 1;
  }
  .token-actions { display: flex; gap: 0.4rem; flex-shrink: 0; }
  .inactive-toggle {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 0.3rem 1rem 0.3rem 1.75rem;
    border: none;
    border-top: 1px solid var(--color-border, #e5e7eb);
    background: var(--color-bg-sunken);
    color: var(--color-text-muted);
    font-size: 0.8rem;
    cursor: pointer;
    text-align: left;
  }
  .inactive-toggle:hover { color: var(--color-text); }
  .token-row-active { background: color-mix(in srgb, var(--color-success, #22c55e) 8%, var(--color-bg-sunken, #f9fafb)); }
  .token-row-selected { background: color-mix(in srgb, var(--color-primary, #4f46e5) 8%, var(--color-bg-sunken, #f9fafb)); }
  .token-checkbox { flex-shrink: 0; }
  .token-status-dot {
    width: 7px; height: 7px; border-radius: 50%;
    background: var(--color-text-muted, #9ca3af); flex-shrink: 0;
  }
  .token-dot-online {
    background: var(--color-success, #22c55e);
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--color-success, #22c55e) 30%, transparent);
  }
  .token-name-text {
    cursor: pointer;
    border-bottom: 1px dashed transparent;
    transition: border-color 0.15s;
  }
  .token-name-text:hover { border-bottom-color: var(--color-text-muted, #999); }
  .token-name-input {
    font-size: 0.85rem; padding: 0.1rem 0.3rem;
    border: 1px solid var(--color-primary, #4f46e5);
    border-radius: 3px; width: 160px;
  }
  .btn-xs { padding: 0.1rem 0.35rem; font-size: 0.75rem; }
  .status-dot {
    width: 8px; height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }
  .status-dot.online  { background: var(--color-success, #22c55e); }
  .status-dot.offline { background: var(--color-text-muted, #9ca3af); }
  .table-wrapper { overflow-x: auto; }
  .toggle-label {
    display: inline-flex; align-items: center; gap: 0.25rem;
    cursor: pointer;
  }
  .pairing-display {
    display: flex; align-items: center; gap: 1rem; flex-wrap: wrap;
  }
  .pairing-code {
    font-family: var(--font-mono);
    font-size: 2rem; font-weight: 700; letter-spacing: 0.2em;
    background: var(--color-bg-sunken); padding: 0.5rem 1rem;
    border-radius: var(--radius-md); border: 2px solid var(--color-primary);
  }
  .pairing-meta { display: flex; flex-direction: column; gap: 0.25rem; }
  .modal-backdrop {
    position: fixed; inset: 0; background: rgba(0,0,0,0.6);
    display: flex; align-items: center; justify-content: center;
    z-index: 1000; padding: 1rem;
  }
  .modal {
    background: var(--color-bg, #fff); border-radius: var(--radius-md);
    padding: 1.5rem; max-width: 560px; width: 100%;
    box-shadow: 0 10px 40px rgba(0,0,0,0.3);
    display: flex; flex-direction: column; gap: 1rem;
  }
  .form-label { display: flex; flex-direction: column; gap: 0.3rem; font-size: 0.875rem; }
  .form-fieldset { border: 1px solid var(--color-border, #e5e7eb); border-radius: var(--radius-sm); padding: 0.75rem; }
  .form-fieldset legend { font-size: 0.85rem; font-weight: 600; padding: 0 0.25rem; }
  .scope-grid {
    display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 0.4rem; margin-top: 0.5rem;
  }
  .scope-label { display: flex; align-items: center; gap: 0.4rem; font-size: 0.8rem; cursor: pointer; }
  .scope-label.dangerous { color: var(--color-error); }
  .modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 0.5rem; }
  .world-notif-panel {
    border-top: 1px solid var(--color-border, #e5e7eb);
    background: var(--color-bg-sunken);
    padding: 0.75rem 1rem;
  }
  .world-notif-form { display: flex; flex-direction: column; gap: 0.6rem; }
  .form-group-row { display: flex; gap: 0.75rem; flex-wrap: wrap; }
  .form-group { display: flex; flex-direction: column; gap: 0.25rem; }
  .form-group .form-label { font-size: 0.8rem; font-weight: 600; color: var(--color-text-muted); }
  .form-group .form-input {
    font-size: 0.85rem; padding: 0.3rem 0.5rem;
    border: 1px solid var(--color-border, #e5e7eb);
    border-radius: var(--radius-sm); background: var(--color-bg);
    color: var(--color-text);
  }
  .flex-1 { flex: 1; min-width: 180px; }
  .notif-toggles { display: flex; gap: 1rem; flex-wrap: wrap; align-items: center; }
  .notif-msg { font-size: 0.8rem; }
  .notif-msg-success { color: var(--color-success, #22c55e); }
  .notif-msg-error { color: var(--color-error, #ef4444); }
</style>
