import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/account': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/financial': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/plaid': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/advisor': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/notification': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
