<script lang="ts">
  import { onMount } from 'svelte';
  import { user, isLoggedIn, clearUser } from '../lib/auth';
  import { logout as apiLogout } from '../lib/api';
  import { validateResetToken } from '../lib/api';
  import { theme, toggleTheme, initTheme } from '../lib/theme';
  import { billingEnabled, headlessEnabled, loadServerConfig } from '../lib/config';
  import AuthForm from './auth/AuthForm.svelte';
  import ForgotPasswordForm from './auth/ForgotPasswordForm.svelte';
  import ResetPasswordForm from './auth/ResetPasswordForm.svelte';
  import Dashboard from './dashboard/Dashboard.svelte';
  import ApiKeysPage from './api-keys/ApiKeysPage.svelte';
  import ConnectionsPage from './connections/ConnectionsPage.svelte';
  import CredentialsPage from './credentials/CredentialsPage.svelte';
  import ActivityLog from './activity/ActivityLog.svelte';
  import { fetchActivity } from '../lib/api';
  import NotificationSettingsForm from './notifications/NotificationSettingsForm.svelte';
  import KeyApprovalPage from './approval/KeyApprovalPage.svelte';
  import PairApprovalPage from './approval/PairApprovalPage.svelte';

  type AuthView = 'auth' | 'forgot' | 'reset';
  type AppView =
    | 'dashboard' | 'api-keys' | 'connections' | 'credentials' | 'remote-requests' | 'notifications' | 'approve' | 'pair';

  let authView = $state<AuthView>('auth');
  let authInitialMode = $state<'signin' | 'signup'>('signin');
  let appView = $state<AppView>('dashboard');
  let resetToken = $state<string | null>(null);
  let approveCode = $state<string | null>(null);
  let pairCode = $state<string | null>(null);
  let mobileMenuOpen = $state(false);

  onMount(() => {
    loadServerConfig();
    initTheme();

    const params = new URLSearchParams(window.location.search);
    const token = params.get('reset-token');

    if (token) {
      validateResetToken(token).then((result) => {
        if (result.ok && result.data.valid) {
          resetToken = token;
          authView = 'reset';
        } else {
          authInitialMode = 'signin';
          authView = 'auth';
        }
        window.history.replaceState({}, document.title, window.location.pathname);
      });
    }

    // Check for /approve/:code and /pair/:code routes
    const path = window.location.pathname;
    const approveMatch = path.match(/^\/approve\/([A-Za-z0-9]+)$/);
    if (approveMatch) {
      approveCode = approveMatch[1];
      appView = 'approve';
    }
    const pairMatch = path.match(/^\/pair\/([A-Za-z0-9]+)$/);
    if (pairMatch) {
      pairCode = pairMatch[1];
      appView = 'pair';
    }
  });

  function handleAuthSuccess() {
    // Preserve pending approval/pair views so users who logged in via the
    // link land directly on the right page rather than the dashboard.
    if (pairCode) {
      appView = 'pair';
    } else if (approveCode) {
      appView = 'approve';
    } else {
      appView = 'dashboard';
    }
  }

  async function handleLogout() {
    // Best-effort: tell the server to invalidate this session token before
    // wiping local state. We don't block on the response — if it fails for
    // any reason, the cleanup loop will eventually expire the row.
    try {
      await apiLogout();
    } catch {
      // ignore
    }
    clearUser();
    authInitialMode = 'signin';
    authView = 'auth';
  }

  function handleAccountDeleted() {
    clearUser();
    authInitialMode = 'signin';
    authView = 'auth';
  }

  function handleResetComplete() {
    resetToken = null;
    authInitialMode = 'signin';
    authView = 'auth';
  }

  function setAuthView(view: AuthView) {
    authView = view;
    mobileMenuOpen = false;
  }

  function setAppView(view: AppView) {
    appView = view;
    mobileMenuOpen = false;
  }


</script>

