<script lang="ts">
  import { user, clearUser } from '../../lib/auth';
  import { regenerateKey } from '../../lib/api';

  let regenLoading = $state(false);
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
    // Wipe the one-time key from memory and force re-login. The backend
    // invalidates ALL dashboard sessions for this user when the master key
    // is rotated, so the current sessionToken is dead — we have to log out.
    oneTimeKey = '';
    showOneTimeModal = false;
    acknowledged = false;
    oneTimeCopied = false;
    clearUser();
    // Hard reload to land on the login screen with a clean state
    window.location.href = '/';
  }

  async function handleRegenerate() {
    if (!confirm('Are you sure you want to regenerate your master API key?\n\nThis will:\n• Invalidate your current master key immediately\n• Delete all scoped API keys\n• Delete all connection tokens (Foundry modules will need to re-pair)\n• Log you out of all dashboard sessions across every device\n• Require all integrations to be reconfigured\n\nThe new key will be shown EXACTLY ONCE. Make sure you can save it before continuing.')) return;

    const email = $user?.email;
    if (!email) return;

    const password = prompt('Enter your password to confirm:');
    if (!password) return;

    regenLoading = true;
    const result = await regenerateKey(email, password);
    regenLoading = false;

    if (result.ok) {
      // Master key is now invalid, ALL sessions invalidated server-side.
      // Show the one-time modal; on dismiss we force re-login.
      oneTimeKey = result.data.apiKey;
      showOneTimeModal = true;
      acknowledged = false;
      oneTimeCopied = false;
    } else {
      alert(result.error);
    }
  }
</script>

{#if $user?.apiKeyRotationRequired}
  <div class="rotation-warning">
    <strong>⚠ Master API key rotation required</strong>
    <p>For security reasons, your master API key must be rotated. API access has been blocked until you regenerate the key. Click the "Regenerate" button below to create a new key.</p>
  </div>
{/if}

<div class="key-row">
  <span class="form-label">Master API Key</span>
  <div class="key-actions">
    <p class="key-status">Never shown after creation. Regenerating issues a new key once — and resets all scoped keys, connections, and active sessions.</p>
    <button class="btn btn-sm btn-primary" onclick={handleRegenerate} disabled={regenLoading}>
      {regenLoading ? 'Regenerating...' : 'Regenerate'}
    </button>
  </div>
</div>

{#if $user?.apiKeyRotationRequired}
{:else}
{/if}

{#if showOneTimeModal}
  <div class="modal-backdrop" role="dialog" aria-modal="true">
    <div class="modal">
      <h2>⚠ Save Your New Master API Key</h2>
      <p class="warning">
        <strong>This is the ONLY time this key will be displayed.</strong>
        Copy it and save it in a secure location (password manager, secrets vault) before dismissing this dialog.
      </p>
      <div class="key-display">
        <code class="full-key">{oneTimeKey}</code>
        <button class="btn btn-primary" onclick={copyOneTime}>
          {oneTimeCopied ? '✓ Copied!' : 'Copy'}
        </button>
      </div>
      <p class="info">
        After dismissing this dialog, the key cannot be retrieved. If you lose it, you must regenerate again — which invalidates this key.
      </p>
      <label class="ack-label">
        <input type="checkbox" bind:checked={acknowledged} />
        I have saved this key in a secure location
      </label>
      <div class="modal-actions">
        <button class="btn btn-primary" onclick={dismissOneTimeModal} disabled={!acknowledged}>
          Done
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .key-row {
    margin-bottom: 1rem;
  }

  .rotation-warning {
    background: var(--color-warning-bg, #fff8e1);
    border: 1px solid var(--color-warning, #f39c12);
    border-radius: var(--radius-md);
    padding: 0.75rem 1rem;
    margin-bottom: 1rem;
    color: var(--color-warning-text, #856404);
  }

  .rotation-warning p {
    margin: 0.25rem 0 0;
    font-size: 0.875rem;
  }

  .key-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
    margin-top: 0.25rem;
  }

  .key-status {
    flex: 1;
    margin: 0;
    font-size: 0.875rem;
    color: var(--color-text-muted, #666);
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
