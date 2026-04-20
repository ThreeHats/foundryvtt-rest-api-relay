<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchCredentials, createCredential, updateCredential, deleteCredential } from '../../lib/api';
  import type { Credential } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';

  let credentials = $state<Credential[]>([]);
  let loading = $state(true);
  let showForm = $state(false);
  let editingCredential = $state<Credential | null>(null);

  let name = $state('');
  let foundryUrl = $state('');
  let foundryUsername = $state('');
  let foundryPassword = $state('');
  let saving = $state(false);
  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');

  let isEdit = $derived(!!editingCredential);

  onMount(() => {
    loadCredentials();
  });

  async function loadCredentials() {
    loading = true;
    const result = await fetchCredentials();
    loading = false;
    if (result.ok) {
      credentials = result.data.credentials || [];
    }
  }

  function handleCreate() {
    editingCredential = null;
    name = '';
    foundryUrl = '';
    foundryUsername = '';
    foundryPassword = '';
    showForm = true;
    message = '';
  }

  function handleEdit(cred: Credential) {
    editingCredential = cred;
    name = cred.name;
    foundryUrl = cred.foundryUrl;
    foundryUsername = cred.foundryUsername;
    foundryPassword = '';
    showForm = true;
    message = '';
  }

  function handleCancel() {
    showForm = false;
    editingCredential = null;
    message = '';
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();

    if (!name.trim() || !foundryUrl.trim() || !foundryUsername.trim()) {
      message = 'Name, URL, and username are required.';
      messageType = 'error';
      return;
    }

    if (!foundryPassword) {
      message = 'Password is required.';
      messageType = 'error';
      return;
    }

    saving = true;
    message = '';

    let result;
    if (isEdit) {
      const body: Record<string, unknown> = {
        name: name.trim(),
        foundryUrl: foundryUrl.trim(),
        foundryUsername: foundryUsername.trim(),
        foundryPassword,
      };
      result = await updateCredential(editingCredential!.id, body);
    } else {
      result = await createCredential({
        name: name.trim(),
        foundryUrl: foundryUrl.trim(),
        foundryUsername: foundryUsername.trim(),
        foundryPassword,
      });
    }

    saving = false;

    if (result.ok) {
      showForm = false;
      editingCredential = null;
      loadCredentials();
    } else {
      message = result.error;
      messageType = 'error';
    }
  }

  let modal = $state<{ open: boolean; action: () => Promise<void> }>({ open: false, action: async () => {} });
  let modalTitle = $state('');
  let modalMessage = $state('');

  async function runModal() { modal.open = false; await modal.action(); }

  function handleDelete(cred: Credential) {
    modalTitle = 'Delete credential';
    modalMessage = `Delete credential "${cred.name}"? This cannot be undone.`;
    modal = {
      open: true,
      action: async () => {
        const result = await deleteCredential(cred.id);
        if (result.ok) { loadCredentials(); } else { message = result.error; messageType = 'error'; }
      },
    };
  }
</script>

<ConfirmModal
  open={modal.open}
  title={modalTitle}
  message={modalMessage}
  confirmLabel="Delete"
  dangerous={true}
  onConfirm={runModal}
  onCancel={() => modal.open = false}
/>

<h2 class="page-title">Foundry Credentials</h2>

<div class="card">
  <div class="card-header">
    <h3 class="card-title">Stored Credentials</h3>
    {#if !showForm}
      <button class="btn btn-primary btn-sm" onclick={handleCreate}>+ Add Credential</button>
    {/if}
  </div>

  {#if showForm}
    <div class="key-form card" style="border: 1px solid var(--color-border-light);">
      <h3 class="card-title mb-2">{isEdit ? 'Edit Credential' : 'Add Credential'}</h3>

      <form onsubmit={handleSubmit}>
        <div class="form-group">
          <label class="form-label" for="cred-name">Name *</label>
          <input class="form-input" type="text" id="cred-name" bind:value={name} placeholder="e.g., My Foundry Server" required />
        </div>

        <div class="form-group">
          <label class="form-label" for="cred-url">Foundry URL *</label>
          <input class="form-input" type="url" id="cred-url" bind:value={foundryUrl} placeholder="https://your-foundry-instance.com" required />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label" for="cred-user">Username *</label>
            <input class="form-input" type="text" id="cred-user" bind:value={foundryUsername} placeholder="GM username" required />
          </div>
          <div class="form-group">
            <label class="form-label" for="cred-pass">Password *</label>
            <input class="form-input" type="password" id="cred-pass" bind:value={foundryPassword} placeholder="Encrypted at rest" required />
          </div>
        </div>

        <div class="flex gap-1 mt-2">
          <button class="btn btn-primary" type="submit" disabled={saving}>
            {saving ? 'Saving...' : (isEdit ? 'Update' : 'Save')}
          </button>
          <button class="btn btn-secondary" type="button" onclick={handleCancel}>Cancel</button>
        </div>

        {#if message}
          <div class="alert mt-1" class:alert-error={messageType === 'error'}>
            {message}
          </div>
        {/if}
      </form>
    </div>
  {/if}

  {#if loading}
    <p class="text-muted">Loading credentials...</p>
  {:else if credentials.length === 0}
    <p class="text-muted" style="font-size: 0.875rem;">No stored credentials. Add one to enable headless Foundry sessions.</p>
  {:else}
    <div class="table-wrapper">
      <table class="table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Foundry URL</th>
            <th>Username</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each credentials as cred (cred.id)}
            <tr>
              <td>{cred.name}</td>
              <td><code class="inline-code">{cred.foundryUrl}</code></td>
              <td>{cred.foundryUsername}</td>
              <td>{new Date(cred.createdAt).toLocaleDateString()}</td>
              <td class="actions-cell">
                <button class="btn btn-sm btn-ghost" onclick={() => handleEdit(cred)}>Edit</button>
                <button class="btn btn-sm btn-danger" onclick={() => handleDelete(cred)}>Delete</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  .table-wrapper {
    overflow-x: auto;
  }

  .actions-cell {
    white-space: nowrap;
  }
</style>
