<script lang="ts">
  import { user } from '../../lib/auth';
  import { createCheckoutSession, createPortalSession } from '../../lib/api';

  let status = $derived($user?.subscriptionStatus ?? 'free');
  let loading = $state(false);

  async function handleSubscribe() {
    loading = true;
    const result = await createCheckoutSession();
    loading = false;
    if (result.ok) {
      window.location.href = result.data.url;
    } else {
      alert(result.error);
    }
  }

  async function handleManage() {
    loading = true;
    const result = await createPortalSession();
    loading = false;
    if (result.ok) {
      window.location.href = result.data.url;
    } else {
      alert(result.error);
    }
  }
</script>

<div class="subscription-row">
  <div class="status-display">
    <span class="form-label">Subscription</span>
    <span class="status-text">
      {#if status === 'active'}
        <span class="badge badge-active">Active</span>
      {:else if status === 'past_due'}
        <span class="badge badge-disabled">Past Due</span>
      {:else}
        <span class="badge badge-expired">Free</span>
      {/if}
    </span>
  </div>
  <div class="subscription-actions">
    {#if status === 'active' || status === 'past_due'}
      <button class="btn btn-secondary btn-sm" onclick={handleManage} disabled={loading}>
        Manage Subscription
      </button>
    {:else}
      <button class="btn btn-primary btn-sm" onclick={handleSubscribe} disabled={loading}>
        {loading ? 'Loading...' : '$5/month - Unlimited Access'}
      </button>
    {/if}
  </div>
</div>

<style>
  .subscription-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .status-display {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .status-text {
    font-size: 0.875rem;
  }
</style>
