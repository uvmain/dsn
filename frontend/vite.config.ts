import Vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [
    Vue({
      include: [/\.vue$/],
    }),
    UnoCSS(),
    AutoImport({
      imports: [
        'vue',
        '@vueuse/core',
        '@vueuse/head',
        'vue-router',
        'pinia',
      ],
      dts: 'auto-imports.d.ts',
      dirs: [
        'src/composables',
        'src/stores',
      ],
      vueTemplate: true,
    }),
    Icons(),
    Components({
      dts: 'components.d.ts',
      dirs: ['src/components'],
      resolvers: [
        IconsResolver({
          prefix: 'icon',
        }),
      ],
    }),
  ],
  resolve: {
    alias: {
      '~/': new URL('./src/', import.meta.url).pathname,
    },
  },
  ssgOptions: {
    script: 'async',
    formatting: 'minify',
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
