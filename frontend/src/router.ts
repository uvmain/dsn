import type { RouteRecordRaw } from 'vue-router'

// Manual route definitions until unplugin-vue-router generates them properly
export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'index',
    component: async () => import('~/pages/index.vue'),
  },
  {
    path: '/login',
    name: 'login',
    component: async () => import('~/pages/login.vue'),
  },
  {
    path: '/register',
    name: 'register',
    component: async () => import('~/pages/register.vue'),
  },
  {
    path: '/notes',
    name: 'notes',
    component: async () => import('~/pages/notes.vue'),
    meta: { requiresAuth: true },
  },
]
