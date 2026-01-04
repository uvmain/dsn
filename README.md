# DSN - Self hosted digital sticky notes

## Features

- Multi-user support with JWT authentication
- Cookie-based session management
- SQLite database for data persistence
- RESTful API built with stdlib net/http
- Admin privileges for first registered user
- Note management (create, read, update, delete)
- Note searching
- Note archiving and pinning
- Tag support (future enhancement)

## Prerequisites

- Go 1.25.5 or higher

## Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd dsn
```

2. Install dependencies:
```bash
task deps
```

3. Set environment variables (optional):
```bash
export JWT_SECRET="your-super-secret-jwt-key"
export PORT="8080"
export DB_PATH="./dsn.db"
```

4. Run the application:
```bash
task run
```

## Build Tasks

This project uses [Task](https://taskfile.dev/) as the build runner.

```bash
task --list              # Show all available tasks
task build               # Build the application
task run                 # Run the application
```

## API Endpoints

### Authentication
- `POST /api/register` - Register a new user
- `POST /api/login` - Login user
- `POST /api/logout` - Logout user

### Notes
- `GET /api/notes` - Get all notes for authenticated user
- `GET /api/notes?archived=true` - Get all notes including archived
- `POST /api/notes` - Create a new note
- `GET /api/notes/{id}` - Get specific note
- `PUT /api/notes/{id}` - Update note
- `DELETE /api/notes/{id}` - Delete note

### User Management (Admin only)
- `GET /api/users` - Get all users
- `DELETE /api/users/{id}` - Delete user

## Request/Response Examples

### Register User
```json
POST /api/register
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

### Login
```json
POST /api/login
{
  "username": "john_doe",
  "password": "securepassword"
}
```

### Create Note
```json
POST /api/notes
{
  "title": "My First Note",
  "content": "This is the content of my note",
  "color": "#ffeb3b",
  "pinned": false,
  "archived": false
}
```

### Update Note
```json
PUT /api/notes/1
{
  "title": "Updated Note Title",
  "content": "Updated content",
  "pinned": true
}
```