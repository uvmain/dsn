package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"dsn/core/database"
	"dsn/core/handlers"
	"dsn/core/services"
	"dsn/core/types"
)

func setupTestServer(t *testing.T) (*httptest.Server, func()) {
	// Create temporary database
	db, err := database.Initialize(":memory:")
	if err != nil {
		t.Fatal("Failed to initialize test database:", err)
	}

	// Initialize services
	authService := services.NewAuthService("test-secret", db)
	userService := services.NewUserService(db)
	noteService := services.NewNoteService(db)
	tagService := services.NewTagService(db)

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("POST /api/register", handlers.RegisterHandler(userService, authService))
	mux.HandleFunc("POST /api/login", handlers.LoginHandler(userService, authService))

	// Protected routes
	mux.Handle("GET /api/notes", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetNotesHandler(noteService))))
	mux.Handle("POST /api/notes", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.CreateNoteHandler(noteService))))
	mux.Handle("GET /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetNoteHandler(noteService))))
	mux.Handle("PUT /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.UpdateNoteHandler(noteService))))
	mux.Handle("PATCH /api/notes/{id}/pin", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.TogglePinHandler(noteService))))
	mux.Handle("PATCH /api/notes/{id}/archive", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.ToggleArchiveHandler(noteService))))
	mux.Handle("DELETE /api/notes/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.DeleteNoteHandler(noteService))))
	mux.Handle("GET /api/notes/search", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.SearchNotesHandler(noteService))))
	mux.Handle("GET /api/tags", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetTagsHandler(tagService))))
	mux.Handle("POST /api/tags", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.CreateTagHandler(tagService))))
	mux.Handle("PUT /api/tags/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.UpdateTagHandler(tagService))))
	mux.Handle("DELETE /api/tags/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.DeleteTagHandler(tagService))))
	mux.Handle("POST /api/notes/{noteId}/tags/{tagId}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.AssignTagToNoteHandler(tagService))))
	mux.Handle("DELETE /api/notes/{noteId}/tags/{tagId}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.RemoveTagFromNoteHandler(tagService))))
	mux.Handle("PUT /api/notes/{id}/tags", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.SetNoteTagsHandler(tagService))))

	// Admin routes
	mux.Handle("GET /api/users", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.GetUsersHandler(userService))))
	mux.Handle("DELETE /api/users/{id}", handlers.AuthMiddleware(authService)(http.HandlerFunc(handlers.DeleteUserHandler(userService))))

	server := httptest.NewServer(mux)

	cleanup := func() {
		server.Close()
		db.Close()
	}

	return server, cleanup
}

// Helper function to register a user and return auth cookie
func registerAndLogin(t *testing.T, server *httptest.Server) *http.Cookie {
	// Test user registration
	registerReq := types.CreateUserRequest{
		Username: "testUser",
		Email:    "test@example.com",
		Password: "testPass123",
	}
	reqBody, _ := json.Marshal(registerReq)

	resp, err := http.Post(server.URL+"/api/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal("Failed to register user:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Extract auth cookie
	var authCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth_token" {
			authCookie = cookie
			break
		}
	}

	if authCookie == nil {
		t.Fatal("No auth cookie received")
	}

	return authCookie
}

