<script setup lang="ts">
const form = reactive({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')
const router = useRouter()
const userStore = useUserStore()
const { success, error: showError } = useNotifications()

async function handleLogin() {
  loading.value = true
  error.value = ''

  try {
    await userStore.login(form.username, form.password)

    success('Welcome back!')
    // Login successful, redirect to notes
    await router.push('/notes')
  }
  catch (err) {
    let errorMessage = 'Login failed. Please try again.'

    if (err instanceof Error) {
      if (err.message.includes('401') || err.message.includes('Unauthorized')) {
        errorMessage = 'Invalid username or password.'
      }
      else if (err.message.includes('network') || err.message.includes('fetch')) {
        errorMessage = 'Network error. Please check your connection.'
      }
    }

    error.value = errorMessage
    showError(errorMessage)
    console.error('Login error:', err)
  }
  finally {
    loading.value = false
  }
}

useHead({
  title: 'Login - DSN',
})
</script>

<template>
  <div class="mx-auto max-w-md">
    <div class="rounded-lg bg-white p-6 shadow-md">
      <h1 class="mb-6 text-center text-2xl font-bold">
        Login to DSN
      </h1>

      <form class="space-y-4" @submit.prevent="handleLogin">
        <div>
          <label for="username" class="mb-1 block text-sm text-gray-700 font-medium">
            Username
          </label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
        </div>

        <div>
          <label for="password" class="mb-1 block text-sm text-gray-700 font-medium">
            Password
          </label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="btn w-full"
        >
          {{ loading ? 'Logging in...' : 'Login' }}
        </button>

        <div v-if="error" class="text-center text-sm text-red-600">
          {{ error }}
        </div>
      </form>

      <div class="mt-4 text-center">
        <RouterLink to="/register" class="text-sm text-primary-600 hover:underline">
          Don't have an account? Register here
        </RouterLink>
      </div>
    </div>
  </div>
</template>
