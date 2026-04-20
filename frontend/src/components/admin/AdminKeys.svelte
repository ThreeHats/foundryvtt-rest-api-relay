<script lang="ts">
  import { onMount } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminApiKeyView } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';

  let keys = $state<AdminApiKeyView[]>([]);
  let total = $state(0);
  let offset = $state(0);
  const limit = 25;
  let error = $state('');

  async function load() {
    try {
      const r = await adminApi.listKeys(offset, limit);
      keys = r.keys;
      total = r.total;
    } catch (e: any) {
      error = e?.message ?? 'Failed';
    }
  }
  onMount(load);

  let modal = $state<{ open: boolean; action: () => Promise<void> }>({ open: false, action: async () => {} });
  let modalMsg = $state('');
  async function runModal() { modal.open = false; await modal.action(); }

  function revoke(id: number) {
    modalMsg = 'Revoke this API key? It will stop working immediately.';
    modal = { open: true, action: async () => { await adminApi.revokeKey(id); await load(); } };
  }
  async function toggle(k: AdminApiKeyView) {
    await adminApi.patchKey(k.id, { enabled: !k.enabled });
    await load();
  }
</script>

<ConfirmModal open={modal.open} title="Confirm" message={modalMsg} confirmLabel="Confirm" dangerous={true} onConfirm={runModal} onCancel={() => modal.open = false} />

<div class="admin-page">
  <h1>All API Keys ({total})</h1>
  {#if error}<p class="error">{error}</p>{/if}
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>User ID</th>
          <th>Name</th>
          <th>Prefix</th>
          <th>Status</th>
          <th>Scopes</th>
          <th>Today</th>
          <th>Created</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each keys as k (k.id)}
          <tr>
            <td>{k.id}</td>
            <td>{k.userId}</td>
            <td>{k.name}</td>
            <td><code>{k.keyPrefix}</code></td>
            <td>{k.enabled ? 'Enabled' : 'Disabled'} {k.isExpired ? '(expired)' : ''}</td>
            <td>{(k.scopes ?? []).join(', ') || '(all)'}</td>
            <td>{k.requestsToday}</td>
            <td>{k.createdAt}</td>
            <td>
              <button onclick={() => toggle(k)}>{k.enabled ? 'Disable' : 'Enable'}</button>
              <button class="danger" onclick={() => revoke(k.id)}>Revoke</button>
            </td>
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
</div>

<style>
  .admin-page { padding: 1.5rem; }
  h1 { margin-top: 0; }
  .table-wrap { overflow-x: auto; border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  table { width: 100%; border-collapse: collapse; font-size: 0.8125rem; }
  th, td { text-align: left; padding: 0.5rem 0.75rem; border-bottom: 1px solid var(--color-border-light); }
  th { background: var(--color-bg-elevated); }
  button { padding: 0.25rem 0.5rem; font-size: 0.75rem; margin-right: 0.25rem; }
  button.danger { background: var(--color-error, #e53935); color: white; border: none; }
  .pagination { display: flex; justify-content: center; gap: 1rem; margin-top: 1rem; align-items: center; }
  .error { color: var(--color-error, #e53935); }
</style>
