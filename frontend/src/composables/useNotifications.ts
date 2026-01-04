import type { Ref } from 'vue'

export interface NotificationOptions {
  type?: 'success' | 'error' | 'info' | 'warning'
  duration?: number
}

interface Notification extends Required<NotificationOptions> {
  id: string
  message: string
}

const notifications: Ref<Notification[]> = ref([])

export function useNotifications() {
  function addNotification(message: string, options: NotificationOptions = {}) {
    const notification: Notification = {
      id: Date.now().toString(),
      message,
      type: options.type ?? 'info',
      duration: options.duration ?? 4000,
    }

    notifications.value.push(notification)

    // Auto remove after duration
    setTimeout(() => {
      removeNotification(notification.id)
    }, notification.duration)

    return notification.id
  }

  function removeNotification(id: string) {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  function success(message: string) {
    return addNotification(message, { type: 'success' })
  }

  function error(message: string) {
    return addNotification(message, { type: 'error' })
  }

  function info(message: string) {
    return addNotification(message, { type: 'info' })
  }

  function warning(message: string) {
    return addNotification(message, { type: 'warning' })
  }

  return {
    notifications: readonly(notifications),
    addNotification,
    removeNotification,
    success,
    error,
    info,
    warning,
  }
}
