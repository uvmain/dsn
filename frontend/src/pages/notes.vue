<script setup lang="ts">
import type { Note } from '~/types'
import { api } from '~/composables/useApi'

const notes = ref<Note[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const showModal = ref(false)
const selectedNote = ref<Note | null>(null)
const showDeleteConfirm = ref(false)
const noteToDelete = ref<Note | null>(null)
const searchQuery = ref('')
const isSearching = ref(false)

async function loadNotes() {
  try {
    loading.value = true
    error.value = null

    if (searchQuery.value.trim()) {
      isSearching.value = true
      notes.value = await api.searchNotes(searchQuery.value.trim())
    }
    else {
      isSearching.value = false
      notes.value = await api.getNotes()
    }
  }
  catch (err) {
    console.error('Failed to load notes:', err)
    error.value = 'Failed to load notes. Please try again.'
  }
  finally {
    loading.value = false
  }
}

function performSearch() {
  loadNotes()
}

function clearSearch() {
  searchQuery.value = ''
  loadNotes()
}

function createNote() {
  selectedNote.value = null
  showModal.value = true
}

function editNote(note: Note) {
  selectedNote.value = note
  showModal.value = true
}

function deleteNote(note: Note) {
  noteToDelete.value = note
  showDeleteConfirm.value = true
}

async function confirmDeleteNote() {
  if (!noteToDelete.value)
    return

  try {
    await api.deleteNote(noteToDelete.value.id)
    notes.value = notes.value.filter(n => n.id !== noteToDelete.value!.id)
    noteToDelete.value = null
    showDeleteConfirm.value = false
  }
  catch (err) {
    console.error('Failed to delete note:', err)
    error.value = 'Failed to delete note. Please try again.'
  }
}

function cancelDeleteNote() {
  noteToDelete.value = null
  showDeleteConfirm.value = false
}

async function saveNote(noteData: { title: string, content: string, color: string }) {
  try {
    error.value = null

    if (selectedNote.value) {
      // Update existing note
      const updatedNote = await api.updateNote(selectedNote.value.id, noteData)
      const index = notes.value.findIndex(n => n.id === selectedNote.value!.id)
      if (index !== -1) {
        notes.value[index] = updatedNote
      }
    }
    else {
      // Create new note
      const newNote = await api.createNote({
        ...noteData,
        pinned: false,
        archived: false,
      })
      notes.value.unshift(newNote)
    }

    closeModal()
  }
  catch (err) {
    console.error('Failed to save note:', err)
    error.value = 'Failed to save note. Please try again.'
  }
}

function closeModal() {
  showModal.value = false
  selectedNote.value = null
}

// Load notes on component mount
onMounted(() => {
  loadNotes()
})

useHead({
  title: 'My Notes - DSN',
})
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl text-gray-800 font-bold">
        My Notes
      </h1>
      <button class="btn" @click="createNote">
        <i class="i-heroicons-plus mr-2 h-4 w-4"></i>
        New Note
      </button>
    </div>

    <!-- Search Bar -->
    <div class="mb-6">
      <div class="relative max-w-md">
        <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
          <i class="i-heroicons-magnifying-glass h-5 w-5 text-gray-400"></i>
        </div>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search notes..."
          class="w-full border border-gray-300 rounded-lg py-2 pl-10 pr-10 focus:border-transparent focus:outline-none focus:ring-2 focus:ring-primary-500"
          @keyup.enter="performSearch"
        >
        <button
          v-if="searchQuery"
          class="absolute inset-y-0 right-0 flex items-center pr-3"
          @click="clearSearch"
        >
          <i class="i-heroicons-x-mark h-5 w-5 text-gray-400 hover:text-gray-600"></i>
        </button>
      </div>
      <div v-if="isSearching" class="mt-2 text-sm text-gray-600">
        Showing search results for "{{ searchQuery }}"
      </div>
    </div>

    <div v-if="loading" class="py-8 text-center">
      <div class="text-gray-500">
        Loading notes...
      </div>
    </div>

    <div v-else-if="error" class="py-8 text-center">
      <div class="mb-4 text-red-500">
        {{ error }}
      </div>
      <button class="btn" @click="loadNotes">
        Try Again
      </button>
    </div>

    <div v-else-if="notes.length === 0" class="py-12 text-center">
      <div class="mb-4 text-gray-500">
        No notes yet
      </div>
      <button class="btn" @click="createNote">
        Create Your First Note
      </button>
    </div>

    <div v-else class="grid gap-4 lg:grid-cols-3 md:grid-cols-2 xl:grid-cols-4">
      <NoteCard
        v-for="note in notes"
        :key="note.id"
        :note="note"
        @edit="editNote"
        @delete="deleteNote"
      />
    </div>

    <!-- Note Modal -->
    <NoteModal
      v-if="showModal"
      :note="selectedNote"
      @save="saveNote"
      @close="closeModal"
    />

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteConfirm" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-30">
      <div class="max-w-sm w-full rounded bg-white p-6 shadow-lg">
        <div class="mb-4 text-lg font-semibold">
          Delete Note
        </div>
        <div class="mb-6 text-gray-700">
          Are you sure you want to delete this note? This action cannot be undone.
        </div>
        <div class="flex justify-end gap-2">
          <button class="btn-secondary btn" @click="cancelDeleteNote">
            Cancel
          </button>
          <button class="btn-danger btn" @click="confirmDeleteNote">
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
