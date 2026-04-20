<script lang="ts">
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';
  import { user, updateUser } from '../../lib/auth';
  import { fetchUserData, fetchSubscriptionStatus, fetchKnownClients, fetchScopedKeys } from '../../lib/api';
  import { billingEnabled } from '../../lib/config';
  import type { KnownClient, ScopedKey } from '../../lib/types';
  import ApiKeyDisplay from './ApiKeyDisplay.svelte';
  import UsageStats from './UsageStats.svelte';
  import SubscriptionStatus from './SubscriptionStatus.svelte';
  import ChangePasswordForm from './ChangePasswordForm.svelte';
  import DataPrivacyActions from './DataPrivacyActions.svelte';

  interface Props {
    onAccountDeleted: () => void;
    onNavigate?: (view: string) => void;
  }

  let { onAccountDeleted, onNavigate }: Props = $props();

  let knownClients = $state<KnownClient[]>([]);
  let scopedKeys = $state<ScopedKey[]>([]);
  let clientsLoading = $state(true);
  let keysLoading = $state(true);

  let onlineCount = $derived(knownClients.filter(c => c.isOnline).length);

  onMount(async () => {
    // Fetch fresh data in background
    const result = await fetchUserData();
    if (result.ok) {
      updateUser(result.data);

      // Also fetch subscription status (only needed when billing is enabled)
      if (get(billingEnabled)) {
        const subResult = await fetchSubscriptionStatus();
        if (subResult.ok) {
          updateUser({
            subscriptionStatus: subResult.data.subscriptionStatus as 'free' | 'active' | 'past_due',
            subscriptionEndsAt: subResult.data.subscriptionEndsAt,
          });
        }
      }
    }

    // Fetch known clients
    const clientsResult = await fetchKnownClients();
    clientsLoading = false;
    if (clientsResult.ok) {
      knownClients = clientsResult.data.clients || [];
    }

    // Fetch scoped keys
    const keysResult = await fetchScopedKeys();
    keysLoading = false;
    if (keysResult.ok) {
      scopedKeys = keysResult.data.keys || [];
    }
  });
</script>

<h2 class="page-title">Dashboard</h2>

<div class="card">
  <div class="account-header">
    <span class="account-email">{$user?.email}</span>
    {#if $billingEnabled}
      <SubscriptionStatus />
    {/if}
  </div>
  <div class="mt-2">
    <UsageStats />
  </div>
</div>

<div class="card">
  <ChangePasswordForm />
</div>

<div class="card">
  <ApiKeyDisplay />
</div>

<div class="card">
  <DataPrivacyActions {onAccountDeleted} />
</div>

<!-- Connected Worlds -->
<div class="card">
  <div class="card-header">
    <h3 class="card-title">Connected Worlds</h3>
    {#if onNavigate}
      <button class="btn btn-sm btn-ghost" onclick={() => onNavigate?.('connections')}>View All</button>
    {/if}
  </div>

  {#if clientsLoading}
    <p class="text-muted" style="font-size: 0.85rem;">Loading...</p>
  {:else if knownClients.length === 0}
    <p class="text-muted" style="font-size: 0.85rem;">No connected worlds yet.</p>
  {:else}
    <p class="text-muted mb-1" style="font-size: 0.8rem;">{onlineCount} online, {knownClients.length} total</p>
    <div class="client-list">
      {#each knownClients.slice(0, 5) as client (client.id)}
        <div class="client-item">
          <span>
            {#if client.isOnline}
              <span class="badge badge-active">online</span>
            {:else}
              <span class="badge badge-disabled">offline</span>
            {/if}
          </span>
          <span class="client-name">{client.customName || client.worldTitle || client.clientId.substring(0, 16)}</span>
          <span class="text-muted" style="font-size: 0.8rem;">{client.systemTitle || client.systemId}</span>
        </div>
      {/each}
      {#if knownClients.length > 5}
        <p class="text-muted" style="font-size: 0.8rem;">...and {knownClients.length - 5} more</p>
      {/if}
    </div>
  {/if}
</div>

<!-- API Keys -->
<div class="card">
  <div class="card-header">
    <h3 class="card-title">API Keys</h3>
    {#if onNavigate}
      <button class="btn btn-sm btn-ghost" onclick={() => onNavigate?.('api-keys')}>View All</button>
    {/if}
  </div>

  {#if keysLoading}
    <p class="text-muted" style="font-size: 0.85rem;">Loading...</p>
  {:else}
    <p class="text-muted" style="font-size: 0.85rem;">
      {scopedKeys.length} scoped key{scopedKeys.length !== 1 ? 's' : ''}
    </p>
    {#if scopedKeys.length > 0}
      <div class="integration-list mt-1">
        {#each scopedKeys.slice(0, 5) as key (key.id)}
          <div class="integration-item">
            <span class="integration-name">{key.name}</span>
            {#if key.isExpired}
              <span class="badge badge-expired">expired</span>
            {/if}
            {#if key.dailyLimit}
              <span class="text-muted" style="font-size: 0.8rem;">{key.requestsToday}/{key.dailyLimit} today</span>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

<div class="card">
  <h3 class="card-title mb-1">Using Your API Key</h3>
  <pre><code>curl -X GET https://foundryrestapi.com/clients \
  -H "x-api-key: YOUR_API_KEY"</code></pre>
</div>

<style>
  .account-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .account-email {
    font-size: 0.9rem;
    color: var(--color-text-muted);
  }

  .client-list {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
  }

  .client-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
  }

  .client-name {
    font-weight: 500;
  }

  .integration-list {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
  }

  .integration-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
  }

  .integration-name {
    font-weight: 500;
  }
</style>
