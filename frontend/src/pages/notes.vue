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

const notes = ref<Note[]>([])
const loading = ref(true)
const showModal = ref(false)
const selectedNote = ref<Note | null>(null)
const showDeleteConfirm = ref(false)
const noteToDelete = ref<Note | null>(null)

async function loadNotes() {
  try {
    // TODO: Implement API call
    console.log('Loading notes...')

    // Simulate API call with dummy data
    await new Promise(resolve => setTimeout(resolve, 1000))

    notes.value = [
      {
        id: 1,
        title: 'Welcome to DSN',
        content: 'This is your first note! You can edit or delete it.',
        color: '#fef3c7',
        pinned: true,
        archived: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      },
    ]
  }
  catch (error) {
    console.error('Failed to load notes:', error)
  }
  finally {
    loading.value = false
  }
}

function createNote() {
  selectedNote.value = null
  showModal.value = true
}

function confirmDeleteNote() {
  if (noteToDelete.value) {
    // TODO: Implement API call
    notes.value = notes.value.filter(n => n.id !== noteToDelete.value!.id)
    noteToDelete.value = null
    showDeleteConfirm.value = false
  }
}

function cancelDeleteNote() {
  noteToDelete.value = null
  showDeleteConfirm.value = false
}

async function saveNote(noteData: any) {
  try {
    // TODO: Implement API call
    console.log('Saving note:', noteData)

    if (selectedNote.value) {
      // Update existing note
      const index = notes.value.findIndex(n => n.id === selectedNote.value!.id)
      if (index !== -1) {
        notes.value[index] = { ...notes.value[index], ...noteData }
      }
    }
    else {
      // Create new note
      const newNote: Note = {
        id: Date.now(),
        ...noteData,
        pinned: false,
        archived: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      notes.value.unshift(newNote)
    }

    closeModal()
  }
  catch (error) {
    console.error('Failed to save note:', error)
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

    <div v-if="loading" class="py-8 text-center">
      <div class="text-gray-500">
        Loading notes...
      </div>
    </div>

    <div v-else-if="notes.length === 0" class="py-12 text-center">
      <div class="mb-4 text-gray-500">
        No notes yet
      </div>
      <button class="btn" @click="createNote">
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
              Are you sure you want to delete this note?
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
      </button>
    </div>
  </div>
</template>
