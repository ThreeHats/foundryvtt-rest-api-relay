<script lang="ts">
  import { user, clearUser } from '../../lib/auth';
  import { regenerateKey } from '../../lib/api';

  let regenLoading = $state(false);

  async function handleReset() {
    if (!confirm('Are you sure you want to reset all credentials?\n\nThis will:\n• Disconnect all paired Foundry modules (they will need to re-pair)\n• Delete all scoped API keys\n• Delete all connection tokens\n• Log you out of all dashboard sessions\n\nThis cannot be undone.')) return;

    const email = $user?.email;
    if (!email) return;

    const password = prompt('Enter your password to confirm:');
    if (!password) return;

    regenLoading = true;
    const result = await regenerateKey(email, password);
    regenLoading = false;

    if (result.ok) {
      alert('All credentials have been reset. You will now be logged out.');
      clearUser();
      window.location.href = '/';
    } else {
      alert(result.error);
    }
  }
</script>

{#if $user?.apiKeyRotationRequired}
  <div class="rotation-warning">
    <strong>⚠ Security action required</strong>
    <p>For security reasons, your credentials must be reset. Click the button below to proceed.</p>
  </div>
{/if}

<div class="reset-row">
  <div>
    <span class="form-label">Reset All Credentials</span>
    <p class="reset-description">Wipes all scoped API keys, connection tokens, and active sessions. Every paired Foundry module will need to re-pair. Use this only if you suspect a security breach.</p>
  </div>
  <button class="btn btn-sm btn-danger" onclick={handleReset} disabled={regenLoading}>
    {regenLoading ? 'Resetting...' : 'Reset Credentials'}
  </button>
</div>

<style>
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

  .reset-row {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .reset-description {
    margin: 0.25rem 0 0;
    font-size: 0.875rem;
    color: var(--color-text-muted, #666);
    max-width: 480px;
  }
</style>
