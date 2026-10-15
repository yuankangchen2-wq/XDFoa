import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api/bff': { target: 'http://localhost:8095', changeOrigin: true },
      '/api/scrm': { target: 'http://localhost:8096', changeOrigin: true },
      '/api/crm': { target: 'http://localhost:8087', changeOrigin: true },
      '/api/opp': { target: 'http://localhost:8088', changeOrigin: true },
      '/api/sales': { target: 'http://localhost:8089', changeOrigin: true },
      '/api/ticket': { target: 'http://localhost:8090', changeOrigin: true },
      '/api/analytics': { target: 'http://localhost:8093', changeOrigin: true }
    }
  }
})
