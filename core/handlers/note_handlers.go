package handlers

import (
	"dsn/core/services"
	"dsn/core/types"
	"encoding/json"
	"net/http"
	"strconv"
)

// GetNotesHandler handles getting all notes for a user
func GetNotesHandler(noteService *services.NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		includeArchived := r.URL.Query().Get("archived") == "true"
		notes, err := noteService.GetByUserID(userID, includeArchived)
		if err != nil {
			http.Error(w, "Failed to get notes", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	}
}

// CreateNoteHandler handles creating a new note
func CreateNoteHandler(noteService *services.NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req types.CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		note, err := noteService.Create(userID, req)
		if err != nil {
			http.Error(w, "Failed to create note", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

// GetNoteHandler handles getting a specific note
func GetNoteHandler(noteService *services.NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		noteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		note, err := noteService.GetByID(noteID, userID)
		if err != nil {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)
	}
}

// UpdateNoteHandler handles updating a note
func UpdateNoteHandler(noteService *services.NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		noteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		var req types.UpdateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		note, err := noteService.Update(noteID, userID, req)
		if err != nil {
			http.Error(w, "Failed to update note", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)
	}
}

// DeleteNoteHandler handles deleting a note
func DeleteNoteHandler(noteService *services.NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		noteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		err = noteService.Delete(noteID, userID)
		if err != nil {
			http.Error(w, "Failed to delete note", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
