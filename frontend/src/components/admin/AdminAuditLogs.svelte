<script lang="ts">
  import { onMount } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminAuditLogEntry } from '../../lib/types';
  import ActivityLog from '../activity/ActivityLog.svelte';

  // ── Admin audit log (admin actions) section ──
  let entries = $state<AdminAuditLogEntry[]>([]);
  let total = $state(0);
  let offset = $state(0);
  const limit = 50;
  let actionFilter = $state('');
  let error = $state('');

  async function load() {
    try {
      const r = await adminApi.listAuditLogs({ offset, limit, action: actionFilter || undefined });
      entries = r.entries;
      total = r.total;
    } catch (e: any) {
      error = e?.message ?? 'Failed';
    }
  }
  onMount(load);

  function nullStr(v: any): string {
    if (!v) return '';
    if (typeof v === 'string') return v;
    if (v.Valid) return v.String;
    return '';
  }

  // Bridge admin activity to ActivityLog's expected { ok, data, error } shape.
  // adminGet throws on error and returns data directly — wrap so ActivityLog can consume it.
  function adminActivityFetch(params: Record<string, unknown>) {
    return adminApi.getActivity(params as Parameters<typeof adminApi.getActivity>[0])
      .then(data => ({ ok: true as const, data, error: null }))
      .catch((e: any) => ({ ok: false as const, data: null, error: e?.message ?? 'Failed to load activity' }));
  }
</script>

<div class="admin-page">
  <!-- ── User Activity Feed ── -->
  <section class="section-activity">
    <h1>Activity Log</h1>
    <ActivityLog adminMode={true} fetchFn={adminActivityFetch} />
  </section>

  <!-- ── Admin Audit Log ── -->
  <section class="section-audit">
    <h2>Admin Actions ({total})</h2>
    <p class="muted">Actions taken by admin users on the dashboard (user disables, key deletions, etc.)</p>

    <div class="filters">
      <input type="text" bind:value={actionFilter} placeholder="Filter by action (e.g. user.disable)" />
      <button onclick={() => { offset = 0; load(); }}>Apply</button>
    </div>

    {#if error}<p class="error">{error}</p>{/if}

    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Time</th>
            <th>Admin</th>
            <th>Action</th>
            <th>Target</th>
            <th>IP</th>
          </tr>
        </thead>
        <tbody>
          {#each entries as e (e.id)}
            <tr>
              <td>{e.id}</td>
              <td>{e.createdAt}</td>
              <td>{e.adminUserId}</td>
              <td>{e.action}</td>
              <td>{e.targetType} {nullStr(e.targetId)}</td>
              <td>{nullStr(e.ipAddress)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <div class="pagination">
      <button disabled={offset === 0} onclick={() => { offset = Math.max(0, offset - limit); load(); }}>‹ Prev</button>
      <span>{offset + 1}–{Math.min(offset + limit, total)} of {total}</span>
      <button disabled={offset + limit >= total} onclick={() => { offset += limit; load(); }}>Next ›</button>
    </div>
  </section>
</div>

<style>
  .admin-page { padding: 1.5rem; }
  .section-activity { margin-bottom: 3rem; }
  .section-audit { }
  h1 { margin-top: 0; }
  h2 { font-size: 1rem; font-weight: 600; margin-bottom: 0.25rem; }
  .muted { color: var(--color-text-muted, #888); font-size: 0.8rem; margin-bottom: 0.75rem; }
  .filters { margin-bottom: 1rem; display: flex; gap: 0.5rem; }
  .filters input { padding: 0.5rem; flex: 1; max-width: 400px; }
  .table-wrap { overflow-x: auto; border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  table { width: 100%; border-collapse: collapse; font-size: 0.8125rem; }
  th, td { text-align: left; padding: 0.5rem 0.75rem; border-bottom: 1px solid var(--color-border-light); }
  th { background: var(--color-bg-elevated); }
  .pagination { display: flex; justify-content: center; gap: 1rem; margin-top: 1rem; align-items: center; }
  .error { color: var(--color-error, #e53935); }
</style>
