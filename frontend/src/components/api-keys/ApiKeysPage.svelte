<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../lib/auth';
  import { fetchScopedKeys } from '../../lib/api';
  import type { ScopedKey } from '../../lib/types';
  import ScopedKeyForm from './ScopedKeyForm.svelte';
  import ScopedKeysTable from './ScopedKeysTable.svelte';

  let keys = $state<ScopedKey[]>([]);
  let loading = $state(true);
  let showForm = $state(false);
  let editingKey = $state<ScopedKey | null>(null);
  let createdKeyValue = $state<string | null>(null);
  let masterKeyMasked = $state(true);
  let copyLabel = $state('Copy');

  let masterKey = $derived($user?.apiKey ?? '');

  onMount(() => {
    loadKeys();
  });

  async function loadKeys() {
    loading = true;
    const result = await fetchScopedKeys();
    loading = false;
    if (result.ok) {
      keys = result.data.keys || [];
    }
  }

  function handleCreate() {
    editingKey = null;
    showForm = true;
    createdKeyValue = null;
  }

  function handleEdit(key: ScopedKey) {
    editingKey = key;
    showForm = true;
    createdKeyValue = null;
  }

  function handleSave(keyValue?: string) {
    showForm = false;
    editingKey = null;
    if (keyValue) createdKeyValue = keyValue;
    loadKeys();
  }

  function handleCancel() {
    showForm = false;
    editingKey = null;
  }

  function maskMasterKey(key: string): string {
    if (!key || key.length < 12) return '••••••••';
    return key.substring(0, 8) + '...';
  }

  async function copyMasterKey() {
    if (!masterKey) return;
    await navigator.clipboard.writeText(masterKey);
    copyLabel = 'Copied!';
    setTimeout(() => { copyLabel = 'Copy'; }, 1500);
  }

  async function copyCreatedKey() {
    if (!createdKeyValue) return;
    await navigator.clipboard.writeText(createdKeyValue);
    alert('Key copied!');
  }
</script>

<h2 class="page-title">API Keys</h2>

<!-- Master Key -->
<div class="card">
  <h3 class="card-title mb-1">Master API Key</h3>
  <p class="text-muted mb-1" style="font-size: 0.85rem;">Your master key has full access. Manage it from the Dashboard tab.</p>
  <div class="key-display">
    <code class="key-value">{masterKeyMasked ? maskMasterKey(masterKey) : masterKey}</code>
    <button class="btn btn-sm btn-secondary" onclick={() => masterKeyMasked = !masterKeyMasked}>
      {masterKeyMasked ? 'Show' : 'Hide'}
    </button>
    <button class="btn btn-sm btn-secondary" onclick={copyMasterKey}>{copyLabel}</button>
  </div>
</div>

<!-- Scoped Keys -->
<div class="card">
  <div class="card-header">
    <h3 class="card-title">Scoped API Keys</h3>
    {#if !showForm}
      <button class="btn btn-primary btn-sm" onclick={handleCreate}>+ Create Scoped Key</button>
    {/if}
  </div>

  {#if showForm}
    <ScopedKeyForm editKey={editingKey} onSave={handleSave} onCancel={handleCancel} />
  {/if}

  {#if createdKeyValue}
    <div class="alert alert-success mb-2" style="border: 2px solid var(--color-success);">
      <strong>Key created!</strong> Copy it now — it won't be shown again.
      <div class="created-key-row">
        <code class="key-value" style="flex: 1;">{createdKeyValue}</code>
        <button class="btn btn-sm btn-secondary" onclick={copyCreatedKey}>Copy</button>
      </div>
    </div>
  {/if}

  {#if loading}
    <p class="text-muted">Loading scoped keys...</p>
  {:else}
    <ScopedKeysTable {keys} onEdit={handleEdit} onRefresh={loadKeys} />
  {/if}
</div>

<style>
  .created-key-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }
</style>
