import { writable } from 'svelte/store';

// Whether billing/subscription features are enabled on this relay instance.
// Set to true only when the server is configured with a Stripe secret key.
// Self-hosted instances default to false — no subscription UI is shown.
export const billingEnabled = writable<boolean>(false);

// Whether headless browser sessions are enabled on this relay instance.
// Defaults to true — the Docker image ships with ALLOW_HEADLESS=true.
// Set ALLOW_HEADLESS=false on the server to disable and hide all headless UI.
export const headlessEnabled = writable<boolean>(true);

// Fetches server configuration from /api/status. Called once on app mount.
export async function loadServerConfig(): Promise<void> {
  try {
    const res = await fetch('/api/status');
    if (res.ok) {
      const data = await res.json();
      billingEnabled.set(data.billingEnabled === true);
      headlessEnabled.set(data.headlessEnabled !== false); // default true if field absent
    }
  } catch {
    // Network error — leave billingEnabled as false, headlessEnabled as true (safe defaults)
  }
}
