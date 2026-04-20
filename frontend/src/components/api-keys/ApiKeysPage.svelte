<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchScopedKeys } from '../../lib/api';
  import type { ScopedKey } from '../../lib/types';
  import ScopedKeyForm from './ScopedKeyForm.svelte';
  import ScopedKeysTable from './ScopedKeysTable.svelte';

  let keys = $state<ScopedKey[]>([]);
  let loading = $state(true);
  let showForm = $state(false);
  let editingKey = $state<ScopedKey | null>(null);
  let createdKeyValue = $state<string | null>(null);

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
    if (keyValue) { createdKeyValue = keyValue; keyBannerLabel = 'Key created!'; copyLabel = 'Copy'; }
    loadKeys();
  }

  function handleRegenerate(keyValue: string) {
    createdKeyValue = keyValue;
    keyBannerLabel = 'Key regenerated!';
    copyLabel = 'Copy';
    loadKeys();
  }

  function handleCancel() {
    showForm = false;
    editingKey = null;
  }

  let keyBannerLabel = $state('Key created!');
  let copyLabel = $state('Copy');

  async function copyCreatedKey() {
    if (!createdKeyValue) return;
    await navigator.clipboard.writeText(createdKeyValue);
    copyLabel = 'Copied!';
    setTimeout(() => { copyLabel = 'Copy'; }, 1500);
  }
</script>

<h2 class="page-title">API Keys</h2>

<div class="card">
  {#if !showForm}
    <div class="mb-2">
      <button class="btn btn-primary btn-sm" onclick={handleCreate}>+ Create Scoped Key</button>
    </div>
  {/if}

  {#if showForm}
    <ScopedKeyForm editKey={editingKey} onSave={handleSave} onCancel={handleCancel} />
  {/if}

  {#if createdKeyValue}
    <div class="alert alert-success mb-2" style="border: 2px solid var(--color-success);">
      <strong>{keyBannerLabel}</strong> Copy it now — it won't be shown again.
      <div class="created-key-row">
        <code class="key-value" style="flex: 1;">{createdKeyValue}</code>
        <button class="btn btn-sm btn-secondary" onclick={copyCreatedKey}>{copyLabel}</button>
      </div>
    </div>
  {/if}

  {#if loading}
    <p class="text-muted">Loading scoped keys...</p>
  {:else if !showForm}
    <ScopedKeysTable {keys} onEdit={handleEdit} onRefresh={loadKeys} onRegenerate={handleRegenerate} />
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
