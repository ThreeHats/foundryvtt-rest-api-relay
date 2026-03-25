import { writable } from 'svelte/store';

export type Theme = 'light' | 'dark';

function getInitialTheme(): Theme {
  if (typeof localStorage !== 'undefined') {
    const stored = localStorage.getItem('theme');
    if (stored === 'light' || stored === 'dark') return stored;
  }
  if (typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    return 'dark';
  }
  return 'light';
}

export const theme = writable<Theme>(getInitialTheme());

export function toggleTheme() {
  theme.update((current) => {
    const next = current === 'light' ? 'dark' : 'light';
    localStorage.setItem('theme', next);
    document.documentElement.setAttribute('data-theme', next);
    return next;
  });
}

export function initTheme() {
  const t = getInitialTheme();
  document.documentElement.setAttribute('data-theme', t);
  theme.set(t);
}
