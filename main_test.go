package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"dsn/core/api"
	"dsn/core/database"
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

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("POST /api/register", api.RegisterHandler(userService, authService))
	mux.HandleFunc("POST /api/login", api.LoginHandler(userService, authService))

	// Protected routes
	mux.Handle("POST /api/notes", api.AuthMiddleware(authService)(http.HandlerFunc(api.CreateNoteHandler(noteService))))
	mux.Handle("GET /api/notes", api.AuthMiddleware(authService)(http.HandlerFunc(api.GetNotesHandler(noteService))))

	server := httptest.NewServer(mux)

	cleanup := func() {
		server.Close()
		db.Close()
	}

	return server, cleanup
}

func TestUserRegistrationAndNoteCreation(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

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

	// Test note creation
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
	noteResp, err := client.Do(req)
	if err != nil {
		t.Fatal("Failed to create note:", err)
	}
	defer noteResp.Body.Close()

	if noteResp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", noteResp.StatusCode)
	}

	fmt.Println("✅ User registration and note creation test passed!")
}
