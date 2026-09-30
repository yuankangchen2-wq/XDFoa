import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import qiankun from 'vite-plugin-qiankun'

export default defineConfig({
  plugins: [vue(), qiankun('app-iam', { useDevMode: true })],
  server: {
    port: 5175,
    proxy: {
      '/api/audit': { target: 'http://localhost:8094', changeOrigin: true },
      '/api': {
        target: 'http://localhost:8082',
        changeOrigin: true
      }
    }
  }
})
