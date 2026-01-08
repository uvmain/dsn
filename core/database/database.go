package database

import (
	"context"
	"database/sql"
	"dsn/core/config"
	"dsn/core/io"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var dbFile = "sqlite.db"
var DB *sql.DB
var err error

func Initialise(ctx context.Context) {
	openDatabase(ctx)
	createTables(ctx)
}

func openDatabase(ctx context.Context) {
	dbFilePath := filepath.Join(config.DatabaseDirectory, dbFile)

	if io.FileExists(dbFilePath) {
		log.Println("Database already exists")
	} else {
		log.Println("Creating new database file")
	}

	dataSource := fmt.Sprintf("file:%s?_journal_mode=WAL&_foreign_keys=on", dbFilePath)
	var err error
	DB, err = sql.Open("sqlite3", dataSource)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	DB.SetMaxIdleConns(5)

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
}

func CleanShutdown() {
	if DB != nil {
		log.Println("Closing database...")
		_, err = DB.Exec("PRAGMA wal_checkpoint(FULL);")
		if err != nil {
			log.Printf("Error committing WAL checkpoint on shutdown: %v", err)
		}
		DB.Close()

		// wait for data to be flushed to disk
		f, err := os.OpenFile(filepath.Join(config.DatabaseDirectory, dbFile), os.O_RDWR, 0660)
		if err == nil {
			f.Sync()
			f.Close()
		}
		log.Println("Database closed.")
	}
}
