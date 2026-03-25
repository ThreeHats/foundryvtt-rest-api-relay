<script lang="ts">
  import { user, clearUser } from '../../lib/auth';
  import { exportData, deleteAccount } from '../../lib/api';

  interface Props {
    onAccountDeleted: () => void;
  }

  let { onAccountDeleted }: Props = $props();

  let exportLoading = $state(false);

  async function handleExport() {
    exportLoading = true;
    const result = await exportData();
    exportLoading = false;

    if (result.ok) {
      const blob = new Blob([JSON.stringify(result.data, null, 2)], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'my-foundry-api-data.json';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } else {
      alert('Failed to export data. Please try again.');
    }
  }

  async function handleDelete() {
    if (!confirm('Are you sure you want to delete your account?\n\nThis action is PERMANENT and cannot be undone.\nAll your data will be deleted.')) return;

    const email = $user?.email;
    const confirmEmail = prompt('Please enter your email address to confirm account deletion:');
    if (!confirmEmail) return;

    if (confirmEmail !== email) {
      alert('Email address does not match. Account deletion cancelled.');
      return;
    }

    const password = prompt('Please enter your password to confirm account deletion:');
    if (!password) return;

    if (!confirm('FINAL WARNING: Your account will be permanently deleted.\n\nClick OK to proceed with deletion.')) return;

    const result = await deleteAccount(confirmEmail, password);

    if (result.ok) {
      alert("Your account has been deleted. We're sorry to see you go.");
      clearUser();
      onAccountDeleted();
    } else {
      alert(result.error);
    }
  }
</script>

<div class="privacy-section">
  <h3 class="section-title">Data & Privacy</h3>
  <p class="text-muted mb-1" style="font-size: 0.85rem;">
    Export your data or delete your account in accordance with GDPR/CCPA.
  </p>
  <div class="flex gap-1 flex-wrap">
    <button class="btn btn-secondary btn-sm" onclick={handleExport} disabled={exportLoading}>
      {exportLoading ? 'Exporting...' : 'Export My Data'}
    </button>
    <button class="btn btn-danger btn-sm" onclick={handleDelete}>
      Delete Account
    </button>
  </div>
</div>

<style>
  .section-title {
    font-size: 0.95rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }
</style>
