import { writable, derived } from 'svelte/store';
import type { UserData } from './types';

const STORAGE_KEY = 'foundryApiUser';

function loadUser(): UserData | null {
  if (typeof localStorage === 'undefined') return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export const user = writable<UserData | null>(loadUser());

export const isLoggedIn = derived(user, ($user) => $user !== null);

export const apiKey = derived(user, ($user) => $user?.apiKey ?? null);

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
