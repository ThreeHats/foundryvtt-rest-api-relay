<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchNotificationSettings, updateNotificationSettings, testNotification } from '../../lib/api';

  let notifyOnNewClientConnect = $state(true);
  let notifyOnConnect = $state(true);
  let notifyOnDisconnect = $state(false);
  let notifyOnMetadataMismatch = $state(true);
  let notifyOnSettingsChange = $state(false);
  let notifyOnExecuteJs = $state(false);
  let notifyOnMacroExecute = $state(false);
  let logCrossWorldRequests = $state(true);
  let discordWebhookUrl = $state('');
  let notifyEmail = $state('');
  let notificationDebounceWindowSecs = $state(0);
  let remoteRequestBatchWindowSecs = $state(300);
  let smtpAvailable = $state(false);

  let loading = $state(true);
  let saving = $state(false);
  let testing = $state(false);
  let message = $state('');
  let messageType = $state<'success' | 'error'>('error');

  onMount(() => {
    loadSettings();
  });

  async function loadSettings() {
    loading = true;
    const result = await fetchNotificationSettings();
    loading = false;
    if (result.ok) {
      notifyOnNewClientConnect = result.data.notifyOnNewClientConnect;
      notifyOnConnect = result.data.notifyOnConnect;
      notifyOnDisconnect = result.data.notifyOnDisconnect;
      notifyOnMetadataMismatch = result.data.notifyOnMetadataMismatch;
      notifyOnSettingsChange = result.data.notifyOnSettingsChange;
      notifyOnExecuteJs = result.data.notifyOnExecuteJs;
      notifyOnMacroExecute = result.data.notifyOnMacroExecute;
      logCrossWorldRequests = result.data.logCrossWorldRequests ?? true;
      discordWebhookUrl = result.data.discordWebhookUrl || '';
      notifyEmail = result.data.notifyEmail || '';
      notificationDebounceWindowSecs = result.data.notificationDebounceWindowSecs ?? 0;
      remoteRequestBatchWindowSecs = result.data.remoteRequestBatchWindowSecs ?? 300;
      smtpAvailable = result.data.smtpAvailable;
    }
  }

  async function handleSave(e: Event) {
    e.preventDefault();
    saving = true;
    message = '';

    const result = await updateNotificationSettings({
      notifyOnNewClientConnect,
      notifyOnConnect,
      notifyOnDisconnect,
      notifyOnMetadataMismatch,
      notifyOnSettingsChange,
      notifyOnExecuteJs,
      notifyOnMacroExecute,
      logCrossWorldRequests,
      discordWebhookUrl: discordWebhookUrl.trim() || '',
      notifyEmail: notifyEmail.trim() || '',
      notificationDebounceWindowSecs,
      remoteRequestBatchWindowSecs,
    });

    saving = false;

    if (result.ok) {
      message = 'Settings saved.';
      messageType = 'success';
      setTimeout(() => { message = ''; }, 3000);
    } else {
      message = result.error;
      messageType = 'error';
    }
  }

  async function handleTest() {
    testing = true;
    message = '';

    // Pass the current (unsaved) form values so the user can test before saving.
    const result = await testNotification({
      discordWebhookUrl: discordWebhookUrl.trim() || undefined,
      notifyEmail: notifyEmail.trim() || undefined,
    });

    testing = false;

    if (result.ok) {
      message = 'Test notification sent!';
      messageType = 'success';
      setTimeout(() => { message = ''; }, 3000);
    } else {
      message = result.error;
      messageType = 'error';
    }
  }
</script>

<h2 class="page-title">Notifications</h2>

