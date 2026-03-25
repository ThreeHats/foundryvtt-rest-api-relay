<script lang="ts">
  import { changePassword } from '../../lib/api';

  let open = $state(false);
  let currentPassword = $state('');
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
      message = 'New passwords do not match.';
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

    const result = await changePassword(currentPassword, newPassword);
    loading = false;

    if (result.ok) {
      message = 'Password changed successfully!';
      messageType = 'success';
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
    } else {
      message = result.error;
      messageType = 'error';
    }
  }
</script>

<div class="collapsible">
  <button class="collapsible-trigger" class:open onclick={() => open = !open}>
    <span class="arrow">&#9654;</span>
    Change Password
  </button>

  {#if open}
    <form class="mt-2" onsubmit={handleSubmit}>
      <div class="form-group">
        <label class="form-label" for="cp-current">Current Password</label>
        <input class="form-input" type="password" id="cp-current" bind:value={currentPassword} required autocomplete="current-password" />
      </div>
      <div class="form-group">
        <label class="form-label" for="cp-new">New Password</label>
        <input class="form-input" type="password" id="cp-new" bind:value={newPassword} required minlength="8" autocomplete="new-password" />
        <p class="form-hint">Min 8 characters, with uppercase, lowercase, and a number</p>
      </div>
      <div class="form-group">
        <label class="form-label" for="cp-confirm">Confirm New Password</label>
        <input class="form-input" type="password" id="cp-confirm" bind:value={confirmPassword} required minlength="8" autocomplete="new-password" />
      </div>
      <button class="btn btn-primary" type="submit" disabled={loading}>
        {loading ? 'Changing...' : 'Change Password'}
      </button>
      {#if message}
        <div class="alert mt-1" class:alert-success={messageType === 'success'} class:alert-error={messageType === 'error'}>
          {message}
        </div>
      {/if}
    </form>
  {/if}
</div>

<style>
  .collapsible {
    max-width: 420px;
  }
</style>
