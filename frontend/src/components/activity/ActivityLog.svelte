<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { fetchActivity } from '../../lib/api';
  import type { ActivityEvent } from '../../lib/types';

  interface Props {
    adminMode?: boolean;
    fetchFn?: (params: Record<string, unknown>) => Promise<{ ok: boolean; data: { events: ActivityEvent[]; total: number; offset: number; limit: number } | null; error: string | null }>;
  }

  let { adminMode = false, fetchFn }: Props = $props();

  // Filters (applied on "Apply" click, not live)
  let filterType = $state('');
  let filterWorld = $state('');
  let filterAction = $state('');
  let filterFailuresOnly = $state(false);
  let filterSince = $state('');
  let filterUntil = $state('');
  let filterUserId = $state('');

  // Pagination
  const PAGE_SIZE = 50;
  let offset = $state(0);
  let total = $state(0);

  // Data
  let events = $state<ActivityEvent[]>([]);
  let loading = $state(false);
  let error = $state('');

  let refreshTimer: ReturnType<typeof setInterval> | null = null;

  async function load(resetOffset = false) {
    if (resetOffset) offset = 0;
    loading = true;
    error = '';

    const params: Record<string, unknown> = {
      limit: PAGE_SIZE,
      offset,
    };
    if (filterType) params.type = filterType;
    if (filterWorld.trim()) params.world = filterWorld.trim();
    if (filterAction.trim()) params.action = filterAction.trim();
    if (filterFailuresOnly) params.success = false;
    if (filterSince) params.since = new Date(filterSince).toISOString();
    if (filterUntil) params.until = new Date(filterUntil + 'T23:59:59').toISOString();
    if (adminMode && filterUserId.trim()) params.userId = parseInt(filterUserId, 10) || undefined;

    try {
      const doFetch = fetchFn ?? (fetchActivity as typeof fetchFn);
      const result = await doFetch!(params);
      if (result.ok && result.data) {
        events = result.data.events ?? [];
        total = result.data.total ?? 0;
      } else {
        error = result.error ?? 'Failed to load activity';
      }
    } catch {
      error = 'Failed to load activity';
    } finally {
      loading = false;
    }
  }

  function resetTimer() {
    if (refreshTimer) clearInterval(refreshTimer);
    refreshTimer = setInterval(() => load(), 30_000);
  }

  onMount(() => {
    load();
    resetTimer();
  });

  onDestroy(() => {
    if (refreshTimer) clearInterval(refreshTimer);
  });

  function applyFilters() {
    load(true);
    resetTimer();
  }

  function prevPage() {
    if (offset > 0) {
      offset = Math.max(0, offset - PAGE_SIZE);
      load();
    }
  }

  function nextPage() {
    if (offset + PAGE_SIZE < total) {
      offset += PAGE_SIZE;
      load();
    }
  }

  function relativeTime(ts: string): string {
    const diff = Date.now() - new Date(ts).getTime();
    const s = Math.floor(diff / 1000);
    if (s < 60) return `${s}s ago`;
    const m = Math.floor(s / 60);
    if (m < 60) return `${m}m ago`;
    const h = Math.floor(m / 60);
    if (h < 24) return `${h}h ago`;
    return `${Math.floor(h / 24)}d ago`;
  }

  function eventIcon(ev: ActivityEvent): string {
    if (ev.type === 'connection') return '🔗';
    if (ev.type === 'remote_request') return '↔';
    return '⚡';
  }

  function eventLabel(ev: ActivityEvent): string {
    if (ev.type === 'connection') {
      const world = ev.worldTitle || ev.clientId;
      const token = ev.tokenName ? ` (token: "${ev.tokenName}")` : '';
      return `${world} connected${token}`;
    }
    if (ev.type === 'remote_request') {
      return `${ev.clientId} → ${ev.targetClientId}: ${ev.action}`;
    }
    // module_event
    const who = ev.actor ? ` by ${ev.actor}` : '';
    return `${ev.action}${who}`;
  }

  const typeOptions = [
    { value: '', label: 'All types' },
    { value: 'connection', label: 'Connections' },
    { value: 'remote_request', label: 'Cross-World' },
    { value: 'module_event', label: 'Module Events' },
  ];
</script>

