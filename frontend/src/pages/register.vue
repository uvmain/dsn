<script setup lang="ts">
const form = reactive({
  username: '',
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref('')
const router = useRouter()
const userStore = useUserStore()
const { success, error: showError } = useNotifications()

async function handleRegister() {
  loading.value = true
  error.value = ''

  try {
    await userStore.register(form.username, form.email, form.password)

    success('Account created successfully! Welcome to DSN!')
    // Registration successful, redirect to notes
    await router.push('/notes')
  }
  catch (err) {
    // Handle different types of errors
    let errorMessage = 'Registration failed. Please try again.'

    if (err instanceof Error) {
      if (err.message.includes('400')) {
        errorMessage = 'Invalid input. Please check your information.'
      }
      else if (err.message.includes('409') || err.message.includes('duplicate')) {
        errorMessage = 'Username or email already exists.'
      }
      else if (err.message.includes('network') || err.message.includes('fetch')) {
        errorMessage = 'Network error. Please check your connection.'
      }
    }

    error.value = errorMessage
    showError(errorMessage)
    console.error('Registration error:', err)
  }
  finally {
    loading.value = false
  }
}

// Form validation
const isFormValid = computed(() => {
  return form.username.trim().length >= 3
    && form.email.includes('@')
    && form.password.length >= 6
})

useHead({
  title: 'Register - DSN',
})
</script>

<template>
  <div class="mx-auto max-w-md">
    <div class="rounded-lg bg-white p-6 shadow-md">
      <h1 class="mb-6 text-center text-2xl font-bold">
        Create Account
      </h1>

      <form class="space-y-4" @submit.prevent="handleRegister">
        <div>
          <label for="username" class="mb-1 block text-sm text-gray-700 font-medium">
            Username
          </label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            required
            minlength="3"
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500"
            placeholder="Enter username (min 3 characters)"
          >
        </div>

        <div>
          <label for="email" class="mb-1 block text-sm text-gray-700 font-medium">
            Email
          </label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500"
            placeholder="Enter your email address"
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
            minlength="6"
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500"
            placeholder="Enter password (min 6 characters)"
          >
        </div>

        <button
          type="submit"
          :disabled="loading || !isFormValid"
          class="btn w-full disabled:opacity-50"
        >
          {{ loading ? 'Creating Account...' : 'Register' }}
        </button>

        <div v-if="error" class="text-center text-sm text-red-600">
          {{ error }}
        </div>
      </form>

      <div class="mt-4 text-center">
        <RouterLink to="/login" class="text-sm text-primary-600 hover:underline">
          Already have an account? Login here
        </RouterLink>
      </div>
    </div>
  </div>
</template>
