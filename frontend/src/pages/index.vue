<script setup lang="ts">
const userStore = useUserStore()

useHead({
  title: 'DSN - Digital Sticky Notes',
  meta: [
    { name: 'description', content: 'A modern note-taking application inspired by Google Keep' },
  ],
})
</script>

<template>
  <div class="mx-auto max-w-md">
    <div class="mb-8 text-center">
      <h1 class="mb-2 text-3xl text-gray-800 font-bold">
        Welcome to DSN
      </h1>
      <p class="text-gray-600">
        Your digital sticky notes, organized and accessible.
      </p>
    </div>

    <div class="space-y-4">
      <template v-if="userStore.isAuthenticated">
        <RouterLink to="/notes" class="block btn w-full text-center">
          View My Notes
        </RouterLink>
        <div class="text-center">
          <p class="text-sm text-gray-500">
            Welcome back, {{ userStore.user?.username }}!
          </p>
        </div>
      </template>
      <template v-else>
        <RouterLink to="/register" class="block btn w-full text-center">
          Get Started
        </RouterLink>
        <RouterLink to="/login" class="block btn btn-orange w-full text-center">
          Login
        </RouterLink>
        <div class="text-center">
          <p class="text-sm text-gray-500">
            Create an account to start organizing your thoughts.
          </p>
        </div>
      </template>
    </div>
  </div>
</template>
