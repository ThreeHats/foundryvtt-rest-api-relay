<script lang="ts">
  import { onMount } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminUserView } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';

  let users = $state<AdminUserView[]>([]);
  let total = $state(0);
  let offset = $state(0);
  const limit = 25;
  let loading = $state(false);
  let error = $state('');

  let selectedIds = $state(new Set<number>());
  const allSelected = $derived(users.length > 0 && selectedIds.size === users.length);

  function toggleSelect(id: number) {
    const next = new Set(selectedIds);
    if (next.has(id)) next.delete(id); else next.add(id);
    selectedIds = next;
  }

  function toggleSelectAll() {
    if (allSelected) {
      selectedIds = new Set();
    } else {
      selectedIds = new Set(users.map(u => u.id));
    }
  }

  async function load() {
    loading = true;
    error = '';
    try {
      const r = await adminApi.listUsers(offset, limit);
      users = r.users;
      total = r.total;
      selectedIds = new Set();
    } catch (e: any) {
      error = e?.message ?? 'Failed to load';
    } finally {
      loading = false;
    }
  }

  onMount(load);

  let modal = $state<{ open: boolean; action: () => Promise<void> }>({ open: false, action: async () => {} });
  let modalTitle = $state('');
  let modalMsg = $state('');
  let modalLabel = $state('Confirm');
  let rotatedKeyPrefix = $state('');
  async function runModal() { modal.open = false; await modal.action(); }

  function disable(id: number) {
    modalTitle = 'Disable user'; modalMsg = 'Disable this user? They will not be able to authenticate until re-enabled.'; modalLabel = 'Disable';
    modal = { open: true, action: async () => { await adminApi.disableUser(id); await load(); } };
  }
  async function enable(id: number) {
    await adminApi.enableUser(id);
    await load();
  }
  function rotate(id: number) {
    modalTitle = 'Rotate API key'; modalMsg = "Rotate this user's API key? Their current key will stop working immediately."; modalLabel = 'Rotate';
    modal = { open: true, action: async () => { const r = await adminApi.rotateUserKey(id); rotatedKeyPrefix = r.keyPrefix; await load(); } };
  }
  function flagRotation(id: number) {
    modalTitle = 'Force key rotation';
    modalMsg = 'This user will be blocked from all API calls until they log in to the dashboard and regenerate their master key. Continue?';
    modalLabel = 'Force Rotation';
    modal = { open: true, action: async () => { await adminApi.flagUserRotation(id); await load(); } };
  }
  function del(id: number) {
    modalTitle = 'Delete user'; modalMsg = 'Permanently delete this user? This cannot be undone.'; modalLabel = 'Delete';
    modal = { open: true, action: async () => { await adminApi.deleteUser(id); await load(); } };
  }

  function bulkDelete() {
    const n = selectedIds.size;
    modalTitle = `Delete ${n} user${n === 1 ? '' : 's'}`;
    modalMsg = `Permanently delete ${n} user account${n === 1 ? '' : 's'}? This cannot be undone.`;
    modalLabel = 'Delete all';
    modal = {
      open: true,
      action: async () => {
        await Promise.all([...selectedIds].map(id => adminApi.deleteUser(id)));
        await load();
      },
    };
  }

  function bulkRotate() {
    const n = selectedIds.size;
    modalTitle = `Rotate ${n} API key${n === 1 ? '' : 's'}`;
    modalMsg = `Rotate API keys for ${n} user${n === 1 ? '' : 's'}? Their current keys will stop working immediately.`;
    modalLabel = 'Rotate all';
    modal = {
      open: true,
      action: async () => {
        await Promise.all([...selectedIds].map(id => adminApi.rotateUserKey(id)));
        await load();
      },
    };
  }
</script>

<ConfirmModal open={modal.open} title={modalTitle} message={modalMsg} confirmLabel={modalLabel} dangerous={true} onConfirm={runModal} onCancel={() => modal.open = false} />

<div class="admin-page">
  <h1>Users</h1>
  {#if rotatedKeyPrefix}
    <div class="info-banner">New key prefix: <code>{rotatedKeyPrefix}</code> <button onclick={() => rotatedKeyPrefix = ''}>✕</button></div>
  {/if}
  {#if error}<p class="error">{error}</p>{/if}
  {#if loading}<p>Loading…</p>{/if}

  {#if selectedIds.size > 0}
    <div class="bulk-bar">
      <span>{selectedIds.size} selected</span>
      <button onclick={bulkRotate}>Rotate Keys ({selectedIds.size})</button>
      <button class="danger" onclick={bulkDelete}>Delete ({selectedIds.size})</button>
      <button onclick={() => selectedIds = new Set()}>Clear</button>
    </div>
  {/if}

  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th><input type="checkbox" checked={allSelected} onchange={toggleSelectAll} /></th>
          <th>ID</th>
          <th>Email</th>
          <th>Role</th>
          <th>Status</th>
          <th>Rotation</th>
          <th>Subscription</th>
          <th>Today</th>
          <th>Month</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each users as u (u.id)}
          <tr class:disabled={u.disabled} class:selected={selectedIds.has(u.id)}>
            <td><input type="checkbox" checked={selectedIds.has(u.id)} onchange={() => toggleSelect(u.id)} /></td>
            <td>{u.id}</td>
            <td>{u.email}</td>
            <td>{u.role}</td>
            <td>{u.disabled ? 'Disabled' : 'Active'}</td>
            <td>{#if u.apiKeyRotationRequired}<span class="rotation-badge" title="User must regenerate their master key">⚠ Required</span>{:else}—{/if}</td>
            <td>{u.subscriptionStatus}</td>
            <td>{u.requestsToday}</td>
            <td>{u.requestsThisMonth}</td>
            <td>
              {#if u.disabled}
                <button onclick={() => enable(u.id)}>Enable</button>
              {:else}
                <button onclick={() => disable(u.id)}>Disable</button>
              {/if}
              <button onclick={() => rotate(u.id)}>Rotate Key</button>
              {#if !u.apiKeyRotationRequired}
                <button onclick={() => flagRotation(u.id)}>Flag Rotation</button>
              {/if}
              <button class="danger" onclick={() => del(u.id)}>Delete</button>
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
  table { width: 100%; border-collapse: collapse; }
  th, td { text-align: left; padding: 0.5rem 0.75rem; border-bottom: 1px solid var(--color-border-light); font-size: 0.875rem; }
  th { background: var(--color-bg-elevated); font-weight: 600; }
  tr.disabled td { opacity: 0.5; }
  .rotation-badge { color: var(--color-warning, #f59e0b); font-size: 0.75rem; font-weight: 600; white-space: nowrap; }
  tr.selected td { background: color-mix(in srgb, var(--color-primary, #4f46e5) 6%, transparent); }
  button { margin-right: 0.25rem; padding: 0.25rem 0.5rem; font-size: 0.75rem; }
  button.danger { background: var(--color-error, #e53935); color: white; border: none; }
  .pagination { display: flex; justify-content: center; gap: 1rem; margin-top: 1rem; align-items: center; }
  .error { color: var(--color-error, #e53935); }
  .bulk-bar {
    display: flex; align-items: center; gap: 0.5rem;
    padding: 0.5rem 0.75rem; margin-bottom: 0.75rem;
    background: var(--color-bg-elevated); border: 1px solid var(--color-border);
    border-radius: var(--radius-md); font-size: 0.875rem;
  }
  .info-banner { padding: 0.5rem 0.75rem; background: var(--color-bg-elevated); border-radius: var(--radius-md); margin-bottom: 0.75rem; font-size: 0.875rem; }
</style>
