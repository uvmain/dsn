import type { AssignTagsToNoteRequest, CreateNoteRequest, CreateTagRequest, CreateUserRequest, LoginRequest, Note, Tag, ToggleArchiveRequest, TogglePinRequest, UpdateNoteRequest, UpdateTagRequest, User } from '~/types'

const BASE_URL = '/api'

class ApiClient {
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${BASE_URL}${endpoint}`
    const config: RequestInit = {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    const response = await fetch(url, config)

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`)
    }

    return response.json()
  }

  // Auth endpoints
  async register(data: CreateUserRequest): Promise<User> {
    return this.request<User>('/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async login(data: LoginRequest): Promise<User> {
    return this.request<User>('/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async logout(): Promise<void> {
    return this.request<void>('/logout', {
      method: 'POST',
    })
  }

  async checkAuth(): Promise<User> {
    return this.request<User>('/auth/check')
  }

  // Note endpoints
  async getNotes(): Promise<Note[]> {
    return this.request<Note[]>('/notes')
  }

  async searchNotes(query: string): Promise<Note[]> {
    return this.request<Note[]>(`/notes/search?q=${encodeURIComponent(query)}`)
  }

  async getNote(id: number): Promise<Note> {
    return this.request<Note>(`/notes/${id}`)
  }

  async createNote(data: CreateNoteRequest): Promise<Note> {
    return this.request<Note>('/notes', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateNote(id: number, data: UpdateNoteRequest): Promise<Note> {
    return this.request<Note>(`/notes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async togglePin(id: number, data: TogglePinRequest): Promise<Note> {
    return this.request<Note>(`/notes/${id}/pin`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  }

  async toggleArchive(id: number, data: ToggleArchiveRequest): Promise<Note> {
    return this.request<Note>(`/notes/${id}/archive`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  }

  async deleteNote(id: number): Promise<void> {
    return this.request<void>(`/notes/${id}`, {
      method: 'DELETE',
    })
  }

  async updateNotesOrder(noteOrders: Record<number, number>): Promise<void> {
    return this.request<void>('/notes/order', {
      method: 'PUT',
      body: JSON.stringify(noteOrders),
    })
  }

  // Tag endpoints
  async getTags(): Promise<Tag[]> {
    return this.request<Tag[]>('/tags')
  }

  async createTag(data: CreateTagRequest): Promise<Tag> {
    return this.request<Tag>('/tags', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateTag(id: number, data: UpdateTagRequest): Promise<Tag> {
    return this.request<Tag>(`/tags/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteTag(id: number): Promise<void> {
    return this.request<void>(`/tags/${id}`, {
      method: 'DELETE',
    })
  }

  async assignTagToNote(noteId: number, tagId: number): Promise<void> {
    return this.request<void>(`/notes/${noteId}/tags/${tagId}`, {
      method: 'POST',
    })
  }

  async removeTagFromNote(noteId: number, tagId: number): Promise<void> {
    return this.request<void>(`/notes/${noteId}/tags/${tagId}`, {
      method: 'DELETE',
    })
  }

  async setNoteTags(noteId: number, data: AssignTagsToNoteRequest): Promise<void> {
    return this.request<void>(`/notes/${noteId}/tags`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  // User endpoints (admin only)
  async getUsers(): Promise<User[]> {
    return this.request<User[]>('/users')
  }

  async deleteUser(id: number): Promise<void> {
    return this.request<void>(`/users/${id}`, {
      method: 'DELETE',
    })
  }
}

export const api = new ApiClient()

// Composable for easier usage
export function useApi() {
  return {
    api,
  }
}
