<script setup lang="ts">
interface Note {
  id: number
  title: string
  content: string
  color: string
  pinned: boolean
  archived: boolean
  created_at: string
  updated_at: string
}

interface Props {
  note?: Note | null
}

interface Emits {
  save: [data: { title: string, content: string, color: string }]
  close: []
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const colors = [
  { name: 'White', value: '#ffffff' },
  { name: 'Yellow', value: '#fef3c7' },
  { name: 'Orange', value: '#fed7aa' },
  { name: 'Red', value: '#fecaca' },
  { name: 'Pink', value: '#f9a8d4' },
  { name: 'Purple', value: '#e9d5ff' },
  { name: 'Blue', value: '#bfdbfe' },
  { name: 'Green', value: '#bbf7d0' },
]

const form = reactive({
  title: props.note?.title || '',
  content: props.note?.content || '',
  color: props.note?.color || '#ffffff',
})

function handleSave() {
  // Emit the save event with form data
  emit('save', {
    title: form.title,
    content: form.content,
    color: form.color,
  })
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
    <div class="mx-4 max-w-md w-full rounded-lg bg-white shadow-xl">
      <div class="flex items-center justify-between border-b p-4">
        <h2 class="text-lg font-semibold">
          {{ note ? 'Edit Note' : 'New Note' }}
        </h2>
        <button
          class="icon-btn"
          @click="$emit('close')"
        >
          <i class="i-heroicons-x-mark h-5 w-5"></i>
        </button>
      </div>

      <form class="p-4 space-y-4" @submit.prevent="handleSave">
        <div>
          <input
            v-model="form.title"
            type="text"
            placeholder="Note title..."
            class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
        </div>

        <div>
          <textarea
            v-model="form.content"
            placeholder="Take a note..."
            rows="6"
            class="w-full resize-none border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-500"
          ></textarea>
        </div>

        <div>
          <label class="mb-2 block text-sm text-gray-700 font-medium">
            Color
          </label>
          <div class="flex space-x-2">
            <button
              v-for="color in colors"
              :key="color.value"
              type="button"
              class="h-8 w-8 border-2 border-gray-300 rounded-full hover:border-gray-400"
              :class="{ 'ring-2 ring-primary-500': form.color === color.value }"
              :style="{ backgroundColor: color.value }"
              @click="form.color = color.value"
            >
              <i
                v-if="form.color === color.value"
                class="i-heroicons-check h-4 w-4 text-gray-700"
              ></i>
            </button>
          </div>
        </div>

        <div class="flex justify-end pt-4 space-x-2">
          <button
            type="button"
            class="px-4 py-2 text-gray-600 hover:text-gray-800"
            @click="$emit('close')"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="btn"
          >
            Save
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
