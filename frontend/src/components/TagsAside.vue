<script setup lang="ts">
import type { Tag } from '~/types'
import { api } from '~/composables/useApi'

interface Props {
  selectedTagIds?: number[]
}

interface Emits {
  filterByTag: [tagId: number | null]
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const tags = ref<Tag[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const showCreateForm = ref(false)
const editingTag = ref<Tag | null>(null)
const newTagName = ref('')
const newTagColor = ref('#e0e0e0')
const editTagName = ref('')
const editTagColor = ref('#e0e0e0')

// Predefined color palette
const colorPalette = [
  '#e0e0e0',
  '#f87171',
  '#fb923c',
  '#fbbf24',
  '#a3e635',
  '#34d399',
  '#22d3ee',
  '#60a5fa',
  '#a78bfa',
  '#f472b6',
  '#64748b',
  '#374151',
  '#1f2937',
]

async function loadTags() {
  try {
    loading.value = true
    error.value = null
    tags.value = await api.getTags()
  }
  catch (err) {
    console.error('Failed to load tags:', err)
    error.value = 'Failed to load tags. Please try again.'
  }
  finally {
    loading.value = false
  }
}

function startCreateTag() {
  newTagName.value = ''
  newTagColor.value = '#e0e0e0'
  showCreateForm.value = true
  editingTag.value = null
}

function cancelCreateTag() {
  showCreateForm.value = false
  newTagName.value = ''
  newTagColor.value = '#e0e0e0'
}

async function createTag() {
  if (!newTagName.value.trim()) {
    return
  }

  try {
    const tag = await api.createTag({
      name: newTagName.value.trim(),
      color: newTagColor.value,
    })
    tags.value.push(tag)
    tags.value.sort((a, b) => a.name.localeCompare(b.name))
    cancelCreateTag()
  }
  catch (err) {
    console.error('Failed to create tag:', err)
    error.value = 'Failed to create tag. Please try again.'
  }
}

function startEditTag(tag: Tag) {
  editingTag.value = tag
  editTagName.value = tag.name
  editTagColor.value = tag.color
  showCreateForm.value = false
}

function cancelEditTag() {
  editingTag.value = null
  editTagName.value = ''
  editTagColor.value = '#e0e0e0'
}

async function updateTag() {
  if (!editingTag.value || !editTagName.value.trim()) {
    return
  }

  try {
    const updatedTag = await api.updateTag(editingTag.value.id, {
      name: editTagName.value.trim(),
      color: editTagColor.value,
    })

    const index = tags.value.findIndex(t => t.id === editingTag.value!.id)
    if (index !== -1) {
      tags.value[index] = updatedTag
      tags.value.sort((a, b) => a.name.localeCompare(b.name))
    }

    cancelEditTag()
  }
  catch (err) {
    console.error('Failed to update tag:', err)
    error.value = 'Failed to update tag. Please try again.'
  }
}

async function deleteTag(tag: Tag) {
  // eslint-disable-next-line no-alert
  if (!confirm(`Are you sure you want to delete the tag "${tag.name}"? This will remove it from all notes.`)) {
    return
  }

  try {
    await api.deleteTag(tag.id)
    tags.value = tags.value.filter(t => t.id !== tag.id)

    // If the deleted tag was selected, clear the filter
    if (props.selectedTagIds?.includes(tag.id)) {
      emit('filterByTag', null)
    }
  }
  catch (err) {
    console.error('Failed to delete tag:', err)
    error.value = 'Failed to delete tag. Please try again.'
  }
}

function toggleTagFilter(tagId: number) {
  const isSelected = props.selectedTagIds?.includes(tagId)
  emit('filterByTag', isSelected ? null : tagId)
}

function clearFilter() {
  emit('filterByTag', null)
}

// Load tags on component mount
onMounted(() => {
  loadTags()
})
</script>

<template>
  <aside class="tags-aside">
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-lg text-gray-800 font-semibold">
        Tags
      </h2>
      <button
        class="icon-btn text-gray-500 hover:text-gray-700"
        title="Create new tag"
        @click="startCreateTag"
      >
        <icon-heroicons-plus class="h-5 w-5" />
      </button>
    </div>

    <!-- Create Tag Form -->
    <div v-if="showCreateForm" class="mb-4 border border-gray-200 rounded-lg bg-gray-50 p-3">
      <div class="mb-3">
        <input
          v-model="newTagName"
          type="text"
          placeholder="Tag name"
          class="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20"
          @keyup.enter="createTag"
        >
      </div>

      <div class="mb-3">
        <div class="mb-2 text-sm text-gray-600">
          Color
        </div>
        <div class="flex flex-wrap gap-1">
          <button
            v-for="color in colorPalette"
            :key="color"
            class="h-6 w-6 border-2 border-gray-300 rounded"
            :class="{ 'border-gray-800': newTagColor === color }"
            :style="{ backgroundColor: color }"
            @click="newTagColor = color"
          />
        </div>
      </div>

      <div class="flex space-x-2">
        <button
          class="btn-sm"
          :disabled="!newTagName.trim()"
          @click="createTag"
        >
          Create
        </button>
        <button
          class="btn-sm-secondary"
          @click="cancelCreateTag"
        >
          Cancel
        </button>
      </div>
    </div>

    <!-- Edit Tag Form -->
    <div v-if="editingTag" class="mb-4 border border-gray-200 rounded-lg bg-gray-50 p-3">
      <div class="mb-3">
        <input
          v-model="editTagName"
          type="text"
          placeholder="Tag name"
          class="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20"
          @keyup.enter="updateTag"
        >
      </div>

      <div class="mb-3">
        <div class="mb-2 text-sm text-gray-600">
          Color
        </div>
        <div class="flex flex-wrap gap-1">
          <button
            v-for="color in colorPalette"
            :key="color"
            class="h-6 w-6 border-2 border-gray-300 rounded"
            :class="{ 'border-gray-800': editTagColor === color }"
            :style="{ backgroundColor: color }"
            @click="editTagColor = color"
          />
        </div>
      </div>

      <div class="flex space-x-2">
        <button
          class="btn-sm"
          :disabled="!editTagName.trim()"
          @click="updateTag"
        >
          Update
        </button>
        <button
          class="btn-sm-secondary"
          @click="cancelEditTag"
        >
          Cancel
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="text-center">
        <div class="mx-auto mb-2 h-6 w-6 animate-spin border-4 border-primary-200 border-t-primary-600 rounded-full"></div>
        <div class="text-sm text-gray-500">
          Loading tags...
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="border border-red-200 rounded-lg bg-red-50 p-4 text-center">
      <div class="mb-2 text-red-600">
        <icon-heroicons-exclamation-triangle class="mx-auto mb-1 h-5 w-5" />
        <div class="text-sm font-medium">
          Error Loading Tags
        </div>
      </div>
      <div class="mb-3 text-sm text-red-700">
        {{ error }}
      </div>
      <button class="btn-sm" @click="loadTags">
        Try Again
      </button>
    </div>

    <!-- Tags List -->
    <div v-else-if="tags.length > 0" class="space-y-1">
      <!-- Clear Filter Button -->
      <button
        v-if="selectedTagIds && selectedTagIds.length > 0"
        class="w-full rounded px-3 py-2 text-left text-sm text-gray-600 hover:bg-gray-100"
        @click="clearFilter"
      >
        <icon-heroicons-x-mark class="mr-2 inline h-4 w-4" />
        Clear filter
      </button>

      <!-- Tag Items -->
      <div
        v-for="tag in tags"
        :key="tag.id"
        class="group flex items-center justify-between rounded px-3 py-2 hover:bg-gray-100"
        :class="{ 'bg-blue-50': selectedTagIds?.includes(tag.id) }"
      >
        <button
          class="flex flex-1 items-center text-left"
          @click="toggleTagFilter(tag.id)"
        >
          <span
            class="mr-2 h-3 w-3 border border-gray-300 rounded-full"
            :style="{ backgroundColor: tag.color }"
          />
          <span class="text-sm text-gray-700">
            {{ tag.name }}
          </span>
        </button>

        <div class="ml-2 flex opacity-0 space-x-1 group-hover:opacity-100">
          <button
            class="icon-btn-sm text-gray-400 hover:text-gray-600"
            title="Edit tag"
            @click.stop="startEditTag(tag)"
          >
            <icon-heroicons-pencil class="h-3 w-3" />
          </button>
          <button
            class="icon-btn-sm text-gray-400 hover:text-red-500"
            title="Delete tag"
            @click.stop="deleteTag(tag)"
          >
            <icon-heroicons-trash class="h-3 w-3" />
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="py-8 text-center">
      <div class="mb-4">
        <icon-heroicons-tag class="mx-auto mb-2 h-8 w-8 text-gray-300" />
        <h3 class="mb-1 text-sm text-gray-900 font-medium">
          No tags yet
        </h3>
        <p class="text-xs text-gray-500">
          Create labels to organize your notes
        </p>
      </div>
      <button class="btn-sm" @click="startCreateTag">
        Create Label
      </button>
    </div>
  </aside>
</template>

<style scoped>
.tags-aside {
  width: 280px;
  min-height: 100vh;
  background-color: #f9fafb;
  border-right: 1px solid #e5e7eb;
  padding: 1.5rem;
}

.icon-btn {
  @apply flex items-center justify-center rounded p-1 transition-colors;
}

.icon-btn:hover {
  @apply bg-gray-200;
}

.icon-btn-sm {
  @apply flex items-center justify-center rounded p-0.5 transition-colors;
}

.btn {
  @apply rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-sm {
  @apply rounded bg-primary-600 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-sm-secondary {
  @apply rounded bg-gray-200 px-3 py-1.5 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-300;
}
</style>
