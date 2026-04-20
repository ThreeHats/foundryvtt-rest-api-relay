<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminSystemHealth } from '../../lib/types';
  import { headlessEnabled } from '../../lib/config';

  let health = $state<AdminSystemHealth | null>(null);
  let intervalId: ReturnType<typeof setInterval> | undefined;

  async function load() {
    try { health = await adminApi.getHealth(); } catch { /* ignore */ }
  }
  onMount(() => { load(); intervalId = setInterval(load, 10000); });
  onDestroy(() => { if (intervalId) clearInterval(intervalId); });

  function fmt(bytes: number): string {
    if (bytes >= 1e9) return (bytes / 1e9).toFixed(2) + ' GB';
    if (bytes >= 1e6) return (bytes / 1e6).toFixed(1) + ' MB';
    if (bytes >= 1e3) return (bytes / 1e3).toFixed(1) + ' KB';
    return bytes + ' B';
  }
</script>

<div class="admin-page">
  <h1>System Health</h1>
  {#if health}
    <div class="cards">
      <div class="card"><h3>Status</h3><p class="big" class:degraded={health.status !== 'ok'}>{health.status}</p></div>
      <div class="card"><h3>Version</h3><p class="big">{health.version}</p></div>
      <div class="card"><h3>WS Connections</h3><p class="big">{health.wsConnectionCount}</p></div>
      <div class="card"><h3>Pending</h3><p class="big">{health.pendingRequestCount ?? 0}</p></div>
      {#if $headlessEnabled}<div class="card"><h3>Headless</h3><p class="big">{health.headlessSessionCount ?? 0}</p></div>{/if}
      <div class="card"><h3>Goroutines</h3><p class="big">{health.goroutines}</p></div>
      <div class="card"><h3>Memory</h3><p class="big">{fmt(health.memoryAllocBytes)}</p></div>
      <div class="card"><h3>Heap</h3><p class="big">{fmt(health.memoryHeapBytes)}</p></div>
      <div class="card"><h3>GC Count</h3><p class="big">{health.gcCount}</p></div>
      <div class="card"><h3>Redis</h3><p class="big">{health.redisStatus ?? 'disabled'}</p></div>
    </div>
  {:else}
    <p>Loading…</p>
  {/if}
</div>

<style>
  .admin-page { padding: 1.5rem; }
  h1 { margin-top: 0; }
  .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 1rem; }
  .card { padding: 1rem; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  .card h3 { margin: 0; font-size: 0.875rem; color: var(--color-text-muted); }
  .big { font-size: 1.75rem; margin: 0.25rem 0 0; font-weight: 700; word-break: break-all; }
  .big.degraded { color: var(--color-error, #e53935); }
</style>
