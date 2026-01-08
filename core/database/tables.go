package database

import (
	"context"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func createTables(ctx context.Context) error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	notesTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL DEFAULT '',
		content TEXT NOT NULL DEFAULT '',
		color TEXT DEFAULT '#ffffff',
		pinned BOOLEAN DEFAULT FALSE,
		archived BOOLEAN DEFAULT FALSE,
		order_position INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);`

	tagsTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		color TEXT DEFAULT '#e0e0e0',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	noteTagsTable := `
	CREATE TABLE IF NOT EXISTS note_tags (
		note_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (note_id, tag_id),
		FOREIGN KEY (note_id) REFERENCES notes (id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
	);`

	tables := []string{usersTable, notesTable, tagsTable, noteTagsTable}
	for _, table := range tables {
		if _, err := DB.ExecContext(ctx, table); err != nil {
			return err
		}
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_notes_pinned ON notes(pinned);",
		"CREATE INDEX IF NOT EXISTS idx_notes_archived ON notes(archived);",
		"CREATE INDEX IF NOT EXISTS idx_notes_order_position ON notes(order_position);",
	}

	for _, index := range indexes {
		if _, err := DB.ExecContext(ctx, index); err != nil {
			return err
		}
	}

	return nil
}
