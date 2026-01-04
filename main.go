package main

import (
	"log"
	"net/http"

	"dsn/core/api"
	"dsn/core/config"
	"dsn/core/database"
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
	mux.HandleFunc("POST /api/register", api.RegisterHandler(userService, authService))
	mux.HandleFunc("POST /api/login", api.LoginHandler(userService, authService))
	mux.HandleFunc("POST /api/logout", api.LogoutHandler())

	// api routes
	mux.Handle("GET /api/notes", api.AuthMiddleware(authService)(http.HandlerFunc(api.GetNotesHandler(noteService))))
	mux.Handle("POST /api/notes", api.AuthMiddleware(authService)(http.HandlerFunc(api.CreateNoteHandler(noteService))))
	mux.Handle("GET /api/notes/{id}", api.AuthMiddleware(authService)(http.HandlerFunc(api.GetNoteHandler(noteService))))
	mux.Handle("PUT /api/notes/{id}", api.AuthMiddleware(authService)(http.HandlerFunc(api.UpdateNoteHandler(noteService))))
	mux.Handle("DELETE /api/notes/{id}", api.AuthMiddleware(authService)(http.HandlerFunc(api.DeleteNoteHandler(noteService))))

	// admin routes
	mux.Handle("GET /api/users", api.AuthMiddleware(authService)(http.HandlerFunc(api.GetUsersHandler(userService))))
	mux.Handle("DELETE /api/users/{id}", api.AuthMiddleware(authService)(http.HandlerFunc(api.DeleteUserHandler(userService))))

	handler := api.CorsMiddleware(mux)

	log.Printf("Server starting on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}