// Helper function to create a test note and return its ID
func createTestNote(t *testing.T, server *httptest.Server, authCookie *http.Cookie) int {
	noteReq := types.CreateNoteRequest{
		Title:   "Test Note",
		Content: "This is a test note",
		Color:   "#ffeb3b",
	}
	noteReqBody, _ := json.Marshal(noteReq)

	req, _ := http.NewRequest("POST", server.URL+"/api/notes", bytes.NewBuffer(noteReqBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(authCookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to create note:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdNote types.Note
	if err := json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
		t.Fatal("Failed to decode created note:", err)
	}

	return createdNote.ID
}

func TestUserRegistrationAndLogin(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	authCookie := registerAndLogin(t, server)

	if authCookie == nil {
		t.Fatal("Expected auth cookie to be returned")
	}

	fmt.Println("☑️ User registration and login test passed!")
}

func TestNoteCRUD(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	authCookie := registerAndLogin(t, server)
	client := &http.Client{}

	// Test note creation
	noteID := createTestNote(t, server, authCookie)

	// Test get all notes
	getReq, _ := http.NewRequest("GET", server.URL+"/api/notes", nil)
	getReq.AddCookie(authCookie)

	getResp, err := client.Do(getReq)
	if err != nil {
		t.Fatal("Failed to get notes:", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", getResp.StatusCode)
	}

	var notes []types.Note
	if err := json.NewDecoder(getResp.Body).Decode(&notes); err != nil {
		t.Fatal("Failed to decode notes:", err)
	}

	if len(notes) != 1 {
		t.Fatalf("Expected 1 note, got %d", len(notes))
	}

	// Test get single note
	getSingleReq, _ := http.NewRequest("GET", server.URL+"/api/notes/"+fmt.Sprintf("%d", noteID), nil)
	getSingleReq.AddCookie(authCookie)

	getSingleResp, err := client.Do(getSingleReq)
	if err != nil {
		t.Fatal("Failed to get single note:", err)
	}
	defer getSingleResp.Body.Close()

	if getSingleResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", getSingleResp.StatusCode)
	}

	// Test update note
	updateReq := types.UpdateNoteRequest{
		Title: stringPtr("Updated Title"),
	}
	updateReqBody, _ := json.Marshal(updateReq)

	updateHttpReq, _ := http.NewRequest("PUT", server.URL+"/api/notes/"+fmt.Sprintf("%d", noteID), bytes.NewBuffer(updateReqBody))
	updateHttpReq.Header.Set("Content-Type", "application/json")
	updateHttpReq.AddCookie(authCookie)

	updateResp, err := client.Do(updateHttpReq)
	if err != nil {
		t.Fatal("Failed to update note:", err)
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", updateResp.StatusCode)
	}

	var updatedNote types.Note
	if err := json.NewDecoder(updateResp.Body).Decode(&updatedNote); err != nil {
		t.Fatal("Failed to decode updated note:", err)
	}

	if updatedNote.Title != "Updated Title" {
		t.Fatalf("Expected title 'Updated Title', got '%s'", updatedNote.Title)
	}

	fmt.Println("☑️ Note CRUD operations test passed!")
}

func TestNoteSearch(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	authCookie := registerAndLogin(t, server)
	createTestNote(t, server, authCookie)

	client := &http.Client{}

	// Test note search
	searchReq, _ := http.NewRequest("GET", server.URL+"/api/notes/search?q=test", nil)
	searchReq.AddCookie(authCookie)

	searchResp, err := client.Do(searchReq)
	if err != nil {
		t.Fatal("Failed to search notes:", err)
	}
	defer searchResp.Body.Close()

	if searchResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for search, got %d", searchResp.StatusCode)
	}

	var searchResults []types.Note
	if err := json.NewDecoder(searchResp.Body).Decode(&searchResults); err != nil {
		t.Fatal("Failed to decode search results:", err)
	}

	if len(searchResults) != 1 {
		t.Fatalf("Expected 1 search result, got %d", len(searchResults))
	}

	if searchResults[0].Title != "Test Note" {
		t.Fatalf("Expected note title 'Test Note', got '%s'", searchResults[0].Title)
	}

	fmt.Println("☑️ Note search test passed!")
}

func TestTagManagement(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	authCookie := registerAndLogin(t, server)
	client := &http.Client{}

	// Test tag creation
	tagReq := types.CreateTagRequest{
		Name:  "Test Tag",
		Color: "#ff0000",
	}
	tagReqBody, _ := json.Marshal(tagReq)

	tagCreateReq, _ := http.NewRequest("POST", server.URL+"/api/tags", bytes.NewBuffer(tagReqBody))
	tagCreateReq.Header.Set("Content-Type", "application/json")
	tagCreateReq.AddCookie(authCookie)

	tagResp, err := client.Do(tagCreateReq)
	if err != nil {
		t.Fatal("Failed to create tag:", err)
	}
	defer tagResp.Body.Close()

	if tagResp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 for tag creation, got %d", tagResp.StatusCode)
	}

	// Test get tags
	getTagsReq, _ := http.NewRequest("GET", server.URL+"/api/tags", nil)
	getTagsReq.AddCookie(authCookie)

	tagsResp, err := client.Do(getTagsReq)
	if err != nil {
		t.Fatal("Failed to get tags:", err)
	}
	defer tagsResp.Body.Close()

	if tagsResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for get tags, got %d", tagsResp.StatusCode)
	}

	var tags []types.Tag
	if err := json.NewDecoder(tagsResp.Body).Decode(&tags); err != nil {
		t.Fatal("Failed to decode tags:", err)
	}

	if len(tags) != 1 {
		t.Fatalf("Expected 1 tag, got %d", len(tags))
	}

	if tags[0].Name != "Test Tag" {
		t.Fatalf("Expected tag name 'Test Tag', got '%s'", tags[0].Name)
	}

	fmt.Println("☑️ Tag management test passed!")
}

func TestNoteToggles(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	authCookie := registerAndLogin(t, server)
	noteID := createTestNote(t, server, authCookie)
	client := &http.Client{}

	// Test toggle pin
	pinReq := types.TogglePinRequest{Pinned: true}
	pinReqBody, _ := json.Marshal(pinReq)

	togglePinReq, _ := http.NewRequest("PATCH", server.URL+"/api/notes/"+fmt.Sprintf("%d", noteID)+"/pin", bytes.NewBuffer(pinReqBody))
	togglePinReq.Header.Set("Content-Type", "application/json")
	togglePinReq.AddCookie(authCookie)

	pinResp, err := client.Do(togglePinReq)
	if err != nil {
		t.Fatal("Failed to toggle pin:", err)
	}
	defer pinResp.Body.Close()

	if pinResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for toggle pin, got %d", pinResp.StatusCode)
	}

	var pinnedNote types.Note
	if err := json.NewDecoder(pinResp.Body).Decode(&pinnedNote); err != nil {
		t.Fatal("Failed to decode pinned note:", err)
	}

	if !pinnedNote.Pinned {
		t.Fatal("Expected note to be pinned")
	}

	// Test toggle archive
	archiveReq := types.ToggleArchiveRequest{Archived: true}
	archiveReqBody, _ := json.Marshal(archiveReq)

	toggleArchiveReq, _ := http.NewRequest("PATCH", server.URL+"/api/notes/"+fmt.Sprintf("%d", noteID)+"/archive", bytes.NewBuffer(archiveReqBody))
	toggleArchiveReq.Header.Set("Content-Type", "application/json")
	toggleArchiveReq.AddCookie(authCookie)

	archiveResp, err := client.Do(toggleArchiveReq)
	if err != nil {
		t.Fatal("Failed to toggle archive:", err)
	}
	defer archiveResp.Body.Close()

	if archiveResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 for toggle archive, got %d", archiveResp.StatusCode)
	}

	var archivedNote types.Note
	if err := json.NewDecoder(archiveResp.Body).Decode(&archivedNote); err != nil {
		t.Fatal("Failed to decode archived note:", err)
	}

	if !archivedNote.Archived {
		t.Fatal("Expected note to be archived")
	}

	fmt.Println("☑️ Note toggles test passed!")
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
