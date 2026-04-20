<script lang="ts">
  import {
    createScopedKey,
    updateScopedKey,
    fetchKnownClients,
    fetchKnownUsers,
    fetchApiKeyNotificationSettings,
    updateApiKeyNotificationSettings,
    testApiKeyNotification,
  } from '../../lib/api';
  import type { ScopedKey, KnownClient, KnownUser } from '../../lib/types';
  import { headlessEnabled } from '../../lib/config';

  interface Props {
    editKey?: ScopedKey | null;
    onSave: (createdKeyValue?: string) => void;
    onCancel: () => void;
  }

  let { editKey = null, onSave, onCancel }: Props = $props();

  let name = $state(editKey?.name ?? '');
  let dailyLimit = $state(editKey?.dailyLimit?.toString() ?? '');
  let expiresAt = $state(editKey?.expiresAt ? new Date(editKey.expiresAt).toISOString().slice(0, 16) : '');
  let scopesOpen = $state(!editKey);
  let notificationsOpen = $state(false);

  // Per-key notification settings (edit mode only)
  let nsLoaded = $state(false);
  let nsDiscordWebhookUrl = $state('');
  let nsNotifyEmail = $state('');
  let nsNotifyOnExecuteJs = $state(false);
  let nsNotifyOnMacroExecute = $state(false);
  let nsNotifyOnRateLimit = $state(false);
  let nsNotifyOnError = $state(false);
  let nsSmtpAvailable = $state(false);
  let nsTesting = $state(false);
  let nsMessage = $state('');
  let nsMessageType = $state<'success' | 'error'>('error');

  const DEFAULT_SCOPES = new Set([
    'entity:read', 'roll:read', 'chat:read', 'encounter:read',
    'macro:list', 'scene:read', 'canvas:read', 'effects:read',
    'user:read', 'file:read', 'world:info', 'clients:read',
    'sheet:read', 'events:subscribe', 'search', 'structure:read', 'dnd5e',
  ]);

  let selectedScopes = $state<Set<string>>(new Set(editKey ? (editKey.scopes ?? []) : DEFAULT_SCOPES));
  let knownClients = $state<KnownClient[]>([]);
  let selectedClientIds = $state<Set<string>>(
    editKey?.scopedClientIds
      ? new Set(
          Array.isArray(editKey.scopedClientIds)
            ? editKey.scopedClientIds
            : (editKey.scopedClientIds as string).split(',').map((s: string) => s.trim()).filter(Boolean)
        )
      : new Set()
  );

  // Per-client user selection: clientId → userId (empty string = unrestricted)
  let selectedUserIds = $state<Record<string, string>>(
    editKey?.scopedUserIds
      ? { ...(editKey.scopedUserIds as Record<string, string>) }
      : {}
  );

  // Loaded KnownUsers per clientId
  let knownUsersByClientId = $state<Record<string, KnownUser[]>>({});

  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');
  let saving = $state(false);

  let isEdit = $derived(!!editKey);

  const scopeGroups: Record<string, string[]> = {
    'Entity': ['entity:read', 'entity:write'],
    'Roll': ['roll:read', 'roll:execute'],
    'Chat': ['chat:read', 'chat:write'],
    'Encounter': ['encounter:read', 'encounter:manage'],
    'Scene': ['scene:read', 'scene:write'],
    'Canvas': ['canvas:read', 'canvas:write'],
    'Effects': ['effects:read', 'effects:write'],
    'Macro': ['macro:list', 'macro:execute', 'macro:write'],
    'User': ['user:read', 'user:write'],
    'File System': ['file:read', 'file:write'],
    'Structure': ['structure:read', 'structure:write'],
    'Sheet': ['sheet:read'],
    'Playlist': ['playlist:control'],
    'World': ['world:info', 'clients:read'],
    'Events': ['events:subscribe'],
    'Search': ['search'],
    'Session': ['session:manage'],
    'D&D 5e': ['dnd5e'],
    'Advanced': ['execute-js'],
  };

  let visibleScopeGroups = $derived(
    Object.entries(scopeGroups).filter(([name]) => name !== 'Session' || $headlessEnabled)
  );

  function toggleScope(scope: string) {
    if (selectedScopes.has(scope)) {
      selectedScopes.delete(scope);
    } else {
      selectedScopes.add(scope);
    }
    selectedScopes = new Set(selectedScopes);
  }

  function selectAll() {
    selectedScopes = new Set(visibleScopeGroups.flatMap(([, scopes]) => scopes));
  }

  function resetDefaults() {
    selectedScopes = new Set(DEFAULT_SCOPES);
  }

  function clearAll() {
    selectedScopes = new Set();
  }

  async function loadUsersForClient(clientId: string) {
    const client = knownClients.find(c => c.clientId === clientId);
    if (!client) return;
    const r = await fetchKnownUsers(client.id);
    if (r.ok) {
      knownUsersByClientId = { ...knownUsersByClientId, [clientId]: r.data.users };
    }
  }

  async function toggleKnownClient(clientId: string) {
    if (selectedClientIds.has(clientId)) {
      selectedClientIds.delete(clientId);
      // Remove userId selection for this client
      const next = { ...selectedUserIds };
      delete next[clientId];
      selectedUserIds = next;
    } else {
      selectedClientIds.add(clientId);
      // Load users for this client if not already loaded
      if (!knownUsersByClientId[clientId]) {
        await loadUsersForClient(clientId);
      }
    }
    selectedClientIds = new Set(selectedClientIds);
  }

  // Load known clients on mount; also load users for pre-selected clients in edit mode
  $effect(() => {
    loadKnownClients();
    if (isEdit && editKey && !nsLoaded) {
      loadNotificationSettings();
    }
  });

  async function loadNotificationSettings() {
    if (!editKey) return;
    const result = await fetchApiKeyNotificationSettings(editKey.id);
    if (result.ok) {
      nsDiscordWebhookUrl = result.data.discordWebhookUrl || '';
      nsNotifyEmail = result.data.notifyEmail || '';
      nsNotifyOnExecuteJs = result.data.notifyOnExecuteJs;
      nsNotifyOnMacroExecute = result.data.notifyOnMacroExecute;
      nsNotifyOnRateLimit = result.data.notifyOnRateLimit;
      nsNotifyOnError = result.data.notifyOnError;
      nsSmtpAvailable = result.data.smtpAvailable;
      nsLoaded = true;
    }
  }

  async function handleTestNotification() {
    if (!editKey) return;
    nsTesting = true;
    nsMessage = '';
    await updateApiKeyNotificationSettings(editKey.id, {
      discordWebhookUrl: nsDiscordWebhookUrl.trim(),
      notifyEmail: nsNotifyEmail.trim(),
      notifyOnExecuteJs: nsNotifyOnExecuteJs,
      notifyOnMacroExecute: nsNotifyOnMacroExecute,
      notifyOnRateLimit: nsNotifyOnRateLimit,
      notifyOnError: nsNotifyOnError,
    });
    const result = await testApiKeyNotification(editKey.id);
    nsTesting = false;
    if (result.ok) {
      nsMessage = 'Test notification sent!';
      nsMessageType = 'success';
      setTimeout(() => { nsMessage = ''; }, 3000);
    } else {
      nsMessage = result.error;
      nsMessageType = 'error';
    }
  }

  async function loadKnownClients() {
    const result = await fetchKnownClients();
    if (result.ok) {
      knownClients = result.data.clients || [];
      // Pre-load users for clients that are already selected (edit mode)
      for (const clientId of selectedClientIds) {
        if (!knownUsersByClientId[clientId]) {
          loadUsersForClient(clientId);
        }
      }
    }
  }

  const roleNames: Record<number, string> = { 0: 'None', 1: 'Player', 2: 'Trusted', 3: 'Asst GM', 4: 'GM' };

  async function handleSubmit(e: Event) {
    e.preventDefault();

    if (!name.trim()) {
      message = 'Name is required.';
      messageType = 'error';
      return;
    }

    if (selectedScopes.size === 0) {
      message = 'At least one scope is required.';
      messageType = 'error';
      return;
    }

    saving = true;
    message = '';

    // Build scopedUserIds — only include entries where a userId was actually selected
    const scopedUserIdsMap = Object.fromEntries(
      Object.entries(selectedUserIds).filter(([, v]) => v !== '')
    );

    const body: Record<string, unknown> = {
      name: name.trim(),
      dailyLimit: dailyLimit !== '' ? parseInt(dailyLimit) : null,
      expiresAt: expiresAt || null,
      scopes: Array.from(selectedScopes),
      scopedClientIds: selectedClientIds.size > 0 ? Array.from(selectedClientIds) : null,
      scopedUserIds: Object.keys(scopedUserIdsMap).length > 0 ? scopedUserIdsMap : null,
    };

    const result = isEdit
      ? await updateScopedKey(editKey!.id, body)
      : await createScopedKey(body);

    if (result.ok && isEdit && editKey) {
      await updateApiKeyNotificationSettings(editKey.id, {
        discordWebhookUrl: nsDiscordWebhookUrl.trim(),
        notifyEmail: nsNotifyEmail.trim(),
        notifyOnExecuteJs: nsNotifyOnExecuteJs,
        notifyOnMacroExecute: nsNotifyOnMacroExecute,
        notifyOnRateLimit: nsNotifyOnRateLimit,
        notifyOnError: nsNotifyOnError,
      });
    }

    saving = false;

    if (result.ok) {
      const createdKey = !isEdit && 'key' in result.data ? (result.data as any).key : undefined;
      onSave(createdKey);
    } else {
      message = result.error;
      messageType = 'error';
    }
  }
