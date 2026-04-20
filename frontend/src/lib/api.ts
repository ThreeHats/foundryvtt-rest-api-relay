import { get } from 'svelte/store';
import { sessionToken } from './auth';
import type { AuthResponse, UserData, ScopedKey, ConnectedClient, Player, ConnectionToken, KnownClient, KnownUser, Credential, NotificationSettings, ApiKeyNotificationSettings, KnownClientNotificationSettings, ConnectionLog, KeyRequestDetails, PairRequestDetails, RemoteRequestLog, ActivityEvent } from './types';

// Dashboard requests authenticate via Authorization: Bearer <session-token>.
// The plaintext master API key is NEVER held by the frontend after the
// one-time registration / regeneration modal — it lives only in the user's
// password manager. The session token is a separate, revocable credential
// minted by /auth/login (or /auth/register) and used for everything else.
function getHeaders(includeAuth = true): Record<string, string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (includeAuth) {
    const token = get(sessionToken);
    if (token) headers['Authorization'] = `Bearer ${token}`;
  }
  return headers;
}

function authHeaders(): Record<string, string> {
  const token = get(sessionToken);
  return token ? { 'Authorization': `Bearer ${token}` } : {};
}

async function handleResponse<T>(response: Response): Promise<{ ok: true; data: T } | { ok: false; error: string }> {
  const data = await response.json().catch(() => ({ error: 'Failed to parse response' }));
  if (response.ok) {
    return { ok: true, data: data as T };
  }
  return { ok: false, error: data.error || `Request failed (${response.status})` };
}

// ==================== Auth ====================

export async function register(email: string, password: string) {
  const body: Record<string, string> = { email, password };
  const res = await fetch('/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  // Returns AuthResponse: includes the one-time `apiKey` field that the
  // SignUpForm extracts for the modal display, plus the persistent
  // sessionToken used for ongoing dashboard auth.
  return handleResponse<AuthResponse>(res);
}

export async function login(email: string, password: string) {
  const body: Record<string, string> = { email, password };
  const res = await fetch('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  // Login returns sessionToken but NOT apiKey — the master key is never
  // re-displayed after the registration / regeneration modal.
  return handleResponse<UserData>(res);
}

export async function logout() {
  const res = await fetch('/auth/logout', {
    method: 'POST',
    headers: authHeaders(),
  });
  return handleResponse<{ success: boolean }>(res);
}

export async function fetchUserData() {
  const res = await fetch('/auth/user-data', { headers: authHeaders() });
  return handleResponse<UserData>(res);
}

export async function fetchSubscriptionStatus() {
  const res = await fetch('/api/subscriptions/status', { headers: authHeaders() });
  return handleResponse<{ subscriptionStatus: string; subscriptionEndsAt?: string }>(res);
}

export async function changePassword(currentPassword: string, newPassword: string) {
  const res = await fetch('/auth/change-password', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ currentPassword, newPassword }),
  });
  return handleResponse<{ message: string }>(res);
}

export async function regenerateKey(email: string, password: string) {
  const res = await fetch('/auth/regenerate-key', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  // Returns the new plaintext apiKey ONCE for the one-time modal. The
  // backend also invalidates all dashboard sessions for this user, so the
  // frontend MUST log out and prompt the user to log in again with their
  // new credentials after closing the modal.
  return handleResponse<{ apiKey: string }>(res);
}

export async function forgotPassword(email: string) {
  const res = await fetch('/auth/forgot-password', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  });
  return handleResponse<{ message: string }>(res);
}

export async function validateResetToken(token: string) {
  const res = await fetch(`/auth/validate-reset-token/${encodeURIComponent(token)}`);
  return handleResponse<{ valid: boolean }>(res);
}

export async function resetPassword(token: string, password: string) {
  const res = await fetch('/auth/reset-password', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ token, password }),
  });
  return handleResponse<{ message: string }>(res);
}

export async function exportData() {
  const res = await fetch('/auth/export-data', { headers: authHeaders() });
  return handleResponse<Record<string, unknown>>(res);
}

