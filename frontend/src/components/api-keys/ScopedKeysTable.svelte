<script lang="ts">
  import type { ScopedKey } from '../../lib/types';
  import { updateScopedKey, deleteScopedKey as apiDeleteKey, regenerateScopedKey } from '../../lib/api';
  import ConfirmModal from '../ui/ConfirmModal.svelte';

  interface Props {
    keys: ScopedKey[];
    onEdit: (key: ScopedKey) => void;
    onRefresh: () => void;
    onRegenerate: (keyValue: string) => void;
  }

  let { keys, onEdit, onRefresh, onRegenerate }: Props = $props();

  // ── Selection (bulk delete) ───────────────────────────────────────────────
  let selectedIds = $state(new Set<number>());

  function toggleSelect(id: number) {
    const next = new Set(selectedIds);
    if (next.has(id)) next.delete(id); else next.add(id);
    selectedIds = next;
  }

  function toggleAll() {
    selectedIds = selectedIds.size === keys.length
      ? new Set()
      : new Set(keys.map(k => k.id));
  }

  // ── Confirmation modal ────────────────────────────────────────────────────
  let modal = $state<{ open: boolean; title: string; message: string; confirmLabel: string; action: () => Promise<void> }>({
    open: false, title: '', message: '', confirmLabel: 'Confirm', action: async () => {},
  });

  function ask(title: string, message: string, confirmLabel: string, action: () => Promise<void>) {
    modal = { open: true, title, message, confirmLabel, action };
  }

  async function runModal() {
    modal.open = false;
    await modal.action();
  }

  // ── Actions ───────────────────────────────────────────────────────────────
  function handleDelete(key: ScopedKey) {
    ask(
      'Delete API key',
      `Delete scoped key "${key.name}"? This cannot be undone.`,
      'Delete',
      async () => { await apiDeleteKey(key.id); onRefresh(); },
    );
  }

  function handleRegenerate(key: ScopedKey) {
    ask(
      'Regenerate API key',
      `Regenerate key "${key.name}"? Anything using the current value will stop working immediately.`,
      'Regenerate',
      async () => {
        const r = await regenerateScopedKey(key.id);
        if (r.ok) onRegenerate(r.data.key);
      },
    );
  }

  function handleBulkDelete() {
    const n = selectedIds.size;
    ask(
      `Delete ${n} key${n === 1 ? '' : 's'}`,
      `Permanently delete ${n} scoped key${n === 1 ? '' : 's'}? This cannot be undone.`,
      'Delete all',
      async () => {
        await Promise.all([...selectedIds].map(id => apiDeleteKey(id)));
        selectedIds = new Set();
        onRefresh();
      },
    );
  }

  function formatRestrictions(key: ScopedKey): string {
    const parts: string[] = [];
    if (key.scopedClientId) parts.push(`Client: ${key.scopedClientId.substring(0, 12)}…`);
    if (key.scopedUserId) parts.push(`User: ${key.scopedUserId}`);
    return parts.join(', ');
  }

  let allSelected = $derived(keys.length > 0 && selectedIds.size === keys.length);
  let someSelected = $derived(selectedIds.size > 0);

  // ── Scope expansion ───────────────────────────────────────────────────────
  let expandedScopeRows = $state(new Set<number>());

  function toggleScopeExpand(id: number) {
    const next = new Set(expandedScopeRows);
    if (next.has(id)) next.delete(id); else next.add(id);
    expandedScopeRows = next;
  }
</script>

<ConfirmModal
  open={modal.open}
  title={modal.title}
  message={modal.message}
  confirmLabel={modal.confirmLabel}
  dangerous={true}
  onConfirm={runModal}
  onCancel={() => modal.open = false}
/>

