import type { App } from 'vue'
import type { Router } from 'vue-router'

export interface UserModule {
  (ctx: {
    app: App
    router: Router
    routes: any[]
    isClient: boolean
    initialState: Record<string, any>
  }): void
}

export interface Note {
  id: number
  title: string
  content: string
  color: string
  pinned: boolean
  archived: boolean
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  username: string
  email: string
  is_admin: boolean
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  username: string
  email: string
  password: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface CreateNoteRequest {
  title: string
  content: string
  color: string
  pinned: boolean
  archived: boolean
}

export interface UpdateNoteRequest {
  title?: string
  content?: string
  color?: string
  pinned?: boolean
  archived?: boolean
}
