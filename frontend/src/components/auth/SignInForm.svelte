<script lang="ts">
  import { login } from '../../lib/api';
  import { saveUser } from '../../lib/auth';

  interface Props {
    onSuccess: () => void;
    onForgotPassword: () => void;
  }

  let { onSuccess, onForgotPassword }: Props = $props();

  let email = $state('');
  let password = $state('');
  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');
  let loading = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    message = '';

    const result = await login(email, password);
    loading = false;

    if (result.ok) {
      message = 'Login successful!';
      messageType = 'success';
      saveUser(result.data);
      onSuccess();
    } else {
      message = result.error;
      messageType = 'error';
    }
  }
</script>

<div class="auth-form">
  <h2 class="page-title">Sign In</h2>
  <form onsubmit={handleSubmit}>
    <div class="form-group">
      <label class="form-label" for="login-email">Email</label>
      <input
        class="form-input"
        type="email"
        id="login-email"
        bind:value={email}
        required
        autocomplete="email"
      />
    </div>
    <div class="form-group">
      <label class="form-label" for="login-password">Password</label>
      <input
        class="form-input"
        type="password"
        id="login-password"
        bind:value={password}
        required
        autocomplete="current-password"
      />
    </div>
    <button class="btn btn-primary" type="submit" disabled={loading}>
      {loading ? 'Signing in...' : 'Sign In'}
    </button>
  </form>
  <p class="mt-2">
    <button class="link-button" onclick={onForgotPassword}>Forgot your password?</button>
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
