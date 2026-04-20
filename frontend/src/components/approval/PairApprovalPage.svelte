<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchPairRequest, approvePairRequest, denyPairRequest } from '../../lib/api';
  import type { PairRequestDetails, KnownClient } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';
  import { scopeDescriptions, dangerDetails } from '../../lib/scopeDescriptions';

  interface Props {
    code: string;
  }

  let { code }: Props = $props();

  let request = $state<PairRequestDetails | null>(null);
  let loading = $state(true);
  let processing = $state(false);
  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');
  let done = $state(false);
  let upgraded = $state(false);

  // Cross-world settings — pre-populated from requestedRemoteScopes / requestedTargetClients
  let selectedScopes = $state<string[]>([]);
  let selectedTargets = $state<string[]>([]);
  let rateLimit = $state(0);
  let showCrossWorld = $state(false);


  onMount(async () => {
    const result = await fetchPairRequest(code);
    loading = false;
    if (result.ok) {
      request = result.data;
      selectedScopes = [...(request.requestedRemoteScopes || [])];
      selectedTargets = [...(request.requestedTargetClients || [])];
      // Pre-expand cross-world section if the world requested any permissions.
      showCrossWorld = selectedScopes.length > 0 || selectedTargets.length > 0;
    } else {
      message = result.error;
      messageType = 'error';
    }
  });

  function toggleScope(scope: string) {
    if (selectedScopes.includes(scope)) {
      selectedScopes = selectedScopes.filter(s => s !== scope);
    } else {
      selectedScopes = [...selectedScopes, scope];
    }
  }

  function toggleTarget(clientId: string) {
    if (selectedTargets.includes(clientId)) {
      selectedTargets = selectedTargets.filter(t => t !== clientId);
    } else {
      selectedTargets = [...selectedTargets, clientId];
    }
  }

  function clientLabel(client: KnownClient): string {
    return client.customName || client.worldTitle || client.clientId;
  }

  async function handleApprove() {
    processing = true;
    message = '';
    const result = await approvePairRequest(code, {
      remoteScopes: selectedScopes,
      allowedTargetClients: selectedTargets,
      remoteRequestsPerHour: rateLimit,
    });
    processing = false;
    if (result.ok) {
      done = true;
      upgraded = !!result.data.upgraded;
      message = upgraded
        ? 'Cross-world permissions updated! Your Foundry module has been upgraded. This tab will close shortly.'
        : 'Pairing approved! Your Foundry module will connect shortly. This tab will close shortly.';
      messageType = 'success';
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
    const result = await denyPairRequest(code);
    processing = false;
    if (result.ok) {
      done = true;
      message = 'Pairing request denied.';
      messageType = 'success';
    } else {
      message = result.error;
      messageType = 'error';
    }
  }
</script>

<ConfirmModal
  open={denyModalOpen}
  title="Deny pairing request"
  message="Deny this pairing request? The Foundry module will not receive a connection token."
  confirmLabel="Deny"
  dangerous={true}
  onConfirm={handleDeny}
  onCancel={() => denyModalOpen = false}
/>

<h2 class="page-title">
  {#if request?.upgradeOnly}
    Upgrade World Permissions
  {:else}
    Pair Foundry World
  {/if}
</h2>

{#if loading}
  <div class="card">
    <p class="text-muted">Loading request details...</p>
  </div>
{:else if !request}
  <div class="card">
    <div class="alert alert-error">
      {message || 'Pair request not found or has expired.'}
    </div>
  </div>
{:else if request.status !== 'pending'}
  <div class="card">
    <div class="alert alert-warning">
      This pair request has already been {request.status}.
    </div>
  </div>
{:else}
  <!-- World Info -->
  <div class="card">
    <h3 class="card-title mb-1">World Details</h3>
    {#if request.upgradeOnly}
      <div class="alert alert-info mb-2" style="font-size: 0.85rem;">
        This world is already paired. Approving will only update its cross-world permissions — no new connection token will be created.
      </div>
    {/if}
    <div class="info-grid">
      <div class="info-row">
        <span class="form-label">World</span>
        <span>{request.worldTitle || request.worldId}</span>
      </div>
      {#if request.systemTitle}
        <div class="info-row">
          <span class="form-label">System</span>
          <span>{request.systemTitle}{request.systemVersion ? ` v${request.systemVersion}` : ''}</span>
        </div>
      {/if}
      {#if request.foundryVersion}
        <div class="info-row">
          <span class="form-label">Foundry</span>
          <span>v{request.foundryVersion}</span>
        </div>
      {/if}
      <div class="info-row">
        <span class="form-label">Expires</span>
        <span>{new Date(request.expiresAt).toLocaleString()}</span>
      </div>
    </div>
  </div>

  <!-- Cross-World Settings -->
  <div class="card">
    <div class="cross-world-header" role="button" tabindex="0"
      onclick={() => showCrossWorld = !showCrossWorld}
      onkeydown={(e) => e.key === 'Enter' && (showCrossWorld = !showCrossWorld)}>
      <h3 class="card-title" style="margin: 0;">
        Cross-World Permissions
        {#if selectedScopes.length > 0 || selectedTargets.length > 0}
          <span class="badge badge-active" style="margin-left: 0.5rem; font-size: 0.7rem; vertical-align: middle;">
            {selectedScopes.length} scope{selectedScopes.length !== 1 ? 's' : ''},
            {selectedTargets.length} target{selectedTargets.length !== 1 ? 's' : ''}
          </span>
        {:else}
          <span class="badge" style="margin-left: 0.5rem; font-size: 0.7rem; vertical-align: middle; background: var(--color-border); color: var(--color-text-muted);">
            none
          </span>
        {/if}
      </h3>
      <span class="chevron" class:open={showCrossWorld}>▾</span>
    </div>

    {#if showCrossWorld}
      <div class="cross-world-body">
        <p class="text-muted" style="font-size: 0.8rem; margin-bottom: 1rem;">
          Allow this world to make relay requests <em>through</em> the relay to other worlds you own.
          These settings are optional — leave both sections empty to skip cross-world access.
        </p>

        <!-- Target worlds -->
        <div class="section-label">Allowed target worlds</div>
        <p class="text-muted" style="font-size: 0.8rem; margin-bottom: 0.5rem;">
          Which of your other worlds may this world send remote requests to?
        </p>
        {#if (request.requestedTargetClients || []).length > 0}
          {@const requestedClients = (request.knownClients || []).filter(c => (request.requestedTargetClients || []).includes(c.clientId))}
          {#if requestedClients.length === 0}
            <p class="text-muted" style="font-size: 0.875rem;">No matching paired worlds found.</p>
          {:else}
            <div class="checkbox-list">
              {#each requestedClients as client (client.id)}
                <label class="checkbox-label">
                  <input type="checkbox"
                    checked={selectedTargets.includes(client.clientId)}
                    onchange={() => toggleTarget(client.clientId)} />
                  <span>
                    {clientLabel(client)}
                    {#if client.isOnline}
                      <span class="badge badge-active" style="margin-left: 0.25rem; font-size: 0.7rem;">online</span>
                    {/if}
                  </span>
                </label>
              {/each}
            </div>
          {/if}
        {:else if (request.knownClients || []).filter(c => c.clientId !== request?.worldId).length === 0}
          <p class="text-muted" style="font-size: 0.875rem;">No other paired worlds found.</p>
        {:else}
          <div class="checkbox-list">
            {#each (request.knownClients || []).filter(c => c.clientId !== request?.worldId) as client (client.id)}
              <label class="checkbox-label">
                <input type="checkbox"
                  checked={selectedTargets.includes(client.clientId)}
                  onchange={() => toggleTarget(client.clientId)} />
                <span>
                  {clientLabel(client)}
                  {#if client.isOnline}
                    <span class="badge badge-active" style="margin-left: 0.25rem; font-size: 0.7rem;">online</span>
                  {/if}
                </span>
              </label>
            {/each}
          </div>
        {/if}

        <!-- Allowed scopes -->
        <div class="section-label" style="margin-top: 1rem;">Allowed scopes</div>
        <p class="text-muted" style="font-size: 0.8rem; margin-bottom: 0.5rem;">
          Operations this world may invoke on its allowed target worlds.
        </p>
        {#if (request.requestedRemoteScopes || []).length === 0}
          <p class="text-muted" style="font-size: 0.875rem;">No specific scopes were requested.</p>
        {:else}
          <div class="scope-list">
            {#each (request.requestedRemoteScopes || []) as scope}
              {@const isDangerous = scope in dangerDetails}
              <label class="scope-label" class:dangerous={isDangerous}>
                <input type="checkbox"
                  checked={selectedScopes.includes(scope)}
                  onchange={() => toggleScope(scope)} />
                <span>
                  {scope}
                  {#if isDangerous}
                    <span class="badge badge-expired" style="margin-left: 0.25rem; font-size: 0.65rem;">dangerous</span>
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
        {/if}

        <!-- Rate limit -->
        <div class="form-group" style="margin-top: 1rem; max-width: 260px;">
          <label class="form-label" for="pr-rate-limit">
            Rate limit (remote requests / hour)
          </label>
          <input class="form-input" type="number" id="pr-rate-limit"
            bind:value={rateLimit} min="0" placeholder="0 = unlimited" />
          <p class="text-muted" style="font-size: 0.75rem; margin-top: 0.25rem;">0 = unlimited</p>
        </div>
      </div>
    {/if}
  </div>

  <!-- Actions -->
  <div class="card">
    {#if done}
      <div class="alert alert-success">
        {message}
      </div>
    {:else}
      {#if !request.upgradeOnly}
        <p class="text-muted" style="font-size: 0.875rem; margin-bottom: 0.75rem;">
          Approving will link <strong>{request.worldTitle || request.worldId}</strong> to your relay account.
          The Foundry module will be able to receive REST API requests forwarded by the relay.
          You can manage or revoke access at any time from the <strong>Connections</strong> page.
        </p>
      {/if}
      <div class="flex gap-1">
        <button class="btn btn-primary" onclick={handleApprove} disabled={processing}>
          {#if processing}
            Processing...
          {:else if request.upgradeOnly}
            Upgrade Permissions
          {:else}
            Approve Pairing
          {/if}
        </button>
        <button class="btn btn-danger" onclick={handleDenyClick} disabled={processing}>
          Deny
        </button>
      </div>
    {/if}

    {#if message && !done}
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

  .mb-1 { margin-bottom: 0.5rem; }
  .mb-2 { margin-bottom: 0.75rem; }

  .cross-world-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
    user-select: none;
  }

  .chevron {
    font-size: 1.1rem;
    color: var(--color-text-muted);
    transition: transform 0.15s ease;
  }

  .chevron.open {
    transform: rotate(180deg);
  }

  .cross-world-body {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
  }

  .section-label {
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
    margin-bottom: 0.25rem;
  }

  .checkbox-list {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
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

  .scope-list {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .scope-label {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.8rem;
    cursor: pointer;
  }

  .scope-label input[type="checkbox"] {
    width: 0.9rem;
    height: 0.9rem;
    accent-color: var(--color-primary);
  }

  .scope-label.dangerous {
    color: var(--color-error);
  }

  .scope-desc {
    font-size: 0.75rem;
    color: var(--color-text-muted);
    margin: 0 0 0.2rem 1.4rem;
  }

  .danger-box {
    font-size: 0.75rem;
    color: var(--color-error);
    background: color-mix(in srgb, var(--color-error) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-error) 30%, transparent);
    border-radius: 4px;
    padding: 0.35rem 0.5rem;
    margin: 0 0 0.4rem 1.4rem;
    display: flex;
    gap: 0.4rem;
    align-items: flex-start;
  }

  .alert-info {
    background: rgba(59, 130, 246, 0.08);
    border: 1px solid rgba(59, 130, 246, 0.3);
    border-radius: 4px;
    padding: 0.5rem 0.75rem;
    color: var(--color-text);
  }
</style>
