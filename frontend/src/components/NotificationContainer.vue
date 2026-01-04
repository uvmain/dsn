<script setup lang="ts">
const { notifications, removeNotification } = useNotifications()
</script>

<template>
  <div
    v-if="notifications.length > 0"
    class="fixed right-4 top-4 z-50 space-y-2"
  >
    <div
      v-for="notification in notifications"
      :key="notification.id"
      class="max-w-sm transform rounded-lg px-4 py-3 shadow-lg transition-all duration-300 ease-in-out"
      :class="{
        'bg-green-500 text-white': notification.type === 'success',
        'bg-red-500 text-white': notification.type === 'error',
        'bg-blue-500 text-white': notification.type === 'info',
        'bg-yellow-500 text-black': notification.type === 'warning',
      }"
    >
      <div class="flex items-center justify-between">
        <span class="text-sm font-medium">{{ notification.message }}</span>
        <button
          class="ml-3 text-lg opacity-70 hover:opacity-100"
          @click="removeNotification(notification.id)"
        >
          ×
        </button>
      </div>
    </div>
  </div>
</template>
