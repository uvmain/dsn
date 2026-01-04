<script setup lang="ts">
const userStore = useUserStore()
const router = useRouter()
const { success } = useNotifications()

async function handleLogout() {
  try {
    await userStore.logout()
    success('Successfully logged out')
    await router.push('/')
  }
  catch (err) {
    console.error('Logout error:', err)
    // Force logout on client side even if API call fails
    userStore.user = null
    await router.push('/')
  }
}
</script>

<template>
  <nav class="border-b bg-white shadow-sm">
    <div class="container mx-auto px-4">
      <div class="flex items-center justify-between py-4">
        <div class="flex items-center space-x-4">
          <RouterLink to="/" class="text-xl text-primary-600 font-bold">
            DSN
          </RouterLink>
          <span class="text-sm text-gray-500">Digital Sticky Notes</span>
        </div>

        <div class="flex items-center space-x-4">
          <template v-if="userStore.isAuthenticated">
            <RouterLink to="/notes" class="text-gray-600 hover:text-primary-600">
              Notes
            </RouterLink>
            <span v-if="userStore.user" class="text-sm text-gray-500">
              Welcome, {{ userStore.user.username }}
            </span>
            <button class="btn" @click="handleLogout">
              Logout
            </button>
          </template>
          <template v-else>
            <RouterLink to="/login" class="text-gray-600 hover:text-primary-600">
              Login
            </RouterLink>
            <RouterLink to="/register" class="btn">
              Register
            </RouterLink>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>
