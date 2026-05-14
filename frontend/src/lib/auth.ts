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

export const sessionToken = derived(user, ($user) => $user?.sessionToken ?? null);

export function saveUser(data: UserData) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
  user.set(data);
}

export function updateUser(partial: Partial<UserData>) {
  user.update((current) => {
    if (!current) return current;
    const updated = { ...current, ...partial };
    localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
    return updated;
  });
}

export function clearUser() {
  localStorage.removeItem(STORAGE_KEY);
  user.set(null);
}
