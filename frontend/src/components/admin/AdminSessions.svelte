<script lang="ts">
  import { onMount } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import ConfirmModal from '../ui/ConfirmModal.svelte';

  let sessions = $state<any[]>([]);
  let total = $state(0);
  let error = $state('');

  async function load() {
    try {
      const r = await adminApi.listHeadlessSessions();
      sessions = r.sessions;
      total = r.total;
    } catch (e: any) { error = e?.message ?? 'Failed'; }
  }
  onMount(load);

  let modal = $state<{ open: boolean; action: () => Promise<void> }>({ open: false, action: async () => {} });
  let modalMsg = $state('');
  async function runModal() { modal.open = false; await modal.action(); }

  function kill(id: string) {
    modalMsg = 'Kill this headless session? The running Foundry process will be terminated.';
    modal = { open: true, action: async () => { await adminApi.killHeadlessSession(id); await load(); } };
  }
</script>

<ConfirmModal open={modal.open} title="Confirm" message={modalMsg} confirmLabel="Kill session" dangerous={true} onConfirm={runModal} onCancel={() => modal.open = false} />

<div class="admin-page">
  <h1>Headless Sessions ({total})</h1>
  {#if error}<p class="error">{error}</p>{/if}
  <div class="table-wrap">
    <table>
      <thead><tr>
        <th>Session ID</th><th>Client ID</th><th>Foundry URL</th><th>World</th><th>Started</th><th>Last Activity</th><th></th>
      </tr></thead>
      <tbody>
        {#each sessions as s}
          <tr>
            <td><code>{s.sessionId}</code></td>
            <td><code>{s.clientId}</code></td>
            <td>{s.foundryUrl}</td>
            <td>{s.worldName ?? '—'}</td>
            <td>{new Date(s.startedAt).toLocaleString()}</td>
            <td>{new Date(s.lastActivity).toLocaleString()}</td>
            <td><button class="danger" onclick={() => kill(s.sessionId)}>Kill</button></td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<style>
  .admin-page { padding: 1.5rem; }
  .table-wrap { overflow-x: auto; border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  table { width: 100%; border-collapse: collapse; font-size: 0.8125rem; }
  th, td { padding: 0.5rem; text-align: left; border-bottom: 1px solid var(--color-border-light); }
  code { font-size: 0.7rem; }
  button.danger { background: var(--color-error, #e53935); color: white; border: none; padding: 0.25rem 0.5rem; }
  .error { color: var(--color-error, #e53935); }
</style>
