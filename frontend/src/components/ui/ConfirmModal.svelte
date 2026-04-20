<script lang="ts">
  interface Props {
    open: boolean;
    title: string;
    message: string;
    confirmLabel?: string;
    dangerous?: boolean;
    onConfirm: () => void;
    onCancel: () => void;
  }

  let { open, title, message, confirmLabel = 'Confirm', dangerous = false, onConfirm, onCancel }: Props = $props();

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') onCancel();
  }
</script>

{#if open}
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="modal-backdrop"
    role="dialog"
    aria-modal="true"
    onkeydown={handleKeydown}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div class="backdrop-click" onclick={onCancel}></div>
    <div class="modal">
      <h3 class="modal-title">{title}</h3>
      <p class="modal-message">{message}</p>
      <div class="modal-actions">
        <button class="btn btn-secondary" onclick={onCancel}>Cancel</button>
        <button class="btn" class:btn-danger={dangerous} class:btn-primary={!dangerous} onclick={onConfirm}>
          {confirmLabel}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
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

  .backdrop-click {
    position: absolute;
    inset: 0;
  }

  .modal {
    position: relative;
    background: var(--color-bg, #fff);
    border-radius: var(--radius-md);
    padding: 1.5rem;
    max-width: 420px;
    width: 100%;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .modal-title {
    font-size: 1rem;
    font-weight: 600;
    margin: 0;
  }

  .modal-message {
    font-size: 0.875rem;
    color: var(--color-text-secondary, #6b7280);
    margin: 0;
    line-height: 1.5;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 0.25rem;
  }
</style>
