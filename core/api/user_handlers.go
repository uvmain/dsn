package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dsn/core/services"
)

// GetUsersHandler handles getting all users (admin only)
func GetUsersHandler(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin, err := getIsAdminFromRequest(r)
		if err != nil || !isAdmin {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		users, err := userService.GetAll()
		if err != nil {
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// DeleteUserHandler handles deleting a user (admin only)
func DeleteUserHandler(userService *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin, err := getIsAdminFromRequest(r)
		if err != nil || !isAdmin {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		userID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		err = userService.Delete(userID)
		if err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
