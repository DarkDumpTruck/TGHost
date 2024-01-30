import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import ElementPlus from 'unplugin-element-plus/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  server: {
    proxy: {
      '/examples': {
        target: 'http://localhost:40123',
        changeOrigin: true
      },
      '/api': {
        target: 'http://localhost:40123',
        changeOrigin: true
      },
      '/ws': {
        target: 'ws://localhost:40123',
        changeOrigin: true,
        ws: true
      }
    }
  },
  plugins: [
    vue(),
    ElementPlus({}),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
