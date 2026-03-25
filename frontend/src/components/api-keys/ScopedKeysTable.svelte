<script lang="ts">
  import type { ScopedKey } from '../../lib/types';
  import { updateScopedKey, deleteScopedKey as apiDeleteKey } from '../../lib/api';

  interface Props {
    keys: ScopedKey[];
    onEdit: (key: ScopedKey) => void;
    onRefresh: () => void;
  }

  let { keys, onEdit, onRefresh }: Props = $props();

  function getStatus(key: ScopedKey): 'active' | 'disabled' | 'expired' {
    if (!key.enabled) return 'disabled';
    if (key.isExpired) return 'expired';
    return 'active';
  }

  async function handleToggle(key: ScopedKey) {
    await updateScopedKey(key.id, { enabled: !key.enabled });
    onRefresh();
  }

  async function handleDelete(key: ScopedKey) {
    if (!confirm(`Delete scoped key "${key.name}"? This cannot be undone.`)) return;
    await apiDeleteKey(key.id);
    onRefresh();
  }

  function formatScope(key: ScopedKey): string {
    const parts: string[] = [];
    if (key.scopedClientId) parts.push(`Client: ${key.scopedClientId.substring(0, 12)}...`);
    if (key.scopedUserId) parts.push(`User: ${key.scopedUserId}`);
    return parts.join(', ') || 'Unrestricted';
  }
</script>

{#if keys.length === 0}
  <p class="text-muted" style="font-size: 0.875rem;">No scoped keys yet. Create one to get started.</p>
{:else}
  <div class="table-wrapper">
    <table class="table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Key</th>
          <th>Status</th>
          <th>Scopes</th>
          <th>Daily</th>
          <th>Expires</th>
          <th>Creds</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each keys as key (key.id)}
          {@const status = getStatus(key)}
          <tr>
            <td>{key.name}</td>
            <td><code class="inline-code">{key.key}</code></td>
            <td>
              <span class="badge" class:badge-active={status === 'active'} class:badge-disabled={status === 'disabled'} class:badge-expired={status === 'expired'}>
                {status}
              </span>
            </td>
            <td class="text-muted" style="font-size: 0.8rem;">{formatScope(key)}</td>
            <td>
              {#if key.dailyLimit}
                {key.requestsToday}/{key.dailyLimit}
              {:else}
                <span class="text-muted">Unlimited</span>
              {/if}
            </td>
            <td>
              {#if key.expiresAt}
                {new Date(key.expiresAt).toLocaleDateString()}
              {:else}
                <span class="text-muted">Never</span>
              {/if}
            </td>
            <td style="text-align: center;">
              {#if key.hasFoundryCredentials}&#10003;{/if}
            </td>
            <td class="actions-cell">
              <button class="btn btn-sm btn-ghost" onclick={() => onEdit(key)}>Edit</button>
              <button class="btn btn-sm btn-ghost" onclick={() => handleToggle(key)}>
                {key.enabled ? 'Disable' : 'Enable'}
              </button>
              <button class="btn btn-sm btn-danger" onclick={() => handleDelete(key)}>Delete</button>
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

  .actions-cell {
    white-space: nowrap;
  }
</style>
