<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchRemoteRequestLogs } from '../../lib/api';
  import type { RemoteRequestLog } from '../../lib/types';

  let logs    = $state<RemoteRequestLog[]>([]);
  let total   = $state(0);
  let offset  = $state(0);
  const limit = 50;
  let loading = $state(true);

  onMount(() => { loadLogs(); });

  async function loadLogs() {
    loading = true;
    const r = await fetchRemoteRequestLogs(limit, offset);
    loading = false;
    if (r.ok) { logs = r.data.logs || []; total = r.data.total || 0; }
  }

  function next() { if (offset + limit < total) { offset += limit; loadLogs(); } }
  function prev() { if (offset > 0) { offset = Math.max(0, offset - limit); loadLogs(); } }
</script>

<h2 class="page-title">Cross-World Operations</h2>

<div class="card">
  <div class="card-header">
    <div>
      <h3 class="card-title">Remote Request Audit Log</h3>
      <p class="text-muted" style="font-size: 0.85rem;">
        Every cross-world operation tunneled through the relay via <code>remote-request</code> is logged here.
      </p>
    </div>
    <button class="btn btn-sm btn-secondary" onclick={loadLogs} disabled={loading}>&#8635; Refresh</button>
  </div>

  {#if loading}
    <p class="text-muted">Loading logs...</p>
  {:else if logs.length === 0}
    <p class="text-muted" style="font-size: 0.875rem;">No cross-world operations recorded yet. These appear when a Foundry module uses <code>remote-request</code> to invoke actions on another world.</p>
  {:else}
    <div class="table-wrapper">
      <table class="table">
        <thead>
          <tr>
            <th>Timestamp</th>
            <th>Source World</th>
            <th>Target World</th>
            <th>Action</th>
            <th>Result</th>
            <th>Source IP</th>
          </tr>
        </thead>
        <tbody>
          {#each logs as log (log.id)}
            <tr class:row-error={!log.success}>
              <td style="white-space: nowrap;">{log.createdAt ? new Date(log.createdAt).toLocaleString() : '—'}</td>
              <td><code class="inline-code" style="font-size: 0.7rem;">{log.sourceClientId}</code></td>
              <td><code class="inline-code" style="font-size: 0.7rem;">{log.targetClientId}</code></td>
              <td><code class="inline-code">{log.action}</code></td>
              <td>
                {#if log.success}
                  <span class="badge badge-active">ok</span>
                {:else}
                  <span class="badge badge-expired" title={log.errorMessage || ''}>
                    ✗ {log.errorMessage || 'failed'}
                  </span>
                {/if}
              </td>
              <td><code class="inline-code" style="font-size: 0.7rem;">{log.sourceIp || '—'}</code></td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <div class="flex items-center justify-between mt-1">
      <span class="text-muted" style="font-size: 0.8rem;">
        Showing {offset + 1}–{offset + logs.length} of {total}
      </span>
      <div class="flex gap-1">
        <button class="btn btn-sm btn-secondary" onclick={prev} disabled={offset === 0}>Previous</button>
        <button class="btn btn-sm btn-secondary" onclick={next} disabled={logs.length < limit}>Next</button>
      </div>
    </div>
  {/if}
</div>

<style>
  .table-wrapper { overflow-x: auto; }
  .row-error td { color: var(--color-error, #ef4444); }
</style>
