package main

import (
	"log"
	"net/http"
	"path/filepath"

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
	tagService := services.NewTagService(db)

	mux := http.NewServeMux()

	// auth routes
	mux.HandleFunc("POST /api/register", handlers.RegisterHandler(userService, authService))
	mux.HandleFunc("POST /api/login", handlers.LoginHandler(userService, authService))
	mux.HandleFunc("POST /api/logout", handlers.LogoutHandler())
	mux.Handle("GET /api/auth/check", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.CheckAuthHandler(userService))))

	// api routes
	mux.Handle("GET /api/notes", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetNotesHandler(noteService))))
	mux.Handle("GET /api/notes/search", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.SearchNotesHandler(noteService))))
	mux.Handle("POST /api/notes", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.CreateNoteHandler(noteService))))
	mux.Handle("GET /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetNoteHandler(noteService))))
	mux.Handle("PUT /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.UpdateNoteHandler(noteService))))
	mux.Handle("PATCH /api/notes/{id}/pin", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.TogglePinHandler(noteService))))
	mux.Handle("PATCH /api/notes/{id}/archive", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.ToggleArchiveHandler(noteService))))
	mux.Handle("PUT /api/notes/order", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.UpdateNotesOrderHandler(noteService))))
	mux.Handle("DELETE /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.DeleteNoteHandler(noteService))))

	// tag routes
	mux.Handle("GET /api/tags", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetTagsHandler(tagService))))
	mux.Handle("POST /api/tags", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.CreateTagHandler(tagService))))
	mux.Handle("PUT /api/tags/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.UpdateTagHandler(tagService))))
	mux.Handle("DELETE /api/tags/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.DeleteTagHandler(tagService))))
	mux.Handle("POST /api/notes/{noteId}/tags/{tagId}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.AssignTagToNoteHandler(tagService))))
	mux.Handle("DELETE /api/notes/{noteId}/tags/{tagId}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.RemoveTagFromNoteHandler(tagService))))
	mux.Handle("PUT /api/notes/{id}/tags", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.SetNoteTagsHandler(tagService))))

	// admin routes
	mux.Handle("GET /api/users", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetUsersHandler(userService))))
	mux.Handle("DELETE /api/users/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.DeleteUserHandler(userService))))

	// Serve uploaded files
	uploadsDir := filepath.Join(config.DataDirectoryPath, "uploads")
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// Upload route
	mux.Handle("POST /api/upload/image", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.UploadImageHandler())))

	handler := handlers.CorsMiddleware(mux)

	log.Printf("Server starting on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}