<div class="activity-log">
  <div class="log-header">
    <h2>Activity Log</h2>
    <button class="refresh-btn" onclick={() => { load(); resetTimer(); }} disabled={loading}>
      {loading ? 'Loading…' : '↻ Refresh'}
    </button>
  </div>

  {#if error}<p class="error">{error}</p>{/if}

  <!-- Filter bar -->
  <div class="filters">
    <label>
      Type
      <select bind:value={filterType}>
        {#each typeOptions as opt}
          <option value={opt.value}>{opt.label}</option>
        {/each}
      </select>
    </label>
    <label>
      World / Client ID
      <input type="text" bind:value={filterWorld} placeholder="client-id" />
    </label>
    <label>
      Action
      <input type="text" bind:value={filterAction} placeholder="entity, roll…" />
    </label>
    {#if adminMode}
      <label>
        User ID
        <input type="number" bind:value={filterUserId} placeholder="0 = all" />
      </label>
    {/if}
    <label class="checkbox-label">
      <input type="checkbox" bind:checked={filterFailuresOnly} />
      Failures only
    </label>
    <label>
      From
      <input type="date" bind:value={filterSince} />
    </label>
    <label>
      To
      <input type="date" bind:value={filterUntil} />
    </label>
    <button class="apply-btn" onclick={applyFilters}>Apply</button>
  </div>

  <!-- Event table -->
  {#if events.length === 0 && !loading}
    <p class="muted">No events match the current filters.</p>
  {:else}
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Time</th>
            <th></th>
            <th>Event</th>
            {#if adminMode}<th>User ID</th>{/if}
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {#each events as ev (ev.type + ':' + ev.id)}
            <tr class:row-failure={ev.type === 'remote_request' && ev.success === false} class:row-flagged={ev.flagged}>
              <td class="col-time">{relativeTime(ev.createdAt)}</td>
              <td class="col-icon">{eventIcon(ev)}</td>
              <td class="col-label">
                <span class="event-main">{eventLabel(ev)}</span>
                {#if ev.flagged}<span class="badge flagged">⚠ flagged</span>{/if}
                {#if ev.errorMessage}
                  <div class="event-error">{ev.errorMessage}</div>
                {/if}
                {#if ev.type === 'module_event' && ev.description}
                  <div class="event-desc">{ev.description}</div>
                {/if}
              </td>
              {#if adminMode}<td class="col-user">{ev.userId ?? '—'}</td>{/if}
              <td class="col-status">
                {#if ev.type === 'remote_request'}
                  <span class="badge" class:badge-ok={ev.success} class:badge-fail={!ev.success}>
                    {ev.success ? '✓' : '✗'}
                  </span>
                {:else if ev.type === 'connection'}
                  <span class="badge badge-conn">conn</span>
                {:else}
                  <span class="badge badge-evt">evt</span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="pagination">
      <button onclick={prevPage} disabled={offset === 0}>‹ Prev</button>
      <span>{offset + 1}–{Math.min(offset + PAGE_SIZE, total)} of {total}</span>
      <button onclick={nextPage} disabled={offset + PAGE_SIZE >= total}>Next ›</button>
    </div>
  {/if}
</div>

<style>
  .activity-log { padding: 0; }
  .log-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem; }
  .log-header h2 { margin: 0; font-size: 1.1rem; }
  .refresh-btn { padding: 0.3rem 0.7rem; font-size: 0.8rem; }
  .error { color: var(--color-error, #e53935); font-size: 0.875rem; }
  .muted { color: var(--color-text-muted, #888); font-size: 0.875rem; }

  /* Filters */
  .filters { display: flex; flex-wrap: wrap; gap: 0.75rem; align-items: flex-end; margin-bottom: 1rem; padding: 0.75rem; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  .filters label { display: flex; flex-direction: column; gap: 0.2rem; font-size: 0.8rem; font-weight: 500; }
  .filters select, .filters input[type="text"], .filters input[type="number"], .filters input[type="date"] { padding: 0.3rem 0.5rem; border: 1px solid var(--color-border); border-radius: var(--radius-sm, 4px); background: var(--color-bg); font-size: 0.8rem; }
  .checkbox-label { flex-direction: row !important; align-items: center; gap: 0.35rem; font-weight: 400 !important; }
  .apply-btn { align-self: flex-end; padding: 0.3rem 0.8rem; font-size: 0.8rem; }

  /* Table */
  .table-wrap { overflow-x: auto; border: 1px solid var(--color-border); border-radius: var(--radius-md); }
  table { width: 100%; border-collapse: collapse; font-size: 0.85rem; }
  th, td { padding: 0.4rem 0.6rem; border-bottom: 1px solid var(--color-border-light); text-align: left; }
  th { background: var(--color-bg-elevated); font-weight: 600; font-size: 0.8rem; }
  .col-time { white-space: nowrap; color: var(--color-text-muted, #888); font-size: 0.75rem; width: 5rem; }
  .col-icon { width: 2rem; text-align: center; font-size: 1rem; }
  .col-label { max-width: 500px; }
  .col-user { white-space: nowrap; font-family: monospace; font-size: 0.8rem; }
  .col-status { width: 4rem; text-align: center; }

  .event-main { display: block; }
  .event-error { margin-top: 0.2rem; font-size: 0.75rem; color: var(--color-error, #e53935); font-family: monospace; }
  .event-desc { margin-top: 0.2rem; font-size: 0.75rem; color: var(--color-text-muted, #666); }

  .row-failure td { background: color-mix(in srgb, var(--color-error, #e53935) 6%, transparent); }
  .row-flagged td { background: color-mix(in srgb, #fb8c00 8%, transparent); }

  .badge { display: inline-block; padding: 0.1rem 0.35rem; border-radius: 3px; font-size: 0.7rem; font-weight: 700; text-transform: uppercase; }
  .badge-ok { background: #43a047; color: white; }
  .badge-fail { background: var(--color-error, #e53935); color: white; }
  .badge-conn { background: #1e88e5; color: white; }
  .badge-evt { background: #8e24aa; color: white; }
  .badge.flagged { background: #fb8c00; color: white; margin-left: 0.4rem; }

  /* Pagination */
  .pagination { display: flex; align-items: center; gap: 0.75rem; padding: 0.5rem 0.75rem; font-size: 0.8rem; color: var(--color-text-muted, #666); justify-content: flex-end; }
  .pagination button { padding: 0.2rem 0.6rem; font-size: 0.8rem; }
  .pagination button:disabled { opacity: 0.4; cursor: default; }
</style>
