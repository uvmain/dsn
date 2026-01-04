import type { UserModule } from '~/types'

export const install: UserModule = ({ router }) => {
  router.beforeEach((to, from, next) => {
    const userStore = useUserStore()

    // Explicitly type meta to avoid 'any' issues
    const meta = to.meta as { requiresAuth?: boolean }

    // Check if route requires authentication
    if (meta.requiresAuth === true && Boolean(userStore?.isAuthenticated) === false) {
      // Redirect to login if trying to access protected route while not authenticated
      next('/login')
      return
    }

    // Redirect authenticated users away from login/register pages
    if ((to.path === '/login' || to.path === '/register') && userStore.isAuthenticated) {
      next('/notes')
      return
    }

    next()
  })
}
