<script lang="ts">
  import { onMount } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminUserView } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';

  const SUBSCRIPTION_STATUSES = [
    'free', 'active', 'trialing', 'past_due', 'canceled', 'unpaid', 'incomplete', 'incomplete_expired',
  ];

  const SORT_COLUMNS = [
    { key: 'id', label: 'ID' },
    { key: 'email', label: 'Email' },
    { key: 'role', label: 'Role' },
    { key: 'disabled', label: 'Status' },
    { key: 'emailVerified', label: 'Verified' },
    { key: 'apiKeyRotationRequired', label: 'Rotation' },
    { key: 'subscriptionStatus', label: 'Subscription' },
    { key: 'requestsToday', label: 'Today' },
    { key: 'requestsThisMonth', label: 'Month' },
  ];

  let users = $state<AdminUserView[]>([]);
  let total = $state(0);
  let offset = $state(0);
  const limit = 25;
  let loading = $state(false);
  let error = $state('');
  let notice = $state('');

  let search = $state('');
  let roleFilter = $state('');
  let statusFilter = $state('');
  let verifiedFilter = $state('');
  let rotationFilter = $state('');
  let subFilter = $state('');
  let sortBy = $state('id');
  let sortDir = $state<'asc' | 'desc'>('asc');

  let selectedIds = $state(new Set<number>());
  const allSelected = $derived(users.length > 0 && selectedIds.size === users.length);

  function toggleSelect(id: number) {
    const next = new Set(selectedIds);
    if (next.has(id)) next.delete(id); else next.add(id);
    selectedIds = next;
  }

  function toggleSelectAll() {
    selectedIds = allSelected ? new Set() : new Set(users.map(u => u.id));
  }

  async function load() {
    loading = true;
    error = '';
    try {
      const r = await adminApi.listUsers({
        offset,
        limit,
        search: search.trim() || undefined,
        role: roleFilter || undefined,
        disabled: statusFilter === '' ? null : statusFilter === 'disabled',
        verified: verifiedFilter === '' ? null : verifiedFilter === 'true',
        rotation: rotationFilter === '' ? null : rotationFilter === 'true',
        subscription: subFilter || undefined,
        sort: sortBy,
        dir: sortDir,
      });
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

  // Reset to the first page and reload whenever a filter or sort changes.
  function applyFilters() {
    offset = 0;
    load();
  }

  let searchTimer: ReturnType<typeof setTimeout>;
  function onSearchInput(e: Event) {
    search = (e.currentTarget as HTMLInputElement).value;
    clearTimeout(searchTimer);
    searchTimer = setTimeout(applyFilters, 300);
  }

  function setSort(key: string) {
    if (sortBy === key) {
      sortDir = sortDir === 'asc' ? 'desc' : 'asc';
    } else {
      sortBy = key;
      sortDir = 'asc';
    }
    applyFilters();
  }

  function clearFilters() {
    search = '';
    roleFilter = '';
    statusFilter = '';
    verifiedFilter = '';
    rotationFilter = '';
    subFilter = '';
    applyFilters();
  }

  let modal = $state<{ open: boolean; action: () => Promise<void> }>({ open: false, action: async () => {} });
  let modalTitle = $state('');
  let modalMsg = $state('');
  let modalLabel = $state('Confirm');
  let modalDanger = $state(true);
  let rotatedKeyPrefix = $state('');

  async function runModal() {
    modal.open = false;
    error = '';
    try {
      await modal.action();
    } catch (e: any) {
      error = e?.message ?? 'Action failed';
    }
  }

  // Wraps a per-row mutation: surfaces errors and reloads so the row reflects
  // the authoritative server state (reverting the inline control on failure).
  async function mutate(fn: () => Promise<unknown>) {
    error = '';
    notice = '';
    try {
      await fn();
    } catch (e: any) {
      error = e?.message ?? 'Action failed';
    }
    await load();
  }

  function changeSubscription(u: AdminUserView, value: string) {
    if (value === u.subscriptionStatus) return;
    mutate(() => adminApi.updateUser(u.id, { subscriptionStatus: value }));
  }

  function changeVerified(u: AdminUserView, value: string) {
    const verified = value === 'true';
    if (verified === u.emailVerified) return;
    mutate(() => adminApi.updateUser(u.id, { emailVerified: verified }));
  }

  function changeStatus(u: AdminUserView, value: string) {
    const disabled = value === 'disabled';
    if (disabled === u.disabled) return;
    mutate(() => (disabled ? adminApi.disableUser(u.id) : adminApi.enableUser(u.id)));
  }

  function rotate(id: number) {
    modalTitle = 'Rotate API key'; modalMsg = "Rotate this user's API key? Their current key will stop working immediately."; modalLabel = 'Rotate'; modalDanger = true;
    modal = { open: true, action: async () => { const r = await adminApi.rotateUserKey(id); rotatedKeyPrefix = r.keyPrefix; await load(); } };
  }
  function flagRotation(id: number) {
    modalTitle = 'Force key rotation';
    modalMsg = 'This user will need to log in to the dashboard and regenerate their relay key before their Foundry module can reconnect. Continue?';
    modalLabel = 'Force Rotation'; modalDanger = true;
    modal = { open: true, action: async () => { await adminApi.flagUserRotation(id); await load(); } };
  }
  function sendReset(u: AdminUserView) {
    modalTitle = 'Send password reset';
    modalMsg = `Send a password reset email to ${u.email}? The reset link expires in 1 hour.`;
    modalLabel = 'Send email'; modalDanger = false;
    modal = { open: true, action: async () => { const r = await adminApi.sendPasswordReset(u.id); notice = r.message; } };
  }
  function del(id: number) {
    modalTitle = 'Delete user'; modalMsg = 'Permanently delete this user? This cannot be undone.'; modalLabel = 'Delete'; modalDanger = true;
    modal = { open: true, action: async () => { await adminApi.deleteUser(id); await load(); } };
  }

  function bulkDelete() {
    const n = selectedIds.size;
    modalTitle = `Delete ${n} user${n === 1 ? '' : 's'}`;
    modalMsg = `Permanently delete ${n} user account${n === 1 ? '' : 's'}? This cannot be undone.`;
    modalLabel = 'Delete all'; modalDanger = true;
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
    modalLabel = 'Rotate all'; modalDanger = true;
    modal = {
      open: true,
      action: async () => {
        await Promise.all([...selectedIds].map(id => adminApi.rotateUserKey(id)));
        await load();
      },
    };
  }
</script>

<ConfirmModal open={modal.open} title={modalTitle} message={modalMsg} confirmLabel={modalLabel} dangerous={modalDanger} onConfirm={runModal} onCancel={() => modal.open = false} />

<div class="admin-page">
  <h1>Users</h1>
  {#if rotatedKeyPrefix}
    <div class="info-banner">New key prefix: <code>{rotatedKeyPrefix}</code> <button onclick={() => rotatedKeyPrefix = ''}>✕</button></div>
  {/if}
  {#if notice}
    <div class="info-banner success">{notice} <button onclick={() => notice = ''}>✕</button></div>
  {/if}
  {#if error}<p class="error">{error}</p>{/if}

  <div class="filters">
    <input class="search" type="search" placeholder="Search ID or email…" value={search} oninput={onSearchInput} />
    <select bind:value={roleFilter} onchange={applyFilters} aria-label="Filter by role">
      <option value="">All roles</option>
      <option value="user">user</option>
      <option value="admin">admin</option>
    </select>
    <select bind:value={statusFilter} onchange={applyFilters} aria-label="Filter by status">
      <option value="">All statuses</option>
      <option value="active">Active</option>
      <option value="disabled">Disabled</option>
    </select>
    <select bind:value={verifiedFilter} onchange={applyFilters} aria-label="Filter by verified">
      <option value="">All verification</option>
      <option value="true">Verified</option>
      <option value="false">Unverified</option>
    </select>
    <select bind:value={rotationFilter} onchange={applyFilters} aria-label="Filter by rotation">
      <option value="">All rotation</option>
      <option value="true">Rotation required</option>
      <option value="false">No rotation</option>
    </select>
    <select bind:value={subFilter} onchange={applyFilters} aria-label="Filter by subscription">
      <option value="">All subscriptions</option>
      {#each SUBSCRIPTION_STATUSES as s}<option value={s}>{s}</option>{/each}
    </select>
    <button onclick={clearFilters}>Clear</button>
    {#if loading}<span class="muted">Loading…</span>{/if}
  </div>

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
          {#each SORT_COLUMNS as col}
            <th class="sortable" class:col-email={col.key === 'email'} onclick={() => setSort(col.key)}>
              {col.label}{#if sortBy === col.key}<span class="arrow">{sortDir === 'asc' ? '▲' : '▼'}</span>{/if}
            </th>
          {/each}
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each users as u (u.id)}
          <tr class:disabled={u.disabled} class:selected={selectedIds.has(u.id)}>
            <td><input type="checkbox" checked={selectedIds.has(u.id)} onchange={() => toggleSelect(u.id)} /></td>
            <td>{u.id}</td>
            <td class="col-email" title={u.email}>{u.email}</td>
            <td>{u.role}</td>
            <td>
              <select value={u.disabled ? 'disabled' : 'active'} onchange={(e) => changeStatus(u, e.currentTarget.value)}>
                <option value="active">Active</option>
                <option value="disabled">Disabled</option>
              </select>
            </td>
            <td>
              <select value={u.emailVerified ? 'true' : 'false'} onchange={(e) => changeVerified(u, e.currentTarget.value)}>
                <option value="true">Yes</option>
                <option value="false">No</option>
              </select>
            </td>
            <td>{#if u.apiKeyRotationRequired}<span class="rotation-badge" title="User must regenerate their relay key">⚠ Required</span>{:else}—{/if}</td>
            <td>
              <select class="sub-select" value={u.subscriptionStatus} onchange={(e) => changeSubscription(u, e.currentTarget.value)}>
                {#if !SUBSCRIPTION_STATUSES.includes(u.subscriptionStatus)}
                  <option value={u.subscriptionStatus}>{u.subscriptionStatus}</option>
                {/if}
                {#each SUBSCRIPTION_STATUSES as s}<option value={s}>{s}</option>{/each}
              </select>
            </td>
            <td>{u.requestsToday}</td>
            <td>{u.requestsThisMonth}</td>
            <td class="actions">
              <button onclick={() => sendReset(u)}>Send Reset</button>
              <button onclick={() => rotate(u.id)}>Rotate Key</button>
              {#if !u.apiKeyRotationRequired}
                <button onclick={() => flagRotation(u.id)}>Flag Rotation</button>
              {/if}
              <button class="danger" onclick={() => del(u.id)}>Delete</button>
            </td>
          </tr>
        {/each}
        {#if users.length === 0 && !loading}
          <tr><td colspan="11" class="empty">No users match the current filters.</td></tr>
        {/if}
      </tbody>
    </table>
  </div>

  <div class="pagination">
    <button disabled={offset === 0} onclick={() => { offset = Math.max(0, offset - limit); load(); }}>‹ Prev</button>
    <span>{total === 0 ? 0 : offset + 1}–{Math.min(offset + limit, total)} of {total}</span>
    <button disabled={offset + limit >= total} onclick={() => { offset += limit; load(); }}>Next ›</button>
  </div>
</div>

<style>
  .admin-page { padding: 1.5rem; }
  h1 { margin-top: 0; }
  .filters { display: flex; flex-wrap: wrap; gap: 0.5rem; align-items: center; margin-bottom: 0.75rem; }
  .filters .search { flex: 1 1 220px; min-width: 180px; padding: 0.35rem 0.5rem; }
  .filters select { padding: 0.35rem 0.5rem; }
  .filters .muted { color: var(--color-text-muted, #888); font-size: 0.875rem; }
  .table-wrap { overflow-x: auto; border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  table { width: 100%; border-collapse: collapse; }
  th, td { text-align: left; padding: 0.5rem 0.75rem; border-bottom: 1px solid var(--color-border-light); font-size: 0.875rem; }
  th { background: var(--color-bg-elevated); font-weight: 600; }
  th.sortable { cursor: pointer; user-select: none; white-space: nowrap; }
  th.sortable:hover { color: var(--color-primary, #4f46e5); }
  .arrow { margin-left: 0.25rem; font-size: 0.7rem; }
  td select { padding: 0.2rem 0.3rem; font-size: 0.8rem; }
  .col-email { max-width: 180px; }
  td.col-email { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .sub-select { max-width: 108px; }
  td.actions { white-space: nowrap; }
  tr.disabled td { opacity: 0.5; }
  tr.disabled td select { opacity: 1; }
  .rotation-badge { color: var(--color-warning, #f59e0b); font-size: 0.75rem; font-weight: 600; white-space: nowrap; }
  tr.selected td { background: color-mix(in srgb, var(--color-primary, #4f46e5) 6%, transparent); }
  td.empty { text-align: center; color: var(--color-text-muted, #888); padding: 1.5rem; }
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
  .info-banner.success { background: color-mix(in srgb, var(--color-success, #22c55e) 12%, transparent); }
</style>
