import { defineConfig } from 'astro/config';
import svelte from '@astrojs/svelte';

export default defineConfig({
  output: 'static',
  integrations: [svelte()],
  outDir: '../public-dist',
  vite: {
    server: {
      proxy: {
        '/auth':       'http://localhost:3010',
        '/api':        'http://localhost:3010',
        '/clients':    'http://localhost:3010',
        '/players':    'http://localhost:3010',
        '/admin/api':  'http://localhost:3010',
        '/admin/auth': 'http://localhost:3010',
      }
    }
  }
});
