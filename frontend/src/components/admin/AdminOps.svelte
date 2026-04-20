<script lang="ts">
  import { onMount } from 'svelte';
  import { adminApi } from '../../lib/adminApi';
  import type { AdminFeatureFlags } from '../../lib/types';
  import ConfirmModal from '../ui/ConfirmModal.svelte';

  let flags = $state<AdminFeatureFlags | null>(null);
  let error = $state('');

  async function load() {
    try { flags = await adminApi.getFeatureFlags(); } catch (e: any) { error = e?.message ?? 'Failed'; }
  }
  onMount(load);

  let modal = $state<{ open: boolean; action: () => Promise<void> }>({ open: false, action: async () => {} });
  let modalMsg = $state('');
  async function runModal() { modal.open = false; await modal.action(); }

  async function toggle(name: keyof AdminFeatureFlags) {
    if (!flags) return;
    if (name === 'maintenance_mode' && !flags[name]) {
      modalMsg = 'Enable maintenance mode? All non-admin API requests will return 503 until it is disabled.';
      modal = { open: true, action: async () => { flags = await adminApi.setFeatureFlag(name, true); } };
      return;
    }
    flags = await adminApi.setFeatureFlag(name, !flags[name]);
  }
</script>

<ConfirmModal open={modal.open} title="Confirm" message={modalMsg} confirmLabel="Enable" dangerous={true} onConfirm={runModal} onCancel={() => modal.open = false} />

<div class="admin-page">
  <h1>Operational Tools</h1>
  {#if error}<p class="error">{error}</p>{/if}

  <section>
    <h2>Feature Flags</h2>
    {#if flags}
      <div class="flags">
        {#each Object.entries(flags) as [name, value]}
          <label class="flag">
            <input type="checkbox" checked={value} onchange={() => toggle(name as keyof AdminFeatureFlags)} />
            <span class="name">{name}</span>
            <span class="value">{value ? 'ON' : 'OFF'}</span>
          </label>
        {/each}
      </div>
    {/if}
  </section>
</div>

<style>
  .admin-page { padding: 1.5rem; }
  .flags { display: flex; flex-direction: column; gap: 0.75rem; max-width: 500px; }
  .flag { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem; background: var(--color-bg-elevated); border-radius: var(--radius-md); }
  .name { flex: 1; font-family: monospace; }
  .value { font-weight: 700; }
  .error { color: var(--color-error, #e53935); }
</style>
