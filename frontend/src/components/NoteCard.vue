<script setup lang="ts">
import type { Note } from '~/types'

interface Props {
  note: Note
}

interface Emits {
  edit: [note: Note]
  delete: [note: Note]
  togglePin: [note: Note]
  toggleArchive: [note: Note]
}

defineProps<Props>()
defineEmits<Emits>()

function formatDate(dateString: string) {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== new Date().getFullYear() ? 'numeric' : undefined,
  })
}
</script>

<template>
  <div
    class="cursor-pointer border rounded-lg bg-white p-4 shadow-sm transition-shadow hover:shadow-md"
    :style="{ backgroundColor: note.color }"
    @click="$emit('edit', note)"
  >
    <div class="mb-2 flex items-start justify-between">
      <h3 v-if="note.title" class="line-clamp-2 text-gray-800 font-semibold">
        {{ note.title }}
      </h3>
      <div class="ml-2 flex space-x-1">
        <button
          class="icon-btn hover:text-yellow-500"
          :class="{ 'text-yellow-500': note.pinned }"
          title="Toggle pin"
          @click.stop="$emit('togglePin', note)"
        >
          <i class="i-heroicons-bookmark-solid h-4 w-4"></i>
        </button>
        <button
          class="icon-btn hover:text-gray-500"
          title="Toggle archive"
          @click.stop="$emit('toggleArchive', note)"
        >
          <i class="i-heroicons-archive-box h-4 w-4"></i>
        </button>
        <button
          class="icon-btn hover:text-red-500"
          title="Delete"
          @click.stop="$emit('delete', note)"
        >
          <i class="i-heroicons-trash h-4 w-4"></i>
        </button>
      </div>
    </div>

    <p v-if="note.content" class="line-clamp-4 mb-3 text-sm text-gray-700">
      {{ note.content }}
    </p>

    <div v-if="note.tags && note.tags.length > 0" class="mb-2 flex flex-wrap gap-1">
      <span
        v-for="tag in note.tags"
        :key="tag.id"
        class="inline-flex items-center rounded-full px-2 py-0.5 text-xs text-white font-medium"
        :style="{ backgroundColor: tag.color }"
      >
        {{ tag.name }}
      </span>
    </div>

    <div class="text-xs text-gray-500">
      {{ formatDate(note.updated_at) }}
    </div>
  </div>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-clamp: 2;
  overflow: hidden;
}

.line-clamp-4 {
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  line-clamp: 4;
  overflow: hidden;
}
</style>
