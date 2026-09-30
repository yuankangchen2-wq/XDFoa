import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import qiankun from 'vite-plugin-qiankun'

export default defineConfig({
  plugins: [vue(), qiankun('app-portal', { useDevMode: true })],
  server: {
    port: 5177,
    proxy: {
      '/api/order': { target: 'http://localhost:8086', changeOrigin: true },
      '/api/mdm': { target: 'http://localhost:8081', changeOrigin: true }
    }
  }
})
