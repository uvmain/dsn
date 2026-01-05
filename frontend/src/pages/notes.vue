<script setup lang="ts">
import type { Note } from '~/types'
import TagsAside from '~/components/TagsAside.vue'
import { api } from '~/composables/useApi'

const notes = ref<Note[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const showModal = ref(false)
const selectedNote = ref<Note | null>(null)
const searchQuery = ref('')
const isSearching = ref(false)
const draggedNote = ref<Note | null>(null)
const selectedTagIds = ref<number[]>([])

// Computed properties to separate pinned and unpinned notes
const pinnedNotes = computed(() => notes.value.filter(note => note.pinned))
const otherNotes = computed(() => notes.value.filter(note => !note.pinned))

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

    // Apply tag filtering
    if (selectedTagIds.value.length > 0) {
      notes.value = notes.value.filter((note: Note) =>
        note.tags?.some(tag => selectedTagIds.value.includes(tag.id)),
      )
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
  selectedTagIds.value = []
  loadNotes()
}

function onFilterByTag(tagId: number | null) {
  if (tagId === null) {
    selectedTagIds.value = []
  }
  else {
    selectedTagIds.value = [tagId]
  }
  loadNotes()
}

function onDragStart(event: DragEvent, note: Note) {
  draggedNote.value = note
  event.dataTransfer!.effectAllowed = 'move'
  event.dataTransfer!.setData('text/html', (event.target as HTMLElement).outerHTML)
}

function onDragOver(event: DragEvent) {
  event.preventDefault()
  event.dataTransfer!.dropEffect = 'move'
  return false
}

function onDragEnter(event: DragEvent) {
  event.preventDefault()
  const target = event.target as HTMLElement
  const noteCard = target.closest('.note-card')
  if (noteCard) {
    noteCard.classList.add('drag-over')
  }
}

function onDragLeave(event: DragEvent) {
  event.preventDefault()
  const target = event.target as HTMLElement
  const noteCard = target.closest('.note-card')
  if (noteCard) {
    noteCard.classList.remove('drag-over')
  }
}

async function onDrop(event: DragEvent, targetNote: Note) {
  event.preventDefault()
  const target = event.target as HTMLElement
  const noteCard = target.closest('.note-card')
  if (noteCard) {
    noteCard.classList.remove('drag-over')
  }

  if (!draggedNote.value || draggedNote.value.id === targetNote.id) {
    return
  }

  // Determine which section the dragged note and target note are in
  const draggedIsPinned = draggedNote.value.pinned
  const targetIsPinned = targetNote.pinned

  // If they're in different sections, toggle the pin status
  if (draggedIsPinned !== targetIsPinned) {
    try {
      await api.togglePin(draggedNote.value.id, { pinned: targetIsPinned })
      // Reload notes to update sections
      await loadNotes()
    }
    catch (err) {
      console.error('Failed to toggle pin:', err)
      error.value = 'Failed to toggle pin status. Please try again.'
    }
    draggedNote.value = null
    return
  }

  // Same section reordering logic
  // Get the notes in the same section
  const sectionNotes = draggedIsPinned ? pinnedNotes.value : otherNotes.value
  const draggedIndex = sectionNotes.findIndex(n => n.id === draggedNote.value!.id)
  const targetIndex = sectionNotes.findIndex(n => n.id === targetNote.id)

  if (draggedIndex === -1 || targetIndex === -1) {
    return
  }

  // Reorder within the section
  const newSectionNotes = [...sectionNotes]
  const [draggedItem] = newSectionNotes.splice(draggedIndex, 1)
  newSectionNotes.splice(targetIndex, 0, draggedItem)

  // Update the global notes array with the new order within the section
  const newNotes = [...notes.value]

  // Replace the section in the global array
  if (draggedIsPinned) {
    newNotes.splice(0, pinnedNotes.value.length, ...newSectionNotes)
  }
  else {
    newNotes.splice(pinnedNotes.value.length, otherNotes.value.length, ...newSectionNotes)
  }

  // Update the order values
  const noteOrders: Record<number, number> = {}
  newNotes.forEach((note, index) => {
    noteOrders[note.id] = index
  })

  // Update local state immediately for responsive UI
  notes.value = newNotes.map((note, index) => ({ ...note, order: index }))

  // Save the new order to the backend
  try {
    await api.updateNotesOrder(noteOrders)
  }
  catch (err) {
    console.error('Failed to update note order:', err)
    // Revert the changes on error
    await loadNotes()
  }

  draggedNote.value = null
}

function onDragEnd() {
  draggedNote.value = null
  // Clean up any remaining drag-over classes
  document.querySelectorAll('.note-card').forEach((card) => {
    card.classList.remove('drag-over')
  })
  document.querySelectorAll('.notes-grid').forEach((grid) => {
    grid.classList.remove('drag-over-grid')
  })
}

function onDragOverGrid(event: DragEvent, targetIsPinned: boolean) {
  // Only allow dropping if we're dragging a note and it's not already in this section
  if (!draggedNote.value || draggedNote.value.pinned === targetIsPinned) {
    return
  }

  event.preventDefault()
  event.dataTransfer!.dropEffect = 'move'

  // Add visual feedback to the grid
  const grid = event.currentTarget as HTMLElement
  grid.classList.add('drag-over-grid')
}

function onDragLeaveGrid(event: DragEvent) {
  // Remove visual feedback when leaving the grid
  const grid = event.currentTarget as HTMLElement
  const relatedTarget = event.relatedTarget as HTMLElement

  // Only remove the class if we're actually leaving the grid (not just moving to a child element)
  if (!grid.contains(relatedTarget)) {
    grid.classList.remove('drag-over-grid')
  }
}

async function onDropOnGrid(event: DragEvent, targetIsPinned: boolean) {
  event.preventDefault()

  if (!draggedNote.value) {
    return
  }

  const draggedIsPinned = draggedNote.value.pinned

  // If dropping in a different section, toggle the pin status
  if (draggedIsPinned !== targetIsPinned) {
    try {
      await api.togglePin(draggedNote.value.id, { pinned: targetIsPinned })
      // Reload notes to update sections
      await loadNotes()
    }
    catch (err) {
      console.error('Failed to toggle pin:', err)
      error.value = 'Failed to toggle pin status. Please try again.'
    }
  }

  draggedNote.value = null
}

function createNote() {
  selectedNote.value = null
  showModal.value = true
}

function editNote(note: Note) {
  selectedNote.value = note
  showModal.value = true
}

async function deleteNote(note: Note) {
  try {
    await api.deleteNote(note.id)
    notes.value = notes.value.filter(n => n.id !== note.id)
  }
  catch (err) {
    console.error('Failed to delete note:', err)
    error.value = 'Failed to delete note. Please try again.'
  }
}

async function togglePin(note: Note) {
  try {
    await api.togglePin(note.id, { pinned: !note.pinned })
    // Reload notes to update sections
    await loadNotes()
  }
  catch (err) {
    console.error('Failed to toggle pin:', err)
    error.value = 'Failed to toggle pin status. Please try again.'
  }
}

async function toggleArchive(note: Note) {
  try {
    await api.toggleArchive(note.id, { archived: !note.archived })
    // Remove the note from the current view if it was archived
    if (!note.archived) {
      notes.value = notes.value.filter(n => n.id !== note.id)
    }
    else {
      // If unarchiving, update the local note
      const index = notes.value.findIndex(n => n.id === note.id)
      if (index !== -1) {
        notes.value[index] = { ...notes.value[index], archived: !note.archived }
      }
    }
  }
  catch (err) {
    console.error('Failed to toggle archive:', err)
    error.value = 'Failed to toggle archive status. Please try again.'
  }
}

async function saveNote(noteData: { title: string, content: string, color: string, tags?: number[], pinned: boolean }) {
  try {
    error.value = null

    if (selectedNote.value) {
      // Update existing note
      const updatedNote = await api.updateNote(selectedNote.value.id, {
        title: noteData.title,
        content: noteData.content,
        color: noteData.color,
      })
      if (updatedNote) {
        const index = notes.value.findIndex(n => n.id === selectedNote.value!.id)
        if (index !== -1) {
          notes.value[index] = updatedNote
        }
      }

      // Update pinned status if changed
      if (noteData.pinned !== selectedNote.value.pinned) {
        await api.togglePin(selectedNote.value.id, { pinned: noteData.pinned })
        // Reload notes to update sections
        await loadNotes()
        return
      }

      // Update tags if provided
      if (noteData.tags !== undefined) {
        await api.setNoteTags(selectedNote.value.id, { tag_ids: noteData.tags })
        // Reload notes to get updated tags
        await loadNotes()
      }
    }
    else {
      // Create new note
      const newNote = await api.createNote({
        title: noteData.title,
        content: noteData.content,
        color: noteData.color,
        pinned: noteData.pinned,
        archived: false,
        order: 0, // New notes get order 0, will be reordered later
      })

      // Set tags if provided
      if (noteData.tags && noteData.tags.length > 0) {
        await api.setNoteTags(newNote.id, { tag_ids: noteData.tags })
        // Reload notes to get tags
        await loadNotes()
      }
      else {
        notes.value.unshift(newNote)
      }
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
  <div class="flex">
    <!-- Tags Sidebar -->
    <TagsAside
      :selected-tag-ids="selectedTagIds"
      @filter-by-tag="onFilterByTag"
    />

    <!-- Main Content -->
    <div class="container mx-auto flex-1 px-8 py-8 lg:px-16">
      <div class="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
        <h1 class="text-2xl text-gray-800 font-bold sm:text-3xl">
          My Notes
        </h1>
        <button class="btn flex flex-row items-center" @click="createNote">
          <icon-heroicons-plus class="mr-2 h-4 w-4" />
          New Note
        </button>
      </div>

      <!-- Search Bar -->
      <div class="mb-6">
        <div class="relative max-w-md">
          <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <icon-heroicons-magnifying-glass class="h-5 w-5 text-gray-400" />
          </div>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search notes..."
            class="w-full border border-gray-300 rounded-lg py-2 pl-10 pr-10 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20"
            @keyup.enter="performSearch"
          >
          <button
            v-if="searchQuery"
            class="absolute inset-y-0 right-0 flex items-center pr-3"
            @click="clearSearch"
          >
            <icon-heroicons-x-mark class="h-5 w-5 text-gray-400 hover:text-gray-600" />
          </button>
        </div>
        <div v-if="isSearching" class="mt-2 text-sm text-gray-600">
          Showing search results for "{{ searchQuery }}"
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <div class="mx-auto mb-4 h-8 w-8 animate-spin border-4 border-primary-200 border-t-primary-600 rounded-full"></div>
          <div class="text-gray-500">
            Loading notes...
          </div>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="border border-red-200 rounded-lg bg-red-50 p-6 text-center">
        <div class="mb-4 text-red-600">
          <icon-heroicons-exclamation-triangle class="mx-auto mb-2 h-12 w-12" />
          <div class="text-lg font-medium">
            Error Loading Notes
          </div>
        </div>
        <div class="mb-4 text-red-700">
          {{ error }}
        </div>
        <button class="btn" @click="loadNotes">
          Try Again
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="notes && notes.length === 0" class="py-12 text-center">
        <div class="mb-6">
          <icon-heroicons-document-text class="mx-auto mb-4 h-16 w-16 text-gray-300" />
          <h3 class="mb-2 text-xl text-gray-900 font-medium">
            No notes yet
          </h3>
          <p class="mx-auto mb-6 max-w-sm text-gray-500">
            Create your first note to get started organizing your thoughts and ideas.
          </p>
          <button class="btn" @click="createNote">
            <icon-heroicons-plus class="mr-2 h-5 w-5" />
            Create Your First Note
          </button>
        </div>
      </div>

      <!-- Notes Sections -->
      <div v-else>
        <!-- Pinned Notes Section -->
        <div v-if="pinnedNotes.length > 0" class="mb-8">
          <h2 class="mb-4 text-lg text-gray-700 font-semibold">
            Pinned
          </h2>
          <div
            class="notes-grid"
            @dragover="onDragOverGrid($event, true)"
            @dragleave="onDragLeaveGrid($event)"
            @drop="onDropOnGrid($event, true)"
          >
            <NoteCard
              v-for="note in pinnedNotes"
              :key="note.id"
              :note="note"
              class="note-card"
              draggable="true"
              @edit="editNote"
              @delete="deleteNote"
              @toggle-pin="togglePin"
              @toggle-archive="toggleArchive"
              @dragstart="onDragStart($event, note)"
              @dragover="onDragOver($event)"
              @dragenter="onDragEnter($event)"
              @dragleave="onDragLeave($event)"
              @drop="onDrop($event, note)"
              @dragend="onDragEnd"
            />
          </div>
        </div>

        <!-- Other Notes Section -->
        <div v-if="otherNotes.length > 0">
          <h2 v-if="pinnedNotes.length > 0" class="mb-4 text-lg text-gray-700 font-semibold">
            Other notes
          </h2>
          <div
            class="notes-grid"
            @dragover="onDragOverGrid($event, false)"
            @dragleave="onDragLeaveGrid($event)"
            @drop="onDropOnGrid($event, false)"
          >
            <NoteCard
              v-for="note in otherNotes"
              :key="note.id"
              :note="note"
              class="note-card"
              draggable="true"
              @edit="editNote"
              @delete="deleteNote"
              @toggle-pin="togglePin"
              @toggle-archive="toggleArchive"
              @dragstart="onDragStart($event, note)"
              @dragover="onDragOver($event)"
              @dragenter="onDragEnter($event)"
              @dragleave="onDragLeave($event)"
              @drop="onDrop($event, note)"
              @dragend="onDragEnd"
            />
          </div>
        </div>
      </div>

      <!-- Note Modal -->
      <NoteModal
        v-if="showModal"
        :note="selectedNote"
        @save="saveNote"
        @close="closeModal"
      />
    </div>
  </div>
</template>

<style scoped>
.notes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
  padding: 0;
}

.note-card {
  transition: all 0.2s ease;
}

.note-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.note-card.drag-over {
  transform: scale(1.02);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.2);
  border: 2px dashed #3b82f6;
}

.note-card[draggable="true"] {
  user-select: none;
}

.note-card[draggable="true"]:active {
  transform: rotate(1deg);
}

.notes-grid {
  position: relative;
  min-height: 120px;
}

.notes-grid.drag-over-grid {
  background-color: rgba(59, 130, 246, 0.1);
  border: 2px dashed #3b82f6;
  border-radius: 8px;
  padding: 1rem;
  margin: -1rem;
}
</style>