</script>

<div class="key-form card">
  <h3 class="card-title mb-2">{isEdit ? 'Edit Scoped Key' : 'Create Scoped Key'}</h3>

  <form onsubmit={handleSubmit}>
    <div class="form-group">
      <label class="form-label" for="sk-name">Name *</label>
      <input class="form-input" type="text" id="sk-name" bind:value={name} placeholder="e.g., Discord Bot, Player App" required />
    </div>

    <div class="form-row">
      <div class="form-group">
        <label class="form-label" for="sk-limit">Daily Limit</label>
        <input class="form-input" type="number" id="sk-limit" bind:value={dailyLimit} placeholder="Unlimited" min="0" />
      </div>
      <div class="form-group">
        <label class="form-label" for="sk-expires">Expires At</label>
        <input class="form-input" type="datetime-local" id="sk-expires" bind:value={expiresAt} />
      </div>
    </div>

    <hr class="divider" />

    <button type="button" class="collapsible-trigger" class:open={scopesOpen} onclick={() => scopesOpen = !scopesOpen}>
      <span class="arrow">&#9654;</span>
      Scopes *
    </button>

    {#if scopesOpen}
      <div class="mt-1 mb-2">
        <div class="scope-header">
          <p class="text-muted" style="font-size: 0.8rem; margin: 0;">
            Choose which API scopes this key can use. Dangerous scopes are highlighted.
          </p>
          <div class="scope-controls">
            <button type="button" class="btn btn-sm btn-ghost" onclick={selectAll}>Select All</button>
            <button type="button" class="btn btn-sm btn-ghost" onclick={resetDefaults}>Reset to Defaults</button>
            <button type="button" class="btn btn-sm btn-ghost" onclick={clearAll}>Clear All</button>
          </div>
        </div>
        {#each visibleScopeGroups as [group, scopes]}
          <div class="scope-group">
            <span class="scope-group-label">{group}</span>
            <div class="scope-items">
              {#each scopes as scope}
                <label class="scope-checkbox" class:dangerous={scope === 'execute-js' || scope === 'macro:execute' || scope === 'macro:write'}>
                  <input type="checkbox" checked={selectedScopes.has(scope)} onchange={() => toggleScope(scope)} />
                  <span>{scope}</span>
                  {#if scope === 'execute-js' || scope === 'macro:execute' || scope === 'macro:write'}
                    <span class="badge badge-expired" style="margin-left: 0.25rem; font-size: 0.65rem;">dangerous</span>
                  {/if}
                </label>
              {/each}
            </div>
          </div>
        {/each}

        {#if knownClients.length > 0}
          <div class="mt-2">
            <span class="scope-group-label">Allowed Clients</span>
            <p class="text-muted mb-1" style="font-size: 0.8rem;">
              Restrict this key to specific Foundry worlds. Leave all unchecked for unrestricted access.
            </p>
            {#each knownClients as client (client.id)}
              <label class="scope-checkbox" style="display: flex; margin-bottom: 0.25rem;">
                <input type="checkbox" checked={selectedClientIds.has(client.clientId)} onchange={() => toggleKnownClient(client.clientId)} />
                <span>
                  {client.customName || client.worldTitle || client.clientId.substring(0, 16)}
                  {#if client.isOnline}
                    <span class="badge badge-active" style="margin-left: 0.25rem;">online</span>
                  {/if}
                </span>
              </label>
            {/each}
          </div>

          {#if selectedClientIds.size > 0}
            <div class="mt-2">
              <span class="scope-group-label">Scoped User per World</span>
              <p class="text-muted mb-1" style="font-size: 0.8rem;">
                Optionally restrict requests to a specific player in each world.
              </p>
              {#each Array.from(selectedClientIds) as clientId (clientId)}
                {@const client = knownClients.find(c => c.clientId === clientId)}
                {@const users = knownUsersByClientId[clientId] ?? []}
                <div class="client-user-row">
                  <span class="client-user-label">
                    {client?.customName || client?.worldTitle || clientId.substring(0, 16)}
                  </span>
                  <select class="form-select" style="font-size: 0.8rem;"
                    bind:value={selectedUserIds[clientId]}>
                    <option value="">Unrestricted (any user)</option>
                    {#each users as u (u.userId)}
                      <option value={u.userId}>{u.name} ({roleNames[u.role] ?? `role ${u.role}`})</option>
                    {/each}
                    {#if users.length === 0}
                      <option disabled value="">No stored users — connect world first</option>
                    {/if}
                  </select>
                </div>
              {/each}
            </div>
          {/if}
        {/if}
      </div>
    {/if}

    {#if isEdit}
      <hr class="divider" />

      <button type="button" class="collapsible-trigger" class:open={notificationsOpen} onclick={() => notificationsOpen = !notificationsOpen}>
        <span class="arrow">&#9654;</span>
        Notification Settings (per-key)
      </button>

      {#if notificationsOpen}
        <div class="mt-1">
          <p class="text-muted mb-1" style="font-size: 0.8rem;">
            Configure notifications that fire only for events on this specific key. These run independently of your account-wide notifications.
          </p>

          <div class="form-group">
            <label class="form-label" for="sk-ns-discord">Discord Webhook URL</label>
            <input class="form-input" type="url" id="sk-ns-discord" bind:value={nsDiscordWebhookUrl} placeholder="https://discord.com/api/webhooks/..." />
          </div>

          {#if nsSmtpAvailable}
            <div class="form-group">
              <label class="form-label" for="sk-ns-email">Notification Email</label>
              <input class="form-input" type="email" id="sk-ns-email" bind:value={nsNotifyEmail} placeholder="you@example.com" />
            </div>
          {/if}

          <div class="form-group">
            <span class="scope-group-label">Notify me when&hellip;</span>
            <div class="ns-checkbox-list">
              <label class="scope-checkbox">
                <input type="checkbox" bind:checked={nsNotifyOnExecuteJs} />
                <span>An execute-js call is made</span>
              </label>
              <label class="scope-checkbox">
                <input type="checkbox" bind:checked={nsNotifyOnMacroExecute} />
                <span>A macro-execute call is made</span>
              </label>
              <label class="scope-checkbox">
                <input type="checkbox" bind:checked={nsNotifyOnRateLimit} />
                <span>A rate limit is hit</span>
              </label>
              <label class="scope-checkbox">
                <input type="checkbox" bind:checked={nsNotifyOnError} />
                <span>An error occurs</span>
              </label>
            </div>
          </div>

          <div class="flex gap-1 mt-1">
            <button class="btn btn-secondary" type="button" onclick={handleTestNotification} disabled={nsTesting}>
              {nsTesting ? 'Sending...' : 'Test Notification'}
            </button>
          </div>

          {#if nsMessage}
            <div class="alert mt-1" class:alert-success={nsMessageType === 'success'} class:alert-error={nsMessageType === 'error'}>
              {nsMessage}
            </div>
          {/if}
        </div>
      {/if}
    {/if}

    <div class="flex gap-1 mt-2">
      <button class="btn btn-primary" type="submit" disabled={saving}>
        {saving ? 'Saving...' : (isEdit ? 'Update Key' : 'Create Key')}
      </button>
      <button class="btn btn-secondary" type="button" onclick={onCancel}>Cancel</button>
    </div>

    {#if message}
      <div class="alert mt-1" class:alert-error={messageType === 'error'}>
        {message}
      </div>
    {/if}
  </form>
</div>

<style>
  .scope-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 0.75rem;
  }

  .scope-controls {
    display: flex;
    gap: 0.375rem;
    flex-shrink: 0;
  }

  .scope-group {
    margin-bottom: 0.75rem;
  }

  .scope-group-label {
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
    display: block;
    margin-bottom: 0.25rem;
  }

  .scope-items {
    display: flex;
    flex-wrap: wrap;
    gap: 0.375rem 1rem;
  }

  .scope-checkbox {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.85rem;
    cursor: pointer;
  }

  .scope-checkbox input[type="checkbox"] {
    width: 0.9rem;
    height: 0.9rem;
    accent-color: var(--color-primary);
  }

  .scope-checkbox.dangerous {
    color: var(--color-error);
  }

  .client-user-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.375rem;
  }

  .client-user-label {
    font-size: 0.8rem;
    color: var(--color-text-muted);
    min-width: 8rem;
    flex-shrink: 0;
  }

  .ns-checkbox-list {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
    margin-top: 0.25rem;
  }
</style>
