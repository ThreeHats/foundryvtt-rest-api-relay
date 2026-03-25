<script lang="ts">
  import { user } from '../../lib/auth';

  let status = $derived($user?.subscriptionStatus ?? 'free');
  let requestsToday = $derived($user?.requestsToday ?? 0);
  let requestsThisMonth = $derived($user?.requestsThisMonth ?? 0);
  let limits = $derived($user?.limits);
  let dailyLimit = $derived(limits?.dailyLimit ?? 1000);
  let monthlyLimit = $derived(limits?.monthlyLimit ?? 100);
  let unlimitedMonthly = $derived(limits?.unlimitedMonthly ?? false);
</script>

<div class="stats-grid">
  <div class="stat-item">
    <span class="stat-label">Daily Requests</span>
    <span class="stat-value">{requestsToday.toLocaleString()} <span class="stat-limit">/ {dailyLimit.toLocaleString()}</span></span>
  </div>
  <div class="stat-item">
    <span class="stat-label">Monthly Requests</span>
    <span class="stat-value">
      {requestsThisMonth.toLocaleString()}
      {#if unlimitedMonthly || status === 'active'}
        <span class="stat-limit">(unlimited)</span>
      {:else}
        <span class="stat-limit">/ {monthlyLimit.toLocaleString()}</span>
      {/if}
    </span>
  </div>
</div>

{#if !unlimitedMonthly && status !== 'active'}
  <div class="alert alert-warning mt-1">
    All users are limited to {dailyLimit.toLocaleString()} requests per day.
    Free accounts are limited to {monthlyLimit.toLocaleString()} requests per month. Subscribe for unlimited monthly access.
  </div>
{:else}
  <div class="alert alert-info mt-1">
    All users are limited to {dailyLimit.toLocaleString()} requests per day.
    You have unlimited monthly access with your subscription.
  </div>
{/if}

<style>
  .stats-grid {
    display: flex;
    gap: 1.5rem;
    flex-wrap: wrap;
  }

  .stat-item {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .stat-label {
    font-size: 0.75rem;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
  }

  .stat-value {
    font-size: 1.25rem;
    font-weight: 600;
  }

  .stat-limit {
    font-size: 0.85rem;
    font-weight: 400;
    color: var(--color-text-muted);
  }
</style>
