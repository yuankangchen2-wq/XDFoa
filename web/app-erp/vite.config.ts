import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import qiankun from 'vite-plugin-qiankun'

export default defineConfig({
  plugins: [vue(), qiankun('app-erp', { useDevMode: true })],
  server: {
    port: 5176,
    proxy: {
      '/api/inventory': { target: 'http://localhost:8083', changeOrigin: true },
      '/api/purchase': { target: 'http://localhost:8084', changeOrigin: true },
      '/api/production': { target: 'http://localhost:8085', changeOrigin: true },
      '/api/finance': { target: 'http://localhost:8091', changeOrigin: true },
      '/api/cost': { target: 'http://localhost:8092', changeOrigin: true }
    }
  }
})
