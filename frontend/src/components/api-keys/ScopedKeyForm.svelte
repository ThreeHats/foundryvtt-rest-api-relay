<script lang="ts">
  import { fetchClients, fetchPlayers, createScopedKey, updateScopedKey } from '../../lib/api';
  import type { ScopedKey, ConnectedClient, Player } from '../../lib/types';

  interface Props {
    editKey?: ScopedKey | null;
    onSave: (createdKeyValue?: string) => void;
    onCancel: () => void;
  }

  let { editKey = null, onSave, onCancel }: Props = $props();

  let name = $state(editKey?.name ?? '');
  let scopedClientId = $state(editKey?.scopedClientId ?? '');
  let scopedUserId = $state(editKey?.scopedUserId ?? '');
  let dailyLimit = $state(editKey?.dailyLimit?.toString() ?? '');
  let expiresAt = $state(editKey?.expiresAt ? new Date(editKey.expiresAt).toISOString().slice(0, 16) : '');
  let foundryUrl = $state(editKey?.foundryUrl ?? '');
  let foundryUsername = $state(editKey?.foundryUsername ?? '');
  let foundryPassword = $state('');
  let credentialsOpen = $state(false);

  let clients = $state<ConnectedClient[]>([]);
  let players = $state<Player[]>([]);
  let clientsLoading = $state(false);
  let playersLoading = $state(false);

  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');
  let saving = $state(false);

  let isEdit = $derived(!!editKey);

  // Load clients on mount
  $effect(() => {
    loadClients();
  });

  async function loadClients() {
    clientsLoading = true;
    const result = await fetchClients();
    clientsLoading = false;
    if (result.ok) {
      const raw = result.data.clients || (result.data as any) || [];
      clients = Array.isArray(raw) ? raw : [];
    }
  }

  async function loadPlayers() {
    if (!scopedClientId) {
      players = [];
      return;
    }
    playersLoading = true;
    const result = await fetchPlayers(scopedClientId);
    playersLoading = false;
    if (result.ok) {
      const raw = (result.data as any).users || (result.data as any).data || (result.data as any).players || [];
      players = Array.isArray(raw) ? raw : [];
    }
  }

  function handleClientChange() {
    loadPlayers();
  }

  const roleNames: Record<number, string> = { 0: 'None', 1: 'Player', 2: 'Trusted', 3: 'Asst GM', 4: 'GM' };

  async function handleSubmit(e: Event) {
    e.preventDefault();

    if (!name.trim()) {
      message = 'Name is required.';
      messageType = 'error';
      return;
    }

    saving = true;
    message = '';

    const body: Record<string, unknown> = {
      name: name.trim(),
      scopedClientId: scopedClientId || null,
      scopedUserId: scopedUserId || null,
      dailyLimit: dailyLimit ? parseInt(dailyLimit) : null,
      expiresAt: expiresAt || null,
      foundryUrl: foundryUrl || null,
      foundryUsername: foundryUsername || null,
    };
    if (foundryPassword) body.foundryPassword = foundryPassword;

    const result = isEdit
      ? await updateScopedKey(editKey!.id, body)
      : await createScopedKey(body);

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
        <label class="form-label" for="sk-client">
          Scoped Client ID
          <button type="button" class="btn-refresh" onclick={loadClients} disabled={clientsLoading} title="Refresh clients">
            &#8635;
          </button>
        </label>
        <select class="form-select" id="sk-client" bind:value={scopedClientId} onchange={handleClientChange}>
          <option value="">Unrestricted (any client)</option>
          {#each clients as c}
            <option value={c.id || c.clientId || c}>{c.customName ? `${c.customName} (${c.id || c.clientId})` : (c.id || c.clientId)}</option>
          {/each}
        </select>
      </div>
      <div class="form-group">
        <label class="form-label" for="sk-user">
          Scoped User ID
          <button type="button" class="btn-refresh" onclick={loadPlayers} disabled={playersLoading} title="Refresh players">
            &#8635;
          </button>
        </label>
        <select class="form-select" id="sk-user" bind:value={scopedUserId}>
          <option value="">Unrestricted (any user)</option>
          {#each players as p}
            <option value={p.name || p.id}>
              {p.name || p.id} ({roleNames[p.role] || `role ${p.role}`}{p.active ? ' - online' : ''})
            </option>
          {/each}
        </select>
        <p class="form-hint">Select a client first, then refresh to load players</p>
        {#if !scopedUserId}
          <div class="alert alert-warning mt-1" style="font-size: 0.8rem;">
            Without a Scoped User ID, all requests through this key run with full GM permissions.
          </div>
        {/if}
      </div>
    </div>

    <div class="form-row">
      <div class="form-group">
        <label class="form-label" for="sk-limit">Daily Limit</label>
        <input class="form-input" type="number" id="sk-limit" bind:value={dailyLimit} placeholder="Unlimited" min="1" />
      </div>
      <div class="form-group">
        <label class="form-label" for="sk-expires">Expires At</label>
        <input class="form-input" type="datetime-local" id="sk-expires" bind:value={expiresAt} />
      </div>
    </div>

    <hr class="divider" />

    <button type="button" class="collapsible-trigger" class:open={credentialsOpen} onclick={() => credentialsOpen = !credentialsOpen}>
      <span class="arrow">&#9654;</span>
      Foundry Credentials (optional)
    </button>

    {#if credentialsOpen}
      <div class="mt-1">
        <p class="text-muted mb-1" style="font-size: 0.8rem;">
          Store Foundry credentials to enable headless sessions without additional parameters.
        </p>
        <div class="alert alert-warning mb-2" style="font-size: 0.8rem;">
          The username and password must be for a GM user in Foundry.
        </div>
        <div class="form-group">
          <label class="form-label" for="sk-furl">Foundry URL</label>
          <input class="form-input" type="url" id="sk-furl" bind:value={foundryUrl} placeholder="https://your-foundry-instance.com" />
        </div>
        <div class="form-group">
          <label class="form-label" for="sk-fuser">Foundry Username</label>
          <input class="form-input" type="text" id="sk-fuser" bind:value={foundryUsername} placeholder="GM username" />
        </div>
        <div class="form-group">
          <label class="form-label" for="sk-fpass">Foundry Password</label>
          <input class="form-input" type="password" id="sk-fpass" bind:value={foundryPassword} placeholder={isEdit ? 'Leave blank to keep existing' : 'Encrypted at rest'} />
        </div>
      </div>
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
  .btn-refresh {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    color: var(--color-primary);
    padding: 0 0.25rem;
    vertical-align: middle;
    transition: transform 0.3s;
  }

  .btn-refresh:hover {
    color: var(--color-primary-hover);
  }

  .btn-refresh:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
