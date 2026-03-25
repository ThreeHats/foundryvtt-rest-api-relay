<script lang="ts">
  import { register } from '../../lib/api';
  import { saveUser } from '../../lib/auth';

  interface Props {
    onSuccess: () => void;
  }

  let { onSuccess }: Props = $props();

  let email = $state('');
  let password = $state('');
  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');
  let loading = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    message = '';

    const result = await register(email, password);
    loading = false;

    if (result.ok) {
      message = 'Account created successfully!';
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
  <h2 class="page-title">Create Account</h2>
  <form onsubmit={handleSubmit}>
    <div class="form-group">
      <label class="form-label" for="signup-email">Email</label>
      <input
        class="form-input"
        type="email"
        id="signup-email"
        bind:value={email}
        required
        autocomplete="email"
      />
    </div>
    <div class="form-group">
      <label class="form-label" for="signup-password">Password</label>
      <input
        class="form-input"
        type="password"
        id="signup-password"
        bind:value={password}
        required
        minlength="8"
        autocomplete="new-password"
      />
      <p class="form-hint">Min 8 characters, with uppercase, lowercase, and a number</p>
    </div>
    <button class="btn btn-primary" type="submit" disabled={loading}>
      {loading ? 'Creating account...' : 'Sign Up'}
    </button>
  </form>
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
</style>