export async function deleteAccount(confirmEmail: string, password: string) {
  const res = await fetch('/auth/account', {
    method: 'DELETE',
    headers: getHeaders(),
    body: JSON.stringify({ confirmEmail, password }),
  });
  return handleResponse<{ message: string }>(res);
}

// ==================== Subscriptions ====================

export async function createCheckoutSession() {
  const res = await fetch('/api/subscriptions/create-checkout-session', {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse<{ url: string }>(res);
}

export async function createPortalSession() {
  const res = await fetch('/api/subscriptions/create-portal-session', {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse<{ url: string }>(res);
}

// ==================== Scoped Keys ====================

export async function fetchScopedKeys() {
  const res = await fetch('/auth/api-keys', { headers: authHeaders() });
  return handleResponse<{ keys: ScopedKey[] }>(res);
}

export async function createScopedKey(body: Record<string, unknown>) {
  const res = await fetch('/auth/api-keys', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(body),
  });
  return handleResponse<{ key: string; id: number }>(res);
}

export async function updateScopedKey(id: number, body: Record<string, unknown>) {
  const res = await fetch(`/auth/api-keys/${id}`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify(body),
  });
  return handleResponse<{ message: string }>(res);
}

export async function deleteScopedKey(id: number) {
  const res = await fetch(`/auth/api-keys/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

export async function regenerateScopedKey(id: number) {
  const res = await fetch(`/auth/api-keys/${id}/regenerate`, {
    method: 'POST',
    headers: authHeaders(),
  });
  return handleResponse<{ key: string }>(res);
}

// ==================== Clients & Players ====================

export async function fetchClients() {
  const res = await fetch('/clients', { headers: authHeaders() });
  return handleResponse<{ clients: ConnectedClient[] }>(res);
}

export async function fetchPlayers(clientId: string) {
  const res = await fetch(`/players?clientId=${encodeURIComponent(clientId)}`, {
    headers: authHeaders(),
  });
  return handleResponse<{ users: Player[] }>(res);
}

// ==================== Connection Tokens ====================

export async function fetchConnectionTokens() {
  const res = await fetch('/auth/connection-tokens', { headers: authHeaders() });
  return handleResponse<{ tokens: ConnectionToken[] }>(res);
}

// generatePairingCode mints a 6-char code that the Foundry module exchanges
// for a connection token. Optional params:
//   - clientId: bind the code to an existing known clientId (used for the
//     "add this browser" flow when a second GM joins an already-paired world)
//   - allowedTargetClients: clientIds the resulting connection token may
//     invoke remote-request operations against (cross-world tunnel opt-in)
//   - remoteScopes: scope strings the resulting token holds for those
//     cross-world operations (e.g. "entity:write", "user:write")
export async function generatePairingCode(opts?: {
  clientId?: string;
  allowedTargetClients?: string[];
  remoteScopes?: string[];
}) {
  const body: Record<string, unknown> = {};
  if (opts?.clientId) body.clientId = opts.clientId;
  if (opts?.allowedTargetClients) body.allowedTargetClients = opts.allowedTargetClients;
  if (opts?.remoteScopes) body.remoteScopes = opts.remoteScopes;
  const res = await fetch('/auth/connection-tokens', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(body),
  });
  return handleResponse<{ code: string; expiresAt: string }>(res);
}

// updateConnectionTokenPermissions updates a connection token's name and/or IP allowlist.
// Cross-world settings are now world-level; use updateKnownClientCrossWorldSettings instead.
export async function updateConnectionTokenPermissions(id: number, data: { name?: string; allowedIps?: string }) {
  const res = await fetch(`/auth/connection-tokens/${id}`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse<{ success: boolean }>(res);
}

// updateKnownClientCrossWorldSettings sets the cross-world tunneling permissions
// for a world (KnownClient). All browsers (connection tokens) for this world
// inherit these settings — targets, scopes, and rate limit are world-level.
export async function updateKnownClientCrossWorldSettings(id: number, data: {
  allowedTargetClients?: string[];
  remoteScopes?: string[];
  remoteRequestsPerHour?: number;
}) {
  const res = await fetch(`/auth/known-clients/${id}`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse<{ success: boolean }>(res);
}

export async function deleteConnectionToken(id: number) {
  const res = await fetch(`/auth/connection-tokens/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

// ==================== Known Clients ====================

export async function fetchKnownUsers(knownClientId: number) {
  const res = await fetch(`/auth/known-clients/${knownClientId}/users`, { headers: authHeaders() });
  return handleResponse<{ users: KnownUser[] }>(res);
}

export async function fetchKnownClients() {
  const res = await fetch('/auth/known-clients', { headers: authHeaders() });
  return handleResponse<{ clients: KnownClient[] }>(res);
}

export async function renameConnectionToken(id: number, name: string) {
  const res = await fetch(`/auth/connection-tokens/${id}`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify({ name }),
  });
  return handleResponse<{ success: boolean }>(res);
}

export async function deleteKnownClient(id: number) {
  const res = await fetch(`/auth/known-clients/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

// setKnownClientAutoStart toggles whether the relay will spawn a headless
// session for this clientId in response to incoming remote-request from a
// sibling client (when this client is currently offline).
export async function setKnownClientAutoStart(id: number, enabled: boolean) {
  const res = await fetch(`/auth/known-clients/${id}/auto-start`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify({ enabled }),
  });
  return handleResponse<{ success: boolean }>(res);
}

// setKnownClientCredential links a Credential to a KnownClient. The headless
// auto-start path uses this credential to log into the target Foundry world.
// Pass null to clear the link (auto-start will fall back to the user's first
// credential if exactly one exists).
export async function setKnownClientCredential(id: number, credentialId: number | null) {
  const res = await fetch(`/auth/known-clients/${id}/credential`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify({ credentialId }),
  });
  return handleResponse<{ success: boolean }>(res);
}

// ==================== Credentials ====================

export async function fetchCredentials() {
  const res = await fetch('/auth/credentials', { headers: authHeaders() });
  return handleResponse<{ credentials: Credential[] }>(res);
}

export async function createCredential(data: { name: string; foundryUrl: string; foundryUsername: string; foundryPassword: string }) {
  const res = await fetch('/auth/credentials', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse<{ credential: Credential }>(res);
}

export async function updateCredential(id: number, data: Record<string, unknown>) {
  const res = await fetch(`/auth/credentials/${id}`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse<{ credential: Credential }>(res);
}

export async function deleteCredential(id: number) {
  const res = await fetch(`/auth/credentials/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

// ==================== Notification Settings ====================

export async function fetchNotificationSettings() {
  const res = await fetch('/auth/notification-settings', { headers: authHeaders() });
  return handleResponse<NotificationSettings>(res);
}

export async function updateNotificationSettings(data: Partial<NotificationSettings>) {
  const res = await fetch('/auth/notification-settings', {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse<{ message: string }>(res);
}

export async function testNotification(params?: { discordWebhookUrl?: string; notifyEmail?: string }) {
  const res = await fetch('/auth/notification-settings/test', {
    method: 'POST',
    headers: getHeaders(),
    body: params ? JSON.stringify(params) : undefined,
  });
  return handleResponse<{ message: string }>(res);
}

// ==================== Per-Key Notification Settings ====================

export async function fetchApiKeyNotificationSettings(keyId: number) {
  const res = await fetch(`/auth/api-keys/${keyId}/notification-settings`, { headers: authHeaders() });
  return handleResponse<ApiKeyNotificationSettings>(res);
}

export async function updateApiKeyNotificationSettings(keyId: number, data: Partial<ApiKeyNotificationSettings>) {
  const res = await fetch(`/auth/api-keys/${keyId}/notification-settings`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse<{ message: string }>(res);
}

export async function deleteApiKeyNotificationSettings(keyId: number) {
  const res = await fetch(`/auth/api-keys/${keyId}/notification-settings`, {
    method: 'DELETE',
    headers: authHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

export async function testApiKeyNotification(keyId: number) {
  const res = await fetch(`/auth/api-keys/${keyId}/notification-settings/test`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

// ==================== World (KnownClient) Notification Settings ====================

export async function fetchWorldNotificationSettings(id: number) {
  const res = await fetch(`/auth/known-clients/${id}/notification-settings`, {
    headers: authHeaders(),
  });
  return handleResponse<KnownClientNotificationSettings>(res);
}

export async function updateWorldNotificationSettings(id: number, data: Partial<KnownClientNotificationSettings>) {
  const res = await fetch(`/auth/known-clients/${id}/notification-settings`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse<KnownClientNotificationSettings>(res);
}

export async function deleteWorldNotificationSettings(id: number) {
  const res = await fetch(`/auth/known-clients/${id}/notification-settings`, {
    method: 'DELETE',
    headers: authHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

export async function testWorldNotification(id: number) {
  const res = await fetch(`/auth/known-clients/${id}/notification-settings/test`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

// ==================== Connection Logs ====================

export async function fetchConnectionLogs(limit = 50, offset = 0) {
  const res = await fetch(`/auth/connection-logs?limit=${limit}&offset=${offset}`, {
    headers: authHeaders(),
  });
  return handleResponse<{ logs: ConnectionLog[]; total: number }>(res);
}

// ==================== Remote Request Logs ====================
//
// Audit trail of cross-world remote-request operations. Each entry records
// who (sourceTokenId/sourceClientId) called what (action) on which target,
// from which IP, with success/failure.

export async function fetchRemoteRequestLogs(limit = 50, offset = 0) {
  const res = await fetch(`/auth/remote-request-logs?limit=${limit}&offset=${offset}`, {
    headers: authHeaders(),
  });
  return handleResponse<{ logs: RemoteRequestLog[]; total: number }>(res);
}

// ==================== Activity Log ====================

export async function fetchActivity(params: {
  type?: string;
  world?: string;
  action?: string;
  success?: boolean | null;
  since?: string;
  until?: string;
  limit?: number;
  offset?: number;
} = {}) {
  const q = new URLSearchParams();
  if (params.type) q.set('type', params.type);
  if (params.world) q.set('world', params.world);
  if (params.action) q.set('action', params.action);
  if (params.success === true) q.set('success', 'true');
  else if (params.success === false) q.set('success', 'false');
  if (params.since) q.set('since', params.since);
  if (params.until) q.set('until', params.until);
  if (params.limit != null) q.set('limit', String(params.limit));
  if (params.offset != null) q.set('offset', String(params.offset));
  const res = await fetch(`/auth/activity?${q}`, { headers: authHeaders() });
  return handleResponse<{ events: ActivityEvent[]; total: number; offset: number; limit: number }>(res);
}

// ==================== Pair Requests ====================

export async function fetchPairRequest(code: string) {
  const res = await fetch(`/auth/pair-request/${encodeURIComponent(code)}`, {
    headers: authHeaders(),
  });
  return handleResponse<PairRequestDetails>(res);
}

export async function approvePairRequest(code: string, body?: {
  remoteScopes?: string[];
  allowedTargetClients?: string[];
  remoteRequestsPerHour?: number;
}) {
  const res = await fetch(`/auth/pair-request/${encodeURIComponent(code)}/approve`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(body ?? {}),
  });
  return handleResponse<{ success: boolean; upgraded?: boolean }>(res);
}

export async function denyPairRequest(code: string) {
  const res = await fetch(`/auth/pair-request/${encodeURIComponent(code)}/deny`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}

// ==================== Key Requests ====================

export async function fetchKeyRequest(code: string) {
  const res = await fetch(`/auth/key-request/${encodeURIComponent(code)}`, {
    headers: authHeaders(),
  });
  return handleResponse<KeyRequestDetails>(res);
}

export async function approveKeyRequest(code: string, data: Record<string, unknown>) {
  const res = await fetch(`/auth/key-request/${encodeURIComponent(code)}/approve`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(data),
  });
  return handleResponse<{ message: string; key?: string }>(res);
}

export async function denyKeyRequest(code: string) {
  const res = await fetch(`/auth/key-request/${encodeURIComponent(code)}/deny`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse<{ message: string }>(res);
}
