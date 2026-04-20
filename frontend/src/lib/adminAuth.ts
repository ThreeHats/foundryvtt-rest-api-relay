/**
 * Admin dashboard auth.
 *
 * Admin auth is COMPLETELY separate from the regular user auth (localStorage + x-api-key):
 * - Uses an HttpOnly __Host-admin_token cookie set by /admin/auth/login
 * - CSRF token returned in response body / header, stored ONLY in memory (never localStorage)
 * - All requests go through `adminFetch` which attaches the CSRF header on mutations
 *
 * On app load, `loadAdminSession()` calls /admin/auth/me — if it returns 200,
 * the cookie is still valid and we mark the user as an active admin.
 */

import { writable, derived, get } from 'svelte/store';
import type { AdminSession } from './types';

const adminSessionStore = writable<AdminSession | null>(null);

export const adminSession = adminSessionStore;
export const isAdmin = derived(adminSessionStore, ($s) => $s !== null);

// CSRF token is held in memory only — never persist to localStorage.
let csrfToken = '';

function setSession(session: AdminSession | null) {
  adminSessionStore.set(session);
  csrfToken = session?.csrfToken ?? '';
}

/**
 * Attempt to restore an admin session from an existing cookie.
 * Called on app startup.
 */
export async function loadAdminSession(): Promise<void> {
  try {
    const response = await fetch('/admin/auth/me', {
      method: 'GET',
      credentials: 'include',
    });
    if (!response.ok) {
      setSession(null);
      return;
    }
    const data = await response.json();
    if (data?.role === 'admin') {
      // /me re-issues the CSRF token in the X-CSRF-Token header so mutations
      // continue to work after a page reload (the token is memory-only).
      const token = response.headers.get('X-CSRF-Token') ?? csrfToken;
      setSession({ email: data.email, role: 'admin', csrfToken: token });
    }
  } catch {
    setSession(null);
  }
}

export async function adminLogin(email: string, password: string): Promise<{ ok: true } | { ok: false; error: string }> {
  try {
    const response = await fetch('/admin/auth/login', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
    if (!response.ok) {
      const data = await response.json().catch(() => ({}));
      return { ok: false, error: data?.error ?? `HTTP ${response.status}` };
    }
    const data = await response.json();
    const token = response.headers.get('X-CSRF-Token') ?? data?.csrfToken ?? '';
    setSession({ email: data.email, role: data.role, csrfToken: token });
    return { ok: true };
  } catch (err: any) {
    return { ok: false, error: err?.message ?? 'Network error' };
  }
}

export async function adminLogout(): Promise<void> {
  try {
    await fetch('/admin/auth/logout', {
      method: 'POST',
      credentials: 'include',
      headers: csrfToken ? { 'X-CSRF-Token': csrfToken } : undefined,
    });
  } finally {
    setSession(null);
  }
}

/**
 * Fetch wrapper for admin endpoints. Attaches credentials and CSRF token.
 */
export async function adminFetch(path: string, options: RequestInit = {}): Promise<Response> {
  const headers = new Headers(options.headers ?? {});
  const method = (options.method ?? 'GET').toUpperCase();
  const isMutation = method !== 'GET' && method !== 'HEAD' && method !== 'OPTIONS';
  if (isMutation && csrfToken) {
    headers.set('X-CSRF-Token', csrfToken);
  }
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }
  return fetch(path, {
    ...options,
    headers,
    credentials: 'include',
  });
}

/**
 * Helper for JSON GET requests.
 */
export async function adminGet<T>(path: string): Promise<T> {
  const r = await adminFetch(path);
  if (!r.ok) throw new Error(`HTTP ${r.status}`);
  return (await r.json()) as T;
}

/**
 * Helper for JSON POST/PATCH/DELETE requests.
 */
export async function adminMutate<T>(method: string, path: string, body?: any): Promise<T> {
  const r = await adminFetch(path, {
    method,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });
  if (!r.ok) {
    const errData = await r.json().catch(() => ({}));
    throw new Error(errData?.error ?? `HTTP ${r.status}`);
  }
  return (await r.json()) as T;
}

export function getCsrfToken(): string {
  return csrfToken;
}

export function clearAdminSession(): void {
  setSession(null);
}

// Re-export helper to access store value
export function getAdminSession(): AdminSession | null {
  return get(adminSessionStore);
}
