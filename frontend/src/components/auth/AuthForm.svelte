<script lang="ts">
  import { login, register } from '../../lib/api';
  import { saveUser } from '../../lib/auth';

  interface Props {
    initialMode?: 'signin' | 'signup';
    onSuccess: () => void;
    onForgotPassword: () => void;
  }

  let { initialMode = 'signin', onSuccess, onForgotPassword }: Props = $props();

  let mode = $state<'signin' | 'signup'>(initialMode);
  let email = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');
  let loading = $state(false);

  function switchMode(newMode: 'signin' | 'signup') {
    mode = newMode;
    message = '';
    password = '';
    confirmPassword = '';
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    message = '';

    if (mode === 'signin') {
      const result = await login(email, password);
      loading = false;
      if (result.ok) {
        saveUser(result.data);
        onSuccess();
      } else {
        message = result.error;
        messageType = 'error';
      }
    } else {
      if (password !== confirmPassword) {
        message = 'Passwords do not match.';
        messageType = 'error';
        loading = false;
        return;
      }
      const result = await register(email, password);
      loading = false;
      if (result.ok) {
        saveUser(result.data);
        onSuccess();
      } else {
        message = result.error;
        messageType = 'error';
      }
    }
  }
</script>

<div class="auth-form">
  <h2 class="page-title">{mode === 'signin' ? 'Sign In' : 'Create Account'}</h2>

  <form onsubmit={handleSubmit}>
    <div class="form-group">
      <label class="form-label" for="auth-email">Email</label>
      <input
        class="form-input"
        type="email"
        id="auth-email"
        bind:value={email}
        required
        autocomplete="email"
      />
    </div>
    <div class="form-group">
      <label class="form-label" for="auth-password">Password</label>
      <input
        class="form-input"
        type="password"
        id="auth-password"
        bind:value={password}
        required
        minlength={mode === 'signup' ? 8 : undefined}
        autocomplete={mode === 'signup' ? 'new-password' : 'current-password'}
      />
      {#if mode === 'signup'}
        <p class="form-hint">Min 8 characters, with uppercase, lowercase, and a number</p>
      {/if}
    </div>

    {#if mode === 'signup'}
      <div class="form-group">
        <label class="form-label" for="auth-confirm-password">Confirm Password</label>
        <input
          class="form-input"
          type="password"
          id="auth-confirm-password"
          bind:value={confirmPassword}
          required
          autocomplete="new-password"
        />
      </div>
    {/if}

    <button class="btn btn-primary" type="submit" disabled={loading}>
      {#if loading}
        {mode === 'signin' ? 'Signing in...' : 'Creating account...'}
      {:else}
        {mode === 'signin' ? 'Sign In' : 'Sign Up'}
      {/if}
    </button>
  </form>

  {#if mode === 'signin'}
    <p class="mt-2">
      <button class="link-button" onclick={onForgotPassword} type="button">Forgot your password?</button>
    </p>
  {/if}

  <p class="mode-toggle">
    {#if mode === 'signin'}
      Don't have an account? <button class="link-button" onclick={() => switchMode('signup')} type="button">Sign up</button>
    {:else}
      Already have an account? <button class="link-button" onclick={() => switchMode('signin')} type="button">Sign in</button>
    {/if}
  </p>

  {#if message}
    <div class="alert mt-2" class:alert-success={messageType === 'success'} class:alert-error={messageType === 'error'}>
      {message}
    </div>
  {/if}
</div>

<style>
  .auth-form {
    max-width: 420px;
  }

  .mode-toggle {
    margin-top: 1.25rem;
    font-size: 0.875rem;
    color: var(--color-text-muted);
  }

  .link-button {
    background: none;
    border: none;
    color: var(--color-primary);
    cursor: pointer;
    font-size: 0.875rem;
    padding: 0;
    text-decoration: underline;
    text-underline-offset: 2px;
  }

  .link-button:hover {
    color: var(--color-primary-hover);
  }

</style>
