import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  base: './',
  build: {
    outDir: '../visual_build',
    assetsDir: 'visual_assets',
    emptyOutDir: true,
  },
});
