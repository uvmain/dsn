<script setup lang="ts">
import type { Note, Tag } from '~/types'
import { api } from '~/composables/useApi'

interface Props {
  note?: Note | null
}

interface Emits {
  save: [data: { title: string, content: string, color: string, tags?: number[] }]
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
  selectedTagIds: props.note?.tags?.map(tag => tag.id) || [],
})

const availableTags = ref<Tag[]>([])
const showTagInput = ref(false)
const newTagName = ref('')
const newTagColor = ref('#3b82f6')

async function loadTags() {
  try {
    availableTags.value = await api.getTags()
  }
  catch (error) {
    console.error('Failed to load tags:', error)
  }
}

function toggleTag(tagId: number) {
  const index = form.selectedTagIds.indexOf(tagId)
  if (index > -1) {
    form.selectedTagIds.splice(index, 1)
  }
  else {
    form.selectedTagIds.push(tagId)
  }
}

async function createNewTag() {
  if (!newTagName.value.trim())
    return

  try {
    const newTag = await api.createTag({
      name: newTagName.value.trim(),
      color: newTagColor.value,
    })
    availableTags.value.push(newTag)
    form.selectedTagIds.push(newTag.id)
    newTagName.value = ''
    showTagInput.value = false
  }
  catch (error) {
    console.error('Failed to create tag:', error)
  }
}

function handleSave() {
  emit('save', {
    title: form.title,
    content: form.content,
    color: form.color,
    tags: form.selectedTagIds,
  })
}

onMounted(() => {
  loadTags()
})
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
          <icon-heroicons-x-mark class="h-5 w-5" />
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
              class="h-8 w-8 border-2 rounded-full transition-transform hover:scale-110"
              :class="{
                'border-gray-300': form.color !== color.value,
                'border-primary-500 ring-2 ring-primary-500': form.color === color.value,
                'bg-white shadow-inner': color.value === '#ffffff',
              }"
              :style="{ backgroundColor: color.value !== '#ffffff' ? color.value : undefined }"
              :title="color.name"
              @click="form.color = color.value"
            >
              <i
                v-if="form.color === color.value"
                class="i-heroicons-check h-4 w-4 text-gray-700"
              ></i>
            </button>
          </div>
        </div>

        <div>
          <label class="mb-2 block text-sm text-gray-700 font-medium">
            Tags
          </label>
          <div class="space-y-2">
            <!-- Selected tags -->
            <div v-if="form.selectedTagIds.length > 0" class="flex flex-wrap gap-1">
              <span
                v-for="tagId in form.selectedTagIds"
                :key="tagId"
                class="inline-flex items-center rounded-full px-2 py-1 text-xs font-medium"
                :style="{ backgroundColor: availableTags.find(t => t.id === tagId)?.color || '#e5e7eb', color: '#000' }"
              >
                {{ availableTags.find(t => t.id === tagId)?.name || 'Unknown' }}
                <button
                  class="ml-1 rounded-full p-0.5 hover:bg-black hover:bg-opacity-20"
                  @click="toggleTag(tagId)"
                >
                  <icon-heroicons-x-mark class="h-3 w-3" />
                </button>
              </span>
            </div>

            <!-- Available tags -->
            <div class="flex flex-wrap gap-1">
              <button
                v-for="tag in availableTags.filter(t => !form.selectedTagIds.includes(t.id))"
                :key="tag.id"
                type="button"
                class="border border-gray-300 rounded-full px-2 py-1 text-xs font-medium hover:border-gray-400"
                :style="{ borderColor: tag.color, color: tag.color }"
                @click="toggleTag(tag.id)"
              >
                + {{ tag.name }}
              </button>
            </div>

            <!-- Create new tag -->
            <div v-if="!showTagInput" class="pt-2">
              <button
                type="button"
                class="text-sm text-primary-600 hover:text-primary-700"
                @click="showTagInput = true"
              >
                + Create new tag
              </button>
            </div>

            <div v-else class="flex gap-2 pt-2">
              <input
                v-model="newTagName"
                type="text"
                placeholder="Tag name"
                class="flex-1 border border-gray-300 rounded px-2 py-1 text-sm focus:outline-none focus:ring-1 focus:ring-primary-500"
                @keyup.enter="createNewTag"
              >
              <input
                v-model="newTagColor"
                type="color"
                class="h-8 w-8 cursor-pointer border border-gray-300 rounded"
              >
              <button
                type="button"
                class="rounded bg-primary-600 px-2 py-1 text-sm text-white hover:bg-primary-700"
                @click="createNewTag"
              >
                Add
              </button>
              <button
                type="button"
                class="px-2 py-1 text-sm text-gray-600 hover:text-gray-800"
                @click="showTagInput = false"
              >
                Cancel
              </button>
            </div>
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
