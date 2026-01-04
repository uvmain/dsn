import { ViteSSG } from 'vite-ssg'
import App from './App.vue'
import { routes } from './router'

import '@unocss/reset/tailwind.css'
import 'uno.css'

export const createApp = ViteSSG(
  App,
  {
    routes,
    base: import.meta.env.BASE_URL,
  },
  (ctx) => {
    // Install all modules under `modules/`
    Object.values(import.meta.glob<{ install: any }>('./modules/*.ts', { eager: true }))
      .forEach(i => i.install?.(ctx))
  },
)
