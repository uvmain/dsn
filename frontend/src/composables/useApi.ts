import type { CreateNoteRequest, CreateUserRequest, LoginRequest, Note, UpdateNoteRequest, User } from '~/types'

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

  // Note endpoints
  async getNotes(): Promise<Note[]> {
    return this.request<Note[]>('/notes')
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

  async deleteNote(id: number): Promise<void> {
    return this.request<void>(`/notes/${id}`, {
      method: 'DELETE',
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
