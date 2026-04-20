<script lang="ts">
  import { onMount } from 'svelte';
  import { theme, toggleTheme, initTheme } from '../../lib/theme';
  import { isAdmin, loadAdminSession, adminLogout } from '../../lib/adminAuth';
  import { billingEnabled, headlessEnabled, loadServerConfig } from '../../lib/config';
  import AdminLogin from './AdminLogin.svelte';
  import AdminUsers from './AdminUsers.svelte';
  import AdminClients from './AdminClients.svelte';
  import AdminKeys from './AdminKeys.svelte';
  import AdminAuditLogs from './AdminAuditLogs.svelte';
  import AdminMetrics from './AdminMetrics.svelte';
  import AdminAlerts from './AdminAlerts.svelte';
  import AdminHealth from './AdminHealth.svelte';
  import AdminSessions from './AdminSessions.svelte';
  import AdminSubscriptions from './AdminSubscriptions.svelte';
  import AdminOps from './AdminOps.svelte';
  import AdminTestRunner from './AdminTestRunner.svelte';

  type AdminView =
    | 'users' | 'clients' | 'keys' | 'logs' | 'metrics' | 'alerts'
    | 'health' | 'sessions' | 'subscriptions' | 'ops' | 'tests';

  let view = $state<AdminView>('users');
  let mobileMenuOpen = $state(false);

  onMount(() => {
    loadServerConfig();
    initTheme();
    loadAdminSession();
  });

  function setView(v: AdminView) {
    view = v;
    mobileMenuOpen = false;
  }

  async function handleAdminLogout() {
    await adminLogout();
  }
</script>

<button class="mobile-menu-toggle" onclick={() => mobileMenuOpen = !mobileMenuOpen} aria-label="Toggle menu">
  <i class="fa-solid fa-bars"></i>
</button>

{#if mobileMenuOpen}
  <button class="sidebar-overlay open" onclick={() => mobileMenuOpen = false} aria-label="Close menu"></button>
{/if}

<div class="app-layout">
  <nav class="sidebar" class:open={mobileMenuOpen}>
    <div class="sidebar-header">
      <h1>Admin Panel</h1>
      <p>Foundry REST API</p>
    </div>

    <div class="sidebar-nav">
      {#if $isAdmin}
        <ul>
          <li>
            <button class="nav-button" class:active={view === 'users'} onclick={() => setView('users')}>
              <i class="fa-solid fa-users-gear icon"></i> Users
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={view === 'clients'} onclick={() => setView('clients')}>
              <i class="fa-solid fa-network-wired icon"></i> Clients
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={view === 'keys'} onclick={() => setView('keys')}>
              <i class="fa-solid fa-key icon"></i> All Keys
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={view === 'logs'} onclick={() => setView('logs')}>
              <i class="fa-solid fa-clipboard-list icon"></i> Audit Logs
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={view === 'metrics'} onclick={() => setView('metrics')}>
              <i class="fa-solid fa-chart-line icon"></i> Metrics
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={view === 'alerts'} onclick={() => setView('alerts')}>
              <i class="fa-solid fa-bell icon"></i> Alerts
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={view === 'health'} onclick={() => setView('health')}>
              <i class="fa-solid fa-heart-pulse icon"></i> Health
            </button>
          </li>
          {#if $headlessEnabled}
          <li>
            <button class="nav-button" class:active={view === 'sessions'} onclick={() => setView('sessions')}>
              <i class="fa-solid fa-window-restore icon"></i> Sessions
            </button>
          </li>
          {/if}
          {#if $billingEnabled}
          <li>
            <button class="nav-button" class:active={view === 'subscriptions'} onclick={() => setView('subscriptions')}>
              <i class="fa-solid fa-credit-card icon"></i> Subscriptions
            </button>
          </li>
          {/if}
          <li>
            <button class="nav-button" class:active={view === 'ops'} onclick={() => setView('ops')}>
              <i class="fa-solid fa-toolbox icon"></i> Ops
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={view === 'tests'} onclick={() => setView('tests')}>
              <i class="fa-solid fa-vial icon"></i> Tests
            </button>
          </li>
          <li>
            <button class="nav-button" onclick={handleAdminLogout}>
              <i class="fa-solid fa-shield-halved icon"></i> Admin Logout
            </button>
          </li>
        </ul>
      {/if}

      <div class="sidebar-spacer"></div>

      <ul>
        <li><a class="nav-link" href="/">
          <i class="fa-solid fa-arrow-left icon"></i> Dashboard
        </a></li>
        <li><a class="nav-link" href="/docs" target="_blank" rel="noopener">
          <i class="fa-solid fa-file-lines icon"></i> Documentation
        </a></li>
        <li><a class="nav-link" href="https://github.com/ThreeHats/foundryvtt-rest-api-relay/issues" target="_blank" rel="noopener">
          <i class="fa-solid fa-circle-exclamation icon"></i> Report an Issue
        </a></li>
        <li><a class="nav-link" href="https://discord.gg/U634xNGRAC" target="_blank" rel="noopener">
          <i class="fa-brands fa-discord icon"></i> Discord
        </a></li>
      </ul>
    </div>

    <div class="sidebar-footer">
      <button class="theme-toggle" onclick={toggleTheme}>
        {#if $theme === 'dark'}
          <i class="fa-solid fa-sun icon"></i> Light Mode
        {:else}
          <i class="fa-solid fa-moon icon"></i> Dark Mode
        {/if}
      </button>
    </div>
  </nav>

  <div class="main-content">
    <div class="content-area">
      {#if $isAdmin}
        {#if view === 'users'}
          <AdminUsers />
        {:else if view === 'clients'}
          <AdminClients />
        {:else if view === 'keys'}
          <AdminKeys />
        {:else if view === 'logs'}
          <AdminAuditLogs />
        {:else if view === 'metrics'}
          <AdminMetrics />
        {:else if view === 'alerts'}
          <AdminAlerts />
        {:else if view === 'health'}
          <AdminHealth />
        {:else if view === 'sessions'}
          <AdminSessions />
        {:else if view === 'subscriptions'}
          <AdminSubscriptions />
        {:else if view === 'ops'}
          <AdminOps />
        {:else if view === 'tests'}
          <AdminTestRunner />
        {/if}
      {:else}
        <AdminLogin onSuccess={() => setView('users')} />
      {/if}
    </div>

    <footer class="site-footer">
      <p>This is an open-source project that is not associated with Foundry Virtual Tabletop, or Foundry Gaming, LLC.</p>
      <p>Foundry Virtual Tabletop is copyright of Foundry Gaming, LLC.</p>
    </footer>
  </div>
</div>
