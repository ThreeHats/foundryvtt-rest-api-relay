<script lang="ts">
  import { forgotPassword } from '../../lib/api';

  interface Props {
    onBackToLogin: () => void;
  }

  let { onBackToLogin }: Props = $props();

  let email = $state('');
  let message = $state('');
  let messageType = $state<'success' | 'error'>('success');
  let loading = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    message = '';

    const result = await forgotPassword(email);
    loading = false;

    message = result.ok
      ? result.data.message || 'If an account with that email exists, a password reset link has been sent.'
      : 'If an account with that email exists, a password reset link has been sent.';
    messageType = 'success';
  }
</script>

<div class="auth-form">
  <h2 class="page-title">Forgot Password</h2>
  <p class="text-muted mb-2">Enter your email address and we'll send you a link to reset your password.</p>
  <form onsubmit={handleSubmit}>
    <div class="form-group">
      <label class="form-label" for="forgot-email">Email</label>
      <input
        class="form-input"
        type="email"
        id="forgot-email"
        bind:value={email}
        required
        autocomplete="email"
      />
    </div>
    <button class="btn btn-primary" type="submit" disabled={loading}>
      {loading ? 'Sending...' : 'Send Reset Link'}
    </button>
  </form>
  <p class="mt-2">
    <button class="link-button" onclick={onBackToLogin}>Back to Sign In</button>
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
