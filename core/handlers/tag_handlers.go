package handlers

import (
	"dsn/core/services"
	"dsn/core/types"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetTagsHandler(tagService *services.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := tagService.GetAll()
		if err != nil {
			http.Error(w, "Failed to get tags", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tags)
	}
}

func CreateTagHandler(tagService *services.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			http.Error(w, "Tag name is required", http.StatusBadRequest)
			return
		}

		tag, err := tagService.Create(req)
		if err != nil {
			http.Error(w, "Failed to create tag", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tag)
	}
}

func UpdateTagHandler(tagService *services.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		var req types.UpdateTagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		tag, err := tagService.Update(tagID, req)
		if err != nil {
			http.Error(w, "Failed to update tag", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tag)
	}
}

func DeleteTagHandler(tagService *services.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		err = tagService.Delete(tagID)
		if err != nil {
			http.Error(w, "Failed to delete tag", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func AssignTagToNoteHandler(tagService *services.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteID, err := strconv.Atoi(r.PathValue("noteId"))
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		tagID, err := strconv.Atoi(r.PathValue("tagId"))
		if err != nil {
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		err = tagService.AssignToNote(noteID, tagID)
		if err != nil {
			http.Error(w, "Failed to assign tag to note", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func RemoveTagFromNoteHandler(tagService *services.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteID, err := strconv.Atoi(r.PathValue("noteId"))
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		tagID, err := strconv.Atoi(r.PathValue("tagId"))
		if err != nil {
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		err = tagService.RemoveFromNote(noteID, tagID)
		if err != nil {
			http.Error(w, "Failed to remove tag from note", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func SetNoteTagsHandler(tagService *services.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		var req types.AssignTagsToNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = tagService.SetNoteTags(noteID, req.TagIDs)
		if err != nil {
			http.Error(w, "Failed to set note tags", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
