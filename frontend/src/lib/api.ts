import { get } from 'svelte/store';
import { apiKey } from './auth';
import type { UserData, ScopedKey, ConnectedClient, Player } from './types';

function getHeaders(includeAuth = true): Record<string, string> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (includeAuth) {
    const key = get(apiKey);
    if (key) headers['x-api-key'] = key;
  }
  return headers;
}

function authHeaders(): Record<string, string> {
  const key = get(apiKey);
  return key ? { 'x-api-key': key } : {};
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
  const res = await fetch('/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  return handleResponse<UserData>(res);
}

export async function login(email: string, password: string) {
  const res = await fetch('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });
  return handleResponse<UserData>(res);
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
