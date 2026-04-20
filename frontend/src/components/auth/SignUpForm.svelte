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

  // One-time API key display modal state
  let showOneTimeModal = $state(false);
  let oneTimeKey = $state('');
  let oneTimeCopied = $state(false);
  let acknowledged = $state(false);

  async function copyOneTime() {
    if (!oneTimeKey) return;
    await navigator.clipboard.writeText(oneTimeKey);
    oneTimeCopied = true;
    setTimeout(() => { oneTimeCopied = false; }, 2000);
  }

  function dismissOneTimeModal() {
    if (!acknowledged) {
      alert('Please acknowledge that you have saved the key.');
      return;
    }
    oneTimeKey = '';
    showOneTimeModal = false;
    acknowledged = false;
    oneTimeCopied = false;
    onSuccess();
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();

    loading = true;
    message = '';

    const result = await register(email, password);
    loading = false;

    if (result.ok) {
      const fullKey = result.data.apiKey ?? '';
      const { apiKey: _omit, ...userData } = result.data;
      saveUser(userData);

      if (result.data.emailVerified === false) {
        message = 'Account created! Check your email for a verification link before using the API.';
      } else {
        message = 'Account created successfully!';
      }
      messageType = 'success';

      // Show one-time key modal — onSuccess() is called when user dismisses
      oneTimeKey = fullKey;
      showOneTimeModal = true;
      acknowledged = false;
      oneTimeCopied = false;
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

{#if showOneTimeModal}
  <div class="modal-backdrop" role="dialog" aria-modal="true">
    <div class="modal">
      <h2>⚠ Save Your Master API Key</h2>
      <p class="warning">
        <strong>This is the ONLY time this key will be displayed.</strong>
        Copy it and save it in a secure location (password manager, secrets vault) before continuing.
      </p>
      <div class="key-display">
        <code class="full-key">{oneTimeKey}</code>
        <button class="btn btn-primary" onclick={copyOneTime}>
          {oneTimeCopied ? '✓ Copied!' : 'Copy'}
        </button>
      </div>
      <p class="info">
        After dismissing this dialog, the key cannot be retrieved. If you lose it, you must regenerate the key — which invalidates this one. The dashboard will only show a truncated form.
      </p>
      <label class="ack-label">
        <input type="checkbox" bind:checked={acknowledged} />
        I have saved this key in a secure location
      </label>
      <div class="modal-actions">
        <button class="btn btn-primary" onclick={dismissOneTimeModal} disabled={!acknowledged}>
          Continue to Dashboard
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .auth-form {
    max-width: 420px;
  }

  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 1rem;
  }

  .modal {
    background: var(--color-bg, #fff);
    border-radius: var(--radius-md);
    padding: 1.5rem;
    max-width: 600px;
    width: 100%;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  }

  .modal h2 {
    margin: 0 0 1rem;
    color: var(--color-warning, #f39c12);
  }

  .warning {
    background: var(--color-warning-bg, #fff8e1);
    border-left: 4px solid var(--color-warning, #f39c12);
    padding: 0.75rem 1rem;
    margin: 1rem 0;
    border-radius: var(--radius-sm);
  }

  .key-display {
    display: flex;
    gap: 0.5rem;
    align-items: stretch;
    margin: 1rem 0;
  }

  .full-key {
    background: var(--color-bg-sunken);
    padding: 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.85rem;
    font-family: monospace;
    word-break: break-all;
    flex: 1;
    border: 2px solid var(--color-warning, #f39c12);
    user-select: all;
  }

  .info {
    font-size: 0.85rem;
    color: var(--color-text-muted, #666);
    margin: 0.5rem 0;
  }

  .ack-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 1rem 0;
    cursor: pointer;
    font-size: 0.9rem;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 1rem;
  }
</style>
