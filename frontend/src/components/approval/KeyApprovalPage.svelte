<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchKeyRequest, approveKeyRequest, denyKeyRequest, fetchKnownClients } from '../../lib/api';
  import type { KeyRequestDetails, KnownClient } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';
  import { headlessEnabled } from '../../lib/config';
  import { scopeDescriptions, dangerDetails } from '../../lib/scopeDescriptions';

  interface Props {
    code: string;
  }

  let { code }: Props = $props();

  let request = $state<KeyRequestDetails | null>(null);
  let knownClients = $state<KnownClient[]>([]);
  let loading = $state(true);
  let processing = $state(false);
  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');
  let approvedKeyValue = $state<string | null>(null);

  // Form state
  let selectedScopes = $state<Set<string>>(new Set());
  let selectedClientIds = $state<Set<string>>(new Set());
  let dailyLimit = $state('');
  let expiresAt = $state('');

  onMount(async () => {
    const [reqResult, clientsResult] = await Promise.all([
      fetchKeyRequest(code),
      fetchKnownClients(),
    ]);

    loading = false;

    if (reqResult.ok) {
      request = reqResult.data;
      selectedScopes = new Set(request.requestedScopes || []);
      selectedClientIds = new Set(request.requestedClientIds || []);
      dailyLimit = request.suggestedDailyLimit ? request.suggestedDailyLimit.toString() : '';
      expiresAt = request.suggestedExpiry ? new Date(request.suggestedExpiry).toISOString().slice(0, 16) : '';
    } else {
      message = reqResult.error;
      messageType = 'error';
    }

    if (clientsResult.ok) {
      knownClients = clientsResult.data.clients || [];
    }
  });

  function toggleScope(scope: string) {
    // Users can only uncheck, not add new scopes beyond what was requested
    if (selectedScopes.has(scope)) {
      selectedScopes.delete(scope);
      selectedScopes = new Set(selectedScopes);
    } else if (request?.requestedScopes.includes(scope)) {
      selectedScopes.add(scope);
      selectedScopes = new Set(selectedScopes);
    }
  }

  function toggleClient(clientId: string) {
    if (selectedClientIds.has(clientId)) {
      selectedClientIds.delete(clientId);
    } else if (!request?.requestedClientIds?.length || request.requestedClientIds.includes(clientId)) {
      selectedClientIds.add(clientId);
    }
    selectedClientIds = new Set(selectedClientIds);
  }

  async function handleApprove() {
    processing = true;
    message = '';

    const result = await approveKeyRequest(code, {
      scopes: Array.from(selectedScopes),
      clientIds: Array.from(selectedClientIds),
      dailyLimit: dailyLimit ? parseInt(dailyLimit) : null,
      expiresAt: expiresAt || null,
    });

    processing = false;

    if (result.ok) {
      message = 'Key request approved! This tab will close shortly.';
      messageType = 'success';
      if (result.data.key) {
        approvedKeyValue = result.data.key;
      }
      setTimeout(() => window.close(), 3000);
    } else {
      message = result.error;
      messageType = 'error';
    }
  }

  let denyModalOpen = $state(false);

  function handleDenyClick() { denyModalOpen = true; }

  async function handleDeny() {
    denyModalOpen = false;
    processing = true;
    message = '';

    const result = await denyKeyRequest(code);

    processing = false;

    if (result.ok) {
      message = 'Key request denied. This tab will close shortly.';
      messageType = 'success';
      setTimeout(() => window.close(), 3000);
    } else {
      message = result.error;
      messageType = 'error';
    }
  }

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

  // Hide Session scope when headless sessions are disabled on this relay.
  let visibleScopeGroups = $derived(
    Object.entries(scopeGroups).filter(([name]) => name !== 'Session' || $headlessEnabled)
  );
</script>

<ConfirmModal
  open={denyModalOpen}
  title="Deny key request"
  message="Deny this key request? The requesting application will not receive a key."
  confirmLabel="Deny"
  dangerous={true}
  onConfirm={handleDeny}
  onCancel={() => denyModalOpen = false}
/>

<h2 class="page-title">Approve API Key Request</h2>

