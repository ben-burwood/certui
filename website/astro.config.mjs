// @ts-check
import { defineConfig } from 'astro/config';
import icon from 'astro-icon';
import tailwindcss from '@tailwindcss/vite';

// Deploys to GitHub Pages at https://ben-burwood.github.io/certui
export default defineConfig({
  site: 'https://ben-burwood.github.io',
  base: '/certui',
  integrations: [icon()],
  vite: {
    plugins: [tailwindcss()],
  },
});
