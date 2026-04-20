<script lang="ts">
  import { adminLogin } from '../../lib/adminAuth';

  interface Props {
    onSuccess: () => void;
  }
  let { onSuccess }: Props = $props();

  let email = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    loading = true;
    const result = await adminLogin(email, password);
    loading = false;
    if (result.ok) {
      onSuccess();
    } else {
      // Use a generic message to avoid leaking server implementation details.
      error = 'Invalid credentials. Please check your email and password.';
    }
  }
</script>

<div class="admin-login">
  <h1><i class="fa-solid fa-shield-halved"></i> Admin Login</h1>
  <p class="hint">Admin sessions use a separate, short-lived authentication cookie. Sessions expire after 15 minutes of inactivity.</p>

  <form onsubmit={submit}>
    <label>
      Email
      <input type="email" bind:value={email} required autocomplete="username" />
    </label>
    <label>
      Password
      <input type="password" bind:value={password} required autocomplete="current-password" />
    </label>
    {#if error}
      <p class="error">{error}</p>
    {/if}
    <button type="submit" disabled={loading}>
      {loading ? 'Signing in…' : 'Sign in to admin'}
    </button>
  </form>
</div>

<style>
  .admin-login {
    max-width: 480px;
    margin: 2rem auto;
    padding: 2rem;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-lg);
  }
  .admin-login h1 {
    margin: 0 0 0.5rem;
    color: var(--color-primary);
  }
  .hint {
    color: var(--color-text-muted);
    margin-bottom: 1.5rem;
    font-size: 0.875rem;
  }
  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  label {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-weight: 500;
  }
  input {
    padding: 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg);
    color: var(--color-text);
  }
  button {
    padding: 0.75rem;
    background: var(--color-primary);
    color: white;
    border: none;
    border-radius: var(--radius-md);
    font-weight: 600;
    cursor: pointer;
  }
  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .error {
    color: var(--color-error, #e53935);
    margin: 0;
  }
</style>
