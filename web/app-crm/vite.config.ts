import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import qiankun from 'vite-plugin-qiankun'

export default defineConfig({
  plugins: [vue(), qiankun('app-crm', { useDevMode: true })],
  server: {
    port: 5178,
    proxy: {
      '/api/crm': { target: 'http://localhost:8087', changeOrigin: true },
      '/api/opp': { target: 'http://localhost:8088', changeOrigin: true },
      '/api/sales': { target: 'http://localhost:8089', changeOrigin: true },
      '/api/ticket': { target: 'http://localhost:8090', changeOrigin: true },
      '/api/analytics': { target: 'http://localhost:8093', changeOrigin: true }
    }
  }
})
