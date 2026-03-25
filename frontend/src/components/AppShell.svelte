<script lang="ts">
  import { onMount } from 'svelte';
  import { user, isLoggedIn, clearUser } from '../lib/auth';
  import { validateResetToken } from '../lib/api';
  import { theme, toggleTheme, initTheme } from '../lib/theme';
  import SignUpForm from './auth/SignUpForm.svelte';
  import SignInForm from './auth/SignInForm.svelte';
  import ForgotPasswordForm from './auth/ForgotPasswordForm.svelte';
  import ResetPasswordForm from './auth/ResetPasswordForm.svelte';
  import Dashboard from './dashboard/Dashboard.svelte';
  import ApiKeysPage from './api-keys/ApiKeysPage.svelte';

  type AuthView = 'signup' | 'signin' | 'forgot' | 'reset';
  type AppView = 'dashboard' | 'api-keys';

  let authView = $state<AuthView>('signup');
  let appView = $state<AppView>('dashboard');
  let resetToken = $state<string | null>(null);
  let mobileMenuOpen = $state(false);

  onMount(() => {
    initTheme();

    const params = new URLSearchParams(window.location.search);
    const token = params.get('reset-token');

    if (token) {
      validateResetToken(token).then((result) => {
        if (result.ok && result.data.valid) {
          resetToken = token;
          authView = 'reset';
        } else {
          authView = 'signin';
        }
        window.history.replaceState({}, document.title, window.location.pathname);
      });
    }
  });

  function handleAuthSuccess() {
    appView = 'dashboard';
  }

  function handleLogout() {
    clearUser();
    authView = 'signup';
  }

  function handleAccountDeleted() {
    clearUser();
    authView = 'signin';
  }

  function handleResetComplete() {
    resetToken = null;
    authView = 'signin';
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
            <button class="nav-button" onclick={handleLogout}>
              <i class="fa-solid fa-right-from-bracket icon"></i> Sign Out
            </button>
          </li>
        </ul>
      {:else}
        <ul>
          <li>
            <button class="nav-button" class:active={authView === 'signup'} onclick={() => setAuthView('signup')}>
              <i class="fa-solid fa-user-plus icon"></i> Sign Up
            </button>
          </li>
          <li>
            <button class="nav-button" class:active={authView === 'signin'} onclick={() => setAuthView('signin')}>
              <i class="fa-solid fa-right-to-bracket icon"></i> Sign In
            </button>
          </li>
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
        <li><a class="nav-link" href="https://github.com/JustAnotherIdea/foundryvtt-rest-api/issues/new" target="_blank" rel="noopener">
          <i class="fa-solid fa-bug icon"></i> Report a Bug
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
          <Dashboard onAccountDeleted={handleAccountDeleted} />
        {:else if appView === 'api-keys'}
          <ApiKeysPage />
        {/if}
      {:else}
        {#if authView === 'signup'}
          <SignUpForm onSuccess={handleAuthSuccess} />
        {:else if authView === 'signin'}
          <SignInForm onSuccess={handleAuthSuccess} onForgotPassword={() => setAuthView('forgot')} />
        {:else if authView === 'forgot'}
          <ForgotPasswordForm onBackToLogin={() => setAuthView('signin')} />
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
