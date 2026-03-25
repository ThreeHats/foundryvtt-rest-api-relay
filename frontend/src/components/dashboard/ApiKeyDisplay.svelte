<script lang="ts">
  import { user, updateUser } from '../../lib/auth';
  import { regenerateKey } from '../../lib/api';

  let masked = $state(true);
  let copyLabel = $state('Copy');
  let regenLoading = $state(false);

  let currentKey = $derived($user?.apiKey ?? '');

  function maskKey(key: string): string {
    if (!key || key.length < 12) return '••••••••••••';
    return key.substring(0, 8) + '••••••••' + key.substring(key.length - 4);
  }

  function toggleVisibility() {
    masked = !masked;
  }

  async function copyKey() {
    if (!currentKey) return;
    await navigator.clipboard.writeText(currentKey);
    copyLabel = 'Copied!';
    setTimeout(() => { copyLabel = 'Copy'; }, 1500);
  }

  async function handleRegenerate() {
    if (!confirm('Are you sure you want to regenerate your API key? This will invalidate your current key.')) return;

    const email = $user?.email;
    if (!email) return;

    const password = prompt('Enter your password to confirm:');
    if (!password) return;

    regenLoading = true;
    const result = await regenerateKey(email, password);
    regenLoading = false;

    if (result.ok) {
      updateUser({ apiKey: result.data.apiKey });
      masked = true;
      alert('API key regenerated successfully! Update any applications using the old key.');
    } else {
      alert(result.error);
    }
  }
</script>

<div class="key-row">
  <span class="form-label">API Key</span>
  <div class="key-actions">
    <code class="key-value">{masked ? maskKey(currentKey) : currentKey}</code>
    <button class="btn btn-sm btn-secondary" onclick={toggleVisibility}>
      {masked ? 'Show' : 'Hide'}
    </button>
    <button class="btn btn-sm btn-secondary" onclick={copyKey}>
      {copyLabel}
    </button>
    <button class="btn btn-sm btn-secondary" onclick={handleRegenerate} disabled={regenLoading}>
      {regenLoading ? 'Regenerating...' : 'Regenerate'}
    </button>
  </div>
</div>

<style>
  .key-row {
    margin-bottom: 1rem;
  }

  .key-actions {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    flex-wrap: wrap;
    margin-top: 0.25rem;
  }

  .key-value {
    background: var(--color-bg-sunken);
    padding: 0.375rem 0.625rem;
    border-radius: var(--radius-sm);
    font-size: 0.8rem;
    word-break: break-all;
    flex: 1;
    min-width: 200px;
  }
</style>
