<script lang="ts">
  import { onMount } from 'svelte';
  import { user, updateUser } from '../../lib/auth';
  import { fetchUserData, fetchSubscriptionStatus } from '../../lib/api';
  import ApiKeyDisplay from './ApiKeyDisplay.svelte';
  import UsageStats from './UsageStats.svelte';
  import SubscriptionStatus from './SubscriptionStatus.svelte';
  import ChangePasswordForm from './ChangePasswordForm.svelte';
  import DataPrivacyActions from './DataPrivacyActions.svelte';

  interface Props {
    onAccountDeleted: () => void;
  }

  let { onAccountDeleted }: Props = $props();

  onMount(async () => {
    // Fetch fresh data in background
    const result = await fetchUserData();
    if (result.ok) {
      updateUser(result.data);

      // Also fetch subscription status
      const subResult = await fetchSubscriptionStatus();
      if (subResult.ok) {
        updateUser({
          subscriptionStatus: subResult.data.subscriptionStatus as 'free' | 'active' | 'past_due',
          subscriptionEndsAt: subResult.data.subscriptionEndsAt,
        });
      }
    }
  });
</script>

<h2 class="page-title">Dashboard</h2>

<div class="card">
  <div class="info-row mb-2">
    <span class="form-label">Email</span>
    <span>{$user?.email}</span>
  </div>

  <ApiKeyDisplay />

  <hr class="divider" />

  <UsageStats />

  <hr class="divider" />

  <SubscriptionStatus />
</div>

<div class="card">
  <ChangePasswordForm />
</div>

<div class="card">
  <DataPrivacyActions {onAccountDeleted} />
</div>

<div class="card">
  <h3 class="card-title mb-1">Using Your API Key</h3>
  <pre><code>curl -X GET https://foundryvtt-rest-api-relay.fly.dev/clients \
  -H "x-api-key: YOUR_API_KEY"</code></pre>
</div>

<style>
  .info-row {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }
</style>
