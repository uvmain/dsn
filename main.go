package main

import (
	"log"

	"dsn/core/config"
	"dsn/core/database"
	"dsn/core/io"
	"dsn/core/services"
)

func main() {
	config.LoadConfig()
	io.CreateDirs()

	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	authService := services.NewAuthService(db)
	userService := services.NewUserService(db)
	noteService := services.NewNoteService(db)
	tagService := services.NewTagService(db)

	StartServer(userService, authService, noteService, tagService)
}
