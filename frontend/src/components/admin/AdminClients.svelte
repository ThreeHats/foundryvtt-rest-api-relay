<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminClientInfo } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';

  let clients = $state<AdminClientInfo[]>([]);
  let total = $state(0);
  let error = $state('');
  let intervalId: ReturnType<typeof setInterval> | undefined;

  async function load() {
    try {
      const r = await adminApi.listClients();
      clients = r.clients;
      total = r.total;
      error = '';
    } catch (e: any) {
      error = e?.message ?? 'Failed to load';
    }
  }

  onMount(() => {
    load();
    intervalId = setInterval(load, 5000);
  });
  onDestroy(() => { if (intervalId) clearInterval(intervalId); });

  let modal = $state<{ open: boolean; action: () => Promise<void> }>({ open: false, action: async () => {} });
  let modalMsg = $state('');
  async function runModal() { modal.open = false; await modal.action(); }

  function disconnect(clientId: string) {
    modalMsg = `Force disconnect client ${clientId}?`;
    modal = { open: true, action: async () => { await adminApi.disconnectClient(clientId); await load(); } };
  }

</script>

<ConfirmModal open={modal.open} title="Confirm" message={modalMsg} confirmLabel="Disconnect" dangerous={true} onConfirm={runModal} onCancel={() => modal.open = false} />

<div class="admin-page">
  <h1>Connected Clients ({total})</h1>
  {#if error}<p class="error">{error}</p>{/if}

  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>Client ID</th>
          <th>World</th>
          <th>System</th>
          <th>Foundry</th>
          <th>IP</th>
          <th>Connected</th>
          <th>Last Seen</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each clients as c (c.clientId)}
          <tr>
            <td><code>{c.clientId}</code></td>
            <td>{c.worldTitle ?? c.worldId ?? '—'}</td>
            <td>{c.systemTitle ?? c.systemId ?? '—'}</td>
            <td>{c.foundryVersion ?? '—'}</td>
            <td>{c.ipAddress ?? '—'}</td>
            <td>{new Date(c.connectedSince).toLocaleString()}</td>
            <td>{new Date(c.lastSeen).toLocaleString()}</td>
            <td><button class="danger" onclick={() => disconnect(c.clientId)}>Disconnect</button></td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<style>
  .admin-page { padding: 1.5rem; }
  h1 { margin-top: 0; }
  .table-wrap { overflow-x: auto; border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  table { width: 100%; border-collapse: collapse; font-size: 0.875rem; }
  th, td { text-align: left; padding: 0.5rem 0.75rem; border-bottom: 1px solid var(--color-border-light); }
  th { background: var(--color-bg-elevated); font-weight: 600; }
  code { font-size: 0.75rem; }
  button.danger { background: var(--color-error, #e53935); color: white; border: none; padding: 0.25rem 0.5rem; }
  .error { color: var(--color-error, #e53935); }
</style>
