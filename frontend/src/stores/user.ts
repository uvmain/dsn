import type { User } from '~/types'
import { api } from '~/composables/useApi'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.is_admin ?? false)

  function setUser(userData: User) {
    user.value = userData
  }

  function clearUser() {
    user.value = null
  }

  async function register(username: string, email: string, password: string) {
    try {
      const userData = await api.register({ username, email, password })
      setUser(userData)
      return userData
    }
    catch (error) {
      console.error('Registration failed:', error)
      throw error
    }
  }

  async function login(username: string, password: string) {
    try {
      const userData = await api.login({ username, password })
      setUser(userData)
      return userData
    }
    catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  async function logout() {
    try {
      await api.logout()
      clearUser()
    }
    catch (error) {
      console.error('Logout failed:', error)
      // Clear user anyway on client side
      clearUser()
    }
  }

  async function checkAuth() {
    try {
      const userData = await api.checkAuth()
      setUser(userData)
      return userData
    }
    catch (error) {
      // User is not authenticated, clear any stale data
      clearUser()
      throw error
    }
  }

  return {
    user: readonly(user),
    isAuthenticated,
    isAdmin,
    setUser,
    clearUser,
    register,
    login,
    logout,
    checkAuth,
  }
})