<div class="card">
  <h3 class="card-title mb-2">Notification Settings</h3>

  {#if loading}
    <p class="text-muted">Loading settings...</p>
  {:else}
    <form onsubmit={handleSave}>
      <div class="form-group">
        <span class="form-label">Notify me when&hellip;</span>
        <div class="checkbox-list">
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={notifyOnNewClientConnect} />
            <span>A new Foundry world connects for the first time</span>
          </label>
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={notifyOnConnect} />
            <span>Any Foundry client connects or reconnects</span>
          </label>
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={notifyOnDisconnect} />
            <span>A Foundry client disconnects</span>
          </label>
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={notifyOnMetadataMismatch} />
            <span>A suspicious connection occurs (metadata mismatch)</span>
          </label>
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={notifyOnSettingsChange} />
            <span>Foundry settings change</span>
          </label>
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={notifyOnExecuteJs} />
            <span>An execute-js call is made</span>
          </label>
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={notifyOnMacroExecute} />
            <span>A macro-execute call is made</span>
          </label>
        </div>
      </div>

      <div class="form-group">
        <span class="form-label">Activity Log</span>
        <div class="checkbox-group">
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={logCrossWorldRequests} />
            <span>Log cross-world requests</span>
          </label>
          <p class="form-hint">Store a record of each cross-world API call made by your connection tokens. Disable to stop logging (existing logs are unaffected).</p>
        </div>
      </div>

      <div class="form-group">
        <label class="form-label" for="ns-discord">Discord Webhook URL</label>
        <input class="form-input" type="url" id="ns-discord" bind:value={discordWebhookUrl} placeholder="https://discord.com/api/webhooks/..." />
        <p class="form-hint">Receive notifications as Discord messages. Create a webhook in your Discord server under Settings &rarr; Integrations &rarr; Webhooks.</p>
      </div>

      <div class="form-group">
        <label class="form-label" for="ns-email" class:text-muted={!smtpAvailable}>Notification Email</label>
        <input
          class="form-input"
          type="email"
          id="ns-email"
          bind:value={notifyEmail}
          placeholder="you@example.com"
          disabled={!smtpAvailable}
        />
        <p class="form-hint">
          {#if smtpAvailable}
            Receive notifications via email.
          {:else}
            Email notifications require SMTP to be configured on the relay server.
          {/if}
        </p>
      </div>

      <div class="form-group">
        <span class="form-label">Notification Frequency</span>
        <div class="frequency-grid">
          <div class="frequency-row">
            <label class="form-label-sm" for="ns-debounce">Debounce window (seconds)</label>
            <input
              class="form-input form-input-sm"
              type="number"
              id="ns-debounce"
              min="0"
              bind:value={notificationDebounceWindowSecs}
            />
            <p class="form-hint">Suppress duplicate notifications of the same event for the same client within this window. Set to 0 to disable.</p>
          </div>
          <div class="frequency-row">
            <label class="form-label-sm" for="ns-batch">Cross-world batch window (seconds)</label>
            <input
              class="form-input form-input-sm"
              type="number"
              id="ns-batch"
              min="0"
              bind:value={remoteRequestBatchWindowSecs}
            />
            <p class="form-hint">Group cross-world activity into a single summary notification after this window expires. Default: 300 (5 minutes).</p>
          </div>
        </div>
      </div>

      <div class="flex gap-1 mt-2">
        <button class="btn btn-primary" type="submit" disabled={saving}>
          {saving ? 'Saving...' : 'Save Settings'}
        </button>
        <button class="btn btn-secondary" type="button" onclick={handleTest} disabled={testing}>
          {testing ? 'Sending...' : 'Send Test'}
        </button>
      </div>

      {#if message}
        <div class="alert mt-1" class:alert-success={messageType === 'success'} class:alert-error={messageType === 'error'}>
          {message}
        </div>
      {/if}
    </form>
  {/if}
</div>

<style>
  .checkbox-label {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
  }

  .checkbox-label input[type="checkbox"] {
    width: 1rem;
    height: 1rem;
    accent-color: var(--color-primary);
  }

  .checkbox-list {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
    margin-top: 0.5rem;
  }

  .frequency-grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-top: 0.5rem;
  }

  .frequency-row {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .form-label-sm {
    font-size: 0.8125rem;
    font-weight: 500;
    color: var(--color-text-secondary, inherit);
  }

  .form-input-sm {
    max-width: 8rem;
  }

  .text-muted {
    opacity: 0.6;
  }
</style>
