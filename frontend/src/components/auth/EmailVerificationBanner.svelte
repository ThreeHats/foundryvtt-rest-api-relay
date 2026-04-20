<script lang="ts">
  import { resendVerificationEmail } from '../../lib/api';
  import { updateUser } from '../../lib/auth';

  let sending = $state(false);
  let sent = $state(false);
  let error = $state('');

  async function handleResend() {
    sending = true;
    error = '';
    const result = await resendVerificationEmail();
    sending = false;
    if (result.ok) {
      sent = true;
    } else {
      error = result.error;
    }
  }
</script>

<div class="verify-banner">
  <span class="banner-icon">&#9888;</span>
  <span class="banner-text">
    Your email address is not verified. Please check your inbox for a verification link.
  </span>
  {#if !sent}
    <button class="btn-resend" onclick={handleResend} disabled={sending}>
      {sending ? 'Sending...' : 'Resend email'}
    </button>
  {:else}
    <span class="sent-confirm">Sent! Check your inbox.</span>
  {/if}
  {#if error}
    <span class="banner-error">{error}</span>
  {/if}
</div>

<style>
  .verify-banner {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    background: var(--color-warning-bg, #fff8e1);
    border-bottom: 2px solid var(--color-warning, #f39c12);
    padding: 0.65rem 1.5rem;
    font-size: 0.9rem;
    flex-wrap: wrap;
  }

  .banner-icon {
    font-size: 1.1rem;
    color: var(--color-warning, #f39c12);
    flex-shrink: 0;
  }

  .banner-text {
    flex: 1;
    color: var(--color-text, #333);
  }

  .btn-resend {
    background: none;
    border: 1px solid var(--color-warning, #f39c12);
    color: var(--color-warning, #c87800);
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-sm, 4px);
    cursor: pointer;
    font-size: 0.85rem;
    white-space: nowrap;
  }

  .btn-resend:hover:not(:disabled) {
    background: var(--color-warning, #f39c12);
    color: #fff;
  }

  .btn-resend:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .sent-confirm {
    color: var(--color-success, #27ae60);
    font-size: 0.85rem;
    font-weight: 500;
  }

  .banner-error {
    color: var(--color-error, #e74c3c);
    font-size: 0.85rem;
  }
</style>
