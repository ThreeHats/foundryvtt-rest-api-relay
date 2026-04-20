<script lang="ts">
  import { onMount } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminUserView } from '../../lib/types';

  let statusCounts = $state<Record<string, number>>({});
  let subscribers = $state<AdminUserView[]>([]);
  let error = $state('');

  async function load() {
    try {
      const r = await adminApi.getSubscriptions();
      statusCounts = r.statusCounts;
      subscribers = r.subscribers;
    } catch (e: any) { error = e?.message ?? 'Failed'; }
  }
  onMount(load);
</script>

<div class="admin-page">
  <h1>Subscriptions</h1>
  {#if error}<p class="error">{error}</p>{/if}

  <div class="cards">
    {#each Object.entries(statusCounts) as [status, count]}
      <div class="card"><h3>{status}</h3><p class="big">{count}</p></div>
    {/each}
  </div>

  <h2>Active Subscribers</h2>
  <table>
    <thead><tr><th>ID</th><th>Email</th><th>Status</th><th>Created</th></tr></thead>
    <tbody>
      {#each subscribers as s}
        <tr><td>{s.id}</td><td>{s.email}</td><td>{s.subscriptionStatus}</td><td>{s.createdAt}</td></tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .admin-page { padding: 1.5rem; }
  .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
  .card { padding: 1rem; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  .card h3 { margin: 0; font-size: 0.875rem; text-transform: uppercase; color: var(--color-text-muted); }
  .big { font-size: 2rem; margin: 0.25rem 0 0; font-weight: 700; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.5rem; text-align: left; border-bottom: 1px solid var(--color-border-light); font-size: 0.875rem; }
  .error { color: var(--color-error, #e53935); }
</style>