{#if loading}
  <div class="card">
    <p class="text-muted">Loading request details...</p>
  </div>
{:else if !request}
  <div class="card">
    <div class="alert alert-error">
      {message || 'Key request not found or has expired.'}
    </div>
  </div>
{:else if request.status !== 'pending'}
  <div class="card">
    <div class="alert alert-warning">
      This key request has already been {request.status}.
    </div>
  </div>
{:else}
  <!-- App Info -->
  <div class="card">
    <h3 class="card-title mb-1">Application Details</h3>
    <div class="info-grid">
      <div class="info-row">
        <span class="form-label">App Name</span>
        <span>{request.appName}</span>
      </div>
      {#if request.appDescription}
        <div class="info-row">
          <span class="form-label">Description</span>
          <span>{request.appDescription}</span>
        </div>
      {/if}
      {#if request.appUrl}
        <div class="info-row">
          <span class="form-label">URL</span>
          <a href={request.appUrl} target="_blank" rel="noopener">{request.appUrl}</a>
        </div>
      {/if}
      <div class="info-row">
        <span class="form-label">Expires</span>
        <span>{new Date(request.expiresAt).toLocaleString()}</span>
      </div>
    </div>
  </div>

  <!-- Scopes -->
  <div class="card">
    <h3 class="card-title mb-1">Requested Scopes</h3>
    <p class="text-muted mb-2" style="font-size: 0.8rem;">
      Uncheck scopes you do not want to grant. You cannot add scopes beyond what was requested.
    </p>

    {#each visibleScopeGroups as [group, scopes]}
      {@const relevantScopes = scopes.filter(s => request?.requestedScopes.includes(s))}
      {#if relevantScopes.length > 0}
        <div class="scope-group">
          <span class="scope-group-label">{group}</span>
          <div class="scope-items">
            {#each relevantScopes as scope}
              {@const isDangerous = scope in dangerDetails}
              <label class="checkbox-label" class:dangerous={isDangerous}>
                <input type="checkbox" checked={selectedScopes.has(scope)} onchange={() => toggleScope(scope)} />
                <span>
                  {scope}
                  {#if isDangerous}
                    <span class="badge badge-expired" style="margin-left: 0.25rem;">dangerous</span>
                  {/if}
                </span>
              </label>
              {#if scopeDescriptions[scope]}
                <p class="scope-desc">{scopeDescriptions[scope]}</p>
              {/if}
              {#if isDangerous}
                <div class="danger-box">
                  <i class="fa-solid fa-triangle-exclamation"></i>
                  {dangerDetails[scope]}
                </div>
              {/if}
            {/each}
          </div>
        </div>
      {/if}
    {/each}
  </div>

  <!-- Client Selector -->
  <div class="card">
    <h3 class="card-title mb-1">Allowed Clients</h3>
    <p class="text-muted mb-2" style="font-size: 0.8rem;">
      {#if request?.requestedClientIds?.length}
        Uncheck worlds you do not want to grant access to. You cannot add worlds beyond what was requested.
      {:else}
        Select which Foundry worlds this key can access.
      {/if}
    </p>

    {#if knownClients.length === 0}
      <p class="text-muted" style="font-size: 0.875rem;">No known clients available.</p>
    {:else}
      {@const visibleClients = (request?.requestedClientIds?.length)
        ? knownClients.filter(c => request!.requestedClientIds.includes(c.clientId))
        : knownClients}
      {#if visibleClients.length === 0}
        <p class="text-muted" style="font-size: 0.875rem;">No matching clients found.</p>
      {:else}
        {#each visibleClients as client (client.id)}
          <label class="checkbox-label mb-1">
            <input type="checkbox" checked={selectedClientIds.has(client.clientId)} onchange={() => toggleClient(client.clientId)} />
            <span>
              {client.customName || client.worldTitle || client.clientId}
              {#if client.isOnline}
                <span class="badge badge-active" style="margin-left: 0.25rem;">online</span>
              {/if}
            </span>
          </label>
        {/each}
      {/if}
    {/if}
  </div>

  <!-- Limits -->
  <div class="card">
    <h3 class="card-title mb-1">Limits</h3>
    <div class="form-row">
      <div class="form-group">
        <label class="form-label" for="kr-limit">Daily Request Limit</label>
        <input class="form-input" type="number" id="kr-limit" bind:value={dailyLimit} placeholder="Unlimited" min="1" />
      </div>
      <div class="form-group">
        <label class="form-label" for="kr-expires">Expiry Date</label>
        <input class="form-input" type="datetime-local" id="kr-expires" bind:value={expiresAt} />
      </div>
    </div>
  </div>

  <!-- Actions -->
  <div class="card">
    {#if approvedKeyValue}
      <div class="alert alert-success mb-2" style="border: 2px solid var(--color-success);">
        <strong>Key approved and created!</strong>
        {#if request.callbackUrl}
          <p style="font-size: 0.85rem; margin-top: 0.25rem;">The key has been sent to the application's callback URL.</p>
        {/if}
      </div>
    {:else}
      <div class="flex gap-1">
        <button class="btn btn-primary" onclick={handleApprove} disabled={processing}>
          {processing ? 'Processing...' : 'Approve'}
        </button>
        <button class="btn btn-danger" onclick={handleDenyClick} disabled={processing}>
          Deny
        </button>
      </div>
    {/if}

    {#if message}
      <div class="alert mt-1" class:alert-success={messageType === 'success'} class:alert-error={messageType === 'error'}>
        {message}
      </div>
    {/if}
  </div>
{/if}

<style>
  .info-grid {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .info-row {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .info-row a {
    color: var(--color-primary);
    text-decoration: underline;
    text-underline-offset: 2px;
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
    flex-direction: column;
    gap: 0.25rem;
  }

  .checkbox-label {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    cursor: pointer;
  }

  .checkbox-label input[type="checkbox"] {
    width: 1rem;
    height: 1rem;
    accent-color: var(--color-primary);
  }

  .checkbox-label.dangerous {
    color: var(--color-error);
  }

  .scope-desc {
    font-size: 0.78rem;
    color: var(--color-text-muted);
    margin: 0 0 0.25rem 1.5rem;
  }

  .danger-box {
    font-size: 0.78rem;
    color: var(--color-error);
    background: color-mix(in srgb, var(--color-error) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-error) 30%, transparent);
    border-radius: 4px;
    padding: 0.4rem 0.6rem;
    margin: 0 0 0.5rem 1.5rem;
    display: flex;
    gap: 0.4rem;
    align-items: flex-start;
  }

  .mb-1 {
    margin-bottom: 0.5rem;
  }
</style>
