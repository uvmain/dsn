import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'

export function createAuthGuard() {
  return (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
    const userStore: { isAuthenticated: boolean } = useUserStore()

    const protectedRoutes = ['/notes']

    const guestOnlyRoutes = ['/login', '/register']

    if (protectedRoutes.includes(to.path) && userStore.isAuthenticated === false) {
      next('/login')
    }
    else if (guestOnlyRoutes.includes(to.path) && userStore.isAuthenticated === true) {
      next('/notes')
    }
    else {
      next()
    }
  }
}
