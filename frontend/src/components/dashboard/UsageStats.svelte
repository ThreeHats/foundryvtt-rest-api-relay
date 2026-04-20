<script lang="ts">
  import { user } from '../../lib/auth';

  let status = $derived($user?.subscriptionStatus ?? 'free');
  let requestsThisMonth = $derived($user?.requestsThisMonth ?? 0);
  let limits = $derived($user?.limits);
  let monthlyLimit = $derived(limits?.monthlyLimit ?? 0);
  let unlimitedMonthly = $derived(limits?.unlimitedMonthly ?? false);
</script>

<div class="usage-row">
  <div class="usage-item">
    <span class="usage-label">This month</span>
    <span class="usage-value">
      {requestsThisMonth.toLocaleString()}{#if monthlyLimit > 0 && !unlimitedMonthly && status !== 'active'}<span class="usage-of"> / {monthlyLimit.toLocaleString()}</span>{/if}
    </span>
  </div>
</div>

<style>
  .usage-row {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .usage-item {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
  }

  .usage-label {
    font-size: 0.7rem;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
  }

  .usage-value {
    font-size: 1.1rem;
    font-weight: 600;
  }

  .usage-of {
    font-size: 0.85rem;
    font-weight: 400;
    color: var(--color-text-muted);
  }
</style>
