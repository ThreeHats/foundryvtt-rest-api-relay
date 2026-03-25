<script lang="ts">
  import { resetPassword } from '../../lib/api';

  interface Props {
    token: string;
    onComplete: () => void;
  }

  let { token, onComplete }: Props = $props();

  let newPassword = $state('');
  let confirmPassword = $state('');
  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');
  let loading = $state(false);

  function validatePassword(pw: string): string | null {
    if (pw.length < 8) return 'Password must be at least 8 characters long';
    if (!/[A-Z]/.test(pw)) return 'Password must contain at least one uppercase letter';
    if (!/[a-z]/.test(pw)) return 'Password must contain at least one lowercase letter';
    if (!/[0-9]/.test(pw)) return 'Password must contain at least one number';
    return null;
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();

    if (newPassword !== confirmPassword) {
      message = 'Passwords do not match.';
      messageType = 'error';
      return;
    }

    const validationError = validatePassword(newPassword);
    if (validationError) {
      message = validationError;
      messageType = 'error';
      return;
    }

    loading = true;
    message = '';

    const result = await resetPassword(token, newPassword);
    loading = false;

    if (result.ok) {
      message = 'Password reset successfully! Redirecting to login...';
      messageType = 'success';
      setTimeout(onComplete, 2000);
    } else {
      message = result.error;
      messageType = 'error';
    }
  }
</script>

<div class="auth-form">
  <h2 class="page-title">Reset Password</h2>
  <form onsubmit={handleSubmit}>
    <div class="form-group">
      <label class="form-label" for="reset-new-pw">New Password</label>
      <input
        class="form-input"
        type="password"
        id="reset-new-pw"
        bind:value={newPassword}
        required
        minlength="8"
        autocomplete="new-password"
      />
      <p class="form-hint">Min 8 characters, with uppercase, lowercase, and a number</p>
    </div>
    <div class="form-group">
      <label class="form-label" for="reset-confirm-pw">Confirm Password</label>
      <input
        class="form-input"
        type="password"
        id="reset-confirm-pw"
        bind:value={confirmPassword}
        required
        minlength="8"
        autocomplete="new-password"
      />
    </div>
    <button class="btn btn-primary" type="submit" disabled={loading}>
      {loading ? 'Resetting...' : 'Reset Password'}
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
