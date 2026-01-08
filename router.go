package main

import (
	"dsn/core/auth"
	"dsn/core/config"
	"dsn/core/handlers"
	"dsn/core/logic"
	"dsn/core/services"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-swiss/compress"
	"github.com/rs/cors"
)

//go:embed all:frontend/dist
var dist embed.FS
var distSubFS fs.FS

func init() {
	var err error
	distSubFS, err = fs.Sub(dist, "frontend/dist")
	if err != nil {
		panic("Failed to create sub filesystem: " + err.Error())
	}
}

func StartServer(userService *services.UserService, authService *services.AuthService, noteService *services.NoteService, tagService *services.TagService) *http.Server {
	mux := http.NewServeMux()

	// auth routes
	mux.HandleFunc("POST /api/register", handlers.RegisterHandler(userService, authService))
	mux.HandleFunc("POST /api/login", handlers.LoginHandler(userService, authService))
	mux.HandleFunc("POST /api/logout", handlers.LogoutHandler())
	mux.Handle("GET /api/auth/check", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.CheckAuthHandler(userService))))

	// api routes
	mux.Handle("GET /api/notes", auth.Middleware(authService)(http.HandlerFunc(handlers.GetNotesHandler(noteService))))
	mux.Handle("GET /api/notes/search", auth.Middleware(authService)(http.HandlerFunc(handlers.SearchNotesHandler(noteService))))
	mux.Handle("POST /api/notes", auth.Middleware(authService)(http.HandlerFunc(handlers.CreateNoteHandler(noteService))))
	mux.Handle("GET /api/notes/{id}", auth.Middleware(authService)(http.HandlerFunc(handlers.GetNoteHandler(noteService))))
	mux.Handle("PUT /api/notes/{id}", auth.Middleware(authService)(http.HandlerFunc(handlers.UpdateNoteHandler(noteService))))
	mux.Handle("PATCH /api/notes/{id}/pin", auth.Middleware(authService)(http.HandlerFunc(handlers.TogglePinHandler(noteService))))
	mux.Handle("PATCH /api/notes/{id}/archive", auth.Middleware(authService)(http.HandlerFunc(handlers.ToggleArchiveHandler(noteService))))
	mux.Handle("PUT /api/notes/order", auth.Middleware(authService)(http.HandlerFunc(handlers.UpdateNotesOrderHandler(noteService))))
	mux.Handle("DELETE /api/notes/{id}", auth.Middleware(authService)(http.HandlerFunc(handlers.DeleteNoteHandler(noteService))))

	// tag routes
	mux.Handle("GET /api/tags", auth.Middleware(authService)(http.HandlerFunc(handlers.GetTagsHandler(tagService))))
	mux.Handle("POST /api/tags", auth.Middleware(authService)(http.HandlerFunc(handlers.CreateTagHandler(tagService))))
	mux.Handle("PUT /api/tags/{id}", auth.Middleware(authService)(http.HandlerFunc(handlers.UpdateTagHandler(tagService))))
	mux.Handle("DELETE /api/tags/{id}", auth.Middleware(authService)(http.HandlerFunc(handlers.DeleteTagHandler(tagService))))
	mux.Handle("POST /api/notes/{noteId}/tags/{tagId}", auth.Middleware(authService)(http.HandlerFunc(handlers.AssignTagToNoteHandler(tagService))))
	mux.Handle("DELETE /api/notes/{noteId}/tags/{tagId}", auth.Middleware(authService)(http.HandlerFunc(handlers.RemoveTagFromNoteHandler(tagService))))
	mux.Handle("PUT /api/notes/{id}/tags", auth.Middleware(authService)(http.HandlerFunc(handlers.SetNoteTagsHandler(tagService))))

	// admin routes
	mux.Handle("GET /api/users", auth.Middleware(authService)(http.HandlerFunc(handlers.GetUsersHandler(userService))))
	mux.Handle("DELETE /api/users/{id}", auth.Middleware(authService)(http.HandlerFunc(handlers.DeleteUserHandler(userService))))

	// Serve uploaded files
	uploadsDir := filepath.Join(config.DataDirectoryPath, "uploads")
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// Upload route
	mux.Handle("POST /api/upload/image", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.UploadImageHandler())))

	// frontend routes
	mux.HandleFunc("/", handleFrontend)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.CorsBaseUrl},
		AllowCredentials: true,
		Debug:            false,
	})

	handler := c.Handler(
		compress.Middleware(mux),
	)

	serverAddress := fmt.Sprintf(":%d", config.Port)

	log.Printf("Application running at %s", serverAddress)

	server := &http.Server{
		Addr:    serverAddress,
		Handler: handler,
	}
	return server
}

func handleFrontend(w http.ResponseWriter, r *http.Request) {
	bootTime := logic.GetBootTime().Truncate(time.Second).UTC()

	cleanPath := path.Clean(r.URL.Path)
	if cleanPath == "/" {
		cleanPath = "/index.html"
	} else {
		cleanPath = strings.TrimPrefix(cleanPath, "/")
	}

	// static content
	file, err := distSubFS.Open(cleanPath)
	if err == nil {
		defer file.Close()

		http.ServeContent(w, r, cleanPath, bootTime, file.(io.ReadSeeker))
		return
	}

	// serve index.html for vue-router content
	indexFile, err := distSubFS.Open("index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusNotFound)
		return
	}
	defer indexFile.Close()

	http.ServeContent(w, r, "index.html", bootTime, indexFile.(io.ReadSeeker))
}
