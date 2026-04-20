import { writable, derived } from 'svelte/store';
import type { UserData } from './types';

const STORAGE_KEY = 'foundryApiUser';

function loadUser(): UserData | null {
  if (typeof localStorage === 'undefined') return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    // Migration: legacy entries had `apiKey`. The new model uses
    // `sessionToken` exclusively. If we find a legacy entry, drop it —
    // the user will be logged out and prompted to sign in again, where
    // they'll get a fresh session token via the new auth flow.
    if (!parsed.sessionToken) return null;
    return parsed;
  } catch {
    return null;
  }
}

export const user = writable<UserData | null>(loadUser());

export const isLoggedIn = derived(user, ($user) => $user !== null);

// Bearer-token used by every authenticated dashboard request. The plaintext
// master API key is NEVER held in this store — only the session token,
// which is a separate credential issued by /auth/login and revocable via
// /auth/logout.
export const sessionToken = derived(user, ($user) => $user?.sessionToken ?? null);

export function saveUser(data: UserData) {
  // Defensive: strip any apiKey field that may have leaked into the user
  // object (the register/regenerate handlers return apiKey for the one-time
  // modal, but it must NEVER be persisted to the user store).
  const { ...rest } = data as UserData & { apiKey?: string };
  delete (rest as { apiKey?: string }).apiKey;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(rest));
  user.set(rest);
}

export function updateUser(partial: Partial<UserData>) {
  user.update((current) => {
    if (!current) return current;
    const updated = { ...current, ...partial };
    // Same defensive scrub on update.
    delete (updated as UserData & { apiKey?: string }).apiKey;
    localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
    return updated;
  });
}

export function clearUser() {
  localStorage.removeItem(STORAGE_KEY);
  user.set(null);
}
