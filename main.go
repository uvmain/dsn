package main

import (
	"log"
	"net/http"

	"dsn/core/config"
	"dsn/core/database"
	"dsn/core/handlers"
	"dsn/core/services"
)

func main() {
	config.LoadConfig()

	db, err := database.Initialize(config.DatabasePath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	authService := services.NewAuthService(config.JwtSecret, db)
	userService := services.NewUserService(db)
	noteService := services.NewNoteService(db)

	mux := http.NewServeMux()

	// auth routes
	mux.HandleFunc("POST /api/register", handlers.RegisterHandler(userService, authService))
	mux.HandleFunc("POST /api/login", handlers.LoginHandler(userService, authService))
	mux.HandleFunc("POST /api/logout", handlers.LogoutHandler())

	// api routes
	mux.Handle("GET /api/notes", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetNotesHandler(noteService))))
	mux.Handle("POST /api/notes", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.CreateNoteHandler(noteService))))
	mux.Handle("GET /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetNoteHandler(noteService))))
	mux.Handle("PUT /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.UpdateNoteHandler(noteService))))
	mux.Handle("DELETE /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.DeleteNoteHandler(noteService))))

	// admin routes
	mux.Handle("GET /api/users", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetUsersHandler(userService))))
	mux.Handle("DELETE /api/users/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.DeleteUserHandler(userService))))

	handler := handlers.CorsMiddleware(mux)

	log.Printf("Server starting on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}
