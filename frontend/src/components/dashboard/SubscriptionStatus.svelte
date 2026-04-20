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

<div class="sub-inline">
  {#if status === 'active'}
    <span class="badge badge-active">Active</span>
    <button class="btn btn-ghost btn-sm" onclick={handleManage} disabled={loading}>Manage</button>
  {:else if status === 'past_due'}
    <span class="badge badge-disabled">Past Due</span>
    <button class="btn btn-ghost btn-sm" onclick={handleManage} disabled={loading}>Manage</button>
  {:else}
    <span class="badge badge-expired">Free plan</span>
    <button class="btn btn-primary btn-sm" onclick={handleSubscribe} disabled={loading}>
      {loading ? 'Loading…' : 'Upgrade · $5/month'}
    </button>
  {/if}
</div>

<style>
  .sub-inline {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }
</style>