{#if keys.length === 0}
  <p class="text-muted" style="font-size: 0.875rem;">No scoped keys yet. Create one to get started.</p>
{:else}
  {#if someSelected}
    <div class="bulk-bar">
      <span class="text-muted" style="font-size: 0.875rem;">{selectedIds.size} selected</span>
      <button class="btn btn-sm btn-danger" onclick={handleBulkDelete}>Delete {selectedIds.size} selected</button>
    </div>
  {/if}

  <div class="table-wrapper">
    <table class="table">
      <thead>
        <tr>
          <th class="col-check">
            <input type="checkbox" checked={allSelected} onchange={toggleAll} title="Select all" />
          </th>
          <th>Name</th>
          <th>Key</th>
          <th>Permissions</th>
          <th>Daily</th>
          <th>Expires</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each keys as key (key.id)}
          {@const restrictions = formatRestrictions(key)}
          <tr class:row-selected={selectedIds.has(key.id)}>
            <td class="col-check">
              <input type="checkbox" checked={selectedIds.has(key.id)} onchange={() => toggleSelect(key.id)} />
            </td>
            <td>
              <div class="name-cell">
                <span>{key.name}</span>
                {#if restrictions}
                  <span class="cell-sub">{restrictions}</span>
                {/if}
                {#if key.hasFoundryCredentials}
                  <span class="cell-sub">Has credentials</span>
                {/if}
              </div>
            </td>
            <td>
              <code class="inline-code key-truncated">{key.key}</code>
            </td>
            <td>
              {#if key.scopes && key.scopes.length > 0}
                {@const expanded = expandedScopeRows.has(key.id)}
                {@const overflow = key.scopes.length - 3}
                <div class="scope-badges">
                  {#each (expanded ? key.scopes : key.scopes.slice(0, 3)) as scope}
                    <span class="scope-badge">{scope}</span>
                  {/each}
                  {#if !expanded && overflow > 0}
                    <button class="scope-badge scope-badge-more" onclick={() => toggleScopeExpand(key.id)}>+{overflow} more</button>
                  {:else if expanded && overflow > 0}
                    <button class="scope-badge scope-badge-more" onclick={() => toggleScopeExpand(key.id)}>show less</button>
                  {/if}
                </div>
              {:else}
                <span class="scope-badge scope-badge-all">All</span>
              {/if}
            </td>
            <td>
              {#if key.dailyLimit}
                {key.requestsToday}/{key.dailyLimit}
              {:else}
                <span class="text-muted">—</span>
              {/if}
            </td>
            <td>
              {#if key.expiresAt}
                {new Date(key.expiresAt).toLocaleDateString()}
              {:else}
                <span class="text-muted">—</span>
              {/if}
            </td>
            <td class="actions-cell">
              <button class="btn btn-sm btn-ghost" onclick={() => onEdit(key)}>Edit</button>
              <button class="btn btn-sm btn-ghost" onclick={() => handleRegenerate(key)}>Regen</button>
              <button class="btn btn-sm btn-danger" onclick={() => handleDelete(key)}>Del</button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}

<style>
  .table-wrapper {
    overflow-x: auto;
  }

  .col-check {
    width: 2rem;
    text-align: center;
  }

  .actions-cell {
    white-space: nowrap;
  }

  .bulk-bar {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.5rem;
  }

  .row-selected {
    background: color-mix(in srgb, var(--color-primary, #4f46e5) 6%, transparent);
  }

  .name-cell {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .cell-sub {
    font-size: 0.75rem;
    color: var(--color-text-muted, #9ca3af);
  }

  .key-truncated {
    font-size: 0.75rem;
  }

  .scope-badges {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .scope-badge {
    display: inline-block;
    padding: 0.125rem 0.4rem;
    font-size: 0.7rem;
    border-radius: 3px;
    background: var(--color-surface-hover, rgba(255, 255, 255, 0.08));
    color: var(--color-text-secondary, #aaa);
    border: 1px solid var(--color-border, rgba(255, 255, 255, 0.1));
    white-space: nowrap;
  }

  .scope-badge-more {
    cursor: pointer;
    background: none;
    border: 1px dashed var(--color-border, rgba(255, 255, 255, 0.1));
    color: var(--color-text-muted, #9ca3af);
    font-style: italic;
  }

  .scope-badge-all {
    background: rgba(234, 179, 8, 0.15);
    color: rgb(234, 179, 8);
    border-color: rgba(234, 179, 8, 0.3);
  }
</style>