{#if $isLoggedIn && ((appView === 'approve' && approveCode) || (appView === 'pair' && pairCode))}
  <!-- Standalone layout for approve/pair flows — no sidebar -->
  <div class="standalone-layout">
    <header class="standalone-header">
      <span class="standalone-brand">Foundry REST API</span>
      <button class="theme-toggle" onclick={toggleTheme}>
        {#if $theme === 'dark'}
          <i class="fa-solid fa-sun icon"></i> Light Mode
        {:else}
          <i class="fa-solid fa-moon icon"></i> Dark Mode
        {/if}
      </button>
    </header>
    <main class="standalone-content">
      {#if appView === 'approve' && approveCode}
        <KeyApprovalPage code={approveCode} />
      {:else if appView === 'pair' && pairCode}
        <PairApprovalPage code={pairCode} />
      {/if}
    </main>
  </div>
{:else}

<button class="mobile-menu-toggle" onclick={() => mobileMenuOpen = !mobileMenuOpen} aria-label="Toggle menu">
  <i class="fa-solid fa-bars"></i>
</button>

{#if mobileMenuOpen}
  <button class="sidebar-overlay open" onclick={() => mobileMenuOpen = false} aria-label="Close menu"></button>
{/if}

<div class="app-layout">
  <nav class="sidebar" class:open={mobileMenuOpen}>
    <div class="sidebar-header">
      <h1>Foundry REST API</h1>
      <p>Connect Foundry VTT to external apps</p>
    </div>

    <div class="sidebar-nav">
      {#if $isLoggedIn}
        <ul>
          <li>
            <button class="nav-button" class:active={appView === 'dashboard'} onclick={() => setAppView('dashboard')}>
              <i class="fa-solid fa-gauge icon"></i> Dashboard
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={appView === 'api-keys'} onclick={() => setAppView('api-keys')}>
              <i class="fa-solid fa-key icon"></i> API Keys
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={appView === 'connections'} onclick={() => setAppView('connections')}>
              <i class="fa-solid fa-link icon"></i> Connections
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={appView === 'remote-requests'} onclick={() => setAppView('remote-requests')}>
              <i class="fa-solid fa-clock-rotate-left icon"></i> Activity
            </button>
          </li>
          {#if $headlessEnabled}
          <li>
            <button class="nav-button" class:active={appView === 'credentials'} onclick={() => setAppView('credentials')}>
              <i class="fa-solid fa-lock icon"></i> Credentials
            </button>
          </li>
          {/if}
          <li>
            <button class="nav-button" class:active={appView === 'notifications'} onclick={() => setAppView('notifications')}>
              <i class="fa-solid fa-bell icon"></i> Notifications
            </button>
          </li>
          <li>
            <button class="nav-button" onclick={handleLogout}>
              <i class="fa-solid fa-right-from-bracket icon"></i> Sign Out
            </button>
          </li>

          {#if $user?.role === 'admin'}
            <li>
              <a class="nav-link" href="/admin">
                <i class="fa-solid fa-shield-halved icon"></i> Admin Panel
              </a>
            </li>
          {/if}
        </ul>
      {/if}

      <div class="sidebar-spacer"></div>

      <ul>
        <li><a class="nav-link" href="https://github.com/JustAnotherIdea/foundryvtt-rest-api" target="_blank" rel="noopener">
          <i class="fa-solid fa-book-open icon"></i> Read Me
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
      {#if $isLoggedIn}
        {#if appView === 'dashboard'}
          <Dashboard onAccountDeleted={handleAccountDeleted} onNavigate={(view) => setAppView(view as AppView)} />
        {:else if appView === 'api-keys'}
          <ApiKeysPage />
        {:else if appView === 'connections'}
          <ConnectionsPage />
        {:else if appView === 'remote-requests'}
          <div style="padding: 1.5rem;">
            <ActivityLog fetchFn={fetchActivity} />
          </div>
        {:else if appView === 'credentials'}
          <CredentialsPage />
        {:else if appView === 'notifications'}
          <NotificationSettingsForm />
        {/if}
      {:else}
        {#if authView === 'auth'}
          <AuthForm
            initialMode={authInitialMode}
            onSuccess={handleAuthSuccess}
            onForgotPassword={() => setAuthView('forgot')}
          />
        {:else if authView === 'forgot'}
          <ForgotPasswordForm onBackToLogin={() => { authInitialMode = 'signin'; setAuthView('auth'); }} />
        {:else if authView === 'reset' && resetToken}
          <ResetPasswordForm token={resetToken} onComplete={handleResetComplete} />
        {/if}
      {/if}
    </div>

    <footer class="site-footer">
      <p>This is an open-source project that is not associated with Foundry Virtual Tabletop, or Foundry Gaming, LLC.</p>
      <p>Foundry Virtual Tabletop is copyright of Foundry Gaming, LLC.</p>
      <p><a href="/privacy">Privacy Policy</a></p>
    </footer>
  </div>
</div>

{/if}
