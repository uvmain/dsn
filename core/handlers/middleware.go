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
			if !config.NoAuthForUserZero {
				claims, err := authService.GetUserFromRequest(r)
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				r.Header.Set("X-User-ID", strconv.Itoa(claims.UserID))
				r.Header.Set("X-Username", claims.Username)
				r.Header.Set("X-Is-Admin", strconv.FormatBool(claims.IsAdmin))
			} else {
				// set default headers for no-auth mode
				r.Header.Set("X-User-ID", "0")
				r.Header.Set("X-Username", "noAuthUser")
				r.Header.Set("X-Is-Admin", "true")
			}

			next.ServeHTTP(w, r)
		})
	}
}

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
