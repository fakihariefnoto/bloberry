import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  build: {
    outDir: 'dist',
  },
  server: {
    port: 5173,
    proxy: {
      // Dev environment talks to the local Go binary (04-environments.md).
      '/v1': 'http://localhost:8080',
      '/s': 'http://localhost:8080',
      '/health': 'http://localhost:8080',
      '/auth': 'http://localhost:8080',
    },
  },
})
