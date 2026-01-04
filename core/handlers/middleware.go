package handlers

import (
	"net/http"
	"strconv"

	"dsn/core/config"
	"dsn/core/services"
)

func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := authService.GetUserFromRequest(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add user info to request headers for easy access in handlers
			r.Header.Set("X-User-ID", strconv.Itoa(claims.UserID))
			r.Header.Set("X-Username", claims.Username)
			r.Header.Set("X-Is-Admin", strconv.FormatBool(claims.IsAdmin))

			next.ServeHTTP(w, r)
		})
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.CorsBaseUrl)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper functions to extract user information from request headers
func getUserIDFromRequest(r *http.Request) (int, error) {
	userIDStr := r.Header.Get("X-User-ID")
	return strconv.Atoi(userIDStr)
}

func getUsernameFromRequest(r *http.Request) string {
	return r.Header.Get("X-Username")
}

func getIsAdminFromRequest(r *http.Request) (bool, error) {
	isAdminStr := r.Header.Get("X-Is-Admin")
	return strconv.ParseBool(isAdminStr)
}
