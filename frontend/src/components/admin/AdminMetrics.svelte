<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminMetricsOverview } from '../../lib/types';

  let overview = $state<AdminMetricsOverview | null>(null);
  let endpoints = $state<Array<{ endpoint: string; count: number }>>([]);
  let topUsers = $state<Array<{ userId: number; count: number }>>([]);
  let intervalId: ReturnType<typeof setInterval> | undefined;

  async function load() {
    try {
      overview = await adminApi.getMetricsOverview();
      const ep = await adminApi.getMetricsByEndpoint();
      endpoints = ep.endpoints.slice(0, 20);
      const top = await adminApi.getTopConsumers(10);
      topUsers = top.users;
    } catch {
      // ignore transient errors
    }
  }

  onMount(() => { load(); intervalId = setInterval(load, 10000); });
  onDestroy(() => { if (intervalId) clearInterval(intervalId); });
</script>

<div class="admin-page">
  <h1>Metrics</h1>

  {#if overview}
    <div class="cards">
      <div class="card"><h3>Per Minute</h3><p class="big">{overview.requestsPerMinute}</p></div>
      <div class="card"><h3>Per Hour</h3><p class="big">{overview.requestsPerHour}</p></div>
      <div class="card"><h3>Per Day</h3><p class="big">{overview.requestsPerDay}</p></div>
      <div class="card"><h3>Errors</h3><p class="big">{overview.errorsTotal}</p></div>
    </div>
  {/if}

  <h2>Top Endpoints</h2>
  <div class="bars">
    {#each endpoints as e}
      <div class="bar-row">
        <span class="label">{e.path}</span>
        <div class="bar" style="width: {Math.min(100, (e.count / (endpoints[0]?.count || 1)) * 100)}%"></div>
        <span class="count">{e.count}</span>
      </div>
    {/each}
  </div>

  <h2>Top Consumers (User IDs)</h2>
  <table>
    <thead><tr><th>User ID</th><th>Requests</th></tr></thead>
    <tbody>
      {#each topUsers as u}
        <tr><td>{u.userId}</td><td>{u.count}</td></tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .admin-page { padding: 1.5rem; }
  h1 { margin-top: 0; }
  .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
  .card { padding: 1rem; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  .card h3 { margin: 0; font-size: 0.875rem; color: var(--color-text-muted); }
  .big { font-size: 2rem; margin: 0.25rem 0 0; font-weight: 700; }
  .bars { display: flex; flex-direction: column; gap: 0.25rem; margin-bottom: 2rem; }
  .bar-row { display: grid; grid-template-columns: 250px 1fr 60px; align-items: center; gap: 0.5rem; font-size: 0.75rem; }
  .label { font-family: monospace; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .bar { background: var(--color-primary); height: 18px; border-radius: 4px; min-width: 2px; }
  .count { text-align: right; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: 0.5rem; text-align: left; border-bottom: 1px solid var(--color-border-light); }
</style>
